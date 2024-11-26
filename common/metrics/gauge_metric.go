package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ GaugeMetric = &gaugeMetric{}

// gaugeMetric is a standard implementation of the GaugeMetric interface via prometheus.
type gaugeMetric struct {
	Metric

	// logger is the logger used to log errors.
	logger logging.Logger

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
	logger logging.Logger,
	registry *prometheus.Registry,
	namespace string,
	name string,
	unit string,
	description string,
	labelTemplate any) (GaugeMetric, error) {

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_%s", name, unit),
		},
		labeler.getKeys(),
	)

	return &gaugeMetric{
		logger:      logger,
		name:        name,
		unit:        unit,
		description: description,
		vec:         vec,
		labeler:     labeler,
	}, nil
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

func (m *gaugeMetric) LabelFields() []string {
	return m.labeler.getKeys()
}

func (m *gaugeMetric) Set(value float64, label ...any) {
	var l any
	if len(label) > 0 {
		l = label[0]
	}

	values, err := m.labeler.extractValues(l)
	if err != nil {
		m.logger.Errorf("failed to extract values from label: %v", err)
		return
	}

	observer := m.vec.WithLabelValues(values...)

	observer.Set(value)
}
