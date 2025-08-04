.PHONY: build run test clean deps lint format help

# Variables
BINARY_NAME=clean-architecture-api
BUILD_DIR=build
MAIN_FILE=cmd/server/main.go

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_FILE)

# Run with SQLite database (no Docker required)
run-sqlite:
	@echo "Running $(BINARY_NAME) with SQLite database..."
	@go run -tags sqlite cmd/server/main_sqlite.go

# Run with in-memory database (for testing)
run-memory:
	@echo "Running $(BINARY_NAME) with in-memory database..."
	@go run cmd/server/main.go

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Running in development mode with hot reload..."
	@air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Lint and fix code with comprehensive formatting
lint-fix:
	@echo "Linting and fixing code..."
	@golangci-lint run --fix
	@make format

# Format code with Go standard tools (NEW - comprehensive formatting)
format:
	@echo "Formatting Go code with standard tools..."
	@echo "Running gofmt..."
	@gofmt -w -s .
	@echo "Running goimports..."
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		if [ -n "$$(go env GOBIN)" ]; then \
			$$(go env GOBIN)/goimports -w .; \
		elif [ -n "$$(go env GOPATH)" ]; then \
			$$(go env GOPATH)/bin/goimports -w .; \
		else \
			echo "Could not find goimports after installation. Please run 'make format-deps' first."; \
			exit 1; \
		fi; \
	fi
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "Running go vet..."
	@go vet ./...
	@echo "Checking for ineffectual assignments..."
	@if command -v ineffassign >/dev/null 2>&1; then \
		ineffassign ./...; \
	else \
		echo "ineffassign not installed, skipping..."; \
	fi
	@echo "Checking for misspellings..."
	@if command -v misspell >/dev/null 2>&1; then \
		misspell -w .; \
	else \
		echo "misspell not installed, skipping..."; \
	fi
	@echo "Code formatting completed!"

# Install formatting dependencies
format-deps:
	@echo "Installing formatting dependencies..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@echo "Formatting dependencies installed!"

# Generate swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g $(MAIN_FILE)

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME)

# Database commands
db-migrate:
	@echo "Running database migrations..."
	@go run $(MAIN_FILE) migrate

db-seed:
	@echo "Seeding database..."
	@go run $(MAIN_FILE) seed

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  lint-fix      - Lint and fix code with formatting"
	@echo "  format        - Format code with Go standard tools (gofmt, goimports, etc.)"
	@echo "  format-deps   - Install formatting tool dependencies"
	@echo "  swagger       - Generate swagger docs (requires swag)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  db-migrate    - Run database migrations"
	@echo "  db-seed       - Seed database"
	@echo "  help          - Show this help message"
