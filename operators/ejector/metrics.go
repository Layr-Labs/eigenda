package ejector

import (
	"fmt"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/codes"
)

type Metrics struct {
	PeriodicEjectionRequests *prometheus.CounterVec
	UrgentEjectionRequests   *prometheus.CounterVec
	OperatorsToEject         *prometheus.CounterVec
	StakeShareToEject        *prometheus.GaugeVec
	EjectionGasUsed          prometheus.Gauge
}

func NewMetrics(reg *prometheus.Registry, logger logging.Logger) *Metrics {
	namespace := "eigenda_ejector"
	metrics := &Metrics{
		// PeriodicEjectionRequests is a more detailed metric than NumRequests, specifically for
		// tracking the ejection calls that are periodically initiated according to the SLA
		// evaluation time window.
		PeriodicEjectionRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "periodic_ejection_requests_total",
				Help:      "the total number of periodic ejection requests",
			},
			[]string{"status"},
		),
		// UrgentEjectionRequests is a more detailed metric than NumRequests, specifically for
		// tracking the ejection calls that are urgently initiated due to bad network health
		// condition.
		UrgentEjectionRequests: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "urgent_ejection_requests_total",
				Help:      "the total number of urgent ejection requests",
			},
			[]string{"status"},
		),
		// The number of operators requested to eject. Note this may be different than the
		// actual number of operators ejected as EjectionManager contract may perform rate
		// limiting.
		OperatorsToEject: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "operators_to_eject",
				Help:      "the total number of operators requested to eject",
			}, []string{"quorum"},
		),
		// The total stake share requested to eject. Note this may be different than the
		// actual stake share ejected as EjectionManager contract may perform rate limiting.
		StakeShareToEject: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "stake_share_to_eject",
				Help:      "the total stake share requested to eject",
			}, []string{"quorum"},
		),
		// The gas used by EjectionManager contract for operator ejection.
		EjectionGasUsed: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "ejection_gas_used",
				Help:      "Gas used for operator ejection",
			},
		),
	}
	return metrics
}

func (g *Metrics) IncrementEjectionRequest(mode Mode, status codes.Code) {
	switch mode {
	case PeriodicMode:
		g.PeriodicEjectionRequests.With(prometheus.Labels{
			"status": status.String(),
		}).Inc()
	case UrgentMode:
		g.UrgentEjectionRequests.With(prometheus.Labels{
			"status": status.String(),
		}).Inc()
	}
}

func (g *Metrics) UpdateEjectionGasUsed(gasUsed uint64) {
	g.EjectionGasUsed.Set(float64(gasUsed))
}

func (g *Metrics) UpdateRequestedOperatorMetric(numOperatorsByQuorum map[uint8]int, stakeShareByQuorum map[uint8]float64) {
	for q, count := range numOperatorsByQuorum {
		for i := 0; i < count; i++ {
			g.OperatorsToEject.With(prometheus.Labels{
				"quorum": fmt.Sprintf("%d", q),
			}).Inc()
		}
	}
	for q, stakeShare := range stakeShareByQuorum {
		g.StakeShareToEject.With(prometheus.Labels{
			"quorum": fmt.Sprintf("%d", q),
		}).Set(stakeShare)
	}
}
