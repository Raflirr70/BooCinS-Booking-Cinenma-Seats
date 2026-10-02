package dto

import (
	"strings"

	"github.com/rafli/boocins/internal/domain/entity"
)

// Halaman Utama
type HalamanUtamaResponse struct {
	UserID uint               `json:"user_id"`
	Name   string             `json:"name"`
	Promo  *PromoHomeResponse `json:"promo"`
	Films  []FilmHomeResponse `json:"films"`
}

// Halaman List Film
type HalamanListFilmResponse struct {
	UserID uint               `json:"user_id"`
	Name   string             `json:"name"`
	Films  []FilmHomeResponse `json:"films"`
}

// Halaman Detaul Films
type HalamanDetailFilmResponse struct {
	UserID uint                `json:"user_id"`
	Name   string              `json:"name"`
	Film   *FilmDetailResponse `json:"film"`
}

type HalamanScheduleMapSeatResponse struct {
	UserID   uint                       `json:"user_id"`
	Name     string                     `json:"name"`
	Film     FilmDetailResponse         `json:"film"`
	Schedule ScheduleSeatDetailResponse `json:"schedule"`
}

type PromoHomeResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"img"`
}

type FilmHomeResponse struct {
	FilmID   uint                   `json:"film_id"`
	Name     string                 `json:"name"`
	Cover    string                 `json:"cover"`
	Duration int                    `json:"duration"`
	Price    float64                `json:"price"`
	Genres   []string               `json:"genres"`
	Schedule []ScheduleHomeResponse `json:"schedule"`
}

type ScheduleHomeResponse struct {
	ScheduleID     uint   `json:"schedule_id"`
	RoomID         uint   `json:"room_id"`
	RoomName       string `json:"room_name"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	NSeats         int    `json:"nseat"`
	BookedSeats    int    `json:"booked_seats"`
	AvailableSeats int    `json:"available_seats"`
	Seats          int    `json:"seats"` //jumlah sheculeseat berdasarkan id schedule
}

func ToHomeResponse(user *entity.User, promo *entity.Promo, films []entity.Film, schedules []entity.Schedule, scheduleSeats map[uint][]entity.ScheduleSeat) HalamanUtamaResponse {
	schedulesByFilm := make(map[uint][]entity.Schedule, len(films))
	for i := range schedules {
		schedulesByFilm[schedules[i].FilmID] = append(schedulesByFilm[schedules[i].FilmID], schedules[i])
	}

	res := HalamanUtamaResponse{
		Promo: ToPromoHomeResponse(promo),
		Films: toHomeFilmResponses(films, schedulesByFilm, scheduleSeats),
	}

	if user != nil {
		res.UserID = user.ID
		res.Name = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return res
}
func ToHalamanListFilm(user *entity.User, films []entity.Film, schedules []entity.Schedule, scheduleSeats map[uint][]entity.ScheduleSeat) HalamanListFilmResponse {
	schedulesByFilm := make(map[uint][]entity.Schedule, len(films))
	for i := range schedules {
		schedulesByFilm[schedules[i].FilmID] = append(schedulesByFilm[schedules[i].FilmID], schedules[i])
	}

	res := HalamanListFilmResponse{
		Films: toHomeFilmResponses(films, schedulesByFilm, scheduleSeats),
	}

	if user != nil {
		res.UserID = user.ID
		res.Name = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return res
}
func ToHalamanDetailFilm(user *entity.User, film *entity.Film, schedules []entity.Schedule, scheduleSeats map[uint][]entity.ScheduleSeat) HalamanDetailFilmResponse {
	var filmResp *FilmDetailResponse

	if film != nil {
		detail := ToFilmDetailResponse(film)

		scheduleResponses := make([]ScheduleHomeResponse, 0, len(schedules))
		for i := range schedules {
			scheduleResponses = append(scheduleResponses, toHomeScheduleResponse(&schedules[i], scheduleSeats[schedules[i].ID]))
		}
		detail.Schedule = scheduleResponses
		filmResp = &detail
	}

	res := HalamanDetailFilmResponse{
		Film: filmResp,
	}

	if user != nil {
		res.UserID = user.ID
		res.Name = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return res
}
func toHomeFilmResponses(films []entity.Film, schedulesByFilm map[uint][]entity.Schedule, scheduleSeats map[uint][]entity.ScheduleSeat) []FilmHomeResponse {
	results := make([]FilmHomeResponse, 0, len(films))
	for i := range films {
		results = append(results, toHomeFilmResponse(&films[i], schedulesByFilm, scheduleSeats))
	}
	return results
}

func toHomeFilmResponse(film *entity.Film, schedulesByFilm map[uint][]entity.Schedule, scheduleSeats map[uint][]entity.ScheduleSeat) FilmHomeResponse {
	genres := make([]string, 0, len(film.Genres))
	for _, g := range film.Genres {
		genres = append(genres, g.Name)
	}

	schedules := schedulesByFilm[film.ID]
	scheduleResponses := make([]ScheduleHomeResponse, 0, len(schedules))
	for i := range schedules {
		scheduleResponses = append(scheduleResponses, toHomeScheduleResponse(&schedules[i], scheduleSeats[schedules[i].ID]))
	}

	return FilmHomeResponse{
		FilmID:   film.ID,
		Name:     film.Title,
		Cover:    film.Cover,
		Duration: film.Duration,
		Price:    film.Price,
		Genres:   genres,
		Schedule: scheduleResponses,
	}
}

func toHomeScheduleResponse(schedule *entity.Schedule, seats []entity.ScheduleSeat) ScheduleHomeResponse {
	nseats := schedule.Room.Capacity
	if len(schedule.Room.Seats) > 0 {
		nseats = len(schedule.Room.Seats)
	}

	booked := 0
	for range seats {
		booked++
	}
	if nseats < booked {
		nseats = booked
	}

	return ScheduleHomeResponse{
		ScheduleID:     schedule.ID,
		RoomID:         schedule.RoomID,
		RoomName:       schedule.Room.Name,
		Date:           schedule.ShowDate,
		Time:           schedule.ShowTime,
		NSeats:         nseats,
		BookedSeats:    booked,
		AvailableSeats: nseats - booked,
		Seats:          len(seats),
	}
}

func ToPromoHomeResponse(promo *entity.Promo) *PromoHomeResponse {
	if promo == nil {
		return nil
	}
	return &PromoHomeResponse{
		ID:          promo.ID,
		Title:       promo.Title,
		Description: promo.Description,
		ImageURL:    promo.ImageURL,
	}
}
func ToHalamanScheduleMapSeat(
	user *entity.User,
	schedule *entity.Schedule,
	film *entity.Film,
	scheduleSeats []entity.ScheduleSeat,
) HalamanScheduleMapSeatResponse {
	existingSeatsMap := make(map[uint]entity.ScheduleSeat, len(scheduleSeats))
	for _, ss := range scheduleSeats {
		existingSeatsMap[ss.SeatID] = ss
	}

	seatsDetail := make([]SeatItemDetail, 0, len(schedule.Room.Seats))
	bookedCount := 0

	for _, seat := range schedule.Room.Seats {
		item := SeatItemDetail{
			SeatID:      seat.ID,
			Label:       seat.Label,
			Number:      seat.Number,
			Status:      "available",
			IsAvailable: true,
		}

		if ss, exists := existingSeatsMap[seat.ID]; exists {
			item.ScheduleSeatID = &ss.ID
			item.Status = ss.Status
			if ss.Status != "available" || ss.LockedAt != nil {
				item.IsAvailable = false
				bookedCount++
			}
		}
		seatsDetail = append(seatsDetail, item)
	}

	totalSeats := len(schedule.Room.Seats)
	if totalSeats == 0 {
		totalSeats = schedule.Room.Capacity
	}
	availableSeats := totalSeats - bookedCount
	if availableSeats < 0 {
		availableSeats = 0
	}

	res := HalamanScheduleMapSeatResponse{
		Film: ToFilmDetailResponse(film),
		Schedule: ScheduleSeatDetailResponse{
			ScheduleID:     schedule.ID,
			ShowDate:       schedule.ShowDate,
			ShowTime:       schedule.ShowTime,
			RoomName:       schedule.Room.Name,
			TotalSeats:     totalSeats,
			AvailableSeats: availableSeats,
			BookedSeats:    bookedCount,
			Seats:          seatsDetail,
		},
	}

	if user != nil {
		res.UserID = user.ID
		res.Name = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return res
}
