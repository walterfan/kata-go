# Variables
PROJECT_NAME = service_monitor
GO_FILES = $(shell find . -name '*.go')
GO = go
GOFLAGS = -v
TEST_FLAGS = -coverprofile=coverage.out

# Default target
all: build

# Build the project
build:
	@echo "Building $(PROJECT_NAME)..."
	@$(GO) build $(GOFLAGS) -o $(PROJECT_NAME) main.go

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
	@rm -f $(PROJECT_NAME) coverage.out

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