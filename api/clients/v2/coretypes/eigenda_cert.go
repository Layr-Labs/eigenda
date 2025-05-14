package coretypes

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	v2_cert_verifier "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	v3CertTypeEncodeArgs abi.Arguments
)

func init() {
	// load the ABI and parse the dummy interface methods used to encode the cert
	// NOTE: the only other way would be defining the certificate using go-ethereum's abi
	// low level types which would require much boiler plate
	certTypesBinding, err := certTypesBinding.ContractIEigenDACertTypeBindingsMetaData.GetAbi()
	if err != nil {
		panic(err)
	}

	v3CertTypeEncodeMethod, ok := certTypesBinding.Methods["dummyVerifyDACertV3"]
	if !ok {
		panic("dummyVerifyDACertV3 not found in IEigenDACertTypes ABI")
	}

	v3CertTypeEncodeArgs = v3CertTypeEncodeMethod.Inputs

}

// CertificateVersion denotes the version of the EigenDA certificate
// and is interpreted from querying the EigenDACertVerifier contract's
// CertVersion() view function
type CertificateVersion = byte

const (
	// starting at two since we never formally defined a V1 cert in the core code
	VersionTwoCert   = 0x2
	VersionThreeCert = 0x3
)

type EigenDACert interface {
	BlobVersion() v2.BlobVersion
	RelayKeys() []v2.RelayKey
	Version() CertificateVersion
	ReferenceBlockNumber() uint64
	QuorumNumbers() []byte

	ComputeBlobKey() (*v2.BlobKey, error)
	Commitments() (*encoding.BlobCommitments, error)
	Serialize() ([]byte, error)
}

var _ EigenDACert = &EigenDACertV2{}
var _ EigenDACert = &EigenDACertV3{}


// This struct represents the composition of a EigenDA V3 certificate, as it would exist in a rollup inbox.
type EigenDACertV3 certTypesBinding.EigenDACertTypesEigenDACertV3

// BuildEigenDACertV3 creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDACertV3(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature,
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

// RelayKeys returns the relay keys used for reading blob contents from disperser relays
func (c *EigenDACertV3) RelayKeys() []v2.RelayKey {
	return c.BlobInclusionInfo.BlobCertificate.RelayKeys
}

// QuorumNumbers returns the quorum numbers requested
func (c *EigenDACertV3) QuorumNumbers() []byte {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers
}

// RBN returns the reference block number
func (c *EigenDACertV3) ReferenceBlockNumber() uint64 {
	return uint64(c.BatchHeader.ReferenceBlockNumber)
}

// BlobVersion returns the blob version of the blob header
func (c *EigenDACertV3) BlobVersion() v2.BlobVersion {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.Version
}

// ComputeBlobKey computes the blob key used for looking up the blob against an EigenDA network retrieval
// entrypoint (e.g, a relay or a validator node)
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
	return v3CertTypeEncodeArgs.Pack(c)
}

// Commitments returns the blob's cryptographic kzg commitments 
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

// Version returns the version of the EigenDA certificate
func (c *EigenDACertV3) Version() CertificateVersion {
	return VersionThreeCert
}

// This struct represents the composition of an EigenDA V2 certificate
// NOTE: This type is hardforked from the V3 type and will no longer
//       be supported for dispersals after the CertV3 hardfork
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
// RelayKeys returns the relay keys used for reading blob contents from disperser relays
func (c *EigenDACertV2) RelayKeys() []v2.RelayKey {
	return c.BlobInclusionInfo.BlobCertificate.RelayKeys
}

// Commitments returns the blob's cryptographic kzg commitments 
func (c *EigenDACertV2) Commitments() (*encoding.BlobCommitments, error) {
	return BlobCommitmentsBindingToInternal(
		&c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment)
}

// RBN returns the reference block number
func (c *EigenDACertV2) ReferenceBlockNumber() uint64 {
	return uint64(c.BatchHeader.ReferenceBlockNumber)
}

// BlobVersion returns the blob version of the blob header
func (c *EigenDACertV2) BlobVersion() v2.BlobVersion {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.Version
}
// QuorumNumbers returns the quorum numbers requested
func (c *EigenDACertV2) QuorumNumbers() []byte {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers
}

// Serialize serializes the EigenDACertV2 to bytes
func (c *EigenDACertV2) Serialize() ([]byte, error) {
	b, err := rlp.EncodeToBytes(c)
	if err != nil {
		return nil, fmt.Errorf("rlp encode v2 cert: %w", err)
	}

	return b, nil
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

// Version returns the version of the EigenDA certificate
func (c *EigenDACertV2) Version() CertificateVersion {
	return VersionTwoCert
}