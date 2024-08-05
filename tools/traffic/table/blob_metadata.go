package table

import "errors"

// BlobMetadata encapsulates various information about a blob written by the traffic generator.
type BlobMetadata struct {
	// Key of the blob, set when the blob is initially uploaded.
	Key []byte

	// BlobIndex of the blob.
	BlobIndex uint

	// Checksum of the blob.
	Checksum [16]byte

	// Size of the blob, in bytes.
	Size uint

	// RemainingReadPermits describes the maximum number of remaining reads permitted against this blob.
	// If -1 then an unlimited number of reads are permitted.
	RemainingReadPermits int
}

// NewBlobMetadata creates a new BlobMetadata instance. The readPermits parameter describes the maximum number of
// remaining reads permitted against this blob. If -1 then an unlimited number of reads are permitted.
func NewBlobMetadata(
	key []byte,
	checksum [16]byte,
	size uint,
	blobIndex uint,
	readPermits int) (*BlobMetadata, error) {

	if readPermits == 0 {
		return nil, errors.New("read permits must not be zero")
	}

	return &BlobMetadata{
		Key:                  key,
		Checksum:             checksum,
		Size:                 size,
		BlobIndex:            blobIndex,
		RemainingReadPermits: readPermits,
	}, nil
}
