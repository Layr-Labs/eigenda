package grpc

import (
	"errors"
	"fmt"
	"net"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
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
	if config.EnableV2 && serverV2 == nil {
		return errors.New("node V2 server is not configured")
	}

	go func() {
		// V1 dispersal service
		for {
			addr := fmt.Sprintf("%s:%s", localhost, config.InternalDispersalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Fatalf("Could not start tcp listener: %v", err)
			}

			opt := grpc.MaxRecvMsgSize(60 * 1024 * 1024 * 1024) // 60 GiB
			gs := grpc.NewServer(opt)

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			pb.RegisterDispersalServer(gs, serverV1)

			healthcheck.RegisterHealthServer("node.Dispersal", gs)

			logger.Info("port", config.InternalDispersalPort, "address", listener.Addr().String(), "GRPC Listening")
			if err := gs.Serve(listener); err != nil {
				logger.Error("dispersal server failed; restarting.", "err", err)
			}
		}
	}()

	// V2 dispersal service
	go func() {
		if !config.EnableV2 {
			logger.Warn("V2 is not enabled, skipping V2 server startup")
			return
		}
		for {
			addr := fmt.Sprintf("%s:%s", localhost, config.V2DispersalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Fatalf("Could not start tcp listener: %v", err)
			}

			opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
			gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			validator.RegisterDispersalServer(gs, serverV2)

			healthcheck.RegisterHealthServer("node.v2.Dispersal", gs)

			logger.Info("port", config.V2DispersalPort, "address", listener.Addr().String(), "GRPC Listening")
			if err := gs.Serve(listener); err != nil {
				logger.Error("dispersal v2 server failed; restarting.", "err", err)
			}
		}
	}()

	// Retrieval service
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
			validator.RegisterRetrievalServer(gs, serverV2)
			healthcheck.RegisterHealthServer("node.Retrieval", gs)

			logger.Info("port", config.InternalRetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
			if err := gs.Serve(listener); err != nil {
				logger.Error("retrieval server failed; restarting.", "err", err)
			}
		}
	}()

	return nil
}
