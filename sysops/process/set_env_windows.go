//go:build windows

package sysops

import (
	"fmt"
	"syscall"
)

func Set_env(key, value string) error {
	// Set the environment variable value using syscall.SetEnvironmentVariable on Windows
	keyptr, err := syscall.UTF16PtrFromString(key)
	if err != nil {
		return err
	}
	valueptr, err := syscall.UTF16PtrFromString(value)
	err = syscall.SetEnvironmentVariable(keyptr, valueptr)
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}
	return nil
}
