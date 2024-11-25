package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _ CountMetric = &countMetric{}

// countMetric a standard implementation of the CountMetric.
type countMetric struct {
	Metric

	// name is the name of the metric.
	name string

	// description is the description of the metric.
	description string

	// counter is the prometheus counter used to report this metric.
	counter prometheus.Counter
}

// newCountMetric creates a new CountMetric instance.
func newCountMetric(name string, description string, vec *prometheus.CounterVec) CountMetric {
	var counter prometheus.Counter
	if vec != nil {
		counter = vec.WithLabelValues() // TODO
	}

	return &countMetric{
		name:        name,
		description: description,
		counter:     counter,
	}
}

func (m *countMetric) Name() string {
	return m.name
}

func (m *countMetric) Unit() string {
	return "count"
}

func (m *countMetric) Description() string {
	return m.description
}

func (m *countMetric) Type() string {
	return "counter"
}

func (m *countMetric) Enabled() bool {
	return m.counter != nil
}

func (m *countMetric) Increment() {
	if m.counter == nil {
		return
	}
	m.counter.Inc()
}

func (m *countMetric) Add(value float64) {
	if m.counter == nil {
		return
	}
	m.counter.Add(value)
}
