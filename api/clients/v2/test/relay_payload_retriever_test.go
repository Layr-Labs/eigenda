package test

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	contractEigenDACertVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	prover2 "github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const g1Path = "../../../../inabox/resources/kzg/g1.point"
const payloadLength = 100

type RelayPayloadRetrieverTester struct {
	Random                *testrandom.TestRandom
	RelayPayloadRetriever *clients.RelayPayloadRetriever
	MockRelayClient       *clientsmock.MockRelayClient
	Codec                 *codecs.DefaultBlobCodec
	G1Srs                 []bn254.G1Affine
}

// buildRelayPayloadRetrieverTester sets up a client with mocks necessary for testing
func buildRelayPayloadRetrieverTester(t *testing.T) RelayPayloadRetrieverTester {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	// the constructor checks that these values aren't empty. we don't need them, though, since we're using mocks
	payloadClientConfig := clients.PayloadClientConfig{
		EigenDACertVerifierAddr: "x",
	}

	clientConfig := clients.RelayPayloadRetrieverConfig{
		PayloadClientConfig: payloadClientConfig,
		RelayTimeout:        50 * time.Millisecond,
	}

	mockRelayClient := clientsmock.MockRelayClient{}
	codec := codecs.NewDefaultBlobCodec()

	random := testrandom.NewTestRandom(t)

	g1Srs, err := kzg.ReadG1Points(g1Path, 5, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	client, err := clients.NewRelayPayloadRetriever(
		logger,
		random.Rand,
		clientConfig,
		&mockRelayClient,
		&codec,
		g1Srs)

	require.NotNil(t, client)
	require.NoError(t, err)

	return RelayPayloadRetrieverTester{
		Random:                random,
		RelayPayloadRetriever: client,
		MockRelayClient:       &mockRelayClient,
		Codec:                 &codec,
		G1Srs:                 g1Srs,
	}
}

// Builds a random blob key, blob bytes, and valid certificate
func buildBlobAndCert(
	t *testing.T,
	tester RelayPayloadRetrieverTester,
	relayKeys []core.RelayKey,
) (core.BlobKey, []byte, *verification.EigenDACert) {

	payloadBytes := tester.Random.Bytes(payloadLength)
	blobBytes, err := tester.Codec.EncodeBlob(payloadBytes)
	require.NoError(t, err)
	require.NotNil(t, blobBytes)

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../../../../inabox/resources/kzg/g1.point",
		G2Path:          "../../../../inabox/resources/kzg/g2.point",
		CacheDir:        "../../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	prover, err := prover2.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	commitments, err := prover.GetCommitmentsForPaddedLength(blobBytes)
	require.NoError(t, err)

	commitmentsProto, err := commitments.ToProtobuf()
	require.NoError(t, err)

	blobHeader := &commonv2.BlobHeader{
		Version:       1,
		QuorumNumbers: make([]uint32, 0),
		PaymentHeader: &commonv2.PaymentHeader{},
		Commitment:    commitmentsProto,
	}

	blobCertificate := &commonv2.BlobCertificate{
		RelayKeys:  relayKeys,
		BlobHeader: blobHeader,
	}

	inclusionInfo := &disperserv2.BlobInclusionInfo{
		BlobCertificate: blobCertificate,
	}

	convertedInclusionInfo, err := verification.InclusionInfoProtoToBinding(inclusionInfo)
	require.NoError(t, err)

	eigenDACert := &verification.EigenDACert{
		BlobInclusionInfo: *convertedInclusionInfo,
	}

	blobKey, err := eigenDACert.ComputeBlobKey()
	require.NoError(t, err)

	return *blobKey, blobBytes, eigenDACert
}

// TestGetPayloadSuccess tests that a blob is received without error in the happy case
func TestGetPayloadSuccess(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)
	relayKeys := make([]core.RelayKey, 1)
	relayKeys[0] = tester.Random.Uint32()
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(blobBytes, nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)

	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestRelayCallTimeout verifies that calls to the relay timeout after the expected duration
func TestRelayCallTimeout(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)
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
			_, _ = tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
		})

	require.Panics(
		t, func() {
			_, _ = tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
		})

	tester.MockRelayClient.AssertExpectations(t)
}

// TestRandomRelayRetries verifies correct behavior when some relays do not respond with the blob,
// requiring the PayloadRetriever to retry with other relays.
func TestRandomRelayRetries(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

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
		payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
		require.NotNil(t, payload)
		require.NoError(t, err)

		requiredTries[failedCallCount] = true
	}

	// with 100 random tries, with possible values between 1 and 100, we can very confidently require that there are at least 10 unique values
	require.Greater(t, len(requiredTries), 10)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestNoRelayResponse tests functionality when none of the relays respond
func TestNoRelayResponse(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, _, blobCert := buildBlobAndCert(t, tester, relayKeys)

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(nil, fmt.Errorf("offline relay"))

	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
	require.Nil(t, payload)
	require.NotNil(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestNoRelays tests that having no relay keys is handled gracefully
func TestNoRelays(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	_, _, blobCert := buildBlobAndCert(t, tester, []core.RelayKey{})

	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
	require.Nil(t, payload)
	require.NotNil(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturns0Len verifies that a 0 length blob returned from a relay is handled gracefully, and that the PayloadRetriever retries after such a failure
func TestGetBlobReturns0Len(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

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
	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturnsDifferentBlob tests what happens when one relay returns a blob that doesn't match the commitment.
// It also tests that the PayloadRetriever retries to get the correct blob from a different relay
func TestGetBlobReturnsDifferentBlob(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)
	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey1, blobBytes1, blobCert1 := buildBlobAndCert(t, tester, relayKeys)
	_, blobBytes2, _ := buildBlobAndCert(t, tester, relayKeys)

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey1).Return(blobBytes2, nil).Once()
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey1).Return(blobBytes1, nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert1)
	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturnsInvalidBlob tests what happens if a relay returns a blob which causes commitment verification to
// throw an error. It verifies that the PayloadRetriever tries again with a different relay after such a failure.
func TestGetBlobReturnsInvalidBlob(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)
	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobKey, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	tooLongBytes := make([]byte, len(blobBytes)+100)
	copy(tooLongBytes[:], blobBytes)

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(tooLongBytes, nil).Once()
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(blobBytes, nil).Once()

	// this will fail the first time, since there isn't enough srs loaded to compute the commitment of the returned bytes
	// it will succeed when the second relay gives the correct bytes
	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)

	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturnsBlobWithInvalidLen check what happens if the blob length doesn't match the length that exists in
// the BlobCommitment
func TestGetBlobReturnsBlobWithInvalidLen(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	relayKeys := make([]core.RelayKey, 1)
	relayKeys[0] = tester.Random.Uint32()

	_, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	blobCert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length = (uint32(len(blobBytes)) / 32) - 1

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, mock.Anything).Return(blobBytes, nil).Once()

	// this will fail, since the length in the BlobCommitment doesn't match the actual blob length
	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)

	require.Nil(t, payload)
	require.Error(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestFailedDecoding verifies that a failed blob decode is handled gracefully
func TestFailedDecoding(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	_, blobBytes, blobCert := buildBlobAndCert(t, tester, relayKeys)

	// intentionally cause the payload header claimed length to differ from the actual length
	binary.BigEndian.PutUint32(blobBytes[2:6], uint32(len(blobBytes)-1))

	// generate a malicious cert, which will verify for the invalid blob
	maliciousCommitment, err := verification.GenerateBlobCommitment(tester.G1Srs, blobBytes)
	require.NoError(t, err)
	require.NotNil(t, maliciousCommitment)

	blobCert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment = contractEigenDACertVerifier.BN254G1Point{
		X: maliciousCommitment.X.BigInt(new(big.Int)),
		Y: maliciousCommitment.Y.BigInt(new(big.Int)),
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, mock.Anything).Return(
		blobBytes,
		nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(context.Background(), blobCert)
	require.Error(t, err)
	require.Nil(t, payload)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestErrorFreeClose tests the happy case, where none of the internal closes yield an error
func TestErrorFreeClose(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	tester.MockRelayClient.On("Close").Return(nil).Once()

	err := tester.RelayPayloadRetriever.Close()
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestErrorClose tests what happens when subcomponents throw errors when being closed
func TestErrorClose(t *testing.T) {
	tester := buildRelayPayloadRetrieverTester(t)

	tester.MockRelayClient.On("Close").Return(fmt.Errorf("close failed")).Once()

	err := tester.RelayPayloadRetriever.Close()
	require.NotNil(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}
