package service

import (
	"context"
	"errors"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
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
	userRepo      user.UserRepository
	tokenService  token.TokenService
	tokenStore    map[string]string // In-memory store for refresh tokens (use Redis in production)
}

func NewAuthService(userRepo user.UserRepository, tokenService token.TokenService) AuthService {
	return &authService{
		userRepo:      userRepo,
		tokenService:  tokenService,
		tokenStore:    make(map[string]string), // Initialize in-memory store
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
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
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

	// Store refresh token (in-memory for now, use Redis in production)
	s.tokenStore[user.ID.String()] = refreshToken

	// Convert user to response DTO
	userResponse := user.ToResponse()

	return &LoginResponse{
		User: userResponse,
		Token: &TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(15 * time.Minute / time.Second), // 15 minutes in seconds
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// In a real application, you would validate the refresh token against your token store
	// and get the associated user ID
	var userID string
	found := false
	for uid, rt := range s.tokenStore {
		if rt == refreshToken {
			userID = uid
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := s.tokenService.GenerateAccessToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Optionally, generate a new refresh token (refresh token rotation)
	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Update the refresh token in the store
	s.tokenStore[userID] = newRefreshToken

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken, // Include only if rotating refresh tokens
		ExpiresIn:    int64(15 * time.Minute / time.Second),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	// Remove the refresh token from the store
	delete(s.tokenStore, userID)
	return nil
}
