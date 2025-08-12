package metrics

import (
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace           = "eigenda"
	accountantSubsystem = "accountant"
)

var (
	gweiFactor = big.NewInt(1000000000)
)

type AccountantMetricer interface {
	RecordCumulativePayment(accountID string, wei *big.Int)
}

type AccountantMetrics struct {
	CumulativePayment *prometheus.GaugeVec
}

func NewAccountantMetrics(registry *prometheus.Registry) AccountantMetricer {
	if registry == nil {
		return &noopAccountantMetricer{}
	}

	return &AccountantMetrics{
		CumulativePayment: promauto.With(registry).NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cumulative_payment",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Current cumulative payment balance (gwei)",
		}, []string{
			"account_id",
		}),
	}
}

func (m *AccountantMetrics) RecordCumulativePayment(accountID string, wei *big.Int) {
	// Convert wei to gwei
	gwei := new(big.Int).Div(wei, gweiFactor)
	gweiFloat, _ := gwei.Float64()
	m.CumulativePayment.WithLabelValues(accountID).Set(gweiFloat)
}

type noopAccountantMetricer struct {
}

var NoopAccountantMetrics AccountantMetricer = new(noopAccountantMetricer)

func (n *noopAccountantMetricer) RecordCumulativePayment(_ string, _ *big.Int) {
}
