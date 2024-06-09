//go:build windows

package sysops

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

func Set_file_owner(path string, owner string) error {
	usernamePtr, err := windows.UTF16PtrFromString(owner)
	if err != nil {
		return fmt.Errorf("failed to convert username to UTF-16: %v", err)
	}

	// Buffers for SID and domain name
	var sid *windows.SID
	var domainName [256]uint16
	var sidSize uint32
	var domainNameSize uint32 = uint32(len(domainName))
	var use uint32

	// First call to LookupAccountName to get the buffer size
	err = windows.LookupAccountName(nil, usernamePtr, nil, &sidSize, &domainName[0], &domainNameSize, &use)
	if err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return fmt.Errorf("failed to lookup account name: %v", err)
	}

	// Allocate the buffer for the SID
	sid = (*windows.SID)(unsafe.Pointer(&make([]byte, sidSize)[0]))

	err = windows.LookupAccountName(nil, usernamePtr, sid, &sidSize, &domainName[0], &domainNameSize, &use)
	if err != nil {
		return fmt.Errorf("failed to lookup account name: %v", err)
	}

	// Use SetNamedSecurityInfo to set the file's owner
	err = windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION,
		sid,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
