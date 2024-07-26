package metrics

import (
	"sync"
	"time"
)

// MockMetrics implements metrics, useful for testing.
type MockMetrics struct {
	// A map from each count metric's description to its count.
	counts map[string]float64
	// A map from each gauge metric's description to its value.
	gauges map[string]float64
	// A map from each latency metric's description to its most recently reported latency.
	latencies map[string]time.Duration

	// Used to ensure thread safety.
	lock sync.Mutex
}

// NewMockMetrics creates a new MockMetrics instance.
func NewMockMetrics() *MockMetrics {
	return &MockMetrics{
		counts:    make(map[string]float64),
		gauges:    make(map[string]float64),
		latencies: make(map[string]time.Duration),
	}
}

// GetCount returns the count of a type of event.
func (m *MockMetrics) GetCount(description string) float64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.counts[description]
}

// GetGaugeValue returns the value of a gauge metric.
func (m *MockMetrics) GetGaugeValue(description string) float64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.gauges[description]
}

// GetLatency returns the most recently reported latency of an operation.
func (m *MockMetrics) GetLatency(description string) time.Duration {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.latencies[description]
}

func (m *MockMetrics) Start() {
	// intentional no-op
}

func (m *MockMetrics) NewLatencyMetric(description string) LatencyMetric {
	return &mockLatencyMetric{
		metrics:     m,
		description: description,
	}
}

func (m *MockMetrics) NewCountMetric(description string) CountMetric {
	return &mockCountMetric{
		metrics:     m,
		description: description,
	}
}

func (m *MockMetrics) NewGaugeMetric(description string) GaugeMetric {
	return &mockGaugeMetric{
		metrics:     m,
		description: description,
	}
}

// mockLatencyMetric implements LatencyMetric, useful for testing.
type mockLatencyMetric struct {
	metrics     *MockMetrics
	description string
}

func (m *mockLatencyMetric) ReportLatency(latency time.Duration) {
	m.metrics.lock.Lock()
	m.metrics.latencies[m.description] = latency
	m.metrics.lock.Unlock()
}

// mockCountMetric implements CountMetric, useful for testing.
type mockCountMetric struct {
	metrics     *MockMetrics
	description string
}

func (m *mockCountMetric) Increment() {
	m.metrics.lock.Lock()
	m.metrics.counts[m.description]++
	m.metrics.lock.Unlock()
}

// mockGaugeMetric implements GaugeMetric, useful for testing.
type mockGaugeMetric struct {
	metrics     *MockMetrics
	description string
}

func (m *mockGaugeMetric) Set(value float64) {
	m.metrics.lock.Lock()
	m.metrics.gauges[m.description] = value
	m.metrics.lock.Unlock()
}
