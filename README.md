# Redfish MCP Server (Go)

[![Go Version](https://img.shields.io/badge/go-1.25.2-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/theoriginalaiexplorer/mcp-redfish-go)](https://github.com/theoriginalaiexplorer/mcp-redfish-go/releases)

A high-performance MCP (Model Context Protocol) server implementation in Go for interacting with Redfish-enabled infrastructure. This server provides AI agents and applications with natural language access to Redfish API endpoints, enabling intelligent management and monitoring of data center infrastructure.

## Overview

The Redfish MCP Server enables AI-driven workflows to interact with Redfish-compliant servers through natural language queries. It provides a standardized interface for accessing hardware information, sensor data, and management capabilities of modern infrastructure components.

### Key Features

- 🚀 **High Performance**: Native Go implementation with compiled binaries and low memory footprint
- 🔍 **Redfish API Integration**: Full support for Redfish v1.x API specifications
- 🤖 **MCP Compatible**: Seamless integration with MCP clients and AI assistants
- 🔐 **Secure Authentication**: Support for Basic and Session-based authentication
- 🌐 **Multiple Transports**: stdio, SSE, and streamable-http transport options
- 🔄 **Concurrent Operations**: Goroutine-based concurrent discovery and requests
- 📊 **Structured Data Access**: JSON-based resource data retrieval with headers
- 🛠️ **Extensible Architecture**: Modular design for easy customization

## Quick Start

```bash
# Clone the repository
git clone https://github.com/theoriginalaiexplorer/mcp-redfish-go.git
cd mcp-redfish-go

# Build the server
go build -o bin/redfish-mcp ./cmd/redfish-mcp

# Configure environment
export REDFISH_HOSTS='[{"address": "192.168.1.100", "username": "admin", "password": "secret"}]'

# Run the server
./bin/redfish-mcp
```

## Installation

### Prerequisites

- Go 1.21 or later
- Access to Redfish-enabled infrastructure

### Building from Source

```bash
# Clone the repository
git clone https://github.com/theoriginalaiexplorer/mcp-redfish-go.git
cd mcp-redfish-go

# Initialize Go module
go mod tidy

# Build for your platform
go build -o bin/redfish-mcp ./cmd/redfish-mcp

# Or build for multiple platforms
make go-build-all
```

### Pre-built Binaries

Download the latest release binaries from the [Releases](https://github.com/theoriginalaiexplorer/mcp-redfish-go/releases) page for:
- Linux (amd64)
- macOS (amd64)
- Windows (amd64)

## MCP Tools

The server provides two primary MCP tools for interacting with Redfish infrastructure:

### `list_servers`
Lists all configured Redfish servers that can be accessed.

**Example usage:**
```
List all available Redfish servers
```

**Response:**
```json
{
  "servers": ["192.168.1.100", "192.168.1.101"]
}
```

### `get_resource_data`
Fetches data from a specific Redfish resource endpoint.

**Parameters:**
- `url`: The Redfish resource URL (e.g., `https://192.168.1.100/redfish/v1/Systems/1`)

**Example usage:**
```
Get the system information from https://192.168.1.100/redfish/v1/Systems/1
```

**Response:**
```json
{
  "headers": {
    "content-type": ["application/json"],
    "etag": ["\"12345\""]
  },
  "data": {
    "@odata.id": "/redfish/v1/Systems/1",
    "Name": "System1",
    "Manufacturer": "Example Corp"
  }
}
```

## Configuration

The server supports configuration through environment variables or JSON files. Configuration is validated at startup to ensure proper setup.

### Configuration Methods

1. **JSON File** (Recommended for complex configurations):
   ```bash
   ./bin/redfish-mcp --config config.json
   ```

2. **Environment Variables** (Default):
   ```bash
   export REDFISH_HOSTS='[{"address": "192.168.1.100"}]'
   ./bin/redfish-mcp
   ```

### JSON Configuration File

Create a `config.json` file with the complete configuration:

```json
{
  "redfish": {
    "hosts": [
      {
        "address": "192.168.1.100",
        "port": 443,
        "username": "admin",
        "password": "secret123",
        "auth_method": "session",
        "tls_server_ca_cert": "/path/to/ca-cert.pem"
      }
    ],
    "port": 443,
    "auth_method": "session",
    "username": "default_user",
    "password": "default_pass",
    "tls_server_ca_cert": "",
    "insecure_skip_verify": false,
    "discovery_enabled": false,
    "discovery_interval": 30
  },
  "mcp": {
    "transport": "stdio",
    "log_level": "INFO"
  }
}
```

### Accessing Servers Without SSL Certificates

For development or testing with Redfish servers that have self-signed or invalid SSL certificates, you can skip certificate verification:

**Environment Variable:**
```bash
export REDFISH_INSECURE_SKIP_VERIFY=true
```

**JSON Configuration:**
```json
{
  "redfish": {
    "insecure_skip_verify": true
  }
}
```

⚠️ **Security Warning:** Only use `insecure_skip_verify` in development or trusted environments. This option disables SSL certificate verification, making connections vulnerable to man-in-the-middle attacks.

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDFISH_CONFIG_FILE` | Path to JSON config file | - | No |
| `REDFISH_HOSTS` | JSON array of host configs | `[{"address":"127.0.0.1"}]` | Yes* |
| `REDFISH_PORT` | Default Redfish port | `443` | No |
| `REDFISH_AUTH_METHOD` | Auth method: `basic` or `session` | `session` | No |
| `REDFISH_USERNAME` | Default username | `""` | No |
| `REDFISH_PASSWORD` | Default password | `""` | No |
| `REDFISH_SERVER_CA_CERT` | CA certificate path | `""` | No |
| `REDFISH_INSECURE_SKIP_VERIFY` | Skip SSL certificate verification | `false` | No |
| `REDFISH_DISCOVERY_ENABLED` | Enable SSDP discovery | `false` | No |
| `REDFISH_DISCOVERY_INTERVAL` | Discovery interval (seconds) | `30` | No |
| `MCP_TRANSPORT` | Transport: `stdio`, `sse`, `streamable-http` | `stdio` | No |
| `MCP_REDFISH_LOG_LEVEL` | Log level: `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL` | `INFO` | No |

*Required when not using JSON config file

### Host Configuration

Each host in `REDFISH_HOSTS` or the JSON file can specify:

- `address` (required): IP or hostname
- `port` (optional): Port number
- `username` (optional): Host-specific username
- `password` (optional): Host-specific password
- `auth_method` (optional): `basic` or `session`
- `tls_server_ca_cert` (optional): Custom CA certificate path

### Validation

The server validates configuration on startup and exits with detailed error messages if invalid.

## Running the Server

### Basic Usage

```bash
# Build the server
go build -o bin/redfish-mcp ./cmd/redfish-mcp

# Set configuration
export REDFISH_HOSTS='[{"address": "192.168.1.100", "username": "admin", "password": "secret"}]'

# Run with stdio transport (default)
./bin/redfish-mcp

# Or use a config file
./bin/redfish-mcp --config redfish-config.json
```

### Transports

The server supports multiple MCP transport mechanisms:

#### stdio Transport (Default)
Standard input/output communication for direct integration with MCP clients.

```bash
export MCP_TRANSPORT="stdio"
./bin/redfish-mcp
```

#### SSE Transport (Server-Sent Events)
Network-based communication over HTTP for remote clients.

```bash
export MCP_TRANSPORT="sse"
./bin/redfish-mcp
# Server will be available at http://localhost:8000/sse
```

#### streamable-http Transport
Alternative HTTP-based transport for specific MCP implementations.

```bash
export MCP_TRANSPORT="streamable-http"
./bin/redfish-mcp
```

### Makefile Targets

Use the provided Makefile for common operations:

```bash
make go-build      # Build the binary
make go-run        # Run the server
make go-test       # Run tests
make go-fmt        # Format code
make go-vet        # Lint code
```

## Integration with MCP Clients

### Claude Desktop

Add the server to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "redfish": {
      "command": "/path/to/bin/redfish-mcp",
      "env": {
        "REDFISH_HOSTS": "[{\"address\": \"192.168.1.100\", \"username\": \"admin\", \"password\": \"secret123\"}]",
        "MCP_TRANSPORT": "stdio"
      }
    }
  }
}
```

### VS Code with GitHub Copilot

Configure in your VS Code settings or `.vscode/mcp.json`:

```json
{
  "mcp": {
    "servers": {
      "redfish": {
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

### MCP Inspector

Test the server with the MCP Inspector:

```bash
npx @modelcontextprotocol/inspector /path/to/bin/redfish-mcp
```

For more information, see the [VS Code documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers).


## Testing

### Unit Tests

Run the Go test suite:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...

# Using Makefile
make go-test
```

### MCP Inspector

Test the server interactively with the MCP Inspector:

```bash
# Build and test
go build -o bin/redfish-mcp ./cmd/redfish-mcp
npx @modelcontextprotocol/inspector ./bin/redfish-mcp
```

### Integration Testing

The server can be tested against real Redfish hardware or emulators. Ensure your `REDFISH_HOSTS` configuration points to accessible endpoints.

## Example Use Cases
- **AI Assistants**: Enable LLMs to fetch infrastructure data via Redfish API.
- **Chatbots & Virtual Agents**: Retrieve data, and personalize responses.

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/theoriginalaiexplorer/mcp-redfish-go.git
cd mcp-redfish-go

# Install dependencies
go mod tidy

# Build
go build ./cmd/redfish-mcp
```

### Development Workflow

```bash
# Format code
go fmt ./...

# Lint code
go vet ./...

# Run tests
go test ./...

# Build for development
go build -o bin/redfish-mcp ./cmd/redfish-mcp

# Run the server
./bin/redfish-mcp
```

### Project Structure

```
.
├── cmd/redfish-mcp/          # Main application
│   └── main.go              # Entry point
├── pkg/                     # Go packages
│   ├── config/              # Configuration management
│   │   ├── config.go        # Config structs and validation
│   │   ├── env.go           # Environment parsing
│   │   └── config_test.go   # Unit tests
│   ├── redfish/             # Redfish client and discovery
│   │   ├── client.go        # HTTP client with retry logic
│   │   ├── discovery.go     # SSDP discovery
│   │   └── types.go         # Type definitions
│   ├── mcp/                 # MCP server implementation
│   │   └── server.go        # MCP server setup and tools
│   └── common/              # Shared utilities
│       └── hosts.go         # Host management
├── .github/workflows/       # CI/CD workflows
│   └── release.yml          # Release automation
├── Makefile                 # Build and development tasks
└── README.md               # This file
```

### Code Quality

- **Formatting**: `go fmt` for consistent formatting
- **Linting**: `go vet` for static analysis
- **Testing**: Comprehensive unit tests with `go test`
- **Modules**: Go modules for dependency management

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `go test ./...` and `go vet ./...`
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) for the protocol specification
- [DMTF Redfish](https://www.dmtf.org/standards/redfish) for the API standard
- Go community for the excellent tooling and libraries
