//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

type UnixUser struct {
	Username string
	Uid      uint32
	Gid      uint32
}

func Set_file_owner(path string) error {
	passwdFd, err := syscall.Open("/etc/passwd", syscall.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to open /etc/passwd: %v", err)
	}
	defer syscall.Close(passwdFd)

	// Read the contents of /etc/passwd
	var passwdContents []byte
	buffer := make([]byte, 4096)
	for {
		n, err := syscall.Read(passwdFd, buffer)
		if err != nil && err != syscall.EINTR {
			return fmt.Errorf("failed to read /etc/passwd: %v", err)
		}
		if n == 0 {
			break
		}
		passwdContents = append(passwdContents, buffer[:n]...)
	}

	// Find the user information in /etc/passwd
	var user UnixUser
	userFound := false
	lines := syscall.SplitString(string(passwdContents), '\n')
	for _, line := range lines {
		fields := syscall.SplitString(line, ':')
		if len(fields) < 4 {
			continue
		}
		if fields[0] == username {
			user.Username = username
			fmt.Sscanf(fields[2], "%d", &user.Uid) // Parse UID
			fmt.Sscanf(fields[3], "%d", &user.Gid) // Parse GID
			userFound = true
			break
		}
	}

	if !userFound {
		return fmt.Errorf("user %s not found in /etc/passwd", username)
	}

	// Change the ownership of the specified file
	err = syscall.Chown(path, int(user.Uid), int(user.Gid))
	if err != nil {
		return fmt.Errorf("failed to change ownership of %s: %v", path, err)
	}

	return nil
}

// TODO linux rebuild & rewrite

// SplitString is a helper function that splits a string by a delimiter into a slice.
func SplitString(s string, delimiter rune) []string {
	var result []string
	current := make([]rune, 0, len(s))

	for _, r := range s {
		if r == delimiter {
			result = append(result, string(current))
			current = current[:0]
		} else {
			current = append(current, r)
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}
