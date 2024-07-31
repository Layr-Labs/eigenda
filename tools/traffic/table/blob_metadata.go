package table

// BlobMetadata encapsulates various information about a blob written by the traffic generator.
type BlobMetadata struct {
	// key of the blob, set when the blob is initially uploaded.
	key *[]byte

	// checksum of the blob.
	checksum *[16]byte

	// batchHeaderHash of the blob in bytes.
	size uint

	// batchHeaderHash of the blob.
	batchHeaderHash *[]byte

	// blobIndex of the blob.
	blobIndex uint

	// remainingReadPermits describes the maximum number of remaining reads permitted against this blob.
	// If -1 then an unlimited number of reads are permitted.
	remainingReadPermits int

	// index describes the position of this blob within the blobTable.
	index uint
}

// NewBlobMetadata creates a new BlobMetadata instance. The readPermits parameter describes the maximum number of
// remaining reads permitted against this blob. If -1 then an unlimited number of reads are permitted.
func NewBlobMetadata(
	key *[]byte,
	checksum *[16]byte,
	size uint,
	batchHeaderHash *[]byte,
	blobIndex uint,
	readPermits int) *BlobMetadata {

	return &BlobMetadata{
		key:                  key,
		checksum:             checksum,
		size:                 size,
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

// Checksum returns the checksum of the blob.
func (blob *BlobMetadata) Checksum() *[16]byte {
	return blob.checksum
}

// Size returns the size of the blob, in bytes.
func (blob *BlobMetadata) Size() uint {
	return blob.size
}

// BatchHeaderHash returns the batchHeaderHash of the blob.
func (blob *BlobMetadata) BatchHeaderHash() *[]byte {
	return blob.batchHeaderHash
}

// BlobIndex returns the blobIndex of the blob.
func (blob *BlobMetadata) BlobIndex() uint {
	return blob.blobIndex
}

// RemainingReadPermits returns the maximum number of remaining reads permitted against this blob.
func (blob *BlobMetadata) RemainingReadPermits() int {
	return blob.remainingReadPermits
}
