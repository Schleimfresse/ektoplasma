//go:build windows

package sysops

import (
	"fmt"
	"strings"
	"syscall"
)

// Exec_cmd executes a command in the shell on Windows systems.
func Exec_cmd(command string, args []string) error {
	// Prepare the command and arguments
	argv := strings.Join(args, " ")
	cmdLine, err := syscall.UTF16PtrFromString("cmd.exe /C " + command + " " + argv)
	if err != nil {
		return err
	}

	// Create structures for process and startup information
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation

	// Create the process
	err = syscall.CreateProcess(
		nil,     // Application name
		cmdLine, // Command line
		nil,     // Process attributes
		nil,     // Thread attributes
		false,   // Inherit handles
		0,       // Creation flags
		nil,     // Environment
		nil,     // Current directory
		&si,     // Startup information
		&pi,     // Process information
	)
	if err != nil {
		return fmt.Errorf("failed to create process: %v", err)
	}
	defer syscall.CloseHandle(pi.Process)
	defer syscall.CloseHandle(pi.Thread)

	// Wait for the process to complete
	_, err = syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	if err != nil {
		return fmt.Errorf("failed to wait for process: %v", err)
	}

	return nil
}
