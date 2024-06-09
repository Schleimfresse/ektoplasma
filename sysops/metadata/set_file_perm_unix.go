//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

// Set_file_perm sets the file permissions on Unix-like systems.
func Set_file_perm(path string, perm uint32) error {
	return syscall.Chmod(path, perm)
}
