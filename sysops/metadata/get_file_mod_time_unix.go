//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"syscall"
	"time"
)

// get_file_mod_time retrieves the last modification time of the file.
func get_file_mod_time(path string) (time.Time, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return time.Time{}, err
	}
	return time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec), nil
}
