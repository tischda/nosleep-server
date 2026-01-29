package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func serve(cfg *Config) {
	// Configure listener
	address := fmt.Sprintf("%s:%d", cfg.address, cfg.port)
	listener, err := net.Listen(cfg.network, address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	// Configure and start ExecStateManager
	manager := &ExecStateManager{listener: listener}
	manager.Start()
	defer manager.Stop()

	// set the initial sleep mode
	if cfg.display {
		if err := manager.Display(ExecStateRequest{}, &ExecStateReply{}); err != nil {
			log.Fatalf("Failed to set initial display state: %v", err)
		}
	} else {
		if err := manager.System(ExecStateRequest{}, &ExecStateReply{}); err != nil {
			log.Fatalf("Failed to set initial system state: %v", err)
		}
	}

	// Register RPC server with ExecStateManager methods
	if err := rpc.Register(manager); err != nil {
		log.Fatalf("Failed to register RPC server: %v", err)
	}

	log.Printf("RPC server listening on %s (%s)", address, cfg.network)
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Exit cleanly if the listener was intentionally closed.
			if errors.Is(err, net.ErrClosed) { // Go 1.16+
				break
			}
			// Other errors are real and should be logged/handled.
			log.Printf("accept error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
	log.Println("RPC server shutdown complete.")
}
