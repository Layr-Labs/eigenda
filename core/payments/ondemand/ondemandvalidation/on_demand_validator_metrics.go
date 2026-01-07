package ondemandvalidation

import (
	"github.com/Layr-Labs/eigenda/common/nameremapping"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Tracks metrics for the [OnDemandPaymentValidator]
type OnDemandValidatorMetrics struct {
	// Although payments internally tracks things in symbols, the consumer of metrics wants to see things in bytes.
	// For a histogram, it's actually not possible to automatically rename bucket labels in grafana, so using
	// symbols here causes dashboards to be less intuitive.
	onDemandBytes              prometheus.Histogram
	onDemandSymbolsTotal       *prometheus.CounterVec
	onDemandDispersalsTotal    *prometheus.CounterVec
	onDemandInsufficientFunds  prometheus.Counter
	onDemandQuorumNotSupported prometheus.Counter
	onDemandUnexpectedErrors   prometheus.Counter
	enablePerAccountMetrics    bool
	userAccountRemapping       map[string]string
}

func NewOnDemandValidatorMetrics(
	registry *prometheus.Registry,
	namespace string,
	subsystem string,
	enablePerAccountMetrics bool,
	userAccountRemapping map[string]string,
) *OnDemandValidatorMetrics {
	if registry == nil {
		return nil
	}

	bytes := promauto.With(registry).NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "on_demand_bytes",
			Subsystem: subsystem,
			Help: "Distribution of byte counts for successful on-demand payments. " +
				"Counts reflect actual dispersed bytes, not billed bytes (which may be higher due to min size).",
			// Buckets chosen to go from min to max blob sizes (128KiB -> 16MiB)
			Buckets: prometheus.ExponentialBuckets(128*units.KiB, 2, 8),
		},
	)

	symbolsTotal := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_symbols_total",
			Subsystem: subsystem,
			Help: "Total number of symbols validated for successful on-demand payments. " +
				"Counts reflect actual dispersed symbols, not billed symbols (which may be higher due to min size).",
		},
		[]string{"account_id"},
	)

	dispersalsTotal := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "on_demand_dispersals_total",
			Subsystem: subsystem,
			Help:      "Total number of dispersals successfully paid for by on-demand.",
		},
		[]string{"account_id"},
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
		onDemandBytes:              bytes,
		onDemandSymbolsTotal:       symbolsTotal,
		onDemandDispersalsTotal:    dispersalsTotal,
		onDemandInsufficientFunds:  insufficientFunds,
		onDemandQuorumNotSupported: quorumNotSupported,
		onDemandUnexpectedErrors:   unexpectedErrors,
		enablePerAccountMetrics:    enablePerAccountMetrics,
		userAccountRemapping:       userAccountRemapping,
	}
}

// Records a successful on-demand payment
func (m *OnDemandValidatorMetrics) RecordSuccess(accountID string, symbolCount uint32) {
	if m == nil {
		return
	}
	m.onDemandBytes.Observe(float64(symbolCount) * encoding.BYTES_PER_SYMBOL)

	labelValue := nameremapping.GetAccountLabel(accountID, m.userAccountRemapping, m.enablePerAccountMetrics)
	m.onDemandSymbolsTotal.WithLabelValues(labelValue).Add(float64(symbolCount))
	m.onDemandDispersalsTotal.WithLabelValues(labelValue).Inc()
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
