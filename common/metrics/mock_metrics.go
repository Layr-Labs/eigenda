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

func (m mockMetrics) Start() error {
	return nil
}

func (m mockMetrics) Stop() error {
	return nil
}

func (m mockMetrics) NewLatencyMetric(name string, label string, quantiles ...*Quantile) (LatencyMetric, error) {
	return newLatencyMetric(name, label, nil), nil
}

func (m mockMetrics) NewCountMetric(name string, label string) (CountMetric, error) {
	return newCountMetric(name, label, nil), nil
}

func (m mockMetrics) NewGaugeMetric(name string, label string) (GaugeMetric, error) {
	return newGaugeMetric(name, label, nil), nil
}

func (m mockMetrics) NewAutoGauge(name string, label string, pollPeriod time.Duration, source func() float64) error {
	return nil
}
