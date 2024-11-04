package relay

import (
	"context"
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"sync"
	"sync/atomic"
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
	metadataCache CachedAccessor[core.BlobKey, blobMetadata]
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

	server := &Server{
		config:        config,
		metadataStore: metadataStore,
		blobStore:     blobStore,
		chunkReader:   chunkReader,
	}

	metadataCache, err :=
		NewCachedAccessor[core.BlobKey, blobMetadata](config.MetadataCacheSize, server.getMetadataForBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}
	server.metadataCache = metadataCache

	return server, nil
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
	_, err := s.getMetadataForBlobs(blobKeys)
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

// getMetadataForBlob retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (s *Server) getMetadataForBlob(key core.BlobKey) (*blobMetadata, error) {
	// Retrieve the metadata from the store.
	fullMetadata, err := s.metadataStore.GetBlobMetadata(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	// TODO check if the blob is in the correct shard, return error as if not found if not in the correct shard

	metadata := &blobMetadata{
		certificate:  nil, // TODO
		fragmentInfo: fullMetadata.FragmentInfo,
	}

	return metadata, nil
}

// getMetadataForBlobs retrieves metadata about multiple blobs in parallel.
func (s *Server) getMetadataForBlobs(keys []core.BlobKey) (map[core.BlobKey]*blobMetadata, error) {

	if len(keys) == 1 {
		// Special case: no need to spawn a goroutine.
		metadata, err := s.metadataCache.Get(keys[0])
		if err != nil {
			return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", keys[0].Hex(), err)
		}

		return map[core.BlobKey]*blobMetadata{keys[0]: metadata}, nil
	}

	metadataMap := make(map[core.BlobKey]*blobMetadata)
	errors := atomic.Bool{}

	wg := sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		// TODO limits on parallelism
		boundKey := key
		go func() {
			defer wg.Done()

			metadata, err := s.metadataCache.Get(boundKey)
			if err != nil {
				// TODO log error
				errors.Store(true)
				return
			}
			metadataMap[boundKey] = metadata
		}()
	}

	wg.Wait()

	if errors.Load() {
		return nil, fmt.Errorf("error retrieving metadata for one or more blobs")
	}

	return metadataMap, nil
}
