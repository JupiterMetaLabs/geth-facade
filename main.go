package main

import (
	"flag"
	"log"
	"math/big"
	"strings"

	"github.com/jupitermetalabs/geth-facade/pkg/jmdtgethfacade"
	"github.com/jupitermetalabs/geth-facade/pkg/memorybackend"
)

func main() {
	// Flags
	chainIDFlag := flag.String("chainid", "11155111", "Chain ID in hex (e.g. 0xaa36a7) or decimal (e.g. 11155111)")
	httpAddrFlag := flag.String("http", ":8545", "HTTP listen address (e.g. :8545 or 0.0.0.0:8545)")
	wsAddrFlag := flag.String("ws", ":8546", "WebSocket listen address (e.g. :8546 or 0.0.0.0:8546)")
	flag.Parse()

	// Parse chain id
	var chainID = new(big.Int)
	if strings.HasPrefix(*chainIDFlag, "0x") || strings.HasPrefix(*chainIDFlag, "0X") {
		chainID.SetString((*chainIDFlag)[2:], 16)
	} else {
		chainID.SetString(*chainIDFlag, 10)
	}

	// Create server configuration with memory backend (for testing/development)
	config := jmdtgethfacade.Config{
		Backend:  memorybackend.NewMemoryBackend(chainID),
		HTTPAddr: *httpAddrFlag,
		WSAddr:   *wsAddrFlag,
	}

	// Create and start server
	server := jmdtgethfacade.NewServer(config)

	log.Printf("Starting JMDT Geth Facade server...")
	log.Printf("Chain ID: %s", chainID.String())
	log.Printf("HTTP JSON-RPC on %s", *httpAddrFlag)
	log.Printf("WebSocket JSON-RPC on %s", *wsAddrFlag)

	if err := server.Start(); err != nil {
		log.Fatal("Server error:", err)
	}
}
