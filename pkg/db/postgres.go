package db

import (
	"base-code-go-gin-clean/internal/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*bun.DB, error) {
	dsn := cfg.DB.GetDSN()

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create Bun DB instance
	db := bun.NewDB(sqlDB, pgdialect.New())

	// Add query hooks for debugging
	if cfg.Server.Environment == "development" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithEnabled(cfg.Server.Environment == "development"),
		))
	}

	return db, nil
}
