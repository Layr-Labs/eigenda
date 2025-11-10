package relay

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/auth"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/relay/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	// config is the configuration for the relay Server.
	config *Config

	// the logger for the server
	logger logging.Logger

	// metadataProvider encapsulates logic for fetching metadata for blobs.
	metadataProvider *metadataProvider

	// blobProvider encapsulates logic for fetching blobs.
	blobProvider *blobProvider

	// chunkProvider encapsulates logic for fetching chunks.
	chunkProvider *chunkProvider

	// blobRateLimiter enforces rate limits on GetBlob and operations.
	blobRateLimiter *limiter.BlobRateLimiter

	// chunkRateLimiter enforces rate limits on GetChunk operations.
	chunkRateLimiter *limiter.ChunkRateLimiter

	// listener is the network listener for the gRPC server.
	listener net.Listener

	// grpcServer is the gRPC server.
	grpcServer *grpc.Server

	// authenticator is used to authenticate requests to the relay service.
	authenticator auth.RequestAuthenticator

	// replayGuardian is used to guard against replay attacks.
	replayGuardian replay.ReplayGuardian

	// chainReader is the core.Reader used to fetch blob parameters.
	chainReader core.Reader

	// metrics encapsulates the metrics for the relay server.
	metrics *metrics.RelayMetrics
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	metricsRegistry *prometheus.Registry,
	logger logging.Logger,
	config *Config,
	metadataStore blobstore.MetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader chunkstore.ChunkReader,
	chainReader core.Reader,
	ics core.IndexedChainState,
	listener net.Listener,
) (*Server, error) {
	if listener == nil {
		return nil, errors.New("listener is required")
	}
	if chainReader == nil {
		return nil, errors.New("chainReader is required")
	}

	blobParams, err := chainReader.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching blob params: %w", err)
	}

	relayMetrics := metrics.NewRelayMetrics(metricsRegistry, logger, config.MetricsPort)

	mp, err := newMetadataProvider(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataMaxConcurrency,
		config.RelayKeys,
		config.Timeouts.InternalGetMetadataTimeout,
		v2.NewBlobVersionParameterMap(blobParams),
		relayMetrics.MetadataCacheMetrics)

	if err != nil {
		return nil, fmt.Errorf("error creating metadata provider: %w", err)
	}

	bp, err := newBlobProvider(
		ctx,
		logger,
		blobStore,
		config.BlobCacheBytes,
		config.BlobMaxConcurrency,
		config.Timeouts.InternalGetBlobTimeout,
		relayMetrics.BlobCacheMetrics)
	if err != nil {
		return nil, fmt.Errorf("error creating blob provider: %w", err)
	}

	cp, err := newChunkProvider(
		ctx,
		logger,
		chunkReader,
		config.ChunkCacheBytes,
		config.ChunkMaxConcurrency,
		config.Timeouts.InternalGetProofsTimeout,
		config.Timeouts.InternalGetCoefficientsTimeout,
		relayMetrics.ChunkCacheMetrics)
	if err != nil {
		return nil, fmt.Errorf("error creating chunk provider: %w", err)
	}

	var authenticator auth.RequestAuthenticator
	if !config.AuthenticationDisabled {
		authenticator, err = auth.NewRequestAuthenticator(ctx, ics, config.AuthenticationKeyCacheSize)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}

	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.GetChunksRequestMaxPastAge,
		config.GetChunksRequestMaxPastAge)

	server := &Server{
		config:           config,
		logger:           logger.With("component", "RelayServer"),
		metadataProvider: mp,
		blobProvider:     bp,
		chunkProvider:    cp,
		blobRateLimiter:  limiter.NewBlobRateLimiter(&config.RateLimits, relayMetrics),
		chunkRateLimiter: limiter.NewChunkRateLimiter(&config.RateLimits, relayMetrics),
		authenticator:    authenticator,
		replayGuardian:   replayGuardian,
		metrics:          relayMetrics,
		chainReader:      chainReader,
		listener:         listener,
	}

	// Setup gRPC server
	opt := grpc.MaxRecvMsgSize(config.MaxGRPCMessageSize)
	keepAliveConfig := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     config.MaxIdleConnectionAge,
		MaxConnectionAge:      config.MaxConnectionAge,
		MaxConnectionAgeGrace: config.MaxConnectionAgeGrace,
	})

	server.grpcServer = grpc.NewServer(opt, relayMetrics.GetGRPCServerOption(), keepAliveConfig)
	reflection.Register(server.grpcServer)
	pb.RegisterRelayServer(server.grpcServer, server)

	// Register Server for Health Checks
	name := pb.Relay_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, server.grpcServer)

	return server, nil
}

// Start starts the server using the listener provided in the constructor.
// This method will block until the server is stopped.
func (s *Server) Start(ctx context.Context) error {
	// Start metrics server if enabled
	if s.config.EnableMetrics {
		s.metrics.Start()
		s.logger.Info("Enabled metrics for relay server", "port", s.config.MetricsPort)
	}

	// Start pprof server if enabled
	if s.config.EnablePprof {
		pprofProfiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", s.config.PprofHttpPort), s.logger)
		go pprofProfiler.Start()
		s.logger.Info("Enabled pprof for relay server", "port", s.config.PprofHttpPort)
	}

	if s.chainReader != nil && s.metadataProvider != nil {
		go func() {
			_ = s.RefreshOnchainState(ctx)
		}()
	}

	// Serve grpc requests
	s.logger.Info("GRPC Listening", "address", s.listener.Addr().String())
	if err := s.grpcServer.Serve(s.listener); err != nil {
		return fmt.Errorf("could not start GRPC server: %w", err)
	}

	return nil
}

func (s *Server) RefreshOnchainState(ctx context.Context) error {
	ticker := time.NewTicker(s.config.OnchainStateRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Info("refreshing onchain state")
			blobParams, err := s.chainReader.GetAllVersionedBlobParams(ctx)
			if err != nil {
				s.logger.Error("error fetching blob params", "err", err)
				continue
			}
			s.metadataProvider.UpdateBlobVersionParameters(v2.NewBlobVersionParameterMap(blobParams))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Stop stops the server.
func (s *Server) Stop() error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.config.EnableMetrics {
		err := s.metrics.Stop()
		if err != nil {
			return fmt.Errorf("error stopping metrics server: %w", err)
		}
	}

	return nil
}
