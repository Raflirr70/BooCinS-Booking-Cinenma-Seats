package postgres

import (
	"context"
	"time"

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
func (r *promoRepository) Delete(ctx context.Context, promo *entity.Promo) error {
	return r.db.WithContext(ctx).Delete(&promo).Error
}
func (r *promoRepository) Create(ctx context.Context, promo *entity.Promo) error {
	return r.db.WithContext(ctx).Create(&promo).Error
}
func (r *promoRepository) FindOverlapping(ctx context.Context, startDate time.Time, endDate time.Time) (*entity.Promo, error) {
	var promo entity.Promo
	err := r.db.WithContext(ctx).Where("start_date <= ?", endDate).Where("end_date >= ?", startDate).First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}
func (r *promoRepository) GetCurrent(ctx context.Context) (*entity.Promo, error) {
	var promo *entity.Promo
	if err := r.db.WithContext(ctx).Where("start_date <= CURRENT_DATE").Where("end_date >= CURRENT_DATE").First(&promo).Error; err != nil {
		return nil, err
	}
	return promo, nil
}
