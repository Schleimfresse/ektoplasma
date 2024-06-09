//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

// Dir_exist checks if a directory exists on Unix-like systems.
func Dir_exist(path string) bool {
	// Use syscall.Stat to check if the directory exists
	_, err := syscall.Stat(path, nil)
	return err == nil
}
