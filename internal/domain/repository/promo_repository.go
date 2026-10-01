package repository

import (
	"context"
	"time"

	"github.com/rafli/boocins/internal/domain/entity"
)

type PromoRepository interface {
	GetByID(ctx context.Context, id uint) (*entity.Promo, error)
	Update(ctx context.Context, promo *entity.Promo) error
	Delete(ctx context.Context, promo *entity.Promo) error

	Create(ctx context.Context, promo *entity.Promo) error
	// GetCurrent(ctx context.Context) (*entity.Promo, error)
	FindOverlapping(ctx context.Context, startDate time.Time, endDate time.Time) (*entity.Promo, error)
	// FindOverlappingExceptID(ctx context.Context, id uint, startDate time.Time, endDate time.Time) (*entity.Promo, error)
	// DeleteExpired(ctx context.Context, date time.Time) error
}
