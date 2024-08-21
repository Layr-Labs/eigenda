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

type MetrisConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type Metrics struct {
	logger   logging.Logger
	registry *prometheus.Registry
	httpPort string

	NumEncodeBlobRequests *prometheus.CounterVec
	BlobSize              *prometheus.CounterVec
	Latency               *prometheus.SummaryVec
}

func NewMetrics(httpPort string, logger logging.Logger) *Metrics {
	reg := prometheus.NewRegistry()
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
				Help:      "the number and size of total encode blob request at server side per state",
			},
			[]string{"type", "state"}, // type is either number or size; state is either success, ratelimited, canceled, or failure
		),
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  "eigenda_encoder",
				Name:       "encoding_latency_ms",
				Help:       "latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"size", "time"}, // size is the bucket of the blob size, time is either encoding or total
		),
	}
}

// IncrementSuccessfulBlobRequestNum increments the number of successful requests
// this counter incrementation is atomic
func (m *Metrics) IncrementSuccessfulBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("number", "success").Inc()
	m.NumEncodeBlobRequests.WithLabelValues("size", "success").Add(float64(blobSize))
}

// IncrementFailedBlobRequestNum increments the number of failed requests
// this counter incrementation is atomic
func (m *Metrics) IncrementFailedBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("number", "failed").Inc()
	m.NumEncodeBlobRequests.WithLabelValues("size", "failed").Add(float64(blobSize))
}

// IncrementRateLimitedBlobRequestNum increments the number of rate limited requests
// this counter incrementation is atomic
func (m *Metrics) IncrementRateLimitedBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("number", "ratelimited").Inc()
	m.NumEncodeBlobRequests.WithLabelValues("size", "ratelimited").Add(float64(blobSize))
}

// IncrementCanceledBlobRequestNum increments the number of canceled requests
// this counter incrementation is atomic
func (m *Metrics) IncrementCanceledBlobRequestNum(blobSize int) {
	m.NumEncodeBlobRequests.WithLabelValues("number", "canceled").Inc()
	m.NumEncodeBlobRequests.WithLabelValues("size", "canceled").Add(float64(blobSize))
}

// BlobSizeBucket maps the blob size into a bucket that's defined according to
// the power of 2.
func BlobSizeBucket(blobSize int) string {
	switch {
	case blobSize <= 32*1024:
		return "32KiB"
	case blobSize <= 64*1024:
		return "64KiB"
	case blobSize <= 128*1024:
		return "128KiB"
	case blobSize <= 256*1024:
		return "256KiB"
	case blobSize <= 512*1024:
		return "512KiB"
	case blobSize <= 1024*1024:
		return "1MiB"
	case blobSize <= 2*1024*1024:
		return "2MiB"
	case blobSize <= 4*1024*1024:
		return "4MiB"
	case blobSize <= 8*1024*1024:
		return "8MiB"
	case blobSize <= 16*1024*1024:
		return "16MiB"
	case blobSize <= 32*1024*1024:
		return "32MiB"
	default:
		return "invalid"
	}
}

func (m *Metrics) TakeLatency(blobSize int, encoding, total time.Duration) {
	size := BlobSizeBucket(blobSize)
	m.Latency.WithLabelValues(size, "encoding").Observe(float64(encoding.Milliseconds()))
	m.Latency.WithLabelValues(size, "total").Observe(float64(total.Milliseconds()))
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
