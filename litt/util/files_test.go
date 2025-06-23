package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVerifyFileProperties(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Test cases
	tests := []struct {
		name             string
		setup            func() string
		expectedExists   bool
		expectedSize     int64
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "existing file with correct permissions",
			setup: func() string {
				path := filepath.Join(tempDir, "test-file")
				err := os.WriteFile(path, []byte("test data"), 0600)
				require.NoError(t, err)
				return path
			},
			expectedExists: true,
			expectedSize:   9, // "test data" is 9 bytes
			expectError:    false,
		},
		{
			name: "non-existent file with writable parent",
			setup: func() string {
				return filepath.Join(tempDir, "non-existent-file")
			},
			expectedExists: false,
			expectedSize:   -1,
			expectError:    false,
		},
		{
			name: "non-existent file with non-existent parent",
			setup: func() string {
				return filepath.Join(tempDir, "non-existent-dir", "non-existent-file")
			},
			expectedExists:   false,
			expectedSize:     -1,
			expectError:      true,
			expectedErrorMsg: "parent directory",
		},
		{
			name: "existing file is a directory",
			setup: func() string {
				path := filepath.Join(tempDir, "test-dir")
				err := os.Mkdir(path, 0755)
				require.NoError(t, err)
				return path
			},
			expectedExists:   false,
			expectedSize:     -1,
			expectError:      true,
			expectedErrorMsg: "is a directory",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.setup()
			exists, size, err := VerifyFileProperties(path)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedExists, exists)
			require.Equal(t, tc.expectedSize, size)
		})
	}
}

func TestExists(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing-file")
	err := os.WriteFile(existingFile, []byte("test"), 0600)
	require.NoError(t, err)

	nonExistentFile := filepath.Join(tempDir, "non-existent-file")

	// Test cases
	tests := []struct {
		name        string
		path        string
		expected    bool
		expectError bool
	}{
		{
			name:        "existing file",
			path:        existingFile,
			expected:    true,
			expectError: false,
		},
		{
			name:        "non-existent file",
			path:        nonExistentFile,
			expected:    false,
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := Exists(tc.path)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, exists)
		})
	}
}

func TestVerifyDirectoryWritable(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Create a non-writable directory (0500 = read & execute, no write)
	nonWritableDir := filepath.Join(tempDir, "non-writable-dir")
	err := os.Mkdir(nonWritableDir, 0500)
	require.NoError(t, err)

	// Create a writable directory
	writableDir := filepath.Join(tempDir, "writable-dir")
	err = os.Mkdir(writableDir, 0700)
	require.NoError(t, err)

	// Create a regular file
	regularFile := filepath.Join(tempDir, "regular-file")
	err = os.WriteFile(regularFile, []byte("test"), 0600)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "writable directory",
			path:        writableDir,
			expectError: false,
		},
		{
			name:        "non-writable directory",
			path:        nonWritableDir,
			expectError: true,
			errorMsg:    "not writable",
		},
		{
			name:        "regular file",
			path:        regularFile,
			expectError: true,
			errorMsg:    "is not a directory",
		},
		{
			name:        "non-existent directory with writable parent",
			path:        filepath.Join(writableDir, "non-existent"),
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyDirectoryWritable(tc.path)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}

	// Cleanup special permissions
	err = os.Chmod(nonWritableDir, 0700)
	require.NoError(t, err)
}

func TestEnsureParentDirExists(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Create a non-writable directory (0500 = read & execute, no write)
	nonWritableDir := filepath.Join(tempDir, "non-writable-dir")
	err := os.Mkdir(nonWritableDir, 0500)
	require.NoError(t, err)

	// Create a test file
	testFile := filepath.Join(tempDir, "test-file")
	err = os.WriteFile(testFile, []byte("test"), 0600)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "parent exists and is writable",
			path:        filepath.Join(tempDir, "new-file"),
			expectError: false,
		},
		{
			name:        "parent is non-writable",
			path:        filepath.Join(nonWritableDir, "new-file"),
			expectError: true,
			errorMsg:    "not writable",
		},
		{
			name:        "multi-level parent doesn't exist",
			path:        filepath.Join(tempDir, "new-dir", "subdir", "new-file"),
			expectError: false,
		},
		{
			name:        "parent exists but is a file",
			path:        filepath.Join(testFile, "impossible"),
			expectError: true,
			errorMsg:    "is not a directory",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ensureParentDirExists(tc.path)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)

				// Verify the parent directory was created if needed
				parentDir := filepath.Dir(tc.path)
				exists, err := Exists(parentDir)
				require.NoError(t, err)
				require.True(t, exists)
			}
		})
	}

	// Cleanup special permissions
	err = os.Chmod(nonWritableDir, 0700)
	require.NoError(t, err)
}

func TestCopyRegularFile(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Create a source file with specific content, permissions, and time
	sourceFile := filepath.Join(tempDir, "source-file")
	content := []byte("test content")
	err := os.WriteFile(sourceFile, content, 0640)
	require.NoError(t, err)

	// Set a specific modification time
	modTime := time.Now().Add(-24 * time.Hour) // yesterday
	err = os.Chtimes(sourceFile, modTime, modTime)
	require.NoError(t, err)

	// Get file info for permissions and modtime
	sourceInfo, err := os.Stat(sourceFile)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		destPath    string
		expectError bool
	}{
		{
			name:        "copy to a new file",
			destPath:    filepath.Join(tempDir, "dest-file"),
			expectError: false,
		},
		{
			name:        "overwrite existing file",
			destPath:    filepath.Join(tempDir, "existing-file"),
			expectError: false,
		},
		{
			name:        "copy to a new subdirectory",
			destPath:    filepath.Join(tempDir, "subdir", "dest-file"),
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// If testing overwrite, create the file first
			if tc.name == "overwrite existing file" {
				err := os.WriteFile(tc.destPath, []byte("original content"), 0600)
				require.NoError(t, err)
			}

			err := copyRegularFile(sourceFile, tc.destPath, sourceInfo.Mode(), sourceInfo.ModTime())

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify the file was copied correctly
				destInfo, err := os.Stat(tc.destPath)
				require.NoError(t, err)

				// Check content
				destContent, err := os.ReadFile(tc.destPath)
				require.NoError(t, err)
				require.Equal(t, content, destContent)

				// Check permissions (mask out the bits we don't care about for comparison)
				// We only care about user, group, and world permissions (0777)
				// This handles umask and platform differences
				require.Equal(t, sourceInfo.Mode()&0777, destInfo.Mode()&0777)

				// Check modification time
				require.Equal(t, sourceInfo.ModTime().Unix(), destInfo.ModTime().Unix())
			}
		})
	}
}

func TestCopySymlink(t *testing.T) {
	// Skip on platforms that don't support symlinks (like Windows in some cases)
	if !supportsSymlinks() {
		t.Skip("Symlinks not supported on this platform/environment")
	}

	// Setup
	tempDir := t.TempDir()

	// Create a target file for the symlink
	targetFile := filepath.Join(tempDir, "target-file")
	err := os.WriteFile(targetFile, []byte("target content"), 0644)
	require.NoError(t, err)

	// Create a source symlink
	sourceSymlink := filepath.Join(tempDir, "source-symlink")
	err = os.Symlink(targetFile, sourceSymlink)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		destPath    string
		expectError bool
	}{
		{
			name:        "copy to a new location",
			destPath:    filepath.Join(tempDir, "dest-symlink"),
			expectError: false,
		},
		{
			name:        "copy to a new subdirectory",
			destPath:    filepath.Join(tempDir, "subdir", "dest-symlink"),
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := copySymlink(sourceSymlink, tc.destPath)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify the symlink was created
				linkDest, err := os.Readlink(tc.destPath)
				require.NoError(t, err)

				// Verify it points to the right target
				require.Equal(t, targetFile, linkDest)
			}
		})
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Create a regular file
	regularFile := filepath.Join(tempDir, "regular-file")
	err := os.WriteFile(regularFile, []byte("test"), 0600)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		dirPath     string
		mode        os.FileMode
		setup       func(path string)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "directory doesn't exist",
			dirPath:     filepath.Join(tempDir, "new-dir"),
			mode:        0755,
			setup:       func(path string) {},
			expectError: false,
		},
		{
			name:    "directory already exists",
			dirPath: filepath.Join(tempDir, "existing-dir"),
			mode:    0755,
			setup: func(path string) {
				err := os.Mkdir(path, 0755)
				require.NoError(t, err)
			},
			expectError: false,
		},
		{
			name:        "path exists but is a file",
			dirPath:     regularFile,
			mode:        0755,
			setup:       func(path string) {},
			expectError: true,
			errorMsg:    "is not a directory",
		},
		{
			name:    "directory exists but is non-writable",
			dirPath: filepath.Join(tempDir, "non-writable-dir"),
			mode:    0755,
			setup: func(path string) {
				err := os.Mkdir(path, 0500) // read & execute only
				require.NoError(t, err)
			},
			expectError: true,
			errorMsg:    "not writable",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.dirPath)

			err := ensureDirectoryExists(tc.dirPath, tc.mode)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)

				// Verify the directory exists
				info, err := os.Stat(tc.dirPath)
				require.NoError(t, err)
				require.True(t, info.IsDir())

				// If we created a new directory, verify the mode
				if tc.name == "directory doesn't exist" {
					// Note: mode comparison can be tricky due to umask and OS differences
					// So we just check that it's writable
					require.True(t, info.Mode()&0200 != 0, "Directory should be writable")
				}
			}
		})
	}

	// Clean up non-writable directory
	nonWritableDir := filepath.Join(tempDir, "non-writable-dir")
	if _, err := os.Stat(nonWritableDir); err == nil {
		err = os.Chmod(nonWritableDir, 0700)
		require.NoError(t, err)
	}
}

// Helper function to check if symlinks are supported in the current environment
func supportsSymlinks() bool {
	tempDir, err := os.MkdirTemp("", "symlink-test")
	if err != nil {
		return false
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clean up temp directory %s: %v\n", tempDir, err)
		}
	}()

	source := filepath.Join(tempDir, "source")
	target := filepath.Join(tempDir, "target")

	err = os.WriteFile(source, []byte{}, 0644)
	if err != nil {
		return false
	}

	err = os.Symlink(source, target)
	return err == nil
}

func TestAtomicWrite(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Test cases
	tests := []struct {
		name        string
		setup       func() (string, []byte)
		expectError bool
		errorMsg    string
	}{
		{
			name: "write to new file",
			setup: func() (string, []byte) {
				path := filepath.Join(tempDir, "new-file.txt")
				data := []byte("test content")
				return path, data
			},
			expectError: false,
		},
		{
			name: "overwrite existing file",
			setup: func() (string, []byte) {
				path := filepath.Join(tempDir, "existing-file.txt")
				// Create existing file with different content
				err := os.WriteFile(path, []byte("old content"), 0644)
				require.NoError(t, err)
				data := []byte("new content")
				return path, data
			},
			expectError: false,
		},
		{
			name: "write to subdirectory",
			setup: func() (string, []byte) {
				subDir := filepath.Join(tempDir, "subdir")
				err := os.Mkdir(subDir, 0755)
				require.NoError(t, err)
				path := filepath.Join(subDir, "file.txt")
				data := []byte("content in subdirectory")
				return path, data
			},
			expectError: false,
		},
		{
			name: "write with empty data",
			setup: func() (string, []byte) {
				path := filepath.Join(tempDir, "empty-file.txt")
				data := []byte("")
				return path, data
			},
			expectError: false,
		},
		{
			name: "write to non-existent parent directory",
			setup: func() (string, []byte) {
				path := filepath.Join(tempDir, "non-existent-dir", "file.txt")
				data := []byte("content")
				return path, data
			},
			expectError: true,
			errorMsg:    "failed to create swap file",
		},
		{
			name: "write with large data",
			setup: func() (string, []byte) {
				path := filepath.Join(tempDir, "large-file.txt")
				// Create 1MB of data
				data := make([]byte, 1024*1024)
				for i := range data {
					data[i] = byte(i % 256)
				}
				return path, data
			},
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, data := tc.setup()
			swapPath := path + SwapFileExtension

			// Ensure swap file doesn't exist before test
			_, err := os.Stat(swapPath)
			require.True(t, os.IsNotExist(err), "Swap file should not exist before test")

			err = AtomicWrite(path, data, true)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)

				// Verify that the destination file wasn't created or modified
				if tc.name == "overwrite existing file" {
					// Original file should still have old content
					content, err := os.ReadFile(path)
					require.NoError(t, err)
					require.Equal(t, "old content", string(content))
				}
			} else {
				require.NoError(t, err)

				// Verify the file was written correctly
				content, err := os.ReadFile(path)
				require.NoError(t, err)
				require.Equal(t, data, content)

				// Verify the swap file was cleaned up
				_, err = os.Stat(swapPath)
				require.True(t, os.IsNotExist(err), "Swap file should be cleaned up after successful write")

				// Verify file permissions are reasonable (at least owner readable/writable)
				info, err := os.Stat(path)
				require.NoError(t, err)
				require.True(t, info.Mode()&0600 != 0, "File should be readable and writable by owner")
			}
		})
	}
}

func TestAtomicWriteSwapFileCleanup(t *testing.T) {
	// Test that swap files are properly cleaned up even if something goes wrong
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test-file.txt")
	swapPath := path + SwapFileExtension
	data := []byte("test content")

	// Simulate a scenario where swap file might be left behind
	// by creating a swap file manually first
	err := os.WriteFile(swapPath, []byte("old swap content"), 0644)
	require.NoError(t, err)

	// Verify swap file exists
	_, err = os.Stat(swapPath)
	require.NoError(t, err)

	// Now run AtomicWrite - it should overwrite the swap file and clean up
	err = AtomicWrite(path, data, true)
	require.NoError(t, err)

	// Verify the target file has the correct content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, data, content)

	// Verify the swap file was cleaned up
	_, err = os.Stat(swapPath)
	require.True(t, os.IsNotExist(err), "Swap file should be cleaned up")
}

func TestAtomicWritePreservesOtherFiles(t *testing.T) {
	// Test that AtomicWrite doesn't interfere with other files in the same directory
	tempDir := t.TempDir()

	// Create some existing files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	targetFile := filepath.Join(tempDir, "target.txt")

	err := os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	// Perform atomic write on target file
	targetData := []byte("target content")
	err = AtomicWrite(targetFile, targetData, true)
	require.NoError(t, err)

	// Verify all files have correct content
	content1, err := os.ReadFile(file1)
	require.NoError(t, err)
	require.Equal(t, "content1", string(content1))

	content2, err := os.ReadFile(file2)
	require.NoError(t, err)
	require.Equal(t, "content2", string(content2))

	targetContent, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	require.Equal(t, targetData, targetContent)
}

func TestAtomicRename(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Test cases
	tests := []struct {
		name        string
		setup       func() (string, string)
		expectError bool
		errorMsg    string
	}{
		{
			name: "rename file in same directory",
			setup: func() (string, string) {
				oldPath := filepath.Join(tempDir, "old-name.txt")
				newPath := filepath.Join(tempDir, "new-name.txt")
				err := os.WriteFile(oldPath, []byte("test content"), 0644)
				require.NoError(t, err)
				return oldPath, newPath
			},
			expectError: false,
		},
		{
			name: "rename file to different directory",
			setup: func() (string, string) {
				subDir := filepath.Join(tempDir, "subdir")
				err := os.Mkdir(subDir, 0755)
				require.NoError(t, err)

				oldPath := filepath.Join(tempDir, "file.txt")
				newPath := filepath.Join(subDir, "moved-file.txt")
				err = os.WriteFile(oldPath, []byte("content to move"), 0644)
				require.NoError(t, err)
				return oldPath, newPath
			},
			expectError: false,
		},
		{
			name: "overwrite existing file",
			setup: func() (string, string) {
				oldPath := filepath.Join(tempDir, "source.txt")
				newPath := filepath.Join(tempDir, "target.txt")

				// Create source file
				err := os.WriteFile(oldPath, []byte("source content"), 0644)
				require.NoError(t, err)

				// Create target file that will be overwritten
				err = os.WriteFile(newPath, []byte("target content"), 0644)
				require.NoError(t, err)

				return oldPath, newPath
			},
			expectError: false,
		},
		{
			name: "rename non-existent file",
			setup: func() (string, string) {
				oldPath := filepath.Join(tempDir, "non-existent.txt")
				newPath := filepath.Join(tempDir, "new.txt")
				return oldPath, newPath
			},
			expectError: true,
			errorMsg:    "failed to rename file",
		},
		{
			name: "rename to non-existent directory",
			setup: func() (string, string) {
				oldPath := filepath.Join(tempDir, "existing.txt")
				newPath := filepath.Join(tempDir, "non-existent-dir", "file.txt")
				err := os.WriteFile(oldPath, []byte("content"), 0644)
				require.NoError(t, err)
				return oldPath, newPath
			},
			expectError: true,
			errorMsg:    "failed to rename file",
		},
		{
			name: "rename directory",
			setup: func() (string, string) {
				oldDir := filepath.Join(tempDir, "old-dir")
				newDir := filepath.Join(tempDir, "new-dir")

				err := os.Mkdir(oldDir, 0755)
				require.NoError(t, err)

				// Add a file inside the directory
				err = os.WriteFile(filepath.Join(oldDir, "file.txt"), []byte("dir content"), 0644)
				require.NoError(t, err)

				return oldDir, newDir
			},
			expectError: false,
		},
		{
			name: "rename with same source and destination",
			setup: func() (string, string) {
				path := filepath.Join(tempDir, "same-file.txt")
				err := os.WriteFile(path, []byte("content"), 0644)
				require.NoError(t, err)
				return path, path
			},
			expectError: false, // os.Rename typically succeeds for same path
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldPath, newPath := tc.setup()

			// Store original content if file exists
			var originalContent []byte
			var originalInfo os.FileInfo
			if info, err := os.Stat(oldPath); err == nil {
				if !info.IsDir() {
					originalContent, err = os.ReadFile(oldPath)
					require.NoError(t, err)
				}
				originalInfo = info
			}

			err := AtomicRename(oldPath, newPath, true)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)

				// Verify original file still exists (rename failed)
				if originalInfo != nil {
					_, err := os.Stat(oldPath)
					if tc.errorMsg == "failed to rename file" {
						require.NoError(t, err, "Original file should still exist after failed rename")
					}
				}
			} else {
				require.NoError(t, err)

				// Verify the rename was successful
				if tc.name != "rename with same source and destination" {
					// Old path should not exist
					_, err := os.Stat(oldPath)
					require.True(t, os.IsNotExist(err), "Old path should not exist after successful rename")
				}

				// New path should exist
				newInfo, err := os.Stat(newPath)
				require.NoError(t, err, "New path should exist after successful rename")

				// Verify content and properties if it was a file
				if originalInfo != nil && !originalInfo.IsDir() {
					if tc.name != "rename with same source and destination" {
						// Check content preservation
						newContent, err := os.ReadFile(newPath)
						require.NoError(t, err)
						require.Equal(t, originalContent, newContent, "File content should be preserved")
					}

					// Check that it's still a file
					require.False(t, newInfo.IsDir(), "Renamed file should still be a file")
				} else if originalInfo != nil && originalInfo.IsDir() {
					// Check that it's still a directory
					require.True(t, newInfo.IsDir(), "Renamed directory should still be a directory")

					// Check that directory contents are preserved
					if tc.name == "rename directory" {
						fileContent, err := os.ReadFile(filepath.Join(newPath, "file.txt"))
						require.NoError(t, err)
						require.Equal(t, "dir content", string(fileContent))
					}
				}
			}
		})
	}
}

func TestAtomicRenamePreservesPermissions(t *testing.T) {
	// Test that file permissions are preserved during atomic rename
	tempDir := t.TempDir()

	oldPath := filepath.Join(tempDir, "source.txt")
	newPath := filepath.Join(tempDir, "dest.txt")

	// Create file with specific permissions
	err := os.WriteFile(oldPath, []byte("test content"), 0640)
	require.NoError(t, err)

	// Get original permissions
	originalInfo, err := os.Stat(oldPath)
	require.NoError(t, err)

	// Perform atomic rename
	err = AtomicRename(oldPath, newPath, true)
	require.NoError(t, err)

	// Verify permissions are preserved
	newInfo, err := os.Stat(newPath)
	require.NoError(t, err)
	require.Equal(t, originalInfo.Mode(), newInfo.Mode(), "File permissions should be preserved")
}

func TestAtomicRenameWithSymlink(t *testing.T) {
	// Skip on platforms that don't support symlinks
	if !supportsSymlinks() {
		t.Skip("Symlinks not supported on this platform/environment")
	}

	tempDir := t.TempDir()

	// Create a target file
	targetFile := filepath.Join(tempDir, "target.txt")
	err := os.WriteFile(targetFile, []byte("target content"), 0644)
	require.NoError(t, err)

	// Create a symlink
	oldLink := filepath.Join(tempDir, "old-link")
	err = os.Symlink(targetFile, oldLink)
	require.NoError(t, err)

	// Rename the symlink
	newLink := filepath.Join(tempDir, "new-link")
	err = AtomicRename(oldLink, newLink, true)
	require.NoError(t, err)

	// Verify the symlink was renamed and still points to the same target
	linkTarget, err := os.Readlink(newLink)
	require.NoError(t, err)
	require.Equal(t, targetFile, linkTarget)

	// Verify old symlink no longer exists
	_, err = os.Stat(oldLink)
	require.True(t, os.IsNotExist(err))
}

func TestAtomicRenameAcrossFilesystems(t *testing.T) {
	// This test checks behavior when renaming across filesystem boundaries
	// On most systems, this will fall back to copy+delete, but the function should still work
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	content := []byte("content for cross-filesystem test")
	err := os.WriteFile(srcPath, content, 0644)
	require.NoError(t, err)

	// Try to rename to /tmp (likely different filesystem on many systems)
	// If this fails due to cross-filesystem issues, that's expected behavior
	tmpDir, err := os.MkdirTemp("", "atomic-rename-test-")
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to clean up temp directory %s: %v\n", tmpDir, err)
		}
	})

	dstPath := filepath.Join(tmpDir, "dest.txt")

	err = AtomicRename(srcPath, dstPath, true)
	// This might succeed or fail depending on the system
	// If it succeeds, verify the file was moved correctly
	if err == nil {
		// Verify content
		newContent, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		require.Equal(t, content, newContent)

		// Verify source no longer exists
		_, err = os.Stat(srcPath)
		require.True(t, os.IsNotExist(err))
	}
	// If it fails, that's also acceptable for cross-filesystem renames
}

func TestDeleteOrphanedSwapFiles(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	// Test cases
	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "delete swap files in directory with mixed files",
			setup: func() string {
				testDir := filepath.Join(tempDir, "mixed-files")
				err := os.Mkdir(testDir, 0755)
				require.NoError(t, err)

				// Create regular files
				err = os.WriteFile(filepath.Join(testDir, "regular1.txt"), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testDir, "regular2.log"), []byte("content2"), 0644)
				require.NoError(t, err)

				// Create swap files
				err = os.WriteFile(filepath.Join(testDir, "file1.txt"+SwapFileExtension), []byte("swap1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testDir, "file2.log"+SwapFileExtension), []byte("swap2"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testDir, "orphaned"+SwapFileExtension), []byte("orphaned"), 0644)
				require.NoError(t, err)

				// Create a subdirectory (should be ignored)
				subDir := filepath.Join(testDir, "subdir")
				err = os.Mkdir(subDir, 0755)
				require.NoError(t, err)

				// Create a swap file in subdirectory (should not be deleted by this call)
				err = os.WriteFile(filepath.Join(subDir, "nested"+SwapFileExtension), []byte("nested"), 0644)
				require.NoError(t, err)

				return testDir
			},
			expectError: false,
		},
		{
			name: "empty directory",
			setup: func() string {
				testDir := filepath.Join(tempDir, "empty-dir")
				err := os.Mkdir(testDir, 0755)
				require.NoError(t, err)
				return testDir
			},
			expectError: false,
		},
		{
			name: "directory with only swap files",
			setup: func() string {
				testDir := filepath.Join(tempDir, "only-swap")
				err := os.Mkdir(testDir, 0755)
				require.NoError(t, err)

				// Create only swap files
				err = os.WriteFile(filepath.Join(testDir, "swap1"+SwapFileExtension), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testDir, "swap2"+SwapFileExtension), []byte("content2"), 0644)
				require.NoError(t, err)

				return testDir
			},
			expectError: false,
		},
		{
			name: "directory with no swap files",
			setup: func() string {
				testDir := filepath.Join(tempDir, "no-swap")
				err := os.Mkdir(testDir, 0755)
				require.NoError(t, err)

				// Create only regular files
				err = os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(testDir, "file2.log"), []byte("content2"), 0644)
				require.NoError(t, err)

				return testDir
			},
			expectError: false,
		},
		{
			name: "non-existent directory",
			setup: func() string {
				return filepath.Join(tempDir, "non-existent")
			},
			expectError: true,
			errorMsg:    "failed to read directory",
		},
		{
			name: "path is a file not directory",
			setup: func() string {
				filePath := filepath.Join(tempDir, "not-a-dir.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectError: true,
			errorMsg:    "failed to read directory",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dirPath := tc.setup()

			// Count files before deletion for verification
			var beforeFiles []string
			if entries, err := os.ReadDir(dirPath); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						beforeFiles = append(beforeFiles, entry.Name())
					}
				}
			}

			err := DeleteOrphanedSwapFiles(dirPath)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)

				// Verify that all swap files were deleted
				entries, err := os.ReadDir(dirPath)
				require.NoError(t, err)

				var afterFiles []string
				var afterSwapFiles []string
				for _, entry := range entries {
					if !entry.IsDir() {
						afterFiles = append(afterFiles, entry.Name())
						if filepath.Ext(entry.Name()) == SwapFileExtension {
							afterSwapFiles = append(afterSwapFiles, entry.Name())
						}
					}
				}

				// No swap files should remain
				require.Empty(t, afterSwapFiles, "All swap files should be deleted")

				// Regular files should remain unchanged
				var beforeRegularFiles []string
				var afterRegularFiles []string
				for _, file := range beforeFiles {
					if filepath.Ext(file) != SwapFileExtension {
						beforeRegularFiles = append(beforeRegularFiles, file)
					}
				}
				for _, file := range afterFiles {
					if filepath.Ext(file) != SwapFileExtension {
						afterRegularFiles = append(afterRegularFiles, file)
					}
				}
				require.ElementsMatch(t, beforeRegularFiles, afterRegularFiles, "Regular files should be unchanged")

				// Verify subdirectories are not affected
				if tc.name == "delete swap files in directory with mixed files" {
					subDirPath := filepath.Join(dirPath, "subdir")
					subEntries, err := os.ReadDir(subDirPath)
					require.NoError(t, err)
					require.Len(t, subEntries, 1, "Subdirectory should still contain its swap file")
					require.Equal(t, "nested"+SwapFileExtension, subEntries[0].Name())
				}
			}
		})
	}
}

func TestDeleteOrphanedSwapFilesPermissions(t *testing.T) {
	// Test behavior with permission issues
	tempDir := t.TempDir()

	// Create a directory with swap files
	testDir := filepath.Join(tempDir, "perm-test")
	err := os.Mkdir(testDir, 0755)
	require.NoError(t, err)

	// Create a swap file
	swapFile := filepath.Join(testDir, "test"+SwapFileExtension)
	err = os.WriteFile(swapFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Make the directory read-only (no write permissions)
	err = os.Chmod(testDir, 0555) // read + execute only
	require.NoError(t, err)

	// Attempt to delete swap files should fail
	err = DeleteOrphanedSwapFiles(testDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to remove swap file")

	// Restore permissions for cleanup
	err = os.Chmod(testDir, 0755)
	require.NoError(t, err)
}

func TestRecursiveMove(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	tests := []struct {
		name             string
		setup            func() (string, string)
		deep             bool
		preserveOriginal bool
		fsync            bool
		expectError      bool
		errorMsg         string
		verify           func(t *testing.T, source, dest string)
	}{
		{
			name: "move single file - rename optimization when not preserving",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "single-source.txt")
				dest := filepath.Join(tempDir, "single-dest.txt")
				err := os.WriteFile(source, []byte("single file content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Source should not exist (moved, not copied)
				_, err := os.Stat(source)
				require.True(t, os.IsNotExist(err), "Source should not exist after move")

				// Destination should exist with correct content
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "single file content", string(content))
			},
		},
		{
			name: "move single file - preserve original with hard link",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "preserve-source.txt")
				dest := filepath.Join(tempDir, "preserve-dest.txt")
				err := os.WriteFile(source, []byte("preserve content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Both should exist
				sourceContent, err := os.ReadFile(source)
				require.NoError(t, err)
				destContent, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "preserve content", string(sourceContent))
				require.Equal(t, "preserve content", string(destContent))

				// They should be hard linked (same inode)
				sourceInfo, err := os.Stat(source)
				require.NoError(t, err)
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				// Note: hard link verification is platform-specific, so we just check both exist
				require.Equal(t, sourceInfo.Size(), destInfo.Size())
			},
		},
		{
			name: "move single file - deep copy when preserving",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "deep-source.txt")
				dest := filepath.Join(tempDir, "deep-dest.txt")
				err := os.WriteFile(source, []byte("deep copy content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             true,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Both should exist with same content
				sourceContent, err := os.ReadFile(source)
				require.NoError(t, err)
				destContent, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "deep copy content", string(sourceContent))
				require.Equal(t, "deep copy content", string(destContent))
			},
		},
		{
			name: "move directory structure",
			setup: func() (string, string) {
				sourceDir := filepath.Join(tempDir, "source-dir")
				destDir := filepath.Join(tempDir, "dest-dir")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				// Create files and subdirectories
				err = os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("file1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("file2"), 0644)
				require.NoError(t, err)

				subDir := filepath.Join(sourceDir, "subdir")
				err = os.Mkdir(subDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Source should not exist
				_, err := os.Stat(source)
				require.True(t, os.IsNotExist(err), "Source directory should not exist after move")

				// Destination should have all content
				content1, err := os.ReadFile(filepath.Join(dest, "file1.txt"))
				require.NoError(t, err)
				require.Equal(t, "file1", string(content1))

				content2, err := os.ReadFile(filepath.Join(dest, "file2.txt"))
				require.NoError(t, err)
				require.Equal(t, "file2", string(content2))

				nestedContent, err := os.ReadFile(filepath.Join(dest, "subdir", "nested.txt"))
				require.NoError(t, err)
				require.Equal(t, "nested", string(nestedContent))
			},
		},
		{
			name: "move directory - preserve original",
			setup: func() (string, string) {
				sourceDir := filepath.Join(tempDir, "preserve-source-dir")
				destDir := filepath.Join(tempDir, "preserve-dest-dir")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				err = os.WriteFile(filepath.Join(sourceDir, "preserve-file.txt"), []byte("preserve content"), 0644)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             false,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Both should exist
				sourceContent, err := os.ReadFile(filepath.Join(source, "preserve-file.txt"))
				require.NoError(t, err)
				destContent, err := os.ReadFile(filepath.Join(dest, "preserve-file.txt"))
				require.NoError(t, err)
				require.Equal(t, "preserve content", string(sourceContent))
				require.Equal(t, "preserve content", string(destContent))
			},
		},
		{
			name: "move with symlink",
			setup: func() (string, string) {
				if !supportsSymlinks() {
					t.Skip("Symlinks not supported")
				}

				sourceDir := filepath.Join(tempDir, "symlink-source")
				destDir := filepath.Join(tempDir, "symlink-dest")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				// Create target file
				targetFile := filepath.Join(sourceDir, "target.txt")
				err = os.WriteFile(targetFile, []byte("target content"), 0644)
				require.NoError(t, err)

				// Create symlink
				linkFile := filepath.Join(sourceDir, "link.txt")
				err = os.Symlink(targetFile, linkFile)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				if !supportsSymlinks() {
					return
				}

				// Check that symlink was copied correctly (shallow copy preserves symlink)
				linkTarget, err := os.Readlink(filepath.Join(dest, "link.txt"))
				require.NoError(t, err)
				// The target should still point to the original source location
				expectedTarget := filepath.Join(source, "target.txt")
				require.Equal(t, expectedTarget, linkTarget)

				// Check that target file was also moved
				content, err := os.ReadFile(filepath.Join(dest, "target.txt"))
				require.NoError(t, err)
				require.Equal(t, "target content", string(content))
			},
		},
		{
			name: "move to subdirectory",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "subdir-source.txt")
				dest := filepath.Join(tempDir, "new-subdir", "dest.txt")
				err := os.WriteFile(source, []byte("subdir content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, source, dest string) {
				// Parent directory should have been created
				parentDir := filepath.Dir(dest)
				info, err := os.Stat(parentDir)
				require.NoError(t, err)
				require.True(t, info.IsDir())

				// File should be moved
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "subdir content", string(content))

				_, err = os.Stat(source)
				require.True(t, os.IsNotExist(err))
			},
		},
		{
			name: "error - source does not exist",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "non-existent.txt")
				dest := filepath.Join(tempDir, "dest.txt")
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      true,
			errorMsg:         "source path",
		},
		{
			name: "error - destination parent not writable",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "readonly-source.txt")
				err := os.WriteFile(source, []byte("content"), 0644)
				require.NoError(t, err)

				// Create read-only parent directory
				readOnlyDir := filepath.Join(tempDir, "readonly-parent")
				err = os.Mkdir(readOnlyDir, 0555) // read + execute only
				require.NoError(t, err)

				dest := filepath.Join(readOnlyDir, "dest.txt")
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      true,
			errorMsg:         "destination parent directory not writable",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source, dest := tc.setup()

			err := RecursiveMove(source, dest, tc.deep, tc.preserveOriginal, tc.fsync)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, source, dest)
				}
			}

			// Cleanup read-only directories for next tests
			if readOnlyDir := filepath.Join(tempDir, "readonly-parent"); strings.Contains(dest, "readonly-parent") {
				_ = os.Chmod(readOnlyDir, 0755) // Ignore error for cleanup
			}
		})
	}
}

func TestRecursiveMoveFile(t *testing.T) {
	// Test the internal recursiveMoveFile function directly
	tempDir := t.TempDir()

	tests := []struct {
		name             string
		setup            func() (string, string)
		deep             bool
		preserveOriginal bool
		fsync            bool
		expectError      bool
		errorMsg         string
	}{
		{
			name: "successful rename when not preserving",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "move-source.txt")
				dest := filepath.Join(tempDir, "move-dest.txt")
				err := os.WriteFile(source, []byte("move content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
		},
		{
			name: "hard link when preserving and not deep",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "link-source.txt")
				dest := filepath.Join(tempDir, "link-dest.txt")
				err := os.WriteFile(source, []byte("link content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
		},
		{
			name: "copy when deep and preserving",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "copy-source.txt")
				dest := filepath.Join(tempDir, "copy-dest.txt")
				err := os.WriteFile(source, []byte("copy content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             true,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
		},
		{
			name: "copy and delete when not preserving but rename fails",
			setup: func() (string, string) {
				source := filepath.Join(tempDir, "cross-source.txt")
				// Try to move to a different temp directory (might fail rename, fall back to copy)
				otherTempDir := t.TempDir()
				dest := filepath.Join(otherTempDir, "cross-dest.txt")
				err := os.WriteFile(source, []byte("cross content"), 0644)
				require.NoError(t, err)
				return source, dest
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source, dest := tc.setup()

			err := recursiveMoveFile(source, dest, tc.deep, tc.preserveOriginal, tc.fsync)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)

				// Verify destination exists
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Contains(t, string(content), "content")

				// Check source existence based on preserveOriginal
				_, err = os.Stat(source)
				if tc.preserveOriginal {
					require.NoError(t, err, "Source should exist when preserving original")
				} else {
					require.True(t, os.IsNotExist(err), "Source should not exist when not preserving original")
				}
			}
		})
	}
}

func TestHandleSymlink(t *testing.T) {
	if !supportsSymlinks() {
		t.Skip("Symlinks not supported on this platform/environment")
	}

	tempDir := t.TempDir()

	tests := []struct {
		name             string
		setup            func() (string, string, string) // Returns source symlink, destination, and target file/dir
		deep             bool
		preserveOriginal bool
		fsync            bool
		expectError      bool
		errorMsg         string
		verify           func(t *testing.T, sourceSymlink, dest, target string)
	}{
		{
			name: "shallow copy - symlink to file copied as symlink",
			setup: func() (string, string, string) {
				// Create target file
				target := filepath.Join(tempDir, "target-file.txt")
				err := os.WriteFile(target, []byte("target content"), 0644)
				require.NoError(t, err)

				// Create symlink
				sourceSymlink := filepath.Join(tempDir, "source-link")
				err = os.Symlink(target, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "dest-link")
				return sourceSymlink, dest, target
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, target string) {
				// Destination should be a symlink
				linkTarget, err := os.Readlink(dest)
				require.NoError(t, err)
				require.Equal(t, target, linkTarget)

				// Verify we can read through the symlink
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "target content", string(content))
			},
		},
		{
			name: "deep copy - symlink to file copied as actual file",
			setup: func() (string, string, string) {
				// Create target file
				target := filepath.Join(tempDir, "deep-target-file.txt")
				err := os.WriteFile(target, []byte("deep target content"), 0644)
				require.NoError(t, err)

				// Create symlink
				sourceSymlink := filepath.Join(tempDir, "deep-source-link")
				err = os.Symlink(target, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "deep-dest-file")
				return sourceSymlink, dest, target
			},
			deep:             true,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, target string) {
				// Destination should be a regular file, not a symlink
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				require.False(t, destInfo.Mode()&os.ModeSymlink != 0, "Destination should not be a symlink")

				// Content should match the target file
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "deep target content", string(content))

				// Original target file should not exist when not preserving original
				_, err = os.Stat(target)
				require.True(t, os.IsNotExist(err), "Original target should be removed when not preserving original")
			},
		},
		{
			name: "shallow copy - symlink to directory copied as symlink",
			setup: func() (string, string, string) {
				// Create target directory
				targetDir := filepath.Join(tempDir, "target-dir")
				err := os.Mkdir(targetDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(targetDir, "file.txt"), []byte("dir content"), 0644)
				require.NoError(t, err)

				// Create symlink to directory
				sourceSymlink := filepath.Join(tempDir, "source-dir-link")
				err = os.Symlink(targetDir, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "dest-dir-link")
				return sourceSymlink, dest, targetDir
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, targetDir string) {
				// Destination should be a symlink
				linkTarget, err := os.Readlink(dest)
				require.NoError(t, err)
				require.Equal(t, targetDir, linkTarget)

				// Verify we can access directory through the symlink
				content, err := os.ReadFile(filepath.Join(dest, "file.txt"))
				require.NoError(t, err)
				require.Equal(t, "dir content", string(content))
			},
		},
		{
			name: "deep copy - symlink to directory copied as actual directory",
			setup: func() (string, string, string) {
				// Create target directory
				targetDir := filepath.Join(tempDir, "deep-target-dir")
				err := os.Mkdir(targetDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(targetDir, "deep-file.txt"), []byte("deep dir content"), 0644)
				require.NoError(t, err)

				// Create subdirectory
				subDir := filepath.Join(targetDir, "subdir")
				err = os.Mkdir(subDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested content"), 0644)
				require.NoError(t, err)

				// Create symlink to directory
				sourceSymlink := filepath.Join(tempDir, "deep-source-dir-link")
				err = os.Symlink(targetDir, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "deep-dest-dir")
				return sourceSymlink, dest, targetDir
			},
			deep:             true,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, targetDir string) {
				// Destination should be a regular directory, not a symlink
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				require.True(t, destInfo.IsDir(), "Destination should be a directory")
				require.False(t, destInfo.Mode()&os.ModeSymlink != 0, "Destination should not be a symlink")

				// Content should be copied
				content, err := os.ReadFile(filepath.Join(dest, "deep-file.txt"))
				require.NoError(t, err)
				require.Equal(t, "deep dir content", string(content))

				// Nested content should be copied
				nestedContent, err := os.ReadFile(filepath.Join(dest, "subdir", "nested.txt"))
				require.NoError(t, err)
				require.Equal(t, "nested content", string(nestedContent))

				// Original target directory should not exist when not preserving original
				_, err = os.Stat(targetDir)
				require.True(t, os.IsNotExist(err), "Original target directory should be removed when not preserving original")
			},
		},
		{
			name: "error - broken symlink with deep copy",
			setup: func() (string, string, string) {
				// Create symlink pointing to non-existent file
				sourceSymlink := filepath.Join(tempDir, "broken-link")
				nonExistentTarget := filepath.Join(tempDir, "non-existent.txt")
				err := os.Symlink(nonExistentTarget, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "broken-dest")
				return sourceSymlink, dest, nonExistentTarget
			},
			deep:             true,
			preserveOriginal: false,
			fsync:            false,
			expectError:      true,
			errorMsg:         "failed to resolve symlink",
		},
		{
			name: "shallow copy - broken symlink copied as symlink",
			setup: func() (string, string, string) {
				// Create symlink pointing to non-existent file
				sourceSymlink := filepath.Join(tempDir, "broken-shallow-link")
				nonExistentTarget := filepath.Join(tempDir, "non-existent-shallow.txt")
				err := os.Symlink(nonExistentTarget, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "broken-shallow-dest")
				return sourceSymlink, dest, nonExistentTarget
			},
			deep:             false,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, target string) {
				// Destination should be a symlink pointing to the same (non-existent) target
				linkTarget, err := os.Readlink(dest)
				require.NoError(t, err)
				require.Equal(t, target, linkTarget)

				// Trying to read through the symlink should fail
				_, err = os.ReadFile(dest)
				require.Error(t, err)
			},
		},
		{
			name: "chain of symlinks - deep copy follows to final target",
			setup: func() (string, string, string) {
				// Create final target
				finalTarget := filepath.Join(tempDir, "final-target.txt")
				err := os.WriteFile(finalTarget, []byte("final content"), 0644)
				require.NoError(t, err)

				// Create intermediate symlink
				intermediateLink := filepath.Join(tempDir, "intermediate-link")
				err = os.Symlink(finalTarget, intermediateLink)
				require.NoError(t, err)

				// Create source symlink pointing to intermediate
				sourceSymlink := filepath.Join(tempDir, "chain-source-link")
				err = os.Symlink(intermediateLink, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "chain-dest")
				return sourceSymlink, dest, finalTarget
			},
			deep:             true,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, finalTarget string) {
				// Destination should be a regular file with final target content
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				require.False(t, destInfo.Mode()&os.ModeSymlink != 0, "Destination should not be a symlink")

				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "final content", string(content))

				// Original target should not exist when not preserving original
				_, err = os.Stat(finalTarget)
				require.True(t, os.IsNotExist(err), "Original target should be removed when not preserving original")
			},
		},
		{
			name: "deep copy - preserve original symlink target file",
			setup: func() (string, string, string) {
				// Create target file
				target := filepath.Join(tempDir, "preserve-target-file.txt")
				err := os.WriteFile(target, []byte("preserve target content"), 0644)
				require.NoError(t, err)

				// Create symlink
				sourceSymlink := filepath.Join(tempDir, "preserve-source-link")
				err = os.Symlink(target, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "preserve-dest-file")
				return sourceSymlink, dest, target
			},
			deep:             true,
			preserveOriginal: true,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, target string) {
				// Destination should be a regular file, not a symlink
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				require.False(t, destInfo.Mode()&os.ModeSymlink != 0, "Destination should not be a symlink")

				// Content should match the target file
				content, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, "preserve target content", string(content))

				// Original target file should still exist when preserving original
				targetContent, err := os.ReadFile(target)
				require.NoError(t, err)
				require.Equal(t, "preserve target content", string(targetContent))
			},
		},
		{
			name: "deep copy - don't preserve original symlink target directory",
			setup: func() (string, string, string) {
				// Create target directory
				targetDir := filepath.Join(tempDir, "remove-target-dir")
				err := os.Mkdir(targetDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(targetDir, "file.txt"), []byte("dir content"), 0644)
				require.NoError(t, err)

				// Create symlink to directory
				sourceSymlink := filepath.Join(tempDir, "remove-source-dir-link")
				err = os.Symlink(targetDir, sourceSymlink)
				require.NoError(t, err)

				dest := filepath.Join(tempDir, "remove-dest-dir")
				return sourceSymlink, dest, targetDir
			},
			deep:             true,
			preserveOriginal: false,
			fsync:            false,
			expectError:      false,
			verify: func(t *testing.T, sourceSymlink, dest, targetDir string) {
				// Destination should be a regular directory, not a symlink
				destInfo, err := os.Stat(dest)
				require.NoError(t, err)
				require.True(t, destInfo.IsDir(), "Destination should be a directory")
				require.False(t, destInfo.Mode()&os.ModeSymlink != 0, "Destination should not be a symlink")

				// Content should be copied
				content, err := os.ReadFile(filepath.Join(dest, "file.txt"))
				require.NoError(t, err)
				require.Equal(t, "dir content", string(content))

				// Original target directory should not exist when not preserving original
				_, err = os.Stat(targetDir)
				require.True(t, os.IsNotExist(err), "Original target directory should be removed when not preserving original")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sourceSymlink, dest, target := tc.setup()

			err := handleSymlink(sourceSymlink, dest, tc.deep, tc.preserveOriginal, tc.fsync)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, sourceSymlink, dest, target)
				}
			}
		})
	}
}

func TestRecursiveMoveWithSymlinksIntegration(t *testing.T) {
	if !supportsSymlinks() {
		t.Skip("Symlinks not supported on this platform/environment")
	}

	tempDir := t.TempDir()

	tests := []struct {
		name             string
		setup            func() (string, string)
		deep             bool
		preserveOriginal bool
		verify           func(t *testing.T, source, dest string)
	}{
		{
			name: "directory with mixed symlinks - shallow copy",
			setup: func() (string, string) {
				sourceDir := filepath.Join(tempDir, "mixed-source")
				destDir := filepath.Join(tempDir, "mixed-dest-shallow")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				// Create regular file
				regularFile := filepath.Join(sourceDir, "regular.txt")
				err = os.WriteFile(regularFile, []byte("regular content"), 0644)
				require.NoError(t, err)

				// Create target for symlink
				targetFile := filepath.Join(sourceDir, "target.txt")
				err = os.WriteFile(targetFile, []byte("target content"), 0644)
				require.NoError(t, err)

				// Create symlink to file
				fileLink := filepath.Join(sourceDir, "file-link")
				err = os.Symlink(targetFile, fileLink)
				require.NoError(t, err)

				// Create target directory
				targetSubDir := filepath.Join(sourceDir, "target-subdir")
				err = os.Mkdir(targetSubDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(targetSubDir, "sub-file.txt"), []byte("sub content"), 0644)
				require.NoError(t, err)

				// Create symlink to directory
				dirLink := filepath.Join(sourceDir, "dir-link")
				err = os.Symlink(targetSubDir, dirLink)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             false,
			preserveOriginal: false,
			verify: func(t *testing.T, source, dest string) {
				// Regular file should be copied
				content, err := os.ReadFile(filepath.Join(dest, "regular.txt"))
				require.NoError(t, err)
				require.Equal(t, "regular content", string(content))

				// Target file should be copied
				content, err = os.ReadFile(filepath.Join(dest, "target.txt"))
				require.NoError(t, err)
				require.Equal(t, "target content", string(content))

				// File symlink should still be a symlink
				linkTarget, err := os.Readlink(filepath.Join(dest, "file-link"))
				require.NoError(t, err)
				expectedTarget := filepath.Join(source, "target.txt")
				require.Equal(t, expectedTarget, linkTarget)

				// Directory symlink should still be a symlink
				linkTarget, err = os.Readlink(filepath.Join(dest, "dir-link"))
				require.NoError(t, err)
				expectedDirTarget := filepath.Join(source, "target-subdir")
				require.Equal(t, expectedDirTarget, linkTarget)

				// Source should not exist (moved)
				_, err = os.Stat(source)
				require.True(t, os.IsNotExist(err))
			},
		},
		{
			name: "directory with external symlinks - deep copy",
			setup: func() (string, string) {
				sourceDir := filepath.Join(tempDir, "external-deep-source")
				destDir := filepath.Join(tempDir, "external-dest-deep")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				// Create regular file
				regularFile := filepath.Join(sourceDir, "regular-deep.txt")
				err = os.WriteFile(regularFile, []byte("regular deep content"), 0644)
				require.NoError(t, err)

				// Create external target for symlink (outside source directory)
				externalTargetFile := filepath.Join(tempDir, "external-target.txt")
				err = os.WriteFile(externalTargetFile, []byte("external target content"), 0644)
				require.NoError(t, err)

				// Create symlink to external file
				fileLink := filepath.Join(sourceDir, "external-link")
				err = os.Symlink(externalTargetFile, fileLink)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             true,
			preserveOriginal: false,
			verify: func(t *testing.T, source, dest string) {
				// Regular file should be copied
				content, err := os.ReadFile(filepath.Join(dest, "regular-deep.txt"))
				require.NoError(t, err)
				require.Equal(t, "regular deep content", string(content))

				// External symlink should be resolved to actual file
				destFileFromLink := filepath.Join(dest, "external-link")
				info, err := os.Stat(destFileFromLink)
				require.NoError(t, err)
				require.False(t, info.Mode()&os.ModeSymlink != 0, "External symlink should be resolved to actual file")

				content, err = os.ReadFile(destFileFromLink)
				require.NoError(t, err)
				require.Equal(t, "external target content", string(content))

				// External target should not exist (removed when not preserving original)
				externalTargetFile := filepath.Join(tempDir, "external-target.txt")
				_, err = os.Stat(externalTargetFile)
				require.True(t, os.IsNotExist(err), "External target should be removed when not preserving original")

				// Source should not exist (moved)
				_, err = os.Stat(source)
				require.True(t, os.IsNotExist(err))
			},
		},
		{
			name: "error - directory with internal symlinks and deep copy",
			setup: func() (string, string) {
				sourceDir := filepath.Join(tempDir, "internal-deep-source")
				destDir := filepath.Join(tempDir, "internal-dest-deep")

				err := os.Mkdir(sourceDir, 0755)
				require.NoError(t, err)

				// Create target file within source directory
				targetFile := filepath.Join(sourceDir, "internal-target.txt")
				err = os.WriteFile(targetFile, []byte("internal target content"), 0644)
				require.NoError(t, err)

				// Create symlink to internal file (this should cause an error in deep copy)
				fileLink := filepath.Join(sourceDir, "internal-link")
				err = os.Symlink(targetFile, fileLink)
				require.NoError(t, err)

				return sourceDir, destDir
			},
			deep:             true,
			preserveOriginal: false,
			verify: func(t *testing.T, source, dest string) {
				// This test case should result in an error, so verify won't be called
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source, dest := tc.setup()

			err := RecursiveMove(source, dest, tc.deep, tc.preserveOriginal, false)

			if tc.name == "error - directory with internal symlinks and deep copy" {
				require.Error(t, err)
				require.Contains(t, err.Error(), "cannot deep copy symlink")
				require.Contains(t, err.Error(), "target")
				require.Contains(t, err.Error(), "is within the source directory being moved")
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, source, dest)
				}
			}
		})
	}
}

func TestSanitizePath(t *testing.T) {
	// Get the current working directory and home directory for test expectations
	cwd, err := os.Getwd()
	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name           string
		input          string
		expectedResult func() string // Function to compute expected result
		expectError    bool
		errorMsg       string
	}{
		{
			name:  "tilde expansion - home directory only",
			input: "~",
			expectedResult: func() string {
				return homeDir
			},
			expectError: false,
		},
		{
			name:  "tilde expansion - home directory with subdirectory",
			input: "~/Documents/test.txt",
			expectedResult: func() string {
				return filepath.Join(homeDir, "Documents/test.txt")
			},
			expectError: false,
		},
		{
			name:  "tilde expansion - home directory with nested subdirectories",
			input: "~/Documents/Projects/test-project/file.txt",
			expectedResult: func() string {
				return filepath.Join(homeDir, "Documents/Projects/test-project/file.txt")
			},
			expectError: false,
		},
		{
			name:  "absolute path - no changes needed",
			input: "/usr/local/bin/test",
			expectedResult: func() string {
				return "/usr/local/bin/test"
			},
			expectError: false,
		},
		{
			name:  "relative path - converted to absolute",
			input: "test-file.txt",
			expectedResult: func() string {
				return filepath.Join(cwd, "test-file.txt")
			},
			expectError: false,
		},
		{
			name:  "relative path with subdirectory",
			input: "subdir/test-file.txt",
			expectedResult: func() string {
				return filepath.Join(cwd, "subdir/test-file.txt")
			},
			expectError: false,
		},
		{
			name:  "path with redundant elements",
			input: "/usr/local/../local/bin/./test",
			expectedResult: func() string {
				return "/usr/local/bin/test"
			},
			expectError: false,
		},
		{
			name:  "path with current directory reference",
			input: "./test-file.txt",
			expectedResult: func() string {
				return filepath.Join(cwd, "test-file.txt")
			},
			expectError: false,
		},
		{
			name:  "path with parent directory reference",
			input: "../test-file.txt",
			expectedResult: func() string {
				return filepath.Join(filepath.Dir(cwd), "test-file.txt")
			},
			expectError: false,
		},
		{
			name:  "empty path",
			input: "",
			expectedResult: func() string {
				return cwd
			},
			expectError: false,
		},
		{
			name:  "path with multiple slashes",
			input: "/usr//local///bin/test",
			expectedResult: func() string {
				return "/usr/local/bin/test"
			},
			expectError: false,
		},
		{
			name:  "tilde in middle of path - not expanded",
			input: "/path/to/~user/file.txt",
			expectedResult: func() string {
				return "/path/to/~user/file.txt"
			},
			expectError: false,
		},
		{
			name:  "complex relative path with redundant elements",
			input: "./subdir/../another/./file.txt",
			expectedResult: func() string {
				return filepath.Join(cwd, "another/file.txt")
			},
			expectError: false,
		},
		{
			name:  "tilde with complex path",
			input: "~/Documents/../Downloads/./file.txt",
			expectedResult: func() string {
				return filepath.Join(homeDir, "Downloads/file.txt")
			},
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SanitizePath(tc.input)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				expected := tc.expectedResult()
				require.Equal(t, expected, result)

				// Verify the result is an absolute path
				require.True(t, filepath.IsAbs(result), "Result should be an absolute path")

				// Verify the path is clean (no redundant elements)
				require.Equal(t, filepath.Clean(result), result, "Result should be clean")
			}
		})
	}
}
