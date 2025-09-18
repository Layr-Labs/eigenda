package controller

import (
	"fmt"
	"sort"
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
	sendChunksRetryCount         *prometheus.GaugeVec
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
	attestation                  *prometheus.GaugeVec
	blobSetSize                  *prometheus.GaugeVec
	batchStageTimer              *common.StageTimer
	sendToValidatorStageTimer    *common.StageTimer
	importantSigningThresholds   []float64
	signatureThresholds          *prometheus.CounterVec
	signedBatchCount             *prometheus.CounterVec
	signedByteCount              *prometheus.CounterVec
	unsignedBatchCount           *prometheus.CounterVec
	unsignedByteCount            *prometheus.CounterVec
	timeoutBatchCount            *prometheus.CounterVec
	timeoutByteCount             *prometheus.CounterVec
	signingLatency               *prometheus.SummaryVec
}

// NewDispatcherMetrics sets up metrics for the dispatcher.
//
// importantSigningThresholds is a list of meaningful thresholds. Thresholds should be between 0.0 and 1.0.
// A count of batches meeting each specified threshold is reported as a metric.
func newDispatcherMetrics(
	registry *prometheus.Registry,
	importantSigningThresholds []float64,
) (*dispatcherMetrics, error) {
	if registry == nil {
		return nil, nil
	}

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	attestation := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: controllerNamespace,
			Name:      "attestation",
			Help:      "number of signers and non-signers for the batch",
		},
		[]string{"type", "quorum"},
	)

	sendChunksRetryCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: controllerNamespace,
			Name:      "send_chunks_retry_count",
			Help:      "The number of times chunks were retried to be sent (part of HandleBatch()).",
		},
		[]string{},
	)

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

	// Verify that thresholds are sane
	for _, threshold := range importantSigningThresholds {
		if threshold < 0 || threshold > 1 {
			return nil, fmt.Errorf("threshold %f is not between 0.0 and 1.0", threshold)
		}
	}
	sort.Float64s(importantSigningThresholds)

	// Add thresholds for 0.0 and 1.0, if missing.
	if len(importantSigningThresholds) == 0 || importantSigningThresholds[0] != 0.0 {
		importantSigningThresholds = append([]float64{0.0}, importantSigningThresholds...)
	}
	if importantSigningThresholds[len(importantSigningThresholds)-1] != 1.0 {
		importantSigningThresholds = append(importantSigningThresholds, 1.0)
	}

	batchSigningThresholdCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "batch_signing_threshold_count",
			Help:      "A count of batches that have reached various signature thresholds.",
		},
		[]string{"quorum", "threshold"},
	)

	signingRateLabels := []string{"id", "quorum"}

	signedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_signed_batch_count",
			Help:      "Total number of batches successfully signed by validators",
		},
		signingRateLabels,
	)

	signedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_signed_byte_count",
			Help:      "Total number of bytes successfully signed by validators",
		},
		signingRateLabels,
	)

	unsignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_unsigned_batch_count",
			Help:      "Total number of batches that validators failed to sign",
		},
		signingRateLabels,
	)

	unsignedByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_unsigned_byte_count",
			Help:      "Total number of bytes that validators failed to sign",
		},
		signingRateLabels,
	)

	timeoutBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_timeout_batch_count",
			Help:      "Total number of batches that validators failed to sign due to timeout",
		},
		signingRateLabels,
	)

	timeoutByteCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_timeout_byte_count",
			Help:      "Total number of bytes that validators failed to sign due to timeout",
		},
		signingRateLabels,
	)

	signingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "validator_signing_latency_seconds",
			Help:       "Latency for validators to sign batches",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		signingRateLabels,
	)

	return &dispatcherMetrics{
		sendChunksRetryCount:         sendChunksRetryCount,
		processSigningMessageLatency: processSigningMessageLatency,
		signingMessageChannelLatency: signingMessageChannelLatency,
		attestationUpdateLatency:     attestationUpdateLatency,
		attestationBuildingLatency:   attestationBuildingLatency,
		thresholdSignedToDoneLatency: thresholdSignedToDoneLatency,
		aggregateSignaturesLatency:   aggregateSignaturesLatency,
		putAttestationLatency:        putAttestationLatency,
		attestationUpdateCount:       attestationUpdateCount,
		updateBatchStatusLatency:     updateBatchStatusLatency,
		blobE2EDispersalLatency:      blobE2EDispersalLatency,
		completedBlobs:               completedBlobs,
		attestation:                  attestation,
		blobSetSize:                  blobSetSize,
		batchStageTimer:              batchStageTimer,
		sendToValidatorStageTimer:    sendToValidatorStageTimer,
		importantSigningThresholds:   importantSigningThresholds,
		signatureThresholds:          batchSigningThresholdCount,
		signedBatchCount:             signedBatchCount,
		signedByteCount:              signedByteCount,
		unsignedBatchCount:           unsignedBatchCount,
		unsignedByteCount:            unsignedByteCount,
		timeoutBatchCount:            timeoutBatchCount,
		timeoutByteCount:             timeoutByteCount,
		signingLatency:               signingLatency,
	}, nil
}

func (m *dispatcherMetrics) reportSendChunksRetryCount(retries float64) {
	if m == nil {
		return
	}
	m.sendChunksRetryCount.WithLabelValues().Set(retries)
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

func (m *dispatcherMetrics) reportAttestation(
	operatorCount map[core.QuorumID]int,
	signerCount map[core.QuorumID]int,
	quorumResults map[core.QuorumID]*core.QuorumResult,
) {

	if m == nil {
		return
	}

	for quorumID, count := range operatorCount {
		quorumStr := fmt.Sprintf("%d", quorumID)
		signers, ok := signerCount[quorumID]
		if !ok {
			continue
		}
		nonSigners := count - signers
		quorumResult, ok := quorumResults[quorumID]
		if !ok {
			continue
		}

		m.attestation.WithLabelValues("signers", quorumStr).Set(float64(signers))
		m.attestation.WithLabelValues("non_signers", quorumStr).Set(float64(nonSigners))
		m.attestation.WithLabelValues("percent_signed", quorumStr).Set(float64(quorumResult.PercentSigned))

		m.reportSigningThreshold(quorumID, float64(quorumResult.PercentSigned)/100.0)
	}
}

func (m *dispatcherMetrics) reportSigningThreshold(quorumID core.QuorumID, signingFraction float64) {
	if m == nil {
		return
	}

	// First, determine the threshold to report. In order to be reported as threshold X, the signing fraction
	// must be greater than or equal to X, but strictly less than the next highest threshold.
	//
	// For example, let's say important thresholds are [0, 0.55, 0.67, 0.80, 1.0]
	// 0.55 signing -> threshold 0.55 (>= 0.55 but < 0.67)
	// 0.56 signing -> threshold 0.55 (>= 0.55 but < 0.67)
	// 0.66 signing -> threshold 0.55 (>= 0.55 but < 0.67)
	// 0.67 signing -> threshold 0.67 (>= 0.67 but < 0.80)

	var threshold float64
	for i := len(m.importantSigningThresholds) - 1; i >= 0; i-- {
		candidateThreshold := m.importantSigningThresholds[i]
		if candidateThreshold <= signingFraction {
			threshold = candidateThreshold
			break
		}
	}

	quorumString := fmt.Sprintf("%d", quorumID)
	thresholdString := fmt.Sprintf("%f", threshold)

	m.signatureThresholds.WithLabelValues(quorumString, thresholdString).Inc()
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
func (m *dispatcherMetrics) ReportSigningSuccess(
	id core.OperatorID,
	batchSize uint64,
	signingLatency time.Duration,
	quorum core.QuorumID) {

	if m == nil {
		return
	}

	label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

	m.signedBatchCount.With(label).Add(1)
	m.signedByteCount.With(label).Add(float64(batchSize))
	m.signingLatency.With(label).Observe(signingLatency.Seconds())
}

// Report a failed signing event for a validator.
func (m *dispatcherMetrics) ReportSigningFailure(
	id core.OperatorID,
	batchSize uint64,
	timeout bool,
	quorum core.QuorumID) {

	if m == nil {
		return
	}

	label := prometheus.Labels{"id": id.Hex(), "quorum": fmt.Sprintf("%d", quorum)}

	m.unsignedBatchCount.With(label).Add(1)
	m.unsignedByteCount.With(label).Add(float64(batchSize))
	if timeout {
		m.timeoutBatchCount.With(label).Add(1)
		m.timeoutByteCount.With(label).Add(float64(batchSize))
	}
}
