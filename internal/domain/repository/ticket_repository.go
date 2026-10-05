package repository

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	UpdateStatus(ctx context.Context, id uint, status string) error
	UpdateSnapToken(ctx context.Context, id uint, snapToken string) error
	WithTx(tx *gorm.DB) TransactionRepository
}

type TicketRepository interface {
	CreateMany(ctx context.Context, tickets []entity.Ticket) error
	CancelByTransaction(ctx context.Context, transactionID uint) error
	WithTx(tx *gorm.DB) TicketRepository
}

type GuestOrderRepository interface {
	Create(ctx context.Context, guest *entity.GuestOrder) error
	SetTransaction(ctx context.Context, guestOrderID uint, transactionID uint) error
	WithTx(tx *gorm.DB) GuestOrderRepository
}
