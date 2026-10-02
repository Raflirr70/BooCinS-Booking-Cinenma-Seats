package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
)

type HomeHandler struct {
	authUC         usecase.AuthUsecase
	filmUC         usecase.FilmUsecase
	scheduleUC     usecase.ScheduleUsecase
	scheduleSeatUC usecase.ScheduleSeatUsecase
	promoUC        usecase.PromoUsecase
}

func NewHomeHandler(
	authUC usecase.AuthUsecase,
	filmUC usecase.FilmUsecase,
	scheduleUC usecase.ScheduleUsecase,
	scheduleSeatUC usecase.ScheduleSeatUsecase,
	promoUC usecase.PromoUsecase,
) *HomeHandler {
	return &HomeHandler{
		authUC:         authUC,
		filmUC:         filmUC,
		scheduleUC:     scheduleUC,
		scheduleSeatUC: scheduleSeatUC,
		promoUC:        promoUC,
	}
}

func (h *HomeHandler) GetHome(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. User (opsional, jika token disertakan)
	var user *entity.User
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok && id > 0 {
			u, err := h.authUC.GetProfile(ctx, id)
			if err == nil {
				user = u
			}
		}
	}

	// 2. Promo aktif
	promo, _ := h.promoUC.GetPromo(ctx)

	// 3. Film + Genres
	films, err := h.filmUC.GetAllWithDetails(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 4. Schedule
	schedules, err := h.scheduleUC.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 5. ScheduleSeats per schedule
	scheduleSeats := make(map[uint][]entity.ScheduleSeat)
	for _, s := range schedules {
		seats, err := h.scheduleSeatUC.GetBySchedule(ctx, s.ID)
		if err == nil && len(seats) > 0 {
			scheduleSeats[s.ID] = seats
		}
	}

	// 6. Response via DTO
	res := dto.ToHomeResponse(user, promo, films, schedules, scheduleSeats)
	c.JSON(http.StatusOK, response.Success("home fetched", res))
}
func (h *HomeHandler) GetFilms(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. User (opsional, jika token disertakan)
	var user *entity.User
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok && id > 0 {
			u, err := h.authUC.GetProfile(ctx, id)
			if err == nil {
				user = u
			}
		}
	}

	// 3. Film + Genres
	films, err := h.filmUC.GetAllWithDetails(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 4. Schedule
	schedules, err := h.scheduleUC.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 5. ScheduleSeats per schedule
	scheduleSeats := make(map[uint][]entity.ScheduleSeat)
	for _, s := range schedules {
		seats, err := h.scheduleSeatUC.GetBySchedule(ctx, s.ID)
		if err == nil && len(seats) > 0 {
			scheduleSeats[s.ID] = seats
		}
	}

	// 6. Response via DTO
	res := dto.ToHalamanListFilm(user, films, schedules, scheduleSeats)
	c.JSON(http.StatusOK, response.Success("home fetched", res))
}
func (h *HomeHandler) GetDetailFilms(c *gin.Context) {
	ctx := c.Request.Context()
	idf, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	// 1. User (opsional, jika token disertakan)
	var user *entity.User
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok && id > 0 {
			u, err := h.authUC.GetProfile(ctx, id)
			if err == nil {
				user = u
			}
		}
	}

	// 3. Film + Genres
	film, err := h.filmUC.GetByID(ctx, uint(idf))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 4. Schedule
	schedules, err := h.scheduleUC.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	// 5. ScheduleSeats per schedule
	scheduleSeats := make(map[uint][]entity.ScheduleSeat)
	for _, s := range schedules {
		seats, err := h.scheduleSeatUC.GetBySchedule(ctx, s.ID)
		if err == nil && len(seats) > 0 {
			scheduleSeats[s.ID] = seats
		}
	}

	// 6. Response via DTO
	res := dto.ToHalamanDetailFilm(user, film, schedules, scheduleSeats)
	c.JSON(http.StatusOK, response.Success("home fetched", res))
}
