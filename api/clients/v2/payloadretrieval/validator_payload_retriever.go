package payloadretrieval

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
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
	metrics         metrics.ClientMetricer
}

var _ clients.PayloadRetriever = &ValidatorPayloadRetriever{}

// NewValidatorPayloadRetriever creates a new ValidatorPayloadRetriever from already constructed objects
func NewValidatorPayloadRetriever(
	logger logging.Logger,
	config ValidatorPayloadRetrieverConfig,
	retrievalClient validator.ValidatorClient,
	g1Srs []bn254.G1Affine,
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
	eigenDACert coretypes.RetrievableEigenDACert,
) (*coretypes.Payload, error) {
	encodedPayload, err := pr.GetEncodedPayload(ctx, eigenDACert)
	if err != nil {
		if pr.metrics != nil {
			pr.metrics.ReportValidatorPayloadRetrieval("failure", "GetPayload", 0, 0)
		}
		return nil, err
	}

	payload, err := encodedPayload.Decode()
	if err != nil {
		// If we successfully compute the blob key, we add it to the error message to help with debugging.
		blobKey, keyErr := eigenDACert.ComputeBlobKey()
		if keyErr == nil {
			err = fmt.Errorf("blob %v: %w", blobKey.Hex(), err)
			if pr.metrics != nil {
				pr.metrics.ReportValidatorPayloadDecodeFailure(blobKey.Hex())
			}
		}
		if pr.metrics != nil {
			pr.metrics.ReportValidatorPayloadRetrieval("failure", "GetPayload", 0, 0)
		}
		return nil, coretypes.ErrBlobDecodingFailedDerivationError.WithMessage(err.Error())
	}

	if pr.metrics != nil {
		payloadSize := uint64(len(payload.Serialize()))
		pr.metrics.ReportValidatorPayloadRetrieval("success", "GetPayload", payloadSize, 0)
	}

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
	eigenDACert coretypes.RetrievableEigenDACert,
) (*coretypes.EncodedPayload, error) {
	var quorumAttempts int

	blobHeader, err := eigenDACert.BlobHeader()
	if err != nil {
		return nil, fmt.Errorf("get blob header from eigenDACert: %w", err)
	}

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key from eigenDACert: %w", err)
	}

	// TODO (litt3): Add a feature which keeps chunks from previous quorums, and just fills in gaps
	for _, quorumID := range blobHeader.QuorumNumbers {
		quorumAttempts++
		quorumIDStr := fmt.Sprintf("%d", quorumID)

		blob, err := pr.retrieveBlobWithTimeout(
			ctx,
			blobHeader,
			uint32(eigenDACert.ReferenceBlockNumber()),
			quorumIDStr)

		if err != nil {
			pr.logger.Error(
				"blob couldn't be retrieved from quorum",
				"blobKey", blobKey.Hex(),
				"quorumId", quorumID,
				"error", err)
			continue
		}

		// TODO (litt3): eventually, we should make GenerateAndCompareBlobCommitment accept a blob, instead of the
		//  serialization of a blob. Commitment generation operates on field elements, which is how a blob is stored
		//  under the hood, so it's actually duplicating work to serialize the blob here. I'm declining to make this
		//  change now, to limit the size of the refactor PR.
		valid, err := verification.GenerateAndCompareBlobCommitment(pr.g1Srs, blob.Serialize(), blobHeader.BlobCommitments.Commitment)
		if err != nil {
			pr.logger.Warn(
				"generate and compare blob commitment",
				"blobKey", blobKey.Hex(), "quorumID", quorumID, "error", err)
			if pr.metrics != nil {
				pr.metrics.ReportValidatorCommitmentVerificationFailure(quorumIDStr)
			}
			continue
		}
		if !valid {
			pr.logger.Warn(
				"generated commitment doesn't match cert commitment",
				"blobKey", blobKey.Hex(), "quorumID", quorumID)
			if pr.metrics != nil {
				pr.metrics.ReportValidatorCommitmentVerificationFailure(quorumIDStr)
			}
			continue
		}

		encodedPayload, err := blob.ToEncodedPayload(pr.config.PayloadPolynomialForm)
		if err != nil {
			// TODO(samlaf): ToEncodedPayload is doing too much decoding. It shouldn't read and validate the payload header.
			// That needs to be left to the rollup's derivation pipeline, such that a failed decoding can be skipped safely.
			// A lot of the logic in blob->encodedPayload prob needs to happen in encodedPayload->payload instead.
			return nil, fmt.Errorf("convert blob to encoded payload failed."+
				" blobKey: %s, quorumID: %v, error: %v", blobKey.Hex(), quorumID, err)
		}

		// Report successful retrieval
		if pr.metrics != nil {
			payloadSize := uint64(len(encodedPayload.Serialize()))
			pr.metrics.ReportValidatorPayloadRetrieval("success", "GetEncodedPayload", payloadSize, quorumAttempts)
		}

		return encodedPayload, nil
	}

	// Report failed retrieval
	if pr.metrics != nil {
		pr.metrics.ReportValidatorPayloadRetrieval("failure", "GetEncodedPayload", 0, quorumAttempts)
	}

	return nil, fmt.Errorf("unable to retrieve encoded payload with blobKey %v from quorums %v", blobKey.Hex(), blobHeader.QuorumNumbers)
}

// retrieveBlobWithTimeout attempts to retrieve a blob from a given quorum, and times out based on config.RetrievalTimeout
func (pr *ValidatorPayloadRetriever) retrieveBlobWithTimeout(
	ctx context.Context,
	header *corev2.BlobHeaderWithHashedPayment,
	referenceBlockNumber uint32,
	quorumIDStr string,
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
		if pr.metrics != nil {
			// Check if it's a timeout error
			if ctx.Err() == context.DeadlineExceeded || timeoutCtx.Err() == context.DeadlineExceeded {
				pr.metrics.ReportValidatorTimeout(quorumIDStr)
			} else {
				pr.metrics.ReportValidatorConnectionError(quorumIDStr)
			}
			pr.metrics.ReportValidatorQuorumRequest(quorumIDStr, "failure")
		}
		return nil, fmt.Errorf("get blob: %w", err)
	}

	blob, err := coretypes.DeserializeBlob(blobBytes, uint32(header.BlobCommitments.Length))
	if err != nil {
		if pr.metrics != nil {
			pr.metrics.ReportValidatorBlobDeserializationFailure(quorumIDStr)
			pr.metrics.ReportValidatorQuorumRequest(quorumIDStr, "failure")
		}
		return nil, fmt.Errorf("deserialize blob: %w", err)
	}

	// Report successful quorum request
	if pr.metrics != nil {
		pr.metrics.ReportValidatorQuorumRequest(quorumIDStr, "success")
	}

	return blob, nil
}

func (pr *ValidatorPayloadRetriever) SetMetrics(metrics metrics.ClientMetricer) {
	pr.metrics = metrics
}
