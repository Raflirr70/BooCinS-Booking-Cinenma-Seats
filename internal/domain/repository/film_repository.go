package repository

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type FilmRepository interface {
	Create(ctx context.Context, film *entity.Film) error
	Update(ctx context.Context, film *entity.Film) error
	Delete(ctx context.Context, film *entity.Film) error

	FindAll(ctx context.Context) ([]entity.Film, error)
	FindByID(ctx context.Context, id uint) (*entity.Film, error)
	FindByGenre(ctx context.Context, genre string) ([]entity.Film, error)

	AddGenre(ctx context.Context, filmID uint, genreID uint) error
	RemoveGenre(ctx context.Context, filmID uint, genreID uint) error
}

type GenreRepository interface {
	Create(ctx context.Context, genre *entity.Genre) error
	Update(ctx context.Context, genre *entity.Genre) error
	Delete(ctx context.Context, genre *entity.Genre) error

	FindAll(ctx context.Context) ([]entity.Genre, error)
	FindByID(ctx context.Context, id uint) (*entity.Genre, error)
}

type MediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	Update(ctx context.Context, media *entity.Media) error
	Delete(ctx context.Context, media *entity.Media) error

	FindByID(ctx context.Context, id uint) (*entity.Media, error)
	FindByFilmID(ctx context.Context, filmID uint) ([]entity.Media, error)
}
