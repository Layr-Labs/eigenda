package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GaugeMetric allows values to be reported.
type GaugeMetric struct {
	metrics     *Metrics
	description string
}

// Set sets the value of a gauge metric.
func (metric GaugeMetric) Set(value float64) {
	metric.metrics.gauge.WithLabelValues(metric.description).Set(value)
}

// NewGaugeMetric creates a collector for gauge metrics.
func buildGaugeCollector(namespace string, registry *prometheus.Registry) *prometheus.GaugeVec {
	return promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gauge",
		}, []string{"label"},
	)
}
