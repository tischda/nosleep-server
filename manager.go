package main

import (
	"log"
	"runtime"
	"sync/atomic"
)

// ExecStateManager controls the ES state on a dedicated OS thread
type ExecStateManager struct {
	previousState uint32
	commandCh     chan uint32
	mgrShutdownCh chan struct{}
	rpcShutdownCh chan struct{}
}

// Start launches the dedicated OS thread goroutine
func (m *ExecStateManager) Start() {
	m.commandCh = make(chan uint32)
	m.mgrShutdownCh = make(chan struct{})

	go func() {
		// Lock goroutine to its current OS thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

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
			case <-m.mgrShutdownCh:
				return
			}
		}
	}()
}

// Clears state. This function is meant to be called via defer() right after Start().
func (m *ExecStateManager) Stop() {
	close(m.mgrShutdownCh)

	if _, err := SetThreadExecutionState(ES_CONTINUOUS); err != nil {
		log.Printf("SetThreadExecutionState error during Stop: %v", err)
	}
	log.Println("ThreadExecutionState cleared.")
}

// getAtomicState atomically returns the previous flags value
func (m *ExecStateManager) getAtomicState() uint32 {
	return atomic.LoadUint32(&m.previousState)
}

// setAtomicState atomically sets the flags value
func (m *ExecStateManager) setAtomicState(flags uint32, reply *ExecStateReply) {
	m.commandCh <- flags
	reply.Flags = m.getAtomicState()
}
