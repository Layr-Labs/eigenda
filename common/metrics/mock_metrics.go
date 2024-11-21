package metrics

import (
	"time"
)

var _ Metrics = &mockMetrics{}

// mockMetrics is a mock implementation of the Metrics interface.
type mockMetrics struct {
}

// NewMockMetrics creates a new mock Metrics instance.
// Suitable for testing or for when you just want to disable all metrics.
func NewMockMetrics() Metrics {
	return &mockMetrics{}
}

func (m *mockMetrics) GenerateMetricsDocumentation() string {
	return ""
}

func (m *mockMetrics) WriteMetricsDocumentation(fileName string) error {
	return nil
}

func (m *mockMetrics) Start() error {
	return nil
}

func (m *mockMetrics) Stop() error {
	return nil
}

func (m *mockMetrics) NewLatencyMetric(
	name string,
	label string,
	description string,
	quantiles ...*Quantile) (LatencyMetric, error) {
	return newLatencyMetric(name, label, description, nil), nil
}

func (m *mockMetrics) NewCountMetric(name string, label string, description string) (CountMetric, error) {
	return newCountMetric(name, label, description, nil), nil
}

func (m *mockMetrics) NewGaugeMetric(name string, label string, unit string, description string) (GaugeMetric, error) {
	return newGaugeMetric(name, label, unit, description, nil), nil
}

func (m *mockMetrics) NewAutoGauge(
	name string,
	label string,
	unit string,
	description string,
	pollPeriod time.Duration,
	source func() float64) error {
	return nil
}
