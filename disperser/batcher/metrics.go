package batcher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FailReason string

const (
	FailBatchHeaderHash        FailReason = "batch_header_hash"
	FailAggregateSignatures    FailReason = "aggregate_signatures"
	FailNoSignatures           FailReason = "no_signatures"
	FailConfirmBatch           FailReason = "confirm_batch"
	FailGetBatchID             FailReason = "get_batch_id"
	FailUpdateConfirmationInfo FailReason = "update_confirmation_info"
	FailNoAggregatedSignature  FailReason = "no_aggregated_signature"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type EncodingStreamerMetrics struct {
	EncodedBlobs *prometheus.GaugeVec
}

type TxnManagerMetrics struct {
	Latency  prometheus.Summary
	GasUsed  prometheus.Gauge
	SpeedUps prometheus.Gauge
	TxQueue  prometheus.Gauge
	NumTx    *prometheus.CounterVec
}

type Metrics struct {
	*EncodingStreamerMetrics
	*TxnManagerMetrics

	registry *prometheus.Registry

	Blob             *prometheus.CounterVec
	Batch            *prometheus.CounterVec
	BatchProcLatency *prometheus.SummaryVec
	Attestation      *prometheus.GaugeVec
	BatchError       *prometheus.CounterVec

	httpPort string
	logger   common.Logger
}

func NewMetrics(httpPort string, logger common.Logger) *Metrics {
	namespace := "eigenda_batcher"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	encodingStreamerMetrics := EncodingStreamerMetrics{
		EncodedBlobs: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "encoded_blobs",
				Help:      "number and size of all encoded blobs",
			},
			[]string{"type"},
		),
	}

	txnManagerMetrics := TxnManagerMetrics{
		Latency: promauto.With(reg).NewSummary(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "txn_manager_latency_ms",
				Help:       "transaction confirmation latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
		),
		GasUsed: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "gas_used",
				Help:      "gas used for onchain batch confirmation",
			},
		),
		SpeedUps: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "speed_ups",
				Help:      "number of times the gas price was increased",
			},
		),
		TxQueue: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "tx_queue",
				Help:      "number of transactions in transaction queue",
			},
		),
		NumTx: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tx_total",
				Help:      "number of transactions processed",
			},
			[]string{"state"},
		),
	}

	metrics := &Metrics{
		EncodingStreamerMetrics: &encodingStreamerMetrics,
		TxnManagerMetrics:       &txnManagerMetrics,
		Blob: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "blobs_total",
				Help:      "the number and unencoded size of total dispersal blobs",
			},
			[]string{"state", "data"}, // state is either success or failure
		),
		Batch: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "batches_total",
				Help:      "the number and unencoded size of total dispersal batches",
			},
			[]string{"data"},
		),
		BatchProcLatency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "batch_process_latency_ms",
				Help:       "batch process latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"stage"},
		),
		Attestation: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "attestation",
				Help:      "number of signers and non-signers for the batch",
			},
			[]string{"type"},
		),
		BatchError: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "batch_error",
				Help:      "number of batch errors",
			},
			[]string{"type"},
		),
		registry: reg,
		httpPort: httpPort,
		logger:   logger,
	}
	return metrics
}

func (g *Metrics) UpdateAttestation(operatorCount, nonSignerCount int, quorumResults map[core.QuorumID]*core.QuorumResult) {
	g.Attestation.WithLabelValues("signers").Set(float64(operatorCount - nonSignerCount))
	g.Attestation.WithLabelValues("non_signers").Set(float64(nonSignerCount))

	for _, quorumResult := range quorumResults {
		label := fmt.Sprintf("quorum_result_%d", quorumResult.QuorumID)
		g.Attestation.WithLabelValues(label).Set(float64(quorumResult.PercentSigned))
	}
}

// UpdateCompletedBlob increments the number and updates size of processed blobs.
func (g *Metrics) UpdateCompletedBlob(size int, status disperser.BlobStatus) {
	switch status {
	case disperser.Confirmed:
		g.Blob.WithLabelValues("confirmed", "number").Inc()
		g.Blob.WithLabelValues("confirmed", "size").Add(float64(size))
	case disperser.Failed:
		g.Blob.WithLabelValues("failed", "number").Inc()
		g.Blob.WithLabelValues("failed", "size").Add(float64(size))
	case disperser.InsufficientSignatures:
		g.Blob.WithLabelValues("insufficient_signature", "number").Inc()
		g.Blob.WithLabelValues("insufficient_signature", "size").Add(float64(size))
	default:
		return
	}

	g.Blob.WithLabelValues("total", "number").Inc()
	g.Blob.WithLabelValues("total", "size").Add(float64(size))
}

func (g *Metrics) IncrementBatchCount(size int64) {
	g.Batch.WithLabelValues("number").Inc()
	g.Batch.WithLabelValues("size").Add(float64(size))
}

func (g *Metrics) UpdateBatchError(errType FailReason, numBlobs int) {
	g.BatchError.WithLabelValues(string(errType)).Add(float64(numBlobs))
}

func (g *Metrics) ObserveLatency(stage string, latencyMs float64) {
	g.BatchProcLatency.WithLabelValues(stage).Observe(latencyMs)
}

func (g *Metrics) Start(ctx context.Context) {
	g.logger.Info("starting metrics server at ", "port", g.httpPort)
	addr := fmt.Sprintf(":%s", g.httpPort)
	go func() {
		log := g.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			g.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("prometheus server failed", "err", err)
	}()
}

func (e *EncodingStreamerMetrics) UpdateEncodedBlobs(count int, size uint64) {
	e.EncodedBlobs.WithLabelValues("size").Set(float64(size))
	e.EncodedBlobs.WithLabelValues("number").Set(float64(count))
}

func (t *TxnManagerMetrics) ObserveLatency(latencyMs float64) {
	t.Latency.Observe(latencyMs)
}

func (t *TxnManagerMetrics) UpdateGasUsed(gasUsed uint64) {
	t.GasUsed.Set(float64(gasUsed))
}

func (t *TxnManagerMetrics) UpdateSpeedUps(speedUps int) {
	t.SpeedUps.Set(float64(speedUps))
}

func (t *TxnManagerMetrics) UpdateTxQueue(txQueue int) {
	t.TxQueue.Set(float64(txQueue))
}

func (t *TxnManagerMetrics) IncrementTxnCount(state string) {
	t.NumTx.WithLabelValues(state).Inc()
}
