package config

import "time"

// TokenConfig holds configuration for the token service
type TokenConfig struct {
	AccessSecret       string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// NewTokenConfig creates a new TokenConfig from the main Config
func NewTokenConfig(cfg *Config) *TokenConfig {
	return &TokenConfig{
		AccessSecret:       cfg.Auth.AccessTokenSecret,
		RefreshSecret:      cfg.Auth.RefreshTokenSecret,
		AccessTokenExpiry:  time.Duration(cfg.Auth.AccessTokenExpiry),
		RefreshTokenExpiry: time.Duration(cfg.Auth.RefreshTokenExpiry),
	}
}
