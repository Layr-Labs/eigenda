package batcher

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/wealdtech/go-merkletree/v2"
	grpc_metadata "google.golang.org/grpc/metadata"
)

const encodingInterval = 2 * time.Second

const operatorStateCacheSize = 32

var errNoEncodedResults = errors.New("no encoded results")

type EncodedSizeNotifier struct {
	mu sync.Mutex

	Notify chan struct{}
	// threshold is the size of the total encoded blob results in bytes that triggers the notifier
	threshold uint64
	// active is set to false after the notifier is triggered to prevent it from triggering again for the same batch
	// This is reset when CreateBatch is called and the encoded results have been consumed
	active bool
}

type StreamerConfig struct {

	// SRSOrder is the order of the SRS used for encoding
	SRSOrder int
	// EncodingRequestTimeout is the timeout for each encoding request
	EncodingRequestTimeout time.Duration

	// ChainStateTimeout is the timeout used for getting the chainstate
	ChainStateTimeout time.Duration

	// EncodingQueueLimit is the maximum number of encoding requests that can be queued
	EncodingQueueLimit int

	// TargetNumChunks is the target number of chunks per encoded blob
	TargetNumChunks uint

	// Maximum number of Blobs to fetch from store
	MaxBlobsToFetchFromStore int

	FinalizationBlockDelay uint
}

type EncodingStreamer struct {
	StreamerConfig

	mu sync.RWMutex

	EncodedBlobstore     *encodedBlobStore
	ReferenceBlockNumber uint
	Pool                 common.WorkerPool
	EncodedSizeNotifier  *EncodedSizeNotifier

	blobStore             disperser.BlobStore
	chainState            core.IndexedChainState
	encoderClient         disperser.EncoderClient
	assignmentCoordinator core.AssignmentCoordinator

	encodingCtxCancelFuncs []context.CancelFunc

	metrics        *EncodingStreamerMetrics
	batcherMetrics *Metrics
	logger         logging.Logger

	// Used to keep track of the last evaluated key for fetching metadatas
	exclusiveStartKey *disperser.BlobStoreExclusiveStartKey

	operatorStateCache *lru.Cache[string, *core.IndexedOperatorState]
}

type batch struct {
	EncodedBlobs []core.EncodedBlob
	BlobMetadata []*disperser.BlobMetadata
	BlobHeaders  []*core.BlobHeader
	BatchHeader  *core.BatchHeader
	State        *core.IndexedOperatorState
	MerkleTree   *merkletree.MerkleTree
}

func NewEncodedSizeNotifier(notify chan struct{}, threshold uint64) *EncodedSizeNotifier {
	return &EncodedSizeNotifier{
		Notify:    notify,
		threshold: threshold,
		active:    true,
	}
}

func NewEncodingStreamer(
	config StreamerConfig,
	blobStore disperser.BlobStore,
	chainState core.IndexedChainState,
	encoderClient disperser.EncoderClient,
	assignmentCoordinator core.AssignmentCoordinator,
	encodedSizeNotifier *EncodedSizeNotifier,
	workerPool common.WorkerPool,
	metrics *EncodingStreamerMetrics,
	batcherMetrics *Metrics,
	logger logging.Logger) (*EncodingStreamer, error) {
	if config.EncodingQueueLimit <= 0 {
		return nil, errors.New("EncodingQueueLimit should be greater than 0")
	}
	operatorStateCache, err := lru.New[string, *core.IndexedOperatorState](operatorStateCacheSize)
	if err != nil {
		return nil, err
	}
	return &EncodingStreamer{
		StreamerConfig:         config,
		EncodedBlobstore:       newEncodedBlobStore(logger),
		ReferenceBlockNumber:   uint(0),
		Pool:                   workerPool,
		EncodedSizeNotifier:    encodedSizeNotifier,
		blobStore:              blobStore,
		chainState:             chainState,
		encoderClient:          encoderClient,
		assignmentCoordinator:  assignmentCoordinator,
		encodingCtxCancelFuncs: make([]context.CancelFunc, 0),
		metrics:                metrics,
		batcherMetrics:         batcherMetrics,
		logger:                 logger.With("component", "EncodingStreamer"),
		exclusiveStartKey:      nil,
		operatorStateCache:     operatorStateCache,
	}, nil
}

func (e *EncodingStreamer) Start(ctx context.Context) error {
	encoderChan := make(chan EncodingResultOrStatus)

	// goroutine for handling blob encoding responses
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case response := <-encoderChan:
				err := e.ProcessEncodedBlobs(ctx, response)
				if err != nil {
					if strings.Contains(err.Error(), context.Canceled.Error()) {
						// ignore canceled errors because canceled encoding requests are normal
						continue
					}
					if strings.Contains(err.Error(), "too many requests") {
						e.logger.Warn("encoding request ratelimited", "err", err)
					} else if strings.Contains(err.Error(), "connection reset by peer") {
						e.logger.Warn("encoder connection reset by peer", "err", err)
					} else if strings.Contains(err.Error(), "error reading from server: EOF") {
						e.logger.Warn("encoder request dropped", "err", err)
					} else if strings.Contains(err.Error(), "connection refused") {
						e.logger.Warn("encoder connection refused", "err", err)
					} else {
						e.logger.Error("error processing encoded blobs", "err", err)
					}
				}
			}
		}
	}()

	// goroutine for making blob encoding requests
	go func() {
		ticker := time.NewTicker(encodingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := e.RequestEncoding(ctx, encoderChan)
				if err != nil {
					e.logger.Warn("error requesting encoding", "err", err)
				}
			}
		}
	}()

	return nil
}

func (e *EncodingStreamer) dedupRequests(metadatas []*disperser.BlobMetadata, referenceBlockNumber uint) []*disperser.BlobMetadata {
	res := make([]*disperser.BlobMetadata, 0)
	for _, meta := range metadatas {
		allQuorumsRequested := true
		// check if the blob has been requested for all quorums
		for _, quorum := range meta.RequestMetadata.SecurityParams {
			if !e.EncodedBlobstore.HasEncodingRequested(meta.GetBlobKey(), quorum.QuorumID, referenceBlockNumber) {
				allQuorumsRequested = false
				break
			}
		}
		if !allQuorumsRequested {
			res = append(res, meta)
		}
	}

	return res
}

func (e *EncodingStreamer) RequestEncoding(ctx context.Context, encoderChan chan EncodingResultOrStatus) error {
	stageTimer := time.Now()
	// pull new blobs and send to encoder
	e.mu.Lock()
	metadatas, newExclusiveStartKey, err := e.blobStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, int32(e.StreamerConfig.MaxBlobsToFetchFromStore), e.exclusiveStartKey)
	e.exclusiveStartKey = newExclusiveStartKey
	e.mu.Unlock()

	if err != nil {
		return fmt.Errorf("error getting blob metadatas: %w", err)
	}
	if len(metadatas) == 0 {
		e.logger.Info("no new metadatas to encode")
		return nil
	}

	// read lock to access e.ReferenceBlockNumber
	e.mu.RLock()
	referenceBlockNumber := e.ReferenceBlockNumber
	e.mu.RUnlock()

	if referenceBlockNumber == 0 {
		// Update the reference block number for the next iteration
		blockNumber, err := e.chainState.GetCurrentBlockNumber()
		if err != nil {
			return fmt.Errorf("failed to get current block number, won't request encoding: %w", err)
		} else {
			if blockNumber > e.FinalizationBlockDelay {
				blockNumber -= e.FinalizationBlockDelay
			}

			e.mu.Lock()
			e.ReferenceBlockNumber = blockNumber
			e.mu.Unlock()
			referenceBlockNumber = blockNumber
		}
	}

	e.logger.Debug("metadata in processing status", "numMetadata", len(metadatas))
	metadatas = e.dedupRequests(metadatas, referenceBlockNumber)
	if len(metadatas) == 0 {
		e.logger.Info("no new metadatas to encode")
		return nil
	}

	waitingQueueSize := e.Pool.WaitingQueueSize()
	numMetadatastoProcess := e.EncodingQueueLimit - waitingQueueSize
	if numMetadatastoProcess > len(metadatas) {
		numMetadatastoProcess = len(metadatas)
	}
	if numMetadatastoProcess <= 0 {
		// encoding queue is full
		e.logger.Warn("worker pool queue is full. skipping this round of encoding requests", "waitingQueueSize", waitingQueueSize, "encodingQueueLimit", e.EncodingQueueLimit)
		return nil
	}
	// only process subset of blobs so it doesn't exceed the EncodingQueueLimit
	// TODO: this should be done at the request time and keep the cursor so that we don't fetch the same metadata every time
	metadatas = metadatas[:numMetadatastoProcess]

	e.logger.Debug("new metadatas to encode", "numMetadata", len(metadatas), "duration", time.Since(stageTimer))

	// Get the operator state

	timeoutCtx, cancel := context.WithTimeout(ctx, e.ChainStateTimeout)
	defer cancel()
	state, err := e.getOperatorState(timeoutCtx, metadatas, referenceBlockNumber)
	if err != nil {
		return fmt.Errorf("error getting operator state: %w", err)
	}
	metadatas = e.validateMetadataQuorums(metadatas, state)

	metadataByKey := make(map[disperser.BlobKey]*disperser.BlobMetadata, 0)
	for _, metadata := range metadatas {
		metadataByKey[metadata.GetBlobKey()] = metadata
	}

	stageTimer = time.Now()
	blobs, err := e.blobStore.GetBlobsByMetadata(ctx, metadatas)
	if err != nil {
		return fmt.Errorf("error getting blobs from blob store: %w", err)
	}
	e.logger.Debug("retrieved blobs to encode", "numBlobs", len(blobs), "duration", time.Since(stageTimer))

	e.logger.Debug("encoding blobs...", "numBlobs", len(blobs), "blockNumber", referenceBlockNumber)

	for i := range metadatas {
		metadata := metadatas[i]

		e.RequestEncodingForBlob(ctx, metadata, blobs[metadata.GetBlobKey()], state, referenceBlockNumber, encoderChan)
	}

	return nil
}

type pendingRequestInfo struct {
	BlobQuorumInfo *core.BlobQuorumInfo
	EncodingParams encoding.EncodingParams
	Assignments    map[core.OperatorID]core.Assignment
}

func (e *EncodingStreamer) RequestEncodingForBlob(ctx context.Context, metadata *disperser.BlobMetadata, blob *core.Blob, state *core.IndexedOperatorState, referenceBlockNumber uint, encoderChan chan EncodingResultOrStatus) {

	// Validate the encoding parameters for each quorum

	blobKey := metadata.GetBlobKey()

	pending := make([]pendingRequestInfo, 0, len(metadata.RequestMetadata.SecurityParams))

	for ind := range metadata.RequestMetadata.SecurityParams {

		quorum := metadata.RequestMetadata.SecurityParams[ind]

		// Check if the blob has already been encoded for this quorum
		if e.EncodedBlobstore.HasEncodingRequested(blobKey, quorum.QuorumID, referenceBlockNumber) {
			continue
		}

		blobLength := encoding.GetBlobLength(metadata.RequestMetadata.BlobSize)

		chunkLength, err := e.assignmentCoordinator.CalculateChunkLength(state.OperatorState, blobLength, e.StreamerConfig.TargetNumChunks, quorum)
		if err != nil {
			e.logger.Error("error calculating chunk length", "err", err)
			continue
		}

		blobQuorumInfo := &core.BlobQuorumInfo{
			SecurityParam: core.SecurityParam{
				QuorumID:              quorum.QuorumID,
				AdversaryThreshold:    quorum.AdversaryThreshold,
				ConfirmationThreshold: quorum.ConfirmationThreshold,
				QuorumRate:            quorum.QuorumRate,
			},
			ChunkLength: chunkLength,
		}
		assignments, info, err := e.assignmentCoordinator.GetAssignments(state.OperatorState, blobLength, blobQuorumInfo)
		if err != nil {
			e.logger.Error("error getting assignments", "err", err)
			continue
		}

		params := encoding.ParamsFromMins(chunkLength, info.TotalChunks)

		err = encoding.ValidateEncodingParamsAndBlobLength(params, uint64(blobLength), uint64(e.SRSOrder))
		if err != nil {
			e.logger.Error("invalid encoding params", "err", err)
			// Cancel the blob
			err := e.blobStore.MarkBlobFailed(ctx, blobKey)
			if err != nil {
				e.logger.Error("error marking blob failed", "err", err)
			}
			return
		}

		pending = append(pending, pendingRequestInfo{
			BlobQuorumInfo: blobQuorumInfo,
			EncodingParams: params,
			Assignments:    assignments,
		})
	}

	if len(pending) > 0 {
		requestTime := time.Unix(0, int64(metadata.RequestMetadata.RequestedAt))
		e.batcherMetrics.ObserveBlobAge("encoding_requested", float64(time.Since(requestTime).Milliseconds()))
	}

	// Execute the encoding requests
	for ind := range pending {
		res := pending[ind]

		// Create a new context for each encoding request
		// This allows us to cancel all outstanding encoding requests when we create a new batch
		// This is necessary because an encoding request is dependent on the reference block number
		// If the reference block number changes, we need to cancel all outstanding encoding requests
		// and re-request them with the new reference block number
		encodingCtx, cancel := context.WithTimeout(ctx, e.EncodingRequestTimeout)
		e.mu.Lock()
		e.encodingCtxCancelFuncs = append(e.encodingCtxCancelFuncs, cancel)
		e.mu.Unlock()

		// Add headers for routing
		md := grpc_metadata.New(map[string]string{
			"content-type":   "application/grpc",
			"x-payload-size": fmt.Sprintf("%d", len(blob.Data)),
		})
		encodingCtx = grpc_metadata.NewOutgoingContext(encodingCtx, md)

		e.Pool.Submit(func() {
			defer cancel()
			start := time.Now()
			commits, chunks, err := e.encoderClient.EncodeBlob(encodingCtx, blob.Data, res.EncodingParams)
			if err != nil {
				encoderChan <- EncodingResultOrStatus{Err: err, EncodingResult: EncodingResult{
					BlobMetadata:   metadata,
					BlobQuorumInfo: res.BlobQuorumInfo,
				}}
				e.metrics.ObserveEncodingLatency("failed", res.BlobQuorumInfo.QuorumID, len(blob.Data), float64(time.Since(start).Milliseconds()))
				return
			}

			encoderChan <- EncodingResultOrStatus{
				EncodingResult: EncodingResult{
					BlobMetadata:         metadata,
					ReferenceBlockNumber: referenceBlockNumber,
					BlobQuorumInfo:       res.BlobQuorumInfo,
					Commitment:           commits,
					ChunksData:           chunks,
					Assignments:          res.Assignments,
				},
				Err: nil,
			}
			e.metrics.ObserveEncodingLatency("success", res.BlobQuorumInfo.QuorumID, len(blob.Data), float64(time.Since(start).Milliseconds()))
		})
		e.EncodedBlobstore.PutEncodingRequest(blobKey, res.BlobQuorumInfo.QuorumID)
	}
}

func (e *EncodingStreamer) ProcessEncodedBlobs(ctx context.Context, result EncodingResultOrStatus) error {
	if result.Err != nil {
		e.EncodedBlobstore.DeleteEncodingRequest(result.BlobMetadata.GetBlobKey(), result.BlobQuorumInfo.QuorumID)
		return fmt.Errorf("error encoding blob: %w", result.Err)
	}

	err := e.EncodedBlobstore.PutEncodingResult(&result.EncodingResult)
	if err != nil {
		return fmt.Errorf("failed to putEncodedBlob: %w", err)
	}

	requestTime := time.Unix(0, int64(result.BlobMetadata.RequestMetadata.RequestedAt))
	e.batcherMetrics.ObserveBlobAge("encoded", float64(time.Since(requestTime).Milliseconds()))
	e.batcherMetrics.IncrementBlobSize("encoded", result.BlobQuorumInfo.QuorumID, int(result.BlobMetadata.RequestMetadata.BlobSize))

	count, encodedSize := e.EncodedBlobstore.GetEncodedResultSize()
	e.metrics.UpdateEncodedBlobs(count, encodedSize)
	if e.EncodedSizeNotifier.threshold > 0 && encodedSize >= e.EncodedSizeNotifier.threshold {
		e.EncodedSizeNotifier.mu.Lock()

		if e.EncodedSizeNotifier.active {
			e.logger.Info("encoded size threshold reached", "size", encodedSize)
			e.EncodedSizeNotifier.Notify <- struct{}{}
			// make sure this doesn't keep triggering before encoded blob store is reset
			e.EncodedSizeNotifier.active = false
		}
		e.EncodedSizeNotifier.mu.Unlock()
	}

	return nil
}

func (e *EncodingStreamer) UpdateReferenceBlock(currentBlockNumber uint) error {
	blockNumber := currentBlockNumber
	if blockNumber > e.FinalizationBlockDelay {
		blockNumber -= e.FinalizationBlockDelay
	}
	if e.ReferenceBlockNumber > blockNumber {
		return fmt.Errorf("reference block number is being updated to a lower value: from %d to %d", e.ReferenceBlockNumber, blockNumber)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ReferenceBlockNumber < blockNumber {
		// Wipe out the encoding results based on previous reference block number
		_ = e.EncodedBlobstore.PopLatestEncodingResults(e.ReferenceBlockNumber)
	}
	e.ReferenceBlockNumber = blockNumber
	return nil
}

// CreateBatch makes a batch from all blobs in the encoded blob store.
// If successful, it returns a batch, and updates the reference block number for next batch to use.
// Otherwise, it returns an error and keeps the blobs in the encoded blob store.
// This function is meant to be called periodically in a single goroutine as it resets the state of the encoded blob store.
func (e *EncodingStreamer) CreateBatch(ctx context.Context) (*batch, error) {
	// lock to update e.ReferenceBlockNumber
	e.mu.Lock()
	defer e.mu.Unlock()
	// Cancel outstanding encoding requests
	// Assumption: `CreateBatch` will be called at an interval longer than time it takes to encode a single blob
	if len(e.encodingCtxCancelFuncs) > 0 {
		e.logger.Info("canceling outstanding encoding requests", "count", len(e.encodingCtxCancelFuncs))
		for _, cancel := range e.encodingCtxCancelFuncs {
			cancel()
		}
		e.encodingCtxCancelFuncs = make([]context.CancelFunc, 0)
	}

	// If there were no requested blobs between the last batch and now, there is no need to create a new batch
	if e.ReferenceBlockNumber == 0 {
		blockNumber, err := e.chainState.GetCurrentBlockNumber()
		if err != nil {
			e.logger.Error("failed to get current block number. will not clean up the encoded blob store.", "err", err)
		} else {
			_ = e.EncodedBlobstore.GetNewAndDeleteStaleEncodingResults(blockNumber)
		}
		return nil, errNoEncodedResults
	}

	// Delete any encoded results that are not from the current batching iteration (i.e. that has different reference block number)
	// If any pending encoded results are discarded here, it will be re-requested in the next iteration
	encodedResults := e.EncodedBlobstore.GetNewAndDeleteStaleEncodingResults(e.ReferenceBlockNumber)

	// Reset the notifier
	e.EncodedSizeNotifier.mu.Lock()
	e.EncodedSizeNotifier.active = true
	e.EncodedSizeNotifier.mu.Unlock()

	e.logger.Info("creating a batch...", "numBlobs", len(encodedResults), "refblockNumber", e.ReferenceBlockNumber)
	if len(encodedResults) == 0 {
		return nil, errNoEncodedResults
	}

	encodedBlobByKey := make(map[disperser.BlobKey]core.EncodedBlob)
	blobQuorums := make(map[disperser.BlobKey][]*core.BlobQuorumInfo)
	blobHeaderByKey := make(map[disperser.BlobKey]*core.BlobHeader)
	metadataByKey := make(map[disperser.BlobKey]*disperser.BlobMetadata)
	for i := range encodedResults {
		// each result represent an encoded result per (blob, quorum param)
		// if the same blob has been dispersed multiple time with different security params,
		// there will be multiple encoded results for that (blob, quorum)
		result := encodedResults[i]
		blobKey := result.BlobMetadata.GetBlobKey()
		if _, ok := encodedBlobByKey[blobKey]; !ok {
			metadataByKey[blobKey] = result.BlobMetadata
			blobQuorums[blobKey] = make([]*core.BlobQuorumInfo, 0)
			blobHeader := &core.BlobHeader{
				BlobCommitments: *result.Commitment,
			}
			blobHeaderByKey[blobKey] = blobHeader
			encodedBlobByKey[blobKey] = core.EncodedBlob{
				BlobHeader:               blobHeader,
				EncodedBundlesByOperator: make(map[core.OperatorID]core.EncodedBundles),
			}
		}

		// Populate the assigned bundles
		for opID, assignment := range result.Assignments {
			bundles, ok := encodedBlobByKey[blobKey].EncodedBundlesByOperator[opID]
			if !ok {
				encodedBlobByKey[blobKey].EncodedBundlesByOperator[opID] = make(core.EncodedBundles)
				bundles = encodedBlobByKey[blobKey].EncodedBundlesByOperator[opID]
			}
			bundles[result.BlobQuorumInfo.QuorumID] = new(core.ChunksData)
			bundles[result.BlobQuorumInfo.QuorumID].Format = result.ChunksData.Format
			bundles[result.BlobQuorumInfo.QuorumID].Chunks = append(bundles[result.BlobQuorumInfo.QuorumID].Chunks, result.ChunksData.Chunks[assignment.StartIndex:assignment.StartIndex+assignment.NumChunks]...)
			bundles[result.BlobQuorumInfo.QuorumID].ChunkLen = result.ChunksData.ChunkLen
		}

		blobQuorums[blobKey] = append(blobQuorums[blobKey], result.BlobQuorumInfo)
	}

	// Populate the blob quorum infos
	for blobKey, encodedBlob := range encodedBlobByKey {
		encodedBlob.BlobHeader.QuorumInfos = blobQuorums[blobKey]
	}

	for blobKey, metadata := range metadataByKey {
		quorumPresent := make(map[core.QuorumID]bool)
		for _, quorum := range blobQuorums[blobKey] {
			quorumPresent[quorum.QuorumID] = true
		}
		// Check if the blob has valid quorums. If any of the quorums are not valid, delete the blobKey
		for _, quorum := range metadata.RequestMetadata.SecurityParams {
			_, ok := quorumPresent[quorum.QuorumID]
			if !ok {
				// Delete the blobKey. These encoded blobs will be automatically removed by the next run of
				// RequestEncoding
				delete(metadataByKey, blobKey)
				break
			}
		}
	}

	if len(metadataByKey) == 0 {
		return nil, errNoEncodedResults
	}

	// Transform maps to slices so orders in different slices match
	encodedBlobs := make([]core.EncodedBlob, 0, len(metadataByKey))
	blobHeaders := make([]*core.BlobHeader, 0, len(metadataByKey))
	metadatas := make([]*disperser.BlobMetadata, 0, len(metadataByKey))
	for key := range metadataByKey {
		err := e.transitionBlobToDispersing(ctx, metadataByKey[key])
		if err != nil {
			continue
		}
		encodedBlobs = append(encodedBlobs, encodedBlobByKey[key])
		blobHeaders = append(blobHeaders, blobHeaderByKey[key])
		metadatas = append(metadatas, metadataByKey[key])
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), e.ChainStateTimeout)
	defer cancel()

	state, err := e.getOperatorState(timeoutCtx, metadatas, e.ReferenceBlockNumber)
	if err != nil {
		for _, metadata := range metadatas {
			_ = e.handleFailedMetadata(ctx, metadata)
		}
		return nil, err
	}

	// Populate the batch header
	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: e.ReferenceBlockNumber,
		BatchRoot:            [32]byte{},
	}

	tree, err := batchHeader.SetBatchRoot(blobHeaders)
	if err != nil {
		for _, metadata := range metadatas {
			_ = e.handleFailedMetadata(ctx, metadata)
		}
		return nil, err
	}

	e.ReferenceBlockNumber = 0

	return &batch{
		EncodedBlobs: encodedBlobs,
		BatchHeader:  batchHeader,
		BlobHeaders:  blobHeaders,
		BlobMetadata: metadatas,
		State:        state,
		MerkleTree:   tree,
	}, nil
}

func (e *EncodingStreamer) handleFailedMetadata(ctx context.Context, metadata *disperser.BlobMetadata) error {
	err := e.blobStore.MarkBlobProcessing(ctx, metadata.GetBlobKey())
	if err != nil {
		e.logger.Error("error marking blob as processing", "err", err)
	}

	return err
}

func (e *EncodingStreamer) transitionBlobToDispersing(ctx context.Context, metadata *disperser.BlobMetadata) error {
	blobKey := metadata.GetBlobKey()
	err := e.blobStore.MarkBlobDispersing(ctx, blobKey)
	if err != nil {
		e.logger.Error("error marking blob as dispersing", "err", err, "blobKey", blobKey.String())
		return err
	}
	// remove encoded blob from storage so we don't disperse it again
	e.RemoveEncodedBlob(metadata)
	return nil
}

func (e *EncodingStreamer) RemoveEncodedBlob(metadata *disperser.BlobMetadata) {
	for _, sp := range metadata.RequestMetadata.SecurityParams {
		e.EncodedBlobstore.DeleteEncodingResult(metadata.GetBlobKey(), sp.QuorumID)
	}
}

// getOperatorState returns the operator state for the blobs that have valid quorums
func (e *EncodingStreamer) getOperatorState(ctx context.Context, metadatas []*disperser.BlobMetadata, blockNumber uint) (*core.IndexedOperatorState, error) {

	quorums := make(map[core.QuorumID]QuorumInfo, 0)
	for _, metadata := range metadatas {
		for _, quorum := range metadata.RequestMetadata.SecurityParams {
			quorums[quorum.QuorumID] = QuorumInfo{}
		}
	}

	quorumIds := make([]core.QuorumID, len(quorums))
	i := 0
	for id := range quorums {
		quorumIds[i] = id
		i++
	}

	cacheKey := computeCacheKey(blockNumber, quorumIds)
	if val, ok := e.operatorStateCache.Get(cacheKey); ok {
		return val, nil
	}
	// GetIndexedOperatorState should return state for valid quorums only
	state, err := e.chainState.GetIndexedOperatorState(ctx, blockNumber, quorumIds)
	if err != nil {
		return nil, fmt.Errorf("error getting operator state at block number %d: %w", blockNumber, err)
	}
	e.operatorStateCache.Add(cacheKey, state)
	return state, nil
}

// It also returns the list of valid blob metadatas (i.e. blobs that have valid quorums)
func (e *EncodingStreamer) validateMetadataQuorums(metadatas []*disperser.BlobMetadata, state *core.IndexedOperatorState) []*disperser.BlobMetadata {
	validMetadata := make([]*disperser.BlobMetadata, 0)
	for _, metadata := range metadatas {
		valid := true
		for _, quorum := range metadata.RequestMetadata.SecurityParams {
			if aggKey, ok := state.AggKeys[quorum.QuorumID]; !ok || aggKey == nil {
				e.logger.Warn("got blob with a quorum without APK. Will skip.", "blobKey", metadata.GetBlobKey(), "quorum", quorum.QuorumID)
				valid = false
			}
		}
		if valid {
			validMetadata = append(validMetadata, metadata)
		} else {
			_, err := e.blobStore.HandleBlobFailure(context.Background(), metadata, 0)
			if err != nil {
				e.logger.Error("error handling blob failure", "err", err)
			}
		}
	}
	return validMetadata
}

func computeCacheKey(blockNumber uint, quorumIDs []uint8) string {
	bytes := make([]byte, 8+len(quorumIDs))
	binary.LittleEndian.PutUint64(bytes, uint64(blockNumber))
	copy(bytes[8:], quorumIDs)
	return string(bytes)
}
