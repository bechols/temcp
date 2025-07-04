package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"bechols/temcp/cmd/mcp-server/clients"
	"bechols/temcp/cmd/mcp-server/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/api/identity/v1"
)

// RegisterServiceAccountTools registers all service account management tools with the MCP server
func RegisterServiceAccountTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_list_service_accounts tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_list_service_accounts",
			mcp.WithDescription("List Temporal Cloud service accounts with pagination"),
			mcp.WithNumber("page_size", mcp.Description("Number of service accounts per page (optional)")),
			mcp.WithString("page_token", mcp.Description("Token for next page (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleListServiceAccounts(ctx, request, clientManager)
		},
	)

	// Register temporal_create_service_account tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_create_service_account",
			mcp.WithDescription("Create a new Temporal Cloud service account with namespace access"),
			mcp.WithString("name", mcp.Description("Service account name"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
			mcp.WithString("permission", mcp.Description("Permission level: admin, write, or read"), mcp.Required()),
			mcp.WithString("description", mcp.Description("Service account description (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleCreateServiceAccount(ctx, request, clientManager)
		},
	)
}

func handleListServiceAccounts(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	// Extract optional pagination parameters
	var pageSize int32 = 50 // default
	if ps, ok := arguments["page_size"].(float64); ok {
		pageSize = int32(ps)
	}

	pageToken := ""
	if token, ok := arguments["page_token"].(string); ok {
		pageToken = token
	}

	getServiceAccountsReq := &cloudservice.GetServiceAccountsRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	// Call GetServiceAccounts through cloud client
	cloudClient := clientManager.GetCloudClient()
	result, err := cloudClient.CloudService().GetServiceAccounts(ctx, getServiceAccountsReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error listing service accounts: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error serializing result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

func handleCreateServiceAccount(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	name, ok := arguments["name"].(string)
	if !ok || name == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: name is required and must be a string",
				},
			},
		}, nil
	}

	namespace, ok := arguments["namespace"].(string)
	if !ok || namespace == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: namespace is required and must be a string",
				},
			},
		}, nil
	}

	permissionStr, ok := arguments["permission"].(string)
	if !ok || permissionStr == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: permission is required and must be a string",
				},
			},
		}, nil
	}

	var permission identity.NamespaceAccess_Permission
	switch permissionStr {
	case "admin":
		permission = identity.NamespaceAccess_PERMISSION_ADMIN
	case "write":
		permission = identity.NamespaceAccess_PERMISSION_WRITE
	case "read":
		permission = identity.NamespaceAccess_PERMISSION_READ
	default:
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: permission must be 'admin', 'write', or 'read'",
				},
			},
		}, nil
	}

	description := ""
	if desc, ok := arguments["description"].(string); ok {
		description = desc
	}

	serviceAccountSpec := &identity.ServiceAccountSpec{
		Name:        name,
		Description: description,
		NamespaceScopedAccess: &identity.NamespaceScopedAccess{
			Namespace: namespace,
			Access: &identity.NamespaceAccess{
				Permission: permission,
			},
		},
	}

	createReq := &cloudservice.CreateServiceAccountRequest{
		Spec: serviceAccountSpec,
	}

	cloudClient := clientManager.GetCloudClient()
	result, err := cloudClient.CloudService().CreateServiceAccount(ctx, createReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error creating service account: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error serializing result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}
