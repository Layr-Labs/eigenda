package metrics

import (
	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	retrievalSubsystem = "retrieval"
)

type RetrievalMetricer interface {
	RecordPayloadSizeBytes(size int)

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

	return &RetrievalMetrics{
		PayloadSize: factory.NewHistogram(prometheus.HistogramOpts{
			Name:      "payload_size_bytes",
			Namespace: namespace,
			Subsystem: retrievalSubsystem,
			Help:      "Size of decoded payloads in bytes",
			Buckets:   blobSizeBuckets,
		}),
		factory: factory,
	}
}

func (m *RetrievalMetrics) RecordPayloadSizeBytes(size int) {
	m.PayloadSize.Observe(float64(size))
}

func (m *RetrievalMetrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopRetrievalMetricer struct {
}

var NoopRetrievalMetrics RetrievalMetricer = new(noopRetrievalMetricer)

func (n *noopRetrievalMetricer) RecordPayloadSizeBytes(_ int) {
}

func (n *noopRetrievalMetricer) Document() []metrics.DocumentedMetric {
	return []metrics.DocumentedMetric{}
}
