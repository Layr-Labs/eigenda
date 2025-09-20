package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PaymentAuthorizationMetrics struct {
	registry                      *prometheus.Registry
	reservationLedgerCacheSize    prometheus.GaugeFunc
	onDemandLedgerCacheSize       prometheus.GaugeFunc
	onDemandPaymentsCount         prometheus.Counter
	reservationPaymentsCount      prometheus.Counter
	onDemandGlobalMeterExhausted  prometheus.Counter
	onDemandInsufficientFunds     prometheus.Counter
	onDemandQuorumNotSupported    prometheus.Counter
	reservationInsufficientFunds  prometheus.Counter
	reservationQuorumNotPermitted prometheus.Counter
	reservationTimeOutOfRange     prometheus.Counter
	reservationTimeMovedBackward  prometheus.Counter
	onDemandUnexpectedErrors      prometheus.Counter
	reservationUnexpectedErrors   prometheus.Counter
	reservationSymbolCount        prometheus.Histogram
	onDemandSymbolCount           prometheus.Histogram
}

func NewPaymentAuthorizationMetrics(registry *prometheus.Registry) *PaymentAuthorizationMetrics {
	if registry == nil {
		return nil
	}

	onDemandPaymentsCount := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_payments_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of successful on-demand payments processed",
		},
	)

	reservationPaymentsCount := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_payments_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of successful reservation payments processed",
		},
	)

	onDemandGlobalMeterExhausted := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_global_meter_exhausted_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of on-demand payments rejected due to global rate limit",
		},
	)

	onDemandInsufficientFunds := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_insufficient_funds_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of on-demand payments rejected due to insufficient funds",
		},
	)

	onDemandQuorumNotSupported := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_quorum_not_supported_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of on-demand payments rejected due to unsupported quorums",
		},
	)

	reservationInsufficientBandwidth := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_insufficient_bandwidth_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of reservation payments rejected due to insufficient bandwidth",
		},
	)

	reservationQuorumNotPermitted := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_quorum_not_permitted_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of reservation payments rejected due to unpermitted quorums",
		},
	)

	reservationTimeOutOfRange := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_time_out_of_range_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of reservation payments rejected due to time out of range",
		},
	)

	reservationTimeMovedBackward := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_time_moved_backward_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of reservation payments rejected due to time moving backwards",
		},
	)

	onDemandUnexpectedErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_unexpected_errors_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of unexpected errors during on-demand payment authorization",
		},
	)

	reservationUnexpectedErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "reservation_unexpected_errors_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Total number of unexpected errors during reservation payment authorization",
		},
	)

	reservationSymbolCount := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "reservation_symbol_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Distribution of symbol counts for successful reservation payments",
			// bucket sizes chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(4096, 2, 8),
		},
	)

	onDemandSymbolCount := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "on_demand_symbol_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Distribution of symbol counts for successful on-demand payments",
			// bucket sizes chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(4096, 2, 8),
		},
	)

	registry.MustRegister(
		onDemandPaymentsCount,
		reservationPaymentsCount,
		onDemandGlobalMeterExhausted,
		onDemandInsufficientFunds,
		onDemandQuorumNotSupported,
		reservationInsufficientBandwidth,
		reservationQuorumNotPermitted,
		reservationTimeOutOfRange,
		reservationTimeMovedBackward,
		onDemandUnexpectedErrors,
		reservationUnexpectedErrors,
		reservationSymbolCount,
		onDemandSymbolCount)

	return &PaymentAuthorizationMetrics{
		registry:                      registry,
		onDemandPaymentsCount:         onDemandPaymentsCount,
		reservationPaymentsCount:      reservationPaymentsCount,
		onDemandGlobalMeterExhausted:  onDemandGlobalMeterExhausted,
		onDemandInsufficientFunds:     onDemandInsufficientFunds,
		onDemandQuorumNotSupported:    onDemandQuorumNotSupported,
		reservationInsufficientFunds:  reservationInsufficientBandwidth,
		reservationQuorumNotPermitted: reservationQuorumNotPermitted,
		reservationTimeOutOfRange:     reservationTimeOutOfRange,
		reservationTimeMovedBackward:  reservationTimeMovedBackward,
		onDemandUnexpectedErrors:      onDemandUnexpectedErrors,
		reservationUnexpectedErrors:   reservationUnexpectedErrors,
		reservationSymbolCount:        reservationSymbolCount,
		onDemandSymbolCount:           onDemandSymbolCount,
	}
}

func (m *PaymentAuthorizationMetrics) RegisterReservationCacheSize(getter func() int) {
	if m == nil || m.registry == nil {
		return
	}

	m.reservationLedgerCacheSize = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "reservation_ledger_cache_size",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Current number of entries in the reservation ledger cache",
		},
		func() float64 {
			return float64(getter())
		},
	)

	m.registry.MustRegister(m.reservationLedgerCacheSize)
}

func (m *PaymentAuthorizationMetrics) RegisterOnDemandCacheSize(getter func() int) {
	if m == nil || m.registry == nil {
		return
	}

	m.onDemandLedgerCacheSize = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "on_demand_ledger_cache_size",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Current number of entries in the on-demand ledger cache",
		},
		func() float64 {
			return float64(getter())
		},
	)

	m.registry.MustRegister(m.onDemandLedgerCacheSize)
}

// Records metrics for a successful on-demand payment.
func (m *PaymentAuthorizationMetrics) RecordOnDemandPaymentSuccess(symbolCount uint32) {
	if m == nil {
		return
	}

	m.onDemandSymbolCount.Observe(float64(symbolCount))
	m.onDemandPaymentsCount.Inc()
}

// Records metrics for a successful reservation payment.
func (m *PaymentAuthorizationMetrics) RecordReservationPaymentSuccess(symbolCount uint32) {
	if m == nil {
		return
	}

	m.reservationSymbolCount.Observe(float64(symbolCount))
	m.reservationPaymentsCount.Inc()
}

// Increments the counter for on-demand payments rejected due to global rate limit.
func (m *PaymentAuthorizationMetrics) IncrementOnDemandGlobalMeterExhausted() {
	if m == nil {
		return
	}

	m.onDemandGlobalMeterExhausted.Inc()
}

// Increments the counter for on-demand payments rejected due to insufficient funds.
func (m *PaymentAuthorizationMetrics) IncrementOnDemandInsufficientFunds() {
	if m == nil {
		return
	}

	m.onDemandInsufficientFunds.Inc()
}

// Increments the counter for on-demand payments rejected due to unsupported quorums.
func (m *PaymentAuthorizationMetrics) IncrementOnDemandQuorumNotSupported() {
	if m == nil {
		return
	}

	m.onDemandQuorumNotSupported.Inc()
}

// Increments the counter for reservation payments rejected due to insufficient bandwidth.
func (m *PaymentAuthorizationMetrics) IncrementReservationInsufficientFunds() {
	if m == nil {
		return
	}

	m.reservationInsufficientFunds.Inc()
}

// Increments the counter for reservation payments rejected due to unpermitted quorums.
func (m *PaymentAuthorizationMetrics) IncrementReservationQuorumNotPermitted() {
	if m == nil {
		return
	}

	m.reservationQuorumNotPermitted.Inc()
}

// Increments the counter for reservation payments rejected due to time out of range.
func (m *PaymentAuthorizationMetrics) IncrementReservationTimeOutOfRange() {
	if m == nil {
		return
	}

	m.reservationTimeOutOfRange.Inc()
}

// Increments the counter for reservation payments rejected due to time moving backwards.
func (m *PaymentAuthorizationMetrics) IncrementReservationTimeMovedBackward() {
	if m == nil {
		return
	}

	m.reservationTimeMovedBackward.Inc()
}

// Increments the counter for unexpected errors during on-demand payment authorization.
func (m *PaymentAuthorizationMetrics) IncrementOnDemandUnexpectedErrors() {
	if m == nil {
		return
	}

	m.onDemandUnexpectedErrors.Inc()
}

// Increments the counter for unexpected errors during reservation payment authorization.
func (m *PaymentAuthorizationMetrics) IncrementReservationUnexpectedErrors() {
	if m == nil {
		return
	}

	m.reservationUnexpectedErrors.Inc()
}
