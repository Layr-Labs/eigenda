package meterer

/* NOTHING TO SEE HERE ; NO METRICS YET */
import (
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type FailReason string

const (
	FailBatchHeaderHash        FailReason = "batch_header_hash"
	FailAggregateSignatures    FailReason = "aggregate_signatures"
	FailNoSignatures           FailReason = "no_signatures"
	FailConfirmBatch           FailReason = "confirm_batch"
	FailGetBatchID             FailReason = "get_batch_id"
	FailUpdateConfirmationInfo FailReason = "update_confirmation_info"
	FailNoAggregatedSignature  FailReason = "no_aggregated_signature"
)

type MetricsConfig struct {
	HTTPPort      string
	EnableMetrics bool
}

type ReservationMetrics struct {
}

type OnDemandPaymentMetrics struct {
}

type Metrics struct {
	*ReservationMetrics
	*OnDemandPaymentMetrics

	registry *prometheus.Registry

	httpPort string
	logger   logging.Logger
}

func NewMetrics(httpPort string, logger logging.Logger) *Metrics {
	// namespace := "eigenda_meterer"
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())
	reservationMetrics := ReservationMetrics{}
	onDemandPaymentMetrics := OnDemandPaymentMetrics{}
	metrics := &Metrics{
		ReservationMetrics:     &reservationMetrics,
		OnDemandPaymentMetrics: &onDemandPaymentMetrics,
		registry:               reg,
		httpPort:               httpPort,
		logger:                 logger.With("component", "BatcherMetrics"),
	}
	return metrics
}

// type noopMetrics struct{}

// var NoopMetrics Metrics = new(noopMetrics)
