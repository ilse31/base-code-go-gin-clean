package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
	"base-code-go-gin-clean/internal/pkg/redis"
	"base-code-go-gin-clean/internal/pkg/token"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*user.UserResponse, error)
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, userID string) error
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginResponse struct {
	User  *user.UserResponse `json:"user"`
	Token *TokenResponse     `json:"token"`
}

type authService struct {
	userRepo     user.UserRepository
	tokenService token.TokenService
	redisRepo    redis.Repository
	cfg          Config
}

type Config struct {
	Auth AuthConfig
}

type AuthConfig struct {
	AccessTokenExpiry int
}

func NewAuthService(userRepo user.UserRepository, tokenService token.TokenService, redisRepo redis.Repository, cfg Config) AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
		redisRepo:    redisRepo,
		cfg:          cfg,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*user.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	newUser := &user.User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}

	// Hash password
	if err := newUser.HashPassword(password); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser.ToResponse(), nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Store refresh token in Redis
	if err := s.redisRepo.Set(ctx, "refresh_token:"+user.ID.String(), refreshToken, 15*time.Minute); err != nil {
		return nil, errors.New("failed to store refresh token")
	}

	// Convert user to response DTO
	userResponse := user.ToResponse()

	return &LoginResponse{
		User: userResponse,
		Token: &TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(s.cfg.Auth.AccessTokenExpiry * 60), // Convert minutes to seconds
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Get user ID from Redis
	userID, err := s.redisRepo.Get(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := s.tokenService.GenerateAccessToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate new refresh token
	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Update the refresh token in Redis
	if err := s.redisRepo.Set(ctx, "refresh_token:"+userID, newRefreshToken, 7*24*time.Hour); err != nil {
		return nil, errors.New("failed to update refresh token")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.cfg.Auth.AccessTokenExpiry * 60),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// Remove the refresh token from Redis
	err := s.redisRepo.Delete(ctx, "refresh_token:"+userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}
