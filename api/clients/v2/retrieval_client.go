package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
)

// RetrievalClient is an object that can retrieve blobs from the DA nodes.
// To retrieve a blob from the relay, use RelayClient instead.
type RetrievalClient interface {
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

type retrievalClient struct {
	logger         logging.Logger
	ethClient      core.Reader
	chainState     core.ChainState
	verifier       encoding.Verifier
	connectionPool *workerpool.WorkerPool
	computePool    *workerpool.WorkerPool
}

var _ RetrievalClient = &retrievalClient{}

// NewRetrievalClient creates a new retrieval client.
func NewRetrievalClient(
	logger logging.Logger,
	ethClient core.Reader,
	chainState core.ChainState,
	verifier encoding.Verifier,
// connectionPoolSize limits the maximum number of concurrent network connections
	connectionPoolSize int,
// computePoolSize limits the maximum number of concurrent compute intensive tasks
	computePoolSize int,
) RetrievalClient {
	return &retrievalClient{
		logger:         logger.With("component", "RetrievalClient"),
		ethClient:      ethClient,
		chainState:     chainState,
		verifier:       verifier,
		connectionPool: workerpool.New(connectionPoolSize),
		computePool:    workerpool.New(computePoolSize),
	}
}

func (r *retrievalClient) GetBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
) ([]byte, error) {
	return r.GetBlobWithProbe(ctx, blobKey, blobVersion, blobCommitments, referenceBlockNumber, quorumID, nil)
}

func (r *retrievalClient) GetBlobWithProbe(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
	probe *common.SequenceProbe,
) ([]byte, error) {

	probe.SetStage("verify_commitment")
	commitmentBatch := []encoding.BlobCommitments{blobCommitments}
	err := r.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_operator_state")
	operatorState, err := r.chainState.GetOperatorStateWithSocket(ctx, uint(referenceBlockNumber), []core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}
	operators, ok := operatorState.Operators[quorumID]
	if !ok {
		return nil, fmt.Errorf("no quorum with ID: %d", quorumID)
	}

	probe.SetStage("get_blob_versions")
	blobVersions, err := r.ethClient.GetAllVersionedBlobParams(ctx)
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

	worker, err := newRetrievalWorker(
		ctx,
		r.logger,
		DefaultValidatorRetrievalConfig(),
		r.connectionPool,
		r.computePool,
		operators,
		operatorState,
		blobParams,
		encodingParams,
		assignments,
		quorumID,
		blobKey,
		r.verifier,
		blobCommitments)
	if err != nil {
		return nil, fmt.Errorf("failed to create retrieval worker: %w", err)
	}

	probe.SetStage("download_and_verify")
	data, err := worker.downloadBlobFromValidators()
	if err != nil {
		return nil, fmt.Errorf("failed to download blob from validators: %w", err)
	}
	return data, nil
}
