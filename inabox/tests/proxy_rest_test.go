package integration_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/memconfig_client"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProxyAPIsEnabledRestALTDA tests to ensure that the enabled APIs expression is
// getting respected by the REST ALTDA Server when wiring up a proxy application instance
// with just `op-generic` mode enabled.
//
// This test has been migrated from api/proxy/test/e2e/server_rest_test.go to use inabox infrastructure.
func TestProxyAPIsEnabledRestALTDA(t *testing.T) {
	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.EnabledRestAPIs = &enablement.RestApisEnabled{
		Admin:               false,
		OpGenericCommitment: true,
		OpKeccakCommitment:  false,
		StandardCommitment:  false,
	}
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	t.Logf("Proxy server started at %s", ts.RestAddress())

	// Test that standard commitment mode is disabled (should return 403)
	standardClient := standard_client.New(&standard_client.Config{
		URL: ts.RestAddress(),
	})

	testBlob := []byte("hello world")
	t.Log("Attempting to set data using standard commitment (should fail with 403)...")
	_, err = standardClient.SetData(ts.Ctx, testBlob)
	require.Error(t, err)
	require.ErrorContains(t, err, "403")

	// Test that op-generic mode works (should succeed)
	opGenericClient := altda.NewDAClient(ts.RestAddress(), false, false)

	t.Log("Setting data using op-generic commitment (should succeed)...")
	daCommit, err := opGenericClient.SetInput(ts.Ctx, testBlob)
	require.NoError(t, err)

	t.Log("Getting data using op-generic commitment (should succeed)...")
	preimage, err := opGenericClient.GetInput(ts.Ctx, daCommit, 0)
	require.NoError(t, err)
	require.Equal(t, testBlob, preimage)

	t.Log("TestProxyAPIsEnabledRestALTDA completed successfully")

	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

// TestProxyClientWriteRead tests that the proxy client can write and read data to the proxy server.
// This is the "basic" proxy test: "is proxy working?"
//
// This test has been migrated from api/proxy/test/e2e/server_rest_test.go to use inabox infrastructure.
func TestProxyClientWriteRead(t *testing.T) {
	t.Parallel()

	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	t.Logf("Proxy server started at %s", ts.RestAddress())
	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(100))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

func TestOptimismClientWithKeccak256Commitment(t *testing.T) {
	t.Parallel()

	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseKeccak256ModeS3 = true
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy REST server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	requireOPClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(100), true)
	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

func TestOptimismClientWithGenericCommitment(t *testing.T) {
	t.Parallel()

	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	t.Logf("Proxy server started at %s", ts.RestAddress())
	requireOPClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(100), false)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.OptimismGenericCommitmentMode)

	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

// TODO(iquidus): Determine why this test is failing due to connection refused
// TestProxyClientServerIntegration tests the proxy client and server integration by setting the data as a single byte,
// many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func TestProxyClientServerIntegration(t *testing.T) {
	t.Parallel()

	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	// single byte preimage set data case
	testPreimage := []byte{1} // single byte preimage
	t.Log("Setting input data on proxy server...")
	_, err = daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	// unicode preimage set data case
	testPreimage = []byte("§§©ˆªªˆ˙√ç®∂§∞¶§ƒ¥√¨¥√¨¥ƒƒ©˙˜ø˜˜˜∫˙∫¥∫√†®®√ç¨ˆ¨˙ï") // many unicode characters
	t.Log("Setting input data on proxy server...")
	_, err = daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	testPreimage = []byte("§") // single unicode character
	t.Log("Setting input data on proxy server...")
	_, err = daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	// empty preimage set data case
	testPreimage = []byte("") // Empty preimage
	t.Log("Setting input data on proxy server...")
	_, err = daClient.SetData(ts.Ctx, testPreimage)
	require.NoError(t, err)

	// get data edge cases
	testCert := []byte("")
	_, err = daClient.GetData(ts.Ctx, testCert)
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

	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

// Ensure that proxy is able to write/read from a cache backend when enabled
func TestProxyCaching(t *testing.T) {
	t.Parallel()

	// Create fresh test harness from global infrastructure
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseS3Caching = true
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(100_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	// Verify the server is still running
	require.NotNil(t, ts.RestServer)
}

// Ensure that fallback location is read from when EigenDA blob is not available.
// This is done by setting the memstore expiration time to 1ms and waiting for the blob to expire
// before attempting to read it.
func TestProxyReadFallback(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseS3Fallback = true
	testCfg.UseMemstore = true
	// ensure that blob memstore eviction times result in near immediate activation
	testCfg.Expiration = time.Millisecond * 1
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)
	expectedBlob := testutils.RandBytes(100_000)
	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, expectedBlob)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	t.Log("Getting input data from proxy server...")
	actualBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, actualBlob)

	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(100_000))
	requireWriteReadSecondary(t, ts.Metrics.SecondaryRequestsTotal, common.S3BackendType)
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestProxyWriteCacheOnMiss(t *testing.T) {
	t.Skip("TODO(iquidus): create s3 bucket in localstack for this test")
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseS3Caching = true
	testCfg.WriteOnCacheMiss = true

	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)
	expectedBlob := testutils.RandBytes(100_000)
	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, expectedBlob)
	require.NoError(t, err)

	_, err = daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)

	exists, err := testutils.ExistsBlobInfotInBucket(proxyConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.True(t, exists)

	t.Log("Erase blob from the cache...")
	err = testutils.RemoveBlobInfoFromBucket(proxyConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	exists, err = testutils.ExistsBlobInfotInBucket(proxyConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.False(t, exists)

	// Blob created in disperser, removed from S3
	t.Log("Getting input data from proxy server...")
	actualBlob, err := daClient.GetData(ts.Ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, expectedBlob, actualBlob)

	exists, err = testutils.ExistsBlobInfotInBucket(proxyConfig.StoreBuilderConfig.S3Config.Bucket, blobInfo)
	require.NoError(t, err)
	require.True(t, exists)
}

func TestErrorOnSecondaryInsertFailureFlagOn(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseMemstore = true
	// Use S3 as fallback with invalid credentials to simulate S3 failure
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag
	// Ensure async writes are disabled (required for flag to work)
	testCfg.WriteThreadCount = 0

	// Create a test suite with invalid S3 config to force secondary write failures
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)
	// Override S3 config with invalid credentials to force write failures
	proxyConfig.StoreBuilderConfig.S3Config = s3.Config{
		Bucket:          "invalid-bucket-name",
		Endpoint:        "invalid-endpoint:9000",
		AccessKeyID:     "invalid-key",
		AccessKeySecret: "invalid-secret",
		EnableTLS:       false,
		CredentialType:  s3.CredentialTypeStatic,
	}

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	testBlob := testutils.RandBytes(100)

	cfg := &standard_client.Config{
		URL: ts.RestAddress(),
	}
	daClient := standard_client.New(cfg)

	// PUT should fail because S3 write fails and flag is ON
	t.Log("Setting data - should fail due to S3 failure with flag enabled")
	_, err = daClient.SetData(ts.Ctx, testBlob)
	require.Error(t, err, "PUT should fail when error-on-secondary-insert-failure=true and S3 fails")

	// Error should indicate it's a server error (5xx)
	require.Contains(t, err.Error(), "500", "Expected HTTP 500 error")
}

func TestErrorOnSecondaryInsertFailureFlagOffPartialFailure(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseMemstore = true
	// Use both cache and fallback - cache will fail, fallback will succeed
	testCfg.UseS3Caching = true
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = false // default: OFF
	testCfg.WriteThreadCount = 0

	// Create a test suite with invalid S3 config to force secondary write failures
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)
	// Override S3 config with invalid credentials to force write failures
	proxyConfig.StoreBuilderConfig.S3Config = s3.Config{
		Bucket:          "invalid-bucket-name",
		Endpoint:        "invalid-endpoint:9000",
		AccessKeyID:     "invalid-key",
		AccessKeySecret: "invalid-secret",
		EnableTLS:       false,
		CredentialType:  s3.CredentialTypeStatic,
	}

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

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

func TestErrorOnSecondaryInsertFailureFlagOnSuccess(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseMemstore = true
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag
	testCfg.WriteThreadCount = 0

	// Use valid S3 config
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

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

// Test can't be ran against inabox backend since read failure case can't be manually triggered
func TestProxyMemConfigClientCanGetAndPatch(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseMemstore = true
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

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

func TestMaxBlobSize(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.UseMemstore = true
	testCfg.MaxBlobLength = "16mib"

	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	// the payload has things added to it during encoding, so it has a slightly lower limit than max blob size
	maxPayloadSize, err := codec.BlobSymbolsToMaxPayloadSize(
		uint32(proxyConfig.StoreBuilderConfig.ClientConfigV2.MaxBlobSizeBytes / encoding.BYTES_PER_SYMBOL))
	require.NoError(t, err)

	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(int(maxPayloadSize)))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestValidatorRetrieverOnly(t *testing.T) {
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), testutils.RandBytes(1000))
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)
}

func TestReservationPayments(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.ClientLedgerMode = clientledger.ClientLedgerModeReservationOnly
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	// Test basic dispersal and retrieval with reservation payments
	blob := testutils.RandBytes(1000)
	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), blob)
	// Verify that dispersal and retrieval succeeded
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	t.Log("Successfully dispersed and retrieved blob using reservation-only payments")
}

func TestOnDemandPayments(t *testing.T) {
	t.Skip("TODO(iquidus): Insufficent on-demand balance currently causes test to fail")
	t.Parallel()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	testCfg.ClientLedgerMode = clientledger.ClientLedgerModeOnDemandOnly
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	// Test basic dispersal and retrieval with on-demand payments
	blob := testutils.RandBytes(1000)
	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), blob)

	// Verify that dispersal and retrieval succeeded
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	t.Log("Successfully dispersed and retrieved blob using on-demand-only payments")
}

// OP contract tests
// Contract Test here refers to https://pactflow.io/blog/what-is-contract-testing/, not evm contracts.
func TestOPContractTestRBNRecentyCheck(t *testing.T) {
	t.Skip("TODO(iquidus): RBN recency check failed, fails")
	t.Parallel()

	var testTable = []struct {
		name                 string
		RBNRecencyWindowSize uint64
		certRBN              uint32
		certL1IBN            uint64
		requireErrorFn       func(t *testing.T, err error)
	}{
		{
			name:                 "RBN recency check failed",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            201,
			requireErrorFn: func(t *testing.T, err error) {
				// expect proxy to return a 418 error which the client converts to this structured error
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t,
					int(coretypes.ErrRecencyCheckFailedDerivationError.StatusCode),
					dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check passed",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            199,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check skipped - Proxy set window size 0",
			RBNRecencyWindowSize: 0,
			certRBN:              100,
			certL1IBN:            201,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name:                 "RBN recency check skipped - client set IBN to 0",
			RBNRecencyWindowSize: 100,
			certRBN:              100,
			certL1IBN:            0,
			requireErrorFn: func(t *testing.T, err error) {
				// After RBN check succeeds, CertVerifier.checkDACert contract call is made,
				// which returns a [verification.CertVerificationFailedError] with StatusCode 2 (inclusion proof
				// invalid). This gets converted to a [eigendav2store.ErrInvalidCertDerivationError] which gets marshalled
				// and returned as the body of a 418 response by the proxy.
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log("Running test: ", tt.name)
			testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
			require.NoError(t, err)
			defer testHarness.Cleanup()

			testCfg := integration.NewProxyTestConfig(globalInfra)
			proxyConfig, err := integration.CreateProxyConfig(testCfg)
			require.NoError(t, err)

			ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
			require.NoError(t, err)
			defer cleanup()

			// Build + Serialize (empty) cert with the given RBN
			certV3 := coretypes.EigenDACertV3{
				BatchHeader: bindings.EigenDATypesV2BatchHeaderV2{
					ReferenceBlockNumber: tt.certRBN,
				},
			}
			serializedCertV3, err := rlp.EncodeToBytes(certV3)
			require.NoError(t, err)
			// altdaCommitment is what is returned by the proxy
			altdaCommitment, err := commitments.EncodeCommitment(
				certs.NewVersionedCert(serializedCertV3, certs.V2VersionByte),
				commitments.OptimismGenericCommitmentMode)
			require.NoError(t, err)
			// the op client expects a typed commitment, so we have to decode the altdaCommitment
			commitmentData, err := altda.DecodeCommitmentData(altdaCommitment)
			require.NoError(t, err)

			daClient := altda.NewDAClient(ts.RestAddress(), false, false)
			_, err = daClient.GetInput(ts.Ctx, commitmentData, tt.certL1IBN)
			tt.requireErrorFn(t, err)
		})
	}
}

// Test that proxy DerivationErrors are correctly parsed as DropCommitmentErrors on op side,
// for parsing and cert validation errors.
func TestOPContractTestValidAndInvalidCertErrors(t *testing.T) {
	t.Skip("TODO(iquidus): connection refused error")
	t.Parallel()

	var testTable = []struct {
		name           string
		certCreationFn func() ([]byte, error)
		requireErrorFn func(t *testing.T, err error)
	}{
		{
			// TODO: need to figure out why this is happening, since ErrNotFound is supposed to be a keccak only error.
			// Seems like op-client allows submitting an empty cert, and because its not a valid cert request, it gets
			// matched by proxy's keccak commitment handler, which returns ErrNotFound (there is no such key in the store).
			// I think this is ok behavior... since it would be a bug to submit an empty cert....?
			// But need to think about this more.
			name: "empty cert returns ErrNotFound",
			certCreationFn: func() ([]byte, error) {
				return []byte{}, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				require.ErrorIs(t, err, altda.ErrNotFound)
			},
		},
		{
			name: "cert parsing error",
			certCreationFn: func() ([]byte, error) {
				cert := make([]byte, 10)
				return cert, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrCertParsingFailedDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
		{
			name: "invalid (default) cert",
			certCreationFn: func() ([]byte, error) {
				// Build + Serialize invalid default cert
				certV3 := coretypes.EigenDACertV3{}
				serializedCertV3, err := rlp.EncodeToBytes(certV3)
				if err != nil {
					return nil, err
				}
				return serializedCertV3, nil
			},
			requireErrorFn: func(t *testing.T, err error) {
				var dropEigenDACommitmentErr altda.DropEigenDACommitmentError
				require.ErrorAs(t, err, &dropEigenDACommitmentErr)
				require.Equal(t, int(coretypes.ErrInvalidCertDerivationError.StatusCode), dropEigenDACommitmentErr.StatusCode)
			},
		},
	}

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := integration.NewProxyTestConfig(globalInfra)
	proxyConfig, err := integration.CreateProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := integration.StartProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Log("Running test: ", tt.name)
			serializedCert, err := tt.certCreationFn()
			require.NoError(t, err)

			altdaCommitment, err := commitments.EncodeCommitment(
				certs.NewVersionedCert(serializedCert, certs.V2VersionByte),
				commitments.OptimismGenericCommitmentMode)
			require.NoError(t, err)
			// the op client expects a typed commitment, so we have to decode the altdaCommitment
			commitmentData, err := altda.DecodeCommitmentData(altdaCommitment)
			require.NoError(t, err)

			daClient := altda.NewDAClient(ts.RestAddress(), false, false)
			_, err = daClient.GetInput(ts.Ctx, commitmentData, 0)

			tt.requireErrorFn(t, err)
		})
	}
}

func TestOPContractTestBlobDecodingErrors(t *testing.T) {
	// Writing this test is a lot more involved... because we need to populate mock relay backends
	// that would return a blob that doesn't decode properly.
	// Probably will require adding this after we've created a better test suite framework for the eigenda clients.
	t.Skip("TODO: implement blob decoding errors test")
}

// requireStandardClientSetGet ... ensures that std proxy client can disperse and read a blob
func requireStandardClientSetGet(t *testing.T, ctx context.Context, restEndpoint string, blob []byte) {
	cfg := &standard_client.Config{
		URL: restEndpoint,
	}
	daClient := standard_client.New(cfg)

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ctx, blob)
	require.NoError(t, err)

	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ctx, blobInfo)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)
}

// requireDispersalRetrievalEigenDA ... ensure that blob was successfully dispersed/read to/from EigenDA
func requireDispersalRetrievalEigenDA(t *testing.T, cm *proxy_metrics.CountMap, mode commitments.CommitmentMode) {
	writeCount, err := cm.Get(string(mode), http.MethodPost)
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(string(mode), http.MethodGet)
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

// requireOPClientSetGet ... ensures that alt-da client can disperse and read a blob
func requireOPClientSetGet(t *testing.T, ctx context.Context, restEndpoint string, blob []byte, precompute bool) {
	daClient := altda.NewDAClient(restEndpoint, false, precompute)

	commit, err := daClient.SetInput(ctx, blob)
	require.NoError(t, err)

	preimage, err := daClient.GetInput(ctx, commit, 0)
	require.NoError(t, err)
	require.Equal(t, blob, preimage)
}

// requireWriteReadSecondary ... ensure that secondary backend was successfully written/read to/from
func requireWriteReadSecondary(t *testing.T, cm *proxy_metrics.CountMap, bt common.BackendType) {
	writeCount, err := cm.Get(http.MethodPut, secondary.Success, bt.String())
	require.NoError(t, err)
	require.True(t, writeCount > 0)

	readCount, err := cm.Get(http.MethodGet, secondary.Success, bt.String())
	require.NoError(t, err)
	require.True(t, readCount > 0)
}

func isNilPtrDerefPanic(err string) bool {
	return strings.Contains(err, "panic") && strings.Contains(err, "SIGSEGV") &&
		strings.Contains(err, "nil pointer dereference")
}
