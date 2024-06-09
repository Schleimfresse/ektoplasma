//go:build windows

package sysops

import "syscall"

// Create_file creates or opens a file on Windows.
func Create_file(path string, flag int, perm uint32) (int, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return -1, err
	}

	// Use syscall.CreateFile to create or open the file
	handle, err := syscall.CreateFile(ptr, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return -1, err
	}
	return int(handle), nil
}
