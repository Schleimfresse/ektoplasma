package sysops

import (
	"syscall"
)

// TODO unix style
func List_env() []string {
	// List all environment variables using syscall.Environ on Windows
	envs := syscall.Environ()
	return envs
}
