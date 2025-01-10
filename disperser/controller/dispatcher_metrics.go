package controller

import (
	"time"

	common "github.com/Layr-Labs/eigenda/common"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const dispatcherNamespace = "eigenda_dispatcher"

// dispatcherMetrics is a struct that holds the metrics for the dispatcher.
type dispatcherMetrics struct {
	handleBatchLatency          *prometheus.SummaryVec
	newBatchLatency             *prometheus.SummaryVec
	getBlobMetadataLatency      *prometheus.SummaryVec
	getOperatorStateLatency     *prometheus.SummaryVec
	getBlobCertificatesLatency  *prometheus.SummaryVec
	buildMerkleTreeLatency      *prometheus.SummaryVec
	putBatchHeaderLatency       *prometheus.SummaryVec
	proofLatency                *prometheus.SummaryVec
	putVerificationInfosLatency *prometheus.SummaryVec
	poolSubmissionLatency       *prometheus.SummaryVec
	putDispersalRequestLatency  *prometheus.SummaryVec
	sendChunksLatency           *prometheus.SummaryVec
	sendChunksRetryCount        *prometheus.GaugeVec
	putDispersalResponseLatency *prometheus.SummaryVec
	handleSignaturesLatency     *prometheus.SummaryVec
	receiveSignaturesLatency    *prometheus.SummaryVec
	aggregateSignaturesLatency  *prometheus.SummaryVec
	putAttestationLatency       *prometheus.SummaryVec
	updateBatchStatusLatency    *prometheus.SummaryVec
	blobE2EDispersalLatency     *prometheus.SummaryVec
	completedBlobs              *prometheus.CounterVec
}

// NewDispatcherMetrics sets up metrics for the dispatcher.
func newDispatcherMetrics(registry *prometheus.Registry) *dispatcherMetrics {
	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	handleBatchLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "handle_batch_latency_ms",
			Help:       "The time required to handle a batch.",
			Objectives: objectives,
		},
		[]string{},
	)

	newBatchLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "new_batch_latency_ms",
			Help:       "The time required to create a new batch (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobMetadataLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "get_blob_metadata_latency_ms",
			Help:       "The time required to get blob metadata (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	getOperatorStateLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "get_operator_state_latency_ms",
			Help:       "The time required to get the operator state (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobCertificatesLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "get_blob_certificates_latency_ms",
			Help:       "The time required to get blob certificates (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	buildMerkleTreeLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "build_merkle_tree_latency_ms",
			Help:       "The time required to build the Merkle tree (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	putBatchHeaderLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "put_batch_header_latency_ms",
			Help:       "The time required to put the batch header (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	proofLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "proof_latency_ms",
			Help:       "The time required to generate the proof (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	putVerificationInfosLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "put_verification_infos_latency_ms",
			Help:       "The time required to put the verification infos (part of NewBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	poolSubmissionLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "pool_submission_latency_ms",
			Help:       "The time required to submit a batch to the worker pool (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	putDispersalRequestLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "put_dispersal_latency_ms",
			Help:       "The time required to put the dispersal request (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	sendChunksLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "send_chunks_latency_ms",
			Help:       "The time required to send chunks (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	sendChunksRetryCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: dispatcherNamespace,
			Name:      "send_chunks_retry_count",
			Help:      "The number of times chunks were retried to be sent (part of HandleBatch()).",
		},
		[]string{},
	)

	putDispersalResponseLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "put_dispersal_response_latency_ms",
			Help:       "The time required to put the dispersal response (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	handleSignaturesLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "handle_signatures_latency_ms",
			Help:       "The time required to handle signatures (part of HandleBatch()).",
			Objectives: objectives,
		},
		[]string{},
	)

	receiveSignaturesLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  dispatcherNamespace,
			Name:       "receive_signatures_latency_ms",
			Help:       "The time required to receive signatures (part of HandleSignatures()).",
			Objectives: objectives,
		},
		[]string{},
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

	return &dispatcherMetrics{
		handleBatchLatency:          handleBatchLatency,
		newBatchLatency:             newBatchLatency,
		getBlobMetadataLatency:      getBlobMetadataLatency,
		getOperatorStateLatency:     getOperatorStateLatency,
		getBlobCertificatesLatency:  getBlobCertificatesLatency,
		buildMerkleTreeLatency:      buildMerkleTreeLatency,
		putBatchHeaderLatency:       putBatchHeaderLatency,
		proofLatency:                proofLatency,
		putVerificationInfosLatency: putVerificationInfosLatency,
		poolSubmissionLatency:       poolSubmissionLatency,
		putDispersalRequestLatency:  putDispersalRequestLatency,
		sendChunksLatency:           sendChunksLatency,
		sendChunksRetryCount:        sendChunksRetryCount,
		putDispersalResponseLatency: putDispersalResponseLatency,
		handleSignaturesLatency:     handleSignaturesLatency,
		receiveSignaturesLatency:    receiveSignaturesLatency,
		aggregateSignaturesLatency:  aggregateSignaturesLatency,
		putAttestationLatency:       putAttestationLatency,
		updateBatchStatusLatency:    updateBatchStatusLatency,
		blobE2EDispersalLatency:     blobE2EDispersalLatency,
		completedBlobs:              completedBlobs,
	}
}

func (m *dispatcherMetrics) reportHandleBatchLatency(duration time.Duration) {
	m.handleBatchLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportNewBatchLatency(duration time.Duration) {
	m.newBatchLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportGetBlobMetadataLatency(duration time.Duration) {
	m.getBlobMetadataLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportGetOperatorStateLatency(duration time.Duration) {
	m.getOperatorStateLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportGetBlobCertificatesLatency(duration time.Duration) {
	m.getBlobCertificatesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportBuildMerkleTreeLatency(duration time.Duration) {
	m.buildMerkleTreeLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutBatchHeaderLatency(duration time.Duration) {
	m.putBatchHeaderLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportProofLatency(duration time.Duration) {
	m.proofLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutVerificationInfosLatency(duration time.Duration) {
	m.putVerificationInfosLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPoolSubmissionLatency(duration time.Duration) {
	m.poolSubmissionLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutDispersalRequestLatency(duration time.Duration) {
	m.putDispersalRequestLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportSendChunksLatency(duration time.Duration) {
	m.sendChunksLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportSendChunksRetryCount(retries float64) {
	m.sendChunksRetryCount.WithLabelValues().Set(retries)
}

func (m *dispatcherMetrics) reportPutDispersalResponseLatency(duration time.Duration) {
	m.putDispersalResponseLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportHandleSignaturesLatency(duration time.Duration) {
	m.handleSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportReceiveSignaturesLatency(duration time.Duration) {
	m.receiveSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportAggregateSignaturesLatency(duration time.Duration) {
	m.aggregateSignaturesLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportPutAttestationLatency(duration time.Duration) {
	m.putAttestationLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportUpdateBatchStatusLatency(duration time.Duration) {
	m.updateBatchStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportE2EDispersalLatency(duration time.Duration) {
	m.blobE2EDispersalLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *dispatcherMetrics) reportCompletedBlob(size int, status dispv2.BlobStatus) {
	switch status {
	case dispv2.Certified:
		m.completedBlobs.WithLabelValues("certified", "number").Inc()
		m.completedBlobs.WithLabelValues("certified", "size").Add(float64(size))
	case dispv2.Failed:
		m.completedBlobs.WithLabelValues("failed", "number").Inc()
		m.completedBlobs.WithLabelValues("failed", "size").Add(float64(size))
	case dispv2.InsufficientSignatures:
		m.completedBlobs.WithLabelValues("insufficient_signature", "number").Inc()
		m.completedBlobs.WithLabelValues("insufficient_signature", "size").Add(float64(size))
	default:
		return
	}

	m.completedBlobs.WithLabelValues("total", "number").Inc()
	m.completedBlobs.WithLabelValues("total", "size").Add(float64(size))
}
