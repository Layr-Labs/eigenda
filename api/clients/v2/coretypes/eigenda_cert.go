package coretypes

import (
	"fmt"

	disperser "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	contractEigenDACertVerifierV2 "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	coreV2 "github.com/Layr-Labs/eigenda/core/v2"
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
type CertificateVersion = uint8

const (
	// starting at two since we never formally defined a V1 cert in the core codebase
	VersionTwoCert   = 0x2
	VersionThreeCert = 0x3
)

// VerificationStatusCode represents the status codes that can be returned by the EigenDACertVerifier.checkDACert contract calls.
// It should match exactly the status codes defined in the contract:
// https://github.com/Layr-Labs/eigenda/blob/1091f460ba762b84019389cbb82d9b04bb2c2bdb/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol#L48-L54
type VerificationStatusCode uint8

const (
	// NULL_ERROR Unused status code. If this is returned, there is a bug in the code.
	StatusNullError VerificationStatusCode = iota
	// SUCCESS Verification succeeded
	StatusSuccess
	// INVALID_INCLUSION_PROOF Merkle inclusion proof is invalid
	StatusInvalidInclusionProof
	// SECURITY_ASSUMPTIONS_NOT_MET Security assumptions not met
	StatusSecurityAssumptionsNotMet
	// BLOB_QUORUMS_NOT_SUBSET Blob quorums not a subset of confirmed quorums
	StatusBlobQuorumsNotSubset
	// REQUIRED_QUORUMS_NOT_SUBSET Required quorums not a subset of blob quorums
	StatusRequiredQuorumsNotSubset
)

// String returns a human-readable representation of the StatusCode.
func (s VerificationStatusCode) String() string {
	switch s {
	case StatusNullError:
		return "Null Error: Unused status code. If this is returned, there is a bug in the code."
	case StatusSuccess:
		return "Success: Verification succeeded"
	case StatusInvalidInclusionProof:
		return "Invalid inclusion proof detected: Merkle inclusion proof for blob batch is invalid"
	case StatusSecurityAssumptionsNotMet:
		return "Security assumptions not met: BLS signer weight is less than the required threshold"
	case StatusBlobQuorumsNotSubset:
		return "Blob quorums are not a subset of the confirmed quorums"
	case StatusRequiredQuorumsNotSubset:
		return "Required quorums are not a subset of the blob quorums"
	default:
		return "Unknown status code"
	}
}

type CertSerializationType byte

const (
	// CertSerializationRLP is the RLP encoding of the certificate
	CertSerializationRLP CertSerializationType = iota
	// CertSerializationABI is the ABI encoding of the certificate
	CertSerializationABI
)

// EigenDACert is a sum type interface returned by the payload disperser
type EigenDACert interface {
	Version() CertificateVersion
}

// RetrievableEigenDACert is an interface that defines data field accessor methods
// used for retrieving the EigenDA certificate from the relay subnet or validator nodes
type RetrievableEigenDACert interface {
	RelayKeys() []coreV2.RelayKey
	QuorumNumbers() []byte
	ReferenceBlockNumber() uint64
	ComputeBlobKey() (*coreV2.BlobKey, error)
	BlobHeader() (*coreV2.BlobHeaderWithHashedPayment, error)
	Commitments() (*encoding.BlobCommitments, error)
	Serialize(ct CertSerializationType) ([]byte, error)
}

var _ EigenDACert = &EigenDACertV2{}
var _ EigenDACert = &EigenDACertV3{}

// This struct represents the composition of a EigenDA V3 certificate, as it would exist in a rollup inbox.
type EigenDACertV3 certTypesBinding.EigenDACertTypesEigenDACertV3

// NewEigenDACertV3 creates a new EigenDACertV3 from a BlobStatusReply, and NonSignerStakesAndSignature
func NewEigenDACertV3(
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
func (c *EigenDACertV3) RelayKeys() []coreV2.RelayKey {
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

// ComputeBlobKey computes the blob key used for looking up the blob against an EigenDA network retrieval
// entrypoint (e.g, a relay or a validator node)
func (c *EigenDACertV3) ComputeBlobKey() (*coreV2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader
	blobCommitments, err := c.Commitments()
	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKeyBytes, err := coreV2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)
	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}
	blobKey, err := coreV2.BytesToBlobKey(blobKeyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("bytes to blob key: %w", err)
	}
	return &blobKey, nil
}

// BlobHeader returns the blob header of the EigenDACertV3
func (c *EigenDACertV3) BlobHeader() (*coreV2.BlobHeaderWithHashedPayment, error) {
	commitments, err := c.Commitments()
	if err != nil {
		return nil, fmt.Errorf("calculate coretype commitments: %w", err)
	}

	blobHeader := &coreV2.BlobHeaderWithHashedPayment{
		BlobVersion:         c.BlobInclusionInfo.BlobCertificate.BlobHeader.Version,
		BlobCommitments:     *commitments,
		QuorumNumbers:       c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
		PaymentMetadataHash: c.BlobInclusionInfo.BlobCertificate.BlobHeader.PaymentHeaderHash,
	}

	return blobHeader, nil
}

func (c *EigenDACertV3) Serialize(ct CertSerializationType) ([]byte, error) {
	switch ct {
	case CertSerializationRLP:
		b, err := rlp.EncodeToBytes(c)
		if err != nil {
			return nil, fmt.Errorf("rlp encode v3 cert: %w", err)
		}
		return b, nil

	case CertSerializationABI:
		b, err := v3CertTypeEncodeArgs.Pack(c)
		if err != nil {
			return nil, fmt.Errorf("abi encode v3 cert: %w", err)
		}
		return b, nil

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)
	}

}

// Commitments returns the blob's cryptographic kzg commitments
func (c *EigenDACertV3) Commitments() (*encoding.BlobCommitments, error) {
	// TODO: figure out how to remove this casting entirely
	commitments := contractEigenDACertVerifierV2.EigenDATypesV2BlobCommitment{
		Commitment: contractEigenDACertVerifierV2.BN254G1Point{
			X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.X,
			Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.Y,
		},
		LengthCommitment: contractEigenDACertVerifierV2.BN254G2Point{
			X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.X,
			Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.Y,
		},
		LengthProof: contractEigenDACertVerifierV2.BN254G2Point{
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
//
//	be supported for dispersals after the CertV3 hardfork
type EigenDACertV2 struct {
	BlobInclusionInfo           contractEigenDACertVerifierV2.EigenDATypesV2BlobInclusionInfo
	BatchHeader                 contractEigenDACertVerifierV2.EigenDATypesV2BatchHeaderV2
	NonSignerStakesAndSignature contractEigenDACertVerifierV2.EigenDATypesV1NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
}

// BuildEigenDAV2Cert creates a new EigenDACertV2 from a BlobStatusReply, and NonSignerStakesAndSignature
func BuildEigenDAV2Cert(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *contractEigenDACertVerifierV2.EigenDATypesV1NonSignerStakesAndSignature,
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
func (c *EigenDACertV2) RelayKeys() []coreV2.RelayKey {
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

// QuorumNumbers returns the quorum numbers requested
func (c *EigenDACertV2) QuorumNumbers() []byte {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers
}

// BlobHeader returns the blob header of the EigenDACertV2
func (c *EigenDACertV2) BlobHeader() (*coreV2.BlobHeaderWithHashedPayment, error) {
	commitments, err := c.Commitments()
	if err != nil {
		return nil, fmt.Errorf("calculate coretype commitments: %w", err)
	}

	blobHeader := &coreV2.BlobHeaderWithHashedPayment{
		BlobVersion:         c.BlobInclusionInfo.BlobCertificate.BlobHeader.Version,
		BlobCommitments:     *commitments,
		QuorumNumbers:       c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
		PaymentMetadataHash: c.BlobInclusionInfo.BlobCertificate.BlobHeader.PaymentHeaderHash,
	}
	return blobHeader, nil
}

// Serialize serializes the EigenDACertV2 to bytes
func (c *EigenDACertV2) Serialize(ct CertSerializationType) ([]byte, error) {
	switch ct {
	case CertSerializationRLP:
		b, err := rlp.EncodeToBytes(c)
		if err != nil {
			return nil, fmt.Errorf("rlp encode v2 cert: %w", err)
		}
		return b, nil

	case CertSerializationABI:
		return nil, fmt.Errorf("abi serialization not supported for v2 cert")

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)

	}
}

// ComputeBlobKey computes the BlobKey of the blob that belongs to the EigenDACertV2
func (c *EigenDACertV2) ComputeBlobKey() (*coreV2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader

	blobCommitments, err := BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKeyBytes, err := coreV2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)

	if err != nil {
		return nil, fmt.Errorf("compute blob key: %w", err)
	}

	blobKey, err := coreV2.BytesToBlobKey(blobKeyBytes[:])
	if err != nil {
		return nil, fmt.Errorf("bytes to blob key: %w", err)
	}

	return &blobKey, nil
}

// Version returns the version of the EigenDA certificate
func (c *EigenDACertV2) Version() CertificateVersion {
	return VersionTwoCert
}

// ToV3 converts an EigenDACertV2 to an EigenDACertV3
func (c *EigenDACertV2) ToV3() (*EigenDACertV3, error) {
	// Convert BlobInclusionInfo from V2 to V3 format
	v3BlobInclusionInfo := certTypesBinding.EigenDATypesV2BlobInclusionInfo{
		BlobCertificate: certTypesBinding.EigenDATypesV2BlobCertificate{
			BlobHeader: certTypesBinding.EigenDATypesV2BlobHeaderV2{
				Version:       c.BlobInclusionInfo.BlobCertificate.BlobHeader.Version,
				QuorumNumbers: c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
				Commitment: certTypesBinding.EigenDATypesV2BlobCommitment{
					Commitment: certTypesBinding.BN254G1Point{
						X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.X,
						Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.Y,
					},
					LengthCommitment: certTypesBinding.BN254G2Point{
						X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.X,
						Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.Y,
					},
					LengthProof: certTypesBinding.BN254G2Point{
						X: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.X,
						Y: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.Y,
					},
					Length: c.BlobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length,
				},
				PaymentHeaderHash: c.BlobInclusionInfo.BlobCertificate.BlobHeader.PaymentHeaderHash,
			},
			Signature: c.BlobInclusionInfo.BlobCertificate.Signature,
			RelayKeys: convertUint32SliceToRelayKeys(c.BlobInclusionInfo.BlobCertificate.RelayKeys),
		},
		BlobIndex:      c.BlobInclusionInfo.BlobIndex,
		InclusionProof: c.BlobInclusionInfo.InclusionProof,
	}

	// Convert BatchHeader from V2 to V3 format
	v3BatchHeader := certTypesBinding.EigenDATypesV2BatchHeaderV2{
		BatchRoot:            c.BatchHeader.BatchRoot,
		ReferenceBlockNumber: c.BatchHeader.ReferenceBlockNumber,
	}

	// Convert NonSignerStakesAndSignature from V2 to V3 format
	v3NonSignerStakesAndSignature := certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: c.NonSignerStakesAndSignature.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             convertV2PubkeysToV3(c.NonSignerStakesAndSignature.NonSignerPubkeys),
		QuorumApks:                   convertV2PubkeysToV3(c.NonSignerStakesAndSignature.QuorumApks),
		ApkG2: certTypesBinding.BN254G2Point{
			X: c.NonSignerStakesAndSignature.ApkG2.X,
			Y: c.NonSignerStakesAndSignature.ApkG2.Y,
		},
		Sigma: certTypesBinding.BN254G1Point{
			X: c.NonSignerStakesAndSignature.Sigma.X,
			Y: c.NonSignerStakesAndSignature.Sigma.Y,
		},
		QuorumApkIndices:      c.NonSignerStakesAndSignature.QuorumApkIndices,
		TotalStakeIndices:     c.NonSignerStakesAndSignature.TotalStakeIndices,
		NonSignerStakeIndices: c.NonSignerStakesAndSignature.NonSignerStakeIndices,
	}

	// Create the V3 certificate
	certV3 := &EigenDACertV3{
		BlobInclusionInfo:           v3BlobInclusionInfo,
		BatchHeader:                 v3BatchHeader,
		NonSignerStakesAndSignature: v3NonSignerStakesAndSignature,
		SignedQuorumNumbers:         c.SignedQuorumNumbers,
	}

	return certV3, nil
}

// convertUint32SliceToRelayKeys converts []uint32 to []coreV2.RelayKey for V3 format
func convertUint32SliceToRelayKeys(relayKeys []uint32) []coreV2.RelayKey {
	result := make([]coreV2.RelayKey, len(relayKeys))
	for i, key := range relayKeys {
		result[i] = coreV2.RelayKey(key)
	}
	return result
}

// convertV2PubkeysToV3 converts V2 pubkeys format to V3 format
func convertV2PubkeysToV3(v2Pubkeys []contractEigenDACertVerifierV2.BN254G1Point) []certTypesBinding.BN254G1Point {
	result := make([]certTypesBinding.BN254G1Point, len(v2Pubkeys))
	for i, pubkey := range v2Pubkeys {
		result[i] = certTypesBinding.BN254G1Point{
			X: pubkey.X,
			Y: pubkey.Y,
		}
	}
	return result
}
