package main

import (
	"log"
	"math/big"

	jmdtgethfacade "github.com/JupiterMetaLabs/geth-facade/pkg/jmdtgethfacade"
	"github.com/JupiterMetaLabs/geth-facade/pkg/memorybackend"
)

func main() {
	// Create a simple server with memory backend for Sepolia testnet
	chainID := big.NewInt(11155111) // Sepolia testnet chain ID

	log.Println("Starting JMDT Geth Facade server...")
	log.Println("HTTP endpoint will be available at: http://localhost:8545")
	log.Println("WebSocket endpoint will be available at: ws://localhost:8546")
	log.Println("Press Ctrl+C to stop the server")

	// Create memory backend and start the server
	memoryBackend := memorybackend.NewMemoryBackend(chainID)
	if err := jmdtgethfacade.QuickStart(memoryBackend); err != nil {
		log.Fatal("Server error:", err)
	}
}
