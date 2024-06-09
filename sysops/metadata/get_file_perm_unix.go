//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import "syscall"

// Get_file_perm retrieves the file permissions on Unix-like systems.
func Get_file_perm(path string) (uint32, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return 0, err
	}
	return stat.Mode & 0777, nil // mask to get only permission bits
}
