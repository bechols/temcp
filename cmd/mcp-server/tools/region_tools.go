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

// RegisterRegionTools registers all region management tools with the MCP server
func RegisterRegionTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_region tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_region",
			mcp.WithDescription("Get information about a specific Temporal Cloud region"),
			mcp.WithString("region_id", mcp.Description("Region ID"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetRegionImpl(ctx, request, clientManager)
		},
	)

	// Register temporal_list_regions tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_list_regions",
			mcp.WithDescription("List all available Temporal Cloud regions"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleListRegionsImpl(ctx, clientManager)
		},
	)
}

func handleGetRegionImpl(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	regionID, ok := arguments["region_id"].(string)
	if !ok || regionID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: region_id is required and must be a string",
				},
			},
		}, nil
	}

	getRegionReq := &cloudservice.GetRegionRequest{
		Region: regionID,
	}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetRegion workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetRegionWorkflowType, getRegionReq)
	} else {
		// Call Cloud API directly
		result, err = clientManager.GetCloudClient().CloudService().GetRegion(ctx, getRegionReq)
	}
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting region: %v", err),
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

func handleListRegionsImpl(ctx context.Context, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	getRegionsReq := &cloudservice.GetRegionsRequest{}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetAllRegions workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetAllRegionsWorkflowType, getRegionsReq)
	} else {
		// Call Cloud API directly
		result, err = clientManager.GetCloudClient().CloudService().GetRegions(ctx, getRegionsReq)
	}

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error listing regions: %v", err),
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
