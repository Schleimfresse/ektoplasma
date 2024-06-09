//go:build linux || darwin || freebsd || openbsd || netbsd

package sysops

// Copy_dir copies a directory on Unix-like systems.
func Copy_dir(src, dst string) error {
	// Create destination directory if it doesn't exist
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	// Get list of files in the source directory
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each file and sub-directory recursively
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err = CopyDirectoryUnix(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
