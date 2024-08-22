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
	chainState            *core.IndexedOperatorState
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
	ics, err := chainState.GetIndexedOperatorState(context.Background(), 0, []core.QuorumID{0, 1})
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

	_, err = c.minibatcher.HandleSingleMinibatch(ctx)
	assert.NoError(t, err)

	// Check the minibatcher state
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(1))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobHeaders, 2)
	assert.ElementsMatch(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BlobMetadata, []*disperser.BlobMetadata{encoded1.BlobMetadata, encoded2.BlobMetadata})

	// Check the dispersal records
	dispersal, err := c.minibatchStore.GetDispersal(ctx, c.minibatcher.CurrentBatchID, 0, opId0)
	assert.NoError(t, err)
	assert.Equal(t, dispersal.BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, dispersal.MinibatchIndex, uint(0))
	assert.Equal(t, dispersal.OperatorID, opId0)
	assert.Equal(t, dispersal.Socket, c.chainState.IndexedOperators[opId0].Socket)
	assert.Equal(t, dispersal.NumBlobs, uint(2))
	assert.NotNil(t, dispersal.RequestedAt)

	dispersal, err = c.minibatchStore.GetDispersal(ctx, c.minibatcher.CurrentBatchID, 0, opId1)
	assert.NoError(t, err)
	assert.Equal(t, dispersal.BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, dispersal.MinibatchIndex, uint(0))
	assert.Equal(t, dispersal.OperatorID, opId1)
	assert.Equal(t, dispersal.Socket, c.chainState.IndexedOperators[opId1].Socket)
	assert.Equal(t, dispersal.NumBlobs, uint(2))
	assert.NotNil(t, dispersal.RequestedAt)

	// Check the blob minibatch mappings
	blobKey1 := encoded1.BlobMetadata.GetBlobKey()
	blobMinibatchMappings, err := c.minibatchStore.GetBlobMinibatchMappings(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Len(t, blobMinibatchMappings, 1)
	mapping1 := blobMinibatchMappings[0]
	assert.Equal(t, mapping1.BlobKey, &blobKey1)
	assert.Equal(t, mapping1.BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, mapping1.MinibatchIndex, uint(0))
	assert.Equal(t, mapping1.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{encoded1.BlobQuorumInfo})
	serializedCommitment1, err := encoded1.Commitment.Commitment.Serialize()
	assert.NoError(t, err)
	expectedCommitment1, err := mapping1.Commitment.Serialize()
	assert.NoError(t, err)
	serializedLengthCommitment1, err := encoded1.Commitment.LengthCommitment.Serialize()
	assert.NoError(t, err)
	expectedLengthCommitment1, err := mapping1.LengthCommitment.Serialize()
	assert.NoError(t, err)
	serializedLengthProof1, err := encoded1.Commitment.LengthProof.Serialize()
	assert.NoError(t, err)
	expectedLengthProof1, err := mapping1.LengthProof.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, serializedCommitment1, expectedCommitment1)
	assert.Equal(t, serializedLengthCommitment1, expectedLengthCommitment1)
	assert.Equal(t, serializedLengthProof1, expectedLengthProof1)
	blobKey2 := encoded1.BlobMetadata.GetBlobKey()
	blobMinibatchMappings, err = c.minibatchStore.GetBlobMinibatchMappings(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Len(t, blobMinibatchMappings, 1)
	mapping2 := blobMinibatchMappings[0]
	assert.Equal(t, mapping2.BlobKey, &blobKey2)
	assert.Equal(t, mapping2.BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, mapping2.MinibatchIndex, uint(0))
	assert.Equal(t, mapping2.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{encoded1.BlobQuorumInfo})
	if mapping1.BlobIndex != 0 {
		assert.Equal(t, mapping2.BlobIndex, uint(1))
	} else if mapping1.BlobIndex != 1 {
		assert.Equal(t, mapping2.BlobIndex, uint(0))
	} else {
		t.Fatal("invalid blob index")
	}
	serializedCommitment2, err := encoded2.Commitment.Commitment.Serialize()
	assert.NoError(t, err)
	expectedCommitment2, err := mapping2.Commitment.Serialize()
	assert.NoError(t, err)
	serializedLengthCommitment2, err := encoded2.Commitment.LengthCommitment.Serialize()
	assert.NoError(t, err)
	expectedLengthCommitment2, err := mapping2.LengthCommitment.Serialize()
	assert.NoError(t, err)
	serializedLengthProof2, err := encoded2.Commitment.LengthProof.Serialize()
	assert.NoError(t, err)
	expectedLengthProof2, err := mapping2.LengthProof.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, serializedCommitment2, expectedCommitment2)
	assert.Equal(t, serializedLengthCommitment2, expectedLengthCommitment2)
	assert.Equal(t, serializedLengthProof2, expectedLengthProof2)

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
	_, err = c.minibatcher.HandleSingleMinibatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(2))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Equal(t, c.minibatcher.Batches[c.minibatcher.CurrentBatchID].BatchID, c.minibatcher.CurrentBatchID)
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

	blobKey3 := encoded3.BlobMetadata.GetBlobKey()
	blobMinibatchMappings, err = c.minibatchStore.GetBlobMinibatchMappings(ctx, blobKey3)
	assert.NoError(t, err)
	assert.Len(t, blobMinibatchMappings, 1)
	mapping3 := blobMinibatchMappings[0]
	assert.Equal(t, mapping3.BlobKey, &blobKey3)
	assert.Equal(t, mapping3.BatchID, c.minibatcher.CurrentBatchID)
	assert.Equal(t, mapping3.MinibatchIndex, uint(1))
	assert.Equal(t, mapping3.BlobIndex, uint(0))

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
	_, err = c.minibatcher.HandleSingleMinibatch(ctx)
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
	assert.Len(t, batchState.BlobHeaders, 3)
	assert.ElementsMatch(t, batchState.BlobMetadata, []*disperser.BlobMetadata{encoded1.BlobMetadata, encoded2.BlobMetadata, encoded3.BlobMetadata})
	assert.NotNil(t, batchState.OperatorState)
	assert.Len(t, c.minibatcher.Batches, 1)
	assert.Nil(t, c.minibatcher.Batches[b.ID])

	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 6)
	dispersals, err := c.minibatchStore.GetDispersalsByMinibatch(ctx, b.ID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersals, 2)
	opIDs := make([]core.OperatorID, 2)
	for i, dispersal := range dispersals {
		assert.Equal(t, dispersal.BatchID, b.ID)
		assert.Equal(t, dispersal.MinibatchIndex, uint(0))
		assert.Equal(t, dispersal.NumBlobs, uint(2))
		assert.NotNil(t, dispersal.Socket)
		assert.NotNil(t, dispersal.RequestedAt)
		opIDs[i] = dispersal.OperatorID
		assert.NotNil(t, dispersal.RespondedAt)
		assert.NoError(t, dispersal.Error)
		assert.Len(t, dispersal.Signatures, 1)
	}
	assert.ElementsMatch(t, opIDs, []core.OperatorID{opId0, opId1})
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

	_, err = c.minibatcher.HandleSingleMinibatch(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, c.minibatcher.CurrentBatchID)
	assert.Equal(t, c.minibatcher.MinibatchIndex, uint(1))
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, initialBlock)

	b, err := c.minibatchStore.GetBatch(ctx, c.minibatcher.CurrentBatchID)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, c.minibatcher.CurrentBatchID, b.ID)
	assert.NotNil(t, b.CreatedAt)
	assert.Equal(t, c.minibatcher.ReferenceBlockNumber, b.ReferenceBlockNumber)

	c.pool.StopWait()
	c.dispatcher.AssertNumberOfCalls(t, "SendBlobsToOperator", 2)
	dispersals, err := c.minibatchStore.GetDispersalsByMinibatch(ctx, c.minibatcher.CurrentBatchID, 0)
	assert.NoError(t, err)
	assert.Len(t, dispersals, 2)
	opIDs := make([]core.OperatorID, 2)
	for i, dispersal := range dispersals {
		assert.Equal(t, dispersal.BatchID, c.minibatcher.CurrentBatchID)
		assert.Equal(t, dispersal.MinibatchIndex, uint(0))
		assert.Equal(t, dispersal.NumBlobs, uint(2))
		assert.NotNil(t, dispersal.Socket)
		assert.NotNil(t, dispersal.RequestedAt)
		opIDs[i] = dispersal.OperatorID
		assert.NotNil(t, dispersal.RespondedAt)
		assert.NoError(t, dispersal.Error)
		assert.Len(t, dispersal.Signatures, 1)
	}
	assert.ElementsMatch(t, opIDs, []core.OperatorID{opId0, opId1})
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

func TestSendBlobsToOperatorWithRetriesCanceled(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	ctx := context.Background()

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	_, _ = queueBlob(t, ctx, &blob, c.blobStore)

	out := make(chan batcher.EncodingResultOrStatus)
	err := c.encodingStreamer.RequestEncoding(ctx, out)
	assert.NoError(t, err)
	err = c.encodingStreamer.ProcessEncodedBlobs(ctx, <-out)
	assert.NoError(t, err)
	batch, err := c.encodingStreamer.CreateMinibatch(ctx)
	assert.NoError(t, err)
	minibatchIndex := uint(12)
	c.dispatcher.On("SendBlobsToOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, context.Canceled)
	c.minibatcher.DisperseBatch(ctx, batch.State, batch.EncodedBlobs, batch.BatchHeader, c.minibatcher.CurrentBatchID, minibatchIndex)
	c.pool.StopWait()
	dispersals, err := c.minibatchStore.GetDispersalsByMinibatch(ctx, c.minibatcher.CurrentBatchID, minibatchIndex)
	assert.NoError(t, err)
	assert.Len(t, dispersals, 2)

	indexedState, err := mockChainState.GetIndexedOperatorState(ctx, initialBlock, []core.QuorumID{0})
	assert.NoError(t, err)
	assert.Len(t, dispersals, len(indexedState.IndexedOperators))
	for _, dispersal := range dispersals {
		assert.ErrorContains(t, dispersal.Error, "context canceled")
		assert.GreaterOrEqual(t, dispersal.RespondedAt, dispersal.RequestedAt)
	}
}

func TestMinibatcherTooManyPendingRequests(t *testing.T) {
	c := newMinibatcher(t, defaultConfig)
	ctx := context.Background()
	mockWorkerPool := &cmock.MockWorkerpool{}
	// minibatcher with mock worker pool
	m, err := batcher.NewMinibatcher(defaultConfig, c.blobStore, c.minibatchStore, c.dispatcher, c.minibatcher.ChainState, c.assignmentCoordinator, c.encodingStreamer, c.ethClient, mockWorkerPool, c.logger)
	assert.NoError(t, err)
	mockWorkerPool.On("WaitingQueueSize").Return(int(defaultConfig.MaxNumConnections + 1)).Once()
	_, err = m.HandleSingleMinibatch(ctx)
	assert.ErrorContains(t, err, "too many pending requests")
}
