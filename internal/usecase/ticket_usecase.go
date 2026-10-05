package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rafli/boocins/config"
	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
	"gorm.io/gorm"
)

const maxSeatsPerBooking = 10

type bookingUsecase struct {
	db               *gorm.DB
	scheduleRepo     repository.ScheduleRepository
	seatRepo         repository.SeatRepository
	scheduleSeatRepo repository.ScheduleSeatRepository
	transactionRepo  repository.TransactionRepository
	ticketRepo       repository.TicketRepository
	guestOrderRepo   repository.GuestOrderRepository
	userRepo         repository.UserRepository
	snapClient       snap.Client
}

func NewBookingUsecase(
	db *gorm.DB,
	scheduleRepo repository.ScheduleRepository,
	seatRepo repository.SeatRepository,
	scheduleSeatRepo repository.ScheduleSeatRepository,
	transactionRepo repository.TransactionRepository,
	ticketRepo repository.TicketRepository,
	guestOrderRepo repository.GuestOrderRepository,
	userRepo repository.UserRepository,
	midtransCfg config.MidtransConfig,
) uc.BookingUsecase {
	env := midtrans.Sandbox
	if midtransCfg.IsProduction {
		env = midtrans.Production
	}
	snapClient := snap.Client{}
	snapClient.New(midtransCfg.ServerKey, env)
	return &bookingUsecase{
		db: db, scheduleRepo: scheduleRepo, seatRepo: seatRepo,
		scheduleSeatRepo: scheduleSeatRepo, transactionRepo: transactionRepo,
		ticketRepo: ticketRepo, guestOrderRepo: guestOrderRepo,
		userRepo: userRepo, snapClient: snapClient,
	}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, in uc.CreateBookingInput) (*uc.CreateBookingOutput, error) {
	seatIDs := dedupeIDs(in.SeatIDs)
	if in.ScheduleID == 0 || len(seatIDs) == 0 || len(seatIDs) > maxSeatsPerBooking {
		return nil, uc.ErrInvalidBookingInput
	}
	if len(seatIDs) != len(in.SeatIDs) {
		return nil, uc.ErrDuplicateSeats
	}
	if u.snapClient.ServerKey == "" {
		return nil, uc.ErrPaymentNotConfigured
	}

	// member → ambil nama dari users;
	// guest → wajib memasukan nama.
	var userID *uint
	var firstName, lastName, email string
	var discount float64
	if in.UserID != nil {
		user, err := u.userRepo.FindByID(ctx, *in.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, uc.ErrBuyerNotFound
			}
			return nil, err
		}
		userID = &user.ID
		firstName, lastName, email = user.FirstName, user.LastName, user.Email
		if user.Membership != nil {
			discount = user.Membership.Discount // FindByID sudah Preload Membership
		}
	} else {
		if strings.TrimSpace(in.FirstName) == "" || strings.TrimSpace(in.LastName) == "" {
			return nil, uc.ErrGuestNameRequired
		}
		firstName, lastName, email = in.FirstName, in.LastName, in.Email
	}

	// order_id unik, 40 char (batas Midtrans 50).
	rawUUID, err := newUUID()
	if err != nil {
		return nil, err
	}
	orderID := "BOOCINS-" + strings.ReplaceAll(rawUUID, "-", "")

	// 2-5. Satu DB transaction: lock → tulis. Tanpa HTTP call di dalamnya.
	var out *uc.CreateBookingOutput
	var lockedSSIDs []uint
	err = u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SET LOCAL lock_timeout = '3s'").Error; err != nil {
			return err
		}

		schedule, err := u.scheduleRepo.WithTx(tx).FindByID(ctx, in.ScheduleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return uc.ErrScheduleNotFound
			}
			return err
		}
		if !schedule.Status {
			return uc.ErrScheduleInactive
		}

		roomSeats, err := u.seatRepo.WithTx(tx).FindByRoom(ctx, schedule.RoomID)
		if err != nil {
			return err
		}
		active := make(map[uint]entity.Seat, len(roomSeats))
		for _, s := range roomSeats {
			if s.Status == "active" {
				active[s.ID] = s
			}
		}
		for _, id := range seatIDs {
			if _, ok := active[id]; !ok {
				return uc.ErrSeatNotInRoom
			}
		}

		txSS := u.scheduleSeatRepo.WithTx(tx)
		if err := txSS.EnsureAvailableSeats(ctx, schedule.ID, seatIDs); err != nil {
			return err
		}
		locked, err := txSS.LockAvailableSeats(ctx, schedule.ID, seatIDs)
		if err != nil {
			return err
		}
		if len(locked) != len(seatIDs) {
			return uc.ErrSeatsUnavailable // seat sudah diambil orang lain
		}

		lockedSSIDs = make([]uint, 0, len(locked))
		for _, ss := range locked {
			lockedSSIDs = append(lockedSSIDs, ss.ID)
		}
		now := time.Now()
		if err := txSS.SetStatus(ctx, lockedSSIDs, "locked", &now); err != nil {
			return err
		}

		//total = jumlah seat × harga film − diskon membership.
		total := math.Round(schedule.Film.Price * float64(len(seatIDs)) * (1 - discount/100))
		if total <= 0 {
			return uc.ErrInvalidTotal
		}

		var guestOrderID *uint
		if userID == nil {
			guest := &entity.GuestOrder{Email: email}
			if err := u.guestOrderRepo.WithTx(tx).Create(ctx, guest); err != nil {
				return err
			}
			guestOrderID = &guest.ID
		}

		pm := "midtrans"
		trx := &entity.Transaction{
			UserID:          userID,
			GuestOrderID:    guestOrderID,
			Status:          "pending",
			TotalPrice:      total,
			PaymentMethod:   &pm,
			Source:          "online",
			MidtransOrderID: &orderID,
		}
		if err := u.transactionRepo.WithTx(tx).Create(ctx, trx); err != nil {
			return err
		}
		if guestOrderID != nil {
			if err := u.guestOrderRepo.WithTx(tx).SetTransaction(ctx, *guestOrderID, trx.ID); err != nil {
				return err
			}
		}

		tickets := make([]entity.Ticket, 0, len(locked))
		seats := make([]uc.BookedSeat, 0, len(locked))
		for _, ss := range locked {
			qr, err := newUUID()
			if err != nil {
				return err
			}
			tickets = append(tickets, entity.Ticket{
				UserID: userID, ScheduleSeatID: ss.ID,
				TransactionID: trx.ID, QRToken: qr, Status: "active",
			})
			s := active[ss.SeatID]
			seats = append(seats, uc.BookedSeat{SeatID: ss.SeatID, Label: s.Label, Number: s.Number})
		}
		if err := u.ticketRepo.WithTx(tx).CreateMany(ctx, tickets); err != nil {
			return err
		}

		out = &uc.CreateBookingOutput{
			TransactionID: trx.ID, OrderID: orderID,
			TotalPrice: total, Seats: seats,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	//Snap SETELAH commit — tidak ada lock yang ditahan selama HTTP call.
	snapResp, mterr := u.snapClient.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: orderID, GrossAmt: int64(out.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: firstName, LName: lastName, Email: email,
		},
	})
	if mterr != nil || snapResp == nil || snapResp.Token == "" {
		if cerr := u.compensate(ctx, out.TransactionID, lockedSSIDs); cerr != nil {
			return nil, fmt.Errorf("snap failed (%v), compensation failed: %w", mterr, cerr)
		}
		return nil, fmt.Errorf("%w: %v", uc.ErrPaymentFailed, mterr)
	}

	if err := u.transactionRepo.UpdateSnapToken(ctx, out.TransactionID, snapResp.Token); err != nil {
		if cerr := u.compensate(ctx, out.TransactionID, lockedSSIDs); cerr != nil {
			return nil, fmt.Errorf("save snap token failed (%v), compensation failed: %w", err, cerr)
		}
		return nil, err
	}

	// Return ke frontend untuk popup Midtrans.
	out.SnapToken = snapResp.Token
	out.RedirectURL = snapResp.RedirectURL
	return out, nil
}

// compensate: rollback manual untuk kegagalan SETELAH commit.
// Snap token: tidak ada uang berpindah, jadi tidak perlu melakukan refund.
func (u *bookingUsecase) compensate(ctx context.Context, transactionID uint, scheduleSeatIDs []uint) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.transactionRepo.WithTx(tx).UpdateStatus(ctx, transactionID, "cancelled"); err != nil {
			return err
		}
		if err := u.ticketRepo.WithTx(tx).CancelByTransaction(ctx, transactionID); err != nil {
			return err
		}
		return u.scheduleSeatRepo.WithTx(tx).SetStatus(ctx, scheduleSeatIDs, "available", nil)
	})
}

func dedupeIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
