//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"fmt"
	"syscall"
)

func Set_env(key, value string) error {
	// Set the environment variable value using syscall.Setenv on Unix
	err := syscall.Setenv(key, value)
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}
	return nil
}
