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
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

type chunkServer struct {
	ctx    context.Context
	logger logging.Logger

	// metadataCache is an LRU cache of blob metadata. Blobs that do not belong to one of the relay shards
	// assigned to this server will not be in the cache.
	frameCache cache.CachedAccessor[blobKeyWithMetadata, []*encoding.Frame]

	// chunkReader is used to read chunks from the chunk store.
	chunkReader chunkstore.ChunkReader

	// pool is a work pool for managing concurrent requests. Used to limit the number of concurrent
	// requests.
	pool *errgroup.Group
}

// blobKeyWithMetadata attaches some additional metadata to a blobKey.
type blobKeyWithMetadata struct {
	blobKey  v2.BlobKey
	metadata blobMetadata
}

func (m *blobKeyWithMetadata) Compare(other *blobKeyWithMetadata) int {
	return bytes.Compare(m.blobKey[:], other.blobKey[:])
}

// newChunkServer creates a new chunkServer.
func newChunkServer(
	ctx context.Context,
	logger logging.Logger,
	chunkReader chunkstore.ChunkReader,
	cacheSize int,
	workPoolSize int) (*chunkServer, error) {

	pool := &errgroup.Group{}
	pool.SetLimit(workPoolSize)

	server := &chunkServer{
		ctx:         ctx,
		logger:      logger,
		chunkReader: chunkReader,
		pool:        pool,
	}

	c, err := cache.NewCachedAccessor[blobKeyWithMetadata, []*encoding.Frame](cacheSize, server.fetchFrames)
	if err != nil {
		return nil, err
	}
	server.frameCache = c

	return server, nil
}

// frameMap is a map of blob keys to frames.
type frameMap map[v2.BlobKey][]*encoding.Frame

// GetFrames retrieves the frames for a blob.
func (s *chunkServer) GetFrames(ctx context.Context, mMap *metadataMap) (*frameMap, error) {

	keys := make([]*blobKeyWithMetadata, 0, len(*mMap))
	for k, v := range *mMap {
		keys = append(keys, &blobKeyWithMetadata{blobKey: k, metadata: *v})
	}

	fMap := make(frameMap, len(keys))
	hadError := atomic.Bool{}

	wg := sync.WaitGroup{}
	wg.Add(len(keys))

	for _, key := range keys {
		boundKey := key
		s.pool.Go(func() error {

			defer wg.Done()

			frames, err := s.frameCache.Get(*boundKey)
			if err != nil {
				s.logger.Error("Failed to get frames for blob %v: %v", boundKey.blobKey, err)
				hadError.Store(true)
			} else {
				fMap[boundKey.blobKey] = *frames
			}

			return nil
		})
	}

	wg.Wait()

	return &fMap, nil
}

// fetchFrames retrieves the frames for a single blob.
func (s *chunkServer) fetchFrames(key blobKeyWithMetadata) (*[]*encoding.Frame, error) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	var proofs []*encoding.Proof
	var proofsErr error

	s.pool.Go(
		func() error {
			defer wg.Done()
			proofs, proofsErr = s.chunkReader.GetChunkProofs(s.ctx, key.blobKey)
			return nil
		})

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

	return &frames, nil
}

// assembleFrames assembles a slice of frames from its composite proofs and coefficients.
func assembleFrames(frames []*rs.Frame, proof []*encoding.Proof) ([]*encoding.Frame, error) {
	if len(frames) != len(proof) {
		return nil, fmt.Errorf("number of frames and proofs must be equal (%d != %d)", len(frames), len(proof))
	}

	assembledFrames := make([]*encoding.Frame, len(frames))
	for i := range frames {
		assembledFrames[i] = &encoding.Frame{
			Proof:  *proof[i],
			Coeffs: frames[i].Coeffs,
		}
	}

	return assembledFrames, nil
}
