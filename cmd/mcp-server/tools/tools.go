package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/clients"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
)

// RegisterAllTools registers all MCP tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, cfg *config.Config) error {
	// Create client manager
	clientManager, err := clients.NewClientManager(cfg)
	if err != nil {
		return err
	}

	// Register user management tools (simple version for Phase 2)
	RegisterUserToolsSimple(mcpServer, cfg, clientManager)
	
	// Register namespace management tools (Phase 3)
	RegisterNamespaceToolsImpl(mcpServer, cfg, clientManager)
	
	// Register region management tools (Phase 3)
	RegisterRegionToolsImpl(mcpServer, cfg, clientManager)
	
	// Register async operation tools (Phase 3)
	RegisterOperationToolsImpl(mcpServer, cfg, clientManager)
	
	// Register export processing tools (Phase 3)
	RegisterExportToolsImpl(mcpServer, cfg, clientManager)
	
	return nil
}