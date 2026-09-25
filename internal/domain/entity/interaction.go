package entity

import "time"

type Bookmark struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_film_bookmark"`
	FilmID    uint      `json:"film_id" gorm:"not null;uniqueIndex:idx_user_film_bookmark"`
	User      User      `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Film      Film      `json:"film" gorm:"foreignKey:FilmID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
}

type Rating struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_film_rating"`
	FilmID    uint      `json:"film_id" gorm:"not null;uniqueIndex:idx_user_film_rating;index"`
	Score     int       `json:"score" gorm:"type:smallint;not null"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Query     string    `json:"query" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}
