//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"fmt"
	"log"
	"time"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/auth"
	"base-code-go-gin-clean/internal/pkg/redis"
	"base-code-go-gin-clean/internal/pkg/telemetry"
	"base-code-go-gin-clean/internal/pkg/token"
	"base-code-go-gin-clean/internal/repository/user"
	"base-code-go-gin-clean/internal/server"
	"base-code-go-gin-clean/internal/service"
	"base-code-go-gin-clean/pkg/logger"

	"github.com/google/wire"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ProvideBunDB(db *config.DB) *bun.DB {
	return db.DB
}

func ProvideEnv(cfg *config.Config) string {
	return cfg.Server.Environment
}

// ProvideServerOptions assembles all server options (handlers)
func ProvideServerOptions(
	userHandler *handler.UserHandler,
	authHandler *auth.AuthHandler,
	emailHandler *handler.EmailHandler,
	tokenConfig *config.TokenConfig,
	db *bun.DB,
	redisRepo redis.Repository,
) *server.ServerOptions {
	return &server.ServerOptions{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		EmailHandler: emailHandler,
		TokenConfig:  tokenConfig,
		DB:           db,
		RedisRepo:    redisRepo,
	}
}

// TokenConfigSet provides the token configuration
var TokenConfigSet = wire.NewSet(
	config.NewTokenConfig,
)

// TokenServiceSet provides the token service with its dependencies
var TokenServiceSet = wire.NewSet(
	token.NewTokenService,
	TokenConfigSet,
)

// AuthServiceSet is a Wire provider set that provides the auth service with its dependencies
var AuthServiceSet = wire.NewSet(
	service.NewAuthService,
	ProvideServiceConfig,
	TokenServiceSet,
	RedisSet,
)

// RedisSet is a Wire provider set that provides Redis client and repository
var RedisSet = wire.NewSet(
	ProvideRedisClient,
	ProvideRedisRepository,
)

// HandlerSet is a Wire provider set that provides all handlers
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewEmailHandler,
	auth.NewAuthHandler,
	ProvideEmailHandler,
)

// ServiceSet is a Wire provider set that provides all services
// TelemetrySet provides telemetry-related dependencies
var TelemetrySet = wire.NewSet(
	ProvideTracerProvider,
)

// ProvideTracerProvider creates a new TracerProvider
func ProvideTracerProvider(cfg *config.Config) (*trace.TracerProvider, func(), error) {
	cleanup, err := telemetry.InitTracer(cfg.Tracing.ServiceName, cfg.Tracing.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Get the global TracerProvider that was set by InitTracer
	tp := otel.GetTracerProvider()

	return tp.(*trace.TracerProvider), func() {
		// Ensure all spans are flushed before shutting down
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.(*trace.TracerProvider).Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		if cleanup != nil {
			cleanup()
		}
	}, nil
}

// ProvideServiceConfig provides the service configuration
func ProvideServiceConfig(cfg *config.Config) service.Config {
	return service.Config{
		Auth: service.AuthConfig{
			AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		},
	}
}

var ServiceSet = wire.NewSet(
	service.NewUserService,
	ProvideEmailService,
	ProvideUserServiceConfig,
	AuthServiceSet,
	ProvideServiceConfig,
)

// RepositorySet is a Wire provider set that provides all repositories
var RepositorySet = wire.NewSet(
	user.NewUserRepository,
	RedisSet,
)

// InitializeServer initializes the application server with all dependencies
func InitializeServer() (*server.Server, func(), error) {
	wire.Build(
		// Core providers
		ProvideConfig,
		ProvideEnv,
		ProvideDB,
		ProvideBunDB,
		TelemetrySet,
		logger.New,

		// Configuration
		config.NewTokenConfig,

		// Redis
		RedisSet,

		// Repositories
		user.NewUserRepository,

		// Services
		ProvideUserServiceConfig,
		service.NewUserService,
		ProvideTokenService,
		service.NewAuthService,
		ProvideEmailService,

		// Handlers
		handler.NewUserHandler,
		ProvideEmailHandler,
		auth.NewAuthHandler,

		// Server options
		wire.Struct(new(server.ServerOptions), "*"),

		// Server
		ProvideServiceConfig, // <--- Tambahkan baris ini
		server.New,
	)
	return nil, nil, nil // This will be replaced by Wire
}
