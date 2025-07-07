package wrappers

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// LittDB has a bunch of metrics it exports. This little utility replicates some of the more important metrics.
// The names for these metrics mirror the names exported by LittDB, so that the metrics show up in the same dashboards.
type basicWrapperMetrics struct {
	tableName string

	// The number of bytes written to disk since startup. Only includes values, not metadata.
	bytesWrittenCounter *prometheus.CounterVec

	// The number of bytes read from disk since startup.
	bytesReadCounter *prometheus.CounterVec

	// Reports on the write latency of the database.
	writeLatency *prometheus.SummaryVec

	// Reports on the latency of a flush operation.
	flushLatency *prometheus.SummaryVec

	// Reports on the read latency of the database. This metric includes both cache hits and cache misses.
	readLatency *prometheus.SummaryVec
}

func newBasicWrapperMetrics(registry *prometheus.Registry, tableName string) *basicWrapperMetrics {
	if registry == nil {
		return nil
	}

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	bytesReadCounter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "litt",
			Name:      "bytes_read",
			Help:      "The number of bytes read from disk since startup.",
		},
		[]string{"table"},
	)

	readLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "litt",
			Name:      "read_latency_ms",
			Help: "Reports on the read latency of the database. " +
				"This metric includes both cache hits and cache misses.",
			Objectives: objectives,
		},
		[]string{"table"},
	)

	bytesWrittenCounter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "litt",
			Name:      "bytes_written",
			Help:      "The number of bytes written to disk since startup. Only includes values, not metadata.",
		},
		[]string{"table"},
	)

	writeLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "litt",
			Name:       "write_latency_ms",
			Help:       "Reports on the write latency of the database.",
			Objectives: objectives,
		},
		[]string{"table"},
	)

	flushLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "litt",
			Name:       "flush_latency_ms",
			Help:       "Reports on the latency of a flush operation.",
			Objectives: objectives,
		},
		[]string{"table"},
	)

	return &basicWrapperMetrics{
		bytesWrittenCounter: bytesWrittenCounter,
		bytesReadCounter:    bytesReadCounter,
		writeLatency:        writeLatency,
		flushLatency:        flushLatency,
		readLatency:         readLatency,
		tableName:           tableName,
	}
}

func (m *basicWrapperMetrics) RecordBytesWritten(bytes uint64) {
	//if m == nil {
	//	return
	//}
	m.bytesWrittenCounter.WithLabelValues(m.tableName).Add(float64(bytes))
}

func (m *basicWrapperMetrics) RecordBytesRead(bytes uint64) {
	//if m == nil {
	//	return
	//}
	m.bytesReadCounter.WithLabelValues(m.tableName).Add(float64(bytes))
}

func (m *basicWrapperMetrics) RecordWriteLatency(latency time.Duration) {
	//if m == nil {
	//	return
	//}
	m.writeLatency.WithLabelValues(m.tableName).Observe(common.ToMilliseconds(latency))
}

func (m *basicWrapperMetrics) RecordFlushLatency(latency time.Duration) {
	//if m == nil {
	//	return
	//}
	m.flushLatency.WithLabelValues(m.tableName).Observe(common.ToMilliseconds(latency))
}

func (m *basicWrapperMetrics) RecordReadLatency(latency time.Duration) {
	//if m == nil {
	//	return
	//}
	m.readLatency.WithLabelValues(m.tableName).Observe(common.ToMilliseconds(latency))
}
