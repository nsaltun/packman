.PHONY: help build run test clean docker-up docker-down docker-logs docker-clean postgres-up fmt lint

# Variables
APP_NAME=packman

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build the Go application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-up      - Start all services with docker-compose"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-logs    - Show docker-compose logs"
	@echo "  make docker-clean   - Stop services and remove volumes"
	@echo "  make postgres-up    - Start PostgreSQL service"
	@echo "  make fmt            - Format Go code"
	@echo "  make lint           - Run golangci-lint (if installed)"

# Build the Go application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/main.go
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application locally
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Start all services with docker-compose
docker-up:
	@echo "Starting services with docker-compose..."
	@docker compose up -d --build
	@echo "Services started successfully"
	@echo "Application: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"

# Stop all services
docker-down:
	@echo "Stopping services..."
	@docker compose down
	@echo "Services stopped"

# Show docker-compose logs
docker-logs:
	@docker compose logs -f

# Clean Docker resources (including volumes)
docker-clean:
	@echo "Stopping services and removing volumes..."
	@docker compose down -v
	@echo "Docker cleanup complete"

postgres-up:
	@echo "Starting PostgreSQL service..."
	@docker compose up -d postgres
	@echo "PostgreSQL started successfully"

# Format Go code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Run golangci-lint (if installed)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
	fi
