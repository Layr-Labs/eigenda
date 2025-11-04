package controller

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const controllerNamespace = "eigenda_dispatcher"

// dispatcherMetrics is a struct that holds the metrics for the dispatcher.
type dispatcherMetrics struct {
	processSigningMessageLatency *prometheus.SummaryVec
	signingMessageChannelLatency *prometheus.SummaryVec
	attestationUpdateLatency     *prometheus.SummaryVec
	attestationBuildingLatency   *prometheus.SummaryVec
	thresholdSignedToDoneLatency *prometheus.SummaryVec
	aggregateSignaturesLatency   *prometheus.SummaryVec
	putAttestationLatency        *prometheus.SummaryVec
	attestationUpdateCount       *prometheus.SummaryVec
	updateBatchStatusLatency     *prometheus.SummaryVec
	blobE2EDispersalLatency      *prometheus.SummaryVec
	completedBlobs               *prometheus.CounterVec
	blobSetSize                  *prometheus.GaugeVec
	batchStageTimer              *common.StageTimer
	sendToValidatorStageTimer    *common.StageTimer

	minimumSigningThreshold float64

	validatorSignedBatchCount   *prometheus.CounterVec
	validatorSignedByteCount    *prometheus.CounterVec
	validatorUnsignedBatchCount *prometheus.CounterVec
	validatorUnsignedByteCount  *prometheus.CounterVec
	validatorTimeoutBatchCount  *prometheus.CounterVec
	validatorTimeoutByteCount   *prometheus.CounterVec
	validatorSigningLatency     *prometheus.SummaryVec

	globalSignedBatchCount   *prometheus.CounterVec
	globalUnsignedBatchCount *prometheus.CounterVec
	globalSignedByteCount    *prometheus.CounterVec
	globalUnsignedByteCount  *prometheus.CounterVec

	globalSigningFractionHistogram *prometheus.HistogramVec

	collectDetailedValidatorMetrics bool
}

// NewDispatcherMetrics sets up metrics for the dispatcher.
//
// importantSigningThresholds is a list of meaningful thresholds. Thresholds should be between 0.0 and 1.0.
// A count of batches meeting each specified threshold is reported as a metric.
func newDispatcherMetrics(
	registry *prometheus.Registry,
	// The minimum fraction of signers for a batch to be considered properly signed. Any fraction greater
	// than or equal to this value is considered a successful signing.
	minimumSigningThreshold float64,
	// If true, collect detailed per-validator metrics. This can be disabled if the volume of data
	// produced is too high.
	collectDetailedValidatorMetrics bool,
) (*dispatcherMetrics, error) {
	if registry == nil {
		return nil, nil
	}

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	processSigningMessageLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "process_signing_message_latency_ms",
			Help:       "The time required to process a single signing message (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	signingMessageChannelLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "signing_message_channel_latency_ms",
			Help:       "The time a signing message sits in the channel waiting to be processed (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationUpdateLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "attestation_update_latency_ms",
			Help:       "The time between the signature receiver yielding attestations (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationBuildingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "attestation_building_latency_ms",
			Help:       "The time it takes for the signature receiver to build and send a single attestation (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationUpdateCount := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "attestation_update_count",
			Help:       "The number of updates to the batch attestation throughout the signature gathering process.",
			Objectives: objectives,
		},
		[]string{},
	)

	thresholdSignedToDoneLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: controllerNamespace,
			Name:      "threshold_signed_to_done_latency_ms",
			Help: "the time elapsed between the signing percentage reaching a configured threshold, and the end " +
				"of signature gathering",
			Objectives: objectives,
		},
		[]string{"quorum"},
	)

	aggregateSignaturesLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "aggregate_signatures_latency_ms",
			Help:       "The time required to aggregate signatures (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	putAttestationLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "put_attestation_latency_ms",
			Help:       "The time required to put the attestation (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	updateBatchStatusLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "update_batch_status_latency_ms",
			Help:       "The time required to update the batch status (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	blobE2EDispersalLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "e2e_dispersal_latency_ms",
			Help:       "The time required to disperse a blob end-to-end.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	completedBlobs := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "completed_blobs_total",
			Help:      "The number and size of completed blobs by status.",
		},
		[]string{"state", "data"},
	)

	blobSetSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: controllerNamespace,
			Name:      "blob_queue_size",
			Help:      "The size of the blob queue used for deduplication.",
		},
		[]string{},
	)

	batchStageTimer := common.NewStageTimer(registry, controllerNamespace, "batch", false)
	sendToValidatorStageTimer := common.NewStageTimer(
		registry,
		controllerNamespace,
		"send_to_validator",
		false)

	signingRateLabels := []string{"id", "quorum"}

	validatorSignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_signed_batch_count",
			Help:      "Total number of batches successfully signed by validators",
		},
		signingRateLabels,
	)

	validatorSignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_signed_byte_count",
			Help:      "Total number of bytes successfully signed by validators",
		},
		signingRateLabels,
	)

	validatorUnsignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_unsigned_batch_count",
			Help:      "Total number of batches that validators failed to sign",
		},
		signingRateLabels,
	)

	validatorUnsignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_unsigned_byte_count",
			Help:      "Total number of bytes that validators failed to sign",
		},
		signingRateLabels,
	)

	validatorTimeoutBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_timeout_batch_count",
			Help:      "Total number of batches that validators failed to sign due to timeout",
		},
		signingRateLabels,
	)

	validatorTimeoutByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_timeout_byte_count",
			Help:      "Total number of bytes that validators failed to sign due to timeout",
		},
		signingRateLabels,
	)

	validatorSigningLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "validator_signing_latency_seconds",
			Help:       "Latency for validators to sign batches",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		signingRateLabels,
	)

	globalSignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "global_signed_batch_count",
			Help:      "Total number of batches successfully signed by a critical mass of validators",
		},
		[]string{"quorum"},
	)

	globalUnsignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "global_unsigned_batch_count",
			Help:      "Total number of batches that were not signed by a critical mass of validators",
		},
		[]string{"quorum"},
	)

	globalSignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "global_signed_byte_count",
			Help:      "Total number of bytes successfully signed by a critical mass of validators",
		},
		[]string{"quorum"},
	)

	globalUnsignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "global_unsigned_byte_count",
			Help:      "Total number of bytes that were not signed by a critical mass of validators",
		},
		[]string{"quorum"},
	)

	globalSigningFractionHistogram := promauto.With(registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: controllerNamespace,
			Name:      "global_signing_fraction_histogram",
			Help:      "Histogram of the fraction of validators that signed each batch",
			Buckets:   prometheus.LinearBuckets(0.0, 0.05, 21),
		},
		[]string{"quorum"},
	)

	return &dispatcherMetrics{
		processSigningMessageLatency:    processSigningMessageLatency,
		signingMessageChannelLatency:    signingMessageChannelLatency,
		attestationUpdateLatency:        attestationUpdateLatency,
		attestationBuildingLatency:      attestationBuildingLatency,
		thresholdSignedToDoneLatency:    thresholdSignedToDoneLatency,
		aggregateSignaturesLatency:      aggregateSignaturesLatency,
		putAttestationLatency:           putAttestationLatency,
		attestationUpdateCount:          attestationUpdateCount,
		updateBatchStatusLatency:        updateBatchStatusLatency,
		blobE2EDispersalLatency:         blobE2EDispersalLatency,
		completedBlobs:                  completedBlobs,
		blobSetSize:                     blobSetSize,
		batchStageTimer:                 batchStageTimer,
		sendToValidatorStageTimer:       sendToValidatorStageTimer,
		minimumSigningThreshold:         minimumSigningThreshold,
		validatorSignedBatchCount:       validatorSignedBatchCount,
		validatorSignedByteCount:        validatorSignedByteCount,
		validatorUnsignedBatchCount:     validatorUnsignedBatchCount,
		validatorUnsignedByteCount:      validatorUnsignedByteCount,
		validatorTimeoutBatchCount:      validatorTimeoutBatchCount,
		validatorTimeoutByteCount:       validatorTimeoutByteCount,
		validatorSigningLatency:         validatorSigningLatency,
		collectDetailedValidatorMetrics: collectDetailedValidatorMetrics,
		globalSignedBatchCount:          globalSignedBatchCount,
		globalUnsignedBatchCount:        globalUnsignedBatchCount,
		globalSignedByteCount:           globalSignedByteCount,
		globalUnsignedByteCount:         globalUnsignedByteCount,
		globalSigningFractionHistogram:  globalSigningFractionHistogram,
	}, nil
}

func (m *dispatcherMetrics) reportProcessSigningMessageLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.processSigningMessageLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportSigningMessageChannelLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.signingMessageChannelLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationUpdateLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.attestationUpdateLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationBuildingLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.attestationBuildingLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportThresholdSignedToDoneLatency(quorumID core.QuorumID, duration time.Duration) {
	if m == nil {
		return
	}
	m.thresholdSignedToDoneLatency.WithLabelValues(fmt.Sprintf("%d", quorumID)).Observe(
		common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAggregateSignaturesLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.aggregateSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutAttestationLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.putAttestationLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationUpdateCount(attestationCount float64) {
	if m == nil {
		return
	}
	m.attestationUpdateCount.WithLabelValues().Observe(attestationCount)
}

func (m *dispatcherMetrics) reportUpdateBatchStatusLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.updateBatchStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportE2EDispersalLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.blobE2EDispersalLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportCompletedBlob(size int, status dispv2.BlobStatus) {
	if m == nil {
		return
	}
	switch status {
	case dispv2.Complete:
		m.completedBlobs.WithLabelValues("complete", "number").Inc()
		m.completedBlobs.WithLabelValues("complete", "size").Add(float64(size))
	case dispv2.Failed:
		m.completedBlobs.WithLabelValues("failed", "number").Inc()
		m.completedBlobs.WithLabelValues("failed", "size").Add(float64(size))
	default:
		return
	}

	m.completedBlobs.WithLabelValues("total", "number").Inc()
	m.completedBlobs.WithLabelValues("total", "size").Add(float64(size))
}

func (m *dispatcherMetrics) reportBlobSetSize(size int) {
	if m == nil {
		return
	}
	m.blobSetSize.WithLabelValues().Set(float64(size))
}

func (m *dispatcherMetrics) reportSigningThreshold(
	quorumID core.QuorumID,
	batchSizeBytes uint64,
	signingFraction float64,
) {
	if m == nil {
		return
	}

	quorumString := fmt.Sprintf("%d", quorumID)
	labels := prometheus.Labels{"quorum": quorumString}

	if signingFraction >= m.minimumSigningThreshold {
		m.globalSignedBatchCount.With(labels).Inc()
		m.globalSignedByteCount.With(labels).Add(float64(batchSizeBytes))
	} else {
		m.globalUnsignedBatchCount.With(labels).Inc()
		m.globalUnsignedByteCount.With(labels).Add(float64(batchSizeBytes))
	}

	m.globalSigningFractionHistogram.With(labels).Observe(signingFraction)
}

func (m *dispatcherMetrics) newBatchProbe() *common.SequenceProbe {
	if m == nil {
		// A sequence probe becomes a no-op when nil.
		return nil
	}

	return m.batchStageTimer.NewSequence()
}

func (m *dispatcherMetrics) newSendToValidatorProbe() *common.SequenceProbe {
	if m == nil {
		// A sequence probe becomes a no-op when nil.
		return nil
	}

	return m.sendToValidatorStageTimer.NewSequence()
}

// Report a successful signing event for a validator.
func (m *dispatcherMetrics) ReportValidatorSigningSuccess(
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
	quorums []core.QuorumID,
) {

	if m == nil || !m.collectDetailedValidatorMetrics {
		return
	}

	for _, quorum := range quorums {
		label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

		m.validatorSignedBatchCount.With(label).Add(1)
		m.validatorSignedByteCount.With(label).Add(float64(batchSize))
		m.validatorSigningLatency.With(label).Observe(signingLatency.Seconds())
	}

}

// Report a failed signing event for a validator.
func (m *dispatcherMetrics) ReportValidatorSigningFailure(
	id core.OperatorID,
	batchSize uint64,
	timeout bool,
	quorums []core.QuorumID,
) {

	if m == nil || !m.collectDetailedValidatorMetrics {
		return
	}

	for _, quorum := range quorums {
		label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

		m.validatorUnsignedBatchCount.With(label).Add(1)
		m.validatorUnsignedByteCount.With(label).Add(float64(batchSize))
		if timeout {
			m.validatorTimeoutBatchCount.With(label).Add(1)
			m.validatorTimeoutByteCount.With(label).Add(float64(batchSize))
		}
	}
}
