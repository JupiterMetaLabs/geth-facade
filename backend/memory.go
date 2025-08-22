package backend

import (
	"context"
	"math/big"
	"time"
)

type mem struct {
	chainID *big.Int
	num     *big.Int
}

func NewMemoryBackend(chainID *big.Int) Backend {
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
	return "geth-facade/0.1 (go)", nil
}
func (m *mem) BlockNumber(ctx context.Context) (*big.Int, error) { return new(big.Int).Set(m.num), nil }
func (m *mem) BlockByNumber(ctx context.Context, num *big.Int, fullTx bool) (*Block, error) {
	return &Block{Number: new(big.Int).Set(m.num), Hash: "0x0", ParentHash: "0x0", Timestamp: uint64(time.Now().Unix())}, nil
}
func (m *mem) Balance(ctx context.Context, addr string, block *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (m *mem) Call(ctx context.Context, msg CallMsg, block *big.Int) ([]byte, error) {
	return []byte{}, nil
}
func (m *mem) EstimateGas(ctx context.Context, msg CallMsg) (uint64, error) { return 21000, nil }
func (m *mem) GasPrice(ctx context.Context) (*big.Int, error)               { return big.NewInt(1_000_000_000), nil }
func (m *mem) SendRawTx(ctx context.Context, rawHex string) (string, error) { return "0xdeadbeef", nil }
func (m *mem) TxByHash(ctx context.Context, hash string) (*Tx, error) {
	return &Tx{Hash: "0xdeadbeef"}, nil
}
func (m *mem) ReceiptByHash(ctx context.Context, hash string) (map[string]any, error) {
	return map[string]any{"transactionHash": "0xdeadbeef", "status": "0x1", "blockNumber": "0x1"}, nil
}
func (m *mem) GetLogs(ctx context.Context, q FilterQuery) ([]Log, error) { return nil, nil }

// naive tick-based channels for demo purposes
func (m *mem) SubscribeNewHeads(ctx context.Context) (<-chan *Block, func(), error) {
	out := make(chan *Block, 1)
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
				out <- &Block{Number: new(big.Int).Set(m.num), Hash: "0x0", ParentHash: "0x0", Timestamp: uint64(time.Now().Unix())}
			}
		}
	}()
	return out, func() { close(stop) }, nil
}
func (m *mem) SubscribeLogs(ctx context.Context, q *FilterQuery) (<-chan Log, func(), error) {
	ch := make(chan Log)
	return ch, func() { close(ch) }, nil
}
func (m *mem) SubscribePendingTxs(ctx context.Context) (<-chan string, func(), error) {
	ch := make(chan string)
	return ch, func() { close(ch) }, nil
}
