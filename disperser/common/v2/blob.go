package v2

import (
	core "github.com/Layr-Labs/eigenda/core/v2"
)

type BlobStatus uint

const (
	Queued BlobStatus = iota
	Encoded
	Certified
	Failed
)

// BlobMetadata is an internal representation of a blob's metadata.
type BlobMetadata struct {
	BlobHeader core.BlobHeader

	// BlobStatus indicates the current status of the blob
	BlobStatus BlobStatus
	// Expiry is Unix timestamp of the blob expiry in seconds from epoch
	Expiry uint64
	// NumRetries is the number of times the blob has been retried
	NumRetries uint
	// BlobSize is the size of the blob in bytes
	BlobSize uint64
	// RequestedAt is the Unix timestamp of when the blob was requested in seconds
	RequestedAt uint64
	// UpdatedAt is the Unix timestamp of when the blob was last updated in _nanoseconds_
	UpdatedAt uint64
}
