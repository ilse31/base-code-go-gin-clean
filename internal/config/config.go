package config

import "fmt"

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	DB      DatabaseConfig
	Tracing TracingConfig
	Email   EmailConfig
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

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        GetEnv("PORT", "8080"),
			Environment: GetEnv("ENVIRONMENT", "development"),
		},
		DB: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "120579"),
			Name:     GetEnv("DB_NAME", "postgres"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Tracing: TracingConfig{
			Enabled:     GetEnv("TRACING_ENABLED", "false") == "true",
			ServiceName: GetEnv("SERVICE_NAME", "base-code-go-gin-clean"),
			Version:     GetEnv("SERVICE_VERSION", "1.0.0"),
			DSN:         GetEnv("UPTRACE_DSN", ""),
		},
		Email: EmailConfig{
			SMTPServer:   GetEnv("SMTP_SERVER", "smtp.gmail.com"),
			SMTPPort:     GetEnv("SMTP_PORT", "587"),
			SMTPUsername: GetEnv("SMTP_USERNAME", ""),
			SMTPPassword: GetEnv("SMTP_PASSWORD", ""),
			From:         GetEnv("EMAIL_FROM", "noreply@example.com"),
		},
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
