// Configuration tests are to test specific configuration/initialization scenarios,
// that aren't specific to any particular API. Tests that are specific to an API
// (op, rest, arb) should go in their respective test files instead.
package e2e

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/stretchr/testify/require"
)

// Tests that a proxy started with V2 EigenDA backend and without a signer private key
// is in read-only mode, meaning that POST routes return 500 errors, while GET routes work as expected.
// TODO(samlaf): Feels a bit dumb to run a simple test like this in e2e framework,
// since it takes 9 seconds, requires an actual eth-rpc (adds ci flakiness), etc.
// We don't really have an alternative however given that the read-only feature is only
// implemented inside the EigenDAV2 store.
func TestProxyV2ReadOnlyMode(t *testing.T) {
	if testutils.GetBackend() == testutils.MemstoreBackend {
		t.Skip("Don't run for memstore backend, since read-only mode is only implemented for eigenda v2 backend")
	}

	// We test against sepolia backend in order to test the client creation code (which reads the signer private key).
	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	tsConfig.SecretConfig.SignerPaymentKey = "" // ensure no signer key is set
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()
	testBlob := []byte("hello world")

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	t.Log("Setting input data on proxy server...")
	_, err := daClient.SetData(ts.Ctx, testBlob)
	require.Error(t, err)
	// expect 500 in read-only mode. Routes are turned off but we don't have an explicit "read-only" mode config,
	// so error return only says "PUT routes are disabled, did you provide a signer private key?".
	require.ErrorContains(t, err, "500")
	require.ErrorContains(t, err, "PUT routes are disabled")

	// We also check that the Get routes are still working.
	// We pass a fake bogus cert which doesn't even parse, so expect a 418 error (indicating to discard cert).
	fakeStdCommitment := []byte{1, 2, 3, 4, 5, 6}
	_, err = daClient.GetData(ts.Ctx, fakeStdCommitment)
	require.Error(t, err)
	require.ErrorContains(t, err, "418")
}
