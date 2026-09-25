package repository

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
}
