package testbed

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
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
	rpcURL    string
}

// AnvilOptions configures the Anvil container
type AnvilOptions struct {
	ExposeHostPort bool   // If true, binds container port 8545 to host port 8545
	HostPort       string // Custom host port to bind to (defaults to "8545" if empty and ExposeHostPort is true)
}

// NewAnvilContainer creates and starts a new Anvil container
func NewAnvilContainer(ctx context.Context) (*AnvilContainer, error) {
	return NewAnvilContainerWithOptions(ctx, AnvilOptions{})
}

// NewAnvilContainerWithOptions creates and starts a new Anvil container with custom options
func NewAnvilContainerWithOptions(ctx context.Context, opts AnvilOptions) (*AnvilContainer, error) {
	// Generate a unique container name using timestamp to avoid conflicts in parallel tests
	uniqueName := fmt.Sprintf("anvil-%d", time.Now().UnixNano())

	req := testcontainers.ContainerRequest{
		Cmd:          []string{"anvil"},
		ExposedPorts: []string{AnvilPort},
		Env:          map[string]string{"ANVIL_IP_ADDR": "0.0.0.0"},
		Image:        AnvilImage,
		Name:         uniqueName,
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("8545/tcp"),
			wait.ForLog("Listening on 0.0.0.0:8545").WithStartupTimeout(30*time.Second),
		),
	}

	// Add host port binding if requested
	if opts.ExposeHostPort {
		hostPort := opts.HostPort
		if hostPort == "" {
			hostPort = "8545"
		}
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"8545/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: hostPort,
					},
				},
			}
		}
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
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	rpcURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	return &AnvilContainer{
		container: container,
		rpcURL:    rpcURL,
	}, nil
}

// RpcURL returns the RPC URL for connecting to the Anvil instance
func (ac *AnvilContainer) RpcURL() string {
	return ac.rpcURL
}

// Terminate stops and removes the container
func (ac *AnvilContainer) Terminate(ctx context.Context) error {
	if ac == nil || ac.container == nil {
		return nil
	}
	if err := ac.container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate Anvil container: %w", err)
	}
	return nil
}
