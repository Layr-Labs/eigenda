package containers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	AnvilImage = "ghcr.io/foundry-rs/foundry:latest"
	AnvilPort  = "8545/tcp"
)

// AnvilContainer wraps testcontainers functionality for Anvil
type AnvilContainer struct {
	container testcontainers.Container
	config    AnvilConfig
	rpcURL    string
}

// NewAnvilContainer creates and starts a new Anvil container
func NewAnvilContainer(ctx context.Context, config AnvilConfig) (*AnvilContainer, error) {
	return NewAnvilContainerWithNetwork(ctx, config, "")
}

// NewAnvilContainerWithNetwork creates and starts a new Anvil container in a specific network
func NewAnvilContainerWithNetwork(
	ctx context.Context, config AnvilConfig, networkName string,
) (*AnvilContainer, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("anvil container is disabled in config")
	}

	args := buildAnvilArgs(config)
	fmt.Printf("DEBUG: Anvil command will be: anvil %s\n", strings.Join(args, " "))

	// Generate a unique container name using timestamp to avoid conflicts in parallel tests
	uniqueName := fmt.Sprintf("anvil-test-%d-%d", config.ChainID, time.Now().UnixNano())

	req := testcontainers.ContainerRequest{
		Image:        AnvilImage,
		Cmd:          append([]string{"anvil"}, args...),
		ExposedPorts: []string{AnvilPort},
		Env:          map[string]string{"ANVIL_IP_ADDR": "0.0.0.0"},
		WaitingFor:   wait.ForListeningPort("8545/tcp"),
		Name:         uniqueName,
	}

	// Add network if specified
	if networkName != "" {
		req.Networks = []string{networkName}
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, "8545")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the host
	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	rpcURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	return &AnvilContainer{
		container: container,
		config:    config,
		rpcURL:    rpcURL,
	}, nil
}

// RPCURL returns the RPC endpoint URL for the Anvil node
func (a *AnvilContainer) RPCURL() string {
	return a.rpcURL
}

// ChainID returns the configured chain ID
func (a *AnvilContainer) ChainID() int {
	return a.config.ChainID
}

// Accounts returns the number of pre-funded accounts
func (a *AnvilContainer) Accounts() int {
	return a.config.Accounts
}

// Mnemonic returns the mnemonic used for account generation
func (a *AnvilContainer) Mnemonic() string {
	return a.config.Mnemonic
}

// ContainerName returns the container name for internal network communication
func (a *AnvilContainer) ContainerName(ctx context.Context) (string, error) {
	name, err := a.container.Name(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container name: %w", err)
	}
	// Remove leading slash that Docker adds to container names
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	return name, nil
}

// InternalRPCURL returns the RPC URL for internal container network communication
func (a *AnvilContainer) InternalRPCURL(ctx context.Context) (string, error) {
	name, err := a.ContainerName(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:8545", name), nil
}

// Terminate stops and removes the container
func (a *AnvilContainer) Terminate(ctx context.Context) error {
	if a.container != nil {
		return a.container.Terminate(ctx)
	}
	return nil
}

// GetContainer returns the underlying testcontainer for debugging
func (a *AnvilContainer) GetContainer() testcontainers.Container {
	return a.container
}

// GetPrivateKey returns the private key for the given account index (0-based)
func (a *AnvilContainer) GetPrivateKey(accountIndex int) (string, error) {
	if accountIndex >= a.config.Accounts {
		return "", fmt.Errorf("account index %d exceeds available accounts %d", accountIndex, a.config.Accounts)
	}

	// Anvil uses deterministic private keys based on the mnemonic
	// For the default mnemonic, these are the well-known private keys
	knownKeys := []string{
		"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		"0x5de4111afa1a4b94908f83103c54a14de4a2e7938c0f15d2e0b1b7bcbe1e7e3b",
		"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
		"0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		"0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
		"0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
		"0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
		"0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
		"0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
	}

	if accountIndex < len(knownKeys) {
		return knownKeys[accountIndex], nil
	}

	return "", fmt.Errorf("private key not available for account index %d", accountIndex)
}

// buildAnvilArgs constructs the command line arguments for Anvil
func buildAnvilArgs(config AnvilConfig) []string {
	args := []string{
		"--host", "0.0.0.0",
		"--port", "8545",
		"--chain-id", strconv.Itoa(config.ChainID),
		"--accounts", strconv.Itoa(config.Accounts),
	}

	if config.BlockTime > 0 {
		args = append(args, "--block-time", strconv.Itoa(config.BlockTime))
	}

	if config.GasLimit > 0 {
		args = append(args, "--gas-limit", strconv.FormatUint(config.GasLimit, 10))
	}

	if config.GasPrice > 0 {
		args = append(args, "--gas-price", strconv.FormatUint(config.GasPrice, 10))
	}

	if config.Mnemonic != "" {
		args = append(args, "--mnemonic", config.Mnemonic)
	}

	if config.Fork != "" {
		args = append(args, "--fork-url", config.Fork)
		if config.ForkBlock > 0 {
			args = append(args, "--fork-block-number", strconv.FormatUint(config.ForkBlock, 10))
		}
	}

	return args
}

// WaitForReady waits for the Anvil node to be ready to accept requests
func (a *AnvilContainer) WaitForReady(ctx context.Context) error {
	// The wait strategy in the container request should handle this,
	// but we can add additional checks here if needed
	return nil
}

// GetLogs returns the container logs for debugging
func (a *AnvilContainer) GetLogs(ctx context.Context) (string, error) {
	if a.container == nil {
		return "", fmt.Errorf("container not started")
	}

	logs, err := a.container.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	buf := make([]byte, 1024*1024) // 1MB buffer
	n, err := logs.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(buf[:n]), nil
}
