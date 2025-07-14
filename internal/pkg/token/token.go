package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"base-code-go-gin-clean/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(tokenString string) (string, error)
}

type tokenService struct {
	accessSecret       string
	refreshSecret      string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewTokenService creates a new TokenService with the given configuration
func NewTokenService(config *config.TokenConfig) TokenService {
	return &tokenService{
		accessSecret:       config.AccessSecret,
		refreshSecret:      config.RefreshSecret,
		accessTokenExpiry:  config.AccessTokenExpiry,
		refreshTokenExpiry: config.RefreshTokenExpiry,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *tokenService) GenerateAccessToken(userID string) (string, error) {
	expirationTime := time.Now().Add(s.accessTokenExpiry)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

func (s *tokenService) GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (s *tokenService) ValidateAccessToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.New("token has expired")
		}
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}
