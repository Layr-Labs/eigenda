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
				Help:      "the number and size of total dispersal blob",
			},
			[]string{"state", "data"}, // state is either success or failure
		),
		Batch: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "batches_total",
				Help:      "the number and size of total dispersal batch",
			},
			[]string{"state", "data"},
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

func (g *Metrics) UpdateAttestation(signers, nonSigners int) {
	g.Attestation.WithLabelValues("signers").Set(float64(signers))
	g.Attestation.WithLabelValues("non_signers").Set(float64(nonSigners))
}

// UpdateFailedBatchAndBlobs updates failed a batch and number of blob within it, it only
// counts the number of blob and batches
func (g *Metrics) UpdateFailedBatchAndBlobs(numBlob int) {
	g.Blob.WithLabelValues("failed", "number").Add(float64(numBlob))
	g.Batch.WithLabelValues("failed", "number").Inc()
}

// UpdateCompletedBatchAndBlobs updates whenever there is a completed batch. it updates both the
// number for batch and blob, and it updates size of data blob. Moreover, it updates the
// time it takes to process the entire batch from "getting the blobs" to "marking as finished"
func (g *Metrics) UpdateCompletedBatchAndBlobs(blobsInBatch []*disperser.BlobMetadata, succeeded []bool) {
	totalBlobSucceeded := 0
	totalBlobFailed := 0
	totalBlobSize := 0

	for ind, metadata := range blobsInBatch {
		if succeeded[ind] {
			totalBlobSucceeded += 1
			totalBlobSize += int(metadata.RequestMetadata.BlobSize)
		} else {
			totalBlobFailed += 1
		}
	}

	// Failed blob
	g.UpdateFailedBatchAndBlobs(totalBlobFailed)

	// Blob
	g.Blob.WithLabelValues("completed", "number").Add(float64(totalBlobSucceeded))
	g.Blob.WithLabelValues("completed", "size").Add(float64(totalBlobSize))
	// Batch
	g.Batch.WithLabelValues("completed", "number").Inc()
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
