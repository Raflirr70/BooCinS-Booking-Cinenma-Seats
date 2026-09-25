package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type RoomUsecase interface {
	GetAll(ctx context.Context) ([]entity.Room, error)
	GetByID(ctx context.Context, id uint) (*entity.Room, error)
	Create(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	Update(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	UpdateSeats(ctx context.Context, roomID uint, row string, newTotal int) error
	Delete(ctx context.Context, room *entity.Room) error
}
type SeatUsecase interface {
}
