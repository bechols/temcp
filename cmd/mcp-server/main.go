package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/tools"
)

func main() {
	log.Println("Temporal Cloud MCP Server starting...")

	// Load configuration from environment
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		cfg.ServerName,
		cfg.ServerVersion,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Register all tool handlers
	if err := tools.RegisterAllTools(mcpServer, cfg); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	log.Printf("Starting MCP server '%s' v%s on stdio...", cfg.ServerName, cfg.ServerVersion)
	server.ServeStdio(mcpServer)
}