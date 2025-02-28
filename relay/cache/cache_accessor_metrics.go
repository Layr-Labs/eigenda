package cache

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

const namespace = "eigenda_relay"

// CacheAccessorMetrics provides metrics for a CacheAccessor.
type CacheAccessorMetrics struct {
	cacheHits        *prometheus.CounterVec
	cacheNearMisses  *prometheus.CounterVec
	cacheMisses      *prometheus.CounterVec
	size             *prometheus.GaugeVec
	weight           *prometheus.GaugeVec
	averageWeight    *prometheus.GaugeVec
	cacheMissLatency *prometheus.SummaryVec
}

// NewCacheAccessorMetrics creates a new CacheAccessorMetrics.
func NewCacheAccessorMetrics(
	registry *prometheus.Registry,
	cacheName string) *CacheAccessorMetrics {

	cacheHits := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_hit_count", cacheName),
			Help:      "Number of cache hits",
		},
		[]string{},
	)

	cacheNearMisses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_near_miss_count", cacheName),
			Help:      "Number of near cache misses (i.e. a lookup is already in progress)",
		},
		[]string{},
	)

	cacheMisses := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_miss_count", cacheName),
			Help:      "Number of cache misses",
		},
		[]string{},
	)

	size := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_size", cacheName),
			Help:      "Number of items in the cache",
		},
		[]string{},
	)

	weight := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_weight", cacheName),
			Help:      "Total weight of items in the cache",
		},
		[]string{},
	)

	averageWeight := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_average_weight", cacheName),
			Help:      "Weight of each item currently in the cache",
		},
		[]string{},
	)

	cacheMissLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       fmt.Sprintf("%s_cache_miss_latency_ms", cacheName),
			Help:       "Latency of cache misses",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.05, 0.99: 0.01},
		},
		[]string{},
	)

	return &CacheAccessorMetrics{
		cacheHits:        cacheHits,
		cacheNearMisses:  cacheNearMisses,
		cacheMisses:      cacheMisses,
		size:             size,
		weight:           weight,
		averageWeight:    averageWeight,
		cacheMissLatency: cacheMissLatency,
	}
}

func (m *CacheAccessorMetrics) ReportCacheHit() {
	m.cacheHits.WithLabelValues().Inc()
}

func (m *CacheAccessorMetrics) ReportCacheNearMiss() {
	m.cacheNearMisses.WithLabelValues().Inc()
}

func (m *CacheAccessorMetrics) ReportCacheMiss() {
	m.cacheMisses.WithLabelValues().Inc()
}

func (m *CacheAccessorMetrics) ReportSize(size int) {
	m.size.WithLabelValues().Set(float64(size))
}

func (m *CacheAccessorMetrics) ReportWeight(weight uint64) {
	m.weight.WithLabelValues().Set(float64(weight))
}

func (m *CacheAccessorMetrics) ReportAverageWeight(averageWeight float64) {
	m.averageWeight.WithLabelValues().Set(averageWeight)
}

func (m *CacheAccessorMetrics) ReportCacheMissLatency(duration time.Duration) {
	m.cacheMissLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}
