package server

import (
	"context"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	domainhttplog "base-code-go-gin-clean/internal/domain/httplog"
	pkghttplog "base-code-go-gin-clean/internal/pkg/httplog"
	"base-code-go-gin-clean/internal/pkg/telemetry"
	"base-code-go-gin-clean/pkg/dbutils"
	"base-code-go-gin-clean/pkg/middleware"
)

// setupMiddlewares configures all middleware for the server
func (s *Server) setupMiddlewares(db *bun.DB) {
	// Add request size limiting middleware (10MB max)
	s.router.MaxMultipartMemory = 10 << 20 // 10MB

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

	// Compression middleware
	s.router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(middleware.CORS())

	// Security headers middleware
	s.router.Use(middleware.Secure())

	// Rate limiting
	s.router.Use(middleware.RateLimitMiddleware())

	// Timeout middleware (10 seconds)
	s.router.Use(TimeoutMiddleware(10 * time.Second))

	// Only set up HTTP log middleware if we have a database connection
	if s.db != nil {
		httpLogRepo := domainhttplog.NewRepository(s.db)
		httpLogService := domainhttplog.NewService(httpLogRepo)

		s.router.Use(pkghttplog.Middleware(pkghttplog.Config{
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

	// Add database connection to the Gin context
	s.router.Use(func(c *gin.Context) {
		c.Set("db_conn", db)
		c.Next()
	})

	// Add transaction middleware
	s.router.Use(dbutils.TransactionMiddleware(db))
}

// TimeoutMiddleware creates a middleware that times out requests after specified duration
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request with timeout context
		c.Request = c.Request.WithContext(ctx)

		// Continue processing
		c.Next()
	}
}
