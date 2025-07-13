package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
	svc "base-code-go-gin-clean/internal/service"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	mockUserRepository struct {
		mock.Mock
	}

	mockRedisRepository struct {
		mock.Mock
	}
)

// Implement redis.Repository interface for mockRedisRepository
func (m *mockRedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *mockRedisRepository) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockRedisRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *mockRedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	args := m.Called(ctx, key, expiration)
	return args.Bool(0), args.Error(1)
}

func (m *mockRedisRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) On(methodName string, arguments ...interface{}) *mock.Call {
	return m.Mock.On(methodName, arguments...)
}

func (m *mockUserRepository) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(mockUserRepository)
	mockRedis := new(mockRedisRepository)

	service := svc.NewUserService(svc.UserServiceConfig{
		UserRepo:  mockRepo,
		RedisRepo: mockRedis,
	})
	ctx := context.Background()

	t.Run("success - from cache", func(t *testing.T) {
		testID := uuid.New()
		cacheKey := "user:" + testID.String()
		expectedUser := &user.User{
			ID:        testID,
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedResponse := expectedUser.ToResponse()
		cachedData, _ := json.Marshal(expectedResponse)

		// Mock cache hit
		mockRedis.On("Get", ctx, cacheKey).Return(string(cachedData), nil)

		result, err := service.GetUserByID(ctx, testID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Name, result.Name)
		assert.Equal(t, expectedResponse.Email, result.Email)
		// Don't compare exact timestamps as they may vary slightly
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
		mockRedis.AssertExpectations(t)
		// Should not call the database when cache hit
		mockRepo.AssertNotCalled(t, "GetByID", ctx, testID)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		result, err := service.GetUserByID(ctx, "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid user ID format")
	})

	t.Run("success - from database", func(t *testing.T) {
		testID := uuid.New()
		cacheKey := "user:" + testID.String()
		expectedUser := &user.User{
			ID:        testID,
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedResponse := expectedUser.ToResponse()

		// Mock cache miss
		mockRedis.On("Get", ctx, cacheKey).Return("", redis.Nil) // redis.Nil is from go-redis/v9

		// Mock database call
		mockRepo.On("GetByID", ctx, testID).Return(expectedUser, nil)

		// ⚠️ Convert to string before mocking Set
		expectedCacheDataBytes, _ := json.Marshal(expectedResponse)
		expectedCacheData := string(expectedCacheDataBytes)

		// Mock cache set (string, not []byte)
		mockRedis.On("Set", ctx, cacheKey, expectedCacheData, mock.AnythingOfType("time.Duration")).Return(nil)

		// Call the service
		result, err := service.GetUserByID(ctx, testID.String())

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Name, result.Name)
		assert.Equal(t, expectedResponse.Email, result.Email)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		testID := uuid.New()
		cacheKey := "user:" + testID.String()

		mockRedis.On("Get", ctx, cacheKey).Return("", redis.Nil)
		mockRepo.On("GetByID", ctx, testID).Return(nil, errors.New("user not found"))

		result, err := service.GetUserByID(ctx, testID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})

}
