package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
)

// FILM

type filmUsecase struct {
	filmRepo  repository.FilmRepository
	genreRepo repository.GenreRepository
}

func NewFilmUsecase(filmRepo repository.FilmRepository, genreRepo repository.GenreRepository) uc.FilmUsecase {
	return &filmUsecase{filmRepo: filmRepo, genreRepo: genreRepo}
}

func (u *filmUsecase) GetAll(ctx context.Context) ([]entity.Film, error) {
	return u.filmRepo.FindAll(ctx)
}

func (u *filmUsecase) GetAllWithDetails(ctx context.Context) ([]entity.Film, error) {
	return u.filmRepo.FindAllWithDetails(ctx)
}

func (u *filmUsecase) Create(ctx context.Context, film *entity.Film) error {
	genres := film.Genres
	film.Genres = nil

	genreIDs := make([]uint, 0, len(genres))
	for _, g := range genres {
		genreIDs = append(genreIDs, g.ID)
	}
	validIDs, err := u.ValidateGenreIDs(ctx, genreIDs)
	if err != nil {
		return err
	}

	if err := u.filmRepo.Create(ctx, film); err != nil {
		return err
	}

	for _, id := range validIDs {
		if err := u.filmRepo.AddGenre(ctx, film.ID, id); err != nil {
			return err
		}
	}

	return nil
}

func (u *filmUsecase) GetByID(ctx context.Context, id uint) (*entity.Film, error) {
	return u.filmRepo.FindByID(ctx, id)
}

func (u *filmUsecase) Update(ctx context.Context, film *entity.Film, genreIDs *[]uint) error {
	var validIDs []uint
	if genreIDs != nil {
		var err error
		validIDs, err = u.ValidateGenreIDs(ctx, *genreIDs)
		if err != nil {
			return err
		}
	}
	if err := u.filmRepo.Update(ctx, film); err != nil {
		return err
	}
	if genreIDs != nil {
		if err := u.filmRepo.ReplaceGenres(ctx, film.ID, validIDs); err != nil {
			return err
		}
	}
	return nil
}

func (u *filmUsecase) Delete(ctx context.Context, id uint) error {
	film := &entity.Film{ID: id}
	return u.filmRepo.Delete(ctx, film)
}

func (u *filmUsecase) ValidateGenreIDs(ctx context.Context, genreIDs []uint) ([]uint, error) {
	seen := make(map[uint]struct{}, len(genreIDs))
	unique := make([]uint, 0, len(genreIDs))
	for _, id := range genreIDs {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}
	if len(unique) == 0 {
		return unique, nil
	}
	found, err := u.genreRepo.FindByIDs(ctx, unique)
	if err != nil {
		return nil, err
	}
	if len(found) != len(unique) {
		return nil, uc.ErrGenreNotFound

	}
	return unique, nil
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

// Media

type mediaUsecase struct {
	mediaRepo repository.MediaRepository
}

func NewMediaUsecase(mediaRepo repository.MediaRepository) uc.MediaUsecase {
	return &mediaUsecase{mediaRepo: mediaRepo}
}
func (u *mediaUsecase) GetByFilm(ctx context.Context, filmID uint) ([]entity.Media, error) {
	return u.mediaRepo.FindByFilmID(ctx, filmID)
}
func (u *mediaUsecase) GetByID(ctx context.Context, id uint) (*entity.Media, error) {
	return u.mediaRepo.FindByID(ctx, id)
}
func (u *mediaUsecase) Update(ctx context.Context, media *entity.Media) error {
	return u.mediaRepo.Update(ctx, media)
}
func (u *mediaUsecase) Delete(ctx context.Context, media *entity.Media) error {
	return u.mediaRepo.Delete(ctx, media)
}
func (u *mediaUsecase) Create(ctx context.Context, media *entity.Media) error {
	return u.mediaRepo.Create(ctx, media)
}
