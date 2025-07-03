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
)

func RegisterUserTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_user tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_user",
			mcp.WithDescription("Get a Temporal Cloud user by ID"),
			mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetUser(ctx, request, clientManager)
		},
	)

	// Register temporal_list_users tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_list_users",
			mcp.WithDescription("List Temporal Cloud users with pagination"),
			mcp.WithNumber("page_size", mcp.Description("Number of users per page (optional)")),
			mcp.WithString("page_token", mcp.Description("Token for next page (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleListUsers(ctx, request, clientManager)
		},
	)

}

func handleGetUser(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	// Extract user_id from request arguments
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

	getUserReq := &cloudservice.GetUserRequest{
		UserId: userID,
	}
	var user interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetUser workflow
		user, err = clientManager.ExecuteWorkflow(ctx, workflows.GetUserWorkflowType, getUserReq)
	} else {
		// Call GetUser activity directly through cloud client
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

	// Convert result to JSON
	var resultData interface{}
	if clientManager.GetTemporalClient() != nil {
		// Workflow returns the user directly
		resultData = user
	} else {
		// Direct API call returns a response with .User field
		if userResponse, ok := user.(*cloudservice.GetUserResponse); ok {
			resultData = userResponse.User
		} else {
			resultData = user
		}
	}

	resultJSON, err := json.MarshalIndent(resultData, "", "  ")
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

func handleListUsers(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
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

	getUsersReq := &cloudservice.GetUsersRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetUsers workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetUsersWorkflowType, getUsersReq)
	} else {
		// Call GetUsers through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().GetUsers(ctx, getUsersReq)
	}
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error listing users: %v", err),
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
