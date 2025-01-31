package clients

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
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

// DistributedPayloadRetriever provides the ability to get payloads from the EigenDA nodes directly
//
// This struct is goroutine safe.
type DistributedPayloadRetriever struct {
	logger          logging.Logger
	config DistributedPayloadRetrieverConfig
	codec           codecs.BlobCodec
	retrievalClient RetrievalClient
	g1Srs           []bn254.G1Affine
}

var _ PayloadRetriever = &DistributedPayloadRetriever{}

// TODO: add basic constructor

// BuildDistributedPayloadRetriever builds a DistributedPayloadRetriever from config structs.
func BuildDistributedPayloadRetriever(
	logger logging.Logger,
	distributedPayloadRetrieverConfig DistributedPayloadRetrieverConfig,
	ethConfig geth.EthClientConfig,
	thegraphConfig thegraph.Config,
	kzgConfig kzg.KzgConfig,
) (*DistributedPayloadRetriever, error) {

	ethClient, err := geth.NewClient(ethConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}

	reader, err := eth.NewReader(
		logger,
		ethClient,
		distributedPayloadRetrieverConfig.BlsOperatorStateRetrieverAddr,
		distributedPayloadRetrieverConfig.EigenDAServiceManagerAddr)

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
		int(distributedPayloadRetrieverConfig.ConnectionCount))

	codec, err := codecs.CreateCodec(
		distributedPayloadRetrieverConfig.PayloadPolynomialForm,
		distributedPayloadRetrieverConfig.BlobEncodingVersion)
	if err != nil {
		return nil, err
	}

	return &DistributedPayloadRetriever{
		logger:          logger,
		config:          distributedPayloadRetrieverConfig,
		codec:           codec,
		retrievalClient: retrievalClient,
		g1Srs:           kzgVerifier.Srs.G1,
	}, nil
}

// GetPayload iteratively attempts to retrieve a given blob from the quorums listed in the EigenDACert.
//
// If the blob is successfully retrieved, then the blob is verified. If the verification succeeds, the blob is decoded
// to yield the payload (the original user data, with no padding or any modification), and the payload is returned.
func (pr *DistributedPayloadRetriever) GetPayload(
	ctx context.Context,
	eigenDACert *verification.EigenDACert,
) ([]byte, error) {

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobHeader := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader
	convertedCommitment, err := verification.BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return nil, fmt.Errorf("convert commitments binding to internal: %w", err)
	}

	// TODO (litt3): Add a feature which keeps chunks from previous quorums, and just fills in gaps
	for _, quorumID := range blobHeader.QuorumNumbers {
		blobBytes, err := pr.getBlobWithTimeout(
			ctx,
			*blobKey,
			blobHeader.Version,
			*convertedCommitment,
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

		err = verification.CheckBlobLength(blobBytes, convertedCommitment.Length)
		if err != nil {
			pr.logger.Warn("check blob length", "blobKey", blobKey.Hex(), "quorumID", quorumID, "error", err)
			continue
		}
		err = verification.VerifyBlobAgainstCert(blobKey, blobBytes, convertedCommitment.Commitment, pr.g1Srs)
		if err != nil {
			pr.logger.Warn("verify blob against cert", "blobKey", blobKey.Hex(), "quorumID", quorumID, "error", err)
			continue
		}

		payload, err := pr.codec.DecodeBlob(blobBytes)
		if err != nil {
			pr.logger.Error(
				`Cert verification was successful, but decode blob failed!
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

// getBlobWithTimeout attempts to get a blob from a given quorum, and times out based on config.FetchTimeout
func (pr *DistributedPayloadRetriever) getBlobWithTimeout(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint32,
	quorumID core.QuorumID) ([]byte, error) {

	timeoutCtx, cancel := context.WithTimeout(ctx, pr.config.RetrievalTimeout)
	defer cancel()

	return pr.retrievalClient.GetBlob(
		timeoutCtx,
		blobKey,
		blobVersion,
		blobCommitments,
		uint64(referenceBlockNumber),
		quorumID)
}
