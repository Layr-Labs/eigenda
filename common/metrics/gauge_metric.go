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

	// unit is the unit of the metric.
	unit string

	// description is the description of the metric.
	description string

	// gauge is the prometheus gauge used to report this metric.
	gauge prometheus.Gauge
}

// newGaugeMetric creates a new GaugeMetric instance.
func newGaugeMetric(
	name string,
	label string,
	unit string,
	description string,
	vec *prometheus.GaugeVec) GaugeMetric {

	var gauge prometheus.Gauge
	if vec != nil {
		gauge = vec.WithLabelValues(label)
	}

	return &gaugeMetric{
		name:        name,
		label:       label,
		unit:        unit,
		description: description,
		gauge:       gauge,
	}
}

func (m *gaugeMetric) Name() string {
	return m.name
}

func (m *gaugeMetric) Label() string {
	return m.label
}

func (m *gaugeMetric) Unit() string {
	return m.unit
}

func (m *gaugeMetric) Description() string {
	return m.description
}

func (m *gaugeMetric) Type() string {
	return "gauge"
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
