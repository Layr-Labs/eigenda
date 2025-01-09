package controller

import (
	"time"

	common "github.com/Layr-Labs/eigenda/common"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const encodingManagerNamespace = "eigenda_encoding_manager"

// encodingManagerMetrics is a struct that holds the metrics for the encoding manager.
type encodingManagerMetrics struct {
	batchSubmissionLatency  *prometheus.SummaryVec
	blobHandleLatency       *prometheus.SummaryVec
	encodingLatency         *prometheus.SummaryVec
	putBlobCertLatency      *prometheus.SummaryVec
	updateBlobStatusLatency *prometheus.SummaryVec
	blobE2EEncodingLatency  *prometheus.SummaryVec
	batchSize               *prometheus.GaugeVec
	batchDataSize           *prometheus.GaugeVec
	batchRetryCount         *prometheus.GaugeVec
	failedSubmissionCount   *prometheus.CounterVec
	completedBlobs          *prometheus.CounterVec
}

// NewEncodingManagerMetrics sets up metrics for the encoding manager.
func newEncodingManagerMetrics(registry *prometheus.Registry) *encodingManagerMetrics {
	batchSubmissionLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "batch_submission_latency_ms",
			Help:       "The time required to submit a blob to the work pool for encoding.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	blobHandleLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "blob_handle_latency_ms",
			Help:       "The total time required to handle a blob.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	encodingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "encoding_latency_ms",
			Help:       "The time required to encode a blob.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	putBlobCertLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "put_blob_cert_latency_ms",
			Help:       "The time required to put a blob certificate.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	updateBlobStatusLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "update_blob_status_latency_ms",
			Help:       "The time required to update a blob status.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	blobE2EEncodingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  encodingManagerNamespace,
			Name:       "e2e_encoding_latency_ms",
			Help:       "The time required to encode a blob end-to-end.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	batchSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: encodingManagerNamespace,
			Name:      "batch_size",
			Help:      "The number of blobs in a batch.",
		},
		[]string{},
	)

	batchDataSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: encodingManagerNamespace,
			Name:      "batch_data_size_bytes",
			Help:      "The size of the data in a batch.",
		},
		[]string{},
	)

	batchRetryCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: encodingManagerNamespace,
			Name:      "batch_retry_count",
			Help:      "The number of retries required to encode a blob.",
		},
		[]string{},
	)

	failSubmissionCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: encodingManagerNamespace,
			Name:      "failed_submission_count",
			Help:      "The number of failed blob submissions (even after retries).",
		},
		[]string{},
	)

	completedBlobs := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: encodingManagerNamespace,
			Name:      "completed_blobs_total",
			Help:      "The number and size of completed blobs by status.",
		},
		[]string{"state", "data"},
	)

	return &encodingManagerMetrics{
		batchSubmissionLatency:  batchSubmissionLatency,
		blobHandleLatency:       blobHandleLatency,
		encodingLatency:         encodingLatency,
		putBlobCertLatency:      putBlobCertLatency,
		updateBlobStatusLatency: updateBlobStatusLatency,
		blobE2EEncodingLatency:  blobE2EEncodingLatency,
		batchSize:               batchSize,
		batchDataSize:           batchDataSize,
		batchRetryCount:         batchRetryCount,
		failedSubmissionCount:   failSubmissionCount,
		completedBlobs:          completedBlobs,
	}
}

func (m *encodingManagerMetrics) reportBatchSubmissionLatency(duration time.Duration) {
	m.batchSubmissionLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportBlobHandleLatency(duration time.Duration) {
	m.blobHandleLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportEncodingLatency(duration time.Duration) {
	m.encodingLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportPutBlobCertLatency(duration time.Duration) {
	m.putBlobCertLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportUpdateBlobStatusLatency(duration time.Duration) {
	m.updateBlobStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportE2EEncodingLatency(duration time.Duration) {
	m.blobE2EEncodingLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *encodingManagerMetrics) reportBatchSize(size int) {
	m.batchSize.WithLabelValues().Set(float64(size))
}

func (m *encodingManagerMetrics) reportBatchDataSize(size uint64) {
	m.batchDataSize.WithLabelValues().Set(float64(size))
}

func (m *encodingManagerMetrics) reportBatchRetryCount(count int) {
	m.batchRetryCount.WithLabelValues().Set(float64(count))
}

func (m *encodingManagerMetrics) reportFailedSubmission() {
	m.failedSubmissionCount.WithLabelValues().Inc()
}

func (m *encodingManagerMetrics) reportCompletedBlob(size int, status dispv2.BlobStatus) {
	switch status {
	case dispv2.Encoded:
		m.completedBlobs.WithLabelValues("encoded", "number").Inc()
		m.completedBlobs.WithLabelValues("encoded", "size").Add(float64(size))
	case dispv2.Failed:
		m.completedBlobs.WithLabelValues("failed", "number").Inc()
		m.completedBlobs.WithLabelValues("failed", "size").Add(float64(size))
	default:
		return
	}

	m.completedBlobs.WithLabelValues("total", "number").Inc()
	m.completedBlobs.WithLabelValues("total", "size").Add(float64(size))
}
