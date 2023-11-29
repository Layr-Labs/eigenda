package batcher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	registry *prometheus.Registry

	Blob             *prometheus.CounterVec
	Batch            *prometheus.CounterVec
	BatchProcLatency *prometheus.SummaryVec
	GasUsed          prometheus.Gauge
	Attestation      *prometheus.GaugeVec

	httpPort string
	logger   common.Logger
}

func NewMetrics(httpPort string, logger common.Logger) *Metrics {
	namespace := "eigenda_batcher"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		Blob: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "blobs_total",
				Help:      "the number and encoded size of total dispersal blobs",
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
		GasUsed: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "gas_used",
				Help:      "gas used for onchain batch confirmation",
			},
		),
		Attestation: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "attestation",
				Help:      "number of signers and non-signers for the batch",
			},
			[]string{"type"},
		),
		registry: reg,
		httpPort: httpPort,
		logger:   logger,
	}
	return metrics
}

func (g *Metrics) UpdateAttestation(operatorCount, nonSignerCount int) {
	g.Attestation.WithLabelValues("signers").Set(float64(operatorCount - nonSignerCount))
	g.Attestation.WithLabelValues("non_signers").Set(float64(nonSignerCount))
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
