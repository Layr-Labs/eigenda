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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequiredV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6101406040523480156200001257600080fd5b5060405162004f8438038062004f84833981016040819052620000359162000265565b6001600160a01b0388811660805287811660a05286811660c05285811660e0528481166101009081529084166101205282516000805460208087015160ff94851661ffff1990931692909217939091169093029190911790558151620000a29160019190840190620000b1565b505050505050505050620003c7565b828054620000bf906200038a565b90600052602060002090601f016020900481019282620000e357600085556200012e565b82601f10620000fe57805160ff19168380011785556200012e565b828001600101855582156200012e579182015b828111156200012e57825182559160200191906001019062000111565b506200013c92915062000140565b5090565b5b808211156200013c576000815560010162000141565b6001600160a01b03811681146200016d57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620001b157620001b162000170565b604052919050565b805160ff81168114620001cb57600080fd5b919050565b600082601f830112620001e257600080fd5b81516001600160401b03811115620001fe57620001fe62000170565b602062000214601f8301601f1916820162000186565b82815285828487010111156200022957600080fd5b60005b83811015620002495785810183015182820184015282016200022c565b838111156200025b5760008385840101525b5095945050505050565b600080600080600080600080888a036101208112156200028457600080fd5b8951620002918162000157565b60208b0151909950620002a48162000157565b60408b0151909850620002b78162000157565b60608b0151909750620002ca8162000157565b60808b0151909650620002dd8162000157565b60a08b0151909550620002f08162000157565b9350604060bf19820112156200030557600080fd5b50604080519081016001600160401b0380821183831017156200032c576200032c62000170565b816040526200033e60c08d01620001b9565b83526200034e60e08d01620001b9565b60208401526101008c0151929450808311156200036a57600080fd5b50506200037a8b828c01620001d0565b9150509295985092959890939650565b600181811c908216806200039f57607f821691505b60208210811415620003c157634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c05160e0516101005161012051614ac5620004bf6000396000818161029b015281816108330152610d16015260008181610235015281816108120152610cf50152600081816102c201528181610684015281816107f101526109b301526000818161039201528181610663015281816107d001526109920152600081816102740152818161060a015261093e0152600081816103d901528181610416015281816104aa01528181610579015281816105e901528181610642015281816107af0152818161091d0152818161097101528181610aa801528181610bc301528181610c380152610caf0152614ac56000f3fe608060405234801561001057600080fd5b506004361061014d5760003560e01c80637d644cad116100c3578063e15234ff1161007c578063e15234ff14610342578063ed0450ae1461034a578063ee6c3bcf1461037a578063efd4532b1461038d578063f25de3f8146103b4578063f8c66814146103d457600080fd5b80637d644cad146102e4578063813c2eb0146102f75780638687feae1461030a578063b74d78711461031f578063bafa910714610327578063ccb7cd0d1461032f57600080fd5b8063415ef61411610115578063415ef6141461020a578063421c02221461021d5780634ca22c3f14610230578063640f65d91461026f5780636d14a9871461029657806372276443146102bd57600080fd5b8063048886d2146101525780631429c7c21461017a578063143eb4d91461019f5780632ecfe72b146101b457806331a3479a146101f7575b600080fd5b610165610160366004612fb7565b6103fb565b60405190151581526020015b60405180910390f35b61018d610188366004612fb7565b61048f565b60405160ff9091168152602001610171565b6101b26101ad366004613136565b61051e565b005b6101c76101c23660046131be565b61053f565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff1690820152606001610171565b6101b261020536600461321d565b6105e4565b61016561021836600461331d565b61063a565b6101b261022b3660046133c6565b6107aa565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610171565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b6101b26102f2366004613429565b610918565b6101b261030536600461331d565b61096c565b610312610aa4565b60405161017191906134f3565b610312610b31565b610312610bbf565b6101b261033d366004613506565b610c1f565b610312610c34565b6000546103609060ff8082169161010090041682565b6040805160ff938416815292909116602083015201610171565b61018d610388366004612fb7565b610c94565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b6103c76103c2366004613531565b610ce6565b6040516101719190613754565b6102577f000000000000000000000000000000000000000000000000000000000000000081565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa158015610465573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104899190613767565b92915050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa1580156104fa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104899190613789565b60008061052b8484610d4b565b915091506105398282610e81565b50505050565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa1580156105c0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061048991906137a6565b6105397f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000086868686610635610c34565b611110565b6000806107757f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006106b2368b90038b018b613845565b6106bb8a6139a8565b6106c48a613c36565b6040805180820190915260005460ff8082168352610100909104166020820152600180546106f190613d5e565b80601f016020809104026020016040519081016040528092919081815260200182805461071d90613d5e565b801561076a5780601f1061073f5761010080835404028352916020019161076a565b820191906000526020600020905b81548152906001019060200180831161074d57829003601f168201915b50505050508b611bd9565b509050600081600681111561078c5761078c613d93565b141561079c5760019150506107a2565b60009150505b949350505050565b6109147f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061085b88613da9565b610864886139a8565b6040805180820190915260005460ff80821683526101009091041660208201526001805461089190613d5e565b80601f01602080910402602001604051908101604052809291908181526020018280546108bd90613d5e565b801561090a5780601f106108df5761010080835404028352916020019161090a565b820191906000526020600020905b8154815290600101906020018083116108ed57829003601f168201915b5050505050611d79565b5050565b6109147f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610967610c34565b611dad565b6105397f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006109e136899003890189613845565b6109ea886139a8565b6109f388613c36565b6040805180820190915260005460ff808216835261010090910416602082015260018054610a2090613d5e565b80601f0160208091040260200160405190810160405280929190818152602001828054610a4c90613d5e565b8015610a995780601f10610a6e57610100808354040283529160200191610a99565b820191906000526020600020905b815481529060010190602001808311610a7c57829003601f168201915b5050505050896124b2565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b04573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610b2c9190810190613ea0565b905090565b60018054610b3e90613d5e565b80601f0160208091040260200160405190810160405280929190818152602001828054610b6a90613d5e565b8015610bb75780601f10610b8c57610100808354040283529160200191610bb7565b820191906000526020600020905b815481529060010190602001808311610b9a57829003601f168201915b505050505081565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b04573d6000803e3d6000fd5b60008061052b610c2e8561053f565b84610d4b565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b04573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf906024016104dd565b610cee612efc565b6000610d437f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610d3e86613da9565b6124d4565b509392505050565b60006060826020015160ff16836000015160ff1611610da857600283600001518460200151604051602001610d9092919060ff92831681529116602082015260400190565b60405160208183030381529060405291509150610e7a565b60208301518351600091610dbb91613f23565b60ff1690506000856020015163ffffffff16866040015160ff1683620f4240610de49190613f5c565b610dee9190613f5c565b610dfa90612710613f70565b610e049190613f87565b8651909150600090610e1890612710613fa6565b63ffffffff169050808210610e455760006040518060200160405280600081525094509450505050610e7a565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050505b9250929050565b6000826006811115610e9557610e95613d93565b1415610e9f575050565b6001826006811115610eb357610eb3613d93565b1415610f0957600080600083806020019051810190610ed29190613fd2565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b6002826006811115610f1d57610f1d613d93565b1415610f655760008082806020019051810190610f3a9190614000565b604051631b00235d60e01b815260ff8084166004830152821660248201529193509150604401610f00565b6003826006811115610f7957610f79613d93565b1415610fcd57600080600083806020019051810190610f989190613fd2565b6040516001626dc9ad60e11b031981526004810184905260248101839052604481018290529295509093509150606401610f00565b6004826006811115610fe157610fe1613d93565b14156110265760008082806020019051810190610ffe919061402f565b604051634a47030360e11b815260048101839052602481018290529193509150604401610f00565b600582600681111561103a5761103a613d93565b141561107f5760008082806020019051810190611057919061402f565b60405163114b085b60e21b815260048101839052602481018290529193509150604401610f00565b600682600681111561109357611093613d93565b14156110d3576000818060200190518101906110af9190614053565b6040516309efaa0b60e41b815263ffffffff82166004820152909150602401610f00565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b6044820152606401610f00565b8382146111af5760405162461bcd60e51b815260206004820152606d6024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a20626c6f62486561646560648201527f727320616e6420626c6f62566572696669636174696f6e50726f6f6673206c6560848201526c0dccee8d040dad2e6dac2e8c6d609b1b60a482015260c401610f00565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156111ef573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526112179190810190613ea0565b905060005b85811015611bce57876001600160a01b031663eccbbfc986868481811061124557611245614070565b90506020028101906112579190614086565b6112659060208101906140a6565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156112a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112cb91906140c3565b61130e8686848181106112e0576112e0614070565b90506020028101906112f29190614086565b6113009060408101906140dc565b611309906140f2565b6126e9565b146113a15760405162461bcd60e51b81526020600482015260636024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a2062617463684d65746160648201527f6461746120646f6573206e6f74206d617463682073746f726564206d6574616460848201526261746160e81b60a482015260c401610f00565b6114e78585838181106113b6576113b6614070565b90506020028101906113c89190614086565b6113d69060608101906141c9565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525089925088915085905081811061141f5761141f614070565b90506020028101906114319190614086565b61143f9060408101906140dc565b611449908061420f565b3561147f8a8a8681811061145f5761145f614070565b9050602002810190611471919061420f565b61147a90614225565b61275a565b60405160200161149191815260200190565b604051602081830303815290604052805190602001208888868181106114b9576114b9614070565b90506020028101906114cb9190614086565b6114dc9060408101906020016140a6565b63ffffffff1661278a565b6115615760405162461bcd60e51b81526020600482015260516024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a20696e636c7573696f6e606482015270081c1c9bdbd9881a5cc81a5b9d985b1a59607a1b608482015260a401610f00565b6000805b88888481811061157757611577614070565b9050602002810190611589919061420f565b611597906060810190614344565b9050811015611b09578888848181106115b2576115b2614070565b90506020028101906115c4919061420f565b6115d2906060810190614344565b828181106115e2576115e2614070565b6115f89260206080909202019081019150612fb7565b60ff1687878581811061160d5761160d614070565b905060200281019061161f9190614086565b61162d9060408101906140dc565b611637908061420f565b6116459060208101906141c9565b89898781811061165757611657614070565b90506020028101906116699190614086565b6116779060808101906141c9565b8581811061168757611687614070565b919091013560f81c90508181106116a0576116a0614070565b9050013560f81c60f81b60f81c60ff161461172c5760405162461bcd60e51b81526020600482015260526024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a2071756f72756d4e756d6064820152710c4cae440c8decae640dcdee840dac2e8c6d60731b608482015260a401610f00565b88888481811061173e5761173e614070565b9050602002810190611750919061420f565b61175e906060810190614344565b8281811061176e5761176e614070565b90506080020160200160208101906117869190612fb7565b60ff1689898581811061179b5761179b614070565b90506020028101906117ad919061420f565b6117bb906060810190614344565b838181106117cb576117cb614070565b90506080020160400160208101906117e39190612fb7565b60ff161161186d5760405162461bcd60e51b815260206004820152605a6024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a207468726573686f6c6460648201527f2070657263656e746167657320617265206e6f742076616c6964000000000000608482015260a401610f00565b8389898581811061188057611880614070565b9050602002810190611892919061420f565b6118a0906060810190614344565b838181106118b0576118b0614070565b6118c69260206080909202019081019150612fb7565b60ff16815181106118d9576118d9614070565b016020015160f81c8989858181106118f3576118f3614070565b9050602002810190611905919061420f565b611913906060810190614344565b8381811061192357611923614070565b905060800201604001602081019061193b9190612fb7565b60ff16101561195c5760405162461bcd60e51b8152600401610f009061438d565b88888481811061196e5761196e614070565b9050602002810190611980919061420f565b61198e906060810190614344565b8281811061199e5761199e614070565b90506080020160400160208101906119b69190612fb7565b60ff168787858181106119cb576119cb614070565b90506020028101906119dd9190614086565b6119eb9060408101906140dc565b6119f5908061420f565b611a039060408101906141c9565b898987818110611a1557611a15614070565b9050602002810190611a279190614086565b611a359060808101906141c9565b85818110611a4557611a45614070565b919091013560f81c9050818110611a5e57611a5e614070565b9050013560f81c60f81b60f81c60ff161015611a8c5760405162461bcd60e51b8152600401610f009061438d565b611af5828a8a86818110611aa257611aa2614070565b9050602002810190611ab4919061420f565b611ac2906060810190614344565b84818110611ad257611ad2614070565b611ae89260206080909202019081019150612fb7565b600160ff919091161b1790565b915080611b0181614408565b915050611565565b50611b1d611b16856127a2565b8281161490565b611bbd5760405162461bcd60e51b81526020600482015260716024820152600080516020614a7083398151915260448201527f7269667944414365727473466f7251756f72756d733a2072657175697265642060648201527f71756f72756d7320617265206e6f74206120737562736574206f662074686520608482015270636f6e6669726d65642071756f72756d7360781b60a482015260c401610f00565b50611bc781614408565b905061121c565b505050505050505050565b60006060611be7888861292f565b90925090506000826006811115611c0057611c00613d93565b14611c0a57611d6b565b611c1c89886000015160400151612a05565b90925090506000826006811115611c3557611c35613d93565b14611c3f57611d6b565b86515151604051632ecfe72b60e01b815261ffff9091166004820152611cba906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa158015611c90573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cb491906137a6565b86610d4b565b90925090506000826006811115611cd357611cd3613d93565b14611cdd57611d6b565b6000611cf98b611cec8b612b3d565b868c602001518b8b612b68565b919450925090506000836006811115611d1457611d14613d93565b14611d1f5750611d6b565b87515160200151600090611d339083612cc9565b919550935090506000846006811115611d4e57611d4e613d93565b14611d5a575050611d6b565b611d648682612d2f565b9350935050505b995099975050505050505050565b600080611d878888886124d4565b91509150611da08b8b8b896000015189878a8a896124b2565b5050505050505050505050565b6001600160a01b03841663eccbbfc9611dc960208501856140a6565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611e0b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e2f91906140c3565b611e3f61130060408501856140dc565b14611ed15760405162461bcd60e51b81526020600482015260626024820152600080516020614a7083398151915260448201527f72696679444143657274466f7251756f72756d733a2062617463684d6574616460648201527f61746120646f6573206e6f74206d617463682073746f726564206d6574616461608482015261746160f01b60a482015260c401610f00565b611f75611ee160608401846141c9565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611f239250505060408501856140dc565b611f2d908061420f565b35611f3a61147a87614225565b604051602001611f4c91815260200190565b604051602081830303815290604052805190602001208560200160208101906114dc91906140a6565b611fee5760405162461bcd60e51b81526020600482015260506024820152600080516020614a7083398151915260448201527f72696679444143657274466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a401610f00565b6000805b611fff6060860186614344565b90508110156123fe576120156060860186614344565b8281811061202557612025614070565b61203b9260206080909202019081019150612fb7565b60ff1661204b60408601866140dc565b612055908061420f565b6120639060208101906141c9565b61207060808801886141c9565b8581811061208057612080614070565b919091013560f81c905081811061209957612099614070565b9050013560f81c60f81b60f81c60ff16146121245760405162461bcd60e51b81526020600482015260516024820152600080516020614a7083398151915260448201527f72696679444143657274466f7251756f72756d733a2071756f72756d4e756d626064820152700cae440c8decae640dcdee840dac2e8c6d607b1b608482015260a401610f00565b6121316060860186614344565b8281811061214157612141614070565b90506080020160200160208101906121599190612fb7565b60ff166121696060870187614344565b8381811061217957612179614070565b90506080020160400160208101906121919190612fb7565b60ff161161221b5760405162461bcd60e51b81526020600482015260596024820152600080516020614a7083398151915260448201527f72696679444143657274466f7251756f72756d733a207468726573686f6c642060648201527f70657263656e746167657320617265206e6f742076616c696400000000000000608482015260a401610f00565b6001600160a01b038716631429c7c26122376060880188614344565b8481811061224757612247614070565b61225d9260206080909202019081019150612fb7565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa15801561229b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122bf9190613789565b60ff166122cf6060870187614344565b838181106122df576122df614070565b90506080020160400160208101906122f79190612fb7565b60ff1610156123185760405162461bcd60e51b8152600401610f0090614423565b6123256060860186614344565b8281811061233557612335614070565b905060800201604001602081019061234d9190612fb7565b60ff1661235d60408601866140dc565b612367908061420f565b6123759060408101906141c9565b61238260808801886141c9565b8581811061239257612392614070565b919091013560f81c90508181106123ab576123ab614070565b9050013560f81c60f81b60f81c60ff1610156123d95760405162461bcd60e51b8152600401610f0090614423565b6123ea82611ac26060880188614344565b9150806123f681614408565b915050611ff2565b5061240b611b16836127a2565b6124aa5760405162461bcd60e51b81526020600482015260706024820152600080516020614a7083398151915260448201527f72696679444143657274466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c401610f00565b505050505050565b6000806124c68b8b8b8b8b8b8b8b8b611bd9565b91509150611da08282610e81565b6124dc612efc565b60606000836020015160000151516001600160401b0381111561250157612501612fdb565b60405190808252806020026020018201604052801561252a578160200160208202803683370190505b50905060005b602085015151518110156125a75761257a856020015160000151828151811061255b5761255b614070565b6020026020010151805160009081526020918201519091526040902090565b82828151811061258c5761258c614070565b60209081029190910101526125a081614408565b9050612530565b5060005b84602001516080015151811015612612578285602001516080015182815181106125d7576125d7614070565b60200260200101516040516020016125f0929190614495565b60405160208183030381529060405292508061260b90614408565b90506125ab565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f749161264d918a9190889088906004016144c7565b600060405180830381865afa15801561266a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052612692919081019061461a565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b6000610489826000015160405160200161270391906146f2565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b60008160405160200161276d9190614752565b604051602081830303815290604052805190602001209050919050565b600083612798868585612d7f565b1495945050505050565b60006101008251111561282b5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610f00565b815161283957506000919050565b6000808360008151811061284f5761284f614070565b0160200151600160f89190911c81901b92505b84518110156129265784818151811061287d5761287d614070565b0160200151600160f89190911c1b91508282116129125760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610f00565b9181179161291f81614408565b9050612862565b50909392505050565b6000606060006129428460000151612e82565b905060008160405160200161295991815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000612996876040015183858a6020015163ffffffff1661278a565b905080156129bd576000604051806020016040528060008152509550955050505050610e7a565b6020808801516040805163ffffffff90921692820192909252908101849052606081018390526001906080016040516020818303038152906040529550955050505050610e7a565b6000606060005b8351811015612b225760006001600160a01b0316856001600160a01b031663b5a872da868481518110612a4157612a41614070565b60200260200101516040518263ffffffff1660e01b8152600401612a71919063ffffffff91909116815260200190565b602060405180830381865afa158015612a8e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ab291906147f7565b6001600160a01b03161415612b12576006848281518110612ad557612ad5614070565b6020026020010151604051602001612af9919063ffffffff91909116815260200190565b6040516020818303038152906040529250925050610e7a565b612b1b81614408565b9050612a0c565b50506040805160208101909152600080825291509250929050565b60008160405160200161276d91908151815260209182015163ffffffff169181019190915260400190565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b8152600401612ba19493929190614820565b600060405180830381865afa158015612bbe573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052612be691908101906148d3565b5090506000915060005b8851811015612ca757856000015160ff1682602001518281518110612c1757612c17614070565b6020026020010151612c29919061496f565b6001600160601b0316606483600001518381518110612c4a57612c4a614070565b60200260200101516001600160601b0316612c659190613f87565b10612c9557612c92838a8381518110612c8057612c80614070565b0160200151600160f89190911c1b1790565b92505b80612c9f81614408565b915050612bf0565b5050604080516020810190915260008082529350915096509650969350505050565b600060606000612cd8856127a2565b9050838116811415612cfd576040805160208101909152600080825293509150612d28565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b600060606000612d3e856127a2565b9050838116811415612d63575050604080516020810190915260008082529150610e7a565b6040805160208101839052908101859052600590606001612af9565b600060208451612d8f9190614995565b15612e165760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610f00565b8260205b85518111612e7957612e2d600285614995565b612e4e57816000528086015160205260406000209150600284049350612e67565b8086015160005281602052604060002091506002840493505b612e726020826149a9565b9050612e1a565b50949350505050565b6000612e918260000151612eaa565b602080840151604080860151905161276d9493016149c1565b6000816000015182602001518360400151604051602001612ecd939291906149f6565b60408051601f19818403018152828252805160209182012060608087015192850191909152918301520161276d565b604051806101000160405280606081526020016060815260200160608152602001612f25612f62565b8152602001612f47604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060400160405280612f75612f87565b8152602001612f82612f87565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff81168114612fb457600080fd5b50565b600060208284031215612fc957600080fd5b8135612fd481612fa5565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561301357613013612fdb565b60405290565b604051606081016001600160401b038111828210171561301357613013612fdb565b604051608081016001600160401b038111828210171561301357613013612fdb565b60405161010081016001600160401b038111828210171561301357613013612fdb565b60405160a081016001600160401b038111828210171561301357613013612fdb565b604051601f8201601f191681016001600160401b03811182821017156130ca576130ca612fdb565b604052919050565b63ffffffff81168114612fb457600080fd5b80356130ef816130d2565b919050565b60006040828403121561310657600080fd5b61310e612ff1565b9050813561311b81612fa5565b8152602082013561312b81612fa5565b602082015292915050565b60008082840360a081121561314a57600080fd5b606081121561315857600080fd5b50613161613019565b833561316c816130d2565b8152602084013561317c816130d2565b6020820152604084013561318f81612fa5565b604082015291506131a384606085016130f4565b90509250929050565b803561ffff811681146130ef57600080fd5b6000602082840312156131d057600080fd5b612fd4826131ac565b60008083601f8401126131eb57600080fd5b5081356001600160401b0381111561320257600080fd5b6020830191508360208260051b8501011115610e7a57600080fd5b6000806000806040858703121561323357600080fd5b84356001600160401b038082111561324a57600080fd5b613256888389016131d9565b9096509450602087013591508082111561326f57600080fd5b5061327c878288016131d9565b95989497509550505050565b60006060828403121561329a57600080fd5b50919050565b60006001600160401b038211156132b9576132b9612fdb565b50601f01601f191660200190565b600082601f8301126132d857600080fd5b81356132eb6132e6826132a0565b6130a2565b81815284602083860101111561330057600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a081121561333457600080fd5b604081121561334257600080fd5b5084935060408501356001600160401b038082111561336057600080fd5b61336c88838901613288565b9450606087013591508082111561338257600080fd5b90860190610180828903121561339757600080fd5b909250608086013590808211156133ad57600080fd5b506133ba878288016132c7565b91505092959194509250565b600080604083850312156133d957600080fd5b82356001600160401b03808211156133f057600080fd5b6133fc86838701613288565b9350602085013591508082111561341257600080fd5b5061341f85828601613288565b9150509250929050565b6000806040838503121561343c57600080fd5b82356001600160401b038082111561345357600080fd5b908401906080828703121561346757600080fd5b9092506020840135908082111561347d57600080fd5b50830160a0818603121561349057600080fd5b809150509250929050565b60005b838110156134b657818101518382015260200161349e565b838111156105395750506000910152565b600081518084526134df81602086016020860161349b565b601f01601f19169290920160200192915050565b602081526000612fd460208301846134c7565b6000806060838503121561351957600080fd5b613522836131ac565b91506131a384602085016130f4565b60006020828403121561354357600080fd5b81356001600160401b0381111561355957600080fd5b6107a284828501613288565b600081518084526020808501945080840160005b8381101561359b57815163ffffffff1687529582019590820190600101613579565b509495945050505050565b600081518084526020808501945080840160005b8381101561359b576135d787835180518252602090810151910152565b60409690960195908201906001016135ba565b8060005b60028110156105395781518452602093840193909101906001016135ee565b6136188282516135ea565b602081015161362a60408401826135ea565b505050565b600082825180855260208086019550808260051b84010181860160005b8481101561367a57601f19868403018952613668838351613565565b9884019892509083019060010161364c565b5090979650505050505050565b6000610180825181855261369d82860182613565565b915050602083015184820360208601526136b782826135a6565b915050604083015184820360408601526136d182826135a6565b91505060608301516136e6606086018261360d565b506080830151805160e08601526020015161010085015260a08301518482036101208601526137158282613565565b91505060c08301518482036101408601526137308282613565565b91505060e083015184820361016086015261374b828261362f565b95945050505050565b602081526000612fd46020830184613687565b60006020828403121561377957600080fd5b81518015158114612fd457600080fd5b60006020828403121561379b57600080fd5b8151612fd481612fa5565b6000606082840312156137b857600080fd5b604051606081018181106001600160401b03821117156137da576137da612fdb565b60405282516137e8816130d2565b815260208301516137f8816130d2565b6020820152604083015161380b81612fa5565b60408201529392505050565b60006040828403121561382957600080fd5b613831612ff1565b905081358152602082013561312b816130d2565b60006040828403121561385757600080fd5b612fd48383613817565b60006040828403121561387357600080fd5b61387b612ff1565b9050813581526020820135602082015292915050565b600082601f8301126138a257600080fd5b6138aa612ff1565b8060408401858111156138bc57600080fd5b845b818110156138d65780358452602093840193016138be565b509095945050505050565b6000608082840312156138f357600080fd5b6138fb612ff1565b90506139078383613891565b815261312b8360408401613891565b60006001600160401b0382111561392f5761392f612fdb565b5060051b60200190565b600082601f83011261394a57600080fd5b8135602061395a6132e683613916565b82815260059290921b8401810191818101908684111561397957600080fd5b8286015b8481101561399d578035613990816130d2565b835291830191830161397d565b509695505050505050565b6000606082360312156139ba57600080fd5b6139c2613019565b82356001600160401b03808211156139d957600080fd5b8185019150606082360312156139ee57600080fd5b6139f6613019565b823582811115613a0557600080fd5b8301368190036101c0811215613a1a57600080fd5b613a2261303b565b613a2b836131ac565b815260208084013586811115613a4057600080fd5b613a4c368287016132c7565b83830152506040610160603f1985011215613a6657600080fd5b613a6e61303b565b9350613a7c36828701613861565b8452613a8b36608087016138e1565b82850152613a9d3661010087016138e1565b81850152610180850135613ab0816130d2565b8060608601525083818401526101a0850135606084015282865281880135945086851115613add57600080fd5b613ae936868a016132c7565b8287015280880135945086851115613b0057600080fd5b613b0c36868a01613939565b81870152858952613b1e828c016130e4565b828a0152808b0135975086881115613b3557600080fd5b613b4136898d016132c7565b90890152509598975050505050505050565b600082601f830112613b6457600080fd5b81356020613b746132e683613916565b82815260069290921b84018101918181019086841115613b9357600080fd5b8286015b8481101561399d57613ba98882613861565b835291830191604001613b97565b600082601f830112613bc857600080fd5b81356020613bd86132e683613916565b82815260059290921b84018101918181019086841115613bf757600080fd5b8286015b8481101561399d5780356001600160401b03811115613c1a5760008081fd5b613c288986838b0101613939565b845250918301918301613bfb565b60006101808236031215613c4957600080fd5b613c5161305d565b82356001600160401b0380821115613c6857600080fd5b613c7436838701613939565b83526020850135915080821115613c8a57600080fd5b613c9636838701613b53565b60208401526040850135915080821115613caf57600080fd5b613cbb36838701613b53565b6040840152613ccd36606087016138e1565b6060840152613cdf3660e08701613861565b6080840152610120850135915080821115613cf957600080fd5b613d0536838701613939565b60a0840152610140850135915080821115613d1f57600080fd5b613d2b36838701613939565b60c0840152610160850135915080821115613d4557600080fd5b50613d5236828601613bb7565b60e08301525092915050565b600181811c90821680613d7257607f821691505b6020821081141561329a57634e487b7160e01b600052602260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b600060608236031215613dbb57600080fd5b613dc3612ff1565b613dcd3684613817565b815260408301356001600160401b0380821115613de957600080fd5b81850191506101208236031215613dff57600080fd5b613e07613080565b823582811115613e1657600080fd5b613e2236828601613b53565b825250602083013582811115613e3757600080fd5b613e4336828601613b53565b602083015250613e563660408501613861565b6040820152613e6836608085016138e1565b606082015261010083013582811115613e8057600080fd5b613e8c36828601613939565b608083015250602084015250909392505050565b600060208284031215613eb257600080fd5b81516001600160401b03811115613ec857600080fd5b8201601f81018413613ed957600080fd5b8051613ee76132e6826132a0565b818152856020838501011115613efc57600080fd5b61374b82602083016020860161349b565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff841680821015613f3d57613f3d613f0d565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b600082613f6b57613f6b613f46565b500490565b600082821015613f8257613f82613f0d565b500390565b6000816000190483118215151615613fa157613fa1613f0d565b500290565b600063ffffffff80831681851681830481118215151615613fc957613fc9613f0d565b02949350505050565b600080600060608486031215613fe757600080fd5b8351925060208401519150604084015190509250925092565b6000806040838503121561401357600080fd5b825161401e81612fa5565b602084015190925061349081612fa5565b6000806040838503121561404257600080fd5b505080516020909101519092909150565b60006020828403121561406557600080fd5b8151612fd4816130d2565b634e487b7160e01b600052603260045260246000fd5b60008235609e1983360301811261409c57600080fd5b9190910192915050565b6000602082840312156140b857600080fd5b8135612fd4816130d2565b6000602082840312156140d557600080fd5b5051919050565b60008235605e1983360301811261409c57600080fd5b60006060823603121561410457600080fd5b61410c613019565b82356001600160401b038082111561412357600080fd5b81850191506080823603121561413857600080fd5b61414061303b565b8235815260208301358281111561415657600080fd5b614162368286016132c7565b60208301525060408301358281111561417a57600080fd5b614186368286016132c7565b6040830152506060830135925061419c836130d2565b826060820152808452505050602083013560208201526141be604084016130e4565b604082015292915050565b6000808335601e198436030181126141e057600080fd5b8301803591506001600160401b038211156141fa57600080fd5b602001915036819003821315610e7a57600080fd5b60008235607e1983360301811261409c57600080fd5b6000608080833603121561423857600080fd5b614240613019565b61424a3685613861565b815260408085013561425b816130d2565b6020818185015260609150818701356001600160401b0381111561427e57600080fd5b870136601f82011261428f57600080fd5b803561429d6132e682613916565b81815260079190911b820183019083810190368311156142bc57600080fd5b928401925b82841015614330578884360312156142d95760008081fd5b6142e161303b565b84356142ec81612fa5565b8152848601356142fb81612fa5565b818701528488013561430c81612fa5565b818901528487013561431d816130d2565b81880152825292880192908401906142c1565b958701959095525093979650505050505050565b6000808335601e1984360301811261435b57600080fd5b8301803591506001600160401b0382111561437557600080fd5b6020019150600781901b3603821315610e7a57600080fd5b6020808252606190820152600080516020614a7083398151915260408201527f7269667944414365727473466f7251756f72756d733a20636f6e6669726d617460608201527f696f6e5468726573686f6c6450657263656e74616765206973206e6f74206d656080820152601d60fa1b60a082015260c00190565b600060001982141561441c5761441c613f0d565b5060010190565b60208082526060908201819052600080516020614a7083398151915260408301527f72696679444143657274466f7251756f72756d733a20636f6e6669726d617469908201527f6f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6574608082015260a00190565b600083516144a781846020880161349b565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526144f560808401866134c7565b838103606085015284518082528286019183019060005b818110156145285783518352928401929184019160010161450c565b50909998505050505050505050565b600082601f83011261454857600080fd5b815160206145586132e683613916565b82815260059290921b8401810191818101908684111561457757600080fd5b8286015b8481101561399d57805161458e816130d2565b835291830191830161457b565b600082601f8301126145ac57600080fd5b815160206145bc6132e683613916565b82815260059290921b840181019181810190868411156145db57600080fd5b8286015b8481101561399d5780516001600160401b038111156145fe5760008081fd5b61460c8986838b0101614537565b8452509183019183016145df565b60006020828403121561462c57600080fd5b81516001600160401b038082111561464357600080fd5b908301906080828603121561465757600080fd5b61465f61303b565b82518281111561466e57600080fd5b61467a87828601614537565b82525060208301518281111561468f57600080fd5b61469b87828601614537565b6020830152506040830151828111156146b357600080fd5b6146bf87828601614537565b6040830152506060830151828111156146d757600080fd5b6146e38782860161459b565b60608301525095945050505050565b6020815281516020820152600060208301516080604084015261471860a08401826134c7565b90506040840151601f1984830301606085015261473582826134c7565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b808410156147e8578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906147a5565b509a9950505050505050505050565b60006020828403121561480957600080fd5b81516001600160a01b0381168114612fd457600080fd5b84815260806020820152600061483960808301866134c7565b63ffffffff8516604084015282810360608401526148578185613687565b979650505050505050565b600082601f83011261487357600080fd5b815160206148836132e683613916565b82815260059290921b840181019181810190868411156148a257600080fd5b8286015b8481101561399d5780516001600160601b03811681146148c65760008081fd5b83529183019183016148a6565b600080604083850312156148e657600080fd5b82516001600160401b03808211156148fd57600080fd5b908401906040828703121561491157600080fd5b614919612ff1565b82518281111561492857600080fd5b61493488828601614862565b82525060208301518281111561494957600080fd5b61495588828601614862565b602083015250809450505050602083015190509250929050565b60006001600160601b0380831681851681830481118215151615613fc957613fc9613f0d565b6000826149a4576149a4613f46565b500690565b600082198211156149bc576149bc613f0d565b500190565b8381526060602082015260006149da60608301856134c7565b82810360408401526149ec8185613565565b9695505050505050565b60006101a061ffff86168352806020840152614a14818401866134c7565b8451805160408601526020015160608501529150614a2f9050565b6020830151614a41608084018261360d565b506040830151614a5561010084018261360d565b5063ffffffff60608401511661018083015294935050505056fe456967656e444143657274566572696669636174696f6e56314c69622e5f7665a2646970667358221220da3f44793938d7fb729024db04aa4a2c45d62ba1c0aaa5e3de1b0e0d9ba2634764736f6c634300080c0033",
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
