package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	dispersalSubsystem = "dispersal"
)

type DispersalMetricer interface {
	RecordBlobSizeBytes(size uint)

	Document() []metrics.DocumentedMetric
}

type DispersalMetrics struct {
	BlobSize     *prometheus.HistogramVec
  
  factory *metrics.Documentor
}

func NewDispersalMetrics(registry *prometheus.Registry) DispersalMetricer {
	if registry == nil {
		return NoopDispersalMetrics
	}

	factory := metrics.With(registry)
	// Define size buckets for payload and blob size measurements
	// Starting from 0 up to 16MiB
	sizeBuckets := []float64{
		0,
		131072,   // 128KiB
		262144,   // 256KiB
		524288,   // 512KiB
		1048576,  // 1MiB
		2097152,  // 2MiB
		4194304,  // 4MiB
		8388608,  // 8MiB
		16777216, // 16MiB
	}

	return &DispersalMetrics{
		BlobSize: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "blob_size_bytes",
			Namespace: namespace,
			Subsystem: dispersalSubsystem,
			Help:      "Size of blobs created from payloads in bytes",
			Buckets:   sizeBuckets,
		}, []string{}),
	}
}

func (m *DispersalMetrics) RecordBlobSizeBytes(size uint) {
	m.BlobSize.WithLabelValues().Observe(float64(size))
}

func (m *DispersalMetrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopDispersalMetricer struct {
}

var NoopDispersalMetrics DispersalMetricer = new(noopDispersalMetricer)

func (n *noopDispersalMetricer) RecordBlobSizeBytes(_ uint) {
}

func (n *noopDispersalMetricer) Document() []metrics.DocumentedMetric {
	return []metrics.DocumentedMetric{}
}
