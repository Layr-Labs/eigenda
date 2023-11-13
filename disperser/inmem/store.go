package inmem

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
)

// BlobStore is an in-memory implementation of the BlobStore interface
type BlobStore struct {
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
	}

	return blobKey, nil
}

func (q *BlobStore) GetBlobContent(ctx context.Context, blobHash disperser.BlobHash) ([]byte, error) {
	if holder, ok := q.Blobs[blobHash]; ok {
		return holder.Data, nil
	} else {
		return nil, disperser.ErrBlobNotFound
	}
}

func (q *BlobStore) MarkBlobConfirmed(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	blobKey := existingMetadata.GetBlobKey()
	if _, ok := q.Metadata[blobKey]; !ok {
		return nil, disperser.ErrBlobNotFound
	}
	newMetadata := *existingMetadata
	newMetadata.BlobStatus = disperser.Confirmed
	newMetadata.ConfirmationInfo = confirmationInfo
	q.Metadata[blobKey] = &newMetadata
	return &newMetadata, nil
}

func (q *BlobStore) MarkBlobInsufficientSignatures(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	blobKey := existingMetadata.GetBlobKey()
	if _, ok := q.Metadata[blobKey]; !ok {
		return nil, disperser.ErrBlobNotFound
	}
	newMetadata := *existingMetadata
	newMetadata.BlobStatus = disperser.InsufficientSignatures
	newMetadata.ConfirmationInfo = confirmationInfo
	q.Metadata[blobKey] = &newMetadata
	return &newMetadata, nil
}

func (q *BlobStore) MarkBlobFinalized(ctx context.Context, blobKey disperser.BlobKey) error {
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Finalized
	return nil
}

func (q *BlobStore) MarkBlobProcessing(ctx context.Context, blobKey disperser.BlobKey) error {
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Processing
	return nil
}

func (q *BlobStore) MarkBlobFailed(ctx context.Context, blobKey disperser.BlobKey) error {
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Failed
	return nil
}

func (q *BlobStore) IncrementBlobRetryCount(ctx context.Context, existingMetadata *disperser.BlobMetadata) error {
	if _, ok := q.Metadata[existingMetadata.GetBlobKey()]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[existingMetadata.GetBlobKey()].NumRetries++
	return nil
}

func (q *BlobStore) GetBlobsByMetadata(ctx context.Context, metadata []*disperser.BlobMetadata) (map[disperser.BlobKey]*core.Blob, error) {
	blobs := make(map[disperser.BlobKey]*core.Blob)
	for _, meta := range metadata {
		if holder, ok := q.Blobs[meta.BlobHash]; ok {
			blobs[meta.GetBlobKey()] = &core.Blob{
				RequestHeader: meta.RequestMetadata.BlobRequestHeader,
				Data:          holder.Data,
			}
		} else {
			return nil, disperser.ErrBlobNotFound
		}
	}
	return blobs, nil
}

func (q *BlobStore) GetBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) ([]*disperser.BlobMetadata, error) {
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.BlobStatus == status {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *BlobStore) GetMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*disperser.BlobMetadata, error) {
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash && meta.ConfirmationInfo.BlobIndex == blobIndex {
			return meta, nil
		}
	}

	return nil, disperser.ErrBlobNotFound
}

func (q *BlobStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*disperser.BlobMetadata, error) {
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *BlobStore) GetBlobMetadata(ctx context.Context, blobKey disperser.BlobKey) (*disperser.BlobMetadata, error) {
	if meta, ok := q.Metadata[blobKey]; ok {
		return meta, nil
	}
	return nil, disperser.ErrBlobNotFound
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
