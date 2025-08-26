# Simple HTTP Server with Gin

A simple HTTP server built with Go and the Gin framework that provides user CRUD operations, command execution, and graceful shutdown capabilities.

## Features

- **Health Check**: GET `/health` endpoint returns server status
- **User Management**: Full CRUD operations for users via `/users` endpoints
- **Command Execution**: POST `/commands` endpoint for executing whitelisted shell commands
- **Data Persistence**: Users are stored in a configurable data directory (default: `/data/users.json`)
- **Configuration**: YAML-based configuration for server port and command whitelist
- **Graceful Shutdown**: Proper signal handling and graceful server shutdown

## API Endpoints

### Health Check
- `GET /health` - Returns `{"status": "ok"}`

### Users
- `GET /users` - Get all users
- `GET /users/:id` - Get a specific user by ID
- `POST /users` - Create a new user
- `PUT /users/:id` - Update an existing user
- `DELETE /users/:id` - Delete a user

### Commands
- `POST /commands` - Execute a whitelisted shell command

## Configuration

The server configuration is stored in `config.yaml`:

```yaml
server:
  port: 8080

data:
  data_dir: "/data"

commands:
  whitelist:
    - "ls -la"
    - "pwd"
    - "whoami"
    - "date"
    - "echo 'Hello World'"
    - "uname -a"
```

## Building and Running

### Prerequisites
- Go 1.21 or later

### Build the application
```bash
go mod tidy
go build -o server main.go
```

### Run the server
```bash
./server
```

The server will start on the configured port (default: 8080).

### Stop the server
Press `Ctrl+C` to gracefully shutdown the server.

## Docker Support

The project includes a multi-stage Dockerfile that:
- Uses `golang:1.22-alpine` for building the application
- Creates a minimal runtime image using `scratch`
- Results in a very small container size
- Includes only the necessary runtime dependencies

### Build Docker image
```bash
make docker-build
# or
docker build -t scratch-verification .
```

### Run with Docker
```bash
make docker-run
# or
docker run -p 8080:8080 --name scratch-server scratch-verification
```

### Run with Docker Compose
```bash
docker-compose up -d
```

### Stop Docker container
```bash
make docker-stop
# or
docker-compose down
```

### Clean up Docker resources
```bash
make docker-clean
# or
docker-compose down --rmi all
```

## Example Usage

### Create a user
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user1",
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }'
```

### Get all users
```bash
curl http://localhost:8080/users
```

### Execute a command
```bash
curl -X POST http://localhost:8080/commands \
  -H "Content-Type: application/json" \
  -d '{"command": "pwd"}'
```

### Health check
```bash
curl http://localhost:8080/health
```

## File Structure

- `main.go` - Main application code
- `config.yaml` - Configuration file
- `go.mod` - Go module definition
- `users.json` - User data storage (created automatically)
- `Dockerfile` - Multi-stage Docker build file
- `docker-compose.yml` - Docker Compose configuration
- `.dockerignore` - Docker build context exclusions
- `Makefile` - Build and run commands
- `test_server.sh` - Comprehensive API testing script
- `README.md` - This file

## Security Notes

- Only whitelisted commands can be executed
- Commands are executed with shell privileges
- User input is validated before processing
- Consider restricting the command whitelist in production environments

## Graceful Shutdown

The server handles SIGINT and SIGTERM signals gracefully:
- Stops accepting new connections
- Waits for existing requests to complete (up to 30 seconds)
- Saves any pending user data
- Exits cleanly
