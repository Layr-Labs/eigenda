package clients

import (
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
)

// EigenDAClientConfigV2 contains configuration values for EigenDAClientV2
type EigenDAClientConfigV2 struct {
	// The blob encoding version to use when writing blobs from the high level interface.
	PutBlobEncodingVersion codecs.BlobEncodingVersion

	// Point verification mode does an IFFT on data before it is written, and does an FFT on data after it is read.
	// This makes it possible to open points on the KZG commitment to prove that the field elements correspond to
	// the commitment. With this mode disabled, you will need to supply the entire blob to perform a verification
	// that any part of the data matches the KZG commitment.
	DisablePointVerificationMode bool
}
