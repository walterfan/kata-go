# GoPostman - REST API Client

A beautiful, modern REST API client built with Go and Fyne, inspired by Postman. Features a clean GUI with support for predefined request collections loaded from YAML configuration.

## ✨ Features

- **Beautiful Modern UI** with emoji icons and intuitive layout
- **Collection Management** - Load predefined requests from `config.yaml`
- **Full HTTP Support** - GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **Headers Management** - Add custom headers with easy syntax
- **Parameters Support** - URL parameters with automatic encoding
- **Request Body** - Support for JSON, XML, and other formats
- **JSON Formatting** - Automatic pretty-printing of JSON responses
- **Response Analysis** - Status codes, timing, and formatted output
- **Copy/Paste Support** - Easy response copying to clipboard

## 🏗️ Installation

### Quick Start
```bash
# Clone or download the project
git clone <your-repo-url>
cd kata-rest

# Build and run using Makefile (recommended)
make build
make run

# Or run directly
make run
```

### Manual Installation
```bash
# Install dependencies
go mod tidy

# Run the application
go run .
```

### Using Makefile (Recommended)
The project includes a comprehensive Makefile for easy development and deployment:

```bash
# Show all available commands
make help

# Build the application
make build

# Run the application
make run

# Development with hot reload
make dev

# Build for multiple platforms
make build-all

# Install to system
make install

# Create release packages
make release
```

## 📁 Configuration

The app loads predefined requests from `config.yaml`. Here's the structure:

```yaml
collections:
  - name: "My API Collection"
    requests:
      - name: "Get Users"
        method: "GET"
        url: "https://api.example.com/users"
        headers:
          Accept: "application/json"
          Authorization: "Bearer your-token"
        parameters:
          page: "1"
          limit: "10"
        body: ""
      
      - name: "Create User"
        method: "POST"
        url: "https://api.example.com/users"
        headers:
          Content-Type: "application/json"
        parameters: {}
        body: |
          {
            "name": "John Doe",
            "email": "john@example.com"
          }
```

## 🎮 Usage

### Loading Predefined Requests

1. **Select Collection** - Choose from collections defined in `config.yaml`
2. **Select Request** - Pick a specific request from the collection
3. **Auto-populate** - All fields (method, URL, headers, parameters, body) are automatically filled
4. **Customize** - Modify any field as needed
5. **Send** - Click the Send button

### Manual Requests

1. **Clear Form** - Click "➕ New Request" to start fresh
2. **Set Method** - Choose HTTP method (GET, POST, etc.)
3. **Enter URL** - Type your API endpoint
4. **Add Headers** (optional) - Format: `Header-Name: value` (one per line)
5. **Add Parameters** (optional) - Format: `key1=value1&key2=value2`
6. **Add Body** (optional) - JSON, XML, or any text
7. **Send Request** - Click Send

### Headers Format
```
Content-Type: application/json
Authorization: Bearer your-token-here
Accept: application/json
```

### Parameters Format
```
page=1&limit=10&sort=name&order=asc
```

### Features

- **🔄 Refresh Config** - Reload `config.yaml` without restarting
- **🎨 Format JSON** - Pretty-print JSON in request body
- **📋 Copy Response** - Copy response to clipboard
- **🗑️ Clear Response** - Clear response area
- **⚡ Status Display** - Shows response code, status, and timing

## 🛠️ Development

### Dependencies
- [Fyne v2](https://fyne.io/) - Modern Go GUI framework
- [Resty](https://github.com/go-resty/resty) - HTTP client library
- [YAML v3](https://gopkg.in/yaml.v3) - YAML parsing

### Project Structure
```
kata-rest/
├── main.go        # Application entry point (main function only)
├── config.go      # Configuration structures and loading
├── app.go         # App struct and core application logic  
├── ui.go          # User interface creation and layout
├── http.go        # HTTP request handling and processing
├── config.yaml    # Request collections configuration
├── Makefile       # Build automation and development tools
├── go.mod         # Go modules
├── go.sum         # Dependencies
└── README.md      # This file
```

### Code Organization

The application is now organized into focused modules:

- **`main.go`** - Clean entry point with just the main function
- **`config.go`** - Configuration structures (`Request`, `Collection`, `Config`) and YAML loading
- **`app.go`** - Main `App` struct and core application logic (request selection, form handling)
- **`ui.go`** - All GUI creation functions (panels, layouts, widgets)
- **`http.go`** - HTTP client functionality (sending requests, parsing responses, headers/parameters)

### Benefits of Modular Structure

✅ **Better Organization** - Related functionality grouped together  
✅ **Easier Maintenance** - Changes isolated to specific modules  
✅ **Enhanced Readability** - Smaller, focused files are easier to understand  
✅ **Improved Testing** - Each module can be tested independently  
✅ **Team Collaboration** - Different developers can work on different modules  
✅ **Code Reusability** - Modules can be extended or reused more easily

## 🛠️ Development Tools

### Makefile Commands

The project includes a comprehensive Makefile with the following commands:

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build` | Build the application for current platform |
| `make build-all` | Build for multiple platforms (Linux, macOS, Windows) |
| `make run` | Run the application |
| `make dev` | Run with hot reload (requires 'air') |
| `make deps` | Install dependencies |
| `make deps-update` | Update dependencies to latest versions |
| `make test` | Run tests |
| `make test-coverage` | Run tests with coverage report |
| `make lint` | Lint the code (requires 'golangci-lint') |
| `make fmt` | Format the code |
| `make clean` | Clean build artifacts |
| `make install` | Install binary to system (/usr/local/bin) |
| `make uninstall` | Uninstall binary from system |
| `make release` | Create release packages for all platforms |
| `make install-tools` | Install development tools (air, golangci-lint) |
| `make check-tools` | Check if development tools are installed |
| `make info` | Show project information |

### Development Workflow

```bash
# 1. Install development tools
make install-tools

# 2. Check your setup
make check-tools

# 3. Install dependencies
make deps

# 4. Start development with hot reload
make dev

# 5. Build and test
make build
make test

# 6. Create release
make release
```

## 🎨 GUI Layout

```
┌─────────────────────────────────────────────────────────────┐
│ GoPostman - REST API Client                            │
├─────────────┬───────────────────────────────────────────────┤
│📚 Collections│ 🌐 Request                                    │
│             │ ┌─────────────────────────────────────────┐   │
│Collection:  │ │[GET] [URL.....................] [🚀Send]│   │
│[Dropdown]   │ └─────────────────────────────────────────┘   │
│             │                                              │
│Request:     │ 📋 Parameters                                │
│[Dropdown]   │ [key1=value1&key2=value2...............]     │
│             │                                              │
│🔄 Refresh   │ 🔖 Headers                                   │
│➕ New       │ [Header-Name: value.....................]     │
│             │                                              │
│             │ 📝 Body                          [🎨Format]  │
│             │ [Request body...........................]     │
│             ├──────────────────────────────────────────────│
│             │ 📨 Response                    ✅ 200 OK     │
│             │ [Response content......................]     │
│             │ [📋 Copy] [🗑️ Clear]                        │
└─────────────┴───────────────────────────────────────────────┘
```

## 🔧 Customization

- Edit `config.yaml` to add your own API collections
- Collections support multiple environments (dev, staging, prod)
- Headers and parameters are automatically URL-encoded
- JSON responses are automatically formatted for readability

## 📝 Example Requests

The included `config.yaml` contains sample requests for:
- **JSONPlaceholder API** - Public testing API
- **Local Development** - Common localhost endpoints  
- **Sample REST APIs** - Various public APIs for testing

Try them out to see the app in action! 🎉 