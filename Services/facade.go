// Package jmdtgethfacade provides a minimal, production-ready JSON-RPC/WebSocket façade
// that mirrors common geth endpoints (eth_*, net_*, web3_*).
//
// This package allows you to wire your own backend implementation to provide
// real blockchain data through a standard Ethereum JSON-RPC interface.
package Services

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jupitermetalabs/geth-facade/Types"
)

// Server represents a geth facade server that can serve both HTTP and WebSocket JSON-RPC endpoints.
type Server struct {
	handlers *Handlers
	backend  Types.Backend
	httpAddr string
	wsAddr   string
}

// Config holds the configuration for the facade server.
type Config struct {
	// Backend is the blockchain backend implementation
	Backend Types.Backend
	// HTTPAddr is the HTTP server address (e.g., ":8545")
	HTTPAddr string
	// WSAddr is the WebSocket server address (e.g., ":8546")
	WSAddr string
}

// NewServer creates a new facade server with the given configuration.
func NewServer(config Config) *Server {
	return &Server{
		handlers: NewHandlers(config.Backend),
		backend:  config.Backend,
		httpAddr: config.HTTPAddr,
		wsAddr:   config.WSAddr,
	}
}

// Start starts both HTTP and WebSocket servers.
// This method blocks until one of the servers encounters an error.
func (s *Server) Start() error {
	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP JSON-RPC server starting on %s", s.httpAddr)
		httpServer := NewHTTPServer(s.handlers)
		if err := httpServer.Serve(s.httpAddr); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start WebSocket server (blocking)
	log.Printf("WebSocket JSON-RPC server starting on %s", s.wsAddr)
	wsServer := NewWSServer(s.handlers, s.backend)
	return wsServer.Serve(s.wsAddr)
}

// StartHTTP starts only the HTTP server.
func (s *Server) StartHTTP() error {
	log.Printf("HTTP JSON-RPC server starting on %s", s.httpAddr)
	httpServer := NewHTTPServer(s.handlers)
	return httpServer.Serve(s.httpAddr)
}

// StartWS starts only the WebSocket server.
func (s *Server) StartWS() error {
	log.Printf("WebSocket JSON-RPC server starting on %s", s.wsAddr)
	wsServer := NewWSServer(s.handlers, s.backend)
	return wsServer.Serve(s.wsAddr)
}

// GetHandlers returns the RPC handlers for custom server implementations.
func (s *Server) GetHandlers() *Handlers {
	return s.handlers
}

// GetBackend returns the backend implementation.
func (s *Server) GetBackend() Types.Backend {
	return s.backend
}

// DefaultConfig returns a default configuration.
// Note: You must provide your own backend implementation.
func DefaultConfig() Config {
	return Config{
		Backend:  nil, // Must be set by user
		HTTPAddr: ":8545",
		WSAddr:   ":8546",
	}
}

// DefaultConfigWithCustomBackend returns a default configuration with a custom Types.
func DefaultConfigWithCustomBackend(be Types.Backend) Config {
	return Config{
		Backend:  be,
		HTTPAddr: ":8545",
		WSAddr:   ":8546",
	}
}

// QuickStart starts a server with default configuration.
// Note: You must provide your own backend implementation.
// For testing, see the examples/memory-backend directory.
func QuickStart(be Types.Backend) error {
	config := DefaultConfigWithCustomBackend(be)
	server := NewServer(config)
	return server.Start()
}

// QuickStartWithBackend starts a server with default configuration and custom Types.
func QuickStartWithBackend(be Types.Backend) error {
	config := DefaultConfigWithCustomBackend(be)
	server := NewServer(config)
	return server.Start()
}

// HealthCheck provides a simple health check endpoint.
// Note: This is now handled by the Gin HTTP server directly
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check if backend is responsive
	_, err := s.backend.BlockNumber(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Backend not available"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ReadyCheck provides a readiness check endpoint.
// Note: This is now handled by the Gin HTTP server directly
func (s *Server) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check if backend is ready
	_, err := s.backend.ChainID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Backend not ready"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
