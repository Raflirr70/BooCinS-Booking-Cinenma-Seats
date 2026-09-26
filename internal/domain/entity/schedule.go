package entity

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Capacity  int            `json:"capacity" gorm:"not null"`
	Seats     []Seat         `json:"seats,omitempty" gorm:"foreignKey:RoomID"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Seat struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	RoomID uint   `json:"room_id" gorm:"not null"`
	Label  string `json:"label" gorm:"type:varchar(10);not null"`
	Number int    `json:"number" gorm:"not null"`
	Status string `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	Room   Room   `json:"-" gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

type Schedule struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FilmID    uint      `json:"film_id" gorm:"not null"`
	RoomID    uint      `json:"room_id" gorm:"not null"`
	ShowDate  string    `json:"show_date" gorm:"type:date;not null"`
	ShowTime  string    `json:"show_time" gorm:"type:time;not null"`
	Status    bool      `json:"status" gorm:"not null;default:true"`
	Film      Film      `json:"film" gorm:"foreignKey:FilmID"`
	Room      Room      `json:"room" gorm:"foreignKey:RoomID"`
	CreatedAt time.Time `json:"created_at"`
}

type ScheduleSeat struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	ScheduleID uint       `json:"schedule_id" gorm:"not null;uniqueIndex:idx_schedule_seat"`
	SeatID     uint       `json:"seat_id" gorm:"not null;uniqueIndex:idx_schedule_seat"`
	Status     string     `json:"status" gorm:"type:varchar(20);not null;default:'available'"`
	LockedAt   *time.Time `json:"locked_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Schedule   Schedule   `json:"-" gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE"`
	Seat       Seat       `json:"-" gorm:"foreignKey:SeatID"`
}
