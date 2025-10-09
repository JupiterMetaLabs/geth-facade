package Services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jupitermetalabs/geth-facade/Types"

	"github.com/gorilla/websocket"
)

// WSServer provides WebSocket JSON-RPC server for real-time subscriptions
// //future: May add connection pooling and rate limiting
// //debugging: Includes connection and subscription logging
type WSServer struct {
	h   *Handlers
	be  Types.Backend
	upg websocket.Upgrader
}

func NewWSServer(h *Handlers, be Types.Backend) *WSServer {
	return &WSServer{
		h: h, be: be,
		upg: websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }},
	}
}

func (s *WSServer) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleWS)
	srv := &http.Server{Addr: addr, Handler: mux}
	return srv.ListenAndServe()
}

type sub struct {
	id   string
	stop func()
}

func (s *WSServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upg.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subs := map[string]*sub{}
	var mu sync.Mutex

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Log incoming WebSocket message
		log.Printf("🔌 WS Message: %s", string(data))

		var req Types.Request
		if err := json.Unmarshal(data, &req); err != nil {
			log.Printf("❌ WS Parse Error: %v", err)
			_ = conn.WriteJSON(Types.RespErr(nil, -32700, "Parse error"))
			continue
		}

		if req.Method == "eth_subscribe" {
			// params: [subscriptionType, (optional) filter]
			if len(req.Params) < 1 {
				_ = conn.WriteJSON(Types.RespErr(req.ID, -32602, "missing subscription type"))
				continue
			}
			typ, _ := req.Params[0].(string)
			sid := "0x" + time.Now().Format("150405.000000") // simple unique id; replace in prod

			switch typ {
			case "newHeads":
				ch, stop, err := s.be.SubscribeNewHeads(ctx)
				if err != nil {
					_ = conn.WriteJSON(Types.RespErr(req.ID, -32000, err.Error()))
					continue
				}
				storeSub(subs, mu, sid, stop)
				_ = conn.WriteJSON(Types.RespOK(req.ID, sid))
				go forwardBlocks(conn, sid, ch)

			case "logs":
				var q Types.FilterQuery
				if len(req.Params) > 1 {
					if qq, err := toFilterQuery(req.Params[1]); err == nil {
						q = *qq
					}
				}
				ch, stop, err := s.be.SubscribeLogs(ctx, &q)
				if err != nil {
					_ = conn.WriteJSON(Types.RespErr(req.ID, -32000, err.Error()))
					continue
				}
				storeSub(subs, mu, sid, stop)
				_ = conn.WriteJSON(Types.RespOK(req.ID, sid))
				go forwardLogs(conn, sid, ch)

			case "newPendingTransactions":
				ch, stop, err := s.be.SubscribePendingTxs(ctx)
				if err != nil {
					_ = conn.WriteJSON(Types.RespErr(req.ID, -32000, err.Error()))
					continue
				}
				storeSub(subs, mu, sid, stop)
				_ = conn.WriteJSON(Types.RespOK(req.ID, sid))
				go forwardPending(conn, sid, ch)

			default:
				_ = conn.WriteJSON(Types.RespErr(req.ID, -32602, "unsupported subscription"))
			}
			continue
		}

		if req.Method == "eth_unsubscribe" {
			if len(req.Params) < 1 {
				_ = conn.WriteJSON(Types.RespErr(req.ID, -32602, "missing id"))
				continue
			}
			id, _ := req.Params[0].(string)
			mu.Lock()
			if s, ok := subs[id]; ok {
				s.stop()
				delete(subs, id)
			}
			mu.Unlock()
			_ = conn.WriteJSON(Types.RespOK(req.ID, true))
			continue
		}

		// regular RPC via WS
		resp, _ := s.h.Handle(ctx, req)
		log.Printf("📤 WS Response: %s -> %+v", req.Method, resp)
		_ = conn.WriteJSON(resp)
	}
}

func storeSub(m map[string]*sub, mu sync.Mutex, id string, stop func()) {
	mu.Lock()
	defer mu.Unlock()
	m[id] = &sub{id: id, stop: stop}
}

type subMsg struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription string `json:"subscription"`
		Result       any    `json:"result"`
	} `json:"params"`
}

func forwardBlocks(conn *websocket.Conn, sid string, ch <-chan *Types.Block) {
	for b := range ch {
		msg := subMsg{Jsonrpc: "2.0", Method: "eth_subscription"}
		msg.Params.Subscription = sid
		msg.Params.Result = marshalBlock(b, false)
		_ = conn.WriteJSON(msg)
	}
}
func forwardLogs(conn *websocket.Conn, sid string, ch <-chan *Types.Log) {
	for l := range ch {
		msg := subMsg{Jsonrpc: "2.0", Method: "eth_subscription"}
		msg.Params.Subscription = sid
		msg.Params.Result = marshalLogs([]*Types.Log{l})[0]
		_ = conn.WriteJSON(msg)
	}
}
func forwardPending(conn *websocket.Conn, sid string, ch <-chan []byte) {
	for h := range ch {
		msg := subMsg{Jsonrpc: "2.0", Method: "eth_subscription"}
		msg.Params.Subscription = sid
		msg.Params.Result = "0x" + hex.EncodeToString(h)
		_ = conn.WriteJSON(msg)
	}
}
