# Vault - Simple HTTP Server with Gin

A simple HTTP server built with Go and the Gin framework that provides user CRUD operations, command execution, and graceful shutdown capabilities.

## Features

- **Health Check**: GET `/health` endpoint returns server status
- **Site Management**: Full CRUD operations for sites via `/sites` endpoints
- **Command Execution**: POST `/commands` endpoint for executing whitelisted shell commands
- **Data Persistence**: Sites are stored in SQLite database (default: `/data/sites.db`)
- **Password Encryption**: Passwords are encrypted using AES-GCM with SHA256 key derivation
- **Configuration**: YAML-based configuration with go-viper for flexible config management
- **Structured Logging**: HTTP request/response logging with zap and lumberjack for log rotation
- **Graceful Shutdown**: Proper signal handling and graceful server shutdown

## API Endpoints

### Health Check
- `GET /health` - Returns `{"status": "ok"}`

### Sites
- `GET /sites` - Get all sites
- `GET /sites/:id` - Get a specific site by ID
- `POST /sites` - Create a new site
- `PUT /sites/:id` - Update an existing site
- `DELETE /sites/:id` - Delete a site

### Commands
- `POST /commands` - Execute a whitelisted shell command

## Configuration

The server configuration is stored in `config.yaml`:

```yaml
server:
  port: 8080

data:
  data_dir: "/data"

logging:
  level: "info"
  file:
    filename: "/data/app.log"
    max_size: 100
    max_age: 30
    max_backups: 10
    compress: true

commands:
  whitelist:
    - "ls -la"
    - "pwd"
    - "whoami"
    - "date"
    - "echo 'Hello World'"
    - "uname -a"
```

### Configuration Options

- **server.port**: HTTP server port (default: 8080)
- **data.data_dir**: Directory for storing user data (default: "/data")
- **logging.level**: Log level - debug, info, warn, error (default: "info")
- **logging.file.filename**: Log file path (default: "/data/app.log")
- **logging.file.max_size**: Maximum log file size in MB (default: 100)
- **logging.file.max_age**: Maximum age of log files in days (default: 30)
- **logging.file.max_backups**: Maximum number of backup files (default: 10)
- **logging.file.compress**: Compress rotated log files (default: true)
- **commands.whitelist**: List of allowed shell commands

**Note**: All configuration values can be overridden using environment variables. For example, `SERVER_PORT=9000` will override the server port.

### Environment Variables

- **AES_KEY**: Required environment variable for password encryption. This string will be hashed with SHA256 to create the AES-256 key.

**Option 1: Using .env file (Recommended for development)**:
```bash
# Create .env file
cp .env.example .env
# Edit .env file with your AES key
```

**Option 2: Using system environment variables**:
```bash
export AES_KEY="my-secret-encryption-key-2024"
```

**Example .env file**:
```env
AES_KEY=my-secret-encryption-key-2024
```

## Building and Running

### Prerequisites
- Go 1.21 or later

### Build the application
```bash
go mod tidy
go build -o vault main.go
```

### Run the vault
```bash
# Option 1: Using .env file (recommended)
cp .env.example .env
# Edit .env file with your AES key
./vault server

# Option 2: Using environment variable
export AES_KEY="my-secret-encryption-key-2024"
./vault server
```

The vault will start on the configured port (default: 8080).

### CLI Commands

The application supports comprehensive CLI commands for password encryption/decryption and site management:

#### Basic Commands
```bash
# Show help
./vault --help

# Encrypt a password
./vault encrypt "mypassword123"

# Decrypt a password
./vault decrypt "encrypted_password_string"

# Start the HTTP server
./vault server
```

#### Sites Management Commands
```bash
# List all sites
./vault sites list

# Get a specific site
./vault sites get "site-id"

# Create a new site
./vault sites create "site-id" "Site Name" "username" "password"

# Update an existing site
./vault sites update "site-id" "New Name" "new_username" "new_password"

# Delete a site
./vault sites delete "site-id"

# Show sites help
./vault sites --help
```

### Stop the vault
Press `Ctrl+C` to gracefully shutdown the vault.

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
# Option 1: Using .env file
cp .env.example .env
# Edit .env file with your AES key
make docker-run

# Option 2: Using environment variable
export AES_KEY="my-secret-encryption-key-2024"
make docker-run

# Or with docker command
docker run -p 8080:8080 -e AES_KEY="my-secret-encryption-key-2024" --name vault-server scratch-verification server
```

### Run with Docker Compose
```bash
# Option 1: Using .env file
cp .env.example .env
# Edit .env file with your AES key
docker-compose up -d

# Option 2: Using environment variable
export AES_KEY="my-secret-encryption-key-2024"
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

### Create a site
```bash
curl -X POST http://localhost:8080/sites \
  -H "Content-Type: application/json" \
  -d '{
    "id": "site1",
    "name": "Example Site",
    "username": "admin",
    "password": "password123"
  }'
```

### Get all sites
```bash
curl http://localhost:8080/sites
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
- `sites.db` - SQLite database file (created automatically)
- `.env.example` - Example environment variables file
- `.env` - Environment variables file (create from .env.example)
- `Dockerfile` - Multi-stage Docker build file
- `docker-compose.yml` - Docker Compose configuration
- `.dockerignore` - Docker build context exclusions
- `Makefile` - Build and run commands
- `test_server.sh` - Comprehensive API testing script
- `test_encryption.sh` - Encryption functionality test script
- `example_sites.json` - Example site data structure
- `README.md` - This file

## Security Notes

- **Password Encryption**: All passwords are encrypted using AES-GCM with SHA256 key derivation
- **Key Management**: The AES key is derived from the `AES_KEY` environment variable using SHA256
- **Command Security**: Only whitelisted commands can be executed
- **Input Validation**: User input is validated before processing
- **Production Considerations**: 
  - Use a strong, unique AES_KEY in production
  - Consider restricting the command whitelist
  - Store the AES_KEY securely (e.g., in a secrets management system)

## Graceful Shutdown

The server handles SIGINT and SIGTERM signals gracefully:
- Stops accepting new connections
- Waits for existing requests to complete (up to 30 seconds)
- Saves any pending user data
- Exits cleanly
