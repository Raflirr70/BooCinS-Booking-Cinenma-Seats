package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type FilmUsecase interface {
	GetAll(ctx context.Context) ([]entity.Film, error)
	GetAllWithDetails(ctx context.Context) ([]entity.Film, error)
	Create(ctx context.Context, film *entity.Film) error
	GetByID(ctx context.Context, id uint) (*entity.Film, error)
	Update(ctx context.Context, film *entity.Film) error
	Delete(ctx context.Context, id uint) error
}

type GenreUsecase interface {
	GetAll(ctx context.Context) ([]entity.Genre, error)
	Create(ctx context.Context, genre *entity.Genre) error
	GetByID(ctx context.Context, id uint) (*entity.Genre, error)
	Update(ctx context.Context, genre *entity.Genre) error
	Delete(ctx context.Context, genre *entity.Genre) error
}
