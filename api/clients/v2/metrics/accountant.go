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
	gweiFactor = 1e9 // gweiFactor is used when converting wei to gwei
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
	// The prometheus.GaugeVec uses a float64. To minimize precision loss when
	// converting from wei, the cumulative payment value is first converted
	// to gwei before reporting the metric. Users can perform transformations
	// on the value via dashboard functions to change denomination.
	gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(gweiFactor))
	gweiFloat64, _ := gwei.Float64()
	m.CumulativePayment.WithLabelValues(accountID).Set(gweiFloat64)
}

type noopAccountantMetricer struct {
}

var NoopAccountantMetrics AccountantMetricer = new(noopAccountantMetricer)

func (n *noopAccountantMetricer) RecordCumulativePayment(_ string, _ *big.Int) {
}
