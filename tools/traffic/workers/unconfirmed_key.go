package workers

import (
	"time"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// UncertifiedKey is a Key that has not yet been certified by the disperser service.
type UncertifiedKey struct {
	// The Key of the blob.
	Key v2.BlobKey
	// The Size of the blob in bytes.
	Size uint
	// The Checksum of the blob.
	Checksum [16]byte
	// The time the blob was submitted to the disperser service.
	SubmissionTime time.Time
}
