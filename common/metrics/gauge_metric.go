package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var _ GaugeMetric = &gaugeMetric{}

// gaugeMetric is a standard implementation of the GaugeMetric interface via prometheus.
type gaugeMetric struct {
	Metric

	// name is the name of the metric.
	name string

	// unit is the unit of the metric.
	unit string

	// description is the description of the metric.
	description string

	// gauge is the prometheus gauge used to report this metric.
	vec *prometheus.GaugeVec

	// labeler is the label maker used to create labels for this metric.
	labeler *labelMaker
}

// newGaugeMetric creates a new GaugeMetric instance.
func newGaugeMetric(
	name string,
	unit string,
	description string,
	vec *prometheus.GaugeVec,
	labeler *labelMaker) GaugeMetric {

	return &gaugeMetric{
		name:        name,
		unit:        unit,
		description: description,
		vec:         vec,
		labeler:     labeler,
	}
}

func (m *gaugeMetric) Name() string {
	return m.name
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

func (m *gaugeMetric) Set(value float64, label ...any) error {
	if m.vec == nil {
		// metric is not enabled
		return nil
	}

	if len(label) > 1 {
		return fmt.Errorf("too many labels provided, expected 1, got %d", len(label))
	}

	var l any
	if len(label) == 1 {
		l = label[0]
	}

	values, err := m.labeler.extractValues(l)
	if err != nil {
		return fmt.Errorf("error extracting values from label for metric %s: %v", m.name, err)
	}

	observer := m.vec.WithLabelValues(values...)

	observer.Set(value)
	return nil
}
