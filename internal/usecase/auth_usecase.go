package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rafli/boocins/config"
	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/domain/repository"
	uc "github.com/rafli/boocins/internal/domain/usecase"
	jwtPkg "github.com/rafli/boocins/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo repository.UserRepository
	redis    *redis.Client
	jwtCfg   config.JWTConfig
}

func NewAuthUsecase(userRepo repository.UserRepository, redisClient *redis.Client, jwtCfg config.JWTConfig) uc.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		redis:    redisClient,
		jwtCfg:   jwtCfg,
	}
}

func (u *authUsecase) Register(ctx context.Context, user *entity.User) error {
	existing, _ := u.userRepo.FindByEmail(ctx, user.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	role, err := u.userRepo.FindRoleByName(ctx, "member")
	if err != nil {
		return errors.New("default role not found")
	}
	user.RoleID = role.ID

	return u.userRepo.Create(ctx, user)
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := jwtPkg.GenerateToken(user.ID, user.Email, user.Role.Name, u.jwtCfg.Secret, u.jwtCfg.Expiration)
	if err != nil {
		return nil, "", err
	}

	redisKey := fmt.Sprintf("session:%d", user.ID)
	u.redis.Set(ctx, redisKey, token, u.jwtCfg.Expiration)

	return user, token, nil
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	claims, err := jwtPkg.ParseToken(token, u.jwtCfg.Secret)
	if err != nil {
		return errors.New("invalid token")
	}

	redisKey := fmt.Sprintf("session:%d", claims.UserID)
	u.redis.Del(ctx, redisKey)

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		blacklistKey := fmt.Sprintf("blacklist:%s", token)
		u.redis.Set(ctx, blacklistKey, "1", ttl)
	}

	return nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID uint) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
