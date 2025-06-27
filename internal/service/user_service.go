package service

import (
	"context"
	"errors"

	"base-code-go-gin-clean/internal/domain"
	"base-code-go-gin-clean/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, idStr string) (*domain.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
