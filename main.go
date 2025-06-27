package main

import (
	"log/slog"
	"os"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler"
	"base-code-go-gin-clean/internal/repository"
	"base-code-go-gin-clean/internal/server"
	"base-code-go-gin-clean/internal/service"
	"base-code-go-gin-clean/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Server.Environment)

	// Initialize database
	dbWrapper, err := config.NewDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	dbConn := dbWrapper.DB

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbConn)

	// Initialize services
	userSvc := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userSvc)

	// Create server
	srv := server.New(cfg, log)

	// Setup routes
	srv.SetupUserRoutes(userHandler)

	// Start server
	if err := srv.Start(); err != nil {
		log.Error("Server error", "error", err)
		os.Exit(1)
	}
}
