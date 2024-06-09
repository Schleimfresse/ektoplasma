//go:build windows

package sysops

import (
	"errors"
	"fmt"
	"syscall"
)

func Get_env(key string) (string, error) {
	// Get the environment variable value using syscall.GetEnvironmentVariable on Windows
	keyptr, err := syscall.UTF16PtrFromString(key)
	if err != nil {
		return "", err
	}
	var buf [syscall.MAX_PATH]uint16
	_, err = syscall.GetEnvironmentVariable(keyptr, &buf[0], syscall.MAX_PATH)
	if err != nil {
		if errors.Is(err, syscall.ERROR_ENVVAR_NOT_FOUND) {
			return "", fmt.Errorf("environment variable not found")
		}
		return "", fmt.Errorf("failed to get environment variable: %v", err)
	}

	return syscall.UTF16ToString(buf[:]), nil
}
