package postgres

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	"gorm.io/gorm"
)

type filmRepository struct {
	db *gorm.DB
}

func NewFilmRepository(db *gorm.DB) repository.FilmRepository {
	return &filmRepository{db: db}
}

func (r *filmRepository) Create(ctx context.Context, film *entity.Film) error {
	return r.db.WithContext(ctx).Create(film).Error
}
func (r *filmRepository) Update(ctx context.Context, film *entity.Film) error {
	return r.db.WithContext(ctx).Where("id = ?", film.ID).Updates(film).Error
}
func (r *filmRepository) Delete(ctx context.Context, film *entity.Film) error {
	return r.db.WithContext(ctx).Delete(film).Error
}
func (r *filmRepository) FindAll(ctx context.Context) ([]entity.Film, error) {
	var films []entity.Film
	err := r.db.WithContext(ctx).Find(&films).Error
	return films, err
}
func (r *filmRepository) FindAllWithDetails(ctx context.Context) ([]entity.Film, error) {
	var films []entity.Film
	err := r.db.WithContext(ctx).Preload("Genres").Preload("Media").Find(&films).Error
	return films, err
}
func (r *filmRepository) FindByID(ctx context.Context, id uint) (*entity.Film, error) {
	var film entity.Film
	err := r.db.WithContext(ctx).First(&film, id).Error
	if err != nil {
		return nil, err
	}
	return &film, nil
}
func (r *filmRepository) FindByGenre(ctx context.Context, genre string) ([]entity.Film, error) {
	var films []entity.Film

	err := r.db.WithContext(ctx).
		Joins("JOIN film_genres ON film_genres.film_id = films.id").
		Joins("JOIN genres ON genres.id = film_genres.genre_id").
		Where("genres.name = ?", genre).
		Find(&films).Error
	return films, err
}
func (r *filmRepository) AddGenre(ctx context.Context, filmID uint, genreID uint) error {
	film := entity.Film{ID: filmID}
	genre := entity.Genre{ID: genreID}
	return r.db.WithContext(ctx).Model(&film).Association("Genres").Append(&genre)
}
func (r *filmRepository) RemoveGenre(ctx context.Context, filmID uint, genreID uint) error {
	film := entity.Film{ID: filmID}
	genre := entity.Genre{ID: genreID}
	return r.db.WithContext(ctx).Model(&film).Association("Genres").Delete(&genre)
}

// Genre

type genreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) repository.GenreRepository {
	return &genreRepository{db: db}
}
func (r *genreRepository) Create(ctx context.Context, genre *entity.Genre) error {
	return r.db.WithContext(ctx).Create(genre).Error
}
func (r *genreRepository) Update(ctx context.Context, genre *entity.Genre) error {
	return r.db.WithContext(ctx).Where("id = ?", genre.ID).Updates(genre).Error
}
func (r *genreRepository) Delete(ctx context.Context, genre *entity.Genre) error {
	return r.db.WithContext(ctx).Delete(genre).Error
}
func (r *genreRepository) FindAll(ctx context.Context) ([]entity.Genre, error) {
	var genres []entity.Genre
	err := r.db.WithContext(ctx).Find(&genres).Error
	return genres, err
}
func (r *genreRepository) FindByID(ctx context.Context, id uint) (*entity.Genre, error) {
	var genre entity.Genre
	err := r.db.WithContext(ctx).First(&genre, id).Error
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

// Media

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) repository.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *entity.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}
func (r *mediaRepository) Update(ctx context.Context, media *entity.Media) error {
	return r.db.WithContext(ctx).Where("id = ?", media.ID).Updates(media).Error
}
func (r *mediaRepository) Delete(ctx context.Context, media *entity.Media) error {
	return r.db.WithContext(ctx).Delete(media).Error
}
func (r *mediaRepository) FindByID(ctx context.Context, id uint) (*entity.Media, error) {
	var media entity.Media
	err := r.db.WithContext(ctx).First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}
func (r *mediaRepository) FindByFilmID(ctx context.Context, filmID uint) ([]entity.Media, error) {
	var medias []entity.Media
	err := r.db.WithContext(ctx).Where("film_id = ?", filmID).Find(&medias).Error
	if err != nil {
		return nil, err
	}
	return medias, nil
}
