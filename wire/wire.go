//go:build wireinject
// +build wireinject

package wire

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/auth"
	"base-code-go-gin-clean/internal/pkg/token"
	"base-code-go-gin-clean/internal/repository"
	"base-code-go-gin-clean/internal/server"
	"base-code-go-gin-clean/internal/service"
	"base-code-go-gin-clean/pkg/logger"

	"github.com/google/wire"
	"github.com/uptrace/bun"
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
) *server.ServerOptions {
	return &server.ServerOptions{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		EmailHandler: emailHandler,
		TokenConfig:  tokenConfig,
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
	TokenServiceSet,
)

// HandlerSet is a Wire provider set that provides all handlers
var HandlerSet = wire.NewSet(
	handler.NewUserHandler,
	handler.NewEmailHandler,
	auth.NewAuthHandler,
	ProvideEmailHandler,
)

// ServiceSet is a Wire provider set that provides all services
var ServiceSet = wire.NewSet(
	service.NewUserService,
	ProvideEmailService,
	AuthServiceSet,
)

// RepositorySet is a Wire provider set that provides all repositories
var RepositorySet = wire.NewSet(
	repository.NewUserRepository,
)

// InitializeServer initializes the application server with all dependencies
func InitializeServer() (*server.Server, func(), error) {
	wire.Build(
		// Core providers
		ProvideConfig,
		ProvideEnv,
		ProvideDB,
		ProvideBunDB,
		logger.New,

		// Configuration
		config.NewTokenConfig,

		// Repositories
		repository.NewUserRepository,

		// Services
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
		server.New,
	)
	return nil, nil, nil // This will be replaced by Wire
}
