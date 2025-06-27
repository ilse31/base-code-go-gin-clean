package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/uptrace-go/uptrace"
)

// DB wraps bun.DB to provide database operations
type DB struct {
	*bun.DB
}

// NewDB creates a new database connection using the provided configuration
func NewDB(cfg *Config) (*DB, error) {
	// Initialize database connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.GetDSN())))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query debug hook if in development
	if cfg.Server.Environment == "development" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	// Configure tracing if enabled
	if cfg.Tracing.Enabled && cfg.Tracing.DSN != "" {
		uptrace.ConfigureOpentelemetry(
			uptrace.WithDSN(cfg.Tracing.DSN),
			uptrace.WithServiceName(cfg.Tracing.ServiceName),
			uptrace.WithServiceVersion(cfg.Tracing.Version),
		)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}
