package service_test

import (
	"context"
	"testing"

	"base-code-go-gin-clean/internal/domain/user"
	svc "base-code-go-gin-clean/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
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
	service := svc.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		testID := uuid.New()
		expectedUser := &user.User{
			ID: testID,
			// Add other required fields
		}
		expectedResponse := expectedUser.ToResponse()

		mockRepo.On("GetByID", ctx, testID).Return(expectedUser, nil)

		result, err := service.GetUserByID(ctx, testID.String())

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		result, err := service.GetUserByID(ctx, "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid user ID format")
	})

	t.Run("user not found", func(t *testing.T) {
		testID := uuid.New()
		mockRepo.On("GetByID", ctx, testID).Return(nil, assert.AnError)

		result, err := service.GetUserByID(ctx, testID.String())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})
}
