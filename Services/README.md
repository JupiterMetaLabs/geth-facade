# Services

This folder contains all service implementations and business logic for the JMDT Geth Facade.

## Files

### `handlers.go`
JSON-RPC method handlers that process incoming requests and return responses:

- **Core Methods**: `eth_chainId`, `eth_blockNumber`, `net_version`, `web3_clientVersion`
- **Block Methods**: `eth_getBlockByNumber`, `eth_getBlockByHash`, `eth_getBlockTransactionCountBy*`
- **Account Methods**: `eth_getBalance`, `eth_getCode`, `eth_getStorageAt`, `eth_getTransactionCount`
- **Transaction Methods**: `eth_gasPrice`, `eth_estimateGas`, `eth_call`, `eth_sendRawTransaction`
- **Transaction Info**: `eth_getTransactionByHash`, `eth_getTransactionReceipt`
- **Network Methods**: `net_peerCount`, `net_listening`, `eth_syncing`
- **Mining Methods**: `eth_mining`, `eth_hashrate`
- **Uncle Methods**: `eth_getUncleCountBy*`, `eth_getUncleBy*`
- **Log Methods**: `eth_getLogs`

### `http_server.go`
HTTP server implementation using Gin framework:

- **Gin Integration**: High-performance HTTP server with middleware
- **CORS Support**: Cross-origin resource sharing configuration
- **Health Endpoints**: `/health` and `/ready` for monitoring
- **JSON-RPC Endpoint**: Main API endpoint at `/`
- **Error Handling**: Proper HTTP status codes and error responses

### `ws_server.go`
WebSocket server for real-time subscriptions:

- **Subscription Support**: `eth_subscribe` and `eth_unsubscribe`
- **Real-time Events**: New block headers, logs, pending transactions
- **Connection Management**: WebSocket upgrade and connection handling
- **Message Forwarding**: Efficient message routing to subscribers

### `memory.go`
In-memory mock backend implementation:

- **Mock Data**: Provides realistic test data for development
- **Full Interface**: Implements all Backend interface methods
- **Subscription Support**: Mock real-time event generation
- **Configurable**: Easy to modify for different test scenarios

### `facade.go`
Main facade service that orchestrates all components:

- **Server Management**: HTTP and WebSocket server coordination
- **Backend Integration**: Connects to various blockchain backends
- **Configuration**: Handles server configuration and startup
- **Health Checks**: Backend connectivity verification

### `doc.go`
Package documentation and examples.

## Key Features

- **Gin Framework**: High-performance HTTP server with middleware
- **WebSocket Support**: Real-time subscriptions for blockchain events
- **Mock Backend**: Complete in-memory implementation for testing
- **Error Handling**: Comprehensive error handling and logging
- **Health Monitoring**: Built-in health and readiness checks

## Architecture

```
Client Request → HTTP/WS Server → Handlers → Backend → Response
```

## Comments

The code includes standardized comments:
- `//debugging`: Debug-related code and logging
- `//future`: Planned features like rate limiting, caching
- `//test`: Test-related code and mock data
- `//conversions`: Data type conversions between formats
