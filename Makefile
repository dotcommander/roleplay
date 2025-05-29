.PHONY: build test clean install lint fmt run help

# Variables
BINARY_NAME=roleplay
GO=go
GOFLAGS=-v

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-s -w \
	-X 'github.com/dotcommander/roleplay/cmd.Version=$(VERSION)' \
	-X 'github.com/dotcommander/roleplay/cmd.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/dotcommander/roleplay/cmd.BuildDate=$(BUILD_DATE)'

# Default target
all: test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install the binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)" .

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Run the application
run: build
	./$(BINARY_NAME)

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Generate mocks (if needed)
mocks:
	@echo "Generating mocks..."
	$(GO) generate ./...

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms with version $(VERSION)..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Run development server with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "air not found. Install it with: go install github.com/cosmtrek/air@latest" && exit 1)
	air

# Database migrations (if applicable)
migrate-up:
	@echo "Running migrations up..."
	# Add migration commands here

migrate-down:
	@echo "Running migrations down..."
	# Add migration commands here

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install the binary"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make run           - Build and run"
	@echo "  make deps          - Update dependencies"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make help          - Show this help message"