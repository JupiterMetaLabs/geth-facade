# Tests

This folder contains all testing scripts and documentation for the JMDT Geth Facade.

## Files

### `test-basic.sh`
Basic API testing script that tests all JSON-RPC endpoints without external dependencies:

- **Health Checks**: Tests `/health` and `/ready` endpoints
- **Basic Info**: Tests `eth_chainId`, `eth_blockNumber`, `net_version`, `web3_clientVersion`
- **Block Operations**: Tests block retrieval and transaction count methods
- **Account Operations**: Tests balance, code, storage, and transaction count methods
- **Transaction Operations**: Tests gas price, estimation, calls, and transaction retrieval
- **Network Operations**: Tests peer count, listening status, and sync status
- **Mining Operations**: Tests mining status and hashrate
- **Uncle Operations**: Tests uncle-related methods
- **Log Operations**: Tests log filtering
- **Error Handling**: Tests invalid method and parameter handling

### `test-apis.sh`
Comprehensive testing script that includes WebSocket functionality:

- **All Basic Tests**: Includes everything from `test-basic.sh`
- **WebSocket Testing**: Tests WebSocket connections, RPC calls, and subscriptions
- **Subscription Testing**: Tests `newHeads`, `logs`, and `pendingTransactions` subscriptions
- **Performance Testing**: Concurrent request testing
- **Dependencies**: Requires Python3 with websockets library

### `test-websocket.sh`
Focused WebSocket testing using Node.js:

- **Connection Testing**: Basic WebSocket connectivity
- **RPC Testing**: JSON-RPC calls over WebSocket
- **Subscription Testing**: Real-time subscription functionality
- **Interactive Mode**: Optional wscat integration for manual testing

### `test-ci.sh`
CI/CD testing script for automated testing:

- **Essential Tests**: Core functionality verification
- **Error Handling**: Comprehensive error testing
- **Performance**: Basic performance validation
- **Automated**: Designed for continuous integration pipelines

### `TESTING.md`
Comprehensive testing documentation:

- **Setup Instructions**: Dependency installation and configuration
- **Usage Examples**: How to run each test script
- **Manual Testing**: Examples for manual API testing
- **Troubleshooting**: Common issues and solutions
- **Performance Testing**: Advanced testing techniques

## Usage

### Prerequisites
```bash
# Basic dependencies
curl, jq

# For comprehensive testing
python3, websockets library

# For WebSocket testing
node, wscat (optional)
```

### Running Tests
```bash
# Start the server
./jmdt-geth-facade -http ":8545" -ws ":8546" &

# Run basic tests
./Tests/test-basic.sh

# Run comprehensive tests
./Tests/test-apis.sh

# Run WebSocket tests
./Tests/test-websocket.sh

# Run CI tests
./Tests/test-ci.sh
```

## Test Coverage

The test suite covers:

- ✅ **All JSON-RPC Methods**: Complete API coverage
- ✅ **Error Scenarios**: Invalid inputs and error handling
- ✅ **WebSocket Functionality**: Real-time subscriptions
- ✅ **Performance**: Concurrent request handling
- ✅ **Health Monitoring**: Health and readiness checks

## Comments

Test scripts include standardized comments:

- `//debugging`: Debug output and verbose logging
- `//future`: Planned test enhancements
- `//test`: Test-specific functionality
- `//conversions`: Data conversion testing
