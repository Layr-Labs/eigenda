package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var _ LatencyMetric = &latencyMetric{}

// latencyMetric is a standard implementation of the LatencyMetric interface via prometheus.
type latencyMetric struct {
	Metric

	// name is the name of the metric.
	name string

	// label is the label of the metric.
	label string

	// observer is the prometheus observer used to report this metric.
	observer prometheus.Observer
}

// newLatencyMetric creates a new LatencyMetric instance.
func newLatencyMetric(name string, label string, vec *prometheus.SummaryVec) LatencyMetric {
	var observer prometheus.Observer
	if vec != nil {
		observer = vec.WithLabelValues(label)
	}

	return &latencyMetric{
		name:     name,
		label:    label,
		observer: observer,
	}
}

func (m *latencyMetric) Name() string {
	return m.name
}

func (m *latencyMetric) Label() string {
	return m.label
}

func (m *latencyMetric) Enabled() bool {
	return m.observer != nil
}

func (m *latencyMetric) ReportLatency(latency time.Duration) {
	if m.observer == nil {
		// this metric has been disabled
		return
	}
	m.observer.Observe(latency.Seconds())
}
