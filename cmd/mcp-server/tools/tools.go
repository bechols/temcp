package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/clients"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
)

func RegisterAllTools(mcpServer *server.MCPServer, cfg *config.Config) error {
	clientManager, err := clients.NewClientManager(cfg)
	if err != nil {
		return err
	}

	RegisterUserTools(mcpServer, cfg, clientManager)

	RegisterAccountAccessTools(mcpServer, cfg, clientManager)

	RegisterNamespaceAccessTools(mcpServer, cfg, clientManager)

	RegisterNamespaceToolsImpl(mcpServer, cfg, clientManager)

	RegisterRegionToolsImpl(mcpServer, cfg, clientManager)

	RegisterOperationToolsImpl(mcpServer, cfg, clientManager)

	RegisterExportToolsImpl(mcpServer, cfg, clientManager)

	return nil
}
