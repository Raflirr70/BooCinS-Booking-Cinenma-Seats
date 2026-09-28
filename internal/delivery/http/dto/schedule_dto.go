package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rafli/boocins/internal/domain/entity"
)

type SeatRows map[string]int

type CreateRoomRequest struct {
	Name     string   `json:"name" binding:"required"`
	SeatRows SeatRows `json:"seat_rows" binding:"required"`
}
type CreateScheduleRequest struct {
	FilmID   uint   `json:"film_id" binding:"required"`
	RoomID   uint   `json:"room_id" binding:"required"`
	ShowDate string `json:"show_date" binding:"required"`
	ShowTime string `json:"show_time" binding:"required"`
	Status   bool   `json:"status" binding:"required"`
}
type ScheduleResponse struct {
	FilmID    uint      `json:"film_id" gorm:"not null"`
	RoomID    uint      `json:"room_id" gorm:"not null"`
	ShowDate  string    `json:"show_date" gorm:"type:date;not null"`
	ShowTime  string    `json:"show_time" gorm:"type:time;not null"`
	Status    bool      `json:"status" gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at"`
}
type RoomResponse struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
type RoomWithDetailResponse struct {
	Name     string        `json:"name"`
	Capacity int           `json:"capacity"`
	Seats    []SeatRespone `json:"seat"`
}
type UpdateRoomRequest struct {
	Name     string   `json:"name"`
	SeatRows SeatRows `json:"seat_rows"`
}
type SeatRespone struct {
	Label  string `json:"label"`
	Number int    `json:"number"`
	Status string `json:"status"`
}

func (r *CreateRoomRequest) ToEntityRoom() *entity.Room {
	return &entity.Room{
		Name: r.Name,
	}
}
func (r *UpdateRoomRequest) ToEntityRoom() *entity.Room {
	return &entity.Room{
		Name: r.Name,
	}
}
func ToSeatListResponse(seat *entity.Seat) SeatRespone {
	return SeatRespone{
		Label:  seat.Label,
		Number: seat.Number,
		Status: seat.Status,
	}
}
func (r *CreateScheduleRequest) ToEntitySchedule() *entity.Schedule {
	return &entity.Schedule{
		FilmID:   r.FilmID,
		RoomID:   r.RoomID,
		ShowDate: r.ShowDate,
		ShowTime: r.ShowTime,
		Status:   r.Status,
	}
}

func ToRoomResponse(room *entity.Room) RoomResponse {
	return RoomResponse{
		Name:     room.Name,
		Capacity: room.Capacity,
	}
}
func ToRoomListResponse(room []entity.Room) []RoomResponse {
	rooms := make([]RoomResponse, len(room))
	for i, room := range room {
		rooms[i] = ToRoomResponse(&room)
	}
	return rooms
}

func ToRoomWithDetailListResponse(room *entity.Room) RoomWithDetailResponse {
	seats := make([]SeatRespone, len(room.Seats))
	for i, seat := range room.Seats {
		seats[i] = ToSeatListResponse(&seat)
	}
	return RoomWithDetailResponse{
		Name:     room.Name,
		Capacity: room.Capacity,
		Seats:    seats,
	}
}

func (m *SeatRows) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if _, err := dec.Token(); err != nil {
		return err
	}
	if *m == nil {
		*m = SeatRows{}
	}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := tok.(string)
		if !ok {
			return fmt.Errorf("seat_rows: kunci baris harus string")
		}
		if _, dup := (*m)[key]; dup {
			return fmt.Errorf("baris seat duplikat: %q", key)
		}
		var val int
		if err := dec.Decode(&val); err != nil {
			return err
		}
		if val <= 0 {
			return fmt.Errorf("jumlah seat baris %q harus lebih dari 0", key)
		}
		(*m)[key] = val
	}
	return nil
}
