package payloadretrieval

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	disperserv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/math"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/test"
	testrandom "github.com/Layr-Labs/eigenda/test/random"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	maxPayloadBytes = 1025 // arbitrary value
	g1Path          = "../../../../resources/srs/g1.point"
)

type RelayPayloadRetrieverTester struct {
	Random                *testrandom.TestRandom
	RelayPayloadRetriever *RelayPayloadRetriever
	MockRelayClient       *clientsmock.MockRelayClient
	G1Srs                 []bn254.G1Affine
}

func (t *RelayPayloadRetrieverTester) PayloadPolynomialForm() codecs.PolynomialForm {
	return t.RelayPayloadRetriever.config.PayloadPolynomialForm
}

// buildRelayPayloadRetrieverTester sets up a client with mocks necessary for testing
func buildRelayPayloadRetrieverTester(t *testing.T) RelayPayloadRetrieverTester {
	logger := test.GetLogger()

	clientConfig := RelayPayloadRetrieverConfig{
		PayloadClientConfig: clients.PayloadClientConfig{},
		RelayTimeout:        50 * time.Millisecond,
	}

	mockRelayClient := clientsmock.MockRelayClient{}
	random := testrandom.NewTestRandom()

	srsPointsToLoad := math.NextPowOf2u32(codec.GetPaddedDataLength(maxPayloadBytes)) / encoding.BYTES_PER_SYMBOL

	g1Srs, err := kzg.ReadG1Points(g1Path, uint64(srsPointsToLoad), uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	client, err := NewRelayPayloadRetriever(
		logger,
		clientConfig,
		&mockRelayClient,
		g1Srs,
		metrics.NoopRetrievalMetrics)

	require.NotNil(t, client)
	require.NoError(t, err)

	return RelayPayloadRetrieverTester{
		Random:                random,
		RelayPayloadRetriever: client,
		MockRelayClient:       &mockRelayClient,
		G1Srs:                 g1Srs,
	}
}

// Builds a random blob and valid certificate
func buildBlobAndCert(
	t *testing.T,
	tester RelayPayloadRetrieverTester,
) (*coretypes.Blob, *coretypes.EigenDACertV3) {

	payloadBytes := tester.Random.Bytes(tester.Random.Intn(maxPayloadBytes))
	blob, err := coretypes.Payload(payloadBytes).ToBlob(tester.PayloadPolynomialForm())
	require.NoError(t, err)
	cert := buildCertFromBlobBytes(t, blob.Serialize(), tester.Random.Uint32())
	return blob, cert
}

// Builds a valid certificate from the given blob bytes.
// It is used to generate a valid cert from a wrongly encoded blob, to test for decoding errors.
func buildCertFromBlobBytes(
	t *testing.T,
	blobBytes []byte,
	relayKey core.RelayKey,
) *coretypes.EigenDACertV3 {

	committerConfig := committer.Config{
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
		SRSNumberToLoad:   4096,
	}

	committer, err := committer.NewFromConfig(committerConfig)
	require.NoError(t, err)

	commitments, err := committer.GetCommitmentsForPaddedLength(blobBytes)
	require.NoError(t, err)

	commitmentsProto, err := commitments.ToProtobuf()
	require.NoError(t, err)

	blobHeader := &commonv2.BlobHeader{
		Version:       1,
		QuorumNumbers: make([]uint32, 0),
		PaymentHeader: &commonv2.PaymentHeader{
			AccountId: gethcommon.Address{1}.Hex(),
		},
		Commitment: commitmentsProto,
	}

	blobCertificate := &commonv2.BlobCertificate{
		RelayKeys:  []core.RelayKey{relayKey},
		BlobHeader: blobHeader,
	}

	inclusionInfo := &disperserv2.BlobInclusionInfo{
		BlobCertificate: blobCertificate,
	}

	convertedInclusionInfo, err := coretypes.InclusionInfoProtoToIEigenDATypesBinding(inclusionInfo)
	require.NoError(t, err)

	return &coretypes.EigenDACertV3{
		BlobInclusionInfo: *convertedInclusionInfo,
	}
}

// TestGetPayloadSuccess tests that a blob is received without error in the happy case
func TestGetPayloadSuccess(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	blob, blobCert := buildBlobAndCert(t, tester)

	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(blob, nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)

	require.NotNil(t, payload)
	require.NoError(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestRelayCallTimeout verifies that calls to the relay timeout after the expected duration
func TestRelayCallTimeout(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	_, blobCert := buildBlobAndCert(t, tester)

	// the timeout should occur before the panic has a chance to be triggered
	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(
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

	// the panic should be triggered, since it happens faster than the configured timeout
	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(
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
			_, _ = tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)
		})

	require.Panics(
		t, func() {
			_, _ = tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)
		})

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturnsError tests that errors from GetBlob are propagated correctly
func TestGetBlobReturnsError(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	_, blobCert := buildBlobAndCert(t, tester)

	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(nil, fmt.Errorf("relay error"))

	payload, err := tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)
	require.Nil(t, payload)
	require.NotNil(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestGetBlobReturnsDifferentBlob tests that when the relay returns a blob that doesn't match the commitment,
// an error is returned.
func TestGetBlobReturnsDifferentBlob(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	_, blobCert := buildBlobAndCert(t, tester)
	wrongBlob, _ := buildBlobAndCert(t, tester)

	// Return a wrong blob that doesn't match the cert commitment
	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(wrongBlob, nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)
	require.Nil(t, payload)
	require.Error(t, err)

	tester.MockRelayClient.AssertExpectations(t)
}

// TestFailedDecoding verifies that decoding errors (caused by corrupted payload headers) are handled gracefully.
func TestFailedDecoding(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	blob, _ := buildBlobAndCert(t, tester)
	blobBytes := blob.Serialize()

	// Corrupt the blob bytes to have an invalid payload header length
	binary.BigEndian.PutUint32(blobBytes[2:6], uint32(len(blobBytes)-1))

	// Build a cert that matches the corrupted blob so commitment verification passes
	maliciousCert := buildCertFromBlobBytes(t, blobBytes, tester.Random.Uint32())
	blobLengthSymbols := maliciousCert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length
	maliciousBlob, err := coretypes.DeserializeBlob(blobBytes, blobLengthSymbols)
	require.NoError(t, err)

	// The mock returns this malicious blob, which passes commitment verification but fails decoding
	tester.MockRelayClient.On("GetBlob", mock.Anything, maliciousCert).Return(maliciousBlob, nil).Once()

	payload, err := tester.RelayPayloadRetriever.GetPayload(ctx, maliciousCert)
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

// TestCommitmentVerifiesButBlobToPayloadFails tests the case where commitment verification succeeds
// but conversion from blob to payload fails. This is a critical edge case that should not be possible
// with valid data, but could indicate malicious dispersed data.
func TestCommitmentVerifiesButBlobToPayloadFails(t *testing.T) {
	ctx := t.Context()

	tester := buildRelayPayloadRetrieverTester(t)
	// We keep the blob in coeff form so that we can manipulate it directly (otherwise it gets IFFT'd)
	tester.RelayPayloadRetriever.config.PayloadPolynomialForm = codecs.PolynomialFormCoeff

	payloadBytes := tester.Random.Bytes(tester.Random.Intn(maxPayloadBytes))
	blob, err := coretypes.Payload(payloadBytes).ToBlob(tester.PayloadPolynomialForm())
	require.NoError(t, err)
	blobBytes := blob.Serialize()
	require.NotNil(t, blobBytes)
	blobBytes[1] = 0xFF // Invalid encoding version - this will cause decode to fail

	blobCert := buildCertFromBlobBytes(t, blobBytes, tester.Random.Uint32())
	blobLengthSymbols := blobCert.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length
	maliciousBlob, err := coretypes.DeserializeBlob(blobBytes, blobLengthSymbols)
	require.NoError(t, err)

	// Mock the relay to return our incorrectly encoded blob
	tester.MockRelayClient.On("GetBlob", mock.Anything, blobCert).Return(maliciousBlob, nil).Once()

	// Try to get the payload - this should fail during blob to payload conversion
	payload, err := tester.RelayPayloadRetriever.GetPayload(ctx, blobCert)
	require.Nil(t, payload)
	require.Error(t, err)

	// Verify it's specifically a DerivationError with status code 4 (blob decoding failed)
	derivationErr := coretypes.DerivationError{}
	require.ErrorAs(t, err, &derivationErr)
	require.Equal(t, coretypes.ErrBlobDecodingFailedDerivationError.StatusCode, derivationErr.StatusCode)

	tester.MockRelayClient.AssertExpectations(t)
}
