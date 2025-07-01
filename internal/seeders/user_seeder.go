package seeders

import (
	"context"
	"time"

	"base-code-go-gin-clean/internal/domain/user"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	db *bun.DB
}

// NewUserSeeder creates a new UserSeeder
func NewUserSeeder(db *bun.DB) *UserSeeder {
	return &UserSeeder{
		db: db,
	}
}

// Seed creates sample users in the database
func (s *UserSeeder) Seed(ctx context.Context) error {
	// Check if users already exist
	count, err := s.db.NewSelect().Model((*user.User)(nil)).Count(ctx)
	if err != nil {
		return err
	}

	// If users already exist, skip seeding
	if count > 0 {
		return nil
	}

	// Sample users data
	users := []*user.User{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name:      "Admin User",
			Email:     "admin@example.com",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert users in a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err = tx.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// SeedAll runs all seeders
func SeedAll(db *bun.DB) error {
	ctx := context.Background()

	// Initialize seeders
	userSeeder := NewUserSeeder(db)

	// Run seeders
	if err := userSeeder.Seed(ctx); err != nil {
		return err
	}

	return nil
}
