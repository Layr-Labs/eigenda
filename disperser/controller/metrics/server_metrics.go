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
	namespace                  = "eigenda_controller"
	authorizePaymentsSubsystem = "authorize_payments"
)

// Encapsulates metrics for the controller GRPC server
type ServerMetrics struct {
	logger           logging.Logger
	grpcServerOption grpc.ServerOption

	paymentAuthorizationStageTimer *common.StageTimer
	paymentAuthorizationFailures   prometheus.Counter
	paymentAuthorizationReplays    prometheus.Counter
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
			Namespace: namespace,
			Name:      "payment_authorization_failure_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Number of AuthorizePayment RPC failures",
		},
	)

	paymentAuthorizationReplays := promauto.With(registry).NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "payment_authorization_replay_count",
			Subsystem: authorizePaymentsSubsystem,
			Help:      "Number of payment authorization requests rejected due to replay detection",
		},
	)

	paymentAuthorizationStageTimer := common.NewStageTimer(registry, namespace, "payment_authorization", false)

	return &ServerMetrics{
		logger:                         logger,
		grpcServerOption:               grpcServerOption,
		paymentAuthorizationStageTimer: paymentAuthorizationStageTimer,
		paymentAuthorizationFailures:   paymentAuthorizationFailures,
		paymentAuthorizationReplays:    paymentAuthorizationReplays,
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

// Creates a new SequenceProbe for tracking payment authorization stages.
func (m *ServerMetrics) NewPaymentAuthorizationProbe() *common.SequenceProbe {
	if m == nil || m.paymentAuthorizationStageTimer == nil {
		return nil
	}
	return m.paymentAuthorizationStageTimer.NewSequence()
}
