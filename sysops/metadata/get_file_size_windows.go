//go:build windows

package sysops

import "syscall"

// Get_file_size retrieves the size of the file on Windows.
func Get_file_size(path string) (int64, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	handle, err := syscall.CreateFile(
		ptr,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(handle)

	var fileinfo syscall.ByHandleFileInformation

	err = syscall.GetFileInformationByHandle(handle, &fileinfo)
	if err != nil {
		return 0, err
	}

	size := int64(fileinfo.FileSizeHigh)<<32 + int64(fileinfo.FileSizeLow)
	return size, nil
}
