package server

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	domainhttplog "base-code-go-gin-clean/internal/domain/httplog"
	"base-code-go-gin-clean/internal/pkg/httplog"
	"base-code-go-gin-clean/internal/pkg/telemetry"
	"base-code-go-gin-clean/pkg/dbutils"
	"base-code-go-gin-clean/pkg/middleware"
)

// setupMiddleware configures global middleware for the server
func (s *Server) setupMiddleware() {
	// Add Jaeger tracing middleware if configured
	if s.config.Tracing.Enabled && s.config.Tracing.DSN != "" {
		cleanup, err := telemetry.InitTracer(
			s.config.Tracing.ServiceName,
			s.config.Tracing.DSN,
		)
		if err != nil {
			s.logger.Error("failed to initialize tracer", "error", err)
		} else {
			// Store the cleanup function to be called on server shutdown
			s.tracerCleanup = cleanup
			// Add tracing middleware
			s.router.Use(telemetry.Middleware(s.config.Tracing.ServiceName))
		}
	}

	// Request ID middleware
	s.router.Use(middleware.RequestID())

	// Logger middleware
	s.router.Use(gin.Logger())

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(middleware.CORS())

	// Security headers middleware
	s.router.Use(middleware.Secure())

	// Rate limiting
	s.router.Use(middleware.RateLimitMiddleware())

	// Only set up HTTP log middleware if we have a database connection
	if s.db != nil {
		httpLogRepo := domainhttplog.NewRepository(s.db)
		httpLogService := domainhttplog.NewService(httpLogRepo)

		s.router.Use(httplog.Middleware(httplog.Config{
			Service:             httpLogService,
			SkipPaths:           []string{"/health", "/metrics"},
			SkipHeaders:         []string{"Authorization", "Cookie"},
			SkipBodyMethods:     map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true},
			MaxBodySize:         1024 * 1024,
			IncludeResponseBody: true,
		}))
	} else {
		s.logger.Warn("No database connection available, HTTP logging will be disabled")
	}
}

// setupDatabaseMiddleware adds database-related middleware to the router
func (s *Server) setupDatabaseMiddleware(db *bun.DB) {
	// Add database connection to the Gin context
	s.router.Use(func(c *gin.Context) {
		c.Set("db_conn", db)
		c.Next()
	})

	// Add transaction middleware
	s.router.Use(dbutils.TransactionMiddleware(db))
}
