package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
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

	registry *prometheus.Registry
	httpPort string
	logger   logging.Logger
}

// newAPIServerV2Metrics creates a new metricsV2 instance.
func newAPIServerV2Metrics(registry *prometheus.Registry, metricsConfig disperser.MetricsConfig, logger logging.Logger) *metricsV2 {
	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	getBlobCommitmentLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_blob_commitment_latency_ms",
			Help:       "The time required to get the blob commitment.",
			Objectives: objectives,
		},
		[]string{},
	)

	getPaymentStateLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "get_payment_state_latency_ms",
			Help:       "The time required to get the payment state.",
			Objectives: objectives,
		},
		[]string{},
	)

	disperseBlobLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "disperse_blob_latency_ms",
			Help:       "The time required to disperse a blob.",
			Objectives: objectives,
		},
		[]string{},
	)

	disperseBlobSize := promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "disperse_blob_size_bytes",
			Help:      "The size of the blob in bytes.",
		},
		[]string{},
	)

	validateDispersalRequestLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "validate_dispersal_request_latency_ms",
			Help:       "The time required to validate a dispersal request.",
			Objectives: objectives,
		},
		[]string{},
	)

	storeBlobLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "store_blob_latency_ms",
			Help:       "The time required to store a blob.",
			Objectives: objectives,
		},
		[]string{},
	)

	getBlobStatusLatency := promauto.With(registry).NewSummaryVec(
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
		registry:                        registry,
		httpPort:                        metricsConfig.HTTPPort,
		logger:                          logger.With("component", "DisperserV2Metrics"),
	}
}

// Start the metrics server
func (m *metricsV2) Start(ctx context.Context) {
	m.logger.Info("Starting metrics server at ", "port", m.httpPort)
	addr := fmt.Sprintf(":%s", m.httpPort)
	go func() {
		log := m.logger
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			m.registry,
			promhttp.HandlerOpts{},
		))
		err := http.ListenAndServe(addr, mux)
		log.Error("Prometheus server failed", "err", err)
	}()
}

func (m *metricsV2) reportGetBlobCommitmentLatency(duration time.Duration) {
	m.getBlobCommitmentLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *metricsV2) reportGetPaymentStateLatency(duration time.Duration) {
	m.getPaymentStateLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *metricsV2) reportDisperseBlobLatency(duration time.Duration) {
	m.disperseBlobLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *metricsV2) reportDisperseBlobSize(size int) {
	m.disperseBlobSize.WithLabelValues().Set(float64(size))
}

func (m *metricsV2) reportValidateDispersalRequestLatency(duration time.Duration) {
	m.validateDispersalRequestLatency.WithLabelValues().Observe(
		common.ToMilliseconds(duration))
}

func (m *metricsV2) reportStoreBlobLatency(duration time.Duration) {
	m.storeBlobLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

func (m *metricsV2) reportGetBlobStatusLatency(duration time.Duration) {
	m.getBlobStatusLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}
