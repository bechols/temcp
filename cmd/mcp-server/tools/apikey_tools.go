package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/clients"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
	cloudservicev1 "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	identityv1 "go.temporal.io/cloud-sdk/api/identity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RegisterApiKeyTools registers API key management tools with the MCP server
func RegisterApiKeyTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_create_api_key tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_create_api_key",
			mcp.WithDescription("Create a new Temporal Cloud API key"),
			mcp.WithString("owner_type", mcp.Description("Type of owner (user, service-account)"), mcp.Required()),
			mcp.WithString("owner_id", mcp.Description("ID of the owner"), mcp.Required()),
			mcp.WithString("display_name", mcp.Description("Display name for the API key"), mcp.Required()),
			mcp.WithString("expiry_time", mcp.Description("Expiry time in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)"), mcp.Required()),
			mcp.WithString("description", mcp.Description("Description for the API key (optional)")),
			mcp.WithBoolean("disabled", mcp.Description("Whether the API key should be disabled (optional, default: false)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleCreateApiKey(ctx, request, clientManager)
		},
	)
}

func handleCreateApiKey(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	// Extract required parameters
	ownerType, ok := arguments["owner_type"].(string)
	if !ok || ownerType == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: owner_type is required and must be a string (user or service-account)",
				},
			},
		}, nil
	}

	ownerID, ok := arguments["owner_id"].(string)
	if !ok || ownerID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: owner_id is required and must be a string",
				},
			},
		}, nil
	}

	displayName, ok := arguments["display_name"].(string)
	if !ok || displayName == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: display_name is required and must be a string",
				},
			},
		}, nil
	}

	expiryTimeStr, ok := arguments["expiry_time"].(string)
	if !ok || expiryTimeStr == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: expiry_time is required and must be a string in ISO 8601 format",
				},
			},
		}, nil
	}

	// Parse expiry time
	expiryTime, err := time.Parse(time.RFC3339, expiryTimeStr)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing expiry_time: %v. Use ISO 8601 format like 2024-12-31T23:59:59Z", err),
				},
			},
		}, nil
	}

	// Extract optional parameters
	description := ""
	if desc, ok := arguments["description"].(string); ok {
		description = desc
	}

	disabled := false
	if dis, ok := arguments["disabled"].(bool); ok {
		disabled = dis
	}

	// Convert owner type to enum - using the constants the backend expects
	var ownerTypeEnum identityv1.OwnerType
	switch ownerType {
	case "user":
		ownerTypeEnum = identityv1.OWNER_TYPE_USER
	case "service-account":
		ownerTypeEnum = identityv1.OWNER_TYPE_SERVICE_ACCOUNT
	default:
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: owner_type must be either 'user' or 'service-account'",
				},
			},
		}, nil
	}

	// Create the API key request
	createReq := &cloudservicev1.CreateApiKeyRequest{
		Spec: &identityv1.ApiKeySpec{
			OwnerId:     ownerID,
			OwnerType:   ownerTypeEnum,
			DisplayName: displayName,
			Description: description,
			ExpiryTime:  timestamppb.New(expiryTime),
			Disabled:    disabled,
		},
		AsyncOperationId: uuid.New().String(),
	}
	
	// Debug: log the enum value being sent (to stderr, not stdout which corrupts JSON)
	log.Printf("DEBUG: OwnerType enum value: %v (%d)", ownerTypeEnum, int32(ownerTypeEnum))

	// Call the Cloud API directly (API keys don't typically have workflows)
	cloudClient := clientManager.GetCloudClient()
	resp, err := cloudClient.CloudService().CreateApiKey(ctx, createReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error creating API key: %v", err),
				},
			},
		}, nil
	}

	// Create result structure
	result := map[string]interface{}{
		"api_key_id":         resp.KeyId,
		"token":              resp.Token,
		"async_operation_id": resp.AsyncOperation.Id,
		"state":              resp.AsyncOperation.State.String(),
		"request_details": map[string]interface{}{
			"owner_type":   ownerType,
			"owner_id":     ownerID,
			"display_name": displayName,
			"description":  description,
			"expiry_time":  expiryTimeStr,
			"disabled":     disabled,
		},
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
