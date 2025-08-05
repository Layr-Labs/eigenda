package relay

import (
	"context"
	"fmt"
	"time"

	cache2 "github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// blobProvider encapsulates logic for fetching blobs. Utilized by the relay Server.
// This struct adds caching and concurrency limitation on top of blobstore.BlobStore.
type blobProvider struct {
	ctx    context.Context
	logger logging.Logger

	// blobStore is used to read blobs from S3.
	blobStore *blobstore.BlobStore

	// blobCache is an LRU cache of blobs.
	blobCache cache.CacheAccessor[v2.BlobKey, []byte]

	// fetchTimeout is the maximum time to wait for a blob fetch operation to complete.
	fetchTimeout time.Duration
}

// newBlobProvider creates a new blobProvider.
func newBlobProvider(
	ctx context.Context,
	logger logging.Logger,
	blobStore *blobstore.BlobStore,
	blobCacheSize uint64,
	maxIOConcurrency int,
	fetchTimeout time.Duration,
	metrics *cache.CacheAccessorMetrics) (*blobProvider, error) {

	server := &blobProvider{
		ctx:          ctx,
		logger:       logger,
		blobStore:    blobStore,
		fetchTimeout: fetchTimeout,
	}

	cacheAccessor, err := cache.NewCacheAccessor[v2.BlobKey, []byte](
		cache2.NewFIFOCache[v2.BlobKey, []byte](blobCacheSize, computeBlobCacheWeight, nil),
		maxIOConcurrency,
		server.fetchBlob,
		metrics)

	if err != nil {
		return nil, fmt.Errorf("error creating blob cache: %w", err)
	}
	server.blobCache = cacheAccessor

	return server, nil
}

// computeChunkCacheWeight computes the 'weight' of the blob for the cache. The weight of a blob
// is equal to its size, in bytes.
func computeBlobCacheWeight(_ v2.BlobKey, value []byte) uint64 {
	return uint64(len(value))
}

// GetBlob retrieves a blob from the blob store.
func (s *blobProvider) GetBlob(ctx context.Context, blobKey v2.BlobKey) ([]byte, error) {
	data, err := s.blobCache.Get(ctx, blobKey)

	if err != nil {
		return nil, fmt.Errorf("error calling blobCache.Get: %v", err)
	}

	return data, nil
}

// fetchBlob retrieves a single blob from the blob store.
func (s *blobProvider) fetchBlob(blobKey v2.BlobKey) ([]byte, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.fetchTimeout)
	defer cancel()

	data, err := s.blobStore.GetBlob(ctx, blobKey)
	if err != nil {
		return nil, fmt.Errorf("error calling blobStore.GetBlob: %v", err)
	}

	return data, nil
}
