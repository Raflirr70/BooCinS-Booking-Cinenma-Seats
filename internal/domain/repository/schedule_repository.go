package repository

import (
	"context"
	"time"

	"github.com/rafli/boocins/internal/domain/entity"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, room *entity.Room) error
	FindAll(ctx context.Context) ([]entity.Room, error)
	FindByID(ctx context.Context, id uint) (*entity.Room, error)
	FindAllWithDetail(ctx context.Context) ([]entity.Room, error)
	WithTx(tx *gorm.DB) RoomRepository
}

type SeatRepository interface {
	Create(ctx context.Context, seat *entity.Seat) error
	Update(ctx context.Context, seat *entity.Seat) error
	Delete(ctx context.Context, seat *entity.Seat) error
	FindByID(ctx context.Context, id uint) (*entity.Seat, error)
	FindByRoom(ctx context.Context, filmID uint) ([]entity.Seat, error)
	FindByRoomAndRow(ctx context.Context, roomID uint, row string) ([]entity.Seat, error)
	WithTx(tx *gorm.DB) SeatRepository
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, schedule *entity.Schedule) error
	FindAll(ctx context.Context) ([]entity.Schedule, error)
	FindByID(ctx context.Context, id uint) (*entity.Schedule, error)
	FindByRoom(ctx context.Context, roomID uint) ([]entity.Schedule, error)
	FindByFilm(ctx context.Context, filmID uint) ([]entity.Schedule, error)
	WithTx(tx *gorm.DB) ScheduleRepository
}

type ScheduleSeatRepository interface {
	Create(ctx context.Context, ss *entity.ScheduleSeat) error
	Update(ctx context.Context, ss *entity.ScheduleSeat) error
	Delete(ctx context.Context, ss *entity.ScheduleSeat) error
	FindByID(ctx context.Context, id uint) (*entity.ScheduleSeat, error)
	FindBySchedule(ctx context.Context, scheduleID uint) ([]entity.ScheduleSeat, error)
	WithTx(tx *gorm.DB) ScheduleSeatRepository
	EnsureAvailableSeats(ctx context.Context, scheduleID uint, seatIDs []uint) error
	LockAvailableSeats(ctx context.Context, scheduleID uint, seatIDs []uint) ([]entity.ScheduleSeat, error)
	SetStatus(ctx context.Context, ids []uint, status string, lockedAt *time.Time) error
}
