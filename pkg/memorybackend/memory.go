// Package memorybackend provides a mock implementation of the backend.Backend interface
// for testing and development purposes. This is NOT suitable for production use.
//
// This example shows how to implement the backend.Backend interface and can be used
// as a reference for creating your own backend implementations.
package memorybackend

import (
	"context"
	"math/big"
	"time"

	"github.com/jupitermetalabs/geth-facade/backend"
)

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
func NewMemoryBackend(chainID *big.Int) backend.Backend {
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

func (m *mem) ChainID(ctx context.Context) (*big.Int, error) { return m.chainID, nil }
func (m *mem) ClientVersion(ctx context.Context) (string, error) {
	return "memory-backend/0.1.0 (mock)", nil
}
func (m *mem) BlockNumber(ctx context.Context) (*big.Int, error) { return new(big.Int).Set(m.num), nil }
func (m *mem) BlockByNumber(ctx context.Context, num *big.Int, fullTx bool) (*backend.Block, error) {
	return &backend.Block{Number: new(big.Int).Set(m.num), Hash: "0x0", ParentHash: "0x0", Timestamp: uint64(time.Now().Unix())}, nil
}
func (m *mem) Balance(ctx context.Context, addr string, block *big.Int) (*big.Int, error) {
	// Mock balances for demonstration purposes
	// In a real implementation, you would query your blockchain node
	switch addr {
	case "0x31fcB3C05F73242AeDd88b024E33d25a81Fe67DB":
		// 100 tokens = 100 * 10^18 wei
		return new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), nil
	case "0xA2902C128D42A64F371457b82BB6aBb05B9b8bf1":
		// 150 tokens = 150 * 10^18 wei
		return new(big.Int).Mul(big.NewInt(150), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), nil
	default:
		return big.NewInt(0), nil
	}
}
func (m *mem) Call(ctx context.Context, msg backend.CallMsg, block *big.Int) ([]byte, error) {
	return []byte{}, nil
}
func (m *mem) EstimateGas(ctx context.Context, msg backend.CallMsg) (uint64, error) {
	return 21000, nil
}
func (m *mem) GasPrice(ctx context.Context) (*big.Int, error)               { return big.NewInt(1_000_000_000), nil }
func (m *mem) SendRawTx(ctx context.Context, rawHex string) (string, error) { return "0xdeadbeef", nil }
func (m *mem) TxByHash(ctx context.Context, hash string) (*backend.Tx, error) {
	return &backend.Tx{Hash: "0xdeadbeef"}, nil
}
func (m *mem) ReceiptByHash(ctx context.Context, hash string) (map[string]any, error) {
	return map[string]any{"transactionHash": "0xdeadbeef", "status": "0x1", "blockNumber": "0x1"}, nil
}
func (m *mem) GetLogs(ctx context.Context, q backend.FilterQuery) ([]backend.Log, error) {
	return nil, nil
}

// Mock subscription methods for demo purposes
func (m *mem) SubscribeNewHeads(ctx context.Context) (<-chan *backend.Block, func(), error) {
	out := make(chan *backend.Block, 1)
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
				out <- &backend.Block{Number: new(big.Int).Set(m.num), Hash: "0x0", ParentHash: "0x0", Timestamp: uint64(time.Now().Unix())}
			}
		}
	}()
	return out, func() { close(stop) }, nil
}
func (m *mem) SubscribeLogs(ctx context.Context, q *backend.FilterQuery) (<-chan backend.Log, func(), error) {
	ch := make(chan backend.Log)
	return ch, func() { close(ch) }, nil
}
func (m *mem) SubscribePendingTxs(ctx context.Context) (<-chan string, func(), error) {
	ch := make(chan string)
	return ch, func() { close(ch) }, nil
}
