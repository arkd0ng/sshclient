# Makefile for SSH Client

# Variables
APP_NAME=sshclient
VERSION=1.2.0
BUILD_DIR=bin
SRC_DIR=src
MAIN_FILE=$(SRC_DIR)/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Output binaries
BINARY_NAME=$(BUILD_DIR)/$(APP_NAME)
BINARY_WINDOWS=$(BUILD_DIR)/$(APP_NAME).exe
BINARY_LINUX=$(BUILD_DIR)/$(APP_NAME)-linux
BINARY_DARWIN=$(BUILD_DIR)/$(APP_NAME)-darwin
BINARY_DARWIN_ARM64=$(BUILD_DIR)/$(APP_NAME)-darwin-arm64

.PHONY: all build clean test deps help install run

# Default target
all: clean build

# Help target
help:
	@echo "SSH Client - Makefile commands:"
	@echo ""
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make windows       - Build for Windows (64-bit)"
	@echo "  make linux         - Build for Linux (64-bit)"
	@echo "  make darwin        - Build for macOS Intel (64-bit)"
	@echo "  make darwin-arm64  - Build for macOS Apple Silicon"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make deps          - Download dependencies"
	@echo "  make install       - Install to $$GOPATH/bin"
	@echo "  make run           - Build and run"
	@echo "  make version       - Show version"
	@echo ""

# Build for current platform
build:
	@echo "Building $(APP_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(SRC_DIR)/*.go
	@echo "Build complete: $(BINARY_NAME)"

# Build for all platforms
build-all: windows linux darwin darwin-arm64
	@echo "All builds complete!"

# Build for Windows
windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) $(SRC_DIR)/*.go
	@echo "Build complete: $(BINARY_WINDOWS)"

# Build for Linux
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_LINUX) $(SRC_DIR)/*.go
	@echo "Build complete: $(BINARY_LINUX)"

# Build for macOS Intel
darwin:
	@echo "Building for macOS (Intel)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DARWIN) $(SRC_DIR)/*.go
	@echo "Build complete: $(BINARY_DARWIN)"

# Build for macOS Apple Silicon
darwin-arm64:
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DARWIN_ARM64) $(SRC_DIR)/*.go
	@echo "Build complete: $(BINARY_DARWIN_ARM64)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./$(SRC_DIR)/...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOGET) -v -t -d ./$(SRC_DIR)/...
	@$(GOMOD) tidy
	@echo "Dependencies updated"

# Install to $GOPATH/bin
install: build
	@echo "Installing to $$GOPATH/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/$(APP_NAME)
	@echo "Install complete: $(GOPATH)/bin/$(APP_NAME)"

# Build and run
run: build
	@echo "Running $(APP_NAME)..."
	@$(BINARY_NAME)

# Show version
version:
	@echo "SSH Client v$(VERSION)"

# Development helpers
.PHONY: dev fmt lint

# Run in development mode (rebuild on changes)
dev:
	@echo "Development mode - building and running..."
	@$(MAKE) build
	@$(BINARY_NAME) -h

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./$(SRC_DIR)/...
	@echo "Format complete"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run ./$(SRC_DIR)/... || echo "Note: golangci-lint not installed"

# Quick build and test
quick: build test
	@echo "Quick build and test complete"
