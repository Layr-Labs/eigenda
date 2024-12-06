package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

const encodingManagerNamespace = "eigenda_encoding_manager"

// encodingManagerMetrics is a struct that holds the metrics for the encoding manager.
type encodingManagerMetrics struct {
	batchSubmissionLatency  *prometheus.SummaryVec
	blobHandleLatency       *prometheus.SummaryVec
	encodingLatency         *prometheus.SummaryVec
	putBlobCertLatency      *prometheus.SummaryVec
	updateBlobStatusLatency *prometheus.SummaryVec
	batchSize               *prometheus.GaugeVec
	batchDataSize           *prometheus.GaugeVec
	batchRetryCount         *prometheus.GaugeVec
	batchSleepTime          *prometheus.GaugeVec
	failedSubmissionCount   *prometheus.CounterVec
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

	batchSleepTime := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: encodingManagerNamespace,
			Name:      "batch_sleep_time_ms",
			Help:      "The time slept during while waiting to retry encoding a blob.",
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

	return &encodingManagerMetrics{
		batchSubmissionLatency:  batchSubmissionLatency,
		blobHandleLatency:       blobHandleLatency,
		encodingLatency:         encodingLatency,
		putBlobCertLatency:      putBlobCertLatency,
		updateBlobStatusLatency: updateBlobStatusLatency,
		batchSize:               batchSize,
		batchDataSize:           batchDataSize,
		batchRetryCount:         batchRetryCount,
		batchSleepTime:          batchSleepTime,
		failedSubmissionCount:   failSubmissionCount,
	}
}

func (m *encodingManagerMetrics) reportBatchSubmissionLatency(duration time.Duration) {
	m.batchSubmissionLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *encodingManagerMetrics) reportBlobHandleLatency(duration time.Duration) {
	m.blobHandleLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *encodingManagerMetrics) reportEncodingLatency(duration time.Duration) {
	m.encodingLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *encodingManagerMetrics) reportPutBlobCertLatency(duration time.Duration) {
	m.putBlobCertLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *encodingManagerMetrics) reportUpdateBlobStatusLatency(duration time.Duration) {
	m.updateBlobStatusLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
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

func (m *encodingManagerMetrics) reportBatchSleepTime(duration time.Duration) {
	m.batchSleepTime.WithLabelValues().Set(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *encodingManagerMetrics) reportFailedSubmission() {
	m.failedSubmissionCount.WithLabelValues().Inc()
}
