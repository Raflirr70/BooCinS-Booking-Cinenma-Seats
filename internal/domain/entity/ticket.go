package entity

import "time"

type GuestOrder struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Email         string    `json:"email" gorm:"type:varchar(255);not null;index"`
	IsVerified    bool      `json:"is_verified" gorm:"not null;default:false"`
	TransactionID *uint     `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type EmailVerification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null;index"`
	OTPCode   string    `json:"otp_code" gorm:"type:varchar(6);not null"`
	IsUsed    bool      `json:"is_used" gorm:"not null;default:false"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        *uint     `json:"user_id"`
	GuestOrderID  *uint     `json:"guest_order_id"`
	Status        string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TotalPrice    float64   `json:"total_price" gorm:"type:decimal(12,2);not null"`
	PaymentMethod *string   `json:"payment_method" gorm:"type:varchar(20)"`
	Source        string    `json:"source" gorm:"type:varchar(20);not null;default:'online'"`
	StaffID       *uint     `json:"staff_id"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Ticket struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         *uint      `json:"user_id"`
	ScheduleSeatID uint       `json:"schedule_seat_id" gorm:"not null"`
	TransactionID  uint       `json:"transaction_id" gorm:"not null"`
	QRToken        string     `json:"qr_token" gorm:"type:uuid;uniqueIndex;not null"`
	Status         string     `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	CheckedInAt    *time.Time `json:"checked_in_at"`
	CheckedInBy    *uint      `json:"checked_in_by"`
	CreatedAt      time.Time  `json:"created_at"`
	ScheduleSeat   ScheduleSeat `json:"schedule_seat" gorm:"foreignKey:ScheduleSeatID"`
	Transaction    Transaction  `json:"transaction" gorm:"foreignKey:TransactionID"`
}
