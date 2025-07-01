//go:build wireinject
// +build wireinject

package wire

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
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
	emailHandler *handler.EmailHandler,
) []server.Option {
	return []server.Option{
		server.WithUserHandler(userHandler),
		server.WithEmailHandler(emailHandler),
	}
}

func InitializeServer() (*server.Server, error) {
	wire.Build(
		// Core
		ProvideConfig,
		ProvideEnv,
		ProvideDB,
		ProvideBunDB,
		logger.New,

		// Repository
		repository.NewUserRepository,

		// Service
		service.NewUserService,
		service.NewEmailService,

		// Handler
		handler.NewUserHandler,
		handler.NewEmailHandler,

		// Server options + server
		ProvideServerOptions,
		server.New,
	)
	return nil, nil
}
