//go:build windows

package sysops

import "syscall"

// Is_dir checks if a given path is a directory on Windows.
func Is_dir(path string) (bool, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attr, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false, err
	}

	return attr&syscall.FILE_ATTRIBUTE_DIRECTORY != 0, nil
}
