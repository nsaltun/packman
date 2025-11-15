.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean logs fmt lint

# Variables
APP_NAME=packman
DOCKER_IMAGE=packman:latest
DOCKER_CONTAINER=packman-container

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build the Go application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-stop    - Stop Docker containers"
	@echo "  make docker-clean   - Remove Docker containers and images"
	@echo "  make logs           - Show Docker container logs"
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

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built successfully"

# Run application in Docker
docker-run: docker-build
	@echo "Starting $(APP_NAME) container..."
	@docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		$(DOCKER_IMAGE)
	@echo "$(APP_NAME) is running on http://localhost:8080"

# Stop Docker containers
docker-stop:
	@echo "Stopping containers..."
	@docker stop $(DOCKER_CONTAINER)  2>/dev/null || true
	@docker rm $(DOCKER_CONTAINER) 2>/dev/null || true
	@echo "Containers stopped"

# Clean Docker resources
docker-clean: docker-stop
	@echo "Removing Docker image..."
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@echo "Docker cleanup complete"

# Show Docker logs
logs:
	@docker logs -f $(DOCKER_CONTAINER)

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
