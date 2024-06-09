package sysops

import (
	"golang.org/x/sys/windows"
)

func Kill_process(pid int) error {
	// Open the process handle
	handle, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return err
	}

	// Terminate the process
	err = windows.TerminateProcess(handle, 1)
	if err != nil {
		return err
	}
	return nil
}
