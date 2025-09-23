package testbed

import (
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestChurnerContainerStartup(t *testing.T) {
	ctx := t.Context()
	logger := test.GetLogger()

	// Create a test network using the newer API
	dockerNetwork, err := network.New(ctx,
		network.WithDriver("bridge"),
		network.WithAttachable(),
	)
	require.NoError(t, err, "Failed to create docker network")
	defer dockerNetwork.Remove(ctx)

	// Start Anvil container first
	anvil, err := NewAnvilContainerWithOptions(ctx, AnvilOptions{
		Network: dockerNetwork,
		Logger:  logger,
	})
	require.NoError(t, err, "Failed to start anvil container")
	defer func() {
		if anvil != nil {
			err := anvil.Terminate(ctx)
			require.NoError(t, err, "Failed to terminate anvil container")
		}
	}()

	// Deploy contracts to Anvil
	deployment, err := DeployContractsToAnvil(anvil.RpcURL(), 1, logger)
	require.NoError(t, err, "Failed to deploy contracts")
	logger.Info("Contracts deployed:")
	logger.Infof("  EigenDADirectory: %s", deployment.EigenDA.EigenDADirectory)
	logger.Infof("  ServiceManager: %s", deployment.EigenDA.ServiceManager)
	logger.Infof("  OperatorStateRetriever: %s", deployment.EigenDA.OperatorStateRetriever)

	// Get Anvil's default key for the churner to use
	anvilDefaultKey, _ := GetAnvilDefaultKeys()

	// Create churner config with required fields
	config := DefaultChurnerConfig()
	config.ChainRPC = "http://anvil:8545"                         // Use the anvil container in the same network
	config.PrivateKey = strings.TrimPrefix(anvilDefaultKey, "0x") // Use Anvil's default account
	config.ServiceManager = deployment.EigenDA.ServiceManager
	config.OperatorStateRetriever = deployment.EigenDA.OperatorStateRetriever
	config.GraphURL = "http://graph:8000" // Required field - use dummy URL for testing (won't actually connect)
	config.StartupTimeout = 30 * time.Second

	// Start the churner container
	churner, err := NewChurnerContainerWithNetwork(ctx, config, dockerNetwork)
	require.NoError(t, err, "Failed to start churner container")
	require.NotNil(t, churner)
	defer func() {
		if churner != nil {
			err := churner.Stop(ctx)
			require.NoError(t, err, "Failed to stop churner container")
		}
	}()

	// Verify the container is running
	state, err := churner.Container.State(ctx)
	require.NoError(t, err)
	require.True(t, state.Running, "Container should be running")

	// Verify we have URLs
	require.NotEmpty(t, churner.URL(), "External URL should not be empty")
	require.NotEmpty(t, churner.InternalURL(), "Internal URL should not be empty")

	// Log the URLs for debugging
	logger.Infof("Churner started successfully:")
	logger.Infof("  External URL: %s", churner.URL())
	logger.Infof("  Internal URL: %s", churner.InternalURL())
	logger.Infof("  Internal IP: %s", churner.InternalIP())
	logger.Infof("  Log path: %s", churner.LogPath())
}
