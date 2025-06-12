package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
func AtomicWrite(destination string, data []byte, fsync bool) error {

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

	if fsync {
		// Ensure the data in the swap file is fully written to disk.
		err = swapFile.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync swap file: %v", err)
		}
	}

	err = swapFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close swap file: %v", err)
	}

	// Rename the swap file to the destination file.
	err = AtomicRename(swapPath, destination, fsync)
	if err != nil {
		return fmt.Errorf("failed to rename swap file: %v", err)
	}

	return nil
}

// AtomicRename renames a file from oldPath to newPath atomically.
func AtomicRename(oldPath string, newPath string, fsync bool) error {
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

	if fsync {
		err = dirFile.Sync()
		if err != nil {
			return fmt.Errorf("failed to sync parent directory %s: %w", parentDirectory, err)
		}
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

// RecursiveMove transfers files/directory trees from the source to the destination.
//
// If delete is true, then the files at the source will be deleted when this method returns.
// If delete is false, then this function will leave behind a copy of the original files at the source.
//
// If deep is false, then this function will prefer hard-linking files instead of copying them. If the source and
// destination are on different filesystems, this will fall back to copying the files instead. If deep is true,
// then files are always copied, even if they are on the same filesystem. Deep=true also influences how symlinks
// are treated. If deep is false, then symlinks are copied as symlinks. If deep is true, then the file the symlink
// points to is copied instead.
//
// If preserveOriginal is true, then the original files at the source will be preserved after the move (they may
// still be hard linked if deep is false). If preserveOriginal is false, then the original files at the source will be
// deleted when this function returns.
//
// The fsync flag is intended to make this function faster for unit tests. If fsync is true, then the function will
// ensure that all file operations are fully flushed to disk before returning. If fsync is false, then the function
// will not perform any fsync operations, which may result in faster execution but less data safety in case of a crash.
func RecursiveMove(
	source string,
	destination string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
) error {
	// Sanitize paths
	source, err := SanitizePath(source)
	if err != nil {
		return fmt.Errorf("failed to sanitize source path: %w", err)
	}

	destination, err = SanitizePath(destination)
	if err != nil {
		return fmt.Errorf("failed to sanitize destination path: %w", err)
	}

	// Verify source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("source path %s does not exist: %w", source, err)
	}

	// Verify destination parent directory is writable
	if err := verifyDirectoryWritable(filepath.Dir(destination)); err != nil {
		return fmt.Errorf("destination parent directory not writable: %w", err)
	}

	// If source is a file, handle it directly
	if !sourceInfo.IsDir() {
		return recursiveMoveFile(source, destination, deep, preserveOriginal, fsync)
	}

	// Source is a directory, handle recursively
	return recursiveMoveDirectory(source, destination, deep, preserveOriginal, fsync)
}

// recursiveMoveFile handles moving a single file
func recursiveMoveFile(source string, destination string, deep bool, preserveOriginal bool, fsync bool) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Ensure parent directory exists
	if err := ensureParentDirExists(destination); err != nil {
		return fmt.Errorf("failed to ensure parent directory exists: %w", err)
	}

	// If not preserving original, try to move the file first (regardless of deep mode)
	if !preserveOriginal {
		// Try simple rename first (works if on same filesystem)
		if err := os.Rename(source, destination); err == nil {
			return nil
		}
		// Rename failed (likely different filesystem), fall back to copy+delete
	}

	// If preserving original or rename failed, use copy-based approach
	if preserveOriginal && !deep {
		// Try hard link if preserving original and not doing deep copy
		if err := os.Link(source, destination); err == nil {
			return nil
		}
		// Hard link failed, fall back to copy
	}

	// Check if source is a symlink
	if sourceInfo.Mode()&os.ModeSymlink != 0 {
		return handleSymlink(source, destination, deep, preserveOriginal, fsync)
	}

	// Copy the file
	if err := copyRegularFile(source, destination, sourceInfo.Mode(), sourceInfo.ModTime()); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Sync if requested
	if fsync {
		if err := syncFile(destination); err != nil {
			return fmt.Errorf("failed to sync destination file: %w", err)
		}
	}

	// Remove source if not preserving original
	if !preserveOriginal {
		if err := os.Remove(source); err != nil {
			return fmt.Errorf("failed to remove source file: %w", err)
		}
	}

	return nil
}

// recursiveMoveDirectory handles moving a directory and its contents
func recursiveMoveDirectory(
	source string,
	destination string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
) error {

	// Create destination directory if it doesn't exist
	if err := ensureDirectoryExists(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through source directory
	err := filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		// Skip the root directory itself
		if path == source {
			return nil
		}

		// Calculate relative path and destination path
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		destPath := filepath.Join(destination, relPath)

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", path, err)
		}

		switch {
		case d.IsDir():
			// Create directory at destination
			if err := ensureDirectoryExists(destPath, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}

		case (info.Mode() & os.ModeSymlink) != 0:
			// Handle symlink
			if err := handleSymlink(path, destPath, deep, preserveOriginal, fsync, source); err != nil {
				return fmt.Errorf("failed to handle symlink: %w", err)
			}

		default:
			// Handle regular file
			if !deep {
				// Try hard link first
				if err := ensureParentDirExists(destPath); err != nil {
					return fmt.Errorf("failed to ensure parent dir exists: %w", err)
				}

				if err := os.Link(path, destPath); err == nil {
					// Hard link succeeded
					if fsync {
						if err := syncFile(destPath); err != nil {
							return fmt.Errorf("failed to sync hard-linked file: %w", err)
						}
					}
					return nil
				}
				// Hard link failed, fall back to copy
			}

			// Copy the file
			if err := copyRegularFile(path, destPath, info.Mode(), info.ModTime()); err != nil {
				return fmt.Errorf("failed to copy regular file: %w", err)
			}

			if fsync {
				if err := syncFile(destPath); err != nil {
					return fmt.Errorf("failed to sync copied file: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Sync destination directory if requested
	if fsync {
		if err := syncDirectory(destination); err != nil {
			return fmt.Errorf("failed to sync destination directory: %w", err)
		}
	}

	// Remove source directory if not preserving original
	if !preserveOriginal {
		if err := os.RemoveAll(source); err != nil {
			return fmt.Errorf("failed to remove source directory: %w", err)
		}
	}

	return nil
}

// syncFile syncs a file to disk
func syncFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file for sync: %w", err)
	}
	defer file.Close()

	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// syncDirectory syncs a directory to disk
func syncDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open directory for sync: %w", err)
	}
	defer dir.Close()

	if err := dir.Sync(); err != nil {
		return fmt.Errorf("failed to sync directory: %w", err)
	}

	return nil
}

// handleSymlink handles symlink copying based on the deep parameter.
// If deep is false, copies the symlink as a symlink.
// If deep is true, copies the target file that the symlink points to.
// sourceRoot parameter is optional and used to check for nested directory conflicts.
func handleSymlink(
	source string,
	destination string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	sourceRoot ...string,
) error {

	if !deep {
		// Copy symlink as symlink
		return copySymlink(source, destination)
	}

	// Deep copy: follow the symlink and copy the target
	target, err := filepath.EvalSymlinks(source)
	if err != nil {
		return fmt.Errorf("failed to resolve symlink %s: %w", source, err)
	}

	// Check if target is within the source root (which would cause conflicts)
	if len(sourceRoot) > 0 {
		// Resolve symlinks and ensure both paths are absolute for comparison
		absSourceRoot, err1 := filepath.Abs(sourceRoot[0])
		realSourceRoot, err2 := filepath.EvalSymlinks(absSourceRoot)
		realTarget, err3 := filepath.EvalSymlinks(target)
		if err1 == nil && err2 == nil && err3 == nil {
			// Check if target is within the source root directory
			relPath, err := filepath.Rel(realSourceRoot, realTarget)
			if err == nil && !strings.HasPrefix(relPath, "..") && relPath != "." {
				// Target is within source root, this would cause conflicts
				return fmt.Errorf(
					"cannot deep copy symlink %s: target %s is within the source directory being moved",
					source, target)
			}
		}
	}

	// Get target file info
	targetInfo, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("failed to stat symlink target %s: %w", target, err)
	}

	if targetInfo.IsDir() {
		// Target is a directory, recursively copy it
		return recursiveMoveDirectory(target, destination, deep, preserveOriginal, fsync)
	} else {
		// Target is a file, copy it
		if err := copyRegularFile(target, destination, targetInfo.Mode(), targetInfo.ModTime()); err != nil {
			return fmt.Errorf("failed to copy symlink target: %w", err)
		}

		// Sync if requested
		if fsync {
			if err := syncFile(destination); err != nil {
				return fmt.Errorf("failed to sync copied symlink target: %w", err)
			}
		}

		// Remove original target file if not preserving original
		if !preserveOriginal {
			if err := os.Remove(target); err != nil {
				return fmt.Errorf("failed to remove original symlink target: %w", err)
			}
		}
	}

	return nil
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
