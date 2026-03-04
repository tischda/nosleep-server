package main

import (
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
)

type execStateCommand struct {
	flags   uint32
	errChan chan error
}

// ExecStateManager controls the ES state on a dedicated OS thread
type ExecStateManager struct {
	previousState uint32
	commandCh     chan execStateCommand
	mgrShutdownCh chan struct{}
	listener      net.Listener
	processesMu   sync.Mutex
	processes     map[int]struct{}
}

// Start launches the dedicated OS thread goroutine
func (m *ExecStateManager) Start() {
	m.commandCh = make(chan execStateCommand)
	m.mgrShutdownCh = make(chan struct{})
	if m.processes == nil {
		m.processes = make(map[int]struct{})
	}

	go func() {
		// Lock goroutine to its current OS thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		for {
			select {
			case cmd := <-m.commandCh:
				// Call Windows API on this thread
				ret, err := SetThreadExecutionState(cmd.flags | ES_CONTINUOUS)
				if err != nil {
					log.Printf("SetThreadExecutionState error: %v", err)
					atomic.StoreUint32(&m.previousState, 0)
				} else {
					// Please note that return value is the PREVIOUS state
					atomic.StoreUint32(&m.previousState, uint32(ret))
				}
				cmd.errChan <- err
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
func (m *ExecStateManager) setAtomicState(flags uint32, reply *ExecStateReply) error {
	errChan := make(chan error)
	m.commandCh <- execStateCommand{flags: flags, errChan: errChan}
	err := <-errChan
	reply.Flags = m.getAtomicState()
	return err
}

func (m *ExecStateManager) getRegisteredProcesses() []int {
	m.processesMu.Lock()
	defer m.processesMu.Unlock()

	var pids []int
	for pid := range m.processes {
		pids = append(pids, pid)
	}
	return pids
}

func (m *ExecStateManager) registerProcess(pid int) {
	m.processesMu.Lock()
	defer m.processesMu.Unlock()

	m.processes[pid] = struct{}{}
}

func (m *ExecStateManager) unregisterProcess(pid int) {
	m.processesMu.Lock()
	defer m.processesMu.Unlock()

	delete(m.processes, pid)
}

func (m *ExecStateManager) hasRegisteredProcesses() bool {
	m.processesMu.Lock()
	defer m.processesMu.Unlock()

	return len(m.processes) > 0
}
