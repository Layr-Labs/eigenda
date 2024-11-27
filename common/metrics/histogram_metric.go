package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ HistogramMetric = &histogramMetric{}

type histogramMetric struct {

	// logger is the logger used to log errors.
	logger logging.Logger

	// name is the name of the metric.
	name string

	// unit is the unit of the metric.
	unit string

	// description is the description of the metric.
	description string

	// vec is the prometheus histogram vector used to report this metric.
	vec *prometheus.HistogramVec

	// lm is the label maker used to create labels for this metric.
	labeler *labelMaker
}

// newHistogramMetric creates a new HistogramMetric instance.
func newHistogramMetric(
	logger logging.Logger,
	registry *prometheus.Registry,
	namespace string,
	name string,
	unit string,
	description string,
	bucketFactor float64,
	labelTemplate any) (HistogramMetric, error) {

	labeler, err := newLabelMaker(labelTemplate)
	if err != nil {
		return nil, err
	}

	vec := promauto.With(registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:                   namespace,
			Name:                        fmt.Sprintf("%s_%s", name, unit),
			Help:                        description,
			NativeHistogramBucketFactor: bucketFactor,
		},
		labeler.getKeys(),
	)

	return &histogramMetric{
		logger:      logger,
		name:        name,
		unit:        unit,
		description: description,
		vec:         vec,
		labeler:     labeler,
	}, nil
}

func (m *histogramMetric) Name() string {
	return m.name
}

func (m *histogramMetric) Unit() string {
	return m.unit
}

func (m *histogramMetric) Description() string {
	return m.description
}

func (m *histogramMetric) Type() string {
	return "histogram"
}

func (m *histogramMetric) LabelFields() []string {
	return m.labeler.getKeys()
}

func (m *histogramMetric) Observe(value float64, label ...any) {
	var l any
	if len(label) > 0 {
		l = label[0]
	}

	values, err := m.labeler.extractValues(l)
	if err != nil {
		m.logger.Errorf("error extracting values from label: %v", err)
	}

	m.vec.WithLabelValues(values...).Observe(value)
}
