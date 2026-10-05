package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
)

type BookingHandler struct {
	bookingUC usecase.BookingUsecase
}

func NewBookingHandler(bookingUC usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUC: bookingUC}
}

// Create: user login — first_name/last_name diambil dari JWT, body diabaikan.
func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	v, exists := c.Get("user_id")
	uid, ok := v.(uint)
	if !exists || !ok || uid == 0 {
		c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
		return
	}
	out, err := h.bookingUC.CreateBooking(c.Request.Context(), usecase.CreateBookingInput{
		ScheduleID: req.ScheduleID, SeatIDs: req.SeatIDs, UserID: &uid,
	})
	if err != nil {
		writeBookingError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Success("booking created", out))
}

// CreateGuest: tanpa login — first_name/last_name wajib di body.
func (h *BookingHandler) CreateGuest(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	out, err := h.bookingUC.CreateBooking(c.Request.Context(), usecase.CreateBookingInput{
		ScheduleID: req.ScheduleID, SeatIDs: req.SeatIDs,
		FirstName: req.FirstName, LastName: req.LastName, Email: req.Email,
	})
	if err != nil {
		writeBookingError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.Success("booking created", out))
}

func writeBookingError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrSeatsUnavailable),
		errors.Is(err, usecase.ErrSeatNotInRoom),
		errors.Is(err, usecase.ErrDuplicateSeats):
		c.JSON(http.StatusConflict, response.Error(err.Error()))
	case errors.Is(err, usecase.ErrScheduleNotFound),
		errors.Is(err, usecase.ErrBuyerNotFound):
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
	case errors.Is(err, usecase.ErrPaymentFailed):
		c.JSON(http.StatusBadGateway, response.Error(err.Error()))
	case errors.Is(err, usecase.ErrInvalidBookingInput),
		errors.Is(err, usecase.ErrGuestNameRequired),
		errors.Is(err, usecase.ErrInvalidTotal):
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}
}
