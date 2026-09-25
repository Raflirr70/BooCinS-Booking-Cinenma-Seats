package postgres

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *entity.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}
func (r *roomRepository) Update(ctx context.Context, room *entity.Room) error {
	return r.db.WithContext(ctx).Where("id = ?", room.ID).Updates(room).Error
}
func (r *roomRepository) Delete(ctx context.Context, room *entity.Room) error {
	return r.db.WithContext(ctx).Delete(room).Error
}
func (r *roomRepository) FindAll(ctx context.Context) ([]entity.Room, error) {
	var rooms []entity.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	return rooms, err
}
func (r *roomRepository) FindByID(ctx context.Context, id uint) (*entity.Room, error) {
	var room entity.Room
	err := r.db.WithContext(ctx).Preload("Seats").First(&room, id).Error
	return &room, err
}
func (r *roomRepository) FindAllWithDetail(ctx context.Context) ([]entity.Room, error) {
	var rooms []entity.Room
	err := r.db.WithContext(ctx).Preload("Seats").Find(&rooms).Error
	return rooms, err
}

// ============================ Seats ============================

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) repository.SeatRepository {
	return &seatRepository{db: db}
}
func (r *seatRepository) Create(ctx context.Context, seat *entity.Seat) error {
	return r.db.WithContext(ctx).Create(seat).Error
}
func (r *seatRepository) Update(ctx context.Context, seat *entity.Seat) error {
	return r.db.WithContext(ctx).Where("id = ?", seat.ID).Updates(seat).Error
}
func (r *seatRepository) Delete(ctx context.Context, seat *entity.Seat) error {
	return r.db.WithContext(ctx).Delete(seat).Error
}
func (r *seatRepository) FindByID(ctx context.Context, id uint) (*entity.Seat, error) {
	var seat entity.Seat
	err := r.db.WithContext(ctx).First(&seat, id).Error
	return &seat, err
}
func (r *seatRepository) FindByRoom(ctx context.Context, roomID uint) ([]entity.Seat, error) {
	var seats []entity.Seat
	err := r.db.WithContext(ctx).Where("room_id = ?", roomID).Find(&seats).Error
	return seats, err
}
func (r *seatRepository) FindByRoomAndRow(ctx context.Context, roomID uint, row string) ([]entity.Seat, error) {
	var seats []entity.Seat
	err := r.db.WithContext(ctx).Where("room_id = ? AND row = ?", roomID, row).Order("number ASC").Find(&seats).Error
	return seats, err
}
