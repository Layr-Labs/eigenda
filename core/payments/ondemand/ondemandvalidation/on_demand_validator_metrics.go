package ondemandvalidation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Tracks metrics for the [OnDemandPaymentValidator]
type OnDemandValidatorMetrics struct {
	onDemandSymbols            prometheus.Histogram
	onDemandSymbolsTotal       prometheus.Counter
	onDemandInsufficientFunds  prometheus.Counter
	onDemandQuorumNotSupported prometheus.Counter
	onDemandUnexpectedErrors   prometheus.Counter
}

func NewOnDemandValidatorMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
) *OnDemandValidatorMetrics {
	if registry == nil {
		return nil
	}

	symbols := promauto.With(registry).NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "on_demand_symbols",
			Subsystem: subsystem,
			Help: "Distribution of symbol counts for successful on-demand payments. " +
				"Counts reflect actual dispersed symbols, not billed symbols (which may be higher due to min size).",
			// Buckets chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(4096, 2, 8),
		},
	)

	symbolsTotal := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_symbols_total",
			Subsystem: subsystem,
			Help: "Total number of symbols validated for successful on-demand payments. " +
				"Counts reflect actual dispersed symbols, not billed symbols (which may be higher due to min size).",
		},
	)

	insufficientFunds := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_insufficient_funds_count",
			Subsystem: subsystem,
			Help:      "Total number of on-demand payments rejected due to insufficient funds",
		},
	)

	quorumNotSupported := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_quorum_not_supported_count",
			Subsystem: subsystem,
			Help:      "Total number of on-demand payments rejected due to unsupported quorums",
		},
	)

	unexpectedErrors := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_unexpected_errors_count",
			Subsystem: subsystem,
			Help:      "Total number of unexpected errors during on-demand payment authorization",
		},
	)

	return &OnDemandValidatorMetrics{
		onDemandSymbols:            symbols,
		onDemandSymbolsTotal:       symbolsTotal,
		onDemandInsufficientFunds:  insufficientFunds,
		onDemandQuorumNotSupported: quorumNotSupported,
		onDemandUnexpectedErrors:   unexpectedErrors,
	}
}

// Records a successful on-demand payment
func (m *OnDemandValidatorMetrics) RecordSuccess(symbolCount uint32) {
	if m == nil {
		return
	}
	m.onDemandSymbols.Observe(float64(symbolCount))
	m.onDemandSymbolsTotal.Add(float64(symbolCount))
}

// Increments the counter for insufficient funds errors
func (m *OnDemandValidatorMetrics) IncrementInsufficientFunds() {
	if m == nil {
		return
	}
	m.onDemandInsufficientFunds.Inc()
}

// Increments the counter for unsupported quorum errors
func (m *OnDemandValidatorMetrics) IncrementQuorumNotSupported() {
	if m == nil {
		return
	}
	m.onDemandQuorumNotSupported.Inc()
}

// Increments the counter for unexpected errors
func (m *OnDemandValidatorMetrics) IncrementUnexpectedErrors() {
	if m == nil {
		return
	}
	m.onDemandUnexpectedErrors.Inc()
}
