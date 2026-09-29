package postgres

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	"gorm.io/gorm"
)

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) repository.PromoRepository {
	return &promoRepository{db: db}
}
func (r *promoRepository) GetByID(ctx context.Context, id uint) (*entity.Promo, error) {
	var promo *entity.Promo
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&promo).Error; err != nil {
		return nil, err
	}
	return promo, nil
}
func (r *promoRepository) Update(ctx context.Context, promo *entity.Promo) error {
	return r.db.WithContext(ctx).Where("id = ?", promo.ID).Updates(&promo).Error
}
