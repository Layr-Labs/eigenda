package util

import (
	"os"
	"path/filepath"
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
			name:        "directory already exists",
			dirPath:     filepath.Join(tempDir, "existing-dir"),
			mode:        0755,
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
			name:        "directory exists but is non-writable",
			dirPath:     filepath.Join(tempDir, "non-writable-dir"),
			mode:        0755,
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

func TestCopyDirectoryRecursively(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	
	// Create a source directory structure
	sourceDir := filepath.Join(tempDir, "source")
	err := os.Mkdir(sourceDir, 0755)
	require.NoError(t, err)
	
	// Create some files in the source directory
	err = os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("file1 content"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("file2 content"), 0644)
	require.NoError(t, err)
	
	// Create a subdirectory
	subDir := filepath.Join(sourceDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	
	// Create files in the subdirectory
	err = os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("file3 content"), 0644)
	require.NoError(t, err)
	
	// Create a symlink if supported
	supportsLinks := supportsSymlinks()
	var symlinkPath string
	if supportsLinks {
		symlinkPath = filepath.Join(sourceDir, "symlink")
		err = os.Symlink(filepath.Join(sourceDir, "file1.txt"), symlinkPath)
		require.NoError(t, err)
	}
	
	// Test cases
	tests := []struct {
		name        string
		destDir     string
		expectError bool
	}{
		{
			name:        "copy to new destination",
			destDir:     filepath.Join(tempDir, "dest"),
			expectError: false,
		},
		{
			name:        "copy to existing destination",
			destDir:     filepath.Join(tempDir, "existing-dest"),
			expectError: false,
		},
	}
	
	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// For the "existing destination" test, create the directory first
			if tc.name == "copy to existing destination" {
				err := os.Mkdir(tc.destDir, 0755)
				require.NoError(t, err)
				
				// Add a pre-existing file
				err = os.WriteFile(filepath.Join(tc.destDir, "existing.txt"), []byte("existing content"), 0644)
				require.NoError(t, err)
			}
			
			err := CopyDirectoryRecursively(sourceDir, tc.destDir)
			
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				
				// Verify the directory was copied correctly
				
				// Check for file1.txt
				content, err := os.ReadFile(filepath.Join(tc.destDir, "file1.txt"))
				require.NoError(t, err)
				require.Equal(t, "file1 content", string(content))
				
				// Check for file2.txt
				content, err = os.ReadFile(filepath.Join(tc.destDir, "file2.txt"))
				require.NoError(t, err)
				require.Equal(t, "file2 content", string(content))
				
				// Check for subdirectory and its file
				subDirPath := filepath.Join(tc.destDir, "subdir")
				info, err := os.Stat(subDirPath)
				require.NoError(t, err)
				require.True(t, info.IsDir())
				
				content, err = os.ReadFile(filepath.Join(subDirPath, "file3.txt"))
				require.NoError(t, err)
				require.Equal(t, "file3 content", string(content))
				
				// Check for symlink if supported
				if supportsLinks {
					linkTarget, err := os.Readlink(filepath.Join(tc.destDir, "symlink"))
					require.NoError(t, err)
					require.Equal(t, filepath.Join(sourceDir, "file1.txt"), linkTarget)
				}
				
				// For the "existing destination" test, verify the pre-existing file is still there
				if tc.name == "copy to existing destination" {
					content, err = os.ReadFile(filepath.Join(tc.destDir, "existing.txt"))
					require.NoError(t, err)
					require.Equal(t, "existing content", string(content))
				}
			}
		})
	}
}

// Helper function to check if symlinks are supported in the current environment
func supportsSymlinks() bool {
	tempDir, err := os.MkdirTemp("", "symlink-test")
	if err != nil {
		return false
	}
	defer os.RemoveAll(tempDir)
	
	source := filepath.Join(tempDir, "source")
	target := filepath.Join(tempDir, "target")
	
	err = os.WriteFile(source, []byte{}, 0644)
	if err != nil {
		return false
	}
	
	err = os.Symlink(source, target)
	return err == nil
}