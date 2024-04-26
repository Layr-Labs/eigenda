package dataapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	registry *prometheus.Registry

	NumRequests *prometheus.CounterVec
	Latency     *prometheus.SummaryVec

	EjectionRequests *prometheus.CounterVec
	OperatorsEjected *prometheus.CounterVec
	StakeEjected     *prometheus.CounterVec
	EjectionGasUsed  prometheus.Gauge

	httpPort string
	logger   logging.Logger
}

func NewMetrics(blobMetadataStore *blobstore.BlobMetadataStore, httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_dataapi"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(NewDynamoDBCollector(blobMetadataStore, logger))
	metrics := &Metrics{
		NumRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests",
				Help:      "the number of requests",
			},
			[]string{"status", "method"},
		),
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"method"},
		),
		// EjectionRequests is a more detailed metric than NumRequests, specifically for tracking
		// the ejection calls.
		// The "mode" could be:
		// - "periodic": periodically initiated ejection; or
		// - "urgent": urgently invoked ejection in case of bad network health condition.
		// The "status" indicates the final processing result of the ejection request.
		EjectionRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "ejection_requests_total",
				Help:      "the total number of ejection requests",
			},
			[]string{"status", "mode"},
		),
		// The "state" could be:
		// - "requested": operator is requested for ejection; or
		// - "ejected": operator is actually ejected
		OperatorsEjected: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "operators_total",
				Help:      "the total number of operators to be ejected or actually ejected",
			}, []string{"quorum", "state"},
		),
		StakeEjected: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "stake_ejected",
				Help:      "the total amount of stake to be ejected or actually ejected",
			}, []string{"quorum", "state"},
		),
		EjectionGasUsed: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "ejection_gas_used",
				Help:      "Gas used for operator ejection",
			},
		),
		registry: reg,
		httpPort: httpPort,
		logger:   logger.With("component", "DataAPIMetrics"),
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, latencyMs float64) {
	g.Latency.WithLabelValues(method).Observe(latencyMs)
}

// IncrementSuccessfulRequestNum increments the number of successful requests
func (g *Metrics) IncrementSuccessfulRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "success",
		"method": method,
	}).Inc()
}

// IncrementFailedRequestNum increments the number of failed requests
func (g *Metrics) IncrementFailedRequestNum(method string) {
	g.NumRequests.With(prometheus.Labels{
		"status": "failed",
		"method": method,
	}).Inc()
}

func (g *Metrics) IncrementEjectionRequest(mode string, status codes.Code) {
	g.EjectionRequests.With(prometheus.Labels{
		"status": status.String(),
		"mode":   mode,
	}).Inc()
}

func (g *Metrics) UpdateRequestedOperatorMetric(numOperatorsByQuorum map[uint8]int, stakeShareByQuorum map[uint8]float64) {
	for q, count := range numOperatorsByQuorum {
		for i := 0; i < count; i++ {
			g.OperatorsEjected.With(prometheus.Labels{
				"quorum": fmt.Sprintf("%d", q),
				"state":  "requested",
			}).Inc()
		}
	}
	for q, stakeShare := range stakeShareByQuorum {
		g.StakeEjected.With(prometheus.Labels{
			"quorum": fmt.Sprintf("%d", q),
			"state":  "requested",
		}).Add(stakeShare)
	}
}

func (g *Metrics) UpdateEjectionGasUsed(gasUsed uint64) {
	g.EjectionGasUsed.Set(float64(gasUsed))
}

// Start starts the metrics server
func (g *Metrics) Start(ctx context.Context) {
	g.logger.Info("Starting metrics server at ", "port", g.httpPort)
	addr := fmt.Sprintf(":%s", g.httpPort)
	go func() {
		log := g.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			g.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("Prometheus server failed", "err", err)
	}()
}

type DynamoDBCollector struct {
	blobMetadataStore *blobstore.BlobMetadataStore
	blobStatusMetric  *prometheus.Desc
	logger            logging.Logger
}

func NewDynamoDBCollector(blobMetadataStore *blobstore.BlobMetadataStore, logger logging.Logger) *DynamoDBCollector {
	return &DynamoDBCollector{
		blobMetadataStore: blobMetadataStore,
		blobStatusMetric: prometheus.NewDesc("dynamodb_blob_metadata_status_count",
			"Number of blobs with specific status in DynamoDB",
			[]string{"status"},
			nil,
		),
		logger: logger,
	}
}

func (collector *DynamoDBCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.blobStatusMetric
}

func (collector *DynamoDBCollector) Collect(ch chan<- prometheus.Metric) {
	count, err := collector.blobMetadataStore.GetBlobMetadataByStatusCount(context.Background(), disperser.Processing)
	if err != nil {
		collector.logger.Error("failed to get count of blob metadata by status", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		collector.blobStatusMetric,
		prometheus.GaugeValue,
		float64(count),
		disperser.Processing.String(),
	)
}
