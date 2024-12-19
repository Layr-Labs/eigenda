package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/e2e"
	"github.com/stretchr/testify/assert"
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

// TestProxyClientServerIntegration tests the proxy client and server integration by setting the data as a single byte,
// many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func TestProxyClientServerIntegration(t *testing.T) {
	t.Parallel()

	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)
	t.Cleanup(kill)

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	t.Run("single byte preimage set data case", func(t *testing.T) {
		t.Parallel()
		testPreimage := []byte{1} // single byte preimage
		t.Log("Setting input data on proxy server...")
		_, err := daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)
	})

	t.Run("unicode preimage set data case", func(t *testing.T) {
		t.Parallel()
		testPreimage := []byte("§§©ˆªªˆ˙√ç®∂§∞¶§ƒ¥√¨¥√¨¥ƒƒ©˙˜ø˜˜˜∫˙∫¥∫√†®®√ç¨ˆ¨˙ï") // many unicode characters
		t.Log("Setting input data on proxy server...")
		_, err := daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)

		testPreimage = []byte("§") // single unicode character
		t.Log("Setting input data on proxy server...")
		_, err = daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)

	})

	t.Run("empty preimage set data case", func(t *testing.T) {
		t.Parallel()
		testPreimage := []byte("") // Empty preimage
		t.Log("Setting input data on proxy server...")
		_, err := daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)
	})

	t.Run("get data edge cases", func(t *testing.T) {
		t.Parallel()
		testCert := []byte("")
		_, err := daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(),
			"404") && !isNilPtrDerefPanic(err.Error()))

		testCert = []byte{1}
		_, err = daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(),
			"400") && !isNilPtrDerefPanic(err.Error()))

		testCert = e2e.RandBytes(10000)
		_, err = daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "400") && !isNilPtrDerefPanic(err.Error()))
	})

}

func TestProxyClient(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	tsConfig := e2e.TestSuiteConfig(e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	testPreimage := e2e.RandBytes(100)

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
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
