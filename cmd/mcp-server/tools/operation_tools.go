package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bechols/temcp/cmd/mcp-server/clients"
	"bechols/temcp/cmd/mcp-server/config"
	"bechols/temcp/workflows"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/api/operation/v1"
)

// RegisterOperationTools registers all async operation management tools with the MCP server
func RegisterOperationTools(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_get_async_operation tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_get_async_operation",
			mcp.WithDescription("Get the status of an async operation"),
			mcp.WithString("operation_id", mcp.Description("Async operation ID"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetAsyncOperationImpl(ctx, request, clientManager)
		},
	)

	// Register temporal_wait_for_operation tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_wait_for_operation",
			mcp.WithDescription("Wait for an async operation to complete with optional timeout"),
			mcp.WithString("operation_id", mcp.Description("Async operation ID"), mcp.Required()),
			mcp.WithNumber("timeout_seconds", mcp.Description("Timeout in seconds (optional, default 300)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleWaitForOperationImpl(ctx, request, clientManager)
		},
	)
}

func handleGetAsyncOperationImpl(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	operationID, ok := arguments["operation_id"].(string)
	if !ok || operationID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: operation_id is required and must be a string",
				},
			},
		}, nil
	}

	getOpReq := &cloudservice.GetAsyncOperationRequest{
		AsyncOperationId: operationID,
	}
	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise call API directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing GetAsyncOperation workflow
		result, err = clientManager.ExecuteWorkflow(ctx, workflows.GetAsyncOperationWorkflowType, getOpReq)
	} else {
		// Call GetAsyncOperation through cloud client
		cloudClient := clientManager.GetCloudClient()
		result, err = cloudClient.CloudService().GetAsyncOperation(ctx, getOpReq)
	}
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting async operation: %v", err),
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

func handleWaitForOperationImpl(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	operationID, ok := arguments["operation_id"].(string)
	if !ok || operationID == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: operation_id is required and must be a string",
				},
			},
		}, nil
	}

	// Extract optional timeout (default to 5 minutes)
	timeoutSeconds := 300.0 // default 5 minutes
	if ts, ok := arguments["timeout_seconds"].(float64); ok {
		timeoutSeconds = ts
	}

	var result interface{}
	var err error

	// Use workflow if Temporal client is available, otherwise implement polling directly
	if clientManager.GetTemporalClient() != nil {
		// Use the existing WaitForAsyncOperation workflow
		waitInput := &workflows.WaitForAsyncOperationInput{
			AsyncOperationID: operationID,
			Timeout:          time.Duration(timeoutSeconds) * time.Second,
		}
		result, err = clientManager.ExecuteWorkflowWithTimeout(ctx, workflows.WaitForAsyncOperationType, waitInput, time.Duration(timeoutSeconds+30)*time.Second)
	} else {
		// Implement polling logic directly
		cloudClient := clientManager.GetCloudClient()
		getOpReq := &cloudservice.GetAsyncOperationRequest{
			AsyncOperationId: operationID,
		}

		// Set up timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
		defer cancel()

		// Poll until complete or timeout
		for {
			opResult, pollErr := cloudClient.CloudService().GetAsyncOperation(timeoutCtx, getOpReq)
			if pollErr != nil {
				err = pollErr
				break
			}

			// Check if operation is complete
			if opResult.AsyncOperation.State != operation.AsyncOperation_STATE_PENDING &&
				opResult.AsyncOperation.State != operation.AsyncOperation_STATE_IN_PROGRESS {
				// Operation completed (success or failure)
				result = opResult
				break
			}

			// Wait before next poll
			select {
			case <-timeoutCtx.Done():
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						mcp.TextContent{
							Type: "text",
							Text: "Timeout waiting for async operation to complete",
						},
					},
				}, nil
			case <-time.After(2 * time.Second):
				// Continue polling
			}
		}
	}
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error waiting for async operation: %v", err),
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
