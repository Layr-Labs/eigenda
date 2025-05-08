package coretypes

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	v2_cert_verifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	v3_cert_verifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV3"
	cert_types_binding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	V3VerifierABI *abi.ABI
)

func init() {
	var err error
	V3VerifierABI, err = v3_cert_verifier.ContractEigenDACertVerifierV3MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
}

type CertificateVersion = byte

const (
	// we never had a proper definition for a version 1 certificate
	// in the eigenda-proxy prefix encoding; version 1 certs are mapped to 0x0
	VersionTwoCert   = 0x1
	VersionThreeCert = 0x2
)

type EigenDACert interface {
	ComputeBlobKey() (*v2.BlobKey, error)
	Serialize() ([]byte, error)
	RelayKeys() []v2.RelayKey
	Commitments() (*encoding.BlobCommitments, error)
}

// This struct represents the composition of a EigenDA V3 certificate, as it would exist in a rollup inbox.
type EigenDACertV3 cert_types_binding.EigenDACertTypesEigenDACertV3

// BuildEigenDACertV3 creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACertV3(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *cert_types_binding.EigenDATypesV1NonSignerStakesAndSignature,
) (*EigenDACertV3, error) {

	bindingInclusionInfo, err := InclusionInfoProtoToIEigenDATypesBinding(blobStatusReply.GetBlobInclusionInfo())
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	signedBatch := blobStatusReply.GetSignedBatch()

	bindingBatchHeader, err := BatchHeaderProtoToIEigenDATypesBinding(signedBatch.GetHeader())
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

func (c *EigenDACertV3) RelayKeys() []v2.RelayKey {
	return c.BlobInclusionInfo.BlobCertificate.RelayKeys
}

func (c *EigenDACertV3) ComputeBlobKey() (*v2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader
	blobCommitments, err := c.Commitments()
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

func (c *EigenDACertV3) Serialize() ([]byte, error) {
	certBytes, err := V3VerifierABI.Methods["dummyFnCertV3"].Inputs.Pack(c)
	if err != nil {
		return nil, fmt.Errorf("encode cert: %w", err)
	}

	return certBytes, nil
}

func (c *EigenDACertV3) Commitments() (*encoding.BlobCommitments, error) {
		// TODO: figure out how to remove this casting entirely
		commitments := v2_cert_verifier.EigenDATypesV2BlobCommitment{
			Commitment: v2_cert_verifier.BN254G1Point{
				X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.X,
				Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.Y,
			},
			LengthCommitment: v2_cert_verifier.BN254G2Point{
				X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.X,
				Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.Y,
			},
			LengthProof: v2_cert_verifier.BN254G2Point{
				X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.X,
				Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.Y,
			},
			Length: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length,
		}
	
		blobCommitments, err := BlobCommitmentsBindingToInternal(&commitments)
		if err != nil {
			return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
		}

		return blobCommitments, nil
}


// This struct represents the composition of an EigenDA V2 certificate
type EigenDACertV2 struct {
	BlobInclusionInfo           v2_cert_verifier.EigenDATypesV2BlobInclusionInfo
	BatchHeader                 v2_cert_verifier.EigenDATypesV2BatchHeaderV2
	NonSignerStakesAndSignature v2_cert_verifier.EigenDATypesV1NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
}

// BuildEigenDAV2Cert creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDAV2Cert(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *v2_cert_verifier.EigenDATypesV1NonSignerStakesAndSignature,
) (*EigenDACertV2, error) {

	bindingInclusionInfo, err := InclusionInfoProtoToV2CertVerifierBinding(blobStatusReply.GetBlobInclusionInfo())
	if err != nil {
		return nil, fmt.Errorf("convert inclusion info to binding: %w", err)
	}

	signedBatch := blobStatusReply.GetSignedBatch()

	bindingBatchHeader, err := BatchHeaderProtoToV2CertVerifierBinding(signedBatch.GetHeader())
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

func (c *EigenDACertV2) RelayKeys() []v2.RelayKey {
	return c.BlobInclusionInfo.BlobCertificate.RelayKeys
}

func (c *EigenDACertV2) Commitments() (*encoding.BlobCommitments, error) {
	return BlobCommitmentsBindingToInternal(
		&c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment)
}

func (c *EigenDACertV2) Serialize() ([]byte, error) {
	panic("not implemented")
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
