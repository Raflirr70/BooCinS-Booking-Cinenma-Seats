package usecase

import (
	"context"

	"github.com/rafli/boocins/internal/domain/entity"
)

type AuthUsecase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email, password string) (*entity.User, string, error)
	Logout(ctx context.Context, token string) error
	GetProfile(ctx context.Context, userID uint) (*entity.User, error)
}
