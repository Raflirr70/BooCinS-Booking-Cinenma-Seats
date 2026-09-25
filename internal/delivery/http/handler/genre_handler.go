package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafli/boocins/internal/delivery/http/dto"
	"github.com/rafli/boocins/internal/domain/usecase"
	"github.com/rafli/boocins/pkg/response"
)

type GenreHandler struct {
	genreUC usecase.GenreUsecase
}

func NewGenreHandler(genreUC usecase.GenreUsecase) *GenreHandler {
	return &GenreHandler{genreUC: genreUC}
}

func (h *GenreHandler) GetAll(c *gin.Context) {
	genres, err := h.genreUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("genres fetched", dto.ToGenreResponses(genres)))
}

func (h *GenreHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	genre, err := h.genreUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("genre not found"))
		return
	}
	c.JSON(http.StatusOK, response.Success("genre fetched", dto.ToGenreResponse(genre)))
}
func (h *GenreHandler) Create(c *gin.Context) {
	var req dto.CreateGenreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	genre := req.ToEntity()
	err := h.genreUC.Create(c.Request.Context(), genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("genre created", dto.ToGenreResponse(genre)))
}

func (h *GenreHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	var req dto.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	genre, err := h.genreUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("genre not found"))
		return
	}
	if &req.Name != nil {
		genre.Name = req.Name
	}
	if err := h.genreUC.Update(c.Request.Context(), genre); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("genre updated", dto.ToGenreResponse(genre)))
}

func (h *GenreHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	genre, err := h.genreUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("genre not found"))
		return
	}
	if err := h.genreUC.Delete(c.Request.Context(), genre); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("genre deleted", nil))
}
