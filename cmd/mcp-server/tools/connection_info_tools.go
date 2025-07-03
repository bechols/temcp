package tools

import (
	"context"

	"bechols/temcp/cmd/mcp-server/clients"
	"bechols/temcp/cmd/mcp-server/config"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterConnectionInfoTools registers all connection info tools with the MCP server
func RegisterConnectionInfoTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_cloud_connection_info tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_cloud_connection_info",
			mcp.WithDescription("How to connect workers and workflows to Temporal Cloud, including API key auth configuration."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleConnectionInfoImpl()
		},
	)
}

func handleConnectionInfoImpl() (*mcp.CallToolResult, error) {
	connectionInfo := `# Temporal Cloud Connection Information

If there is both workflow and worker code present, make sure to update the connection information for both.

## Official Documentation
- **Temporal Cloud Connection Guide (has links to Go, Typescript, Java, Python, .NET, and PHP documentation)**: https://docs.temporal.io/develop/go/temporal-clients#connect-to-temporal-cloud
- **Sample Go Code with API Key**: https://github.com/temporalio/samples-go/tree/main/helloworld-apiKey

## Key Connection Requirements

### 1. API Key Authentication
- Set TEMPORAL_CLOUD_API_KEY environment variable
- Use temporal.ClientOptions with the API key

### 2. Namespace Configuration
- Specify your Temporal Cloud namespace
- Format: namespace.account-id.tmprl.cloud

### 3. TLS Configuration
- Temporal Cloud requires TLS connections
- Use temporal.ClientOptions with TLS settings

## Sample Connection Code
` + "```go" + `
client, err := temporal.NewNamespaceClient(temporal.ClientOptions{
    HostPort:  "namespace.account-id.tmprl.cloud:7233",
    Namespace: "your-namespace",
    ConnectionOptions: temporal.ConnectionOptions{
        TLS: &tls.Config{
            ServerName: "namespace.account-id.tmprl.cloud",
        },
    },
    Credentials: temporal.NewAPIKeyStaticCredentials("your-api-key"),
})
` + "```" + `

Refer to the documentation links above for complete setup instructions.`

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: connectionInfo,
			},
		},
	}, nil
}
