package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CountMetric tracks the count of a type of event.
type CountMetric struct {
	metrics     *Metrics
	description string
}

// Increment increments the count of a type of event.
func (metric *CountMetric) Increment() {
	metric.metrics.count.WithLabelValues(metric.description).Inc()
}

// NewCountMetric creates a new prometheus collector for counting metrics.
func buildCounterCollector(namespace string, registry *prometheus.Registry) *prometheus.CounterVec {
	return promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "event_count",
		},
		[]string{"label"},
	)
}
