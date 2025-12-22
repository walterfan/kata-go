# GoPostman - REST API Client Makefile
# Author: Walter Fan
# Description: Build and development automation for GoPostman

# Variables
APP_NAME := gopostman
BINARY_NAME := $(APP_NAME)
VERSION := 1.0.0
BUILD_DIR := build
DIST_DIR := dist
GO_FILES := $(wildcard *.go)
CONFIG_FILE := config.yaml

# Go build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"
BUILD_FLAGS := -trimpath -mod=readonly

# Platform detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    PLATFORM := linux
    EXT := 
endif
ifeq ($(UNAME_S),Darwin)
    PLATFORM := darwin
    EXT := 
endif
ifeq ($(OS),Windows_NT)
    PLATFORM := windows
    EXT := .exe
endif

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)$(EXT)

$(BUILD_DIR)/$(BINARY_NAME)$(EXT): $(GO_FILES) $(CONFIG_FILE)
	@echo "🔨 Building $(APP_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) .
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(EXT)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	@echo "Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64$(EXT) .
	
	# macOS AMD64
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64$(EXT) .
	
	# macOS ARM64 (Apple Silicon)
	@echo "Building for macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64$(EXT) .
	
	# Windows AMD64
	@echo "Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ Multi-platform build completed in $(DIST_DIR)/"

# Run the application
.PHONY: run
run: $(CONFIG_FILE)
	@echo "🚀 Running $(APP_NAME)..."
	@go run .

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "🔥 Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ 'air' not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "💡 Falling back to regular run..."; \
		$(MAKE) run; \
	fi

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "🔄 Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Lint the code
.PHONY: lint
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ 'golangci-lint' not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "💡 Running basic go vet instead..."; \
		go vet ./...; \
	fi

# Format the code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

# Install the binary to system
.PHONY: install
install: build
	@echo "📦 Installing $(APP_NAME) to system..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME)$(EXT) /usr/local/bin/
	@echo "✅ $(APP_NAME) installed to /usr/local/bin/"

# Uninstall the binary from system
.PHONY: uninstall
uninstall:
	@echo "🗑️ Uninstalling $(APP_NAME) from system..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)$(EXT)
	@echo "✅ $(APP_NAME) uninstalled"

# Create release packages
.PHONY: release
release: build-all
	@echo "📦 Creating release packages..."
	@mkdir -p $(DIST_DIR)/packages
	
	# Create tar.gz packages
	@for binary in $(DIST_DIR)/$(BINARY_NAME)-*; do \
		platform=$$(basename $$binary | sed 's/$(BINARY_NAME)-//' | sed 's/\.exe//'); \
		tar -czf $(DIST_DIR)/packages/$(BINARY_NAME)-v$(VERSION)-$$platform.tar.gz -C $(DIST_DIR) $(BINARY_NAME)-$$platform$(EXT) $(CONFIG_FILE) README.md; \
		echo "Created: $(BINARY_NAME)-v$(VERSION)-$$platform.tar.gz"; \
	done
	
	@echo "✅ Release packages created in $(DIST_DIR)/packages/"

# Show help
.PHONY: help
help:
	@echo "🚀 GoPostman - REST API Client"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build the application for current platform"
	@echo "  build-all     - Build for multiple platforms (Linux, macOS, Windows)"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires 'air')"
	@echo "  deps          - Install dependencies"
	@echo "  deps-update   - Update dependencies to latest versions"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Lint the code (requires 'golangci-lint')"
	@echo "  fmt           - Format the code"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install binary to system (/usr/local/bin)"
	@echo "  uninstall     - Uninstall binary from system"
	@echo "  release       - Create release packages for all platforms"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build    - Build for current platform"
	@echo "  make run      - Run the application"
	@echo "  make dev      - Development with hot reload"
	@echo "  make clean    - Clean all build artifacts"

# Development tools installation
.PHONY: install-tools
install-tools:
	@echo "🛠️ Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Check if required tools are installed
.PHONY: check-tools
check-tools:
	@echo "🔍 Checking development tools..."
	@echo -n "Go: "
	@go version
	@echo -n "Air: "
	@if command -v air > /dev/null; then echo "✅ Installed"; else echo "❌ Not installed (run 'make install-tools')"; fi
	@echo -n "golangci-lint: "
	@if command -v golangci-lint > /dev/null; then echo "✅ Installed"; else echo "❌ Not installed (run 'make install-tools')"; fi

# Show project info
.PHONY: info
info:
	@echo "📋 Project Information:"
	@echo "  Name: $(APP_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Platform: $(PLATFORM)"
	@echo "  Go Version: $(shell go version)"
	@echo "  Build Directory: $(BUILD_DIR)"
	@echo "  Distribution Directory: $(DIST_DIR)"
	@echo "  Source Files: $(words $(GO_FILES)) Go files"
