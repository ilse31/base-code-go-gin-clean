package wire

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/domain/user"
	emailDomain "base-code-go-gin-clean/internal/domain/email"
	emailHandler "base-code-go-gin-clean/internal/handler/email"
	"base-code-go-gin-clean/internal/pkg/redis"
	"base-code-go-gin-clean/internal/service"
	emailService "base-code-go-gin-clean/internal/service/email"
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

// ProvideRedisClient creates a new Redis client
func ProvideRedisClient(cfg *config.Config) (*redis.Client, error) {
	return redis.NewClient(&cfg.Redis)
}

// ProvideRedisRepository creates a new Redis repository
func ProvideRedisRepository(client *redis.Client) redis.Repository {
	return redis.NewRepository(client)
}

// ProvideUserServiceConfig creates a new user service configuration
func ProvideUserServiceConfig(
	userRepo user.UserRepository,
	redisRepo redis.Repository,
) service.UserServiceConfig {
	return service.UserServiceConfig{
		UserRepo:  userRepo,
		RedisRepo: redisRepo,
		// Use default cache TTL
	}
}
