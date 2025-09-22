package ondemand

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Tracks metrics for the [OnDemandLedgerCache]
type OnDemandCacheMetrics struct {
	registry    *prometheus.Registry
	namespace   string
	subsystem   string
	cacheSize   prometheus.GaugeFunc
	evictions   prometheus.Counter
	cacheMisses prometheus.Counter
}

func NewOnDemandCacheMetrics(registry *prometheus.Registry, namespace string, subsystem string) *OnDemandCacheMetrics {
	if registry == nil {
		return nil
	}

	evictions := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_ledger_cache_evictions",
			Subsystem: subsystem,
			Help:      "Total number of evictions from the on-demand ledger cache",
		},
	)

	cacheMisses := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_ledger_cache_misses",
			Subsystem: subsystem,
			Help:      "Total number of cache misses in the on-demand ledger cache",
		},
	)

	registry.MustRegister(evictions, cacheMisses)

	return &OnDemandCacheMetrics{
		registry:    registry,
		namespace:   namespace,
		subsystem:   subsystem,
		evictions:   evictions,
		cacheMisses: cacheMisses,
	}
}

// Registers a gauge for cache size at runtime
//
// This should be called after the cache is initialized
func (m *OnDemandCacheMetrics) RegisterSizeGauge(sizeGetter func() int) {
	if m == nil || m.registry == nil || m.cacheSize != nil {
		return
	}

	m.cacheSize = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Name:      "on_demand_ledger_cache_size",
			Subsystem: m.subsystem,
			Help:      "Current number of entries in the on-demand ledger cache",
		},
		func() float64 {
			return float64(sizeGetter())
		},
	)

	m.registry.MustRegister(m.cacheSize)
}

// Increments the evictions counter
func (m *OnDemandCacheMetrics) IncrementEvictions() {
	if m == nil {
		return
	}
	m.evictions.Inc()
}

// Increments the cache misses counter
func (m *OnDemandCacheMetrics) IncrementCacheMisses() {
	if m == nil {
		return
	}
	m.cacheMisses.Inc()
}
