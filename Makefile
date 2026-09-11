# Variables

APP_NAME := trade-now
GO := go
BINARY_NAME := $(APP_NAME)


SHELL := /bin/bash
SWAG := swag



# Default target
all: build


# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) ./cmd/server

# Run the application
run: build
	./$(BINARY_NAME)

# Generate swagger docs
docs:
	$(SWAG) init -g main.go --dir cmd/server,internal/user,internal/auth,internal/http -o ./docs


# Run all tests
test:
	@echo "Running tests..."
	$(GO) test ./...

# Run tests with verbose output
test-v:
	@echo "Running tests (verbose)..."
	$(GO) test -v ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Phony targets (do not conflict with file names)
.PHONY: all build run test test-v test-coverage clean docs