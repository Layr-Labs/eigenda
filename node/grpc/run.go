package grpc

import (
	"errors"
	"fmt"
	"net"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const localhost = "0.0.0.0"

func RunServers(serverV1 *Server, serverV2 *ServerV2, config *node.Config, logger logging.Logger) error {
	if serverV1 == nil {
		return errors.New("node V1 server is not configured")
	}
	if serverV2 == nil {
		return errors.New("node V2 server is not configured")
	}

	go func() {
		for {
			addr := fmt.Sprintf("%s:%s", localhost, config.InternalDispersalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Fatalf("Could not start tcp listener: %v", err)
			}

			opt := grpc.MaxRecvMsgSize(60 * 1024 * 1024 * 1024) // 60 GiB
			gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			pb.RegisterDispersalServer(gs, serverV1)
			pbv2.RegisterDispersalServer(gs, serverV2)

			healthcheck.RegisterHealthServer("node.Dispersal", gs)

			logger.Info("port", config.InternalDispersalPort, "address", listener.Addr().String(), "GRPC Listening")
			if err := gs.Serve(listener); err != nil {
				logger.Error("dispersal server failed; restarting.", "err", err)
			}
		}
	}()

	go func() {
		for {
			addr := fmt.Sprintf("%s:%s", localhost, config.InternalRetrievalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Fatalf("Could not start tcp listener: %v", err)
			}

			opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
			gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			pb.RegisterRetrievalServer(gs, serverV1)
			pbv2.RegisterRetrievalServer(gs, serverV2)
			healthcheck.RegisterHealthServer("node.Retrieval", gs)

			logger.Info("port", config.InternalRetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
			if err := gs.Serve(listener); err != nil {
				logger.Error("retrieval server failed; restarting.", "err", err)
			}
		}
	}()

	return nil
}
