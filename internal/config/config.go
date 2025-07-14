package config

import "fmt"

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	DB      DatabaseConfig
	Tracing TracingConfig
	Email   EmailConfig
	Auth    AuthConfig
	Redis   RedisConfig
}

// AuthConfig holds authentication related configuration
type AuthConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  int // in minutes
	RefreshTokenExpiry int // in hours
	Issuer             string
}

type ServerConfig struct {
	Port        string
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Version     string
	DSN         string
}

type EmailConfig struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	From         string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        GetEnv("PORT", "8080"),
			Environment: GetEnv("ENVIRONMENT", "development"),
		},
		DB: DatabaseConfig{
			Host:     GetEnv("DB_HOST", ""),
			Port:     GetEnv("DB_PORT", ""),
			User:     GetEnv("DB_USER", ""),
			Password: GetEnv("DB_PASSWORD", ""),
			Name:     GetEnv("DB_NAME", ""),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     GetEnv("REDIS_HOST", ""),
			Port:     GetEnv("REDIS_PORT", ""),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		Auth: AuthConfig{
			AccessTokenSecret:  GetEnv("ACCESS_TOKEN_SECRET", ""),
			RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", ""),
			AccessTokenExpiry:  GetEnvAsInt("ACCESS_TOKEN_EXPIRY_MINUTES", 15),
			RefreshTokenExpiry: GetEnvAsInt("REFRESH_TOKEN_EXPIRY_HOURS", 24),
			Issuer:             GetEnv("JWT_ISSUER", "base-code-go-gin-clean"),
		},
		Tracing: TracingConfig{
			Enabled:     GetEnv("TRACING_ENABLED", "false") == "true",
			ServiceName: GetEnv("SERVICE_NAME", "base-code-go-gin-clean"),
			Version:     GetEnv("SERVICE_VERSION", "1.0.0"),
			DSN:         GetEnv("UPTRACE_DSN", ""),
		},
		Email: EmailConfig{
			SMTPServer:   GetEnv("SMTP_SERVER", ""),
			SMTPPort:     GetEnv("SMTP_PORT", ""),
			SMTPUsername: GetEnv("SMTP_USERNAME", ""),
			SMTPPassword: GetEnv("SMTP_PASSWORD", ""),
			From:         GetEnv("EMAIL_FROM", ""),
		},
	}

	if cfg.DB.Host == "" || cfg.DB.Port == "" || cfg.DB.User == "" || cfg.DB.Password == "" || cfg.DB.Name == "" {
		return nil, fmt.Errorf("database credentials are not set")
	}

	if cfg.Redis.Host == "" || cfg.Redis.Port == "" {
		return nil, fmt.Errorf("redis credentials are not set")
	}

	if cfg.Auth.AccessTokenSecret == "" || cfg.Auth.RefreshTokenSecret == "" {
		return nil, fmt.Errorf("JWT secrets are not set")
	}

	if cfg.Email.SMTPServer == "" || cfg.Email.SMTPPort == "" || cfg.Email.SMTPUsername == "" || cfg.Email.SMTPPassword == "" || cfg.Email.From == "" {
		return nil, fmt.Errorf("email credentials are not set")
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}
