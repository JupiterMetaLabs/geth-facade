#!/bin/bash

# JMDT Geth Facade API Test Script
# This script tests all JSON-RPC APIs and WebSocket functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
HTTP_URL="http://localhost:8545"
WS_URL="ws://localhost:8546"
TEST_ADDRESS="0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB"
TEST_ADDRESS_2="0xA2902C128D42A64F371457b82BB6aBb05B9b8bf1"

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

# Function to make JSON-RPC requests
make_request() {
    local method="$1"
    local params="$2"
    local id="$3"
    
    curl -s -X POST "$HTTP_URL/" \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"$method\",\"params\":$params,\"id\":$id}" \
        | jq .
}

# Function to test basic connectivity
test_connectivity() {
    print_header "Testing Basic Connectivity"
    
    print_test "Health Check"
    if curl -s "$HTTP_URL/health" | jq -e '.status == "healthy"' > /dev/null; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        return 1
    fi
    
    print_test "Ready Check"
    if curl -s "$HTTP_URL/ready" | jq -e '.status == "ready"' > /dev/null; then
        print_success "Ready check passed"
    else
        print_error "Ready check failed"
        return 1
    fi
}

# Function to test basic blockchain info
test_basic_info() {
    print_header "Testing Basic Blockchain Info"
    
    print_test "web3_clientVersion"
    make_request "web3_clientVersion" "[]" 1
    
    print_test "net_version"
    make_request "net_version" "[]" 2
    
    print_test "eth_chainId"
    make_request "eth_chainId" "[]" 3
    
    print_test "eth_blockNumber"
    make_request "eth_blockNumber" "[]" 4
}

# Function to test block operations
test_block_operations() {
    print_header "Testing Block Operations"
    
    print_test "eth_getBlockByNumber (latest, false)"
    make_request "eth_getBlockByNumber" "[\"latest\", false]" 5
    
    print_test "eth_getBlockByNumber (latest, true)"
    make_request "eth_getBlockByNumber" "[\"latest\", true]" 6
    
    print_test "eth_getBlockByNumber (0x0, false)"
    make_request "eth_getBlockByNumber" "[\"0x0\", false]" 7
    
    print_test "eth_getBlockByHash"
    make_request "eth_getBlockByHash" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\", false]" 8
    
    print_test "eth_getBlockTransactionCountByNumber"
    make_request "eth_getBlockTransactionCountByNumber" "[\"latest\"]" 9
    
    print_test "eth_getBlockTransactionCountByHash"
    make_request "eth_getBlockTransactionCountByHash" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\"]" 10
}

# Function to test account operations
test_account_operations() {
    print_header "Testing Account Operations"
    
    print_test "eth_getBalance"
    make_request "eth_getBalance" "[\"$TEST_ADDRESS\", \"latest\"]" 11
    
    print_test "eth_getBalance (second address)"
    make_request "eth_getBalance" "[\"$TEST_ADDRESS_2\", \"latest\"]" 12
    
    print_test "eth_getCode"
    make_request "eth_getCode" "[\"$TEST_ADDRESS\", \"latest\"]" 13
    
    print_test "eth_getStorageAt"
    make_request "eth_getStorageAt" "[\"$TEST_ADDRESS\", \"0x0\", \"latest\"]" 14
    
    print_test "eth_getTransactionCount"
    make_request "eth_getTransactionCount" "[\"$TEST_ADDRESS\", \"latest\"]" 15
}

# Function to test transaction operations
test_transaction_operations() {
    print_header "Testing Transaction Operations"
    
    print_test "eth_gasPrice"
    make_request "eth_gasPrice" "[]" 16
    
    print_test "eth_estimateGas"
    make_request "eth_estimateGas" "[{\"from\":\"$TEST_ADDRESS\",\"to\":\"$TEST_ADDRESS_2\",\"value\":\"0x1\"}]" 17
    
    print_test "eth_call"
    make_request "eth_call" "[{\"from\":\"$TEST_ADDRESS\",\"to\":\"$TEST_ADDRESS_2\",\"value\":\"0x1\"}, \"latest\"]" 18
    
    print_test "eth_getTransactionByHash"
    make_request "eth_getTransactionByHash" "[\"0xdeadbeef\"]" 19
    
    print_test "eth_getTransactionByBlockNumberAndIndex"
    make_request "eth_getTransactionByBlockNumberAndIndex" "[\"latest\", \"0x0\"]" 20
    
    print_test "eth_getTransactionByBlockHashAndIndex"
    make_request "eth_getTransactionByBlockHashAndIndex" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\", \"0x0\"]" 21
    
    print_test "eth_getTransactionReceipt"
    make_request "eth_getTransactionReceipt" "[\"0xdeadbeef\"]" 22
}

# Function to test network operations
test_network_operations() {
    print_header "Testing Network Operations"
    
    print_test "net_peerCount"
    make_request "net_peerCount" "[]" 23
    
    print_test "net_listening"
    make_request "net_listening" "[]" 24
    
    print_test "eth_syncing"
    make_request "eth_syncing" "[]" 25
}

# Function to test mining operations
test_mining_operations() {
    print_header "Testing Mining Operations"
    
    print_test "eth_mining"
    make_request "eth_mining" "[]" 26
    
    print_test "eth_hashrate"
    make_request "eth_hashrate" "[]" 27
}

# Function to test uncle operations
test_uncle_operations() {
    print_header "Testing Uncle Operations"
    
    print_test "eth_getUncleCountByBlockNumber"
    make_request "eth_getUncleCountByBlockNumber" "[\"latest\"]" 28
    
    print_test "eth_getUncleCountByBlockHash"
    make_request "eth_getUncleCountByBlockHash" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\"]" 29
    
    print_test "eth_getUncleByBlockNumberAndIndex"
    make_request "eth_getUncleByBlockNumberAndIndex" "[\"latest\", \"0x0\"]" 30
    
    print_test "eth_getUncleByBlockHashAndIndex"
    make_request "eth_getUncleByBlockHashAndIndex" "[\"0x0000000000000000000000000000000000000000000000000000000000000000\", \"0x0\"]" 31
}

# Function to test log operations
test_log_operations() {
    print_header "Testing Log Operations"
    
    print_test "eth_getLogs"
    make_request "eth_getLogs" "[{\"fromBlock\":\"latest\",\"toBlock\":\"latest\",\"address\":[\"$TEST_ADDRESS\"]}]" 32
}

# Function to test WebSocket functionality
test_websocket() {
    print_header "Testing WebSocket Functionality"
    
    print_test "WebSocket Connection Test"
    
    # Create a temporary Python script for WebSocket testing
    cat > /tmp/ws_test.py << 'EOF'
#!/usr/bin/env python3
import asyncio
import websockets
import json
import sys

async def test_websocket():
    uri = "ws://localhost:8546"
    try:
        async with websockets.connect(uri) as websocket:
            print("✅ WebSocket connection established")
            
            # Test basic RPC call
            request = {
                "jsonrpc": "2.0",
                "method": "eth_chainId",
                "params": [],
                "id": 1
            }
            
            await websocket.send(json.dumps(request))
            response = await websocket.recv()
            result = json.loads(response)
            
            if "result" in result:
                print(f"✅ RPC call successful: {result['result']}")
            else:
                print(f"❌ RPC call failed: {result}")
                return False
            
            # Test subscription
            sub_request = {
                "jsonrpc": "2.0",
                "method": "eth_subscribe",
                "params": ["newHeads"],
                "id": 2
            }
            
            await websocket.send(json.dumps(sub_request))
            sub_response = await websocket.recv()
            sub_result = json.loads(sub_response)
            
            if "result" in sub_result:
                print(f"✅ Subscription successful: {sub_result['result']}")
                subscription_id = sub_result['result']
                
                # Wait for a subscription message
                try:
                    await asyncio.wait_for(websocket.recv(), timeout=10.0)
                    print("✅ Received subscription message")
                except asyncio.TimeoutError:
                    print("⚠️  No subscription message received within 10 seconds")
                
                # Unsubscribe
                unsub_request = {
                    "jsonrpc": "2.0",
                    "method": "eth_unsubscribe",
                    "params": [subscription_id],
                    "id": 3
                }
                
                await websocket.send(json.dumps(unsub_request))
                unsub_response = await websocket.recv()
                unsub_result = json.loads(unsub_response)
                
                if "result" in unsub_result and unsub_result["result"]:
                    print("✅ Unsubscription successful")
                else:
                    print("❌ Unsubscription failed")
                    return False
            else:
                print(f"❌ Subscription failed: {sub_result}")
                return False
            
            return True
            
    except Exception as e:
        print(f"❌ WebSocket test failed: {e}")
        return False

if __name__ == "__main__":
    result = asyncio.run(test_websocket())
    sys.exit(0 if result else 1)
EOF
    
    if python3 /tmp/ws_test.py; then
        print_success "WebSocket functionality test passed"
    else
        print_error "WebSocket functionality test failed"
    fi
    
    # Clean up
    rm -f /tmp/ws_test.py
}

# Function to test error handling
test_error_handling() {
    print_header "Testing Error Handling"
    
    print_test "Invalid method"
    make_request "invalid_method" "[]" 99
    
    print_test "Invalid parameters"
    make_request "eth_getBalance" "[]" 98
    
    print_test "Invalid block tag"
    make_request "eth_getBlockByNumber" "[\"invalid\", false]" 97
}

# Function to run performance test
test_performance() {
    print_header "Testing Performance"
    
    print_test "Concurrent requests test"
    
    # Create multiple concurrent requests
    for i in {1..10}; do
        make_request "eth_blockNumber" "[]" $((100 + i)) &
    done
    
    # Wait for all background jobs to complete
    wait
    
    print_success "Concurrent requests completed"
}

# Main execution
main() {
    echo -e "${GREEN}🚀 Starting JMDT Geth Facade API Tests${NC}"
    echo -e "${GREEN}HTTP URL: $HTTP_URL${NC}"
    echo -e "${GREEN}WebSocket URL: $WS_URL${NC}"
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        print_error "jq is required but not installed. Please install jq first."
        exit 1
    fi
    
    # Check if Python3 is installed (for WebSocket tests)
    if ! command -v python3 &> /dev/null; then
        print_error "python3 is required for WebSocket tests but not installed."
        exit 1
    fi
    
    # Check if websockets library is available
    if ! python3 -c "import websockets" 2>/dev/null; then
        print_error "websockets library is required. Install with: pip3 install websockets"
        exit 1
    fi
    
    # Run all tests
    test_connectivity
    test_basic_info
    test_block_operations
    test_account_operations
    test_transaction_operations
    test_network_operations
    test_mining_operations
    test_uncle_operations
    test_log_operations
    test_websocket
    test_error_handling
    test_performance
    
    print_header "All Tests Completed"
    echo -e "${GREEN}🎉 All API tests have been executed!${NC}"
    echo -e "${YELLOW}Note: Some tests may show expected errors (like missing transactions) since this is a mock backend.${NC}"
}

# Run main function
main "$@"
