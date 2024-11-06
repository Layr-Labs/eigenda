package relay

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

// blobServer encapsulates logic for fetching blobs. Utilized by the relay Server.
type blobServer struct {
	ctx    context.Context
	logger logging.Logger

	// blobStore can be used to read blobs from S3.
	blobStore *blobstore.BlobStore

	// blobCache is an LRU cache of blobs.
	blobCache CachedAccessor[v2.BlobKey, []byte]

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

	cache, err := NewCachedAccessor[v2.BlobKey, []byte](blobCacheSize, server.fetchBlob)
	if err != nil {
		return nil, fmt.Errorf("error creating blob cache: %w", err)
	}
	server.blobCache = cache

	return server, nil
}

// GetBlobs fetches blobs from the cache or the blob store.
func (s *blobServer) GetBlobs(
	metadataMap *map[v2.BlobKey]*blobMetadata) (*map[v2.BlobKey][]byte, error) {

	dataMap := make(map[v2.BlobKey][]byte)
	hadError := atomic.Bool{}

	wg := sync.WaitGroup{}
	wg.Add(len(*metadataMap))

	for blobKey := range *metadataMap {
		boundKey := blobKey
		s.pool.Go(func() error {

			data, err := s.blobCache.Get(boundKey)
			if err != nil {
				s.logger.Error("Failed to fetch blob: %v", err)
				hadError.Store(true)
				return err
			}

			dataMap[boundKey] = *data

			return nil
		})
	}

	wg.Wait()
	if hadError.Load() {
		return nil, errors.New("failed to fetch one or more blobs")
	}

	return &dataMap, nil
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
