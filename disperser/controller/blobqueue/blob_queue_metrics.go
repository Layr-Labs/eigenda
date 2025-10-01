package blobqueue

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const prefix = "eigenda_dispatcher_"

// Encapsulates all metrics for the blob queue.
type blobQueueMetrics struct {
	blobQueueSize        *prometheus.GaugeVec
	dedupSetSize         *prometheus.GaugeVec
	uniqueBlobCounter    *prometheus.CounterVec
	duplicateBlobCounter *prometheus.CounterVec
	timeSpentInQueue     *prometheus.SummaryVec
	pollLatency          *prometheus.SummaryVec
}

// Create new metrics for the blob queue.
func newBlobQueueMetrics(registry *prometheus.Registry) *blobQueueMetrics {
	if registry == nil {
		return nil
	}

	blobQueueSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Name:      "blob_queue_size",
			Help:      "The number of blobs waiting to be put into a batch for dispersal.",
		},
		[]string{},
	)

	dedupSetSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Name:      "dedup_set_size",
			Help:      "The size of the deduplication set.",
		},
		[]string{},
	)

	uniqueBlobCounter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "unique_blobs",
			Help:      "The total number of unique blobs added to the queue.",
		},
		[]string{},
	)

	duplicateBlobCounter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "duplicate_blobs",
			Help:      "The total number of duplicate blobs seen when polling the blob source.",
		},
		[]string{},
	)

	timeSpentInQueue := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "time_in_queue_ms",
			Help:      "The time blobs spent in the queue before being acquired by the controller",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{},
	)

	pollLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "blob_source_poll_latency_ms",
			Help:      "The time taken to poll the blob source for new blobs.",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{},
	)

	return &blobQueueMetrics{
		blobQueueSize:        blobQueueSize,
		dedupSetSize:         dedupSetSize,
		uniqueBlobCounter:    uniqueBlobCounter,
		duplicateBlobCounter: duplicateBlobCounter,
		timeSpentInQueue:     timeSpentInQueue,
		pollLatency:          pollLatency,
	}
}

// Report the number of blobs waiting in the queue.
func (m *blobQueueMetrics) reportQueueSize(size uint64) {
	if m == nil {
		return
	}

	m.blobQueueSize.WithLabelValues().Set(float64(size))
}

// This should be called each time we poll the blob source. The first count should be the total blob count, not
// just unique blobs. The second count should be the number of duplicate blobs (i.e. blobs that have been seen
// before). The duration is how long the poll took.
func (m *blobQueueMetrics) reportBlobSourcePoll(blobCount uint64, duplicateBlobCount uint64, duration time.Duration) {
	if m == nil {
		return
	}

	uniqueBlobs := blobCount - duplicateBlobCount

	m.uniqueBlobCounter.WithLabelValues().Add(float64(uniqueBlobs))
	m.duplicateBlobCounter.WithLabelValues().Add(float64(duplicateBlobCount))
	m.pollLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// Report the size of the deduplication set.
func (m *blobQueueMetrics) reportDedupSetSize(size uint64) {
	if m == nil {
		return
	}

	m.dedupSetSize.WithLabelValues().Set(float64(size))
}

// Report the time a blob spent in the queue before being acquired by the controller for batching.
func (m *blobQueueMetrics) reportTimeInQueue(duration time.Duration) {
	if m == nil {
		return
	}

	m.timeSpentInQueue.WithLabelValues().Observe(common.ToMilliseconds(duration))
}
