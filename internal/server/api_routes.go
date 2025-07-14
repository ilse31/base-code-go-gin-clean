package server

import (
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/handler/health"
	"base-code-go-gin-clean/internal/middleware"
	"base-code-go-gin-clean/internal/pkg/token"
	"base-code-go-gin-clean/internal/routes"
	"base-code-go-gin-clean/internal/service"

	"github.com/jmoiron/sqlx"
)

func (s *Server) setupRoutes(opts *ServerOptions) {
	var healthHandler *health.HealthHandler
	if s.db != nil {
		sqlxDB := sqlx.NewDb(s.db.DB, "postgres")
		healthHandler = health.NewHealthHandler(service.NewHealthService(sqlxDB))
	}

	// API v1 routes
	apiV1 := s.router.Group("/api/v1")
	{
		// Public routes
		public := apiV1.Group("")
		{
			public.GET("/health", healthHandler.HealthCheck)
			public.GET("/ping", handler.Ping)
		}

		// Protected routes (can be with or without auth)
		protected := apiV1.Group("")

		// If token config provided, apply auth middleware
		if opts.TokenConfig != nil {
			tokenService := token.NewTokenService(opts.TokenConfig)
			protected.Use(middleware.AuthMiddleware(tokenService))
		}

		// Setup user routes even if TokenConfig is nil
		if opts.UserHandler != nil {
			routes.SetupUserRoutes(protected, opts.UserHandler)
		}

		// Setup email routes
		if opts.EmailHandler != nil {
			routes.SetupEmailRoutes(apiV1, opts.EmailHandler)
		}

		// Setup auth routes (requires TokenConfig)
		if opts.AuthHandler != nil && opts.TokenConfig != nil {
			routes.SetupAuthRoutes(apiV1, opts.AuthHandler, opts.TokenConfig)
		}
	}
}
