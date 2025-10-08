package mcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/nokia/mcp-redfish-go/pkg/common"
	"github.com/nokia/mcp-redfish-go/pkg/config"
	"github.com/nokia/mcp-redfish-go/pkg/redfish"
)

// Server wraps the MCP server with Redfish-specific functionality
type Server struct {
	mcpServer   *mcp.Server
	config      *config.Config
	hostManager *common.HostManager
	logger      *slog.Logger
}

// NewServer creates a new Redfish MCP server
func NewServer(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	if logger == nil {
		logger = slog.Default()
	}

	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "redfish-mcp",
			Version: "0.1.0",
		},
		&mcp.ServerOptions{
			// Configure based on transport
		},
	)

	// Create host manager
	hostManager := common.NewHostManager(logger)

	server := &Server{
		mcpServer:   mcpServer,
		config:      cfg,
		hostManager: hostManager,
		logger:      logger,
	}

	// Register tools
	if err := server.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	logger.Info("Redfish MCP server created successfully")
	return server, nil
}

// GetResourceInput represents input for the get_resource_data tool
type GetResourceInput struct {
	URL string `json:"url" jsonschema:"Redfish resource URL"`
}

// GetResourceOutput represents output for the get_resource_data tool
type GetResourceOutput struct {
	Headers map[string][]string `json:"headers"`
	Data    interface{}         `json:"data"`
}

// registerTools registers the MCP tools
func (s *Server) registerTools() error {
	// Register list_servers tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "list_servers",
		Description: "List all Redfish servers that can be accessed",
	}, s.handleListServers)

	// Register get_resource_data tool
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_resource_data",
		Description: "Fetch data from a specific Redfish resource",
	}, s.handleGetResourceData)

	s.logger.Info("MCP tools registered successfully")
	return nil
}

// ListServersOutput represents the output for the list_servers tool
type ListServersOutput struct {
	Servers []string `json:"servers"`
}

// handleListServers handles the list_servers tool
func (s *Server) handleListServers(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, ListServersOutput, error) {
	s.logger.Info("Handling list_servers request")

	addresses := s.hostManager.GetAddresses()

	return nil, ListServersOutput{Servers: addresses}, nil
}

// handleGetResourceData handles the get_resource_data tool
func (s *Server) handleGetResourceData(ctx context.Context, req *mcp.CallToolRequest, input GetResourceInput) (*mcp.CallToolResult, GetResourceOutput, error) {
	s.logger.Info("Handling get_resource_data request")

	// Parse the URL to extract server address and resource path
	serverAddr, resourcePath, err := s.parseRedfishURL(input.URL)
	if err != nil {
		return nil, GetResourceOutput{}, fmt.Errorf("invalid Redfish URL: %w", err)
	}

	// Find the server configuration
	hostConfig, found := s.hostManager.GetHostByAddress(serverAddr)
	if !found {
		return nil, GetResourceOutput{}, fmt.Errorf("server %s not found in configuration", serverAddr)
	}

	// Create Redfish client
	clientConfig := s.createClientConfig(hostConfig)
	client := redfish.NewClient(clientConfig, s.logger)

	// Login and fetch data
	if err := client.Login(); err != nil {
		return nil, GetResourceOutput{}, fmt.Errorf("failed to login to Redfish server: %w", err)
	}
	defer client.Close()

	// Get resource data with headers
	response, err := client.GetWithHeaders(resourcePath)
	if err != nil {
		return nil, GetResourceOutput{}, fmt.Errorf("failed to get resource data: %w", err)
	}

	return nil, GetResourceOutput{
		Headers: response.Headers,
		Data:    response.Data,
	}, nil
}

// parseRedfishURL parses a Redfish URL to extract server address and resource path
func (s *Server) parseRedfishURL(url string) (string, string, error) {
	// This is a simplified parser - in production, use proper URL parsing
	// Expected format: https://server:port/redfish/v1/resource/path

	if len(url) < 8 || url[:8] != "https://" {
		return "", "", fmt.Errorf("URL must use HTTPS")
	}

	// Remove https:// prefix
	withoutScheme := url[8:]

	// Find the first / after the host
	hostEnd := -1
	for i, char := range withoutScheme {
		if char == '/' {
			hostEnd = i
			break
		}
	}

	if hostEnd == -1 {
		return "", "", fmt.Errorf("invalid URL format")
	}

	serverAddr := withoutScheme[:hostEnd]
	resourcePath := withoutScheme[hostEnd:]

	// Basic validation
	if serverAddr == "" {
		return "", "", fmt.Errorf("empty server address")
	}

	if resourcePath == "" {
		resourcePath = "/"
	}

	return serverAddr, resourcePath, nil
}

// createClientConfig creates a Redfish client config from host config
func (s *Server) createClientConfig(hostConfig config.HostConfig) *redfish.ClientConfig {
	config := redfish.DefaultClientConfig()

	config.Address = hostConfig.Address
	if hostConfig.Port != 0 {
		config.Port = hostConfig.Port
	} else {
		config.Port = s.config.Redfish.Port
	}

	config.Username = hostConfig.Username
	if config.Username == "" {
		config.Username = s.config.Redfish.Username
	}

	config.Password = hostConfig.Password
	if config.Password == "" {
		config.Password = s.config.Redfish.Password
	}

	config.AuthMethod = redfish.AuthMethod(hostConfig.AuthMethod)
	if config.AuthMethod == "" {
		config.AuthMethod = redfish.AuthMethod(s.config.Redfish.AuthMethod)
	}

	config.TLSServerCACert = hostConfig.TLSServerCACert
	if config.TLSServerCACert == "" {
		config.TLSServerCACert = s.config.Redfish.TLSServerCACert
	}

	return config
}

// Start starts the MCP server with the specified transport
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting Redfish MCP server",
		"transport", s.config.MCP.Transport)

	// For now, we'll implement stdio transport
	// Other transports can be added later
	switch s.config.MCP.Transport {
	case config.MCPTransportStdio:
		return s.startStdio(ctx)
	default:
		return fmt.Errorf("unsupported transport: %s", s.config.MCP.Transport)
	}
}

// startStdio starts the server with stdio transport
func (s *Server) startStdio(ctx context.Context) error {
	transport := &mcp.StdioTransport{}
	return s.mcpServer.Run(ctx, transport)
}

// GetMCPServer returns the underlying MCP server
func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcpServer
}
