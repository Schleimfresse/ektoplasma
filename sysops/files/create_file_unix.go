//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import "syscall"

// Create_file creates or opens a file on Unix-like systems.
func Create_file(path string, flag int, perm uint32) (int, error) {
	// Use syscall.Open to create or open the file
	fd, err := syscall.Open(path, flag|syscall.O_CREAT, perm)
	return fd, err
}
