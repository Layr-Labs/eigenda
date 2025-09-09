package controller

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	pb.UnimplementedControllerServiceServer

	logger                      logging.Logger
	server                      *grpc.Server
	grpcPort                    string
	listener                    net.Listener
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
}

func NewGrpcServer(
	logger logging.Logger,
	grpcPort string,
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler,
	registry *prometheus.Registry,
) (*GrpcServer, error) {
	if grpcPort == "" {
		return nil, fmt.Errorf("grpc port is required")
	}

	return &GrpcServer{
		grpcPort:                    grpcPort,
		logger:                      logger,
		paymentAuthorizationHandler: paymentAuthorizationHandler,
	}, nil
}

func (s *GrpcServer) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%s", s.grpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener: %w", err)
	}
	s.listener = listener

	// Create gRPC server with metrics interceptor
	grpcMetrics := grpcprom.NewServerMetrics()
	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	reflection.Register(s.server)
	pb.RegisterControllerServiceServer(s.server, s)

	name := pb.ControllerService_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, s.server)

	s.logger.Debugf("gRPC server listening", "address", listener.Addr().String())

	// blocks until the server is stopped
	return s.server.Serve(listener)
}

func (s *GrpcServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *GrpcServer) AuthorizePayment(
	ctx context.Context,
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	if s.paymentAuthorizationHandler == nil {
		return nil, fmt.Errorf("payment authorization handler not configured")
	}

	return s.paymentAuthorizationHandler.AuthorizePayment(ctx, request)
}
