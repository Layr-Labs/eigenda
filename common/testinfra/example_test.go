package testinfra_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testinfra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInfrastructureMinimal demonstrates basic usage with Anvil and LocalStack only
func TestInfrastructureMinimal(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start minimal infrastructure (Anvil + LocalStack)
	manager, result, err := testinfra.StartMinimal(ctx)
	require.NoError(t, err)
	defer manager.Stop(ctx)

	// Verify Anvil is running
	assert.NotEmpty(t, result.AnvilRPC)
	assert.Equal(t, 31337, result.AnvilChainID)

	// Verify LocalStack is running
	assert.NotEmpty(t, result.LocalStackURL)

	// Test that we can get private keys from Anvil
	anvil := manager.GetAnvil()
	require.NotNil(t, anvil)

	privateKey, err := anvil.GetPrivateKey(0)
	assert.NoError(t, err)
	assert.NotEmpty(t, privateKey)

	// Test LocalStack AWS config
	localstack := manager.GetLocalStack()
	require.NotNil(t, localstack)

	awsConfig := localstack.GetAWSConfig()
	assert.Equal(t, "test", awsConfig["AWS_ACCESS_KEY_ID"])
	assert.Equal(t, result.LocalStackURL, awsConfig["AWS_ENDPOINT_URL"])
}

// TestCustomConfiguration demonstrates custom configuration
func TestCustomConfiguration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Create custom config
	config := testinfra.DefaultConfig()
	config.Anvil.ChainID = 1337
	config.Anvil.Accounts = 5
	config.LocalStack.Debug = true
	config.GraphNode.Enabled = false // Keep disabled for this test

	manager, result, err := testinfra.StartCustom(ctx, config)
	require.NoError(t, err)
	defer manager.Stop(ctx)

	// Verify custom configuration was applied
	assert.Equal(t, 1337, result.AnvilChainID)

	anvil := manager.GetAnvil()
	require.NotNil(t, anvil)
	assert.Equal(t, 5, anvil.Accounts())
}

// TestInfrastructureFull demonstrates full infrastructure setup (disabled by default due to complexity)
func TestInfrastructureFull(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start full infrastructure (includes Graph Node)
	manager, result, err := testinfra.StartFull(ctx)
	require.NoError(t, err)
	defer manager.Stop(ctx)

	// Verify all components are running
	assert.NotEmpty(t, result.AnvilRPC)
	assert.NotEmpty(t, result.LocalStackURL)
	assert.NotEmpty(t, result.GraphNodeURL)

	graphNode := manager.GetGraphNode()
	require.NotNil(t, graphNode)
	assert.NotEmpty(t, graphNode.HTTPURL())
}

// TestInfrastructureLifecycle demonstrates proper cleanup
func TestInfrastructureLifecycle(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start infrastructure
	manager, result, err := testinfra.StartMinimal(ctx)
	require.NoError(t, err)

	// Verify it's running
	assert.NotEmpty(t, result.AnvilRPC)

	// Explicitly stop infrastructure
	err = manager.Stop(ctx)
	assert.NoError(t, err)

	// Attempting to use stopped infrastructure should fail gracefully
	anvil := manager.GetAnvil()
	assert.Nil(t, anvil) // Should be nil after cleanup
}
