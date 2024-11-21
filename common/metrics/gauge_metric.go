package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _ GaugeMetric = &gaugeMetric{}

// gaugeMetric is a standard implementation of the GaugeMetric interface via prometheus.
type gaugeMetric struct {
	Metric

	// name is the name of the metric.
	name string

	// label is the label of the metric.
	label string

	// gauge is the prometheus gauge used to report this metric.
	gauge prometheus.Gauge
}

// newGaugeMetric creates a new GaugeMetric instance.
func newGaugeMetric(name string, label string, vec *prometheus.GaugeVec) GaugeMetric {
	var gauge prometheus.Gauge
	if vec != nil {
		gauge = vec.WithLabelValues(label)
	}

	return &gaugeMetric{
		name:  name,
		label: label,
		gauge: gauge,
	}
}

func (m *gaugeMetric) Name() string {
	return m.name
}

func (m *gaugeMetric) Label() string {
	return m.label
}

func (m *gaugeMetric) Enabled() bool {
	return m.gauge != nil
}

func (m *gaugeMetric) Set(value float64) {
	if m.gauge == nil {
		return
	}
	m.gauge.Set(value)
}
