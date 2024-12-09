package apiserver

import (
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

const namespace = "eigenda_disperser_api"

// metricsV2 encapsulates the metrics for the v2 API server.
type metricsV2 struct {
	grpcServerOption grpc.ServerOption

	getBlobCommitmentLatency        *prometheus.SummaryVec
	getPaymentStateLatency          *prometheus.SummaryVec
	disperseBlobLatency             *prometheus.SummaryVec
	disperseBlobSize                *prometheus.GaugeVec
	validateDispersalRequestLatency *prometheus.SummaryVec
	storeBlobLatency                *prometheus.SummaryVec
	getBlobStatusLatency            *prometheus.SummaryVec
}

// newAPIServerV2Metrics creates a new metricsV2 instance.
func newAPIServerV2Metrics(registry *prometheus.Registry) *metricsV2 {
	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)

	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	getBlobCommitmentLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_commitment_latency_ms",
			Help:       "The time required to get the blob commitment.",
			Objectives: objectives,
		},
		[]string{},
	)

	getPaymentStateLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_payment_state_latency_ms",
			Help:       "The time required to get the payment state.",
			Objectives: objectives,
		},
		[]string{},
	)

	disperseBlobLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "disperse_blob_latency_ms",
			Help:       "The time required to disperse a blob.",
			Objectives: objectives,
		},
		[]string{},
	)

	disperseBlobSize := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "disperse_blob_size_bytes",
			Help:      "The size of the blob in bytes.",
		},
		[]string{},
	)

	validateDispersalRequestLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "validate_dispersal_request_latency_ms",
			Help:       "The time required to validate a dispersal request.",
			Objectives: objectives,
		},
		[]string{},
	)

	storeBlobLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "store_blob_latency_ms",
			Help:       "The time required to store a blob.",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobStatusLatency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_status_latency_ms",
			Help:       "The time required to get the blob status.",
			Objectives: objectives,
		},
		[]string{},
	)

	return &metricsV2{
		grpcServerOption:                grpcServerOption,
		getBlobCommitmentLatency:        getBlobCommitmentLatency,
		getPaymentStateLatency:          getPaymentStateLatency,
		disperseBlobLatency:             disperseBlobLatency,
		disperseBlobSize:                disperseBlobSize,
		validateDispersalRequestLatency: validateDispersalRequestLatency,
		storeBlobLatency:                storeBlobLatency,
		getBlobStatusLatency:            getBlobStatusLatency,
	}
}

func (m *metricsV2) reportGetBlobCommitmentLatency(duration time.Duration) {
	m.getBlobCommitmentLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *metricsV2) reportGetPaymentStateLatency(duration time.Duration) {
	m.getPaymentStateLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *metricsV2) reportDisperseBlobLatency(duration time.Duration) {
	m.disperseBlobLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *metricsV2) reportDisperseBlobSize(size int) {
	m.disperseBlobSize.WithLabelValues().Set(float64(size))
}

func (m *metricsV2) reportValidateDispersalRequestLatency(duration time.Duration) {
	m.validateDispersalRequestLatency.WithLabelValues().Observe(
		float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *metricsV2) reportStoreBlobLatency(duration time.Duration) {
	m.storeBlobLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}

func (m *metricsV2) reportGetBlobStatusLatency(duration time.Duration) {
	m.getBlobStatusLatency.WithLabelValues().Observe(float64(duration.Nanoseconds()) / float64(time.Millisecond))
}
