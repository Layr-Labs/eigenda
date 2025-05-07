package coretypes

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	v2_cert_verifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	cert_types_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

type CertificateVersion = byte

const (
	// we never had a proper definition for a version 1 certificate
	// in the eigenda-proxy prefix encoding; version 1 certs are mapped to 0x0
	VersionTwoCert   = 0x1
	VersionThreeCert = 0x2
)

// This struct represents the composition of a EigenDA V3 certificate, as it would exist in a rollup inbox.
type EigenDACertV3 = cert_types_binding.EigenDACertTypesEigenDACertV3

// BuildEigenDACertV3 creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACertV3(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *cert_types_binding.EigenDATypesV1NonSignerStakesAndSignature,
) (*EigenDACertV3, error) {

	bindingInclusionInfo, err := InclusionInfoProtoToBinding(blobStatusReply.GetBlobInclusionInfo())
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	signedBatch := blobStatusReply.GetSignedBatch()

	bindingBatchHeader, err := BatchHeaderProtoToBinding(signedBatch.GetHeader())
	if err != nil {
		return nil, fmt.Errorf("convert batch header to binding: %w", err)
	}

	quorumNumbers, err := QuorumNumbersUint32ToUint8(signedBatch.GetAttestation().GetQuorumNumbers())
	if err != nil {
		return nil, fmt.Errorf("convert quorum numbers to uint8: %w", err)
	}

	return &EigenDACertV3{
		BlobInclusionInfo:           *bindingInclusionInfo,
		BatchHeader:                 *bindingBatchHeader,
		NonSignerStakesAndSignature: *nonSignerStakesAndSignature,
		SignedQuorumNumbers:         quorumNumbers,
	}, nil
}

// This struct represents the composition of an EigenDA V2 certificate, as it would exist in a rollup inbox.
type EigenDACertV2 struct {
	BlobInclusionInfo           v2_cert_verifier.EigenDATypesV2BlobInclusionInfo
	BatchHeader                 v2_cert_verifier.EigenDATypesV2BatchHeaderV2
	NonSignerStakesAndSignature v2_cert_verifier.EigenDATypesV1NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
}

// BuildEigenDACertV2 creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACertV2(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *v2_cert_verifier.EigenDATypesV1NonSignerStakesAndSignature,
) (*EigenDACertV2, error) {

	bindingInclusionInfo, err := InclusionInfoProtoToBinding(blobStatusReply.GetBlobInclusionInfo())
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	signedBatch := blobStatusReply.GetSignedBatch()

	bindingBatchHeader, err := BatchHeaderProtoToBinding(signedBatch.GetHeader())
	if err != nil {
		return nil, fmt.Errorf("convert batch header to binding: %w", err)
	}

	quorumNumbers, err := QuorumNumbersUint32ToUint8(signedBatch.GetAttestation().GetQuorumNumbers())
	if err != nil {
		return nil, fmt.Errorf("convert quorum numbers to uint8: %w", err)
	}

	return &EigenDACertV2{
		BlobInclusionInfo:           *bindingInclusionInfo,
		BatchHeader:                 *bindingBatchHeader,
		NonSignerStakesAndSignature: *nonSignerStakesAndSignature,
		SignedQuorumNumbers:         quorumNumbers,
	}, nil
}

// ComputeBlobKey computes the BlobKey of the blob that belongs to the EigenDACertV2
func (c *EigenDACertV2) ComputeBlobKey() (*v2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader

	blobCommitments, err := BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKeyBytes, err := v2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)

	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobKey, err := v2.BytesToBlobKey(blobKeyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("bytes to blob key: %w", err)
	}

	return &blobKey, nil
}
