//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import "syscall"

// Is_dir checks if a given path is a directory on Unix-like systems.
func Is_dir(path string) (bool, error) {
	var stat syscall.Stat_t

	// Use syscall.Stat to get file status
	err := syscall.Stat(path, &stat)
	if err != nil {
		return false, err
	}

	// Check if the mode indicates a directory
	return stat.Mode&syscall.S_IFDIR != 0, nil
}
