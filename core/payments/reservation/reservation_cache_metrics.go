package reservation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Tracks metrics for the [ReservationLedgerCache]
type ReservationCacheMetrics struct {
	registry           *prometheus.Registry
	namespace          string
	subsystem          string
	cacheSize          prometheus.GaugeFunc
	evictions          prometheus.Counter
	prematureEvictions prometheus.Counter
	resizes            prometheus.Counter
	cacheMisses        prometheus.Counter
}

func NewReservationCacheMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
) *ReservationCacheMetrics {
	if registry == nil {
		return nil
	}

	evictions := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_ledger_cache_evictions",
			Subsystem: subsystem,
			Help:      "Total number of evictions from the reservation ledger cache",
		},
	)

	prematureEvictions := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_ledger_cache_premature_evictions",
			Subsystem: subsystem,
			Help:      "Total number of premature evictions (non-empty bucket) from the reservation ledger cache",
		},
	)

	resizes := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_ledger_cache_resizes",
			Subsystem: subsystem,
			Help:      "Total number of times the reservation ledger cache was resized",
		},
	)

	cacheMisses := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_ledger_cache_misses",
			Subsystem: subsystem,
			Help:      "Total number of cache misses in the reservation ledger cache",
		},
	)

	return &ReservationCacheMetrics{
		registry:           registry,
		namespace:          namespace,
		subsystem:          subsystem,
		evictions:          evictions,
		prematureEvictions: prematureEvictions,
		resizes:            resizes,
		cacheMisses:        cacheMisses,
	}
}

// Registers a gauge for cache size at runtime
//
// This should be called after the cache is initialized
func (m *ReservationCacheMetrics) RegisterSizeGauge(sizeGetter func() int) {
	if m == nil || m.registry == nil || m.cacheSize != nil {
		return
	}

	m.cacheSize = promauto.With(m.registry).NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Name:      "reservation_ledger_cache_size",
			Subsystem: m.subsystem,
			Help:      "Current number of entries in the reservation ledger cache",
		},
		func() float64 {
			return float64(sizeGetter())
		},
	)
}

// Increments the evictions counter
func (m *ReservationCacheMetrics) IncrementEvictions() {
	if m == nil {
		return
	}
	m.evictions.Inc()
}

// Increments the premature evictions counter
func (m *ReservationCacheMetrics) IncrementPrematureEvictions() {
	if m == nil {
		return
	}
	m.prematureEvictions.Inc()
}

// Increments the counter tracking number of cache resizes
func (m *ReservationCacheMetrics) IncrementResizes() {
	if m == nil {
		return
	}
	m.resizes.Inc()
}

// Increments the cache misses counter
func (m *ReservationCacheMetrics) IncrementCacheMisses() {
	if m == nil {
		return
	}
	m.cacheMisses.Inc()
}
