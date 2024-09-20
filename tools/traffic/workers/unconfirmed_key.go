package workers

import "time"

// UnconfirmedKey is a Key that has not yet been confirmed by the disperser service.
type UnconfirmedKey struct {
	// The Key of the blob.
	Key []byte
	// The Size of the blob in bytes.
	Size uint
	// The Checksum of the blob.
	Checksum [16]byte
	// The time the blob was submitted to the disperser service.
	SubmissionTime time.Time
}
