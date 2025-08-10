package main

import (
	"syscall"
)

const (
	// Enables away mode. This value must be specified with ES_CONTINUOUS.
	// Away mode should be used only by media-recording and media-distribution
	// applications that must perform critical background processing on desktop
	// computers while the computer appears to be sleeping.
	ES_AWAYMODE_REQUIRED = 0x00000040

	// Informs the system that the state being set should remain in effect
	// until the next call that uses ES_CONTINUOUS and one of the other state
	// flags is cleared.
	ES_CONTINUOUS = 0x80000000

	// Forces the display to be on by resetting the display idle timer.
	ES_DISPLAY_REQUIRED = 0x00000002

	// Forces the system to be in the working state by resetting the system idle timer.
	ES_SYSTEM_REQUIRED = 0x00000001

	// This value is not supported. If ES_USER_PRESENT is combined with other esFlags
	// values, the call will fail and none of the specified states will be set.
	ES_USER_PRESENT = 0x00000004
)

var (
	modkernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecutionState = modkernel32.NewProc("SetThreadExecutionState")
)

// SetThreadExecutionState sets the thread's execution state using the Windows API.
// The flags parameter should be a combination of ES_CONTINUOUS, ES_SYSTEM_REQUIRED, ES_DISPLAY_REQUIRED, etc.
//
// If the function succeeds, the return value is the previous thread execution state.
// If the function fails, the return value is 0.
//
// See: https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-setthreadexecutionstate
func SetThreadExecutionState(flags uint32) (uint32, error) {
	ret, _, err := procSetThreadExecutionState.Call(uintptr(flags))
	if ret == 0 {
		return 0, err
	}
	return uint32(ret), nil
}
