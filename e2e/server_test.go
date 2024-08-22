package e2e_test

import (
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/stretchr/testify/require"
)

func useMemory() bool {
	return !runTestnetIntegrationTests
}

func TestOptimismClientWithS3Backend(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), true)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, true)

	testPreimage := []byte(e2e.RandString(100))

	commit, err := daClient.SetInput(ts.Ctx, testPreimage)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestOptimismClientWithEigenDABackend(t *testing.T) {
	// this test asserts that the data can be posted/read to EigenDA with a concurrent S3 backend configured

	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), true)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, false)

	testPreimage := []byte(e2e.RandString(100))

	t.Log("Setting input data on proxy server...")
	commit, err := daClient.SetInput(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClient(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	testPreimage := []byte(e2e.RandString(100))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClientWithLargeBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  2MB blob
	testPreimage := []byte(e2e.RandString(4_000_000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClientWithOversizedBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, useMemory(), false)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  17MB blob
	testPreimage := []byte(e2e.RandString(17_000_0000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.Empty(t, blobInfo)
	require.Error(t, err)

	oversizedError := false
	if strings.Contains(err.Error(), "blob is larger than max blob size") {
		oversizedError = true
	}

	if strings.Contains(err.Error(), "blob size cannot exceed") {
		oversizedError = true
	}

	require.True(t, oversizedError)

}
