//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"syscall"
)

// DeleteDirectory deletes a directory on Unix-like systems.
func RmDir(path string) error {
	return syscall.Rmdir(path)
}
