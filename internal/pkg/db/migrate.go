package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs all database migrations in the specified direction
func RunMigrations(db *sql.DB, dbName string, down bool) error {
	// Create a temporary directory to store the migrations
	tmpDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write the embedded migrations to the temp directory
	err = fs.WalkDir(migrationsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := migrationsFS.ReadFile(path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tmpDir, filepath.Base(path))
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to write migrations to temp dir: %w", err)
	}

	// Create a new Postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+tmpDir,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations in the specified direction
	if down {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to rollback migrations: %w", err)
		}
		fmt.Println("Successfully rolled back migrations")
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		fmt.Println("Successfully applied migrations")
	}

	return nil
}
