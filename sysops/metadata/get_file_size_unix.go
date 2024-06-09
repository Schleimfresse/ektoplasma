//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

// Get_file_size retrieves the size of the file on Unix-like systems.
func Get_file_size(path string) (int64, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return 0, err
	}
	return stat.Size, nil
}
