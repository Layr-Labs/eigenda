package signingrate

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// It's unfortunate that this isn't "controller",
// but at this point in time the effort to change the dashboards is extremely high.
//
// This is a duplicate definition to avoid an import cycle.
const ControllerMetricsNamespace = "eigenda_dispatcher"

// Encapsulates metrics about the signing rate of validators.
type SigningRateMetrics struct {
	signedBatchCount   *prometheus.CounterVec
	signedByteCount    *prometheus.CounterVec
	unsignedBatchCount *prometheus.CounterVec
	unsignedByteCount  *prometheus.CounterVec
	timeoutBatchCount  *prometheus.CounterVec
	timeoutByteCount   *prometheus.CounterVec
	signingLatency     *prometheus.SummaryVec
}

// NewSigningRateMetrics creates a new SigningRateMetrics instance.
func NewSigningRateMetrics(registry *prometheus.Registry) *SigningRateMetrics {
	if registry == nil {
		return nil
	}

	labels := []string{"id", "quorum"}

	signedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_signed_batch_count",
			Help:      "Total number of batches successfully signed by validators",
		},
		labels,
	)

	signedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_signed_byte_count",
			Help:      "Total number of bytes successfully signed by validators",
		},
		labels,
	)

	unsignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_unsigned_batch_count",
			Help:      "Total number of batches that validators failed to sign",
		},
		labels,
	)

	unsignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_unsigned_byte_count",
			Help:      "Total number of bytes that validators failed to sign",
		},
		labels,
	)

	timeoutBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_timeout_batch_count",
			Help:      "Total number of batches that validators failed to sign due to timeout",
		},
		labels,
	)

	timeoutByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ControllerMetricsNamespace,
			Name:      "validator_timeout_byte_count",
			Help:      "Total number of bytes that validators failed to sign due to timeout",
		},
		labels,
	)

	signingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  ControllerMetricsNamespace,
			Name:       "validator_signing_latency_seconds",
			Help:       "Latency for validators to sign batches",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		labels,
	)

	return &SigningRateMetrics{
		signedBatchCount:   signedBatchCount,
		signedByteCount:    signedByteCount,
		unsignedBatchCount: unsignedBatchCount,
		unsignedByteCount:  unsignedByteCount,
		timeoutBatchCount:  timeoutBatchCount,
		timeoutByteCount:   timeoutByteCount,
		signingLatency:     signingLatency,
	}
}

// Report a successful signing event for a validator.
func (s *SigningRateMetrics) ReportSuccess(
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
	quorums []core.QuorumID) {

	if s == nil {
		return
	}

	for _, quorum := range quorums {
		label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

		s.signedBatchCount.With(label).Add(1)
		s.signedByteCount.With(label).Add(float64(batchSize))
		s.signingLatency.With(label).Observe(signingLatency.Seconds())
	}

}

// Report a failed signing event for a validator.
func (s *SigningRateMetrics) ReportFailure(
	id core.OperatorID,
	batchSize uint64,
	timeout bool,
	quorums []core.QuorumID) {

	if s == nil {
		return
	}

	for _, quorum := range quorums {
		label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

		s.unsignedBatchCount.With(label).Add(1)
		s.unsignedByteCount.With(label).Add(float64(batchSize))
		if timeout {
			s.timeoutBatchCount.With(label).Add(1)
			s.timeoutByteCount.With(label).Add(float64(batchSize))
		}
	}
}
