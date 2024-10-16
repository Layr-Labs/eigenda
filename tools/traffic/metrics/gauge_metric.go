package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// GaugeMetric allows specific values to be reported.
type GaugeMetric interface {
	// Set sets the value of a gauge metric.
	Set(value float64)
}

// gaugeMetric is a standard implementation of the GaugeMetric interface via prometheus.
type gaugeMetric struct {
	metrics     *metrics
	description string
	// disabled specifies whether the metrics should behave as a no-op
	disabled bool
}

// Set sets the value of a gauge metric.
func (metric *gaugeMetric) Set(value float64) {
	if metric.disabled {
		return
	}
	metric.metrics.gauge.WithLabelValues(metric.description).Set(value)
}

// buildGaugeCollector creates a collector for gauge metrics.
func buildGaugeCollector(namespace string, registry *prometheus.Registry) *prometheus.GaugeVec {
	return promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gauge",
		}, []string{"label"},
	)
}
