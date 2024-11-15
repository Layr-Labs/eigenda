package relay

import (
	"context"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"time"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	// config is the configuration for the relay Server.
	config *Config

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
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader chunkstore.ChunkReader) (*Server, error) {

	ms, err := newMetadataProvider(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataMaxConcurrency,
		config.RelayIDs)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata server: %w", err)
	}

	bs, err := newBlobProvider(
		ctx,
		logger,
		blobStore,
		config.BlobCacheSize,
		config.BlobMaxConcurrency)
	if err != nil {
		return nil, fmt.Errorf("error creating blob server: %w", err)
	}

	cs, err := newChunkProvider(
		ctx,
		logger,
		chunkReader,
		config.ChunkCacheSize,
		config.ChunkMaxConcurrency)
	if err != nil {
		return nil, fmt.Errorf("error creating chunk server: %w", err)
	}

	return &Server{
		config:           config,
		metadataProvider: ms,
		blobProvider:     bs,
		chunkProvider:    cs,
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

	err = s.blobRateLimiter.RequestGetBlobBandwidth(time.Now(), metadata.blobSizeBytes) // TODO make sure this field is populated
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

	clientID := fmt.Sprintf("%d", request.RequesterId) //TODO
	err := s.chunkRateLimiter.BeginGetChunkOperation(time.Now(), clientID)
	if err != nil {
		return nil, err
	}
	defer s.chunkRateLimiter.FinishGetChunkOperation(clientID)

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

	mMap, err := s.metadataProvider.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blob, check if blob exists and is assigned to this relay: %w", err)
	}

	requiredBandwidth := 0 // TODO calculate this
	err = s.chunkRateLimiter.RequestGetChunkBandwidth(time.Now(), clientID, requiredBandwidth)
	if err != nil {
		return nil, err
	}

	frames, err := s.chunkProvider.GetFrames(ctx, mMap)
	if err != nil {
		return nil, fmt.Errorf("error fetching frames: %w", err)
	}

	bytesToSend := make([][]byte, 0, len(keys))

	// return data in the order that it was requested
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

	return &pb.GetChunksReply{
		Data: bytesToSend,
	}, nil
}
