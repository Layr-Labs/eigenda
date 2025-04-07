package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/clients/memconfig_client"
	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/testutils"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requireDispersalRetrievalEigenDA ... ensure that blob was successfully dispersed/read to/from EigenDA
func requireDispersalRetrievalEigenDA(t *testing.T, cm *metrics.CountMap, mode commitments.CommitmentMode) {
	writeCount, err := cm.Get(string(mode), http.MethodPost)
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(string(mode), http.MethodGet)
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireDispersalRetrievalEigenDACounts ensures that blob were successfully dispersed/read to/from EigenDA, the
// expected number of times
func requireDispersalRetrievalEigenDACounts(
	t *testing.T,
	cm *metrics.CountMap,
	mode commitments.CommitmentMode,
	expectedWriteCount uint64,
	expectedReadCount uint64,
) {
	actualWriteCount, err := cm.Get(string(mode), http.MethodPost)
	require.NoError(t, err)
	require.Equal(t, expectedWriteCount, actualWriteCount)

	actualReadCount, err := cm.Get(string(mode), http.MethodGet)
	require.NoError(t, err)
	require.Equal(t, expectedReadCount, actualReadCount)
}

// requireWriteReadSecondary ... ensure that secondary backend was successfully written/read to/from
func requireWriteReadSecondary(t *testing.T, cm *metrics.CountMap, bt common.BackendType) {
	writeCount, err := cm.Get(http.MethodPut, store.Success, bt.String())
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(http.MethodGet, store.Success, bt.String())
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireStandardClientSetGet ... ensures that std proxy client can disperse and read a blob
func requireStandardClientSetGet(t *testing.T, ts testutils.TestSuite, blob []byte) {
	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, blob)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)

}

// requireOPClientSetGet ... ensures that alt-da client can disperse and read a blob
func requireOPClientSetGet(t *testing.T, ts testutils.TestSuite, blob []byte, precompute bool) {
	daClient := altda.NewDAClient(ts.Address(), false, precompute)

	commit, err := daClient.SetInput(ts.Ctx, blob)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)
}

func TestOptimismClientWithKeccak256CommitmentV1(t *testing.T) {
	testOptimismClientWithKeccak256Commitment(t, false)
}

func TestOptimismClientWithKeccak256CommitmentV2(t *testing.T) {
	testOptimismClientWithKeccak256Commitment(t, true)
}

func testOptimismClientWithKeccak256Commitment(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireOPClientSetGet(t, ts, testutils.RandBytes(100), true)
}

func TestOptimismClientWithGenericCommitmentV1(t *testing.T) {
	testOptimismClientWithGenericCommitment(t, false)
}

func TestOptimismClientWithGenericCommitmentV2(t *testing.T) {
	testOptimismClientWithGenericCommitment(t, true)
}

/*
this test asserts that the data can be posted/read to EigenDA
with a concurrent S3 backend configured
*/
func testOptimismClientWithGenericCommitment(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireOPClientSetGet(t, ts, testutils.RandBytes(100), false)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.OptimismGeneric)
}

func TestProxyClientServerIntegrationV1(t *testing.T) {
	testProxyClientServerIntegration(t, false)
}

func TestProxyClientServerIntegrationV2(t *testing.T) {
	testProxyClientServerIntegration(t, true)
}

// TestProxyClientServerIntegration tests the proxy client and server integration by setting the data as a single byte,
// many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func testProxyClientServerIntegration(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	t.Cleanup(kill)

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)

	t.Run(
		"single byte preimage set data case", func(t *testing.T) {
			t.Parallel()
			testPreimage := []byte{1} // single byte preimage
			t.Log("Setting input data on proxy server...")
			_, err := daClient.SetData(ts.Ctx, testPreimage)
			require.NoError(t, err)
		})

	t.Run(
		"unicode preimage set data case", func(t *testing.T) {
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

	t.Run(
		"empty preimage set data case", func(t *testing.T) {
			t.Parallel()
			testPreimage := []byte("") // Empty preimage
			t.Log("Setting input data on proxy server...")
			_, err := daClient.SetData(ts.Ctx, testPreimage)
			require.NoError(t, err)
		})

	t.Run(
		"get data edge cases", func(t *testing.T) {
			t.Parallel()
			testCert := []byte("")
			_, err := daClient.GetData(ts.Ctx, testCert)
			require.Error(t, err)
			assert.True(
				t, strings.Contains(
					err.Error(),
					"404") && !isNilPtrDerefPanic(err.Error()))

			testCert = []byte{2}
			_, err = daClient.GetData(ts.Ctx, testCert)
			require.Error(t, err)
			assert.True(
				t, strings.Contains(
					err.Error(),
					"400") && !isNilPtrDerefPanic(err.Error()))

			testCert = testutils.RandBytes(10000)
			_, err = daClient.GetData(ts.Ctx, testCert)
			require.Error(t, err)
			assert.True(t, strings.Contains(err.Error(), "400") && !isNilPtrDerefPanic(err.Error()))
		})
}

func TestProxyClientV1(t *testing.T) {
	testProxyClient(t, false)
}

func TestProxyClientV2(t *testing.T) {
	testProxyClient(t, true)
}

func testProxyClient(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)

	testPreimage := testutils.RandBytes(100)

	t.Log("Setting input data on proxy server...")
	daCommitment, err := daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.Ctx, daCommitment)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestProxyClientWriteReadV1(t *testing.T) {
	testProxyClientWriteRead(t, false)
}

func TestProxyClientWriteReadV2(t *testing.T) {
	testProxyClientWriteRead(t, true)
}

func testProxyClientWriteRead(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.MemstoreBackend, disperseToV2)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(100))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyCachingV1(t *testing.T) {
	testProxyCaching(t, false)
}

func TestProxyCachingV2(t *testing.T) {
	testProxyCaching(t, true)
}

/*
Ensure that proxy is able to write/read from a cache backend when enabled
*/
func testProxyCaching(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	testCfg.UseS3Caching = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyCachingWithRedisV1(t *testing.T) {
	testProxyCachingWithRedis(t, false)
}

func TestProxyCachingWithRedisV2(t *testing.T) {
	testProxyCachingWithRedis(t, true)
}

func testProxyCachingWithRedis(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	testCfg.UseRedisCaching = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.RedisBackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyReadFallbackV1(t *testing.T) {
	testProxyReadFallback(t, false)
}

func TestProxyReadFallbackV2(t *testing.T) {
	testProxyReadFallback(t, true)
}

/*
Ensure that fallback location is read from when EigenDA blob is not available.
This is done by setting the memstore expiration time to 1ms and waiting for the blob to expire
before attempting to read it.
*/
func testProxyReadFallback(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip(`test only runs with memstore, since fallback relies on blob fetch failing, and it won't fail
						against actual eigen DA`)
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	testCfg.UseS3Fallback = true
	// ensure that blob memstore eviction times result in near immediate activation
	testCfg.Expiration = time.Millisecond * 1

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)
	expectedBlob := testutils.RandBytes(1_000_000)
	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, expectedBlob)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	t.Log("Getting input data from proxy server...")
	actualBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, actualBlob)

	requireStandardClientSetGet(t, ts, testutils.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestProxyMemConfigClientCanGetAndPatchV1(t *testing.T) {
	testProxyMemConfigClientCanGetAndPatch(t, false)
}

func TestProxyMemConfigClientCanGetAndPatchV2(t *testing.T) {
	testProxyMemConfigClientCanGetAndPatch(t, true)
}

func testProxyMemConfigClientCanGetAndPatch(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	useMemstore := testutils.GetBackend() == testutils.MemstoreBackend
	if !useMemstore {
		t.Skip("test can't be run against holesky since read failure case can't be manually triggered")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	memClient := memconfig_client.New(
		&memconfig_client.Config{
			URL: "http://" + ts.Server.Endpoint(),
		})

	// 1 - ensure cfg can be read from memconfig handlers
	cfg, err := memClient.GetConfig(ts.Ctx)
	require.NoError(t, err)

	// 2 - update PutLatency field && ensure that newly fetched config reflects change
	expectedChange := time.Second * 420
	cfg.PutLatency = expectedChange

	cfg, err = memClient.UpdateConfig(ts.Ctx, cfg)
	require.NoError(t, err)

	require.Equal(t, cfg.PutLatency, expectedChange)

	// 3 - get cfg again to verify that memconfig state update is now reflected on server
	cfg, err = memClient.GetConfig(ts.Ctx)

	require.NoError(t, err)
	require.Equal(t, cfg.PutLatency, expectedChange)
}

// TestInterleavedVersions alternately disperses payloads to v1 and v2, and then retrieves them.
func TestInterleavedVersions(t *testing.T) {
	t.Parallel()
	testRandom := random.NewTestRandom()

	testCfg := testutils.NewTestConfig(
		testutils.GetBackend(),
		false,
		common.V1EigenDABackend,
		common.V2EigenDABackend)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	testSuite, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	client := standard_client.New(
		&standard_client.Config{
			URL: testSuite.Address(),
		})

	// disperse a payload to v1
	payload1a := testRandom.Bytes(1000)
	cert1a, err := client.SetData(testSuite.Ctx, payload1a)
	require.NoError(t, err)

	// disperse a payload to v2
	testSuite.Server.SetDisperseToV2(true)
	payload2a := testRandom.Bytes(1000)
	cert2a, err := client.SetData(testSuite.Ctx, payload2a)
	require.NoError(t, err)

	// disperse another payload to v1
	testSuite.Server.SetDisperseToV2(false)
	payload1b := testRandom.Bytes(1000)
	cert1b, err := client.SetData(testSuite.Ctx, payload1b)
	require.NoError(t, err)

	// disperse another payload to v2
	testSuite.Server.SetDisperseToV2(true)
	payload2b := testRandom.Bytes(1000)
	cert2b, err := client.SetData(testSuite.Ctx, payload2b)
	require.NoError(t, err)

	// fetch in reverse order, because why not
	fetchedPayload2b, err := client.GetData(testSuite.Ctx, cert2b)
	require.NoError(t, err)
	fetchedPayload1b, err := client.GetData(testSuite.Ctx, cert1b)
	require.NoError(t, err)
	fetchedPayload2a, err := client.GetData(testSuite.Ctx, cert2a)
	require.NoError(t, err)
	fetchedPayload1a, err := client.GetData(testSuite.Ctx, cert1a)
	require.NoError(t, err)

	require.Equal(t, payload1a, fetchedPayload1a)
	require.Equal(t, payload2a, fetchedPayload2a)
	require.Equal(t, payload1b, fetchedPayload1b)
	require.Equal(t, payload2b, fetchedPayload2b)

	requireStandardClientSetGet(t, testSuite, testRandom.Bytes(100))
	requireDispersalRetrievalEigenDA(t, testSuite.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}

func TestMaxBlobSizeV1(t *testing.T) {
	if testutils.GetBackend() == testutils.PreprodBackend {
		t.Skip("Preprod for v1 has a stricter blob size than normal.")
	}

	testMaxBlobSize(t, false)
}

func TestMaxBlobSizeV2(t *testing.T) {
	testMaxBlobSize(t, true)
}

func testMaxBlobSize(t *testing.T, disperseToV2 bool) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), disperseToV2)
	testCfg.MaxBlobLength = "16mib"
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	// the payload has things added to it during encoding, so it has a slightly lower limit than max blob size
	maxPayloadSize, err := codec.GetMaxPermissiblePayloadLength(
		uint32(tsConfig.EigenDAConfig.ClientConfigV2.MaxBlobSizeBytes / 32))
	require.NoError(t, err)

	requireStandardClientSetGet(t, ts, testutils.RandBytes(int(maxPayloadSize)))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.Standard)
}
