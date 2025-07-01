package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *bun.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:120579@127.0.0.1:5432/postgres?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Verify the connection
	err := db.Ping()
	assert.NoError(t, err, "Failed to connect to test database")

	// You might want to run migrations here if needed

	return db
}

// TeardownTestDB closes the database connection
func TeardownTestDB(t *testing.T, db *bun.DB) {
	if db != nil {
		err := db.Close()
		assert.NoError(t, err, "Failed to close test database connection")
	}
}
