package dto

import (
	"github.com/rafli/boocins/internal/domain/entity"
)

// Request DTOs
// ============================ CREATE REQUEST ============================
type CreateFilmRequest struct {
	Title       string               `json:"title" binding:"required"`
	Cover       string               `json:"cover" binding:"required"`
	Synopsis    string               `json:"synopsis"`
	Description string               `json:"description"`
	Director    string               `json:"director"`
	Duration    int                  `json:"duration" binding:"required,gt=0"`
	Price       float64              `json:"price" binding:"required,gt=0"`
	Status      string               `json:"status" binding:"required,oneof=coming_soon now_showing finished"`
	GenreIDs    []uint               `json:"genre_ids"`
	Media       []CreateMediaRequest `json:"media"`
}
type CreateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}
type CreateMediaRequest struct {
	Type string `json:"type" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

// ============================ UPDATE REQUEST ============================
type UpdateFilmRequest struct {
	Title       *string               `json:"title"`
	Cover       *string               `json:"cover"`
	Synopsis    *string               `json:"synopsis"`
	Description *string               `json:"description"`
	Director    *string               `json:"director"`
	Duration    *int                  `json:"duration" binding:"omitempty,gt=0"`
	Price       *float64              `json:"price" binding:"omitempty,gt=0"`
	Status      *string               `json:"status" binding:"omitempty,oneof=coming_soon now_showing finished"`
	GenreIDs    *[]uint               `json:"genre_ids"`
	Media       []*UpdateMediaRequest `json:"media"`
}
type UpdateMediaRequest struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ============================ RESPONSE ============================
// FilmListResponse — untuk admin (ringkas)
type FilmListResponse struct {
	ID       uint    `json:"id"`
	Cover    string  `json:"cover"`
	Title    string  `json:"title"`
	Duration int     `json:"duration"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
}

// FilmDetailResponse — untuk visitor (lengkap)
type FilmDetailResponse struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Cover       string          `json:"cover"`
	Synopsis    string          `json:"synopsis"`
	Description string          `json:"description"`
	Director    string          `json:"director"`
	Duration    int             `json:"duration"`
	Price       float64         `json:"price"`
	Status      string          `json:"status"`
	Genres      []GenreResponse `json:"genres"`
	Media       []MediaResponse `json:"media"`
}
type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type MediaResponse struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ============================ ToEntity ============================
func (r *CreateFilmRequest) ToEntityFilm() *entity.Film {
	genres := make([]entity.Genre, len(r.GenreIDs))
	for i, id := range r.GenreIDs {
		genres[i] = entity.Genre{ID: id}
	}
	media := make([]entity.Media, len(r.Media))
	for i, m := range r.Media {
		media[i] = entity.Media{Type: m.Type, URL: m.URL}
	}

	return &entity.Film{
		Title:       r.Title,
		Cover:       r.Cover,
		Synopsis:    r.Synopsis,
		Description: r.Description,
		Director:    r.Director,
		Duration:    r.Duration,
		Price:       r.Price,
		Status:      r.Status,
		Genres:      genres,
		Media:       media,
	}
}
func (r *CreateGenreRequest) ToEntityGenre() *entity.Genre {
	return &entity.Genre{
		Name: r.Name,
	}
}
func (r *CreateMediaRequest) ToEntityMedia() *entity.Media {
	return &entity.Media{
		Type: r.Type,
		URL:  r.URL,
	}
}

// ============================ ToResponse ============================
func ToFilmListResponse(film *entity.Film) FilmListResponse {
	return FilmListResponse{
		ID:       film.ID,
		Cover:    film.Cover,
		Title:    film.Title,
		Duration: film.Duration,
		Price:    film.Price,
		Status:   film.Status,
	}
}
func ToFilmListResponses(films []entity.Film) []FilmListResponse {
	results := make([]FilmListResponse, len(films))
	for i, film := range films {
		results[i] = ToFilmListResponse(&film)
	}
	return results
}
func ToFilmDetailResponse(film *entity.Film) FilmDetailResponse {
	genres := make([]GenreResponse, len(film.Genres))
	for i, g := range film.Genres {
		genres[i] = GenreResponse{ID: g.ID, Name: g.Name}
	}
	media := make([]MediaResponse, len(film.Media))
	for i, m := range film.Media {
		media[i] = MediaResponse{ID: m.ID, Type: m.Type, URL: m.URL}
	}
	return FilmDetailResponse{
		ID:          film.ID,
		Title:       film.Title,
		Cover:       film.Cover,
		Synopsis:    film.Synopsis,
		Description: film.Description,
		Director:    film.Director,
		Duration:    film.Duration,
		Price:       film.Price,
		Status:      film.Status,
		Genres:      genres,
		Media:       media,
	}
}
func ToFilmDetailResponses(films []entity.Film) []FilmDetailResponse {
	results := make([]FilmDetailResponse, len(films))
	for i, film := range films {
		results[i] = ToFilmDetailResponse(&film)
	}
	return results
}
func ToGenreResponse(genre *entity.Genre) GenreResponse {
	return GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}
}
func ToGenreResponses(genre []entity.Genre) []GenreResponse {
	results := make([]GenreResponse, len(genre))
	for i, genre := range genre {
		results[i] = ToGenreResponse(&genre)
	}
	return results
}
func ToMediaResponse(media *entity.Media) MediaResponse {
	return MediaResponse{
		ID:   media.ID,
		Type: media.Type,
		URL:  media.URL,
	}
}
func ToMediaResponses(media []entity.Media) []MediaResponse {
	results := make([]MediaResponse, len(media))
	for i, media := range media {
		results[i] = ToMediaResponse(&media)
	}
	return results
}
