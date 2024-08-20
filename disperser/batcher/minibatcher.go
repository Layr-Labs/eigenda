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

type BatchState struct {
	BatchID              uuid.UUID
	ReferenceBlockNumber uint
	BlobHeaders          []*core.BlobHeader
	BlobMetadata         []*disperser.BlobMetadata
	OperatorState        *core.IndexedOperatorState
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
	Batches              map[uuid.UUID]*BatchState
	ReferenceBlockNumber uint
	CurrentBatchID       uuid.UUID
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

		Batches:              make(map[uuid.UUID]*BatchState),
		ReferenceBlockNumber: 0,
		CurrentBatchID:       uuid.Nil,
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
		cancelFuncs := make([]context.CancelFunc, 0)
		for {
			select {
			case <-ctx.Done():
				for _, cancel := range cancelFuncs {
					cancel()
				}
				return
			case <-ticker.C:
				cancel, err := b.HandleSingleMinibatch(ctx)
				if err != nil {
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Warn("no encoded results to make a batch with")
					} else {
						b.logger.Error("failed to process a batch", "err", err)
					}
				}
				if cancel != nil {
					cancelFuncs = append(cancelFuncs, cancel)
				}
			}
		}
	}()

	return nil
}

func (b *Minibatcher) PopBatchState(batchID uuid.UUID) *BatchState {
	batchState, ok := b.Batches[batchID]
	if !ok {
		return nil
	}
	delete(b.Batches, batchID)
	return batchState
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

func (b *Minibatcher) HandleSingleMinibatch(ctx context.Context) (context.CancelFunc, error) {
	log := b.logger
	// If too many dispersal requests are pending, skip an iteration
	if pending := b.Pool.WaitingQueueSize(); pending > int(b.MaxNumConnections) {
		return nil, fmt.Errorf("too many pending requests %d with max number of connections %d. skipping minibatch iteration", pending, b.MaxNumConnections)
	}
	stageTimer := time.Now()
	// All blobs in this batch are marked as DISPERSING
	minibatch, err := b.EncodingStreamer.CreateMinibatch(ctx)
	if err != nil {
		return nil, err
	}
	log.Debug("CreateMinibatch took", "duration", time.Since(stageTimer).String())
	// Processing new full batch
	if b.ReferenceBlockNumber < minibatch.BatchHeader.ReferenceBlockNumber {
		// Update status of the previous batch
		if b.CurrentBatchID != uuid.Nil {
			err = b.MinibatchStore.MarkBatchFormed(ctx, b.CurrentBatchID, b.MinibatchIndex)
			if err != nil {
				_ = b.handleFailure(ctx, minibatch.BlobMetadata, FailReason("error updating batch status"))
				return nil, fmt.Errorf("error updating batch status: %w", err)
			}
		}

		// Reset local batch state and create new batch
		b.CurrentBatchID, err = uuid.NewV7()
		if err != nil {
			_ = b.handleFailure(ctx, minibatch.BlobMetadata, FailReason("error generating batch UUID"))
			return nil, fmt.Errorf("error generating batch ID: %w", err)
		}
		b.MinibatchIndex = 0
		b.ReferenceBlockNumber = minibatch.BatchHeader.ReferenceBlockNumber
		err = b.MinibatchStore.PutBatch(ctx, &BatchRecord{
			ID:                   b.CurrentBatchID,
			CreatedAt:            time.Now().UTC(),
			ReferenceBlockNumber: b.ReferenceBlockNumber,
			Status:               BatchStatusPending,
			NumMinibatches:       0,
		})
		if err != nil {
			_ = b.handleFailure(ctx, minibatch.BlobMetadata, FailReason("error storing batch record"))
			return nil, fmt.Errorf("error storing batch record: %w", err)
		}
		b.Batches[b.CurrentBatchID] = &BatchState{
			BatchID:              b.CurrentBatchID,
			ReferenceBlockNumber: b.ReferenceBlockNumber,
			BlobHeaders:          make([]*core.BlobHeader, 0),
			BlobMetadata:         make([]*disperser.BlobMetadata, 0),
			OperatorState:        minibatch.State,
		}
	}

	// Accumulate batch metadata
	batchState := b.Batches[b.CurrentBatchID]
	batchState.BlobHeaders = append(batchState.BlobHeaders, minibatch.BlobHeaders...)
	batchState.BlobMetadata = append(batchState.BlobMetadata, minibatch.BlobMetadata...)

	// Dispatch encoded batch
	log.Debug("Dispatching encoded batch...", "batchID", b.CurrentBatchID, "minibatchIndex", b.MinibatchIndex, "referenceBlockNumber", b.ReferenceBlockNumber, "numBlobs", len(minibatch.EncodedBlobs))
	stageTimer = time.Now()
	dispersalCtx, cancelDispersal := context.WithCancel(ctx)
	storeMappingsChan := make(chan error)
	// Store the blob minibatch mappings in parallel
	go func() {
		err := b.createBlobMinibatchMappings(ctx, b.CurrentBatchID, b.MinibatchIndex, minibatch.BlobMetadata, minibatch.BlobHeaders)
		storeMappingsChan <- err
	}()
	b.DisperseBatch(dispersalCtx, minibatch.State, minibatch.EncodedBlobs, minibatch.BatchHeader, b.CurrentBatchID, b.MinibatchIndex)
	log.Debug("DisperseBatch took", "duration", time.Since(stageTimer).String())

	h, err := minibatch.State.OperatorState.Hash()
	if err != nil {
		log.Error("error getting operator state hash", "err", err)
	}
	hStr := make([]string, 0, len(h))
	for q, hash := range h {
		hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
	}
	log.Info("Successfully dispatched minibatch", "operatorStateHash", hStr)

	b.MinibatchIndex++

	// Wait for the blob minibatch mappings to be stored then cancel the dispersal process if there was an error
	storeMappingsErr := <-storeMappingsChan
	if storeMappingsErr != nil {
		cancelDispersal()
		return nil, fmt.Errorf("error storing blob minibatch mappings: %w", storeMappingsErr)
	}

	return cancelDispersal, nil
}

func (b *Minibatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader, batchID uuid.UUID, minibatchIndex uint) {
	for id, op := range state.IndexedOperators {
		opInfo := op
		opID := id
		req := &MinibatchDispersal{
			BatchID:        batchID,
			MinibatchIndex: minibatchIndex,
			OperatorID:     opID,
			Socket:         op.Socket,
			NumBlobs:       uint(len(blobs)),
			RequestedAt:    time.Now().UTC(),
		}
		err := b.MinibatchStore.PutDispersal(ctx, req)
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
			err = b.MinibatchStore.UpdateDispersalResponse(ctx, req, &DispersalResponse{
				Signatures:  signatures,
				RespondedAt: time.Now().UTC(),
				Error:       err,
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
	blobMessages := make([]*core.EncodedBlobMessage, 0)
	hasAnyBundles := false
	for _, blob := range blobs {
		if _, ok := blob.EncodedBundlesByOperator[opID]; ok {
			hasAnyBundles = true
		}
		blobMessages = append(blobMessages, &core.EncodedBlobMessage{
			BlobHeader: blob.BlobHeader,
			// Bundles will be empty if the operator is not in the quorums blob is dispersed on
			EncodedBundles: blob.EncodedBundlesByOperator[opID],
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
		signatures, err = b.Dispatcher.SendBlobsToOperator(ctxWithTimeout, blobMessages, batchHeader, op)
		cancel()
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

// createBlobMinibatchMappings creates a mapping between blob metadata and blob headers
// and stores it in the minibatch store. It assumes that the blob metadata and blob headers
// are ordered by blob index.
func (b *Minibatcher) createBlobMinibatchMappings(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, blobMetadatas []*disperser.BlobMetadata, blobHeaders []*core.BlobHeader) error {
	if len(blobMetadatas) != len(blobHeaders) {
		return fmt.Errorf("number of blob metadatas and blob headers do not match")
	}

	mappings := make([]*BlobMinibatchMapping, len(blobMetadatas))
	for i, blobMetadata := range blobMetadatas {
		blobKey := blobMetadata.GetBlobKey()
		blobHeader := blobHeaders[i]
		mappings[i] = &BlobMinibatchMapping{
			BlobKey:        &blobKey,
			BatchID:        batchID,
			MinibatchIndex: minibatchIndex,
			BlobIndex:      uint(i),
			BlobHeader: core.BlobHeader{
				BlobCommitments: blobHeader.BlobCommitments,
				QuorumInfos:     blobHeader.QuorumInfos,
				AccountID:       blobHeader.AccountID,
			},
		}
	}
	return b.MinibatchStore.PutBlobMinibatchMappings(ctx, mappings)
}
