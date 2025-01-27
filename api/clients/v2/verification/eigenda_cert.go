package verification

import (
	contractEigenDABlobVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
)

// EigenDACert contains all data necessary to retrieve and validate a blob
//
// This struct represents the composition of a eigenDA blob certificate, as it would exist in a rollup inbox.
type EigenDACert struct {
	BlobInclusionInfo           contractEigenDABlobVerifier.BlobVerificationProofV2
	BatchHeader                 contractEigenDABlobVerifier.BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDABlobVerifier.NonSignerStakesAndSignature
}
