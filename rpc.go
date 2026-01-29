package main

import (
	"log"
)

// Request types for RPC
type ExecStateRequest struct {
	Flags uint32
}

type ExecStateReply struct {
	Flags uint32
}

// IMPORTANT: All methods return error to comply with net/rpc requirements

// Clears all sleep flags and returns the previous flags in the reply.
func (m *ExecStateManager) Clear(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Clear — Clearing sleep flags")
	return m.setAtomicState(0, reply)
}

// Sets the execution state to keep the system and display on, and returns the previous flags.
func (m *ExecStateManager) Display(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Display — Forcing display ON")
	return m.setAtomicState(ES_SYSTEM_REQUIRED|ES_DISPLAY_REQUIRED, reply)
}

// Sets the execution state to keep the system on, and returns the previous flags.
func (m *ExecStateManager) System(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.System — Forcing system ON")
	return m.setAtomicState(ES_SYSTEM_REQUIRED, reply)
}

// Sets the execution state to keep the system on and enable away mode, and returns the previous flags.
func (m *ExecStateManager) Critical(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Critical — Forcing system critical ON")
	return m.setAtomicState(ES_SYSTEM_REQUIRED|ES_AWAYMODE_REQUIRED, reply)
}

// Returns the previous execution state flags in the reply.
func (m *ExecStateManager) Read(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Read — Returning previous flags")
	reply.Flags = m.getAtomicState()
	return nil
}

// Shuts down the RPC server
func (m *ExecStateManager) Shutdown(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Shutdown - Shutting down RPC server")
	close(m.rpcShutdownCh)
	return nil
}
