package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

// HealthService provides methods for health checking
type HealthService interface {
	// CheckDBHealth verifies database connectivity
	CheckDBHealth(ctx context.Context) error
}

// HealthServiceImpl implements HealthService
type HealthServiceImpl struct {
	db *sqlx.DB
}

// NewHealthService creates a new HealthService
func NewHealthService(db *sqlx.DB) HealthService {
	return &HealthServiceImpl{
		db: db,
	}
}

// CheckDBHealth implements database health check
func (s *HealthServiceImpl) CheckDBHealth(ctx context.Context) error {
	// Set a timeout for the health check
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute a simple query to check database connectivity
	var result int
	err := s.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Database is responsive but returned no rows
		}
		return err
	}

	// Verify we got the expected result
	if result != 1 {
		return errors.New("unexpected database health check result")
	}

	return nil
}
