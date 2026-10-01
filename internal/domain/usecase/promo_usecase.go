package usecase

import (
	"context"
	"errors"

	"github.com/rafli/boocins/internal/domain/entity"
)

var ErrInvalidPromoID = errors.New("invalid promo id")
var ErrInvalidPromo = errors.New("invalid promo")
var ErrPromoNotFound = errors.New("promo not found")

type PromoUsecase interface {
	GetByID(ctx context.Context, id uint) (*entity.Promo, error)
	Update(ctx context.Context, promo *entity.Promo) error
	Delete(ctx context.Context, promor *entity.Promo) error
}
