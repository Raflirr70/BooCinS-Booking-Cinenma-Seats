package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
)

// FILM

type filmUsecase struct {
	filmRepo repository.FilmRepository
}

func NewFilmUsecase(filmRepo repository.FilmRepository) uc.FilmUsecase {
	return &filmUsecase{filmRepo: filmRepo}
}

func (u *filmUsecase) GetAll(ctx context.Context) ([]entity.Film, error) {
	return u.filmRepo.FindAll(ctx)
}

func (u *filmUsecase) GetAllWithDetails(ctx context.Context) ([]entity.Film, error) {
	return u.filmRepo.FindAllWithDetails(ctx)
}

func (u *filmUsecase) Create(ctx context.Context, film *entity.Film) error {
	return u.filmRepo.Create(ctx, film)
}

func (u *filmUsecase) GetByID(ctx context.Context, id uint) (*entity.Film, error) {
	return u.filmRepo.FindByID(ctx, id)
}

func (u *filmUsecase) Update(ctx context.Context, film *entity.Film) error {
	return u.filmRepo.Update(ctx, film)
}

func (u *filmUsecase) Delete(ctx context.Context, id uint) error {
	film := &entity.Film{ID: id}
	return u.filmRepo.Delete(ctx, film)
}

// Genre

type genreUsecase struct {
	genreRepo repository.GenreRepository
}

func NewGenreUsecase(genreRepo repository.GenreRepository) uc.GenreUsecase {
	return &genreUsecase{genreRepo: genreRepo}
}

func (u *genreUsecase) GetAll(ctx context.Context) ([]entity.Genre, error) {
	return u.genreRepo.FindAll(ctx)
}
func (u *genreUsecase) Create(ctx context.Context, genre *entity.Genre) error {
	return u.genreRepo.Create(ctx, genre)
}
func (u *genreUsecase) GetByID(ctx context.Context, id uint) (*entity.Genre, error) {
	return u.genreRepo.FindByID(ctx, id)
}
func (u *genreUsecase) Update(ctx context.Context, genre *entity.Genre) error {
	return u.genreRepo.Update(ctx, genre)
}
func (u *genreUsecase) Delete(ctx context.Context, genre *entity.Genre) error {
	return u.genreRepo.Delete(ctx, genre)
}
