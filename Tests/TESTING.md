# JMDT Geth Facade Testing Guide

This document describes how to test the JMDT Geth Facade API and WebSocket functionality.

## Prerequisites

### Required Dependencies
- `curl` - For HTTP requests
- `jq` - For JSON parsing (optional but recommended)

### Optional Dependencies
- `python3` with `websockets` library - For comprehensive WebSocket testing
- `node` - For WebSocket testing with JavaScript
- `wscat` - For interactive WebSocket testing

### Installing Dependencies

#### macOS (using Homebrew)
```bash
brew install curl jq python3 node
pip3 install websockets
npm install -g wscat
```

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install curl jq python3 python3-pip nodejs npm
pip3 install websockets
sudo npm install -g wscat
```

## Test Scripts

### 1. Basic API Tests (`test-basic.sh`)

Tests all JSON-RPC APIs without WebSocket dependencies.

**Usage:**
```bash
./test-basic.sh
```

**What it tests:**
- ✅ Health and readiness endpoints
- ✅ Basic blockchain info (chainId, blockNumber, etc.)
- ✅ Block operations (getBlockByNumber, getBlockByHash, etc.)
- ✅ Account operations (getBalance, getCode, getStorageAt, etc.)
- ✅ Transaction operations (gasPrice, estimateGas, call, etc.)
- ✅ Network operations (peerCount, listening, syncing)
- ✅ Mining operations (mining, hashrate)
- ✅ Uncle operations (uncle count, get uncle)
- ✅ Log operations (getLogs)
- ✅ Error handling

### 2. Comprehensive Tests (`test-apis.sh`)

Tests all APIs including WebSocket functionality.

**Usage:**
```bash
./test-apis.sh
```

**Requirements:**
- Python3 with websockets library
- All dependencies from basic tests

**What it tests:**
- ✅ Everything from basic tests
- ✅ WebSocket connection
- ✅ WebSocket RPC calls
- ✅ WebSocket subscriptions (newHeads, logs, pending transactions)
- ✅ WebSocket unsubscription
- ✅ Performance testing (concurrent requests)

### 3. WebSocket Tests (`test-websocket.sh`)

Focused WebSocket testing with Node.js.

**Usage:**
```bash
./test-websocket.sh
```

**Requirements:**
- Node.js
- Optional: wscat for interactive testing

## Running Tests

### 1. Start the Server

```bash
# Start the server in the background
./jmdt-geth-facade -http ":8545" -ws ":8546" &

# Or start in foreground (Ctrl+C to stop)
./jmdt-geth-facade -http ":8545" -ws ":8546"
```

### 2. Run Tests

```bash
# Basic API tests (no dependencies)
./test-basic.sh

# Comprehensive tests (requires Python)
./test-apis.sh

# WebSocket tests (requires Node.js)
./test-websocket.sh
```

### 3. Stop the Server

```bash
# If running in background
pkill -f jmdt-geth-facade

# Or find and kill the process
ps aux | grep jmdt-geth-facade
kill <PID>
```

## Manual Testing

### HTTP Endpoints

#### Health Checks
```bash
# Health check
curl http://localhost:8545/health

# Readiness check
curl http://localhost:8545/ready
```

#### JSON-RPC Examples
```bash
# Get chain ID
curl -X POST http://localhost:8545/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'

# Get block number
curl -X POST http://localhost:8545/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":2}'

# Get balance
curl -X POST http://localhost:8545/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB","latest"],"id":3}'

# Get block by number
curl -X POST http://localhost:8545/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest",false],"id":4}'
```

### WebSocket Examples

#### Using wscat
```bash
# Connect to WebSocket
wscat -c ws://localhost:8546

# Send RPC request
{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}

# Subscribe to new blocks
{"jsonrpc":"2.0","method":"eth_subscribe","params":["newHeads"],"id":2}

# Unsubscribe
{"jsonrpc":"2.0","method":"eth_unsubscribe","params":["subscription_id"],"id":3}
```

#### Using Python
```python
import asyncio
import websockets
import json

async def test_websocket():
    uri = "ws://localhost:8546"
    async with websockets.connect(uri) as websocket:
        # Send RPC request
        request = {
            "jsonrpc": "2.0",
            "method": "eth_chainId",
            "params": [],
            "id": 1
        }
        await websocket.send(json.dumps(request))
        
        # Receive response
        response = await websocket.recv()
        print(json.loads(response))

asyncio.run(test_websocket())
```

## Expected Results

### Successful Responses
- ✅ Health checks return `{"status":"healthy","timestamp":...}`
- ✅ JSON-RPC calls return proper JSON-RPC 2.0 responses
- ✅ WebSocket connections establish successfully
- ✅ Subscriptions work and return subscription IDs

### Expected Errors (Normal for Mock Backend)
- ⚠️ Some transaction operations may return empty results
- ⚠️ Uncle operations may return "no uncles available"
- ⚠️ Some storage operations may return errors for invalid keys

### Error Codes
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32000`: Server error
- `-32700`: Parse error

## Troubleshooting

### Common Issues

1. **Server not starting**
   - Check if ports 8545/8546 are available
   - Ensure the binary is executable: `chmod +x jmdt-geth-facade`

2. **Connection refused**
   - Verify server is running: `ps aux | grep jmdt-geth-facade`
   - Check if server is listening: `netstat -an | grep 8545`

3. **WebSocket tests failing**
   - Ensure Python websockets library is installed: `pip3 install websockets`
   - Check if Node.js is available: `node --version`

4. **JSON parsing errors**
   - Install jq: `brew install jq` or `apt-get install jq`
   - Or use `python3 -m json.tool` instead of jq

### Debug Mode

To see detailed server logs, run the server in foreground:
```bash
./jmdt-geth-facade -http ":8545" -ws ":8546"
```

## Performance Testing

The comprehensive test script includes basic performance testing with concurrent requests. For more advanced performance testing, consider:

1. **Load testing with Apache Bench:**
```bash
ab -n 1000 -c 10 -p request.json -T application/json http://localhost:8545/
```

2. **WebSocket load testing with artillery:**
```bash
npm install -g artillery
artillery quick --count 10 --num 10 ws://localhost:8546
```

## Contributing

When adding new APIs or features:

1. Add tests to the appropriate test script
2. Update this documentation
3. Ensure tests pass with both mock and real backends
4. Add error handling tests for new functionality

## Support

For issues or questions:
1. Check the server logs for error messages
2. Verify all dependencies are installed
3. Test with the basic script first
4. Check the GitHub issues for known problems
