package main

import (
	"log/slog"
	"os"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/server"
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

	// Create and start server
	srv := server.New(cfg, log)
	if err := srv.Start(); err != nil {
		log.Error("Server error", "error", err)
		os.Exit(1)
	}
}
