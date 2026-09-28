package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
	"gorm.io/gorm"
)

type roomUsecase struct {
	roomRepo repository.RoomRepository
	seatRepo repository.SeatRepository
	db       *gorm.DB
}

func NewRoomUsecase(roomRepo repository.RoomRepository, seatRepo repository.SeatRepository, db *gorm.DB) uc.RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo, seatRepo: seatRepo, db: db}
}
func (u roomUsecase) GetAll(ctx context.Context) ([]entity.Room, error) {
	return u.roomRepo.FindAll(ctx)
}
func (u roomUsecase) GetByID(ctx context.Context, id uint) (*entity.Room, error) {
	return u.roomRepo.FindByID(ctx, id)
}
func (u roomUsecase) GetSeats(ctx context.Context, id uint) (*entity.Room, error) {
	room, err := u.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	seats, err := u.seatRepo.FindByRoom(ctx, id)
	if err != nil {
		return nil, err
	}
	active := make([]entity.Seat, 0, len(seats))
	for _, seat := range seats {
		if seat.Status == "active" {
			active = append(active, seat)
		}
	}
	room.Seats = active
	return room, nil
}
func (u roomUsecase) Create(ctx context.Context, room *entity.Room, seatRows map[string]int) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roomRepo := u.roomRepo.WithTx(tx)
		seatRepo := u.seatRepo.WithTx(tx)
		// 1. Buat room
		if err := roomRepo.Create(ctx, room); err != nil {
			return err
		}
		// 2. Buat seat berdasarkan konfigurasi
		for row, total := range seatRows {
			for number := 1; number <= total; number++ {
				seat := &entity.Seat{
					RoomID: room.ID,
					Label:  row,
					Number: number,
				}
				if err := seatRepo.Create(ctx, seat); err != nil {
					return err
				}
			}
		}
		active, err := countActiveSeats(ctx, seatRepo, room.ID)
		if err != nil {
			return err
		}
		room.Capacity = active
		return roomRepo.Update(ctx, room)
	})
}
func (u roomUsecase) Update(ctx context.Context, room *entity.Room, seatRows map[string]int) error {
	if len(seatRows) > 0 {
		total := 0
		for _, n := range seatRows {
			total += n
		}
		room.Capacity = total
	}
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roomRepo := u.roomRepo.WithTx(tx)
		seatRepo := u.seatRepo.WithTx(tx)
		if err := roomRepo.Update(ctx, room); err != nil {
			return err
		}
		for row, total := range seatRows {
			if err := u.UpdateSeats(ctx, seatRepo, room.ID, row, total); err != nil {
				return err
			}
		}
		// nonaktifkan kursi dari baris yang tidak ada di request
		keep := make(map[string]struct{}, len(seatRows))
		for row := range seatRows {
			keep[row] = struct{}{}
		}
		allSeats, err := seatRepo.FindByRoom(ctx, room.ID)
		if err != nil {
			return err
		}
		for _, seat := range allSeats {
			if _, ok := keep[seat.Label]; !ok && seat.Status != "inactive" {
				seat.Status = "inactive"
				if err := seatRepo.Update(ctx, &seat); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
func (u roomUsecase) UpdateSeats(ctx context.Context, seatRepo repository.SeatRepository, roomID uint, row string, newTotal int) error {
	seats, err := seatRepo.FindByRoomAndRow(ctx, roomID, row)
	if err != nil {
		return err
	}
	// Aktifkan / nonaktifkan seat yang sudah ada
	for _, seat := range seats {
		var newStatus string
		if seat.Number <= newTotal {
			newStatus = "active"
		} else {
			newStatus = "inactive"
		}
		if seat.Status != newStatus {
			seat.Status = newStatus
			if err := seatRepo.Update(ctx, &seat); err != nil {
				return err
			}
		}
	}
	// Cari nomor seat terbesar yang sudah pernah dibuat
	maxNumber := 0
	for _, seat := range seats {
		if seat.Number > maxNumber {
			maxNumber = seat.Number
		}
	}
	// Buat seat baru jika jumlah baru melebihi seat yang pernah ada
	for number := maxNumber + 1; number <= newTotal; number++ {
		seat := &entity.Seat{
			RoomID: roomID,
			Label:  row,
			Number: number,
			Status: "active",
		}
		if err := seatRepo.Create(ctx, seat); err != nil {
			return err
		}
	}
	return nil
}

func (u roomUsecase) Delete(ctx context.Context, room *entity.Room) error {
	return u.roomRepo.Delete(ctx, room)
}
func countActiveSeats(ctx context.Context, seatRepo repository.SeatRepository, roomID uint) (int, error) {
	seats, err := seatRepo.FindByRoom(ctx, roomID)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, s := range seats {
		if s.Status == "active" {
			n++
		}
	}
	return n, nil
}

type seatUsecase struct {
	seatRepo repository.SeatRepository
}

func NewSeatUsecase(seatRepo repository.SeatRepository) uc.SeatUsecase {
	return &seatUsecase{seatRepo: seatRepo}
}

func (u seatUsecase) GetByRoom(ctx context.Context, roomID uint) ([]entity.Seat, error) {
	return u.seatRepo.FindByRoom(ctx, roomID)
}
func (u seatUsecase) Create(ctx context.Context, seat *entity.Seat) error {
	return u.seatRepo.Create(ctx, seat)
}
func (u seatUsecase) Update(ctx context.Context, seat *entity.Seat) error {
	return u.seatRepo.Update(ctx, seat)
}
func (u seatUsecase) GetByID(ctx context.Context, id uint) (*entity.Seat, error) {
	return u.seatRepo.FindByID(ctx, id)
}
func (u seatUsecase) Delete(ctx context.Context, seat *entity.Seat) error {
	return u.seatRepo.Delete(ctx, seat)
}
