// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDABlobVerifier

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
	RelayKeys  []uint32
}

// BlobCommitment is an auto generated low-level Go binding around an user-defined struct.
type BlobCommitment struct {
	Commitment       BN254G1Point
	LengthCommitment BN254G2Point
	LengthProof      BN254G2Point
	DataLength       uint32
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
	Salt              uint32
	PaymentHeaderHash [32]byte
}

// BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// BlobVerificationProofV2 is an auto generated low-level Go binding around an user-defined struct.
type BlobVerificationProofV2 struct {
	BlobCertificate BlobCertificate
	BlobIndex       uint32
	InclusionProof  []byte
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

// ContractEigenDABlobVerifierMetaData contains all meta data concerning the ContractEigenDABlobVerifier contract.
var ContractEigenDABlobVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"salt\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"salt\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6101406040523480156200001257600080fd5b506040516200446f3803806200446f83398101604081905262000035916200007e565b6001600160a01b0395861660805293851660a05291841660c052831660e052821661010052166101205262000112565b6001600160a01b03811681146200007b57600080fd5b50565b60008060008060008060c087890312156200009857600080fd5b8651620000a58162000065565b6020880151909650620000b88162000065565b6040880151909550620000cb8162000065565b6060880151909450620000de8162000065565b6080880151909350620000f18162000065565b60a0880151909250620001048162000065565b809150509295509295509295565b60805160a05160c05160e0516101005161012051614273620001fc6000396000818161025f0152818161062a01526109db0152600081816101e60152818161060901526109ba015260008181610286015281816105e80152610a4f015260008181610347015281816105c70152610a2e0152600081816102380152818161077201526107d001526000818161038e015281816103de0152818161048801528181610536015281816105a6015281816106c301528181610751015281816107af0152818161080a0152818161086a015281816108e10152818161092e0152610a0d01526142736000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c80638d67b909116100b8578063ee6c3bcf1161007c578063ee6c3bcf14610306578063ef63552914610319578063efd4532b14610342578063f25de3f814610369578063f8c6681414610389578063fa00e5bd146103b057600080fd5b80638d67b909146102bd57806392ce4ab2146102d0578063b60e9662146102e3578063bafa9107146102f6578063e15234ff146102fe57600080fd5b80636179f7d9116100ff5780636179f7d914610220578063640f65d9146102335780636d14a9871461025a57806372276443146102815780638687feae146102a857600080fd5b8063048886d21461013c578063127af44d146101645780631429c7c2146101795780632ecfe72b1461019e5780634ca22c3f146101e1575b600080fd5b61014f61014a36600461266e565b6103c3565b60405190151581526020015b60405180910390f35b6101776101723660046127e2565b610457565b005b61018c61018736600461266e565b61046d565b60405160ff909116815260200161015b565b6101b16101ac366004612816565b6104fc565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff169082015260600161015b565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161015b565b61017761022e366004612849565b6105a1565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b6102b06106bf565b60405161015b9190612904565b6101776102cb366004612917565b61074c565b6101776102de3660046129a6565b6107a0565b6101776102f1366004612a5e565b6107aa565b6102b0610806565b6102b0610866565b61018c61031436600461266e565b6108c6565b610321610918565b60408051825160ff908116825260209384015116928101929092520161015b565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b61037c610377366004612ac9565b6109ad565b60405161015b9190612cec565b6102087f000000000000000000000000000000000000000000000000000000000000000081565b6101776103be366004612cff565b610a08565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa15801561042d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190612d82565b92915050565b610469610463836104fc565b82610af8565b5050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa1580156104d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190612da4565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa15801561057d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190612dc1565b6104697f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061065288613010565b61065b8861317f565b610663610918565b61066d8a80613316565b6106779080613336565b61068590602081019061334d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610cae92505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561071f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107479190810190613393565b905090565b6104697f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000848461079b610866565b610cdd565b6104698282610af8565b6108007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000868686866107fb610866565b61121d565b50505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa15801561071f573d6000803e3d6000fd5b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561071f573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf906024016104bb565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa158015610989573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107479190613400565b6109b56125b3565b6104517f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610a0385613010565b611a6b565b610af37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610a7d36889003880188613441565b610a868761317f565b610a8f876134dc565b610a97610918565b610aa18a80613316565b610aab9080613336565b610ab990602081019061334d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611c8192505050565b505050565b806020015160ff16816000015160ff1611610bb25760405162461bcd60e51b8152602060048201526075602482015260008051602061421e83398151915260448201527f72696679426c6f625365637572697479506172616d733a20636f6e6669726d6160648201527f74696f6e5468726573686f6c64206d7573742062652067726561746572207468608482015274185b8818591d995c9cd85c9e551a1c995cda1bdb19605a1b60a482015260c4015b60405180910390fd5b60208101518151600091610bc59161361a565b60ff1690506000836020015163ffffffff16846040015160ff1683620f4240610bee9190613653565b610bf89190613653565b610c0490612710613667565b610c0e919061367e565b8451909150610c1f9061271061369d565b63ffffffff168110156108005760405162461bcd60e51b8152602060048201526058602482015260008051602061421e83398151915260448201527f72696679426c6f625365637572697479506172616d733a20736563757269747960648201527f20617373756d7074696f6e7320617265206e6f74206d65740000000000000000608482015260a401610ba9565b6000610cbb878787611a6b565b9050610cd18a8a8a886000015188868989611c81565b50505050505050505050565b6001600160a01b03841663eccbbfc9610cf960208501856136c9565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015610d3b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d5f91906136e6565b610d7d610d6f60408501856136ff565b610d7890613715565b612093565b14610d9a5760405162461bcd60e51b8152600401610ba9906137ec565b610e4e610daa606084018461334d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610dec9250505060408501856136ff565b610df6908061385e565b35610e08610e0387613874565b612104565b604051602001610e1a91815260200190565b60405160208183030381529060405280519060200120856020016020810190610e4391906136c9565b63ffffffff16612134565b610e6a5760405162461bcd60e51b8152600401610ba990613993565b6000805b610e7b60608601866139f5565b90508110156111e557610e9160608601866139f5565b82818110610ea157610ea1613a3e565b610eb7926020608090920201908101915061266e565b60ff16610ec760408601866136ff565b610ed1908061385e565b610edf90602081019061334d565b610eec608088018861334d565b85818110610efc57610efc613a3e565b919091013560f81c9050818110610f1557610f15613a3e565b9050013560f81c60f81b60f81c60ff1614610f425760405162461bcd60e51b8152600401610ba990613a54565b610f4f60608601866139f5565b82818110610f5f57610f5f613a3e565b9050608002016020016020810190610f77919061266e565b60ff16610f8760608701876139f5565b83818110610f9757610f97613a3e565b9050608002016040016020810190610faf919061266e565b60ff1611610fcf5760405162461bcd60e51b8152600401610ba990613ab7565b6001600160a01b038716631429c7c2610feb60608801886139f5565b84818110610ffb57610ffb613a3e565b611011926020608090920201908101915061266e565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa15801561104f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110739190612da4565b60ff1661108360608701876139f5565b8381811061109357611093613a3e565b90506080020160400160208101906110ab919061266e565b60ff1610156110cc5760405162461bcd60e51b8152600401610ba990613b28565b6110d960608601866139f5565b828181106110e9576110e9613a3e565b9050608002016040016020810190611101919061266e565b60ff1661111160408601866136ff565b61111b908061385e565b61112990604081019061334d565b611136608088018861334d565b8581811061114657611146613a3e565b919091013560f81c905081811061115f5761115f613a3e565b9050013560f81c60f81b60f81c60ff16101561118d5760405162461bcd60e51b8152600401610ba990613b28565b6111d18261119e60608801886139f5565b848181106111ae576111ae613a3e565b6111c4926020608090920201908101915061266e565b600160ff919091161b1790565b9150806111dd81613b99565b915050610e6e565b506111f96111f28361214c565b8281161490565b6112155760405162461bcd60e51b8152600401610ba990613bb4565b505050505050565b8382146112ba5760405162461bcd60e51b815260206004820152606b602482015260008051602061421e83398151915260448201527f72696679426c6f6273466f7251756f72756d733a20626c6f624865616465727360648201527f20616e6420626c6f62566572696669636174696f6e50726f6f6673206c656e6760848201526a0e8d040dad2e6dac2e8c6d60ab1b60a482015260c401610ba9565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156112fa573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526113229190810190613393565b905060005b85811015611a6057876001600160a01b031663eccbbfc986868481811061135057611350613a3e565b90506020028101906113629190613c3c565b6113709060208101906136c9565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156113b2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113d691906136e6565b61140b8686848181106113eb576113eb613a3e565b90506020028101906113fd9190613c3c565b610d6f9060408101906136ff565b146114285760405162461bcd60e51b8152600401610ba9906137ec565b61155e85858381811061143d5761143d613a3e565b905060200281019061144f9190613c3c565b61145d90606081019061334d565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508992508891508590508181106114a6576114a6613a3e565b90506020028101906114b89190613c3c565b6114c69060408101906136ff565b6114d0908061385e565b356115018a8a868181106114e6576114e6613a3e565b90506020028101906114f8919061385e565b610e0390613874565b60405160200161151391815260200190565b6040516020818303038152906040528051906020012088888681811061153b5761153b613a3e565b905060200281019061154d9190613c3c565b610e439060408101906020016136c9565b61157a5760405162461bcd60e51b8152600401610ba990613993565b6000805b88888481811061159057611590613a3e565b90506020028101906115a2919061385e565b6115b09060608101906139f5565b9050811015611a26578888848181106115cb576115cb613a3e565b90506020028101906115dd919061385e565b6115eb9060608101906139f5565b828181106115fb576115fb613a3e565b611611926020608090920201908101915061266e565b60ff1687878581811061162657611626613a3e565b90506020028101906116389190613c3c565b6116469060408101906136ff565b611650908061385e565b61165e90602081019061334d565b89898781811061167057611670613a3e565b90506020028101906116829190613c3c565b61169090608081019061334d565b858181106116a0576116a0613a3e565b919091013560f81c90508181106116b9576116b9613a3e565b9050013560f81c60f81b60f81c60ff16146116e65760405162461bcd60e51b8152600401610ba990613a54565b8888848181106116f8576116f8613a3e565b905060200281019061170a919061385e565b6117189060608101906139f5565b8281811061172857611728613a3e565b9050608002016020016020810190611740919061266e565b60ff1689898581811061175557611755613a3e565b9050602002810190611767919061385e565b6117759060608101906139f5565b8381811061178557611785613a3e565b905060800201604001602081019061179d919061266e565b60ff16116117bd5760405162461bcd60e51b8152600401610ba990613ab7565b838989858181106117d0576117d0613a3e565b90506020028101906117e2919061385e565b6117f09060608101906139f5565b8381811061180057611800613a3e565b611816926020608090920201908101915061266e565b60ff168151811061182957611829613a3e565b016020015160f81c89898581811061184357611843613a3e565b9050602002810190611855919061385e565b6118639060608101906139f5565b8381811061187357611873613a3e565b905060800201604001602081019061188b919061266e565b60ff1610156118ac5760405162461bcd60e51b8152600401610ba990613b28565b8888848181106118be576118be613a3e565b90506020028101906118d0919061385e565b6118de9060608101906139f5565b828181106118ee576118ee613a3e565b9050608002016040016020810190611906919061266e565b60ff1687878581811061191b5761191b613a3e565b905060200281019061192d9190613c3c565b61193b9060408101906136ff565b611945908061385e565b61195390604081019061334d565b89898781811061196557611965613a3e565b90506020028101906119779190613c3c565b61198590608081019061334d565b8581811061199557611995613a3e565b919091013560f81c90508181106119ae576119ae613a3e565b9050013560f81c60f81b60f81c60ff1610156119dc5760405162461bcd60e51b8152600401610ba990613b28565b611a12828a8a868181106119f2576119f2613a3e565b9050602002810190611a04919061385e565b61119e9060608101906139f5565b915080611a1e81613b99565b91505061157e565b50611a336111f28561214c565b611a4f5760405162461bcd60e51b8152600401610ba990613bb4565b50611a5981613b99565b9050611327565b505050505050505050565b611a736125b3565b602082015151516000906001600160401b03811115611a9457611a946126a9565b604051908082528060200260200182016040528015611abd578160200160208202803683370190505b50905060005b60208401515151811015611b3a57611b0d8460200151600001518281518110611aee57611aee613a3e565b6020026020010151805160009081526020918201519091526040902090565b828281518110611b1f57611b1f613a3e565b6020908102919091010152611b3381613b99565b9050611ac3565b50606060005b84602001516080015151811015611ba757818560200151608001518281518110611b6c57611b6c613a3e565b6020026020010151604051602001611b85929190613c52565b604051602081830303815290604052915080611ba090613b99565b9050611b40565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491611be2918a919087908990600401613c84565b600060405180830381865afa158015611bff573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611c279190810190613dd7565b805185526020958601805151878701528051870151604080880191909152815160609081015181890152915181015160808801529682015160a08701529581015160c0860152949094015160e08401525090949350505050565b611cd384604001518660000151611c9b87600001516122d9565b604051602001611cad91815260200190565b60405160208183030381529060405280519060200120876020015163ffffffff16612134565b611d4c5760405162461bcd60e51b8152602060048201526050602482015260008051602061421e83398151915260448201527f72696679426c6f625632466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a401610ba9565b600080886001600160a01b0316636efb4636611d67896122fc565b885151602090810151908b01516040516001600160e01b031960e086901b168152611d99939291908b90600401613eaf565b600060405180830381865afa158015611db6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611dde9190810190613f62565b91509150611df488876000015160200151612327565b85515151604051632ecfe72b60e01b815261ffff9091166004820152611e6f906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa158015611e45573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e699190612dc1565b85610af8565b6000805b8751516020015151811015611fda57856000015160ff1684602001518281518110611ea057611ea0613a3e565b6020026020010151611eb29190613ffe565b6001600160601b0316606485600001518381518110611ed357611ed3613a3e565b60200260200101516001600160601b0316611eee919061367e565b1015611f955760405162461bcd60e51b8152602060048201526076602482015260008051602061421e83398151915260448201527f72696679426c6f625632466f7251756f72756d733a207369676e61746f72696560648201527f7320646f206e6f74206f776e206174206c65617374207468726573686f6c642060848201527570657263656e74616765206f6620612071756f72756d60501b60a482015260c401610ba9565b875151602001518051611fc691849184908110611fb457611fb4613a3e565b0160200151600160f89190911c1b1790565b915080611fd281613b99565b915050611e73565b50611fe76111f28561214c565b6120865760405162461bcd60e51b8152602060048201526070602482015260008051602061421e83398151915260448201527f72696679426c6f625632466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c401610ba9565b5050505050505050505050565b600061045182600001516040516020016120ad9190614024565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b6000816040516020016121179190614084565b604051602081830303815290604052805190602001209050919050565b60008361214286858561245a565b1495945050505050565b6000610100825111156121d55760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610ba9565b81516121e357506000919050565b600080836000815181106121f9576121f9613a3e565b0160200151600160f89190911c81901b92505b84518110156122d05784818151811061222757612227613a3e565b0160200151600160f89190911c1b91508282116122bc5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610ba9565b918117916122c981613b99565b905061220c565b50909392505050565b60006122e8826000015161255d565b602080840151604051612117939201614129565b60008160405160200161211791908151815260209182015163ffffffff169181019190915260400190565b60005b8151811015610af35760006001600160a01b0316836001600160a01b031663b5a872da84848151811061235f5761235f613a3e565b60200260200101516040518263ffffffff1660e01b815260040161238f919063ffffffff91909116815260200190565b602060405180830381865afa1580156123ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123d09190614142565b6001600160a01b0316141561244a5760405162461bcd60e51b8152602060048201526046602482015260008051602061421e83398151915260448201527f7269667952656c61794b6579735365743a2072656c6179206b6579206973206e6064820152651bdd081cd95d60d21b608482015260a401610ba9565b61245381613b99565b905061232a565b60006020845161246a919061416b565b156124f15760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610ba9565b8260205b855181116125545761250860028561416b565b61252957816000528086015160205260406000209150600284049350612542565b8086015160005281602052604060002091506002840493505b61254d60208261417f565b90506124f5565b50949350505050565b600081600001518260200151836040015184606001516040516020016125869493929190614197565b60408051601f19818403018152828252805160209182012060808601519184015290820152606001612117565b6040518061010001604052806060815260200160608152602001606081526020016125dc612619565b81526020016125fe604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b604051806040016040528061262c61263e565b815260200161263961263e565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff8116811461266b57600080fd5b50565b60006020828403121561268057600080fd5b813561268b8161265c565b9392505050565b803561ffff811681146126a457600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156126e1576126e16126a9565b60405290565b604051606081016001600160401b03811182821017156126e1576126e16126a9565b60405160a081016001600160401b03811182821017156126e1576126e16126a9565b604051608081016001600160401b03811182821017156126e1576126e16126a9565b60405161010081016001600160401b03811182821017156126e1576126e16126a9565b604051601f8201601f191681016001600160401b0381118282101715612798576127986126a9565b604052919050565b6000604082840312156127b257600080fd5b6127ba6126bf565b905081356127c78161265c565b815260208201356127d78161265c565b602082015292915050565b600080606083850312156127f557600080fd5b6127fe83612692565b915061280d84602085016127a0565b90509250929050565b60006020828403121561282857600080fd5b61268b82612692565b60006060828403121561284357600080fd5b50919050565b6000806040838503121561285c57600080fd5b82356001600160401b038082111561287357600080fd5b61287f86838701612831565b9350602085013591508082111561289557600080fd5b506128a285828601612831565b9150509250929050565b60005b838110156128c75781810151838201526020016128af565b838111156108005750506000910152565b600081518084526128f08160208601602086016128ac565b601f01601f19169290920160200192915050565b60208152600061268b60208301846128d8565b6000806040838503121561292a57600080fd5b82356001600160401b038082111561294157600080fd5b908401906080828703121561295557600080fd5b9092506020840135908082111561296b57600080fd5b50830160a0818603121561297e57600080fd5b809150509250929050565b63ffffffff8116811461266b57600080fd5b80356126a481612989565b60008082840360a08112156129ba57600080fd5b60608112156129c857600080fd5b506129d16126e7565b83356129dc81612989565b815260208401356129ec81612989565b602082015260408401356129ff8161265c565b6040820152915061280d84606085016127a0565b60008083601f840112612a2557600080fd5b5081356001600160401b03811115612a3c57600080fd5b6020830191508360208260051b8501011115612a5757600080fd5b9250929050565b60008060008060408587031215612a7457600080fd5b84356001600160401b0380821115612a8b57600080fd5b612a9788838901612a13565b90965094506020870135915080821115612ab057600080fd5b50612abd87828801612a13565b95989497509550505050565b600060208284031215612adb57600080fd5b81356001600160401b03811115612af157600080fd5b612afd84828501612831565b949350505050565b600081518084526020808501945080840160005b83811015612b3b57815163ffffffff1687529582019590820190600101612b19565b509495945050505050565b600081518084526020808501945080840160005b83811015612b3b57612b7787835180518252602090810151910152565b6040969096019590820190600101612b5a565b8060005b6002811015610800578151845260209384019390910190600101612b8e565b612bb8828251612b8a565b6020810151610af36040840182612b8a565b600081518084526020808501808196508360051b8101915082860160005b85811015612c12578284038952612c00848351612b05565b98850198935090840190600101612be8565b5091979650505050505050565b60006101808251818552612c3582860182612b05565b91505060208301518482036020860152612c4f8282612b46565b91505060408301518482036040860152612c698282612b46565b9150506060830151612c7e6060860182612bad565b506080830151805160e08601526020015161010085015260a0830151848203610120860152612cad8282612b05565b91505060c0830151848203610140860152612cc88282612b05565b91505060e0830151848203610160860152612ce38282612bca565b95945050505050565b60208152600061268b6020830184612c1f565b60008060008385036080811215612d1557600080fd5b6040811215612d2357600080fd5b5083925060408401356001600160401b0380821115612d4157600080fd5b612d4d87838801612831565b93506060860135915080821115612d6357600080fd5b5084016101808187031215612d7757600080fd5b809150509250925092565b600060208284031215612d9457600080fd5b8151801515811461268b57600080fd5b600060208284031215612db657600080fd5b815161268b8161265c565b600060608284031215612dd357600080fd5b604051606081018181106001600160401b0382111715612df557612df56126a9565b6040528251612e0381612989565b81526020830151612e1381612989565b60208201526040830151612e268161265c565b60408201529392505050565b600060408284031215612e4457600080fd5b612e4c6126bf565b90508135815260208201356127d781612989565b60006001600160401b03821115612e7957612e796126a9565b5060051b60200190565b600060408284031215612e9557600080fd5b612e9d6126bf565b9050813581526020820135602082015292915050565b600082601f830112612ec457600080fd5b81356020612ed9612ed483612e60565b612770565b82815260069290921b84018101918181019086841115612ef857600080fd5b8286015b84811015612f1c57612f0e8882612e83565b835291830191604001612efc565b509695505050505050565b600082601f830112612f3857600080fd5b612f406126bf565b806040840185811115612f5257600080fd5b845b81811015612f6c578035845260209384019301612f54565b509095945050505050565b600060808284031215612f8957600080fd5b612f916126bf565b9050612f9d8383612f27565b81526127d78360408401612f27565b600082601f830112612fbd57600080fd5b81356020612fcd612ed483612e60565b82815260059290921b84018101918181019086841115612fec57600080fd5b8286015b84811015612f1c57803561300381612989565b8352918301918301612ff0565b60006060823603121561302257600080fd5b61302a6126bf565b6130343684612e32565b815260408301356001600160401b038082111561305057600080fd5b8185019150610120823603121561306657600080fd5b61306e612709565b82358281111561307d57600080fd5b61308936828601612eb3565b82525060208301358281111561309e57600080fd5b6130aa36828601612eb3565b6020830152506130bd3660408501612e83565b60408201526130cf3660808501612f77565b6060820152610100830135828111156130e757600080fd5b6130f336828601612fac565b608083015250602084015250909392505050565b60006001600160401b03821115613120576131206126a9565b50601f01601f191660200190565b600082601f83011261313f57600080fd5b813561314d612ed482613107565b81815284602083860101111561316257600080fd5b816020850160208301376000918101602001919091529392505050565b60006060823603121561319157600080fd5b6131996126e7565b82356001600160401b03808211156131b057600080fd5b818501915060408083360312156131c657600080fd5b6131ce6126bf565b8335838111156131dd57600080fd5b8401368190036101e08112156131f257600080fd5b6131fa612709565b61320383612692565b81526020808401358781111561321857600080fd5b6132243682870161312e565b8383015250610160603f198401121561323c57600080fd5b61324461272b565b925061325236878601612e83565b83526132613660808601612f77565b81840152613273366101008601612f77565b8684015261018084013561328681612989565b606084015281860183905261329e6101a0850161299b565b60608301526101c08401356080830152818552808801359350868411156132c457600080fd5b6132d036858a01612fac565b818601528489526132e2818c0161299b565b9089015250505050858101359250818311156132fd57600080fd5b6133093684880161312e565b9084015250909392505050565b60008235603e1983360301811261332c57600080fd5b9190910192915050565b600082356101de1983360301811261332c57600080fd5b6000808335601e1984360301811261336457600080fd5b8301803591506001600160401b0382111561337e57600080fd5b602001915036819003821315612a5757600080fd5b6000602082840312156133a557600080fd5b81516001600160401b038111156133bb57600080fd5b8201601f810184136133cc57600080fd5b80516133da612ed482613107565b8181528560208385010111156133ef57600080fd5b612ce38260208301602086016128ac565b60006040828403121561341257600080fd5b61341a6126bf565b82516134258161265c565b815260208301516134358161265c565b60208201529392505050565b60006040828403121561345357600080fd5b61268b8383612e32565b600082601f83011261346e57600080fd5b8135602061347e612ed483612e60565b82815260059290921b8401810191818101908684111561349d57600080fd5b8286015b84811015612f1c5780356001600160401b038111156134c05760008081fd5b6134ce8986838b0101612fac565b8452509183019183016134a1565b600061018082360312156134ef57600080fd5b6134f761274d565b82356001600160401b038082111561350e57600080fd5b61351a36838701612fac565b8352602085013591508082111561353057600080fd5b61353c36838701612eb3565b6020840152604085013591508082111561355557600080fd5b61356136838701612eb3565b60408401526135733660608701612f77565b60608401526135853660e08701612e83565b608084015261012085013591508082111561359f57600080fd5b6135ab36838701612fac565b60a08401526101408501359150808211156135c557600080fd5b6135d136838701612fac565b60c08401526101608501359150808211156135eb57600080fd5b506135f83682860161345d565b60e08301525092915050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff84168082101561363457613634613604565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826136625761366261363d565b500490565b60008282101561367957613679613604565b500390565b600081600019048311821515161561369857613698613604565b500290565b600063ffffffff808316818516818304811182151516156136c0576136c0613604565b02949350505050565b6000602082840312156136db57600080fd5b813561268b81612989565b6000602082840312156136f857600080fd5b5051919050565b60008235605e1983360301811261332c57600080fd5b60006060823603121561372757600080fd5b61372f6126e7565b82356001600160401b038082111561374657600080fd5b81850191506080823603121561375b57600080fd5b61376361272b565b8235815260208301358281111561377957600080fd5b6137853682860161312e565b60208301525060408301358281111561379d57600080fd5b6137a93682860161312e565b604083015250606083013592506137bf83612989565b826060820152808452505050602083013560208201526137e16040840161299b565b604082015292915050565b6020808252606090820181905260008051602061421e83398151915260408301527f72696679426c6f62466f7251756f72756d733a2062617463684d657461646174908201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608082015260a00190565b60008235607e1983360301811261332c57600080fd5b6000608080833603121561388757600080fd5b61388f6126e7565b6138993685612e83565b81526040808501356138aa81612989565b6020818185015260609150818701356001600160401b038111156138cd57600080fd5b870136601f8201126138de57600080fd5b80356138ec612ed482612e60565b81815260079190911b8201830190838101903683111561390b57600080fd5b928401925b8284101561397f578884360312156139285760008081fd5b61393061272b565b843561393b8161265c565b81528486013561394a8161265c565b818701528488013561395b8161265c565b818901528487013561396c81612989565b8188015282529288019290840190613910565b958701959095525093979650505050505050565b6020808252604e9082015260008051602061421e83398151915260408201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260608201526d1bdbd9881a5cc81a5b9d985b1a5960921b608082015260a00190565b6000808335601e19843603018112613a0c57600080fd5b8301803591506001600160401b03821115613a2657600080fd5b6020019150600781901b3603821315612a5757600080fd5b634e487b7160e01b600052603260045260246000fd5b6020808252604f9082015260008051602061421e83398151915260408201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260608201526e040c8decae640dcdee840dac2e8c6d608b1b608082015260a00190565b602080825260579082015260008051602061421e83398151915260408201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560608201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608082015260a00190565b6020808252605e9082015260008051602061421e83398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b6000600019821415613bad57613bad613604565b5060010190565b6020808252606e9082015260008051602061421e83398151915260408201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60608201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60808201526d6669726d65642071756f72756d7360901b60a082015260c00190565b60008235609e1983360301811261332c57600080fd5b60008351613c648184602088016128ac565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff86168184015260806040840152613cb260808401866128d8565b838103606085015284518082528286019183019060005b81811015613ce557835183529284019291840191600101613cc9565b50909998505050505050505050565b600082601f830112613d0557600080fd5b81516020613d15612ed483612e60565b82815260059290921b84018101918181019086841115613d3457600080fd5b8286015b84811015612f1c578051613d4b81612989565b8352918301918301613d38565b600082601f830112613d6957600080fd5b81516020613d79612ed483612e60565b82815260059290921b84018101918181019086841115613d9857600080fd5b8286015b84811015612f1c5780516001600160401b03811115613dbb5760008081fd5b613dc98986838b0101613cf4565b845250918301918301613d9c565b600060208284031215613de957600080fd5b81516001600160401b0380821115613e0057600080fd5b9083019060808286031215613e1457600080fd5b613e1c61272b565b825182811115613e2b57600080fd5b613e3787828601613cf4565b825250602083015182811115613e4c57600080fd5b613e5887828601613cf4565b602083015250604083015182811115613e7057600080fd5b613e7c87828601613cf4565b604083015250606083015182811115613e9457600080fd5b613ea087828601613d58565b60608301525095945050505050565b848152608060208201526000613ec860808301866128d8565b63ffffffff851660408401528281036060840152613ee68185612c1f565b979650505050505050565b600082601f830112613f0257600080fd5b81516020613f12612ed483612e60565b82815260059290921b84018101918181019086841115613f3157600080fd5b8286015b84811015612f1c5780516001600160601b0381168114613f555760008081fd5b8352918301918301613f35565b60008060408385031215613f7557600080fd5b82516001600160401b0380821115613f8c57600080fd5b9084019060408287031215613fa057600080fd5b613fa86126bf565b825182811115613fb757600080fd5b613fc388828601613ef1565b825250602083015182811115613fd857600080fd5b613fe488828601613ef1565b602083015250809450505050602083015190509250929050565b60006001600160601b03808316818516818304811182151516156136c0576136c0613604565b6020815281516020820152600060208301516080604084015261404a60a08401826128d8565b90506040840151601f1984830301606085015261406782826128d8565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b8084101561411a578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906140d7565b509a9950505050505050505050565b828152604060208201526000612afd6040830184612b05565b60006020828403121561415457600080fd5b81516001600160a01b038116811461268b57600080fd5b60008261417a5761417a61363d565b500690565b6000821982111561419257614192613604565b500190565b60006101c061ffff871683528060208401526141b5818401876128d8565b85518051604086015260200151606085015291506141d09050565b60208401516141e26080840182612bad565b5060408401516141f6610100840182612bad565b506060939093015163ffffffff908116610180830152919091166101a0909101529291505056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a2646970667358221220c30d156ddbf7d6bc6084f92df06805851f15524b92d975e46b0733d3aba4e65864736f6c634300080c0033",
}

// ContractEigenDABlobVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.ABI instead.
var ContractEigenDABlobVerifierABI = ContractEigenDABlobVerifierMetaData.ABI

// ContractEigenDABlobVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.Bin instead.
var ContractEigenDABlobVerifierBin = ContractEigenDABlobVerifierMetaData.Bin

// DeployContractEigenDABlobVerifier deploys a new Ethereum contract, binding an instance of ContractEigenDABlobVerifier to it.
func DeployContractEigenDABlobVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDABatchMetadataStorage common.Address, _eigenDASignatureVerifier common.Address, _eigenDARelayRegistry common.Address, _operatorStateRetriever common.Address, _registryCoordinator common.Address) (common.Address, *types.Transaction, *ContractEigenDABlobVerifier, error) {
	parsed, err := ContractEigenDABlobVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDABlobVerifierBin), backend, _eigenDAThresholdRegistry, _eigenDABatchMetadataStorage, _eigenDASignatureVerifier, _eigenDARelayRegistry, _operatorStateRetriever, _registryCoordinator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDABlobVerifier{ContractEigenDABlobVerifierCaller: ContractEigenDABlobVerifierCaller{contract: contract}, ContractEigenDABlobVerifierTransactor: ContractEigenDABlobVerifierTransactor{contract: contract}, ContractEigenDABlobVerifierFilterer: ContractEigenDABlobVerifierFilterer{contract: contract}}, nil
}

// ContractEigenDABlobVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDABlobVerifier struct {
	ContractEigenDABlobVerifierCaller     // Read-only binding to the contract
	ContractEigenDABlobVerifierTransactor // Write-only binding to the contract
	ContractEigenDABlobVerifierFilterer   // Log filterer for contract events
}

// ContractEigenDABlobVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDABlobVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDABlobVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDABlobVerifierSession struct {
	Contract     *ContractEigenDABlobVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractEigenDABlobVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDABlobVerifierCallerSession struct {
	Contract *ContractEigenDABlobVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractEigenDABlobVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDABlobVerifierTransactorSession struct {
	Contract     *ContractEigenDABlobVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractEigenDABlobVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierRaw struct {
	Contract *ContractEigenDABlobVerifier // Generic contract binding to access the raw methods on
}

// ContractEigenDABlobVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierCallerRaw struct {
	Contract *ContractEigenDABlobVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDABlobVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDABlobVerifierTransactorRaw struct {
	Contract *ContractEigenDABlobVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDABlobVerifier creates a new instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifier(address common.Address, backend bind.ContractBackend) (*ContractEigenDABlobVerifier, error) {
	contract, err := bindContractEigenDABlobVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifier{ContractEigenDABlobVerifierCaller: ContractEigenDABlobVerifierCaller{contract: contract}, ContractEigenDABlobVerifierTransactor: ContractEigenDABlobVerifierTransactor{contract: contract}, ContractEigenDABlobVerifierFilterer: ContractEigenDABlobVerifierFilterer{contract: contract}}, nil
}

// NewContractEigenDABlobVerifierCaller creates a new read-only instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDABlobVerifierCaller, error) {
	contract, err := bindContractEigenDABlobVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierCaller{contract: contract}, nil
}

// NewContractEigenDABlobVerifierTransactor creates a new write-only instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDABlobVerifierTransactor, error) {
	contract, err := bindContractEigenDABlobVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierTransactor{contract: contract}, nil
}

// NewContractEigenDABlobVerifierFilterer creates a new log filterer instance of ContractEigenDABlobVerifier, bound to a specific deployed contract.
func NewContractEigenDABlobVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDABlobVerifierFilterer, error) {
	contract, err := bindContractEigenDABlobVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierFilterer{contract: contract}, nil
}

// bindContractEigenDABlobVerifier binds a generic wrapper to an already deployed contract.
func bindContractEigenDABlobVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDABlobVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.ContractEigenDABlobVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDABlobVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDABlobVerifier.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDABatchMetadataStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDABatchMetadataStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDABatchMetadataStorage is a free data retrieval call binding the contract method 0x640f65d9.
//
// Solidity: function eigenDABatchMetadataStorage() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDABatchMetadataStorage() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDABatchMetadataStorage(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDARelayRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDARelayRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDARelayRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDARelayRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDABlobVerifier.CallOpts)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(VersionedBlobParams)).(*VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDABlobVerifier.Contract.GetBlobParams(&_ContractEigenDABlobVerifier.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDABlobVerifier.Contract.GetBlobParams(&_ContractEigenDABlobVerifier.CallOpts, version)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetDefaultSecurityThresholdsV2(opts *bind.CallOpts) (SecurityThresholds, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getDefaultSecurityThresholdsV2")

	if err != nil {
		return *new(SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(SecurityThresholds)).(*SecurityThresholds)

	return out0, err

}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDABlobVerifier.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDABlobVerifier.CallOpts)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDABlobVerifier.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDABlobVerifier.CallOpts)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDABlobVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDABlobVerifier.Contract.GetIsQuorumRequired(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDABlobVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDABlobVerifier.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDABlobVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDABlobVerifier.CallOpts, signedBatch)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDABlobVerifier.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDABlobVerifier.CallOpts, quorumNumber)
}

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) OperatorStateRetriever(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "operatorStateRetriever")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) OperatorStateRetriever() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.OperatorStateRetriever(&_ContractEigenDABlobVerifier.CallOpts)
}

// OperatorStateRetriever is a free data retrieval call binding the contract method 0x4ca22c3f.
//
// Solidity: function operatorStateRetriever() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) OperatorStateRetriever() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.OperatorStateRetriever(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDABlobVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDABlobVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDABlobVerifier.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.RegistryCoordinator(&_ContractEigenDABlobVerifier.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDABlobVerifier.Contract.RegistryCoordinator(&_ContractEigenDABlobVerifier.CallOpts)
}

// VerifyBlobSecurityParams is a free data retrieval call binding the contract method 0x127af44d.
//
// Solidity: function verifyBlobSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobSecurityParams(opts *bind.CallOpts, version uint16, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobSecurityParams", version, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobSecurityParams is a free data retrieval call binding the contract method 0x127af44d.
//
// Solidity: function verifyBlobSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobSecurityParams(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobSecurityParams(&_ContractEigenDABlobVerifier.CallOpts, version, securityThresholds)
}

// VerifyBlobSecurityParams is a free data retrieval call binding the contract method 0x127af44d.
//
// Solidity: function verifyBlobSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobSecurityParams(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobSecurityParams(&_ContractEigenDABlobVerifier.CallOpts, version, securityThresholds)
}

// VerifyBlobSecurityParams0 is a free data retrieval call binding the contract method 0x92ce4ab2.
//
// Solidity: function verifyBlobSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobSecurityParams0(opts *bind.CallOpts, blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobSecurityParams0", blobParams, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobSecurityParams0 is a free data retrieval call binding the contract method 0x92ce4ab2.
//
// Solidity: function verifyBlobSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobSecurityParams0(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobSecurityParams0(&_ContractEigenDABlobVerifier.CallOpts, blobParams, securityThresholds)
}

// VerifyBlobSecurityParams0 is a free data retrieval call binding the contract method 0x92ce4ab2.
//
// Solidity: function verifyBlobSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobSecurityParams0(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobSecurityParams0(&_ContractEigenDABlobVerifier.CallOpts, blobParams, securityThresholds)
}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyBlobV1 is a free data retrieval call binding the contract method 0x8d67b909.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0xfa00e5bd.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2", batchHeader, blobVerificationProof, nonSignerStakesAndSignature)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0xfa00e5bd.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0xfa00e5bd.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x6179f7d9.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch", signedBatch, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x6179f7d9.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x6179f7d9.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),uint32,bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobsV1 is a free data retrieval call binding the contract method 0xb60e9662.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobsV1(opts *bind.CallOpts, blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobsV1 is a free data retrieval call binding the contract method 0xb60e9662.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobsV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyBlobsV1 is a free data retrieval call binding the contract method 0xb60e9662.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobsV1(&_ContractEigenDABlobVerifier.CallOpts, blobHeaders, blobVerificationProofs)
}

// ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated)
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
func (it *ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 SecurityThresholds
	NewDefaultSecurityThresholdsV2      SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractEigenDABlobVerifier.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated)
				if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractEigenDABlobVerifierDefaultSecurityThresholdsV2Updated)
	if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractEigenDABlobVerifier.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDABlobVerifierQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractEigenDABlobVerifier.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDABlobVerifierQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated)
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
func (it *ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierQuorumNumbersRequiredUpdatedIterator{contract: _ContractEigenDABlobVerifier.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDABlobVerifier.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated)
				if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated, error) {
	event := new(ContractEigenDABlobVerifierQuorumNumbersRequiredUpdated)
	if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator struct {
	Event *ContractEigenDABlobVerifierVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDABlobVerifierVersionedBlobParamsAdded)
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
		it.Event = new(ContractEigenDABlobVerifierVersionedBlobParamsAdded)
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
func (it *ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDABlobVerifierVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractEigenDABlobVerifier contract.
type ContractEigenDABlobVerifierVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDABlobVerifier.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDABlobVerifierVersionedBlobParamsAddedIterator{contract: _ContractEigenDABlobVerifier.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDABlobVerifierVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDABlobVerifier.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDABlobVerifierVersionedBlobParamsAdded)
				if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractEigenDABlobVerifierVersionedBlobParamsAdded, error) {
	event := new(ContractEigenDABlobVerifierVersionedBlobParamsAdded)
	if err := _ContractEigenDABlobVerifier.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
