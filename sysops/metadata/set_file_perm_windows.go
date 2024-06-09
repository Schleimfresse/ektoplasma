//go:build windows

package sysops

import "syscall"

// Set_file_perm sets the file permissions on Windows.
func Set_file_perm(path string, perm uint32) error {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(ptr, perm)
}
