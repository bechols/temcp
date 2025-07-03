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

Relevant if you are updating code to connect to Temporal Cloud. When updating code to connect to Temporal Cloud, make sure to update both the workflow and worker code to use the correct connection information (endpoint and namespace) and the correct authentication method.

Temporal Cloud acts as the Temporal Server, so you don't need to run a Temporal Server locally or update its configuration.

## Official Documentation for connecting to Temporal Cloud
- https://docs.temporal.io/develop/go/temporal-clients#connect-to-temporal-cloud (**Temporal Cloud Connection Guide (has links to Go, Typescript, Java, Python, .NET, and PHP documentation)**) 
- https://github.com/temporalio/samples-go/tree/main/helloworld-apiKey (**Sample Go Code with API Key**: )

## Key Connection Requirements

### 1. Use the right endpoint
- Specify the correct Temporal Cloud endpoint to connect to
- Format: region.api.temporal.io:7233
- Replace region with the actual region of your Temporal Cloud namespace

### 2. Namespace Configuration
- Specify your Temporal Cloud namespace
- Format: namespace.account-id.tmprl.cloud
- Replace namespace and account-id with your actual values

### 3. API Key Authentication
- Set TEMPORAL_CLOUD_API_KEY environment variable
- Use temporal.ClientOptions with the API key

## Sample Connection Code
` + "```go" + `
clientOptions := client.Options{
    HostPort: <endpoint>,
    Namespace: <namespace_id>.<account_id>,
    ConnectionOptions: client.ConnectionOptions{TLS: &tls.Config{}},
    Credentials: client.NewAPIKeyStaticCredentials(apiKey),
}
c, err := client.Dial(clientOptions)
` + "```" + `

Different languages have different ways to connect to Temporal Cloud.

Refer to the documentation links above for complete setup instructions and sample code.`

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: connectionInfo,
			},
		},
	}, nil
}
