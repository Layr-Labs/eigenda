package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/proxy/clients/memconfig_client"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	enabled_apis "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProxyAPIsEnabledRestALTDA tests to ensure that the enabled APIs expression is
// is getting respected by the REST ALTDA Server when wiring up a proxy application instance
// with just `op-generic` mode enabled.
func TestProxyAPIsEnabledRestALTDA(t *testing.T) {
	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip(`test only runs with memstore, since code paths being asserted upon aren't 
				network specific. running this in multiple envs would be unnecessary and provide
				no further guarantees.`)
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)
	testCfg.EnabledRestAPIs = &enabled_apis.RestApisEnabled{
		OpGenericCommitment: true,
	}
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()
	testBlob := []byte("hello world")

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg) // standard commitment mode (should fail given disabled)

	t.Log("Setting input data on proxy server...")
	_, err := daClient.SetData(ts.Ctx, testBlob)
	require.Error(t, err)
	require.ErrorContains(t, err, "403")

	opGenericClient := altda.NewDAClient(ts.RestAddress(),
		false, false) // now op-generic mode (should work e2e given enabled)

	daCommit, err := opGenericClient.SetInput(ts.Ctx, testBlob)
	require.NoError(t, err)

	preimage, err := opGenericClient.GetInput(ts.Ctx, daCommit, 0)
	require.NoError(t, err)
	require.Equal(t, testBlob, preimage)
}

func TestProxyClientWriteReadV1(t *testing.T) {
	testProxyClientWriteRead(t, common.V1EigenDABackend)
}

func TestProxyClientWriteReadV2(t *testing.T) {
	testProxyClientWriteRead(t, common.V2EigenDABackend)
}

// TestProxyClientWriteRead tests that the proxy client can write and read data to the proxy server.
//
// This is the "basic" proxy test: "is proxy working?"
func testProxyClientWriteRead(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(100))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestOptimismClientWithKeccak256CommitmentV1(t *testing.T) {
	testOptimismClientWithKeccak256Commitment(t, common.V1EigenDABackend)
}

func TestOptimismClientWithKeccak256CommitmentV2(t *testing.T) {
	testOptimismClientWithKeccak256Commitment(t, common.V2EigenDABackend)
}

func testOptimismClientWithKeccak256Commitment(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireOPClientSetGet(t, ts, testutils.RandBytes(100), true)
}

func TestOptimismClientWithGenericCommitmentV1(t *testing.T) {
	testOptimismClientWithGenericCommitment(t, common.V1EigenDABackend)
}

func TestOptimismClientWithGenericCommitmentV2(t *testing.T) {
	testOptimismClientWithGenericCommitment(t, common.V2EigenDABackend)
}

/*
this test asserts that the data can be posted/read to EigenDA
with a concurrent S3 backend configured
*/
func testOptimismClientWithGenericCommitment(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireOPClientSetGet(t, ts, testutils.RandBytes(100), false)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.OptimismGenericCommitmentMode)
}

func TestProxyClientServerIntegrationV1(t *testing.T) {
	testProxyClientServerIntegration(t, common.V1EigenDABackend)
}

func TestProxyClientServerIntegrationV2(t *testing.T) {
	testProxyClientServerIntegration(t, common.V2EigenDABackend)
}

// TestProxyClientServerIntegration tests the proxy client and server integration by setting the data as a single byte,
// many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func testProxyClientServerIntegration(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	t.Cleanup(kill)

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
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

			testCert = []byte{3}
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

func TestProxyCachingV1(t *testing.T) {
	testProxyCaching(t, common.V1EigenDABackend)
}

func TestProxyCachingV2(t *testing.T) {
	testProxyCaching(t, common.V2EigenDABackend)
}

/*
Ensure that proxy is able to write/read from a cache backend when enabled
*/
func testProxyCaching(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseS3Caching = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(1_000_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestProxyReadFallbackV1(t *testing.T) {
	testProxyReadFallback(t, common.V1EigenDABackend)
}

func TestProxyReadFallbackV2(t *testing.T) {
	testProxyReadFallback(t, common.V2EigenDABackend)
}

/*
Ensure that fallback location is read from when EigenDA blob is not available.
This is done by setting the memstore expiration time to 1ms and waiting for the blob to expire
before attempting to read it.
*/
func testProxyReadFallback(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip(`test only runs with memstore, since fallback relies on blob fetch failing, and it won't fail
						against actual eigen DA`)
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseS3Fallback = true
	// ensure that blob memstore eviction times result in near immediate activation
	testCfg.Expiration = time.Millisecond * 1

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
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
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestProxyWriteCacheOnMissV1(t *testing.T) {
	testProxyWriteCacheOnMiss(t, common.V1EigenDABackend)
}

func TestProxyWriteCacheOnMissV2(t *testing.T) {
	testProxyWriteCacheOnMiss(t, common.V2EigenDABackend)
}

func testProxyWriteCacheOnMiss(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseS3Caching = true
	testCfg.WriteOnCacheMiss = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)
	expectedBlob := testutils.RandBytes(1_000_000)
	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, expectedBlob)
	require.NoError(t, err)

	_, err = daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)

	exists, err := testutils.ExistsBlobInfotInBucket(tsConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.True(t, exists)

	t.Log("Erase blob from the cache...")
	err = testutils.RemoveBlobInfoFromBucket(tsConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	exists, err = testutils.ExistsBlobInfotInBucket(tsConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.False(t, exists)

	// Blob created in disperser, removed from S3
	t.Log("Getting input data from proxy server...")
	actualBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, actualBlob)

	exists, err = testutils.ExistsBlobInfotInBucket(tsConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.True(t, exists)
}

// TestErrorOnSecondaryInsertFailureFlagOnV1 verifies that when the flag is ON,
// secondary storage write failures cause the PUT to return HTTP 500.
func TestErrorOnSecondaryInsertFailureFlagOnV1(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOn(t, common.V1EigenDABackend)
}

func TestErrorOnSecondaryInsertFailureFlagOnV2(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOn(t, common.V2EigenDABackend)
}

func testErrorOnSecondaryInsertFailureFlagOn(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip("test only runs with memstore backend")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	// Use S3 as fallback with invalid credentials to simulate S3 failure
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag

	// Ensure async writes are disabled (required for flag to work)
	testCfg.WriteThreadCount = 0

	// Create a test suite with invalid S3 config to force secondary write failures
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	// Override S3 config with invalid credentials to force write failures
	tsConfig.StoreBuilderConfig.S3Config = s3.Config{
		Bucket:          "invalid-bucket-name",
		Endpoint:        "invalid-endpoint:9000",
		AccessKeyID:     "invalid-key",
		AccessKeySecret: "invalid-secret",
		EnableTLS:       false,
		CredentialType:  s3.CredentialTypeStatic,
	}

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	testBlob := testutils.RandBytes(100)

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	// PUT should fail because S3 write fails and flag is ON
	t.Log("Setting data - should fail due to S3 failure with flag enabled")
	_, err := daClient.SetData(ts.Ctx, testBlob)
	require.Error(t, err, "PUT should fail when error-on-secondary-insert-failure=true and S3 fails")

	// Error should indicate it's a server error (5xx)
	require.Contains(t, err.Error(), "500", "Expected HTTP 500 error")
}

// TestErrorOnSecondaryInsertFailureFlagOffPartialFailureV1 verifies that when the flag is OFF (default),
// partial secondary storage failures are tolerated - PUT succeeds if at least one backend succeeds.
func TestErrorOnSecondaryInsertFailureFlagOffPartialFailureV1(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOffPartialFailure(t, common.V1EigenDABackend)
}

func TestErrorOnSecondaryInsertFailureFlagOffPartialFailureV2(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOffPartialFailure(t, common.V2EigenDABackend)
}

func testErrorOnSecondaryInsertFailureFlagOffPartialFailure(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip("test only runs with memstore backend")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	// Use both cache and fallback - cache will fail, fallback will succeed
	testCfg.UseS3Caching = true
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = false // default: OFF
	testCfg.WriteThreadCount = 0

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	// Override with invalid S3 config to force all secondary write failures
	tsConfig.StoreBuilderConfig.S3Config = s3.Config{
		Bucket:          "invalid-bucket-name",
		Endpoint:        "invalid-endpoint:9000",
		AccessKeyID:     "invalid-key",
		AccessKeySecret: "invalid-secret",
		EnableTLS:       false,
		CredentialType:  s3.CredentialTypeStatic,
	}

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	testBlob := testutils.RandBytes(100)
	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	// With flag OFF, secondary failures are logged but not returned as errors
	// PUT should succeed because primary storage (EigenDA) succeeds
	t.Log("Setting data - should succeed because flag OFF means secondary failures are tolerated")
	blobInfo, err := daClient.SetData(ts.Ctx, testBlob)
	require.NoError(t, err, "PUT should succeed when flag OFF even if all secondaries fail")

	// Verify data can be read back from primary storage
	retrievedBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testBlob, retrievedBlob)
}

// TestErrorOnSecondaryInsertFailureFlagOnSuccessV1 verifies that when the flag is ON
// and all secondary writes succeed, PUT succeeds normally (happy path).
func TestErrorOnSecondaryInsertFailureFlagOnSuccessV1(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOnSuccess(t, common.V1EigenDABackend)
}

func TestErrorOnSecondaryInsertFailureFlagOnSuccessV2(t *testing.T) {
	testErrorOnSecondaryInsertFailureFlagOnSuccess(t, common.V2EigenDABackend)
}

func testErrorOnSecondaryInsertFailureFlagOnSuccess(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	if testutils.GetBackend() != testutils.MemstoreBackend {
		t.Skip("test only runs with memstore backend")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag
	testCfg.WriteThreadCount = 0

	// Use valid S3 config
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	testBlob := testutils.RandBytes(100)
	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	// PUT should succeed because all backends (primary + S3) work
	t.Log("Setting data - should succeed with valid S3 config and flag ON")
	blobInfo, err := daClient.SetData(ts.Ctx, testBlob)
	require.NoError(t, err, "PUT should succeed when flag ON and all writes succeed")

	// Verify data can be read back
	t.Log("Getting data back to verify")
	retrievedBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, testBlob, retrievedBlob)
}

func TestProxyMemConfigClientCanGetAndPatchV1(t *testing.T) {
	testProxyMemConfigClientCanGetAndPatch(t, common.V1EigenDABackend)
}

func TestProxyMemConfigClientCanGetAndPatchV2(t *testing.T) {
	testProxyMemConfigClientCanGetAndPatch(t, common.V2EigenDABackend)
}

func testProxyMemConfigClientCanGetAndPatch(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	useMemstore := testutils.GetBackend() == testutils.MemstoreBackend
	if !useMemstore {
		t.Skip("test can't be run against holesky since read failure case can't be manually triggered")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	memClient := memconfig_client.New(
		&memconfig_client.Config{
			URL: "http://" + ts.RestServer.Endpoint(),
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
		common.V1EigenDABackend,
		[]common.EigenDABackend{common.V1EigenDABackend, common.V2EigenDABackend})
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	testSuite, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	client := standard_client.New(
		&standard_client.Config{
			URL: testSuite.RestAddress(),
		})

	// disperse a payload to v1
	payload1a := testRandom.Bytes(1000)
	cert1a, err := client.SetData(testSuite.Ctx, payload1a)
	require.NoError(t, err)

	// disperse a payload to v2
	testSuite.RestServer.SetDispersalBackend(common.V2EigenDABackend)
	payload2a := testRandom.Bytes(1000)
	cert2a, err := client.SetData(testSuite.Ctx, payload2a)
	require.NoError(t, err)

	// disperse another payload to v1
	testSuite.RestServer.SetDispersalBackend(common.V1EigenDABackend)
	payload1b := testRandom.Bytes(1000)
	cert1b, err := client.SetData(testSuite.Ctx, payload1b)
	require.NoError(t, err)

	// disperse another payload to v2
	testSuite.RestServer.SetDispersalBackend(common.V2EigenDABackend)
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
	requireDispersalRetrievalEigenDA(t, testSuite.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestMaxBlobSizeV1(t *testing.T) {
	if testutils.GetBackend() == testutils.PreprodBackend {
		t.Skip("Preprod for v1 has a stricter blob size than normal.")
	}

	testMaxBlobSize(t, common.V1EigenDABackend)
}

func TestMaxBlobSizeV2(t *testing.T) {
	testMaxBlobSize(t, common.V2EigenDABackend)
}

func testMaxBlobSize(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.MaxBlobLength = "16mib"
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	// the payload has things added to it during encoding, so it has a slightly lower limit than max blob size
	maxPayloadSize, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(tsConfig.StoreBuilderConfig.ClientConfigV2.MaxBlobSizeBytes / encoding.BYTES_PER_SYMBOL))
	require.NoError(t, err)

	requireStandardClientSetGet(t, ts, testutils.RandBytes(int(maxPayloadSize)))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

// TestV2ValidatorRetrieverOnly tests that retrieval works when only the validator retriever is enabled
func TestV2ValidatorRetrieverOnly(t *testing.T) {
	if testutils.GetBackend() == testutils.MemstoreBackend {
		t.Skip("Don't run for memstore backend, since memstore tests don't actually hit the retrievers")
	}

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), common.V2EigenDABackend, nil)
	// Modify the test config to only use the validator retriever
	testCfg.Retrievers = []common.RetrieverType{common.ValidatorRetrieverType}

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	requireStandardClientSetGet(t, ts, testutils.RandBytes(1000))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

// requireDispersalRetrievalEigenDA ... ensure that blob was successfully dispersed/read to/from EigenDA
func requireDispersalRetrievalEigenDA(t *testing.T, cm *metrics.CountMap, mode commitments.CommitmentMode) {
	writeCount, err := cm.Get(string(mode), http.MethodPost)
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(string(mode), http.MethodGet)
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireWriteReadSecondary ... ensure that secondary backend was successfully written/read to/from
func requireWriteReadSecondary(t *testing.T, cm *metrics.CountMap, bt common.BackendType) {
	writeCount, err := cm.Get(http.MethodPut, secondary.Success, bt.String())
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(http.MethodGet, secondary.Success, bt.String())
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireStandardClientSetGet ... ensures that std proxy client can disperse and read a blob
func requireStandardClientSetGet(t *testing.T, ts testutils.TestSuite, blob []byte) {
	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
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
	daClient := altda.NewDAClient(ts.RestAddress(), false, precompute)

	commit, err := daClient.SetInput(ts.Ctx, blob)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ts.Ctx, commit, 0)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)
}
