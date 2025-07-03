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

Make sure you're using an up to date version of the Temporal SDK. As of July 3 2025, the latest version is 1.34.0.

## Key Connection Requirements

### 1. Use the right endpoint
- **For API Key authentication**: Use temporal.io endpoints (e.g., us-east-1.api.temporal.io:7233)
- **For mTLS authentication**: Use tmprl.cloud endpoints (e.g., namespace.account-id.tmprl.cloud:7233)

### 2. Namespace Configuration
- Specify your Temporal Cloud namespace
- Format: namespace.account-id.tmprl.cloud
- Replace namespace and account-id with your actual values

## Authentication Methods

### API Key Authentication
Use temporal.io endpoints with API key credentials:

` + "```go" + `
clientOpts := client.Options{
    HostPort:  "us-east-1.api.temporal.io:7233", // Use temporal.io endpoint
    Namespace: "namespace.account-id.tmprl.cloud",
    Credentials: client.NewAPIKeyStaticCredentials(apiKey),
    ConnectionOptions: client.ConnectionOptions{
        TLS: &tls.Config{},
    },
}
c, err := client.Dial(clientOpts)
` + "```" + `

### mTLS Authentication
Use tmprl.cloud endpoints with client certificates:

` + "```go" + `
cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
if err != nil {
    return fmt.Errorf("failed loading key pair: %w", err)
}

clientOpts := client.Options{
    HostPort:  "namespace.account-id.tmprl.cloud:7233", // Use tmprl.cloud endpoint
    Namespace: "namespace.account-id.tmprl.cloud",
    ConnectionOptions: client.ConnectionOptions{
        TLS: &tls.Config{
            Certificates: []tls.Certificate{cert},
        },
    },
}
c, err := client.Dial(clientOpts)
` + "```" + `

## Important Notes
- **Warning**: Using an API key with tmprl.cloud endpoints is invalid and will cause connection failures
- **Warning**: Using mTLS certificates with temporal.io endpoints is invalid and will cause connection failures
- Always match the endpoint domain with the correct authentication method
- The namespace format is the same for both authentication methods

Different languages have different ways to connect to Temporal Cloud.

Refer to the documentation links below for additional setup instructions and sample code for other languages.

## Official Documentation for connecting to Temporal Cloud
- https://docs.temporal.io/develop/go/temporal-clients#connect-to-temporal-cloud (**Temporal Cloud Connection Guide (has links to Go, Typescript, Java, Python, .NET, and PHP documentation)**) 
- https://github.com/temporalio/samples-go/tree/main/helloworld-apiKey (**Sample Go Code with API Key**: )
`

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: connectionInfo,
			},
		},
	}, nil
}
