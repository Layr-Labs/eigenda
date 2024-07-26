package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
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

// NewMetrics creates a new Metrics instance.
func NewMetrics(httpPort string, logger logging.Logger) *Metrics {
	namespace := "eigenda_generator"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &Metrics{
		count:    buildCounterCollector(namespace, reg),
		latency:  buildLatencyCollector(namespace, reg),
		gauge:    buildGaugeCollector(namespace, reg),
		registry: reg,
		httpPort: httpPort,
		logger:   logger.With("component", "GeneratorMetrics"),
	}
	return metrics
}

// Start starts the metrics server.
func (metrics *Metrics) Start() {
	metrics.logger.Info("Starting metrics server at ", "port", metrics.httpPort)
	addr := fmt.Sprintf(":%s", metrics.httpPort)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			metrics.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		panic(fmt.Sprintf("Prometheus server failed: %s", err))
	}()
}

// NewLatencyMetric creates a new LatencyMetric instance.
func (metrics *Metrics) NewLatencyMetric(description string) LatencyMetric {
	return LatencyMetric{
		metrics:     metrics,
		description: description,
	}
}

// NewCountMetric creates a new CountMetric instance.
func (metrics *Metrics) NewCountMetric(description string) CountMetric {
	return CountMetric{
		metrics:     metrics,
		description: description,
	}
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (metrics *Metrics) NewGaugeMetric(description string) GaugeMetric {
	return GaugeMetric{
		metrics:     metrics,
		description: description,
	}
}
