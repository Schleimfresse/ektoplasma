//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import "syscall"

// Get_file_owner retrieves the owner (UID and GID) of the file.
func Get_file_owner(path string) (owner string, err error) {
	var stat syscall.Stat_t
	err := syscall.Stat(filename, &stat)
	if err != nil {
		return "", err
	}

	// Extract
	ownerUID := stat.Uid

	// Get the owner's username
	pwd, err := syscall.Getpwuid(int(ownerUID))
	if err != nil {
		return "", err
	}
	return pwd.Name, nil
}
