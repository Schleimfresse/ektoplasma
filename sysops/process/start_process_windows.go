//go:build windows

package sysops

import (
	"fmt"
	"strings"
	"syscall"
)

func Start_proc(cmd string, args []string) (int, error) {

	// Convert command and arguments to C-style strings
	cmdStr, err := syscall.UTF16PtrFromString(cmd)
	if err != nil {
		return 0, fmt.Errorf("failed to convert command to UTF-16: %v", err)
	}
	argStr, err := syscall.UTF16PtrFromString(strings.Join(args, " "))
	if err != nil {
		return 0, fmt.Errorf("failed to convert arguments to UTF-16: %v", err)
	}

	// Initialize startup information and process information structures
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation

	// Start the process
	err = syscall.CreateProcess(
		cmdStr, // Command
		argStr, // Arguments
		nil,    // Process attributes
		nil,    // Thread attributes
		false,  // Inherit handles
		0,      // Creation flags
		nil,    // Environment
		nil,    // Current directory
		&si,    // Startup information
		&pi,    // Process information
	)
	if err != nil {
		return 0, fmt.Errorf("failed to start process: %v", err)
	}

	return int(pi.ProcessId), nil
}
