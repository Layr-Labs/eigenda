package e2e_test

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda-proxy/e2e"
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

	tsConfig := e2e.TestSuiteConfig(testCfg)
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()
	requireOPClientSetGet(t, ts, e2e.RandBytes(100), true)
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

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	requireOPClientSetGet(t, ts, e2e.RandBytes(100), false)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.OptimismGeneric)
}

func TestProxyClientWriteRead(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, e2e.RandBytes(100))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyWithMaximumSizedBlob(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, e2e.RandBytes(16_000_000))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

/*
Ensure that proxy is able to write/read from a cache backend when enabled
*/
func TestProxyCaching(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseS3Caching = true

	tsConfig := e2e.TestSuiteConfig(testCfg)
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, e2e.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyCachingWithRedis(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseRedisCaching = true

	tsConfig := e2e.TestSuiteConfig(testCfg)
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, e2e.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.RedisBackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

/*
	Ensure that fallback location is read from when EigenDA blob is not available.
	This is done by setting the memstore expiration time to 1ms and waiting for the blob to expire
	before attempting to read it.
*/

func TestProxyReadFallback(t *testing.T) {
	// test can't be ran against holesky since read failure case can't be manually triggered
	if !runIntegrationTests || runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION env var not set")
	}

	t.Parallel()

	// setup server with S3 as a fallback option
	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseS3Fallback = true
	// ensure that blob memstore eviction times result in near immediate activation
	testCfg.Expiration = time.Millisecond * 1

	tsConfig := e2e.TestSuiteConfig(testCfg)
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	expectedBlob := e2e.RandBytes(1_000_000)
	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, expectedBlob)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	t.Log("Getting input data from proxy server...")
	actualBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, actualBlob)

	requireStandardClientSetGet(t, ts, e2e.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}
