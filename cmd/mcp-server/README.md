# Temporal Cloud MCP Server

This MCP (Model Context Protocol) server exposes Temporal Cloud operations as tools that can be used by AI assistants to read and write Temporal Cloud data.

## Features

The server provides the following tool categories:

### User Management Tools
- `temporal_get_user` - Get a user by ID
- `temporal_list_users` - List users with pagination
- `temporal_get_all_users` - Get all users
- `temporal_find_user_by_email` - Find user by email
- `temporal_create_user` - Create a new user
- `temporal_update_user` - Update an existing user
- `temporal_delete_user` - Delete a user
- `temporal_reconcile_users` - Bulk reconcile users

### Namespace Management Tools
- `temporal_get_namespace` - Get a namespace by name
- `temporal_list_namespaces` - List namespaces with pagination
- `temporal_get_all_namespaces` - Get all namespaces
- `temporal_create_namespace` - Create a new namespace
- `temporal_update_namespace` - Update an existing namespace
- `temporal_delete_namespace` - Delete a namespace
- `temporal_reconcile_namespaces` - Bulk reconcile namespaces

### Region Management Tools
- `temporal_get_region` - Get region information
- `temporal_list_regions` - List all available regions

### Async Operation Tools
- `temporal_get_async_operation` - Get operation status
- `temporal_wait_for_operation` - Wait for operation completion

### Export Processing Tools
- `temporal_process_export` - Process exported workflow histories

## Configuration

The server is configured using environment variables:

### Required
- `TEMPORAL_CLOUD_API_KEY` - Your Temporal Cloud API key

### Optional
- `TEMPORAL_CLOUD_NAMESPACE` - Default namespace for operations
- `TEMPORAL_CLOUD_NAMESPACE_API_KEY` - Namespace-specific API key
- `TEMPORAL_CLOUD_NAMESPACE_TLS_CERT` - Path to mTLS certificate
- `TEMPORAL_CLOUD_NAMESPACE_TLS_KEY` - Path to mTLS private key
- `MCP_SERVER_NAME` - Server name (default: "temporal-cloud-mcp-server")
- `MCP_SERVER_VERSION` - Server version (default: "1.0.0")

## Usage

### Building
```bash
go build -o temporal-cloud-mcp-server ./cmd/mcp-server
```

### Running
```bash
export TEMPORAL_CLOUD_API_KEY="your-api-key"
./temporal-cloud-mcp-server
```

### Using with Claude Desktop
Add to your Claude Desktop configuration:
```json
{
  "mcpServers": {
    "temporal-cloud": {
      "command": "/path/to/temporal-cloud-mcp-server",
      "env": {
        "TEMPORAL_CLOUD_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Development Status

🚧 **Under Development** 🚧

This is the Phase 1 implementation with basic project structure and placeholder files. 

### Completed
- ✅ Project structure setup
- ✅ Configuration management
- ✅ Tool placeholders
- ✅ Basic server bootstrap

### TODO
- 🔄 Implement tool handlers
- 🔄 Integrate with existing workflows
- 🔄 Add error handling and validation
- 🔄 Add testing
- 🔄 Complete documentation