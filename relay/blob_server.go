package relay

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"golang.org/x/sync/errgroup"
)

// blobServer encapsulates logic for fetching blobs. Utilized by the relay Server.
// This struct adds caching (and perhaps eventually threading) on top of blobstore.BlobStore.
type blobServer struct {
	ctx    context.Context
	logger logging.Logger

	// blobStore can be used to read blobs from S3.
	blobStore *blobstore.BlobStore

	// blobCache is an LRU cache of blobs.
	blobCache cache.CachedAccessor[v2.BlobKey, []byte]

	// pool is a work pool for managing concurrent worker goroutines.
	pool *errgroup.Group
}

// NewBlobServer creates a new blobServer.
func NewBlobServer(
	ctx context.Context,
	logger logging.Logger,
	blobStore *blobstore.BlobStore,
	blobCacheSize int,
	workPoolSize int) (*blobServer, error) {

	pool, _ := errgroup.WithContext(ctx)
	pool.SetLimit(workPoolSize)

	server := &blobServer{
		ctx:       ctx,
		logger:    logger,
		blobStore: blobStore,
		pool:      pool,
	}

	cache, err := cache.NewCachedAccessor[v2.BlobKey, []byte](blobCacheSize, server.fetchBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating blob cache: %w", err)
	}
	server.blobCache = cache

	return server, nil
}

// GetBlob retrieves a blob from the blob store.
func (s *blobServer) GetBlob(blobKey v2.BlobKey) ([]byte, error) {
	data, err := s.blobCache.Get(blobKey)
	if err != nil {
		s.logger.Error("Failed to fetch blob: %v", err)
		return nil, err
	}

	return *data, nil
}

// fetchBlob retrieves a single blob from the blob store.
func (s *blobServer) fetchBlob(blobKey v2.BlobKey) (*[]byte, error) {
	data, err := s.blobStore.GetBlob(s.ctx, blobKey)
	if err != nil {
		s.logger.Error("Failed to fetch blob: %v", err)
		return nil, err
	}

	return &data, nil
}
