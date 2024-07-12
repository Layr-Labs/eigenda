package batcher_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	cmock "github.com/Layr-Labs/eigenda/common/mock"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/inmem"
	dinmem "github.com/Layr-Labs/eigenda/disperser/common/inmem"
	dmock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	opId0, _          = core.OperatorIDFromHex("e22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311")
	opId1, _          = core.OperatorIDFromHex("e23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312")
	mockChainState, _ = coremock.NewChainDataMock(map[uint8]map[core.OperatorID]int{
		0: {
			opId0: 1,
			opId1: 1,
		},
		1: {
			opId0: 1,
		},
	})
	defaultConfig = batcher.MinibatcherConfig{
		PullInterval:              1 * time.Second,
		MaxNumConnections:         3,
		MaxNumRetriesPerDispersal: 3,
	}
)

const (
	initialBlock = uint(10)
)

type minibatcherComponents struct {
	minibatcher           *batcher.Minibatcher
	blobStore             disperser.BlobStore
	minibatchStore        batcher.MinibatchStore
	dispatcher            *dmock.Dispatcher
	chainState            *coremock.MockIndexedChainState
	assignmentCoordinator core.AssignmentCoordinator
	encodingStreamer      *batcher.EncodingStreamer
	pool                  *workerpool.WorkerPool
	ethClient             *commonmock.MockEthClient
	logger                logging.Logger
}

func newMinibatcher(t *testing.T, config batcher.MinibatcherConfig) *minibatcherComponents {
	logger := logging.NewNoopLogger()
	blobStore := dinmem.NewBlobStore()
	minibatchStore := inmem.NewMinibatchStore(logger)
	chainState, err := coremock.NewChainDataMock(mockChainState.Stakes)
	assert.NoError(t, err)
	state := chainState.GetTotalOperatorState(context.Background(), 0)
	dispatcher := dmock.NewDispatcher(state)
	ics := &coremock.MockIndexedChainState{}
	streamerConfig := batcher.StreamerConfig{
		SRSOrder:                 3000,
		EncodingRequestTimeout:   5 * time.Second,
		EncodingQueueLimit:       10,
		TargetNumChunks:          8092,
		MaxBlobsToFetchFromStore: 10,
		FinalizationBlockDelay:   0,
		ChainStateTimeout:        5 * time.Second,
	}
	encodingWorkerPool := workerpool.New(10)
	p, err := makeTestProver()
	assert.NoError(t, err)
	encoderClient := disperser.NewLocalEncoderClient(p)
	asgn := &core.StdAssignmentCoordinator{}
	chainState.On("GetCurrentBlockNumber").Return(initialBlock, nil)
	metrics := batcher.NewMetrics("9100", logger)
	trigger := batcher.NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		10*1024*1024,
	)
	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, chainState, encoderClient, asgn, trigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
	assert.NoError(t, err)
	ethClient := &cmock.MockEthClient{}
	pool := workerpool.New(int(config.MaxNumConnections))
	m, err := batcher.NewMinibatcher(config, blobStore, minibatchStore, dispatcher, chainState, asgn, encodingStreamer, ethClient, pool, logger)
	assert.NoError(t, err)
	return &minibatcherComponents{
		minibatcher:           m,
		blobStore:             blobStore,
		minibatchStore:        minibatchStore,
		dispatcher:            dispatcher,
		chainState:            ics,
		assignmentCoordinator: asgn,
		encodingStreamer:      encodingStreamer,
		pool:                  pool,
		ethClient:             ethClient,
		logger:                logger,
	}
}

func TestDisperseMinibatch(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	var X, Y fp.Element
	X = *X.SetBigInt(big.NewInt(1))
	Y = *Y.SetBigInt(big.NewInt(2))
	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*core.Signature{
		{
			G1Point: &core.G1Point{
				G1Affine: &bn254.G1Affine{
					X: X,
					Y: Y,
				},
			},
		},
	}, nil)
	ctx := context.Background()

	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 100,
	}})
	_, _ = queueBlob(t, ctx, &blob1, c.blobStore)
	_, _ = queueBlob(t, ctx, &blob2, c.blobStore)

	out := make(chan batcher.EncodingResultOrStatus)
	err := c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	encoded1 := <-out
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, encoded1)
	assert.NoError(t, err)
	encoded2 := <-out
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, encoded2)
	assert.NoError(t, err)
	count, _ := c.encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, 2, count)

	err = c.minibatcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(1))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].NumMinibatches, uint(1))
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobHeaders, 2)
	assert.ElementsMatch(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobMetadata, []*disperser.BlobMetadata{encoded1.BlobMetadata, encoded2.BlobMetadata})

	// Second minibatch
	blob3 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	_, _ = queueBlob(t, ctx, &blob3, c.blobStore)
	err = c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	encoded3 := <-out
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, encoded3)
	assert.NoError(t, err)
	err = c.minibatcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(2))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].NumMinibatches, uint(2))
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobHeaders, 3)
	assert.ElementsMatch(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobMetadata, []*disperser.BlobMetadata{encoded1.BlobMetadata, encoded2.BlobMetadata, encoded3.BlobMetadata})
	assert.NotNil(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].OperatorState)

	b, err := c.minibatchStore.GetBatch(ctx, c.minibatcher.CurrentBatchID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, c.minibatcher.CurrentBatchID, b.ID)
	assert.NotNil(t, b.CreatedAt)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, b.ReferenceBlockNumber)
	mb, err := c.minibatchStore.GetMinibatch(ctx, c.minibatcher.CurrentBatchID, 0)
	assert.NoError(t, err)
	assert.NotNil(t, mb)
	assert.Equal(t, c.minibatcher.CurrentBatchID, mb.BatchID)
	assert.Equal(t, uint(0), mb.MinibatchIndex)
	assert.Len(t, mb.BlobHeaderHashes, 2)
	assert.Equal(t, uint64(12800), mb.BatchSize)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, mb.ReferenceBlockNumber)
	mb, err = c.minibatchStore.GetMinibatch(ctx, c.minibatcher.CurrentBatchID, 1)
	assert.NoError(t, err)
	assert.NotNil(t, mb)
	assert.Equal(t, c.minibatcher.CurrentBatchID, mb.BatchID)
	assert.Equal(t, uint(1), mb.MinibatchIndex)
	assert.Len(t, mb.BlobHeaderHashes, 1)
	assert.Equal(t, uint64(7680), mb.BatchSize)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, mb.ReferenceBlockNumber)

	// Create a new minibatch with increased reference block number
	// Test that the previous batch is marked as formed and that the new batch is created with the correct reference block number
	_, _ = queueBlob(t, ctx, &blob1, c.blobStore)
	_, _ = queueBlob(t, ctx, &blob2, c.blobStore)

	err = c.encodingStreamer.UpdateReferenceBlock(initialBlock + 10)
	assert.NoError(t, err)
	err = c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	encoded4 := <-out
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, encoded4)
	assert.NoError(t, err)
	encoded5 := <-out
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, encoded5)
	assert.NoError(t, err)
	err = c.minibatcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)

	c.pool.StopWait()

	// previous batch should be marked as formed
	b, err = c.minibatchStore.GetBatch(ctx, b.ID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.Status, batcher.BatchStatusFormed)

	// new batch should be created
	assert.NotEqual(t, c.minibatcher.CurrentBatchID, b.ID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(1))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock+10)
	assert.Len(t, c.minibatcher.Batches, 2)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].NumMinibatches, uint(1))
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].ReferenceBlockNumber, initialBlock+10)
	assert.Len(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobHeaders, 2)
	assert.ElementsMatch(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobMetadata, []*disperser.BlobMetadata{encoded4.BlobMetadata, encoded5.BlobMetadata})
	assert.NotNil(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].OperatorState)

	newBatch, err := c.minibatchStore.GetBatch(ctx, c.minibatcher.CurrentBatchID)
	assert.NoError(t, err)
	assert.NotNil(t, newBatch)
	assert.Equal(t, newBatch.ReferenceBlockNumber, initialBlock+10)
	assert.Equal(t, newBatch.Status, batcher.BatchStatusPending)

	// Test PopBatchState
	batchState := c.minibatcher.PopBatchState(b.ID)
	assert.NotNil(t, batchState)
	assert.Equal(t, batchState.BatchID, b.ID)
	assert.Equal(t, batchState.ReferenceBlockNumber, initialBlock)
	assert.Equal(t, batchState.NumMinibatches, uint(2))
	assert.Len(t, batchState.BlobHeaders, 3)
	assert.ElementsMatch(t, batchState.BlobMetadata, []*disperser.BlobMetadata{encoded1.BlobMetadata, encoded2.BlobMetadata, encoded3.BlobMetadata})
	assert.NotNil(t, batchState.OperatorState)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Nil(t, c.minibatcher.Batches[b.ID])

	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 6)
	dispersalRequests, err := c.minibatchStore.GetMinibatchDispersalRequests(ctx, b.ID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersalRequests, 2)
	opIDs := make([]core.OperatorID, 2)
	for i, req := range dispersalRequests {
		assert.Equal(t, req.BatchID, b.ID)
		assert.Equal(t, req.MinibatchIndex, uint(0))
		assert.Equal(t, req.NumBlobs, uint(2))
		assert.NotNil(t, req.Socket)
		assert.NotNil(t, req.RequestedAt)
		opIDs[i] = req.OperatorID
	}
	assert.ElementsMatch(t, opIDs, []core.OperatorID{opId0, opId1})

	dispersalResponses, err := c.minibatchStore.GetMinibatchDispersalResponses(ctx, b.ID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersalResponses, 2)
	for _, resp := range dispersalResponses {
		assert.Equal(t, resp.BatchID, b.ID)
		assert.Equal(t, resp.MinibatchIndex, uint(0))
		assert.NotNil(t, resp.RespondedAt)
		assert.NoError(t, resp.Error)
		assert.Len(t, resp.Signatures, 1)
	}
}

func TestDisperseMinibatchFailure(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	var X, Y fp.Element
	X = *X.SetBigInt(big.NewInt(1))
	Y = *Y.SetBigInt(big.NewInt(2))
	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*core.Signature{
		{
			G1Point: &core.G1Point{
				G1Affine: &bn254.G1Affine{
					X: X,
					Y: Y,
				},
			},
		},
	}, nil)
	ctx := context.Background()

	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 100,
	}})
	_, _ = queueBlob(t, ctx, &blob1, c.blobStore)
	_, _ = queueBlob(t, ctx, &blob2, c.blobStore)

	// Start the batcher
	out := make(chan batcher.EncodingResultOrStatus)
	err := c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, <-out)
	assert.NoError(t, err)
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, <-out)
	assert.NoError(t, err)
	count, _ := c.encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, 2, count)

	err = c.minibatcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(1))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)

	b, err := c.minibatchStore.GetBatch(ctx, c.minibatcher.CurrentBatchID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, c.minibatcher.CurrentBatchID, b.ID)
	assert.NotNil(t, b.HeaderHash)
	assert.NotNil(t, b.CreatedAt)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, b.ReferenceBlockNumber)
	mb, err := c.minibatchStore.GetMinibatch(ctx, c.minibatcher.CurrentBatchID, 0)
	assert.NoError(t, err)
	assert.NotNil(t, mb)
	assert.Equal(t, c.minibatcher.CurrentBatchID, mb.BatchID)
	assert.Equal(t, uint(0), mb.MinibatchIndex)
	assert.Len(t, mb.BlobHeaderHashes, 2)
	assert.Equal(t, uint64(12800), mb.BatchSize)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, mb.ReferenceBlockNumber)

	c.pool.StopWait()
	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 2)
	dispersalRequests, err := c.minibatchStore.GetMinibatchDispersalRequests(ctx, c.minibatcher.CurrentBatchID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersalRequests, 2)
	opIDs := make([]core.OperatorID, 2)
	for i, req := range dispersalRequests {
		assert.Equal(t, req.BatchID, c.minibatcher.CurrentBatchID)
		assert.Equal(t, req.MinibatchIndex, uint(0))
		assert.Equal(t, req.NumBlobs, uint(2))
		assert.NotNil(t, req.Socket)
		assert.NotNil(t, req.RequestedAt)
		opIDs[i] = req.OperatorID
	}
	assert.ElementsMatch(t, opIDs, []core.OperatorID{opId0, opId1})

	dispersalResponses, err := c.minibatchStore.GetMinibatchDispersalResponses(ctx, c.minibatcher.CurrentBatchID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersalResponses, 2)
	for _, resp := range dispersalResponses {
		assert.Equal(t, resp.BatchID, c.minibatcher.CurrentBatchID)
		assert.Equal(t, resp.MinibatchIndex, uint(0))
		assert.NotNil(t, resp.RespondedAt)
		assert.NoError(t, resp.Error)
		assert.Len(t, resp.Signatures, 1)
	}
}

func TestSendBlobsToOperatorWithRetries(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	var X, Y fp.Element
	X = *X.SetBigInt(big.NewInt(1))
	Y = *Y.SetBigInt(big.NewInt(2))
	sig := &core.Signature{
		G1Point: &core.G1Point{
			G1Affine: &bn254.G1Affine{
				X: X,
				Y: Y,
			},
		},
	}
	ctx := context.Background()

	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 100,
	}})
	_, _ = queueBlob(t, ctx, &blob1, c.blobStore)
	_, _ = queueBlob(t, ctx, &blob2, c.blobStore)

	// Start the batcher
	out := make(chan batcher.EncodingResultOrStatus)
	err := c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, <-out)
	assert.NoError(t, err)
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, <-out)
	assert.NoError(t, err)
	count, _ := c.encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, 2, count)
	batch, err := c.encodingStreamer.CreateMinibatch(ctx)
	assert.NoError(t, err)

	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail")).Twice()
	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*core.Signature{sig}, nil).Once()
	signatures, err := c.minibatcher.SendBlobsToOperatorWithRetries(ctx, batch.EncodedBlobs, batch.BatchHeader, batch.State.IndexedOperators[opId0], opId0, 3)
	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 3)
	assert.NoError(t, err)
	assert.Len(t, signatures, 1)

	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fail")).Times(3)
	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*core.Signature{sig}, nil).Once()
	signatures, err = c.minibatcher.SendBlobsToOperatorWithRetries(ctx, batch.EncodedBlobs, batch.BatchHeader, batch.State.IndexedOperators[opId1], opId1, 3)
	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 6)
	assert.ErrorContains(t, err, "failed to send chunks to operator")
	assert.Nil(t, signatures)
}

func TestMinibatcherTooManyPendingRequests(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	ctx := context.Background()
	mockWorkerPool := &cmock.MockWorkerpool{}
	// minibatcher with mock worker pool
	m, err := batcher.NewMinibatcher(defaultConfig, c.blobStore, c.minibatchStore, c.dispatcher, c.minibatcher.ChainState, c.assignmentCoordinator, c.encodingStreamer, c.ethClient, mockWorkerPool, c.logger)
	assert.NoError(t, err)
	mockWorkerPool.On("WaitingQueueSize").Return(int(defaultConfig.MaxNumConnections + 1)).Once()
	err = m.HandleSingleBatch(ctx)
	assert.ErrorContains(t, err, "too many pending requests")
}
