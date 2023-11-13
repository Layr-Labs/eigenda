package encoder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/common"
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
	logger   common.Logger
	registry *prometheus.Registry
	httpPort string

	NumEncodeBlobRequests *prometheus.CounterVec
	Latency               *prometheus.SummaryVec
}

func NewMetrics(httpPort string, logger common.Logger) *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	return &Metrics{
		logger:   logger,
		registry: reg,
		httpPort: httpPort,
		NumEncodeBlobRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "eigenda_encoder",
				Name:      "request_total",
				Help:      "the number and size of total encode blob request at server side per state",
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
			[]string{"time"},
		),
	}
}

// IncrementSuccessfulBlobRequestNum increments the number of successful requests
// this counter incrementation is atomic
func (m *Metrics) IncrementSuccessfulBlobRequestNum() {
	m.NumEncodeBlobRequests.WithLabelValues("success").Inc()
}

// IncrementFailedBlobRequestNum increments the number of failed requests
// this counter incrementation is atomic
func (m *Metrics) IncrementFailedBlobRequestNum() {
	m.NumEncodeBlobRequests.WithLabelValues("failed").Inc()
}

// IncrementRateLimitedBlobRequestNum increments the number of rate limited requests
// this counter incrementation is atomic
func (m *Metrics) IncrementRateLimitedBlobRequestNum() {
	m.NumEncodeBlobRequests.WithLabelValues("ratelimited").Inc()
}

// IncrementCanceledBlobRequestNum increments the number of canceled requests
// this counter incrementation is atomic
func (m *Metrics) IncrementCanceledBlobRequestNum() {
	m.NumEncodeBlobRequests.WithLabelValues("canceled").Inc()
}

func (m *Metrics) TakeLatency(encoding, total time.Duration) {
	m.Latency.WithLabelValues("encoding").Observe(float64(encoding.Milliseconds()))
	m.Latency.WithLabelValues("total").Observe(float64(total.Milliseconds()))
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
