//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import "syscall"

func File_open(path string, flag int, perm uint32) (int, error) {
	// Use syscall.Open to open the file
	fd, err := syscall.Open(path, flag, perm)
	return fd, err
}
