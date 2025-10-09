#!/bin/bash

# WebSocket Test Script using wscat (if available)
# This script tests WebSocket functionality using wscat

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
WS_URL="ws://localhost:8546"

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

# Check if wscat is available
check_wscat() {
    if ! command -v wscat &> /dev/null; then
        print_error "wscat is not installed. Please install it with:"
        echo "  npm install -g wscat"
        echo "  or"
        echo "  brew install wscat"
        return 1
    fi
    return 0
}

# Test WebSocket connection
test_websocket_connection() {
    print_header "Testing WebSocket Connection"
    
    print_test "Basic WebSocket Connection"
    
    # Create a temporary script for wscat
    cat > /tmp/ws_test.js << 'EOF'
// WebSocket test script for wscat
const WebSocket = require('ws');

const ws = new WebSocket('ws://localhost:8546');

ws.on('open', function open() {
    console.log('✅ WebSocket connection established');
    
    // Test basic RPC call
    const request = {
        jsonrpc: "2.0",
        method: "eth_chainId",
        params: [],
        id: 1
    };
    
    ws.send(JSON.stringify(request));
});

ws.on('message', function message(data) {
    const response = JSON.parse(data);
    console.log('📥 Received:', JSON.stringify(response, null, 2));
    
    if (response.result) {
        console.log('✅ RPC call successful');
    } else if (response.error) {
        console.log('❌ RPC call failed:', response.error);
    }
    
    // Test subscription
    const subRequest = {
        jsonrpc: "2.0",
        method: "eth_subscribe",
        params: ["newHeads"],
        id: 2
    };
    
    console.log('📤 Sending subscription request...');
    ws.send(JSON.stringify(subRequest));
    
    // Close after a short delay
    setTimeout(() => {
        ws.close();
    }, 2000);
});

ws.on('close', function close() {
    console.log('🔌 WebSocket connection closed');
});

ws.on('error', function error(err) {
    console.log('❌ WebSocket error:', err.message);
});
EOF
    
    if node /tmp/ws_test.js; then
        print_success "WebSocket test completed"
    else
        print_error "WebSocket test failed"
    fi
    
    # Clean up
    rm -f /tmp/ws_test.js
}

# Test with wscat if available
test_with_wscat() {
    print_header "Testing with wscat (Interactive)"
    
    print_test "Starting wscat session"
    echo -e "${YELLOW}You can now test WebSocket manually. Try these commands:${NC}"
    echo -e "${YELLOW}1. Basic RPC: {\"jsonrpc\":\"2.0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}${NC}"
    echo -e "${YELLOW}2. Subscribe: {\"jsonrpc\":\"2.0\",\"method\":\"eth_subscribe\",\"params\":[\"newHeads\"],\"id\":2}${NC}"
    echo -e "${YELLOW}3. Unsubscribe: {\"jsonrpc\":\"2.0\",\"method\":\"eth_unsubscribe\",\"params\":[\"subscription_id\"],\"id\":3}${NC}"
    echo -e "${YELLOW}Press Ctrl+C to exit${NC}"
    
    wscat -c "$WS_URL"
}

# Main execution
main() {
    echo -e "${GREEN}🚀 Starting WebSocket Tests${NC}"
    echo -e "${GREEN}WebSocket URL: $WS_URL${NC}"
    
    # Check if Node.js is available
    if ! command -v node &> /dev/null; then
        print_error "Node.js is required for WebSocket testing but not installed."
        echo "Please install Node.js first."
        exit 1
    fi
    
    # Run tests
    test_websocket_connection
    
    # Ask if user wants to test with wscat
    if check_wscat; then
        echo -e "\n${YELLOW}Would you like to test with wscat interactively? (y/n)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            test_with_wscat
        fi
    fi
    
    print_header "WebSocket Tests Completed"
    echo -e "${GREEN}🎉 WebSocket tests have been executed!${NC}"
}

# Run main function
main "$@"
