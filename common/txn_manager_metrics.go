package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type TxnManagerMetrics struct {
	Latency  *prometheus.SummaryVec
	GasUsed  prometheus.Gauge
	SpeedUps prometheus.Gauge
	TxQueue  prometheus.Gauge
	NumTx    *prometheus.CounterVec
}

func NewTxnManagerMetrics(namespace string, reg *prometheus.Registry) *TxnManagerMetrics {
	return &TxnManagerMetrics{
		Latency: promauto.With(reg).NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Name:       "txn_manager_latency_ms",
				Help:       "transaction confirmation latency summary in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			},
			[]string{"stage"},
		),
		GasUsed: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "gas_used",
				Help:      "gas used for onchain batch confirmation",
			},
		),
		SpeedUps: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "speed_ups",
				Help:      "number of times the gas price was increased",
			},
		),
		TxQueue: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "tx_queue",
				Help:      "number of transactions in transaction queue",
			},
		),
		NumTx: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tx_total",
				Help:      "number of transactions processed",
			},
			[]string{"state"},
		),
	}
}

func (t *TxnManagerMetrics) ObserveLatency(stage string, latencyMs float64) {
	t.Latency.WithLabelValues(stage).Observe(latencyMs)
}

func (t *TxnManagerMetrics) UpdateGasUsed(gasUsed uint64) {
	t.GasUsed.Set(float64(gasUsed))
}

func (t *TxnManagerMetrics) UpdateSpeedUps(speedUps int) {
	t.SpeedUps.Set(float64(speedUps))
}

func (t *TxnManagerMetrics) UpdateTxQueue(txQueue int) {
	t.TxQueue.Set(float64(txQueue))
}

func (t *TxnManagerMetrics) IncrementTxnCount(state string) {
	t.NumTx.WithLabelValues(state).Inc()
}
