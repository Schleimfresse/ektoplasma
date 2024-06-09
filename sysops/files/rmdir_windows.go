//go:build windows

package sysops

import (
	"syscall"
)

// DeleteDirectory deletes a directory on Windows.
func RmDir(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	// Call RemoveDirectory function from kernel32
	if err := syscall.RemoveDirectory(pathPtr); err != nil {
		return err
	}
	return nil
}
