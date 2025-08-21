package rpc

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"jmdt-geth-facade/backend"
)

type Handlers struct{ be backend.Backend }

func NewHandlers(be backend.Backend) *Handlers { return &Handlers{be: be} }

func (h *Handlers) Handle(ctx context.Context, req Request) (Response, error) {
	switch req.Method {
	case "web3_clientVersion":
		v, err := h.be.ClientVersion(ctx); return finish(req, v, err)
	case "net_version":
		id, err := h.be.ChainID(ctx); return finish(req, id.String(), err)
	case "eth_chainId":
		id, err := h.be.ChainID(ctx); 
		if err != nil { return finish(req, nil, err) }
		return finish(req, "0x"+id.Text(16), nil)

	case "eth_blockNumber":
		n, err := h.be.BlockNumber(ctx); return finish(req, "0x"+n.Text(16), err)

	case "eth_getBlockByNumber":
		// params: [blockTag, fullTx(bool)]
		if len(req.Params) < 1 { return invalidParams(req, "missing block tag") }
		tag, _ := req.Params[0].(string)
		full := false
		if len(req.Params) > 1 {
			if b, ok := req.Params[1].(bool); ok { full = b }
		}
		num, err := parseBlockTag(ctx, h.be, tag)
		if err != nil { return finish(req, nil, err) }
		b, err := h.be.BlockByNumber(ctx, num, full)
		if err != nil { return finish(req, nil, err) }
		return finish(req, marshalBlock(b, full), nil)

	case "eth_getBalance":
		if len(req.Params) < 2 { return invalidParams(req, "need address and block tag") }
		addr, _ := req.Params[0].(string)
		num, err := parseBlockTag(ctx, h.be, mustString(req.Params[1]))
		if err != nil { return finish(req, nil, err) }
		bal, err := h.be.Balance(ctx, addr, num); 
		if err != nil { return finish(req, nil, err) }
		return finish(req, "0x"+bal.Text(16), nil)

	case "eth_call":
		if len(req.Params) < 1 { return invalidParams(req, "missing call object") }
		msg, err := toCallMsg(req.Params[0])
		if err != nil { return finish(req, nil, err) }
		var num *big.Int
		if len(req.Params) > 1 {
			num, err = parseBlockTag(ctx, h.be, mustString(req.Params[1]))
			if err != nil { return finish(req, nil, err) }
		}
		out, err := h.be.Call(ctx, msg, num)
		if err != nil { return finish(req, nil, err) }
		return finish(req, "0x"+hex.EncodeToString(out), nil)

	case "eth_estimateGas":
		if len(req.Params) < 1 { return invalidParams(req, "missing tx object") }
		msg, err := toCallMsg(req.Params[0])
		if err != nil { return finish(req, nil, err) }
		g, err := h.be.EstimateGas(ctx, msg)
		if err != nil { return finish(req, nil, err) }
		return finish(req, "0x"+big.NewInt(int64(g)).Text(16), nil)

	case "eth_gasPrice":
		p, err := h.be.GasPrice(ctx)
		return finish(req, "0x"+p.Text(16), err)

	case "eth_sendRawTransaction":
		if len(req.Params) < 1 { return invalidParams(req, "missing raw tx") }
		raw, _ := req.Params[0].(string)
		txh, err := h.be.SendRawTx(ctx, raw)
		return finish(req, txh, err)

	case "eth_getTransactionByHash":
		if len(req.Params) < 1 { return invalidParams(req, "missing tx hash") }
		tx, err := h.be.TxByHash(ctx, mustString(req.Params[0]))
		if err != nil { return finish(req, nil, err) }
		return finish(req, marshalTx(tx), nil)

	case "eth_getTransactionReceipt":
		if len(req.Params) < 1 { return invalidParams(req, "missing tx hash") }
		rcpt, err := h.be.ReceiptByHash(ctx, mustString(req.Params[0]))
		return finish(req, rcpt, err)

	case "eth_getLogs":
		if len(req.Params) < 1 { return invalidParams(req, "missing filter") }
		q, err := toFilterQuery(req.Params[0])
		if err != nil { return finish(req, nil, err) }
		logs, err := h.be.GetLogs(ctx, *q)
		if err != nil { return finish(req, nil, err) }
		return finish(req, marshalLogs(logs), nil)

	default:
		return RespErr(req.ID, -32601, "Method not found"), nil
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
	if err != nil { return RespErr(req.ID, -32000, err.Error()), nil }
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
	// TODO: implement this properly based on how call objects are structured
	return backend.CallMsg{}, nil
}

func toFilterQuery(p any) (*backend.FilterQuery, error) {
	// TODO: implement this properly
	return &backend.FilterQuery{}, nil
}

func marshalBlock(b *backend.Block, full bool) map[string]any {
	// TODO: implement this
	return map[string]any{}
}

func marshalTx(tx *backend.Tx) map[string]any {
	// TODO: implement this
	return map[string]any{}
}

func marshalLogs(logs []backend.Log) []map[string]any {
	// TODO: implement this
	return []map[string]any{}
}
