package integration_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/memconfig_client"
	"github.com/Layr-Labs/eigenda/api/proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/config"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/proxy/test/testutils"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigensdk-go/logging"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ProxyTestConfig struct {
	UseMemstore      bool
	EnabledRestAPIs  *enablement.RestApisEnabled
	Expiration       time.Duration
	MaxBlobLength    string
	WriteThreadCount int
	WriteOnCacheMiss bool
	// at most one of the below options should be true
	UseKeccak256ModeS3            bool
	UseS3Caching                  bool
	UseS3Fallback                 bool
	ErrorOnSecondaryInsertFailure bool

	ClientLedgerMode     clientledger.ClientLedgerMode
	VaultMonitorInterval time.Duration
}

// NewTestConfig returns a new TestConfig
func NewTestConfig() ProxyTestConfig {
	return ProxyTestConfig{
		UseMemstore: false,
		EnabledRestAPIs: &enablement.RestApisEnabled{
			Admin:               false,
			OpGenericCommitment: true,
			OpKeccakCommitment:  true,
			StandardCommitment:  true,
		},
		Expiration:                    14 * 24 * time.Hour,
		UseKeccak256ModeS3:            false,
		UseS3Caching:                  false,
		UseS3Fallback:                 false,
		WriteThreadCount:              0,
		WriteOnCacheMiss:              false,
		ErrorOnSecondaryInsertFailure: false,
		ClientLedgerMode:              clientledger.ClientLedgerModeReservationOnly,
		VaultMonitorInterval:          30 * time.Second,
	}
}

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

	testCfg := NewTestConfig()
	testCfg.EnabledRestAPIs = &enablement.RestApisEnabled{
		Admin:               false,
		OpGenericCommitment: true,
		OpKeccakCommitment:  false,
		StandardCommitment:  false,
	}
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseKeccak256ModeS3 = true
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy REST server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseS3Caching = true
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseS3Fallback = true
	testCfg.UseMemstore = true
	// ensure that blob memstore eviction times result in near immediate activation
	testCfg.Expiration = time.Millisecond * 1
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

// TODO(iquidus): create s3 bucket in localstack for this test ?
/*
func TestProxyWriteCacheOnMiss(t *testing.T) {
	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := NewTestConfig()
	testCfg.UseS3Caching = true
	testCfg.WriteOnCacheMiss = true

	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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
*/

func TestErrorOnSecondaryInsertFailureFlagOn(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := NewTestConfig()
	testCfg.UseMemstore = true
	// Use S3 as fallback with invalid credentials to simulate S3 failure
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag
	// Ensure async writes are disabled (required for flag to work)
	testCfg.WriteThreadCount = 0

	// Create a test suite with invalid S3 config to force secondary write failures
	proxyConfig, err := createProxyConfig(testCfg)
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

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseMemstore = true
	// Use both cache and fallback - cache will fail, fallback will succeed
	testCfg.UseS3Caching = true
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = false // default: OFF
	testCfg.WriteThreadCount = 0

	// Create a test suite with invalid S3 config to force secondary write failures
	proxyConfig, err := createProxyConfig(testCfg)
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

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseMemstore = true
	testCfg.UseS3Fallback = true
	testCfg.ErrorOnSecondaryInsertFailure = true // Enable flag
	testCfg.WriteThreadCount = 0

	// Use valid S3 config
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseMemstore = true
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.UseMemstore = true
	testCfg.MaxBlobLength = "16mib"

	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	// Start proxy server
	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
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

	testCfg := NewTestConfig()
	testCfg.ClientLedgerMode = clientledger.ClientLedgerModeReservationOnly
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	// Test basic dispersal and retrieval with reservation payments
	blob := testutils.RandBytes(1000)
	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), blob)
	// Verify that dispersal and retrieval succeeded
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	t.Log("Successfully dispersed and retrieved blob using reservation-only payments")
}

// TODO(iquidus): Insufficent on-demand balance currently causes test to fail
/*
func TestOnDemandPayments(t *testing.T) {
	t.Parallel()

	testCfg := NewTestConfig()
	testCfg.ClientLedgerMode = clientledger.ClientLedgerModeOnDemandOnly
	proxyConfig, err := createProxyConfig(testCfg)
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	// Test basic dispersal and retrieval with on-demand payments
	blob := testutils.RandBytes(1000)
	requireStandardClientSetGet(t, ts.Ctx, ts.RestAddress(), blob)

	// Verify that dispersal and retrieval succeeded
	requireDispersalRetrievalEigenDA(t, ts.Metrics.HTTPServerRequestsTotal, commitments.StandardCommitmentMode)

	t.Log("Successfully dispersed and retrieved blob using on-demand-only payments")
}
*/

func TestArbCustomDAGetSupportedHeaderBytesMethod(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	proxyConfig.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	ethClient, err := geth.SafeDial(t.Context(), ts.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	var supportedHeaderBytesResult *arbitrum_altda.SupportedHeaderBytesResult
	err = rpcClient.Call(&supportedHeaderBytesResult,
		arbitrum_altda.MethodGetSupportedHeaderBytes)
	require.NoError(t, err)
	require.Equal(t, supportedHeaderBytesResult.HeaderBytes[0][0], uint8(commitments.ArbCustomDAHeaderByte))
}

func TestArbCustomDAStoreAndRecoverMethods(t *testing.T) {
	t.Parallel()

	testHarness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err)
	defer testHarness.Cleanup()

	testCfg := NewTestConfig()
	proxyConfig, err := createProxyConfig(testCfg)
	proxyConfig.EnabledServersConfig = &enablement.EnabledServersConfig{
		Metric:        false,
		ArbCustomDA:   true,
		RestAPIConfig: enablement.RestApisEnabled{},
	}
	require.NoError(t, err)

	ts, cleanup, err := startProxyServer(context.Background(), globalInfra.Logger, proxyConfig)
	require.NoError(t, err)
	defer cleanup()

	ethClient, err := geth.SafeDial(t.Context(), ts.ArbAddress())
	require.NoError(t, err)
	rpcClient := ethClient.Client()

	var storeResult *arbitrum_altda.StoreResult
	seqMessageArg := "0xDEADBEEF"
	timeoutArg := hexutil.Uint(200)

	err = rpcClient.Call(&storeResult, arbitrum_altda.MethodStore,
		seqMessageArg,
		timeoutArg)
	require.NoError(t, err)

	var recoverPayloadResult *arbitrum_altda.PayloadResult
	batchNum := hexutil.Uint(0)
	batchBlockHash := gethcommon.HexToHash("0x43")

	// pad 40 bytes for "message header"
	seqMessage := hexutil.Bytes(make([]byte, 40))
	seqMessage = append(seqMessage, storeResult.SerializedDACert...)

	err = rpcClient.Call(&recoverPayloadResult, arbitrum_altda.MethodRecoverPayload,
		batchNum,
		batchBlockHash,
		seqMessage,
	)
	require.NoError(t, err)

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
func requireDispersalRetrievalEigenDA(t *testing.T, cm *metrics.CountMap, mode commitments.CommitmentMode) {
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
func requireWriteReadSecondary(t *testing.T, cm *metrics.CountMap, bt common.BackendType) {
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

// createProxyConfig creates a proxy configuration that connects to the inabox disperser
func createProxyConfig(
	testCfg ProxyTestConfig,
) (config.AppConfig, error) {
	payloadClientConfig := clientsv2.PayloadClientConfig{
		PayloadPolynomialForm: codecs.PolynomialFormEval,
		BlobVersion:           0,
	}

	// Get the Ethereum RPC URL from the test config
	ethRPCURL := globalInfra.TestConfig.Deployers[0].RPC

	// Get the disperser API server address from the infrastructure
	disperserHostname, disperserPort, err := net.SplitHostPort(globalInfra.DisperserHarness.APIServerAddress)
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("invalid disperser API server address: %w", err)
	}

	// Get SRS paths using the utility function
	g1Path, g2Path, g2TrailingPath, err := integration.GetSRSPaths()
	if err != nil {
		return config.AppConfig{}, fmt.Errorf("failed to determine SRS file paths: %w", err)
	}

	// Construct cache directory path from g1Path
	srsDir := filepath.Dir(g1Path)
	cacheDir := filepath.Join(srsDir, "SRSTables")

	// Define max blob length
	maxBlobLengthBytes := uint64(16 * 1024 * 1024) // 16 MiB

	builderConfig := builder.Config{
		StoreConfig: store.Config{
			AsyncPutWorkers:               testCfg.WriteThreadCount,
			BackendsToEnable:              []common.EigenDABackend{common.V2EigenDABackend},
			DispersalBackend:              common.V2EigenDABackend,
			WriteOnCacheMiss:              testCfg.WriteOnCacheMiss,
			ErrorOnSecondaryInsertFailure: testCfg.ErrorOnSecondaryInsertFailure,
		},
		ClientConfigV2: common.ClientConfigV2{
			DisperserClientCfg: dispersal.DisperserClientConfig{
				Hostname:          disperserHostname,
				Port:              disperserPort,
				UseSecureGrpcFlag: false, // inabox uses insecure gRPC
			},
			PayloadDisperserCfg: dispersal.PayloadDisperserConfig{
				PayloadClientConfig:    payloadClientConfig,
				DisperseBlobTimeout:    5 * time.Minute,
				BlobCompleteTimeout:    5 * time.Minute,
				BlobStatusPollInterval: 1 * time.Second,
				ContractCallTimeout:    5 * time.Second,
			},
			RelayPayloadRetrieverCfg: payloadretrieval.RelayPayloadRetrieverConfig{
				PayloadClientConfig: payloadClientConfig,
				RelayTimeout:        5 * time.Second,
			},
			PutTries:                           3,
			MaxBlobSizeBytes:                   maxBlobLengthBytes,
			EigenDACertVerifierOrRouterAddress: globalInfra.TestConfig.EigenDA.CertVerifierRouter,
			EigenDADirectory:                   globalInfra.TestConfig.EigenDA.EigenDADirectory,
			EigenDANetwork:                     "", // Empty for inabox (custom network)
			RetrieversToEnable:                 []common.RetrieverType{common.RelayRetrieverType, common.ValidatorRetrieverType},
			ClientLedgerMode:                   testCfg.ClientLedgerMode,
			VaultMonitorInterval:               testCfg.VaultMonitorInterval,
		},
		KzgConfig: kzg.KzgConfig{
			G1Path:          g1Path,
			G2Path:          g2Path,
			G2TrailingPath:  g2TrailingPath,
			CacheDir:        cacheDir,
			SRSOrder:        encoding.SRSOrder,
			NumWorker:       uint64(runtime.GOMAXPROCS(0)), // #nosec G115
			SRSNumberToLoad: maxBlobLengthBytes / 32,
			LoadG2Points:    true,
		},
		MemstoreConfig: memconfig.NewSafeConfig(
			memconfig.Config{
				BlobExpiration:   testCfg.Expiration,
				MaxBlobSizeBytes: maxBlobLengthBytes,
			}),
		MemstoreEnabled: testCfg.UseMemstore,
		VerifierConfigV1: verify.Config{
			VerifyCerts:          true,
			RPCURL:               ethRPCURL,
			EthConfirmationDepth: 1,
			WaitForFinalization:  false,
			MaxBlobSizeBytes:     maxBlobLengthBytes,
		},
	}

	localstack := globalInfra.DisperserHarness.LocalStack
	awsConfig := localstack.GetAWSClientConfig()
	awsEndpoint := strings.TrimPrefix(awsConfig.EndpointURL, "http://")
	s3Config := s3.Config{
		Bucket:          globalInfra.DisperserHarness.S3Buckets.BlobStore,
		Path:            "",
		Endpoint:        awsEndpoint,
		EnableTLS:       false,
		AccessKeySecret: awsConfig.SecretAccessKey,
		AccessKeyID:     awsConfig.AccessKey,
		CredentialType:  s3.CredentialTypeStatic,
	}

	switch {
	case testCfg.UseKeccak256ModeS3:
		builderConfig.S3Config = s3Config
	case testCfg.UseS3Caching:
		builderConfig.StoreConfig.CacheTargets = []string{"S3"}
		builderConfig.S3Config = s3Config
	case testCfg.UseS3Fallback:
		builderConfig.StoreConfig.FallbackTargets = []string{"S3"}
		builderConfig.S3Config = s3Config
	}

	secretConfig := common.SecretConfigV2{
		SignerPaymentKey: integration.GetDefaultTestPayloadDisperserConfig().PrivateKey,
		EthRPCURL:        ethRPCURL,
	}

	return config.AppConfig{
		StoreBuilderConfig: builderConfig,
		SecretConfig:       secretConfig,
		EnabledServersConfig: &enablement.EnabledServersConfig{
			Metric:        false,
			ArbCustomDA:   false,
			RestAPIConfig: *testCfg.EnabledRestAPIs,
		},
		MetricsSvrConfig: proxy_metrics.Config{
			Host: "127.0.0.1",
			Port: 0, // let OS assign a free port
		},
		RestSvrCfg: rest.Config{
			Host:        "127.0.0.1",
			Port:        0, // let OS assign a free port
			APIsEnabled: testCfg.EnabledRestAPIs,
		},
		ArbCustomDASvrCfg: arbitrum_altda.Config{
			Host: "127.0.0.1",
			Port: 0, // let OS assign a free port
		},
	}, nil
}

// startProxyServer starts a proxy REST server with the given configuration
func startProxyServer(
	ctx context.Context,
	logger logging.Logger,
	appConfig config.AppConfig,
) (*testutils.TestSuite, func(), error) {
	if err := appConfig.Check(); err != nil {
		return nil, nil, fmt.Errorf("invalid app config: %w", err)
	}

	var (
		restServer *rest.Server
		arbServer  *arbitrum_altda.Server
		metrics    = proxy_metrics.NewEmulatedMetricer()
	)

	// Build the eth client for contract interactions
	ethClient, chainID, err := common.BuildEthClient(
		ctx,
		logger,
		appConfig.SecretConfig.EthRPCURL,
		appConfig.StoreBuilderConfig.ClientConfigV2.EigenDANetwork,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("build eth client: %w", err)
	}

	// Build storage managers
	certMgr, keccakMgr, err := builder.BuildManagers(
		ctx,
		logger,
		metrics,
		appConfig.StoreBuilderConfig,
		appConfig.SecretConfig,
		nil,
		ethClient,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("build storage managers: %w", err)
	}

	// Create compatibility config
	compatibilityCfg, err := common.NewCompatibilityConfig(
		"test",
		chainID,
		appConfig.StoreBuilderConfig.ClientConfigV2,
		false, // readOnlyMode
		appConfig.EnabledServersConfig.ToAPIStrings(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("new compatibility config: %w", err)
	}

	if appConfig.EnabledServersConfig.RestAPIConfig.DAEndpointEnabled() {
		// Create and start REST server
		appConfig.RestSvrCfg.CompatibilityCfg = compatibilityCfg
		restServer = rest.NewServer(appConfig.RestSvrCfg, certMgr, keccakMgr, logger, metrics)
		router := mux.NewRouter()
		restServer.RegisterRoutes(router)
		if appConfig.StoreBuilderConfig.MemstoreEnabled {
			memconfig.NewHandlerHTTP(logger, appConfig.StoreBuilderConfig.MemstoreConfig).
				RegisterMemstoreConfigHandlers(router)
		}

		if err := restServer.Start(router); err != nil {
			return nil, nil, fmt.Errorf("start proxy server: %w", err)
		}
	}

	if appConfig.EnabledServersConfig.ArbCustomDA {
		arbHandlers := arbitrum_altda.NewHandlers(certMgr, logger, compatibilityCfg)
		arbServer, err = arbitrum_altda.NewServer(ctx, &appConfig.ArbCustomDASvrCfg, arbHandlers)
		if err != nil {
			return nil, nil, fmt.Errorf("create arbitrum server: %v", err.Error())
		}

		if err := arbServer.Start(); err != nil {
			return nil, nil, fmt.Errorf("start arbitrum server: %v", err.Error())
		}
	}

	cleanup := func() {
		if appConfig.EnabledServersConfig.RestAPIConfig.DAEndpointEnabled() {
			if err := restServer.Stop(); err != nil {
				logger.Error("failed to stop proxy server", "err", err)
			}
		}

		if appConfig.EnabledServersConfig.ArbCustomDA {
			if err := arbServer.Stop(); err != nil {
				logger.Error("failed to stop arb server", "err", err)
			}
		}
	}

	return &testutils.TestSuite{
		Ctx:        ctx,
		Log:        logger,
		RestServer: restServer,
		Metrics:    metrics,
		ArbServer:  arbServer,
	}, cleanup, nil
}
