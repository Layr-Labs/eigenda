// This file should have been placed under disperser/dataapi/v2.
// The reason it's placed here in "dataapi" package is to avoid circular dependency
// (the "v2" already has dependency on "dataapi").
// Note the reason there is a "v2" package in the first place is to enable the separation of
// swagger docs for v1 and v2 APIs.

package dataapi

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "eigenda_dataapi"

type FeedCacheMetrics struct {
	// Time range metrics
	cacheTimeRangeSeconds      prometheus.Gauge
	cacheSegmentStartTimestamp prometheus.Gauge
	cacheSegmentEndTimestamp   prometheus.Gauge

	// Cache hit metrics
	cacheHitRatePercent prometheus.Gauge
}

func NewFeedCacheMetrics(name string, registry *prometheus.Registry) *FeedCacheMetrics {
	cacheTimeRangeSeconds := promauto.With(registry).NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_time_range_seconds", name),
			Help:      "Time range in seconds currently covered by the cache",
		},
	)

	cacheSegmentStartTimestamp := promauto.With(registry).NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_segment_start_timestamp_seconds", name),
			Help:      "Unix timestamp of the earliest item in the cache",
		},
	)

	cacheSegmentEndTimestamp := promauto.With(registry).NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_segment_end_timestamp_seconds", name),
			Help:      "Unix timestamp of the latest item in the cache",
		},
	)

	cacheHitRatePercent := promauto.With(registry).NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      fmt.Sprintf("%s_cache_hit_rate_percent", name),
			Help:      "Percentage of items served from cache vs total items requested",
		},
	)

	return &FeedCacheMetrics{
		cacheTimeRangeSeconds:      cacheTimeRangeSeconds,
		cacheSegmentStartTimestamp: cacheSegmentStartTimestamp,
		cacheSegmentEndTimestamp:   cacheSegmentEndTimestamp,
		cacheHitRatePercent:        cacheHitRatePercent,
	}
}

// UpdateHitRate updates the hit rate metric based on accumulated hits and misses.
func (m *FeedCacheMetrics) UpdateHitRate(hits, misses int) {
	total := hits + misses
	if total > 0 {
		hitRate := float64(hits) / float64(total) * 100.0
		m.cacheHitRatePercent.Set(hitRate)
	}
}

// RecordCacheUpdate updates metrics after a cache update operation.
func (m *FeedCacheMetrics) RecordCacheUpdate(
	cacheTimeStart time.Time,
	cacheTimeEnd time.Time,
) {
	if cacheTimeStart.IsZero() || cacheTimeEnd.IsZero() || !cacheTimeEnd.After(cacheTimeStart) {
		// Invalid time range, don't update metrics
		return
	}

	// Update cache time range metric
	cacheRangeSeconds := cacheTimeEnd.Sub(cacheTimeStart).Seconds()
	m.cacheTimeRangeSeconds.Set(cacheRangeSeconds)

	// Update cache segment timestamp gauges
	m.cacheSegmentStartTimestamp.Set(float64(cacheTimeStart.Unix()))
	m.cacheSegmentEndTimestamp.Set(float64(cacheTimeEnd.Unix()))
}
