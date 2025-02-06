package relay

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/core"
	"sync"
	"time"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type chunkProvider struct {
	ctx    context.Context
	logger logging.Logger

	// frameCache contains encoding.Frame objects in a serialized form. This is much more memory efficient than
	// storing the frames in their parsed form. These frames can be deserialized via rs.DeserializeBinaryFrame().
	frameCache cache.CacheAccessor[blobKeyWithMetadata, *core.ChunksData]

	// chunkReader is used to read chunks from the chunk store.
	chunkReader chunkstore.ChunkReader

	// fetchTimeout is the maximum time to wait for a chunk proof fetch operation to complete.
	proofFetchTimeout time.Duration

	// coefficientFetchTimeout is the maximum time to wait for a chunk coefficient fetch operation to complete.
	coefficientFetchTimeout time.Duration
}

// blobKeyWithMetadata attaches some additional metadata to a blobKey.
type blobKeyWithMetadata struct {
	blobKey  v2.BlobKey
	metadata blobMetadata
}

func (m *blobKeyWithMetadata) Compare(other *blobKeyWithMetadata) int {
	return bytes.Compare(m.blobKey[:], other.blobKey[:])
}

// newChunkProvider creates a new chunkProvider.
func newChunkProvider(
	ctx context.Context,
	logger logging.Logger,
	chunkReader chunkstore.ChunkReader,
	cacheSize uint64,
	maxIOConcurrency int,
	proofFetchTimeout time.Duration,
	coefficientFetchTimeout time.Duration,
	metrics *cache.CacheAccessorMetrics) (*chunkProvider, error) {

	server := &chunkProvider{
		ctx:                     ctx,
		logger:                  logger,
		chunkReader:             chunkReader,
		proofFetchTimeout:       proofFetchTimeout,
		coefficientFetchTimeout: coefficientFetchTimeout,
	}

	cacheAccessor, err := cache.NewCacheAccessor[blobKeyWithMetadata, *core.ChunksData](
		cache.NewFIFOCache[blobKeyWithMetadata, *core.ChunksData](cacheSize, server.computeFramesCacheWeight),
		maxIOConcurrency,
		server.fetchFrames,
		metrics)
	if err != nil {
		return nil, err
	}
	server.frameCache = cacheAccessor

	return server, nil
}

// frameMap is a map of blob keys to binary frames.
type frameMap map[v2.BlobKey]*core.ChunksData

// computeFramesCacheWeight computes the 'weight' of the frames for the cache. The weight of a list of frames
// is equal to the size required to store the data, in bytes.
func (s *chunkProvider) computeFramesCacheWeight(_ blobKeyWithMetadata, frames *core.ChunksData) uint64 {
	return frames.Size()
}

// GetFrames retrieves the frames for a blob.
func (s *chunkProvider) GetFrames(ctx context.Context, mMap metadataMap) (frameMap, error) {

	if len(mMap) == 0 {
		return nil, fmt.Errorf("no metadata provided")
	}

	keys := make([]*blobKeyWithMetadata, 0, len(mMap))
	for k, v := range mMap {
		keys = append(keys, &blobKeyWithMetadata{blobKey: k, metadata: *v})
	}

	type framesResult struct {
		key  v2.BlobKey
		data *core.ChunksData
		err  error
	}

	// Channel for results.
	completionChannel := make(chan *framesResult, len(keys))

	for _, key := range keys {

		boundKey := key
		go func() {
			frames, err := s.frameCache.Get(ctx, *boundKey)
			if err != nil {
				s.logger.Errorf("Failed to get frames for blob %v: %v", boundKey.blobKey.Hex(), err)
				completionChannel <- &framesResult{
					key: boundKey.blobKey,
					err: err,
				}
			} else {
				completionChannel <- &framesResult{
					key:  boundKey.blobKey,
					data: frames,
				}
			}

		}()
	}

	fMap := make(frameMap, len(keys))
	for len(fMap) < len(keys) {
		result := <-completionChannel
		if result.err != nil {
			return nil, fmt.Errorf("error fetching frames for blob %v: %w", result.key.Hex(), result.err)
		}
		fMap[result.key] = result.data
	}

	return fMap, nil
}

// fetchFrames retrieves the frames for a single blob.
func (s *chunkProvider) fetchFrames(key blobKeyWithMetadata) (*core.ChunksData, error) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	var proofs [][]byte
	var proofsErr error

	go func() {
		ctx, cancel := context.WithTimeout(s.ctx, s.proofFetchTimeout)
		defer func() {
			wg.Done()
			cancel()
		}()

		proofs, proofsErr = s.chunkReader.GetBinaryFrameProofs(ctx, key.blobKey)
	}()

	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: key.metadata.totalChunkSizeBytes,
		FragmentSizeBytes:   key.metadata.fragmentSizeBytes,
	}

	ctx, cancel := context.WithTimeout(s.ctx, s.coefficientFetchTimeout)
	defer cancel()

	elementCount, coefficients, err := s.chunkReader.GetBinaryFrameCoefficients(ctx, key.blobKey, fragmentInfo)
	if err != nil {
		return nil, err
	}

	wg.Wait()
	if proofsErr != nil {
		return nil, proofsErr
	}

	frames, err := rs.BuildChunksData(proofs, int(elementCount), coefficients)
	if err != nil {
		return nil, err
	}

	return frames, nil
}
