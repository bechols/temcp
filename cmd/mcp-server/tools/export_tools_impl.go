package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/clients"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
	"github.com/temporalio/cloud-samples-go/export"
)

// RegisterExportToolsImpl registers all export processing tools with the MCP server
func RegisterExportToolsImpl(mcpServer *server.MCPServer, cfg *config.Config, clientManager *clients.ClientManager) {
	// Register temporal_process_export tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_process_export",
			mcp.WithDescription("Process an exported Temporal workflow history file"),
			mcp.WithString("export_file_path", mcp.Description("Path to the exported workflow history file"), mcp.Required()),
			mcp.WithBoolean("format_readable", mcp.Description("Format output as human-readable JSON (optional)")),
			mcp.WithBoolean("include_metadata", mcp.Description("Include workflow metadata summary (optional)")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleProcessExportImpl(ctx, request, clientManager)
		},
	)

	// Register temporal_analyze_export tool
	mcpServer.AddTool(
		mcp.NewTool("temporal_analyze_export",
			mcp.WithDescription("Analyze exported workflow history and extract summary information"),
			mcp.WithString("export_file_path", mcp.Description("Path to the exported workflow history file"), mcp.Required()),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleAnalyzeExportImpl(ctx, request, clientManager)
		},
	)
}

func handleProcessExportImpl(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	filePath, ok := arguments["export_file_path"].(string)
	if !ok || filePath == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: export_file_path is required and must be a string",
				},
			},
		}, nil
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: export file not found at path: %s", filePath),
				},
			},
		}, nil
	}

	// Extract optional formatting options
	formatReadable := false
	if fr, ok := arguments["format_readable"].(bool); ok {
		formatReadable = fr
	}

	includeMetadata := false
	if im, ok := arguments["include_metadata"].(bool); ok {
		includeMetadata = im
	}

	// Read the export file
	exportData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error reading export file: %v", err),
				},
			},
		}, nil
	}

	// Deserialize the export
	workflowExecutions, err := export.DeserializeExportedWorkflows(exportData)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error deserializing export: %v", err),
				},
			},
		}, nil
	}

	// Process the data based on formatting options
	var result interface{}
	var resultText string

	if formatReadable && len(workflowExecutions.Items) > 0 {
		// Format as human-readable for the first workflow
		firstWorkflow := workflowExecutions.Items[0]
		resultText = export.FormatWorkflow(firstWorkflow)

		if includeMetadata {
			metadata, err := export.GetExportedWorkflowInformation(firstWorkflow)
			if err == nil {
				resultText = fmt.Sprintf("Workflow Information:\n%s\n\nFormatted Workflow:\n%s", metadata, resultText)
			}
		}
	} else {
		// Return as structured JSON
		result = workflowExecutions
		if includeMetadata && len(workflowExecutions.Items) > 0 {
			// Add metadata summary
			metadata := make(map[string]interface{})
			metadata["total_workflows"] = len(workflowExecutions.Items)

			if len(workflowExecutions.Items) > 0 {
				firstWorkflow := workflowExecutions.Items[0]
				workflowInfo, err := export.GetExportedWorkflowInformation(firstWorkflow)
				if err == nil {
					metadata["first_workflow_info"] = workflowInfo
				}
			}

			// Wrap result with metadata
			result = map[string]interface{}{
				"metadata":  metadata,
				"workflows": workflowExecutions,
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
		resultText = string(resultJSON)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func handleAnalyzeExportImpl(ctx context.Context, request mcp.CallToolRequest, clientManager *clients.ClientManager) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	filePath, ok := arguments["export_file_path"].(string)
	if !ok || filePath == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Error: export_file_path is required and must be a string",
				},
			},
		}, nil
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: export file not found at path: %s", filePath),
				},
			},
		}, nil
	}

	// Read the export file
	exportData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error reading export file: %v", err),
				},
			},
		}, nil
	}

	// Deserialize the export
	workflowExecutions, err := export.DeserializeExportedWorkflows(exportData)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error deserializing export: %v", err),
				},
			},
		}, nil
	}

	// Analyze the export and create summary
	analysis := map[string]interface{}{
		"file_path":       filePath,
		"file_size_bytes": len(exportData),
		"total_workflows": len(workflowExecutions.Items),
		"workflows":       []map[string]interface{}{},
	}

	// Analyze each workflow
	for i, workflow := range workflowExecutions.Items {
		workflowInfo, err := export.GetExportedWorkflowInformation(workflow)
		workflowAnalysis := map[string]interface{}{
			"index": i,
		}

		if err == nil {
			workflowAnalysis["info"] = workflowInfo
		} else {
			workflowAnalysis["info_error"] = err.Error()
		}

		// Count events in history
		if history := workflow.GetHistory(); history != nil {
			workflowAnalysis["event_count"] = len(history.GetEvents())
		}

		analysis["workflows"] = append(analysis["workflows"].([]map[string]interface{}), workflowAnalysis)
	}

	resultJSON, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error serializing analysis: %v", err),
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
