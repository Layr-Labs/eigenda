package meterer

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Tracks metrics for the [OnDemandMeterer]
type OnDemandMetererMetrics struct {
	onDemandGlobalMeterExhaustedRequests prometheus.Counter
	onDemandGlobalMeterExhaustedSymbols  prometheus.Counter
}

func NewOnDemandMetererMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
) *OnDemandMetererMetrics {
	if registry == nil {
		return nil
	}

	onDemandGlobalMeterExhaustedRequests := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_exhausted_requests_count",
			Subsystem: subsystem,
			Help:      "Total number of requests rejected due to global rate limit",
		},
	)

	onDemandGlobalMeterExhaustedSymbols := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_exhausted_symbols_count",
			Subsystem: subsystem,
			Help:      "Total number of symbols rejected due to global rate limit",
		},
	)

	registry.MustRegister(onDemandGlobalMeterExhaustedRequests, onDemandGlobalMeterExhaustedSymbols)

	return &OnDemandMetererMetrics{
		onDemandGlobalMeterExhaustedRequests: onDemandGlobalMeterExhaustedRequests,
		onDemandGlobalMeterExhaustedSymbols:  onDemandGlobalMeterExhaustedSymbols,
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
