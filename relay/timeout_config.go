package relay

import "time"

// TimeoutConfig encapsulates the timeout configuration for the relay server.
type TimeoutConfig struct {

	// The maximum time permitted for a GetChunks GRPC to complete. If zero then no timeout is enforced.
	GetChunksTimeout time.Duration

	// The maximum time permitted for a GetBlob GRPC to complete. If zero then no timeout is enforced.
	GetBlobTimeout time.Duration

	// The maximum time permitted for a single request to the metadata store to fetch the metadata
	// for an individual blob.
	InternalGetMetadataTimeout time.Duration

	// The maximum time permitted for a single request to the blob store to fetch a blob.
	InternalGetBlobTimeout time.Duration

	// The maximum time permitted for a single request to the chunk store to fetch chunk proofs.
	InternalGetProofsTimeout time.Duration

	// The maximum time permitted for a single request to the chunk store to fetch chunk coefficients.
	InternalGetCoefficientsTimeout time.Duration
}
