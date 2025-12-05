package coretypes

import (
	"encoding/json"
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
	v4CertTypeEncodeArgs abi.Arguments
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

	v4CertTypeEncodeMethod, ok := certTypesBinding.Methods["dummyVerifyDACertV4"]
	if !ok {
		panic("dummyVerifyDACertV4 not found in IEigenDACertTypes ABI")
	}

	v4CertTypeEncodeArgs = v4CertTypeEncodeMethod.Inputs

}

// CertificateVersion denotes the version of the EigenDA certificate
// and is interpreted from querying the EigenDACertVerifier contract's
// CertVersion() view function
type CertificateVersion = uint8

const (
	// starting at two since we never formally defined a V1 cert in the core codebase
	VersionTwoCert   = 0x2
	VersionThreeCert = 0x3
	VersionFourCert  = 0x4
)

type CertSerializationType byte

const (
	// CertSerializationRLP is the RLP encoding of the certificate
	CertSerializationRLP CertSerializationType = iota
	// CertSerializationABI is the ABI encoding of the certificate
	CertSerializationABI
)

// EigenDACert is an interface that defines data field accessor methods
// used for retrieving the EigenDA certificate from the relay subnet or validator nodes
type EigenDACert interface {
	RelayKeys() []coreV2.RelayKey
	QuorumNumbers() []byte
	ReferenceBlockNumber() uint64
	ComputeBlobKey() (coreV2.BlobKey, error)
	BlobHeader() (*coreV2.BlobHeaderWithHashedPayment, error)
	Commitments() (*encoding.BlobCommitments, error)
	Serialize(ct CertSerializationType) ([]byte, error)
	// isEigenDACert is an unexported method that restricts
	// which types can implement this interface to only those
	// defined in this package
	//
	// For the theoretical reasoning behind this choice, see
	// https://www.tedinski.com/2018/02/27/the-expression-problem.html
	isEigenDACert()
}

// DeserializeEigenDACert deserializes raw bytes into an EigenDACert
// based on the provided version and serialization type
func DeserializeEigenDACert(data []byte, version CertificateVersion, ct CertSerializationType) (EigenDACert, error) {
	switch version {
	case VersionTwoCert:
		return DeserializeEigenDACertV2(data, ct)
	case VersionThreeCert:
		return DeserializeEigenDACertV3(data, ct)
	case VersionFourCert:
		return DeserializeEigenDACertV4(data, ct)
	default:
		return nil, fmt.Errorf("unsupported certificate version: %d", version)
	}
}

var _ EigenDACert = &EigenDACertV2{}
var _ EigenDACert = &EigenDACertV3{}
var _ EigenDACert = &EigenDACertV4{}

// This struct represents the composition of a EigenDA V4 certificate, as it would exist in a rollup inbox.
type EigenDACertV4 certTypesBinding.EigenDACertTypesEigenDACertV4

// NewEigenDACertV4 creates a new EigenDACertV4 from a BlobStatusReply, NonSignerStakesAndSignature and
// offchainDerivationVersion. A V4 cert is an extension of a V3 cert with the addition of offchainDerivationVersion.
func NewEigenDACertV4(
	blobStatusReply *disperser.BlobStatusReply,
	nonSignerStakesAndSignature *certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature,
	offchainDerivationVersion uint16,
) (*EigenDACertV4, error) {
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

	return &EigenDACertV4{
		BlobInclusionInfo:           *bindingInclusionInfo,
		BatchHeader:                 *bindingBatchHeader,
		NonSignerStakesAndSignature: *nonSignerStakesAndSignature,
		SignedQuorumNumbers:         quorumNumbers,
		OffchainDerivationVersion:   offchainDerivationVersion,
	}, nil
}

// RelayKeys returns the relay keys used for reading blob contents from disperser relays
func (c *EigenDACertV4) RelayKeys() []coreV2.RelayKey {
	return c.BlobInclusionInfo.BlobCertificate.RelayKeys
}

// QuorumNumbers returns the quorum numbers requested
func (c *EigenDACertV4) QuorumNumbers() []byte {
	return c.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers
}

// RBN returns the reference block number
func (c *EigenDACertV4) ReferenceBlockNumber() uint64 {
	return uint64(c.BatchHeader.ReferenceBlockNumber)
}

// ComputeBlobKey computes the blob key used for looking up the blob against an EigenDA network retrieval
// entrypoint (e.g, a relay or a validator node)
func (c *EigenDACertV4) ComputeBlobKey() (coreV2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader
	blobCommitments, err := c.Commitments()
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKey, err := coreV2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("compute blob key: %w", err)
	}
	return blobKey, nil
}

// BlobHeader returns the blob header of the EigenDACertV4
func (c *EigenDACertV4) BlobHeader() (*coreV2.BlobHeaderWithHashedPayment, error) {
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

func (c *EigenDACertV4) Serialize(ct CertSerializationType) ([]byte, error) {
	switch ct {
	case CertSerializationRLP:
		b, err := rlp.EncodeToBytes(c)
		if err != nil {
			return nil, fmt.Errorf("rlp encode v4 cert: %w", err)
		}
		return b, nil

	case CertSerializationABI:
		b, err := v4CertTypeEncodeArgs.Pack(c)
		if err != nil {
			return nil, fmt.Errorf("abi encode v4 cert: %w", err)
		}
		return b, nil

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)
	}
}

// DeserializeEigenDACertV4 deserializes raw bytes into an EigenDACertV4 provided the serialization
// standard being used
func DeserializeEigenDACertV4(data []byte, ct CertSerializationType) (*EigenDACertV4, error) {
	switch ct {
	case CertSerializationRLP:
		var cert EigenDACertV4
		if err := rlp.DecodeBytes(data, &cert); err != nil {
			return nil, fmt.Errorf("rlp decode v4 cert: %w", err)
		}
		return &cert, nil

	case CertSerializationABI:
		abiMap := make(map[string]interface{})
		err := v4CertTypeEncodeArgs.UnpackIntoMap(abiMap, data)
		if err != nil {
			return nil, fmt.Errorf("unpacking from encoding ABI: %w", err)
		}

		// use json as intermediary to cast abstract type to bytes to
		// then deserialize into structured certificate type
		bytes, err := json.Marshal(abiMap["cert"])
		if err != nil {
			return nil, fmt.Errorf("marshalling ABI arg into bytes: %w", err)
		}

		var cert *EigenDACertV4
		err = json.Unmarshal(bytes, &cert)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal v4 cert: %w", err)
		}

		return cert, nil

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)
	}
}

// Commitments returns the blob's cryptographic kzg commitments
func (c *EigenDACertV4) Commitments() (*encoding.BlobCommitments, error) {
	return commitments(&c.BlobInclusionInfo)
}

// isEigenDACert is an unexported method that restricts which types can implement this interface to only those
// defined in this package
func (c *EigenDACertV4) isEigenDACert() {}

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
func (c *EigenDACertV3) ComputeBlobKey() (coreV2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader
	blobCommitments, err := c.Commitments()
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKey, err := coreV2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("compute blob key: %w", err)
	}
	return blobKey, nil
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

// DeserializeEigenDACertV2 deserializes raw bytes into an EigenDACertV2
func DeserializeEigenDACertV2(data []byte, ct CertSerializationType) (*EigenDACertV3, error) {
	switch ct {
	case CertSerializationRLP:
		var cert EigenDACertV2
		if err := rlp.DecodeBytes(data, &cert); err != nil {
			return nil, fmt.Errorf("rlp decode v2 cert: %w", err)
		}

		return cert.ToV3(), nil

	case CertSerializationABI:
		return nil, fmt.Errorf("abi encoding is not supported for legacy v2 cert")

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)
	}
}

// DeserializeEigenDACertV3 deserializes raw bytes into an EigenDACertV3 provided the serialization
// standard being used
func DeserializeEigenDACertV3(data []byte, ct CertSerializationType) (*EigenDACertV3, error) {
	switch ct {
	case CertSerializationRLP:
		var cert EigenDACertV3
		if err := rlp.DecodeBytes(data, &cert); err != nil {
			return nil, fmt.Errorf("rlp decode v3 cert: %w", err)
		}
		return &cert, nil

	case CertSerializationABI:
		abiMap := make(map[string]interface{})
		err := v3CertTypeEncodeArgs.UnpackIntoMap(abiMap, data)
		if err != nil {
			return nil, fmt.Errorf("unpacking from encoding ABI: %w", err)
		}

		// use json as intermediary to cast abstract type to bytes to
		// then deserialize into structured certificate type
		bytes, err := json.Marshal(abiMap["cert"])
		if err != nil {
			return nil, fmt.Errorf("marshalling ABI arg into bytes: %w", err)
		}

		var cert *EigenDACertV3
		err = json.Unmarshal(bytes, &cert)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal v3 cert: %w", err)
		}

		return cert, nil

	default:
		return nil, fmt.Errorf("unknown serialization type: %d", ct)
	}
}

// Commitments returns the blob's cryptographic kzg commitments
func (c *EigenDACertV3) Commitments() (*encoding.BlobCommitments, error) {
	return commitments(&c.BlobInclusionInfo)
}

// isEigenDACert is an unexported method that restricts which types can implement this interface to only those
// defined in this package
func (c *EigenDACertV3) isEigenDACert() {}

// This struct represents the composition of an EigenDA V2 certificate
// NOTE: This type is hardforked from the V3 type and will no longer
// be supported for dispersals after the CertV3 hardfork
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
func (c *EigenDACertV2) ComputeBlobKey() (coreV2.BlobKey, error) {
	blobHeader := c.BlobInclusionInfo.BlobCertificate.BlobHeader

	blobCommitments, err := BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	blobKey, err := coreV2.ComputeBlobKey(
		blobHeader.Version,
		*blobCommitments,
		blobHeader.QuorumNumbers,
		blobHeader.PaymentHeaderHash,
	)
	if err != nil {
		return coreV2.BlobKey{}, fmt.Errorf("compute blob key: %w", err)
	}
	return blobKey, nil
}

// isEigenDACert is an unexported method that restricts which types can implement this interface to only those
// defined in this package
func (c *EigenDACertV2) isEigenDACert() {}

// ToV3 converts an EigenDACertV2 to an EigenDACertV3
func (c *EigenDACertV2) ToV3() *EigenDACertV3 {
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
	return &EigenDACertV3{
		BlobInclusionInfo:           v3BlobInclusionInfo,
		BatchHeader:                 v3BatchHeader,
		NonSignerStakesAndSignature: v3NonSignerStakesAndSignature,
		SignedQuorumNumbers:         c.SignedQuorumNumbers,
	}
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

func commitments(blobInclusionInfo *certTypesBinding.EigenDATypesV2BlobInclusionInfo) (*encoding.BlobCommitments, error) {
	// TODO: figure out how to remove this casting entirely
	commitments := contractEigenDACertVerifierV2.EigenDATypesV2BlobCommitment{
		Commitment: contractEigenDACertVerifierV2.BN254G1Point{
			X: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.X,
			Y: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Commitment.Y,
		},
		LengthCommitment: contractEigenDACertVerifierV2.BN254G2Point{
			X: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.X,
			Y: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthCommitment.Y,
		},
		LengthProof: contractEigenDACertVerifierV2.BN254G2Point{
			X: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.X,
			Y: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.LengthProof.Y,
		},
		Length: blobInclusionInfo.BlobCertificate.BlobHeader.Commitment.Length,
	}

	blobCommitments, err := BlobCommitmentsBindingToInternal(&commitments)
	if err != nil {
		return nil, fmt.Errorf("blob commitments from protobuf: %w", err)
	}

	return blobCommitments, nil
}
