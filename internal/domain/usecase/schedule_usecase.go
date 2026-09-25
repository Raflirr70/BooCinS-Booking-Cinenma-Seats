package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
)

type RoomUsecase interface {
	GetAll(ctx context.Context) ([]entity.Room, error)
	GetByID(ctx context.Context, id uint) (*entity.Room, error)
	Create(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	Update(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	UpdateSeats(ctx context.Context, seatRepo repository.SeatRepository, roomID uint, row string, newTotal int) error
	Delete(ctx context.Context, room *entity.Room) error
}
type SeatUsecase interface {
	GetByRoom(ctx context.Context, roomID uint) ([]entity.Seat, error)
	Create(ctx context.Context, seat *entity.Seat) error
	Update(ctx context.Context, seat *entity.Seat) error
	GetByID(ctx context.Context, id uint) (*entity.Seat, error)
	Delete(ctx context.Context, seat *entity.Seat) error
}
