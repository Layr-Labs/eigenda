package relay

import (
	"context"
	"fmt"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"sync"
	"sync/atomic"
)

// TODO unit tests

// metadataServer encapsulates logic for fetching metadata for blobs. Utilized by the relay Server.
type metadataServer struct {
	ctx           context.Context
	metadataStore *blobstore.BlobMetadataStore

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	metadataCache CachedAccessor[core.BlobKey, blobMetadata]
}

func newMetadataServer(
	ctx context.Context,
	metadataStore *blobstore.BlobMetadataStore,
	metadataCacheSize int) (*metadataServer, error) {

	server := &metadataServer{
		ctx:           ctx,
		metadataStore: metadataStore,
	}

	metadataCache, err := NewCachedAccessor[core.BlobKey, blobMetadata](metadataCacheSize, server.getMetadataForBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}

	server.metadataCache = metadataCache

	return server, nil
}

// getMetadataForBlob retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (m *metadataServer) getMetadataForBlob(key core.BlobKey) (*blobMetadata, error) {
	// Retrieve the metadata from the store.
	fullMetadata, err := m.metadataStore.GetBlobMetadata(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	// TODO check if the blob is in the correct shard, return error as if not found if not in the correct shard

	metadata := &blobMetadata{
		// TODO
	}

	return metadata, nil
}

// getMetadataForBlobs retrieves metadata about multiple blobs in parallel.
func (m *metadataServer) getMetadataForBlobs(keys []core.BlobKey) (map[core.BlobKey]*blobMetadata, error) {

	if len(keys) == 1 {
		// Special case: no need to spawn a goroutine.
		metadata, err := m.metadataCache.Get(keys[0])
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

			metadata, err := m.metadataCache.Get(boundKey)
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
