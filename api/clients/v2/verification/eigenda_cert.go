package verification

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"

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

// BuildEigenDACert creates a new EigenDACert from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACert(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *contractEigenDABlobVerifier.NonSignerStakesAndSignature,
) (*EigenDACert, error) {

	bindingInclusionInfo, err := InclusionInfoProtoToBinding(blobStatusReply.GetBlobInclusionInfo())
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	bindingBatchHeader, err := BatchHeaderProtoToBinding(blobStatusReply.GetSignedBatch().GetHeader())
	if err != nil {
		return nil, fmt.Errorf("convert batch header to binding: %w", err)
	}

	return &EigenDACert{
		BlobInclusionInfo:           *bindingInclusionInfo,
		BatchHeader:                 *bindingBatchHeader,
		NonSignerStakesAndSignature: *nonSignerStakesAndSignature,
	}, nil
}
