.PHONY: build run clean tidy

# Default city for weather query
CITY ?= Hefei

# Build the binary
build:
	go build -o llm-function-call .

# Run the program with optional city argument
run:
	go run . $(CITY)

# Run with a specific city (Chinese)
run-hefei:
	go run . 合肥

run-wuhu:
	go run . 芜湖

# Clean build artifacts
clean:
	rm -f llm-function-call

# Download and tidy dependencies
tidy:
	go mod tidy

# Install dependencies
deps:
	go mod download

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  run        - Run the program (use CITY=<city> to specify city)"
	@echo "  run-hefei  - Run with 合肥 (Hefei)"
	@echo "  run-wuhu   - Run with 芜湖 (Wuhu)"
	@echo "  clean      - Remove build artifacts"
	@echo "  tidy       - Tidy go.mod"
	@echo "  deps       - Download dependencies"
	@echo ""
	@echo "Example:"
	@echo "  make run CITY=Shanghai"

