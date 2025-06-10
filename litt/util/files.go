package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// SwapFileExtension is the file extension used for temporary swap files created during atomic writes.
const SwapFileExtension = ".swap"

// DeleteOrphanedSwapFiles deletes any swap files in the given directory, i.e. files that end with ".swap".
func DeleteOrphanedSwapFiles(directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", directory, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == SwapFileExtension {
			swapFilePath := filepath.Join(directory, entry.Name())
			if err := os.Remove(swapFilePath); err != nil {
				return fmt.Errorf("failed to remove swap file %s: %w", swapFilePath, err)
			}
		}
	}

	return nil
}

// SanitizePath returns a sanitized version of the given path, doing things like expanding
// "~" to the user's home directory, converting to absolute path, normalizing slashes, etc.
func SanitizePath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}

		if len(path) == 1 {
			path = homeDir
		} else if len(path) > 1 && path[1] == '/' {
			path = homeDir + path[1:]
		}
	}

	path = filepath.Clean(path)
	path = filepath.ToSlash(path)
	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return path, nil
}

// AtomicWrite writes data to a file atomically. The parent directory must exist and be writable.
// If the destination file already exists, it will be overwritten.
//
// This method creates a temporary swap file in the same directory as the destination, but with SwapFileExtension
// appended to the filename. If there is a crash during this method's execution, it may leave this swap file behind.
func AtomicWrite(destination string, data []byte) error {

	swapPath := destination + ".swap"

	// Write the data into the swap file.
	swapFile, err := os.Create(swapPath)
	if err != nil {
		return fmt.Errorf("failed to create swap file: %v", err)
	}

	_, err = swapFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to swap file: %v", err)
	}

	// Ensure the data in the swap file is fully written to disk.
	err = swapFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync swap file: %v", err)
	}

	err = swapFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close swap file: %v", err)
	}

	// Rename the swap file to the destination file.
	err = AtomicRename(swapPath, destination)
	if err != nil {
		return fmt.Errorf("failed to rename swap file: %v", err)
	}

	return nil
}

// AtomicRename renames a file from oldPath to newPath atomically.
func AtomicRename(oldPath string, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	parentDirectory := filepath.Dir(newPath)

	// Ensure that the rename is committed to disk.
	dirFile, err := os.Open(parentDirectory)
	if err != nil {
		return fmt.Errorf("failed to open parent directory %s: %w", parentDirectory, err)
	}

	err = dirFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync parent directory %s: %w", parentDirectory, err)
	}

	err = dirFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close parent directory %s: %w", parentDirectory, err)
	}

	return nil
}

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

// CopyDirectoryRecursively creates a deep copy of the directory tree rooted at source and writes it to destination.
// It preserves file permissions, timestamps, and properly handles symlinks.
//
// The function performs a recursive copy of all files and directories, maintaining the same
// relative path structure and file metadata. If the destination directory exists, it will
// merge the source content into it, potentially overwriting files with the same names.
//
// The function checks that the destination has appropriate write permissions before starting the copy.
// If the destination directory doesn't exist, it verifies the parent directory has appropriate permissions.
// For existing directories, it ensures they have write permissions before attempting to copy files into them.
func CopyDirectoryRecursively(source string, destination string) error {
	// Verify the destination is writable (or can be created)
	if err := verifyDirectoryWritable(filepath.Dir(destination)); err != nil {
		return err
	}

	return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		// Compute the path relative to source, then build the destination path
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path from %s to %s: %w", source, path, err)
		}
		target := filepath.Join(destination, rel)

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		switch {
		case d.IsDir():
			return ensureDirectoryExists(target, info.Mode())

		case (info.Mode() & os.ModeSymlink) != 0:
			return copySymlink(path, target)

		default:
			return copyRegularFile(path, target, info.Mode(), info.ModTime())
		}
	})
}

// verifyDirectoryWritable checks if a directory exists and is writable.
// Returns nil if the directory is writable, or an error explaining why it's not.
// If the directory doesn't exist but its parent is writable, returns nil.
func verifyDirectoryWritable(dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist, check parent permissions
			parentDir := filepath.Dir(dirPath)
			return verifyDirectoryWritable(parentDir)
		}
		return fmt.Errorf("failed to access path '%s': %w", dirPath, err)
	}

	// Path exists, verify it's a directory with write permissions
	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", dirPath)
	}

	if info.Mode()&0200 == 0 {
		return fmt.Errorf("directory '%s' is not writable", dirPath)
	}

	return nil
}

// ensureParentDirExists ensures the parent directory of the given path exists and is writable.
// Creates parent directories if they don't exist.
func ensureParentDirExists(path string) error {
	parentDir := filepath.Dir(path)

	// Check if parent exists
	info, err := os.Stat(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Create parent directories if they don't exist
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory %s: %w", parentDir, err)
			}
			return nil
		}
		return fmt.Errorf("failed to check parent directory %s: %w", parentDir, err)
	}

	// Parent exists, verify it's a directory with write permissions
	if !info.IsDir() {
		return fmt.Errorf("parent path %s is not a directory", parentDir)
	}

	if info.Mode()&0200 == 0 {
		return fmt.Errorf("parent directory %s is not writable", parentDir)
	}

	return nil
}

// copyRegularFile copies a regular file from src to dst, preserving permissions and timestamps.
func copyRegularFile(src string, dst string, fileMode os.FileMode, modTime time.Time) error {
	// Ensure parent directory exists
	if err := ensureParentDirExists(dst); err != nil {
		return err
	}

	// Open source file
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer in.Close()

	// If there is already a file at the destination, remove it.
	// This ensures we don't have issues with file permissions or existing symlinks
	if _, err := os.Stat(dst); err == nil {
		// File exists, remove it
		if err := os.Remove(dst); err != nil {
			return fmt.Errorf("failed to remove existing destination file %s: %w", dst, err)
		}
	}

	// Create destination file
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer out.Close()

	// Copy content
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file content from %s to %s: %w", src, dst, err)
	}

	// Preserve timestamps
	if err := os.Chtimes(dst, modTime, modTime); err != nil {
		return fmt.Errorf("failed to preserve timestamps for %s: %w", dst, err)
	}

	return nil
}

// copySymlink copies a symlink from src to dst, preserving the link destination.
func copySymlink(src string, dst string) error {
	// Ensure parent directory exists
	if err := ensureParentDirExists(dst); err != nil {
		return err
	}

	// Read the symlink target
	linkDest, err := os.Readlink(src)
	if err != nil {
		return fmt.Errorf("failed to read symlink %s: %w", src, err)
	}

	// Create the symlink
	if err := os.Symlink(linkDest, dst); err != nil {
		return fmt.Errorf("failed to create symlink at %s pointing to %s: %w", dst, linkDest, err)
	}

	return nil
}

// ensureDirectoryExists ensures a directory exists with the given permissions.
// If the directory already exists, it verifies it has write permissions.
func ensureDirectoryExists(dirPath string, mode os.FileMode) error {
	// Check if directory already exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist, create it
			if err := os.MkdirAll(dirPath, mode); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
			}
			return nil
		}
		return fmt.Errorf("failed to check directory %s: %w", dirPath, err)
	}

	// Directory exists, verify it's actually a directory and has write permissions
	if !info.IsDir() {
		return fmt.Errorf("path %s exists but is not a directory", dirPath)
	}

	if info.Mode()&0200 == 0 {
		return fmt.Errorf("directory %s is not writable", dirPath)
	}

	return nil
}
