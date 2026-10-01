package usecase

import (
	"context"
	"errors"

	"github.com/rafli/boocins/internal/domain/entity"
)

var (
	ErrInvalidPromo      = errors.New("invalid promo")
	ErrInvalidPromoID    = errors.New("invalid promo id")
	ErrPromoNotFound     = errors.New("promo not found")
	ErrPromoDateConflict = errors.New("promo date conflicts with another promo")
	ErrInvalidPromoDate  = errors.New("start date must be before or equal to end date")
	ErrPromoExpired      = errors.New("promo has expired")
)

type PromoUsecase interface {
	GetByID(ctx context.Context, id uint) (*entity.Promo, error)
	Update(ctx context.Context, promo *entity.Promo) error
	Delete(ctx context.Context, promo *entity.Promo) error

	Create(ctx context.Context, promo *entity.Promo) error
	GetPromo(ctx context.Context) (*entity.Promo, error)
}
