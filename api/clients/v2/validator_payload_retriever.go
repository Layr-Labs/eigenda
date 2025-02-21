package clients

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ValidatorPayloadRetriever provides the ability to get payloads from the EigenDA validator nodes directly
//
// This struct is goroutine safe.
type ValidatorPayloadRetriever struct {
	logger          logging.Logger
	config          ValidatorPayloadRetrieverConfig
	retrievalClient RetrievalClient
	g1Srs           []bn254.G1Affine
}

var _ PayloadRetriever = &ValidatorPayloadRetriever{}

// BuildValidatorPayloadRetriever builds a ValidatorPayloadRetriever from config structs.
func BuildValidatorPayloadRetriever(
	logger logging.Logger,
	validatorPayloadRetrieverConfig ValidatorPayloadRetrieverConfig,
	ethConfig geth.EthClientConfig,
	thegraphConfig thegraph.Config,
	kzgConfig kzg.KzgConfig,
) (*ValidatorPayloadRetriever, error) {
	err := validatorPayloadRetrieverConfig.checkAndSetDefaults()
	if err != nil {
		return nil, fmt.Errorf("check and set config defaults: %w", err)
	}

	ethClient, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	reader, err := eth.NewReader(
		logger,
		ethClient,
		validatorPayloadRetrieverConfig.BlsOperatorStateRetrieverAddr,
		validatorPayloadRetrieverConfig.EigenDAServiceManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	chainState := eth.NewChainState(reader, ethClient)
	indexedChainState := thegraph.MakeIndexedChainState(thegraphConfig, chainState, logger)

	kzgVerifier, err := verifier.NewVerifier(&kzgConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("new verifier: %w", err)
	}

	retrievalClient := NewRetrievalClient(
		logger,
		reader,
		indexedChainState,
		kzgVerifier,
		int(validatorPayloadRetrieverConfig.MaxConnectionCount))

	return &ValidatorPayloadRetriever{
		logger:          logger,
		config:          validatorPayloadRetrieverConfig,
		retrievalClient: retrievalClient,
		g1Srs:           kzgVerifier.Srs.G1,
	}, nil
}

// NewValidatorPayloadRetriever creates a new ValidatorPayloadRetriever from already constructed objects
func NewValidatorPayloadRetriever(
	logger logging.Logger,
	config ValidatorPayloadRetrieverConfig,
	retrievalClient RetrievalClient,
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
func (pr *ValidatorPayloadRetriever) GetPayload(
	ctx context.Context,
	eigenDACert *verification.EigenDACert,
) (*codecs.Payload, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobHeader := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader
	commitment, err := verification.BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return nil, fmt.Errorf("convert commitments binding to internal: %w", err)
	}

	// TODO (litt3): Add a feature which keeps chunks from previous quorums, and just fills in gaps
	for _, quorumID := range blobHeader.QuorumNumbers {
		blob, err := pr.getBlobWithTimeout(
			ctx,
			*blobKey,
			blobHeader.Version,
			*commitment,
			eigenDACert.BatchHeader.ReferenceBlockNumber,
			quorumID)

		if err != nil {
			pr.logger.Warn(
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
		valid, err := verification.GenerateAndCompareBlobCommitment(pr.g1Srs, blob.Serialize(), commitment.Commitment)
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

		payload, err := blob.ToPayload(pr.config.PayloadPolynomialForm)
		if err != nil {
			pr.logger.Error(
				`Cert verification was successful, but conversion from blob to payload failed!
					This is likely a problem with the local blob codec configuration,
					but could potentially indicate a maliciously generated blob certificate.
					It should not be possible for an honestly generated certificate to verify
					for an invalid blob!`,
				"blobKey", blobKey.Hex(), "quorumID", quorumID, "eigenDACert", eigenDACert, "error", err)
			return nil, fmt.Errorf("decode blob: %w", err)
		}

		return payload, nil
	}

	return nil, fmt.Errorf("unable to retrieve payload from quorums %v", blobHeader.QuorumNumbers)
}

// getBlobWithTimeout attempts to get a blob from a given quorum, and times out based on config.RetrievalTimeout
func (pr *ValidatorPayloadRetriever) getBlobWithTimeout(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint32,
	quorumID core.QuorumID) (*codecs.Blob, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RetrievalTimeout)
	defer cancel()

	blobBytes, err := pr.retrievalClient.GetBlob(
		timeoutCtx,
		blobKey,
		blobVersion,
		blobCommitments,
		uint64(referenceBlockNumber),
		quorumID)

	if err != nil {
		return nil, fmt.Errorf("get blob: %w", err)
	}

	blob, err := codecs.DeserializeBlob(blobBytes, uint32(blobCommitments.Length))
	if err != nil {
		return nil, fmt.Errorf("deserialize blob: %w", err)
	}

	return blob, nil
}
