package inmem

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common"
)

// BlobStore is an in-memory implementation of the BlobStore interface
type BlobStore struct {
	mu       sync.RWMutex
	Blobs    map[disperser.BlobHash]*BlobHolder
	Metadata map[disperser.BlobKey]*disperser.BlobMetadata
}

// BlobHolder stores the blob along with its status and any other metadata
type BlobHolder struct {
	Data []byte
}

var _ disperser.BlobStore = (*BlobStore)(nil)

// NewBlobStore creates an empty BlobStore
func NewBlobStore() disperser.BlobStore {
	return &BlobStore{
		Blobs:    make(map[disperser.BlobHash]*BlobHolder),
		Metadata: make(map[disperser.BlobKey]*disperser.BlobMetadata),
	}
}

func (q *BlobStore) StoreBlob(ctx context.Context, blob *core.Blob, requestedAt uint64) (disperser.BlobKey, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	blobKey := disperser.BlobKey{}
	// Generate the blob key
	blobHash, err := q.getNewBlobHash()
	if err != nil {
		return blobKey, err
	}
	blobKey.BlobHash = blobHash
	blobKey.MetadataHash = getMetadataHash(requestedAt)

	// Add the blob to the queue
	q.Blobs[blobHash] = &BlobHolder{
		Data: blob.Data,
	}

	q.Metadata[blobKey] = &disperser.BlobMetadata{
		BlobHash:     blobHash,
		MetadataHash: blobKey.MetadataHash,
		BlobStatus:   disperser.Processing,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          uint(len(blob.Data)),
			RequestedAt:       requestedAt,
		},
		Expiry: requestedAt + uint64(time.Hour),
	}

	return blobKey, nil
}

func (q *BlobStore) GetBlobContent(ctx context.Context, blobHash disperser.BlobHash) ([]byte, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if holder, ok := q.Blobs[blobHash]; ok {
		return holder.Data, nil
	} else {
		return nil, common.ErrBlobNotFound
	}
}

func (q *BlobStore) MarkBlobConfirmed(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// TODO (ian-shim): remove this check once we are sure that the metadata is never overwritten
	refreshedMetadata, err := q.GetBlobMetadata(ctx, existingMetadata.GetBlobKey())
	if err != nil {
		return nil, err
	}
	alreadyConfirmed, _ := refreshedMetadata.IsConfirmed()
	if alreadyConfirmed {
		return refreshedMetadata, nil
	}
	blobKey := existingMetadata.GetBlobKey()
	if _, ok := q.Metadata[blobKey]; !ok {
		return nil, common.ErrBlobNotFound
	}
	newMetadata := *existingMetadata
	newMetadata.BlobStatus = disperser.Confirmed
	newMetadata.ConfirmationInfo = confirmationInfo
	q.Metadata[blobKey] = &newMetadata
	return &newMetadata, nil
}

func (q *BlobStore) MarkBlobDispersing(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return common.ErrBlobNotFound
	}
	q.Metadata[blobKey].BlobStatus = disperser.Dispersing
	return nil
}

func (q *BlobStore) MarkBlobInsufficientSignatures(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	blobKey := existingMetadata.GetBlobKey()
	if _, ok := q.Metadata[blobKey]; !ok {
		return nil, common.ErrBlobNotFound
	}
	newMetadata := *existingMetadata
	newMetadata.BlobStatus = disperser.InsufficientSignatures
	newMetadata.ConfirmationInfo = confirmationInfo
	q.Metadata[blobKey] = &newMetadata
	return &newMetadata, nil
}

func (q *BlobStore) MarkBlobFinalized(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return common.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Finalized
	return nil
}

func (q *BlobStore) MarkBlobProcessing(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return common.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Processing
	return nil
}

func (q *BlobStore) MarkBlobFailed(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return common.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Failed
	return nil
}

func (q *BlobStore) IncrementBlobRetryCount(ctx context.Context, existingMetadata *disperser.BlobMetadata) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[existingMetadata.GetBlobKey()]; !ok {
		return common.ErrBlobNotFound
	}

	q.Metadata[existingMetadata.GetBlobKey()].NumRetries++
	return nil
}

func (q *BlobStore) UpdateConfirmationBlockNumber(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationBlockNumber uint32) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[existingMetadata.GetBlobKey()]; !ok {
		return common.ErrBlobNotFound
	}

	if q.Metadata[existingMetadata.GetBlobKey()].ConfirmationInfo == nil {
		return fmt.Errorf("cannot update confirmation block number for blob without confirmation info: %s", existingMetadata.GetBlobKey().String())
	}

	q.Metadata[existingMetadata.GetBlobKey()].ConfirmationInfo.ConfirmationBlockNumber = confirmationBlockNumber
	return nil
}

func (q *BlobStore) GetBlobsByMetadata(ctx context.Context, metadata []*disperser.BlobMetadata) (map[disperser.BlobKey]*core.Blob, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	blobs := make(map[disperser.BlobKey]*core.Blob)
	for _, meta := range metadata {
		if holder, ok := q.Blobs[meta.BlobHash]; ok {
			blobs[meta.GetBlobKey()] = &core.Blob{
				RequestHeader: meta.RequestMetadata.BlobRequestHeader,
				Data:          holder.Data,
			}
		} else {
			return nil, common.ErrBlobNotFound
		}
	}
	return blobs, nil
}

func (q *BlobStore) GetBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) ([]*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.BlobStatus == status {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *BlobStore) GetBlobMetadataByStatusWithPagination(ctx context.Context, status disperser.BlobStatus, limit int32, exclusiveStartKey *disperser.BlobStoreExclusiveStartKey) ([]*disperser.BlobMetadata, *disperser.BlobStoreExclusiveStartKey, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	foundStart := exclusiveStartKey == nil

	keys := make([]disperser.BlobKey, len(q.Metadata))
	i := 0
	for k := range q.Metadata {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return q.Metadata[keys[i]].Expiry < q.Metadata[keys[j]].Expiry
	})
	for _, key := range keys {
		meta := q.Metadata[key]
		if meta.BlobStatus == status {
			if foundStart {
				metas = append(metas, meta)
				if len(metas) == int(limit) {
					return metas, &disperser.BlobStoreExclusiveStartKey{
						BlobStatus: int32(meta.BlobStatus),
						Expiry:     int64(meta.Expiry),
					}, nil
				}
			} else if meta.BlobStatus == disperser.BlobStatus(exclusiveStartKey.BlobStatus) && meta.Expiry > uint64(exclusiveStartKey.Expiry) {
				foundStart = true // Found the starting point, start appending metas from next item
				metas = append(metas, meta)
				if len(metas) == int(limit) {
					return metas, &disperser.BlobStoreExclusiveStartKey{
						BlobStatus: int32(meta.BlobStatus),
						Expiry:     int64(meta.Expiry),
					}, nil
				}
			}
		}
	}

	// Return all the metas if limit is not reached
	return metas, nil, nil
}

func (q *BlobStore) GetMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash && meta.ConfirmationInfo.BlobIndex == blobIndex {
			return meta, nil
		}
	}

	return nil, common.ErrBlobNotFound
}

func (q *BlobStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *BlobStore) GetAllBlobMetadataByBatchWithPagination(ctx context.Context, batchHeaderHash [32]byte, limit int32, exclusiveStartKey *disperser.BatchIndexExclusiveStartKey) ([]*disperser.BlobMetadata, *disperser.BatchIndexExclusiveStartKey, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	foundStart := exclusiveStartKey == nil

	keys := make([]disperser.BlobKey, 0, len(q.Metadata))
	for k, v := range q.Metadata {
		if v.ConfirmationInfo != nil && v.ConfirmationInfo.BatchHeaderHash == batchHeaderHash {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return q.Metadata[keys[i]].ConfirmationInfo.BlobIndex < q.Metadata[keys[j]].ConfirmationInfo.BlobIndex
	})

	for _, key := range keys {
		meta := q.Metadata[key]
		if foundStart {
			metas = append(metas, meta)
			if len(metas) == int(limit) {
				return metas, &disperser.BatchIndexExclusiveStartKey{
					BatchHeaderHash: meta.ConfirmationInfo.BatchHeaderHash[:],
					BlobIndex:       meta.ConfirmationInfo.BlobIndex,
				}, nil
			}
		} else if exclusiveStartKey != nil && meta.ConfirmationInfo.BlobIndex > uint32(exclusiveStartKey.BlobIndex) {
			foundStart = true
			metas = append(metas, meta)
			if len(metas) == int(limit) {
				return metas, &disperser.BatchIndexExclusiveStartKey{
					BatchHeaderHash: meta.ConfirmationInfo.BatchHeaderHash[:],
					BlobIndex:       meta.ConfirmationInfo.BlobIndex,
				}, nil
			}
		}
	}

	// Return all the metas if limit is not reached
	return metas, nil, nil
}

func (q *BlobStore) GetBlobMetadata(ctx context.Context, blobKey disperser.BlobKey) (*disperser.BlobMetadata, error) {
	if meta, ok := q.Metadata[blobKey]; ok {
		return meta, nil
	}
	return nil, common.ErrBlobNotFound
}

func (q *BlobStore) GetBulkBlobMetadata(ctx context.Context, blobKeys []disperser.BlobKey) ([]*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, len(blobKeys))
	for i, key := range blobKeys {
		if meta, ok := q.Metadata[key]; ok {
			metas[i] = meta
		}
	}
	return metas, nil
}

func (q *BlobStore) HandleBlobFailure(ctx context.Context, metadata *disperser.BlobMetadata, maxRetry uint) (bool, error) {
	if metadata.NumRetries < maxRetry {
		if err := q.MarkBlobProcessing(ctx, metadata.GetBlobKey()); err != nil {
			return true, err
		}
		return true, q.IncrementBlobRetryCount(ctx, metadata)
	} else {
		return false, q.MarkBlobFailed(ctx, metadata.GetBlobKey())
	}
}

// getNewBlobHash generates a new blob key
func (q *BlobStore) getNewBlobHash() (disperser.BlobHash, error) {
	var key disperser.BlobHash
	for {
		buf := [32]byte{}
		// then we can call rand.Read.
		_, err := rand.Read(buf[:])
		if err != nil {
			return "", err
		}

		key = disperser.BlobHash(hex.EncodeToString(buf[:]))
		// If the key is already in use, try again
		if _, used := q.Blobs[key]; !used {
			break
		}
	}

	return key, nil
}

func getMetadataHash(requestedAt uint64) string {
	return strconv.FormatUint(requestedAt, 10)
}
