package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DispersalServerV2 struct {
	pb.UnimplementedDisperserServer

	serverConfig disperser.ServerConfig
	logger       logging.Logger
}

// NewDispersalServerV2 creates a new Server struct with the provided parameters.
func NewDispersalServerV2(
	serverConfig disperser.ServerConfig,
	_logger logging.Logger,
) *DispersalServerV2 {
	logger := _logger.With("component", "DispersalServerV2")

	return &DispersalServerV2{
		serverConfig: serverConfig,
		logger:       logger,
	}
}

func (s *DispersalServerV2) DisperseBlob(ctx context.Context, req *pb.DisperseBlobRequest) (*pb.DisperseBlobReply, error) {
	return &pb.DisperseBlobReply{}, api.NewErrorUnimplemented()
}

func (s *DispersalServerV2) GetBlobStatus(ctx context.Context, req *pb.BlobStatusRequest) (*pb.BlobStatusReply, error) {
	return &pb.BlobStatusReply{}, api.NewErrorUnimplemented()
}

func (s *DispersalServerV2) GetBlobCommitment(ctx context.Context, req *pb.BlobCommitmentRequest) (*pb.BlobCommitmentReply, error) {
	return &pb.BlobCommitmentReply{}, api.NewErrorUnimplemented()
}

func (s *DispersalServerV2) Start(ctx context.Context) error {
	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.serverConfig.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.New("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB

	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	s.logger.Info("GRPC Listening", "port", s.serverConfig.GrpcPort, "address", listener.Addr().String())

	if err := gs.Serve(listener); err != nil {
		return errors.New("could not start GRPC server")
	}

	return nil
}
