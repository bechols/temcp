#!/bin/bash

export TEMPORAL_CLOUD_API_KEY="eyJhbGciOiJFUzI1NiIsICJraWQiOiJXdnR3YUEifQ.eyJhY2NvdW50X2lkIjoiYTJkZDYiLCAiYXVkIjpbInRlbXBvcmFsLmlvIl0sICJleHAiOjE3NTM2NzYzMDksICJpc3MiOiJ0ZW1wb3JhbC5pbyIsICJqdGkiOiJjdEdMa0RDOWh6SzNjRTFRMHphd2d2ckVQRGRKTmlwbCIsICJrZXlfaWQiOiJjdEdMa0RDOWh6SzNjRTFRMHphd2d2ckVQRGRKTmlwbCIsICJzdWIiOiI1YWRmOWI5ODllOGU0MjkyYmU2NzFiNjIxMjJkYjQ2YSJ9.aGmcKBmq9vi16iScnmnfz3Ui76EmKuuyRWMbMbyBhYbCJ7wS3BUP3btEjkidSuA0SqxiWU5-LhImYzbmc-gZBw"

echo "Testing MCP Server Tools..."

# Initialize the server
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}'

# List available tools
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'

# Test get user tool (with invalid ID to test error handling)
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "temporal_get_user", "arguments": {"user_id": "test-user-id"}}}'

# Test list users tool
echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "temporal_list_users", "arguments": {"page_size": 10}}}'