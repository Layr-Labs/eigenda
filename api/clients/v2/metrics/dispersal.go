package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	dispersalSubsystem = "dispersal"
)

type DispersalMetricer interface {
	RecordBlobSize(size uint)
	RecordSymbolLength(length uint)

	Document() []DocumentedMetric
}

type DispersalMetrics struct {
	BlobSize     *prometheus.HistogramVec
	SymbolLength *prometheus.HistogramVec

	factory *Documentor
}

func NewDispersalMetrics(registry *prometheus.Registry) DispersalMetricer {
	if registry == nil {
		return NoopDispersalMetrics
	}

	factory := With(registry)
	// Define size buckets for payload and blob size measurements
	// Starting from 1KB up to 16MB with exponential growth
	sizeBuckets := []float64{
		1024,     // 1KB
		4096,     // 4KB
		16384,    // 16KB
		65536,    // 64KB
		262144,   // 256KB
		1048576,  // 1MB
		4194304,  // 4MB
		16777216, // 16MB
	}

	return &DispersalMetrics{
		BlobSize: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "blob_size_bytes",
			Namespace: namespace,
			Subsystem: dispersalSubsystem,
			Help:      "Size of blobs created from payloads in bytes",
			Buckets:   sizeBuckets,
		}, []string{}),
		SymbolLength: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "blob_size_symbols",
			Namespace: namespace,
			Subsystem: dispersalSubsystem,
			Help:      "Size of blobs created from payloads in symbols",
			Buckets:   sizeBuckets,
		}, []string{}),
		factory: factory,
	}
}

func (m *DispersalMetrics) RecordBlobSize(size uint) {
	m.BlobSize.WithLabelValues().Observe(float64(size))
}

func (m *DispersalMetrics) RecordSymbolLength(length uint) {
	m.SymbolLength.WithLabelValues().Observe(float64(length))
}

func (m *DispersalMetrics) Document() []DocumentedMetric {
	return m.factory.Document()
}

type noopDispersalMetricer struct {
}

var NoopDispersalMetrics DispersalMetricer = new(noopDispersalMetricer)

func (n *noopDispersalMetricer) RecordBlobSize(_ uint) {
}

func (n *noopDispersalMetricer) RecordSymbolLength(_ uint) {
}

func (n *noopDispersalMetricer) Document() []DocumentedMetric {
	return []DocumentedMetric{}
}
