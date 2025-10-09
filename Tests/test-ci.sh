#!/bin/bash

# CI/CD Test Script
# This script runs essential tests for continuous integration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
HTTP_URL="http://localhost:8545"
SERVER_PID=""

# Helper functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_test() {
    echo -e "\n${YELLOW}Testing: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to start server
start_server() {
    print_test "Starting server"
    ./jmdt-geth-facade -http ":8545" -ws ":8546" &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 3
    
    # Check if server is running
    if kill -0 $SERVER_PID 2>/dev/null; then
        print_success "Server started (PID: $SERVER_PID)"
    else
        print_error "Failed to start server"
        exit 1
    fi
}

# Function to stop server
stop_server() {
    if [ ! -z "$SERVER_PID" ]; then
        print_test "Stopping server"
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
        print_success "Server stopped"
    fi
}

# Function to make JSON-RPC requests
make_request() {
    local method="$1"
    local params="$2"
    local id="$3"
    
    curl -s -X POST "$HTTP_URL/" \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"$method\",\"params\":$params,\"id\":$id}" \
        | jq -e '.result != null or .error != null' > /dev/null
}

# Function to test essential functionality
test_essential() {
    print_header "Testing Essential Functionality"
    
    # Test health endpoints
    print_test "Health endpoints"
    if curl -s "$HTTP_URL/health" | grep -q "healthy"; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        return 1
    fi
    
    if curl -s "$HTTP_URL/ready" | grep -q "ready"; then
        print_success "Ready check passed"
    else
        print_error "Ready check failed"
        return 1
    fi
    
    # Test basic RPC methods
    print_test "Basic RPC methods"
    
    if make_request "eth_chainId" "[]" 1; then
        print_success "eth_chainId works"
    else
        print_error "eth_chainId failed"
        return 1
    fi
    
    if make_request "eth_blockNumber" "[]" 2; then
        print_success "eth_blockNumber works"
    else
        print_error "eth_blockNumber failed"
        return 1
    fi
    
    if make_request "net_version" "[]" 3; then
        print_success "net_version works"
    else
        print_error "net_version failed"
        return 1
    fi
    
    if make_request "web3_clientVersion" "[]" 4; then
        print_success "web3_clientVersion works"
    else
        print_error "web3_clientVersion failed"
        return 1
    fi
}

# Function to test block operations
test_blocks() {
    print_header "Testing Block Operations"
    
    print_test "Block operations"
    
    if make_request "eth_getBlockByNumber" "[\"latest\", false]" 5; then
        print_success "eth_getBlockByNumber works"
    else
        print_error "eth_getBlockByNumber failed"
        return 1
    fi
    
    if make_request "eth_getBlockByHash" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\", false]" 6; then
        print_success "eth_getBlockByHash works"
    else
        print_error "eth_getBlockByHash failed"
        return 1
    fi
}

# Function to test account operations
test_accounts() {
    print_header "Testing Account Operations"
    
    print_test "Account operations"
    
    if make_request "eth_getBalance" "[\"0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB\", \"latest\"]" 7; then
        print_success "eth_getBalance works"
    else
        print_error "eth_getBalance failed"
        return 1
    fi
    
    if make_request "eth_getCode" "[\"0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB\", \"latest\"]" 8; then
        print_success "eth_getCode works"
    else
        print_error "eth_getCode failed"
        return 1
    fi
    
    if make_request "eth_getTransactionCount" "[\"0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB\", \"latest\"]" 9; then
        print_success "eth_getTransactionCount works"
    else
        print_error "eth_getTransactionCount failed"
        return 1
    fi
}

# Function to test error handling
test_errors() {
    print_header "Testing Error Handling"
    
    print_test "Error handling"
    
    # Test invalid method
    response=$(curl -s -X POST "$HTTP_URL/" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"invalid_method","params":[],"id":99}')
    
    if echo "$response" | jq -e '.error.code == -32601' > /dev/null; then
        print_success "Invalid method error handling works"
    else
        print_error "Invalid method error handling failed"
        return 1
    fi
    
    # Test invalid params
    response=$(curl -s -X POST "$HTTP_URL/" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":[],"id":98}')
    
    if echo "$response" | jq -e '.error.code == -32602' > /dev/null; then
        print_success "Invalid params error handling works"
    else
        print_error "Invalid params error handling failed"
        return 1
    fi
}

# Function to test performance
test_performance() {
    print_header "Testing Performance"
    
    print_test "Concurrent requests"
    
    # Run 10 concurrent requests
    for i in {1..10}; do
        make_request "eth_blockNumber" "[]" $((100 + i)) &
    done
    
    # Wait for all background jobs to complete
    wait
    
    print_success "Concurrent requests completed"
}

# Cleanup function
cleanup() {
    stop_server
}

# Set up trap for cleanup
trap cleanup EXIT

# Main execution
main() {
    echo -e "${GREEN}🚀 Starting CI/CD Tests${NC}"
    
    # Check if jq is available
    if ! command -v jq &> /dev/null; then
        print_error "jq is required but not installed. Please install jq first."
        exit 1
    fi
    
    # Start server
    start_server
    
    # Run tests
    test_essential
    test_blocks
    test_accounts
    test_errors
    test_performance
    
    print_header "All CI/CD Tests Passed"
    echo -e "${GREEN}🎉 All essential tests have passed!${NC}"
}

# Run main function
main "$@"
