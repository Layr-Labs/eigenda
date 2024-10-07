package batcher_test

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

var (
	streamerConfig = batcher.StreamerConfig{
		SRSOrder:                 300000,
		EncodingRequestTimeout:   5 * time.Second,
		EncodingQueueLimit:       100,
		MaxBlobsToFetchFromStore: 10,
		FinalizationBlockDelay:   75,
	}
)

const numOperators = 10

type components struct {
	blobStore     disperser.BlobStore
	chainDataMock *coremock.ChainDataMock
	encoderClient *disperser.LocalEncoderClient
}

func createEncodingStreamer(t *testing.T, initialBlockNumber uint, batchThreshold uint64, streamerConfig batcher.StreamerConfig) (*batcher.EncodingStreamer, *components) {
	logger := logging.NewNoopLogger()
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	assert.Nil(t, err)
	p, err := makeTestProver()
	assert.Nil(t, err)
	encoderClient := disperser.NewLocalEncoderClient(p)
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), batchThreshold)
	workerpool := workerpool.New(5)
	metrics := batcher.NewMetrics("9100", logger)
	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, workerpool, metrics.EncodingStreamerMetrics, metrics, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = initialBlockNumber

	return encodingStreamer, &components{
		blobStore:     blobStore,
		chainDataMock: cst,
		encoderClient: encoderClient,
	}
}

func TestEncodingQueueLimit(t *testing.T) {
	logger := logging.NewNoopLogger()
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	assert.Nil(t, err)
	encoderClient := mock.NewMockEncoderClient()
	encoderClient.On("EncodeBlob", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil, nil)
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), 100000)
	pool := &cmock.MockWorkerpool{}
	metrics := batcher.NewMetrics("9100", logger)
	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, pool, metrics.EncodingStreamerMetrics, metrics, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = 10

	securityParams := []*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}}
	blobData := []byte{1, 2, 3, 4, 5}
	blob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: blobData,
	}

	pool.On("Submit", tmock.Anything).Run(func(args tmock.Arguments) {
		args.Get(0).(func())()
	})

	// assume that encoding queue is already full
	pool.On("WaitingQueueSize").Return(streamerConfig.EncodingQueueLimit).Once()

	ctx := context.Background()
	key, err := blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	out := make(chan batcher.EncodingResultOrStatus, 1)
	// This should return without making a request since encoding queue was already full
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	encoderClient.AssertNotCalled(t, "EncodeBlob")
	select {
	case <-out:
		t.Fatal("did not expect any encoding results")
	default:
	}
	// assume that encoding queue opens up
	pool.On("WaitingQueueSize").Return(0).Once()

	// retry
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	encoderClient.AssertNumberOfCalls(t, "EncodeBlob", 1)
	encoderClient.AssertCalled(t, "EncodeBlob", tmock.Anything, blobData, tmock.Anything)
	var encodingResult batcher.EncodingResultOrStatus
	select {
	case encodingResult = <-out:
	default:
		t.Fatal("did not expect any encoding results")
	}

	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), encodingResult)
	assert.Nil(t, err)
	res, err := encodingStreamer.EncodedBlobstore.GetEncodingResult(key, 0)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestBatchTrigger(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 30_000, streamerConfig)

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	ctx := context.Background()
	_, err := c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	out := make(chan batcher.EncodingResultOrStatus)
	// Request encoding
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	count, size := encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 1)
	assert.Equal(t, size, uint64(26630))

	// try encode the same blobs again at different block (this happens when the blob is retried)
	encodingStreamer.ReferenceBlockNumber = 11
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)

	count, size = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 1)
	assert.Equal(t, size, uint64(26630))

	// don't notify yet
	select {
	case <-encodingStreamer.EncodedSizeNotifier.Notify:
		t.Fatal("expected not to be notified")
	default:
	}

	// Request encoding once more
	_, err = c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)

	count, size = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 2)
	assert.Equal(t, size, uint64(26630)*2)

	// notify
	select {
	case <-encodingStreamer.EncodedSizeNotifier.Notify:
	default:
		t.Fatal("expected to be notified")
	}
}

func TestStreamingEncoding(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 0, 1e12, streamerConfig)

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}})
	ctx := context.Background()
	metadataKey, err := c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata, err := c.blobStore.GetBlobMetadata(ctx, metadataKey)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata.BlobStatus)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)

	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	count, size := encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 0)
	assert.Equal(t, size, uint64(0))

	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	encodedResult, err := encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey, core.QuorumID(0))
	assert.Nil(t, err)
	assert.NotNil(t, encodedResult)
	assert.Equal(t, metadata, encodedResult.BlobMetadata)
	assert.Equal(t, uint(10), encodedResult.ReferenceBlockNumber)
	assert.Equal(t, &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:              0,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
		ChunkLength: 16,
	}, encodedResult.BlobQuorumInfo)
	assert.NotNil(t, encodedResult.Commitment)
	assert.NotNil(t, encodedResult.Commitment.Commitment)
	assert.NotNil(t, encodedResult.Commitment.LengthProof)
	assert.Greater(t, encodedResult.Commitment.Length, uint(0))
	assert.Len(t, encodedResult.Assignments, numOperators)
	assert.Len(t, encodedResult.ChunksData.Chunks, 32)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	count, size = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 1)
	assert.Equal(t, size, uint64(26630))

	// Cancel previous blob so it doesn't get reencoded.
	err = c.blobStore.MarkBlobFailed(ctx, metadataKey)
	assert.Nil(t, err)

	encodingStreamer.ReferenceBlockNumber = 11
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 11)
	assert.False(t, isRequested)
	// Request another blob again
	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err = c.blobStore.StoreBlob(ctx, &blob, requestedAt)
	assert.Nil(t, err)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	encodedResult, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey, core.QuorumID(0))
	assert.Nil(t, err)
	assert.NotNil(t, encodedResult)
	// This should delete the stale results but keep the new encoded results
	results := encodingStreamer.EncodedBlobstore.GetNewAndDeleteStaleEncodingResults(uint(11))
	assert.Len(t, results, 1)
	encodedResult, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey, core.QuorumID(0))
	assert.Nil(t, err)
	assert.NotNil(t, encodedResult)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 11)
	assert.True(t, isRequested)
	count, size = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 1)
	assert.Equal(t, size, uint64(26630))

	// Request the same blob, which should be dedupped
	_, err = c.blobStore.StoreBlob(ctx, &blob, requestedAt)
	assert.Nil(t, err)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	assert.Equal(t, len(out), 0)
	// It should not have been added to the encoded blob store
	count, size = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, count, 1)
	assert.Equal(t, size, uint64(26630))
}

func TestEncodingFailure(t *testing.T) {
	logger := logging.NewNoopLogger()
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	assert.Nil(t, err)
	encoderClient := mock.NewMockEncoderClient()
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), 1e12)
	workerpool := workerpool.New(5)
	streamerConfig := batcher.StreamerConfig{
		SRSOrder:                 300000,
		EncodingRequestTimeout:   5 * time.Second,
		EncodingQueueLimit:       100,
		MaxBlobsToFetchFromStore: 10,
	}
	metrics := batcher.NewMetrics("9100", logger)
	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, workerpool, metrics.EncodingStreamerMetrics, metrics, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = 10

	ctx := context.Background()

	// put a blob in the blobstore
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 100,
	}})

	metadataKey, err := blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)

	cst.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)
	encoderClient.On("EncodeBlob", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil, errors.New("errrrr"))
	// request encoding
	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(1), 10)
	assert.True(t, isRequested)

	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.NotNil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.NotNil(t, err)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 9)
	assert.False(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.False(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 11)
	assert.False(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(1), 10)
	assert.False(t, isRequested)
}

func TestPartialBlob(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 1e12, streamerConfig)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)

	out := make(chan batcher.EncodingResultOrStatus)

	ctx := context.Background()

	// put in first blob and request encoding
	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}})

	metadataKey1, err := c.blobStore.StoreBlob(ctx, &blob1, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata1, err := c.blobStore.GetBlobMetadata(ctx, metadataKey1)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata1.BlobStatus)

	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(0), 10)
	assert.True(t, isRequested)

	// Put in second blob and request encoding
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              1,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              2,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 95,
	}})
	metadataKey2, err := c.blobStore.StoreBlob(ctx, &blob2, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata2, err := c.blobStore.GetBlobMetadata(ctx, metadataKey2)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata2.BlobStatus)

	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(1), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(2), 10)
	assert.True(t, isRequested)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)

	// The second quorum doesn't complete
	<-out
	encodingStreamer.Pool.StopWait()

	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(1), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(2), 10)
	assert.True(t, isRequested)

	// get batch
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(10))
	batch, err := encodingStreamer.CreateBatch(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(0))

	// Check BatchHeader
	assert.NotNil(t, batch.BatchHeader)
	assert.Greater(t, len(batch.BatchHeader.BatchRoot), 0)
	assert.Equal(t, batch.BatchHeader.ReferenceBlockNumber, uint(10))

	// Check BatchMetadata
	assert.NotNil(t, batch.State)
	assert.ElementsMatch(t, batch.BlobMetadata[0].RequestMetadata.SecurityParams, blob1.RequestHeader.SecurityParams)

	// Check EncodedBlobs
	assert.Len(t, batch.EncodedBlobs, 1)
	assert.Len(t, batch.EncodedBlobs[0].EncodedBundlesByOperator, numOperators)

	encodedBlob1 := batch.EncodedBlobs[0]
	assert.NotNil(t, encodedBlob1)
	assert.NotNil(t, encodedBlob1.BlobHeader)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments.Commitment)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments.LengthProof)
	assert.Equal(t, encodedBlob1.BlobHeader.BlobCommitments.Length, uint(48))
	assert.Len(t, encodedBlob1.BlobHeader.QuorumInfos, 1)
	assert.ElementsMatch(t, encodedBlob1.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{{
		SecurityParam: core.SecurityParam{
			QuorumID:              0,
			AdversaryThreshold:    75,
			ConfirmationThreshold: 100,
		},
		ChunkLength: 8,
	}})

	assert.Contains(t, batch.BlobHeaders, encodedBlob1.BlobHeader)
	assert.Len(t, encodedBlob1.EncodedBundlesByOperator, numOperators)
	for _, bundles := range encodedBlob1.EncodedBundlesByOperator {
		assert.Len(t, bundles, 1)
		assert.Greater(t, len(bundles[0].Chunks), 0)
		break
	}

	assert.Len(t, batch.BlobHeaders, 1)
	assert.Len(t, batch.BlobMetadata, 1)
	assert.Contains(t, batch.BlobMetadata, metadata1)
}

func TestIncorrectParameters(t *testing.T) {

	ctx := context.Background()

	streamerConfig := batcher.StreamerConfig{
		SRSOrder:                 3000,
		EncodingRequestTimeout:   5 * time.Second,
		EncodingQueueLimit:       100,
		MaxBlobsToFetchFromStore: 10,
	}

	encodingStreamer, c := createEncodingStreamer(t, 0, 1e12, streamerConfig)

	// put a blob in the blobstore

	// The blob size is acceptable with the first security parameter but too large with the second
	// security parameter. Thus, the entire blob should be rejected.

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    50,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              1,
		AdversaryThreshold:    90,
		ConfirmationThreshold: 100,
	}})
	blob.Data = make([]byte, 10000)
	_, err := rand.Read(blob.Data)
	assert.NoError(t, err)

	metadataKey, err := c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)

	// request encoding
	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.False(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(1), 10)
	assert.False(t, isRequested)

	stats, err := c.blobStore.GetBlobMetadata(ctx, metadataKey)
	assert.NoError(t, err)
	assert.Equal(t, disperser.Failed, stats.BlobStatus)

}

func TestInvalidQuorum(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 1e12, streamerConfig)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)

	out := make(chan batcher.EncodingResultOrStatus)

	ctx := context.Background()

	// this blob should not be encoded because the quorum does not exist
	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              99, // this quorum does not exist
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}})

	// this blob should be encoded
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              1,
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}})

	metadataKey1, err := c.blobStore.StoreBlob(ctx, &blob1, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadataKey2, err := c.blobStore.StoreBlob(ctx, &blob2, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)

	// request encoding
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(0), 10)
	assert.False(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(99), 10)
	assert.False(t, isRequested)

	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(1), 10)
	assert.True(t, isRequested)

	stats, err := c.blobStore.GetBlobMetadata(ctx, metadataKey1)
	assert.NoError(t, err)
	assert.Equal(t, disperser.Failed, stats.BlobStatus)
}

func TestGetBatch(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 1e12, streamerConfig)
	ctx := context.Background()

	// put 2 blobs in the blobstore
	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              0,
		AdversaryThreshold:    80,
		ConfirmationThreshold: 100,
	}, {
		QuorumID:              1,
		AdversaryThreshold:    70,
		ConfirmationThreshold: 95,
	}})
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:              2,
		AdversaryThreshold:    75,
		ConfirmationThreshold: 100,
	}})
	metadataKey1, err := c.blobStore.StoreBlob(ctx, &blob1, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata1, err := c.blobStore.GetBlobMetadata(ctx, metadataKey1)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata1.BlobStatus)
	metadataKey2, err := c.blobStore.StoreBlob(ctx, &blob2, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata2, err := c.blobStore.GetBlobMetadata(ctx, metadataKey2)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata2.BlobStatus)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10)+encodingStreamer.FinalizationBlockDelay, nil)

	// request encoding
	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(1), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(2), 10)
	assert.True(t, isRequested)

	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	encodingStreamer.Pool.StopWait()

	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey1, core.QuorumID(1), 10)
	assert.True(t, isRequested)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey2, core.QuorumID(2), 10)
	assert.True(t, isRequested)

	// get batch
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(10))
	batch, err := encodingStreamer.CreateBatch(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(0))
	metadata1, err = c.blobStore.GetBlobMetadata(ctx, metadataKey1)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Dispersing, metadata1.BlobStatus)
	metadata2, err = c.blobStore.GetBlobMetadata(ctx, metadataKey2)
	assert.Equal(t, disperser.Dispersing, metadata2.BlobStatus)
	assert.Nil(t, err)
	res, err := encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey1, core.QuorumID(0))
	assert.Nil(t, res)
	assert.ErrorContains(t, err, "GetEncodedBlob: no such key")
	res, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey1, core.QuorumID(1))
	assert.Nil(t, res)
	assert.ErrorContains(t, err, "GetEncodedBlob: no such key")
	res, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey2, core.QuorumID(0))
	assert.Nil(t, res)
	assert.ErrorContains(t, err, "GetEncodedBlob: no such key")

	// Check BatchHeader
	assert.NotNil(t, batch.BatchHeader)
	assert.Greater(t, len(batch.BatchHeader.BatchRoot), 0)
	assert.Equal(t, batch.BatchHeader.ReferenceBlockNumber, uint(10))

	// Check State
	assert.NotNil(t, batch.State)

	// Check EncodedBlobs
	assert.Len(t, batch.EncodedBlobs, 2)
	assert.Len(t, batch.EncodedBlobs[0].EncodedBundlesByOperator, numOperators)

	var encodedBlob1 core.EncodedBlob
	var encodedBlob2 core.EncodedBlob
	for i := range batch.BlobHeaders {
		blobHeader := batch.BlobHeaders[i]
		if len(blobHeader.QuorumInfos) > 1 {
			encodedBlob1 = batch.EncodedBlobs[i]
			// batch.EncodedBlobs and batch.BlobMetadata should be in the same order
			assert.ElementsMatch(t, batch.BlobMetadata[i].RequestMetadata.SecurityParams, blob1.RequestHeader.SecurityParams)
		} else {
			encodedBlob2 = batch.EncodedBlobs[i]
			assert.ElementsMatch(t, batch.BlobMetadata[i].RequestMetadata.SecurityParams, blob2.RequestHeader.SecurityParams)
		}
	}
	assert.NotNil(t, encodedBlob1)
	assert.NotNil(t, encodedBlob2)

	assert.NotNil(t, encodedBlob1.BlobHeader)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments.Commitment)
	assert.NotNil(t, encodedBlob1.BlobHeader.BlobCommitments.LengthProof)
	assert.Equal(t, encodedBlob1.BlobHeader.BlobCommitments.Length, uint(48))
	assert.Len(t, encodedBlob1.BlobHeader.QuorumInfos, 2)
	assert.ElementsMatch(t, encodedBlob1.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{
		{
			SecurityParam: core.SecurityParam{
				QuorumID:              0,
				AdversaryThreshold:    80,
				ConfirmationThreshold: 100,
			},
			ChunkLength: 16,
		},
		{
			SecurityParam: core.SecurityParam{
				QuorumID:              1,
				AdversaryThreshold:    70,
				ConfirmationThreshold: 95,
			},
			ChunkLength: 8,
		},
	})

	assert.Contains(t, batch.BlobHeaders, encodedBlob1.BlobHeader)
	for _, bundles := range encodedBlob1.EncodedBundlesByOperator {
		assert.Len(t, bundles, 2)
		assert.Greater(t, len(bundles[0].Chunks), 0)
		assert.Greater(t, len(bundles[1].Chunks), 0)
		break
	}

	assert.NotNil(t, encodedBlob2.BlobHeader)
	assert.NotNil(t, encodedBlob2.BlobHeader.BlobCommitments)
	assert.NotNil(t, encodedBlob2.BlobHeader.BlobCommitments.Commitment)
	assert.NotNil(t, encodedBlob2.BlobHeader.BlobCommitments.LengthProof)
	assert.Equal(t, encodedBlob2.BlobHeader.BlobCommitments.Length, uint(48))
	assert.Len(t, encodedBlob2.BlobHeader.QuorumInfos, 1)
	assert.ElementsMatch(t, encodedBlob2.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{{
		SecurityParam: core.SecurityParam{
			QuorumID:              2,
			AdversaryThreshold:    75,
			ConfirmationThreshold: 100,
		},
		ChunkLength: 8,
	}})
	for _, bundles := range encodedBlob2.EncodedBundlesByOperator {
		assert.Len(t, bundles, 1)
		assert.Greater(t, len(bundles[core.QuorumID(2)].Chunks), 0)
		break
	}
	assert.Len(t, batch.BlobHeaders, 2)
	assert.Len(t, batch.BlobMetadata, 2)
	assert.Contains(t, batch.BlobMetadata, metadata1)
	assert.Contains(t, batch.BlobMetadata, metadata2)
}
