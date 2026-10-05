package usecase

import (
	"context"
	"errors"
)

var (
	ErrInvalidBookingInput  = errors.New("schedule_id dan seat_ids wajib diisi (maks 10 seat)")
	ErrDuplicateSeats       = errors.New("seat_ids tidak boleh duplikat")
	ErrBuyerNotFound        = errors.New("buyer not found")
	ErrGuestNameRequired    = errors.New("first_name dan last_name wajib diisi untuk guest")
	ErrScheduleInactive     = errors.New("schedule tidak aktif")
	ErrSeatNotInRoom        = errors.New("seat tidak termasuk room schedule ini / tidak aktif")
	ErrSeatsUnavailable     = errors.New("satu atau lebih seat sudah tidak tersedia")
	ErrInvalidTotal         = errors.New("total harga tidak valid")
	ErrPaymentNotConfigured = errors.New("midtrans server key belum dikonfigurasi")
	ErrPaymentFailed        = errors.New("gagal membuat transaksi midtrans")
)

type CreateBookingInput struct {
	ScheduleID uint
	SeatIDs    []uint
	UserID     *uint // nil = guest
	FirstName  string
	LastName   string
	Email      string
}

type BookedSeat struct {
	SeatID uint   `json:"seat_id"`
	Label  string `json:"label"`
	Number int    `json:"number"`
}

type CreateBookingOutput struct {
	TransactionID uint         `json:"transaction_id"`
	OrderID       string       `json:"order_id"`
	SnapToken     string       `json:"snap_token"`
	RedirectURL   string       `json:"redirect_url"`
	TotalPrice    float64      `json:"total_price"`
	Seats         []BookedSeat `json:"seats"`
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, in CreateBookingInput) (*CreateBookingOutput, error)
}
