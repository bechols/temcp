# temcp: Temporal Cloud MCP Server

An MCP (Model Context Protocol) server that exposes Temporal Cloud operations as tools for AI assistants.

### Building the MCP Server

```bash
go build -o mcp-server ./cmd/mcp-server
```

### Available Commands

The MCP server provides 13 tools for managing Temporal Cloud resources:

**User Management:**
- `temporal_get_user` - Get user details by ID
- `temporal_list_users` - List users with pagination

**Namespace Management:**
- `temporal_get_namespace` - Get namespace details by name
- `temporal_list_namespaces` - List namespaces with pagination
- `temporal_create_namespace` - Create a new namespace
- `temporal_update_namespace` - Update an existing namespace
- `temporal_delete_namespace` - Delete a namespace

**Region Management:**
- `temporal_get_region` - Get region details by ID
- `temporal_list_regions` - List all available regions

**Async Operations:**
- `temporal_get_async_operation` - Get async operation status
- `temporal_wait_for_operation` - Wait for async operation completion

**Export Processing:**
- `temporal_process_export` - Process exported workflow history files
- `temporal_analyze_export` - Analyze exported workflows and extract summaries

### Configuration

Set the following environment variables:
- `TEMPORAL_CLOUD_API_KEY` - Your Temporal Cloud API key
- `TEMPORAL_CLOUD_NAMESPACE` - Default namespace (optional)

### Usage with Cursor IDE

Configure `.cursor/mcp.json`:
```json
{
  "mcpServers": {
    "temporal-cloud": {
      "command": "/path/to/mcp-server",
      "args": [],
      "env": {
        "TEMPORAL_CLOUD_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Cloud Temporal Client

Sample Go temporal client for a namespace hosted on Temporal Cloud: [client.go](client/temporal/client.go) 

## Cloud API Workflows

Sample Workflow applications that exercises the Temporal Cloud APIs: [README.md](cmd/worker/README.md)

