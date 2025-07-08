package server

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/auth"
	emailHandler "base-code-go-gin-clean/internal/handler/email"
	"github.com/uptrace/bun"
)

// ServerOptions contains options for creating a new Server

type ServerOptions struct {
	UserHandler  *handler.UserHandler
	AuthHandler  *auth.AuthHandler
	EmailHandler *emailHandler.EmailHandler
	TokenConfig  *config.TokenConfig
	DB           *bun.DB // Add database connection to options
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
