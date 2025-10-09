package Types

import (
	"context"
	"math/big"
)

// Block represents an Ethereum block with proper Geth-compatible structure
// //conversions: Uses []byte for hashes to match Geth's internal representation
// //future: May add additional fields for future EIPs
type Block struct {
	Header          *BlockHeader   `json:"header"`
	Transactions    []*Transaction `json:"transactions"`
	Ommers          [][]byte       `json:"ommers"`
	WithdrawalsRoot []byte         `json:"withdrawalsroot"`
	Withdrawals     []*Withdrawal  `json:"withdrawals"`
	BlobGasUsed     []byte         `json:"blobgasused"`
	ExcessBlobGas   []byte         `json:"excessblobgas"`
}

// BlockHeader represents the header of an Ethereum block
// //conversions: All hash fields use []byte for consistency with Geth
// //future: Additional fields may be added for new EIPs
type BlockHeader struct {
	ParentHash          []byte `json:"parenthash"`
	StateRoot           []byte `json:"stateroot"`
	ReceiptsRoot        []byte `json:"receiptsroot"`
	LogsBloom           []byte `json:"logsbloom"`
	Miner               []byte `json:"miner"`
	Number              uint64 `json:"number"`
	GasLimit            uint64 `json:"gaslimit"`
	GasUsed             uint64 `json:"gasused"`
	Timestamp           uint64 `json:"timestamp"`
	MixHashOrPrevRandao []byte `json:"mixhashorprevrandao"`
	BaseFee             []byte `json:"basefee"`
	BlobGasUsedField    uint64 `json:"blobgasusedfield"`
	ExcessBlobGasField  uint64 `json:"excessblobgasfield"`
	ExtraData           []byte `json:"extradata"`
	Hash                []byte `json:"hash"`
}

// Transaction represents an Ethereum transaction with modern features
// //conversions: Addresses and hashes use []byte for Geth compatibility
// //future: Additional transaction types may be added
type Transaction struct {
	Hash                 []byte      `json:"hash"`
	From                 []byte      `json:"from"`
	To                   []byte      `json:"to"`
	Input                []byte      `json:"input"`
	Nonce                uint64      `json:"nonce"`
	Value                []byte      `json:"value"`
	Gas                  uint64      `json:"gas"`
	GasPrice             []byte      `json:"gasprice"`
	Type                 uint32      `json:"type"`
	R                    []byte      `json:"r"`
	S                    []byte      `json:"s"`
	V                    uint32      `json:"v"`
	AccessList           *AccessList `json:"accesslist"`
	MaxFeePerGas         []byte      `json:"maxfeepergas"`
	MaxPriorityFeePerGas []byte      `json:"maxpriorityfeepergas"`
	MaxFeePerBlobGas     []byte      `json:"maxfeeperblobgas"`
	BlobVersionedHashes  [][]byte    `json:"blobversionedhashes"`
}

// AccessList represents EIP-2930 access list
type AccessList struct {
	AccessTuples []*AccessTuple `json:"accesstuples"`
}

// AccessTuple represents an access list entry
type AccessTuple struct {
	Address     []byte   `json:"address"`
	StorageKeys [][]byte `json:"storagekeys"`
}

// Withdrawal represents EIP-4895 withdrawal
type Withdrawal struct {
	Index          uint64 `json:"index"`
	ValidatorIndex uint64 `json:"validatorindex"`
	Address        []byte `json:"address"`
	Amount         uint64 `json:"amount"`
}

// Receipt represents a transaction receipt
type Receipt struct {
	TxHash            []byte `json:"txhash"`
	Status            uint64 `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativegasused"`
	GasUsed           uint64 `json:"gasused"`
	Logs              []*Log `json:"logs"`
	ContractAddress   []byte `json:"contractaddress"`
	Type              uint32 `json:"type"`
	BlockHash         []byte `json:"blockhash"`
	BlockNumber       uint64 `json:"blocknumber"`
	TransactionIndex  uint64 `json:"transactionindex"`
}

// Log represents an Ethereum log with complete fields
type Log struct {
	Address     []byte   `json:"address"`
	Topics      [][]byte `json:"topics"`
	Data        []byte   `json:"data"`
	BlockNumber uint64   `json:"blocknumber"`
	BlockHash   []byte   `json:"blockhash"`
	TxIndex     uint64   `json:"txindex"`
	TxHash      []byte   `json:"txhash"`
	LogIndex    uint64   `json:"logindex"`
	Removed     bool     `json:"removed"`
}

type CallMsg struct {
	From, To      string
	Data          []byte
	Value         *big.Int
	Gas, GasPrice *big.Int
}

type FilterQuery struct {
	FromBlock, ToBlock *big.Int
	Addresses          [][]byte
	Topics             [][]byte
	BlockHash          []byte
}

type Backend interface {
	// Basic blockchain info
	ChainID(ctx context.Context) (*big.Int, error)
	ClientVersion(ctx context.Context) (string, error)
	BlockNumber(ctx context.Context) (*big.Int, error)

	// Block operations
	BlockByNumber(ctx context.Context, num *big.Int, fullTx bool) (*Block, error)
	BlockByHash(ctx context.Context, hash []byte, fullTx bool) (*Block, error)
	BlockTransactionCountByNumber(ctx context.Context, blockNum *big.Int) (uint64, error)
	BlockTransactionCountByHash(ctx context.Context, blockHash []byte) (uint64, error)

	// Account operations
	Balance(ctx context.Context, addr []byte, block *big.Int) (*big.Int, error)
	GetCode(ctx context.Context, addr []byte, block *big.Int) ([]byte, error)
	GetStorageAt(ctx context.Context, addr []byte, key []byte, block *big.Int) ([]byte, error)
	GetTransactionCount(ctx context.Context, addr []byte, block *big.Int) (uint64, error)

	// Transaction operations
	Call(ctx context.Context, msg CallMsg, block *big.Int) ([]byte, error)
	EstimateGas(ctx context.Context, msg CallMsg) (uint64, error)
	GasPrice(ctx context.Context) (*big.Int, error)
	SendRawTx(ctx context.Context, rawHex string) ([]byte, error)
	TxByHash(ctx context.Context, hash []byte) (*Transaction, error)
	TxByBlockNumberAndIndex(ctx context.Context, blockNum *big.Int, index uint64) (*Transaction, error)
	TxByBlockHashAndIndex(ctx context.Context, blockHash []byte, index uint64) (*Transaction, error)
	ReceiptByHash(ctx context.Context, hash []byte) (*Receipt, error)

	// Log operations
	GetLogs(ctx context.Context, q FilterQuery) ([]*Log, error)

	// Network operations
	PeerCount(ctx context.Context) (uint64, error)
	Listening(ctx context.Context) (bool, error)
	Syncing(ctx context.Context) (map[string]any, error)

	// Mining operations (for PoW chains)
	Mining(ctx context.Context) (bool, error)
	Hashrate(ctx context.Context) (uint64, error)

	// Uncle operations (for PoW chains)
	UncleCountByBlockNumber(ctx context.Context, blockNum *big.Int) (uint64, error)
	UncleCountByBlockHash(ctx context.Context, blockHash []byte) (uint64, error)
	UncleByBlockNumberAndIndex(ctx context.Context, blockNum *big.Int, index uint64) (*Block, error)
	UncleByBlockHashAndIndex(ctx context.Context, blockHash []byte, index uint64) (*Block, error)

	// Streaming (for WS subscriptions)
	SubscribeNewHeads(ctx context.Context) (<-chan *Block, func(), error)
	SubscribeLogs(ctx context.Context, q *FilterQuery) (<-chan *Log, func(), error)
	SubscribePendingTxs(ctx context.Context) (<-chan []byte, func(), error)
}
