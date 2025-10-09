# Utils

This folder contains utility functions and helper code used throughout the JMDT Geth Facade.

## Purpose

The Utils folder provides reusable utility functions that support the main application functionality:

- **Data Conversion**: Helper functions for converting between data types
- **Validation**: Input validation and sanitization utilities
- **Formatting**: Data formatting and serialization helpers
- **Common Operations**: Frequently used operations across the codebase

## Current Status

Currently, utility functions are embedded within service files. This folder is prepared for extracting and centralizing common utilities such as:

- **Hex Conversion**: Hexadecimal string to byte array conversions
- **Address Validation**: Ethereum address format validation
- **Hash Operations**: Hash generation and validation
- **JSON Utilities**: JSON marshaling/unmarshaling helpers

## Planned Utilities

### `conversions.go`
- `HexToBytes()`: Convert hex strings to byte arrays
- `BytesToHex()`: Convert byte arrays to hex strings
- `AddressToBytes()`: Convert Ethereum addresses to bytes
- `BigIntToHex()`: Convert big.Int to hex string

### `validation.go`
- `ValidateAddress()`: Validate Ethereum address format
- `ValidateHash()`: Validate hash format and length
- `ValidateBlockTag()`: Validate block tag strings
- `ValidateHexString()`: Validate hex string format

### `formatting.go`
- `FormatResponse()`: Standardize JSON-RPC responses
- `FormatError()`: Standardize error responses
- `FormatBlock()`: Format block data for responses
- `FormatTransaction()`: Format transaction data for responses

## Usage

Utilities will be imported and used throughout the codebase to ensure consistency and reduce code duplication.

## Comments

All utility functions will include standardized comments:

- `//debugging`: Debug-related utility functions
- `//future`: Planned utility enhancements
- `//test`: Test-related utility functions
- `//conversions`: Data type conversion utilities
