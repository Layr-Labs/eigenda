package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	retrievalSubsystem = "retrieval"
)

type RetrievalMetricer interface {
	RecordPayloadSizeBytes(size uint)

	Document() []metrics.DocumentedMetric
}

type RetrievalMetrics struct {
	PayloadSize prometheus.Histogram

	factory *metrics.Documentor
}

func NewRetrievalMetrics(registry *prometheus.Registry) RetrievalMetricer {
	if registry == nil {
		return NoopRetrievalMetrics
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

	return &RetrievalMetrics{
		PayloadSize: factory.NewHistogram(prometheus.HistogramOpts{
			Name:      "payload_size_bytes",
			Namespace: namespace,
			Subsystem: retrievalSubsystem,
			Help:      "Size of decoded payloads in bytes",
			Buckets:   sizeBuckets,
		}),
		factory: factory,
	}
}

func (m *RetrievalMetrics) RecordPayloadSizeBytes(size uint) {
	m.PayloadSize.Observe(float64(size))
}

func (m *RetrievalMetrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopRetrievalMetricer struct {
}

var NoopRetrievalMetrics RetrievalMetricer = new(noopRetrievalMetricer)

func (n *noopRetrievalMetricer) RecordPayloadSizeBytes(_ uint) {
}

func (n *noopRetrievalMetricer) Document() []metrics.DocumentedMetric {
	return []metrics.DocumentedMetric{}
}
