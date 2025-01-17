package verification

import (
	contractEigenDABlobVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
)

// EigenDACert contains all data necessary to verify a Blob
//
// This struct represents the composition of a eigenDA blob certificate, as it would exist in a rollup inbox.
// TODO: (litt3) it is possible that the exact types contained in this struct will change in the near future during
//  client v2 development.
type EigenDACert struct {
	BlobVerificationProof       contractEigenDABlobVerifier.BlobVerificationProofV2
	BatchHeader                 contractEigenDABlobVerifier.BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDABlobVerifier.NonSignerStakesAndSignature
}
