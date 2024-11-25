package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var _ LatencyMetric = &latencyMetric{}

// latencyMetric is a standard implementation of the LatencyMetric interface via prometheus.
type latencyMetric struct {
	Metric

	// name is the name of the metric.
	name string

	// description is the description of the metric.
	description string

	// vec is the prometheus summary vector used to report this metric.
	vec *prometheus.SummaryVec

	// lm is the label maker used to create labels for this metric.
	labeler *labelMaker
}

// newLatencyMetric creates a new LatencyMetric instance.
func newLatencyMetric(
	name string,
	description string,
	vec *prometheus.SummaryVec,
	labeler *labelMaker) LatencyMetric {

	return &latencyMetric{
		name:        name,
		description: description,
		vec:         vec,
		labeler:     labeler,
	}
}

func (m *latencyMetric) Name() string {
	return m.name
}

func (m *latencyMetric) Unit() string {
	return "seconds"
}

func (m *latencyMetric) Description() string {
	return m.description
}

func (m *latencyMetric) Type() string {
	return "latency"
}

func (m *latencyMetric) ReportLatency(latency time.Duration, label ...*struct{}) error {
	if m.vec == nil {
		// metric is not enabled
		return nil
	}

	if len(label) > 1 {
		return fmt.Errorf("too many labels provided, expected 1, got %d", len(label))
	}

	var l *struct{}
	if len(label) == 1 {
		l = label[0]
	}

	values, err := m.labeler.extractValues(l)
	if err != nil {
		return fmt.Errorf("error extracting values from label for metric %s: %v", m.name, err)
	}

	observer := m.vec.WithLabelValues(values...)

	nanoseconds := float64(latency.Nanoseconds())
	milliseconds := nanoseconds / 1e6
	observer.Observe(milliseconds)

	return nil
}
