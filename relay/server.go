package relay

import (
	"bytes"
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
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/auth"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigenda/relay/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	// config is the configuration for the relay Server.
	config *RelayConfig

	// the logger for the server
	logger logging.Logger

	// metadataProvider encapsulates logic for fetching metadata for blobs.
	metadataProvider *metadataProvider

	// blobProvider encapsulates logic for fetching blobs.
	blobProvider *blobProvider

	// legacyChunkProvider encapsulates logic for fetching chunks using the old-style get by index pattern.
	legacyChunkProvider *chunkProvider

	// Provides direct access to the chunk reader client.
	chunkReader chunkstore.ChunkReader

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
	config *RelayConfig,
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

	replayGuardian, err := replay.NewReplayGuardian(
		time.Now,
		config.GetChunksRequestMaxPastAge,
		config.GetChunksRequestMaxPastAge)
	if err != nil {
		return nil, fmt.Errorf("failed to create replay guardian: %w", err)
	}

	server := &Server{
		config:              config,
		logger:              logger.With("component", "RelayServer"),
		metadataProvider:    mp,
		blobProvider:        bp,
		legacyChunkProvider: cp,
		chunkReader:         chunkReader,
		blobRateLimiter:     limiter.NewBlobRateLimiter(&config.RateLimits, relayMetrics),
		chunkRateLimiter:    limiter.NewChunkRateLimiter(&config.RateLimits, relayMetrics),
		authenticator:       authenticator,
		replayGuardian:      replayGuardian,
		metrics:             relayMetrics,
		chainReader:         chainReader,
		listener:            listener,
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

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {
	start := time.Now()

	if s.config.Timeouts.GetBlobTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetBlobTimeout)
		defer cancel()
	}

	// Validate the request params before any further processing (as validation is cheaper)
	key, err := v2.BytesToBlobKey(request.GetBlobKey())
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
		if strings.Contains(err.Error(), blobstore.ErrMetadataNotFound.Error()) {
			// nolint:wrapcheck
			return nil, api.NewErrorNotFound(
				fmt.Sprintf("blob %s not found, check if blob exists and is assigned to this relay", key.Hex()))
		}
		// nolint:wrapcheck
		return nil, api.NewErrorInternal(fmt.Sprintf("error fetching metadata for blob: %v", err))
	}
	metadata := mMap[v2.BlobKey(request.GetBlobKey())]
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
		if strings.Contains(err.Error(), blobstore.ErrBlobNotFound.Error()) {
			return nil, api.NewErrorNotFound(fmt.Sprintf("blob %s not found", key.Hex()))
		} else {
			s.logger.Errorf("error fetching blob %s: %v", key.Hex(), err)
			return nil, api.NewErrorInternal(
				fmt.Sprintf("relay encountered errors while attempting to fetch blob %s", key.Hex()))
		}
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
	if len(request.GetChunkRequests()) == 0 {
		return api.NewErrorInvalidArg("no chunk requests provided")
	}
	if len(request.GetChunkRequests()) > s.config.MaxKeysPerGetChunksRequest {
		return api.NewErrorInvalidArg(fmt.Sprintf(
			"too many chunk requests provided, max is %d", s.config.MaxKeysPerGetChunksRequest))
	}

	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() == nil && chunkRequest.GetByRange() == nil {
			return api.NewErrorInvalidArg("chunk request must be either by index or by range")
		}
	}

	return nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {
	start := time.Now()

	if s.config.Timeouts.GetChunksTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.Timeouts.GetChunksTimeout)
		defer cancel()
	}
	err := s.validateGetChunksRequest(request)
	if err != nil {
		return nil, err
	}

	s.metrics.ReportChunkKeyCount(len(request.GetChunkRequests()))

	if s.authenticator != nil {
		client, ok := peer.FromContext(ctx)
		if !ok {
			return nil, api.NewErrorInvalidArg("could not get peer information")
		}
		clientAddress := client.Addr.String()

		hash, err := s.authenticator.AuthenticateGetChunksRequest(ctx, request)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			s.logger.Debug("rejected GetChunks request", "client", clientAddress)
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("auth failed: %v", err))
		}

		timestamp := time.Unix(int64(request.GetTimestamp()), 0)
		err = s.replayGuardian.VerifyRequest(hash, timestamp)
		if err != nil {
			s.metrics.ReportChunkAuthFailure()
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("failed to verify request: %v", err))
		}

		s.logger.Debug("received authenticated GetChunks request", "client", clientAddress)
	}

	finishedAuthenticating := time.Now()
	if s.authenticator != nil {
		s.metrics.ReportChunkAuthenticationLatency(finishedAuthenticating.Sub(start))
	}

	clientID := string(request.GetOperatorId())
	err = s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
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
		if strings.Contains(err.Error(), blobstore.ErrMetadataNotFound.Error()) {
			// nolint:wrapcheck
			return nil, api.NewErrorNotFound(
				fmt.Sprintf("blob not found, check if blob exists and is assigned to this relay:: %v", keys))
		}
		// nolint:wrapcheck
		return nil, api.NewErrorInternal(fmt.Sprintf("error fetching metadata for blob: %v", err))
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

	// Determine whether to use legacy chunk provider or new chunk provider. We have to use the legacy chunk
	// provider if there are any requests that use the "by index" query pattern.
	useLegacyChunkProvider := false
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			useLegacyChunkProvider = true
			break
		}
	}

	var bytesToSend [][]byte

	if useLegacyChunkProvider {
		frames, err := s.legacyChunkProvider.GetFrames(ctx, mMap)
		if err != nil {
			// nolint:wrapcheck
			return nil, api.NewErrorInternal(fmt.Sprintf("error fetching frames: %v", err))
		}

		bytesToSend, err = gatherChunkDataToSendLegacy(frames, request)
		if err != nil {
			// nolint:wrapcheck
			return nil, api.NewErrorInternal(fmt.Sprintf("error gathering chunk data: %v", err))
		}
	} else {
		var found bool
		bytesToSend, found, err = s.gatherChunkDataToSend(ctx, mMap, request)
		if err != nil {
			// nolint:wrapcheck
			return nil, api.NewErrorInternal(fmt.Sprintf("error gathering chunk data: %v", err))
		}
		if !found {
			// nolint:wrapcheck
			return nil, api.NewErrorNotFound("requested chunks not found")
		}
	}

	s.metrics.ReportChunkDataLatency(time.Since(finishedFetchingMetadata))
	s.metrics.ReportChunkLatency(time.Since(start))

	return &pb.GetChunksReply{
		Data: bytesToSend,
	}, nil
}

// getKeysFromChunkRequest gathers a slice of blob keys from a GetChunks request.
func getKeysFromChunkRequest(request *pb.GetChunksRequest) ([]v2.BlobKey, error) {
	keys := make([]v2.BlobKey, 0, len(request.GetChunkRequests()))

	for _, chunkRequest := range request.GetChunkRequests() {
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

// Used to pass status of downloads from goroutines up to controlling function.
type downloadResult struct {
	key   v2.BlobKey
	found bool
}

// Download and compile the chunk data to send back to the client.
func (s *Server) gatherChunkDataToSend(
	ctx context.Context,
	metadataMap map[v2.BlobKey]*blobMetadata,
	request *pb.GetChunksRequest,
) ([][]byte, bool, error) {

	coefficients, proofs, found, err := s.downloadDataFromRelays(ctx, metadataMap, request)
	if err != nil {
		return nil, false, fmt.Errorf("error downloading chunk data from relays: %w", err)
	}
	if !found {
		return nil, false, nil
	}

	chunkDataObjects, err := combineProofsAndCoefficients(
		proofs,
		coefficients,
		request,
		metadataMap)
	if err != nil {
		return nil, false, fmt.Errorf("error building chunk data: %w", err)
	}

	bytesToSend, err := buildBinaryChunkData(chunkDataObjects, request)
	if err != nil {
		return nil, false, fmt.Errorf("error building binary chunk data: %w", err)
	}

	return bytesToSend, true, nil
}

// Download all data from relays needed to fulfill a GetChunks request.
func (s *Server) downloadDataFromRelays(
	ctx context.Context,
	metadataMap map[v2.BlobKey]*blobMetadata,
	request *pb.GetChunksRequest,
) (coefficients [][][]byte, proofs [][][]byte, allDataFound bool, err error) {
	requestCount := len(request.GetChunkRequests())

	coefficients = make([][][]byte, requestCount)
	proofs = make([][][]byte, requestCount)

	results := make(chan downloadResult, requestCount*2)

	runner, ctx := errgroup.WithContext(ctx)

	// Fan out and make requests in parallel
	for i, chunkRequest := range request.GetChunkRequests() {
		blobKey := v2.BlobKey(chunkRequest.GetByRange().GetBlobKey())
		metadata := metadataMap[blobKey]

		// Download proofs
		runner.Go(func() error {
			data, found, err := s.chunkReader.GetBinaryChunkProofsRange(
				ctx,
				blobKey,
				chunkRequest.GetByRange().GetStartIndex(),
				chunkRequest.GetByRange().GetEndIndex(),
			)

			proofs[i] = data
			if err != nil {
				return fmt.Errorf("failed to download proofs: %w", err)
			}
			results <- downloadResult{key: blobKey, found: found}

			return nil
		})
		// Download coefficients
		runner.Go(func() error {
			data, found, err := s.chunkReader.GetBinaryChunkCoefficientRange(
				ctx,
				blobKey,
				chunkRequest.GetByRange().GetStartIndex(),
				chunkRequest.GetByRange().GetEndIndex(),
				metadata.symbolsPerFrame,
			)

			coefficients[i] = data
			if err != nil {
				return fmt.Errorf("failed to download coefficients: %w", err)
			}
			results <- downloadResult{key: blobKey, found: found}

			return nil
		})
	}

	// Await results
	if err := runner.Wait(); err != nil {
		return nil, nil, false, fmt.Errorf("error downloading chunk data: %w", err)
	}

	// Handle the situation where some data couldn't be found
	for i := 0; i < requestCount*2; i++ {
		result := <-results
		if !result.found {
			return nil, nil, false, nil
		}
	}

	return coefficients, proofs, true, nil
}

// Convert the disparate proofs and coefficients into unified "ChunkData" objects
// (or "chunks" or "frames" or other names, depending on what part of the code you are looking at)
func combineProofsAndCoefficients(
	proofs [][][]byte,
	coefficients [][][]byte,
	request *pb.GetChunksRequest,
	metadataMap map[v2.BlobKey]*blobMetadata,
) ([]*core.ChunksData, error) {

	requestCount := len(request.GetChunkRequests())

	chunkDataObjects := make([]*core.ChunksData, requestCount)
	for i := 0; i < requestCount; i++ {
		blobKey := v2.BlobKey(request.GetChunkRequests()[i].GetByRange().GetBlobKey())
		metadata := metadataMap[blobKey]
		chunkData, err := buildChunksData(proofs[i], int(metadata.symbolsPerFrame), coefficients[i])
		if err != nil {
			return nil, fmt.Errorf("error building chunk data: %w", err)
		}
		chunkDataObjects[i] = chunkData
	}

	return chunkDataObjects, nil
}

// Take the chunk data objects and build the final byte arrays to send back to the client.
func buildBinaryChunkData(
	chunkDataObjects []*core.ChunksData,
	request *pb.GetChunksRequest,
) ([][]byte, error) {

	bytesToSend := make([][]byte, 0, len(request.GetChunkRequests()))
	for requestIndex := 0; requestIndex < len(request.GetChunkRequests()); requestIndex++ {
		nextRequest := request.GetChunkRequests()[requestIndex]
		targetKey := nextRequest.GetByRange().GetBlobKey()

		chunkDataToSend := chunkDataObjects[requestIndex]

		// Validator verification logic expects all chunks for the same blob to be grouped together.
		// This is easy to do with an index request, since an index request allows non-contiguous chunks
		// to be fetched via a single request. But range queries require contiguous chunks, so we may receive
		// multiple range requests for the same blob. In order to avoid breaking tricky validation logic,
		// it is simpler to just group all range requests for the same blob together into a single "bundle"
		// (aka a binary object that encodes a list of chunks).

		// If there are multiple requests for the same blob, combine them.
		for i := requestIndex + 1; i < len(request.GetChunkRequests()); i++ {
			followingRequest := request.GetChunkRequests()[i].GetByRange()

			nextKey := followingRequest.GetBlobKey()
			if !bytes.Equal(targetKey, nextKey) {
				// Next request is for a different blob, don't combine.
				break
			}

			followingChunkData := chunkDataObjects[i]
			chunkDataToSend.Chunks = append(chunkDataToSend.Chunks, followingChunkData.Chunks...)

			// Bump the counter for the outer loop since this iteration handles it
			requestIndex++
		}

		bundleBytes, err := chunkDataToSend.FlattenToBundle()
		if err != nil {
			return nil, fmt.Errorf("error serializing chunk subset: %w", err)
		}

		bytesToSend = append(bytesToSend, bundleBytes)
	}

	return bytesToSend, nil
}

// gatherChunkDataToSendLegacy takes the chunk data and narrows it down to the data requested in the GetChunks request.
// Required for requests that use the old "by index" query pattern.
func gatherChunkDataToSendLegacy(
	frames map[v2.BlobKey]*core.ChunksData,
	request *pb.GetChunksRequest) ([][]byte, error) {

	bytesToSend := make([][]byte, 0, len(request.GetChunkRequests()))

	for requestIndex := 0; requestIndex < len(request.GetChunkRequests()); requestIndex++ {
		nextRequest := request.GetChunkRequests()[requestIndex]

		var framesSubset *core.ChunksData
		var err error

		if nextRequest.GetByIndex() != nil {
			framesSubset, err = selectFrameSubsetByIndex(nextRequest.GetByIndex(), frames)
		} else {
			// Validator verification logic expects all chunks for the same blob to be grouped together.
			// This is easy to do with an index request, since an index request allows non-contiguous chunks
			// to be fetched via a single request. But range queries require contiguous chunks, so we may receive
			// multiple range requests for the same blob. In order to avoid breaking tricky validation logic,
			// it is simpler to just group all range requests for the same blob together into a single "bundle"
			// (aka a binary object that encodes a list of chunks).

			rangeRequests := make([]*pb.ChunkRequestByRange, 0)
			rangeRequests = append(rangeRequests, nextRequest.GetByRange())

			targetKey := nextRequest.GetByRange().GetBlobKey()

			// If there are multiple range requests for the same blob, combine them.
			for i := requestIndex + 1; i < len(request.GetChunkRequests()); i++ {
				followingRequest := request.GetChunkRequests()[i]
				followingRangeRequest := followingRequest.GetByRange()
				if followingRangeRequest == nil {
					// Following request is not by range, don't combine.
					break
				}

				nextKey := followingRangeRequest.GetBlobKey()
				if bytes.Equal(targetKey, nextKey) == false {
					// Next request is for a different blob, don't combine.
					break
				}

				rangeRequests = append(rangeRequests, followingRangeRequest)
				// Bump the counter for the outer loop since this iteration will handle it
				requestIndex++
			}

			framesSubset, err = selectFrameSubsetByRange(rangeRequests, frames)
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
	// One or more requests for chunks from the same blob
	requests []*pb.ChunkRequestByRange,
	allFrames map[v2.BlobKey]*core.ChunksData,
) (*core.ChunksData, error) {

	key := v2.BlobKey(requests[0].GetBlobKey())
	frames, ok := allFrames[key]
	if !ok {
		return nil, fmt.Errorf("frames not found for key %s", key.Hex())
	}

	chunkCount := 0
	for _, request := range requests {
		chunkCount += int(request.GetEndIndex() - request.GetStartIndex())
	}
	chunks := make([][]byte, 0, chunkCount)

	for _, request := range requests {
		startIndex := request.GetStartIndex()
		endIndex := request.GetEndIndex()

		if startIndex > endIndex {
			return nil, fmt.Errorf(
				"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
				startIndex, endIndex, key.Hex())
		}
		if endIndex > uint32(len(frames.Chunks)) {
			return nil, fmt.Errorf(
				"chunk range %d-%d is invalid for key %s, chunk count %d",
				startIndex, endIndex, key, len(frames.Chunks))
		}

		chunks = append(chunks, frames.Chunks[startIndex:endIndex]...)
	}

	framesSubset := &core.ChunksData{
		Chunks:   chunks,
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

	if len(request.GetChunkIndices()) > len(frames.Chunks) {
		return nil, fmt.Errorf("too many requested chunks for key %s, chunk count %d",
			key.Hex(), len(frames.Chunks))
	}

	framesSubset := &core.ChunksData{
		Format:   frames.Format,
		ChunkLen: frames.ChunkLen,
		Chunks:   make([][]byte, 0, len(request.GetChunkIndices())),
	}

	for _, index := range request.GetChunkIndices() {
		if index >= uint32(len(frames.Chunks)) {
			return nil, fmt.Errorf(
				"chunk index %d out of range for key %s, chunk count %d",
				index, key.Hex(), len(frames.Chunks))
		}

		framesSubset.Chunks = append(framesSubset.Chunks, frames.Chunks[index])
	}

	return framesSubset, nil
}

// computeChunkRequestRequiredBandwidth computes the bandwidth required to fulfill a GetChunks request.
func computeChunkRequestRequiredBandwidth(
	request *pb.GetChunksRequest,
	mMap map[v2.BlobKey]*blobMetadata,
) (uint32, error) {
	requiredBandwidth := uint32(0)
	for _, req := range request.GetChunkRequests() {
		var metadata *blobMetadata
		var key v2.BlobKey
		var requestedChunks uint32

		if req.GetByIndex() != nil {
			key = v2.BlobKey(req.GetByIndex().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = uint32(len(req.GetByIndex().GetChunkIndices()))
		} else {
			key = v2.BlobKey(req.GetByRange().GetBlobKey())
			metadata = mMap[key]

			if req.GetByRange().GetEndIndex() < req.GetByRange().GetStartIndex() {
				return 0, fmt.Errorf(
					"chunk range %d-%d is invalid for key %s, start index must be less than or equal to end index",
					req.GetByRange().GetStartIndex(), req.GetByRange().GetEndIndex(), key.Hex())
			}

			requestedChunks = req.GetByRange().GetEndIndex() - req.GetByRange().GetStartIndex()
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
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			chunkCount += len(chunkRequest.GetByIndex().GetChunkIndices())
		} else {
			chunkCount += int(chunkRequest.GetByRange().GetEndIndex() - chunkRequest.GetByRange().GetStartIndex())
		}
	}

	blobCount := len(request.GetChunkRequests())

	return api.NewErrorResourceExhausted(fmt.Sprintf("unable to serve data (%d blobs, %d chunks, %d bytes): %v",
		blobCount, chunkCount, requiredBandwidth, originalError))
}

// Retrieves all chunks allocated to a validator.
// The relay computes which chunks to return based on the deterministic chunk allocation algorithm.
//
// This endpoint will eventually replace `GetChunks`. It is being added as a separate endpoint for the sake of
// backwards compatibility
func (s *Server) GetValidatorChunks(
	ctx context.Context,
	request *pb.GetValidatorChunksRequest,
) (*pb.GetChunksReply, error) {
	// TODO(litt3): this logic will be implemented in a future PR.
	return nil, status.Errorf(codes.Unimplemented, "method GetValidatorChunks not implemented")
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
