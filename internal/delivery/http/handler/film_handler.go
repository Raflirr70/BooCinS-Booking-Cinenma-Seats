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

type FilmHandler struct {
	filmUC usecase.FilmUsecase
}

func NewFilmHandler(filmUC usecase.FilmUsecase) *FilmHandler {
	return &FilmHandler{filmUC: filmUC}
}

func (h *FilmHandler) GetAll(c *gin.Context) {
	films, err := h.filmUC.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("films fetched", dto.ToFilmListResponses(films)))
}

func (h *FilmHandler) GetAllWithDetails(c *gin.Context) {
	films, err := h.filmUC.GetAllWithDetails(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("films fetched", dto.ToFilmDetailResponses(films)))
}

func (h *FilmHandler) Create(c *gin.Context) {
	var req dto.CreateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	film := req.ToEntityFilm()
	if err := h.filmUC.Create(c.Request.Context(), film); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("film created", dto.ToFilmDetailResponse(film)))
}

func (h *FilmHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	film, err := h.filmUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("film not found"))
		return
	}
	c.JSON(http.StatusOK, response.Success("film fetched", dto.ToFilmDetailResponse(film)))
}

func (h *FilmHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req dto.UpdateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	film := &entity.Film{ID: uint(id)}
	if req.Title != nil {
		film.Title = *req.Title
	}
	if req.Cover != nil {
		film.Cover = *req.Cover
	}
	if req.Synopsis != nil {
		film.Synopsis = *req.Synopsis
	}
	if req.Description != nil {
		film.Description = *req.Description
	}
	if req.Director != nil {
		film.Director = *req.Director
	}
	if req.Duration != nil {
		film.Duration = *req.Duration
	}
	if req.Price != nil {
		film.Price = *req.Price
	}
	if req.Status != nil {
		film.Status = *req.Status
	}

	if err := h.filmUC.Update(c.Request.Context(), film); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("film updated", dto.ToFilmDetailResponse(film)))
}

func (h *FilmHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.filmUC.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("film deleted", nil))
}
