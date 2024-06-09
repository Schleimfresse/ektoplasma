//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"fmt"
	"syscall"
)

func Kill_process(pid int) error {
	// Send SIGKILL signal to the process on Unix
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}
