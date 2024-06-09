//go:build windows

package sysops

import "syscall"

// Dir_exists checks if a directory exists on Windows.
func Dir_exists(path string) bool {
	ptr, err2 := syscall.UTF16PtrFromString(path)
	if err2 != nil {
		return false
	}

	// Use syscall.GetFileAttributes to check if the directory exists
	attr, err := syscall.GetFileAttributes(ptr)
	return err == nil && attr&syscall.FILE_ATTRIBUTE_DIRECTORY != 0
}
