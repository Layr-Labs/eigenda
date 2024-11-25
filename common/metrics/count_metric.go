package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	vec *prometheus.CounterVec

	// labeler is the label maker used to create labels for this metric.
	labeler *labelMaker
}

// newCountMetric creates a new CountMetric instance.
func newCountMetric(
	registry *prometheus.Registry,
	namespace string,
	name string,
	description string,
	labelTemplate any) (CountMetric, error) {

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_count", name),
		},
		labeler.getKeys(),
	)

	return &countMetric{
		name:        name,
		description: description,
		vec:         vec,
		labeler:     labeler,
	}, nil
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

func (m *countMetric) LabelFields() []string {
	return m.labeler.getKeys()
}

func (m *countMetric) Increment(label ...any) error {
	return m.Add(1, label...)
}

func (m *countMetric) Add(value float64, label ...any) error {
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
	observer.Add(value)

	return nil
}
