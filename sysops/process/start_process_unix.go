//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"fmt"
	"syscall"
)

func Start_proc(cmd string, args []string) (int, error) {
	// Start the process using syscall.ForkExec on Unix
	pid, err := syscall.ForkExec(cmd, args, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to start process: %v", err)
	}
	return pid, nil
}
