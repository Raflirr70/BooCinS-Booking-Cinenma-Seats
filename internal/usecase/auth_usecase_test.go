package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rafli/boocins/config"
	"github.com/rafli/boocins/internal/domain/entity"
	"github.com/rafli/boocins/internal/usecase"
	"github.com/rafli/boocins/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func generateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func getTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 1 * time.Hour,
	}
}

func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	user := &entity.User{
		FirstName: "Rafli",
		LastName:  "Setiawan",
		Email:     "rafli@test.com",
		Password:  "password123",
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "rafli@test.com").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindRoleByName", mock.Anything, "member").Return(&entity.Role{ID: 5, Name: "member"}, nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	err := uc.Register(context.Background(), user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, uint(5), user.RoleID)
	mockUserRepo.AssertExpectations(t)
}

func TestRegister_EmailExists(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	existing := &entity.User{ID: 1, Email: "rafli@test.com"}
	mockUserRepo.On("FindByEmail", mock.Anything, "rafli@test.com").Return(existing, nil)

	user := &entity.User{
		FirstName: "Rafli",
		LastName:  "Setiawan",
		Email:     "rafli@test.com",
		Password:  "password123",
	}

	err := uc.Register(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	user := &entity.User{
		ID:        1,
		Email:     "rafli@test.com",
		Password:  "$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ12", // bcrypt hash of "password123"
		FirstName: "Rafli",
		LastName:  "Setiawan",
		Role:      entity.Role{ID: 5, Name: "member"},
	}

	// Generate a real bcrypt hash for test
	hashed, _ := generateHash("password123")
	user.Password = string(hashed)

	mockUserRepo.On("FindByEmail", mock.Anything, "rafli@test.com").Return(user, nil)

	returnedUser, token, err := uc.Login(context.Background(), "rafli@test.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, returnedUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, uint(1), returnedUser.ID)
	mockUserRepo.AssertExpectations(t)

	// Verify token is stored in Redis
	redisKey := "session:1"
	val, err := redisClient.Get(context.Background(), redisKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, token, val)
}

func TestLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	mockUserRepo.On("FindByEmail", mock.Anything, "unknown@test.com").Return(nil, errors.New("record not found"))

	user, token, err := uc.Login(context.Background(), "unknown@test.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	hashed, _ := generateHash("password123")
	user := &entity.User{
		ID:       1,
		Email:    "rafli@test.com",
		Password: string(hashed),
		Role:     entity.Role{ID: 5, Name: "member"},
	}

	mockUserRepo.On("FindByEmail", mock.Anything, "rafli@test.com").Return(user, nil)

	returnedUser, token, err := uc.Login(context.Background(), "rafli@test.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, returnedUser)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	// First login to get a valid token
	hashed, _ := generateHash("password123")
	user := &entity.User{
		ID:       1,
		Email:    "rafli@test.com",
		Password: string(hashed),
		Role:     entity.Role{ID: 5, Name: "member"},
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "rafli@test.com").Return(user, nil)

	_, token, err := uc.Login(context.Background(), "rafli@test.com", "password123")
	assert.NoError(t, err)

	// Then logout
	err = uc.Logout(context.Background(), token)
	assert.NoError(t, err)

	// Verify session removed from Redis
	redisKey := "session:1"
	_, err = redisClient.Get(context.Background(), redisKey).Result()
	assert.Error(t, err) // should not exist
}

func TestGetProfile_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	redisClient := setupTestRedis(t)
	jwtCfg := getTestJWTConfig()

	uc := usecase.NewAuthUsecase(mockUserRepo, redisClient, jwtCfg)

	expected := &entity.User{
		ID:        1,
		Email:     "rafli@test.com",
		FirstName: "Rafli",
		LastName:  "Setiawan",
	}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(expected, nil)

	user, err := uc.GetProfile(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Rafli", user.FirstName)
	mockUserRepo.AssertExpectations(t)
}
