package apiserver

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	healthcheck "github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/disperser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var errSystemBlobRateLimit = fmt.Errorf("request ratelimited: system blob limit")
var errSystemThroughputRateLimit = fmt.Errorf("request ratelimited: system throughput limit")
var errAccountBlobRateLimit = fmt.Errorf("request ratelimited: account blob limit")
var errAccountThroughputRateLimit = fmt.Errorf("request ratelimited: account throughput limit")

const systemAccountKey = "system"

const maxBlobSize = 2 * 1024 * 1024 // 2 MiB

type DispersalServer struct {
	pb.UnimplementedDisperserServer
	mu *sync.Mutex

	config disperser.ServerConfig

	blobStore   disperser.BlobStore
	tx          core.Transactor
	quorumCount uint8

	rateConfig    RateConfig
	ratelimiter   common.RateLimiter
	authenticator core.BlobRequestAuthenticator

	metrics *disperser.Metrics

	logger common.Logger
}

// NewServer creates a new Server struct with the provided parameters.
//
// Note: The Server's chunks store will be created at config.DbPath+"/chunk".
func NewDispersalServer(
	config disperser.ServerConfig,
	store disperser.BlobStore,
	tx core.Transactor,
	logger common.Logger,
	metrics *disperser.Metrics,
	ratelimiter common.RateLimiter,
	rateConfig RateConfig,
) *DispersalServer {
	for ip, rateInfoByQuorum := range rateConfig.Allowlist {
		for quorumID, rateInfo := range rateInfoByQuorum {
			logger.Info("[Allowlist]", "ip", ip, "quorumID", quorumID, "throughput", rateInfo.Throughput, "blobRate", rateInfo.BlobRate)
		}
	}

	authenticator := auth.NewAuthenticator(auth.AuthConfig{})

	return &DispersalServer{
		config:        config,
		blobStore:     store,
		tx:            tx,
		quorumCount:   0,
		metrics:       metrics,
		logger:        logger,
		ratelimiter:   ratelimiter,
		authenticator: authenticator,
		rateConfig:    rateConfig,
		mu:            &sync.Mutex{},
	}
}

func (s *DispersalServer) Start(ctx context.Context) error {
	s.logger.Trace("Entering Start function...")
	defer s.logger.Trace("Exiting Start function...")

	// Serve grpc requests
	addr := fmt.Sprintf("%s:%s", disperser.Localhost, s.config.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener")
	}

	opt := grpc.MaxRecvMsgSize(1024 * 1024 * 300) // 300 MiB
	gs := grpc.NewServer(opt)
	reflection.Register(gs)
	pb.RegisterDisperserServer(gs, s)

	// Register Server for Health Checks
	name := pb.Disperser_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, gs)

	s.logger.Info("port", s.config.GrpcPort, "address", listener.Addr().String(), "GRPC Listening")
	if err := gs.Serve(listener); err != nil {
		return fmt.Errorf("could not start GRPC server")
	}

	return nil
}

func (s *DispersalServer) updateQuorumCount(ctx context.Context) error {
	currentBlock, err := s.tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}
	count, err := s.tx.GetQuorumCount(ctx, currentBlock)
	if err != nil {
		return err
	}

	s.logger.Debug("updating quorum count", "currentBlock", currentBlock, "count", count)
	s.mu.Lock()
	s.quorumCount = count
	s.mu.Unlock()
	return nil
}
