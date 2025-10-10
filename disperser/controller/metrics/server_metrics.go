package metrics

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

const (
	Namespace                  = "eigenda_controller"
	AuthorizePaymentsSubsystem = "authorize_payments"
)

// Encapsulates metrics for the controller GRPC server
type ServerMetrics struct {
	logger           logging.Logger
	grpcServerOption grpc.ServerOption

	paymentAuthorizationStageTimer *common.StageTimer
	paymentAuthorizationFailures   prometheus.Counter
	paymentAuthorizationReplays    prometheus.Counter
	refundPaymentFailures          prometheus.Counter
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

	paymentAuthorizationFailures := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "payment_authorization_failure_count",
			Subsystem: AuthorizePaymentsSubsystem,
			Help:      "Number of AuthorizePayment RPC failures",
		},
	)

	paymentAuthorizationReplays := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "payment_authorization_replay_count",
			Subsystem: AuthorizePaymentsSubsystem,
			Help:      "Number of payment authorization requests rejected due to replay detection",
		},
	)

	refundPaymentFailures := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "refund_payment_failure_count",
			Subsystem: AuthorizePaymentsSubsystem,
			Help:      "Number of RefundPayment RPC failures",
		},
	)

	paymentAuthorizationStageTimer := common.NewStageTimer(registry, Namespace, "payment_authorization", false)

	return &ServerMetrics{
		logger:                         logger,
		grpcServerOption:               grpcServerOption,
		paymentAuthorizationStageTimer: paymentAuthorizationStageTimer,
		paymentAuthorizationFailures:   paymentAuthorizationFailures,
		paymentAuthorizationReplays:    paymentAuthorizationReplays,
		refundPaymentFailures:          refundPaymentFailures,
	}
}

// Returns the gRPC server option that enables automatic GRPC metrics collection.
func (m *ServerMetrics) GetGRPCServerOption() grpc.ServerOption {
	if m == nil {
		return nil
	}

	return m.grpcServerOption
}

// Increments the auth failure counter for AuthorizePayment.
func (m *ServerMetrics) ReportAuthorizePaymentFailure() {
	if m == nil {
		return
	}

	m.paymentAuthorizationFailures.Inc()
}

// Increments the payment auth replay protection failure counter.
func (m *ServerMetrics) ReportPaymentAuthReplayProtectionFailure() {
	if m == nil {
		return
	}

	m.paymentAuthorizationReplays.Inc()
}

// Increments the refund failure counter for RefundPayment.
func (m *ServerMetrics) ReportRefundPaymentFailure() {
	if m == nil {
		return
	}

	m.refundPaymentFailures.Inc()
}

// Creates a new SequenceProbe for tracking payment authorization stages.
func (m *ServerMetrics) NewPaymentAuthorizationProbe() *common.SequenceProbe {
	if m == nil || m.paymentAuthorizationStageTimer == nil {
		return nil
	}
	return m.paymentAuthorizationStageTimer.NewSequence()
}
