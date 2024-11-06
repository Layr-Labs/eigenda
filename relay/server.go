package relay

import (
	"context"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	config      *Config
	blobStore   *blobstore.BlobStore
	chunkReader *chunkstore.ChunkReader

	// metadataServer encapsulates logic for fetching metadata for blobs.
	metadataServer *metadataServer

	// blobServer encapsulates logic for fetching blobs.
	blobServer *blobServer
}

// NewServer creates a new relay Server.
func NewServer(
	ctx context.Context,
	logger logging.Logger,
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader *chunkstore.ChunkReader) (*Server, error) {

	ms, err := NewMetadataServer(
		ctx,
		logger,
		metadataStore,
		config.MetadataCacheSize,
		config.MetadataWorkPoolSize,
		config.Shards)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata server: %w", err)
	}

	bs, err := NewBlobServer(
		ctx,
		logger,
		blobStore,
		config.BlobCacheSize,
		config.BlobWorkPoolSize)
	if err != nil {
		return nil, fmt.Errorf("error creating blob server: %w", err)
	}

	return &Server{
		config:         config,
		metadataServer: ms,
		blobServer:     bs,
		blobStore:      blobStore,
		chunkReader:    chunkReader,
	}, nil
}

// GetBlobs retrieves blobs stored by the relay.
func (s *Server) GetBlobs(ctx context.Context, request *pb.GetBlobsRequest) (*pb.GetBlobsReply, error) {

	// TODO:
	//  - global throttle requests / sec
	//  - per-connection throttle requests / sec
	//  - timeouts

	if len(request.BlobKeys) > s.config.MaximumBlobKeyLimit {
		return nil, fmt.Errorf("request touches too many blobs, limit is %d", s.config.MaximumBlobKeyLimit)
	}

	keys := make([]v2.BlobKey, len(request.BlobKeys))
	for i, keyBytes := range request.BlobKeys {
		keys[i] = v2.BlobKey(keyBytes)
	}

	metadataMap, err := s.metadataServer.GetMetadataForBlobs(keys)
	if err != nil {
		return nil, fmt.Errorf("error fetching metadata for blobs, check if blobs exist and are assigned to this relay: %w", err)
	}

	// TODO:
	//  - global bytes / sec throttle
	//  - per-connection bytes / sec throttle
	//  - maximum bytes per request limit

	dataMap, err := s.blobServer.GetBlobs(metadataMap)
	if err != nil {
		return nil, fmt.Errorf("error fetching blobs: %w", err)
	}

	requestedBlobs := make([]*pb.RequestedBlob, len(request.BlobKeys))
	for i, key := range request.BlobKeys {

		// TODO return individual errors per requested blob maybe

		blobKey := v2.BlobKey(key)
		data := (*dataMap)[blobKey]

		requestedBlobs[i] = &pb.RequestedBlob{
			Data: &pb.RequestedBlob_Blob{
				Blob: data,
			},
		}
	}

	reply := &pb.GetBlobsReply{
		Blobs: requestedBlobs,
	}

	return reply, nil
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// Future work: rate limiting
	// Future work: authentication
	// TODO: max request size
	// TODO: limit parallelism

	return nil, nil // TODO
}
