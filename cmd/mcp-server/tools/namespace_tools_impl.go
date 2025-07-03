package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/clients"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
	"github.com/temporalio/cloud-samples-go/workflows"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/api/namespace/v1"
)

// RegisterNamespaceToolsImpl registers all namespace management tools with the MCP server
func RegisterNamespaceToolsImpl(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_namespace tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_namespace",
			mcp.WithDescription("Get a Temporal Cloud namespace by name"),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetNamespace(ctx, request, clientManager)
		},
	)

	// Register temporal_list_namespaces tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_list_namespaces",
			mcp.WithDescription("List Temporal Cloud namespaces with pagination"),
			mcp.WithNumber("page_size", mcp.Description("Number of namespaces per page (optional)")),
			mcp.WithString("page_token", mcp.Description("Token for next page (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleListNamespaces(ctx, request, clientManager)
		},
	)


	// Register temporal_create_namespace tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_create_namespace",
			mcp.WithDescription("Create a new Temporal Cloud namespace"),
			mcp.WithObject("namespace_spec", mcp.Description("Namespace specification object"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleCreateNamespace(ctx, request, clientManager)
		},
	)

	// Register temporal_update_namespace tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_update_namespace",
			mcp.WithDescription("Update an existing Temporal Cloud namespace"),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
			mcp.WithObject("namespace_updates", mcp.Description("Namespace updates object"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleUpdateNamespace(ctx, request, clientManager)
		},
	)

	// Register temporal_delete_namespace tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_delete_namespace",
			mcp.WithDescription("Delete a Temporal Cloud namespace"),
			mcp.WithString("namespace", mcp.Description("Namespace name"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleDeleteNamespace(ctx, request, clientManager)
		},
	)
}

func handleGetNamespace(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	namespaceName, ok := arguments["namespace"].(string)
	if !ok || namespaceName == "" {
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

	// Call GetNamespace through cloud client
	cloudClient := clientManager.GetCloudClient()
	getNamespaceReq := &cloudservice.GetNamespaceRequest{
		Namespace: namespaceName,
	}

	result, err := cloudClient.CloudService().GetNamespace(ctx, getNamespaceReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting namespace: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result.Namespace, "", "  ")
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

func handleListNamespaces(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
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

	// Call GetNamespaces through cloud client
	cloudClient := clientManager.GetCloudClient()
	getNamespacesReq := &cloudservice.GetNamespacesRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	result, err := cloudClient.CloudService().GetNamespaces(ctx, getNamespacesReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error listing namespaces: %v", err),
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


func handleCreateNamespace(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	namespaceSpecRaw, ok := arguments["namespace_spec"]
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: namespace_spec is required",
				},
			},
		}, nil
	}

	// Convert namespace_spec to proper type
	namespaceSpecJSON, err := json.Marshal(namespaceSpecRaw)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace_spec: %v", err),
				},
			},
		}, nil
	}

	var namespaceSpec namespace.NamespaceSpec
	if err := json.Unmarshal(namespaceSpecJSON, &namespaceSpec); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace specification: %v", err),
				},
			},
		}, nil
	}

	// Use the existing CreateNamespace workflow
	result, err := clientManager.ExecuteWorkflow(ctx, workflows.CreateNamespaceWorkflowType, &namespaceSpec)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error creating namespace: %v", err),
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

func handleUpdateNamespace(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	namespaceName, ok := arguments["namespace"].(string)
	if !ok || namespaceName == "" {
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

	namespaceUpdatesRaw, ok := arguments["namespace_updates"]
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: namespace_updates is required",
				},
			},
		}, nil
	}

	// Convert to proper update request
	updatesJSON, err := json.Marshal(namespaceUpdatesRaw)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace_updates: %v", err),
				},
			},
		}, nil
	}

	var updateReq cloudservice.UpdateNamespaceRequest
	if err := json.Unmarshal(updatesJSON, &updateReq); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace updates: %v", err),
				},
			},
		}, nil
	}
	updateReq.Namespace = namespaceName

	// Use the existing UpdateNamespace workflow
	result, err := clientManager.ExecuteWorkflow(ctx, workflows.UpdateNamespaceWorkflowType, &updateReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error updating namespace: %v", err),
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

func handleDeleteNamespace(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	namespaceName, ok := arguments["namespace"].(string)
	if !ok || namespaceName == "" {
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

	// Use the existing DeleteNamespace workflow
	deleteReq := &cloudservice.DeleteNamespaceRequest{
		Namespace: namespaceName,
	}

	result, err := clientManager.ExecuteWorkflow(ctx, workflows.DeleteNamespaceWorkflowType, deleteReq)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error deleting namespace: %v", err),
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