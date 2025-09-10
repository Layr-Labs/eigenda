package grpcserver

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller/v1"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/disperser/controller/metrics"
	"github.com/Layr-Labs/eigenda/disperser/controller/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// The controller GRPC server
type Server struct {
	pb.UnimplementedControllerServiceServer

	config                      Config
	logger                      logging.Logger
	server                      *grpc.Server
	listener                    net.Listener
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler
	metrics                     *metrics.ServerMetrics
	replayGuardian              replay.ReplayGuardian
	disperserPublicKey          *ecdsa.PublicKey
}

func NewServer(
	ctx context.Context,
	config Config,
	awsClientConfig aws.ClientConfig,
	disperserKMSKeyID string,
	logger logging.Logger,
	metricsRegistry *prometheus.Registry,
	paymentAuthorizationHandler *payments.PaymentAuthorizationHandler,
) (*Server, error) {
	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.AuthorizationRequestMaxPastAge,
		config.AuthorizationRequestMaxFutureAge)

	var kmsClient *kms.Client
	if awsClientConfig.EndpointURL != "" {
		endpoint := awsClientConfig.EndpointURL
		kmsClient = kms.New(kms.Options{
			Region:       awsClientConfig.Region,
			BaseEndpoint: &endpoint,
		})
	} else {
		awsConfig, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(awsClientConfig.Region))
		if err != nil {
			return nil, fmt.Errorf("load AWS config: %w", err)
		}
		kmsClient = kms.NewFromConfig(awsConfig)
	}

	disperserPublicKey, err := aws.LoadPublicKeyKMS(ctx, kmsClient, disperserKMSKeyID)
	if err != nil {
		return nil, fmt.Errorf("get disperser public key from KMS: %w", err)
	}

	return &Server{
		config:                      config,
		logger:                      logger,
		metrics:                     metrics.NewServerMetrics(metricsRegistry, logger),
		paymentAuthorizationHandler: paymentAuthorizationHandler,
		replayGuardian:              replayGuardian,
		disperserPublicKey:          disperserPublicKey,
	}, nil
}

// Start the server. Blocks until the server is stopped.
func (s *Server) Start() error {
	if !s.config.EnableServer {
		return fmt.Errorf("controller gRPC server is disabled")
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
	request *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentResponse, error) {
	start := time.Now()

	if s.paymentAuthorizationHandler == nil {
		s.metrics.ReportAuthorizePaymentAuthFailure()
		//nolint:wrapcheck
		return nil, status.Error(codes.FailedPrecondition, "payment authorization handler not configured")
	}

	requestHash, err := hashing.HashAuthorizePaymentRequest(request)
	if err != nil {
		s.metrics.ReportAuthorizePaymentAuthFailure()
		return nil, status.Errorf(codes.Internal, "failed to hash request: %v", err)
	}

	err = s.verifyDisperserSignature(requestHash, request.GetDisperserSignature())
	if err != nil {
		s.metrics.ReportAuthorizePaymentSignatureFailure()
		return nil, err
	}
	s.metrics.ReportAuthorizePaymentSignatureLatency(time.Since(start))

	timestamp := time.Unix(0, request.GetBlobHeader().GetPaymentHeader().GetTimestamp())
	err = s.replayGuardian.VerifyRequest(requestHash, timestamp)
	if err != nil {
		s.metrics.ReportAuthorizePaymentAuthFailure()
		return nil, status.Errorf(codes.InvalidArgument, "replay protection check failed: %v", err)
	}

	reply, err := s.paymentAuthorizationHandler.AuthorizePayment(ctx, request.GetBlobHeader())
	if err != nil {
		//nolint:wrapcheck // payment handler returns properly formatted grpc errors
		return nil, err
	}
	s.metrics.ReportAuthorizePaymentLatency(time.Since(start))

	return reply, nil
}

// Verifies the disperser's signature on a request
func (s *Server) verifyDisperserSignature(requestHash []byte, signature []byte) error {
	if len(signature) != 65 {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature length %d", len(signature))
	}

	// Remove the recovery ID (last byte) for verification
	valid := crypto.VerifySignature(crypto.FromECDSAPub(s.disperserPublicKey), requestHash, signature[:64])
	if !valid {
		return status.Errorf(codes.Unauthenticated, "invalid disperser signature: %v", signature)
	}

	return nil
}
