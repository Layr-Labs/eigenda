package metrics

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

// Encapsulates metrics for the controller GRPC server
type ServerMetrics struct {
	logger           logging.Logger
	grpcServerOption grpc.ServerOption

	// AuthorizePayment metrics
	authorizePaymentLatency      *prometheus.SummaryVec
	authorizePaymentAuthFailures *prometheus.CounterVec
}

func NewServerMetrics(registry *prometheus.Registry, logger logging.Logger) *ServerMetrics {
	if registry == nil {
		return nil
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
			Help:       "Total latency of the AuthorizePayment RPC",
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

	return &ServerMetrics{
		logger:                       logger,
		grpcServerOption:             grpcServerOption,
		authorizePaymentLatency:      authorizePaymentLatency,
		authorizePaymentAuthFailures: authorizePaymentAuthFailures,
	}
}

// Returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *ServerMetrics) GetGRPCServerOption() grpc.ServerOption {
	if m == nil {
		return nil
	}

	return m.grpcServerOption
}

// Reports the total latency of an AuthorizePayment RPC.
func (m *ServerMetrics) ReportAuthorizePaymentLatency(duration time.Duration) {
	if m == nil {
		return
	}

	m.authorizePaymentLatency.WithLabelValues().Observe(common.ToMilliseconds(duration))
}

// Increments the auth failure counter for AuthorizePayment.
func (m *ServerMetrics) ReportAuthorizePaymentAuthFailure() {
	if m == nil {
		return
	}

	m.authorizePaymentAuthFailures.WithLabelValues().Inc()
}
