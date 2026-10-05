package postgres

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type transactionRepository struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}
func (r *transactionRepository) WithTx(tx *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{db: tx}
}
func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(transaction).Error
}
func (r *transactionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("id = ?", id).
		Update("status", status).Error
}
func (r *transactionRepository) UpdateSnapToken(ctx context.Context, id uint, snapToken string) error {
	return r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("id = ?", id).
		Update("snap_token", snapToken).Error
}

type ticketRepository struct{ db *gorm.DB }

func NewTicketRepository(db *gorm.DB) repository.TicketRepository {
	return &ticketRepository{db: db}
}
func (r *ticketRepository) WithTx(tx *gorm.DB) repository.TicketRepository {
	return &ticketRepository{db: tx}
}
func (r *ticketRepository) CreateMany(ctx context.Context, tickets []entity.Ticket) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(&tickets).Error
}
func (r *ticketRepository) CancelByTransaction(ctx context.Context, transactionID uint) error {
	return r.db.WithContext(ctx).Model(&entity.Ticket{}).Where("transaction_id = ?", transactionID).
		Update("status", "cancelled").Error
}

type guestOrderRepository struct{ db *gorm.DB }

func NewGuestOrderRepository(db *gorm.DB) repository.GuestOrderRepository {
	return &guestOrderRepository{db: db}
}
func (r *guestOrderRepository) WithTx(tx *gorm.DB) repository.GuestOrderRepository {
	return &guestOrderRepository{db: tx}
}
func (r *guestOrderRepository) Create(ctx context.Context, guest *entity.GuestOrder) error {
	return r.db.WithContext(ctx).Create(guest).Error
}
func (r *guestOrderRepository) SetTransaction(ctx context.Context, guestOrderID uint, transactionID uint) error {
	return r.db.WithContext(ctx).Model(&entity.GuestOrder{}).Where("id = ?", guestOrderID).
		Update("transaction_id", transactionID).Error
}
