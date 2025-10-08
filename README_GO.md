# Redfish MCP Server - Go Implementation

This is the Go implementation of the Redfish MCP Server, providing the same functionality as the Python version with improved performance and concurrency.

## Features

- **Native Performance**: Compiled Go binary with better memory usage and CPU performance
- **Concurrent Discovery**: Goroutine-based SSDP discovery for better scalability
- **Type Safety**: Strong typing with compile-time guarantees
- **Simplified Deployment**: Single binary deployment with no runtime dependencies
- **Full MCP Compatibility**: Implements the same MCP tools and protocols as the Python version

## Quick Start

### Prerequisites

- Go 1.21 or later
- Access to Redfish-enabled infrastructure

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd mcp-redfish

# Build the Go binary
make go-build
# or directly: go build -o bin/redfish-mcp ./cmd/redfish-mcp
```

### Configuration

The Go implementation uses the same environment variables as the Python version:

```bash
export REDFISH_HOSTS='[{"address": "192.168.1.100", "username": "admin", "password": "secret123"}]'
export MCP_TRANSPORT="stdio"
export MCP_REDFISH_LOG_LEVEL="INFO"
```

### Running

```bash
# Run the server
./bin/redfish-mcp

# Or use the Makefile
make go-run
```

## MCP Tools

The Go implementation provides the same MCP tools as the Python version:

### `list_servers`
Lists all configured Redfish servers.

**Input**: None
**Output**: Array of server addresses

### `get_resource_data`
Fetches data from a specific Redfish resource.

**Input**:
```json
{
  "url": "https://server.example.com/redfish/v1/Systems/1"
}
```

**Output**:
```json
{
  "headers": {...},
  "data": {...}
}
```

## Development

### Building

```bash
# Build for current platform
make go-build

# Build for all platforms
make go-build-all

# Cross-platform builds
make go-build-linux    # Linux amd64
make go-build-darwin   # macOS amd64
make go-build-windows  # Windows amd64
```

### Testing

```bash
# Run all tests
make go-test

# Run tests with verbose output
make go-test-verbose

# Run tests with race detection
make go-test-race
```

### Code Quality

```bash
# Format code
make go-fmt

# Run vet (static analysis)
make go-vet

# Tidy modules
make go-mod-tidy

# Development setup (all of the above)
make go-dev-setup
```

## Architecture

### Project Structure

```
cmd/redfish-mcp/
├── main.go                 # Entry point
pkg/
├── config/
│   ├── config.go           # Configuration structs and validation
│   └── env.go              # Environment variable parsing
├── redfish/
│   ├── client.go           # Redfish HTTP client with retry logic
│   ├── discovery.go        # SSDP discovery implementation
│   └── types.go            # Redfish-specific types
├── mcp/
│   ├── server.go           # MCP server setup and tool handlers
└── common/
    ├── hosts.go            # Host management
```

### Key Components

1. **Configuration Management** (`pkg/config/`): Environment variable parsing and validation
2. **Redfish Client** (`pkg/redfish/client.go`): HTTP client with authentication and retry logic
3. **SSDP Discovery** (`pkg/redfish/discovery.go`): Network discovery of Redfish endpoints
4. **MCP Server** (`pkg/mcp/server.go`): MCP protocol implementation and tool registration
5. **Host Management** (`pkg/common/hosts.go`): Static and discovered host management

### Differences from Python Implementation

- **Concurrency**: Uses goroutines instead of threads for background tasks
- **Error Handling**: Go's explicit error handling instead of exceptions
- **Memory Management**: Go's garbage collection and value semantics
- **HTTP Client**: Standard library `net/http` with proper TLS configuration
- **Retry Logic**: Uses `github.com/avast/retry-go` instead of `tenacity`
- **Type Safety**: Compile-time type checking for all data structures

## Integration

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "redfish-go": {
      "command": "/path/to/bin/redfish-mcp",
      "env": {
        "REDFISH_HOSTS": "[{\"address\": \"192.168.1.100\", \"username\": \"admin\", \"password\": \"secret123\"}]",
        "MCP_TRANSPORT": "stdio"
      }
    }
  }
}
```

### VS Code

Add to your `.vscode/mcp.json` or `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "redfish-go": {
        "type": "stdio",
        "command": "/path/to/bin/redfish-mcp",
        "env": {
          "REDFISH_HOSTS": "[{\"address\": \"192.168.1.100\", \"username\": \"admin\", \"password\": \"secret123\"}]",
          "MCP_TRANSPORT": "stdio"
        }
      }
    }
  }
}
```

## Performance

The Go implementation offers several performance advantages over Python:

- **Startup Time**: Faster cold start due to compiled binary
- **Memory Usage**: Lower memory footprint with efficient garbage collection
- **Concurrency**: Better handling of concurrent requests with goroutines
- **CPU Usage**: More efficient CPU utilization for I/O operations

## Compatibility

The Go implementation is fully compatible with the Python version:

- Same environment variable configuration
- Same MCP tool interfaces
- Same Redfish API interactions
- Same authentication methods
- Same error handling semantics

You can switch between implementations without changing your configuration or integration setup.