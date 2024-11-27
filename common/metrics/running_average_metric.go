package metrics

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
	"time"
)

var _ RunningAverageMetric = &runningAverageMetric{}

type runningAverageMetric struct {

	// logger is the logger used to log errors.
	logger logging.Logger

	// name is the name of the metric.
	name string

	// unit is the unit of the metric.
	unit string

	// description is the description of the metric.
	description string

	// vec is the prometheus gauge vector used to store the metric.
	vec *prometheus.GaugeVec

	// lm is the label maker used to create labels for this metric.
	labeler *labelMaker

	// runningAverage is the running average used to calculate the average of the metric.
	runningAverage *RunningAverage

	// timeWindow is the time window used to calculate the running average.
	timeWindow time.Duration

	// lock is used to provide thread safety for the running average calculator.
	lock sync.Mutex
}

// newRunningAverageMetric creates a new RunningAverageMetric instance.
func newRunningAverageMetric(
	logger logging.Logger,
	registry *prometheus.Registry,
	namespace string,
	name string,
	unit string,
	description string,
	timeWindow time.Duration,
	labelTemplate any) (RunningAverageMetric, error) {

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

	return &runningAverageMetric{
		logger:         logger,
		name:           name,
		unit:           unit,
		description:    description,
		vec:            vec,
		labeler:        labeler,
		runningAverage: NewRunningAverage(timeWindow),
		timeWindow:     timeWindow,
	}, nil
}

func (m *runningAverageMetric) Name() string {
	return m.name
}

func (m *runningAverageMetric) Unit() string {
	return m.unit
}

func (m *runningAverageMetric) Description() string {
	return m.description
}

func (m *runningAverageMetric) Type() string {
	return "running average"
}

func (m *runningAverageMetric) LabelFields() []string {
	return m.labeler.getKeys()
}

func (m *runningAverageMetric) Update(value float64, label ...any) {
	var l any
	if len(label) > 0 {
		l = label[0]
	}

	values, err := m.labeler.extractValues(l)
	if err != nil {
		m.logger.Errorf("error extracting values from label: %v", err)
	}

	m.lock.Lock()
	average := m.runningAverage.Update(time.Now(), value)
	m.lock.Unlock()
	m.vec.WithLabelValues(values...).Set(average)
}

func (m *runningAverageMetric) GetTimeWindow() time.Duration {
	return m.timeWindow
}
