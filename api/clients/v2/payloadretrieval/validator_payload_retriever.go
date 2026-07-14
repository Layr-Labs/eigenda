package payloadretrieval

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/math"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// ValidatorPayloadRetriever provides the ability to get payloads from the EigenDA validator nodes directly
//
// This struct is goroutine safe.
type ValidatorPayloadRetriever struct {
	logger          logging.Logger
	config          ValidatorPayloadRetrieverConfig
	retrievalClient validator.ValidatorClient
	g1Srs           []bn254.G1Affine
	metrics         metrics.RetrievalMetricer
}

var _ clients.PayloadRetriever = &ValidatorPayloadRetriever{}

// NewValidatorPayloadRetriever creates a new ValidatorPayloadRetriever from already constructed objects
func NewValidatorPayloadRetriever(
	logger logging.Logger,
	config ValidatorPayloadRetrieverConfig,
	retrievalClient validator.ValidatorClient,
	g1Srs []bn254.G1Affine,
	metrics metrics.RetrievalMetricer,
) (*ValidatorPayloadRetriever, error) {
	err := config.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set config defaults: %w", err)
	}

	return &ValidatorPayloadRetriever{
		logger:          logger,
		config:          config,
		retrievalClient: retrievalClient,
		g1Srs:           g1Srs,
		metrics:         metrics,
	}, nil
}

// GetPayload iteratively attempts to retrieve a given blob from the quorums listed in the EigenDACert.
//
// If the blob is successfully retrieved, then the blob verified against the EigenDACert. If the verification succeeds,
// the blob is decoded to yield the payload (the original user data, with no padding or any modification), and the
// payload is returned.
//
// This method does NOT verify the eigenDACert on chain: it is assumed that the input eigenDACert has already been
// verified prior to calling this method.
func (pr *ValidatorPayloadRetriever) GetPayload(
	ctx context.Context,
	eigenDACert coretypes.EigenDACert,
) (coretypes.Payload, error) {

	encodedPayload, err := pr.GetEncodedPayload(ctx, eigenDACert)
	if err != nil {
		return nil, err
	}

	payload, err := encodedPayload.Decode()
	if err != nil {
		// If we successfully compute the blob key, we add it to the error message to help with debugging.
		blobKey, keyErr := eigenDACert.ComputeBlobKey()
		if keyErr == nil {
			err = fmt.Errorf("blob %v: %w", blobKey.Hex(), err)
		}
		return nil, coretypes.ErrBlobDecodingFailedDerivationError.WithMessage(err.Error())
	}

	pr.metrics.RecordPayloadSizeBytes(len(payload))

	return payload, nil
}

// GetEncodedPayload iteratively attempts to retrieve a given blob from the quorums
// listed in the EigenDACert.
//
// If the blob is successfully retrieved, then the blob is verified against the EigenDACert.
// If the verification succeeds, the blob is converted to an encoded payload form and returned.
//
// This method does NOT verify the eigenDACert on chain: it is assumed that the input
// eigenDACert has already been verified prior to calling this method.
func (pr *ValidatorPayloadRetriever) GetEncodedPayload(
	ctx context.Context,
	eigenDACert coretypes.EigenDACert,
) (*coretypes.EncodedPayload, error) {

	blobHeader, err := eigenDACert.BlobHeader()
	if err != nil {
		return nil, fmt.Errorf("get blob header from eigenDACert: %w", err)
	}

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key from eigenDACert: %w", err)
	}

	blobLengthSymbols := uint32(blobHeader.BlobCommitments.Length)
	// TODO(samlaf): are there more properties of the Cert that should lead to [coretypes.MaliciousOperatorsError]s?
	if !math.IsPowerOfTwo(blobLengthSymbols) {
		return nil, coretypes.ErrCertCommitmentBlobLengthNotPowerOf2MaliciousOperatorsError.WithBlobKey(blobKey.Hex())
	}

	// TODO (litt3): Add a feature which keeps chunks from previous quorums, and just fills in gaps
	for _, quorumID := range blobHeader.QuorumNumbers {
		blob, err := pr.retrieveBlobWithTimeout(
			ctx,
			blobHeader,
			uint32(eigenDACert.ReferenceBlockNumber()))

		if err != nil {
			pr.logger.Error(
				"blob couldn't be retrieved from quorum",
				"blobKey", blobKey.Hex(),
				"quorumId", quorumID,
				"error", err)
			continue
		}

		valid, err := verification.GenerateAndCompareBlobCommitment(
			pr.g1Srs, blob, blobHeader.BlobCommitments.Commitment)
		if err != nil {
			pr.logger.Warn(
				"generate and compare blob commitment",
				"blobKey", blobKey.Hex(), "quorumID", quorumID, "error", err)
			continue
		}
		if !valid {
			pr.logger.Warn(
				"generated commitment doesn't match cert commitment",
				"blobKey", blobKey.Hex(), "quorumID", quorumID)
			continue
		}

		return blob.ToEncodedPayloadUnchecked(pr.config.PayloadPolynomialForm), nil
	}

	return nil, fmt.Errorf("unable to retrieve encoded payload with blobKey %v from quorums %v", blobKey.Hex(), blobHeader.QuorumNumbers)
}

// retrieveBlobWithTimeout attempts to retrieve a blob from a given quorum,
// and times out based on [ValidatorPayloadRetrieverConfig.RetrievalTimeout].
//
// blobLengthSymbols MUST be taken from the eigenDACert for the blob being retrieved,
// and MUST be a power of 2.
func (pr *ValidatorPayloadRetriever) retrieveBlobWithTimeout(
	ctx context.Context,
	header *corev2.BlobHeaderWithHashedPayment,
	referenceBlockNumber uint32,
) (*coretypes.Blob, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RetrievalTimeout)
	defer cancel()

	// TODO (litt3): eventually, we should make GetBlob return an actual blob object, instead of the serialized bytes.
	blobBytes, err := pr.retrievalClient.GetBlob(
		timeoutCtx,
		header,
		uint64(referenceBlockNumber),
	)

	if err != nil {
		return nil, fmt.Errorf("get blob: %w", err)
	}

	blob, err := coretypes.DeserializeBlob(blobBytes, uint32(header.BlobCommitments.Length))
	if errors.Is(err, coretypes.ErrBlobLengthSymbolsNotPowerOf2) {
		// In a better language I would write this as a debug assert.
		pr.logger.Error("BROKEN INVARIANT: retrieveBlobWithTimeout: blobLengthSymbols is not power of 2: "+
			"this is a major broken invariant, that should have been checked by the validators, "+
			"and the caller (GetEncodedPayload) should already have checked this invariant "+
			"and returned a MaliciousOperatorsError. Returning the same MaliciousOperatorsError "+
			"to be safe, but this code should be fixed.", "err", err)
		blobKey, _ := header.BlobKey() // discard error since returning the below error is most important
		return nil, coretypes.ErrCertCommitmentBlobLengthNotPowerOf2MaliciousOperatorsError.WithBlobKey(blobKey.Hex())
	}
	if err != nil {
		return nil, fmt.Errorf("deserialize blob: %w", err)
	}

	return blob, nil
}
