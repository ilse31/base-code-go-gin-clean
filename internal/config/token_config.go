package config

// TokenConfig holds configuration for the token service
type TokenConfig struct {
	AccessSecret  string
	RefreshSecret string
}

// NewTokenConfig creates a new TokenConfig from the main Config
func NewTokenConfig(cfg *Config) *TokenConfig {
	return &TokenConfig{
		AccessSecret:  cfg.Auth.AccessTokenSecret,
		RefreshSecret: cfg.Auth.RefreshTokenSecret,
	}
}
