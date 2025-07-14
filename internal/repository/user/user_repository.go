package user

import (
	"context"

	"base-code-go-gin-clean/internal/domain/user"
	"base-code-go-gin-clean/internal/pkg/telemetry"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) user.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	user := new(user.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	user := new(user.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)

	if err != nil {
		span.RecordError(err)
	}

	return err
}
