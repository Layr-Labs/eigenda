package relay

import (
	"context"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
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
}

// Metadata about a blob. The relay only needs a small subset of a blob's metadata.
type blobMetadata struct {
	// the size of the file containing the encoded chunks
	totalChunkSizeBytes uint32
	// the fragment size used for uploading the encoded chunks
	fragmentSizeBytes uint32
}

// NewServer creates a new relay Server.
func NewServer(
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader *chunkstore.ChunkReader) (*Server, error) {

	ms, err := newMetadataServer(context.Background(), metadataStore, config.MetadataCacheSize)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata server: %w", err)
	}

	return &Server{
		config:         config,
		metadataServer: ms,
		blobStore:      blobStore,
		chunkReader:    chunkReader,
	}, nil
}

// GetBlobs retrieves blobs stored by the relay.
func (s *Server) GetBlobs(ctx context.Context, request *pb.GetBlobsRequest) (*pb.GetBlobsReply, error) {

	// Future work: rate limiting
	// TODO: max request size
	// TODO: limit parallelism

	// TODO better way to do this conversion?
	blobKeyBytes := request.BlobKeys
	blobKeys := make([]core.BlobKey, len(blobKeyBytes))
	for i, keyBytes := range blobKeyBytes {
		blobKey := core.BlobKey(keyBytes)
		blobKeys[i] = blobKey
	}

	// Fetch metadata for the blobs. This fails if any of the blobs do not exist, or if any of the blobs
	// are assigned to a shard that is not managed by this relay.
	_, err := s.metadataServer.getMetadataForBlobs(blobKeys)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata for blobs, check if blobs exist and are assigned to this relay: %w", err)
	}

	// TODO continue here

	return nil, nil // TODO
}

// GetChunks retrieves chunks from blobs stored by the relay.
func (s *Server) GetChunks(context.Context, *pb.GetChunksRequest) (*pb.GetChunksReply, error) {

	// Future work: rate limiting
	// Future work: authentication
	// TODO: max request size
	// TODO: limit parallelism

	return nil, nil // TODO
}
