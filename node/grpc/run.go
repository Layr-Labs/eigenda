package grpc

import (
	"errors"
	"fmt"
	"net"
	"sync"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const localhost = "0.0.0.0"

// ServerRunner manages the lifecycle of gRPC servers
type ServerRunner struct {
	servers []*grpc.Server
	done    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	logger  logging.Logger
}

// Stop gracefully shuts down all running servers
func (r *ServerRunner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Signal all goroutines to stop
	close(r.done)

	// Gracefully stop all servers
	for _, s := range r.servers {
		s.GracefulStop()
	}

	// Wait for all goroutines to finish
	r.wg.Wait()
	r.logger.Info("All gRPC servers stopped")
}

func RunServers(
	serverV1 *Server, serverV2 *ServerV2, config *node.Config, logger logging.Logger,
) (*ServerRunner, error) {
	if config.EnableV1 && serverV1 == nil {
		return nil, errors.New("node v1 server is not configured")
	}
	if config.EnableV2 && serverV2 == nil {
		return nil, errors.New("node v2 server is not configured")
	}
	if !config.EnableV1 && !config.EnableV2 {
		return nil, errors.New("node is not configured to run any servers")
	}

	runner := &ServerRunner{
		servers: make([]*grpc.Server, 0, 4),
		done:    make(chan struct{}),
		logger:  logger,
	}

	// V1 dispersal service
	if config.EnableV1 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()

			addr := fmt.Sprintf("%s:%s", localhost, config.InternalDispersalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Errorf("Could not start tcp listener for v1 dispersal: %v", err)
				return
			}

			opt := grpc.MaxRecvMsgSize(60 * 1024 * 1024 * 1024) // 60 GiB
			gs := grpc.NewServer(opt)

			runner.mu.Lock()
			runner.servers = append(runner.servers, gs)
			runner.mu.Unlock()

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			pb.RegisterDispersalServer(gs, serverV1)

			healthcheck.RegisterHealthServer("node.Dispersal", gs)

			logger.Info("v1 dispersal enabled on port", config.InternalDispersalPort, "address", listener.Addr().String(), "GRPC Listening")

			// Run server in background
			go func() {
				if err := gs.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					logger.Error("dispersal server failed", "err", err)
				}
			}()

			// Wait for shutdown signal
			<-runner.done
			logger.Info("Shutting down v1 dispersal server")
		}()
	}

	// V2 dispersal service
	if config.EnableV2 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()

			addr := fmt.Sprintf("%s:%s", localhost, config.InternalV2DispersalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Errorf("Could not start tcp listener for v2 dispersal: %v", err)
				return
			}

			opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
			gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

			runner.mu.Lock()
			runner.servers = append(runner.servers, gs)
			runner.mu.Unlock()

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			validator.RegisterDispersalServer(gs, serverV2)

			healthcheck.RegisterHealthServer("node.v2.Dispersal", gs)

			logger.Info("v2 dispersal enabled on port", config.InternalV2DispersalPort, "address", listener.Addr().String(), "GRPC Listening")

			// Run server in background
			go func() {
				if err := gs.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					logger.Error("dispersal v2 server failed", "err", err)
				}
			}()

			// Wait for shutdown signal
			<-runner.done
			logger.Info("Shutting down v2 dispersal server")
		}()
	}

	// v1 Retrieval service
	if config.EnableV1 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()

			addr := fmt.Sprintf("%s:%s", localhost, config.InternalRetrievalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Errorf("Could not start tcp listener for v1 retrieval: %v", err)
				return
			}

			opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
			gs := grpc.NewServer(opt)

			runner.mu.Lock()
			runner.servers = append(runner.servers, gs)
			runner.mu.Unlock()

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			pb.RegisterRetrievalServer(gs, serverV1)
			healthcheck.RegisterHealthServer("node.Retrieval", gs)

			logger.Info("v1 retrieval enabled on port", config.InternalRetrievalPort, "address", listener.Addr().String(), "GRPC Listening")

			// Run server in background
			go func() {
				if err := gs.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					logger.Error("retrieval server failed", "err", err)
				}
			}()

			// Wait for shutdown signal
			<-runner.done
			logger.Info("Shutting down v1 retrieval server")
		}()
	}

	// v2 Retrieval service
	if config.EnableV2 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()

			addr := fmt.Sprintf("%s:%s", localhost, config.InternalV2RetrievalPort)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				logger.Errorf("Could not start tcp listener for v2 retrieval: %v", err)
				return
			}

			opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
			gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

			runner.mu.Lock()
			runner.servers = append(runner.servers, gs)
			runner.mu.Unlock()

			// Register reflection service on gRPC server
			// This makes "grpcurl -plaintext localhost:9000 list" command work
			reflection.Register(gs)

			validator.RegisterRetrievalServer(gs, serverV2)

			healthcheck.RegisterHealthServer("node.v2.Retrieval", gs)

			logger.Info("v2 retrieval enabled on port", config.InternalV2RetrievalPort, "address", listener.Addr().String(), "GRPC Listening")

			// Run server in background
			go func() {
				if err := gs.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
					logger.Error("retrieval v2 server failed", "err", err)
				}
			}()

			// Wait for shutdown signal
			<-runner.done
			logger.Info("Shutting down v2 retrieval server")
		}()
	}

	return runner, nil
}
