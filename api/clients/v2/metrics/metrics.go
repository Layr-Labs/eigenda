package metrics

import (
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	accountantSubsystem = "accountant"
	disperserSubsystem  = "disperser_client"
)

type ClientMetricer interface {
	// accountant
	ReportPaymentUsed(accountID string, method string, symbols uint64)
	ReportOnDemandPayment(accountID string, amount *big.Int)
	ReportCumulativePayment(accountID string, amount *big.Int)
	ReportPaymentFailure(accountID string, reason string)
}

type ClientMetrics struct {
	// accountant
	PaymentMethodUsed *prometheus.CounterVec
	OnDemandPayment   *prometheus.CounterVec
	CumulativePayment *prometheus.GaugeVec
	PaymentFailures   *prometheus.CounterVec
}

func NewClientMetrics(namespace string, factory metrics.Factory) ClientMetrics {
	return ClientMetrics{
		PaymentMethodUsed: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payment_method_used_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total number of payments by method (reservation or on-demand)",
		}, []string{
			"account_id",
			"method",
		}),
		OnDemandPayment: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "on_demand_payment_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total on-demand payments made",
		}, []string{
			"account_id",
		}),
		CumulativePayment: factory.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cumulative_payment",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Current cumulative payment balance",
		}, []string{
			"account_id",
		}),
		PaymentFailures: factory.NewCounterVec(prometheus.CounterOpts{
			Name:      "payment_failures_total",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Total payment failures by reason",
		}, []string{
			"account_id",
			"reason",
		}),
	}
}

func (m *ClientMetrics) ReportPaymentUsed(accountID string, method string, symbols uint64) {
	m.PaymentMethodUsed.WithLabelValues(accountID, method).Add(float64(symbols))
}

// ReportOnDemandPayment records an on-demand payment amount in wei
func (m *ClientMetrics) ReportOnDemandPayment(accountID string, amount *big.Int) {
	// Convert big.Int to float64 for prometheus
	weiFloat, _ := amount.Float64() // TODO(iquidus)
	m.OnDemandPayment.WithLabelValues(accountID).Add(weiFloat)
}

func (m *ClientMetrics) ReportCumulativePayment(accountID string, amount *big.Int) {
	// Convert big.Int to float64 for prometheus
	weiFloat, _ := amount.Float64() // TODO(iquidus)
	m.CumulativePayment.WithLabelValues(accountID).Set(weiFloat)
}

func (m *ClientMetrics) ReportPaymentFailure(accountID string, reason string) {
	m.PaymentFailures.WithLabelValues(accountID, reason).Inc()
}

type NoopAccountantMetricer struct {
}

var NoopAccountantMetrics ClientMetricer = new(NoopAccountantMetricer)

func (m *NoopAccountantMetricer) ReportPaymentUsed(_ string, _ string, _ uint64) {
}

func (m *NoopAccountantMetricer) ReportOnDemandPayment(_ string, _ *big.Int) {
}

func (m *NoopAccountantMetricer) ReportCumulativePayment(_ string, _ *big.Int) {
}

func (m *NoopAccountantMetricer) ReportPaymentFailure(_ string, _ string) {
}
