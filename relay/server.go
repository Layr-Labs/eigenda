package relay

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/auth"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/relay/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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

	// grpcServer is the gRPC server.
	grpcServer *grpc.Server

	// authenticator is used to authenticate requests to the relay service.
	authenticator auth.RequestAuthenticator

	// chainReader is the core.Reader used to fetch blob parameters.
	chainReader core.Reader

	// metrics encapsulates the metrics for the relay server.
	metrics *metrics.RelayMetrics
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader chunkstore.ChunkReader,
	chainReader core.Reader,
	ics core.IndexedChainState,
) (*Server, error) {

	if chainReader == nil {
		return nil, errors.New("chainReader is required")
	}

	blobParams, err := chainReader.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching blob params: %w", err)
	}

	relayMetrics := metrics.NewRelayMetrics(logger, config.MetricsPort)

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
		authenticator, err = auth.NewRequestAuthenticator(
			ctx,
			ics,
			config.AuthenticationKeyCacheSize,
			config.AuthenticationTimeout)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}

	return &Server{
		config:           config,
		logger:           logger.With("component", "RelayServer"),
		metadataProvider: mp,
		blobProvider:     bp,
		chunkProvider:    cp,
		blobRateLimiter:  limiter.NewBlobRateLimiter(&config.RateLimits, relayMetrics),
		chunkRateLimiter: limiter.NewChunkRateLimiter(&config.RateLimits, relayMetrics),
		authenticator:    authenticator,
		metrics:          relayMetrics,
	}, nil
}

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {
	start := time.Now()

	if s.config.Timeouts.GetBlobTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetBlobTimeout)
		defer cancel()
	}

	err := s.blobRateLimiter.BeginGetBlobOperation(time.Now())
	if err != nil {
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("rate limit exceeded: %v", err))
	}
	defer s.blobRateLimiter.FinishGetBlobOperation()

	key, err := v2.BytesToBlobKey(request.BlobKey)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob key: %v", err))
	}

	s.logger.Debug("GetBlob request received", "key", key.Hex())

	keys := []v2.BlobKey{key}
	mMap, err := s.metadataProvider.GetMetadataForBlobs(ctx, keys)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %v", err))
	}
	metadata := mMap[v2.BlobKey(request.BlobKey)]
	if metadata == nil {
		return nil, api.NewErrorNotFound("blob not found")
	}

	finishedFetchingMetadata := time.Now()
	s.metrics.ReportBlobMetadataLatency(finishedFetchingMetadata.Sub(start))

	s.metrics.ReportBlobRequestedBandwidthUsage(int(metadata.blobSizeBytes))
	err = s.blobRateLimiter.RequestGetBlobBandwidth(time.Now(), metadata.blobSizeBytes)
	if err != nil {
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("bandwidth limit exceeded: %v", err))
	}

	data, err := s.blobProvider.GetBlob(ctx, key)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error fetching blob %s: %v", key.Hex(), err))
	}

	s.metrics.ReportBlobBandwidthUsage(len(data))
	s.metrics.ReportBlobDataLatency(time.Since(finishedFetchingMetadata))
	s.metrics.ReportBlobLatency(time.Since(start))

	reply := &pb.GetBlobReply{
		Blob: data,
	}
	return reply, nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	start := time.Now()

	if s.config.Timeouts.GetChunksTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetChunksTimeout)
		defer cancel()
	}

	if len(request.ChunkRequests) <= 0 {
		return nil, api.NewErrorInvalidArg("no chunk requests provided")
	}
	if len(request.ChunkRequests) > s.config.MaxKeysPerGetChunksRequest {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf(
			"too many chunk requests provided, max is %d", s.config.MaxKeysPerGetChunksRequest))
	}
	s.metrics.ReportChunkKeyCount(len(request.ChunkRequests))

	if s.authenticator != nil {
		client, ok := peer.FromContext(ctx)
		if !ok {
			return nil, api.NewErrorInternal("could not get peer information")
		}
		clientAddress := client.Addr.String()

		err := s.authenticator.AuthenticateGetChunksRequest(ctx, clientAddress, request, time.Now())
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			s.logger.Debug("rejected GetChunks request", "client", clientAddress)
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("auth failed: %v", err))
		}
		s.logger.Debug("received authenticated GetChunks request", "client", clientAddress)
	}

	finishedAuthenticating := time.Now()
	if s.authenticator != nil {
		s.metrics.ReportChunkAuthenticationLatency(finishedAuthenticating.Sub(start))
	}

	clientID := string(request.OperatorId)
	err := s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
	if err != nil {
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("rate limit exceeded: %v", err))
	}
	defer s.chunkRateLimiter.FinishGetChunkOperation(clientID)

	// keys might contain duplicate keys
	keys, err := getKeysFromChunkRequest(request)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid request: %v", err))
	}

	mMap, err := s.metadataProvider.GetMetadataForBlobs(ctx, keys)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %v", err))
	}

	finishedFetchingMetadata := time.Now()
	s.metrics.ReportChunkMetadataLatency(finishedFetchingMetadata.Sub(finishedAuthenticating))

	requiredBandwidth, err := computeChunkRequestRequiredBandwidth(request, mMap)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error computing required bandwidth: %v", err))
	}
	s.metrics.ReportGetChunksRequestedBandwidthUsage(requiredBandwidth)
	err = s.chunkRateLimiter.RequestGetChunkBandwidth(time.Now(), clientID, requiredBandwidth)
	if err != nil {
		if strings.Contains(err.Error(), "internal error") {
			return nil, api.NewErrorInternal(err.Error())
		}
		return nil, buildInsufficientGetChunksBandwidthError(request, requiredBandwidth, err)
	}
	s.metrics.ReportGetChunksBandwidthUsage(requiredBandwidth)

	frames, err := s.chunkProvider.GetFrames(ctx, mMap)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error fetching frames: %v", err))
	}

	bytesToSend, err := gatherChunkDataToSend(frames, request)
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error gathering chunk data: %v", err))
	}

	s.metrics.ReportChunkDataLatency(time.Since(finishedFetchingMetadata))
	s.metrics.ReportChunkLatency(time.Since(start))

	return &pb.GetChunksReply{
		Data: bytesToSend,
	}, nil
}

// getKeysFromChunkRequest gathers a slice of blob keys from a GetChunks request.
func getKeysFromChunkRequest(request *pb.GetChunksRequest) ([]v2.BlobKey, error) {
	keys := make([]v2.BlobKey, 0, len(request.ChunkRequests))

	for _, chunkRequest := range request.ChunkRequests {
		var key v2.BlobKey
		if chunkRequest.GetByIndex() != nil {
			var err error
			key, err = v2.BytesToBlobKey(chunkRequest.GetByIndex().GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("invalid blob key: %w", err)
			}
		} else {
			var err error
			key, err = v2.BytesToBlobKey(chunkRequest.GetByRange().GetBlobKey())
			if err != nil {
				return nil, fmt.Errorf("invalid blob key: %w", err)
			}
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// gatherChunkDataToSend takes the chunk data and narrows it down to the data requested in the GetChunks request.
func gatherChunkDataToSend(
	frames map[v2.BlobKey][]*encoding.Frame,
	request *pb.GetChunksRequest) ([][]byte, error) {

	bytesToSend := make([][]byte, 0, len(request.ChunkRequests))

	for _, chunkRequest := range request.ChunkRequests {

		framesToSend := make([]*encoding.Frame, 0)

		if chunkRequest.GetByIndex() != nil {
			key := v2.BlobKey(chunkRequest.GetByIndex().GetBlobKey())
			blobFrames := (frames)[key]

			for index := range chunkRequest.GetByIndex().ChunkIndices {

				if index >= len(blobFrames) {
					return nil, fmt.Errorf(
						"chunk index %d out of range for key %s, chunk count %d",
						index, key.Hex(), len(blobFrames))
				}

				framesToSend = append(framesToSend, blobFrames[index])
			}

		} else {
			key := v2.BlobKey(chunkRequest.GetByRange().GetBlobKey())
			startIndex := chunkRequest.GetByRange().StartIndex
			endIndex := chunkRequest.GetByRange().EndIndex

			blobFrames := (frames)[key]

			if startIndex > endIndex {
				return nil, fmt.Errorf(
					"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
					startIndex, endIndex, key.Hex())
			}
			if endIndex > uint32(len((frames)[key])) {
				return nil, fmt.Errorf(
					"chunk range %d-%d is invald for key %s, chunk count %d",
					chunkRequest.GetByRange().StartIndex, chunkRequest.GetByRange().EndIndex, key, len(blobFrames))
			}

			framesToSend = append(framesToSend, blobFrames[startIndex:endIndex]...)
		}

		bundle := core.Bundle(framesToSend)
		bundleBytes, err := bundle.Serialize()
		if err != nil {
			return nil, fmt.Errorf("error serializing bundle: %w", err)
		}
		bytesToSend = append(bytesToSend, bundleBytes)
	}

	return bytesToSend, nil
}

// computeChunkRequestRequiredBandwidth computes the bandwidth required to fulfill a GetChunks request.
func computeChunkRequestRequiredBandwidth(request *pb.GetChunksRequest, mMap metadataMap) (int, error) {
	requiredBandwidth := 0
	for _, req := range request.ChunkRequests {
		var metadata *blobMetadata
		var key v2.BlobKey
		var requestedChunks int

		if req.GetByIndex() != nil {
			key = v2.BlobKey(req.GetByIndex().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = len(req.GetByIndex().ChunkIndices)
		} else {
			key = v2.BlobKey(req.GetByRange().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = int(req.GetByRange().EndIndex - req.GetByRange().StartIndex)
		}

		if metadata == nil {
			return 0, fmt.Errorf("metadata not found for key %s", key.Hex())
		}

		requiredBandwidth += requestedChunks * int(metadata.chunkSizeBytes)
	}

	return requiredBandwidth, nil
}

// buildInsufficientBandwidthError builds an informative error message for when there is insufficient
// bandwidth to serve a GetChunks() request.
func buildInsufficientGetChunksBandwidthError(
	request *pb.GetChunksRequest,
	requiredBandwidth int,
	originalError error) error {

	chunkCount := 0
	for _, chunkRequest := range request.ChunkRequests {
		if chunkRequest.GetByIndex() != nil {
			chunkCount += len(chunkRequest.GetByIndex().ChunkIndices)
		} else {
			chunkCount += int(chunkRequest.GetByRange().EndIndex - chunkRequest.GetByRange().StartIndex)
		}
	}

	blobCount := len(request.ChunkRequests)

	return api.NewErrorResourceExhausted(fmt.Sprintf("unable to serve data (%d blobs, %d chunks, %d bytes): %v",
		blobCount, chunkCount, requiredBandwidth, originalError))
}

// Start starts the server listening for requests. This method will block until the server is stopped.
func (s *Server) Start(ctx context.Context) error {
	s.metrics.Start()

	if s.chainReader != nil && s.metadataProvider != nil {
		go func() {
			_ = s.RefreshOnchainState(ctx)
		}()
	}

	// Serve grpc requests
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener on %s: %w", addr, err)
	}

	opt := grpc.MaxRecvMsgSize(s.config.MaxGRPCMessageSize)

	s.grpcServer = grpc.NewServer(opt, s.metrics.GetGRPCServerOption())
	reflection.Register(s.grpcServer)
	pb.RegisterRelayServer(s.grpcServer, s)

	// Register Server for Health Checks
	name := pb.Relay_ServiceDesc.ServiceName
	healthcheck.RegisterHealthServer(name, s.grpcServer)

	s.logger.Info("GRPC Listening", "port", s.config.GRPCPort, "address", listener.Addr().String())
	if err = s.grpcServer.Serve(listener); err != nil {
		return errors.New("could not start GRPC server")
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

	err := s.metrics.Stop()
	if err != nil {
		return fmt.Errorf("error stopping metrics server: %w", err)
	}

	return nil
}
