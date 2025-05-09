package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
)

// ValidatorClient is an object that can retrieve blobs from the validator nodes.
// To retrieve a blob from the relay, use RelayClient instead.
type ValidatorClient interface {
	// GetBlob downloads chunks of a blob from operator network and reconstructs the blob.
	GetBlob(
		ctx context.Context,
		blobKey corev2.BlobKey,
		blobVersion corev2.BlobVersion,
		blobCommitments encoding.BlobCommitments,
		referenceBlockNumber uint64,
		quorumID core.QuorumID,
	) ([]byte, error)
}

type validatorClient struct {
	logger         logging.Logger
	ethClient      core.Reader
	chainState     core.ChainState
	verifier       encoding.Verifier
	config         *ValidatorClientConfig
	connectionPool *workerpool.WorkerPool
	computePool    *workerpool.WorkerPool
	metrics        *ValidatorClientMetrics
}

var _ ValidatorClient = &validatorClient{}

// NewValidatorClient creates a new retrieval client.
func NewValidatorClient(
	logger logging.Logger,
	ethClient core.Reader,
	chainState core.ChainState,
	verifier encoding.Verifier,
	config *ValidatorClientConfig,
	metrics *ValidatorClientMetrics,
) ValidatorClient {

	if config.ConnectionPoolSize <= 0 {
		config.ConnectionPoolSize = 1
	}
	if config.ComputePoolSize <= 0 {
		config.ComputePoolSize = 1
	}

	return &validatorClient{
		logger:         logger.With("component", "ValidatorClient"),
		ethClient:      ethClient,
		chainState:     chainState,
		verifier:       verifier,
		config:         config,
		connectionPool: workerpool.New(config.ConnectionPoolSize),
		computePool:    workerpool.New(config.ComputePoolSize),
		metrics:        metrics,
	}
}

func (c *validatorClient) GetBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
) ([]byte, error) {

	probe := c.metrics.newGetBlobProbe()
	defer probe.End()

	probe.SetStage("verify_commitment")
	commitmentBatch := []encoding.BlobCommitments{blobCommitments}
	err := c.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_operator_state")
	operatorState, err := c.chainState.GetOperatorStateWithSocket(
		ctx,
		uint(referenceBlockNumber),
		[]core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_blob_versions")
	blobVersions, err := c.ethClient.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, err
	}

	blobParams, ok := blobVersions[blobVersion]
	if !ok {
		return nil, fmt.Errorf("invalid blob version %d", blobVersion)
	}

	probe.SetStage("get_encoding_params")
	encodingParams, err := corev2.GetEncodingParams(blobCommitments.Length, blobParams)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_assignments")
	assignments, err := corev2.GetAssignments(operatorState, blobParams, quorumID)
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	totalChunkCount := uint32(encodingParams.NumChunks)
	minimumChunkCount := uint32(encodingParams.NumChunks) / blobParams.CodingRate

	worker, err := newRetrievalWorker(
		ctx,
		c.logger,
		c.config,
		c.connectionPool,
		c.computePool,
		c.config.UnsafeValidatorGRPCManagerFactory(c.logger, operatorState.Operators),
		c.config.UnsafeChunkDeserializerFactory(assignments, c.verifier),
		c.config.UnsafeBlobDecoderFactory(c.verifier),
		assignments,
		totalChunkCount,
		minimumChunkCount,
		&encodingParams,
		quorumID,
		blobKey,
		&blobCommitments,
		probe)
	if err != nil {
		return nil, fmt.Errorf("failed to create retrieval worker: %w", err)
	}

	data, err := worker.retrieveBlobFromValidators()
	if err != nil {
		return nil, fmt.Errorf("failed to download blob from validators: %w", err)
	}
	return data, nil
}
