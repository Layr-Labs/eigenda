package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Metrics allows the creation of metrics for the traffic generator.
type Metrics interface {
	// Start starts the metrics server.
	Start()
	// NewLatencyMetric creates a new LatencyMetric instance. Useful for reporting the latency of an operation.
	NewLatencyMetric(description string) LatencyMetric
	// NewCountMetric creates a new CountMetric instance. Useful for tracking the count of a type of event.
	NewCountMetric(description string) CountMetric
	// NewGaugeMetric creates a new GaugeMetric instance. Useful for reporting specific values.
	NewGaugeMetric(description string) GaugeMetric
}

// metrics is a standard implementation of the Metrics interface via prometheus.
type metrics struct {
	registry *prometheus.Registry

	count   *prometheus.CounterVec
	latency *prometheus.SummaryVec
	gauge   *prometheus.GaugeVec

	httpPort string
	logger   logging.Logger
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(httpPort string, logger logging.Logger) Metrics {
	namespace := "eigenda_generator"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	metrics := &metrics{
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
func (metrics *metrics) Start() {
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
func (metrics *metrics) NewLatencyMetric(description string) LatencyMetric {
	return &latencyMetric{
		metrics:     metrics,
		description: description,
	}
}

// NewCountMetric creates a new CountMetric instance.
func (metrics *metrics) NewCountMetric(description string) CountMetric {
	return &countMetric{
		metrics:     metrics,
		description: description,
	}
}

// NewGaugeMetric creates a new GaugeMetric instance.
func (metrics *metrics) NewGaugeMetric(description string) GaugeMetric {
	return &gaugeMetric{
		metrics:     metrics,
		description: description,
	}
}
