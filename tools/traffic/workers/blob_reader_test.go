package workers

import (
	"context"
	"crypto/md5"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/stretchr/testify/assert"
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

	lock := sync.Mutex{}
	chainClient := newMockChainClient(&lock)
	retrievalClient := newMockRetrievalClient(t, &lock)

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

		retrievalClient.AddBlob(blobMetadata, blobData)

		blobTable.Add(blobMetadata)
	}

	// Do a bunch of reads.
	expectedTotalReads := uint(readPermits * blobCount)
	for i := uint(0); i < expectedTotalReads; i++ {
		blobReader.randomRead()

		tu.AssertEventuallyTrue(t, func() bool {
			return retrievalClient.RetrieveBlobChunksCount == i+1 &&
				retrievalClient.CombineChunksCount == i+1 &&
				chainClient.Count == i+1
		}, time.Second)

		remainingPermits := uint(0)
		for _, metadata := range blobTable.GetAll() {
			remainingPermits += uint(metadata.RemainingReadPermits)
		}
		assert.Equal(t, remainingPermits, expectedTotalReads-i-1)

		tu.AssertEventuallyTrue(t, func() bool {
			return uint(readerMetrics.GetCount("read_success")) == i+1 &&
				uint(readerMetrics.GetCount("fetch_batch_header_success")) == i+1 &&
				uint(readerMetrics.GetCount("recombination_success")) == i+1
		}, time.Second)
	}

	expectedInvalidBlobs := uint(invalidBlobCount * readPermits)
	expectedValidBlobs := expectedTotalReads - expectedInvalidBlobs
	tu.AssertEventuallyEquals(t, expectedValidBlobs,
		func() any {
			return uint(readerMetrics.GetCount("valid_blob"))
		}, time.Second)
	tu.AssertEventuallyEquals(t, expectedInvalidBlobs,
		func() any {
			return uint(readerMetrics.GetCount("invalid_blob"))
		}, time.Second)

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
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
