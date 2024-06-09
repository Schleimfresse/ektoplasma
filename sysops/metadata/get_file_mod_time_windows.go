//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

import (
	"syscall"
	"time"
)

// get_file_mod_time retrieves the last modification time of the file.
func get_file_mod_time(path string) (time.Time, error) {
	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(path),
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return time.Time{}, err
	}
	defer syscall.CloseHandle(handle)

	var fileInfo syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(handle, &fileInfo)
	if err != nil {
		return time.Time{}, err
	}
	ft := fileInfo.LastWriteTime
	return time.Unix(0, ft.Nanoseconds()), nil
}
