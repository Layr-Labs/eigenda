package relay

import "time"

// TimeoutConfig encapsulates the timeout configuration for the relay server.
type TimeoutConfig struct {

	// Timeout for GetChunks()
	GetChunksTimeout time.Duration

	// Timeout for GetBlob()
	GetBlobTimeout time.Duration

	// Timeout for internal metadata fetch
	InternalGetMetadataTimeout time.Duration

	// Timeout for internal blob fetch
	InternalGetBlobTimeout time.Duration

	// Timeout for internal proofs fetch
	InternalGetProofsTimeout time.Duration

	// Timeout for internal coefficients fetch
	InternalGetCoefficientsTimeout time.Duration
}
