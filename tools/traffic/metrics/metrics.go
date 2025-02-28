package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics allows the creation of metrics for the traffic generator.
type Metrics interface {
	// Start starts the metrics server.
	Start() error
	// Shutdown shuts down the metrics server.
	Shutdown() error
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

	shutdown func() error
}

// NewMetrics creates a new Metrics instance.
func NewMetrics(
	httpPort string,
	logger logging.Logger,
) Metrics {

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

func (metrics *metrics) Start() error {
	metrics.logger.Info("Starting metrics server", "port", metrics.httpPort)
	addr := fmt.Sprintf(":%s", metrics.httpPort)
	// Create mux and add /metrics handler
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		metrics.registry,
		promhttp.HandlerOpts{},
	))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			metrics.logger.Error("Prometheus server failed", "err", err)
		}
	}()

	// Store shutdown function
	metrics.shutdown = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		metrics.logger.Info("Shutting down metrics server")
		if err := srv.Shutdown(ctx); err != nil {
			metrics.logger.Error("Metrics server shutdown failed", "err", err)
			return err
		}
		return nil
	}

	return nil
}

func (metrics *metrics) Shutdown() error {
	if metrics.shutdown != nil {
		return metrics.shutdown()
	}
	return nil
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
