package metrics

import (
	"github.com/Layr-Labs/eigenda/litt"
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

	// The size of individual tables in the database.
	tableSizeInBytes *prometheus.GaugeVec

	// The number of keys in the database.
	keyCount *prometheus.GaugeVec

	// The number of keys in individual tables in the database.
	tableKeyCount *prometheus.GaugeVec
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

	tableSizeInBytes := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "table_size_bytes",
			Help:      "The size of individual tables in the database in bytes.",
		},
		[]string{"table"},
	)

	keyCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "key_count",
			Help:      "The total number of keys in the database.",
		},
		[]string{},
	)

	tableKeyCount := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "table_key_count",
			Help:      "The number of keys in individual tables in the database.",
		},
		[]string{"table"},
	)

	return &LittDBMetrics{
		sizeInBytes:      sizeInBytes,
		tableSizeInBytes: tableSizeInBytes,
		keyCount:         keyCount,
		tableKeyCount:    tableKeyCount,
	}
}

// CollectPeriodicMetrics is a method that is periodically called to collect metrics. Tables are not permitted to be
// added or dropped while this method is running.
func (m *LittDBMetrics) CollectPeriodicMetrics(db litt.DB, tables map[string]litt.ManagedTable) {
	if m == nil {
		return
	}

	totalSize := uint64(0)
	totalKeyCount := uint64(0)

	for _, table := range tables {
		tableName := table.Name()

		tableSize := table.Size()
		totalSize += tableSize
		m.tableSizeInBytes.WithLabelValues(tableName).Set(float64(tableSize))

		tableKeyCount := table.KeyCount()
		totalKeyCount += tableKeyCount
		m.tableKeyCount.WithLabelValues(tableName).Set(float64(tableKeyCount))
	}

	m.sizeInBytes.WithLabelValues().Set(float64(totalSize))
	m.keyCount.WithLabelValues().Set(float64(totalKeyCount))

}
