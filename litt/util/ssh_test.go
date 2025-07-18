package util

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

func TestSSHSession_NewSSHSession(t *testing.T) {
	t.Parallel()

	container := SetupSSHTestContainer(t, "")
	defer func() { _ = container.Cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	// Test successful connection
	session, err := NewSSHSession(
		logger,
		container.GetUser(),
		container.GetHost(),
		container.GetSSHPort(),
		container.GetPrivateKeyPath(),
		true)
	require.NoError(t, err)
	require.NotNil(t, session)
	defer func() { _ = session.Close() }()

	// Test with non-existent key
	_, err = NewSSHSession(
		logger,
		container.GetUser(),
		container.GetHost(),
		container.GetSSHPort(),
		"/nonexistent/key",
		false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "private key does not exist")

	// Test with wrong user
	_, err = NewSSHSession(
		logger,
		"wronguser",
		container.GetHost(),
		container.GetSSHPort(),
		container.GetPrivateKeyPath(),
		false)
	require.Error(t, err)
}

func TestSSHSession_Mkdirs(t *testing.T) {
	t.Parallel()

	dataDir := t.TempDir()

	container := SetupSSHTestContainer(t, dataDir)
	defer func() { _ = container.Cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		container.GetUser(),
		container.GetHost(),
		container.GetSSHPort(),
		container.GetPrivateKeyPath(),
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Test creating directory
	testDir := path.Join(container.GetDataDir(), "foo", "bar", "baz")
	err = session.Mkdirs(testDir)
	require.NoError(t, err)

	// Verify directories were created
	exists, err := Exists(path.Join(dataDir, "foo"))
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = Exists(path.Join(dataDir, "foo", "bar"))
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = Exists(path.Join(dataDir, "foo", "bar", "baz"))
	require.NoError(t, err)
	require.True(t, exists)

	// Recreating the same directory should not error.
	err = session.Mkdirs(testDir)
	require.NoError(t, err)
}

func TestSSHSession_FindFiles(t *testing.T) {
	t.Parallel()

	container := SetupSSHTestContainer(t, "")
	defer func() { _ = container.Cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		container.GetUser(),
		container.GetHost(),
		container.GetSSHPort(),
		container.GetPrivateKeyPath(),
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Create test directory structure
	err = session.Mkdirs("/mnt/test/search")
	require.NoError(t, err)

	// Create test files using the mounted directory
	mountDir := filepath.Join(container.GetTempDir(), "ssh_mount", "search")
	err = os.MkdirAll(mountDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(mountDir, "test.txt"), []byte("test content"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(mountDir, "test.log"), []byte("log content"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(mountDir, "other.dat"), []byte("data content"), 0644)
	require.NoError(t, err)

	// Test finding files with specific extensions
	files, err := session.FindFiles("/mnt/test/search", []string{".txt", ".log"})
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Contains(t, files, "/mnt/test/search/test.txt")
	require.Contains(t, files, "/mnt/test/search/test.log")

	// Test with non-existent directory
	files, err = session.FindFiles("/nonexistent", []string{".txt"})
	require.NoError(t, err)
	require.Empty(t, files)
}

func TestSSHSession_Rsync(t *testing.T) {
	t.Parallel()

	// Create a temporary data directory for testing
	dataDir := t.TempDir()
	container := SetupSSHTestContainer(t, dataDir)
	defer func() { _ = container.Cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		container.GetUser(),
		container.GetHost(),
		container.GetSSHPort(),
		container.GetPrivateKeyPath(),
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Create local test file
	localFile := filepath.Join(container.GetTempDir(), "test_rsync.txt")
	testContent := []byte("This is test content for rsync")
	err = os.WriteFile(localFile, testContent, 0644)
	require.NoError(t, err)

	// Test rsync without throttling - sync to data directory
	remoteFile := filepath.Join(container.GetDataDir(), "remote_file.txt")
	err = session.Rsync(localFile, remoteFile, 0)
	require.NoError(t, err)

	// Verify file was transferred via the mounted data directory
	transferredFile := filepath.Join(dataDir, "remote_file.txt")
	transferredContent, err := os.ReadFile(transferredFile)
	require.NoError(t, err)
	require.Equal(t, testContent, transferredContent)

	// Test rsync with throttling
	localFile2 := filepath.Join(container.GetTempDir(), "test_rsync2.txt")
	throttledContent := []byte("throttled content")
	err = os.WriteFile(localFile2, throttledContent, 0644)
	require.NoError(t, err)

	remoteFile2 := filepath.Join(container.GetDataDir(), "throttled_file.txt")
	err = session.Rsync(localFile2, remoteFile2, 1.0) // 1MB/s throttle
	require.NoError(t, err)

	// Verify throttled file was transferred via the mounted data directory
	transferredFile2 := filepath.Join(dataDir, "throttled_file.txt")
	transferredContent2, err := os.ReadFile(transferredFile2)
	require.NoError(t, err)
	require.Equal(t, throttledContent, transferredContent2)
}
