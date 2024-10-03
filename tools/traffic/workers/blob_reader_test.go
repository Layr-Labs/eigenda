package workers

import (
	"context"
	"crypto/md5"
	"github.com/Layr-Labs/eigenda/api/clients"
	apiMock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	retrieverMock "github.com/Layr-Labs/eigenda/retriever/mock"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

// TestBlobReaderNoOptionalReads tests the BlobReader's basic functionality'
func TestBlobReader(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)

	blobTable := table.NewBlobStore()

	readerMetrics := metrics.NewMockMetrics()

	chainClient := &retrieverMock.MockChainClient{}
	chainClient.On(
		"FetchBatchHeader",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(&binding.IEigenDAServiceManagerBatchHeader{}, nil)
	retrievalClient := &apiMock.MockRetrievalClient{}

	blobReader := NewBlobReader(
		&ctx,
		&waitGroup,
		logger,
		&config.WorkerConfig{},
		retrievalClient,
		chainClient,
		blobTable,
		readerMetrics)

	blobSize := 1024
	readPermits := 2
	blobCount := 100

	invalidBlobCount := 0

	// Insert some blobs into the table.
	for i := 0; i < blobCount; i++ {

		key := make([]byte, 32)
		_, err = rand.Read(key)
		assert.Nil(t, err)

		blobData := make([]byte, blobSize)
		_, err = rand.Read(blobData)
		assert.Nil(t, err)

		var checksum [16]byte
		if i%10 == 0 {
			// Simulate an invalid blob
			invalidBlobCount++
			_, err = rand.Read(checksum[:])
			assert.Nil(t, err)
		} else {
			checksum = md5.Sum(blobData)
		}

		batchHeaderHash := [32]byte{}
		_, err = rand.Read(batchHeaderHash[:])
		assert.Nil(t, err)

		blobMetadata, err := table.NewBlobMetadata(
			key,
			checksum,
			uint(blobSize),
			uint(i),
			batchHeaderHash,
			readPermits)
		assert.Nil(t, err)

		// Simplify tracking by hijacking the BlobHeaderLength field to store the blob index,
		// which is used as a unique identifier within this test.
		chunks := &clients.BlobChunks{BlobHeaderLength: blobMetadata.BlobIndex}
		retrievalClient.On("RetrieveBlobChunks",
			blobMetadata.BatchHeaderHash,
			uint32(blobMetadata.BlobIndex),
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(chunks, nil)
		retrievalClient.On("CombineChunks", chunks).Return(blobData, nil)

		blobTable.Add(blobMetadata)
	}

	// Do a bunch of reads.
	expectedTotalReads := uint(readPermits * blobCount)
	for i := uint(0); i < expectedTotalReads; i++ {
		blobReader.randomRead()

		chainClient.AssertNumberOfCalls(t, "FetchBatchHeader", int(i+1))
		retrievalClient.AssertNumberOfCalls(t, "RetrieveBlobChunks", int(i+1))
		retrievalClient.AssertNumberOfCalls(t, "CombineChunks", int(i+1))

		remainingPermits := uint(0)
		for _, metadata := range blobTable.GetAll() {
			remainingPermits += uint(metadata.RemainingReadPermits)
		}
		assert.Equal(t, remainingPermits, expectedTotalReads-i-1)

		assert.Equal(t, i+1, uint(readerMetrics.GetCount("read_success")))
		assert.Equal(t, i+1, uint(readerMetrics.GetCount("fetch_batch_header_success")))
		assert.Equal(t, i+1, uint(readerMetrics.GetCount("recombination_success")))
	}

	expectedInvalidBlobs := uint(invalidBlobCount * readPermits)
	expectedValidBlobs := expectedTotalReads - expectedInvalidBlobs

	assert.Equal(t, expectedValidBlobs, uint(readerMetrics.GetCount("valid_blob")))
	assert.Equal(t, expectedInvalidBlobs, uint(readerMetrics.GetCount("invalid_blob")))
	assert.Equal(t, uint(0), uint(readerMetrics.GetGaugeValue("required_read_pool_size")))
	assert.Equal(t, uint(0), uint(readerMetrics.GetGaugeValue("optional_read_pool_size")))

	// Table is empty, so doing a random read should have no effect.
	blobReader.randomRead()

	// Give the system a moment to attempt to do work. This should not result in any reads.
	time.Sleep(time.Second / 10)
	assert.Equal(t, expectedTotalReads, uint(readerMetrics.GetCount("read_success")))
	assert.Equal(t, expectedTotalReads, uint(readerMetrics.GetCount("fetch_batch_header_success")))
	assert.Equal(t, expectedTotalReads, uint(readerMetrics.GetCount("recombination_success")))
	assert.Equal(t, expectedValidBlobs, uint(readerMetrics.GetCount("valid_blob")))
	assert.Equal(t, expectedInvalidBlobs, uint(readerMetrics.GetCount("invalid_blob")))

	cancel()
}
