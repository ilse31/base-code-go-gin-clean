package server

import "base-code-go-gin-clean/internal/handler"

// ServerOptions contains options for creating a new Server
type ServerOptions struct {
	UserHandler  *handler.UserHandler
	RolesHandler *handler.RolesHandler
	EmailHandler *handler.EmailHandler
	// Add other handlers as needed
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
func WithEmailHandler(h *handler.EmailHandler) Option {
	return func(opts *ServerOptions) {
		opts.EmailHandler = h
	}
}
