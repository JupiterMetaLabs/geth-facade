// Package jmdtgethfacade provides a minimal, production-ready JSON-RPC/WebSocket façade
// that mirrors common geth endpoints (eth_*, net_*, web3_*).
//
// This package allows you to wire your own backend implementation to provide
// real blockchain data through a standard Ethereum JSON-RPC interface.
//
// Quick Start
//
//	package main
//
//	import (
//		"log"
//		"math/big"
//
//		jmdtgethfacade "github.com/saishibu/jmdt-geth-facade/pkg/jmdtgethfacade"
//	)
//
//	func main() {
//		// You must provide your own backend implementation
//		// For testing, see examples/memory-backend
//		var myBackend backend.Backend = // your implementation
//		if err := jmdtgethfacade.QuickStart(myBackend); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// Custom Configuration
//
//	package main
//
//	import (
//		"log"
//		"math/big"
//
//		jmdtgethfacade "github.com/saishibu/jmdt-geth-facade/pkg/jmdtgethfacade"
//		"github.com/saishibu/jmdt-geth-facade/backend"
//	)
//
//	func main() {
//		// You must provide your own backend implementation
//		var myBackend backend.Backend = // your implementation
//		config := jmdtgethfacade.Config{
//			Backend:  myBackend,
//			HTTPAddr: ":8545",
//			WSAddr:   ":8546",
//		}
//
//		server := jmdtgethfacade.NewServer(config)
//		if err := server.Start(); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Custom Backend Implementation
//
// To integrate with your own blockchain node, implement the backend.Backend interface:
//
//	type MyBackend struct {
//		// Your blockchain client
//	}
//
//	func (b *MyBackend) ChainID(ctx context.Context) (*big.Int, error) {
//		// Return your chain ID
//	}
//
//	func (b *MyBackend) BlockNumber(ctx context.Context) (*big.Int, error) {
//		// Return current block number from your node
//	}
//
//	// ... implement all other required methods
//
// See the examples/ directory for complete examples.
package jmdtgethfacade
