package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nokia/mcp-redfish-go/pkg/config"
	"github.com/nokia/mcp-redfish-go/pkg/mcp"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "", "Path to Redfish config JSON file")
	flag.Parse()

	// Set config file environment variable if provided
	if *configFile != "" {
		os.Setenv("REDFISH_CONFIG_FILE", *configFile)
	}

	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Configuration loaded successfully",
		"redfish_hosts", len(cfg.Redfish.Hosts),
		"transport", cfg.MCP.Transport)

	// Create MCP server
	server, err := mcp.NewServer(cfg, logger)
	if err != nil {
		logger.Error("Failed to create MCP server", "error", err)
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the server
	logger.Info("Starting Redfish MCP server")
	if err := server.Start(ctx); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	logger.Info("Redfish MCP server stopped")
}
