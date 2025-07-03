# Temporal Cloud MCP Server Implementation Plan

## Overview
This plan outlines the integration of the existing Temporal Cloud Ops API functionality with an MCP (Model Context Protocol) server using mcp-go. The goal is to expose all Temporal Cloud operations as MCP tools that can be used by AI assistants to read and write Temporal Cloud data.

## Current State Analysis
The repository already contains comprehensive Temporal Cloud API functionality:
- **User Management**: CRUD operations, bulk reconciliation
- **Namespace Management**: Full lifecycle management
- **Region Management**: Query available regions
- **Async Operation Tracking**: Monitor long-running operations
- **Export Processing**: Parse workflow history exports
- **Multiple Authentication Methods**: API keys, mTLS, namespace-specific auth

## Implementation Plan

### Phase 1: Project Structure Setup
1. **Add mcp-go dependency**
   - Add `github.com/mark3labs/mcp-go` to go.mod
   - Study the simple_client example for integration patterns

2. **Create MCP server directory structure**
   ```
   cmd/mcp-server/
   ├── main.go              # MCP server entry point
   ├── tools/               # MCP tool implementations
   │   ├── user_tools.go    # User management tools
   │   ├── namespace_tools.go # Namespace management tools
   │   ├── region_tools.go  # Region query tools
   │   ├── operation_tools.go # Async operation tools
   │   └── export_tools.go  # Export processing tools
   ├── config/
   │   └── config.go        # Configuration management
   └── README.md            # MCP server documentation
   ```

### Phase 2: Core MCP Server Implementation

#### 2.1 Server Bootstrap (`cmd/mcp-server/main.go`)
- Initialize MCP server using mcp-go
- Load configuration from environment variables
- Register all tool handlers
- Set up proper logging and error handling
- Handle graceful shutdown

#### 2.2 Configuration Management (`cmd/mcp-server/config/config.go`)
- Environment variable configuration:
  - `TEMPORAL_CLOUD_API_KEY` - Primary Cloud API key
  - `TEMPORAL_CLOUD_NAMESPACE` - Default namespace
  - `TEMPORAL_CLOUD_NAMESPACE_API_KEY` - Namespace-specific API key
  - `TEMPORAL_CLOUD_NAMESPACE_TLS_CERT` - mTLS certificate path
  - `TEMPORAL_CLOUD_NAMESPACE_TLS_KEY` - mTLS key path
  - `MCP_SERVER_NAME` - Server identification
- Validation of required configuration
- Support for multiple authentication methods

### Phase 3: MCP Tool Implementation

#### 3.1 User Management Tools (`cmd/mcp-server/tools/user_tools.go`)

**Tools to implement:**
1. **`temporal_get_user`**
   - Input: `user_id` (string)
   - Output: User details (JSON)
   - Uses: `GetUser` workflow

2. **`temporal_list_users`**
   - Input: `page_size` (optional int), `next_page_token` (optional string)
   - Output: Users list with pagination info
   - Uses: `GetUsers` workflow

3. **`temporal_get_all_users`**
   - Input: None
   - Output: All users (handles pagination internally)
   - Uses: `GetAllUsers` workflow

4. **`temporal_find_user_by_email`**
   - Input: `email` (string) 
   - Output: User details or not found
   - Uses: `GetUserWithEmail` workflow

5. **`temporal_create_user`**
   - Input: User specification (JSON)
   - Output: Created user details
   - Uses: `CreateUser` workflow

6. **`temporal_update_user`**
   - Input: `user_id` (string), user updates (JSON)
   - Output: Updated user details
   - Uses: `UpdateUser` workflow

7. **`temporal_delete_user`**
   - Input: `user_id` (string)
   - Output: Deletion confirmation
   - Uses: `DeleteUser` workflow

8. **`temporal_reconcile_users`**
   - Input: Users specification array (JSON), `delete_unspecified` (optional bool)
   - Output: Reconciliation results
   - Uses: `ReconcileUsers` workflow

#### 3.2 Namespace Management Tools (`cmd/mcp-server/tools/namespace_tools.go`)

**Tools to implement:**
1. **`temporal_get_namespace`**
   - Input: `namespace` (string)
   - Output: Namespace details
   - Uses: `GetNamespace` workflow

2. **`temporal_list_namespaces`**
   - Input: `page_size` (optional int), `next_page_token` (optional string)
   - Output: Namespaces list with pagination
   - Uses: `GetNamespaces` workflow

3. **`temporal_get_all_namespaces`**
   - Input: None
   - Output: All namespaces
   - Uses: `GetAllNamespaces` workflow

4. **`temporal_create_namespace`**
   - Input: Namespace specification (JSON)
   - Output: Created namespace details
   - Uses: `CreateNamespace` workflow

5. **`temporal_update_namespace`**
   - Input: `namespace` (string), updates (JSON)
   - Output: Updated namespace details
   - Uses: `UpdateNamespace` workflow

6. **`temporal_delete_namespace`**
   - Input: `namespace` (string)
   - Output: Deletion confirmation
   - Uses: `DeleteNamespace` workflow

7. **`temporal_reconcile_namespaces`**
   - Input: Namespaces specification array (JSON)
   - Output: Reconciliation results
   - Uses: `ReconcileNamespaces` workflow

#### 3.3 Region and Operation Tools

**Region Tools (`cmd/mcp-server/tools/region_tools.go`):**
1. **`temporal_get_region`**
   - Input: `region_id` (string)
   - Output: Region details
   - Uses: `GetRegion` workflow

2. **`temporal_list_regions`**
   - Input: None
   - Output: All available regions
   - Uses: `GetAllRegions` workflow

**Operation Tools (`cmd/mcp-server/tools/operation_tools.go`):**
1. **`temporal_get_async_operation`**
   - Input: `operation_id` (string)
   - Output: Operation status and details
   - Uses: `GetAsyncOperation` workflow

2. **`temporal_wait_for_operation`**
   - Input: `operation_id` (string), `timeout_seconds` (optional int)
   - Output: Final operation result
   - Uses: `WaitForAsyncOperation` workflow

#### 3.4 Export Processing Tools (`cmd/mcp-server/tools/export_tools.go`)

1. **`temporal_process_export`**
   - Input: `export_file_path` (string)
   - Output: Processed export data (JSON)
   - Uses: Export processing functionality from `/export/`

### Phase 4: Integration Layer

#### 4.1 Workflow Client Integration
- Create helper functions to instantiate Temporal clients
- Reuse existing `client/api/client.go` for Cloud API operations
- Reuse existing `client/temporal/client.go` for workflow execution
- Handle authentication method selection based on configuration

#### 4.2 Error Handling and Validation
- Implement consistent error response format for MCP tools
- Use existing validation from `/internal/validator/`
- Map Temporal errors to appropriate MCP error codes
- Provide helpful error messages for common issues

#### 4.3 Logging and Monitoring
- Integrate with existing logging patterns
- Add MCP-specific logging for tool invocations
- Include request/response logging for debugging
- Add metrics for tool usage (optional)

### Phase 5: Testing and Documentation

#### 5.1 Testing Strategy
- Unit tests for each MCP tool
- Integration tests with mock Temporal Cloud API
- End-to-end tests with real Temporal Cloud instance
- Test error conditions and edge cases

#### 5.2 Documentation
- Update main README with MCP server usage
- Create MCP server specific README in `cmd/mcp-server/`
- Document all available tools with input/output schemas
- Provide configuration examples
- Add troubleshooting guide

### Phase 6: Deployment and Usage

#### 6.1 Build and Distribution
- Add MCP server build target to existing build process
- Consider creating Docker image for easy deployment
- Document installation and setup process

#### 6.2 Integration Examples
- Provide example of using the MCP server with Claude
- Show common usage patterns for Temporal Cloud management
- Create sample workflows for typical operations

## Implementation Notes

### Authentication Strategy
The MCP server will support the same authentication methods as the existing codebase:
1. **Primary**: Cloud API key for all Cloud operations
2. **Namespace-specific**: API key or mTLS for Temporal client connections
3. **Fallback**: Environment variable configuration with validation

### Tool Design Principles
1. **Consistency**: All tools follow the same input/output patterns
2. **Validation**: Input validation using existing validator framework
3. **Error Handling**: Consistent error response format
4. **Documentation**: Self-documenting with clear schemas
5. **Performance**: Reuse existing client connections and workflows

### Development Phases
- **Phase 1-2**: Basic server setup (1-2 days)
- **Phase 3**: Core tool implementation (3-5 days) 
- **Phase 4**: Integration and polish (2-3 days)
- **Phase 5**: Testing and documentation (2-3 days)
- **Phase 6**: Deployment preparation (1 day)

**Total Estimated Time**: 9-16 days depending on testing depth and documentation detail.

## Success Criteria
1. MCP server successfully exposes all Temporal Cloud operations as tools
2. AI assistants can use tools to read and write Temporal Cloud data
3. Authentication works with multiple methods
4. Error handling provides clear, actionable feedback
5. Performance is acceptable for interactive use
6. Documentation enables easy setup and usage