// Package memorybackend provides a mock implementation of the Types.Backend interface
// for testing and development purposes. This is NOT suitable for production use.
//
// This example shows how to implement the Types.Backend interface and can be used
// as a reference for creating your own backend implementations.
package Services

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/jupitermetalabs/geth-facade/Types"
)

// mem provides an in-memory mock backend implementation
// //test: Mock implementation for testing and development
// //debugging: Provides predictable test data
type mem struct {
	chainID *big.Int
	num     *big.Int
}

// NewMemoryBackend creates a new mock backend for testing and development.
// This implementation simulates a blockchain with:
// - Incrementing block numbers every 6 seconds
// - Mock balances for demonstration
// - Stub responses for all RPC methods
//
// WARNING: This is for testing only. Do not use in production.
func NewMemoryBackend(chainID *big.Int) Types.Backend {
	ci := big.NewInt(11155111)
	if chainID != nil {
		ci = new(big.Int).Set(chainID)
	}
	m := &mem{chainID: ci, num: big.NewInt(0)}
	go func() {
		t := time.NewTicker(6 * time.Second)
		for range t.C {
			m.num = new(big.Int).Add(m.num, big.NewInt(1))
		}
	}()
	return m
}

// Basic blockchain info
func (m *mem) ChainID(ctx context.Context) (*big.Int, error) { return m.chainID, nil }
func (m *mem) ClientVersion(ctx context.Context) (string, error) {
	return "memory-backend/0.1.0 (mock)", nil
}
func (m *mem) BlockNumber(ctx context.Context) (*big.Int, error) { return new(big.Int).Set(m.num), nil }

// Block operations
func (m *mem) BlockByNumber(ctx context.Context, num *big.Int, fullTx bool) (*Types.Block, error) {
	header := &Types.BlockHeader{
		Number:     m.num.Uint64(),
		Hash:       []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		ParentHash: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   30000000,
		GasUsed:    21000,
	}
	return &Types.Block{Header: header, Transactions: []*Types.Transaction{}}, nil
}
func (m *mem) BlockByHash(ctx context.Context, hash []byte, fullTx bool) (*Types.Block, error) {
	header := &Types.BlockHeader{
		Number:     m.num.Uint64(),
		Hash:       hash,
		ParentHash: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   30000000,
		GasUsed:    21000,
	}
	return &Types.Block{Header: header, Transactions: []*Types.Transaction{}}, nil
}
func (m *mem) BlockTransactionCountByNumber(ctx context.Context, blockNum *big.Int) (uint64, error) {
	return 0, nil // Mock: no transactions in blocks
}
func (m *mem) BlockTransactionCountByHash(ctx context.Context, blockHash []byte) (uint64, error) {
	return 0, nil // Mock: no transactions in blocks
}

// Account operations
func (m *mem) Balance(ctx context.Context, addr []byte, block *big.Int) (*big.Int, error) {
	// Mock balances for demonstration purposes
	// In a real implementation, you would query your blockchain node
	addrStr := hex.EncodeToString(addr)
	switch addrStr {
	case "31fcb3c05f73242aedd88b024e33d25a81fe67db":
		// 100 tokens = 100 * 10^18 wei
		return new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), nil
	case "a2902c128d42a64f371457b82bb6abb05b9b8bf1":
		// 150 tokens = 150 * 10^18 wei
		return new(big.Int).Mul(big.NewInt(150), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), nil
	default:
		return big.NewInt(0), nil
	}
}
func (m *mem) GetCode(ctx context.Context, addr []byte, block *big.Int) ([]byte, error) {
	// Mock: return empty code for all addresses
	return []byte{}, nil
}
func (m *mem) GetStorageAt(ctx context.Context, addr []byte, key []byte, block *big.Int) ([]byte, error) {
	// Mock: return zero storage for all addresses
	return make([]byte, 32), nil // 32 bytes of zeros
}
func (m *mem) GetTransactionCount(ctx context.Context, addr []byte, block *big.Int) (uint64, error) {
	// Mock: return 0 nonce for all addresses
	return 0, nil
}

// Transaction operations
func (m *mem) Call(ctx context.Context, msg Types.CallMsg, block *big.Int) ([]byte, error) {
	return []byte{}, nil
}
func (m *mem) EstimateGas(ctx context.Context, msg Types.CallMsg) (uint64, error) {
	return 21000, nil
}
func (m *mem) GasPrice(ctx context.Context) (*big.Int, error) { return big.NewInt(1_000_000_000), nil }
func (m *mem) SendRawTx(ctx context.Context, rawHex string) ([]byte, error) {
	return []byte{0xde, 0xad, 0xbe, 0xef}, nil
}
func (m *mem) TxByHash(ctx context.Context, hash []byte) (*Types.Transaction, error) {
	return &Types.Transaction{Hash: []byte{0xde, 0xad, 0xbe, 0xef}}, nil
}
func (m *mem) TxByBlockNumberAndIndex(ctx context.Context, blockNum *big.Int, index uint64) (*Types.Transaction, error) {
	return &Types.Transaction{Hash: []byte{0xde, 0xad, 0xbe, 0xef}}, nil
}
func (m *mem) TxByBlockHashAndIndex(ctx context.Context, blockHash []byte, index uint64) (*Types.Transaction, error) {
	return &Types.Transaction{Hash: []byte{0xde, 0xad, 0xbe, 0xef}}, nil
}
func (m *mem) ReceiptByHash(ctx context.Context, hash []byte) (*Types.Receipt, error) {
	return &Types.Receipt{
		TxHash:      []byte{0xde, 0xad, 0xbe, 0xef},
		Status:      1,
		GasUsed:     21000,
		BlockNumber: 1,
		Logs:        []*Types.Log{},
	}, nil
}

// Log operations
func (m *mem) GetLogs(ctx context.Context, q Types.FilterQuery) ([]*Types.Log, error) {
	return nil, nil
}

// Network operations
func (m *mem) PeerCount(ctx context.Context) (uint64, error) {
	return 0, nil // Mock: no peers
}
func (m *mem) Listening(ctx context.Context) (bool, error) {
	return true, nil // Mock: always listening
}
func (m *mem) Syncing(ctx context.Context) (map[string]any, error) {
	return map[string]any{"syncing": false}, nil // Mock: not syncing
}

// Mining operations (for PoW chains)
func (m *mem) Mining(ctx context.Context) (bool, error) {
	return false, nil // Mock: not mining
}
func (m *mem) Hashrate(ctx context.Context) (uint64, error) {
	return 0, nil // Mock: no hashrate
}

// Uncle operations (for PoW chains)
func (m *mem) UncleCountByBlockNumber(ctx context.Context, blockNum *big.Int) (uint64, error) {
	return 0, nil // Mock: no uncles
}
func (m *mem) UncleCountByBlockHash(ctx context.Context, blockHash []byte) (uint64, error) {
	return 0, nil // Mock: no uncles
}
func (m *mem) UncleByBlockNumberAndIndex(ctx context.Context, blockNum *big.Int, index uint64) (*Types.Block, error) {
	return nil, errors.New("no uncles available") // Mock: no uncles
}
func (m *mem) UncleByBlockHashAndIndex(ctx context.Context, blockHash []byte, index uint64) (*Types.Block, error) {
	return nil, errors.New("no uncles available") // Mock: no uncles
}

// Mock subscription methods for demo purposes
func (m *mem) SubscribeNewHeads(ctx context.Context) (<-chan *Types.Block, func(), error) {
	out := make(chan *Types.Block, 1)
	stop := make(chan struct{})
	go func() {
		t := time.NewTicker(6 * time.Second)
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case <-t.C:
				header := &Types.BlockHeader{
					Number:     m.num.Uint64(),
					Hash:       []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					ParentHash: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					Timestamp:  uint64(time.Now().Unix()),
					GasLimit:   30000000,
					GasUsed:    21000,
				}
				out <- &Types.Block{Header: header, Transactions: []*Types.Transaction{}}
			}
		}
	}()
	return out, func() { close(stop) }, nil
}
func (m *mem) SubscribeLogs(ctx context.Context, q *Types.FilterQuery) (<-chan *Types.Log, func(), error) {
	ch := make(chan *Types.Log)
	return ch, func() { close(ch) }, nil
}
func (m *mem) SubscribePendingTxs(ctx context.Context) (<-chan []byte, func(), error) {
	ch := make(chan []byte)
	return ch, func() { close(ch) }, nil
}
