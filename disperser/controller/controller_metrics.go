package controller

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/nameremapping"
	"github.com/Layr-Labs/eigenda/core"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// "dispatcher" is an unfortunate prefix, but since changing it will break many dashboards and alerts,
// we will keep it for now.
const controllerNamespace = "eigenda_dispatcher"

// controllerMetrics is a struct that holds the metrics for the controller.
type controllerMetrics struct {
	processSigningMessageLatency *prometheus.SummaryVec
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
	staleDispersalCount          prometheus.Counter
	batchStageTimer              *common.StageTimer
	sendToValidatorStageTimer    *common.StageTimer

	minimumSigningThreshold float64

	validatorSignedBatchCount   *prometheus.CounterVec
	validatorSignedByteCount    *prometheus.CounterVec
	validatorUnsignedBatchCount *prometheus.CounterVec
	validatorUnsignedByteCount  *prometheus.CounterVec
	validatorSigningLatency     *prometheus.SummaryVec

	globalSignedBatchCount   *prometheus.CounterVec
	globalUnsignedBatchCount *prometheus.CounterVec
	globalSignedByteCount    *prometheus.CounterVec
	globalUnsignedByteCount  *prometheus.CounterVec

	globalSigningFractionHistogram *prometheus.HistogramVec

	collectDetailedValidatorMetrics bool
	enablePerAccountMetrics         bool
	userAccountRemapping            map[string]string
	validatorIdRemapping            map[string]string
}

// Sets up metrics for the controller.
func newControllerMetrics(
	registry *prometheus.Registry,
	// The minimum fraction of signers for a batch to be considered properly signed. Any fraction greater
	// than or equal to this value is considered a successful signing.
	minimumSigningThreshold float64,
	// If true, collect detailed per-validator metrics. This can be disabled if the volume of data
	// produced is too high.
	collectDetailedValidatorMetrics bool,
	// If false, per-account blob completion metrics will be aggregated under "0x0" to reduce cardinality.
	enablePerAccountMetrics bool,
	// Maps account IDs to user-friendly names.
	userAccountRemapping map[string]string,
	// Maps validator IDs to validator names.
	validatorIdRemapping map[string]string,
) (*controllerMetrics, error) {
	if registry == nil {
		return nil, nil
	}

	if minimumSigningThreshold < 0.0 || minimumSigningThreshold > 1.0 {
		return nil, fmt.Errorf("invalid minimum signing threshold: %f", minimumSigningThreshold)
	}

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	// This metric is a loaded footgun, since it obscures quite a lot of information about what's happening
	// in the system. New metrics replace this, however we need to keep it around until alerts and dashboards
	// are configured to use the new metrics.
	attestation := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: controllerNamespace,
			Name:      "attestation",
			Help:      "number of signers and non-signers for the batch",
		},
		[]string{"type", "quorum"},
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

	attestationUpdateLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: controllerNamespace,
			Name:      "attestation_update_latency_ms",
			Help: "The time between the signature receiver yielding " +
				"attestations (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
	)

	attestationBuildingLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: controllerNamespace,
			Name:      "attestation_building_latency_ms",
			Help: "The time it takes for the signature receiver to build and " +
				"send a single attestation (part of HandleSignatures()).",
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
			Help:      "The number and size of completed blobs by status and account.",
		},
		[]string{"state", "data", "account_id"},
	)

	staleDispersalCount := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "stale_dispersal_count",
			Help:      "Total number of dispersals discarded due to being stale.",
		},
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
			Help: "Total number of bytes successfully signed by validators, " +
				"equal to size of signed batch times stake fraction",
		},
		signingRateLabels,
	)

	validatorUnsignedBatchCount := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: controllerNamespace,
			Name:      "validator_unsigned_batch_count",
			Help: "Total number of batches that validators failed to sign, " +
				"equal to size of unsigned batch times stake fraction",
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

	validatorSigningLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  controllerNamespace,
			Name:       "validator_signing_latency_ms",
			Help:       "The latency of signing messages for each validator.",
			Objectives: objectives,
		},
		[]string{"id"},
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

	return &controllerMetrics{
		processSigningMessageLatency:    processSigningMessageLatency,
		attestationUpdateLatency:        attestationUpdateLatency,
		attestationBuildingLatency:      attestationBuildingLatency,
		thresholdSignedToDoneLatency:    thresholdSignedToDoneLatency,
		aggregateSignaturesLatency:      aggregateSignaturesLatency,
		putAttestationLatency:           putAttestationLatency,
		attestationUpdateCount:          attestationUpdateCount,
		updateBatchStatusLatency:        updateBatchStatusLatency,
		blobE2EDispersalLatency:         blobE2EDispersalLatency,
		completedBlobs:                  completedBlobs,
		attestation:                     attestation,
		staleDispersalCount:             staleDispersalCount,
		batchStageTimer:                 batchStageTimer,
		sendToValidatorStageTimer:       sendToValidatorStageTimer,
		minimumSigningThreshold:         minimumSigningThreshold,
		validatorSignedBatchCount:       validatorSignedBatchCount,
		validatorSignedByteCount:        validatorSignedByteCount,
		validatorUnsignedBatchCount:     validatorUnsignedBatchCount,
		validatorUnsignedByteCount:      validatorUnsignedByteCount,
		validatorSigningLatency:         validatorSigningLatency,
		collectDetailedValidatorMetrics: collectDetailedValidatorMetrics,
		enablePerAccountMetrics:         enablePerAccountMetrics,
		userAccountRemapping:            userAccountRemapping,
		validatorIdRemapping:            validatorIdRemapping,
		globalSignedBatchCount:          globalSignedBatchCount,
		globalUnsignedBatchCount:        globalUnsignedBatchCount,
		globalSignedByteCount:           globalSignedByteCount,
		globalUnsignedByteCount:         globalUnsignedByteCount,
		globalSigningFractionHistogram:  globalSigningFractionHistogram,
	}, nil
}

func (m *controllerMetrics) reportProcessSigningMessageLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.processSigningMessageLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportAttestationUpdateLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.attestationUpdateLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportAttestationBuildingLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.attestationBuildingLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportThresholdSignedToDoneLatency(quorumID core.QuorumID, duration time.Duration) {
	if m == nil {
		return
	}
	m.thresholdSignedToDoneLatency.WithLabelValues(fmt.Sprintf("%d", quorumID)).Observe(
		common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportAggregateSignaturesLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.aggregateSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportPutAttestationLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.putAttestationLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportAttestationUpdateCount(attestationCount float64) {
	if m == nil {
		return
	}
	m.attestationUpdateCount.WithLabelValues().Observe(attestationCount)
}

func (m *controllerMetrics) reportUpdateBatchStatusLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.updateBatchStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportE2EDispersalLatency(duration time.Duration) {
	if m == nil {
		return
	}
	m.blobE2EDispersalLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *controllerMetrics) reportCompletedBlob(size int, status dispv2.BlobStatus, accountID string) {
	if m == nil {
		return
	}

	accountLabel := nameremapping.GetAccountLabel(accountID, m.userAccountRemapping, m.enablePerAccountMetrics)

	switch status {
	case dispv2.Complete:
		m.completedBlobs.WithLabelValues("complete", "number", accountLabel).Inc()
		m.completedBlobs.WithLabelValues("complete", "size", accountLabel).Add(float64(size))
	case dispv2.Failed:
		m.completedBlobs.WithLabelValues("failed", "number", accountLabel).Inc()
		m.completedBlobs.WithLabelValues("failed", "size", accountLabel).Add(float64(size))
	default:
		return
	}

	m.completedBlobs.WithLabelValues("total", "number", accountLabel).Inc()
	m.completedBlobs.WithLabelValues("total", "size", accountLabel).Add(float64(size))
}

func (m *controllerMetrics) reportStaleDispersal() {
	if m == nil {
		return
	}
	m.staleDispersalCount.Inc()
}

func (m *controllerMetrics) reportLegacyAttestation(
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
	}
}

func (m *controllerMetrics) ReportGlobalSigningThreshold(
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

func (m *controllerMetrics) newBatchProbe() *common.SequenceProbe {
	if m == nil {
		// A sequence probe becomes a no-op when nil.
		return nil
	}

	return m.batchStageTimer.NewSequence()
}

func (m *controllerMetrics) newSendToValidatorProbe() *common.SequenceProbe {
	if m == nil {
		// A sequence probe becomes a no-op when nil.
		return nil
	}

	return m.sendToValidatorStageTimer.NewSequence()
}

// Report the result of an attempted signing event for a validator.
func (m *controllerMetrics) ReportValidatorSigningResult(
	id core.OperatorID,
	stakeFraction float64,
	batchSize uint64,
	quorum core.QuorumID,
	success bool,
) {
	if m == nil || !m.collectDetailedValidatorMetrics {
		return
	}

	idLabel := nameremapping.GetAccountLabel(
		"0x"+id.Hex(),
		m.validatorIdRemapping,
		m.collectDetailedValidatorMetrics)
	label := prometheus.Labels{"id": idLabel, "quorum": fmt.Sprintf("%d", quorum)}

	if success {
		m.validatorSignedBatchCount.With(label).Add(1)
		m.validatorSignedByteCount.With(label).Add(float64(batchSize) * stakeFraction)
	} else {
		m.validatorUnsignedBatchCount.With(label).Add(1)
		m.validatorUnsignedByteCount.With(label).Add(float64(batchSize) * stakeFraction)
	}
}

// Report the signing latency for a validator. Should only be used for validators that successfully signed a batch.
func (m *controllerMetrics) ReportValidatorSigningLatency(id core.OperatorID, latency time.Duration) {
	if m == nil || !m.collectDetailedValidatorMetrics {
		return
	}

	idLabel := nameremapping.GetAccountLabel(
		"0x"+id.Hex(),
		m.validatorIdRemapping,
		m.collectDetailedValidatorMetrics)
	m.validatorSigningLatency.WithLabelValues(idLabel).Observe(common.ToMilliseconds(latency))
}
