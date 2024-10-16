package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// LatencyMetric allows the latency of an operation to be tracked.
type LatencyMetric interface {
	ReportLatency(latency time.Duration)
}

// latencyMetric is a standard implementation of the LatencyMetric interface via prometheus.
type latencyMetric struct {
	metrics     *metrics
	description string
	// disabled specifies whether the metrics should behave as a no-op
	disabled bool
}

// ReportLatency reports the latency of an operation.
func (metric *latencyMetric) ReportLatency(latency time.Duration) {
	if metric.disabled {
		return
	}
	metric.metrics.latency.WithLabelValues(metric.description).Observe(latency.Seconds())
}

// buildLatencyCollector creates a new prometheus collector for latency metrics.
func buildLatencyCollector(namespace string, registry *prometheus.Registry) *prometheus.SummaryVec {
	return promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "latency_s",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
		},
		[]string{"label"},
	)
}
