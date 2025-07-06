package service_test

import (
	"context"
	"testing"

	"base-code-go-gin-clean/internal/domain/user"
	"base-code-go-gin-clean/internal/pkg/token"
	svc "base-code-go-gin-clean/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockTokenService struct {
	token.TokenService
}

func (m *mockTokenService) GenerateAccessToken(userID string) (string, error) {
	return "test-access-token", nil
}

func (m *mockTokenService) GenerateRefreshToken() (string, error) {
	return "test-refresh-token", nil
}

func (m *mockTokenService) ValidateAccessToken(tokenString string) (string, error) {
	return "test-user-id", nil
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(mockUserRepository)
	mockTokenSvc := &mockTokenService{}
	service := svc.NewAuthService(mockRepo, mockTokenSvc)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "Test User"
		email := "test@example.com"
		password := "password123"

		mockRepo.On("GetByEmail", ctx, email).Return((*user.User)(nil), nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)

		userResp, err := service.Register(ctx, name, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.Equal(t, name, userResp.Name)
		assert.Equal(t, email, userResp.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		name := "Test User"
		email := "existing@example.com"
		password := "password123"

		existingUser := &user.User{
			ID:    uuid.New(),
			Name:  "Existing User",
			Email: email,
		}

		mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil)

		userResp, err := service.Register(ctx, name, email, password)

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Equal(t, "user with this email already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(mockUserRepository)
	mockTokenSvc := &mockTokenService{}
	service := svc.NewAuthService(mockRepo, mockTokenSvc)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &user.User{
			ID:       uuid.New(),
			Name:     "Test User",
			Email:    email,
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

		userResp, err := service.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.Equal(t, user.ID, userResp.User.ID)
		assert.Equal(t, user.Name, userResp.User.Name)
		assert.Equal(t, user.Email, userResp.User.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"

		user := &user.User{
			ID:       uuid.New(),
			Name:     "Test User",
			Email:    email,
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // hash for "password"
		}

		mockRepo.On("GetByEmail", ctx, email).Return(user, nil)

		userResp, err := service.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		email := "nonexistent@example.com"
		password := "password123"

		mockRepo.On("GetByEmail", ctx, email).Return((*user.User)(nil), assert.AnError)

		userResp, err := service.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Equal(t, "invalid email or password", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
