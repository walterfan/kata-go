# Encoding Tool

## Quick Start

```shell

mkdir github.com/walterfan/confucius
cd github.com/walterfan/confucius
go mod init github.com/walterfan/confucius
go install github.com/spf13/cobra-cli@latest
go get github.com/google/uuid
cobra-cli init
cobra-cli add convert
cobra-cli add generate
```

## Build

```shell
go build -o github.com/walterfan/confucius main.go
```

## Usage

```shell

After adding these commands, you can use them like this:

### Convert commands:
```bash
# Base64 encoding
./github.com/walterfan/confucius convert base64encode "hello world"
# Output: aGVsbG8gd29ybGQ=

# Base64 decoding
./github.com/walterfan/confucius convert base64decode "aGVsbG8gd29ybGQ="
# Output: hello world

# URL encoding
./github.com/walterfan/confucius convert urlencode "hello world & more"
# Output: hello+world+%26+more

# URL decoding
./github.com/walterfan/confucius convert urldecode "hello+world+%26+more"
# Output: hello world & more
```

### Generate commands:
```bash
# Generate UUID
./github.com/walterfan/confucius generate uuid
# Output: 123e4567-e89b-12d3-a456-426614174000

# Generate random string (default 16 chars, letters only)
./github.com/walterfan/confucius generate random

# Generate random string with numbers (20 chars)
./github.com/walterfan/confucius generate random 20 -n

# Generate random string with numbers and symbols (32 chars)
./github.com/walterfan/confucius generate random 32 -n -s
```

## Project Structure Update

```
github.com/walterfan/confucius/
├── cmd/
│   ├── root.go
│   ├── convert.go
│   └── generate.go
├── main.go
├── go.mod
└── go.sum
```

