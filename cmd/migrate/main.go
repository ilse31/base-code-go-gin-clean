package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun/driver/pgdriver"

	"base-code-go-gin-clean/internal/config"
	db2 "base-code-go-gin-clean/internal/pkg/db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// CLI flags
	down := flag.Bool("down", false, "Rollback the last migration")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Load .env
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		log.Fatalf("❌ Failed to resolve .env path: %v", err)
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️ Warning: Could not load .env file at %s: %v", envPath, err)
	} else {
		log.Printf("✅ Loaded .env from: %s", envPath)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to DB
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.GetDSN())))
	defer func() {
		if err := sqldb.Close(); err != nil {
			log.Printf("⚠️ Warning: Failed to close DB: %v", err)
		}
	}()

	// Run migration
	if err := db2.RunMigrations(sqldb, "postgres", *down); err != nil {
		log.Fatalf("❌ Migration error: %v", err)
	}

	log.Println("✅ Migration completed successfully")
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run .         → Apply all pending migrations")
	fmt.Println("  go run . -down   → Rollback the last migration")
	fmt.Println("  go run . -h      → Show this help message")
}
