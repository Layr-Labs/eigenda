package metering

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	pb.UnimplementedControllerServer

	config   *Config
	logger   logging.Logger
	server   *grpc.Server
	listener net.Listener
}

type Config struct {
	GrpcPort string
}

func NewGrpcServer(config *Config, logger logging.Logger, registry *prometheus.Registry) (*GrpcServer, error) {
	if config.GrpcPort == "" {
		return nil, fmt.Errorf("grpc port is required")
	}

	return &GrpcServer{
		config: config,
		logger: logger.With("component", "ControllerGrpcServer"),
	}, nil
}

func (s *GrpcServer) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%s", s.config.GrpcPort)
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

	// Register services
	reflection.Register(s.server)
	pb.RegisterControllerServer(s.server, s)
	
	// Register health check
	name := pb.Controller_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, s.server)

	s.logger.Info("gRPC server listening", "address", listener.Addr().String())
	
	// This blocks until the server is stopped
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

func (s *GrpcServer) AuthorizePayment(ctx context.Context, req *pb.AuthorizePaymentRequest) (*pb.AuthorizePaymentReply, error) {
	// TODO: Implement actual payment authorization logic
	// For now, always return success
	s.logger.Debug("AuthorizePayment called", "request", req)
	
	return &pb.AuthorizePaymentReply{}, nil
}