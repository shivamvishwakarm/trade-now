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


# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)

# Phony targets (do not conflict with file names)
.PHONY: all build run test clean docs