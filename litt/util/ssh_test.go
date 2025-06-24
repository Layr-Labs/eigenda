package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

func TestSSHSession_NewSSHSession(t *testing.T) {
	t.Parallel()

	container := setupSSHTestContainer(t)
	defer func() { _ = container.cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	// Test successful connection
	session, err := NewSSHSession(
		logger,
		"testuser",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		true)
	require.NoError(t, err)
	require.NotNil(t, session)
	defer func() { _ = session.Close() }()

	// Test with non-existent key
	_, err = NewSSHSession(
		logger,
		"testuser",
		"localhost",
		parsePort(container.sshPort),
		"/nonexistent/key",
		false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "private key does not exist")

	// Test with wrong user
	_, err = NewSSHSession(
		logger,
		"wronguser",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		false)
	require.Error(t, err)
}

func TestSSHSession_Ls(t *testing.T) {
	t.Parallel()

	container := setupSSHTestContainer(t)
	defer func() { _ = container.cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		"testuse"+
			"r",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Test listing home directory
	files, err := session.Ls("/home/testuser")
	require.NoError(t, err)
	require.Contains(t, files, ".ssh")

	// Test listing non-existent directory
	_, err = session.Ls("/nonexistent")
	require.Error(t, err)
}

func TestSSHSession_Mkdirs(t *testing.T) {
	t.Parallel()

	container := setupSSHTestContainer(t)
	defer func() { _ = container.cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		"testuser",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Test creating directory
	testDir := "/mnt/test/newdir/subdir"
	err = session.Mkdirs(testDir)
	require.NoError(t, err)

	// Verify directory was created
	files, err := session.Ls("/mnt/test/newdir")
	require.NoError(t, err)
	require.Contains(t, files, "subdir")
}

func TestSSHSession_FindFiles(t *testing.T) {
	t.Parallel()

	container := setupSSHTestContainer(t)
	defer func() { _ = container.cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		"testuser",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Create test directory structure
	err = session.Mkdirs("/mnt/test/search")
	require.NoError(t, err)

	// Create test files using the mounted directory
	mountDir := filepath.Join(container.tempDir, "ssh_mount", "search")
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

	container := setupSSHTestContainer(t)
	defer func() { _ = container.cleanup() }()

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	session, err := NewSSHSession(
		logger,
		"testuser",
		"localhost",
		parsePort(container.sshPort),
		container.privateKey,
		true)
	require.NoError(t, err)
	defer func() { _ = session.Close() }()

	// Create test directory on remote
	err = session.Mkdirs("/mnt/test/rsync")
	require.NoError(t, err)

	// Create local test file
	localFile := filepath.Join(container.tempDir, "test_rsync.txt")
	testContent := []byte("This is test content for rsync")
	err = os.WriteFile(localFile, testContent, 0644)
	require.NoError(t, err)

	// Test rsync without throttling
	err = session.Rsync(localFile, "/mnt/test/rsync/remote_file.txt", 0)
	require.NoError(t, err)

	// Verify file was transferred (via mounted directory)
	mountedFile := filepath.Join(container.tempDir, "ssh_mount", "rsync", "remote_file.txt")
	transferredContent, err := os.ReadFile(mountedFile)
	require.NoError(t, err)
	require.Equal(t, testContent, transferredContent)

	// Test rsync with throttling
	localFile2 := filepath.Join(container.tempDir, "test_rsync2.txt")
	err = os.WriteFile(localFile2, []byte("throttled content"), 0644)
	require.NoError(t, err)

	err = session.Rsync(localFile2, "/mnt/test/rsync/throttled_file.txt", 1.0) // 1MB/s throttle
	require.NoError(t, err)

	// Verify throttled file was transferred
	mountedFile2 := filepath.Join(container.tempDir, "ssh_mount", "rsync", "throttled_file.txt")
	_, err = os.Stat(mountedFile2)
	require.NoError(t, err)
}
