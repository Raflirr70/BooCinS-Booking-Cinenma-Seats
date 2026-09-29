package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type PromoUsecase interface {
	GetByID(ctx context.Context, id uint) (*entity.Promo, error)
	Update(ctx context.Context, promo *entity.Promo) error
}
