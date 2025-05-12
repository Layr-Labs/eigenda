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
		blobHeader *corev2.BlobHeaderWithHashedPayment,
		referenceBlockNumber uint64,
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

	if config.DownloadPessimism < 1 {
		logger.Warnf(
			"Download pessimism %f is less than 1, setting download pessimism to 1", config.DownloadPessimism)
		config.DownloadPessimism = 1
	}
	if config.VerificationPessimism < 1 {
		logger.Warnf(
			"Verification pessimism %f is less than 1, setting verification pessimism to 1",
			config.VerificationPessimism)
		config.VerificationPessimism = 1
	}

	if config.DownloadPessimism < config.VerificationPessimism {
		logger.Warnf(
			"Download pessimism %f is less than verification pessimism %f, setting download pessimism to %f",
			config.DownloadPessimism, config.VerificationPessimism, config.VerificationPessimism)
		config.DownloadPessimism = config.VerificationPessimism
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
	blobHeader *corev2.BlobHeaderWithHashedPayment,
	referenceBlockNumber uint64,
) ([]byte, error) {

	probe := c.metrics.newGetBlobProbe()

	probe.SetStage("verify_commitment")
	commitmentBatch := []encoding.BlobCommitments{blobHeader.BlobCommitments}
	err := c.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_operator_state")
	operatorState, err := c.chainState.GetOperatorStateWithSocket(
		ctx,
		uint(referenceBlockNumber),
		blobHeader.QuorumNumbers)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_blob_versions")
	blobVersions, err := c.ethClient.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, err
	}

	blobParams, ok := blobVersions[blobHeader.BlobVersion]
	if !ok {
		return nil, fmt.Errorf("invalid blob version %d", blobHeader.BlobVersion)
	}

	probe.SetStage("get_encoding_params")
	encodingParams, err := corev2.GetEncodingParams(blobHeader.BlobCommitments.Length, blobParams)
	if err != nil {
		return nil, err
	}

	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_assignments")
	assignments, err := corev2.GetAssignments(operatorState, blobParams, blobHeader.QuorumNumbers, blobKey[:])
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	minimumChunkCount := uint32(encodingParams.NumChunks) / blobParams.CodingRate

	sockets := GetFlattenedOperatorSockets(operatorState.Operators)

	worker, err := newRetrievalWorker(
		ctx,
		c.logger,
		c.config,
		c.connectionPool,
		c.computePool,
		c.config.UnsafeValidatorGRPCManagerFactory(c.logger, sockets),
		c.config.UnsafeChunkDeserializerFactory(assignments, c.verifier),
		c.config.UnsafeBlobDecoderFactory(c.verifier),
		assignments,
		minimumChunkCount,
		&encodingParams,
		blobHeader,
		blobKey,
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

func GetFlattenedOperatorSockets(operatorsMap map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo) map[core.OperatorID]core.OperatorSocket {

	operatorSockets := make(map[core.OperatorID]core.OperatorSocket)
	for _, quorumOperators := range operatorsMap {
		for opID, operator := range quorumOperators {
			if _, ok := operatorSockets[opID]; !ok {
				operatorSockets[opID] = operator.Socket
			}
		}
	}
	return operatorSockets
}
