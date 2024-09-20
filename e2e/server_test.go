package e2e_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"

	"github.com/Layr-Labs/eigenda-proxy/e2e"
	"github.com/Layr-Labs/eigenda-proxy/store"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func useMemory() bool {
	return !runTestnetIntegrationTests
}

func isNilPtrDerefPanic(err string) bool {
	return strings.Contains(err, "panic") && strings.Contains(err, "SIGSEGV") &&
		strings.Contains(err, "nil pointer dereference")
}

// TestOpClientKeccak256MalformedInputs tests the NewDAClient from op_plasma by setting and getting against []byte("")
// preimage. It sets the precompute option to false on the NewDAClient.
func TestOpClientKeccak256MalformedInputs(t *testing.T) {
	if !runIntegrationTests || runTestnetIntegrationTests {
		t.Skip("Skipping test as TESTNET env set or INTEGRATION var not set")
	}

	t.Parallel()
	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseKeccak256ModeS3 = true
	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
	defer kill()

	// nil commitment. Should return an error but currently is not. This needs to be fixed by OP
	// Ref: https://github.com/ethereum-optimism/optimism/issues/11987
	// daClient := op_plasma.NewDAClient(ts.Address(), false, true)
	// t.Run("nil commitment case", func(t *testing.T) {
	//	var commit op_plasma.CommitmentData
	//	_, err := daClient.GetInput(ts.Ctx, commit)
	//	require.Error(t, err)
	//	assert.True(t, !isPanic(err.Error()))
	// })

	daClientPcFalse := op_plasma.NewDAClient(ts.Address(), false, false)

	t.Run("input bad data to SetInput & GetInput", func(t *testing.T) {
		testPreimage := []byte("") // Empty preimage
		_, err := daClientPcFalse.SetInput(ts.Ctx, testPreimage)
		require.Error(t, err)

		// should fail with proper error message as is now, and cannot contain panics or nils
		assert.True(t, strings.Contains(err.Error(), "invalid input") && !isNilPtrDerefPanic(err.Error()))

		// The below test panics silently.
		input := op_plasma.NewGenericCommitment([]byte(""))
		_, err = daClientPcFalse.GetInput(ts.Ctx, input)
		require.Error(t, err)

		// Should not fail on slice bounds out of range. This needs to be fixed by OP.
		// Refer to issue: https://github.com/ethereum-optimism/optimism/issues/11987
		// assert.False(t, strings.Contains(err.Error(), ": EOF") && !isPanic(err.Error()))
	})

}

func TestOptimismClientWithKeccak256Commitment(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, true)

	t.Run("normal case", func(t *testing.T) {
		testPreimage := []byte(e2e.RandString(100))

		commit, err := daClient.SetInput(ts.Ctx, testPreimage)
		require.NoError(t, err)

		preimage, err := daClient.GetInput(ts.Ctx, commit)
		require.NoError(t, err)
		require.Equal(t, testPreimage, preimage)
	})
}

func TestKeccak256CommitmentRequestErrorsWhenS3NotSet(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	tsConfig.S3Config.Endpoint = ""
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
	defer kill()

	daClient := op_plasma.NewDAClient(ts.Address(), false, true)

	testPreimage := []byte(e2e.RandString(100))

	_, err := daClient.SetInput(ts.Ctx, testPreimage)
	// TODO: the server currently returns an internal server error. Should it return a 400 instead?
	require.Error(t, err)
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

	tsConfig := e2e.TestSuiteConfig(t, e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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

// TestProxyClientServerIntegration tests the proxy client and server integration by setting the data as a single byte,
// many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func TestProxyClientServerIntegration(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	tsConfig := e2e.TestSuiteConfig(t, e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)

	t.Run("single byte preimage set data case", func(t *testing.T) {
		testPreimage := []byte{1} // single byte preimage
		t.Log("Setting input data on proxy server...")
		_, err := daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)
	})

	t.Run("unicode preimage set data case", func(t *testing.T) {
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
		testPreimage := []byte("") // Empty preimage
		t.Log("Setting input data on proxy server...")
		_, err := daClient.SetData(ts.Ctx, testPreimage)
		require.NoError(t, err)
	})

	t.Run("get data edge cases", func(t *testing.T) {
		testCert := []byte("")
		_, err := daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(),
			"commitment is too short") && !isNilPtrDerefPanic(err.Error()))

		testCert = []byte{1}
		_, err = daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(),
			"commitment is too short") && !isNilPtrDerefPanic(err.Error()))

		testCert = []byte(e2e.RandString(10000))
		_, err = daClient.GetData(ts.Ctx, testCert)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(),
			"failed to decode DA cert to RLP format: rlp: expected input list for verify.Certificate") &&
			!isNilPtrDerefPanic(err.Error()))
	})

}

func TestProxyClient(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	tsConfig := e2e.TestSuiteConfig(t, e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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

	tsConfig := e2e.TestSuiteConfig(t, e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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

	tsConfig := e2e.TestSuiteConfig(t, e2e.TestConfig(useMemory()))
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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
	require.Contains(t, err.Error(), fmt.Sprint(http.StatusBadRequest))

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

	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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

func TestProxyServerCachingWithRedis(t *testing.T) {
	if !runIntegrationTests && !runTestnetIntegrationTests {
		t.Skip("Skipping test as INTEGRATION or TESTNET env var not set")
	}

	t.Parallel()

	testCfg := e2e.TestConfig(useMemory())
	testCfg.UseRedisCaching = true

	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
	defer kill()

	cfg := &client.Config{
		URL: ts.Address(),
	}
	daClient := client.New(cfg)
	//  10 kb blob
	testPreimage := []byte(e2e.RandString(10_000))

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NotEmpty(t, blobInfo)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)

	// ensure that read was from cache
	redStats, err := ts.Server.GetStoreStats(store.Redis)
	require.NoError(t, err)

	require.Equal(t, 1, redStats.Reads)
	require.Equal(t, 1, redStats.Entries)

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

	tsConfig := e2e.TestSuiteConfig(t, testCfg)
	ts, kill := e2e.CreateTestSuite(t, tsConfig)
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
