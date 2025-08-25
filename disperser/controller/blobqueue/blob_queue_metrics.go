package blobqueue

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Encapsulates all metrics for the blob queue.
type blobQueueMetrics struct {
}

// Create new metrics for the blob queue.
func newBlobQueueMetrics(registry *prometheus.Registry) *blobQueueMetrics {
	return nil
}

// Report the number of blobs waiting in the queue.
func (m *blobQueueMetrics) reportQueueSize(size uint64) {
	if m == nil {
		return
	}

	// TODO
}

// This should be called each time we poll the blob source.
func (m *blobQueueMetrics) reportBlobSourcePoll(blobCount uint64, duplicateBlobCount uint64, duration time.Duration) {
	if m == nil {
		return
	}

	// TODO
}

// Report the size of the deduplication set.
func (m *blobQueueMetrics) reportDedupSetSize(size uint64) {
	if m == nil {
		return
	}

	// TODO
}

// Report the time a blob spent in the queue before being acquired by the controller for batching.
func (m *blobQueueMetrics) reportTimeInQueue(duration time.Duration) {
	if m == nil {
		return
	}

	// TODO
}
