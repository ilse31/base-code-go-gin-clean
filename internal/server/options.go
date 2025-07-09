package server

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/auth"
	emailHandler "base-code-go-gin-clean/internal/handler/email"
	"base-code-go-gin-clean/internal/pkg/redis"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/sdk/trace"
)

// ServerOptions contains options for creating a new Server

type ServerOptions struct {
	UserHandler  *handler.UserHandler
	AuthHandler  *auth.AuthHandler
	EmailHandler *emailHandler.EmailHandler
	TokenConfig  *config.TokenConfig
	DB           *bun.DB // Add database connection to options
	RedisRepo    redis.Repository // Add Redis repository to options
	TracerProvider *trace.TracerProvider // Add TracerProvider for distributed tracing
}

// Option configures how we set up the server
type Option func(*ServerOptions)

// WithUserHandler is an option to set the user handler
func WithUserHandler(h *handler.UserHandler) Option {
	return func(opts *ServerOptions) {
		opts.UserHandler = h
	}
}

// WithRolesHandler is an option to set the roles handler
// func WithRolesHandler(h *handler.RolesHandler) Option {
// 	return func(opts *ServerOptions) {
// 		opts.RolesHandler = h
// 	}
// }

// WithEmailHandler is an option to set the email handler
func WithEmailHandler(h *emailHandler.EmailHandler) Option {
	return func(opts *ServerOptions) {
		opts.EmailHandler = h
	}
}

// WithAuthHandler is an option to set the auth handler
func WithAuthHandler(h *auth.AuthHandler) Option {
	return func(opts *ServerOptions) {
		opts.AuthHandler = h
	}
}

// WithTokenConfig is an option to set the token configuration
func WithTokenConfig(cfg *config.TokenConfig) Option {
	return func(opts *ServerOptions) {
		opts.TokenConfig = cfg
	}
}

// WithRedisRepo is an option to set the Redis repository
func WithRedisRepo(repo redis.Repository) Option {
	return func(opts *ServerOptions) {
		opts.RedisRepo = repo
	}
}

// WithTracerProvider is an option to set the TracerProvider for distributed tracing
func WithTracerProvider(tp *trace.TracerProvider) Option {
	return func(opts *ServerOptions) {
		opts.TracerProvider = tp
	}
}
