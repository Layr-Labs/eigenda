package relay

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	cachecommon "github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/common/tracing"
	"github.com/Layr-Labs/eigenda/core"
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

	cacheAccessor, err := cache.NewCacheAccessor[blobKeyWithMetadata, []*encoding.Frame](
		cache2.NewFIFOCache[blobKeyWithMetadata, []*encoding.Frame](cacheSize, server.computeFramesCacheWeight),
		maxIOConcurrency,
		server.fetchFrames,
		metrics)
	if err != nil {
		return nil, err
	}

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
	ctx, span := tracing.TraceOperation(ctx, "chunkProvider.GetFrames")
	defer span.End()

	if len(mMap) == 0 {
		return nil, fmt.Errorf("no metadata provided")
	}

	// Add span for key preparation
	_, keysSpan := tracing.TraceOperation(ctx, "chunkProvider.GetFrames.prepareKeys")
	keys := make([]*blobKeyWithMetadata, 0, len(mMap))
	for k, v := range mMap {
		keys = append(keys, &blobKeyWithMetadata{blobKey: k, metadata: *v})
	}
	keysSpan.End()

	type framesResult struct {
		key  v2.BlobKey
		data *core.ChunksData
		err  error
	}

	// Channel for results.
	completionChannel := make(chan *framesResult, len(keys))

	// Add span for goroutine launch phase
	_, launchSpan := tracing.TraceOperation(ctx, "chunkProvider.GetFrames.launchGoroutines")
	for _, key := range keys {
		boundKey := key
		go func() {
			routineCtx, routineSpan := tracing.TraceOperation(ctx, "chunkProvider.GetFrames.goroutine")
			defer routineSpan.End()

			// Add cache operation span with sub-spans for different phases
			cacheCtx, cacheSpan := tracing.TraceOperation(routineCtx, "chunkProvider.GetFrames.cacheOperation")
			defer cacheSpan.End()

			frames, err := s.frameCache.Get(cacheCtx, *boundKey)

			// Add error handling span if there's an error
			if err != nil {
				_, errSpan := tracing.TraceOperation(routineCtx, "chunkProvider.GetFrames.errorHandling")
				s.logger.Errorf("Failed to get frames for blob %v: %v", boundKey.blobKey.Hex(), err)
				completionChannel <- &framesResult{
					key: boundKey.blobKey,
					err: err,
				}
				errSpan.End()
			} else {
				_, resultSpan := tracing.TraceOperation(routineCtx, "chunkProvider.GetFrames.sendResult")
				completionChannel <- &framesResult{
					key:  boundKey.blobKey,
					data: frames,
				}
				resultSpan.End()
			}
		}()
	}
	launchSpan.End()

	// Add span for result collection phase
	_, collectSpan := tracing.TraceOperation(ctx, "chunkProvider.GetFrames.collectResults")
	fMap := make(frameMap, len(keys))
	for len(fMap) < len(keys) {
		result := <-completionChannel
		if result.err != nil {
			collectSpan.End()
			return nil, fmt.Errorf("error fetching frames for blob %v: %w", result.key.Hex(), result.err)
		}
		fMap[result.key] = result.data
	}
	collectSpan.End()

	return fMap, nil
}

// fetchFrames retrieves the frames for a single blob.
func (s *chunkProvider) fetchFrames(key blobKeyWithMetadata) (*core.ChunksData, error) {
	ctx, span := tracing.TraceOperation(s.ctx, "chunkProvider.fetchFrames")
	defer span.End()

	// Add span for setup phase
	_, setupSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.setup")
	wg := sync.WaitGroup{}
	wg.Add(1)
	var proofs [][]byte
	var proofsErr error
	setupSpan.End()

	// Add span for proof fetching goroutine launch
	_, proofLaunchSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.launchProofFetch")
	go func() {
		proofCtx, cancel := context.WithTimeout(ctx, s.proofFetchTimeout)
		proofCtx, proofSpan := tracing.TraceOperation(proofCtx, "chunkProvider.fetchProofs")
		defer func() {
			proofSpan.End()
			wg.Done()
			cancel()
		}()

		// Add span for the actual proof fetching operation
		_, fetchSpan := tracing.TraceOperation(proofCtx, "chunkProvider.fetchProofs.getBinaryChunkProofs")
		proofs, proofsErr = s.chunkReader.GetBinaryChunkProofs(proofCtx, key.blobKey)
		fetchSpan.End()
	}()
	proofLaunchSpan.End()

	// Add span for fragment info preparation
	_, fragSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.prepareFragmentInfo")
	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: key.metadata.totalChunkSizeBytes,
		FragmentSizeBytes:   key.metadata.fragmentSizeBytes,
	}
	fragSpan.End()

	// Add spans for coefficient fetching with timeout
	coeffCtx, cancel := context.WithTimeout(ctx, s.coefficientFetchTimeout)
	coeffCtx, coeffSpan := tracing.TraceOperation(coeffCtx, "chunkProvider.fetchCoefficients")
	defer func() {
		coeffSpan.End()
		cancel()
	}()

	// Add span for the actual coefficient fetching operation
	_, coeffFetchSpan := tracing.TraceOperation(coeffCtx, "chunkProvider.fetchCoefficients.getBinaryChunkCoefficients")
	elementCount, coefficients, err := s.chunkReader.GetBinaryChunkCoefficients(coeffCtx, key.blobKey, fragmentInfo)
	coeffFetchSpan.End()

	if err != nil {
		_, errSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.coefficientError")
		defer errSpan.End()
		return nil, err
	}

	// Add span for waiting on proof fetch completion
	_, waitSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.waitForProofs")
	wg.Wait()
	waitSpan.End()

	if proofsErr != nil {
		_, errSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.proofError")
		defer errSpan.End()
		return nil, proofsErr
	}

	// Add span for building chunks data
	_, buildSpan := tracing.TraceOperation(ctx, "chunkProvider.buildChunksData")
	frames, err := rs.BuildChunksData(proofs, int(elementCount), coefficients)
	buildSpan.End()

	if err != nil {
		_, errSpan := tracing.TraceOperation(ctx, "chunkProvider.fetchFrames.buildError")
		defer errSpan.End()
		return nil, err
	}

	return frames, nil
}
