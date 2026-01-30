package main

import (
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Disable logging for tests to keep output clean.
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

// setupTestServer initializes a new ExecStateManager, starts a listener on a random port,
// and returns the manager, listener, and a client connected to the server.
func setupTestServer(t *testing.T) (*ExecStateManager, net.Listener, *rpc.Client) {
	t.Helper()

	// Ensure we have a fresh RPC server for each test to avoid registration conflicts.
	rpc.DefaultServer = rpc.NewServer()

	// Use an in-memory listener for testing
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	manager := &ExecStateManager{
		listener: listener,
	}
	manager.Start()

	if err := rpc.Register(manager); err != nil {
		t.Fatalf("rpc.Register failed: %v", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return // Listener closed
			}
			go rpc.ServeConn(conn)
		}
	}()

	client, err := rpc.Dial(listener.Addr().Network(), listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial RPC server: %v", err)
	}

	return manager, listener, client
}

func TestRPCMethods(t *testing.T) {
	manager, listener, client := setupTestServer(t)
	defer listener.Close()
	defer client.Close()
	defer manager.Stop()

	req := ExecStateRequest{}

	// Set an initial state to have a predictable starting point.
	var initialReply ExecStateReply
	if err := client.Call("ExecStateManager.Clear", req, &initialReply); err != nil {
		t.Fatalf("Initial RPC call to Clear failed: %v", err)
	}

	testCases := []struct {
		name       string
		method     string
		stateToSet uint32
	}{
		{"System", "ExecStateManager.System", ES_SYSTEM_REQUIRED | ES_CONTINUOUS},
		{"Display", "ExecStateManager.Display", ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED | ES_CONTINUOUS},
		{"Critical", "ExecStateManager.Critical", ES_SYSTEM_REQUIRED | ES_AWAYMODE_REQUIRED | ES_CONTINUOUS},
		{"Clear", "ExecStateManager.Clear", ES_CONTINUOUS},
	}

	// The first operation is Clear, so the previous state is what Clear returns.
	previousFlag := initialReply.Flags

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var currentReply ExecStateReply
			err := client.Call(tc.method, req, &currentReply)
			if err != nil {
				t.Fatalf("RPC call to %s failed: %v", tc.method, err)
			}
			// The reply contains the state from the *previous* call.
			if currentReply.Flags != previousFlag {
				t.Errorf("Expected previous flags to be 0x%X, but got 0x%X", previousFlag, currentReply.Flags)
			}
			previousFlag = tc.stateToSet
		})
	}
}

func TestRPCShutdown(t *testing.T) {
	manager, listener, client := setupTestServer(t)
	defer manager.Stop()

	var reply ExecStateReply
	req := ExecStateRequest{}

	// Call Shutdown
	err := client.Call("ExecStateManager.Shutdown", req, &reply)
	if err != nil && err != rpc.ErrShutdown {
		// We might get rpc.ErrShutdown which is fine. Any other error is a problem.
		t.Fatalf("Shutdown RPC call failed with unexpected error: %v", err)
	}

	// The listener should be closed now. Further accepts should fail.
	listener.Close() // Close the listener to ensure the server stops accepting new connections.
	client.Close()   // Close the client to ensure it's not connected anymore.

	// Existing client should see a connection error.
	err = client.Call("ExecStateManager.Read", req, &reply)
	if err != rpc.ErrShutdown {
		t.Errorf("Expected an error when calling RPC on a shutdown server, but got %v", err)
	}
}
