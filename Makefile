# Check for .env file and source it
ifneq (,$(wildcard ./.env))
    include .env
    export
else
    $(error .env file not found. Please create one before running make commands)
endif

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/app/main.go

# Run the application
run: build
	@go run cmd/app/main.go

# Run docker containers
docker-run:
	@if docker compose ps 2>/dev/null | grep -q "Up"; then \
		echo "Attaching to existing Docker Compose session"; \
		docker compose logs -f; \
	else \
		echo "Starting new Docker Compose session"; \
		docker compose up; \
	fi

# Shutdown docker containers
docker-down:
	@docker compose down

# Clean containers
docker-clean:
	@docker compose down -v

# Database migration status
db-status:
	@go run cmd/migrate/main.go -command status

# Database apply migrations
db-migrate:
	@go run cmd/migrate/main.go -command up

# Database Re-apply last migration
db-migrate-redo:
	@go run cmd/migrate/main.go -command redo

# Database downgrade
db-downgrade:
	@go run cmd/migrate/main.go -command down

# Database seed
db-seed:
	@go run cmd/migrate/main.go -command seed

# Test the application
test:
	@echo "Testing..."
	@go test -count=1 ./... -v

# Clean project
clean:
	@echo "Cleaning..."
	@docker compose down -v
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching..."; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching..."; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all setup build run docker-run docker-down db-status db-migrate db-migrate-redo db-downgrade test clean watch
