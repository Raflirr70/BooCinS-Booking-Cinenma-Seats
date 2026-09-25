package repository

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, room *entity.Room) error
	FindAll(ctx context.Context) ([]entity.Room, error)
	FindByID(ctx context.Context, id uint) (*entity.Room, error)
	FindAllWithDetail(ctx context.Context) ([]entity.Room, error)
}

type SeatRepository interface {
	Create(ctx context.Context, seat *entity.Seat) error
	Update(ctx context.Context, seat *entity.Seat) error
	Delete(ctx context.Context, seat *entity.Seat) error
	FindByID(ctx context.Context, id uint) (*entity.Seat, error)
	FindByRoom(ctx context.Context, filmID uint) ([]entity.Seat, error)
	FindByRoomAndRow(ctx context.Context, roomID uint, row string) ([]entity.Seat, error)
}
