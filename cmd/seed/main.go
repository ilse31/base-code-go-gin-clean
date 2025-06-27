package main

import (
	"flag"
	"fmt"
	"log"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/seeders"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Parse command line flags
	flag.Parse()

	// Run seeders
	err = seeders.SeedAll(db.DB)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Println("✅ Database seeded successfully!")
}
