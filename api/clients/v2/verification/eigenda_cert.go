package verification

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"

	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	contractEigenDABlobVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
)

// EigenDACert contains all data necessary to retrieve and validate a Blob
//
// This struct represents the composition of a eigenDA blob certificate, as it would exist in a rollup inbox.
type EigenDACert struct {
	BlobVerificationProof       contractEigenDABlobVerifier.BlobVerificationProofV2
	BatchHeader                 contractEigenDABlobVerifier.BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDABlobVerifier.NonSignerStakesAndSignature
}

// BuildEigenDACert creates a new EigenDACert from a BlobInclusionInfo, BatchHeader, and NonSignerStakesAndSignature
//
// For convenience, this function accepts arguments as protobufs where applicable, since that's the form the caller
// will have.
func BuildEigenDACert(
	blobInclusionInfo *disperser.BlobInclusionInfo,
	batchHeader *commonv2.BatchHeader,
	nonSignerStakesAndSignature *contractEigenDABlobVerifier.NonSignerStakesAndSignature,
) (*EigenDACert, error) {

	bindingVerificationProof, err := VerificationProofProtoToBinding(blobInclusionInfo)
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	bindingBatchHeader, err := BatchHeaderProtoToBinding(batchHeader)
	if err != nil {
		return nil, fmt.Errorf("convert batch header to binding: %w", err)
	}

	return &EigenDACert{
		BlobVerificationProof:       *bindingVerificationProof,
		BatchHeader:                 *bindingBatchHeader,
		NonSignerStakesAndSignature: *nonSignerStakesAndSignature,
	}, nil
}
