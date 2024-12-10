package cache

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common/metrics"
)

// CacheAccessorMetrics provides metrics for a CacheAccessor.
type CacheAccessorMetrics struct {
	cacheHits        metrics.CountMetric
	cacheMisses      metrics.CountMetric
	size             metrics.GaugeMetric
	weight           metrics.GaugeMetric
	averageWeight    metrics.GaugeMetric
	cacheMissLatency metrics.LatencyMetric
}

// NewCacheAccessorMetrics creates a new CacheAccessorMetrics.
func NewCacheAccessorMetrics(
	server metrics.Metrics,
	cacheName string) (*CacheAccessorMetrics, error) {

	cacheHits, err := server.NewCountMetric(
		fmt.Sprintf("%s_cache_hit", cacheName),
		fmt.Sprintf("Number of cache hits in the %s cache", cacheName),
		nil)
	if err != nil {
		return nil, err
	}

	cacheMisses, err := server.NewCountMetric(
		fmt.Sprintf("%s_cache_miss", cacheName),
		fmt.Sprintf("Number of cache misses in the %s cache", cacheName),
		nil)
	if err != nil {
		return nil, err
	}

	size, err := server.NewGaugeMetric(
		fmt.Sprintf("%s_cache", cacheName),
		"size",
		fmt.Sprintf("Number of items in the %s cache", cacheName),
		nil)
	if err != nil {
		return nil, err
	}

	weight, err := server.NewGaugeMetric(
		fmt.Sprintf("%s_cache", cacheName),
		"weight",
		fmt.Sprintf("Total weight of items in the %s cache", cacheName),
		nil)
	if err != nil {
		return nil, err
	}

	averageWeight, err := server.NewGaugeMetric(
		fmt.Sprintf("%s_cache_item", cacheName),
		"weight",
		fmt.Sprintf("Weight of each item currently in the %s cache", cacheName),
		nil)
	if err != nil {
		return nil, err
	}

	cacheMissLatency, err := server.NewLatencyMetric(
		fmt.Sprintf("%s_cache_miss_latency", cacheName),
		fmt.Sprintf("Latency of cache misses in the %s cache", cacheName),
		nil,
		&metrics.Quantile{Quantile: 0.5, Error: 0.05},
		&metrics.Quantile{Quantile: 0.9, Error: 0.05},
		&metrics.Quantile{Quantile: 0.99, Error: 0.05})
	if err != nil {
		return nil, err
	}

	return &CacheAccessorMetrics{
		cacheHits:        cacheHits,
		cacheMisses:      cacheMisses,
		size:             size,
		weight:           weight,
		averageWeight:    averageWeight,
		cacheMissLatency: cacheMissLatency,
	}, nil
}
