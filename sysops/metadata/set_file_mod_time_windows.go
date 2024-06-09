//go:build windows

package sysops

import (
	"syscall"
	"time"
)

// Set_file_modification_time sets the last modification time of the file.
func Set_file_modification_time(path string, mtime time.Time) error {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	handle, err := syscall.CreateFile(
		ptr,
		syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return err
	}
	defer func(handle syscall.Handle) error {
		err := syscall.CloseHandle(handle)
		if err != nil {
			return err
		}
		return nil
	}(handle)

	ft := syscall.NsecToFiletime(mtime.UnixNano())
	return syscall.SetFileTime(handle, nil, nil, &ft)
}
