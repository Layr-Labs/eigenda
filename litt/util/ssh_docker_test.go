package util

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
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

// setupSSHTestContainer creates and starts a Docker container with SSH server
func setupSSHTestContainer(t *testing.T) *SSHTestContainer {
	ctx := context.Background()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	// Generate SSH key pair
	tempDir := t.TempDir()
	privateKeyPath := filepath.Join(tempDir, "test_ssh_key")
	publicKeyPath := filepath.Join(tempDir, "test_ssh_key.pub")

	err = generateSSHKeyPair(privateKeyPath, publicKeyPath)
	require.NoError(t, err)

	publicKeyContent, err := os.ReadFile(publicKeyPath)
	require.NoError(t, err)

	// Create mount directory for file operations
	mountDir := filepath.Join(tempDir, "ssh_mount")
	err = os.MkdirAll(mountDir, 0755)
	require.NoError(t, err)

	// Build Docker image
	imageName := "ssh-test:latest"
	err = buildSSHTestImage(ctx, cli, tempDir, imageName, string(publicKeyContent))
	require.NoError(t, err)

	// Start container
	containerID, sshPort, err := startSSHContainer(ctx, cli, imageName, mountDir)
	require.NoError(t, err)

	// Wait for SSH to be ready
	waitForSSH(t, sshPort, privateKeyPath)

	return &SSHTestContainer{
		client:      cli,
		containerID: containerID,
		sshPort:     sshPort,
		tempDir:     tempDir,
		privateKey:  privateKeyPath,
		publicKey:   publicKeyPath,
	}
}

// cleanup removes the Docker container and cleans up resources
func (c *SSHTestContainer) cleanup() error {
	ctx := context.Background()

	// Stop and remove container
	err := c.client.ContainerStop(ctx, c.containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	err = c.client.ContainerRemove(ctx, c.containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// buildSSHTestImage builds the SSH test image with the provided public key
func buildSSHTestImage(
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
	buildCtx, err := archiveDirectory(buildContext)
	if err != nil {
		return fmt.Errorf("failed to create build context archive: %w", err)
	}
	defer buildCtx.Close()

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
	defer response.Body.Close()

	// Read build output (helpful for debugging)
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return fmt.Errorf("failed to read build response: %w", err)
	}

	return nil
}

// startSSHContainer starts the SSH container and returns container ID and SSH port
func startSSHContainer(
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

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
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

// archiveDirectory creates a tar.gz archive of a directory for Docker build context
func archiveDirectory(srcDir string) (io.ReadCloser, error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		gw := gzip.NewWriter(pw)
		defer gw.Close()

		tw := tar.NewWriter(gw)
		defer tw.Close()

		filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
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
			defer file.Close()

			_, err = io.Copy(tw, file)
			return err
		})
	}()

	return pr, nil
}
