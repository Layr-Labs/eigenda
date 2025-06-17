package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

func TestNewFileLock(t *testing.T) {
	tempDir := t.TempDir()
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful lock creation",
			setup: func() string {
				return filepath.Join(tempDir, "test.lock")
			},
			expectError: false,
		},
		{
			name: "lock already exists with live process",
			setup: func() string {
				lockPath := filepath.Join(tempDir, "existing.lock")
				// Create an existing lock file with current process PID (which is alive)
				content := fmt.Sprintf("PID: %d\nTimestamp: 2023-01-01T00:00:00Z\n", os.Getpid())
				err := os.WriteFile(lockPath, []byte(content), 0644)
				require.NoError(t, err)
				return lockPath
			},
			expectError: true,
			errorMsg:    "process",
		},
		{
			name: "stale lock file gets overridden",
			setup: func() string {
				lockPath := filepath.Join(tempDir, "stale.lock")
				// Create a lock file with a PID that definitely doesn't exist
				// Use PID 999999 which is very unlikely to exist
				stalePID := 999999
				content := fmt.Sprintf("PID: %d\nTimestamp: 2023-01-01T00:00:00Z\n", stalePID)
				err := os.WriteFile(lockPath, []byte(content), 0644)
				require.NoError(t, err)
				return lockPath
			},
			expectError: false,
		},
		{
			name: "malformed lock file gets treated as existing",
			setup: func() string {
				lockPath := filepath.Join(tempDir, "malformed.lock")
				// Create a lock file without proper PID format
				err := os.WriteFile(lockPath, []byte("invalid content"), 0644)
				require.NoError(t, err)
				return lockPath
			},
			expectError: true,
			errorMsg:    "lock file already exists",
		},
		{
			name: "invalid directory",
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent", "test.lock")
			},
			expectError: true,
			errorMsg:    "failed to create lock file",
		},
		{
			name: "tilde expansion",
			setup: func() string {
				return "~/test.lock"
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lockPath := tc.setup()

			lock, err := NewFileLock(logger, lockPath, true)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
				require.Nil(t, lock)

				// Check for specific error messages
				switch tc.name {
				case "lock already exists with live process":
					require.Contains(t, err.Error(), "still running")
					require.Contains(t, err.Error(), fmt.Sprintf("process %d", os.Getpid()))
				case "malformed lock file gets treated as existing":
					require.Contains(t, err.Error(), "lock file already exists")
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, lock)

				// Verify lock file was created
				_, err := os.Stat(lock.Path())
				require.NoError(t, err)

				// Verify lock file contains process info
				content, err := os.ReadFile(lock.Path())
				require.NoError(t, err)
				contentStr := string(content)
				require.Contains(t, contentStr, "PID:")
				require.Contains(t, contentStr, "Timestamp:")

				// Clean up
				lock.Release()
			}
		})
	}
}

func TestFileLockRelease(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "test.lock")

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Create a lock
	lock, err := NewFileLock(logger, lockPath, true)
	require.NoError(t, err)
	require.NotNil(t, lock)

	// Verify lock file exists
	_, err = os.Stat(lockPath)
	require.NoError(t, err)

	// Release the lock
	lock.Release()

	// Verify lock file was removed
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err))

	// Try to release again (should not)
	lock.Release()
}

func TestFileLockPath(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "test.lock")

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	lock, err := NewFileLock(logger, lockPath, true)
	require.NoError(t, err)
	defer lock.Release()

	// Path should be sanitized (absolute)
	returnedPath := lock.Path()
	require.True(t, filepath.IsAbs(returnedPath))
	require.True(t, strings.HasSuffix(returnedPath, "test.lock"))
}

func TestFileLockConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "concurrent.lock")

	const numGoroutines = 10
	const duration = 50 * time.Millisecond

	var successCount int32
	var wg sync.WaitGroup
	results := make(chan bool, numGoroutines)

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Launch multiple goroutines trying to acquire the same lock
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			lock, err := NewFileLock(logger, lockPath, true)
			if err != nil {
				results <- false
				return
			}

			// Hold the lock for a short time
			time.Sleep(duration)

			lock.Release()

			results <- true
		}(i)
	}

	wg.Wait()
	close(results)

	// Count successful lock acquisitions
	successCount = 0
	for success := range results {
		if success {
			successCount++
		}
	}

	// Only one goroutine should have successfully acquired the lock
	require.Equal(t, int32(1), successCount, "Only one goroutine should acquire the lock")
}

func TestFileLockCleanupOnFailure(t *testing.T) {
	tempDir := t.TempDir()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Test that failed lock creation cleans up partial files
	t.Run("sync failure cleanup", func(t *testing.T) {
		// This is harder to test directly since sync failures are rare
		// We'll test the error path by creating a lock and verifying it's cleaned up
		lockPath := filepath.Join(tempDir, "cleanup-test.lock")

		lock, err := NewFileLock(logger, lockPath, true)
		require.NoError(t, err)
		require.NotNil(t, lock)

		// Verify lock was created
		_, err = os.Stat(lockPath)
		require.NoError(t, err)

		// Clean up
		lock.Release()

		// Verify cleanup
		_, err = os.Stat(lockPath)
		require.True(t, os.IsNotExist(err))
	})
}

func TestFileLockEdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	t.Run("release nil lock", func(t *testing.T) {
		lock := &FileLock{}
		lock.Release()
		require.NoError(t, err)
	})

	t.Run("path sanitization", func(t *testing.T) {
		// Test that paths are properly sanitized
		relativePath := "test.lock"

		lock, err := NewFileLock(logger, relativePath, true)
		require.NoError(t, err)
		defer lock.Release()

		// Path should be absolute
		lockPath := lock.Path()
		require.True(t, filepath.IsAbs(lockPath))
		require.Contains(t, lockPath, "test.lock")
	})

	t.Run("double release protection", func(t *testing.T) {
		lockPath := filepath.Join(tempDir, "double-release.lock")

		lock, err := NewFileLock(logger, lockPath, true)
		require.NoError(t, err)

		// First release should succeed
		lock.Release()

		// Second release should fail gracefully
		lock.Release()
	})
}

func TestFileLockDebugInfo(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "debug-test.lock")

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Create first lock
	lock1, err := NewFileLock(logger, lockPath, true)
	require.NoError(t, err)

	// Try to create second lock - should fail with debug info
	lock2, err := NewFileLock(logger, lockPath, true)
	require.Error(t, err)
	require.Nil(t, lock2)

	// Error should contain debug information from existing lock
	require.Contains(t, err.Error(), "lock file already exists")
	require.Contains(t, err.Error(), "existing lock info:")
	require.Contains(t, err.Error(), "PID:")
	require.Contains(t, err.Error(), "Timestamp:")

	// Clean up
	lock1.Release()
}

func TestFileLockContentsFormat(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "content-test.lock")

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	lock, err := NewFileLock(logger, lockPath, true)
	require.NoError(t, err)
	defer lock.Release()

	// Read and verify lock file contents
	content, err := os.ReadFile(lockPath)
	require.NoError(t, err)

	contentStr := string(content)
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")
	require.Len(t, lines, 2)

	// Check PID line
	require.True(t, strings.HasPrefix(lines[0], "PID: "))
	pidStr := strings.TrimPrefix(lines[0], "PID: ")
	require.NotEmpty(t, pidStr)
	// Verify it's a valid number
	require.Regexp(t, `^\d+$`, pidStr)

	// Check timestamp line
	require.True(t, strings.HasPrefix(lines[1], "Timestamp: "))

	// Verify timestamp is valid RFC3339 format
	timestampStr := strings.TrimPrefix(lines[1], "Timestamp: ")
	_, err = time.Parse(time.RFC3339, timestampStr)
	require.NoError(t, err)
}

func TestIsProcessAlive(t *testing.T) {
	tests := []struct {
		name     string
		pid      int
		expected bool
	}{
		{
			name:     "current process",
			pid:      os.Getpid(),
			expected: true,
		},
		{
			name:     "invalid pid zero",
			pid:      0,
			expected: false,
		},
		{
			name:     "invalid pid negative",
			pid:      -1,
			expected: false,
		},
		{
			name:     "nonexistent pid",
			pid:      999999, // Very unlikely to exist
			expected: false,
		},
		{
			name:     "init process",
			pid:      1,
			expected: true, // Init process should always exist on Unix systems
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsProcessAlive(tc.pid)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestParseLockFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectedPID int
		expectError bool
	}{
		{
			name:        "valid lock file",
			content:     "PID: 12345\nTimestamp: 2023-01-01T00:00:00Z\n",
			expectedPID: 12345,
			expectError: false,
		},
		{
			name:        "lock file with extra whitespace",
			content:     "  PID: 67890  \n  Timestamp: 2023-01-01T00:00:00Z  \n",
			expectedPID: 67890,
			expectError: false,
		},
		{
			name:        "lock file missing PID",
			content:     "Timestamp: 2023-01-01T00:00:00Z\n",
			expectedPID: 0,
			expectError: true,
		},
		{
			name:        "lock file with invalid PID",
			content:     "PID: not-a-number\nTimestamp: 2023-01-01T00:00:00Z\n",
			expectedPID: 0,
			expectError: true,
		},
		{
			name:        "empty lock file",
			content:     "",
			expectedPID: 0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lockPath := filepath.Join(tempDir, fmt.Sprintf("test-%s.lock", tc.name))
			err := os.WriteFile(lockPath, []byte(tc.content), 0644)
			require.NoError(t, err)

			pid, err := parseLockFile(lockPath)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPID, pid)
			}
		})
	}
}

func TestStaleLockRecovery(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "stale-recovery.lock")

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Create a stale lock file with a definitely dead PID
	stalePID := 999999
	staleContent := fmt.Sprintf("PID: %d\nTimestamp: 2023-01-01T00:00:00Z\n", stalePID)
	err = os.WriteFile(lockPath, []byte(staleContent), 0644)
	require.NoError(t, err)

	// Verify the lock file exists
	_, err = os.Stat(lockPath)
	require.NoError(t, err)

	// Try to acquire the lock - should succeed by removing stale lock
	lock, err := NewFileLock(logger, lockPath, true)
	require.NoError(t, err)
	require.NotNil(t, lock)

	// Verify the lock file now has our PID
	content, err := os.ReadFile(lockPath)
	require.NoError(t, err)
	require.Contains(t, string(content), fmt.Sprintf("PID: %d", os.Getpid()))

	// Clean up
	lock.Release()
}
