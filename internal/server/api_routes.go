package server

import (
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/health"
	"base-code-go-gin-clean/internal/middleware"
	"base-code-go-gin-clean/internal/pkg/token"
	"base-code-go-gin-clean/internal/routes"
)

// setupRoutes configures all API routes
func (s *Server) setupRoutes(opts *ServerOptions) {

	healthHandler := health.NewHealthHandler()
	// API v1 routes
	apiV1 := s.router.Group("/api/v1")
	{
		// Public routes
		public := apiV1.Group("")
		{

			public.GET("/health", healthHandler.HealthCheck)
			public.GET("/ping", handler.Ping)
		}

		// Protected routes group (requires authentication)
		protected := apiV1.Group("")
		if opts.TokenConfig != nil {
			tokenService := token.NewTokenService(opts.TokenConfig)
			protected.Use(middleware.AuthMiddleware(tokenService))

			// Setup user routes under protected group if user handler is provided
			if opts.UserHandler != nil {
				routes.SetupUserRoutes(protected, opts.UserHandler)
			}
		}

		// Setup roles routes if roles handler is provided
		// if opts.RolesHandler != nil {
		// 	routes.SetupRolesRoutes(apiV1, opts.RolesHandler)
		// }

		// Setup email routes if email handler is provided
		if opts.EmailHandler != nil {
			routes.SetupEmailRoutes(apiV1, opts.EmailHandler)
		}

		// Setup auth routes if auth handler is provided
		if opts.AuthHandler != nil && opts.TokenConfig != nil {
			routes.SetupAuthRoutes(apiV1, opts.AuthHandler, opts.TokenConfig)
		}

		// Protected routes group
		// protected := apiV1.Group("")
		// protected.Use(middleware.AuthMiddleware()) /
		// {
		// 	// Add protected routes here
		// }
	}
}
