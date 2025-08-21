package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HTTPServer struct {
	h *Handlers
}

func NewHTTPServer(h *Handlers) *HTTPServer { return &HTTPServer{h: h} }

func (s *HTTPServer) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleJSONRPC)
	srv := &http.Server{
		Addr:              addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *HTTPServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		write(w, RespErr(nil, -32700, "Parse error")); return
	}
	resp, _ := s.h.Handle(context.Background(), req)
	write(w, resp)
}

func write(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		if r.Method == http.MethodOptions { w.WriteHeader(204); return }
		next.ServeHTTP(w, r)
	})
}
