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
	created, err := h.filmUC.GetByID(c.Request.Context(), film.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.Success("film created", dto.ToFilmDetailResponse(created)))
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
	if err := h.filmUC.Update(c.Request.Context(), film, req.GenreIDs); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	updated, err := h.filmUC.GetByID(c.Request.Context(), film.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("film updated", dto.ToFilmDetailResponse(updated)))
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

//Genre

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
	genre := req.ToEntityGenre()
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

// Media

type MediaHandler struct {
	mediaUC usecase.MediaUsecase
}

func NewMediaHandler(mediaUC usecase.MediaUsecase) *MediaHandler {
	return &MediaHandler{mediaUC: mediaUC}
}

func (h *MediaHandler) GetByFilm(c *gin.Context) {
	filmID, err := strconv.ParseUint(c.Param("film_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid film id"))
		return
	}
	medias, err := h.mediaUC.GetByFilm(c.Request.Context(), uint(filmID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("medias fetched", dto.ToMediaResponses(medias)))
}

func (h *MediaHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	media, err := h.mediaUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("media not found"))
		return
	}
	c.JSON(http.StatusOK, response.Success("media fetched", dto.ToMediaResponse(media)))
}

func (h *MediaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	var req dto.UpdateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}
	media, err := h.mediaUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("media not found"))
		return
	}
	if &req.Type != nil {
		media.Type = req.Type
	}
	if &req.URL != nil {
		media.URL = req.URL
	}
	if err := h.mediaUC.Update(c.Request.Context(), media); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("media updated", dto.ToMediaResponse(media)))
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}
	media, err := h.mediaUC.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("media not found"))
		return
	}
	if err := h.mediaUC.Delete(c.Request.Context(), media); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success("media deleted", nil))
}
