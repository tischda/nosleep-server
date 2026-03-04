package main

import (
	"log"
)

// Request types for RPC (make sure to keep them in sync with the client)
type ExecStateRequest struct {
	Process int
}

type ExecStateReply struct {
	Flags     uint32
	Processes []int
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
	reply.Processes = m.getRegisteredProcesses()
	return nil
}

// Registers a process.
func (m *ExecStateManager) Register(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Register — Register process:", req.Process)
	m.registerProcess(req.Process)
	return nil
}

// Unregisters a process.
func (m *ExecStateManager) Unregister(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Unregister — Unregister process:", req.Process)
	m.unregisterProcess(req.Process)
	if !m.hasRegisteredProcesses() {
		log.Println("ExecStateManager.Unregister — All processes unregistered")
		return m.Shutdown(req, reply)
	}
	return nil
}

// Shuts down the RPC server.
func (m *ExecStateManager) Shutdown(req ExecStateRequest, reply *ExecStateReply) error {
	log.Println("ExecStateManager.Shutdown - Shutting down RPC server")

	// Close the listener to stop accepting new connections, assuming
	// ExcecStateManager.Stop will be called via defer in main()
	return m.listener.Close()
}
