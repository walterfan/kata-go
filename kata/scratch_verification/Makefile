.PHONY: build run test clean

# Build the vault
build:
	go mod tidy
	go build -o vault main.go

# Run the vault
run: build
	./vault server

# Run tests (if you add test files later)
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f vault
	rm -f users.json

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run the test script (requires server to be running)
test-api:
	./test_server.sh

# Run the CLI test script
test-cli:
	./test_cli.sh

# Run the sites CLI test script
test-sites-cli:
	./test_sites_cli.sh

# Docker commands
docker-build:
	docker build --build-arg GOPROXY=$(GOPROXY) -t scratch-verification .

docker-run:
	docker run -p 8080:8080 --name vault-server scratch-verification server

docker-stop:
	docker stop vault-server || true
	docker rm vault-server || true

docker-clean:
	docker rmi scratch-verification || true

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the vault binary"
	@echo "  run          - Build and run the vault"
	@echo "  test         - Run Go tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  test-api     - Run API tests (requires vault running)"
	@echo "  test-cli     - Run CLI commands tests"
	@echo "  test-sites-cli - Run sites CLI CRUD tests"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-stop  - Stop and remove Docker container"
	@echo "  docker-clean - Remove Docker image"
	@echo "  help         - Show this help message"
