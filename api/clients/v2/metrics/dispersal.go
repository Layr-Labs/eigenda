package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	dispersalSubsystem = "dispersal"
)

type DispersalMetricer interface {
	RecordBlobSizeBytes(size uint)
}

type DispersalMetrics struct {
	BlobSize *prometheus.HistogramVec
}

func NewDispersalMetrics(registry *prometheus.Registry) DispersalMetricer {
	if registry == nil {
		return NoopDispersalMetrics
	}

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
		BlobSize: promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
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

type noopDispersalMetricer struct {
}

var NoopDispersalMetrics DispersalMetricer = new(noopDispersalMetricer)

func (n *noopDispersalMetricer) RecordBlobSizeBytes(_ uint) {
}
