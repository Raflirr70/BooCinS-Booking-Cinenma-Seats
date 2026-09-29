package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
)

type promoUsecase struct {
	promoRepo repository.PromoRepository
}

func NewPromoUsecase(promoRepo repository.PromoRepository) uc.PromoUsecase {
	return &promoUsecase{promoRepo: promoRepo}
}

func (u *promoUsecase) GetByID(ctx context.Context, id uint) (*entity.Promo, error) {
	return u.promoRepo.GetByID(ctx, id)
}
func (u *promoUsecase) Update(ctx context.Context, promo *entity.Promo) error {
	return u.promoRepo.Update(ctx, promo)
}
