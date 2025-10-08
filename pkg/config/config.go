// Package config provides configuration management for the Redfish MCP server.
package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// AuthMethod represents Redfish authentication methods
type AuthMethod string

const (
	AuthMethodBasic   AuthMethod = "basic"
	AuthMethodSession AuthMethod = "session"
)

// MCPTransport represents MCP transport types
type MCPTransport string

const (
	MCPTransportStdio          MCPTransport = "stdio"
	MCPTransportSSE            MCPTransport = "sse"
	MCPTransportStreamableHTTP MCPTransport = "streamable-http"
)

// HostConfig represents configuration for a single Redfish host
type HostConfig struct {
	Address         string `json:"address"`
	Port            int    `json:"port,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	AuthMethod      string `json:"auth_method,omitempty"`
	TLSServerCACert string `json:"tls_server_ca_cert,omitempty"`
}

// Validate validates the host configuration
func (h *HostConfig) Validate() error {
	if h.Address == "" {
		return errors.New("host address cannot be empty")
	}

	if h.Port != 0 && (h.Port < 1 || h.Port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", h.Port)
	}

	if h.AuthMethod != "" && h.AuthMethod != string(AuthMethodBasic) && h.AuthMethod != string(AuthMethodSession) {
		return fmt.Errorf("invalid auth_method: %s. Must be one of: %s, %s", h.AuthMethod, AuthMethodBasic, AuthMethodSession)
	}

	return nil
}

// RedfishConfig represents complete Redfish configuration
type RedfishConfig struct {
	Hosts             []HostConfig `json:"hosts"`
	Port              int          `json:"port"`
	AuthMethod        string       `json:"auth_method"`
	Username          string       `json:"username"`
	Password          string       `json:"password"`
	TLSServerCACert   string       `json:"tls_server_ca_cert,omitempty"`
	DiscoveryEnabled  bool         `json:"discovery_enabled"`
	DiscoveryInterval int          `json:"discovery_interval"`
}

// Validate validates the Redfish configuration
func (r *RedfishConfig) Validate() error {
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", r.Port)
	}

	if r.AuthMethod != string(AuthMethodBasic) && r.AuthMethod != string(AuthMethodSession) {
		return fmt.Errorf("invalid auth_method: %s. Must be one of: %s, %s", r.AuthMethod, AuthMethodBasic, AuthMethodSession)
	}

	if r.DiscoveryInterval < 1 {
		return fmt.Errorf("discovery interval must be positive, got: %d", r.DiscoveryInterval)
	}

	for i, host := range r.Hosts {
		if err := host.Validate(); err != nil {
			return fmt.Errorf("invalid host configuration at index %d: %w", i, err)
		}
	}

	return nil
}

// MCPConfig represents MCP server configuration
type MCPConfig struct {
	Transport MCPTransport `json:"transport"`
	LogLevel  string       `json:"log_level"`
}

// Validate validates the MCP configuration
func (m *MCPConfig) Validate() error {
	validTransports := []MCPTransport{MCPTransportStdio, MCPTransportSSE, MCPTransportStreamableHTTP}
	if !slices.Contains(validTransports, m.Transport) {
		return fmt.Errorf("invalid transport: %s. Must be one of: %v", m.Transport, validTransports)
	}

	validLogLevels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}
	if !slices.Contains(validLogLevels, strings.ToUpper(m.LogLevel)) {
		return fmt.Errorf("invalid log_level: %s. Must be one of: %v", m.LogLevel, validLogLevels)
	}

	m.LogLevel = strings.ToUpper(m.LogLevel)
	return nil
}

// Config represents the complete application configuration
type Config struct {
	Redfish *RedfishConfig `json:"redfish"`
	MCP     *MCPConfig     `json:"mcp"`
}

// Validate validates the complete configuration
func (c *Config) Validate() error {
	if c.Redfish == nil {
		return errors.New("redfish configuration is required")
	}
	if err := c.Redfish.Validate(); err != nil {
		return fmt.Errorf("redfish config validation failed: %w", err)
	}

	if c.MCP == nil {
		return errors.New("mcp configuration is required")
	}
	if err := c.MCP.Validate(); err != nil {
		return fmt.Errorf("mcp config validation failed: %w", err)
	}

	return nil
}
