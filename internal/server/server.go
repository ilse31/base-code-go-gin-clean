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
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
	logger *slog.Logger
	// Add other dependencies
}

func New(cfg *config.Config, log *slog.Logger) *Server {
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &Server{
		router: gin.New(),
		config: cfg,
		logger: log,
	}

	srv.SetupRoutes()
	return srv
}

func (s *Server) SetupRoutes() {
	s.router.Use(
		middleware.RequestID(),
		gin.Logger(),
		gin.Recovery(),
	)

	// Health check
	s.router.GET("/health", handler.HealthCheck)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", handler.Ping)
	}
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server in goroutine
	go func() {
		s.logger.Info("Server starting", "port", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error", "error", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("Server shutdown error", "error", err)
		}
		close(quit)
	}()

	<-quit
	return nil
}
