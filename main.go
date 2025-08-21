package main

import (
	"flag"
	"log"
	"math/big"
	"strings"

	"jmdt-geth-facade/backend"
	"jmdt-geth-facade/rpc"
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

	// TODO: swap memory backend with your real adapter
	be := backend.NewMemoryBackend(chainID)

	h := rpc.NewHandlers(be)

	go func() {
		log.Println("HTTP JSON-RPC on", *httpAddrFlag)
		if err := rpc.NewHTTPServer(h).Serve(*httpAddrFlag); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("WS JSON-RPC on", *wsAddrFlag)
	if err := rpc.NewWSServer(h, be).Serve(*wsAddrFlag); err != nil {
		log.Fatal(err)
	}
}
