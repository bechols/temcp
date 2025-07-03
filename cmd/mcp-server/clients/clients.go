package clients

import (
	"context"
	"time"

	"github.com/temporalio/cloud-samples-go/client/api"
	"github.com/temporalio/cloud-samples-go/cmd/mcp-server/config"
	"github.com/temporalio/cloud-samples-go/workflows"
	"github.com/temporalio/cloud-samples-go/workflows/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// ClientManager manages Temporal Cloud API and workflow clients
type ClientManager struct {
	config         *config.Config
	cloudClient    *api.Client
	temporalClient client.Client
	worker         worker.Worker
	workflows      workflows.Workflows
	activities     *activities.Activities
}

// NewClientManager creates a new client manager with the given configuration
func NewClientManager(cfg *config.Config) (*ClientManager, error) {
	cm := &ClientManager{
		config: cfg,
	}

	// Initialize Cloud API client
	cloudClient, err := api.NewConnectionWithAPIKey(cfg.CloudAPIKey)
	if err != nil {
		return nil, err
	}
	cm.cloudClient = cloudClient

	// Initialize workflows and activities
	cm.workflows = workflows.NewWorkflows()
	cm.activities = workflows.NewActivities(cloudClient)

	// Initialize Temporal client if namespace auth is configured
	if cfg.HasNamespaceAuth() {
		// Create a simple local Temporal client for workflow execution
		// In production, this would connect to Temporal Cloud
		temporalClient, err := client.Dial(client.Options{
			// For now, use local development setup
			// TODO: Add proper Temporal Cloud connection
		})
		if err != nil {
			return nil, err
		}
		cm.temporalClient = temporalClient

		// Create and start worker
		cm.worker = worker.New(temporalClient, "mcp-task-queue", worker.Options{})
		workflows.Register(cm.worker, cm.workflows, cm.activities)
		
		// Start worker in background
		go func() {
			err := cm.worker.Run(worker.InterruptCh())
			if err != nil {
				// TODO: Add proper error handling
			}
		}()
	}

	return cm, nil
}

// GetCloudClient returns the Temporal Cloud API client
func (cm *ClientManager) GetCloudClient() *api.Client {
	return cm.cloudClient
}

// GetTemporalClient returns the Temporal workflow client
func (cm *ClientManager) GetTemporalClient() client.Client {
	return cm.temporalClient
}

// GetWorkflows returns the workflows interface
func (cm *ClientManager) GetWorkflows() workflows.Workflows {
	return cm.workflows
}

// Close closes all client connections
func (cm *ClientManager) Close() error {
	if cm.worker != nil {
		cm.worker.Stop()
	}
	if cm.temporalClient != nil {
		cm.temporalClient.Close()
	}
	return nil
}

// ExecuteWorkflow executes a Temporal workflow with the given parameters
func (cm *ClientManager) ExecuteWorkflow(ctx context.Context, workflowType string, args interface{}) (interface{}, error) {
	if cm.temporalClient == nil {
		return nil, nil // Direct workflow execution without Temporal client
	}

	// Create workflow options
	options := client.StartWorkflowOptions{
		TaskQueue: "mcp-task-queue",
	}

	// Start workflow
	workflowRun, err := cm.temporalClient.ExecuteWorkflow(ctx, options, workflowType, args)
	if err != nil {
		return nil, err
	}

	// Wait for result
	var result interface{}
	err = workflowRun.Get(ctx, &result)
	return result, err
}

// ExecuteWorkflowWithTimeout executes a workflow with a timeout
func (cm *ClientManager) ExecuteWorkflowWithTimeout(ctx context.Context, workflowType string, args interface{}, timeout time.Duration) (interface{}, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return cm.ExecuteWorkflow(timeoutCtx, workflowType, args)
}