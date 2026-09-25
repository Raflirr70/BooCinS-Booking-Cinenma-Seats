package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
)

type roomUsecase struct {
	roomRepo repository.RoomRepository
	seatRepo repository.SeatRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) uc.RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo}
}

func (u roomUsecase) GetAll(ctx context.Context) ([]entity.Room, error) {
	return u.roomRepo.FindAll(ctx)
}
func (u roomUsecase) GetByID(ctx context.Context, id uint) (*entity.Room, error) {
	return u.roomRepo.FindByID(ctx, id)
}
func (u roomUsecase) Create(ctx context.Context, room *entity.Room, seatRows map[string]int) error {
	// 1. Buat room
	if err := u.roomRepo.Create(ctx, room); err != nil {
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
			if err := u.seatRepo.Create(ctx, seat); err != nil {
				return err
			}
		}
	}

	return nil
}
func (u roomUsecase) Update(ctx context.Context, room *entity.Room, seatRows map[string]int) error {
	if err := u.roomRepo.Update(ctx, room); err != nil {
		return err
	}
	for row, total := range seatRows {
		if err := u.UpdateSeats(ctx, room.ID, row, total); err != nil {
			return err
		}
	}
	return nil
}
func (u roomUsecase) UpdateSeats(ctx context.Context, roomID uint, row string, newTotal int) error {
	seats, err := u.seatRepo.FindByRoomAndRow(ctx, roomID, row)
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
			if err := u.seatRepo.Update(ctx, &seat); err != nil {
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
		if err := u.seatRepo.Create(ctx, seat); err != nil {
			return err
		}
	}
	return nil
}

func (u roomUsecase) Delete(ctx context.Context, room *entity.Room) error {
	return u.roomRepo.Delete(ctx, room)
}
