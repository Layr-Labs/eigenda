package verification

import (
	contractEigenDACertVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
)

// EigenDACert contains all data necessary to retrieve and validate a blob
//
// This struct represents the composition of a eigenDA blob certificate, as it would exist in a rollup inbox.
type EigenDACert struct {
	BlobInclusionInfo           contractEigenDACertVerifier.BlobInclusionInfo
	BatchHeader                 contractEigenDACertVerifier.BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDACertVerifier.NonSignerStakesAndSignature
}
