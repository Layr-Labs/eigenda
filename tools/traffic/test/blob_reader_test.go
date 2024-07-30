package test

import (
	"context"
	"crypto/md5"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/tools/traffic/metrics"
	"github.com/Layr-Labs/eigenda/tools/traffic/table"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"sync"
	"testing"
	"time"
)

// TestBlobReaderNoOptionalReads tests the BlobReader with only required reads.
func TestBlobReaderNoOverflow(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := newMockTicker(startTime)

	config := &workers.Config{
		ReadOverflowTableSize: 0,
	}

	blobTable := table.NewBlobTable()

	readerMetrics := metrics.NewMockMetrics()

	lock := sync.Mutex{}
	chainClient := newMockChainClient(&lock)
	retrievalClient := newMockRetrievalClient(t, &lock)

	blobReader := workers.NewBlobReader(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		retrievalClient,
		chainClient,
		&blobTable,
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

		var checksum [16]byte
		if i%10 == 0 {
			// Simulate an invalid blob
			invalidBlobCount++
			_, err = rand.Read(checksum[:])
		} else {
			checksum = md5.Sum(blobData)
		}

		batchHeaderHash := make([]byte, 32)
		_, err = rand.Read(batchHeaderHash)
		assert.Nil(t, err)

		blobMetadata := table.NewBlobMetadata(
			&key,
			&checksum,
			uint(blobSize),
			&batchHeaderHash,
			uint(i),
			readPermits)

		retrievalClient.AddBlob(blobMetadata, blobData)

		blobTable.Add(blobMetadata)
	}

	blobReader.Start()

	// Do a bunch of reads.
	expectedTotalReads := uint(readPermits * blobCount)
	for i := uint(0); i < expectedTotalReads; i++ {
		ticker.Tick(time.Second)

		tu.AssertEventuallyTrue(t, func() bool {
			return retrievalClient.RetrieveBlobChunksCount == i+1 &&
				retrievalClient.CombineChunksCount == i+1 &&
				chainClient.Count == i+1
		}, time.Second)

		remainingPermits := uint(0)
		for j := uint(0); j < blobTable.Size(); j++ {
			blob := blobTable.Get(j)
			if blob.RemainingReadPermits() != 2 {
			}
			remainingPermits += uint(blob.RemainingReadPermits())
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

	// Table is empty, so ticking time forward should not result in any reads.
	ticker.Tick(time.Second)
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

// TestBlobReaderWithOverflow tests the BlobReader with a non-zero sized overflow table.
func TestBlobReaderWithOverflow(t *testing.T) {
	tu.InitializeRandom()

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := sync.WaitGroup{}
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.Nil(t, err)
	startTime := time.Unix(rand.Int63()%2_000_000_000, 0)
	ticker := newMockTicker(startTime)

	blobCount := 100
	overflowTableSize := uint(rand.Intn(blobCount-1) + 1)

	config := &workers.Config{
		ReadOverflowTableSize: overflowTableSize,
	}

	blobTable := table.NewBlobTable()

	readerMetrics := metrics.NewMockMetrics()

	lock := sync.Mutex{}
	chainClient := newMockChainClient(&lock)
	retrievalClient := newMockRetrievalClient(t, &lock)

	blobReader := workers.NewBlobReader(
		&ctx,
		&waitGroup,
		logger,
		ticker,
		config,
		retrievalClient,
		chainClient,
		&blobTable,
		readerMetrics)

	blobSize := 1024
	readPermits := 2

	invalidBlobCount := 0

	// Insert some blobs into the table.
	for i := 0; i < blobCount; i++ {

		key := make([]byte, 32)
		_, err = rand.Read(key)
		assert.Nil(t, err)

		blobData := make([]byte, blobSize)
		_, err = rand.Read(blobData)

		var checksum [16]byte
		if i%10 == 0 {
			// Simulate an invalid blob
			invalidBlobCount++
			_, err = rand.Read(checksum[:])
		} else {
			checksum = md5.Sum(blobData)
		}

		batchHeaderHash := make([]byte, 32)
		_, err = rand.Read(batchHeaderHash)
		assert.Nil(t, err)

		blobMetadata := table.NewBlobMetadata(
			&key,
			&checksum,
			uint(blobSize),
			&batchHeaderHash,
			uint(i),
			readPermits)

		retrievalClient.AddBlob(blobMetadata, blobData)

		blobTable.Add(blobMetadata)
	}

	blobReader.Start()

	// Do a bunch of reads.
	expectedTotalReads := uint(readPermits * blobCount)
	for i := uint(0); i < expectedTotalReads; i++ {
		ticker.Tick(time.Second)

		tu.AssertEventuallyTrue(t, func() bool {
			return retrievalClient.RetrieveBlobChunksCount == i+1 &&
				retrievalClient.CombineChunksCount == i+1 &&
				chainClient.Count == i+1
		}, time.Second)

		remainingPermits := uint(0)
		for j := uint(0); j < blobTable.Size(); j++ {
			blob := blobTable.Get(j)
			if blob.RemainingReadPermits() != 2 {
			}
			remainingPermits += uint(blob.RemainingReadPermits())
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
	assert.Equal(t, overflowTableSize, uint(readerMetrics.GetGaugeValue("optional_read_pool_size")))

	// Do an additional read. We should be reading from the overflow table.
	ticker.Tick(time.Second)
	tu.AssertEventuallyEquals(t, expectedTotalReads+1, func() any {
		return uint(readerMetrics.GetCount("read_success"))
	}, time.Second)
	tu.AssertEventuallyTrue(t, func() bool {
		return uint(readerMetrics.GetCount("valid_blob")) == expectedValidBlobs+1 ||
			uint(readerMetrics.GetCount("invalid_blob")) == expectedInvalidBlobs+1
	}, time.Second)

	cancel()
	tu.ExecuteWithTimeout(func() {
		waitGroup.Wait()
	}, time.Second)
}
