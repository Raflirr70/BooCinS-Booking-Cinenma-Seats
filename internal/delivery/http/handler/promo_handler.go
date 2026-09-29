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
	}
	c.JSON(http.StatusOK, response.Success("fetch promo", dto.ToPromoResponse(promo)))
}
