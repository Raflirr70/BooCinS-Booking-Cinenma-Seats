package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/usecase"
	uc "github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
)

type PromoHandler struct {
	promoUC usecase.PromoUsecase
}

func NewPromoHandler(promoUC usecase.PromoUsecase) *PromoHandler {
	return &PromoHandler{promoUC: promoUC}
}
func (h *PromoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	var req dto.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	promo := &entity.Promo{ID: uint(id)}

	if err := h.promoUC.Update(c.Request.Context(), promo); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("Updated Promo", nil))
}
func (h *PromoHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	promo, err := h.promoUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("fetch promo", dto.ToPromoResponse(promo)))
}
func (h *PromoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	promo, err := h.promoUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	if err := h.promoUC.Delete(c.Request.Context(), promo); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("Deleted Promo", nil))
}
func (h *PromoHandler) Create(c *gin.Context) {
	var req dto.CreatePromoRequest
	fmt.Print("test")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	promo, err := req.ToEntityPromo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	if err := h.promoUC.Create(c.Request.Context(), promo); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("Created Promo", nil))
}
func (h *PromoHandler) GetCurrent(c *gin.Context) {
	promo, err := h.promoUC.GetPromo(c.Request.Context())
	if err != nil {
		if errors.Is(err, uc.ErrPromoNotFound) {
			c.JSON(http.StatusNotFound, response.Error("Promo not Found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("fetch promo", dto.ToPromoResponse(promo)))
}
