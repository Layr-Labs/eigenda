package batcher_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	cmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/inmem"
	"github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

var (
	streamerConfig = batcher.StreamerConfig{
		SRSOrder:               300000,
		EncodingRequestTimeout: 5 * time.Second,
		EncodingQueueLimit:     100,
		PoolSize:               5,
	}
)

const numOperators = 10

type components struct {
	blobStore     disperser.BlobStore
	chainDataMock *coremock.ChainDataMock
	encoderClient *disperser.LocalEncoderClient
}

func createEncodingStreamer(t *testing.T, initialBlockNumber uint, batchThreshold uint, streamerConfig batcher.StreamerConfig) (*batcher.EncodingStreamer, *components) {
	logger := &cmock.Logger{}
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.NewChainDataMock(numOperators)
	assert.Nil(t, err)
	enc, err := makeTestEncoder()
	assert.Nil(t, err)
	encoderClient := disperser.NewLocalEncoderClient(enc)
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), batchThreshold)

	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = initialBlockNumber

	return encodingStreamer, &components{
		blobStore:     blobStore,
		chainDataMock: cst,
		encoderClient: encoderClient,
	}
}

func TestEncodingQueueLimit(t *testing.T) {
	logger := &cmock.Logger{}
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.NewChainDataMock(numOperators)
	assert.Nil(t, err)
	encoderClient := mock.NewMockEncoderClient()
	wait := make(chan time.Time)
	encoderClient.On("EncodeBlob", tmock.Anything, tmock.Anything, tmock.Anything).WaitUntil(wait).Return(nil, nil, nil)
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), 100000)

	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = 10

	securityParams := []*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
	}}
	blob1Data := []byte{1, 2, 3, 4, 5}
	blob1 := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: blob1Data,
	}

	ctx := context.Background()
	key1, err := blobStore.StoreBlob(ctx, &blob1, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	blob2Data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	blob2 := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: blob2Data,
	}
	key2, err := blobStore.StoreBlob(ctx, &blob2, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	// EncodeBlob still running, so this should return without making a request
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)

	// EncodeBlob call returns
	wait <- time.Now()
	// second blob should not have been encoded
	encoderClient.AssertNumberOfCalls(t, "EncodeBlob", 1)
	encoderClient.AssertCalled(t, "EncodeBlob", tmock.Anything, blob1Data, tmock.Anything)
	encoderClient.AssertNotCalled(t, "EncodeBlob", tmock.Anything, blob2Data, tmock.Anything)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	res, err := encodingStreamer.EncodedBlobstore.GetEncodingResult(key1, 0)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	res, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(key2, 0)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// retry
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	wait <- time.Now()

	encoderClient.AssertNumberOfCalls(t, "EncodeBlob", 2)
	encoderClient.AssertCalled(t, "EncodeBlob", tmock.Anything, blob2Data, tmock.Anything)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	res, err = encodingStreamer.EncodedBlobstore.GetEncodingResult(key2, 0)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestBatchTrigger(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 200_000, streamerConfig)

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
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
	total := encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(131584))

	// don't notify yet
	assert.Len(t, encodingStreamer.EncodedSizeNotifier.Notify, 0)

	// Request encoding once more
	_, err = c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)

	total = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(131584)*2)

	// notify
	assert.Len(t, encodingStreamer.EncodedSizeNotifier.Notify, 1)
}

func TestStreamingEncoding(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 0, 1e12, streamerConfig)

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
	}})
	ctx := context.Background()
	metadataKey, err := c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)
	metadata, err := c.blobStore.GetBlobMetadata(ctx, metadataKey)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Processing, metadata.BlobStatus)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10), nil)

	out := make(chan batcher.EncodingResultOrStatus)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	isRequested := encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	total := encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(0))

	err = encodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.Nil(t, err)
	encodedResult, err := encodingStreamer.EncodedBlobstore.GetEncodingResult(metadataKey, core.QuorumID(0))
	assert.Nil(t, err)
	assert.NotNil(t, encodedResult)
	assert.Equal(t, encodedResult.BlobMetadata, metadata)
	assert.Equal(t, encodedResult.ReferenceBlockNumber, uint(10))
	assert.Equal(t, encodedResult.BlobQuorumInfo, &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:           0,
			AdversaryThreshold: 80,
			QuorumThreshold:    100,
		},
		QuantizationFactor: batcher.QuantizationFactor,
		EncodedBlobLength:  320,
	})
	assert.NotNil(t, encodedResult.Commitment)
	assert.NotNil(t, encodedResult.Commitment.Commitment)
	assert.NotNil(t, encodedResult.Commitment.LengthProof)
	assert.Greater(t, encodedResult.Commitment.Length, uint(0))
	assert.Len(t, encodedResult.Assignments, numOperators)
	assert.Len(t, encodedResult.Chunks, 16)
	isRequested = encodingStreamer.EncodedBlobstore.HasEncodingRequested(metadataKey, core.QuorumID(0), 10)
	assert.True(t, isRequested)
	total = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(131584))

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
	total = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(131584))

	// Request the same blob, which should be dedupped
	_, err = c.blobStore.StoreBlob(ctx, &blob, requestedAt)
	assert.Nil(t, err)
	err = encodingStreamer.RequestEncoding(context.Background(), out)
	assert.Nil(t, err)
	assert.Equal(t, len(out), 0)
	// It should not have been added to the encoded blob store
	total = encodingStreamer.EncodedBlobstore.GetEncodedResultSize()
	assert.Equal(t, total, uint(131584))
}

func TestEncodingFailure(t *testing.T) {
	logger := &cmock.Logger{}
	blobStore := inmem.NewBlobStore()
	cst, err := coremock.NewChainDataMock(numOperators)
	assert.Nil(t, err)
	encoderClient := mock.NewMockEncoderClient()
	asgn := &core.StdAssignmentCoordinator{}
	sizeNotifier := batcher.NewEncodedSizeNotifier(make(chan struct{}, 1), 1e12)
	streamerConfig := batcher.StreamerConfig{
		SRSOrder:               300000,
		EncodingRequestTimeout: 5 * time.Second,
		EncodingQueueLimit:     100,
		PoolSize:               5,
	}

	encodingStreamer, err := batcher.NewEncodingStreamer(streamerConfig, blobStore, cst, encoderClient, asgn, sizeNotifier, logger)
	assert.Nil(t, err)
	encodingStreamer.ReferenceBlockNumber = 10

	ctx := context.Background()

	// put a blob in the blobstore
	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
	}, {
		QuorumID:           1,
		AdversaryThreshold: 70,
		QuorumThreshold:    100,
	}})

	metadataKey, err := blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)

	cst.On("GetCurrentBlockNumber").Return(uint(10), nil)
	encoderClient.On("EncodeBlob", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil, fmt.Errorf("errrrr"))
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

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10), nil)

	out := make(chan batcher.EncodingResultOrStatus)

	ctx := context.Background()

	// put in first blob and request encoding
	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 75,
		QuorumThreshold:    100,
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
		QuorumID:           1,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
	}, {
		QuorumID:           2,
		AdversaryThreshold: 70,
		QuorumThreshold:    95,
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
	batch, err := encodingStreamer.CreateBatch()
	assert.Nil(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(0))

	// Check BatchHeader
	assert.NotNil(t, batch.BatchHeader)
	assert.Greater(t, len(batch.BatchHeader.BatchRoot), 0)
	assert.Equal(t, batch.BatchHeader.ReferenceBlockNumber, uint(10))

	// Check BatchMetadata
	assert.NotNil(t, batch.BatchMetadata)
	assert.Len(t, batch.BatchMetadata.QuorumInfos, 1)
	assert.Len(t, batch.BatchMetadata.QuorumInfos[0].Assignments, numOperators)
	assert.Equal(t, batch.BatchMetadata.QuorumInfos[0].QuantizationFactor, batcher.QuantizationFactor)

	assert.Equal(t, batch.BatchMetadata.QuorumInfos[0].Info.TotalChunks, uint(15))
	assert.ElementsMatch(t, batch.BlobMetadata[0].RequestMetadata.SecurityParams, blob1.RequestHeader.SecurityParams)

	// Check EncodedBlobs
	assert.Len(t, batch.EncodedBlobs, 1)
	assert.Len(t, batch.EncodedBlobs[0], numOperators)

	encodedBlob1 := batch.EncodedBlobs[0]
	assert.NotNil(t, encodedBlob1)

	for _, blobMessage := range encodedBlob1 {
		assert.NotNil(t, blobMessage)
		assert.NotNil(t, blobMessage.BlobHeader)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.Commitment)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.LengthProof)
		assert.Equal(t, blobMessage.BlobHeader.BlobCommitments.Length, uint(48))
		assert.Len(t, blobMessage.BlobHeader.QuorumInfos, 1)
		assert.ElementsMatch(t, blobMessage.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{{
			SecurityParam: core.SecurityParam{
				QuorumID:           0,
				AdversaryThreshold: 75,
				QuorumThreshold:    100,
			},
			QuantizationFactor: batcher.QuantizationFactor,
			EncodedBlobLength:  160,
		}})

		assert.Contains(t, batch.BlobHeaders, blobMessage.BlobHeader)
		assert.Len(t, blobMessage.Bundles, 1)
		assert.Greater(t, len(blobMessage.Bundles[0]), 0)
		break
	}

	assert.Len(t, batch.BlobHeaders, 1)
	assert.Len(t, batch.BlobMetadata, 1)
	assert.Contains(t, batch.BlobMetadata, metadata1)
}

func TestIncorrectParameters(t *testing.T) {

	ctx := context.Background()

	streamerConfig := batcher.StreamerConfig{
		SRSOrder:               3000,
		EncodingRequestTimeout: 5 * time.Second,
		EncodingQueueLimit:     100,
		PoolSize:               5,
	}

	encodingStreamer, c := createEncodingStreamer(t, 0, 1e12, streamerConfig)

	// put a blob in the blobstore

	// The blob size is acceptable with the first security parameter but too large with the second
	// security parameter. Thus, the entire blob should be rejected.

	blob := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 50,
		QuorumThreshold:    100,
	}, {
		QuorumID:           1,
		AdversaryThreshold: 90,
		QuorumThreshold:    100,
	}})
	blob.Data = make([]byte, 10000)
	_, err := rand.Read(blob.Data)
	assert.NoError(t, err)

	metadataKey, err := c.blobStore.StoreBlob(ctx, &blob, uint64(time.Now().UnixNano()))
	assert.Nil(t, err)

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10), nil)

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

func TestGetBatch(t *testing.T) {
	encodingStreamer, c := createEncodingStreamer(t, 10, 1e12, streamerConfig)
	ctx := context.Background()

	// put 2 blobs in the blobstore
	blob1 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           0,
		AdversaryThreshold: 80,
		QuorumThreshold:    100,
	}, {
		QuorumID:           1,
		AdversaryThreshold: 70,
		QuorumThreshold:    95,
	}})
	blob2 := makeTestBlob([]*core.SecurityParam{{
		QuorumID:           2,
		AdversaryThreshold: 75,
		QuorumThreshold:    100,
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

	c.chainDataMock.On("GetCurrentBlockNumber").Return(uint(10), nil)

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
	batch, err := encodingStreamer.CreateBatch()
	assert.Nil(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, encodingStreamer.ReferenceBlockNumber, uint(0))

	// Check BatchHeader
	assert.NotNil(t, batch.BatchHeader)
	assert.Greater(t, len(batch.BatchHeader.BatchRoot), 0)
	assert.Equal(t, batch.BatchHeader.ReferenceBlockNumber, uint(10))

	// Check BatchMetadata
	assert.NotNil(t, batch.BatchMetadata)
	assert.Len(t, batch.BatchMetadata.QuorumInfos, 3)
	for quorumID := uint8(0); quorumID < 3; quorumID++ {
		assert.Len(t, batch.BatchMetadata.QuorumInfos[quorumID].Assignments, numOperators)
		assert.Equal(t, batch.BatchMetadata.QuorumInfos[quorumID].QuantizationFactor, batcher.QuantizationFactor)
	}
	assert.Equal(t, batch.BatchMetadata.QuorumInfos[0].Info.TotalChunks, uint(15))
	assert.Equal(t, batch.BatchMetadata.QuorumInfos[1].Info.TotalChunks, uint(15))
	assert.Equal(t, batch.BatchMetadata.QuorumInfos[2].Info.TotalChunks, uint(15))

	// Check EncodedBlobs
	assert.Len(t, batch.EncodedBlobs, 2)
	assert.Len(t, batch.EncodedBlobs[0], numOperators)

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
	for _, blobMessage := range encodedBlob1 {
		assert.NotNil(t, blobMessage)
		assert.NotNil(t, blobMessage.BlobHeader)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.Commitment)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.LengthProof)
		assert.Equal(t, blobMessage.BlobHeader.BlobCommitments.Length, uint(48))
		assert.Len(t, blobMessage.BlobHeader.QuorumInfos, 2)
		assert.ElementsMatch(t, blobMessage.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           0,
					AdversaryThreshold: 80,
					QuorumThreshold:    100,
				},
				QuantizationFactor: batcher.QuantizationFactor,
				EncodedBlobLength:  320,
			},
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           1,
					AdversaryThreshold: 70,
					QuorumThreshold:    95,
				},
				QuantizationFactor: batcher.QuantizationFactor,
				EncodedBlobLength:  160,
			},
		})

		assert.Contains(t, batch.BlobHeaders, blobMessage.BlobHeader)
		assert.Len(t, blobMessage.Bundles, 2)
		assert.Greater(t, len(blobMessage.Bundles[0]), 0)
		assert.Greater(t, len(blobMessage.Bundles[1]), 0)
		break
	}

	for _, blobMessage := range encodedBlob2 {
		assert.NotNil(t, blobMessage)
		assert.NotNil(t, blobMessage.BlobHeader)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.Commitment)
		assert.NotNil(t, blobMessage.BlobHeader.BlobCommitments.LengthProof)
		assert.Equal(t, blobMessage.BlobHeader.BlobCommitments.Length, uint(48))
		assert.Len(t, blobMessage.BlobHeader.QuorumInfos, 1)
		assert.ElementsMatch(t, blobMessage.BlobHeader.QuorumInfos, []*core.BlobQuorumInfo{{
			SecurityParam: core.SecurityParam{
				QuorumID:           2,
				AdversaryThreshold: 75,
				QuorumThreshold:    100,
			},
			QuantizationFactor: batcher.QuantizationFactor,
			EncodedBlobLength:  160,
		}})

		assert.Len(t, blobMessage.Bundles, 1)
		assert.Greater(t, len(blobMessage.Bundles[core.QuorumID(2)]), 0)
		break
	}
	assert.Len(t, batch.BlobHeaders, 2)
	assert.Len(t, batch.BlobMetadata, 2)
	assert.Contains(t, batch.BlobMetadata, metadata1)
	assert.Contains(t, batch.BlobMetadata, metadata2)
}
