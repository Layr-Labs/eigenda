package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CountMetric allows the count of a type of event to be tracked.
type CountMetric interface {
	Increment()
}

// countMetric a standard implementation of the CountMetric interface via prometheus.
type countMetric struct {
	metrics     *metrics
	description string
	// disabled specifies whether the metrics should behave as a no-op
	disabled bool
}

// Increment increments the count of a type of event.
func (metric *countMetric) Increment() {
	if metric.disabled {
		return
	}
	metric.metrics.count.WithLabelValues(metric.description).Inc()
}

// buildCounterCollector creates a new prometheus collector for counting metrics.
func buildCounterCollector(namespace string, registry *prometheus.Registry) *prometheus.CounterVec {
	return promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "event_count",
		},
		[]string{"label"},
	)
}
