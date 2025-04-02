package payloaddispersal

import (
	"time"

	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PayloadDisperserMetrics encapsulates the metrics for the PayloadDisperser.
type PayloadDisperserMetrics struct {
	statusCount   *prometheus.GaugeVec
	statusLatency *prometheus.SummaryVec
}

// NewPayloadDisperserMetrics creates a new PayloadDisperserMetrics object.
func NewPayloadDisperserMetrics(namespace string, registry *prometheus.Registry) *PayloadDisperserMetrics {
	if registry == nil {
		return nil
	}

	statusLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "status_latency_ms",
			Help:       "The time a blob spends in a given status.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"status"},
	)

	statusCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "blob_status_count",
			Help:      "The number of blobs with a given status.",
		},
		[]string{"status"},
	)

	return &PayloadDisperserMetrics{
		statusLatency: statusLatency,
		statusCount:   statusCount,
	}
}

// reportBlobStatusChange reports a change in the status of a blob.
func (m *PayloadDisperserMetrics) reportBlobStatusChange(
	oldStatus *v2.BlobStatus,
	newStatus *v2.BlobStatus,
	timeInOldStatus time.Duration) {

	if m == nil {
		return
	}

	if newStatus != nil {
		terminal := *newStatus == v2.BlobStatus_COMPLETE ||
			*newStatus == v2.BlobStatus_FAILED ||
			*newStatus == v2.BlobStatus_UNKNOWN

		if !terminal {
			// Only observe non-terminal statuses. The purpose of this metric is to show the status of blobs in flight,
			// not to track the total number of blobs that have succeeded or failed.
			m.statusCount.WithLabelValues(newStatus.String()).Inc()
		}
	}

	if oldStatus != nil {
		m.statusCount.WithLabelValues(oldStatus.String()).Dec()
		m.statusLatency.WithLabelValues(oldStatus.String()).Observe(common.ToMilliseconds(timeInOldStatus))
	}
}
