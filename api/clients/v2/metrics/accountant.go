package metrics

import (
	"math/big"

	"github.com/Layr-Labs/eigenda/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	accountantSubsystem = "accountant"
)

var (
	gweiFactor = 1e9 // gweiFactor is used when converting wei to gwei
)

type AccountantMetricer interface {
	RecordCumulativePayment(accountID string, wei *big.Int)
	RecordReservationPayment(remainingCapacity float64)

	Document() []metrics.DocumentedMetric
}

type AccountantMetrics struct {
	CumulativePayment            *prometheus.GaugeVec
	ReservationRemainingCapacity prometheus.Gauge

	factory *metrics.Documentor
}

func NewAccountantMetrics(registry *prometheus.Registry) AccountantMetricer {
	if registry == nil {
		return &noopAccountantMetricer{}
	}

	factory := metrics.With(registry)

	return &AccountantMetrics{
		CumulativePayment: factory.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cumulative_payment",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Current cumulative payment balance (gwei)",
		}, []string{
			"account_id",
		}),
		ReservationRemainingCapacity: factory.NewGauge(prometheus.GaugeOpts{
			Name:      "reservation_remaining_capacity",
			Namespace: namespace,
			Subsystem: accountantSubsystem,
			Help:      "Remaining capacity in reservation bucket (symbols)",
		}),
		factory: factory,
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

func (m *AccountantMetrics) RecordReservationPayment(remainingCapacity float64) {
	m.ReservationRemainingCapacity.Set(remainingCapacity)
}

func (m *AccountantMetrics) Document() []metrics.DocumentedMetric {
	return m.factory.Document()
}

type noopAccountantMetricer struct {
}

var NoopAccountantMetrics AccountantMetricer = new(noopAccountantMetricer)

func (n *noopAccountantMetricer) RecordCumulativePayment(_ string, _ *big.Int) {
}

func (n *noopAccountantMetricer) RecordReservationPayment(_ float64) {
}

func (n *noopAccountantMetricer) Document() []metrics.DocumentedMetric {
	return []metrics.DocumentedMetric{}
}
