package rpc

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"encoding/json"
	"log"

	"github.com/saishibu/jmdt-geth-facade/backend"
)

type Handlers struct{ be backend.Backend }

func NewHandlers(be backend.Backend) *Handlers { return &Handlers{be: be} }

func (h *Handlers) Handle(ctx context.Context, req Request) (Response, error) {
	// Log incoming request
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
		addr, _ := req.Params[0].(string)
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
		resp, _ := finish(req, txh, err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err

	case "eth_getTransactionByHash":
		if len(req.Params) < 1 {
			resp, _ := invalidParams(req, "missing tx hash")
			log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
			return resp, nil
		}
		tx, err := h.be.TxByHash(ctx, mustString(req.Params[0]))
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
		rcpt, err := h.be.ReceiptByHash(ctx, mustString(req.Params[0]))
		resp, _ := finish(req, rcpt, err)
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, err

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

	default:
		resp := RespErr(req.ID, -32601, "Method not found")
		log.Printf("📤 RPC Response: %s -> %+v", req.Method, resp)
		return resp, nil
	}
}

func parseBlockTag(ctx context.Context, be backend.Backend, tag string) (*big.Int, error) {
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

func finish(req Request, v any, err error) (Response, error) {
	if err != nil {
		return RespErr(req.ID, -32000, err.Error()), nil
	}
	return RespOK(req.ID, v), nil
}

func invalidParams(req Request, msg string) (Response, error) {
	return RespErr(req.ID, -32602, msg), nil
}

func mustString(v any) string {
	s, _ := v.(string)
	return s
}

func toCallMsg(p any) (backend.CallMsg, error) {
	// Parse call object from JSON-RPC params
	if callObj, ok := p.(map[string]any); ok {
		msg := backend.CallMsg{}

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
	return backend.CallMsg{}, errors.New("invalid call object")
}

func toFilterQuery(p any) (*backend.FilterQuery, error) {
	// Parse filter object from JSON-RPC params
	if filterObj, ok := p.(map[string]any); ok {
		query := &backend.FilterQuery{}

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
			query.Addresses = make([]string, len(addresses))
			for i, addr := range addresses {
				if addrStr, ok := addr.(string); ok {
					query.Addresses[i] = addrStr
				}
			}
		}
		if topics, ok := filterObj["topics"].([]any); ok {
			query.Topics = make([][]string, len(topics))
			for i, topic := range topics {
				if topicArr, ok := topic.([]any); ok {
					query.Topics[i] = make([]string, len(topicArr))
					for j, t := range topicArr {
						if topicStr, ok := t.(string); ok {
							query.Topics[i][j] = topicStr
						}
					}
				} else if topicStr, ok := topic.(string); ok {
					query.Topics[i] = []string{topicStr}
				}
			}
		}

		return query, nil
	}
	return &backend.FilterQuery{}, errors.New("invalid filter object")
}

func marshalBlock(b *backend.Block, full bool) map[string]any {
	result := map[string]any{
		"number":       "0x" + b.Number.Text(16),
		"hash":         b.Hash,
		"parentHash":   b.ParentHash,
		"timestamp":    "0x" + new(big.Int).SetUint64(b.Timestamp).Text(16),
		"transactions": []any{},
	}

	if full && len(b.Transactions) > 0 {
		txs := make([]any, len(b.Transactions))
		for i, tx := range b.Transactions {
			txs[i] = marshalTx(&tx)
		}
		result["transactions"] = txs
	} else if len(b.Transactions) > 0 {
		txHashes := make([]string, len(b.Transactions))
		for i, tx := range b.Transactions {
			txHashes[i] = tx.Hash
		}
		result["transactions"] = txHashes
	}

	return result
}

func marshalTx(tx *backend.Tx) map[string]any {
	result := map[string]any{
		"hash":     tx.Hash,
		"from":     tx.From,
		"to":       tx.To,
		"input":    "0x" + hex.EncodeToString(tx.Input),
		"value":    "0x" + tx.Value.Text(16),
		"nonce":    "0x" + new(big.Int).SetUint64(tx.Nonce).Text(16),
		"gas":      "0x" + tx.Gas.Text(16),
		"gasPrice": "0x" + tx.GasPrice.Text(16),
	}
	return result
}

func marshalLogs(logs []backend.Log) []map[string]any {
	result := make([]map[string]any, len(logs))
	for i, log := range logs {
		result[i] = map[string]any{
			"address":         log.Address,
			"topics":          log.Topics,
			"data":            "0x" + hex.EncodeToString(log.Data),
			"blockNumber":     "0x" + log.BlockNumber.Text(16),
			"transactionHash": log.TxHash,
			"logIndex":        "0x" + new(big.Int).SetUint64(uint64(log.LogIndex)).Text(16),
		}
	}
	return result
}
