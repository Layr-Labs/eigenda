package testbed_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/stretchr/testify/require"
)

// TestDeployWithAnvilContainer demonstrates deploying contracts using Docker-based Anvil
func TestDeployWithAnvilContainer(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()

	// Start Anvil container
	anvil, err := testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
	})
	require.NoError(t, err)
	defer anvil.Terminate(ctx)

	// Deploy contracts to Anvil with 4 operators
	result, err := testbed.DeployContractsToAnvil(anvil.RpcURL(), 4, logger)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify all contract addresses were deployed
	require.NotEmpty(t, result.EigenDA.EigenDADirectory)
	require.NotEmpty(t, result.EigenDA.ServiceManager)
	require.NotEmpty(t, result.EigenDA.OperatorStateRetriever)
	require.NotEmpty(t, result.EigenDA.BlsApkRegistry)
	require.NotEmpty(t, result.EigenDA.RegistryCoordinator)
	require.NotEmpty(t, result.EigenDA.CertVerifierLegacy)
	require.NotEmpty(t, result.EigenDA.CertVerifier)
	require.NotEmpty(t, result.EigenDA.CertVerifierRouter)

	// Verify V1 Cert Verifier address was deployed
	require.NotEmpty(t, result.EigenDAV1CertVerifier)
}
