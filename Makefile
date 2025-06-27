# Migration directory relative to the project root
MIGRATIONS_DIR=internal/pkg/db/migrations

# Go migration tool binary
MIGRATE_CMD=migrate

.PHONY: migrate migrate-up migrate-down migrate-create migrate-help

# Run all pending migrations (default)
migrate: migrate-up

# Run all pending migrations
migrate-up:
	@echo "Running pending migrations..."
	@cd cmd/migrate && go run .

# Rollback the last migration
migrate-down:
	@echo "Rolling back the last migration..."
	@cd cmd/migrate && go run . -down

# Create new migration files
# Usage: make migrate-create name=create_users_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	@mkdir -p $(MIGRATIONS_DIR)
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	up="$(MIGRATIONS_DIR)/$${timestamp}_$(name).up.sql"; \
	down="$(MIGRATIONS_DIR)/$${timestamp}_$(name).down.sql"; \
	touch "$$up" "$$down"; \
	echo "-- +goose Up\n-- +goose StatementBegin\n-- Add your SQL here\n-- +goose StatementEnd" > "$$up"; \
	echo "-- +goose Down\n-- +goose StatementBegin\n-- Add your rollback SQL here\n-- +goose StatementEnd" > "$$down"; \
	echo "Created migration files:"; \
	ls -1 "$$up" "$$down"

# Show migration help
migrate-help:
	@echo "\nMigration Commands:"
	@echo "  make migrate             # Run all pending migrations (same as make migrate-up)"
	@echo "  make migrate-up          # Apply all pending migrations"
	@echo "  make migrate-down        # Rollback the last migration"
	@echo "  make migrate-create name=create_table_name  # Create new migration files"
	@echo "  make migrate-help        # Show this help message"
	@echo "\nExample:"
	@echo "  make migrate-create name=create_users_table"
	@echo ""
