//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

// Exec_cmd executes a command in the shell on Unix systems.
func Exec_cmd(command string) error {
	// Prepare the command and arguments
	args := []string{"sh", "-c", command}

	// Fork a child process to execute the command
	pid, err := syscall.ForkExec("/bin/sh", args, &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr)},
	})
	if err != nil {
		return fmt.Errorf("failed to fork/exec command: %v", err)
	}

	// Wait for the child process to complete
	var ws syscall.WaitStatus
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return fmt.Errorf("failed to wait for process: %v", err)
	}

	if !ws.Exited() {
		return fmt.Errorf("command did not exit normally")
	}

	return nil
}
