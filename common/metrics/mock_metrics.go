package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
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
	description string,
	templateLabel any,
	quantiles ...*Quantile) (LatencyMetric, error) {
	return &mockLatencyMetric{}, nil
}

func (m *mockMetrics) NewCountMetric(name string, description string, templateLabel any) (CountMetric, error) {
	return &mockCountMetric{}, nil
}

func (m *mockMetrics) NewGaugeMetric(
	name string,
	unit string,
	description string,
	labelTemplate any) (GaugeMetric, error) {
	return &mockGaugeMetric{}, nil
}

func (m *mockMetrics) NewAutoGauge(
	name string,
	unit string,
	description string,
	pollPeriod time.Duration,
	source func() float64,
	label ...any) error {
	return nil
}

func (m *mockMetrics) NewRunningAverageMetric(
	name string,
	unit string,
	description string,
	timeWindow time.Duration,
	labelTemplate any) (RunningAverageMetric, error) {
	return &mockRunningAverageMetric{}, nil
}

func (m *mockMetrics) RegisterExternalMetrics(collectors ...prometheus.Collector) {

}

var _ CountMetric = &mockCountMetric{}

type mockCountMetric struct {
}

func (m *mockCountMetric) Name() string {
	return ""
}

func (m *mockCountMetric) Unit() string {
	return ""
}

func (m *mockCountMetric) Description() string {
	return ""
}

func (m *mockCountMetric) Type() string {
	return ""
}

func (m *mockCountMetric) LabelFields() []string {
	return make([]string, 0)
}

func (m *mockCountMetric) Increment(label ...any) {

}

func (m *mockCountMetric) Add(value float64, label ...any) {

}

var _ GaugeMetric = &mockGaugeMetric{}

type mockGaugeMetric struct {
}

func (m *mockGaugeMetric) Name() string {
	return ""
}

func (m *mockGaugeMetric) Unit() string {
	return ""
}

func (m *mockGaugeMetric) Description() string {
	return ""
}

func (m *mockGaugeMetric) Type() string {
	return ""
}

func (m *mockGaugeMetric) LabelFields() []string {
	return make([]string, 0)
}

func (m *mockGaugeMetric) Set(value float64, label ...any) {

}

var _ LatencyMetric = &mockLatencyMetric{}

type mockLatencyMetric struct {
}

func (m *mockLatencyMetric) Name() string {
	return ""
}

func (m *mockLatencyMetric) Unit() string {
	return ""
}

func (m *mockLatencyMetric) Description() string {
	return ""
}

func (m *mockLatencyMetric) Type() string {
	return ""
}

func (m *mockLatencyMetric) LabelFields() []string {
	return make([]string, 0)
}

func (m *mockLatencyMetric) ReportLatency(latency time.Duration, label ...any) {

}

var _ RunningAverageMetric = &mockRunningAverageMetric{}

type mockRunningAverageMetric struct {
}

func (m *mockRunningAverageMetric) Name() string {
	return ""
}

func (m *mockRunningAverageMetric) Unit() string {
	return ""
}

func (m *mockRunningAverageMetric) Description() string {
	return ""
}

func (m *mockRunningAverageMetric) Type() string {
	return ""
}

func (m *mockRunningAverageMetric) LabelFields() []string {
	return make([]string, 0)
}

func (m *mockRunningAverageMetric) Update(value float64, label ...any) {

}

func (m *mockRunningAverageMetric) GetTimeWindow() time.Duration {
	return 0
}
