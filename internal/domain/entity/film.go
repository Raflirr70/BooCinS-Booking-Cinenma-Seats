package entity

import (
	"time"

	"gorm.io/gorm"
)

type Film struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"type:varchar(255);not null"`
	Cover       string         `json:"cover" gorm:"type:varchar(255);not null"`
	Synopsis    string         `json:"synopsis" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	Director    string         `json:"director" gorm:"type:varchar(150)"`
	Duration    int            `json:"duration" gorm:"not null"`
	Price       float64        `json:"price" gorm:"type:decimal(12,2);not null"`
	Status      string         `json:"status" gorm:"type:varchar(20);not null;default:'coming_soon'"`
	Genres      []Genre        `json:"genres,omitempty" gorm:"many2many:film_genres;"`
	Media       []Media        `json:"media,omitempty" gorm:"foreignKey:FilmID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Genre struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Media struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FilmID      uint      `json:"film_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"type:varchar(20);not null;default:'poster'"`
	URL         string    `json:"url" gorm:"type:text;not null"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   time.Time `json:"created_at"`
}
