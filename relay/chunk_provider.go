package relay

import (
	"bytes"
	"context"
	"fmt"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/relay/cache"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"sync"
)

type chunkProvider struct {
	ctx    context.Context
	logger logging.Logger

	// metadataCache is an LRU cache of blob metadata. Each relay is authorized to serve data assigned to one or more
	// relay IDs. Blobs that do not belong to one of the relay IDs assigned to this server will not be in the cache.
	frameCache cache.CachedAccessor[blobKeyWithMetadata, []*encoding.Frame]

	// chunkReader is used to read chunks from the chunk store.
	chunkReader chunkstore.ChunkReader
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
	cacheSize int,
	maxIOConcurrency int) (*chunkProvider, error) {

	server := &chunkProvider{
		ctx:         ctx,
		logger:      logger,
		chunkReader: chunkReader,
	}

	c, err := cache.NewCachedAccessor[blobKeyWithMetadata, []*encoding.Frame](
		cacheSize,
		maxIOConcurrency,
		server.fetchFrames)
	if err != nil {
		return nil, err
	}
	server.frameCache = c

	return server, nil
}

// frameMap is a map of blob keys to frames.
type frameMap map[v2.BlobKey][]*encoding.Frame

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
		data []*encoding.Frame
		err  error
	}

	// Channel for results.
	completionChannel := make(chan *framesResult, len(keys))

	for _, key := range keys {

		boundKey := key
		go func() {
			frames, err := s.frameCache.Get(*boundKey)
			if err != nil {
				s.logger.Errorf("Failed to get frames for blob %v: %v", boundKey.blobKey, err)
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
			return nil, fmt.Errorf("error fetching frames for blob %v: %w", result.key, result.err)
		}
		fMap[result.key] = result.data
	}

	return fMap, nil
}

// fetchFrames retrieves the frames for a single blob.
func (s *chunkProvider) fetchFrames(key blobKeyWithMetadata) ([]*encoding.Frame, error) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	var proofs []*encoding.Proof
	var proofsErr error

	go func() {
		defer func() {
			wg.Done()
		}()
		proofs, proofsErr = s.chunkReader.GetChunkProofs(s.ctx, key.blobKey)
	}()

	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: key.metadata.totalChunkSizeBytes,
		FragmentSizeBytes:   key.metadata.fragmentSizeBytes,
	}

	coefficients, err := s.chunkReader.GetChunkCoefficients(s.ctx, key.blobKey, fragmentInfo)
	if err != nil {
		return nil, err
	}

	wg.Wait()
	if proofsErr != nil {
		return nil, proofsErr
	}

	frames, err := assembleFrames(coefficients, proofs)
	if err != nil {
		return nil, err
	}

	return frames, nil
}

// assembleFrames assembles a slice of frames from its composite proofs and coefficients.
func assembleFrames(frames []*rs.Frame, proofs []*encoding.Proof) ([]*encoding.Frame, error) {
	if len(frames) != len(proofs) {
		return nil, fmt.Errorf("number of frames and proofs must be equal (%d != %d)", len(frames), len(proofs))
	}

	assembledFrames := make([]*encoding.Frame, len(frames))
	for i := range frames {
		assembledFrames[i] = &encoding.Frame{
			Proof:  *proofs[i],
			Coeffs: frames[i].Coeffs,
		}
	}

	return assembledFrames, nil
}
