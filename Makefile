.PHONY: help build run test test-race vet fmt clean docker-build docker-run docker-stop lint install

# Default target
help:
	@echo "Golang Birmingham Website - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run        Run the development server"
	@echo "  make build      Build the application"
	@echo "  make test       Run tests"
	@echo "  make test-race  Run tests with race detection"
	@echo "  make vet        Run go vet"
	@echo "  make fmt        Format code with gofmt"
	@echo "  make lint       Run linter"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build    Build Docker image"
	@echo "  make docker-run      Run with Docker Compose"
	@echo "  make docker-stop     Stop Docker Compose"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean      Clean build artifacts"
	@echo "  make install    Install dependencies"

# Development
run:
	@echo "Starting development server..."
	go run ./cmd/web

build:
	@echo "Building application..."
	go build -o bin/web ./cmd/web

test:
	@echo "Running tests..."
	go test ./...

test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

vet:
	@echo "Running go vet..."
	go vet ./...

fmt:
	@echo "Formatting code..."
	gofmt -s -w .

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.59.1"; \
	fi

# Dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t golang-birmingham-website .

docker-run:
	@echo "Starting with Docker Compose..."
	docker compose up -d

docker-stop:
	@echo "Stopping Docker Compose..."
	docker compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker compose logs -f

# Production
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/web-linux ./cmd/web

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f web
	docker compose down --volumes --remove-orphans 2>/dev/null || true

# Development helpers
dev-setup: install
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.59.1; \
	fi

# Content validation
validate-content:
	@echo "Validating content files..."
	go test ./internal/content -run TestLoadRepositoryContent -count=1

# Pre-commit checks
pre-commit: fmt vet test lint
	@echo "Pre-commit checks passed!"

# Release helpers
release-check: test vet lint
	@echo "Release checks completed successfully!"

# Quick development start
dev: dev-setup run

# All checks
check-all: fmt vet test-race lint
	@echo "All checks completed!"
