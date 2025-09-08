package metering

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigensdk-go/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedControllerServiceServer

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
	pb.RegisterControllerServiceServer(s.server, s)

	// Register health check
	name := pb.ControllerService_ServiceDesc.ServiceName
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

func (s *GrpcServer) AuthorizePayment(
	ctx context.Context,
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	// TODO: Implement actual payment authorization logic

	// Example: Simulate insufficient balance error with structured metadata
	// In real implementation, this would be based on actual balance checks
	simulateInsufficientBalance := false

	if simulateInsufficientBalance {
		accountID := request.GetBlobHeader().GetPaymentHeader().GetAccountId()
		currentBalance := uint64(100) // In production, get actual balance
		requiredCost := uint64(150)   // In production, calculate actual cost

		return nil, s.newInsufficientBalanceError(accountID, currentBalance, requiredCost)
	}

	return &pb.AuthorizePaymentReply{}, nil
}

// newInsufficientBalanceError creates a structured gRPC error for insufficient balance
// with detailed metadata about the account balance and required cost
func (s *GrpcServer) newInsufficientBalanceError(accountID string, currentBalance, requiredCost uint64) error {
	deficit := uint64(0)
	if requiredCost > currentBalance {
		deficit = requiredCost - currentBalance
	}

	st := status.New(codes.FailedPrecondition, "insufficient balance for blob dispersal")

	// Add structured error details with metadata
	// TODO: make this match the same metadata that is returned from the structure payments error. Maybe even make
	// a method on the structure payments error, which knows how to wrap it up as a GRPC error???
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: "INSUFFICIENT_BALANCE",
		Domain: "payment",
		Metadata: map[string]string{
			"account_id":      accountID,
			"current_balance": fmt.Sprintf("%d", currentBalance),
			"required_cost":   fmt.Sprintf("%d", requiredCost),
			"deficit":         fmt.Sprintf("%d", deficit),
		},
	})

	if err != nil {
		// If we can't add details, return the basic error
		s.logger.Error("failed to add error details", "error", err)
		return st.Err()
	}

	return st.Err()
}
