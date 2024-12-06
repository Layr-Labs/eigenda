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
	"google.golang.org/grpc/codes"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	registry *prometheus.Registry

	NumBlobRequests *prometheus.CounterVec
	NumRpcRequests  *prometheus.CounterVec
	BlobSize        *prometheus.GaugeVec
	BlobLatency     *prometheus.GaugeVec
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

func NewMetrics(reg *prometheus.Registry, httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_disperser"
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		// TODO: revamp this metric -- it'll focus on quorum tracking, which is relevant
		// only for the Disperser.DisperserBlob API.
		NumBlobRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "requests_total",
				Help:      "the number of blob requests",
			},
			[]string{"status_code", "status", "quorum", "method"},
		),
		NumRpcRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "grpc_requests_total",
				Help:      "the number of gRPC requests",
			},
			[]string{"status_code", "status_detail", "method"},
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
		BlobLatency: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "blob_latency_ms",
				Help:      "blob dispersal or retrieval latency by size",
			},
			[]string{"method", "size_bucket"},
		),

		registry: reg,
		httpPort: httpPort,
		logger:   logger.With("component", "DisperserMetrics"),
	}
	return metrics
}

// ObserveLatency observes the latency of a stage in 'stage
func (g *Metrics) ObserveLatency(method string, latencyMs float64) {
	g.Latency.WithLabelValues(method).Observe(latencyMs)
}

func (g *Metrics) HandleSuccessfulRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.OK.String(),
		"status_detail": "",
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleInvalidArgRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.InvalidArgument.String(),
		"status_detail": "",
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleNotFoundRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.NotFound.String(),
		"status_detail": "",
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleSystemRateLimitedRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.ResourceExhausted.String(),
		"status_detail": SystemRateLimitedFailure,
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleAccountRateLimitedRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.ResourceExhausted.String(),
		"status_detail": AccountRateLimitedFailure,
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleRateLimitedRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.ResourceExhausted.String(),
		"status_detail": "",
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleInternalFailureRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.Internal.String(),
		"status_detail": "",
		"method":        method,
	}).Inc()
}

func (g *Metrics) HandleStoreFailureRpcRequest(method string) {
	g.NumRpcRequests.With(prometheus.Labels{
		"status_code":   codes.Internal.String(),
		"status_detail": StoreBlobFailure,
		"method":        method,
	}).Inc()
}

// IncrementSuccessfulBlobRequestNum increments the number of successful blob requests
func (g *Metrics) IncrementSuccessfulBlobRequestNum(quorum string, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": codes.OK.String(),
		"status":      "success",
		"quorum":      quorum,
		"method":      method,
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
func (g *Metrics) IncrementFailedBlobRequestNum(statusCode string, quorum string, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": statusCode,
		"status":      "failed",
		"quorum":      quorum,
		"method":      method,
	}).Inc()
}

// HandleFailedRequest updates the number of failed requests and the size of the blob
func (g *Metrics) HandleFailedRequest(statusCode string, quorum string, blobBytes int, method string) {
	g.IncrementFailedBlobRequestNum(statusCode, quorum, method)
	g.BlobSize.With(prometheus.Labels{
		"status": "failed",
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// HandleBlobStoreFailedRequest updates the number of requests failed to store blob and the size of the blob
func (g *Metrics) HandleBlobStoreFailedRequest(quorum string, blobBytes int, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": codes.Internal.String(),
		"status":      StoreBlobFailure,
		"quorum":      quorum,
		"method":      method,
	}).Inc()
	g.BlobSize.With(prometheus.Labels{
		"status": StoreBlobFailure,
		"quorum": quorum,
		"method": method,
	}).Add(float64(blobBytes))
}

// HandleInvalidArgRequest updates the number of invalid argument requests
func (g *Metrics) HandleInvalidArgRequest(method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": codes.InvalidArgument.String(),
		"status":      "failed",
		"quorum":      "",
		"method":      method,
	}).Inc()
}

// HandleInvalidArgRequest updates the number of invalid argument requests
func (g *Metrics) HandleNotFoundRequest(method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": codes.NotFound.String(),
		"status":      "failed",
		"quorum":      "",
		"method":      method,
	}).Inc()
}

// HandleSystemRateLimitedRequest updates the number of system rate limited requests and the size of the blob
func (g *Metrics) HandleSystemRateLimitedRequest(quorum string, blobBytes int, method string) {
	g.NumBlobRequests.With(prometheus.Labels{
		"status_code": codes.ResourceExhausted.String(),
		"status":      SystemRateLimitedFailure,
		"quorum":      quorum,
		"method":      method,
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
		"status_code": codes.ResourceExhausted.String(),
		"status":      AccountRateLimitedFailure,
		"quorum":      quorum,
		"method":      method,
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
