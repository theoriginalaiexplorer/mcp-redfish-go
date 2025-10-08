package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	// Set up test environment variables
	os.Setenv("REDFISH_HOSTS", `[{"address": "test.example.com", "port": 8443}]`)
	os.Setenv("REDFISH_AUTH_METHOD", "session")
	os.Setenv("REDFISH_USERNAME", "testuser")
	os.Setenv("REDFISH_PASSWORD", "testpass")
	os.Setenv("MCP_TRANSPORT", "stdio")
	os.Setenv("MCP_REDFISH_LOG_LEVEL", "INFO")

	defer func() {
		// Clean up
		os.Unsetenv("REDFISH_HOSTS")
		os.Unsetenv("REDFISH_AUTH_METHOD")
		os.Unsetenv("REDFISH_USERNAME")
		os.Unsetenv("REDFISH_PASSWORD")
		os.Unsetenv("MCP_TRANSPORT")
		os.Unsetenv("MCP_REDFISH_LOG_LEVEL")
	}()

	config, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if config == nil {
		t.Fatal("Config is nil")
	}

	if config.Redfish == nil {
		t.Fatal("Redfish config is nil")
	}

	if len(config.Redfish.Hosts) != 1 {
		t.Fatalf("Expected 1 host, got %d", len(config.Redfish.Hosts))
	}

	host := config.Redfish.Hosts[0]
	if host.Address != "test.example.com" {
		t.Errorf("Expected address 'test.example.com', got '%s'", host.Address)
	}

	if host.Port != 8443 {
		t.Errorf("Expected port 8443, got %d", host.Port)
	}

	if config.Redfish.AuthMethod != "session" {
		t.Errorf("Expected auth method 'session', got '%s'", config.Redfish.AuthMethod)
	}

	if config.MCP == nil {
		t.Fatal("MCP config is nil")
	}

	if config.MCP.Transport != MCPTransportStdio {
		t.Errorf("Expected transport 'stdio', got '%s'", config.MCP.Transport)
	}
}

func TestConfigValidation(t *testing.T) {
	config := &Config{
		Redfish: &RedfishConfig{
			Hosts: []HostConfig{
				{Address: "valid.example.com"},
			},
			Port:              443,
			AuthMethod:        "session",
			Username:          "user",
			Password:          "pass",
			DiscoveryEnabled:  false,
			DiscoveryInterval: 30,
		},
		MCP: &MCPConfig{
			Transport: MCPTransportStdio,
			LogLevel:  "INFO",
		},
	}

	err := config.Validate()
	if err != nil {
		t.Fatalf("Valid config failed validation: %v", err)
	}
}

func TestInvalidConfig(t *testing.T) {
	// Test invalid transport
	config := &Config{
		Redfish: &RedfishConfig{
			Hosts: []HostConfig{
				{Address: "valid.example.com"},
			},
			Port:              443,
			AuthMethod:        "session",
			Username:          "user",
			Password:          "pass",
			DiscoveryEnabled:  false,
			DiscoveryInterval: 30,
		},
		MCP: &MCPConfig{
			Transport: "invalid_transport",
			LogLevel:  "INFO",
		},
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Invalid config passed validation")
	}
}
