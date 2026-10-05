package postgres

import (
	"context"
	"time"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (r *roomRepository) WithTx(tx *gorm.DB) repository.RoomRepository {
	return &roomRepository{db: tx}
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
	err := r.db.WithContext(ctx).Where("room_id = ? AND label = ?", roomID, row).Order("number ASC").Find(&seats).Error
	return seats, err
}
func (r *seatRepository) WithTx(tx *gorm.DB) repository.SeatRepository {
	return &seatRepository{db: tx}
}

// ============================ Schedule ============================

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) repository.ScheduleRepository {
	return &scheduleRepository{db: db}
}
func (r *scheduleRepository) Create(ctx context.Context, schedule *entity.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}
func (r *scheduleRepository) Update(ctx context.Context, schedule *entity.Schedule) error {
	return r.db.WithContext(ctx).Where("id = ?", schedule.ID).Updates(schedule).Error
}
func (r *scheduleRepository) Delete(ctx context.Context, schedule *entity.Schedule) error {
	return r.db.WithContext(ctx).Delete(schedule).Error
}
func (r *scheduleRepository) FindAll(ctx context.Context) ([]entity.Schedule, error) {
	var schedules []entity.Schedule
	err := r.db.WithContext(ctx).Preload("Room").Preload("Film").Find(&schedules).Error
	return schedules, err
}
func (r *scheduleRepository) FindByID(ctx context.Context, id uint) (*entity.Schedule, error) {
	var schedule entity.Schedule
	err := r.db.WithContext(ctx).Preload("Film").Preload("Room.Seats").Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
func (r *scheduleRepository) FindByRoom(ctx context.Context, roomID uint) ([]entity.Schedule, error) {
	var schedules []entity.Schedule
	err := r.db.WithContext(ctx).Preload("Room").Preload("Film").Where("room_id = ?", roomID).Find(&schedules).Error
	return schedules, err
}
func (r *scheduleRepository) FindByFilm(ctx context.Context, filmID uint) ([]entity.Schedule, error) {
	var schedules []entity.Schedule
	err := r.db.WithContext(ctx).Preload("Room").Preload("Film").Where("film_id = ?", filmID).Find(&schedules).Error
	return schedules, err
}
func (r *scheduleRepository) WithTx(tx *gorm.DB) repository.ScheduleRepository {
	return &scheduleRepository{db: tx}
}

// ============================ ScheduleSeats ============================

type scheduleSeatRepository struct {
	db *gorm.DB
}

func NewScheduleSeatRepository(db *gorm.DB) repository.ScheduleSeatRepository {
	return &scheduleSeatRepository{db: db}
}
func (r *scheduleSeatRepository) Create(ctx context.Context, ss *entity.ScheduleSeat) error {
	return r.db.WithContext(ctx).Create(ss).Error
}
func (r *scheduleSeatRepository) Update(ctx context.Context, ss *entity.ScheduleSeat) error {
	return r.db.WithContext(ctx).Where("id = ?", ss.ID).Updates(ss).Error
}
func (r *scheduleSeatRepository) Delete(ctx context.Context, ss *entity.ScheduleSeat) error {
	return r.db.WithContext(ctx).Delete(ss).Error
}
func (r *scheduleSeatRepository) FindByID(ctx context.Context, id uint) (*entity.ScheduleSeat, error) {
	var scheduleSeat *entity.ScheduleSeat
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&scheduleSeat).Error
	return scheduleSeat, err
}
func (r *scheduleSeatRepository) FindBySchedule(ctx context.Context, scheduleID uint) ([]entity.ScheduleSeat, error) {
	var scheduleSeats []entity.ScheduleSeat
	err := r.db.WithContext(ctx).Where("schedule_id = ?", scheduleID).Find(&scheduleSeats).Error
	return scheduleSeats, err
}
func (r *scheduleSeatRepository) WithTx(tx *gorm.DB) repository.ScheduleSeatRepository {
	return &scheduleSeatRepository{db: tx}
}
func (r *scheduleSeatRepository) EnsureAvailableSeats(ctx context.Context, scheduleID uint, seatIDs []uint) error {
	if len(seatIDs) == 0 {
		return nil
	}
	rows := make([]entity.ScheduleSeat, 0, len(seatIDs))
	for _, id := range seatIDs {
		rows = append(rows, entity.ScheduleSeat{ScheduleID: scheduleID, SeatID: id})
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).
		Omit(clause.Associations).Create(&rows).Error
}
func (r *scheduleSeatRepository) LockAvailableSeats(ctx context.Context, scheduleID uint, seatIDs []uint) ([]entity.ScheduleSeat, error) {
	var out []entity.ScheduleSeat
	if len(seatIDs) == 0 {
		return out, nil
	}
	err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("schedule_id = ? AND seat_id IN ? AND status = ?", scheduleID, seatIDs, "available").
		Find(&out).Error
	return out, err
}
func (r *scheduleSeatRepository) SetStatus(ctx context.Context, ids []uint, status string, lockedAt *time.Time) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&entity.ScheduleSeat{}).Where("id IN ?", ids).
		Updates(map[string]interface{}{"status": status, "locked_at": lockedAt}).Error
}
