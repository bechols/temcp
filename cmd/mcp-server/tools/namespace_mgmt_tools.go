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
	"go.temporal.io/cloud-sdk/api/namespace/v1"
)

// RegisterNamespaceMgmtTools registers all namespace management tools with the MCP server
func RegisterNamespaceMgmtTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
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
			mcp.WithObject("namespace_spec", mcp.Description("Namespace specification object with required fields: name (string), regions (array of strings), retention_days (number), and optional fields like ca_certificate_base64, codec_server_endpoint, custom_search_attributes"), mcp.Required()),
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

	getNamespaceReq := &cloudservice.GetNamespaceRequest{
		Namespace: namespaceName,
	}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetNamespace workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetNamespaceWorkflowType, getNamespaceReq)
	} else {
		// Call GetNamespace through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().GetNamespace(ctx, getNamespaceReq)
	}
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

	// Convert result to JSON
	var resultData interface{}
	if clientManager.GetTemporalClient() != nil {
		// Workflow returns the namespace directly
		resultData = result
	} else {
		// Direct API call returns a response with .Namespace field
		if nsResponse, ok := result.(*cloudservice.GetNamespaceResponse); ok {
			resultData = nsResponse.Namespace
		} else {
			resultData = result
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

	getNamespacesReq := &cloudservice.GetNamespacesRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetNamespaces workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetNamespacesWorkflowType, getNamespacesReq)
	} else {
		// Call GetNamespaces through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().GetNamespaces(ctx, getNamespacesReq)
	}
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
	namespaceSpecJSON, marshalErr := json.Marshal(namespaceSpecRaw)
	if marshalErr != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace_spec: %v", marshalErr),
				},
			},
		}, nil
	}

	var namespaceSpec namespace.NamespaceSpec
	if unmarshalErr := json.Unmarshal(namespaceSpecJSON, &namespaceSpec); unmarshalErr != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace specification: %v", unmarshalErr),
				},
			},
		}, nil
	}

	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing CreateNamespace workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.CreateNamespaceWorkflowType, &namespaceSpec)
	} else {
		// Call CreateNamespace through cloud client
		cloudClient := clientManager.GetCloudClient()
		createReq := &cloudservice.CreateNamespaceRequest{
			Spec: &namespaceSpec,
		}
		result, err = cloudClient.CloudService().CreateNamespace(ctx, createReq)
	}
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
	updatesJSON, marshalErr := json.Marshal(namespaceUpdatesRaw)
	if marshalErr != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace_updates: %v", marshalErr),
				},
			},
		}, nil
	}

	var updateReq cloudservice.UpdateNamespaceRequest
	if unmarshalErr := json.Unmarshal(updatesJSON, &updateReq); unmarshalErr != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing namespace updates: %v", unmarshalErr),
				},
			},
		}, nil
	}
	updateReq.Namespace = namespaceName

	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing UpdateNamespace workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.UpdateNamespaceWorkflowType, &updateReq)
	} else {
		// Call UpdateNamespace through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().UpdateNamespace(ctx, &updateReq)
	}
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

	// First, get the namespace to obtain its resource version
	getNamespaceReq := &cloudservice.GetNamespaceRequest{
		Namespace: namespaceName,
	}
	var getResult interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetNamespace workflow
		getResult, err = clientManager.ExecuteWorkflow(ctx, workflows.GetNamespaceWorkflowType, getNamespaceReq)
	} else {
		// Call GetNamespace through cloud client
		cloudClient := clientManager.GetCloudClient()
		getResult, err = cloudClient.CloudService().GetNamespace(ctx, getNamespaceReq)
	}
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting namespace before deletion: %v", err),
				},
			},
		}, nil
	}

	// Extract the resource version
	var resourceVersion string
	if clientManager.GetTemporalClient() != nil {
		// Workflow returns a GetNamespaceResponse
		if nsResponse, ok := getResult.(*cloudservice.GetNamespaceResponse); ok {
			resourceVersion = nsResponse.Namespace.ResourceVersion
		} else {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: "Error: unable to extract namespace from workflow result",
					},
				},
			}, nil
		}
	} else {
		// Direct API call returns a response with .Namespace field
		if nsResponse, ok := getResult.(*cloudservice.GetNamespaceResponse); ok {
			resourceVersion = nsResponse.Namespace.ResourceVersion
		} else {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: "Error: unable to extract namespace from API response",
					},
				},
			}, nil
		}
	}

	// Now delete the namespace with the resource version
	deleteReq := &cloudservice.DeleteNamespaceRequest{
		Namespace:       namespaceName,
		ResourceVersion: resourceVersion,
	}
	var result interface{}

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing DeleteNamespace workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.DeleteNamespaceWorkflowType, deleteReq)
	} else {
		// Call DeleteNamespace through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().DeleteNamespace(ctx, deleteReq)
	}
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
