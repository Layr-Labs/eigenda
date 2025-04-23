// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifier

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Attestation is an auto generated low-level Go binding around an user-defined struct.
type Attestation struct {
	NonSignerPubkeys []BN254G1Point
	QuorumApks       []BN254G1Point
	Sigma            BN254G1Point
	ApkG2            BN254G2Point
	QuorumNumbers    []uint32
}

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// BatchHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type BatchHeaderV2 struct {
	BatchRoot            [32]byte
	ReferenceBlockNumber uint32
}

// BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type BatchMetadata struct {
	BatchHeader             BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// BlobCertificate is an auto generated low-level Go binding around an user-defined struct.
type BlobCertificate struct {
	BlobHeader BlobHeaderV2
	Signature  []byte
	RelayKeys  []uint32
}

// BlobCommitment is an auto generated low-level Go binding around an user-defined struct.
type BlobCommitment struct {
	Commitment       BN254G1Point
	LengthCommitment BN254G2Point
	LengthProof      BN254G2Point
	Length           uint32
}

// BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []QuorumBlobParam
}

// BlobHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type BlobHeaderV2 struct {
	Version           uint16
	QuorumNumbers     []byte
	Commitment        BlobCommitment
	PaymentHeaderHash [32]byte
}

// BlobInclusionInfo is an auto generated low-level Go binding around an user-defined struct.
type BlobInclusionInfo struct {
	BlobCertificate BlobCertificate
	BlobIndex       uint32
	InclusionProof  []byte
}

// BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// NonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type NonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// SignedBatch is an auto generated low-level Go binding around an user-defined struct.
type SignedBatch struct {
	BatchHeader BatchHeaderV2
	Attestation Attestation
}

// VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// ContractEigenDACertVerifierMetaData contains all meta data concerning the ContractEigenDACertVerifier contract.
var ContractEigenDACertVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetrieverV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequiredV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6101606040523480156200001257600080fd5b5060405162004de738038062004de783398101604081905262000035916200027f565b6001600160a01b03888116608081905288821660a05260c05286811660e052858116610100908152858216610120529084166101405282516000805460208087015160ff94851661ffff19909316929092179390911690930291909117905581518991889188918891889188918891620000b591600191840190620000cb565b50505050505050505050505050505050620003e1565b828054620000d990620003a4565b90600052602060002090601f016020900481019282620000fd576000855562000148565b82601f106200011857805160ff191683800117855562000148565b8280016001018555821562000148579182015b82811115620001485782518255916020019190600101906200012b565b50620001569291506200015a565b5090565b5b808211156200015657600081556001016200015b565b6001600160a01b03811681146200018757600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620001cb57620001cb6200018a565b604052919050565b805160ff81168114620001e557600080fd5b919050565b600082601f830112620001fc57600080fd5b81516001600160401b038111156200021857620002186200018a565b60206200022e601f8301601f19168201620001a0565b82815285828487010111156200024357600080fd5b60005b838110156200026357858101830151828201840152820162000246565b83811115620002755760008385840101525b5095945050505050565b600080600080600080600080888a036101208112156200029e57600080fd5b8951620002ab8162000171565b60208b0151909950620002be8162000171565b60408b0151909850620002d18162000171565b60608b0151909750620002e48162000171565b60808b0151909650620002f78162000171565b60a08b01519095506200030a8162000171565b9350604060bf19820112156200031f57600080fd5b50604080519081016001600160401b0380821183831017156200034657620003466200018a565b816040526200035860c08d01620001d3565b83526200036860e08d01620001d3565b60208401526101008c0151929450808311156200038457600080fd5b5050620003948b828c01620001ea565b9150509295985092959890939650565b600181811c90821680620003b957607f821691505b60208210811415620003db57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c05160e051610100516101205161014051614911620004d6600039600081816102e901526107a80152600081816102c201526107870152600081816103360152818161069401528181610766015261088d0152600081816101b90152818161067301528181610745015261086c0152600081816101f801528181610406015281816104bb01528181610567015281816105f90152818161065201528181610724015281816107f70152818161084b015281816108e701528181610a0201528181610a770152610af60152600081816103720152818161061a01526108180152600061029b01526149116000f3fe608060405234801561001057600080fd5b506004361061014d5760003560e01c80635fafa482116100c3578063b74d78711161007c578063b74d787114610394578063bafa91071461039c578063ccb7cd0d146103a4578063e15234ff146103b7578063ed0450ae146103bf578063ee6c3bcf146103ef57600080fd5b80635fafa482146102e45780637d644cad1461030b578063813c2eb01461031e57806382c216e7146103315780638687feae14610358578063a9c823e11461036d57600080fd5b80632ecfe72b116101155780632ecfe72b1461021a57806331a3479a1461025d578063415ef61414610270578063421c0222146102835780634cff90c4146102965780635df1f618146102bd57600080fd5b8063048886d2146101525780631429c7c21461017a578063143eb4d91461019f578063154b9e86146101b457806317f3578e146101f3575b600080fd5b610165610160366004612e57565b610402565b60405190151581526020015b60405180910390f35b61018d610188366004612e57565b610498565b60405160ff9091168152602001610171565b6101b26101ad366004612fd6565b610527565b005b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610171565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b61022d61022836600461305e565b610548565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff1690820152606001610171565b6101b261026b3660046130bd565b6105f4565b61016561027e3660046131bd565b61064a565b6101b2610291366004613266565b61071f565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b6101b26103193660046132c9565b6107f2565b6101b261032c3660046131bd565b610846565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b6103606108e3565b6040516101719190613393565b6101db7f000000000000000000000000000000000000000000000000000000000000000081565b610360610970565b6103606109fe565b6101b26103b23660046133a6565b610a5e565b610360610a73565b6000546103d59060ff8082169161010090041682565b6040805160ff938416815292909116602083015201610171565b61018d6103fd366004612e57565b610ad3565b60007f0000000000000000000000000000000000000000000000000000000000000000604051630244436960e11b815260ff841660048201526001600160a01b03919091169063048886d290602401602060405180830381865afa15801561046e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061049291906133d1565b92915050565b604051630a14e3e160e11b815260ff821660048201526000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690631429c7c2906024015b602060405180830381865afa158015610503573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061049291906133f3565b6000806105348484610b25565b915091506105428282610c5b565b50505050565b60408051606081018252600080825260208201819052918101919091527f0000000000000000000000000000000000000000000000000000000000000000604051632ecfe72b60e01b815261ffff841660048201526001600160a01b039190911690632ecfe72b90602401606060405180830381865afa1580156105d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104929190613410565b6105427f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000086868686610645610a73565b610eea565b6000806106ea7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006106c2368b90038b018b6134af565b6106cb8a613612565b6106d48a6138a0565b6106dc6119b3565b6106e46119e7565b8b611a79565b5090506000816006811115610701576107016139c8565b1415610711576001915050610717565b60009150505b949350505050565b6107ee7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006107d0886139de565b6107d988613612565b6107e16119b3565b6107e96119e7565b611c19565b5050565b6107ee7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610841610a73565b611c4d565b6105427f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006108bb368990038901896134af565b6108c488613612565b6108cd886138a0565b6108d56119b3565b6108dd6119e7565b89612352565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610943573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261096b9190810190613ad5565b905090565b6001805461097d90613b4b565b80601f01602080910402602001604051908101604052809291908181526020018280546109a990613b4b565b80156109f65780601f106109cb576101008083540402835291602001916109f6565b820191906000526020600020905b8154815290600101906020018083116109d957829003601f168201915b505050505081565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610943573d6000803e3d6000fd5b600080610534610a6d85610548565b84610b25565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610943573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ee6c3bcf906024016104e6565b60006060826020015160ff16836000015160ff1611610b8257600283600001518460200151604051602001610b6a92919060ff92831681529116602082015260400190565b60405160208183030381529060405291509150610c54565b60208301518351600091610b9591613b96565b60ff1690506000856020015163ffffffff16866040015160ff1683620f4240610bbe9190613bcf565b610bc89190613bcf565b610bd490612710613be3565b610bde9190613bfa565b8651909150600090610bf290612710613c19565b63ffffffff169050808210610c1f5760006040518060200160405280600081525094509450505050610c54565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050505b9250929050565b6000826006811115610c6f57610c6f6139c8565b1415610c79575050565b6001826006811115610c8d57610c8d6139c8565b1415610ce357600080600083806020019051810190610cac9190613c45565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b6002826006811115610cf757610cf76139c8565b1415610d3f5760008082806020019051810190610d149190613c73565b604051631b00235d60e01b815260ff8084166004830152821660248201529193509150604401610cda565b6003826006811115610d5357610d536139c8565b1415610da757600080600083806020019051810190610d729190613c45565b6040516001626dc9ad60e11b031981526004810184905260248101839052604481018290529295509093509150606401610cda565b6004826006811115610dbb57610dbb6139c8565b1415610e005760008082806020019051810190610dd89190613ca2565b604051634a47030360e11b815260048101839052602481018290529193509150604401610cda565b6005826006811115610e1457610e146139c8565b1415610e595760008082806020019051810190610e319190613ca2565b60405163114b085b60e21b815260048101839052602481018290529193509150604401610cda565b6006826006811115610e6d57610e6d6139c8565b1415610ead57600081806020019051810190610e899190613cc6565b6040516309efaa0b60e41b815263ffffffff82166004820152909150602401610cda565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b6044820152606401610cda565b838214610f895760405162461bcd60e51b815260206004820152606d60248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a20626c6f62486561646560648201527f727320616e6420626c6f62566572696669636174696f6e50726f6f6673206c6560848201526c0dccee8d040dad2e6dac2e8c6d609b1b60a482015260c401610cda565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610fc9573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610ff19190810190613ad5565b905060005b858110156119a857876001600160a01b031663eccbbfc986868481811061101f5761101f613ce3565b90506020028101906110319190613cf9565b61103f906020810190613d19565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611081573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110a59190613d36565b6110e88686848181106110ba576110ba613ce3565b90506020028101906110cc9190613cf9565b6110da906040810190613d4f565b6110e390613d65565b612374565b1461117b5760405162461bcd60e51b815260206004820152606360248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a2062617463684d65746160648201527f6461746120646f6573206e6f74206d617463682073746f726564206d6574616460848201526261746160e81b60a482015260c401610cda565b6112c185858381811061119057611190613ce3565b90506020028101906111a29190613cf9565b6111b0906060810190613e3c565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508992508891508590508181106111f9576111f9613ce3565b905060200281019061120b9190613cf9565b611219906040810190613d4f565b6112239080613e82565b356112598a8a8681811061123957611239613ce3565b905060200281019061124b9190613e82565b61125490613e98565b6123e5565b60405160200161126b91815260200190565b6040516020818303038152906040528051906020012088888681811061129357611293613ce3565b90506020028101906112a59190613cf9565b6112b6906040810190602001613d19565b63ffffffff16612415565b61133b5760405162461bcd60e51b815260206004820152605160248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a20696e636c7573696f6e606482015270081c1c9bdbd9881a5cc81a5b9d985b1a59607a1b608482015260a401610cda565b6000805b88888481811061135157611351613ce3565b90506020028101906113639190613e82565b611371906060810190613fb7565b90508110156118e35788888481811061138c5761138c613ce3565b905060200281019061139e9190613e82565b6113ac906060810190613fb7565b828181106113bc576113bc613ce3565b6113d29260206080909202019081019150612e57565b60ff168787858181106113e7576113e7613ce3565b90506020028101906113f99190613cf9565b611407906040810190613d4f565b6114119080613e82565b61141f906020810190613e3c565b89898781811061143157611431613ce3565b90506020028101906114439190613cf9565b611451906080810190613e3c565b8581811061146157611461613ce3565b919091013560f81c905081811061147a5761147a613ce3565b9050013560f81c60f81b60f81c60ff16146115065760405162461bcd60e51b815260206004820152605260248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a2071756f72756d4e756d6064820152710c4cae440c8decae640dcdee840dac2e8c6d60731b608482015260a401610cda565b88888481811061151857611518613ce3565b905060200281019061152a9190613e82565b611538906060810190613fb7565b8281811061154857611548613ce3565b90506080020160200160208101906115609190612e57565b60ff1689898581811061157557611575613ce3565b90506020028101906115879190613e82565b611595906060810190613fb7565b838181106115a5576115a5613ce3565b90506080020160400160208101906115bd9190612e57565b60ff16116116475760405162461bcd60e51b815260206004820152605a60248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a207468726573686f6c6460648201527f2070657263656e746167657320617265206e6f742076616c6964000000000000608482015260a401610cda565b8389898581811061165a5761165a613ce3565b905060200281019061166c9190613e82565b61167a906060810190613fb7565b8381811061168a5761168a613ce3565b6116a09260206080909202019081019150612e57565b60ff16815181106116b3576116b3613ce3565b016020015160f81c8989858181106116cd576116cd613ce3565b90506020028101906116df9190613e82565b6116ed906060810190613fb7565b838181106116fd576116fd613ce3565b90506080020160400160208101906117159190612e57565b60ff1610156117365760405162461bcd60e51b8152600401610cda90614000565b88888481811061174857611748613ce3565b905060200281019061175a9190613e82565b611768906060810190613fb7565b8281811061177857611778613ce3565b90506080020160400160208101906117909190612e57565b60ff168787858181106117a5576117a5613ce3565b90506020028101906117b79190613cf9565b6117c5906040810190613d4f565b6117cf9080613e82565b6117dd906040810190613e3c565b8989878181106117ef576117ef613ce3565b90506020028101906118019190613cf9565b61180f906080810190613e3c565b8581811061181f5761181f613ce3565b919091013560f81c905081811061183857611838613ce3565b9050013560f81c60f81b60f81c60ff1610156118665760405162461bcd60e51b8152600401610cda90614000565b6118cf828a8a8681811061187c5761187c613ce3565b905060200281019061188e9190613e82565b61189c906060810190613fb7565b848181106118ac576118ac613ce3565b6118c29260206080909202019081019150612e57565b600160ff919091161b1790565b9150806118db8161407b565b91505061133f565b506118f76118f08561242d565b8281161490565b6119975760405162461bcd60e51b815260206004820152607160248201526000805160206148bc83398151915260448201527f7269667944414365727473466f7251756f72756d733a2072657175697265642060648201527f71756f72756d7320617265206e6f74206120737562736574206f662074686520608482015270636f6e6669726d65642071756f72756d7360781b60a482015260c401610cda565b506119a18161407b565b9050610ff6565b505050505050505050565b6040805180820182526000808252602091820181905282518084019093525460ff8082168452610100909104169082015290565b6060600180546119f690613b4b565b80601f0160208091040260200160405190810160405280929190818152602001828054611a2290613b4b565b8015611a6f5780601f10611a4457610100808354040283529160200191611a6f565b820191906000526020600020905b815481529060010190602001808311611a5257829003601f168201915b5050505050905090565b60006060611a8788886125ba565b90925090506000826006811115611aa057611aa06139c8565b14611aaa57611c0b565b611abc89886000015160400151612690565b90925090506000826006811115611ad557611ad56139c8565b14611adf57611c0b565b86515151604051632ecfe72b60e01b815261ffff9091166004820152611b5a906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa158015611b30573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b549190613410565b86610b25565b90925090506000826006811115611b7357611b736139c8565b14611b7d57611c0b565b6000611b998b611b8c8b6127c8565b868c602001518b8b6127f3565b919450925090506000836006811115611bb457611bb46139c8565b14611bbf5750611c0b565b87515160200151600090611bd39083612954565b919550935090506000846006811115611bee57611bee6139c8565b14611bfa575050611c0b565b611c0486826129ba565b9350935050505b995099975050505050505050565b600080611c27888888612a0a565b91509150611c408b8b8b896000015189878a8a89612352565b5050505050505050505050565b6001600160a01b03841663eccbbfc9611c696020850185613d19565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611cab573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ccf9190613d36565b611cdf6110da6040850185613d4f565b14611d715760405162461bcd60e51b815260206004820152606260248201526000805160206148bc83398151915260448201527f72696679444143657274466f7251756f72756d733a2062617463684d6574616460648201527f61746120646f6573206e6f74206d617463682073746f726564206d6574616461608482015261746160f01b60a482015260c401610cda565b611e15611d816060840184613e3c565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611dc3925050506040850185613d4f565b611dcd9080613e82565b35611dda61125487613e98565b604051602001611dec91815260200190565b604051602081830303815290604052805190602001208560200160208101906112b69190613d19565b611e8e5760405162461bcd60e51b815260206004820152605060248201526000805160206148bc83398151915260448201527f72696679444143657274466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a401610cda565b6000805b611e9f6060860186613fb7565b905081101561229e57611eb56060860186613fb7565b82818110611ec557611ec5613ce3565b611edb9260206080909202019081019150612e57565b60ff16611eeb6040860186613d4f565b611ef59080613e82565b611f03906020810190613e3c565b611f106080880188613e3c565b85818110611f2057611f20613ce3565b919091013560f81c9050818110611f3957611f39613ce3565b9050013560f81c60f81b60f81c60ff1614611fc45760405162461bcd60e51b815260206004820152605160248201526000805160206148bc83398151915260448201527f72696679444143657274466f7251756f72756d733a2071756f72756d4e756d626064820152700cae440c8decae640dcdee840dac2e8c6d607b1b608482015260a401610cda565b611fd16060860186613fb7565b82818110611fe157611fe1613ce3565b9050608002016020016020810190611ff99190612e57565b60ff166120096060870187613fb7565b8381811061201957612019613ce3565b90506080020160400160208101906120319190612e57565b60ff16116120bb5760405162461bcd60e51b815260206004820152605960248201526000805160206148bc83398151915260448201527f72696679444143657274466f7251756f72756d733a207468726573686f6c642060648201527f70657263656e746167657320617265206e6f742076616c696400000000000000608482015260a401610cda565b6001600160a01b038716631429c7c26120d76060880188613fb7565b848181106120e7576120e7613ce3565b6120fd9260206080909202019081019150612e57565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa15801561213b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061215f91906133f3565b60ff1661216f6060870187613fb7565b8381811061217f5761217f613ce3565b90506080020160400160208101906121979190612e57565b60ff1610156121b85760405162461bcd60e51b8152600401610cda90614096565b6121c56060860186613fb7565b828181106121d5576121d5613ce3565b90506080020160400160208101906121ed9190612e57565b60ff166121fd6040860186613d4f565b6122079080613e82565b612215906040810190613e3c565b6122226080880188613e3c565b8581811061223257612232613ce3565b919091013560f81c905081811061224b5761224b613ce3565b9050013560f81c60f81b60f81c60ff1610156122795760405162461bcd60e51b8152600401610cda90614096565b61228a8261189c6060880188613fb7565b9150806122968161407b565b915050611e92565b506122ab6118f08361242d565b61234a5760405162461bcd60e51b815260206004820152607060248201526000805160206148bc83398151915260448201527f72696679444143657274466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c401610cda565b505050505050565b6000806123668b8b8b8b8b8b8b8b8b611a79565b91509150611c408282610c5b565b6000610492826000015160405160200161238e9190614108565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b6000816040516020016123f89190614168565b604051602081830303815290604052805190602001209050919050565b600083612423868585612c1f565b1495945050505050565b6000610100825111156124b65760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610cda565b81516124c457506000919050565b600080836000815181106124da576124da613ce3565b0160200151600160f89190911c81901b92505b84518110156125b15784818151811061250857612508613ce3565b0160200151600160f89190911c1b915082821161259d5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610cda565b918117916125aa8161407b565b90506124ed565b50909392505050565b6000606060006125cd8460000151612d22565b90506000816040516020016125e491815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000612621876040015183858a6020015163ffffffff16612415565b90508015612648576000604051806020016040528060008152509550955050505050610c54565b6020808801516040805163ffffffff90921692820192909252908101849052606081018390526001906080016040516020818303038152906040529550955050505050610c54565b6000606060005b83518110156127ad5760006001600160a01b0316856001600160a01b031663b5a872da8684815181106126cc576126cc613ce3565b60200260200101516040518263ffffffff1660e01b81526004016126fc919063ffffffff91909116815260200190565b602060405180830381865afa158015612719573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061273d919061420d565b6001600160a01b0316141561279d57600684828151811061276057612760613ce3565b6020026020010151604051602001612784919063ffffffff91909116815260200190565b6040516020818303038152906040529250925050610c54565b6127a68161407b565b9050612697565b50506040805160208101909152600080825291509250929050565b6000816040516020016123f891908151815260209182015163ffffffff169181019190915260400190565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b815260040161282c9493929190614355565b600060405180830381865afa158015612849573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261287191908101906144c2565b5090506000915060005b885181101561293257856000015160ff16826020015182815181106128a2576128a2613ce3565b60200260200101516128b4919061455e565b6001600160601b03166064836000015183815181106128d5576128d5613ce3565b60200260200101516001600160601b03166128f09190613bfa565b106129205761291d838a838151811061290b5761290b613ce3565b0160200151600160f89190911c1b1790565b92505b8061292a8161407b565b91505061287b565b5050604080516020810190915260008082529350915096509650969350505050565b6000606060006129638561242d565b90508381168114156129885760408051602081019091526000808252935091506129b3565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b6000606060006129c98561242d565b90508381168114156129ee575050604080516020810190915260008082529150610c54565b6040805160208101839052908101859052600590606001612784565b612a12612d9c565b60606000836020015160000151516001600160401b03811115612a3757612a37612e7b565b604051908082528060200260200182016040528015612a60578160200160208202803683370190505b50905060005b60208501515151811015612add57612ab08560200151600001518281518110612a9157612a91613ce3565b6020026020010151805160009081526020918201519091526040902090565b828281518110612ac257612ac2613ce3565b6020908102919091010152612ad68161407b565b9050612a66565b5060005b84602001516080015151811015612b4857828560200151608001518281518110612b0d57612b0d613ce3565b6020026020010151604051602001612b26929190614584565b604051602081830303815290604052925080612b419061407b565b9050612ae1565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491612b83918a9190889088906004016145b6565b600060405180830381865afa158015612ba0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052612bc89190810190614709565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b600060208451612c2f91906147e1565b15612cb65760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610cda565b8260205b85518111612d1957612ccd6002856147e1565b612cee57816000528086015160205260406000209150600284049350612d07565b8086015160005281602052604060002091506002840493505b612d126020826147f5565b9050612cba565b50949350505050565b6000612d318260000151612d4a565b60208084015160408086015190516123f894930161480d565b6000816000015182602001518360400151604051602001612d6d93929190614842565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152016123f8565b604051806101000160405280606081526020016060815260200160608152602001612dc5612e02565b8152602001612de7604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060400160405280612e15612e27565b8152602001612e22612e27565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff81168114612e5457600080fd5b50565b600060208284031215612e6957600080fd5b8135612e7481612e45565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715612eb357612eb3612e7b565b60405290565b604051606081016001600160401b0381118282101715612eb357612eb3612e7b565b604051608081016001600160401b0381118282101715612eb357612eb3612e7b565b60405161010081016001600160401b0381118282101715612eb357612eb3612e7b565b60405160a081016001600160401b0381118282101715612eb357612eb3612e7b565b604051601f8201601f191681016001600160401b0381118282101715612f6a57612f6a612e7b565b604052919050565b63ffffffff81168114612e5457600080fd5b8035612f8f81612f72565b919050565b600060408284031215612fa657600080fd5b612fae612e91565b90508135612fbb81612e45565b81526020820135612fcb81612e45565b602082015292915050565b60008082840360a0811215612fea57600080fd5b6060811215612ff857600080fd5b50613001612eb9565b833561300c81612f72565b8152602084013561301c81612f72565b6020820152604084013561302f81612e45565b604082015291506130438460608501612f94565b90509250929050565b803561ffff81168114612f8f57600080fd5b60006020828403121561307057600080fd5b612e748261304c565b60008083601f84011261308b57600080fd5b5081356001600160401b038111156130a257600080fd5b6020830191508360208260051b8501011115610c5457600080fd5b600080600080604085870312156130d357600080fd5b84356001600160401b03808211156130ea57600080fd5b6130f688838901613079565b9096509450602087013591508082111561310f57600080fd5b5061311c87828801613079565b95989497509550505050565b60006060828403121561313a57600080fd5b50919050565b60006001600160401b0382111561315957613159612e7b565b50601f01601f191660200190565b600082601f83011261317857600080fd5b813561318b61318682613140565b612f42565b8181528460208386010111156131a057600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a08112156131d457600080fd5b60408112156131e257600080fd5b5084935060408501356001600160401b038082111561320057600080fd5b61320c88838901613128565b9450606087013591508082111561322257600080fd5b90860190610180828903121561323757600080fd5b9092506080860135908082111561324d57600080fd5b5061325a87828801613167565b91505092959194509250565b6000806040838503121561327957600080fd5b82356001600160401b038082111561329057600080fd5b61329c86838701613128565b935060208501359150808211156132b257600080fd5b506132bf85828601613128565b9150509250929050565b600080604083850312156132dc57600080fd5b82356001600160401b03808211156132f357600080fd5b908401906080828703121561330757600080fd5b9092506020840135908082111561331d57600080fd5b50830160a0818603121561333057600080fd5b809150509250929050565b60005b8381101561335657818101518382015260200161333e565b838111156105425750506000910152565b6000815180845261337f81602086016020860161333b565b601f01601f19169290920160200192915050565b602081526000612e746020830184613367565b600080606083850312156133b957600080fd5b6133c28361304c565b91506130438460208501612f94565b6000602082840312156133e357600080fd5b81518015158114612e7457600080fd5b60006020828403121561340557600080fd5b8151612e7481612e45565b60006060828403121561342257600080fd5b604051606081018181106001600160401b038211171561344457613444612e7b565b604052825161345281612f72565b8152602083015161346281612f72565b6020820152604083015161347581612e45565b60408201529392505050565b60006040828403121561349357600080fd5b61349b612e91565b9050813581526020820135612fcb81612f72565b6000604082840312156134c157600080fd5b612e748383613481565b6000604082840312156134dd57600080fd5b6134e5612e91565b9050813581526020820135602082015292915050565b600082601f83011261350c57600080fd5b613514612e91565b80604084018581111561352657600080fd5b845b81811015613540578035845260209384019301613528565b509095945050505050565b60006080828403121561355d57600080fd5b613565612e91565b905061357183836134fb565b8152612fcb83604084016134fb565b60006001600160401b0382111561359957613599612e7b565b5060051b60200190565b600082601f8301126135b457600080fd5b813560206135c461318683613580565b82815260059290921b840181019181810190868411156135e357600080fd5b8286015b848110156136075780356135fa81612f72565b83529183019183016135e7565b509695505050505050565b60006060823603121561362457600080fd5b61362c612eb9565b82356001600160401b038082111561364357600080fd5b81850191506060823603121561365857600080fd5b613660612eb9565b82358281111561366f57600080fd5b8301368190036101c081121561368457600080fd5b61368c612edb565b6136958361304c565b8152602080840135868111156136aa57600080fd5b6136b636828701613167565b83830152506040610160603f19850112156136d057600080fd5b6136d8612edb565b93506136e6368287016134cb565b84526136f5366080870161354b565b8285015261370736610100870161354b565b8185015261018085013561371a81612f72565b8060608601525083818401526101a085013560608401528286528188013594508685111561374757600080fd5b61375336868a01613167565b828701528088013594508685111561376a57600080fd5b61377636868a016135a3565b81870152858952613788828c01612f84565b828a0152808b013597508688111561379f57600080fd5b6137ab36898d01613167565b90890152509598975050505050505050565b600082601f8301126137ce57600080fd5b813560206137de61318683613580565b82815260069290921b840181019181810190868411156137fd57600080fd5b8286015b848110156136075761381388826134cb565b835291830191604001613801565b600082601f83011261383257600080fd5b8135602061384261318683613580565b82815260059290921b8401810191818101908684111561386157600080fd5b8286015b848110156136075780356001600160401b038111156138845760008081fd5b6138928986838b01016135a3565b845250918301918301613865565b600061018082360312156138b357600080fd5b6138bb612efd565b82356001600160401b03808211156138d257600080fd5b6138de368387016135a3565b835260208501359150808211156138f457600080fd5b613900368387016137bd565b6020840152604085013591508082111561391957600080fd5b613925368387016137bd565b6040840152613937366060870161354b565b60608401526139493660e087016134cb565b608084015261012085013591508082111561396357600080fd5b61396f368387016135a3565b60a084015261014085013591508082111561398957600080fd5b613995368387016135a3565b60c08401526101608501359150808211156139af57600080fd5b506139bc36828601613821565b60e08301525092915050565b634e487b7160e01b600052602160045260246000fd5b6000606082360312156139f057600080fd5b6139f8612e91565b613a023684613481565b815260408301356001600160401b0380821115613a1e57600080fd5b81850191506101208236031215613a3457600080fd5b613a3c612f20565b823582811115613a4b57600080fd5b613a57368286016137bd565b825250602083013582811115613a6c57600080fd5b613a78368286016137bd565b602083015250613a8b36604085016134cb565b6040820152613a9d366080850161354b565b606082015261010083013582811115613ab557600080fd5b613ac1368286016135a3565b608083015250602084015250909392505050565b600060208284031215613ae757600080fd5b81516001600160401b03811115613afd57600080fd5b8201601f81018413613b0e57600080fd5b8051613b1c61318682613140565b818152856020838501011115613b3157600080fd5b613b4282602083016020860161333b565b95945050505050565b600181811c90821680613b5f57607f821691505b6020821081141561313a57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff841680821015613bb057613bb0613b80565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b600082613bde57613bde613bb9565b500490565b600082821015613bf557613bf5613b80565b500390565b6000816000190483118215151615613c1457613c14613b80565b500290565b600063ffffffff80831681851681830481118215151615613c3c57613c3c613b80565b02949350505050565b600080600060608486031215613c5a57600080fd5b8351925060208401519150604084015190509250925092565b60008060408385031215613c8657600080fd5b8251613c9181612e45565b602084015190925061333081612e45565b60008060408385031215613cb557600080fd5b505080516020909101519092909150565b600060208284031215613cd857600080fd5b8151612e7481612f72565b634e487b7160e01b600052603260045260246000fd5b60008235609e19833603018112613d0f57600080fd5b9190910192915050565b600060208284031215613d2b57600080fd5b8135612e7481612f72565b600060208284031215613d4857600080fd5b5051919050565b60008235605e19833603018112613d0f57600080fd5b600060608236031215613d7757600080fd5b613d7f612eb9565b82356001600160401b0380821115613d9657600080fd5b818501915060808236031215613dab57600080fd5b613db3612edb565b82358152602083013582811115613dc957600080fd5b613dd536828601613167565b602083015250604083013582811115613ded57600080fd5b613df936828601613167565b60408301525060608301359250613e0f83612f72565b82606082015280845250505060208301356020820152613e3160408401612f84565b604082015292915050565b6000808335601e19843603018112613e5357600080fd5b8301803591506001600160401b03821115613e6d57600080fd5b602001915036819003821315610c5457600080fd5b60008235607e19833603018112613d0f57600080fd5b60006080808336031215613eab57600080fd5b613eb3612eb9565b613ebd36856134cb565b8152604080850135613ece81612f72565b6020818185015260609150818701356001600160401b03811115613ef157600080fd5b870136601f820112613f0257600080fd5b8035613f1061318682613580565b81815260079190911b82018301908381019036831115613f2f57600080fd5b928401925b82841015613fa357888436031215613f4c5760008081fd5b613f54612edb565b8435613f5f81612e45565b815284860135613f6e81612e45565b8187015284880135613f7f81612e45565b8189015284870135613f9081612f72565b8188015282529288019290840190613f34565b958701959095525093979650505050505050565b6000808335601e19843603018112613fce57600080fd5b8301803591506001600160401b03821115613fe857600080fd5b6020019150600781901b3603821315610c5457600080fd5b60208082526061908201526000805160206148bc83398151915260408201527f7269667944414365727473466f7251756f72756d733a20636f6e6669726d617460608201527f696f6e5468726573686f6c6450657263656e74616765206973206e6f74206d656080820152601d60fa1b60a082015260c00190565b600060001982141561408f5761408f613b80565b5060010190565b602080825260609082018190526000805160206148bc83398151915260408301527f72696679444143657274466f7251756f72756d733a20636f6e6669726d617469908201527f6f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6574608082015260a00190565b6020815281516020820152600060208301516080604084015261412e60a0840182613367565b90506040840151601f1984830301606085015261414b8282613367565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b808410156141fe578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906141bb565b509a9950505050505050505050565b60006020828403121561421f57600080fd5b81516001600160a01b0381168114612e7457600080fd5b600081518084526020808501945080840160005b8381101561426c57815163ffffffff168752958201959082019060010161424a565b509495945050505050565b600081518084526020808501945080840160005b8381101561426c576142a887835180518252602090810151910152565b604096909601959082019060010161428b565b8060005b60028110156105425781518452602093840193909101906001016142bf565b6142e98282516142bb565b60208101516142fb60408401826142bb565b505050565b600081518084526020808501808196508360051b8101915082860160005b85811015614348578284038952614336848351614236565b9885019893509084019060010161431e565b5091979650505050505050565b84815260806020820152600061436e6080830186613367565b63ffffffff851660408401528281036060840152610180845181835261439682840182614236565b915050602085015182820360208401526143b08282614277565b915050604085015182820360408401526143ca8282614277565b91505060608501516143df60608401826142de565b506080850151805160e08401526020015161010083015260a085015182820361012084015261440e8282614236565b91505060c08501518282036101408401526144298282614236565b91505060e08501518282036101608401526144448282614300565b9998505050505050505050565b600082601f83011261446257600080fd5b8151602061447261318683613580565b82815260059290921b8401810191818101908684111561449157600080fd5b8286015b848110156136075780516001600160601b03811681146144b55760008081fd5b8352918301918301614495565b600080604083850312156144d557600080fd5b82516001600160401b03808211156144ec57600080fd5b908401906040828703121561450057600080fd5b614508612e91565b82518281111561451757600080fd5b61452388828601614451565b82525060208301518281111561453857600080fd5b61454488828601614451565b602083015250809450505050602083015190509250929050565b60006001600160601b0380831681851681830481118215151615613c3c57613c3c613b80565b6000835161459681846020880161333b565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526145e46080840186613367565b838103606085015284518082528286019183019060005b81811015614617578351835292840192918401916001016145fb565b50909998505050505050505050565b600082601f83011261463757600080fd5b8151602061464761318683613580565b82815260059290921b8401810191818101908684111561466657600080fd5b8286015b8481101561360757805161467d81612f72565b835291830191830161466a565b600082601f83011261469b57600080fd5b815160206146ab61318683613580565b82815260059290921b840181019181810190868411156146ca57600080fd5b8286015b848110156136075780516001600160401b038111156146ed5760008081fd5b6146fb8986838b0101614626565b8452509183019183016146ce565b60006020828403121561471b57600080fd5b81516001600160401b038082111561473257600080fd5b908301906080828603121561474657600080fd5b61474e612edb565b82518281111561475d57600080fd5b61476987828601614626565b82525060208301518281111561477e57600080fd5b61478a87828601614626565b6020830152506040830151828111156147a257600080fd5b6147ae87828601614626565b6040830152506060830151828111156147c657600080fd5b6147d28782860161468a565b60608301525095945050505050565b6000826147f0576147f0613bb9565b500690565b6000821982111561480857614808613b80565b500190565b8381526060602082015260006148266060830185613367565b82810360408401526148388185614236565b9695505050505050565b60006101a061ffff8616835280602084015261486081840186613367565b845180516040860152602001516060850152915061487b9050565b602083015161488d60808401826142de565b5060408301516148a16101008401826142de565b5063ffffffff60608401511661018083015294935050505056fe456967656e444143657274566572696669636174696f6e56314c69622e5f7665a2646970667358221220f7d5d773ba63891f072350d2a910acc3c010a302368cdaa426cf698a9364266764736f6c634300080c0033",
}

// ContractEigenDACertVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierMetaData.ABI instead.
var ContractEigenDACertVerifierABI = ContractEigenDACertVerifierMetaData.ABI

// ContractEigenDACertVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierMetaData.Bin instead.
var ContractEigenDACertVerifierBin = ContractEigenDACertVerifierMetaData.Bin

// DeployContractEigenDACertVerifier deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifier to it.
func DeployContractEigenDACertVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDABatchMetadataStorage common.Address, _eigenDASignatureVerifier common.Address, _eigenDARelayRegistry common.Address, _operatorStateRetriever common.Address, _registryCoordinator common.Address, _securityThresholdsV2 SecurityThresholds, _quorumNumbersRequiredV2 []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifier, error) {
	parsed, err := ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierBin), backend, _eigenDAThresholdRegistry, _eigenDABatchMetadataStorage, _eigenDASignatureVerifier, _eigenDARelayRegistry, _operatorStateRetriever, _registryCoordinator, _securityThresholdsV2, _quorumNumbersRequiredV2)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifier{ContractEigenDACertVerifierCaller: ContractEigenDACertVerifierCaller{contract: contract}, ContractEigenDACertVerifierTransactor: ContractEigenDACertVerifierTransactor{contract: contract}, ContractEigenDACertVerifierFilterer: ContractEigenDACertVerifierFilterer{contract: contract}}, nil
}

// ContractEigenDACertVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifier struct {
	ContractEigenDACertVerifierCaller     // Read-only binding to the contract
	ContractEigenDACertVerifierTransactor // Write-only binding to the contract
	ContractEigenDACertVerifierFilterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierSession struct {
	Contract     *ContractEigenDACertVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierCallerSession struct {
	Contract *ContractEigenDACertVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractEigenDACertVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierTransactorSession struct {
	Contract     *ContractEigenDACertVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRaw struct {
	Contract *ContractEigenDACertVerifier // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCallerRaw struct {
	Contract *ContractEigenDACertVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactorRaw struct {
	Contract *ContractEigenDACertVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifier creates a new instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifier(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifier, error) {
	contract, err := bindContractEigenDACertVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifier{ContractEigenDACertVerifierCaller: ContractEigenDACertVerifierCaller{contract: contract}, ContractEigenDACertVerifierTransactor: ContractEigenDACertVerifierTransactor{contract: contract}, ContractEigenDACertVerifierFilterer: ContractEigenDACertVerifierFilterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierCaller creates a new read-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierCaller, error) {
	contract, err := bindContractEigenDACertVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierCaller{contract: contract}, nil
}

// NewContractEigenDACertVerifierTransactor creates a new write-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierTransactor, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierTransactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierFilterer creates a new log filterer instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierFilterer, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierFilterer{contract: contract}, nil
}

// bindContractEigenDACertVerifier binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDABatchMetadataStorageV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDABatchMetadataStorageV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDARelayRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDARelayRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDASignatureVerifierV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDASignatureVerifierV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDAThresholdRegistryV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistryV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDAThresholdRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifier.CallOpts)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(VersionedBlobParams)).(*VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDACertVerifier.Contract.GetBlobParams(&_ContractEigenDACertVerifier.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDACertVerifier.Contract.GetBlobParams(&_ContractEigenDACertVerifier.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifier.CallOpts, quorumNumber)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) OperatorStateRetrieverV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "operatorStateRetrieverV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifier.CallOpts)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumNumbersRequiredV2(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumNumbersRequiredV2")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifier.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) RegistryCoordinatorV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "registryCoordinatorV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifier.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) SecurityThresholdsV2(opts *bind.CallOpts) (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "securityThresholdsV2")

	outstruct := new(struct {
		ConfirmationThreshold uint8
		AdversaryThreshold    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConfirmationThreshold = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.AdversaryThreshold = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifier.CallOpts)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertSecurityParams(opts *bind.CallOpts, blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertSecurityParams", blobParams, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifier.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifier.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertSecurityParams0(opts *bind.CallOpts, version uint16, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertSecurityParams0", version, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifier.CallOpts, version, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifier.CallOpts, version, securityThresholds)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV1(&_ContractEigenDACertVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV1(&_ContractEigenDACertVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2(&_ContractEigenDACertVerifier.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2(&_ContractEigenDACertVerifier.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertV2ForZKProof(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertV2ForZKProof", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertV2ForZKProof(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifier.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertV2ForZKProof(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifier.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertV2FromSignedBatch", signedBatch, blobInclusionInfo)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifier.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifier.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "verifyDACertsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifier.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifier.CallOpts, blobHeaders, blobVerificationProofs)
}

// ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 SecurityThresholds
	NewDefaultSecurityThresholdsV2      SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractEigenDACertVerifier.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated)
				if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDefaultSecurityThresholdsV2Updated is a log parse operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractEigenDACertVerifierDefaultSecurityThresholdsV2Updated)
	if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractEigenDACertVerifier.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseQuorumAdversaryThresholdPercentagesUpdated is a log parse operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDACertVerifierQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractEigenDACertVerifier.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseQuorumConfirmationThresholdPercentagesUpdated is a log parse operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDACertVerifierQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractEigenDACertVerifierQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierQuorumNumbersRequiredUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEigenDACertVerifierQuorumNumbersRequiredUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierQuorumNumbersRequiredUpdatedIterator{contract: _ContractEigenDACertVerifier.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDACertVerifier.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierQuorumNumbersRequiredUpdated)
				if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseQuorumNumbersRequiredUpdated is a log parse operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractEigenDACertVerifierQuorumNumbersRequiredUpdated, error) {
	event := new(ContractEigenDACertVerifierQuorumNumbersRequiredUpdated)
	if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierVersionedBlobParamsAddedIterator struct {
	Event *ContractEigenDACertVerifierVersionedBlobParamsAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractEigenDACertVerifierVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierVersionedBlobParamsAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractEigenDACertVerifierVersionedBlobParamsAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractEigenDACertVerifierVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractEigenDACertVerifierVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDACertVerifier.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierVersionedBlobParamsAddedIterator{contract: _ContractEigenDACertVerifier.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDACertVerifier.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierVersionedBlobParamsAdded)
				if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVersionedBlobParamsAdded is a log parse operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractEigenDACertVerifierVersionedBlobParamsAdded, error) {
	event := new(ContractEigenDACertVerifierVersionedBlobParamsAdded)
	if err := _ContractEigenDACertVerifier.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
