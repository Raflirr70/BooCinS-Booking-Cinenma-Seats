package usecase

import (
	"context"
	"errors"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
)

var ErrFilmNotFound = errors.New("film not found")
var ErrRoomNotFound = errors.New("room not found")

type RoomUsecase interface {
	GetAll(ctx context.Context) ([]entity.Room, error)
	GetByID(ctx context.Context, id uint) (*entity.Room, error)
	Create(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	Update(ctx context.Context, room *entity.Room, seatRows map[string]int) error
	UpdateSeats(ctx context.Context, seatRepo repository.SeatRepository, roomID uint, row string, newTotal int) error
	Delete(ctx context.Context, room *entity.Room) error

	GetSeats(ctx context.Context, id uint) (*entity.Room, error)
}
type SeatUsecase interface {
	GetByRoom(ctx context.Context, roomID uint) ([]entity.Seat, error)
	Create(ctx context.Context, seat *entity.Seat) error
	Update(ctx context.Context, seat *entity.Seat) error
	GetByID(ctx context.Context, id uint) (*entity.Seat, error)
	Delete(ctx context.Context, seat *entity.Seat) error
}
type ScheduleUsecase interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, schedule *entity.Schedule) error
	GetByID(ctx context.Context, id uint) (*entity.Schedule, error)
	GetByFilm(ctx context.Context, filmID uint) ([]entity.Schedule, error)
	GetByRoom(ctx context.Context, roomID uint) ([]entity.Schedule, error)
}

type ScheduleSeatUsecase interface {
	Create(ctx context.Context, ss *entity.ScheduleSeat) error
	Update(ctx context.Context, ss *entity.ScheduleSeat) error
	Delete(ctx context.Context, ss *entity.ScheduleSeat) error
	GetBySchedule(ctx context.Context, scheduleID uint) ([]entity.ScheduleSeat, error)
	GetByID(ctx context.Context, id uint) (*entity.ScheduleSeat, error)
}
