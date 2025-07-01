package repository_test

import (
	"context"
	"testing"

	"base-code-go-gin-clean/internal/domain/user"
	repo "base-code-go-gin-clean/internal/repository"
	"base-code-go-gin-clean/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByID(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.TeardownTestDB(t, db)

	repo := repo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		// Create a test user
		testUser := &user.User{
			ID:       uuid.New(),
			Email:    "test@example.com",
			Name:     "Test User",
			Password: "password",
		}
		_, err := db.NewInsert().Model(testUser).Exec(context.Background())
		assert.NoError(t, err)

		// Test GetByID
		result, err := repo.GetByID(context.Background(), testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testUser.ID, result.ID)
	})

	t.Run("not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		result, err := repo.GetByID(context.Background(), nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
