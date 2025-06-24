package util

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// SSHTestContainer manages a Docker container with SSH server for testing
type SSHTestContainer struct {
	client      *client.Client
	containerID string
	sshPort     string
	tempDir     string
	privateKey  string
	publicKey   string
}

// GetSSHPort returns the SSH port of the test container
func (c *SSHTestContainer) GetSSHPort() string {
	return c.sshPort
}

// GetPrivateKeyPath returns the path to the private key file
func (c *SSHTestContainer) GetPrivateKeyPath() string {
	return c.privateKey
}

// GetPublicKeyPath returns the path to the public key file
func (c *SSHTestContainer) GetPublicKeyPath() string {
	return c.publicKey
}

// GetTempDir returns the temporary directory used by the container
func (c *SSHTestContainer) GetTempDir() string {
	return c.tempDir
}

// Cleanup removes the Docker container and cleans up resources
func (c *SSHTestContainer) Cleanup() error {
	ctx := context.Background()

	// Stop and remove container
	err := c.client.ContainerStop(ctx, c.containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	err = c.client.ContainerRemove(ctx, c.containerID, container.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// ParsePort converts string port to uint64
func ParsePort(port string) uint64 {
	var p uint64
	_, _ = fmt.Sscanf(port, "%d", &p)
	return p
}

// GenerateSSHKeyPair creates an RSA key pair for testing
func GenerateSSHKeyPair(privateKeyPath, publicKeyPath string) error {
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

// WaitForSSH waits for the SSH server to be ready
func WaitForSSH(t *testing.T, sshPort, privateKeyPath string) {
	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	require.NoError(t, err)

	// Try to connect multiple times with backoff
	for i := 0; i < 30; i++ {
		session, err := NewSSHSession(
			logger,
			"testuser",
			"localhost",
			ParsePort(sshPort),
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

// SetupSSHTestContainer creates and starts a Docker container with SSH server
func SetupSSHTestContainer(t *testing.T) *SSHTestContainer {
	ctx := context.Background()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	// Generate SSH key pair
	tempDir := t.TempDir()
	privateKeyPath := filepath.Join(tempDir, "test_ssh_key")
	publicKeyPath := filepath.Join(tempDir, "test_ssh_key.pub")

	err = GenerateSSHKeyPair(privateKeyPath, publicKeyPath)
	require.NoError(t, err)

	publicKeyContent, err := os.ReadFile(publicKeyPath)
	require.NoError(t, err)

	// Create mount directory for file operations
	mountDir := filepath.Join(tempDir, "ssh_mount")
	err = os.MkdirAll(mountDir, 0755)
	require.NoError(t, err)

	// Build Docker image
	imageName := "ssh-test:latest"
	err = BuildSSHTestImage(ctx, cli, tempDir, imageName, string(publicKeyContent))
	require.NoError(t, err)

	// Start container
	containerID, sshPort, err := StartSSHContainer(ctx, cli, imageName, mountDir)
	require.NoError(t, err)

	// Wait for SSH to be ready
	WaitForSSH(t, sshPort, privateKeyPath)

	return &SSHTestContainer{
		client:      cli,
		containerID: containerID,
		sshPort:     sshPort,
		tempDir:     tempDir,
		privateKey:  privateKeyPath,
		publicKey:   publicKeyPath,
	}
}

// BuildSSHTestImage builds the SSH test image with the provided public key
func BuildSSHTestImage(
	ctx context.Context,
	cli *client.Client,
	tempDir string,
	imageName string,
	publicKey string,
) error {

	// Get the Dockerfile path
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}
	dockerfilePath := filepath.Join(filepath.Dir(currentFile), "testdata", "ssh-test.Dockerfile")

	// Create build context directory
	buildContext := filepath.Join(tempDir, "docker_build")
	err := os.MkdirAll(buildContext, 0755)
	if err != nil {
		return fmt.Errorf("failed to create build context: %w", err)
	}

	// Copy Dockerfile to build context
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	// Add the public key setup to the Dockerfile
	publicKeySetup := fmt.Sprintf(
		"\n# Add test SSH public key\n"+
			"RUN echo '%s' > /home/testuser/.ssh/authorized_keys\n"+
			"RUN chmod 600 /home/testuser/.ssh/authorized_keys\n"+
			"RUN chown testuser:testuser /home/testuser/.ssh/authorized_keys\n", strings.TrimSpace(publicKey))
	modifiedDockerfile := string(dockerfileContent) + publicKeySetup

	err = os.WriteFile(filepath.Join(buildContext, "Dockerfile"), []byte(modifiedDockerfile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write modified Dockerfile: %w", err)
	}

	// Create tar archive for build context
	buildCtx, err := ArchiveDirectory(buildContext)
	if err != nil {
		return fmt.Errorf("failed to create build context archive: %w", err)
	}
	defer func() { _ = buildCtx.Close() }()

	// Build the image
	buildOptions := types.ImageBuildOptions{
		Tags:        []string{imageName},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
	}

	response, err := cli.ImageBuild(ctx, buildCtx, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	// Read build output (helpful for debugging)
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return fmt.Errorf("failed to read build response: %w", err)
	}

	return nil
}

// StartSSHContainer starts the SSH container and returns container ID and SSH port
func StartSSHContainer(
	ctx context.Context,
	cli *client.Client,
	imageName string,
	mountDir string,
) (string, string, error) {

	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"22/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"22/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0", // Let Docker assign a random port
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountDir,
				Target: "/mnt/test",
			},
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"")
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %w", err)
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to start container: %w", err)
	}

	// Get the assigned SSH port
	containerInfo, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to inspect container: %w", err)
	}

	sshPort := containerInfo.NetworkSettings.Ports["22/tcp"][0].HostPort

	return resp.ID, sshPort, nil
}

// ArchiveDirectory creates a tar.gz archive of a directory for Docker build context
func ArchiveDirectory(srcDir string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer func() { _ = pw.Close() }()

		gw := gzip.NewWriter(pw)
		defer func() { _ = gw.Close() }()

		tw := tar.NewWriter(gw)
		defer func() { _ = tw.Close() }()

		_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}

			// Skip the root directory itself
			if relPath == "." {
				return nil
			}

			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = relPath

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() { _ = file.Close() }()

			_, err = io.Copy(tw, file)
			return err
		})
	}()

	return pr, nil
}