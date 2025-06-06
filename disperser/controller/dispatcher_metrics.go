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

const dispatcherNamespace = "eigenda_dispatcher"

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
}

// NewDispatcherMetrics sets up metrics for the dispatcher.
//
// importantSigningThresholds is a list of meaningful thresholds. Thresholds should be between 0.0 and 1.0.
// A count of batches meeting each specified threshold is reported as a metric.
func newDispatcherMetrics(
	registry *prometheus.Registry,
	importantSigningThresholds []float64,
) (*dispatcherMetrics, error) {

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	attestation := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: dispatcherNamespace,
			Name:      "attestation",
			Help:      "number of signers and non-signers for the batch",
		},
		[]string{"type", "quorum"},
	)

	sendChunksRetryCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: dispatcherNamespace,
			Name:      "send_chunks_retry_count",
			Help:      "The number of times chunks were retried to be sent (part of HandleBatch()).",
		},
		[]string{},
	)

	processSigningMessageLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "process_signing_message_latency_ms",
			Help:       "The time required to process a single signing message (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	signingMessageChannelLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "signing_message_channel_latency_ms",
			Help:       "The time a signing message sits in the channel waiting to be processed (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationUpdateLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "attestation_update_latency_ms",
			Help:       "The time between the signature receiver yielding attestations (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationBuildingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "attestation_building_latency_ms",
			Help:       "The time it takes for the signature receiver to build and send a single attestation (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationUpdateCount := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "attestation_update_count",
			Help:       "The number of updates to the batch attestation throughout the signature gathering process.",
			Objectives: objectives,
		},
		[]string{},
	)

	thresholdSignedToDoneLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: dispatcherNamespace,
			Name:      "threshold_signed_to_done_latency_ms",
			Help: "the time elapsed between the signing percentage reaching a configured threshold, and the end " +
				"of signature gathering",
			Objectives: objectives,
		},
		[]string{"quorum"},
	)

	aggregateSignaturesLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "aggregate_signatures_latency_ms",
			Help:       "The time required to aggregate signatures (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	putAttestationLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "put_attestation_latency_ms",
			Help:       "The time required to put the attestation (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	updateBatchStatusLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "update_batch_status_latency_ms",
			Help:       "The time required to update the batch status (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	blobE2EDispersalLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "e2e_dispersal_latency_ms",
			Help:       "The time required to disperse a blob end-to-end.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{},
	)

	completedBlobs := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: dispatcherNamespace,
			Name:      "completed_blobs_total",
			Help:      "The number and size of completed blobs by status.",
		},
		[]string{"state", "data"},
	)

	blobSetSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: dispatcherNamespace,
			Name:      "blob_queue_size",
			Help:      "The size of the blob queue used for deduplication.",
		},
		[]string{},
	)

	batchStageTimer := common.NewStageTimer(registry, dispatcherNamespace, "batch", false)
	sendToValidatorStageTimer := common.NewStageTimer(
		registry,
		dispatcherNamespace,
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
			Namespace: dispatcherNamespace,
			Name:      "batch_signing_threshold_count",
			Help:      "A count of batches that have reached various signature thresholds.",
		},
		[]string{"quorum", "threshold"},
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
	}, nil
}

func (m *dispatcherMetrics) reportSendChunksRetryCount(retries float64) {
	m.sendChunksRetryCount.WithLabelValues().Set(retries)
}

func (m *dispatcherMetrics) reportProcessSigningMessageLatency(duration time.Duration) {
	m.processSigningMessageLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportSigningMessageChannelLatency(duration time.Duration) {
	m.signingMessageChannelLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationUpdateLatency(duration time.Duration) {
	m.attestationUpdateLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationBuildingLatency(duration time.Duration) {
	m.attestationBuildingLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportThresholdSignedToDoneLatency(quorumID core.QuorumID, duration time.Duration) {
	m.thresholdSignedToDoneLatency.WithLabelValues(fmt.Sprintf("%d", quorumID)).Observe(
		common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAggregateSignaturesLatency(duration time.Duration) {
	m.aggregateSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutAttestationLatency(duration time.Duration) {
	m.putAttestationLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAttestationUpdateCount(attestationCount float64) {
	m.attestationUpdateCount.WithLabelValues().Observe(attestationCount)
}

func (m *dispatcherMetrics) reportUpdateBatchStatusLatency(duration time.Duration) {
	m.updateBatchStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportE2EDispersalLatency(duration time.Duration) {
	m.blobE2EDispersalLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportCompletedBlob(size int, status dispv2.BlobStatus) {
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
	m.blobSetSize.WithLabelValues().Set(float64(size))
}

func (m *dispatcherMetrics) reportAttestation(
	operatorCount map[core.QuorumID]int,
	signerCount map[core.QuorumID]int,
	quorumResults map[core.QuorumID]*core.QuorumResult) {

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
	return m.batchStageTimer.NewSequence()
}

func (m *dispatcherMetrics) newSendToValidatorProbe() *common.SequenceProbe {
	return m.sendToValidatorStageTimer.NewSequence()
}
