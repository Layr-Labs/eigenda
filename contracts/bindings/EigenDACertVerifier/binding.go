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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequiredV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620081683803806200816883398181016040528101906200003891906200069b565b8773ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508673ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250508573ffffffffffffffffffffffffffffffffffffffff1660c08173ffffffffffffffffffffffffffffffffffffffff16815250508473ffffffffffffffffffffffffffffffffffffffff1660e08173ffffffffffffffffffffffffffffffffffffffff16815250508373ffffffffffffffffffffffffffffffffffffffff166101008173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff166101208173ffffffffffffffffffffffffffffffffffffffff1681525050816000808201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548160ff021916908360ff1602179055509050508060019080519060200190620001d1929190620001e0565b505050505050505050620007e9565b828054620001ee90620007b3565b90600052602060002090601f0160209004810192826200021257600085556200025e565b82601f106200022d57805160ff19168380011785556200025e565b828001600101855582156200025e579182015b828111156200025d57825182559160200191906001019062000240565b5b5090506200026d919062000271565b5090565b5b808211156200028c57600081600090555060010162000272565b5090565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620002d182620002a4565b9050919050565b6000620002e582620002c4565b9050919050565b620002f781620002d8565b81146200030357600080fd5b50565b6000815190506200031781620002ec565b92915050565b60006200032a82620002c4565b9050919050565b6200033c816200031d565b81146200034857600080fd5b50565b6000815190506200035c8162000331565b92915050565b60006200036f82620002c4565b9050919050565b620003818162000362565b81146200038d57600080fd5b50565b600081519050620003a18162000376565b92915050565b6000620003b482620002c4565b9050919050565b620003c681620003a7565b8114620003d257600080fd5b50565b600081519050620003e681620003bb565b92915050565b6000620003f982620002c4565b9050919050565b6200040b81620003ec565b81146200041757600080fd5b50565b6000815190506200042b8162000400565b92915050565b60006200043e82620002c4565b9050919050565b620004508162000431565b81146200045c57600080fd5b50565b600081519050620004708162000445565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620004c6826200047b565b810181811067ffffffffffffffff82111715620004e857620004e76200048c565b5b80604052505050565b6000620004fd62000290565b90506200050b8282620004bb565b919050565b600060ff82169050919050565b620005288162000510565b81146200053457600080fd5b50565b60008151905062000548816200051d565b92915050565b60006040828403121562000567576200056662000476565b5b620005736040620004f1565b90506000620005858482850162000537565b60008301525060206200059b8482850162000537565b60208301525092915050565b600080fd5b600080fd5b600067ffffffffffffffff821115620005cf57620005ce6200048c565b5b620005da826200047b565b9050602081019050919050565b60005b8381101562000607578082015181840152602081019050620005ea565b8381111562000617576000848401525b50505050565b6000620006346200062e84620005b1565b620004f1565b905082815260208101848484011115620006535762000652620005ac565b5b62000660848285620005e7565b509392505050565b600082601f83011262000680576200067f620005a7565b5b8151620006928482602086016200061d565b91505092915050565b600080600080600080600080610120898b031215620006bf57620006be6200029a565b5b6000620006cf8b828c0162000306565b9850506020620006e28b828c016200034b565b9750506040620006f58b828c0162000390565b9650506060620007088b828c01620003d5565b95505060806200071b8b828c016200041a565b94505060a06200072e8b828c016200045f565b93505060c0620007418b828c016200054e565b92505061010089015167ffffffffffffffff8111156200076657620007656200029f565b5b620007748b828c0162000668565b9150509295985092959890939650565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620007cc57607f821691505b60208210811415620007e357620007e262000784565b5b50919050565b60805160a05160c05160e0516101005161012051617887620008e1600039600081816108350152818161098d0152610f2c015260008181610814015281816109450152610f0b01526000818161071e015281816107f3015281816109b10152610a720152600081816106fd015281816107d201528181610a510152610eda0152600081816106880152818161096901526109f90152600081816104690152818161050c015281816105c301528181610667015281816106dc015281816107b1015281816109d801528181610a3001528181610b9801528181610cc101528181610d7201528181610e390152610f6801526178876000f3fe608060405234801561001057600080fd5b506004361061014d5760003560e01c80637d644cad116100c3578063e15234ff1161007c578063e15234ff1461038c578063ed0450ae146103aa578063ee6c3bcf146103c9578063efd4532b146103f9578063f25de3f814610417578063f8c66814146104475761014d565b80637d644cad146102de578063813c2eb0146102fa5780638687feae14610316578063b74d787114610334578063bafa910714610352578063ccb7cd0d146103705761014d565b8063415ef61411610115578063415ef6141461021a578063421c02221461024a5780634ca22c3f14610266578063640f65d9146102845780636d14a987146102a257806372276443146102c05761014d565b8063048886d2146101525780631429c7c214610182578063143eb4d9146101b25780632ecfe72b146101ce57806331a3479a146101fe575b600080fd5b61016c60048036038101906101679190612d1c565b610465565b6040516101799190612d64565b60405180910390f35b61019c60048036038101906101979190612d1c565b610508565b6040516101a99190612d8e565b60405180910390f35b6101cc60048036038101906101c79190612f2f565b6105ab565b005b6101e860048036038101906101e39190612fa9565b6105b9565b6040516101f59190613036565b60405180910390f35b6102186004803603810190610213919061310c565b610662565b005b610234600480360381019061022f91906132a5565b6106be565b6040516102419190612d64565b60405180910390f35b610264600480360381019061025f919061337f565b6107ac565b005b61026e610943565b60405161027b9190613476565b60405180910390f35b61028c610967565b60405161029991906134b2565b60405180910390f35b6102aa61098b565b6040516102b791906134ee565b60405180910390f35b6102c86109af565b6040516102d5919061352a565b60405180910390f35b6102f860048036038101906102f39190613583565b6109d3565b005b610314600480360381019061030f91906132a5565b610a2b565b005b61031e610b94565b60405161032b9190613683565b60405180910390f35b61033c610c2f565b6040516103499190613683565b60405180910390f35b61035a610cbd565b6040516103679190613683565b60405180910390f35b61038a600480360381019061038591906136a5565b610d58565b005b610394610d6e565b6040516103a19190613683565b60405180910390f35b6103b2610e09565b6040516103c09291906136e5565b60405180910390f35b6103e360048036038101906103de9190612d1c565b610e35565b6040516103f09190612d8e565b60405180910390f35b610401610ed8565b60405161040e919061372f565b60405180910390f35b610431600480360381019061042c919061374a565b610efc565b60405161043e9190613b9c565b60405180910390f35b61044f610f66565b60405161045c9190613bdf565b60405180910390f35b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663048886d2836040518263ffffffff1660e01b81526004016104c09190612d8e565b602060405180830381865afa1580156104dd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105019190613c26565b9050919050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631429c7c2836040518263ffffffff1660e01b81526004016105639190612d8e565b602060405180830381865afa158015610580573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a49190613c68565b9050919050565b6105b58282610f8a565b5050565b6105c1612bec565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632ecfe72b836040518263ffffffff1660e01b815260040161061a9190613ca4565b606060405180830381865afa158015610637573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061065b9190613d38565b9050919050565b6106b87f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000868686866106b3610d6e565b61109c565b50505050565b600073__$b5e4fc231d81cbd1f60cac3db25e9b5865$__637c45e9fd7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000898989600060018b6040518a63ffffffff1660e01b815260040161076999989796959493929190614897565b60006040518083038186803b15801561078157600080fd5b505af4925050508015610792575060015b61079f57600090506107a4565b600190505b949350505050565b61093f7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008761085e90614db6565b876108689061502d565b60006040518060400160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900460ff1660ff1660ff1681525050600180546108bc90614786565b80601f01602080910402602001604051908101604052809291908181526020018280546108e890614786565b80156109355780601f1061090a57610100808354040283529160200191610935565b820191906000526020600020905b81548152906001019060200180831161091857829003601f168201915b5050505050611a5f565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b610a277f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610a22610d6e565b611a93565b5050565b610b8e7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000087803603810190610aa29190615040565b87610aac9061502d565b87610ab6906152c3565b60006040518060400160405290816000820160009054906101000a900460ff1660ff1660ff1681526020016000820160019054906101000a900460ff1660ff1660ff168152505060018054610b0a90614786565b80601f0160208091040260200160405190810160405280929190818152602001828054610b3690614786565b8015610b835780601f10610b5857610100808354040283529160200191610b83565b820191906000526020600020905b815481529060010190602001808311610b6657829003601f168201915b505050505089612141565b50505050565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610c01573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610c2a9190615346565b905090565b60018054610c3c90614786565b80601f0160208091040260200160405190810160405280929190818152602001828054610c6890614786565b8015610cb55780601f10610c8a57610100808354040283529160200191610cb5565b820191906000526020600020905b815481529060010190602001808311610c9857829003601f168201915b505050505081565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610d2a573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610d539190615346565b905090565b610d6a610d64836105b9565b82610f8a565b5050565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610ddb573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610e049190615346565b905090565b60008060000160009054906101000a900460ff16908060000160019054906101000a900460ff16905082565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663ee6c3bcf836040518263ffffffff1660e01b8152600401610e909190612d8e565b602060405180830381865afa158015610ead573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ed19190613c68565b9050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b610f04612c1c565b6000610f5a7f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000085610f5590614db6565b61249a565b50905080915050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b806020015160ff16816000015160ff1611610fda576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fd19061545e565b60405180910390fd5b600081602001518260000151610ff091906154ad565b60ff1690506000836020015163ffffffff16846040015160ff1683620f42406110199190615510565b6110239190615510565b6127106110309190615541565b61103a9190615575565b9050612710846000015161104e91906155cf565b63ffffffff16811015611096576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161108d906156a5565b60405180910390fd5b50505050565b8282905085859050146110e4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110db90615783565b60405180910390fd5b60008773ffffffffffffffffffffffffffffffffffffffff1663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015611131573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061115a9190615346565b905060005b86869050811015611a54578773ffffffffffffffffffffffffffffffffffffffff1663eccbbfc9868684818110611199576111986157a3565b5b90506020028101906111ab91906157e1565b60000160208101906111bd9190615809565b6040518263ffffffff1660e01b81526004016111d99190615845565b602060405180830381865afa1580156111f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061121a9190615875565b61125f8686848181106112305761122f6157a3565b5b905060200281019061124291906157e1565b806040019061125191906158a2565b61125a906159fa565b6126e9565b1461129f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161129690615acb565b60405180910390fd5b6113fe8585838181106112b5576112b46157a3565b5b90506020028101906112c791906157e1565b80606001906112d69190615aeb565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505086868481811061132c5761132b6157a3565b5b905060200281019061133e91906157e1565b806040019061134d91906158a2565b806000019061135c9190615b4e565b600001356113968a8a86818110611376576113756157a3565b5b90506020028101906113889190615b76565b61139190615d59565b61272f565b6040516020016113a69190615d8d565b604051602081830303815290604052805190602001208888868181106113cf576113ce6157a3565b5b90506020028101906113e191906157e1565b60200160208101906113f39190615809565b63ffffffff1661275f565b61143d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161143490615e40565b60405180910390fd5b600080600090505b888884818110611458576114576157a3565b5b905060200281019061146a9190615b76565b80606001906114799190615e60565b90508110156119f057888884818110611495576114946157a3565b5b90506020028101906114a79190615b76565b80606001906114b69190615e60565b828181106114c7576114c66157a3565b5b90506080020160000160208101906114df9190612d1c565b60ff168787858181106114f5576114f46157a3565b5b905060200281019061150791906157e1565b806040019061151691906158a2565b80600001906115259190615b4e565b80602001906115349190615aeb565b898987818110611547576115466157a3565b5b905060200281019061155991906157e1565b80608001906115689190615aeb565b85818110611579576115786157a3565b5b9050013560f81c60f81b60f81c60ff16818110611599576115986157a3565b5b9050013560f81c60f81b60f81c60ff16146115e9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115e090615f5b565b60405180910390fd5b8888848181106115fc576115fb6157a3565b5b905060200281019061160e9190615b76565b806060019061161d9190615e60565b8281811061162e5761162d6157a3565b5b90506080020160200160208101906116469190612d1c565b60ff1689898581811061165c5761165b6157a3565b5b905060200281019061166e9190615b76565b806060019061167d9190615e60565b8381811061168e5761168d6157a3565b5b90506080020160400160208101906116a69190612d1c565b60ff16116116e9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016116e090616013565b60405180910390fd5b838989858181106116fd576116fc6157a3565b5b905060200281019061170f9190615b76565b806060019061171e9190615e60565b8381811061172f5761172e6157a3565b5b90506080020160000160208101906117479190612d1c565b60ff168151811061175b5761175a6157a3565b5b602001015160f81c60f81b60f81c60ff1689898581811061177f5761177e6157a3565b5b90506020028101906117919190615b76565b80606001906117a09190615e60565b838181106117b1576117b06157a3565b5b90506080020160400160208101906117c99190612d1c565b60ff16101561180d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611804906160f1565b60405180910390fd5b8888848181106118205761181f6157a3565b5b90506020028101906118329190615b76565b80606001906118419190615e60565b82818110611852576118516157a3565b5b905060800201604001602081019061186a9190612d1c565b60ff168787858181106118805761187f6157a3565b5b905060200281019061189291906157e1565b80604001906118a191906158a2565b80600001906118b09190615b4e565b80604001906118bf9190615aeb565b8989878181106118d2576118d16157a3565b5b90506020028101906118e491906157e1565b80608001906118f39190615aeb565b85818110611904576119036157a3565b5b9050013560f81c60f81b60f81c60ff16818110611924576119236157a3565b5b9050013560f81c60f81b60f81c60ff161015611975576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161196c906160f1565b60405180910390fd5b6119db828a8a8681811061198c5761198b6157a3565b5b905060200281019061199e9190615b76565b80606001906119ad9190615e60565b848181106119be576119bd6157a3565b5b90506080020160000160208101906119d69190612d1c565b612778565b915080806119e890616111565b915050611445565b50611a036119fd8561278c565b826128b3565b611a42576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a3990616218565b60405180910390fd5b5080611a4d90616111565b905061115f565b505050505050505050565b600080611a6d88888861249a565b91509150611a868b8b8b896000015189878a8a89612141565b5050505050505050505050565b8373ffffffffffffffffffffffffffffffffffffffff1663eccbbfc9836000016020810190611ac29190615809565b6040518263ffffffff1660e01b8152600401611ade9190615845565b602060405180830381865afa158015611afb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b1f9190615875565b611b40838060400190611b3291906158a2565b611b3b906159fa565b6126e9565b14611b80576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611b77906162f6565b60405180910390fd5b611c4f828060600190611b939190615aeb565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050838060400190611be691906158a2565b8060000190611bf59190615b4e565b60000135611c0b86611c0690615d59565b61272f565b604051602001611c1b9190615d8d565b60405160208183030381529060405280519060200120856020016020810190611c449190615809565b63ffffffff1661275f565b611c8e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611c85906163ae565b60405180910390fd5b600080600090505b848060600190611ca69190615e60565b90508110156120e757848060600190611cbf9190615e60565b82818110611cd057611ccf6157a3565b5b9050608002016000016020810190611ce89190612d1c565b60ff16848060400190611cfb91906158a2565b8060000190611d0a9190615b4e565b8060200190611d199190615aeb565b868060800190611d299190615aeb565b85818110611d3a57611d396157a3565b5b9050013560f81c60f81b60f81c60ff16818110611d5a57611d596157a3565b5b9050013560f81c60f81b60f81c60ff1614611daa576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611da190616466565b60405180910390fd5b848060600190611dba9190615e60565b82818110611dcb57611dca6157a3565b5b9050608002016020016020810190611de39190612d1c565b60ff16858060600190611df69190615e60565b83818110611e0757611e066157a3565b5b9050608002016040016020810190611e1f9190612d1c565b60ff1611611e62576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611e599061651e565b60405180910390fd5b8673ffffffffffffffffffffffffffffffffffffffff16631429c7c2868060600190611e8e9190615e60565b84818110611e9f57611e9e6157a3565b5b9050608002016000016020810190611eb79190612d1c565b6040518263ffffffff1660e01b8152600401611ed39190612d8e565b602060405180830381865afa158015611ef0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f149190613c68565b60ff16858060600190611f279190615e60565b83818110611f3857611f376157a3565b5b9050608002016040016020810190611f509190612d1c565b60ff161015611f94576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f8b906165d6565b60405180910390fd5b848060600190611fa49190615e60565b82818110611fb557611fb46157a3565b5b9050608002016040016020810190611fcd9190612d1c565b60ff16848060400190611fe091906158a2565b8060000190611fef9190615b4e565b8060400190611ffe9190615aeb565b86806080019061200e9190615aeb565b8581811061201f5761201e6157a3565b5b9050013560f81c60f81b60f81c60ff1681811061203f5761203e6157a3565b5b9050013560f81c60f81b60f81c60ff161015612090576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612087906165d6565b60405180910390fd5b6120d2828680606001906120a49190615e60565b848181106120b5576120b46157a3565b5b90506080020160000160208101906120cd9190612d1c565b612778565b915080806120df90616111565b915050611c96565b506120fa6120f48361278c565b826128b3565b612139576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612130906166b4565b60405180910390fd5b505050505050565b6121918560400151876000015161215b88600001516128c2565b60405160200161216b9190615d8d565b60405160208183030381529060405280519060200120886020015163ffffffff1661275f565b6121d0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016121c79061676c565b60405180910390fd5b6000808973ffffffffffffffffffffffffffffffffffffffff16636efb46366121f88a61290a565b858b602001518a6040518563ffffffff1660e01b815260040161221e949392919061679b565b600060405180830381865afa15801561223b573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190612264919061697d565b9150915061227a8988600001516040015161293a565b6123098b73ffffffffffffffffffffffffffffffffffffffff16632ecfe72b896000015160000151600001516040518263ffffffff1660e01b81526004016122c29190613ca4565b606060405180830381865afa1580156122df573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123039190613d38565b86610f8a565b600080600090505b84518110156123d757866000015160ff1684602001518281518110612339576123386157a3565b5b602002602001015161234b91906169d9565b6bffffffffffffffffffffffff16606485600001518381518110612372576123716157a3565b5b60200260200101516bffffffffffffffffffffffff166123929190615575565b106123c4576123c1828683815181106123ae576123ad6157a3565b5b602001015160f81c60f81b60f81c612778565b91505b80806123cf90616111565b915050612311565b5060006123ef8960000151600001516020015161278c565b90506123fb81836128b3565b61243a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161243190616add565b60405180910390fd5b61244c6124468761278c565b826128b3565b61248b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161248290616bbb565b60405180910390fd5b50505050505050505050505050565b6124a2612c1c565b606060008360200151600001515167ffffffffffffffff8111156124c9576124c8612dbf565b5b6040519080825280602002602001820160405280156124f75781602001602082028036833780820191505090505b50905060005b8460200151600001515181101561256957612539856020015160000151828151811061252c5761252b6157a3565b5b6020026020010151612a5f565b82828151811061254c5761254b6157a3565b5b6020026020010181815250508061256290616111565b90506124fd565b5060005b846020015160800151518110156125d55782856020015160800151828151811061259a576125996157a3565b5b60200260200101516040516020016125b3929190616c4d565b6040516020818303038152906040529250806125ce90616111565b905061256d565b5060008673ffffffffffffffffffffffffffffffffffffffff16634f739f748787600001516020015186866040518563ffffffff1660e01b815260040161261f9493929190616d33565b600060405180830381865afa15801561263c573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906126659190616fba565b905080600001518460000181905250846020015160000151846020018190525084602001516020015184604001819052508460200151606001518460600181905250846020015160400151846080018190525080602001518460a0018190525080604001518460c0018190525080606001518460e001819052505050935093915050565b6000612728826000015160405160200161270391906170b7565b6040516020818303038152906040528051906020012083602001518460400151612a7a565b9050919050565b600081604051602001612742919061722d565b604051602081830303815290604052805190602001209050919050565b60008361276d868585612ab0565b149050949350505050565b60008160ff166001901b8317905092915050565b6000610100825111156127d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016127cb906172e7565b60405180910390fd5b6000825114156127e757600090506128ae565b600080836000815181106127fe576127fd6157a3565b5b602001015160f81c60f81b60f81c60ff166001901b91506000600190505b84518110156128a757848181518110612838576128376157a3565b5b602001015160f81c60f81b60f81c60ff166001901b9150828211612891576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016128889061739f565b60405180910390fd5b8183179250806128a090616111565b905061281c565b5081925050505b919050565b60008282841614905092915050565b60006128d18260000151612b80565b826020015183604001516040516020016128ed9392919061742e565b604051602081830303815290604052805190602001209050919050565b60008160405160200161291d91906174a2565b604051602081830303815290604052805190602001209050919050565b60005b8151811015612a5a57600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1663b5a872da84848151811061298d5761298c6157a3565b5b60200260200101516040518263ffffffff1660e01b81526004016129b19190615845565b602060405180830381865afa1580156129ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129f291906174fb565b73ffffffffffffffffffffffffffffffffffffffff161415612a49576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612a40906175c0565b60405180910390fd5b80612a5390616111565b905061293d565b505050565b60008151600052816020015160205260406000209050919050565b6000838383604051602001612a9193929190617616565b6040516020818303038152906040528051906020012090509392505050565b60008060208551612ac19190617653565b14612b01576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612af89061771c565b60405180910390fd5b60008390506000602090505b85518111612b74576000600285612b249190617653565b1415612b4757816000528086015160205260406000209150600284049350612b60565b8086015160005281602052604060002091506002840493505b602081612b6d919061773c565b9050612b0d565b50809150509392505050565b6000816000015182602001518360400151604051602001612ba3939291906177e9565b604051602081830303815290604052805190602001208260600151604051602001612bcf929190617828565b604051602081830303815290604052805190602001209050919050565b6040518060600160405280600063ffffffff168152602001600063ffffffff168152602001600060ff1681525090565b604051806101000160405280606081526020016060815260200160608152602001612c45612c6d565b8152602001612c52612c93565b81526020016060815260200160608152602001606081525090565b6040518060400160405280612c80612cad565b8152602001612c8d612cad565b81525090565b604051806040016040528060008152602001600081525090565b6040518060400160405280600290602082028036833780820191505090505090565b6000604051905090565b600080fd5b600080fd5b600060ff82169050919050565b612cf981612ce3565b8114612d0457600080fd5b50565b600081359050612d1681612cf0565b92915050565b600060208284031215612d3257612d31612cd9565b5b6000612d4084828501612d07565b91505092915050565b60008115159050919050565b612d5e81612d49565b82525050565b6000602082019050612d796000830184612d55565b92915050565b612d8881612ce3565b82525050565b6000602082019050612da36000830184612d7f565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b612df782612dae565b810181811067ffffffffffffffff82111715612e1657612e15612dbf565b5b80604052505050565b6000612e29612ccf565b9050612e358282612dee565b919050565b600080fd5b600063ffffffff82169050919050565b612e5881612e3f565b8114612e6357600080fd5b50565b600081359050612e7581612e4f565b92915050565b600060608284031215612e9157612e90612da9565b5b612e9b6060612e1f565b90506000612eab84828501612e66565b6000830152506020612ebf84828501612e66565b6020830152506040612ed384828501612d07565b60408301525092915050565b600060408284031215612ef557612ef4612da9565b5b612eff6040612e1f565b90506000612f0f84828501612d07565b6000830152506020612f2384828501612d07565b60208301525092915050565b60008060a08385031215612f4657612f45612cd9565b5b6000612f5485828601612e7b565b9250506060612f6585828601612edf565b9150509250929050565b600061ffff82169050919050565b612f8681612f6f565b8114612f9157600080fd5b50565b600081359050612fa381612f7d565b92915050565b600060208284031215612fbf57612fbe612cd9565b5b6000612fcd84828501612f94565b91505092915050565b612fdf81612e3f565b82525050565b612fee81612ce3565b82525050565b60608201600082015161300a6000850182612fd6565b50602082015161301d6020850182612fd6565b5060408201516130306040850182612fe5565b50505050565b600060608201905061304b6000830184612ff4565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261307657613075613051565b5b8235905067ffffffffffffffff81111561309357613092613056565b5b6020830191508360208202830111156130af576130ae61305b565b5b9250929050565b60008083601f8401126130cc576130cb613051565b5b8235905067ffffffffffffffff8111156130e9576130e8613056565b5b6020830191508360208202830111156131055761310461305b565b5b9250929050565b6000806000806040858703121561312657613125612cd9565b5b600085013567ffffffffffffffff81111561314457613143612cde565b5b61315087828801613060565b9450945050602085013567ffffffffffffffff81111561317357613172612cde565b5b61317f878288016130b6565b925092505092959194509250565b600080fd5b6000604082840312156131a8576131a761318d565b5b81905092915050565b6000606082840312156131c7576131c661318d565b5b81905092915050565b600061018082840312156131e7576131e661318d565b5b81905092915050565b600080fd5b600067ffffffffffffffff8211156132105761320f612dbf565b5b61321982612dae565b9050602081019050919050565b82818337600083830152505050565b6000613248613243846131f5565b612e1f565b905082815260208101848484011115613264576132636131f0565b5b61326f848285613226565b509392505050565b600082601f83011261328c5761328b613051565b5b813561329c848260208601613235565b91505092915050565b60008060008060a085870312156132bf576132be612cd9565b5b60006132cd87828801613192565b945050604085013567ffffffffffffffff8111156132ee576132ed612cde565b5b6132fa878288016131b1565b935050606085013567ffffffffffffffff81111561331b5761331a612cde565b5b613327878288016131d0565b925050608085013567ffffffffffffffff81111561334857613347612cde565b5b61335487828801613277565b91505092959194509250565b6000606082840312156133765761337561318d565b5b81905092915050565b6000806040838503121561339657613395612cd9565b5b600083013567ffffffffffffffff8111156133b4576133b3612cde565b5b6133c085828601613360565b925050602083013567ffffffffffffffff8111156133e1576133e0612cde565b5b6133ed858286016131b1565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600061343c613437613432846133f7565b613417565b6133f7565b9050919050565b600061344e82613421565b9050919050565b600061346082613443565b9050919050565b61347081613455565b82525050565b600060208201905061348b6000830184613467565b92915050565b600061349c82613443565b9050919050565b6134ac81613491565b82525050565b60006020820190506134c760008301846134a3565b92915050565b60006134d882613443565b9050919050565b6134e8816134cd565b82525050565b600060208201905061350360008301846134df565b92915050565b600061351482613443565b9050919050565b61352481613509565b82525050565b600060208201905061353f600083018461351b565b92915050565b60006080828403121561355b5761355a61318d565b5b81905092915050565b600060a0828403121561357a5761357961318d565b5b81905092915050565b6000806040838503121561359a57613599612cd9565b5b600083013567ffffffffffffffff8111156135b8576135b7612cde565b5b6135c485828601613545565b925050602083013567ffffffffffffffff8111156135e5576135e4612cde565b5b6135f185828601613564565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561363557808201518184015260208101905061361a565b83811115613644576000848401525b50505050565b6000613655826135fb565b61365f8185613606565b935061366f818560208601613617565b61367881612dae565b840191505092915050565b6000602082019050818103600083015261369d818461364a565b905092915050565b600080606083850312156136bc576136bb612cd9565b5b60006136ca85828601612f94565b92505060206136db85828601612edf565b9150509250929050565b60006040820190506136fa6000830185612d7f565b6137076020830184612d7f565b9392505050565b600061371982613443565b9050919050565b6137298161370e565b82525050565b60006020820190506137446000830184613720565b92915050565b6000602082840312156137605761375f612cd9565b5b600082013567ffffffffffffffff81111561377e5761377d612cde565b5b61378a84828501613360565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b60006137cb8383612fd6565b60208301905092915050565b6000602082019050919050565b60006137ef82613793565b6137f9818561379e565b9350613804836137af565b8060005b8381101561383557815161381c88826137bf565b9750613827836137d7565b925050600181019050613808565b5085935050505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6000819050919050565b6138818161386e565b82525050565b60408201600082015161389d6000850182613878565b5060208201516138b06020850182613878565b50505050565b60006138c28383613887565b60408301905092915050565b6000602082019050919050565b60006138e682613842565b6138f0818561384d565b93506138fb8361385e565b8060005b8381101561392c57815161391388826138b6565b975061391e836138ce565b9250506001810190506138ff565b5085935050505092915050565b600060029050919050565b600081905092915050565b6000819050919050565b60006139658383613878565b60208301905092915050565b6000602082019050919050565b61398781613939565b6139918184613944565b925061399c8261394f565b8060005b838110156139cd5781516139b48782613959565b96506139bf83613971565b9250506001810190506139a0565b505050505050565b6080820160008201516139eb600085018261397e565b5060208201516139fe604085018261397e565b50505050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6000613a3c83836137e4565b905092915050565b6000602082019050919050565b6000613a5c82613a04565b613a668185613a0f565b935083602082028501613a7885613a20565b8060005b85811015613ab45784840389528151613a958582613a30565b9450613aa083613a44565b925060208a01995050600181019050613a7c565b50829750879550505050505092915050565b6000610180830160008301518482036000860152613ae482826137e4565b91505060208301518482036020860152613afe82826138db565b91505060408301518482036040860152613b1882826138db565b9150506060830151613b2d60608601826139d5565b506080830151613b4060e0860182613887565b5060a0830151848203610120860152613b5982826137e4565b91505060c0830151848203610140860152613b7482826137e4565b91505060e0830151848203610160860152613b8f8282613a51565b9150508091505092915050565b60006020820190508181036000830152613bb68184613ac6565b905092915050565b6000613bc982613443565b9050919050565b613bd981613bbe565b82525050565b6000602082019050613bf46000830184613bd0565b92915050565b613c0381612d49565b8114613c0e57600080fd5b50565b600081519050613c2081613bfa565b92915050565b600060208284031215613c3c57613c3b612cd9565b5b6000613c4a84828501613c11565b91505092915050565b600081519050613c6281612cf0565b92915050565b600060208284031215613c7e57613c7d612cd9565b5b6000613c8c84828501613c53565b91505092915050565b613c9e81612f6f565b82525050565b6000602082019050613cb96000830184613c95565b92915050565b600081519050613cce81612e4f565b92915050565b600060608284031215613cea57613ce9612da9565b5b613cf46060612e1f565b90506000613d0484828501613cbf565b6000830152506020613d1884828501613cbf565b6020830152506040613d2c84828501613c53565b60408301525092915050565b600060608284031215613d4e57613d4d612cd9565b5b6000613d5c84828501613cd4565b91505092915050565b613d6e81613bbe565b82525050565b613d7d8161370e565b82525050565b613d8c81613509565b82525050565b6000819050919050565b613da581613d92565b8114613db057600080fd5b50565b600081359050613dc281613d9c565b92915050565b6000613dd76020840184613db3565b905092915050565b613de881613d92565b82525050565b6000613dfd6020840184612e66565b905092915050565b613e0e81612e3f565b82525050565b60408201613e256000830183613dc8565b613e326000850182613ddf565b50613e406020830183613dee565b613e4d6020850182613e05565b50505050565b600080fd5b600082356001606003833603038112613e7457613e73613e53565b5b82810191505092915050565b6000823560016101c003833603038112613e9d57613e9c613e53565b5b82810191505092915050565b6000613eb86020840184612f94565b905092915050565b613ec981612f6f565b82525050565b600080fd5b600080fd5b60008083356001602003843603038112613ef657613ef5613e53565b5b83810192508235915060208301925067ffffffffffffffff821115613f1e57613f1d613ecf565b5b600182023603841315613f3457613f33613ed4565b5b509250929050565b600082825260208201905092915050565b6000613f598385613f3c565b9350613f66838584613226565b613f6f83612dae565b840190509392505050565b600082905092915050565b600082905092915050565b613f998161386e565b8114613fa457600080fd5b50565b600081359050613fb681613f90565b92915050565b6000613fcb6020840184613fa7565b905092915050565b613fdc8161386e565b82525050565b60408201613ff36000830183613fbc565b6140006000850182613fd3565b5061400e6020830183613fbc565b61401b6020850182613fd3565b50505050565b600082905092915050565b600082905092915050565b61404360408383613226565b5050565b60808201614058600083018361402c565b6140656000850182614037565b50614073604083018361402c565b6140806040850182614037565b50505050565b61016082016140986000830183613f85565b6140a56000850182613fe2565b506140b36040830183614021565b6140c06040850182614047565b506140ce60c0830183614021565b6140db60c0850182614047565b506140ea610140830183613dee565b6140f8610140850182613e05565b50505050565b60006101c083016141126000840184613ea9565b61411f6000860182613ec0565b5061412d6020840184613ed9565b8583036020870152614140838284613f4d565b925050506141516040840184613f7a565b61415e6040860182614086565b5061416d6101a0840184613dc8565b61417b6101a0860182613ddf565b508091505092915050565b600080833560016020038436030381126141a3576141a2613e53565b5b83810192508235915060208301925067ffffffffffffffff8211156141cb576141ca613ecf565b5b6020820236038413156141e1576141e0613ed4565b5b509250929050565b600082825260208201905092915050565b6000819050919050565b60006142108383613e05565b60208301905092915050565b6000602082019050919050565b600061423583856141e9565b9350614240826141fa565b8060005b85811015614279576142568284613dee565b6142608882614204565b975061426b8361421c565b925050600181019050614244565b5085925050509392505050565b6000606083016142996000840184613e80565b84820360008601526142ab82826140fe565b9150506142bb6020840184613ed9565b85830360208701526142ce838284613f4d565b925050506142df6040840184614186565b85830360408701526142f2838284614229565b925050508091505092915050565b6000606083016143136000840184613e58565b84820360008601526143258282614286565b9150506143356020840184613dee565b6143426020860182613e05565b506143506040840184613ed9565b8583036040870152614363838284613f4d565b925050508091505092915050565b6000808335600160200384360303811261438e5761438d613e53565b5b83810192508235915060208301925067ffffffffffffffff8211156143b6576143b5613ecf565b5b6040820236038413156143cc576143cb613ed4565b5b509250929050565b600082825260208201905092915050565b6000819050919050565b60006143fb8383613fe2565b60408301905092915050565b6000604082019050919050565b600061442083856143d4565b935061442b826143e5565b8060005b85811015614464576144418284613f85565b61444b88826143ef565b975061445683614407565b92505060018101905061442f565b5085925050509392505050565b6000808335600160200384360303811261448e5761448d613e53565b5b83810192508235915060208301925067ffffffffffffffff8211156144b6576144b5613ecf565b5b6020820236038413156144cc576144cb613ed4565b5b509250929050565b600082825260208201905092915050565b6000819050919050565b60006144fc848484614229565b90509392505050565b6000602082019050919050565b600061451e83856144d4565b935083602084028501614530846144e5565b8060005b8781101561457657848403895261454b8284614186565b6145568682846144ef565b955061456184614505565b935060208b019a505050600181019050614534565b50829750879450505050509392505050565b6000610180830161459c6000840184614186565b85830360008701526145af838284614229565b925050506145c06020840184614371565b85830360208701526145d3838284614414565b925050506145e46040840184614371565b85830360408701526145f7838284614414565b925050506146086060840184614021565b6146156060860182614047565b5061462360e0840184613f85565b61463060e0860182613fe2565b5061463f610120840184614186565b858303610120870152614653838284614229565b92505050614665610140840184614186565b858303610140870152614679838284614229565b9250505061468b610160840184614471565b85830361016087015261469f838284614512565b925050508091505092915050565b60008160001c9050919050565b600060ff82169050919050565b60006146da6146d5836146ad565b6146ba565b9050919050565b6146ea81612ce3565b82525050565b60008160081c9050919050565b600061471061470b836146f0565b6146ba565b9050919050565b60408201600080830154905061472c816146c7565b61473960008601826146e1565b50614743816146fd565b61475060208601826146e1565b5050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061479e57607f821691505b602082108114156147b2576147b1614757565b5b50919050565b600082825260208201905092915050565b60008190508160005260206000209050919050565b600081546147eb81614786565b6147f581866147b8565b94506001821660008114614810576001811461482257614855565b60ff1983168652602086019350614855565b61482b856147c9565b60005b8381101561484d5781548189015260018201915060208101905061482e565b808801955050505b50505092915050565b6000614869826135fb565b61487381856147b8565b9350614883818560208601613617565b61488c81612dae565b840191505092915050565b6000610160820190506148ad600083018c613d65565b6148ba602083018b613d74565b6148c7604083018a613d83565b6148d46060830189613e14565b81810360a08301526148e68188614300565b905081810360c08301526148fa8187614588565b905061490960e0830186614717565b81810361012083015261491c81856147de565b9050818103610140830152614931818461485e565b90509a9950505050505050505050565b60006040828403121561495757614956612da9565b5b6149616040612e1f565b9050600061497184828501613db3565b600083015250602061498584828501612e66565b60208301525092915050565b600067ffffffffffffffff8211156149ac576149ab612dbf565b5b602082029050602081019050919050565b6000604082840312156149d3576149d2612da9565b5b6149dd6040612e1f565b905060006149ed84828501613fa7565b6000830152506020614a0184828501613fa7565b60208301525092915050565b6000614a20614a1b84614991565b612e1f565b90508083825260208201905060408402830185811115614a4357614a4261305b565b5b835b81811015614a6c5780614a5888826149bd565b845260208401935050604081019050614a45565b5050509392505050565b600082601f830112614a8b57614a8a613051565b5b8135614a9b848260208601614a0d565b91505092915050565b600067ffffffffffffffff821115614abf57614abe612dbf565b5b602082029050919050565b6000614add614ad884614aa4565b612e1f565b90508060208402830185811115614af757614af661305b565b5b835b81811015614b205780614b0c8882613fa7565b845260208401935050602081019050614af9565b5050509392505050565b600082601f830112614b3f57614b3e613051565b5b6002614b4c848285614aca565b91505092915050565b600060808284031215614b6b57614b6a612da9565b5b614b756040612e1f565b90506000614b8584828501614b2a565b6000830152506040614b9984828501614b2a565b60208301525092915050565b600067ffffffffffffffff821115614bc057614bbf612dbf565b5b602082029050602081019050919050565b6000614be4614bdf84614ba5565b612e1f565b90508083825260208201905060208402830185811115614c0757614c0661305b565b5b835b81811015614c305780614c1c8882612e66565b845260208401935050602081019050614c09565b5050509392505050565b600082601f830112614c4f57614c4e613051565b5b8135614c5f848260208601614bd1565b91505092915050565b60006101208284031215614c7f57614c7e612da9565b5b614c8960a0612e1f565b9050600082013567ffffffffffffffff811115614ca957614ca8612e3a565b5b614cb584828501614a76565b600083015250602082013567ffffffffffffffff811115614cd957614cd8612e3a565b5b614ce584828501614a76565b6020830152506040614cf9848285016149bd565b6040830152506080614d0d84828501614b55565b60608301525061010082013567ffffffffffffffff811115614d3257614d31612e3a565b5b614d3e84828501614c3a565b60808301525092915050565b600060608284031215614d6057614d5f612da9565b5b614d6a6040612e1f565b90506000614d7a84828501614941565b600083015250604082013567ffffffffffffffff811115614d9e57614d9d612e3a565b5b614daa84828501614c68565b60208301525092915050565b6000614dc23683614d4a565b9050919050565b60006101608284031215614de057614ddf612da9565b5b614dea6080612e1f565b90506000614dfa848285016149bd565b6000830152506040614e0e84828501614b55565b60208301525060c0614e2284828501614b55565b604083015250610140614e3784828501612e66565b60608301525092915050565b60006101c08284031215614e5a57614e59612da9565b5b614e646080612e1f565b90506000614e7484828501612f94565b600083015250602082013567ffffffffffffffff811115614e9857614e97612e3a565b5b614ea484828501613277565b6020830152506040614eb884828501614dc9565b6040830152506101a0614ecd84828501613db3565b60608301525092915050565b600060608284031215614eef57614eee612da9565b5b614ef96060612e1f565b9050600082013567ffffffffffffffff811115614f1957614f18612e3a565b5b614f2584828501614e43565b600083015250602082013567ffffffffffffffff811115614f4957614f48612e3a565b5b614f5584828501613277565b602083015250604082013567ffffffffffffffff811115614f7957614f78612e3a565b5b614f8584828501614c3a565b60408301525092915050565b600060608284031215614fa757614fa6612da9565b5b614fb16060612e1f565b9050600082013567ffffffffffffffff811115614fd157614fd0612e3a565b5b614fdd84828501614ed9565b6000830152506020614ff184828501612e66565b602083015250604082013567ffffffffffffffff81111561501557615014612e3a565b5b61502184828501613277565b60408301525092915050565b60006150393683614f91565b9050919050565b60006040828403121561505657615055612cd9565b5b600061506484828501614941565b91505092915050565b600067ffffffffffffffff82111561508857615087612dbf565b5b602082029050602081019050919050565b60006150ac6150a78461506d565b612e1f565b905080838252602082019050602084028301858111156150cf576150ce61305b565b5b835b8181101561511657803567ffffffffffffffff8111156150f4576150f3613051565b5b8086016151018982614c3a565b855260208501945050506020810190506150d1565b5050509392505050565b600082601f83011261513557615134613051565b5b8135615145848260208601615099565b91505092915050565b6000610180828403121561516557615164612da9565b5b615170610100612e1f565b9050600082013567ffffffffffffffff8111156151905761518f612e3a565b5b61519c84828501614c3a565b600083015250602082013567ffffffffffffffff8111156151c0576151bf612e3a565b5b6151cc84828501614a76565b602083015250604082013567ffffffffffffffff8111156151f0576151ef612e3a565b5b6151fc84828501614a76565b604083015250606061521084828501614b55565b60608301525060e0615224848285016149bd565b60808301525061012082013567ffffffffffffffff81111561524957615248612e3a565b5b61525584828501614c3a565b60a08301525061014082013567ffffffffffffffff81111561527a57615279612e3a565b5b61528684828501614c3a565b60c08301525061016082013567ffffffffffffffff8111156152ab576152aa612e3a565b5b6152b784828501615120565b60e08301525092915050565b60006152cf368361514e565b9050919050565b60006152e96152e4846131f5565b612e1f565b905082815260208101848484011115615305576153046131f0565b5b615310848285613617565b509392505050565b600082601f83011261532d5761532c613051565b5b815161533d8482602086016152d6565b91505092915050565b60006020828403121561535c5761535b612cd9565b5b600082015167ffffffffffffffff81111561537a57615379612cde565b5b61538684828501615318565b91505092915050565b600082825260208201905092915050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f726966794441436572745365637572697479506172616d733a20636f6e66697260208201527f6d6174696f6e5468726573686f6c64206d75737420626520677265617465722060408201527f7468616e206164766572736172795468726573686f6c64000000000000000000606082015250565b600061544860778361538f565b9150615453826153a0565b608082019050919050565b600060208201905081810360008301526154778161543b565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006154b882612ce3565b91506154c383612ce3565b9250828210156154d6576154d561547e565b5b828203905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061551b8261386e565b91506155268361386e565b925082615536576155356154e1565b5b828204905092915050565b600061554c8261386e565b91506155578361386e565b92508282101561556a5761556961547e565b5b828203905092915050565b60006155808261386e565b915061558b8361386e565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156155c4576155c361547e565b5b828202905092915050565b60006155da82612e3f565b91506155e583612e3f565b92508163ffffffff04831182151516156156025761560161547e565b5b828202905092915050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f726966794441436572745365637572697479506172616d733a2073656375726960208201527f747920617373756d7074696f6e7320617265206e6f74206d6574000000000000604082015250565b600061568f605a8361538f565b915061569a8261560d565b606082019050919050565b600060208201905081810360008301526156be81615682565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a20626c6f62486561646560208201527f727320616e6420626c6f62566572696669636174696f6e50726f6f6673206c6560408201527f6e677468206d69736d6174636800000000000000000000000000000000000000606082015250565b600061576d606d8361538f565b9150615778826156c5565b608082019050919050565b6000602082019050818103600083015261579c81615760565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600080fd5b600080fd5b600080fd5b60008235600160a0038336030381126157fd576157fc6157d2565b5b80830191505092915050565b60006020828403121561581f5761581e612cd9565b5b600061582d84828501612e66565b91505092915050565b61583f81612e3f565b82525050565b600060208201905061585a6000830184615836565b92915050565b60008151905061586f81613d9c565b92915050565b60006020828403121561588b5761588a612cd9565b5b600061589984828501615860565b91505092915050565b6000823560016060038336030381126158be576158bd6157d2565b5b80830191505092915050565b6000608082840312156158e0576158df612da9565b5b6158ea6080612e1f565b905060006158fa84828501613db3565b600083015250602082013567ffffffffffffffff81111561591e5761591d612e3a565b5b61592a84828501613277565b602083015250604082013567ffffffffffffffff81111561594e5761594d612e3a565b5b61595a84828501613277565b604083015250606061596e84828501612e66565b60608301525092915050565b6000606082840312156159905761598f612da9565b5b61599a6060612e1f565b9050600082013567ffffffffffffffff8111156159ba576159b9612e3a565b5b6159c6848285016158ca565b60008301525060206159da84828501613db3565b60208301525060406159ee84828501612e66565b60408301525092915050565b6000615a06368361597a565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a2062617463684d65746160208201527f6461746120646f6573206e6f74206d617463682073746f726564206d6574616460408201527f6174610000000000000000000000000000000000000000000000000000000000606082015250565b6000615ab560638361538f565b9150615ac082615a0d565b608082019050919050565b60006020820190508181036000830152615ae481615aa8565b9050919050565b60008083356001602003843603038112615b0857615b076157d2565b5b80840192508235915067ffffffffffffffff821115615b2a57615b296157d7565b5b602083019250600182023603831315615b4657615b456157dc565b5b509250929050565b600082356001608003833603038112615b6a57615b696157d2565b5b80830191505092915050565b600082356001608003833603038112615b9257615b916157d2565b5b80830191505092915050565b600067ffffffffffffffff821115615bb957615bb8612dbf565b5b602082029050602081019050919050565b600060808284031215615be057615bdf612da9565b5b615bea6080612e1f565b90506000615bfa84828501612d07565b6000830152506020615c0e84828501612d07565b6020830152506040615c2284828501612d07565b6040830152506060615c3684828501612e66565b60608301525092915050565b6000615c55615c5084615b9e565b612e1f565b90508083825260208201905060808402830185811115615c7857615c7761305b565b5b835b81811015615ca15780615c8d8882615bca565b845260208401935050608081019050615c7a565b5050509392505050565b600082601f830112615cc057615cbf613051565b5b8135615cd0848260208601615c42565b91505092915050565b600060808284031215615cef57615cee612da9565b5b615cf96060612e1f565b90506000615d09848285016149bd565b6000830152506040615d1d84828501612e66565b602083015250606082013567ffffffffffffffff811115615d4157615d40612e3a565b5b615d4d84828501615cab565b60408301525092915050565b6000615d653683615cd9565b9050919050565b6000819050919050565b615d87615d8282613d92565b615d6c565b82525050565b6000615d998284615d76565b60208201915081905092915050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a20696e636c7573696f6e60208201527f2070726f6f6620697320696e76616c6964000000000000000000000000000000604082015250565b6000615e2a60518361538f565b9150615e3582615da8565b606082019050919050565b60006020820190508181036000830152615e5981615e1d565b9050919050565b60008083356001602003843603038112615e7d57615e7c6157d2565b5b80840192508235915067ffffffffffffffff821115615e9f57615e9e6157d7565b5b602083019250608082023603831315615ebb57615eba6157dc565b5b509250929050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a2071756f72756d4e756d60208201527f62657220646f6573206e6f74206d617463680000000000000000000000000000604082015250565b6000615f4560528361538f565b9150615f5082615ec3565b606082019050919050565b60006020820190508181036000830152615f7481615f38565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a207468726573686f6c6460208201527f2070657263656e746167657320617265206e6f742076616c6964000000000000604082015250565b6000615ffd605a8361538f565b915061600882615f7b565b606082019050919050565b6000602082019050818103600083015261602c81615ff0565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a20636f6e6669726d617460208201527f696f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6560408201527f7400000000000000000000000000000000000000000000000000000000000000606082015250565b60006160db60618361538f565b91506160e682616033565b608082019050919050565b6000602082019050818103600083015261610a816160ce565b9050919050565b600061611c8261386e565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561614f5761614e61547e565b5b600182019050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667944414365727473466f7251756f72756d733a2072657175697265642060208201527f71756f72756d7320617265206e6f74206120737562736574206f66207468652060408201527f636f6e6669726d65642071756f72756d73000000000000000000000000000000606082015250565b600061620260718361538f565b915061620d8261615a565b608082019050919050565b60006020820190508181036000830152616231816161f5565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a2062617463684d6574616460208201527f61746120646f6573206e6f74206d617463682073746f726564206d657461646160408201527f7461000000000000000000000000000000000000000000000000000000000000606082015250565b60006162e060628361538f565b91506162eb82616238565b608082019050919050565b6000602082019050818103600083015261630f816162d3565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a20696e636c7573696f6e2060208201527f70726f6f6620697320696e76616c696400000000000000000000000000000000604082015250565b600061639860508361538f565b91506163a382616316565b606082019050919050565b600060208201905081810360008301526163c78161638b565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a2071756f72756d4e756d6260208201527f657220646f6573206e6f74206d61746368000000000000000000000000000000604082015250565b600061645060518361538f565b915061645b826163ce565b606082019050919050565b6000602082019050818103600083015261647f81616443565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a207468726573686f6c642060208201527f70657263656e746167657320617265206e6f742076616c696400000000000000604082015250565b600061650860598361538f565b915061651382616486565b606082019050919050565b60006020820190508181036000830152616537816164fb565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a20636f6e6669726d61746960208201527f6f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6574604082015250565b60006165c060608361538f565b91506165cb8261653e565b606082019050919050565b600060208201905081810360008301526165ef816165b3565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f72696679444143657274466f7251756f72756d733a207265717569726564207160208201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360408201527f6f6e6669726d65642071756f72756d7300000000000000000000000000000000606082015250565b600061669e60708361538f565b91506166a9826165f6565b608082019050919050565b600060208201905081810360008301526166cd81616691565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f726966794441436572745632466f7251756f72756d733a20696e636c7573696f60208201527f6e2070726f6f6620697320696e76616c69640000000000000000000000000000604082015250565b600061675660528361538f565b9150616761826166d4565b606082019050919050565b6000602082019050818103600083015261678581616749565b9050919050565b61679581613d92565b82525050565b60006080820190506167b0600083018761678c565b81810360208301526167c2818661364a565b90506167d16040830185615836565b81810360608301526167e38184613ac6565b905095945050505050565b600067ffffffffffffffff82111561680957616808612dbf565b5b602082029050602081019050919050565b60006bffffffffffffffffffffffff82169050919050565b61683b8161681a565b811461684657600080fd5b50565b60008151905061685881616832565b92915050565b600061687161686c846167ee565b612e1f565b905080838252602082019050602084028301858111156168945761689361305b565b5b835b818110156168bd57806168a98882616849565b845260208401935050602081019050616896565b5050509392505050565b600082601f8301126168dc576168db613051565b5b81516168ec84826020860161685e565b91505092915050565b60006040828403121561690b5761690a612da9565b5b6169156040612e1f565b9050600082015167ffffffffffffffff81111561693557616934612e3a565b5b616941848285016168c7565b600083015250602082015167ffffffffffffffff81111561696557616964612e3a565b5b616971848285016168c7565b60208301525092915050565b6000806040838503121561699457616993612cd9565b5b600083015167ffffffffffffffff8111156169b2576169b1612cde565b5b6169be858286016168f5565b92505060206169cf85828601615860565b9150509250929050565b60006169e48261681a565b91506169ef8361681a565b9250816bffffffffffffffffffffffff0483118215151615616a1457616a1361547e565b5b828202905092915050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f726966794441436572745632466f7251756f72756d733a20626c6f622071756f60208201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60408201527f6669726d65642071756f72756d73000000000000000000000000000000000000606082015250565b6000616ac7606e8361538f565b9150616ad282616a1f565b608082019050919050565b60006020820190508181036000830152616af681616aba565b9050919050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f726966794441436572745632466f7251756f72756d733a20726571756972656460208201527f2071756f72756d7320617265206e6f74206120737562736574206f662074686560408201527f20626c6f622071756f72756d7300000000000000000000000000000000000000606082015250565b6000616ba5606d8361538f565b9150616bb082616afd565b608082019050919050565b60006020820190508181036000830152616bd481616b98565b9050919050565b600081905092915050565b6000616bf1826135fb565b616bfb8185616bdb565b9350616c0b818560208601613617565b80840191505092915050565b60008160f81b9050919050565b6000616c2f82616c17565b9050919050565b616c47616c4282612ce3565b616c24565b82525050565b6000616c598285616be6565b9150616c658284616c36565b6001820191508190509392505050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b616caa81613d92565b82525050565b6000616cbc8383616ca1565b60208301905092915050565b6000602082019050919050565b6000616ce082616c75565b616cea8185616c80565b9350616cf583616c91565b8060005b83811015616d26578151616d0d8882616cb0565b9750616d1883616cc8565b925050600181019050616cf9565b5085935050505092915050565b6000608082019050616d4860008301876134df565b616d556020830186615836565b8181036040830152616d67818561364a565b90508181036060830152616d7b8184616cd5565b905095945050505050565b6000616d99616d9484614ba5565b612e1f565b90508083825260208201905060208402830185811115616dbc57616dbb61305b565b5b835b81811015616de55780616dd18882613cbf565b845260208401935050602081019050616dbe565b5050509392505050565b600082601f830112616e0457616e03613051565b5b8151616e14848260208601616d86565b91505092915050565b6000616e30616e2b8461506d565b612e1f565b90508083825260208201905060208402830185811115616e5357616e5261305b565b5b835b81811015616e9a57805167ffffffffffffffff811115616e7857616e77613051565b5b808601616e858982616def565b85526020850194505050602081019050616e55565b5050509392505050565b600082601f830112616eb957616eb8613051565b5b8151616ec9848260208601616e1d565b91505092915050565b600060808284031215616ee857616ee7612da9565b5b616ef26080612e1f565b9050600082015167ffffffffffffffff811115616f1257616f11612e3a565b5b616f1e84828501616def565b600083015250602082015167ffffffffffffffff811115616f4257616f41612e3a565b5b616f4e84828501616def565b602083015250604082015167ffffffffffffffff811115616f7257616f71612e3a565b5b616f7e84828501616def565b604083015250606082015167ffffffffffffffff811115616fa257616fa1612e3a565b5b616fae84828501616ea4565b60608301525092915050565b600060208284031215616fd057616fcf612cd9565b5b600082015167ffffffffffffffff811115616fee57616fed612cde565b5b616ffa84828501616ed2565b91505092915050565b600082825260208201905092915050565b600061701f826135fb565b6170298185617003565b9350617039818560208601613617565b61704281612dae565b840191505092915050565b60006080830160008301516170656000860182616ca1565b506020830151848203602086015261707d8282617014565b915050604083015184820360408601526170978282617014565b91505060608301516170ac6060860182612fd6565b508091505092915050565b600060208201905081810360008301526170d1818461704d565b905092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b60808201600082015161711b6000850182612fe5565b50602082015161712e6020850182612fe5565b5060408201516171416040850182612fe5565b5060608201516171546060850182612fd6565b50505050565b60006171668383617105565b60808301905092915050565b6000602082019050919050565b600061718a826170d9565b61719481856170e4565b935061719f836170f5565b8060005b838110156171d05781516171b7888261715a565b97506171c283617172565b9250506001810190506171a3565b5085935050505092915050565b60006080830160008301516171f56000860182613887565b5060208301516172086040860182612fd6565b5060408301518482036060860152617220828261717f565b9150508091505092915050565b6000602082019050818103600083015261724781846171dd565b905092915050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f7264657265644279746573417272617920697320746f6f2060208201527f6c6f6e6700000000000000000000000000000000000000000000000000000000604082015250565b60006172d160448361538f565b91506172dc8261724f565b606082019050919050565b60006020820190508181036000830152617300816172c4565b9050919050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f72646572656442797465734172726179206973206e6f742060208201527f6f72646572656400000000000000000000000000000000000000000000000000604082015250565b600061738960478361538f565b915061739482617307565b606082019050919050565b600060208201905081810360008301526173b88161737c565b9050919050565b600082825260208201905092915050565b60006173db82613793565b6173e581856173bf565b93506173f0836137af565b8060005b8381101561742157815161740888826137bf565b9750617413836137d7565b9250506001810190506173f4565b5085935050505092915050565b6000606082019050617443600083018661678c565b8181036020830152617455818561364a565b9050818103604083015261746981846173d0565b9050949350505050565b6040820160008201516174896000850182616ca1565b50602082015161749c6020850182612fd6565b50505050565b60006040820190506174b76000830184617473565b92915050565b60006174c8826133f7565b9050919050565b6174d8816174bd565b81146174e357600080fd5b50565b6000815190506174f5816174cf565b92915050565b60006020828403121561751157617510612cd9565b5b600061751f848285016174e6565b91505092915050565b7f456967656e444143657274566572696669636174696f6e5574696c732e5f766560008201527f7269667952656c61794b6579735365743a2072656c6179206b6579206973206e60208201527f6f74207365740000000000000000000000000000000000000000000000000000604082015250565b60006175aa60468361538f565b91506175b582617528565b606082019050919050565b600060208201905081810360008301526175d98161759d565b9050919050565b60008160e01b9050919050565b60006175f8826175e0565b9050919050565b61761061760b82612e3f565b6175ed565b82525050565b60006176228286615d76565b6020820191506176328285615d76565b60208201915061764282846175ff565b600482019150819050949350505050565b600061765e8261386e565b91506176698361386e565b925082617679576176786154e1565b5b828206905092915050565b7f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360008201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60208201527f7469706c65206f66203332000000000000000000000000000000000000000000604082015250565b6000617706604b8361538f565b915061771182617684565b606082019050919050565b60006020820190508181036000830152617735816176f9565b9050919050565b60006177478261386e565b91506177528361386e565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156177875761778661547e565b5b828201905092915050565b610160820160008201516177a96000850182613887565b5060208201516177bc60408501826139d5565b5060408201516177cf60c08501826139d5565b5060608201516177e3610140850182612fd6565b50505050565b60006101a0820190506177ff6000830186613c95565b8181036020830152617811818561364a565b90506178206040830184617792565b949350505050565b600060408201905061783d600083018561678c565b61784a602083018461678c565b939250505056fea264697066735822122079554fcf3cee3b67b040558b36b09d1ee8707efee72cd7ea180079d2ea69106764736f6c634300080c0033",
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

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDABatchMetadataStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDABatchMetadataStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDARelayRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDARelayRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDARelayRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDARelayRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
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

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifier.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifier.CallOpts, signedBatch)
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

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) OperatorStateRetriever(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "operatorStateRetriever")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) OperatorStateRetriever() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.OperatorStateRetriever(&_ContractEigenDACertVerifier.CallOpts)
}

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) OperatorStateRetriever() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.OperatorStateRetriever(&_ContractEigenDACertVerifier.CallOpts)
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

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.RegistryCoordinator(&_ContractEigenDACertVerifier.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.RegistryCoordinator(&_ContractEigenDACertVerifier.CallOpts)
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
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
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
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifier.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifier.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
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
