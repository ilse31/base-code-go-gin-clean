package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"base-code-go-gin-clean/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

// Server represents the HTTP server
type Server struct {
	router        *gin.Engine
	config        *config.Config
	logger        *slog.Logger
	server        *http.Server
	tracerCleanup func()  // Function to clean up tracer resources
	db            *bun.DB // Database connection
}

// New creates a new Server instance
func New(cfg *config.Config, log *slog.Logger, opts *ServerOptions) *Server {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize options
	options := opts
	if options == nil {
		options = &ServerOptions{}
	}

	// Initialize Gin router
	router := gin.New()

	// Add Swagger route
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	srv := &Server{
		router: router,
		config: cfg,
		logger: log,
		db:     options.DB, // Store the database connection
	}

	// Setup server
	srv.setupMiddlewares(options.DB)
	srv.setupRoutes(options)

	return srv
}

// Start starts the HTTP server
// Start runs the HTTP server and supports graceful shutdown via context
func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Run server in goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		s.logger.Info("Server starting", "port", s.config.Server.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
		close(serverErrChan)
	}()

	// Shutdown when context is done
	<-ctx.Done()
	s.logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tracer cleanup if applicable
	if s.tracerCleanup != nil {
		s.logger.Info("Cleaning up tracer...")
		s.tracerCleanup()
	}

	err := s.server.Shutdown(shutdownCtx)
	if err != nil {
		s.logger.Error("Forced to shutdown", "error", err)
		return err
	}

	s.logger.Info("Server exited gracefully")
	return <-serverErrChan // if server exited with error
}

// GetRouter returns the underlying gin.Engine instance
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetHTTPServer returns the underlying http.Server instance
func (s *Server) GetHTTPServer() *http.Server {
	return s.server
}
