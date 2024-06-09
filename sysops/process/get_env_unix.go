//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"fmt"
	"syscall"
)

func Get_env(key string) (string, error) {
	// Get the environment variable value using syscall.Getenv on Unix
	value, found := syscall.Getenv(key)
	if !found {
		return "", fmt.Errorf("environment variable not found")
	}
	return value, nil
}
