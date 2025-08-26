# Blog Generator Makefile
# A comprehensive Makefile for the Go blog generation project

# Variables
BINARY_NAME=bloggen
BINARY_UNIX=$(BINARY_NAME)_unix
BUILD_DIR=build
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go related variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_DIR=$(BUILD_DIR)/$(BINARY_NAME)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# Default target
.DEFAULT_GOAL := build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v
	@echo "Build complete!"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Build for current platform
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v
	
	# Build for Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 -v
	
	# Build for macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v
	
	# Build for Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe -v
	
	@echo "Multi-platform build complete!"

# Build with version information
.PHONY: build-version
build-version:
	@echo "Building $(BINARY_NAME) with version info..."
	$(GOBUILD) -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)" -o $(BINARY_NAME) -v
	@echo "Versioned build complete!"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run main.go

# Run with custom parameters
.PHONY: run-custom
run-custom:
	@echo "Running $(BINARY_NAME) with custom parameters..."
	$(GOCMD) run main.go --idea "Building AI Applications" --location "Shanghai" --model "gpt-4o"

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f blog-*.md
	@echo "Clean complete!"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed!"

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy
	@echo "Dependencies updated!"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Code formatted!"

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GOCMD) vet ./...
	@echo "Linting complete!"

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./...
	@echo "Documentation generated!"

# Install the binary
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Uninstall the binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete!"

# Run weather tool test
.PHONY: test-weather
test-weather:
	@echo "Testing weather tool..."
	$(GOCMD) run examples/test_weather_parsing.go

# Test API connection
.PHONY: test-api
test-api:
	@echo "Testing API connection..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Please create one with your API key:"; \
		echo "cp env.example .env"; \
		echo "Then edit .env and add your LLM_API_KEY"; \
		exit 1; \
	fi
	@echo "✅ .env file found"
	@echo "Testing API configuration..."
	@$(GOCMD) run main.go --help > /dev/null 2>&1 || echo "❌ API configuration test failed"

# Run the weather tool server test
.PHONY: test-weather-server
test-weather-server:
	@echo "Testing weather tool server..."
	./examples/test_weather.sh

# Generate a sample blog
.PHONY: generate-blog
generate-blog:
	@echo "Generating sample blog..."
	$(GOCMD) run main.go --idea "Testing Blog Generation" --location "Hefei"

# Generate blog with streaming
.PHONY: generate-blog-stream
generate-blog-stream:
	@echo "Generating blog with streaming..."
	LLM_STREAM=true $(GOCMD) run main.go --idea "Streaming Blog Generation" --location "Beijing"

# Show help
.PHONY: help
help:
	@echo "Blog Generator Makefile"
	@echo "======================"
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-all          - Build for multiple platforms"
	@echo "  build-version      - Build with version information"
	@echo "  run                - Run the application"
	@echo "  run-custom         - Run with custom parameters"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-weather       - Test weather parsing"
	@echo "  test-api           - Test API connection"
	@echo "  test-weather-server- Test weather tool server"
	@echo "  generate-blog      - Generate a sample blog"
	@echo "  generate-blog-stream- Generate blog with streaming"
	@echo "  deps               - Install dependencies"
	@echo "  deps-update        - Update dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  docs               - Generate documentation"
	@echo "  install            - Install binary"
	@echo "  uninstall          - Uninstall binary"
	@echo "  clean              - Clean build artifacts"
	@echo "  help               - Show this help"

# Development workflow
.PHONY: dev
dev: deps fmt lint test build
	@echo "Development workflow complete!"

# Release workflow
.PHONY: release
release: clean deps fmt lint test build-all
	@echo "Release workflow complete!"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "Docker build complete!"

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it $(BINARY_NAME):$(VERSION)

# Check if .env file exists
.PHONY: check-env
check-env:
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found. Copy env.example to .env and configure your settings."; \
		echo "cp env.example .env"; \
	fi

# Setup development environment
.PHONY: setup
setup: check-env deps
	@echo "Development environment setup complete!"
	@echo "Don't forget to configure your .env file with your OpenAI API key."

# Show project info
.PHONY: info
info:
	@echo "Blog Generator Project Info"
	@echo "=========================="
	@echo "Binary Name: $(BINARY_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell go version)"
	@echo "OS/Arch: $(shell go env GOOS)/$(shell go env GOARCH)" 