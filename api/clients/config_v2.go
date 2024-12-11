package clients

import (
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
)

// VerificationMode is an enum that represents the different ways that a blob may be encoded/decoded between
// the client and the disperser.
type VerificationMode uint

const (
	// TODO: write good docs here for IFFT and NoIFFT (I need to update my understanding to be able to write this)
	IFFT VerificationMode = iota
	NoIFFT
)

// EigenDAClientConfigV2 contains configuration values for EigenDAClientV2
type EigenDAClientConfigV2 struct {
	// The blob encoding version to use when writing and reading blobs
	BlobEncodingVersion codecs.BlobEncodingVersion

	// If PointVerificationMode is IFFT, then the client codec will do an IFFT on blobs before they are dispersed, and
	// will do an FFT on blobs after receiving them. This makes it possible to open points on the KZG commitment to prove
	// that the field elements correspond to the commitment.
	//
	// If PointVerificationMode is NoIFFT, the blob must be supplied in its entirety, to perform a verification
	// that any part of the data matches the KZG commitment.
	PointVerificationMode VerificationMode
}
