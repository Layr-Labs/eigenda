package relay

import (
	"context"
	"fmt"
	"sync/atomic"

	"time"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Metadata about a blob. The relay only needs a small subset of a blob's metadata.
// This struct adds caching and threading on top of blobstore.BlobMetadataStore.
type blobMetadata struct {
	// the size of the blob in bytes
	blobSizeBytes uint32
	// the size of each encoded chunk
	chunkSizeBytes uint32
	// the size of the file containing the encoded chunks
	totalChunkSizeBytes uint32
	// the fragment size used for uploading the encoded chunks
	fragmentSizeBytes uint32
}

// metadataProvider encapsulates logic for fetching metadata for blobs. Utilized by the relay Server.
type metadataProvider struct {
	ctx    context.Context
	logger logging.Logger

	// metadataStore can be used to read blob metadata from dynamoDB.
	metadataStore *blobstore.BlobMetadataStore

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	metadataCache cache.CacheAccessor[v2.BlobKey, *blobMetadata]

	// relayIDSet is the set of relay IDs assigned to this relay. This relay will refuse to serve metadata for blobs
	// that are not assigned to one of these IDs.
	relayIDSet map[v2.RelayKey]struct{}

	// fetchTimeout is the maximum time to wait for a metadata fetch operation to complete.
	fetchTimeout time.Duration

	// blobParamsMap is a map of blob version to blob version parameters.
	blobParamsMap atomic.Pointer[v2.BlobVersionParameterMap]
}

// newMetadataProvider creates a new metadataProvider.
func newMetadataProvider(
	ctx context.Context,
	logger logging.Logger,
	metadataStore *blobstore.BlobMetadataStore,
	metadataCacheSize int,
	maxIOConcurrency int,
	relayIDs []v2.RelayKey,
	fetchTimeout time.Duration,
	blobParamsMap *v2.BlobVersionParameterMap,
	metrics *cache.CacheAccessorMetrics) (*metadataProvider, error) {

	relayIDSet := make(map[v2.RelayKey]struct{}, len(relayIDs))
	for _, id := range relayIDs {
		relayIDSet[id] = struct{}{}
	}

	server := &metadataProvider{
		ctx:           ctx,
		logger:        logger,
		metadataStore: metadataStore,
		relayIDSet:    relayIDSet,
		fetchTimeout:  fetchTimeout,
	}
	server.blobParamsMap.Store(blobParamsMap)

	metadataCache, err := cache.NewCacheAccessor[v2.BlobKey, *blobMetadata](
		cache.NewFIFOCache[v2.BlobKey, *blobMetadata](uint64(metadataCacheSize), nil),
		maxIOConcurrency,
		server.fetchMetadata,
		metrics)
	if err != nil {
		return nil, fmt.Errorf("error creating metadata cache: %w", err)
	}

	server.metadataCache = metadataCache

	return server, nil
}

// metadataMap is a map of blob keys to metadata.
type metadataMap map[v2.BlobKey]*blobMetadata

// GetMetadataForBlobs retrieves metadata about multiple blobs in parallel.
// If any of the blobs do not exist, an error is returned.
// Note that resulting metadata map may not have the same length as the input
// keys slice if the input keys slice has duplicate items.
func (m *metadataProvider) GetMetadataForBlobs(ctx context.Context, keys []v2.BlobKey) (metadataMap, error) {

	// blobMetadataResult is the result of a metadata fetch operation.
	type blobMetadataResult struct {
		key      v2.BlobKey
		metadata *blobMetadata
		err      error
	}

	// Completed operations will send a result to this channel.
	completionChannel := make(chan *blobMetadataResult, len(keys))

	// Set when the first error is encountered. Useful for preventing new operations from starting.
	hadError := atomic.Bool{}

	mMap := make(metadataMap)
	for _, key := range keys {
		mMap[key] = nil
	}

	for key := range mMap {
		if hadError.Load() {
			// Don't bother starting new operations if we've already encountered an error.
			break
		}

		boundKey := key
		go func() {
			metadata, err := m.metadataCache.Get(ctx, boundKey)
			if err != nil {
				// Intentionally log at debug level. External users can force this condition to trigger
				// by requesting metadata for a blob that does not exist, and so it's important to avoid
				// allowing hooligans to spam the logs in production environments.
				m.logger.Debugf("error retrieving metadata for blob %s: %v", boundKey.Hex(), err)
				hadError.Store(true)
				completionChannel <- &blobMetadataResult{
					key: boundKey,
					err: err,
				}
			}

			completionChannel <- &blobMetadataResult{
				key:      boundKey,
				metadata: metadata,
			}
		}()
	}

	for range mMap {
		result := <-completionChannel
		if result.err != nil {
			return nil, fmt.Errorf("error fetching metadata for blob %s: %w", result.key.Hex(), result.err)
		}
		mMap[result.key] = result.metadata
	}

	return mMap, nil
}

func (m *metadataProvider) UpdateBlobVersionParameters(blobParamsMap *v2.BlobVersionParameterMap) {
	m.blobParamsMap.Store(blobParamsMap)
}

// fetchMetadata retrieves metadata about a blob. Fetches from the cache if available, otherwise from the store.
func (m *metadataProvider) fetchMetadata(key v2.BlobKey) (*blobMetadata, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.fetchTimeout)
	defer cancel()

	blobParamsMap := m.blobParamsMap.Load()
	if blobParamsMap == nil {
		return nil, fmt.Errorf("blob version parameters is nil")
	}

	// Retrieve the metadata from the store.
	cert, fragmentInfo, err := m.metadataStore.GetBlobCertificate(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving metadata for blob %s: %w", key.Hex(), err)
	}

	if len(m.relayIDSet) > 0 {
		validShard := false
		for _, shard := range cert.RelayKeys {
			if _, ok := m.relayIDSet[shard]; ok {
				validShard = true
				break
			}
		}

		if !validShard {
			return nil, fmt.Errorf("blob %s is not assigned to this relay", key.Hex())
		}
	}

	// TODO(cody-littley): blob size is not correct https://github.com/Layr-Labs/eigenda/pull/906#discussion_r1847396530
	blobSize := uint32(cert.BlobHeader.BlobCommitments.Length) * encoding.BYTES_PER_SYMBOL
	blobParams, ok := blobParamsMap.Get(cert.BlobHeader.BlobVersion)
	if !ok {
		return nil, fmt.Errorf("blob version %d not found in blob params map", cert.BlobHeader.BlobVersion)
	}
	chunkSize, err := v2.GetChunkLength(blobSize, blobParams)
	chunkSize *= encoding.BYTES_PER_SYMBOL
	if err != nil {
		return nil, fmt.Errorf("error getting chunk length: %w", err)
	}

	metadata := &blobMetadata{
		blobSizeBytes:       blobSize,
		chunkSizeBytes:      chunkSize,
		totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
		fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
	}

	return metadata, nil
}
