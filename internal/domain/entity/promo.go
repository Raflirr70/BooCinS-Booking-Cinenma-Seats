package entity

import "time"

type Promo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	ImageURL    string    `json:"image_url" gorm:"type:text"`
	Price       *float64  `json:"price" gorm:"type:decimal(12,2)"`
	Discount    float64   `json:"discount" gorm:"type:decimal(5,2);default:0.00"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	StartDate   *string   `json:"start_date" gorm:"type:date"`
	EndDate     *string   `json:"end_date" gorm:"type:date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *uint     `json:"user_id"`
	Action    string    `json:"action" gorm:"type:varchar(20);not null"`
	TableName string    `json:"table_name" gorm:"type:varchar(100);not null;index"`
	RecordID  *uint     `json:"record_id"`
	OldData   *string   `json:"old_data" gorm:"type:jsonb"`
	NewData   *string   `json:"new_data" gorm:"type:jsonb"`
	IPAddress *string   `json:"ip_address" gorm:"type:varchar(45)"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}
