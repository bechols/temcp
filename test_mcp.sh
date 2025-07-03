#!/bin/bash

# Sample usage: export TEMPORAL_CLOUD_API_KEY="api-key" && ./test_mcp.sh │ list_namespaces | ./mcp-server       

# Export API key

export TEMPORAL_CLOUD_API_KEY="api-key-here"

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

# Operation tools
test_get_operation() {
    init_server
    echo '{"jsonrpc": "2.0", "id": 9, "method": "tools/call", "params": {"name": "temporal_get_async_operation", "arguments": {"operation_id": "test-operation-id"}}}'
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
    echo "  test_get_operation - Test temporal_get_async_operation tool"
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
    "test_get_operation")
        test_get_operation
        ;;
    "all")
        run_all_tests
        ;;
    *)
        show_usage
        ;;
esac