//go:build windows

package sysops

import (
	"golang.org/x/sys/windows"
)

// Get_file_owner retrieves the owner (SID) of the file.
func Get_file_owner(path string) (string, error) {
	// Get file security information
	var sid *windows.SID
	sd, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION,
	)
	if err != nil {
		return "", err
	}
	sid, _, err = sd.Owner()
	if err != nil {
		return "", err
	}

	username, domain, _, err := sid.LookupAccount("")
	if err != nil {
		return "", err
	}

	return domain + "\\" + username, nil
}
