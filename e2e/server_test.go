package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/stretchr/testify/require"
)

func useMemory() bool {
	return !runTestnetIntegrationTests
}

func TestOptimismClientWithKeccak256Commitment(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseKeccak256ModeS3 = true

	ts, kill := e2e.CreateTestSuite(t, testCfg)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, true)

	testPreimage := []byte(e2e.RandString(100))

	commit, err := daClient.SetInput(ts.Ctx, testPreimage)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

/*
this test asserts that the data can be posted/read to EigenDA
with a concurrent S3 backend configured
*/
func TestOptimismClientWithGenericCommitment(t *testing.T) {

	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, e2e.TestConfig(useMemory()))
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

	ts, kill := e2e.CreateTestSuite(t, e2e.TestConfig(useMemory()))
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

func TestProxyServerWithLargeBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, e2e.TestConfig(useMemory()))
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  16MB blob
	testPreimage := []byte(e2e.RandString(16_000_000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyServerWithOversizedBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	ts, kill := e2e.CreateTestSuite(t, e2e.TestConfig(useMemory()))
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

/*
Ensure that proxy is able to write/read from a cache backend when enabled
*/
func TestProxyServerCaching(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseS3Caching = true

	ts, kill := e2e.CreateTestSuite(t, testCfg)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  1mb blob
	testPreimage := []byte(e2e.RandString(1_0000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NotEmpty(t, blobInfo)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)

	// ensure that read was from cache
	s3Stats := ts.Server.GetS3Stats()
	require.Equal(t, 1, s3Stats.Reads)
	require.Equal(t, 1, s3Stats.Entries)

	if useMemory() { // ensure that eigenda was not read from
		memStats := ts.Server.GetEigenDAStats()
		require.Equal(t, 0, memStats.Reads)
		require.Equal(t, 1, memStats.Entries)
	}
}

/*
	Ensure that fallback location is read from when EigenDA blob is not available.
	This is done by setting the memstore expiration time to 1ms and waiting for the blob to expire
	before attempting to read it.
*/

func TestProxyServerReadFallback(t *testing.T) {
	// test can't be ran against holesky since read failure case can't be manually triggered
	if !runIntegrationTests || runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseS3Fallback = true
	testCfg.Expiration = time.Millisecond * 1

	ts, kill := e2e.CreateTestSuite(t, testCfg)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  1mb blob
	testPreimage := []byte(e2e.RandString(1_0000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NotEmpty(t, blobInfo)
	require.NoError(t, err)

	time.Sleep(time.Second * 1)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)

	// ensure that read was from fallback target location (i.e, S3 for this test)
	s3Stats := ts.Server.GetS3Stats()
	require.Equal(t, 1, s3Stats.Reads)
	require.Equal(t, 1, s3Stats.Entries)

	if useMemory() { // ensure that an eigenda read was attempted with zero data available
		memStats := ts.Server.GetEigenDAStats()
		require.Equal(t, 1, memStats.Reads)
		require.Equal(t, 0, memStats.Entries)
	}
}
