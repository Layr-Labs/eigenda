package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// VerifyFileProperties checks if a file has read/write permissions and is a regular file (if it exists),
// returning an error if it does not if the file permissions or file type is not as expected.
// Also returns a boolean indicating if the file exists and its size (to save on additional os.Stat calls).
//
// A file is considered to have the correct permissions/type if:
// - it exists and is a standard file with read+write permissions
// - if it does not exist but its parent directory has read+write permissions.
//
// The arguments for the function are the result of os.Stat(path). There is no need to do error checking on the
// result of os.Stat in the calling context (this method does it for you).
func VerifyFileProperties(path string) (exists bool, size int64, err error) {
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

			return false, -1, nil
		} else {
			return false, 0, fmt.Errorf("failed to stat path %s: %w", path, err)
		}
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

// Exists checks if a file or directory exists at the given path. More aesthetically pleasant than os.Stat.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("error checking if path %s exists: %w", path, err)
}

// CopyDirectoryRecursively creates a deep copy of the directory tree rooted at src and writes it to dst.
// It preserves file permissions, timestamps, and properly handles symlinks.
//
// The function performs a recursive copy of all files and directories, maintaining the same
// relative path structure and file metadata. If the destination directory exists, it will
// merge the source content into it, potentially overwriting files with the same names.
func CopyDirectoryRecursively(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		// Compute the path relative to src, then build the destination path
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path from %s to %s: %w", src, path, err)
		}
		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		switch {
		case d.IsDir():
			// Create directory (and parents) with same mode
			if err := os.MkdirAll(target, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
			return nil

		case (info.Mode() & os.ModeSymlink) != 0:
			// Replicate symlink
			linkDest, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("failed to read symlink %s: %w", path, err)
			}
			if err := os.Symlink(linkDest, target); err != nil {
				return fmt.Errorf("failed to create symlink at %s pointing to %s: %w", target, linkDest, err)
			}
			return nil

		default:
			// Regular file: copy contents
			in, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open source file %s: %w", path, err)
			}
			defer in.Close()

			out, err := os.OpenFile(target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				info.Mode())
			if err != nil {
				return fmt.Errorf("failed to create destination file %s: %w", target, err)
			}
			defer out.Close()

			if _, err := io.Copy(out, in); err != nil {
				return fmt.Errorf("failed to copy file content from %s to %s: %w", path, target, err)
			}
			
			// Preserve timestamps
			if err := os.Chtimes(target, info.ModTime(), info.ModTime()); err != nil {
				return fmt.Errorf("failed to preserve timestamps for %s: %w", target, err)
			}
			return nil
		}
	})
}
