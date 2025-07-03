package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"bechols/temcp/cmd/mcp-server/clients"
	"bechols/temcp/cmd/mcp-server/config"
	"bechols/temcp/workflows"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/api/identity/v1"
)

// RegisterNamespaceAccessTools registers namespace access tools with the MCP server
func RegisterNamespaceAccessTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_user_namespace_access tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_user_namespace_access",
			mcp.WithDescription("Get a user's access level for a specific namespace - for users only, not service accounts"),
			mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetUserNamespaceAccess(ctx, request, clientManager)
		},
	)

	// Register temporal_set_user_namespace_access tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_set_user_namespace_access",
			mcp.WithDescription("Set or update a user's access level for a specific namespace - for users only, not service accounts"),
			mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
			mcp.WithString("permission", mcp.Description("Permission level: ADMIN, WRITE, or READ"), mcp.Required()),
			mcp.WithString("resource_version", mcp.Description("Resource version for optimistic concurrency (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleSetUserNamespaceAccess(ctx, request, clientManager)
		},
	)
}

func handleGetUserNamespaceAccess(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	userID, ok := arguments["user_id"].(string)
	if !ok || userID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: user_id is required and must be a string",
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

	// Get the user details which includes namespace access information
	getUserReq := &cloudservice.GetUserRequest{
		UserId: userID,
	}

	var user interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		user, err = clientManager.ExecuteWorkflow(ctx, workflows.GetUserWorkflowType, getUserReq)
	} else {
		cloudClient := clientManager.GetCloudClient()
		user, err = cloudClient.CloudService().GetUser(ctx, getUserReq)
	}

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting user: %v", err),
				},
			},
		}, nil
	}

	// Extract user data and namespace access
	var userData interface{}
	if clientManager.GetTemporalClient() != nil {
		userData = user
	} else {
		if userResponse, ok := user.(*cloudservice.GetUserResponse); ok {
			userData = userResponse.User
		} else {
			userData = user
		}
	}

	// Convert to JSON and extract namespace access
	userJSON, err := json.Marshal(userData)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error serializing user data: %v", err),
				},
			},
		}, nil
	}

	// Parse the user data to extract namespace access
	var userMap map[string]interface{}
	if err := json.Unmarshal(userJSON, &userMap); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing user data: %v", err),
				},
			},
		}, nil
	}

	// Extract namespace access information
	result := map[string]interface{}{
		"user_id":    userID,
		"namespace":  namespace,
		"access":     nil,
		"has_access": false,
	}

	// Navigate through the user structure to find namespace access
	if spec, ok := userMap["spec"].(map[string]interface{}); ok {
		if access, ok := spec["access"].(map[string]interface{}); ok {
			if namespaceAccesses, ok := access["namespace_accesses"].(map[string]interface{}); ok {
				if namespaceAccess, exists := namespaceAccesses[namespace]; exists {
					result["access"] = namespaceAccess
					result["has_access"] = true
				}
			}
		}
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

func handleSetUserNamespaceAccess(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	userID, ok := arguments["user_id"].(string)
	if !ok || userID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: user_id is required and must be a string",
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
					Text: "Error: permission is required and must be one of ADMIN, WRITE, or READ",
				},
			},
		}, nil
	}

	// Convert permission string to the appropriate enum value
	var permission identity.NamespaceAccess_Permission
	switch permissionStr {
	case "ADMIN":
		permission = identity.NamespaceAccess_PERMISSION_ADMIN
	case "WRITE":
		permission = identity.NamespaceAccess_PERMISSION_WRITE
	case "READ":
		permission = identity.NamespaceAccess_PERMISSION_READ
	default:
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: permission must be one of ADMIN, WRITE, or READ",
				},
			},
		}, nil
	}

	// Extract optional resource_version
	resourceVersion := ""
	if rv, ok := arguments["resource_version"].(string); ok {
		resourceVersion = rv
	}

	// Create the SetUserNamespaceAccess request
	setAccessReq := &cloudservice.SetUserNamespaceAccessRequest{
		UserId:    userID,
		Namespace: namespace,
		Access: &identity.NamespaceAccess{
			Permission: permission,
		},
		ResourceVersion: resourceVersion,
	}

	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use SetUserNamespaceAccess workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.SetUserNamespaceAccessWorkflowType, setAccessReq)
	} else {
		// Call SetUserNamespaceAccess through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().SetUserNamespaceAccess(ctx, setAccessReq)
	}

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error setting user namespace access: %v", err),
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
