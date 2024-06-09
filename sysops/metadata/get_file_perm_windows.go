//go:build windows

package sysops

import "syscall"

// Get_file_perms retrieves the file permissions on Windows.
func Get_file_perms(path string) (uint32, error) {
	var sa syscall.SecurityAttributes
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	handle, err := syscall.CreateFile(
		ptr,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		&sa,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(handle)

	var fileInfo syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(handle, &fileInfo)
	if err != nil {
		return 0, err
	}
	return fileInfo.FileAttributes, nil
}
