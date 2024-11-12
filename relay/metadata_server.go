package relay

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

// Metadata about a blob. The relay only needs a small subset of a blob's metadata.
// This struct adds caching and threading on top of blobstore.BlobMetadataStore.
type blobMetadata struct {
	// the size of the blob in bytes
	blobSizeBytes uint32
	// the size of the file containing the encoded chunks
	totalChunkSizeBytes uint32
	// the fragment size used for uploading the encoded chunks
	fragmentSizeBytes uint32
}

// metadataServer encapsulates logic for fetching metadata for blobs. Utilized by the relay Server.
type metadataServer struct {
	ctx    context.Context
	logger logging.Logger

	// metadataStore can be used to read blob metadata from dynamoDB.
	metadataStore *blobstore.BlobMetadataStore

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	metadataCache cache.CachedAccessor[v2.BlobKey, blobMetadata]

	// shardSet is the set of shards assigned to this relay. This relay will refuse to serve metadata for blobs
	// that are not assigned to one of these shards.
	shardSet map[v2.RelayKey]struct{}

	// pool is a work pool for managing concurrent metadata requests. Used to limit the number of concurrent
	// requests to the metadata store.
	pool *errgroup.Group
}

// newMetadataServer creates a new metadataServer.
func newMetadataServer(
	ctx context.Context,
	logger logging.Logger,
	metadataStore *blobstore.BlobMetadataStore,
	metadataCacheSize int,
	workPoolSize int,
	shards []v2.RelayKey) (*metadataServer, error) {

	shardSet := make(map[v2.RelayKey]struct{}, len(shards))
	for _, shard := range shards {
		shardSet[shard] = struct{}{}
	}

	pool, _ := errgroup.WithContext(ctx)
	pool.SetLimit(workPoolSize)

	server := &metadataServer{
		ctx:           ctx,
		logger:        logger,
		metadataStore: metadataStore,
		shardSet:      shardSet,
		pool:          pool,
	}

	metadataCache, err := cache.NewCachedAccessor[v2.BlobKey, blobMetadata](metadataCacheSize, server.fetchMetadata)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}

	server.metadataCache = metadataCache

	return server, nil
}

// metadataMap is a map of blob keys to metadata.
type metadataMap map[v2.BlobKey]*blobMetadata

// GetMetadataForBlobs retrieves metadata about multiple blobs in parallel.
func (m *metadataServer) GetMetadataForBlobs(keys []v2.BlobKey) (*metadataMap, error) {

	// TODO figure out how timeouts are going to work here

	mapLock := sync.Mutex{}
	mMap := make(metadataMap)
	hadError := atomic.Bool{}

	wg := sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		boundKey := key
		m.pool.Go(func() error {
			defer wg.Done()

			metadata, err := m.metadataCache.Get(boundKey)
			if err != nil {
				// Intentionally log at debug level. External users can force this condition to trigger
				// by requesting metadata for a blob that does not exist, and so it's important to avoid
				// allowing hooligans to spam the logs in production environments.
				m.logger.Debug("error retrieving metadata for blob %s: %v", boundKey.Hex(), err)
				hadError.Store(true)
				return nil
			}
			mapLock.Lock()
			mMap[boundKey] = metadata
			mapLock.Unlock()

			return nil
		})
	}
	wg.Wait()

	if hadError.Load() {
		return nil, fmt.Errorf("error retrieving metadata for one or more blobs")
	}

	return &mMap, nil
}

// fetchMetadata retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (m *metadataServer) fetchMetadata(key v2.BlobKey) (*blobMetadata, error) {
	// Retrieve the metadata from the store.
	cert, fragmentInfo, err := m.metadataStore.GetBlobCertificate(m.ctx, v2.BlobKey(key))
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	if len(m.shardSet) > 0 {
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
	}

	metadata := &blobMetadata{
		blobSizeBytes:       0, /* TODO: populate this once it is added to the metadata store */
		totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
		fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
	}

	return metadata, nil
}
