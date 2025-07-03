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

func RegisterNamespaceServiceAccountAccessTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_service_account_namespace_access",
			mcp.WithDescription("Get namespace access permissions for a service account - for service accounts only, not users"),
			mcp.WithString("service_account_id", mcp.Description("Service account ID"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetServiceAccountNamespaceAccess(ctx, request, clientManager)
		},
	)

	mcpServer.AddTool(
		mcp.NewTool("temporal_set_service_account_namespace_access",
			mcp.WithDescription("Set namespace access permissions for a service account - for service accounts only, not users"),
			mcp.WithString("service_account_id", mcp.Description("Service account ID"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
			mcp.WithString("permission", mcp.Description("Permission level: admin, write, or read"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleSetServiceAccountNamespaceAccess(ctx, request, clientManager)
		},
	)
}

func handleGetServiceAccountNamespaceAccess(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	serviceAccountID, ok := arguments["service_account_id"].(string)
	if !ok || serviceAccountID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: service_account_id is required and must be a string",
				},
			},
		}, nil
	}

	getServiceAccountReq := &cloudservice.GetServiceAccountRequest{
		ServiceAccountId: serviceAccountID,
	}

	cloudClient := clientManager.GetCloudClient()
	result, err := cloudClient.CloudService().GetServiceAccount(ctx, getServiceAccountReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting service account: %v", err),
				},
			},
		}, nil
	}

	if result.ServiceAccount == nil || result.ServiceAccount.Spec == nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Service account or specification not found",
				},
			},
		}, nil
	}

	namespaceAccess := map[string]interface{}{
		"service_account_id": serviceAccountID,
		"namespace_access":   []map[string]interface{}{},
	}

	if result.ServiceAccount.Spec.NamespaceScopedAccess != nil {
		accessInfo := map[string]interface{}{
			"namespace": result.ServiceAccount.Spec.NamespaceScopedAccess.Namespace,
		}

		if result.ServiceAccount.Spec.NamespaceScopedAccess.Access != nil {
			accessInfo["permission"] = result.ServiceAccount.Spec.NamespaceScopedAccess.Access.Permission.String()
		}

		namespaceAccess["namespace_access"] = []map[string]interface{}{accessInfo}
	}

	resultJSON, err := json.MarshalIndent(namespaceAccess, "", "  ")
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

func handleSetServiceAccountNamespaceAccess(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	serviceAccountID, ok := arguments["service_account_id"].(string)
	if !ok || serviceAccountID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: service_account_id is required and must be a string",
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

	cloudClient := clientManager.GetCloudClient()
	getResult, err := cloudClient.CloudService().GetServiceAccount(ctx, &cloudservice.GetServiceAccountRequest{
		ServiceAccountId: serviceAccountID,
	})
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting service account: %v", err),
				},
			},
		}, nil
	}

	if getResult.ServiceAccount == nil || getResult.ServiceAccount.Spec == nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Service account or specification not found",
				},
			},
		}, nil
	}

	updatedSpec := &identity.ServiceAccountSpec{
		Name:        getResult.ServiceAccount.Spec.Name,
		Description: getResult.ServiceAccount.Spec.Description,
		Access: &identity.Access{
			AccountAccess: getResult.ServiceAccount.Spec.Access.AccountAccess,
			NamespaceAccesses: map[string]*identity.NamespaceAccess{
				namespace: {
					Permission: permission,
				},
			},
		},
	}

	updateReq := &cloudservice.UpdateServiceAccountRequest{
		ServiceAccountId: serviceAccountID,
		Spec:             updatedSpec,
		ResourceVersion:  getResult.ServiceAccount.ResourceVersion,
	}

	result, err := cloudClient.CloudService().UpdateServiceAccount(ctx, updateReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error updating service account namespace access: %v", err),
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
