package encoder

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	logger   logging.Logger
	registry *prometheus.Registry
	httpPort string

	NumEncodeBlobRequests *prometheus.CounterVec
	BlobSizeTotal         *prometheus.CounterVec
	Latency               *prometheus.SummaryVec
	BlobQueue             *prometheus.GaugeVec
	QueueCapacity         prometheus.Gauge
	QueueUtilization      prometheus.Gauge
}

func NewMetrics(reg *prometheus.Registry, httpPort string, logger logging.Logger) *Metrics {
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	return &Metrics{
		logger:   logger.With("component", "EncoderMetrics"),
		registry: reg,
		httpPort: httpPort,
		NumEncodeBlobRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "eigenda_encoder",
				Name:      "request_total",
				Help:      "the number of total encode blob request at server side per state",
			},
			[]string{"state"}, // state is either success, ratelimited, canceled, or failure
		),
		BlobSizeTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "eigenda_encoder",
				Name:      "blob_size_total",
				Help:      "the size in bytes of total blob requests at server side per state",
			},
			[]string{"state"}, // state is either success, ratelimited, canceled, or failure
		),
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  "eigenda_encoder",
				Name:       "encoding_latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"time"}, // time is either encoding or total
		),
		BlobQueue: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "eigenda_encoder",
				Name:      "blob_queue",
				Help:      "the number of blobs in the queue for encoding",
			},
			[]string{"size_bucket"},
		),
		QueueCapacity: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: "eigenda_encoder",
				Name:      "request_pool_capacity",
				Help:      "The maximum capacity of the request pool",
			},
		),
		QueueUtilization: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: "eigenda_encoder",
				Name:      "request_pool_utilization",
				Help:      "Current utilization of request pool (total across all buckets)",
			},
		),
	}
}

// IncrementSuccessfulBlobRequestNum increments the number of successful requests
// this counter incrementation is atomic
func (m *Metrics) IncrementSuccessfulBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("success").Inc()
	m.BlobSizeTotal.WithLabelValues("success").Add(float64(blobSize))
}

// IncrementFailedBlobRequestNum increments the number of failed requests
// this counter incrementation is atomic
func (m *Metrics) IncrementFailedBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("failed").Inc()
	m.BlobSizeTotal.WithLabelValues("failed").Add(float64(blobSize))
}

// IncrementRateLimitedBlobRequestNum increments the number of rate limited requests
// this counter incrementation is atomic
func (m *Metrics) IncrementRateLimitedBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("ratelimited").Inc()
	m.BlobSizeTotal.WithLabelValues("ratelimited").Add(float64(blobSize))
}

// IncrementCanceledBlobRequestNum increments the number of canceled requests
// this counter incrementation is atomic
func (m *Metrics) IncrementCanceledBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("canceled").Inc()
	m.BlobSizeTotal.WithLabelValues("canceled").Add(float64(blobSize))
}

func (m *Metrics) ObserveLatency(stage string, duration time.Duration) {
	m.Latency.WithLabelValues(stage).Observe(float64(duration.Milliseconds()))
}

func (m *Metrics) ObserveQueue(queueStats map[string]int) {
	total := 0
	for bucket, num := range queueStats {
		m.BlobQueue.With(prometheus.Labels{"size_bucket": bucket}).Set(float64(num))
		total += num
	}
	m.QueueUtilization.Set(float64(total))
}

func (m *Metrics) SetQueueCapacity(capacity int) {
	m.QueueCapacity.Set(float64(capacity))
}

func (m *Metrics) Start(ctx context.Context) {
	m.logger.Info("Starting metrics server at ", "port", m.httpPort)

	addr := fmt.Sprintf(":%s", m.httpPort)
	log := m.logger

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))

	server := &http.Server{Addr: addr, Handler: mux}
	errc := make(chan error, 1)

	go func() {
		errc <- server.ListenAndServe()
	}()
	go func() {
		select {
		case <-ctx.Done():
			m.shutdown(server)
			return
		case err := <-errc:
			log.Error("Prometheus server failed", "err", err)
		}
	}()
}

func (m *Metrics) shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
