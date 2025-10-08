package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// ConfigError represents configuration validation errors
type ConfigError struct {
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	redfishConfig, err := loadRedfishConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load redfish config: %w", err)
	}

	mcpConfig, err := loadMCPConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mcp config: %w", err)
	}

	config := &Config{
		Redfish: redfishConfig,
		MCP:     mcpConfig,
	}

	if err := config.Validate(); err != nil {
		return nil, &ConfigError{Message: "configuration validation failed", Cause: err}
	}

	return config, nil
}

func loadRedfishConfig() (*RedfishConfig, error) {
	// Check if config file is specified
	if configFile := os.Getenv("REDFISH_CONFIG_FILE"); configFile != "" {
		// Read config from JSON file
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, &ConfigError{
				Message: fmt.Sprintf("failed to read config file %s", configFile),
				Cause:   err,
			}
		}

		var config RedfishConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, &ConfigError{
				Message: fmt.Sprintf("invalid JSON in config file %s", configFile),
				Cause:   err,
			}
		}

		// Validate the config
		if err := config.Validate(); err != nil {
			return nil, &ConfigError{
				Message: "config validation failed",
				Cause:   err,
			}
		}

		return &config, nil
	}

	// Fallback to environment variables
	// Parse hosts from JSON
	hostsJSON := getEnv("REDFISH_HOSTS", `[{"address": "127.0.0.1"}]`)
	var hosts []HostConfig
	if err := json.Unmarshal([]byte(hostsJSON), &hosts); err != nil {
		return nil, &ConfigError{
			Message: fmt.Sprintf("invalid JSON in REDFISH_HOSTS: %s", hostsJSON),
			Cause:   err,
		}
	}

	// Validate each host
	for i, host := range hosts {
		if err := host.Validate(); err != nil {
			return nil, &ConfigError{
				Message: fmt.Sprintf("invalid host configuration at index %d", i),
				Cause:   err,
			}
		}
	}

	port, err := getEnvInt("REDFISH_PORT", 443, 1, 65535)
	if err != nil {
		return nil, err
	}

	discoveryInterval, err := getEnvInt("REDFISH_DISCOVERY_INTERVAL", 30, 1, 3600)
	if err != nil {
		return nil, err
	}

	config := &RedfishConfig{
		Hosts:              hosts,
		Port:               port,
		AuthMethod:         getEnv("REDFISH_AUTH_METHOD", string(AuthMethodSession)),
		Username:           getEnv("REDFISH_USERNAME", ""),
		Password:           getEnv("REDFISH_PASSWORD", ""),
		TLSServerCACert:    getEnv("REDFISH_SERVER_CA_CERT", ""),
		InsecureSkipVerify: getEnvBool("REDFISH_INSECURE_SKIP_VERIFY", false),
		DiscoveryEnabled:   getEnvBool("REDFISH_DISCOVERY_ENABLED", false),
		DiscoveryInterval:  discoveryInterval,
	}

	return config, nil
}

func loadMCPConfig() (*MCPConfig, error) {
	transportStr := getEnv("MCP_TRANSPORT", string(MCPTransportStdio))
	var transport MCPTransport
	switch transportStr {
	case string(MCPTransportStdio):
		transport = MCPTransportStdio
	case string(MCPTransportSSE):
		transport = MCPTransportSSE
	case string(MCPTransportStreamableHTTP):
		transport = MCPTransportStreamableHTTP
	default:
		return nil, &ConfigError{
			Message: fmt.Sprintf("invalid transport: %s", transportStr),
		}
	}

	config := &MCPConfig{
		Transport: transport,
		LogLevel:  getEnv("MCP_REDFISH_LOG_LEVEL", "INFO"),
	}

	return config, nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue, minVal, maxVal int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return 0, &ConfigError{
			Message: fmt.Sprintf("environment variable %s must be an integer", key),
			Cause:   err,
		}
	}

	if intVal < minVal {
		return 0, &ConfigError{
			Message: fmt.Sprintf("environment variable %s must be >= %d, got: %d", key, minVal, intVal),
		}
	}

	if intVal > maxVal {
		return 0, &ConfigError{
			Message: fmt.Sprintf("environment variable %s must be <= %d, got: %d", key, maxVal, intVal),
		}
	}

	return intVal, nil
}
