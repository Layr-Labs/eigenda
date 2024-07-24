package traffic

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// Metrics encapsulates metrics for the traffic generator.
type Metrics struct {
	registry *prometheus.Registry

	count   *prometheus.CounterVec
	latency *prometheus.SummaryVec
	gauge   *prometheus.GaugeVec

	httpPort string
	logger   logging.Logger
}

// LatencyMetric tracks the latency of an operation.
type LatencyMetric struct {
	metrics     *Metrics
	description string
}

// CountMetric tracks the count of a type of event.
type CountMetric struct {
	metrics     *Metrics
	description string
}

type GaugeMetric struct {
	metrics     *Metrics
	description string
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_generator"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		count: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "event_count",
			},
			[]string{"label"},
		),
		latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "latency_s",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"label"},
		),
		gauge: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "gauge",
			}, []string{"label"}),
		registry: reg,
		httpPort: httpPort,
		logger:   logger.With("component", "GeneratorMetrics"),
	}
	return metrics
}

// Start starts the metrics server.
func (metrics *Metrics) Start(ctx context.Context) { // TODO context?
	metrics.logger.Info("Starting metrics server at ", "port", metrics.httpPort)
	addr := fmt.Sprintf(":%s", metrics.httpPort)
	go func() {
		log := metrics.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			metrics.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("Prometheus server failed", "err", err)
	}()
}

// NewLatencyMetric creates a new LatencyMetric instance.
func (metrics *Metrics) NewLatencyMetric(description string) LatencyMetric {
	return LatencyMetric{
		metrics:     metrics,
		description: description,
	}
}

// ReportLatency reports the latency of an operation.
func (metric *LatencyMetric) ReportLatency(latency time.Duration) {
	metric.metrics.latency.WithLabelValues(metric.description).Observe(latency.Seconds())
}

// InvokeAndReportLatency performs an operation. If the operation does not produce an error, then the latency
// of the operation is reported to the metrics framework.
func InvokeAndReportLatency[T any](metric *LatencyMetric, operation func() (T, error)) (T, error) {
	start := time.Now()

	t, err := operation()

	if err == nil {
		end := time.Now()
		duration := end.Sub(start)
		metric.ReportLatency(duration)
	}

	return t, err
}

// NewCountMetric creates a new CountMetric instance.
func (metrics *Metrics) NewCountMetric(description string) CountMetric {
	return CountMetric{
		metrics:     metrics,
		description: description,
	}
}

// Increment increments the count of a type of event.
func (metric *CountMetric) Increment() {
	metric.metrics.count.WithLabelValues(metric.description).Inc()
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (metrics *Metrics) NewGaugeMetric(description string) GaugeMetric {
	return GaugeMetric{
		metrics:     metrics,
		description: description,
	}
}

// Set sets the value of a gauge metric.
func (metric GaugeMetric) Set(value float64) {
	metric.metrics.gauge.WithLabelValues(metric.description).Set(value)
}
