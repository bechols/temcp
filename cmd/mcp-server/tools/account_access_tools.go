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

// RegisterAccountAccessTools registers account access tools with the MCP server
func RegisterAccountAccessTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_account_access tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_account_access",
			mcp.WithDescription("Get a user's account-level access role (owner, admin, developer, finance_admin, read) - for users only, not service accounts"),
			mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetAccountAccess(ctx, request, clientManager)
		},
	)
}

func handleGetAccountAccess(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
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

	var accountAccess interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise get user and extract access
	if clientManager.GetTemporalClient() != nil {
		// Use the GetAccountAccess workflow
		accountAccess, err = clientManager.ExecuteWorkflow(ctx, workflows.GetAccountAccessWorkflowType, userID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Error getting account access: %v", err),
					},
				},
			}, nil
		}
	} else {
		// Call GetUser directly and extract account access
		cloudClient := clientManager.GetCloudClient()
		getUserReq := &cloudservice.GetUserRequest{
			UserId: userID,
		}
		userResp, err := cloudClient.CloudService().GetUser(ctx, getUserReq)
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

		if userResp.User == nil || userResp.User.Spec == nil || userResp.User.Spec.Access == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("User %s has no access information", userID),
					},
				},
			}, nil
		}

		accountAccess = userResp.User.Spec.Access.GetAccountAccess()
	}

	// Create a more user-friendly result structure
	result := map[string]interface{}{
		"user_id":        userID,
		"account_access": accountAccess,
	}

	// If we have AccountAccess, add human-readable role information
	if aa, ok := accountAccess.(*identity.AccountAccess); ok && aa != nil {
		roleStr := aa.GetRole().String()
		roleDescription := getRoleDescription(aa.GetRole())

		result["role"] = roleStr
		result["role_description"] = roleDescription

		// Include deprecated role if present for compatibility
		if aa.GetRoleDeprecated() != "" {
			result["role_deprecated"] = aa.GetRoleDeprecated()
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

func getRoleDescription(role identity.AccountAccess_Role) string {
	switch role {
	case identity.AccountAccess_ROLE_OWNER:
		return "Gives full access to the account, including users, namespaces, and billing"
	case identity.AccountAccess_ROLE_ADMIN:
		return "Gives full access to the account, including users and namespaces"
	case identity.AccountAccess_ROLE_DEVELOPER:
		return "Gives access to create namespaces on the account"
	case identity.AccountAccess_ROLE_FINANCE_ADMIN:
		return "Gives read only access and write access for billing"
	case identity.AccountAccess_ROLE_READ:
		return "Gives read only access to the account"
	case identity.AccountAccess_ROLE_UNSPECIFIED:
		return "Role is not specified"
	default:
		return "Unknown role"
	}
}
