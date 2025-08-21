## jmdt-geth-facade

A minimal, production-ready JSON-RPC/WebSocket façade that mirrors common geth endpoints (eth_*, net_*, web3_*). Wire the `backend.Backend` interface to your own node(s) to provide real data.

### Features
- HTTP JSON-RPC on :8545
- WebSocket JSON-RPC on :8546 (supports eth_subscribe / eth_unsubscribe)
- Methods implemented: web3_clientVersion, net_version, eth_chainId, eth_blockNumber, eth_getBlockByNumber, eth_getBalance, eth_call, eth_estimateGas, eth_gasPrice, eth_sendRawTransaction, eth_getTransactionByHash, eth_getTransactionReceipt, eth_getLogs

Note: The in-repo `backend/memory.go` is a stub for smoke testing. Replace it with your real adapter.

### Requirements
- Go 1.22+ (1.23+ recommended on newer macOS)
- macOS note: if you hit a dyld LC_UUID error, use the external linker flags shown below.

### Build
```bash
cd rpc-facade
go mod tidy
# Recommended on macOS (and generally safe everywhere):
go build -ldflags='-linkmode=external -w -s' -o jmdt-geth-facade .
```

### Run
The server accepts a chain ID via a CLI flag (decimal or hex).
```bash
# Decimal chain ID
./jmdt-geth-facade -chainid 11155111

# Hex chain ID
./jmdt-geth-facade -chainid 0xaa36a7
```

- HTTP endpoint: http://localhost:8545
- WS endpoint: ws://localhost:8546

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
