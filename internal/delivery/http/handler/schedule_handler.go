package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
	"gorm.io/gorm"
)

type RoomHandler struct {
	roomUC usecase.RoomUsecase
}

func NewRoomHandler(roomUC usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{roomUC: roomUC}
}

func (h *RoomHandler) Create(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	room := req.ToEntityRoom()
	if err := h.roomUC.Create(c.Request.Context(), room, req.SeatRows); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("room create", room))
}

func (h *RoomHandler) GetAll(c *gin.Context) {
	rooms, err := h.roomUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	c.JSON(http.StatusOK, response.Success("Room feched", dto.ToRoomListResponse(rooms)))
}

func (h *RoomHandler) GetWithDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	room, err := h.roomUC.GetSeats(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.Error("room not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("seats fetched", dto.ToRoomWithDetailListResponse(room)))
}

func (h *RoomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	room := req.ToEntityRoom()
	room.ID = uint(id)
	if err := h.roomUC.Update(c.Request.Context(), room, req.SeatRows); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("room updated", nil))
}
func (h *RoomHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	room, err := h.roomUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	if err := h.roomUC.Delete(c.Request.Context(), room); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("Room Delete", nil))
}

type ScheduleHandler struct {
	scheduleUC usecase.ScheduleUsecase
}

func NewScheduleHandler(scheduleUC usecase.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{scheduleUC: scheduleUC}
}
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	schedule := req.ToEntitySchedule()
	if err := h.scheduleUC.Create(c.Request.Context(), schedule); err != nil {
		if errors.Is(err, usecase.ErrFilmNotFound) || errors.Is(err, usecase.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("Created Schedule", nil))
}
func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	var req dto.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	schedule, err := h.scheduleUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	if &req.FilmID != nil {
		schedule.FilmID = req.FilmID
	}
	if &req.RoomID != nil {
		schedule.RoomID = req.RoomID
	}
	if &req.ShowDate != nil {
		schedule.ShowDate = req.ShowDate
	}
	if &req.ShowTime != nil {
		schedule.ShowTime = req.ShowTime
	}
	if err := h.scheduleUC.Update(c.Request.Context(), schedule); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("Updated schedule", dto.ToScheduleResponse(schedule)))
}
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	schedule, err := h.scheduleUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	if err := h.scheduleUC.Delete(c.Request.Context(), schedule); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	c.JSON(http.StatusOK, response.Success("Deleted Schedule", nil))
}
func (h *ScheduleHandler) GetAll(c *gin.Context) {
	schedule, err := h.scheduleUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
	res := make([]dto.ScheduleAdminResponse, 0, len(schedule))
	for i := range schedule {
		item := dto.ToScheduleAdminListResponse(&schedule[i])
		res = append(res, item)
	}
	c.JSON(http.StatusOK, response.Success("Feched Schedule", res))
}
func (h *ScheduleHandler) GetByFilm(c *gin.Context) {
	filmID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	schedule, err := h.scheduleUC.GetByFilm(c.Request.Context(), uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	res := make([]dto.ScheduleAdminResponse, 0, len(schedule))
	for i := range schedule {
		item := dto.ToScheduleAdminListResponse(&schedule[i])
		res = append(res, item)
	}
	c.JSON(http.StatusOK, response.Success("fecth schedule film", res))
}
func (h *ScheduleHandler) GetByRoom(c *gin.Context) {
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	schedule, err := h.scheduleUC.GetByRoom(c.Request.Context(), uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	res := make([]dto.ScheduleAdminResponse, 0, len(schedule))
	for i := range schedule {
		item := dto.ToScheduleAdminListResponse(&schedule[i])
		res = append(res, item)
	}
	c.JSON(http.StatusOK, response.Success("fecth schedule film", res))
}
