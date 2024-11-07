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
)

// blobServer encapsulates logic for fetching blobs. Utilized by the relay Server.
// This struct adds caching and threading on top of blobstore.BlobStore.
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

// newBlobServer creates a new blobServer.
func newBlobServer(
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

	// Even though we don't need extra parallelism here, we still use the work pool to ensure that we don't
	// permit too many concurrent requests to the blob store.

	wg := sync.WaitGroup{}
	wg.Add(1)

	var data *[]byte
	var err error

	s.pool.Go(func() error {
		defer wg.Done()
		data, err = s.blobCache.Get(blobKey)
		return nil
	})

	wg.Wait()

	if err != nil {
		// It should not be possible for external users to force an error here since we won't
		// even call this method if the blob key is invalid (so it's ok to have a noisy log here).
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
