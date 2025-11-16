package main

import (
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
)

// ExecStateManager controls the ES state on a dedicated OS thread
type ExecStateManager struct {
	previousState  uint32
	commandCh      chan uint32
	managerStopCh  chan struct{}
	managerStarted atomic.Bool
	rpcShutdownCh  chan struct{}
}

// Start launches the dedicated OS thread goroutine
func (m *ExecStateManager) Start() {
	m.commandCh = make(chan uint32)
	m.managerStopCh = make(chan struct{})

	go func() {
		// Lock goroutine to its current OS thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// We want to ensure that the manager is started
		// before we start processing commands.
		m.managerStarted.Store(true)

		for {
			select {
			case flags := <-m.commandCh:
				// Call Windows API on this thread
				ret, err := SetThreadExecutionState(flags | ES_CONTINUOUS)
				if err != nil {
					log.Printf("SetThreadExecutionState error: %v", err)
					atomic.StoreUint32(&m.previousState, 0)
					continue
				}
				// Please note that return value is the PREVIOUS state
				atomic.StoreUint32(&m.previousState, uint32(ret))

			case <-m.managerStopCh:
				// Clear state on exit
				if _, err := SetThreadExecutionState(ES_CONTINUOUS); err != nil {
					log.Printf("SetThreadExecutionState error during Stop: %v", err)
				}
				// We assume RPC Shutdown has already been requested,
				// so we don't need to save the previous state here.
				return
			}
		}
	}()
}

// Stop signals the ExecStateManager goroutine to exit
func (m *ExecStateManager) Stop() {
	if !m.managerStarted.Load() {
		log.Println("ERROR: ExecStateManager not started")
	} else {
		close(m.managerStopCh)
	}
}

// getAtomicState atomically returns the previous flags value
func (m *ExecStateManager) getAtomicState() uint32 {
	if !m.managerStarted.Load() {
		log.Println("ERROR: ExecStateManager not started")
	}
	return atomic.LoadUint32(&m.previousState)
}

// setAtomicState atomically sets the flags value
func (m *ExecStateManager) setAtomicState(flags uint32, reply *ExecStateReply) error {
	if !m.managerStarted.Load() {
		return fmt.Errorf("ExecStateManager not started")
	}
	m.commandCh <- flags
	reply.Flags = m.getAtomicState()
	return nil
}
