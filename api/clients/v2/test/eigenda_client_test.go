package test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const g1Path = "../../../../inabox/resources/kzg/g1.point"
const payloadLength = 100

type ClientTester struct {
	Random          *testrandom.TestRandom
	Client          *clients.EigenDAClient
	MockRelayClient *clientsmock.MockRelayClient
	Codec           *codecs.DefaultBlobCodec
	G1Srs           []bn254.G1Affine
}

func (c *ClientTester) requireExpectations(t *testing.T) {
	c.MockRelayClient.AssertExpectations(t)
	// TODO: get rid of this
}

// buildClientTester sets up a client with mocks necessary for testing
func buildClientTester(t *testing.T) ClientTester {
	logger := logging.NewNoopLogger()
	clientConfig := &clients.EigenDAClientConfig{
		RelayTimeout: 50 * time.Millisecond,
	}

	mockRelayClient := clientsmock.MockRelayClient{}
	codec := codecs.NewDefaultBlobCodec()

	random := testrandom.NewTestRandom(t)

	g1Srs, err := kzg.ReadG1Points(g1Path, uint64(payloadLength), uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	client, err := clients.NewEigenDAClient(
		logger,
		random.Rand,
		clientConfig,
		&mockRelayClient,
		&codec,
		g1Srs)

	require.NotNil(t, client)
	require.NoError(t, err)

	return ClientTester{
		Random:          random,
		Client:          client,
		MockRelayClient: &mockRelayClient,
		Codec:           &codec,
		G1Srs:           g1Srs,
	}
}

func buildBlobAndCert(
	t *testing.T,
	tester ClientTester,
	relayKeys []core.RelayKey) (core.BlobKey, []byte, *core.BlobCertificate) {

	blobKey := core.BlobKey(tester.Random.Bytes(32))
	payloadBytes := tester.Random.Bytes(payloadLength)
	blobBytes, err := tester.Codec.EncodeBlob(payloadBytes)
	require.NoError(t, err)
	require.NotNil(t, blobBytes)

	kzgCommitment, err := verification.GenerateBlobCommitment(tester.G1Srs, blobBytes)
	require.NoError(t, err)
	require.NotNil(t, kzgCommitment)

	commitments := encoding.BlobCommitments{
		Commitment: kzgCommitment,
		Length:     uint(len(blobBytes)),
	}

	blobHeader := &core.BlobHeader{
		BlobCommitments: commitments,
	}

	return blobKey, blobBytes, &core.BlobCertificate{
		RelayKeys:  relayKeys,
		BlobHeader: blobHeader,
	}
}

// TestGetPayloadSuccess tests that a blob is received without error in the happy case
func TestGetPayloadSuccess(t *testing.T) {
	tester := buildClientTester(t)
	relayKeys := make([]core.RelayKey, 1)
	relayKeys[0] = tester.Random.Uint32()
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(blobBytes, nil).Once()

	payload, err := tester.Client.GetPayload(
		context.Background(),
		blobKey,
		blobCert)

	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.requireExpectations(t)
}

// TestRelayCallTimeout verifies that calls to the relay timeout after the expected duration
func TestRelayCallTimeout(t *testing.T) {
	tester := buildClientTester(t)
	relayKeys := make([]core.RelayKey, 1)
	relayKeys[0] = tester.Random.Uint32()
	blobKey, _, blobCert := buildBlobAndCert(t, tester, relayKeys)

	// the timeout should occur before the panic has a chance to be triggered
	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(
		nil, errors.New("timeout")).Once().Run(
		func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			select {
			case <-ctx.Done():
				// this is the expected case
				return
			case <-time.After(time.Second):
				panic("call should have timed out first")
			}
		})

	// the panic should be triggered, since it happens faster than the configured timout
	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(
		nil, errors.New("timeout")).Once().Run(
		func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond):
				// this is the expected case
				panic("call should not have timed out")
			}
		})

	require.NotPanics(
		t, func() {
			_, _ = tester.Client.GetPayload(context.Background(), blobKey, blobCert)
		})

	require.Panics(
		t, func() {
			_, _ = tester.Client.GetPayload(context.Background(), blobKey, blobCert)
		})

	tester.requireExpectations(t)
}

// TestRandomRelayRetries verifies correct behavior when some relays do not respond with the blob,
// requiring the client to retry with other relays.
func TestRandomRelayRetries(t *testing.T) {
	tester := buildClientTester(t)

	relayCount := 100
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	// for this test, only a single relay is online
	// we will be requiring that it takes a different amount of retries to dial this relay, since the array of relay keys to try is randomized
	onlineRelayKey := relayKeys[tester.Random.Intn(len(relayKeys))]

	offlineKeyMatcher := func(relayKey core.RelayKey) bool { return relayKey != onlineRelayKey }
	onlineKeyMatcher := func(relayKey core.RelayKey) bool { return relayKey == onlineRelayKey }
	var failedCallCount int
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(offlineKeyMatcher), blobKey).Return(
		nil,
		fmt.Errorf("offline relay")).Run(
		func(args mock.Arguments) {
			failedCallCount++
		})
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(onlineKeyMatcher), blobKey).Return(
		blobBytes,
		nil)

	// keep track of how many tries various blob retrievals require
	// this allows us to require that there is variability, i.e. that relay call order is actually random
	requiredTries := map[int]bool{}

	for i := 0; i < relayCount; i++ {
		failedCallCount = 0
		blob, err := tester.Client.GetPayload(context.Background(), blobKey, blobCert)
		require.NotNil(t, blob)
		require.NoError(t, err)

		requiredTries[failedCallCount] = true
	}

	// with 100 random tries, with possible values between 1 and 100, we can very confidently require that there are at least 10 unique values
	require.Greater(t, len(requiredTries), 10)

	tester.requireExpectations(t)
}

// TestNoRelayResponse tests functionality when none of the relays respond
func TestNoRelayResponse(t *testing.T) {
	tester := buildClientTester(t)

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, _, blobCert := buildBlobAndCert(t, tester, relayKeys)

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(nil, fmt.Errorf("offline relay"))

	blob, err := tester.Client.GetPayload(
		context.Background(),
		blobKey,
		blobCert)
	require.Nil(t, blob)
	require.NotNil(t, err)

	tester.requireExpectations(t)
}

// TestNoRelays tests that having no relay keys is handled gracefully
func TestNoRelays(t *testing.T) {
	tester := buildClientTester(t)
	blobKey, _, blobCert := buildBlobAndCert(t, tester, []core.RelayKey{})

	blob, err := tester.Client.GetPayload(context.Background(), blobKey, blobCert)
	require.Nil(t, blob)
	require.NotNil(t, err)

	tester.requireExpectations(t)
}

// TestGetBlobReturns0Len verifies that a 0 length blob returned from a relay is handled gracefully, and that the client retries after such a failure
func TestGetBlobReturns0Len(t *testing.T) {
	tester := buildClientTester(t)

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	// the first GetBlob will return a 0 len blob
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return([]byte{}, nil).Once()
	// the second call will return blob bytes
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(
		blobBytes,
		nil).Once()

	// the call to the first relay will fail with a 0 len blob returned. the call to the second relay will succeed
	blob, err := tester.Client.GetPayload(context.Background(), blobKey, blobCert)
	require.NotNil(t, blob)
	require.NoError(t, err)

	tester.requireExpectations(t)
}

// TestFailedDecoding verifies that a failed blob decode is handled gracefully, and that the client retries after such a failure
func TestFailedDecoding(t *testing.T) {
	tester := buildClientTester(t)

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	// payloadLength random bytes are guaranteed to be an invalid blob
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(
		tester.Random.Bytes(payloadLength),
		nil).Once()
	// the second call will return blob bytes
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(
		blobBytes,
		nil).Once()

	// decoding will fail the first time, but succeed the second time
	blob, err := tester.Client.GetPayload(context.Background(), blobKey, blobCert)
	require.NotNil(t, blob)
	require.NoError(t, err)

	tester.requireExpectations(t)
}

// TestErrorFreeClose tests the happy case, where none of the internal closes yield an error
func TestErrorFreeClose(t *testing.T) {
	tester := buildClientTester(t)

	tester.MockRelayClient.On("Close").Return(nil).Once()

	err := tester.Client.Close()
	require.NoError(t, err)

	tester.requireExpectations(t)
}

// TestErrorClose tests what happens when subcomponents throw errors when being closed
func TestErrorClose(t *testing.T) {
	tester := buildClientTester(t)

	tester.MockRelayClient.On("Close").Return(fmt.Errorf("close failed")).Once()

	err := tester.Client.Close()
	require.NotNil(t, err)

	tester.requireExpectations(t)
}

// TestGetCodec checks that the codec used in construction is returned by GetCodec
func TestGetCodec(t *testing.T) {
	tester := buildClientTester(t)

	require.Equal(t, tester.Codec, tester.Client.GetCodec())

	tester.requireExpectations(t)
}

// TestBuilder tests that the method that builds the client from config doesn't throw any obvious errors
func TestBuilder(t *testing.T) {
	clientConfig := &clients.EigenDAClientConfig{
		BlobEncodingVersion: codecs.DefaultBlobEncoding,
		BlobPolynomialForm:  codecs.Coeff,
		RelayTimeout:        500 * time.Millisecond,
	}

	sockets := make(map[core.RelayKey]string)
	sockets[core.RelayKey(44)] = "socketVal"

	relayClientConfig := &clients.RelayClientConfig{
		Sockets:           sockets,
		UseSecureGrpcFlag: true,
	}

	client, err := clients.BuildEigenDAClient(
		logging.NewNoopLogger(),
		clientConfig,
		relayClientConfig,
		[]bn254.G1Affine{})

	require.NotNil(t, client)
	require.NoError(t, err)

	require.NotNil(t, client.GetCodec())
}
