package service

import (
	"context"
	"errors"

	"base-code-go-gin-clean/internal/domain/user"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*user.UserResponse, error)
}

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, idStr string) (*user.UserResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}
