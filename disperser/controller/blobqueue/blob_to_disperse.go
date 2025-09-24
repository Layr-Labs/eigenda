package blobqueue

import (
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
)

// A wrapper around v2.BlobMetadata. Provides extra metrics/metadata.
type BlobToDisperse struct {
	// Encapsulates all metrics for the blob queue.
	metrics *blobQueueMetrics

	// The actual blob metadata.
	metadata *v2.BlobMetadata

	// The blob key of the blob, so we don't have to hash it multiple times.
	blobKey corev2.BlobKey

	// The time the blob was enqueued.
	enqueueTime time.Time

	// Tracks whether metrics have been reported for this blob yet. Metrics are not reported until
	// the controller calls GetBlobMetadata() for the first time.
	metricsReported bool
}

// Wrap the given blob metadata in a BlobToDisperse struct.
func newBlobToDisperse(metrics *blobQueueMetrics, metadata *v2.BlobMetadata, blobKey corev2.BlobKey) *BlobToDisperse {
	return &BlobToDisperse{
		metrics:         metrics,
		metadata:        metadata,
		blobKey:         blobKey,
		enqueueTime:     time.Now(),
		metricsReported: false,
	}
}

// GetBlobMetadata returns the underlying blob metadata. Also reports some metrics under the hood.
//
// This method is not thread-safe, and should only be called by one goroutine at a time.
func (b *BlobToDisperse) GetBlobMetadata() *v2.BlobMetadata {
	if !b.metricsReported {
		elapsed := time.Since(b.enqueueTime)
		b.metrics.reportTimeInQueue(elapsed)
		b.metricsReported = true
	}

	return b.metadata
}

// GetBlobKey returns the blob key of the blob.
func (b *BlobToDisperse) GetBlobKey() corev2.BlobKey {
	return b.blobKey
}
