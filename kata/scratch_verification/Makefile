.PHONY: build run test clean

# Build the server
build:
	go mod tidy
	go build -o server main.go

# Run the server
run: build
	./server

# Run tests (if you add test files later)
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f server
	rm -f users.json

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run the test script (requires server to be running)
test-api:
	./test_server.sh

# Docker commands
docker-build:
	docker build --build-arg GOPROXY=$(GOPROXY) -t scratch-verification .

docker-run:
	docker run -p 8080:8080 --name scratch-server scratch-verification

docker-stop:
	docker stop scratch-server || true
	docker rm scratch-server || true

docker-clean:
	docker rmi scratch-verification || true

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the server binary"
	@echo "  run          - Build and run the server"
	@echo "  test         - Run Go tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  test-api     - Run API tests (requires server running)"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-stop  - Stop and remove Docker container"
	@echo "  docker-clean - Remove Docker image"
	@echo "  help         - Show this help message"
