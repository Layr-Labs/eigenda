package verification

import (
	contractEigenDABlobVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
)

// EigenDACert contains all data necessary to verify a Blob
type EigenDACert struct {
	BlobVerificationProof       contractEigenDABlobVerifier.BlobVerificationProofV2
	BatchHeader                 contractEigenDABlobVerifier.BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDABlobVerifier.NonSignerStakesAndSignature
}
