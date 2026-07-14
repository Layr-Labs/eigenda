package grpc

import (
	"errors"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/grpc/middleware"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServers(serverV2 *ServerV2, config *node.Config, logger logging.Logger) error {
	if serverV2 == nil {
		return errors.New("node v2 server is not configured")
	}

	// V2 dispersal service
	go func() {

		listener := serverV2.dispersalListener

		opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				serverV2.metrics.GetGRPCUnaryInterceptor(),
				middleware.StoreChunksDisperserAuthAndRateLimitInterceptor(serverV2.rateLimiter, serverV2.chunkAuthenticator),
			),
		)

		// Register reflection service on gRPC server
		// This makes "grpcurl -plaintext localhost:9000 list" command work
		reflection.Register(gs)

		validator.RegisterDispersalServer(gs, serverV2)

		healthcheck.RegisterHealthServer("node.v2.Dispersal", gs)

		logger.Info("v2 dispersal enabled on port", config.InternalV2DispersalPort, "address", listener.Addr().String(), "GRPC Listening")
		if err := gs.Serve(listener); err != nil {
			logger.Error("dispersal v2 server failed", "err", err)
		}
	}()

	// v2 Retrieval service
	go func() {

		listener := serverV2.retrievalListener

		opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				serverV2.metrics.GetGRPCUnaryInterceptor(),
				middleware.StoreChunksDisperserAuthAndRateLimitInterceptor(serverV2.rateLimiter, serverV2.chunkAuthenticator),
			),
		)

		// Register reflection service on gRPC server
		// This makes "grpcurl -plaintext localhost:9000 list" command work
		reflection.Register(gs)

		validator.RegisterRetrievalServer(gs, serverV2)

		healthcheck.RegisterHealthServer("node.v2.Retrieval", gs)

		logger.Info("v2 retrieval enabled on port", config.InternalV2RetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
		if err := gs.Serve(listener); err != nil {
			logger.Error("retrieval v2 server failed", "err", err)
		}
	}()

	return nil
}
