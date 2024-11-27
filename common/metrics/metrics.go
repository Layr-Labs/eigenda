package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// Metrics provides a convenient interface for reporting metrics.
type Metrics interface {
	// Start starts the metrics server.
	Start() error

	// Stop stops the metrics server.
	Stop() error

	// GenerateMetricsDocumentation generates documentation for all currently registered metrics.
	// Documentation is returned as a string in markdown format.
	GenerateMetricsDocumentation() string

	// WriteMetricsDocumentation writes documentation for all currently registered metrics to a file.
	// Documentation is written in markdown format.
	WriteMetricsDocumentation(fileName string) error

	// NewLatencyMetric creates a new LatencyMetric instance. Useful for reporting the latency of an operation.
	// Metric name and label may only contain alphanumeric characters and underscores.
	//
	// The labelTemplate parameter is the label type that will be used for this metric. Each field becomes a label for
	// the metric. Each field type must be a string. If no labels are needed, pass nil.
	NewLatencyMetric(
		name string,
		description string,
		labelTemplate any,
		quantiles ...*Quantile) (LatencyMetric, error)

	// NewCountMetric creates a new CountMetric instance. Useful for tracking the count of a type of event.
	// Metric name and label may only contain alphanumeric characters and underscores.
	//
	// The labelTemplate parameter is the label type that will be used for this metric. Each field becomes a label for
	// the metric. Each field type must be a string. If no labels are needed, pass nil.
	NewCountMetric(
		name string,
		description string,
		labelTemplate any) (CountMetric, error)

	// NewGaugeMetric creates a new GaugeMetric instance. Useful for reporting specific values.
	// Metric name and label may only contain alphanumeric characters and underscores.
	//
	// The labelTemplate parameter is the label type that will be used for this metric. Each field becomes a label for
	// the metric. Each field type must be a string. If no labels are needed, pass nil.
	NewGaugeMetric(
		name string,
		unit string,
		description string,
		labelTemplate any) (GaugeMetric, error)

	// NewAutoGauge creates a new GaugeMetric instance that is automatically updated by the given source function.
	// The function is polled at the given period. This produces a gauge type metric internally.
	// Metric name and label may only contain alphanumeric characters and underscores.
	//
	// The label parameter accepts zero or one label.
	NewAutoGauge(
		name string,
		unit string,
		description string,
		pollPeriod time.Duration,
		source func() float64,
		label ...any) error

	// NewRunningAverageMetric creates a new GaugeMetric instance that keeps track of the average of a series of values
	// over a given time window. Each value within the window is given equal weight.
	NewRunningAverageMetric(
		name string,
		unit string,
		description string,
		timeWindow time.Duration,
		labelTemplate any) (RunningAverageMetric, error)

	// RegisterExternalMetrics registers prometheus collectors created outside the metrics framework.
	RegisterExternalMetrics(collectors ...prometheus.Collector)
}

// Metric represents a metric that can be reported.
type Metric interface {

	// Name returns the name of the metric.
	Name() string

	// Unit returns the unit of the metric.
	Unit() string

	// Description returns the description of the metric. Should be a one or two sentence human-readable description.
	Description() string

	// Type returns the type of the metric.
	Type() string

	// LabelFields returns the fields of the label template.
	LabelFields() []string
}

// GaugeMetric allows specific values to be reported.
type GaugeMetric interface {
	Metric

	// Set sets the value of a gauge metric.
	//
	// The label parameter accepts zero or one label. If the label type does not match the template label type provided
	// when creating the metric, an error will be returned.
	Set(value float64, label ...any)
}

// CountMetric allows the count of a type of event to be tracked.
type CountMetric interface {
	Metric

	// Increment increments the count by 1.
	//
	// The label parameter accepts zero or one label. If the label type does not match the template label type provided
	// when creating the metric, an error will be returned.
	Increment(label ...any)

	// Add increments the count by the given value.
	//
	// The label parameter accepts zero or one label. If the label type does not match the template label type provided
	// when creating the metric, an error will be returned.
	Add(value float64, label ...any)
}

// Quantile describes a quantile of a latency metric that should be reported. For a description of how
// to interpret a quantile, see the prometheus documentation
// https://github.com/prometheus/client_golang/blob/v1.20.5/prometheus/summary.go#L126
type Quantile struct {
	Quantile float64
	Error    float64
}

// NewQuantile creates a new Quantile instance. Error is set to 1% of the quantile.
func NewQuantile(quantile float64) *Quantile {
	return &Quantile{
		Quantile: quantile,
		Error:    quantile / 100.0,
	}
}

// LatencyMetric allows the latency of an operation to be tracked. Similar to a gauge metric, but specialized for time.
//
// The label parameter accepts zero or one label. If the label type does not match the template label type provided
// when creating the metric, an error will be returned.
type LatencyMetric interface {
	Metric

	// ReportLatency reports a latency value.
	//
	// The label parameter accepts zero or one label. If the label type does not match the template label type provided
	// when creating the metric, an error will be returned.
	ReportLatency(latency time.Duration, label ...any)
}

// RunningAverageMetric tracks the average of a series of values over a given time window.
type RunningAverageMetric interface {
	Metric

	// Update adds a new value to the RunningAverage.
	Update(value float64, label ...any)
}
