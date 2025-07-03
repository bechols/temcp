package tools

import (
	"bechols/temcp/cmd/mcp-server/clients"
	"bechols/temcp/cmd/mcp-server/config"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterAllTools(mcpServer *server.MCPServer, cfg *config.Config) error {
	clientManager, err := clients.NewClientManager(cfg)
	if err != nil {
		return err
	}

	RegisterUserTools(mcpServer, cfg, clientManager)

	RegisterAccountAccessTools(mcpServer, cfg, clientManager)

	RegisterNamespaceAccessTools(mcpServer, cfg, clientManager)

	RegisterNamespaceMgmtTools(mcpServer, cfg, clientManager)

	RegisterRegionTools(mcpServer, cfg, clientManager)

	RegisterOperationTools(mcpServer, cfg, clientManager)

	RegisterExportTools(mcpServer, cfg, clientManager)

	RegisterApiKeyTools(mcpServer, cfg, clientManager)

	RegisterServiceAccountTools(mcpServer, cfg, clientManager)

	RegisterNamespaceServiceAccountAccessTools(mcpServer, cfg, clientManager)

	RegisterConnectionInfoTools(mcpServer, cfg, clientManager)

	return nil
}
