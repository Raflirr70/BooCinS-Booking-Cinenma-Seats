package usecase

import (
	"context"
	"errors"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
	"gorm.io/gorm"
)

type promoUsecase struct {
	promoRepo repository.PromoRepository
}

func NewPromoUsecase(promoRepo repository.PromoRepository) uc.PromoUsecase {
	return &promoUsecase{promoRepo: promoRepo}
}

func (u *promoUsecase) GetByID(ctx context.Context, id uint) (*entity.Promo, error) {
	if id == 0 {
		return nil, uc.ErrPromoNotFound
	}
	promo, err := u.promoRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, uc.ErrPromoNotFound
		}
		return nil, err
	}
	return promo, nil
}
func (u *promoUsecase) Update(ctx context.Context, promo *entity.Promo) error {
	if promo == nil {
		return uc.ErrInvalidPromo
	}
	if promo.ID == 0 {
		return uc.ErrInvalidPromoID
	}
	_, err := u.promoRepo.GetByID(ctx, promo.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uc.ErrPromoNotFound
		}
	}
	return u.promoRepo.Update(ctx, promo)
}
func (u *promoUsecase) Delete(ctx context.Context, promo *entity.Promo) error {
	if promo == nil {
		return uc.ErrInvalidPromo
	}
	if promo.ID == 0 {
		return uc.ErrInvalidPromoID
	}
	_, err := u.promoRepo.GetByID(ctx, promo.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uc.ErrPromoNotFound
		}
	}
	return u.promoRepo.Delete(ctx, promo)
}
func (u *promoUsecase) Create(ctx context.Context, promo *entity.Promo) error {
	if promo == nil {
		return uc.ErrInvalidPromo
	}
	if promo.StartDate.After(promo.EndDate) {
		return uc.ErrInvalidPromoDate
	}
	_, err := u.promoRepo.FindOverlapping(ctx, promo.StartDate, promo.EndDate)
	if err == nil {
		return uc.ErrPromoDateConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return u.promoRepo.Create(ctx, promo)
}
