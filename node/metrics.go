package node

import (
	"context"
	"strconv"

	"github.com/Layr-Labs/eigensdk-go/logging"
	eigenmetrics "github.com/Layr-Labs/eigensdk-go/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace = "node"
)

type Metrics struct {
	logger logging.Logger

	// Whether the node is registered.
	Registered prometheus.Gauge
	// Accumulated number of RPC requests received.
	AccNumRequests *prometheus.CounterVec
	// The latency (in ms) to process the request.
	RequestLatency *prometheus.SummaryVec
	// Accumulated number and size of batches processed by their statuses.
	// TODO: this metrics should be removed in future release.
	AccuBatchesDeprecated *prometheus.CounterVec
	// Accumulated number and size of batches processed by their statuses.
	AccuBatches *prometheus.CounterVec
	// Accumulated number and size of batches that have been removed from the Node.
	AccuRemovedBatches *prometheus.CounterVec
	// Accumulated number and size of blobs processed by quorums.
	AccuBlobs *prometheus.CounterVec
	// Total number of changes in the node's socket address.
	AccuSocketUpdates prometheus.Counter
	// avs node spec eigen_ metrics: https://eigen.nethermind.io/docs/spec/metrics/metrics-prom-spec
	EigenMetrics eigenmetrics.Metrics

	registry *prometheus.Registry
	// socketAddr is the address at which the metrics server will be listening.
	// should be in format ip:port
	socketAddr string
}

func NewMetrics(eigenMetrics eigenmetrics.Metrics, reg *prometheus.Registry, logger logging.Logger, socketAddr string) *Metrics {

	// Add Go module collectors
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		Registered: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "registered",
				Help:      "indicator about if DA node is registered",
			},
		),
		// The "status" label has values: success, failure.
		AccNumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_rpc_requests_total",
				Help:      "the total number of requests processed by the DA node",
			},
			[]string{"method", "status"},
		),
		RequestLatency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  Namespace,
				Name:       "request_latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"method", "stage"},
		),
		AccuBatchesDeprecated: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_batches_total",
				Help:      "the total number and size of batches processed by the DA node",
			},
			[]string{"type", "status"},
		),
		AccuBlobs: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_blobs_total",
				Help:      "the total number and size of blobs processed by the DA node",
			},
			[]string{"type", "quorum"},
		),
		// The "status" label has values: received, validated, stored, signed.
		// These are the lifecycle of a batch at the DA Node.
		AccuBatches: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_processed_batches_total",
				Help:      "the total number and size of batches processed by the DA node",
			},
			[]string{"type", "status"},
		),
		AccuRemovedBatches: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_removed_batches_total",
				Help:      "the total number and size of batches that have been removed by the DA node",
			},
			[]string{"type"},
		),
		AccuSocketUpdates: promauto.With(reg).NewCounter(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "eigenda_node_socket_updates_total",
				Help:      "the total number of node's socket address updates",
			},
		),
		EigenMetrics: eigenMetrics,
		logger:       logger,
		registry:     reg,
		socketAddr:   socketAddr,
	}

	return metrics
}

func (g *Metrics) Start() {
	_ = g.EigenMetrics.Start(context.Background(), g.registry)
}

func (g *Metrics) RecordRPCRequest(method string, status string) {
	g.AccNumRequests.WithLabelValues(method, status).Inc()
}

func (g *Metrics) RecordSocketAddressChange() {
	g.AccuSocketUpdates.Inc()
}

func (g *Metrics) ObserveLatency(method, stage string, latencyMs float64) {
	g.RequestLatency.WithLabelValues(method, stage).Observe(latencyMs)
}

func (g *Metrics) RemoveNCurrentBatch(numBatches int, totalBatchSize int64) {
	for i := 0; i < numBatches; i++ {
		g.AccuRemovedBatches.WithLabelValues("number").Inc()
	}
	g.AccuRemovedBatches.WithLabelValues("size").Add(float64(totalBatchSize))
}

func (g *Metrics) AcceptBlobs(quorumId uint8, blobSize int64) {
	quorum := strconv.Itoa(int(quorumId))
	g.AccuBlobs.WithLabelValues("number", quorum).Inc()
	g.AccuBlobs.WithLabelValues("size", quorum).Add(float64(blobSize))
}

func (g *Metrics) AcceptBatches(status string, batchSize int64) {
	g.AccuBatchesDeprecated.WithLabelValues("number", status).Inc()
	g.AccuBatchesDeprecated.WithLabelValues("size", status).Add(float64(batchSize))
	g.AccuBatches.WithLabelValues("number", status).Inc()
	g.AccuBatches.WithLabelValues("size", status).Add(float64(batchSize))
}
