package grpcserver

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

const namespace = "eigenda_controller"

// Encapsulates metrics for the controller GRPC service
type Metrics struct {
	logger           logging.Logger
	grpcServerOption grpc.ServerOption

	// AuthorizePayment metrics
	authorizePaymentLatency           *prometheus.SummaryVec
	authorizePaymentSignatureLatency  *prometheus.SummaryVec
	authorizePaymentAuthFailures      *prometheus.CounterVec
	authorizePaymentSignatureFailures *prometheus.CounterVec
}

func NewMetrics(registry *prometheus.Registry, logger logging.Logger) *Metrics {
	if registry == nil {
		registry = prometheus.NewRegistry()
	}

	grpcMetrics := grpcprom.NewServerMetrics()
	registry.MustRegister(grpcMetrics)
	grpcServerOption := grpc.UnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
	)

	objectives := map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

	authorizePaymentLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "authorize_payment_latency_ms",
			Help:       "Latency of the AuthorizePayment RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	authorizePaymentSignatureLatency := promauto.With(registry).NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       "authorize_payment_signature_latency_ms",
			Help:       "Latency of signature verification in AuthorizePayment RPC",
			Objectives: objectives,
		},
		[]string{},
	)

	authorizePaymentAuthFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "authorize_payment_auth_failure_count",
			Help:      "Number of AuthorizePayment RPC authentication failures",
		},
		[]string{},
	)

	authorizePaymentSignatureFailures := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "authorize_payment_signature_failure_count",
			Help:      "Number of AuthorizePayment RPC signature verification failures",
		},
		[]string{},
	)

	return &Metrics{
		logger:                            logger,
		grpcServerOption:                  grpcServerOption,
		authorizePaymentLatency:           authorizePaymentLatency,
		authorizePaymentSignatureLatency:  authorizePaymentSignatureLatency,
		authorizePaymentAuthFailures:      authorizePaymentAuthFailures,
		authorizePaymentSignatureFailures: authorizePaymentSignatureFailures,
	}
}

// GetGRPCServerOption returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *Metrics) GetGRPCServerOption() grpc.ServerOption {
	return m.grpcServerOption
}

// ReportAuthorizePaymentLatency reports the total latency of an AuthorizePayment RPC.
func (m *Metrics) ReportAuthorizePaymentLatency(duration time.Duration) {
	m.authorizePaymentLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// ReportAuthorizePaymentSignatureLatency reports the latency of signature verification in AuthorizePayment.
func (m *Metrics) ReportAuthorizePaymentSignatureLatency(duration time.Duration) {
	m.authorizePaymentSignatureLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// ReportAuthorizePaymentAuthFailure increments the auth failure counter for AuthorizePayment.
func (m *Metrics) ReportAuthorizePaymentAuthFailure() {
	m.authorizePaymentAuthFailures.WithLabelValues().Inc()
}

// ReportAuthorizePaymentSignatureFailure increments the signature failure counter for AuthorizePayment.
func (m *Metrics) ReportAuthorizePaymentSignatureFailure() {
	m.authorizePaymentSignatureFailures.WithLabelValues().Inc()
}
