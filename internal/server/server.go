package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"base-code-go-gin-clean/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	router        *gin.Engine
	config        *config.Config
	logger        *slog.Logger
	server        *http.Server
	tracerCleanup func() // Function to clean up tracer resources
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
	}

	// Setup server
	srv.setupMiddleware()

	// Setup database middleware if DB is provided
	if options.DB != nil {
		srv.setupDatabaseMiddleware(options.DB)
	}

	srv.setupRoutes(options)

	return srv
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server in goroutine
	go func() {
		s.logger.Info("Server starting", "port", s.config.Server.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error", "error", err)
		}
	}()

	return s.handleShutdown()
}

// handleShutdown handles graceful shutdown of the server
func (s *Server) handleShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal
	<-quit

	s.logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call tracer cleanup if it was initialized
	if s.tracerCleanup != nil {
		s.logger.Info("Cleaning up tracer...")
		s.tracerCleanup()
	}

	// Shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	s.logger.Info("Server exited gracefully")
	return nil
}

// GetRouter returns the underlying gin.Engine instance
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetHTTPServer returns the underlying http.Server instance
func (s *Server) GetHTTPServer() *http.Server {
	return s.server
}
