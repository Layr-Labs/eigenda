package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedControllerServiceServer

	config                      Config
	logger                      logging.Logger
	server                      *grpc.Server
	listener                    net.Listener
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
	metrics                     *Metrics
}

func NewServer(
	config Config,
	metricsRegistry *prometheus.Registry,
	logger logging.Logger,
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler,
) (*Server, error) {
	return &Server{
		config:                      config,
		logger:                      logger,
		paymentAuthorizationHandler: paymentAuthorizationHandler,
		metrics:                     NewMetrics(metricsRegistry, logger),
	}, nil
}

func (s *Server) Start() error {
	if !s.config.Enable {
		s.logger.Info("Controller gRPC server is disabled")
		return nil
	}

	addr := fmt.Sprintf("0.0.0.0:%s", s.config.GrpcPort)
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
	pb.RegisterControllerServiceServer(s.server, s)
	healthcheck.RegisterHealthServer(pb.ControllerService_ServiceDesc.ServiceName, s.server)

	s.logger.Info("gRPC server listening", "port", s.config.GrpcPort, "address", listener.Addr().String())

	// blocks until the server is stopped
	return s.server.Serve(listener)
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *Server) AuthorizePayment(
	ctx context.Context,
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	start := time.Now()

	if s.paymentAuthorizationHandler == nil {
		s.metrics.ReportAuthorizePaymentAuthFailure()
		return nil, status.Error(codes.FailedPrecondition, "payment authorization handler not configured")
	}

	reply, err := s.paymentAuthorizationHandler.AuthorizePayment(ctx, request)
	if err != nil {
		// Check error type and report appropriate failure metric
		if status.Code(err) == codes.Unauthenticated {
			s.metrics.ReportAuthorizePaymentSignatureFailure()
		} else {
			s.metrics.ReportAuthorizePaymentAuthFailure()
		}
		// Return the error as-is if it's already a gRPC status error, otherwise wrap it as Internal
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "authorize payment: %v", err)
	}

	s.metrics.ReportAuthorizePaymentLatency(time.Since(start))
	return reply, nil
}
