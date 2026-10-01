package dto

import (
	"time"

	"github.com/rafli/boocins/internal/domain/entity"
)

type UpdatePromoRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	Price       *float64 `json:"price"`
	Discount    float64  `json:"discount"`
	IsActive    bool     `json:"is_active"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
}
type CreatePromoRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	Price       *float64 `json:"price"`
	Discount    float64  `json:"discount"`
	IsActive    bool     `json:"is_active"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
}
type PromoResponse struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`
	Price       *float64 `json:"price"`
	Discount    float64  `json:"discount"`
	IsActive    bool     `json:"is_active"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
}

func (r *UpdatePromoRequest) ToEntityPromo() (*entity.Promo, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, err
	}

	return &entity.Promo{
		Title:       r.Title,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       r.Price,
		Discount:    r.Discount,
		IsActive:    r.IsActive,
		StartDate:   startDate,
		EndDate:     endDate,
	}, err
}
func (r *CreatePromoRequest) ToEntityPromo() (*entity.Promo, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, err
	}
	return &entity.Promo{
		Title:       r.Title,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       r.Price,
		Discount:    r.Discount,
		IsActive:    r.IsActive,
		StartDate:   startDate,
		EndDate:     endDate,
	}, err
}
func ToPromoResponse(promo *entity.Promo) PromoResponse {
	return PromoResponse{
		Title:       promo.Title,
		Description: promo.Description,
		ImageURL:    promo.ImageURL,
		Price:       promo.Price,
		Discount:    promo.Discount,
		IsActive:    promo.IsActive,
		StartDate:   promo.StartDate.Format("2006-01-02"),
		EndDate:     promo.EndDate.Format("2006-01-02"),
	}
}
