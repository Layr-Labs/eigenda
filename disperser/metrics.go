package disperser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/eigensdk-go/logging"
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

	NumBlobRequests *prometheus.CounterVec
	BlobSize        *prometheus.GaugeVec
	Latency         *prometheus.SummaryVec

	httpPort string
	logger   logging.Logger
}

// The error space of dispersal requests.
const (
	StoreBlobFailure          string = "store-blob-failed"   // Fail to store the blob (S3 or DynamoDB)
	SystemRateLimitedFailure  string = "ratelimited-system"  // The request rate limited at system level
	AccountRateLimitedFailure string = "ratelimited-account" // The request rate limited at account level
)

func NewMetrics(httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_disperser"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		NumBlobRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests_total",
				Help:      "the number of blob requests",
			},
			[]string{"status", "quorum", "method"},
		),
		BlobSize: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "blob_size_bytes",
				Help:      "the size of the blob in bytes",
			},
			[]string{"status", "quorum", "method"},
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
		registry: reg,
		httpPort: httpPort,
		logger:   logger,
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, latencyMs float64) {
	g.Latency.WithLabelValues(method).Observe(latencyMs)
}

// IncrementSuccessfulBlobRequestNum increments the number of successful blob requests
func (g *Metrics) IncrementSuccessfulBlobRequestNum(quorum string, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status": "success",
		"quorum": quorum,
		"method": method,
	}).Inc()
}

// HandleSuccessfulRequest updates the number of successful blob requests and the size of the blob
func (g *Metrics) HandleSuccessfulRequest(quorum string, blobBytes int, method string) {
	g.IncrementSuccessfulBlobRequestNum(quorum, method)
	g.BlobSize.With(prometheus.Labels{
		"status": "success",
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// IncrementFailedBlobRequestNum increments the number of failed blob requests
func (g *Metrics) IncrementFailedBlobRequestNum(quorum string, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status": "failed",
		"quorum": quorum,
		"method": method,
	}).Inc()
}

// HandleFailedRequest updates the number of failed requests and the size of the blob
func (g *Metrics) HandleFailedRequest(quorum string, blobBytes int, method string) {
	g.IncrementFailedBlobRequestNum(quorum, method)
	g.BlobSize.With(prometheus.Labels{
		"status": "failed",
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// HandleBlobStoreFailedRequest updates the number of requests failed to store blob and the size of the blob
func (g *Metrics) HandleBlobStoreFailedRequest(quorum string, blobBytes int, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status": StoreBlobFailure,
		"quorum": quorum,
		"method": method,
	}).Inc()
	g.BlobSize.With(prometheus.Labels{
		"status": StoreBlobFailure,
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// HandleSystemRateLimitedRequest updates the number of system rate limited requests and the size of the blob
func (g *Metrics) HandleSystemRateLimitedRequest(quorum string, blobBytes int, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status": SystemRateLimitedFailure,
		"quorum": quorum,
		"method": method,
	}).Inc()
	g.BlobSize.With(prometheus.Labels{
		"status": SystemRateLimitedFailure,
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// HandleAccountRateLimitedRequest updates the number of account rate limited requests and the size of the blob
func (g *Metrics) HandleAccountRateLimitedRequest(quorum string, blobBytes int, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status": AccountRateLimitedFailure,
		"quorum": quorum,
		"method": method,
	}).Inc()
	g.BlobSize.With(prometheus.Labels{
		"status": AccountRateLimitedFailure,
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
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
