package table

// BlobMetadata encapsulates various information about a blob written by the traffic generator.
type BlobMetadata struct {
	// key of the blob, set when the blob is initially uploaded.
	key []byte

	// checksum of the blob.
	checksum *[16]byte

	// size of the blob, in bytes.
	size uint

	// blobIndex of the blob.
	blobIndex uint

	// remainingReadPermits describes the maximum number of remaining reads permitted against this blob.
	// If -1 then an unlimited number of reads are permitted.
	remainingReadPermits int
}

// NewBlobMetadata creates a new BlobMetadata instance. The readPermits parameter describes the maximum number of
// remaining reads permitted against this blob. If -1 then an unlimited number of reads are permitted.
func NewBlobMetadata(
	key []byte,
	checksum *[16]byte,
	size uint,
	blobIndex uint,
	readPermits int) *BlobMetadata {

	if readPermits == 0 {
		panic("readPermits must be greater than 0, or -1 for unlimited reads")
	}

	return &BlobMetadata{
		key:                  key,
		checksum:             checksum,
		size:                 size,
		blobIndex:            blobIndex,
		remainingReadPermits: readPermits,
	}
}

// Key returns the key of the blob.
func (blob *BlobMetadata) Key() []byte {
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

// BlobIndex returns the blobIndex of the blob.
func (blob *BlobMetadata) BlobIndex() uint {
	return blob.blobIndex
}

// RemainingReadPermits returns the maximum number of remaining reads permitted against this blob.
func (blob *BlobMetadata) RemainingReadPermits() int {
	return blob.remainingReadPermits
}
