package metrics

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	dispersalSubsystem = "dispersal"
)

type DispersalMetricer interface {
	RecordBlobSizeBytes(size int)
	RecordDisperserReputationScore(disperserID uint32, score float64)

	Document() []metrics.DocumentedMetric
}

type DispersalMetrics struct {
	BlobSize                 prometheus.Histogram
	DisperserReputationScore *prometheus.GaugeVec

	factory *metrics.Documentor
}

func NewDispersalMetrics(registry *prometheus.Registry) DispersalMetricer {
	if registry == nil {
		return NoopDispersalMetrics
	}

	factory := metrics.With(registry)

	return &DispersalMetrics{
		BlobSize: factory.NewHistogram(prometheus.HistogramOpts{
			Name:      "blob_size_bytes",
			Namespace: namespace,
			Subsystem: dispersalSubsystem,
			Help:      "Size of blobs created from payloads in bytes",
			Buckets:   blobSizeBuckets,
		}),
		DisperserReputationScore: factory.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "disperser_reputation_score",
			Namespace: namespace,
			Subsystem: dispersalSubsystem,
			Help:      "Current reputation score for each disperser",
		}, []string{"disperser_id"}),
		factory: factory,
	}
}

func (m *DispersalMetrics) RecordBlobSizeBytes(size int) {
	m.BlobSize.Observe(float64(size))
}

func (m *DispersalMetrics) RecordDisperserReputationScore(disperserID uint32, score float64) {
	m.DisperserReputationScore.WithLabelValues(fmt.Sprintf("%d", disperserID)).Set(score)
}

func (m *DispersalMetrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopDispersalMetricer struct {
}

var NoopDispersalMetrics DispersalMetricer = new(noopDispersalMetricer)

func (n *noopDispersalMetricer) RecordBlobSizeBytes(_ int) {
}

func (n *noopDispersalMetricer) RecordDisperserReputationScore(_ uint32, _ float64) {
}

func (n *noopDispersalMetricer) Document() []metrics.DocumentedMetric {
	return []metrics.DocumentedMetric{}
}
