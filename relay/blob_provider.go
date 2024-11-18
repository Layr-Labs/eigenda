package relay

import (
	"context"
	"fmt"
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
	blobCache cache.CachedAccessor[v2.BlobKey, []byte]
}

// newBlobProvider creates a new blobProvider.
func newBlobProvider(
	ctx context.Context,
	logger logging.Logger,
	blobStore *blobstore.BlobStore,
	blobCacheSize int,
	maxIOConcurrency int) (*blobProvider, error) {

	server := &blobProvider{
		ctx:       ctx,
		logger:    logger,
		blobStore: blobStore,
	}

	c, err := cache.NewCachedAccessor[v2.BlobKey, []byte](blobCacheSize, maxIOConcurrency, server.fetchBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating blob cache: %w", err)
	}
	server.blobCache = c

	return server, nil
}

// GetBlob retrieves a blob from the blob store.
func (s *blobProvider) GetBlob(blobKey v2.BlobKey) ([]byte, error) {

	data, err := s.blobCache.Get(blobKey)

	if err != nil {
		// It should not be possible for external users to force an error here since we won't
		// even call this method if the blob key is invalid (so it's ok to have a noisy log here).
		s.logger.Errorf("Failed to fetch blob: %v", err)
		return nil, err
	}

	return data, nil
}

// fetchBlob retrieves a single blob from the blob store.
func (s *blobProvider) fetchBlob(blobKey v2.BlobKey) ([]byte, error) {
	data, err := s.blobStore.GetBlob(s.ctx, blobKey)
	if err != nil {
		s.logger.Errorf("Failed to fetch blob: %v", err)
		return nil, err
	}

	return data, nil
}
