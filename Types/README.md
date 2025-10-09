# Types

This folder contains all data structures and type definitions used throughout the JMDT Geth Facade.

## Files

### `backend.go`
Contains the core data structures that mirror the official Geth implementation:

- **Block**: Complete block structure with header, transactions, ommers, withdrawals, and blob gas fields
- **BlockHeader**: Detailed block header with all Ethereum fields including EIP-1559, EIP-4844, EIP-4895 support
- **Transaction**: Comprehensive transaction structure supporting all transaction types
- **Receipt**: Transaction receipt with logs and status information
- **Log**: Event log structure with topics and data
- **AccessList**: EIP-2930 access list support
- **AccessTuple**: Individual access list entry
- **Withdrawal**: EIP-4895 withdrawal structure
- **CallMsg**: Message structure for eth_call and eth_estimateGas
- **FilterQuery**: Log filtering parameters

### `types.go`
Contains JSON-RPC specific types:

- **Request**: Incoming JSON-RPC request structure
- **Response**: JSON-RPC response structure
- **Error**: JSON-RPC error structure
- **Subscription**: WebSocket subscription management

## Key Features

- **Geth Compatibility**: Structures match official Geth implementation
- **EIP Support**: Full support for EIP-1559, EIP-2930, EIP-4844, EIP-4895
- **Type Safety**: Consistent use of `[]byte` for hashes and addresses
- **JSON Tags**: Proper JSON serialization tags for all fields

## Usage

These types are used throughout the codebase to ensure consistency and compatibility with Ethereum tooling. All backend implementations must conform to these interfaces.

## Comments

The code includes standardized comments:
- `//debugging`: Debug-related code
- `//future`: Planned features or improvements
- `//test`: Test-related code
- `//conversions`: Data type conversions
