package ondemand

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Tracks metrics for the [OnDemandPaymentValidator]
type OnDemandValidatorMetrics struct {
	onDemandPaymentsCount      prometheus.Counter
	onDemandSymbolCount        prometheus.Histogram
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

	paymentsCount := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_payments_count",
			Subsystem: subsystem,
			Help:      "Total number of successful on-demand payments processed",
		},
	)

	symbolCount := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "on_demand_symbol_count",
			Subsystem: subsystem,
			Help:      "Distribution of symbol counts for successful on-demand payments",
			// Buckets chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(4096, 2, 8),
		},
	)

	insufficientFunds := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_insufficient_funds_count",
			Subsystem: subsystem,
			Help:      "Total number of on-demand payments rejected due to insufficient funds",
		},
	)

	quorumNotSupported := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_quorum_not_supported_count",
			Subsystem: subsystem,
			Help:      "Total number of on-demand payments rejected due to unsupported quorums",
		},
	)

	unexpectedErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_unexpected_errors_count",
			Subsystem: subsystem,
			Help:      "Total number of unexpected errors during on-demand payment authorization",
		},
	)

	registry.MustRegister(
		paymentsCount,
		symbolCount,
		insufficientFunds,
		quorumNotSupported,
		unexpectedErrors,
	)

	return &OnDemandValidatorMetrics{
		onDemandPaymentsCount:      paymentsCount,
		onDemandSymbolCount:        symbolCount,
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
	m.onDemandSymbolCount.Observe(float64(symbolCount))
	m.onDemandPaymentsCount.Inc()
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
