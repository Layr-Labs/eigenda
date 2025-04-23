package validator

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
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

	// GetBlobWithProbe is the same as GetBlob, but it also takes a probe object to capture metrics. No metrics are
	// captured if the probe is nil.
	GetBlobWithProbe(
		ctx context.Context,
		blobKey corev2.BlobKey,
		blobVersion corev2.BlobVersion,
		blobCommitments encoding.BlobCommitments,
		referenceBlockNumber uint64,
		quorumID core.QuorumID,
		probe *common.SequenceProbe,
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
}

var _ ValidatorClient = &validatorClient{}

// NewValidatorClient creates a new retrieval client.
func NewValidatorClient(
	logger logging.Logger,
	ethClient core.Reader,
	chainState core.ChainState,
	verifier encoding.Verifier,
	config *ValidatorClientConfig,
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
	return c.GetBlobWithProbe(ctx, blobKey, blobVersion, blobCommitments, referenceBlockNumber, quorumID, nil)
}

func (c *validatorClient) GetBlobWithProbe(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
	probe *common.SequenceProbe,
) ([]byte, error) {

	worker, err := newRetrievalWorker(
		ctx,
		c.logger,
		c.config,
		c.connectionPool,
		c.computePool,
		c.chainState,
		referenceBlockNumber,
		blobVersion,
		c.ethClient,
		quorumID,
		blobKey,
		c.verifier,
		blobCommitments,
		probe)
	if err != nil {
		return nil, fmt.Errorf("failed to create retrieval worker: %w", err)
	}

	data, err := worker.downloadBlobFromValidators()
	if err != nil {
		return nil, fmt.Errorf("failed to download blob from validators: %w", err)
	}
	return data, nil
}
