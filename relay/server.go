package relay

import (
	"context"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	lru "github.com/hashicorp/golang-lru/v2"
)

var _ pb.RelayServer = &Server{}

// Server implements the Relay service defined in api/proto/relay/relay.proto
type Server struct {
	pb.UnimplementedRelayServer

	config        *Config
	metadataStore *blobstore.BlobMetadataStore
	blobStore     *blobstore.BlobStore
	chunkReader   *chunkstore.ChunkReader

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	metadataCache *lru.Cache[core.BlobKey, *blobMetadata]

	// TODO chunk cache
	// TODO blob cache
}

// Metadata about a blob.
type blobMetadata struct {

	// TODO only store what is absolutely needed

	certificate  *core.BlobCertificate
	fragmentInfo *encoding.FragmentInfo
}

func NewServer(
	config *Config,
	metadataStore *blobstore.BlobMetadataStore,
	blobStore *blobstore.BlobStore,
	chunkReader *chunkstore.ChunkReader) (*Server, error) {

	metadataCache, err := lru.New[core.BlobKey, *blobMetadata](config.MetadataCacheSize)

	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}

	return &Server{
		config:        config,
		metadataStore: metadataStore,
		blobStore:     blobStore,
		chunkReader:   chunkReader,
		metadataCache: metadataCache,
	}, nil
}

// GetBlobs retrieves blobs stored by the relay.
func (s *Server) GetBlobs(context.Context, *pb.GetBlobsRequest) (*pb.GetBlobsReply, error) {

	// Future work: rate limiting
	// TODO: max request size
	// TODO: limit parallelism

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

// getBlobMetadata retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (s *Server) getBlobMetadata(ctx context.Context, key core.BlobKey) (*blobMetadata, error) {
	// Check the cache first.
	if metadata, ok := s.metadataCache.Get(key); ok {
		return metadata, nil
	}

	// TODO protect against parallel requests for the same blob

	// Retrieve the metadata from the store.
	fullMetadata, err := s.metadataStore.GetBlobMetadata(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	// TODO check if the blob is in the correct shard, return error as if not found if not in the correct shard

	metadata := &blobMetadata{
		certificate:  nil, // TODO
		fragmentInfo: fullMetadata.FragmentInfo,
	}

	// Cache the metadata.
	s.metadataCache.Add(key, metadata)

	return metadata, nil
}
