package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// VerifyFilePermissions checks if a file has read/write permissions and is a regular file (if it exists),
// returning an error if it does not if the file permissions or file type is not as expected.
// Also returns a boolean indicating if the file exists and its size (to save on additional os.Stat calls).
//
// A file is considered to have the correct permissions/type if:
// - it exists and is a standard file with read+write permissions
// - if it does not exist but its parent directory has read+write permissions.
//
// The arguments for the function are the result of os.Stat(path). There is no need to do error checking on the
// result of os.Stat in the calling context (this method does it for you).
func VerifyFilePermissions(path string) (exists bool, size int64, err error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// The file does not exist. Check the parent.
			parentPath := filepath.Dir(path)
			parentInfo, err := os.Stat(parentPath)
			if err != nil {
				if os.IsNotExist(err) {
					return false, -1, fmt.Errorf("parent directory %s does not exist", parentPath)
				}
				return false, -1, fmt.Errorf("failed to stat parent directory %s: %w", parentPath, err)
			}

			if !parentInfo.IsDir() {
				return false, -1, fmt.Errorf("parent directory %s is not a directory", parentPath)
			}

			if parentInfo.Mode()&0700 != 0700 {
				return false, -1, fmt.Errorf("parent directory %s has insufficent permissions", parentPath)
			}
		}

		return false, 0, nil
	}

	// File exists. Check if it is a regular file and that it is readable+writeable.
	if info.IsDir() {
		return false, -1, fmt.Errorf("file %s is a directory", path)
	}
	if info.Mode()&0600 != 0600 {
		return false, -1, fmt.Errorf("file %s has insufficent permissions", path)
	}

	return true, info.Size(), nil
}
