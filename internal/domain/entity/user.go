package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}

type RolePermission struct {
	RoleID       uint       `json:"role_id" gorm:"primaryKey"`
	PermissionID uint       `json:"permission_id" gorm:"primaryKey"`
	Role         Role       `json:"-" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

type Membership struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	MinScore int     `json:"min_score" gorm:"not null;default:0"`
	Discount float64 `json:"discount" gorm:"type:decimal(5,2);not null;default:0.00"`
}

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password     string         `json:"-" gorm:"type:varchar(255);not null"`
	FirstName    string         `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName     string         `json:"last_name" gorm:"type:varchar(100);not null"`
	Score        int            `json:"score" gorm:"not null;default:0"`
	RoleID       uint           `json:"role_id" gorm:"not null"`
	MembershipID *uint          `json:"membership_id"`
	Role         Role           `json:"role" gorm:"foreignKey:RoleID"`
	Membership   *Membership    `json:"membership,omitempty" gorm:"foreignKey:MembershipID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
