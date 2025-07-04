#!/bin/bash

# Sample usage: export TEMPORAL_CLOUD_API_KEY="api-key" && ./test_mcp.sh │ list_namespaces | ./mcp-server

# Function to run initialization
init_server() {
    echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}'
}

# Function to list tools
list_tools() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
}

# Region tools
test_list_regions() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "temporal_list_regions", "arguments": {}}}'
}

test_get_region() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "temporal_get_region", "arguments": {"region_id": "us-east-1"}}}'
}

# User tools
test_list_users() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "temporal_list_users", "arguments": {"page_size": 10}}}'
}

test_get_user() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 6, "method": "tools/call", "params": {"name": "temporal_get_user", "arguments": {"user_id": "test-user-id"}}}'
}

# Namespace tools
test_list_namespaces() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 7, "method": "tools/call", "params": {"name": "temporal_list_namespaces", "arguments": {"page_size": 10}}}'
}

test_get_namespace() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 8, "method": "tools/call", "params": {"name": "temporal_get_namespace", "arguments": {"namespace": "test-namespace"}}}'
}

test_create_namespace() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 9, "method": "tools/call", "params": {"name": "temporal_create_namespace", "arguments": {"namespace_spec": {"name": "test-ns", "regions": ["us-east-1"], "retention_days": 30}}}}'
}

# Operation tools
test_get_operation() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 9, "method": "tools/call", "params": {"name": "temporal_get_async_operation", "arguments": {"operation_id": "test-operation-id"}}}'
}

# Connection info tools
test_connection_info() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 10, "method": "tools/call", "params": {"name": "temporal_cloud_connection_info", "arguments": {}}}'
}

# Service account tools
test_create_service_account() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 11, "method": "tools/call", "params": {"name": "temporal_create_service_account", "arguments": {"name": "test-sa", "namespace": "test-namespace", "permission": "admin", "description": "Test service account"}}}'
}

# Show usage if no arguments
show_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Available commands:"
    echo "  list_tools         - List all available tools"
    echo "  test_list_regions  - Test temporal_list_regions tool"
    echo "  test_get_region    - Test temporal_get_region tool"
    echo "  test_list_users    - Test temporal_list_users tool"
    echo "  test_get_user      - Test temporal_get_user tool"
    echo "  test_list_namespaces - Test temporal_list_namespaces tool"
    echo "  test_get_namespace - Test temporal_get_namespace tool"
    echo "  test_create_namespace - Test temporal_create_namespace tool (with API key auth default)"
    echo "  test_get_operation - Test temporal_get_async_operation tool"
    echo "  test_connection_info - Test temporal_cloud_connection_info tool"
    echo "  test_create_service_account - Test temporal_create_service_account tool"
    echo "  all               - Run all tests"
    echo ""
    echo "Example usage:"
    echo "  $0 test_list_regions | ./mcp-server"
    echo "  $0 list_tools | ./mcp-server"
}

# Run all tests
run_all_tests() {
    echo "# Testing MCP Server Tools - All Tests"
    init_server
    echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
    echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "temporal_list_regions", "arguments": {}}}'
    echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "temporal_get_region", "arguments": {"region_id": "us-east-1"}}}'
    echo '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "temporal_list_users", "arguments": {"page_size": 10}}}'
    echo '{"jsonrpc": "2.0", "id": 7, "method": "tools/call", "params": {"name": "temporal_list_namespaces", "arguments": {"page_size": 10}}}'
}

# Main script logic
case "${1:-}" in
    "list_tools")
        list_tools
        ;;
    "test_list_regions")
        test_list_regions
        ;;
    "test_get_region")
        test_get_region
        ;;
    "test_list_users")
        test_list_users
        ;;
    "test_get_user")
        test_get_user
        ;;
    "test_list_namespaces")
        test_list_namespaces
        ;;
    "test_get_namespace")
        test_get_namespace
        ;;
    "test_create_namespace")
        test_create_namespace
        ;;
    "test_get_operation")
        test_get_operation
        ;;
    "test_connection_info")
        test_connection_info
        ;;
    "test_create_service_account")
        test_create_service_account
        ;;
    "all")
        run_all_tests
        ;;
    *)
        show_usage
        ;;
esac