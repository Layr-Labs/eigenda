package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// generateSSHKeyPair creates an RSA key pair for testing
func generateSSHKeyPair(privateKeyPath, publicKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Save private key
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer func() { _ = privateKeyFile.Close() }()

	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	err = os.Chmod(privateKeyPath, 0600)
	if err != nil {
		return fmt.Errorf("failed to set private key permissions: %w", err)
	}

	// Save public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to create SSH public key: %w", err)
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)
	err = os.WriteFile(publicKeyPath, publicKeyBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

// waitForSSH waits for the SSH server to be ready
func waitForSSH(t *testing.T, sshPort, privateKeyPath string) {
	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	// Try to connect multiple times with backoff
	for i := 0; i < 30; i++ {
		session, err := NewSSHSession(
			logger,
			"testuser",
			"localhost",
			parsePort(sshPort),
			privateKeyPath,
			false)
		if err == nil {
			_ = session.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}

	require.Fail(t, "SSH server did not become ready in time")
}

// parsePort converts string port to uint64
func parsePort(port string) uint64 {
	var p uint64
	_, _ = fmt.Sscanf(port, "%d", &p)
	return p
}

// Test functions

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
