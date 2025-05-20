# Variables
PROJECT_NAME = coder-helper
GO_FILES = $(shell find . -name '*.go')
GO = go
GOFLAGS = -v
TEST_FLAGS = -coverprofile=coverage.out

X86_64_OUTPUT=$(PROJECT_NAME)-amd64
ARM64_OUTPUT=$(PROJECT_NAME)-arm64

# Default target
all: build

# Build the project
build:
	@echo "Building $(PROJECT_NAME)..."
	@$(GO) build $(GOFLAGS) -o $(PROJECT_NAME) main.go

# Build for X86_64 (Linux)
build-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(X86_64_OUTPUT)

# Build for ARM64 (Linux)
build-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(ARM64_OUTPUT)

# Optional: Cross-platform build for macOS (x86_64)
build-macos-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos-amd64

# Optional: Cross-platform build for macOS (ARM64/M1)
build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-macos-arm64


# Run tests
test:
	@echo "Running tests..."
	@$(GO) test $(TEST_FLAGS) ./...

# Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@$(GO) test -race $(TEST_FLAGS) ./...

# Generate test coverage report
coverage:
	@echo "Generating test coverage report..."
	@$(GO) tool cover -func=coverage.out

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(PROJECT_NAME) coverage.out $(X86_64_OUTPUT) $(ARM64_OUTPUT) $(PROJECT_NAME)-macos-*

# Run the application
run:
	@echo "Running $(PROJECT_NAME)..."
	@./$(PROJECT_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the project"
	@echo "  make test     - Run tests"
	@echo "  make test-race - Run tests with race detector"
	@echo "  make coverage - Generate test coverage report"
	@echo "  make clean    - Clean up build artifacts"
	@echo "  make run      - Run the application"
	@echo " 