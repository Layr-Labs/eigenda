//nolint:wrapcheck // Directly returning errors from the api package is the correct pattern
package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	"github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// The controller GRPC server
type Server struct {
	controller.UnimplementedControllerServiceServer

	config                      Config
	logger                      logging.Logger
	server                      *grpc.Server
	listener                    net.Listener
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
	metrics                     *metrics.ServerMetrics
	replayGuardian              replay.ReplayGuardian
}

func NewServer(
	ctx context.Context,
	config Config,
	logger logging.Logger,
	metricsRegistry *prometheus.Registry,
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler,
) (*Server, error) {
	replayGuardian := replay.NewReplayGuardian(time.Now, config.RequestMaxPastAge, config.RequestMaxFutureAge)

	return &Server{
		config:                      config,
		logger:                      logger,
		metrics:                     metrics.NewServerMetrics(metricsRegistry, logger),
		paymentAuthorizationHandler: paymentAuthorizationHandler,
		replayGuardian:              replayGuardian,
	}, nil
}

// Start the server. Blocks until the server is stopped.
func (s *Server) Start() error {
	if !s.config.EnableServer {
		return fmt.Errorf("controller gRPC server is disabled")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("start tcp listener: %w", err)
	}
	s.listener = listener

	var opts []grpc.ServerOption
	opts = append(opts, s.metrics.GetGRPCServerOption())

	if s.config.MaxGRPCMessageSize > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(s.config.MaxGRPCMessageSize))
	}

	if s.config.MaxIdleConnectionAge > 0 {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.config.MaxIdleConnectionAge,
		}))
	}

	s.server = grpc.NewServer(opts...)
	reflection.Register(s.server)
	controller.RegisterControllerServiceServer(s.server, s)
	healthcheck.RegisterHealthServer(controller.ControllerService_ServiceDesc.ServiceName, s.server)

	s.logger.Infof("gRPC server listening at %v", listener.Addr().String())

	err = s.server.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			s.logger.Errorf("close listener: %w", err)
		}
	}
}

// Handles an AuthorizePaymentRequest
func (s *Server) AuthorizePayment(
	ctx context.Context,
	request *controller.AuthorizePaymentRequest,
) (*controller.AuthorizePaymentResponse, error) {
	if s.paymentAuthorizationHandler == nil {
		return nil, api.NewErrorInternal(fmt.Sprintf(
			"payment authorization handler not configured, request=%s", request.String()))
	}

	probe := s.metrics.NewPaymentAuthorizationProbe()
	success := false
	defer func() {
		probe.End()
		if !success {
			s.metrics.ReportAuthorizePaymentFailure()
		}
	}()

	probe.SetStage("hash_authorize_payment_request")

	requestHash, err := hashing.HashAuthorizePaymentRequest(request)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("failed to hash request: %v, request=%s", err, request.String()))
	}

	probe.SetStage("replay_protection")

	timestamp := time.Unix(0, request.GetBlobHeader().GetPaymentHeader().GetTimestamp())
	err = s.replayGuardian.VerifyRequest(requestHash, timestamp)
	if err != nil {
		s.metrics.ReportPaymentAuthReplayProtectionFailure()
		return nil, api.NewErrorInvalidArg(fmt.Sprintf(
			"replay protection check failed: %v, request=%s", err, request.String()))
	}

	response, err := s.paymentAuthorizationHandler.AuthorizePayment(
		ctx, request.GetBlobHeader(), request.GetClientSignature(), probe)
	if err != nil {
		return nil, err
	}

	success = true
	return response, nil
}
