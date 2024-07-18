package traffic

// BlobMetadata encapsulates various information about a blob written by the traffic generator.
type BlobMetadata struct {
	// key of the blob, set when the blob is initially uploaded.
	key *[]byte

	// batchHeaderHash of the blob.
	batchHeaderHash *[]byte

	// blobIndex of the blob.
	blobIndex uint32

	// remainingReadPermits describes the maximum number of remaining reads permitted against this blob.
	// If -1 then an unlimited number of reads are permitted.
	remainingReadPermits int32

	// index describes the position of this blob within the blobTable.
	index uint32
}

// NewBlobMetadata creates a new BlobMetadata instance. The readPermits parameter describes the maximum number of
// remaining reads permitted against this blob. If -1 then an unlimited number of reads are permitted.
func NewBlobMetadata(
	key *[]byte,
	batchHeaderHash *[]byte,
	blobIndex uint32,
	readPermits int32) *BlobMetadata {

	return &BlobMetadata{
		key:                  key,
		batchHeaderHash:      batchHeaderHash,
		blobIndex:            blobIndex,
		remainingReadPermits: readPermits,
		index:                0,
	}
}

// Key returns the key of the blob.
func (blob *BlobMetadata) Key() *[]byte {
	return blob.key
}

// BatchHeaderHash returns the batchHeaderHash of the blob.
func (blob *BlobMetadata) BatchHeaderHash() *[]byte {
	return blob.batchHeaderHash
}

// BlobIndex returns the blobIndex of the blob.
func (blob *BlobMetadata) BlobIndex() uint32 {
	return blob.blobIndex
}

// TODO method for decrementing read permits
