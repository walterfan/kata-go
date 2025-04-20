# Encoding Tool

## Quick Start

```shell

mkdir encoding_tool
cd encoding_tool
go mod init encoding_tool
go install github.com/spf13/cobra-cli@latest
go get github.com/google/uuid
cobra-cli init
cobra-cli add convert
cobra-cli add generate
```

## Build

```shell
go build -o encoding_tool main.go
```

## Usage

```shell

After adding these commands, you can use them like this:

### Convert commands:
```bash
# Base64 encoding
./encoding_tool convert base64encode "hello world"
# Output: aGVsbG8gd29ybGQ=

# Base64 decoding
./encoding_tool convert base64decode "aGVsbG8gd29ybGQ="
# Output: hello world

# URL encoding
./encoding_tool convert urlencode "hello world & more"
# Output: hello+world+%26+more

# URL decoding
./encoding_tool convert urldecode "hello+world+%26+more"
# Output: hello world & more
```

### Generate commands:
```bash
# Generate UUID
./encoding_tool generate uuid
# Output: 123e4567-e89b-12d3-a456-426614174000

# Generate random string (default 16 chars, letters only)
./encoding_tool generate random

# Generate random string with numbers (20 chars)
./encoding_tool generate random 20 -n

# Generate random string with numbers and symbols (32 chars)
./encoding_tool generate random 32 -n -s
```

## Project Structure Update

```
encoding_tool/
├── cmd/
│   ├── root.go
│   ├── convert.go
│   └── generate.go
├── main.go
├── go.mod
└── go.sum
```

