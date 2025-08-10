# Makefile for backend commands

# This target ensures that the .env file exists.
# It will be run automatically before any migrate-* command.
.env:
	@if [ ! -f .env ]; then \
		echo ".env file not found. Please create it from .env.example."; \
		exit 1; \
	fi

# Include the .env file to load its variables into the Makefile's environment.
# This is a more reliable method for loading credentials.
include .env
# Export the variables so they are available to sub-shells (like the one running 'migrate').
export

# Database connection string for migration tool.
# It now uses the variables loaded and exported above.
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Create a new migration file with a given name.
# Example: make migrate-create name=add_decks_table
migrate-create: .env
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir db/migrations -seq $(name)

# Apply all up migrations.
migrate-up: .env
	@echo "Applying migrations..."
	@migrate -path db/migrations -database "$(DB_URL)" -verbose up

# Rollback the last migration.
migrate-down: .env
	@echo "Rolling back last migration..."
	@migrate -path db/migrations -database "$(DB_URL)" -verbose down

.PHONY: migrate-create migrate-up migrate-down

