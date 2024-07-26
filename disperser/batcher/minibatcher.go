package batcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
)

type MinibatcherConfig struct {
	PullInterval              time.Duration
	MaxNumConnections         uint
	MaxNumRetriesPerBlob      uint
	MaxNumRetriesPerDispersal uint
}

type Minibatcher struct {
	MinibatcherConfig

	BlobStore             disperser.BlobStore
	MinibatchStore        MinibatchStore
	Dispatcher            disperser.Dispatcher
	ChainState            core.IndexedChainState
	AssignmentCoordinator core.AssignmentCoordinator
	EncodingStreamer      *EncodingStreamer
	Pool                  common.WorkerPool

	// local state
	ReferenceBlockNumber uint
	BatchID              uuid.UUID
	MinibatchIndex       uint

	ethClient common.EthClient
	logger    logging.Logger
}

func NewMinibatcher(
	config MinibatcherConfig,
	blobStore disperser.BlobStore,
	minibatchStore MinibatchStore,
	dispatcher disperser.Dispatcher,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	// aggregator core.SignatureAggregator,
	encodingStreamer *EncodingStreamer,
	ethClient common.EthClient,
	workerpool common.WorkerPool,
	logger logging.Logger,
) (*Minibatcher, error) {
	return &Minibatcher{
		MinibatcherConfig:     config,
		BlobStore:             blobStore,
		MinibatchStore:        minibatchStore,
		Dispatcher:            dispatcher,
		ChainState:            chainState,
		AssignmentCoordinator: assignmentCoordinator,
		EncodingStreamer:      encodingStreamer,
		Pool:                  workerpool,

		ReferenceBlockNumber: 0,
		BatchID:              uuid.Nil,
		MinibatchIndex:       0,

		ethClient: ethClient,
		logger:    logger.With("component", "Minibatcher"),
	}, nil
}

func (b *Minibatcher) Start(ctx context.Context) error {
	err := b.ChainState.Start(ctx)
	if err != nil {
		return err
	}
	// Wait for few seconds for indexer to index blockchain
	// This won't be needed when we switch to using Graph node
	time.Sleep(indexerWarmupDelay)
	go func() {
		ticker := time.NewTicker(b.PullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := b.HandleSingleBatch(ctx); err != nil {
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Warn("no encoded results to make a batch with")
					} else {
						b.logger.Error("failed to process a batch", "err", err)
					}
				}
			}
		}
	}()

	return nil
}

func (b *Minibatcher) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	numPermanentFailures := 0
	for _, metadata := range blobMetadatas {
		b.EncodingStreamer.RemoveEncodedBlob(metadata)
		retry, err := b.BlobStore.HandleBlobFailure(ctx, metadata, b.MaxNumRetriesPerBlob)
		if err != nil {
			b.logger.Error("HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}

		if retry {
			continue
		}
		numPermanentFailures++
	}

	// Return the error(s)
	return result.ErrorOrNil()
}

func (b *Minibatcher) HandleSingleBatch(ctx context.Context) error {
	log := b.logger
	// If too many dispersal requests are pending, skip an iteration
	if pending := b.Pool.WaitingQueueSize(); pending > int(b.MaxNumConnections) {
		return fmt.Errorf("too many pending requests %d with max number of connections %d. skipping minibatch iteration", pending, b.MaxNumConnections)
	}
	stageTimer := time.Now()
	// All blobs in this batch are marked as DISPERSING
	batch, err := b.EncodingStreamer.CreateMinibatch(ctx)
	if err != nil {
		return err
	}
	log.Debug("CreateMinibatch took", "duration", time.Since(stageTimer).String())

	// Processing new full batch
	if b.ReferenceBlockNumber < batch.BatchHeader.ReferenceBlockNumber {
		// Update status of the previous batch
		if b.BatchID != uuid.Nil {
			err = b.MinibatchStore.UpdateBatchStatus(ctx, b.BatchID, BatchStatusFormed)
			if err != nil {
				_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error updating batch status"))
				return fmt.Errorf("error updating batch status: %w", err)
			}
		}

		// Create new batch
		b.BatchID, err = uuid.NewV7()
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error generating batch UUID"))
			return fmt.Errorf("error generating batch ID: %w", err)
		}
		batchHeaderHash, err := batch.BatchHeader.GetBatchHeaderHash()
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error getting batch header hash"))
			return fmt.Errorf("error getting batch header hash: %w", err)
		}
		b.MinibatchIndex = 0
		b.ReferenceBlockNumber = batch.BatchHeader.ReferenceBlockNumber
		err = b.MinibatchStore.PutBatch(ctx, &BatchRecord{
			ID:                   b.BatchID,
			CreatedAt:            time.Now().UTC(),
			ReferenceBlockNumber: b.ReferenceBlockNumber,
			HeaderHash:           batchHeaderHash,
		})
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error storing batch record"))
			return fmt.Errorf("error storing batch record: %w", err)
		}
	}

	// Store minibatch record
	blobHeaderHashes := make([][32]byte, 0, len(batch.EncodedBlobs))
	batchSize := int64(0)
	for _, blob := range batch.EncodedBlobs {
		h, err := blob.BlobHeader.GetBlobHeaderHash()
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error getting blob header hash"))
			return fmt.Errorf("error getting blob header hash: %w", err)
		}
		blobHeaderHashes = append(blobHeaderHashes, h)
		batchSize += blob.BlobHeader.EncodedSizeAllQuorums()
	}
	err = b.MinibatchStore.PutMinibatch(ctx, &MinibatchRecord{
		BatchID:              b.BatchID,
		MinibatchIndex:       b.MinibatchIndex,
		BlobHeaderHashes:     blobHeaderHashes,
		BatchSize:            uint64(batchSize),
		ReferenceBlockNumber: b.ReferenceBlockNumber,
	})
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailReason("error storing minibatch record"))
		return fmt.Errorf("error storing minibatch record: %w", err)
	}

	// Dispatch encoded batch
	log.Debug("Dispatching encoded batch...", "batchID", b.BatchID, "minibatchIndex", b.MinibatchIndex, "referenceBlockNumber", b.ReferenceBlockNumber, "numBlobs", len(batch.EncodedBlobs))
	stageTimer = time.Now()
	b.DisperseBatch(ctx, batch.State, batch.EncodedBlobs, batch.BatchHeader, b.BatchID, b.MinibatchIndex)
	log.Debug("DisperseBatch took", "duration", time.Since(stageTimer).String())

	h, err := batch.State.OperatorState.Hash()
	if err != nil {
		log.Error("error getting operator state hash", "err", err)
	}
	hStr := make([]string, 0, len(h))
	for q, hash := range h {
		hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
	}
	log.Info("Successfully dispatched minibatch", "operatorStateHash", hStr)

	b.MinibatchIndex++

	return nil
}

func (b *Minibatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader, batchID uuid.UUID, minibatchIndex uint) {
	for id, op := range state.IndexedOperators {
		opInfo := op
		opID := id
		req := &DispersalRequest{
			BatchID:        batchID,
			MinibatchIndex: minibatchIndex,
			OperatorID:     opID,
			Socket:         op.Socket,
			NumBlobs:       uint(len(blobs)),
			RequestedAt:    time.Now().UTC(),
		}
		err := b.MinibatchStore.PutDispersalRequest(ctx, req)
		if err != nil {
			b.logger.Error("failed to put dispersal request", "err", err)
			continue
		}
		b.Pool.Submit(func() {
			signatures, err := b.SendBlobsToOperatorWithRetries(ctx, blobs, batchHeader, opInfo, opID, int(b.MaxNumRetriesPerDispersal))
			if err != nil {
				b.logger.Errorf("failed to send blobs to operator %s: %v", opID.Hex(), err)
			}
			// Update the minibatch state
			err = b.MinibatchStore.PutDispersalResponse(ctx, &DispersalResponse{
				DispersalRequest: *req,
				Signatures:       signatures,
				RespondedAt:      time.Now().UTC(),
				Error:            err,
			})
			if err != nil {
				b.logger.Error("failed to put dispersal response", "err", err)
			}
		})
	}
}

func (b *Minibatcher) SendBlobsToOperatorWithRetries(
	ctx context.Context,
	blobs []core.EncodedBlob,
	batchHeader *core.BatchHeader,
	op *core.IndexedOperatorInfo,
	opID core.OperatorID,
	maxNumRetries int,
) ([]*core.Signature, error) {
	blobMessages := make([]*core.BlobMessage, 0)
	hasAnyBundles := false
	for _, blob := range blobs {
		if _, ok := blob.BundlesByOperator[opID]; ok {
			hasAnyBundles = true
		}
		blobMessages = append(blobMessages, &core.BlobMessage{
			BlobHeader: blob.BlobHeader,
			// Bundles will be empty if the operator is not in the quorums blob is dispersed on
			Bundles: blob.BundlesByOperator[opID],
		})
	}
	if !hasAnyBundles {
		// Operator is not part of any quorum, no need to send chunks
		return nil, fmt.Errorf("operator %s is not part of any quorum", opID.Hex())
	}

	numRetries := 0
	// initially set the timeout equal to the pull interval with exponential backoff
	timeout := b.PullInterval
	var signatures []*core.Signature
	var err error
	for numRetries < maxNumRetries {
		requestedAt := time.Now()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		signatures, err = b.Dispatcher.SendBlobsToOperator(ctxWithTimeout, blobMessages, batchHeader, op)
		latencyMs := float64(time.Since(requestedAt).Milliseconds())
		if err != nil {
			b.logger.Error("error sending chunks to operator", "operator", opID.Hex(), "err", err, "timeout", timeout.String(), "numRetries", numRetries, "maxNumRetries", maxNumRetries)
			numRetries++
			timeout *= 2
			continue
		}
		b.logger.Info("sent chunks to operator", "operator", opID.Hex(), "latencyMs", latencyMs)
		break
	}

	if signatures == nil && err != nil {
		return nil, fmt.Errorf("failed to send chunks to operator %s: %w", opID.Hex(), err)
	}

	return signatures, nil
}
