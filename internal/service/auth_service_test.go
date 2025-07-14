package service_test

import (
	"context"
	"testing"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
	"base-code-go-gin-clean/internal/service"
	"base-code-go-gin-clean/test/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockTokenSvc := &mocks.MockTokenService{}
	redisRepo := &mocks.MockRedisRepository{}
	cfg := service.Config{
		Auth: service.AuthConfig{
			AccessTokenExpiry: 15,
		},
	}
	service := service.NewAuthService(mockRepo, mockTokenSvc, redisRepo, cfg)
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

	t.Run("user already exists", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, "exist@test.com").Return(&user.User{}, nil)
		_, err := service.Register(ctx, "Test", "exist@test.com", "password")
		assert.Error(t, err)
		assert.Equal(t, "user with this email already exists", err.Error())
	})

	t.Run("password hashing error", func(t *testing.T) {
		email := "hash123@test.com"

		// Simulasikan bahwa user belum ada
		mockRepo.On("GetByEmail", ctx, email).Return((*user.User)(nil), nil)

		// Gunakan password yang sangat panjang (atau cara buatan untuk memicu error hashing)
		invalidPassword := string(make([]byte, 1<<20)) // 1MB string bisa memicu error bcrypt

		_, err := service.Register(ctx, "Test", email, invalidPassword)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash password")
	})

}

func TestAuthService_Login(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenService := &mocks.MockTokenService{}
	redisRepo := &mocks.MockRedisRepository{}
	cfg := service.Config{
		Auth: service.AuthConfig{
			AccessTokenExpiry: 15,
		},
	}
	service := service.NewAuthService(userRepo, tokenService, redisRepo, cfg)
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

		userRepo.On("GetByEmail", ctx, email).Return(user, nil)

		tokenService.On("GenerateAccessToken", mock.AnythingOfType("string")).Return("access-token-123", nil)
		tokenService.On("GenerateRefreshToken", mock.AnythingOfType("string")).Return("refresh-token-123", nil)

		redisRepo.On(
			"Set",
			mock.Anything, // ctx
			"refresh_token:"+user.ID.String(),
			"refresh-token-123",
			15*time.Minute,
		).Return(nil)

		userResp, err := service.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.Equal(t, user.ID, userResp.User.ID)
		assert.Equal(t, "access-token-123", userResp.Token.AccessToken)
		assert.Equal(t, "refresh-token-123", userResp.Token.RefreshToken)

		userRepo.AssertExpectations(t)
		tokenService.AssertExpectations(t)
		redisRepo.AssertExpectations(t)
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

		userRepo.On("GetByEmail", ctx, email).Return(user, nil)

		userResp, err := service.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Equal(t, "invalid email or password", err.Error())

		userRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		email := "nonexistent@example.com"
		password := "password123"

		userRepo.On("GetByEmail", ctx, email).Return((*user.User)(nil), assert.AnError)

		userResp, err := service.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Nil(t, userResp)
		assert.Equal(t, "invalid email or password", err.Error())

		userRepo.AssertExpectations(t)
	})

	t.Run("password mismatch", func(t *testing.T) {
		testUser := &user.User{Password: "$2a$10$hashed"} // Hashed password doesn't match "wrong"
		userRepo.On("GetByEmail", ctx, "mismatch@test.com").Return(testUser, nil)
		_, err := service.Login(ctx, "mismatch@test.com", "wrong")
		assert.Error(t, err)
		assert.Equal(t, "invalid email or password", err.Error())
	})

}

func TestAuthService_RefreshToken(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenService := &mocks.MockTokenService{}
	redisRepo := &mocks.MockRedisRepository{}
	cfg := service.Config{
		Auth: service.AuthConfig{
			AccessTokenExpiry: 15,
		},
	}
	service := service.NewAuthService(userRepo, tokenService, redisRepo, cfg)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Setup Redis Get expectation
		redisRepo.On("Get", mock.Anything, "valid-refresh-token").Return("user123", nil)

		// Setup TokenService expectations
		tokenService.On("GenerateAccessToken", "user123").Return("new-access-token", nil)
		tokenService.On("GenerateRefreshToken").Return("new-refresh-token", nil)

		// Setup Redis Set expectation for new token
		redisRepo.On("Set", mock.Anything, "refresh_token:user123",
			"new-refresh-token", 7*24*time.Hour).Return(nil)

		// Test and assertions
		tokenResp, err := service.RefreshToken(ctx, "valid-refresh-token")
		assert.NoError(t, err)
		assert.NotNil(t, tokenResp)
		assert.Equal(t, "new-access-token", tokenResp.AccessToken)
		assert.Equal(t, "new-refresh-token", tokenResp.RefreshToken)

		tokenService.AssertExpectations(t)
		redisRepo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	tokenService := &mocks.MockTokenService{}
	redisRepo := &mocks.MockRedisRepository{}
	cfg := service.Config{
		Auth: service.AuthConfig{
			AccessTokenExpiry: 15,
		},
	}
	service := service.NewAuthService(userRepo, tokenService, redisRepo, cfg)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Expect Redis Delete call
		redisRepo.On("Delete", mock.Anything, "refresh_token:user123").Return(nil)

		err := service.Logout(ctx, "user123")
		assert.NoError(t, err)
		redisRepo.AssertExpectations(t)
	})
}
