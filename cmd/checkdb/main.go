package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"base-code-go-gin-clean/internal/config"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create new DB connection
	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to create DB connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	fmt.Println("Pinging database...")
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Try a simple query
	var result int
	fmt.Println("Executing test query...")
	err = db.NewSelect().ColumnExpr("1").Scan(ctx, &result)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	fmt.Println("✅ Database connection successful!")
	fmt.Printf("Test query result: %d\n", result)

	// Print database info
	dbUser := cfg.DB.User
	dbHost := cfg.DB.Host
	dbPort := cfg.DB.Port
	dbName := cfg.DB.Name

	dbUser = fmt.Sprintf("user=%s", dbUser)
	dbHost = fmt.Sprintf("host=%s", dbHost)
	dbPort = fmt.Sprintf("port=%s", dbPort)
	dbName = fmt.Sprintf("dbname=%s", dbName)

	if cfg.DB.SSLMode != "" {
		fmt.Printf("SSL Mode: %s\n", cfg.DB.SSLMode)
	}

	fmt.Printf("\nConnected to database:\n")
	fmt.Printf("User:     %s\n", dbUser)
	fmt.Printf("Host:     %s\n", dbHost)
	fmt.Printf("Port:     %s\n", dbPort)
	fmt.Printf("Database: %s\n", dbName)
}
