//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"syscall"
	"time"
)

// SetFileModificationTimeUnix sets the last modification time of the file.
func Set_file_modification_time(path string, mtime time.Time) error {
	// Convert time.Time to syscall.Timespec
	times := []syscall.Timespec{
		syscall.Timespec{Sec: mtime.Unix(), Nsec: int64(mtime.Nanosecond())},
		syscall.Timespec{Sec: mtime.Unix(), Nsec: int64(mtime.Nanosecond())}, // Set atime and mtime to the same value
	}
	return syscall.UtimesNano(path, times)
}
