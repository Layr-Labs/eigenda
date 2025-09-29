package grpc

import (
	"context"
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

// Runner manages the lifecycle of gRPC servers
type Runner struct {
	servers []*grpc.Server
	logger  logging.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
}

// Stop gracefully shuts down all gRPC servers
func (r *Runner) Stop() {
	r.logger.Info("Initiating graceful shutdown of gRPC servers...")

	// Cancel the context to stop accepting new connections
	r.cancel()

	// Gracefully stop all servers
	r.mu.Lock()
	for _, server := range r.servers {
		// GracefulStop waits for all RPCs to finish
		go server.GracefulStop()
	}
	r.mu.Unlock()

	// Wait for all servers to stop
	r.wg.Wait()

	r.logger.Info("All gRPC servers stopped successfully")
}

func RunServers(serverV1 *Server, serverV2 *ServerV2, config *node.Config, logger logging.Logger) (*Runner, error) {
	if config.EnableV1 && serverV1 == nil {
		return nil, errors.New("node v1 server is not configured")
	}
	if config.EnableV2 && serverV2 == nil {
		return nil, errors.New("node v2 server is not configured")
	}
	if !config.EnableV1 && !config.EnableV2 {
		return nil, errors.New("node is not configured to run any servers")
	}

	ctx, cancel := context.WithCancel(context.Background())
	runner := &Runner{
		servers: make([]*grpc.Server, 0, 4), // up to 4 servers (v1/v2 dispersal and retrieval)
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}

	// V1 dispersal service
	if config.EnableV1 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()
			for {
				select {
				case <-runner.ctx.Done():
					logger.Info("v1 dispersal server stopping")
					return
				default:
				}

				addr := fmt.Sprintf("%s:%s", localhost, config.InternalDispersalPort)
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					logger.Fatalf("Could not start tcp listener: %v", err)
				}

				opt := grpc.MaxRecvMsgSize(60 * 1024 * 1024 * 1024) // 60 GiB
				gs := grpc.NewServer(opt)

				// Track the server
				runner.mu.Lock()
				runner.servers = append(runner.servers, gs)
				runner.mu.Unlock()

				// Register reflection service on gRPC server
				// This makes "grpcurl -plaintext localhost:9000 list" command work
				reflection.Register(gs)

				pb.RegisterDispersalServer(gs, serverV1)

				healthcheck.RegisterHealthServer("node.Dispersal", gs)

				logger.Info("v1 dispersal enabled on port", config.InternalDispersalPort, "address", listener.Addr().String(), "GRPC Listening")
				if err := gs.Serve(listener); err != nil && err != grpc.ErrServerStopped {
					select {
					case <-runner.ctx.Done():
						logger.Info("v1 dispersal server stopping due to context cancellation")
						return
					default:
						logger.Error("dispersal server failed; restarting.", "err", err)
					}
				}
			}
		}()
	}

	// V2 dispersal service
	if config.EnableV2 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()
			for {
				select {
				case <-runner.ctx.Done():
					logger.Info("v2 dispersal server stopping")
					return
				default:
				}

				addr := fmt.Sprintf("%s:%s", localhost, config.InternalV2DispersalPort)
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					logger.Fatalf("Could not start tcp listener: %v", err)
				}

				opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
				gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

				// Track the server
				runner.mu.Lock()
				runner.servers = append(runner.servers, gs)
				runner.mu.Unlock()

				// Register reflection service on gRPC server
				// This makes "grpcurl -plaintext localhost:9000 list" command work
				reflection.Register(gs)

				validator.RegisterDispersalServer(gs, serverV2)

				healthcheck.RegisterHealthServer("node.v2.Dispersal", gs)

				logger.Info("v2 dispersal enabled on port", config.InternalV2DispersalPort, "address", listener.Addr().String(), "GRPC Listening")
				if err := gs.Serve(listener); err != nil && err != grpc.ErrServerStopped {
					select {
					case <-runner.ctx.Done():
						logger.Info("v2 dispersal server stopping due to context cancellation")
						return
					default:
						logger.Error("dispersal v2 server failed; restarting.", "err", err)
					}
				}
			}
		}()
	}

	// v1 Retrieval service
	if config.EnableV1 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()
			for {
				select {
				case <-runner.ctx.Done():
					logger.Info("v1 retrieval server stopping")
					return
				default:
				}

				addr := fmt.Sprintf("%s:%s", localhost, config.InternalRetrievalPort)
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					logger.Fatalf("Could not start tcp listener: %v", err)
				}

				opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
				gs := grpc.NewServer(opt)

				// Track the server
				runner.mu.Lock()
				runner.servers = append(runner.servers, gs)
				runner.mu.Unlock()

				// Register reflection service on gRPC server
				// This makes "grpcurl -plaintext localhost:9000 list" command work
				reflection.Register(gs)

				pb.RegisterRetrievalServer(gs, serverV1)
				healthcheck.RegisterHealthServer("node.Retrieval", gs)

				logger.Info("v1 retrieval enabled on port", config.InternalRetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
				if err := gs.Serve(listener); err != nil && err != grpc.ErrServerStopped {
					select {
					case <-runner.ctx.Done():
						logger.Info("v1 retrieval server stopping due to context cancellation")
						return
					default:
						logger.Error("retrieval server failed; restarting.", "err", err)
					}
				}
			}
		}()
	}

	// v2 Retrieval service
	if config.EnableV2 {
		runner.wg.Add(1)
		go func() {
			defer runner.wg.Done()
			for {
				select {
				case <-runner.ctx.Done():
					logger.Info("v2 retrieval server stopping")
					return
				default:
				}

				addr := fmt.Sprintf("%s:%s", localhost, config.InternalV2RetrievalPort)
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					logger.Fatalf("Could not start tcp listener: %v", err)
				}
				opt := grpc.MaxRecvMsgSize(config.GRPCMsgSizeLimitV2)
				gs := grpc.NewServer(opt, serverV2.metrics.GetGRPCServerOption())

				// Track the server
				runner.mu.Lock()
				runner.servers = append(runner.servers, gs)
				runner.mu.Unlock()

				// Register reflection service on gRPC server
				// This makes "grpcurl -plaintext localhost:9000 list" command work
				reflection.Register(gs)

				validator.RegisterRetrievalServer(gs, serverV2)

				healthcheck.RegisterHealthServer("node.v2.Retrieval", gs)

				logger.Info("v2 retrieval enabled on port", config.InternalV2RetrievalPort, "address", listener.Addr().String(), "GRPC Listening")
				if err := gs.Serve(listener); err != nil && err != grpc.ErrServerStopped {
					select {
					case <-runner.ctx.Done():
						logger.Info("v2 retrieval server stopping due to context cancellation")
						return
					default:
						logger.Error("retrieval v2 server failed; restarting.", "err", err)
					}
				}
			}
		}()
	}

	return runner, nil
}
