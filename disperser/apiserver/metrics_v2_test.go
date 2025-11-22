package apiserver

import (
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestTimestampDriftMetrics(t *testing.T) {
	logger := test.GetLogger()
	registry := prometheus.NewRegistry()
	metricsConfig := disperser.MetricsConfig{
		HTTPPort:      "9100",
		EnableMetrics: true,
	}

	metrics := newAPIServerV2Metrics(registry, metricsConfig, logger)

	// Test reporting timestamp drift for accepted requests from account1
	account1 := "0x1234567890123456789012345678901234567890"
	metrics.reportDispersalTimestampDrift(1.5, "accepted", account1)
	metrics.reportDispersalTimestampDrift(2.3, "accepted", account1)
	metrics.reportDispersalTimestampDrift(-0.5, "accepted", account1) // future timestamp

	// Test reporting timestamp drift for rejected requests from account2
	account2 := "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
	metrics.reportDispersalTimestampDrift(65.0, "rejected", account2)  // too old
	metrics.reportDispersalTimestampDrift(-50.0, "rejected", account2) // too far in future

	// Gather metrics
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	// Find the drift histogram
	var driftMetric *dto.MetricFamily
	for _, mf := range metricFamilies {
		if mf.GetName() == "eigenda_disperser_api_dispersal_timestamp_drift_seconds" {
			driftMetric = mf
			break
		}
	}

	require.NotNil(t, driftMetric, "drift metric should exist")
	require.Equal(t, dto.MetricType_HISTOGRAM, driftMetric.GetType())

	// Verify we have histograms for both accepted and rejected (per account)
	require.Len(t, driftMetric.GetMetric(), 2, "should have metrics for 2 account+status combinations")

	// Verify counts and account labels
	for _, m := range driftMetric.GetMetric() {
		status := ""
		accountID := ""
		for _, label := range m.GetLabel() {
			if label.GetName() == "status" {
				status = label.GetValue()
			}
			if label.GetName() == "account_id" {
				accountID = label.GetValue()
			}
		}

		histogram := m.GetHistogram()
		if status == "accepted" {
			require.Equal(t, account1, accountID, "accepted metrics should be from account1")
			require.Equal(t, uint64(3), histogram.GetSampleCount(), "accepted should have 3 samples")
		} else if status == "rejected" {
			require.Equal(t, account2, accountID, "rejected metrics should be from account2")
			require.Equal(t, uint64(2), histogram.GetSampleCount(), "rejected should have 2 samples")
		}
	}
}

func TestTimestampConfigMetrics(t *testing.T) {
	logger := test.GetLogger()
	registry := prometheus.NewRegistry()
	metricsConfig := disperser.MetricsConfig{
		HTTPPort:      "9100",
		EnableMetrics: true,
	}

	metrics := newAPIServerV2Metrics(registry, metricsConfig, logger)

	// Set config values
	maxAge := 45.0   // 45 seconds
	maxFuture := 5.0 // 5 seconds
	metrics.setDispersalTimestampConfig(maxAge, maxFuture)

	// Gather metrics
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	// Find the config gauges
	var maxAgeMetric, maxFutureMetric *dto.MetricFamily
	for _, mf := range metricFamilies {
		if mf.GetName() == "eigenda_disperser_api_dispersal_timestamp_max_age_seconds" {
			maxAgeMetric = mf
		}
		if mf.GetName() == "eigenda_disperser_api_dispersal_timestamp_max_future_seconds" {
			maxFutureMetric = mf
		}
	}

	require.NotNil(t, maxAgeMetric, "max age metric should exist")
	require.NotNil(t, maxFutureMetric, "max future metric should exist")

	// Verify values
	require.Equal(t, maxAge, maxAgeMetric.GetMetric()[0].GetGauge().GetValue())
	require.Equal(t, maxFuture, maxFutureMetric.GetMetric()[0].GetGauge().GetValue())
}

func TestTimestampRejectionMetrics(t *testing.T) {
	logger := test.GetLogger()
	registry := prometheus.NewRegistry()
	metricsConfig := disperser.MetricsConfig{
		HTTPPort:      "9100",
		EnableMetrics: true,
	}

	metrics := newAPIServerV2Metrics(registry, metricsConfig, logger)

	// Report rejections
	metrics.reportDispersalTimestampRejected("stale")
	metrics.reportDispersalTimestampRejected("stale")
	metrics.reportDispersalTimestampRejected("future")

	// Gather metrics
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	// Find the rejection counter
	var rejectionMetric *dto.MetricFamily
	for _, mf := range metricFamilies {
		if mf.GetName() == "eigenda_disperser_api_dispersal_timestamp_rejections_total" {
			rejectionMetric = mf
			break
		}
	}

	require.NotNil(t, rejectionMetric, "rejection metric should exist")
	require.Equal(t, dto.MetricType_COUNTER, rejectionMetric.GetType())

	// Verify counts by reason
	for _, m := range rejectionMetric.GetMetric() {
		reason := ""
		for _, label := range m.GetLabel() {
			if label.GetName() == "reason" {
				reason = label.GetValue()
			}
		}

		if reason == "stale" {
			require.Equal(t, 2.0, m.GetCounter().GetValue())
		} else if reason == "future" {
			require.Equal(t, 1.0, m.GetCounter().GetValue())
		}
	}
}

func TestTimestampValidationWithMetrics(t *testing.T) {
	// This test simulates the full timestamp validation flow with metrics
	logger := test.GetLogger()
	registry := prometheus.NewRegistry()
	metricsConfig := disperser.MetricsConfig{
		HTTPPort:      "9100",
		EnableMetrics: true,
	}

	metrics := newAPIServerV2Metrics(registry, metricsConfig, logger)

	// Set config
	maxAgeSeconds := 45.0
	maxFutureSeconds := 5.0
	metrics.setDispersalTimestampConfig(maxAgeSeconds, maxFutureSeconds)

	// Simulate various timestamp scenarios
	now := time.Now()
	account := "0x1234567890123456789012345678901234567890"

	// Scenario 1: Valid timestamp (2 seconds old)
	drift1 := 2.0
	metrics.reportDispersalTimestampDrift(drift1, "accepted", account)

	// Scenario 2: Too old (60 seconds)
	drift2 := 60.0
	metrics.reportDispersalTimestampDrift(drift2, "rejected", account)
	metrics.reportDispersalTimestampRejected("stale")

	// Scenario 3: Too far in future (-10 seconds)
	drift3 := -10.0
	metrics.reportDispersalTimestampDrift(drift3, "rejected", account)
	metrics.reportDispersalTimestampRejected("future")

	// Verify all metrics were recorded
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	var foundDrift, foundConfig, foundRejection bool
	for _, mf := range metricFamilies {
		switch mf.GetName() {
		case "eigenda_disperser_api_dispersal_timestamp_drift_seconds":
			foundDrift = true
		case "eigenda_disperser_api_dispersal_timestamp_max_age_seconds":
			foundConfig = true
		case "eigenda_disperser_api_dispersal_timestamp_rejections_total":
			foundRejection = true
		}
	}

	require.True(t, foundDrift, "drift metric should be present")
	require.True(t, foundConfig, "config metric should be present")
	require.True(t, foundRejection, "rejection metric should be present")

	_ = now // suppress unused variable warning
}
