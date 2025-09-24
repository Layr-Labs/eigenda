package testbed

import (
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestChurnerGoroutineStartup(t *testing.T) {
	ctx := t.Context()
	logger := test.GetLogger()

	// Create a test network using the newer API
	dockerNetwork, err := network.New(ctx,
		network.WithDriver("bridge"),
		network.WithAttachable(),
	)
	require.NoError(t, err, "Failed to create docker network")
	defer func() { _ = dockerNetwork.Remove(ctx) }()

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
	config.ChainRPC = anvil.RpcURL()                              // Use the anvil RPC URL directly
	config.PrivateKey = strings.TrimPrefix(anvilDefaultKey, "0x") // Use Anvil's default account
	config.ServiceManager = deployment.EigenDA.ServiceManager
	config.OperatorStateRetriever = deployment.EigenDA.OperatorStateRetriever
	config.GraphURL = "http://graph:8000" // Required field - use dummy URL for testing (won't actually connect)
	config.StartupTimeout = 30 * time.Second

	// Start the churner goroutine
	churner, err := StartChurnerGoroutine(config, logger)
	require.NoError(t, err, "Failed to start churner goroutine")
	require.NotNil(t, churner)
	defer func() {
		if churner != nil {
			churner.Stop(ctx)
		}
	}()

	// Verify we have a URL
	require.NotEmpty(t, churner.URL(), "URL should not be empty")

	// Log the URL for debugging
	logger.Infof("Churner started successfully:")
	logger.Infof("  URL: %s", churner.URL())
}
