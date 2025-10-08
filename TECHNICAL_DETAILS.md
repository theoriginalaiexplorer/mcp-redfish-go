# Redfish MCP Server - Technical Documentation

## Project Overview

The Redfish MCP Server is a Python-based Model Context Protocol (MCP) server that provides a natural language interface for managing infrastructure via Redfish APIs. It enables AI agents to interact with Redfish-enabled infrastructure components through conversational queries.

### Key Features
- **Natural Language Queries**: AI agents can query infrastructure data using natural language
- **MCP Integration**: Seamless integration with MCP clients
- **Full Redfish Support**: Wraps the Python Redfish library
- **Multi-transport Support**: stdio, SSE, and streamable-http transports
- **SSDP Discovery**: Automatic endpoint discovery capabilities
- **Comprehensive Configuration**: Environment-based configuration with validation

## Architecture

### Core Components

#### 1. Main Entry Point (`src/main.py`)
- **Purpose**: Server initialization and lifecycle management
- **Responsibilities**:
  - Logging configuration
  - SSDP discovery thread management
  - MCP server startup with configured transport
- **Key Classes**:
  - `RedfishMCPServer`: Main server class handling discovery and MCP initialization

#### 2. MCP Server (`src/common/server.py`)
- **Purpose**: FastMCP server initialization
- **Implementation**: Uses FastMCP framework for MCP protocol handling
- **Error Handling**: Comprehensive error handling during server initialization

#### 3. Configuration Management (`src/common/config.py`)
- **Purpose**: Environment variable loading and validation
- **Features**:
  - JSON-based host configuration
  - Environment variable validation
  - Legacy configuration fallback
  - Type-safe configuration objects
- **Validation**: Comprehensive validation with deprecation warnings

#### 4. MCP Tools (`src/tools/`)
- **Tool Registration**: Automatic registration via module imports
- **Available Tools**:
  - `list_servers()`: Lists accessible Redfish servers
  - `get_resource_data()`: Fetches data from specific Redfish resources

### Data Flow
1. **Configuration Loading**: Environment variables → Validation → Configuration objects
2. **Server Initialization**: FastMCP server setup with transport configuration
3. **Tool Registration**: Import-time registration of MCP tools
4. **Request Processing**: MCP client requests → Tool execution → Redfish API calls → Response formatting

## Technology Stack

### Core Dependencies
- **Python**: 3.13+ (primary development and runtime version)
- **FastMCP**: MCP protocol implementation framework
- **Python Redfish Library**: Redfish API client library
- **Tenacity**: Retry logic for resilient API calls
- **python-dotenv**: Environment variable management

### Development Dependencies
- **pytest**: Testing framework with async support
- **mypy**: Static type checking
- **ruff**: Linting and code formatting
- **pre-commit**: Git hooks for code quality
- **bandit**: Security scanning
- **coverage**: Code coverage reporting

### Build and Packaging
- **uv**: Fast Python package manager
- **Hatchling**: Build backend for Python packages
- **Docker**: Containerization support

## Configuration Management

### Environment Variables

#### Redfish Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDFISH_HOSTS` | JSON array of Redfish endpoints | `[{"address":"127.0.0.1"}]` | Yes |
| `REDFISH_PORT` | Default port for Redfish API | `443` | No |
| `REDFISH_AUTH_METHOD` | Authentication method (`basic`/`session`) | `session` | No |
| `REDFISH_USERNAME` | Default username | `""` | No |
| `REDFISH_PASSWORD` | Default password | `""` | No |
| `REDFISH_SERVER_CA_CERT` | Path to CA certificate | `None` | No |
| `REDFISH_DISCOVERY_ENABLED` | Enable SSDP discovery | `false` | No |
| `REDFISH_DISCOVERY_INTERVAL` | Discovery interval (seconds) | `30` | No |

#### MCP Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MCP_TRANSPORT` | Transport method (`stdio`/`sse`/`streamable-http`) | `stdio` | No |
| `MCP_REDFISH_LOG_LEVEL` | Logging level | `INFO` | No |

### Host Configuration Format
```json
[
  {
    "address": "192.168.1.100",
    "port": 443,
    "username": "admin",
    "password": "password123",
    "auth_method": "session",
    "tls_server_ca_cert": "/path/to/ca-cert.pem"
  }
]
```

### Configuration Validation
- **JSON Syntax**: Validates `REDFISH_HOSTS` JSON structure
- **Required Fields**: Ensures each host has an `address` field
- **Port Ranges**: Validates ports are between 1-65535
- **Authentication Methods**: Accepts only `basic` or `session`
- **Transport Types**: Accepts only `stdio`, `sse`, or `streamable-http`
- **Log Levels**: Standard Python logging levels

## MCP Tools

### 1. `list_servers()`
**Purpose**: Lists all configured Redfish servers
**Input**: None
**Output**: List of server addresses
**Implementation**: Reads from validated configuration

### 2. `get_resource_data(url: str)`
**Purpose**: Fetches data from a specific Redfish resource
**Input**: Full Redfish URL (e.g., `https://server/redfish/v1/Systems/1`)
**Output**: Dictionary with `headers` and `data` fields
**Implementation**:
- URL parsing and validation
- Server configuration lookup
- Redfish client initialization
- API call with error handling
- Response formatting

## Testing Strategy

### Test Structure
```
test/
├── conftest.py              # Test configuration and fixtures
├── common/                  # Unit tests for common modules
├── tools/                   # Tool-specific tests
├── integration/             # Integration tests
├── complex/                 # Complex scenario tests
├── property/                # Property-based tests
└── e2e/                     # End-to-end tests
```

### Test Types
- **Unit Tests**: Individual function/component testing
- **Integration Tests**: Component interaction testing
- **End-to-End Tests**: Full system testing with Redfish emulator
- **Property-Based Tests**: Hypothesis-driven testing

### Test Configuration
- **pytest**: Main testing framework
- **pytest-asyncio**: Async test support
- **pytest-cov**: Coverage reporting
- **pytest-xdist**: Parallel test execution
- **hypothesis**: Property-based testing

### Coverage Requirements
- **Minimum Coverage**: 70%
- **Coverage Tools**: pytest-cov with XML reporting
- **CI Integration**: Codecov integration

### End-to-End Testing
- **Redfish Interface Emulator**: DMTF emulator for realistic testing
- **SSL/TLS Support**: HTTPS testing with self-signed certificates
- **MCP Inspector Integration**: Visual debugging and testing
- **Docker-based**: Containerized test environment

## CI/CD Pipeline

### GitHub Actions Workflows

#### 1. CI/CD Pipeline (`ci-cd.yml`)
**Triggers**: Pull requests to main branch
**Jobs**:
- **Test Suite**: Unit tests with coverage
- **Code Quality**: Linting, formatting, type checking
- **Security Scan**: Bandit security analysis
- **Container Build**: Docker image building and testing
- **E2E Tests**: Full integration testing with emulator

#### 2. CodeQL Analysis (`codeql.yml`)
**Purpose**: Automated security vulnerability detection
**Triggers**: Push to main, scheduled weekly
**Languages**: Python security analysis

#### 3. Dependency Updates (`dependency-updates.yml`)
**Purpose**: Automated dependency management
**Triggers**: Weekly schedule
**Behavior**: Creates PRs with updated `uv.lock` file

### Quality Gates
- **Tests**: Must pass all unit tests
- **Coverage**: Must maintain 70% minimum coverage
- **Quality**: Must pass linting, formatting, and type checking
- **Security**: Security scans must complete (warnings allowed)
- **Container**: Docker build must succeed

## Containerization

### Docker Support
- **Base Image**: `python:3.13-slim`
- **Package Manager**: uv for fast dependency installation
- **Virtual Environment**: Isolated Python environment
- **Multi-platform**: Supports both Docker and Podman

### Container Runtimes
- **Auto-detection**: Automatically detects available runtime
- **Docker**: Optimized with BuildKit cache mounts
- **Podman**: Compatible fallback without cache mounts
- **Manual Override**: `CONTAINER_RUNTIME` environment variable

### Build Optimization
- **Layer Caching**: Efficient Docker layer caching
- **Dependency Locking**: Uses `uv.lock` for reproducible builds
- **Minimal Image**: Slim base image for smaller footprint

## Security Considerations

### Authentication
- **Session-based Auth**: Preferred authentication method
- **Basic Auth**: Supported for legacy systems
- **Per-host Credentials**: Individual credentials per endpoint

### TLS/SSL
- **Certificate Validation**: Optional CA certificate configuration
- **Self-signed Support**: E2E testing with self-signed certificates
- **Secure Defaults**: HTTPS default for Redfish communication

### Security Scanning
- **Bandit**: Automated security vulnerability scanning
- **CodeQL**: GitHub's semantic code analysis
- **Dependency Updates**: Automated security patch management

### Environment Security
- **No Secrets in Code**: Environment variables for sensitive data
- **Secure Defaults**: Conservative default configurations
- **Validation**: Input validation and sanitization

## Development Workflow

### Code Quality Tools
- **ruff**: Fast Python linter and formatter
- **mypy**: Static type checker with strict settings
- **pre-commit**: Git hooks for quality enforcement
- **bandit**: Security vulnerability scanner

### Code Standards
- **Line Length**: 88 characters (Black-compatible)
- **Quote Style**: Double quotes
- **Import Sorting**: ruff-based import organization
- **Type Hints**: Comprehensive type annotations

### Development Commands
```bash
# Setup
make install-dev    # Install development dependencies
make pre-commit-install  # Setup pre-commit hooks

# Code Quality
make lint          # Run ruff linting
make format        # Format code with ruff
make type-check    # Run mypy type checking
make security      # Run bandit security scan

# Testing
make test          # Run unit tests
make test-cov      # Run tests with coverage
make e2e           # Run end-to-end tests

# Server Development
make run-stdio    # Run with stdio transport
make inspect      # Run with MCP Inspector
```

### Project Structure
```
src/
├── main.py              # Entry point and server lifecycle
├── common/              # Shared utilities and configuration
│   ├── client.py        # Redfish API client
│   ├── config.py        # Configuration management
│   ├── discovery.py     # SSDP discovery
│   ├── hosts.py         # Host management
│   ├── server.py        # MCP server initialization
│   └── validation.py    # Configuration validation
└── tools/               # MCP tool implementations
    ├── get.py           # Resource data fetching
    └── servers.py       # Server listing

test/
├── common/              # Unit tests for common modules
├── tools/               # Tool-specific tests
├── integration/         # Integration tests
└── e2e/                 # End-to-end tests
```

## Performance Considerations

### Retry Logic
- **Tenacity Integration**: Configurable retry behavior
- **Exponential Backoff**: Intelligent retry delays
- **Jitter**: Randomized delays to prevent thundering herd
- **Fast Test Configuration**: Reduced delays for testing

### Resource Management
- **Client Lifecycle**: Proper login/logout handling
- **Connection Pooling**: Efficient Redfish client usage
- **Error Handling**: Comprehensive error recovery

### Scalability
- **Async Support**: Asynchronous tool implementations
- **Threading**: Background SSDP discovery
- **Memory Management**: Efficient data structures

## Deployment Options

### Transport Methods
1. **stdio**: Standard input/output (default, most compatible)
2. **SSE**: Server-Sent Events (network-based)
3. **streamable-http**: HTTP streaming (advanced clients)

### Integration Points
- **Claude Desktop**: Native MCP support
- **VS Code**: GitHub Copilot integration
- **Custom Clients**: Any MCP-compatible client

### Production Considerations
- **Environment Variables**: Secure configuration management
- **Logging**: Configurable log levels and formats
- **Monitoring**: Integration with monitoring systems
- **Scaling**: Multiple server instances support

## Future Enhancements

### Planned Features
- **Enhanced Discovery**: Improved SSDP and service discovery
- **Caching**: Response caching for performance
- **Metrics**: Prometheus metrics integration
- **Web UI**: Administrative web interface

### Architecture Improvements
- **Plugin System**: Extensible tool architecture
- **Configuration UI**: Web-based configuration management
- **Multi-tenancy**: Support for multiple isolated environments

This technical documentation provides a comprehensive overview of the Redfish MCP Server's architecture, implementation details, and operational aspects. The modular design and comprehensive testing strategy ensure maintainability and reliability for production deployments.