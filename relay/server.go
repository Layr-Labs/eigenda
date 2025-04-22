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
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/common/tracing"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/auth"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/relay/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

	// replayGuardian is used to guard against replay attacks.
	replayGuardian replay.ReplayGuardian

	// chainReader is the core.Reader used to fetch blob parameters.
	chainReader core.Reader

	// metrics encapsulates the metrics for the relay server.
	metrics *metrics.RelayMetrics

	// telemetryShutdown is the function to call to shutdown telemetry
	telemetryShutdown func(context.Context) error
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
		authenticator, err = auth.NewRequestAuthenticator(ctx, ics, config.AuthenticationKeyCacheSize)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}

	replayGuardian := replay.NewReplayGuardian(
		time.Now,
		config.GetChunksRequestMaxPastAge,
		config.GetChunksRequestMaxPastAge)

	return &Server{
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
	}, nil
}

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {
	ctx, span := tracing.TraceOperation(ctx, "RelayServer.GetBlob")
	defer span.End()

	start := time.Now()

	if s.config.Timeouts.GetBlobTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetBlobTimeout)
		defer cancel()
	}

	// Validate the request params before any further processing (as validation is cheaper)
	key, err := v2.BytesToBlobKey(request.BlobKey)
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid blob key: %v", err))
	}
	s.logger.Debug("GetBlob request received", "key", key.Hex())

	err = s.blobRateLimiter.BeginGetBlobOperation(time.Now())
	if err != nil {
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("rate limit exceeded: %v", err))
	}
	defer s.blobRateLimiter.FinishGetBlobOperation()

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

func (s *Server) validateGetChunksRequest(request *pb.GetChunksRequest) error {
	if request == nil {
		return api.NewErrorInvalidArg("request is nil")
	}
	if len(request.ChunkRequests) == 0 {
		return api.NewErrorInvalidArg("no chunk requests provided")
	}
	if len(request.ChunkRequests) > s.config.MaxKeysPerGetChunksRequest {
		return api.NewErrorInvalidArg(fmt.Sprintf(
			"too many chunk requests provided, max is %d", s.config.MaxKeysPerGetChunksRequest))
	}

	for _, chunkRequest := range request.ChunkRequests {
		if chunkRequest.GetByIndex() == nil && chunkRequest.GetByRange() == nil {
			return api.NewErrorInvalidArg("chunk request must be either by index or by range")
		}
	}

	return nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	ctx, span := tracing.TraceOperation(ctx, "RelayServer.GetChunks")
	defer span.End()

	start := time.Now()

	if s.config.Timeouts.GetChunksTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetChunksTimeout)
		defer cancel()
	}

	// Validation span
	_, validateSpan := tracing.TraceOperation(ctx, "GetChunks.Validate")
	err := s.validateGetChunksRequest(request)
	validateSpan.End()
	if err != nil {
		return nil, err
	}

	s.metrics.ReportChunkKeyCount(len(request.ChunkRequests))

	// Authentication span
	_, authSpan := tracing.TraceOperation(ctx, "GetChunks.Authenticate")
	if s.authenticator != nil {
		client, ok := peer.FromContext(ctx)
		if !ok {
			authSpan.End()
			return nil, api.NewErrorInvalidArg("could not get peer information")
		}
		clientAddress := client.Addr.String()

		hash, err := s.authenticator.AuthenticateGetChunksRequest(ctx, request)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			s.logger.Debug("rejected GetChunks request", "client", clientAddress)
			authSpan.End()
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("auth failed: %v", err))
		}

		timestamp := time.Unix(int64(request.Timestamp), 0)
		err = s.replayGuardian.VerifyRequest(hash, timestamp)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			authSpan.End()
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request: %v", err))
		}

		s.logger.Debug("received authenticated GetChunks request", "client", clientAddress)
	}
	authSpan.End()

	finishedAuthenticating := time.Now()
	if s.authenticator != nil {
		s.metrics.ReportChunkAuthenticationLatency(finishedAuthenticating.Sub(start))
	}

	// Rate limiting span
	_, rateLimitSpan := tracing.TraceOperation(ctx, "GetChunks.RateLimit")
	clientID := string(request.OperatorId)
	err = s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
	rateLimitSpan.End()
	if err != nil {
		return nil, api.NewErrorResourceExhausted(fmt.Sprintf("rate limit exceeded: %v", err))
	}
	defer s.chunkRateLimiter.FinishGetChunkOperation(clientID)

	// Key extraction span
	_, keySpan := tracing.TraceOperation(ctx, "GetChunks.ExtractKeys")
	keys, err := getKeysFromChunkRequest(request)
	keySpan.End()
	if err != nil {
		return nil, api.NewErrorInvalidArg(fmt.Sprintf("invalid request: %v", err))
	}

	// Metadata fetch span
	_, metadataSpan := tracing.TraceOperation(ctx, "GetChunks.FetchMetadata")
	mMap, err := s.metadataProvider.GetMetadataForBlobs(ctx, keys)
	metadataSpan.End()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %v", err))
	}

	finishedFetchingMetadata := time.Now()
	s.metrics.ReportChunkMetadataLatency(finishedFetchingMetadata.Sub(finishedAuthenticating))

	// Bandwidth computation span
	_, bandwidthSpan := tracing.TraceOperation(ctx, "GetChunks.ComputeBandwidth")
	requiredBandwidth, err := computeChunkRequestRequiredBandwidth(request, mMap)
	bandwidthSpan.End()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error computing required bandwidth: %v", err))
	}
	s.metrics.ReportGetChunksRequestedBandwidthUsage(requiredBandwidth)

	// Bandwidth rate limiting span
	_, bandwidthLimitSpan := tracing.TraceOperation(ctx, "GetChunks.BandwidthRateLimit")
	err = s.chunkRateLimiter.RequestGetChunkBandwidth(time.Now(), clientID, requiredBandwidth)
	bandwidthLimitSpan.End()
	if err != nil {
		if strings.Contains(err.Error(), "internal error") {
			return nil, api.NewErrorInternal(err.Error())
		}
		return nil, buildInsufficientGetChunksBandwidthError(request, requiredBandwidth, err)
	}
	s.metrics.ReportGetChunksBandwidthUsage(requiredBandwidth)

	// Frame fetching span
	_, frameSpan := tracing.TraceOperation(ctx, "GetChunks.FetchFrames")
	frames, err := s.chunkProvider.GetFrames(ctx, mMap)
	frameSpan.End()
	if err != nil {
		return nil, api.NewErrorInternal(fmt.Sprintf("error fetching frames: %v", err))
	}

	// Response processing span
	_, processSpan := tracing.TraceOperation(ctx, "GetChunks.ProcessResponse")
	bytesToSend, err := gatherChunkDataToSend(frames, request)
	processSpan.End()
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
	frames map[v2.BlobKey]*core.ChunksData,
	request *pb.GetChunksRequest) ([][]byte, error) {

	bytesToSend := make([][]byte, 0, len(request.ChunkRequests))

	for _, chunkRequest := range request.ChunkRequests {
		var framesSubset *core.ChunksData
		var err error

		if chunkRequest.GetByIndex() != nil {
			framesSubset, err = selectFrameSubsetByIndex(chunkRequest.GetByIndex(), frames)
		} else {
			framesSubset, err = selectFrameSubsetByRange(chunkRequest.GetByRange(), frames)
		}

		if err != nil {
			return nil, fmt.Errorf("error selecting frame subset: %v", err)
		}

		subsetBytes, err := framesSubset.FlattenToBundle()
		if err != nil {
			return nil, fmt.Errorf("error serializing frame subset: %v", err)
		}

		bytesToSend = append(bytesToSend, subsetBytes)
	}

	return bytesToSend, nil
}

// selectFrameSubsetByRange selects a subset of frames from a BinaryFrames object based on a range
func selectFrameSubsetByRange(
	request *pb.ChunkRequestByRange,
	allFrames map[v2.BlobKey]*core.ChunksData) (*core.ChunksData, error) {

	key := v2.BlobKey(request.GetBlobKey())
	startIndex := request.StartIndex
	endIndex := request.EndIndex

	frames, ok := allFrames[key]
	if !ok {
		return nil, fmt.Errorf("frames not found for key %s", key.Hex())
	}

	if startIndex > endIndex {
		return nil, fmt.Errorf(
			"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
			startIndex, endIndex, key.Hex())
	}
	if endIndex > uint32(len(frames.Chunks)) {
		return nil, fmt.Errorf(
			"chunk range %d-%d is invald for key %s, chunk count %d",
			startIndex, endIndex, key, len(frames.Chunks))
	}

	framesSubset := &core.ChunksData{
		Chunks:   frames.Chunks[startIndex:endIndex],
		Format:   frames.Format,
		ChunkLen: frames.ChunkLen,
	}

	return framesSubset, nil
}

// selectFrameSubsetByIndex selects a subset of frames from a BinaryFrames object based on a list of indices
func selectFrameSubsetByIndex(
	request *pb.ChunkRequestByIndex,
	allFrames map[v2.BlobKey]*core.ChunksData) (*core.ChunksData, error) {

	key := v2.BlobKey(request.GetBlobKey())
	frames, ok := allFrames[key]
	if !ok {
		return nil, fmt.Errorf("frames not found for key %s", key.Hex())
	}

	if len(request.ChunkIndices) > len(frames.Chunks) {
		return nil, fmt.Errorf("too many requested chunks for key %s, chunk count %d",
			key.Hex(), len(frames.Chunks))
	}

	framesSubset := &core.ChunksData{
		Format:   frames.Format,
		ChunkLen: frames.ChunkLen,
		Chunks:   make([][]byte, 0, len(request.ChunkIndices)),
	}

	for index := range request.ChunkIndices {
		if index >= len(frames.Chunks) {
			return nil, fmt.Errorf(
				"chunk index %d out of range for key %s, chunk count %d",
				index, key.Hex(), len(frames.Chunks))
		}

		framesSubset.Chunks = append(framesSubset.Chunks, frames.Chunks[index])
	}

	return framesSubset, nil
}

// computeChunkRequestRequiredBandwidth computes the bandwidth required to fulfill a GetChunks request.
func computeChunkRequestRequiredBandwidth(request *pb.GetChunksRequest, mMap metadataMap) (uint32, error) {
	requiredBandwidth := uint32(0)
	for _, req := range request.ChunkRequests {
		var metadata *blobMetadata
		var key v2.BlobKey
		var requestedChunks uint32

		if req.GetByIndex() != nil {
			key = v2.BlobKey(req.GetByIndex().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = uint32(len(req.GetByIndex().ChunkIndices))
		} else {
			key = v2.BlobKey(req.GetByRange().GetBlobKey())
			metadata = mMap[key]

			if req.GetByRange().EndIndex < req.GetByRange().StartIndex {
				return 0, fmt.Errorf(
					"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
					req.GetByRange().StartIndex, req.GetByRange().EndIndex, key.Hex())
			}

			requestedChunks = req.GetByRange().EndIndex - req.GetByRange().StartIndex
		}

		if metadata == nil {
			return 0, fmt.Errorf("metadata not found for key %s", key.Hex())
		}

		requiredBandwidth += requestedChunks * metadata.chunkSizeBytes
	}

	return requiredBandwidth, nil
}

// buildInsufficientBandwidthError builds an informative error message for when there is insufficient
// bandwidth to serve a GetChunks() request.
func buildInsufficientGetChunksBandwidthError(
	request *pb.GetChunksRequest,
	requiredBandwidth uint32,
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
	// Initialize tracing if enabled
	if err := s.InitTracing(ctx); err != nil {
		s.logger.Error("Failed to initialize tracing", "err", err)
		// Continue with startup even if tracing fails
	}

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
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener on %s: %w", addr, err)
	}

	opt := grpc.MaxRecvMsgSize(s.config.MaxGRPCMessageSize)
	otelHandler := grpc.StatsHandler(otelgrpc.NewServerHandler())

	s.grpcServer = grpc.NewServer(opt, s.metrics.GetGRPCServerOption(), otelHandler)
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

	if s.config.EnableMetrics {
		err := s.metrics.Stop()
		if err != nil {
			return fmt.Errorf("error stopping metrics server: %w", err)
		}
	}

	// Shutdown tracing if enabled
	if s.telemetryShutdown != nil {
		if err := s.telemetryShutdown(context.Background()); err != nil {
			s.logger.Error("Failed to shutdown telemetry", "err", err)
		}
	}

	return nil
}

// InitTracing initializes tracing if enabled in the config
func (s *Server) InitTracing(ctx context.Context) error {
	// Initialize tracing if enabled
	if s.config.Tracing.Enabled {
		tracingCfg := tracing.TracingConfig{
			Enabled:     s.config.Tracing.Enabled,
			ServiceName: s.config.Tracing.ServiceName,
			Endpoint:    s.config.Tracing.Endpoint,
			SampleRatio: s.config.Tracing.SampleRatio,
		}

		telemetryShutdown, err := tracing.InitTelemetry(ctx, tracingCfg)
		if err != nil {
			s.logger.Error("Failed to initialize tracing", "err", err)
			// Continue with startup even if tracing fails
		} else {
			s.logger.Info("Enabled tracing for Relay Server", "endpoint", s.config.Tracing.Endpoint)
			s.telemetryShutdown = telemetryShutdown

			// Add cleanup handler for tracing when server shuts down
			go func() {
				<-ctx.Done()
				if err := telemetryShutdown(context.Background()); err != nil {
					s.logger.Error("Failed to shutdown telemetry", "err", err)
				}
			}()
		}
	}
	return nil
}
