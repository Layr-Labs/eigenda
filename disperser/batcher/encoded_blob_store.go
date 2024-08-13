package batcher

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type requestID string

type encodedBlobStore struct {
	mu sync.RWMutex

	requested map[requestID]struct{}
	encoded   map[requestID]*EncodingResult
	// encodedResultSize is the total size of all the chunks in the encoded results in bytes
	encodedResultSize uint64

	logger logging.Logger
}

// EncodingResult contains information about the encoding of a blob
type EncodingResult struct {
	BlobMetadata         *disperser.BlobMetadata
	ReferenceBlockNumber uint
	BlobQuorumInfo       *core.BlobQuorumInfo
	Commitment           *encoding.BlobCommitments
	ChunksData           *core.ChunksData
	Assignments          map[core.OperatorID]core.Assignment
}

// EncodingResultOrStatus is a wrapper for EncodingResult that also contains an error
type EncodingResultOrStatus struct {
	EncodingResult
	// Err is set if there was an error during encoding
	Err error
}

func newEncodedBlobStore(logger logging.Logger) *encodedBlobStore {
	return &encodedBlobStore{
		requested:         make(map[requestID]struct{}),
		encoded:           make(map[requestID]*EncodingResult),
		encodedResultSize: 0,
		logger:            logger,
	}
}

func (e *encodedBlobStore) PutEncodingRequest(blobKey disperser.BlobKey, quorumID core.QuorumID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	requestID := getRequestID(blobKey, quorumID)
	e.requested[requestID] = struct{}{}
}

func (e *encodedBlobStore) HasEncodingRequested(blobKey disperser.BlobKey, quorumID core.QuorumID, referenceBlockNumber uint) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	requestID := getRequestID(blobKey, quorumID)
	if _, ok := e.requested[requestID]; ok {
		return true
	}

	res, ok := e.encoded[requestID]
	if ok && res.ReferenceBlockNumber == referenceBlockNumber {
		return true
	}
	return false
}

func (e *encodedBlobStore) DeleteEncodingRequest(blobKey disperser.BlobKey, quorumID core.QuorumID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	requestID := getRequestID(blobKey, quorumID)
	if _, ok := e.requested[requestID]; !ok {
		return
	}

	delete(e.requested, requestID)
}

func (e *encodedBlobStore) PutEncodingResult(result *EncodingResult) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	blobKey := disperser.BlobKey{
		BlobHash:     result.BlobMetadata.BlobHash,
		MetadataHash: result.BlobMetadata.MetadataHash,
	}
	requestID := getRequestID(blobKey, result.BlobQuorumInfo.QuorumID)
	if _, ok := e.requested[requestID]; !ok {
		return fmt.Errorf("PutEncodedBlob: no such key (%s) in requested set", requestID)
	}

	if _, ok := e.encoded[requestID]; !ok {
		e.encodedResultSize += getChunksSize(result)
	}
	e.encoded[requestID] = result
	delete(e.requested, requestID)

	return nil
}

func (e *encodedBlobStore) GetEncodingResult(blobKey disperser.BlobKey, quorumID core.QuorumID) (*EncodingResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	requestID := getRequestID(blobKey, quorumID)
	if _, ok := e.encoded[requestID]; !ok {
		return nil, fmt.Errorf("GetEncodedBlob: no such key (%s) in encoded set", requestID)
	}

	return e.encoded[requestID], nil
}

func (e *encodedBlobStore) DeleteEncodingResult(blobKey disperser.BlobKey, quorumID core.QuorumID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	requestID := getRequestID(blobKey, quorumID)
	encodedResult, ok := e.encoded[requestID]
	if !ok {
		return
	}

	delete(e.encoded, requestID)
	e.encodedResultSize -= getChunksSize(encodedResult)
}

// PopLatestEncodingResults returns all the encoded results that are pending dispersal and deletes them along with stale results that are older than the given reference block
func (e *encodedBlobStore) PopLatestEncodingResults(refBlockNumber uint) []*EncodingResult {
	e.mu.Lock()
	defer e.mu.Unlock()

	fetched := make([]*EncodingResult, 0)
	staleCount := 0
	for k, encodedResult := range e.encoded {
		if encodedResult.ReferenceBlockNumber == refBlockNumber {
			fetched = append(fetched, encodedResult)
			// this is safe: https://go.dev/doc/effective_go#for
			delete(e.encoded, k)
			e.encodedResultSize -= getChunksSize(encodedResult)
		} else if encodedResult.ReferenceBlockNumber < refBlockNumber {
			delete(e.encoded, k)
			staleCount++
			e.encodedResultSize -= getChunksSize(encodedResult)
		} else {
			e.logger.Error("unexpected case", "refBlockNumber", encodedResult.ReferenceBlockNumber, "refBlockNumber", refBlockNumber)
		}
	}
	e.logger.Debug("consumed encoded results", "fetched", len(fetched), "stale", staleCount, "refBlockNumber", refBlockNumber, "encodedSize", e.encodedResultSize)

	return fetched
}

// GetNewAndDeleteStaleEncodingResults returns all the fresh encoded results that are pending dispersal, and deletes all the stale results that are older than the given block number
func (e *encodedBlobStore) GetNewAndDeleteStaleEncodingResults(blockNumber uint) []*EncodingResult {
	e.mu.Lock()
	defer e.mu.Unlock()
	fetched := make([]*EncodingResult, 0)
	staleCount := 0
	pendingConfirmation := 0
	for k, encodedResult := range e.encoded {
		if encodedResult.ReferenceBlockNumber == blockNumber {
			fetched = append(fetched, encodedResult)
		} else if encodedResult.ReferenceBlockNumber < blockNumber {
			// this is safe: https://go.dev/doc/effective_go#for
			delete(e.encoded, k)
			staleCount++
			e.encodedResultSize -= getChunksSize(encodedResult)
		} else {
			e.logger.Error("unexpected case", "refBlockNumber", encodedResult.ReferenceBlockNumber, "blockNumber", blockNumber)
		}
	}
	e.logger.Debug("consumed encoded results", "fetched", len(fetched), "stale", staleCount, "pendingConfirmation", pendingConfirmation, "blockNumber", blockNumber, "encodedSize", e.encodedResultSize)

	return fetched
}

// GetEncodedResultSize returns the total size of all the chunks in the encoded results in bytes
func (e *encodedBlobStore) GetEncodedResultSize() (int, uint64) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return len(e.encoded), e.encodedResultSize
}

func getRequestID(key disperser.BlobKey, quorumID core.QuorumID) requestID {
	return requestID(fmt.Sprintf("%s-%d", key.String(), quorumID))
}

// getChunksSize returns the total size of all the chunks in the encoded result in bytes
func getChunksSize(result *EncodingResult) uint64 {
	return result.ChunksData.Size()
}
