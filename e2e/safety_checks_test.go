package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/testutils"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func isNilPtrDerefPanic(err string) bool {
	return strings.Contains(err, "panic") && strings.Contains(err, "SIGSEGV") &&
		strings.Contains(err, "nil pointer dereference")
}

func TestOpClientKeccak256MalformedInputsV1(t *testing.T) {
	testOpClientKeccak256MalformedInputs(t, common.V1EigenDABackend)
}

func TestOpClientKeccak256MalformedInputsV2(t *testing.T) {
	testOpClientKeccak256MalformedInputs(t, common.V2EigenDABackend)
}

// TestOpClientKeccak256MalformedInputs tests the NewDAClient from altda by setting and getting against []byte("")
// preimage. It sets the precompute option to false on the NewDAClient.
func testOpClientKeccak256MalformedInputs(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)

	testCfg.UseKeccak256ModeS3 = true
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	// nil commitment. Should return an error but currently is not. This needs to be fixed by OP
	// Ref: https://github.com/ethereum-optimism/optimism/issues/11987
	// daClient := altda.NewDAClient(ts.Address(), false, true)
	// t.Run("nil commitment case", func(t *testing.T) {
	//	var commit altda.CommitmentData
	//	_, err := daClient.GetInput(ts.Ctx, commit)
	//	require.Error(t, err)
	//	assert.True(t, !isPanic(err.Error()))
	// })

	daClientPcFalse := altda.NewDAClient(ts.Address(), false, false)

	t.Run(
		"input bad data to SetInput & GetInput", func(t *testing.T) {
			t.Parallel()

			testPreimage := []byte("") // Empty preimage
			_, err := daClientPcFalse.SetInput(ts.Ctx, testPreimage)
			require.Error(t, err)

			// should fail with proper error message as is now, and cannot contain panics or nils
			assert.True(
				t,
				strings.Contains(err.Error(), "invalid input") && !isNilPtrDerefPanic(err.Error()))

			// The below test panics silently.
			input := altda.NewGenericCommitment([]byte(""))
			_, err = daClientPcFalse.GetInput(ts.Ctx, input, 0)
			require.Error(t, err)

			// Should not fail on slice bounds out of range. This needs to be fixed by OP.
			// Refer to issue: https://github.com/ethereum-optimism/optimism/issues/11987
			// assert.False(t, strings.Contains(err.Error(), ": EOF") && !isPanic(err.Error()))
		})
}

func TestProxyClientMalformedInputCasesV1(t *testing.T) {
	testProxyClientMalformedInputCases(t, common.V1EigenDABackend)
}

func TestProxyClientMalformedInputCasesV2(t *testing.T) {
	testProxyClientMalformedInputCases(t, common.V2EigenDABackend)
}

// TestProxyClientMalformedInputCases tests the proxy client and server integration by setting the data as a single
// byte, many unicode characters, single unicode character and an empty preimage. It then tries to get the data from the
// proxy server with empty byte, single byte and random string.
func testProxyClientMalformedInputCases(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()
	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	t.Run(
		"single byte preimage set data case", func(t *testing.T) {
			t.Parallel()
			ts, kill := testutils.CreateTestSuite(tsConfig)
			defer kill()

			cfg := &standard_client.Config{
				URL: ts.Address(),
			}
			daClient := standard_client.New(cfg)
			testPreimage := []byte{1} // single byte preimage
			t.Log("Setting input data on proxy server...")
			_, err := daClient.SetData(ts.Ctx, testPreimage)
			require.NoError(t, err)
		})

	t.Run(
		"unicode preimage set data case", func(t *testing.T) {
			t.Parallel()
			ts, kill := testutils.CreateTestSuite(tsConfig)
			defer kill()

			cfg := &standard_client.Config{
				URL: ts.Address(),
			}
			daClient := standard_client.New(cfg)
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
			ts, kill := testutils.CreateTestSuite(tsConfig)
			defer kill()

			cfg := &standard_client.Config{
				URL: ts.Address(),
			}
			daClient := standard_client.New(cfg)
			testPreimage := []byte("") // Empty preimage
			t.Log("Setting input data on proxy server...")
			_, err := daClient.SetData(ts.Ctx, testPreimage)
			require.NoError(t, err)
		})

	t.Run(
		"get data edge cases - unsupported version byte 02", func(t *testing.T) {
			t.Parallel()
			ts, kill := testutils.CreateTestSuite(tsConfig)
			defer kill()

			cfg := &standard_client.Config{
				URL: ts.Address(),
			}
			daClient := standard_client.New(cfg)
			testCert := []byte{2}
			_, err := daClient.GetData(ts.Ctx, testCert)
			require.Error(t, err)
			assert.True(
				t,
				strings.Contains(
					err.Error(),
					"unsupported version byte 02") && !isNilPtrDerefPanic(err.Error()))
		})

	// TODO: what exactly is this test testing? What is the edge case?
	// Error tested doesn't seem related to the cert being huge.
	t.Run(
		"get data edge cases - huge cert", func(t *testing.T) {
			t.Parallel()
			ts, kill := testutils.CreateTestSuite(tsConfig)
			defer kill()

			cfg := &standard_client.Config{
				URL: ts.Address(),
			}
			daClient := standard_client.New(cfg)
			// TODO: we need to add the 0 version byte at the beginning.
			// should this not be done automatically by the std_commitment client?
			testCert := append([]byte{0}, testutils.RandBytes(10000)...)
			_, err := daClient.GetData(ts.Ctx, testCert)
			require.Error(t, err)
			// Commenting as this error is not returned by memstore but this test is also run
			// against memstore when running `make test-e2e-local`.
			// assert.True(t, !isNilPtrDerefPanic(err.Error()) &&
			// 	strings.Contains(err.Error(),
			// 		"failed to decode DA cert to RLP format: rlp: expected input list for verify.Certificate"),
			// 	"error: %s", err.Error())
		})
}

func TestKeccak256CommitmentRequestErrorsWhenS3NotSetV1(t *testing.T) {
	testKeccak256CommitmentRequestErrorsWhenS3NotSet(t, common.V1EigenDABackend)
}

func TestKeccak256CommitmentRequestErrorsWhenS3NotSetV2(t *testing.T) {
	testKeccak256CommitmentRequestErrorsWhenS3NotSet(t, common.V2EigenDABackend)
}

// TestKeccak256CommitmentRequestErrorsWhenS3NotSet ensures that the proxy returns a client error in the event
// that an OP Keccak commitment mode is provided when S3 is non-configured server side
func testKeccak256CommitmentRequestErrorsWhenS3NotSet(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	testCfg.UseKeccak256ModeS3 = true

	tsConfig := testutils.BuildTestSuiteConfig(testCfg)
	tsConfig.EigenDAConfig.StorageConfig.S3Config.Endpoint = "localhost:1234"

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	daClient := altda.NewDAClient(ts.Address(), false, true)

	testPreimage := testutils.RandBytes(100)

	_, err := daClient.SetInput(ts.Ctx, testPreimage)
	// TODO: the server currently returns an internal server error. Should it return a 400 instead?
	require.Error(t, err)
}

func TestOversizedBlobRequestErrorsV1(t *testing.T) {
	testOversizedBlobRequestErrors(t, common.V1EigenDABackend)
}

func TestOversizedBlobRequestErrorsV2(t *testing.T) {
	testOversizedBlobRequestErrors(t, common.V2EigenDABackend)
}

func testOversizedBlobRequestErrors(t *testing.T, dispersalBackend common.EigenDABackend) {
	t.Parallel()

	testCfg := testutils.NewTestConfig(testutils.GetBackend(), dispersalBackend, nil)
	tsConfig := testutils.BuildTestSuiteConfig(testCfg)

	ts, kill := testutils.CreateTestSuite(tsConfig)
	defer kill()

	cfg := &standard_client.Config{
		URL: ts.Address(),
	}
	daClient := standard_client.New(cfg)
	//  17MB blob
	testPreimage := testutils.RandBytes(17_000_0000)

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.Ctx, testPreimage)
	require.Empty(t, blobInfo)
	require.Error(t, err)

	oversizedError := false
	// error returned from EigenDA V1 disperser
	if strings.Contains(err.Error(), "blob size cannot exceed") {
		oversizedError = true
	}

	// error caught within proxy
	if strings.Contains(err.Error(), "blob is larger than max blob size") {
		oversizedError = true
	}

	// error caught within proxy
	if strings.Contains(err.Error(), "http: request body too large") {
		oversizedError = true
	}

	require.True(t, oversizedError)
	require.Contains(t, err.Error(), fmt.Sprint(http.StatusBadRequest))
}
