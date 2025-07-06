package wire

import (
	"base-code-go-gin-clean/internal/config"
	emailDomain "base-code-go-gin-clean/internal/domain/email"
	emailService "base-code-go-gin-clean/internal/service/email"
	emailHandler "base-code-go-gin-clean/internal/handler/email"
	"base-code-go-gin-clean/internal/pkg/token"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideDB(cfg *config.Config) (*config.DB, error) {
	return config.NewDB(cfg)
}

func ProvideTokenService(cfg *config.Config) (token.TokenService, error) {
	tokenConfig := config.NewTokenConfig(cfg)
	tokenService := token.NewTokenService(tokenConfig)
	return tokenService, nil
}

// ProvideEmailService creates a new email service
func ProvideEmailService(cfg *config.Config) emailDomain.EmailService {
	return emailService.NewEmailService(cfg)
}

// ProvideEmailHandler creates a new email handler
func ProvideEmailHandler(emailSvc emailDomain.EmailService) *emailHandler.EmailHandler {
	return emailHandler.NewEmailHandler(emailSvc)
}
