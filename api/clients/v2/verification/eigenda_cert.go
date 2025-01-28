package verification

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"

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

// BuildEigenDACert creates a new EigenDACert from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACert(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *contractEigenDACertVerifier.NonSignerStakesAndSignature,
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
