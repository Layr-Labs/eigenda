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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620043413803806200434183398101604081905262000035916200007e565b6001600160a01b0395861660805293851660a05291841660c052831660e052821661010052166101205262000112565b6001600160a01b03811681146200007b57600080fd5b50565b60008060008060008060c087890312156200009857600080fd5b8651620000a58162000065565b6020880151909650620000b88162000065565b6040880151909550620000cb8162000065565b6060880151909450620000de8162000065565b6080880151909350620000f18162000065565b60a0880151909250620001048162000065565b809150509295509295509295565b60805160a05160c05160e0516101005161012051614145620001fc6000396000818161025f015281816108720152610a2d0152600081816101f9015281816108510152610a0c0152600081816102860152818161049e015261083001526000818161035a0152818161047d015261080f015260008181610238015281816106f501526107530152600081816103a1015281816103de0152818161045c01528181610529015281816105d701528181610646015281816106d4015281816107320152818161078d015281816107ee015281816108bc01528181610933015261098001526141456000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c80638d67b909116100b8578063e15234ff1161007c578063e15234ff14610311578063ee6c3bcf14610319578063ef6355291461032c578063efd4532b14610355578063f25de3f81461037c578063f8c668141461039c57600080fd5b80638d67b909146102bd57806392ce4ab2146102d0578063b60e9662146102e3578063bafa9107146102f6578063ca7fd71d146102fe57600080fd5b80634ca22c3f116100ff5780634ca22c3f146101f4578063640f65d9146102335780636d14a9871461025a57806372276443146102815780638687feae146102a857600080fd5b8063048886d21461013c5780630a4715f414610164578063127af44d146101795780631429c7c21461018c5780632ecfe72b146101b1575b600080fd5b61014f61014a3660046125a4565b6103c3565b60405190151581526020015b60405180910390f35b6101776101723660046125e0565b610457565b005b6101776101873660046127b3565b6104f8565b61019f61019a3660046125a4565b61050e565b60405160ff909116815260200161015b565b6101c46101bf3660046127e7565b61059d565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff169082015260600161015b565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161015b565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b6102b0610642565b60405161015b919061285a565b6101776102cb36600461286d565b6106cf565b6101776102de3660046128fc565b610723565b6101776102f13660046129b4565b61072d565b6102b0610789565b61017761030c366004612a1f565b6107e9565b6102b06108b8565b61019f6103273660046125a4565b610918565b61033461096a565b60408051825160ff908116825260209384015116928101929092520161015b565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b61038f61038a366004612a82565b6109ff565b60405161015b9190612ca5565b61021b7f000000000000000000000000000000000000000000000000000000000000000081565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa15801561042d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190612cb8565b92915050565b6104f37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006104cc36889003880188612d08565b6104d587612ee8565b6104de87613151565b6104e661096a565b6104ee6108b8565b610a5a565b505050565b61050a6105048361059d565b82610e77565b5050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa158015610579573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190613279565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa15801561061e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190613296565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156106a2573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106ca9190810190613307565b905090565b61050a7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000848461071e6108b8565b611028565b61050a8282610e77565b6107837f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008686868661077e6108b8565b611561565b50505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156106a2573d6000803e3d6000fd5b61050a7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061089a88613374565b6108a388612ee8565b6108ab61096a565b6108b36108b8565b611daf565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156106a2573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf9060240161055c565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa1580156109db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106ca919061346b565b610a076124e9565b6104517f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610a5585613374565b611dde565b60408401518551855151610aab929190610a7390611ff4565b604051602001610a8591815260200190565b60405160208183030381529060405280519060200120876020015163ffffffff1661205f565b610b295760405162461bcd60e51b815260206004820152605060248201526000805160206140f083398151915260448201527f72696679426c6f625632466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a4015b60405180910390fd5b600080886001600160a01b0316636efb4636610b4489612077565b885151602090810151908b01516040516001600160e01b031960e086901b168152610b76939291908b906004016134ac565b600060405180830381865afa158015610b93573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610bbb919081019061355f565b91509150610bd1888760000151602001516120a2565b85515151604051632ecfe72b60e01b815261ffff9091166004820152610c4c906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa158015610c22573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c469190613296565b85610e77565b6000805b8751516020015151811015610db757856000015160ff1684602001518281518110610c7d57610c7d6135fb565b6020026020010151610c8f9190613627565b6001600160601b0316606485600001518381518110610cb057610cb06135fb565b60200260200101516001600160601b0316610ccb9190613656565b1015610d725760405162461bcd60e51b815260206004820152607660248201526000805160206140f083398151915260448201527f72696679426c6f625632466f7251756f72756d733a207369676e61746f72696560648201527f7320646f206e6f74206f776e206174206c65617374207468726573686f6c642060848201527570657263656e74616765206f6620612071756f72756d60501b60a482015260c401610b20565b875151602001518051610da391849184908110610d9157610d916135fb565b0160200151600160f89190911c1b1790565b915080610daf81613675565b915050610c50565b50610dcb610dc4856121d5565b8281161490565b610e6a5760405162461bcd60e51b815260206004820152607060248201526000805160206140f083398151915260448201527f72696679426c6f625632466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c401610b20565b5050505050505050505050565b806020015160ff16816000015160ff1611610f2c5760405162461bcd60e51b815260206004820152607560248201526000805160206140f083398151915260448201527f72696679426c6f625365637572697479506172616d733a20636f6e6669726d6160648201527f74696f6e5468726573686f6c64206d7573742062652067726561746572207468608482015274185b8818591d995c9cd85c9e551a1c995cda1bdb19605a1b60a482015260c401610b20565b60208101518151600091610f3f91613690565b60ff1690506000836020015163ffffffff16846040015160ff1683620f4240610f6891906136c9565b610f7291906136c9565b610f7e906127106136dd565b610f889190613656565b8451909150610f99906127106136f4565b63ffffffff168110156107835760405162461bcd60e51b815260206004820152605860248201526000805160206140f083398151915260448201527f72696679426c6f625365637572697479506172616d733a20736563757269747960648201527f20617373756d7074696f6e7320617265206e6f74206d65740000000000000000608482015260a401610b20565b6001600160a01b03841663eccbbfc96110446020850185613717565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611086573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110aa9190613734565b6110c86110ba604085018561374d565b6110c39061376d565b612362565b146110e55760405162461bcd60e51b8152600401610b2090613844565b6111996110f560608401846138b6565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061113792505050604085018561374d565b61114190806138fc565b3561115361114e87613912565b6123d3565b60405160200161116591815260200190565b6040516020818303038152906040528051906020012085602001602081019061118e9190613717565b63ffffffff1661205f565b6111b55760405162461bcd60e51b8152600401610b2090613a31565b6000805b6111c66060860186613a93565b9050811015611530576111dc6060860186613a93565b828181106111ec576111ec6135fb565b61120292602060809092020190810191506125a4565b60ff16611212604086018661374d565b61121c90806138fc565b61122a9060208101906138b6565b61123760808801886138b6565b85818110611247576112476135fb565b919091013560f81c9050818110611260576112606135fb565b9050013560f81c60f81b60f81c60ff161461128d5760405162461bcd60e51b8152600401610b2090613adc565b61129a6060860186613a93565b828181106112aa576112aa6135fb565b90506080020160200160208101906112c291906125a4565b60ff166112d26060870187613a93565b838181106112e2576112e26135fb565b90506080020160400160208101906112fa91906125a4565b60ff161161131a5760405162461bcd60e51b8152600401610b2090613b3f565b6001600160a01b038716631429c7c26113366060880188613a93565b84818110611346576113466135fb565b61135c92602060809092020190810191506125a4565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa15801561139a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113be9190613279565b60ff166113ce6060870187613a93565b838181106113de576113de6135fb565b90506080020160400160208101906113f691906125a4565b60ff1610156114175760405162461bcd60e51b8152600401610b2090613bb0565b6114246060860186613a93565b82818110611434576114346135fb565b905060800201604001602081019061144c91906125a4565b60ff1661145c604086018661374d565b61146690806138fc565b6114749060408101906138b6565b61148160808801886138b6565b85818110611491576114916135fb565b919091013560f81c90508181106114aa576114aa6135fb565b9050013560f81c60f81b60f81c60ff1610156114d85760405162461bcd60e51b8152600401610b2090613bb0565b61151c826114e96060880188613a93565b848181106114f9576114f96135fb565b61150f92602060809092020190810191506125a4565b600160ff919091161b1790565b91508061152881613675565b9150506111b9565b5061153d610dc4836121d5565b6115595760405162461bcd60e51b8152600401610b2090613c21565b505050505050565b8382146115fe5760405162461bcd60e51b815260206004820152606b60248201526000805160206140f083398151915260448201527f72696679426c6f6273466f7251756f72756d733a20626c6f624865616465727360648201527f20616e6420626c6f62566572696669636174696f6e50726f6f6673206c656e6760848201526a0e8d040dad2e6dac2e8c6d60ab1b60a482015260c401610b20565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa15801561163e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526116669190810190613307565b905060005b85811015611da457876001600160a01b031663eccbbfc9868684818110611694576116946135fb565b90506020028101906116a69190613ca9565b6116b4906020810190613717565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156116f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061171a9190613734565b61174f86868481811061172f5761172f6135fb565b90506020028101906117419190613ca9565b6110ba90604081019061374d565b1461176c5760405162461bcd60e51b8152600401610b2090613844565b6118a2858583818110611781576117816135fb565b90506020028101906117939190613ca9565b6117a19060608101906138b6565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508992508891508590508181106117ea576117ea6135fb565b90506020028101906117fc9190613ca9565b61180a90604081019061374d565b61181490806138fc565b356118458a8a8681811061182a5761182a6135fb565b905060200281019061183c91906138fc565b61114e90613912565b60405160200161185791815260200190565b6040516020818303038152906040528051906020012088888681811061187f5761187f6135fb565b90506020028101906118919190613ca9565b61118e906040810190602001613717565b6118be5760405162461bcd60e51b8152600401610b2090613a31565b6000805b8888848181106118d4576118d46135fb565b90506020028101906118e691906138fc565b6118f4906060810190613a93565b9050811015611d6a5788888481811061190f5761190f6135fb565b905060200281019061192191906138fc565b61192f906060810190613a93565b8281811061193f5761193f6135fb565b61195592602060809092020190810191506125a4565b60ff1687878581811061196a5761196a6135fb565b905060200281019061197c9190613ca9565b61198a90604081019061374d565b61199490806138fc565b6119a29060208101906138b6565b8989878181106119b4576119b46135fb565b90506020028101906119c69190613ca9565b6119d49060808101906138b6565b858181106119e4576119e46135fb565b919091013560f81c90508181106119fd576119fd6135fb565b9050013560f81c60f81b60f81c60ff1614611a2a5760405162461bcd60e51b8152600401610b2090613adc565b888884818110611a3c57611a3c6135fb565b9050602002810190611a4e91906138fc565b611a5c906060810190613a93565b82818110611a6c57611a6c6135fb565b9050608002016020016020810190611a8491906125a4565b60ff16898985818110611a9957611a996135fb565b9050602002810190611aab91906138fc565b611ab9906060810190613a93565b83818110611ac957611ac96135fb565b9050608002016040016020810190611ae191906125a4565b60ff1611611b015760405162461bcd60e51b8152600401610b2090613b3f565b83898985818110611b1457611b146135fb565b9050602002810190611b2691906138fc565b611b34906060810190613a93565b83818110611b4457611b446135fb565b611b5a92602060809092020190810191506125a4565b60ff1681518110611b6d57611b6d6135fb565b016020015160f81c898985818110611b8757611b876135fb565b9050602002810190611b9991906138fc565b611ba7906060810190613a93565b83818110611bb757611bb76135fb565b9050608002016040016020810190611bcf91906125a4565b60ff161015611bf05760405162461bcd60e51b8152600401610b2090613bb0565b888884818110611c0257611c026135fb565b9050602002810190611c1491906138fc565b611c22906060810190613a93565b82818110611c3257611c326135fb565b9050608002016040016020810190611c4a91906125a4565b60ff16878785818110611c5f57611c5f6135fb565b9050602002810190611c719190613ca9565b611c7f90604081019061374d565b611c8990806138fc565b611c979060408101906138b6565b898987818110611ca957611ca96135fb565b9050602002810190611cbb9190613ca9565b611cc99060808101906138b6565b85818110611cd957611cd96135fb565b919091013560f81c9050818110611cf257611cf26135fb565b9050013560f81c60f81b60f81c60ff161015611d205760405162461bcd60e51b8152600401610b2090613bb0565b611d56828a8a86818110611d3657611d366135fb565b9050602002810190611d4891906138fc565b6114e9906060810190613a93565b915080611d6281613675565b9150506118c2565b50611d77610dc4856121d5565b611d935760405162461bcd60e51b8152600401610b2090613c21565b50611d9d81613675565b905061166b565b505050505050505050565b6000611dbc878787611dde565b9050611dd28a8a8a886000015188868989610a5a565b50505050505050505050565b611de66124e9565b602082015151516000906001600160401b03811115611e0757611e0761267a565b604051908082528060200260200182016040528015611e30578160200160208202803683370190505b50905060005b60208401515151811015611ead57611e808460200151600001518281518110611e6157611e616135fb565b6020026020010151805160009081526020918201519091526040902090565b828281518110611e9257611e926135fb565b6020908102919091010152611ea681613675565b9050611e36565b50606060005b84602001516080015151811015611f1a57818560200151608001518281518110611edf57611edf6135fb565b6020026020010151604051602001611ef8929190613cbf565b604051602081830303815290604052915080611f1390613675565b9050611eb3565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491611f55918a919087908990600401613cf1565b600060405180830381865afa158015611f72573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611f9a9190810190613e44565b805185526020958601805151878701528051870151604080880191909152815160609081015181890152915181015160808801529682015160a08701529581015160c0860152949094015160e08401525090949350505050565b600081600001518260200151836040015160405160200161201793929190613f1c565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152015b604051602081830303815290604052805190602001209050919050565b60008361206d8685856123e6565b1495945050505050565b60008160405160200161204291908151815260209182015163ffffffff169181019190915260400190565b60005b81518110156104f35760006001600160a01b0316836001600160a01b031663b5a872da8484815181106120da576120da6135fb565b60200260200101516040518263ffffffff1660e01b815260040161210a919063ffffffff91909116815260200190565b602060405180830381865afa158015612127573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061214b9190613f95565b6001600160a01b031614156121c55760405162461bcd60e51b815260206004820152604660248201526000805160206140f083398151915260448201527f7269667952656c61794b6579735365743a2072656c6179206b6579206973206e6064820152651bdd081cd95d60d21b608482015260a401610b20565b6121ce81613675565b90506120a5565b60006101008251111561225e5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610b20565b815161226c57506000919050565b60008083600081518110612282576122826135fb565b0160200151600160f89190911c81901b92505b8451811015612359578481815181106122b0576122b06135fb565b0160200151600160f89190911c1b91508282116123455760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610b20565b9181179161235281613675565b9050612295565b50909392505050565b6000610451826000015160405160200161237c9190613fbe565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b600081604051602001612042919061401e565b6000602084516123f691906140c3565b1561247d5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610b20565b8260205b855181116124e0576124946002856140c3565b6124b5578160005280860151602052604060002091506002840493506124ce565b8086015160005281602052604060002091506002840493505b6124d96020826140d7565b9050612481565b50949350505050565b60405180610100016040528060608152602001606081526020016060815260200161251261254f565b8152602001612534604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060400160405280612562612574565b815260200161256f612574565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff811681146125a157600080fd5b50565b6000602082840312156125b657600080fd5b81356125c181612592565b9392505050565b6000606082840312156125da57600080fd5b50919050565b600080600083850360808112156125f657600080fd5b604081121561260457600080fd5b5083925060408401356001600160401b038082111561262257600080fd5b61262e878388016125c8565b9350606086013591508082111561264457600080fd5b508401610180818703121561265857600080fd5b809150509250925092565b803561ffff8116811461267557600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156126b2576126b261267a565b60405290565b604051606081016001600160401b03811182821017156126b2576126b261267a565b604051608081016001600160401b03811182821017156126b2576126b261267a565b60405161010081016001600160401b03811182821017156126b2576126b261267a565b60405160a081016001600160401b03811182821017156126b2576126b261267a565b604051601f8201601f191681016001600160401b03811182821017156127695761276961267a565b604052919050565b60006040828403121561278357600080fd5b61278b612690565b9050813561279881612592565b815260208201356127a881612592565b602082015292915050565b600080606083850312156127c657600080fd5b6127cf83612663565b91506127de8460208501612771565b90509250929050565b6000602082840312156127f957600080fd5b6125c182612663565b60005b8381101561281d578181015183820152602001612805565b838111156107835750506000910152565b60008151808452612846816020860160208601612802565b601f01601f19169290920160200192915050565b6020815260006125c1602083018461282e565b6000806040838503121561288057600080fd5b82356001600160401b038082111561289757600080fd5b90840190608082870312156128ab57600080fd5b909250602084013590808211156128c157600080fd5b50830160a081860312156128d457600080fd5b809150509250929050565b63ffffffff811681146125a157600080fd5b8035612675816128df565b60008082840360a081121561291057600080fd5b606081121561291e57600080fd5b506129276126b8565b8335612932816128df565b81526020840135612942816128df565b6020820152604084013561295581612592565b604082015291506127de8460608501612771565b60008083601f84011261297b57600080fd5b5081356001600160401b0381111561299257600080fd5b6020830191508360208260051b85010111156129ad57600080fd5b9250929050565b600080600080604085870312156129ca57600080fd5b84356001600160401b03808211156129e157600080fd5b6129ed88838901612969565b90965094506020870135915080821115612a0657600080fd5b50612a1387828801612969565b95989497509550505050565b60008060408385031215612a3257600080fd5b82356001600160401b0380821115612a4957600080fd5b612a55868387016125c8565b93506020850135915080821115612a6b57600080fd5b50612a78858286016125c8565b9150509250929050565b600060208284031215612a9457600080fd5b81356001600160401b03811115612aaa57600080fd5b612ab6848285016125c8565b949350505050565b600081518084526020808501945080840160005b83811015612af457815163ffffffff1687529582019590820190600101612ad2565b509495945050505050565b600081518084526020808501945080840160005b83811015612af457612b3087835180518252602090810151910152565b6040969096019590820190600101612b13565b8060005b6002811015610783578151845260209384019390910190600101612b47565b612b71828251612b43565b60208101516104f36040840182612b43565b600081518084526020808501808196508360051b8101915082860160005b85811015612bcb578284038952612bb9848351612abe565b98850198935090840190600101612ba1565b5091979650505050505050565b60006101808251818552612bee82860182612abe565b91505060208301518482036020860152612c088282612aff565b91505060408301518482036040860152612c228282612aff565b9150506060830151612c376060860182612b66565b506080830151805160e08601526020015161010085015260a0830151848203610120860152612c668282612abe565b91505060c0830151848203610140860152612c818282612abe565b91505060e0830151848203610160860152612c9c8282612b83565b95945050505050565b6020815260006125c16020830184612bd8565b600060208284031215612cca57600080fd5b815180151581146125c157600080fd5b600060408284031215612cec57600080fd5b612cf4612690565b90508135815260208201356127a8816128df565b600060408284031215612d1a57600080fd5b6125c18383612cda565b60006001600160401b03821115612d3d57612d3d61267a565b50601f01601f191660200190565b600082601f830112612d5c57600080fd5b8135612d6f612d6a82612d24565b612741565b818152846020838601011115612d8457600080fd5b816020850160208301376000918101602001919091529392505050565b600060408284031215612db357600080fd5b612dbb612690565b9050813581526020820135602082015292915050565b600082601f830112612de257600080fd5b612dea612690565b806040840185811115612dfc57600080fd5b845b81811015612e16578035845260209384019301612dfe565b509095945050505050565b600060808284031215612e3357600080fd5b612e3b612690565b9050612e478383612dd1565b81526127a88360408401612dd1565b60006001600160401b03821115612e6f57612e6f61267a565b5060051b60200190565b600082601f830112612e8a57600080fd5b81356020612e9a612d6a83612e56565b82815260059290921b84018101918181019086841115612eb957600080fd5b8286015b84811015612edd578035612ed0816128df565b8352918301918301612ebd565b509695505050505050565b600060608236031215612efa57600080fd5b612f026126b8565b82356001600160401b0380821115612f1957600080fd5b81850191506040808336031215612f2f57600080fd5b612f37612690565b833583811115612f4657600080fd5b8401368190036101c0811215612f5b57600080fd5b612f636126da565b612f6c83612663565b815260208084013587811115612f8157600080fd5b612f8d36828701612d4b565b8383015250610160603f1984011215612fa557600080fd5b612fad6126da565b9250612fbb36878601612da1565b8352612fca3660808601612e21565b81840152612fdc366101008601612e21565b86840152610180840135612fef816128df565b8060608501525082868301526101a084013560608301528185528088013593508684111561301c57600080fd5b61302836858a01612e79565b8186015284895261303a818c016128f1565b90890152505050508581013592508183111561305557600080fd5b61306136848801612d4b565b9084015250909392505050565b600082601f83011261307f57600080fd5b8135602061308f612d6a83612e56565b82815260069290921b840181019181810190868411156130ae57600080fd5b8286015b84811015612edd576130c48882612da1565b8352918301916040016130b2565b600082601f8301126130e357600080fd5b813560206130f3612d6a83612e56565b82815260059290921b8401810191818101908684111561311257600080fd5b8286015b84811015612edd5780356001600160401b038111156131355760008081fd5b6131438986838b0101612e79565b845250918301918301613116565b6000610180823603121561316457600080fd5b61316c6126fc565b82356001600160401b038082111561318357600080fd5b61318f36838701612e79565b835260208501359150808211156131a557600080fd5b6131b13683870161306e565b602084015260408501359150808211156131ca57600080fd5b6131d63683870161306e565b60408401526131e83660608701612e21565b60608401526131fa3660e08701612da1565b608084015261012085013591508082111561321457600080fd5b61322036838701612e79565b60a084015261014085013591508082111561323a57600080fd5b61324636838701612e79565b60c084015261016085013591508082111561326057600080fd5b5061326d368286016130d2565b60e08301525092915050565b60006020828403121561328b57600080fd5b81516125c181612592565b6000606082840312156132a857600080fd5b604051606081018181106001600160401b03821117156132ca576132ca61267a565b60405282516132d8816128df565b815260208301516132e8816128df565b602082015260408301516132fb81612592565b60408201529392505050565b60006020828403121561331957600080fd5b81516001600160401b0381111561332f57600080fd5b8201601f8101841361334057600080fd5b805161334e612d6a82612d24565b81815285602083850101111561336357600080fd5b612c9c826020830160208601612802565b60006060823603121561338657600080fd5b61338e612690565b6133983684612cda565b815260408301356001600160401b03808211156133b457600080fd5b818501915061012082360312156133ca57600080fd5b6133d261271f565b8235828111156133e157600080fd5b6133ed3682860161306e565b82525060208301358281111561340257600080fd5b61340e3682860161306e565b6020830152506134213660408501612da1565b60408201526134333660808501612e21565b60608201526101008301358281111561344b57600080fd5b61345736828601612e79565b608083015250602084015250909392505050565b60006040828403121561347d57600080fd5b613485612690565b825161349081612592565b815260208301516134a081612592565b60208201529392505050565b8481526080602082015260006134c5608083018661282e565b63ffffffff8516604084015282810360608401526134e38185612bd8565b979650505050505050565b600082601f8301126134ff57600080fd5b8151602061350f612d6a83612e56565b82815260059290921b8401810191818101908684111561352e57600080fd5b8286015b84811015612edd5780516001600160601b03811681146135525760008081fd5b8352918301918301613532565b6000806040838503121561357257600080fd5b82516001600160401b038082111561358957600080fd5b908401906040828703121561359d57600080fd5b6135a5612690565b8251828111156135b457600080fd5b6135c0888286016134ee565b8252506020830151828111156135d557600080fd5b6135e1888286016134ee565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006001600160601b038083168185168183048111821515161561364d5761364d613611565b02949350505050565b600081600019048311821515161561367057613670613611565b500290565b600060001982141561368957613689613611565b5060010190565b600060ff821660ff8416808210156136aa576136aa613611565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826136d8576136d86136b3565b500490565b6000828210156136ef576136ef613611565b500390565b600063ffffffff8083168185168183048111821515161561364d5761364d613611565b60006020828403121561372957600080fd5b81356125c1816128df565b60006020828403121561374657600080fd5b5051919050565b60008235605e1983360301811261376357600080fd5b9190910192915050565b60006060823603121561377f57600080fd5b6137876126b8565b82356001600160401b038082111561379e57600080fd5b8185019150608082360312156137b357600080fd5b6137bb6126da565b823581526020830135828111156137d157600080fd5b6137dd36828601612d4b565b6020830152506040830135828111156137f557600080fd5b61380136828601612d4b565b60408301525060608301359250613817836128df565b82606082015280845250505060208301356020820152613839604084016128f1565b604082015292915050565b602080825260609082018190526000805160206140f083398151915260408301527f72696679426c6f62466f7251756f72756d733a2062617463684d657461646174908201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608082015260a00190565b6000808335601e198436030181126138cd57600080fd5b8301803591506001600160401b038211156138e757600080fd5b6020019150368190038213156129ad57600080fd5b60008235607e1983360301811261376357600080fd5b6000608080833603121561392557600080fd5b61392d6126b8565b6139373685612da1565b8152604080850135613948816128df565b6020818185015260609150818701356001600160401b0381111561396b57600080fd5b870136601f82011261397c57600080fd5b803561398a612d6a82612e56565b81815260079190911b820183019083810190368311156139a957600080fd5b928401925b82841015613a1d578884360312156139c65760008081fd5b6139ce6126da565b84356139d981612592565b8152848601356139e881612592565b81870152848801356139f981612592565b8189015284870135613a0a816128df565b81880152825292880192908401906139ae565b958701959095525093979650505050505050565b6020808252604e908201526000805160206140f083398151915260408201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260608201526d1bdbd9881a5cc81a5b9d985b1a5960921b608082015260a00190565b6000808335601e19843603018112613aaa57600080fd5b8301803591506001600160401b03821115613ac457600080fd5b6020019150600781901b36038213156129ad57600080fd5b6020808252604f908201526000805160206140f083398151915260408201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260608201526e040c8decae640dcdee840dac2e8c6d608b1b608082015260a00190565b60208082526057908201526000805160206140f083398151915260408201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560608201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608082015260a00190565b6020808252605e908201526000805160206140f083398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b6020808252606e908201526000805160206140f083398151915260408201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60608201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60808201526d6669726d65642071756f72756d7360901b60a082015260c00190565b60008235609e1983360301811261376357600080fd5b60008351613cd1818460208801612802565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff86168184015260806040840152613d1f608084018661282e565b838103606085015284518082528286019183019060005b81811015613d5257835183529284019291840191600101613d36565b50909998505050505050505050565b600082601f830112613d7257600080fd5b81516020613d82612d6a83612e56565b82815260059290921b84018101918181019086841115613da157600080fd5b8286015b84811015612edd578051613db8816128df565b8352918301918301613da5565b600082601f830112613dd657600080fd5b81516020613de6612d6a83612e56565b82815260059290921b84018101918181019086841115613e0557600080fd5b8286015b84811015612edd5780516001600160401b03811115613e285760008081fd5b613e368986838b0101613d61565b845250918301918301613e09565b600060208284031215613e5657600080fd5b81516001600160401b0380821115613e6d57600080fd5b9083019060808286031215613e8157600080fd5b613e896126da565b825182811115613e9857600080fd5b613ea487828601613d61565b825250602083015182811115613eb957600080fd5b613ec587828601613d61565b602083015250604083015182811115613edd57600080fd5b613ee987828601613d61565b604083015250606083015182811115613f0157600080fd5b613f0d87828601613dc5565b60608301525095945050505050565b60006101a061ffff86168352806020840152613f3a8184018661282e565b8451805160408601526020015160608501529150613f559050565b6020830151613f676080840182612b66565b506040830151613f7b610100840182612b66565b5063ffffffff606084015116610180830152949350505050565b600060208284031215613fa757600080fd5b81516001600160a01b03811681146125c157600080fd5b60208152815160208201526000602083015160806040840152613fe460a084018261282e565b90506040840151601f19848303016060850152614001828261282e565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b808410156140b4578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190614071565b509a9950505050505050505050565b6000826140d2576140d26136b3565b500690565b600082198211156140ea576140ea613611565b50019056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a26469706673582212203f154c301d7c18e3b0787b279c2bcc5e5fc53d56092249a2942f16b3dc54bb4164736f6c634300080c0033",
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

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x0a4715f4.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2", batchHeader, blobVerificationProof, nonSignerStakesAndSignature)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x0a4715f4.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x0a4715f4.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xca7fd71d.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch", signedBatch, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xca7fd71d.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xca7fd71d.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
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
