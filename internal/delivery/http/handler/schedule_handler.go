package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
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
	totalSeats := 0
	for _, n := range req.SeatRows {
		totalSeats += n
	}
	room.Capacity = totalSeats
	if err := h.roomUC.Create(c.Request.Context(), room, req.SeatRows); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("room create", room))
}
