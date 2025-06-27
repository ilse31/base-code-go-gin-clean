package repository

import (
	"context"

	"base-code-go-gin-clean/internal/domain"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}
