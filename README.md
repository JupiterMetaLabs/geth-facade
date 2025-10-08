## jmdt-geth-facade

A minimal, production-ready JSON-RPC/WebSocket façade that mirrors common geth endpoints (eth_*, net_*, web3_*). Wire the `backend.Backend` interface to your own node(s) to provide real data.

This project can be used both as a standalone application and as a Go package for integration into your own applications.

### Features
- HTTP JSON-RPC on :8545 (configurable)
- WebSocket JSON-RPC on :8546 (configurable, supports eth_subscribe / eth_unsubscribe)
- Methods implemented: web3_clientVersion, net_version, eth_chainId, eth_blockNumber, eth_getBlockByNumber, eth_getBalance, eth_call, eth_estimateGas, eth_gasPrice, eth_sendRawTransaction, eth_getTransactionByHash, eth_getTransactionReceipt, eth_getLogs


Note: The in-repo `backend/memory.go` includes hardcoded JMDT balances for testing. Replace it with your real adapter for production use.

### Requirements
- Go 1.22+ (1.23+ recommended on newer macOS)
- macOS note: if you hit a dyld LC_UUID error, use the external linker flags shown below.

## Usage as a Go Package

### Installation

```bash
go get github.com/saishibu/jmdt-geth-facade
```

## Step-by-Step Integration Guide

### Step 1: Create Your Backend Implementation

First, you need to implement the `backend.Backend` interface to connect to your blockchain node:

```go
package main

import (
    "context"
    "math/big"
    
    "github.com/saishibu/jmdt-geth-facade/backend"
)

type MyBlockchainBackend struct {
    // Your blockchain client (e.g., ethclient.Client, custom RPC client, etc.)
    client interface{} // Replace with your actual client type
}

// Implement all required methods
func (b *MyBlockchainBackend) ChainID(ctx context.Context) (*big.Int, error) {
    // Return your blockchain's chain ID
    return big.NewInt(1), nil // Example: Ethereum mainnet
}

func (b *MyBlockchainBackend) ClientVersion(ctx context.Context) (string, error) {
    // Return your client version
    return "my-blockchain-client/1.0.0", nil
}

func (b *MyBlockchainBackend) BlockNumber(ctx context.Context) (*big.Int, error) {
    // Query current block number from your node
    // Example: return b.client.BlockNumber(ctx)
    return big.NewInt(18000000), nil
}

func (b *MyBlockchainBackend) BlockByNumber(ctx context.Context, num *big.Int, fullTx bool) (*backend.Block, error) {
    // Fetch block by number from your node
    return &backend.Block{
        Number:     num,
        Hash:       "0x...", // Get from your node
        ParentHash: "0x...", // Get from your node
        Timestamp:  uint64(time.Now().Unix()),
    }, nil
}

func (b *MyBlockchainBackend) Balance(ctx context.Context, addr string, block *big.Int) (*big.Int, error) {
    // Query balance from your node
    // Example: return b.client.BalanceAt(ctx, common.HexToAddress(addr), block)
    return big.NewInt(1000000000000000000), nil // 1 ETH in wei
}

// ... implement all other required methods (Call, EstimateGas, GasPrice, etc.)
```

### Step 2: Create Your Server

```go
package main

import (
    "log"
    "math/big"
    
    jmdtgethfacade "github.com/saishibu/jmdt-geth-facade/pkg/jmdtgethfacade"
)

func main() {
    // Create your backend
    myBackend := &MyBlockchainBackend{
        // Initialize your blockchain client
    }
    
    // Configure the facade server
    config := jmdtgethfacade.Config{
        Backend:  myBackend,
        HTTPAddr: ":8545", // HTTP JSON-RPC port
        WSAddr:   ":8546", // WebSocket JSON-RPC port
    }
    
    // Create and start the server
    server := jmdtgethfacade.NewServer(config)
    
    log.Println("Starting blockchain facade server...")
    if err := server.Start(); err != nil {
        log.Fatal("Server error:", err)
    }
}
```

### Step 3: Test Your Integration

Test your facade with curl:

```bash
# Test basic connectivity
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545

# Test block number
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# Test balance query
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x1234567890123456789012345678901234567890","latest"],"id":1}' \
  http://localhost:8545
```

### Step 4: Add Health Checks (Optional)

```go
// Add health check endpoints to your HTTP server
http.HandleFunc("/health", server.HealthCheck)
http.HandleFunc("/ready", server.ReadyCheck)
```

### Step 5: Production Considerations

1. **Error Handling**: Implement proper error handling in your backend methods
2. **Logging**: Add structured logging for debugging
3. **Metrics**: Add metrics collection for monitoring
4. **Rate Limiting**: Consider adding rate limiting for production use
5. **Security**: Implement proper authentication/authorization if needed
6. **Configuration**: Use environment variables for configuration

### Step 6: Testing with Memory Backend

For development and testing, you can use the included memory backend:

```go
import "github.com/saishibu/jmdt-geth-facade/pkg/memorybackend"

// Use memory backend for testing
testBackend := memorybackend.NewMemoryBackend(big.NewInt(11155111))
config := jmdtgethfacade.Config{
    Backend:  testBackend,
    HTTPAddr: ":8545",
    WSAddr:   ":8546",
}
```

## Common Integration Patterns

### Pattern 1: Ethereum Client Integration

```go
import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
)

type EthereumBackend struct {
    client *ethclient.Client
}

func (e *EthereumBackend) Balance(ctx context.Context, addr string, block *big.Int) (*big.Int, error) {
    address := common.HexToAddress(addr)
    return e.client.BalanceAt(ctx, address, block)
}

func (e *EthereumBackend) BlockNumber(ctx context.Context) (*big.Int, error) {
    return e.client.BlockNumber(ctx)
}
```

### Pattern 2: Custom RPC Client

```go
type CustomRPCBackend struct {
    rpcURL string
    httpClient *http.Client
}

func (c *CustomRPCBackend) Call(ctx context.Context, msg backend.CallMsg, block *big.Int) ([]byte, error) {
    // Make RPC call to your custom blockchain
    payload := map[string]interface{}{
        "method": "eth_call",
        "params": []interface{}{msg, block},
    }
    
    // Send HTTP request to your RPC endpoint
    // ... implementation details
}
```

### Pattern 3: Multi-Chain Support

```go
type MultiChainBackend struct {
    chains map[int64]backend.Backend
}

func (m *MultiChainBackend) GetBackend(chainID *big.Int) backend.Backend {
    return m.chains[chainID.Int64()]
}
```

## Troubleshooting

### Common Issues

1. **"Method not found" errors**
   - Ensure all required backend methods are implemented
   - Check that your backend returns proper data types

2. **Connection timeouts**
   - Verify your blockchain node is running and accessible
   - Check network connectivity and firewall settings

3. **Invalid response format**
   - Ensure your backend returns data in the expected format
   - Check that hex values are properly formatted (e.g., "0x1234")

4. **WebSocket connection issues**
   - Verify WebSocket support in your backend
   - Check that subscription methods return proper channels

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
import "log"

// Add this to see all RPC requests/responses
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Quick Start

```go
package main

import (
    "log"
    "math/big"
    
    jmdtgethfacade "github.com/saishibu/jmdt-geth-facade/pkg/jmdtgethfacade"
    "github.com/saishibu/jmdt-geth-facade/backend"
)

func main() {
    // You must provide your own backend implementation
    var myBackend backend.Backend = // your implementation
    
    if err := jmdtgethfacade.QuickStart(myBackend); err != nil {
        log.Fatal(err)
    }
}
```

### Custom Configuration

```go
package main

import (
    "log"
    "math/big"
    
    jmdtgethfacade "github.com/saishibu/jmdt-geth-facade/pkg/jmdtgethfacade"
    "github.com/saishibu/jmdt-geth-facade/backend"
)

func main() {
    // You must provide your own backend implementation
    var myBackend backend.Backend = // your implementation
    
    // Configure server
    config := jmdtgethfacade.Config{
        Backend:  myBackend,
        HTTPAddr: ":8545",
        WSAddr:   ":8546",
    }
    
    // Create and start server
    server := jmdtgethfacade.NewServer(config)
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Custom Backend Implementation

To integrate with your own blockchain node, implement the `backend.Backend` interface:

```go
type MyBackend struct {
    // Your blockchain client
}

func (b *MyBackend) ChainID(ctx context.Context) (*big.Int, error) {
    // Return your chain ID
}

func (b *MyBackend) BlockNumber(ctx context.Context) (*big.Int, error) {
    // Return current block number from your node
}

// ... implement all other required methods
```

See the `examples/` directory for complete examples.

For testing and development, you can use the memory backend:
```go
import "github.com/saishibu/jmdt-geth-facade/pkg/memorybackend"

// Create a mock backend for testing
backend := memorybackend.NewMemoryBackend(big.NewInt(11155111))
```

### Package API

The main package provides:

- `jmdtgethfacade.Config` - Server configuration
- `jmdtgethfacade.NewServer(config)` - Create a new server
- `jmdtgethfacade.QuickStart(chainID)` - Start with default config
- `server.Start()` - Start both HTTP and WebSocket servers
- `server.StartHTTP()` - Start only HTTP server
- `server.StartWS()` - Start only WebSocket server
- `server.HealthCheck()` - Health check endpoint
- `server.ReadyCheck()` - Readiness check endpoint

## Examples

The repository includes several examples to help you get started:

### Simple Example (with memory backend for testing)
```bash
cd examples/simple
go run main.go
```

### Custom Backend Example
```bash
cd examples/custom-backend
go run main.go
```

### Standalone Application
```bash
go build -o jmdt-geth-facade .
./jmdt-geth-facade -chainid 11155111
```

## Standalone Application Usage

### Build
```bash
go mod tidy
# Recommended on macOS (and generally safe everywhere):
go build -ldflags='-linkmode=external -w -s' -o jmdt-geth-facade .
```

### Run
The server accepts chain ID and port configuration via CLI flags.
```bash
# Basic run with default ports (8545, 8546)
./jmdt-geth-facade -chainid 11155111

# Custom ports
./jmdt-geth-facade -chainid 11155111 -http :8547 -ws :8548

# Hex chain ID
./jmdt-geth-facade -chainid 0xaa36a7
```

- HTTP endpoint: http://localhost:8545 (or custom port)
- WS endpoint: ws://localhost:8546 (or custom port)

### Supported RPCs
- web3_clientVersion
- net_version
- eth_chainId
- eth_blockNumber
- eth_getBlockByNumber
- eth_getBalance
- eth_call
- eth_estimateGas
- eth_gasPrice
- eth_sendRawTransaction
- eth_getTransactionByHash
- eth_getTransactionReceipt
- eth_getLogs

WebSocket subscriptions:
- eth_subscribe (newHeads, logs, newPendingTransactions)
- eth_unsubscribe

### Quick tests with curl
Replace `localhost:8545` if running remotely.

- web3_clientVersion
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' \
  http://localhost:8545
```

- eth_chainId
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545
```

- net_version
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
  http://localhost:8545
```

- eth_blockNumber
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

- eth_getBalance (returns hardcoded JMDT balances or 0x0)
```bash
# Test hardcoded JMDT balances
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB","latest"],"id":1}' \
  http://localhost:8545

# Test another hardcoded address
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xA2902C128D42A64F371457b82BB6aBb05B9b8bf1","latest"],"id":1}' \
  http://localhost:8545

# Test random address (returns 0x0)
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x1234567890123456789012345678901234567890","latest"],"id":1}' \
  http://localhost:8545
```

- eth_gasPrice
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}' \
  http://localhost:8545
```

- eth_estimateGas (stub returns 0x5208 for simple tx)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{}],"id":1}' \
  http://localhost:8545
```

- eth_sendRawTransaction (stub returns 0xdeadbeef)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0x01"],"id":1}' \
  http://localhost:8545
```

- eth_getTransactionByHash (stub)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xdeadbeef"],"id":1}' \
  http://localhost:8545
```

- eth_getTransactionReceipt (stub)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xdeadbeef"],"id":1}' \
  http://localhost:8545
```

- eth_getLogs (returns empty with memory backend)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getLogs","params":[{}],"id":1}' \
  http://localhost:8545
```

- eth_getBlockByNumber (structure is stubbed in memory backend)
```bash
curl -s -X POST -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}' \
  http://localhost:8545
```

### WebSocket subscription test
Using npx wscat:
```bash
npx wscat -c ws://localhost:8546
# Then send this message in the wscat prompt:
{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["newHeads"]}
```
You should receive periodic newHeads notifications from the memory backend.

### Next steps
- Replace `backend/memory.go` with your real `Backend` implementation.
- Map your node’s APIs into Ethereum-compatible responses (hex quantities, receipt fields, logs/topics).
- Harden WS server with ping/pong and timeouts if deploying behind NATs/LBs.
