package verification

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	contractEigenDACertVerifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
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

// ComputeBlobKey computes the BlobKey of the blob that belongs to the EigenDACert
func (c *EigenDACert) ComputeBlobKey() (*v2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader

	blobCommitmentsProto := BlobCommitmentBindingToProto(&blobHeader.Commitment)
	blobCommitments, err := encoding.BlobCommitmentsFromProtobuf(blobCommitmentsProto)
	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKeyBytes, err := v2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
		blobHeader.Salt)

	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobKey, err := v2.BytesToBlobKey(blobKeyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("bytes to blob key: %w", err)
	}

	return &blobKey, nil
}
