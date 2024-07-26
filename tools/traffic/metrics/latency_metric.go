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
}

// ReportLatency reports the latency of an operation.
func (metric *latencyMetric) ReportLatency(latency time.Duration) {
	metric.metrics.latency.WithLabelValues(metric.description).Observe(latency.Seconds())
}

// InvokeAndReportLatency performs an operation. If the operation does not produce an error, then the latency
// of the operation is reported to the metrics framework.
func InvokeAndReportLatency[T any](metric LatencyMetric, operation func() (T, error)) (T, error) {
	start := time.Now()

	t, err := operation()

	if err == nil {
		end := time.Now()
		duration := end.Sub(start)
		metric.ReportLatency(duration)
	}

	return t, err
}

// NewLatencyMetric creates a new prometheus collector for latency metrics.
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
