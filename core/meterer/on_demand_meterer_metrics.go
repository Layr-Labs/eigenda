package meterer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Tracks metrics for the [OnDemandMeterer]
type OnDemandMetererMetrics struct {
	onDemandGlobalMeterExhaustedRequests  prometheus.Counter
	onDemandGlobalMeterExhaustedSymbols   prometheus.Counter
	onDemandGlobalMeterThroughputRequests prometheus.Counter
	onDemandGlobalMeterThroughputSymbols  prometheus.Counter
}

func NewOnDemandMetererMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
) *OnDemandMetererMetrics {
	if registry == nil {
		return nil
	}

	onDemandGlobalMeterExhaustedRequests := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_exhausted_requests_count",
			Subsystem: subsystem,
			Help:      "Total number of requests rejected due to global rate limit",
		},
	)

	onDemandGlobalMeterExhaustedSymbols := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_exhausted_symbols_count",
			Subsystem: subsystem,
			Help:      "Total number of symbols rejected due to global rate limit",
		},
	)

	onDemandGlobalMeterThroughputRequests := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_throughput_requests_count",
			Subsystem: subsystem,
			Help:      "Total number of requests successfully metered for on-demand dispersals",
		},
	)

	onDemandGlobalMeterThroughputSymbols := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_throughput_symbols_count",
			Subsystem: subsystem,
			Help:      "Total number of symbols successfully metered for on-demand dispersals",
		},
	)

	return &OnDemandMetererMetrics{
		onDemandGlobalMeterExhaustedRequests:  onDemandGlobalMeterExhaustedRequests,
		onDemandGlobalMeterExhaustedSymbols:   onDemandGlobalMeterExhaustedSymbols,
		onDemandGlobalMeterThroughputRequests: onDemandGlobalMeterThroughputRequests,
		onDemandGlobalMeterThroughputSymbols:  onDemandGlobalMeterThroughputSymbols,
	}
}

// RecordGlobalMeterExhaustion records a request rejection due to global rate limit
func (m *OnDemandMetererMetrics) RecordGlobalMeterExhaustion(symbolCount uint32) {
	if m == nil {
		return
	}
	m.onDemandGlobalMeterExhaustedRequests.Inc()
	m.onDemandGlobalMeterExhaustedSymbols.Add(float64(symbolCount))
}

// RecordGlobalMeterThroughput records successful metering for on-demand dispersals
func (m *OnDemandMetererMetrics) RecordGlobalMeterThroughput(symbolCount uint32) {
	if m == nil {
		return
	}
	m.onDemandGlobalMeterThroughputRequests.Inc()
	m.onDemandGlobalMeterThroughputSymbols.Add(float64(symbolCount))
}
