package main

import (
	"log"
	"syscall"
)

// https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-setthreadexecutionstate
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
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	setThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

func ForceDisplayOn() {
	ret, _, err := setThreadExecutionState.Call(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
	if ret == 0 {
		log.Println("Failed to force display on:", err)
	}
}

func ForceSystemOn() {
	ret, _, err := setThreadExecutionState.Call(ES_CONTINUOUS | ES_SYSTEM_REQUIRED)
	if ret == 0 {
		log.Println("Failed to force system on:", err)
	}
}

func ForceSystemCriticalOn() {
	ret, _, err := setThreadExecutionState.Call(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_AWAYMODE_REQUIRED)
	if ret == 0 {
		log.Println("Failed to force system critical on:", err)
	}
}

func ClearSleepFlags() {
	ret, _, err := setThreadExecutionState.Call(ES_CONTINUOUS)
	if ret == 0 {
		log.Println("Failed to clean sleep flags:", err)
	}
}

// WARNING: Microsoft does not provide an API to reliably read the currentSetThreadExecutionState
// flags. The documentation states that calling the function with zero doesn't set any state,
// but returns the prior value, which is not always meaningful.
func ReadFlags() uint32 {
	ret, _, err := setThreadExecutionState.Call(0)
	if ret == 0 {
		log.Println("Failed to retrieve flag value:", err)
	}
	return uint32(ret)
}
