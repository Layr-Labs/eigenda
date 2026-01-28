package payloadretrieval

import (
	"context"
	"fmt"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// RelayPayloadRetriever provides the ability to get payloads from the relay subsystem.
//
// This struct is goroutine safe.
type RelayPayloadRetriever struct {
	log         logging.Logger
	config      RelayPayloadRetrieverConfig
	relayClient relay.RelayClient
	g1Srs       []bn254.G1Affine
	metrics     metrics.RetrievalMetricer
}

var _ clients.PayloadRetriever = &RelayPayloadRetriever{}

// NewRelayPayloadRetriever assembles a RelayPayloadRetriever from subcomponents that have already been constructed and
// initialized.
func NewRelayPayloadRetriever(
	log logging.Logger,
	relayPayloadRetrieverConfig RelayPayloadRetrieverConfig,
	relayClient relay.RelayClient,
	g1Srs []bn254.G1Affine,
	metrics metrics.RetrievalMetricer) (*RelayPayloadRetriever, error) {

	err := relayPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set RelayPayloadRetrieverConfig config: %w", err)
	}

	return &RelayPayloadRetriever{
		log:         log,
		config:      relayPayloadRetrieverConfig,
		relayClient: relayClient,
		g1Srs:       g1Srs,
		metrics:     metrics,
	}, nil
}

// GetPayload retrieves a blob from the relay specified in the EigenDACert.
//
// If the blob is successfully retrieved, then the blob is verified against the certificate. If the verification
// succeeds, the blob is decoded to yield the payload (the original user data, with no padding or any modification),
// and the payload is returned.
//
// This method does NOT verify the eigenDACert on chain: it is assumed that the input eigenDACert has already been
// verified prior to calling this method.
func (pr *RelayPayloadRetriever) GetPayload(
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

// GetEncodedPayload retrieves a blob from the relay specified in the EigenDACert.
//
// If the blob is successfully retrieved, then the blob is verified against the EigenDACert.
// If the verification succeeds, the blob is converted to an encoded payload form and returned.
//
// This method does NOT verify the eigenDACert on chain: it is assumed that the input
// eigenDACert has already been verified prior to calling this method.
func (pr *RelayPayloadRetriever) GetEncodedPayload(
	ctx context.Context,
	eigenDACert coretypes.EigenDACert,
) (*coretypes.EncodedPayload, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobCommitments, err := eigenDACert.Commitments()
	if err != nil {
		return nil, fmt.Errorf("blob %s: get commitments from eigenDACert: %w", blobKey.Hex(), err)
	}

	// TODO(samlaf): are there more properties of the Cert that should lead to [coretypes.MaliciousOperatorsError]s?
	if !math.IsPowerOfTwo(blobCommitments.Length) {
		return nil, coretypes.ErrCertCommitmentBlobLengthNotPowerOf2MaliciousOperatorsError.WithBlobKey(blobKey.Hex())
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RelayTimeout)
	defer cancel()
	blob, err := pr.relayClient.GetBlob(timeoutCtx, eigenDACert)
	if err != nil {
		return nil, fmt.Errorf("blob %s: get blob from relay: %w", blobKey.Hex(), err)
	}

	// TODO (litt3): eventually, we should make GenerateAndCompareBlobCommitment accept a blob, instead of the
	//  serialization of a blob. Commitment generation operates on field elements, which is how a blob is stored
	//  under the hood, so it's actually duplicating work to serialize the blob here. I'm declining to make this
	//  change now, to limit the size of the refactor PR.
	valid, err := verification.GenerateAndCompareBlobCommitment(pr.g1Srs, blob.Serialize(), blobCommitments.Commitment)
	if err != nil {
		return nil, fmt.Errorf("blob %s: generate and compare blob commitment: %w", blobKey.Hex(), err)
	}
	if !valid {
		return nil, fmt.Errorf("blob %s: commitment mismatch with cert", blobKey.Hex())
	}

	return blob.ToEncodedPayloadUnchecked(pr.config.PayloadPolynomialForm), nil
}

// Close is responsible for calling close on all internal clients. This method will do its best to close all internal
// clients, even if some closes fail.
//
// Any and all errors returned from closing internal clients will be joined and returned.
//
// This method should only be called once.
func (pr *RelayPayloadRetriever) Close() error {
	err := pr.relayClient.Close()
	if err != nil {
		return fmt.Errorf("close relay client: %w", err)
	}

	return nil
}
