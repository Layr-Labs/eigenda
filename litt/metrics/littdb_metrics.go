package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TODO
//  - total disk used
//  - total disk used, broken down by root
//  - total disk used, broken down by table
//  - disk available on each root
//  - disk available in total
//  - number of segments in each table
//  - number of keys in each table
//  - read latency
//  - write latency
//  - flush latency
//  - keymap flush latency
//  - read throughput, bytes per second
//  - read throughput, keys per second
//  - write throughput, bytes per second
//  - write throughput, keys per second
//  - control loop idle fraction
//    - main control loop
//    - flush loop
//    - shard control loops
//    - keyfile control loop
//  - average segment span
//  - cache size, entry count
//  - cache size, byte count
//  - average time spent in cache
//  - cache hit rate
//  - cache miss rate
//  - gc latency

// LittDBMetrics encapsulates metrics for a LittDB.
type LittDBMetrics struct {
	// The total disk used by the database.
	sizeInBytes *prometheus.GaugeVec
}

// NewLittDBMetrics creates a new LittDBMetrics instance.
func NewLittDBMetrics(registry *prometheus.Registry, namespace string) *LittDBMetrics {

	sizeInBytes := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "size_bytes",
			Help:      "The total size of the database in bytes.",
		},
		[]string{},
	)

	return &LittDBMetrics{
		sizeInBytes: sizeInBytes,
	}
}

// ReportSizeInBytes reports the total disk used by the database.
func (m *LittDBMetrics) ReportSizeInBytes(sizeInBytes uint64) {
	if m == nil {
		return
	}
	m.sizeInBytes.WithLabelValues().Set(float64(sizeInBytes))
}
