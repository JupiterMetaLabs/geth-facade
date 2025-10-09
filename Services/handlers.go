package Services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/jupitermetalabs/geth-facade/Types"
)

// Handlers manages JSON-RPC request handling
// //debugging: Includes request/response logging for debugging
// //future: May add rate limiting and caching
type Handlers struct{ be Types.Backend }

func NewHandlers(be Types.Backend) *Handlers { return &Handlers{be: be} }

func (h *Handlers) Handle(ctx context.Context, req Types.Request) (Types.Response, error) {
	// //debugging: Log incoming request for debugging
	reqJSON, _ := json.Marshal(req)
	log.Printf("📥 RPC Request: %s", string(reqJSON))

	switch req.Method {
	case "web3_clientVersion":
		v, err := h.be.ClientVersion(ctx)
		resp, _ := finish(req, v, err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err
	case "net_version":
		id, err := h.be.ChainID(ctx)
		resp, _ := finish(req, id.String(), err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err
	case "eth_chainId":
		id, err := h.be.ChainID(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+id.Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_blockNumber":
		n, err := h.be.BlockNumber(ctx)
		resp, _ := finish(req, "0x"+n.Text(16), err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err

	case "eth_getBlockByNumber":
		// params: [blockTag, fullTx(bool)]
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		tag, _ := req.Params[0].(string)
		full := false
		if len(req.Params) > 1 {
			if b, ok := req.Params[1].(bool); ok {
				full = b
			}
		}
		num, err := parseBlockTag(ctx, h.be, tag)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		b, err := h.be.BlockByNumber(ctx, num, full)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalBlock(b, full), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getBalance":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need address and block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		addrStr, _ := req.Params[0].(string)
		addr, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		bal, err := h.be.Balance(ctx, addr, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+bal.Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_call":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing call object")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		msg, err := toCallMsg(req.Params[0])
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		var num *big.Int
		if len(req.Params) > 1 {
			num, err = parseBlockTag(ctx, h.be, mustString(req.Params[1]))
			if err != nil {
				resp, _ := finish(req, nil, err)
				log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
				return resp, err
			}
		}
		out, err := h.be.Call(ctx, msg, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+hex.EncodeToString(out), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_estimateGas":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing tx object")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		msg, err := toCallMsg(req.Params[0])
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		g, err := h.be.EstimateGas(ctx, msg)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+big.NewInt(int64(g)).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_gasPrice":
		p, err := h.be.GasPrice(ctx)
		resp, _ := finish(req, "0x"+p.Text(16), err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err

	case "eth_sendRawTransaction":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing raw tx")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		raw, _ := req.Params[0].(string)
		txh, err := h.be.SendRawTx(ctx, raw)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+hex.EncodeToString(txh), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getTransactionByHash":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing tx hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		tx, err := h.be.TxByHash(ctx, hash)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalTx(tx), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getTransactionReceipt":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing tx hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		rcpt, err := h.be.ReceiptByHash(ctx, hash)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalReceipt(rcpt), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getLogs":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing filter")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		q, err := toFilterQuery(req.Params[0])
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		logs, err := h.be.GetLogs(ctx, *q)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalLogs(logs), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Block operations by hash
	case "eth_getBlockByHash":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		full := false
		if len(req.Params) > 1 {
			if b, ok := req.Params[1].(bool); ok {
				full = b
			}
		}
		b, err := h.be.BlockByHash(ctx, hash, full)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalBlock(b, full), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Block transaction count operations
	case "eth_getBlockTransactionCountByNumber":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[0]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		count, err := h.be.BlockTransactionCountByNumber(ctx, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getBlockTransactionCountByHash":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		count, err := h.be.BlockTransactionCountByHash(ctx, hash)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Account operations
	case "eth_getCode":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need address and block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		addrStr, _ := req.Params[0].(string)
		addr, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		code, err := h.be.GetCode(ctx, addr, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+hex.EncodeToString(code), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getStorageAt":
		if len(req.Params) < 3 {
			resp, _ := invalidParams(req, "need address, key, and block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		addrStr, _ := req.Params[0].(string)
		addr, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		keyStr, _ := req.Params[1].(string)
		key, err := hex.DecodeString(strings.TrimPrefix(keyStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[2]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		storage, err := h.be.GetStorageAt(ctx, addr, key, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+hex.EncodeToString(storage), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getTransactionCount":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need address and block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		addrStr, _ := req.Params[0].(string)
		addr, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		count, err := h.be.GetTransactionCount(ctx, addr, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Transaction operations by block and index
	case "eth_getTransactionByBlockNumberAndIndex":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need block tag and index")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[0]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		index, err := parseHexUint64(mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		tx, err := h.be.TxByBlockNumberAndIndex(ctx, num, index)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalTx(tx), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getTransactionByBlockHashAndIndex":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need block hash and index")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		index, err := parseHexUint64(mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		tx, err := h.be.TxByBlockHashAndIndex(ctx, hash, index)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalTx(tx), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Network operations
	case "net_peerCount":
		count, err := h.be.PeerCount(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "net_listening":
		listening, err := h.be.Listening(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, listening, nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Sync operations
	case "eth_syncing":
		sync, err := h.be.Syncing(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, sync, nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Mining operations (for PoW chains)
	case "eth_mining":
		mining, err := h.be.Mining(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, mining, nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_hashrate":
		hashrate, err := h.be.Hashrate(ctx)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(hashrate).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	// Uncle operations (for PoW chains)
	case "eth_getUncleCountByBlockNumber":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block tag")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[0]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		count, err := h.be.UncleCountByBlockNumber(ctx, num)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getUncleCountByBlockHash":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing block hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		count, err := h.be.UncleCountByBlockHash(ctx, hash)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, "0x"+new(big.Int).SetUint64(count).Text(16), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getUncleByBlockNumberAndIndex":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need block tag and index")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[0]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		index, err := parseHexUint64(mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		uncle, err := h.be.UncleByBlockNumberAndIndex(ctx, num, index)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalBlock(uncle, false), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	case "eth_getUncleByBlockHashAndIndex":
		if len(req.Params) < 2 {
			resp, _ := invalidParams(req, "need block hash and index")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		hashStr, _ := req.Params[0].(string)
		hash, err := hex.DecodeString(strings.TrimPrefix(hashStr, "0x"))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		index, err := parseHexUint64(mustString(req.Params[1]))
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		uncle, err := h.be.UncleByBlockHashAndIndex(ctx, hash, index)
		if err != nil {
			resp, _ := finish(req, nil, err)
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, err
		}
		resp, _ := finish(req, marshalBlock(uncle, false), nil)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil

	default:
		resp := Types.RespErr(req.ID, -32601, "Method not found")
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil
	}
}

func parseBlockTag(ctx context.Context, be Types.Backend, tag string) (*big.Int, error) {
	switch strings.ToLower(tag) {
	case "latest", "":
		return be.BlockNumber(ctx)
	case "pending":
		// map to latest for now; refine if you track pending state
		return be.BlockNumber(ctx)
	default:
		if strings.HasPrefix(tag, "0x") {
			n := new(big.Int)
			n.SetString(tag[2:], 16)
			return n, nil
		}
		return nil, errors.New("unsupported block tag")
	}
}

func finish(req Types.Request, v any, err error) (Types.Response, error) {
	if err != nil {
		return Types.RespErr(req.ID, -32000, err.Error()), nil
	}
	return Types.RespOK(req.ID, v), nil
}

func invalidParams(req Types.Request, msg string) (Types.Response, error) {
	return Types.RespErr(req.ID, -32602, msg), nil
}

func mustString(v any) string {
	s, _ := v.(string)
	return s
}

func parseHexUint64(hexStr string) (uint64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	bigInt := new(big.Int)
	_, ok := bigInt.SetString(hexStr, 16)
	if !ok {
		return 0, errors.New("invalid hex string")
	}
	return bigInt.Uint64(), nil
}

func toCallMsg(p any) (Types.CallMsg, error) {
	// Parse call object from JSON-RPC params
	if callObj, ok := p.(map[string]any); ok {
		msg := Types.CallMsg{}

		if from, ok := callObj["from"].(string); ok {
			msg.From = from
		}
		if to, ok := callObj["to"].(string); ok {
			msg.To = to
		}
		if data, ok := callObj["data"].(string); ok {
			if strings.HasPrefix(data, "0x") {
				msg.Data, _ = hex.DecodeString(data[2:])
			} else {
				msg.Data, _ = hex.DecodeString(data)
			}
		}
		if value, ok := callObj["value"].(string); ok {
			if strings.HasPrefix(value, "0x") {
				bigVal := new(big.Int)
				bigVal.SetString(value[2:], 16)
				msg.Value = bigVal
			}
		}
		if gas, ok := callObj["gas"].(string); ok {
			if strings.HasPrefix(gas, "0x") {
				bigGas := new(big.Int)
				bigGas.SetString(gas[2:], 16)
				msg.Gas = bigGas
			}
		}
		if gasPrice, ok := callObj["gasPrice"].(string); ok {
			if strings.HasPrefix(gasPrice, "0x") {
				bigGasPrice := new(big.Int)
				bigGasPrice.SetString(gasPrice[2:], 16)
				msg.GasPrice = bigGasPrice
			}
		}

		return msg, nil
	}
	return Types.CallMsg{}, errors.New("invalid call object")
}

func toFilterQuery(p any) (*Types.FilterQuery, error) {
	// Parse filter object from JSON-RPC params
	if filterObj, ok := p.(map[string]any); ok {
		query := &Types.FilterQuery{}

		if fromBlock, ok := filterObj["fromBlock"].(string); ok {
			if strings.HasPrefix(fromBlock, "0x") {
				bigFromBlock := new(big.Int)
				bigFromBlock.SetString(fromBlock[2:], 16)
				query.FromBlock = bigFromBlock
			}
		}
		if toBlock, ok := filterObj["toBlock"].(string); ok {
			if strings.HasPrefix(toBlock, "0x") {
				bigToBlock := new(big.Int)
				bigToBlock.SetString(toBlock[2:], 16)
				query.ToBlock = bigToBlock
			}
		}
		if addresses, ok := filterObj["address"].([]any); ok {
			query.Addresses = make([][]byte, len(addresses))
			for i, addr := range addresses {
				if addrStr, ok := addr.(string); ok {
					addrBytes, err := hex.DecodeString(strings.TrimPrefix(addrStr, "0x"))
					if err == nil {
						query.Addresses[i] = addrBytes
					}
				}
			}
		}
		if topics, ok := filterObj["topics"].([]any); ok {
			query.Topics = make([][]byte, len(topics))
			for i, topic := range topics {
				if topicArr, ok := topic.([]any); ok {
					// For now, just take the first topic if it's an array
					if len(topicArr) > 0 {
						if topicStr, ok := topicArr[0].(string); ok {
							topicBytes, err := hex.DecodeString(strings.TrimPrefix(topicStr, "0x"))
							if err == nil {
								query.Topics[i] = topicBytes
							}
						}
					}
				} else if topicStr, ok := topic.(string); ok {
					topicBytes, err := hex.DecodeString(strings.TrimPrefix(topicStr, "0x"))
					if err == nil {
						query.Topics[i] = topicBytes
					}
				}
			}
		}

		return query, nil
	}
	return &Types.FilterQuery{}, errors.New("invalid filter object")
}

// marshalBlock converts a Block to JSON-RPC format
// //conversions: Converts []byte fields to hex strings for JSON-RPC
// //debugging: Used for block data serialization
func marshalBlock(b *Types.Block, full bool) map[string]any {
	result := map[string]any{
		"number":        "0x" + new(big.Int).SetUint64(b.Header.Number).Text(16),
		"hash":          "0x" + hex.EncodeToString(b.Header.Hash),
		"parentHash":    "0x" + hex.EncodeToString(b.Header.ParentHash),
		"stateRoot":     "0x" + hex.EncodeToString(b.Header.StateRoot),
		"receiptsRoot":  "0x" + hex.EncodeToString(b.Header.ReceiptsRoot),
		"logsBloom":     "0x" + hex.EncodeToString(b.Header.LogsBloom),
		"miner":         "0x" + hex.EncodeToString(b.Header.Miner),
		"gasLimit":      "0x" + new(big.Int).SetUint64(b.Header.GasLimit).Text(16),
		"gasUsed":       "0x" + new(big.Int).SetUint64(b.Header.GasUsed).Text(16),
		"timestamp":     "0x" + new(big.Int).SetUint64(b.Header.Timestamp).Text(16),
		"mixHash":       "0x" + hex.EncodeToString(b.Header.MixHashOrPrevRandao),
		"baseFeePerGas": "0x" + hex.EncodeToString(b.Header.BaseFee),
		"extraData":     "0x" + hex.EncodeToString(b.Header.ExtraData),
		"transactions":  []any{},
		"uncles":        []any{},
		"withdrawals":   []any{},
	}

	// Add blob gas fields if present
	if len(b.BlobGasUsed) > 0 {
		result["blobGasUsed"] = "0x" + hex.EncodeToString(b.BlobGasUsed)
	}
	if len(b.ExcessBlobGas) > 0 {
		result["excessBlobGas"] = "0x" + hex.EncodeToString(b.ExcessBlobGas)
	}

	// Add withdrawals if present
	if len(b.Withdrawals) > 0 {
		withdrawals := make([]any, len(b.Withdrawals))
		for i, w := range b.Withdrawals {
			withdrawals[i] = map[string]any{
				"index":          "0x" + new(big.Int).SetUint64(w.Index).Text(16),
				"validatorIndex": "0x" + new(big.Int).SetUint64(w.ValidatorIndex).Text(16),
				"address":        "0x" + hex.EncodeToString(w.Address),
				"amount":         "0x" + new(big.Int).SetUint64(w.Amount).Text(16),
			}
		}
		result["withdrawals"] = withdrawals
	}

	// Add transactions
	if full && len(b.Transactions) > 0 {
		txs := make([]any, len(b.Transactions))
		for i, tx := range b.Transactions {
			txs[i] = marshalTx(tx)
		}
		result["transactions"] = txs
	} else if len(b.Transactions) > 0 {
		txHashes := make([]string, len(b.Transactions))
		for i, tx := range b.Transactions {
			txHashes[i] = "0x" + hex.EncodeToString(tx.Hash)
		}
		result["transactions"] = txHashes
	}

	// Add uncles
	if len(b.Ommers) > 0 {
		uncles := make([]string, len(b.Ommers))
		for i, uncle := range b.Ommers {
			uncles[i] = "0x" + hex.EncodeToString(uncle)
		}
		result["uncles"] = uncles
	}

	return result
}

func marshalTx(tx *Types.Transaction) map[string]any {
	result := map[string]any{
		"hash":     "0x" + hex.EncodeToString(tx.Hash),
		"from":     "0x" + hex.EncodeToString(tx.From),
		"to":       "0x" + hex.EncodeToString(tx.To),
		"input":    "0x" + hex.EncodeToString(tx.Input),
		"value":    "0x" + hex.EncodeToString(tx.Value),
		"nonce":    "0x" + new(big.Int).SetUint64(tx.Nonce).Text(16),
		"gas":      "0x" + new(big.Int).SetUint64(tx.Gas).Text(16),
		"gasPrice": "0x" + hex.EncodeToString(tx.GasPrice),
		"type":     "0x" + new(big.Int).SetUint64(uint64(tx.Type)).Text(16),
		"r":        "0x" + hex.EncodeToString(tx.R),
		"s":        "0x" + hex.EncodeToString(tx.S),
		"v":        "0x" + new(big.Int).SetUint64(uint64(tx.V)).Text(16),
	}

	// Add EIP-1559 fields if present
	if len(tx.MaxFeePerGas) > 0 {
		result["maxFeePerGas"] = "0x" + hex.EncodeToString(tx.MaxFeePerGas)
	}
	if len(tx.MaxPriorityFeePerGas) > 0 {
		result["maxPriorityFeePerGas"] = "0x" + hex.EncodeToString(tx.MaxPriorityFeePerGas)
	}

	// Add EIP-4844 blob fields if present
	if len(tx.MaxFeePerBlobGas) > 0 {
		result["maxFeePerBlobGas"] = "0x" + hex.EncodeToString(tx.MaxFeePerBlobGas)
	}
	if len(tx.BlobVersionedHashes) > 0 {
		hashes := make([]string, len(tx.BlobVersionedHashes))
		for i, hash := range tx.BlobVersionedHashes {
			hashes[i] = "0x" + hex.EncodeToString(hash)
		}
		result["blobVersionedHashes"] = hashes
	}

	// Add access list if present
	if tx.AccessList != nil && len(tx.AccessList.AccessTuples) > 0 {
		accessList := make([]any, len(tx.AccessList.AccessTuples))
		for i, tuple := range tx.AccessList.AccessTuples {
			storageKeys := make([]string, len(tuple.StorageKeys))
			for j, key := range tuple.StorageKeys {
				storageKeys[j] = "0x" + hex.EncodeToString(key)
			}
			accessList[i] = map[string]any{
				"address":     "0x" + hex.EncodeToString(tuple.Address),
				"storageKeys": storageKeys,
			}
		}
		result["accessList"] = accessList
	}

	return result
}

func marshalLogs(logs []*Types.Log) []map[string]any {
	result := make([]map[string]any, len(logs))
	for i, log := range logs {
		topics := make([]string, len(log.Topics))
		for j, topic := range log.Topics {
			topics[j] = "0x" + hex.EncodeToString(topic)
		}

		result[i] = map[string]any{
			"address":          "0x" + hex.EncodeToString(log.Address),
			"topics":           topics,
			"data":             "0x" + hex.EncodeToString(log.Data),
			"blockNumber":      "0x" + new(big.Int).SetUint64(log.BlockNumber).Text(16),
			"blockHash":        "0x" + hex.EncodeToString(log.BlockHash),
			"transactionHash":  "0x" + hex.EncodeToString(log.TxHash),
			"transactionIndex": "0x" + new(big.Int).SetUint64(log.TxIndex).Text(16),
			"logIndex":         "0x" + new(big.Int).SetUint64(log.LogIndex).Text(16),
			"removed":          log.Removed,
		}
	}
	return result
}

func marshalReceipt(receipt *Types.Receipt) map[string]any {
	result := map[string]any{
		"transactionHash":   "0x" + hex.EncodeToString(receipt.TxHash),
		"status":            "0x" + new(big.Int).SetUint64(receipt.Status).Text(16),
		"cumulativeGasUsed": "0x" + new(big.Int).SetUint64(receipt.CumulativeGasUsed).Text(16),
		"gasUsed":           "0x" + new(big.Int).SetUint64(receipt.GasUsed).Text(16),
		"blockNumber":       "0x" + new(big.Int).SetUint64(receipt.BlockNumber).Text(16),
		"blockHash":         "0x" + hex.EncodeToString(receipt.BlockHash),
		"transactionIndex":  "0x" + new(big.Int).SetUint64(receipt.TransactionIndex).Text(16),
		"type":              "0x" + new(big.Int).SetUint64(uint64(receipt.Type)).Text(16),
		"logs":              marshalLogs(receipt.Logs),
	}

	// Add contract address if present
	if len(receipt.ContractAddress) > 0 {
		result["contractAddress"] = "0x" + hex.EncodeToString(receipt.ContractAddress)
	}

	return result
}
