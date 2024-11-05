package relay

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/core/v2"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"sync"
	"sync/atomic"
)

// TODO unit tests

// metadataServer encapsulates logic for fetching metadata for blobs. Utilized by the relay Server.
type metadataServer struct {
	ctx    context.Context
	logger logging.Logger

	// metadataStore can be used to read blob metadata from dynamoDB.
	metadataStore *blobstore.BlobMetadataStore

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	metadataCache CachedAccessor[core.BlobKey, blobMetadata]

	// shardSet is the set of shards assigned to this relay. This relay will refuse to serve metadata for blobs
	// that are not assigned to one of these shards.
	shardSet map[core.RelayKey]struct{}
}

func newMetadataServer(
	ctx context.Context,
	logger logging.Logger,
	metadataStore *blobstore.BlobMetadataStore,
	metadataCacheSize int,
	shards []core.RelayKey) (*metadataServer, error) {

	shardSet := make(map[core.RelayKey]struct{}, len(shards))
	for _, shard := range shards {
		shardSet[shard] = struct{}{}
	}

	server := &metadataServer{
		ctx:           ctx,
		logger:        logger,
		metadataStore: metadataStore,
		shardSet:      shardSet,
	}

	metadataCache, err := NewCachedAccessor[core.BlobKey, blobMetadata](metadataCacheSize, server.getMetadataForBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}

	server.metadataCache = metadataCache

	return server, nil
}

// getMetadataForBlob retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (m *metadataServer) getMetadataForBlob(key v2.BlobKey) (*blobMetadata, error) {
	// Retrieve the metadata from the store.
	cert, fragmentInfo, err := m.metadataStore.GetBlobCertificate(m.ctx, core.BlobKey(key))
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	validShard := false
	for _, shard := range cert.RelayKeys {
		if _, ok := m.shardSet[shard]; ok {
			validShard = true
			break
		}
	}

	if !validShard {
		return nil, fmt.Errorf("blob %s is not assigned to this relay", key.Hex())
	}

	metadata := &blobMetadata{
		totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
		fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
	}

	return metadata, nil
}

// getMetadataForBlobs retrieves metadata about multiple blobs in parallel.
func (m *metadataServer) getMetadataForBlobs(keys []core.BlobKey) (map[core.BlobKey]*blobMetadata, error) {

	// TODO figure out how timeouts are going to work here

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
				// Intentionally log at debug level. External users can force this condition to trigger
				// by requesting metadata for a blob that does not exist, and so it's important to avoid
				// allowing hooligans to spam the logs in production environments.
				m.logger.Debug("error retrieving metadata for blob %s: %v", boundKey.Hex(), err)
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
