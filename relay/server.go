package relay

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
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
}

type Config struct {

	// RelayIDs contains the IDs of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	RelayIDs []v2.RelayKey

	// GRPCPort is the port that the relay server listens on.
	GRPCPort int

	// MaxGRPCMessageSize is the maximum size of a gRPC message that the server will accept.
	MaxGRPCMessageSize int

	// MetadataCacheSize is the maximum number of items in the metadata cache.
	MetadataCacheSize int

	// MetadataMaxConcurrency puts a limit on the maximum number of concurrent metadata fetches actively running on
	// goroutines.
	MetadataMaxConcurrency int

	// BlobCacheSize is the maximum number of items in the blob cache.
	BlobCacheSize int

	// BlobMaxConcurrency puts a limit on the maximum number of concurrent blob fetches actively running on goroutines.
	BlobMaxConcurrency int

	// ChunkCacheSize is the maximum number of items in the chunk cache.
	ChunkCacheSize int

	// ChunkMaxConcurrency is the size of the work pool for fetching chunks. Note that this does not
	// impact concurrency utilized by the s3 client to upload/download fragmented files.
	ChunkMaxConcurrency int

	// MaxKeysPerGetChunksRequest is the maximum number of keys that can be requested in a single GetChunks request.
	MaxKeysPerGetChunksRequest int

	// RateLimits contains configuration for rate limiting.
	RateLimits limiter.Config
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader chunkstore.ChunkReader) (*Server, error) {

	mp, err := newMetadataProvider(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataMaxConcurrency,
		config.RelayIDs)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata provider: %w", err)
	}

	bp, err := newBlobProvider(
		ctx,
		logger,
		blobStore,
		config.BlobCacheSize,
		config.BlobMaxConcurrency)
	if err != nil {
		return nil, fmt.Errorf("error creating blob provider: %w", err)
	}

	cp, err := newChunkProvider(
		ctx,
		logger,
		chunkReader,
		config.ChunkCacheSize,
		config.ChunkMaxConcurrency)
	if err != nil {
		return nil, fmt.Errorf("error creating chunk provider: %w", err)
	}

	return &Server{
		config:           config,
		logger:           logger,
		metadataProvider: mp,
		blobProvider:     bp,
		chunkProvider:    cp,
		blobRateLimiter:  limiter.NewBlobRateLimiter(&config.RateLimits),
		chunkRateLimiter: limiter.NewChunkRateLimiter(&config.RateLimits),
	}, nil
}

// GetBlob retrieves a blob stored by the relay.
func (s *Server) GetBlob(ctx context.Context, request *pb.GetBlobRequest) (*pb.GetBlobReply, error) {

	// Future work:
	//  - timeouts

	err := s.blobRateLimiter.BeginGetBlobOperation(time.Now())
	if err != nil {
		return nil, err
	}

	key, err := v2.BytesToBlobKey(request.BlobKey)
	if err != nil {
		return nil, fmt.Errorf("invalid blob key: %w", err)
	}

	keys := []v2.BlobKey{key}
	mMap, err := s.metadataProvider.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %w", err)
	}
	metadata := mMap[v2.BlobKey(request.BlobKey)]
	if metadata == nil {
		return nil, fmt.Errorf("blob not found")
	}

	err = s.blobRateLimiter.RequestGetBlobBandwidth(time.Now(), metadata.blobSizeBytes)
	if err != nil {
		return nil, err
	}

	data, err := s.blobProvider.GetBlob(key)
	if err != nil {
		return nil, fmt.Errorf("error fetching blob %s: %w", key.Hex(), err)
	}

	reply := &pb.GetBlobReply{
		Blob: data,
	}

	return reply, nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(ctx context.Context, request *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// Future work:
	//  - authentication
	//  - timeouts

	if len(request.ChunkRequests) <= 0 {
		return nil, fmt.Errorf("no chunk requests provided")
	}
	if len(request.ChunkRequests) > s.config.MaxKeysPerGetChunksRequest {
		return nil, fmt.Errorf(
			"too many chunk requests provided, max is %d", s.config.MaxKeysPerGetChunksRequest)
	}

	// Future work: client IDs will be fixed when authentication is implemented
	clientID := fmt.Sprintf("%d", request.RequesterId)
	err := s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
	if err != nil {
		return nil, err
	}
	defer s.chunkRateLimiter.FinishGetChunkOperation(clientID)

	keys, err := getKeysFromChunkRequest(request)
	if err != nil {
		return nil, err
	}

	mMap, err := s.metadataProvider.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %w", err)
	}

	requiredBandwidth := computeChunkRequestRequiredBandwidth(request, mMap)
	err = s.chunkRateLimiter.RequestGetChunkBandwidth(time.Now(), clientID, requiredBandwidth)
	if err != nil {
		return nil, err
	}

	frames, err := s.chunkProvider.GetFrames(ctx, mMap)
	if err != nil {
		return nil, fmt.Errorf("error fetching frames: %w", err)
	}

	bytesToSend, err := gatherChunkDataToSend(frames, request)
	if err != nil {
		return nil, fmt.Errorf("error gathering chunk data: %w", err)
	}

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

	bytesToSend := make([][]byte, 0, len(frames))

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
func computeChunkRequestRequiredBandwidth(request *pb.GetChunksRequest, mMap metadataMap) int {
	requiredBandwidth := 0
	for _, req := range request.ChunkRequests {
		var metadata *blobMetadata
		var requestedChunks int

		if req.GetByIndex() != nil {
			key := v2.BlobKey(req.GetByIndex().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = len(req.GetByIndex().ChunkIndices)
		} else {
			key := v2.BlobKey(req.GetByRange().GetBlobKey())
			metadata = mMap[key]
			requestedChunks = int(req.GetByRange().EndIndex - req.GetByRange().StartIndex)
		}

		requiredBandwidth += requestedChunks * int(metadata.chunkSizeBytes)
	}

	return requiredBandwidth

}

// Start starts the server listening for requests. This method will block until the server is stopped.
func (s *Server) Start() error {
	// Serve grpc requests
	addr := fmt.Sprintf("0.0.0.0:%d", s.config.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not start tcp listener on %s: %w", addr, err)
	}

	opt := grpc.MaxRecvMsgSize(s.config.MaxGRPCMessageSize)

	s.grpcServer = grpc.NewServer(opt)
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

// Stop stops the server.
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
}
