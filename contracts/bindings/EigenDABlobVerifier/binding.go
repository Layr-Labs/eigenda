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
	NonSignerPubkeys     []BN254G1Point
	QuorumApks           []BN254G1Point
	Sigma                BN254G1Point
	ApkG2                BN254G2Point
	QuorumNumbers        []uint32
	ReferenceBlockNumber uint32
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
	BlobHeader           BlobHeaderV2
	ReferenceBlockNumber uint32
	RelayKeys            []uint32
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"securityThreshold\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple[]\",\"internalType\":\"structSecurityThresholds[]\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b5060405162004a8238038062004a82833981016040819052620000359162000078565b6001600160a01b0394851660805292841660a05290831660c052821660e0521661010052620000f8565b6001600160a01b03811681146200007557600080fd5b50565b600080600080600060a086880312156200009157600080fd5b85516200009e816200005f565b6020870151909550620000b1816200005f565b6040870151909450620000c4816200005f565b6060870151909350620000d7816200005f565b6080870151909250620000ea816200005f565b809150509295509295909350565b60805160a05160c05160e0516101005161488f620001f3600039600081816102a601526109ec01526000818161021a01526109cb01526000818161039a015281816105680152818161068e015281816106ed0152610bb701526000818161027f015281816107f301528181610847015281816108ce0152610b3e0152600081816103d401528181610424015281816104ce01528181610547015281816105fd0152818161066d015281816106cc01528181610744015281816107d201528181610826015281816108ad015281816109020152818161096201528181610a4c01528181610a9901528181610b1d0152610b96015261488f6000f3fe608060405234801561001057600080fd5b50600436106101585760003560e01c80638f3a8f32116100c3578063ee6c3bcf1161007c578063ee6c3bcf14610359578063ef6355291461036c578063efd4532b14610395578063f12afea6146103bc578063f8c66814146103cf578063fe727205146103f657600080fd5b80638f3a8f32146102f057806392ce4ab214610303578063b60e966214610316578063bafa910714610329578063e15234ff14610331578063ed4d5a471461033957600080fd5b8063579e958c11610115578063579e958c146102545780635f44b41814610267578063640f65d91461027a5780636d14a987146102a15780638687feae146102c85780638d67b909146102dd57600080fd5b8063048886d21461015d578063127af44d146101855780631429c7c21461019a5780632229cfdb146101bf5780632ecfe72b146101d25780634ca22c3f14610215575b600080fd5b61017061016b366004612914565b610409565b60405190151581526020015b60405180910390f35b610198610193366004612a3c565b61049d565b005b6101ad6101a8366004612914565b6104b3565b60405160ff909116815260200161017c565b6101986101cd366004612af5565b610542565b6101e56101e0366004612ba6565b6105c3565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff169082015260600161017c565b61023c7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161017c565b610198610262366004612bc1565b610668565b610198610275366004612c84565b6106c7565b61023c7f000000000000000000000000000000000000000000000000000000000000000081565b61023c7f000000000000000000000000000000000000000000000000000000000000000081565b6102d0610740565b60405161017c9190612df5565b6101986102eb366004612e2c565b6107cd565b6101986102fe366004612e8f565b610821565b610198610311366004612f36565b61089e565b610198610324366004612fe7565b6108a8565b6102d06108fe565b6102d061095e565b61034c610347366004613046565b6109be565b60405161017c91906131cb565b6101ad610367366004612914565b610a31565b610374610a83565b60408051825160ff908116825260209384015116928101929092520161017c565b61023c7f000000000000000000000000000000000000000000000000000000000000000081565b6101986103ca3660046132b3565b610b18565b61023c7f000000000000000000000000000000000000000000000000000000000000000081565b61019861040436600461332d565b610b91565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa158015610473573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061049791906133a1565b92915050565b6104af6104a9836105c3565b82610bf3565b5050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa15801561051e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061049791906133c3565b6105bb7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008888888861059361095e565b89896040516020016105a7939291906133e0565b604051602081830303815290604052610da9565b505050505050565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa158015610644573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104979190613408565b6106c07f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008787876106b8610a83565b61059361095e565b5050505050565b6105bb7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008888888861071861095e565b898960405160200161072c939291906133e0565b60405160208183030381529060405261112f565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156107a0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107c891908101906134a0565b905090565b6104af7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000848461081c61095e565b6114ca565b6108987f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000868661087061095e565b8787604051602001610884939291906133e0565b6040516020818303038152906040526114ca565b50505050565b6104af8282610bf3565b6108987f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000868686866108f961095e565b6119f0565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156107a0573d6000803e3d6000fd5b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156107a0573d6000803e3d6000fd5b6109c6612859565b610a2a7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610a186040870187613516565b610a256020870187613537565b61223e565b9392505050565b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf90602401610501565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa158015610af4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107c8919061357d565b6105bb7f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000088888888610b6961095e565b8989604051602001610b7d939291906133e0565b6040516020818303038152906040526119f0565b610bee7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000858585610be1610a83565b610be961095e565b610da9565b505050565b806020015160ff16816000015160ff1611610cad5760405162461bcd60e51b8152602060048201526075602482015260008051602061483a83398151915260448201527f72696679426c6f625365637572697479506172616d733a20636f6e6669726d6160648201527f74696f6e5468726573686f6c64206d7573742062652067726561746572207468608482015274185b8818591d995c9cd85c9e551a1c995cda1bdb19605a1b60a482015260c4015b60405180910390fd5b60208101518151600091610cc0916135d4565b60ff1690506000836020015163ffffffff16846040015160ff1683620f4240610ce9919061360d565b610cf3919061360d565b610cff90612710613621565b610d099190613638565b8451909150610d1a90612710613657565b63ffffffff168110156108985760405162461bcd60e51b8152602060048201526058602482015260008051602061483a83398151915260448201527f72696679426c6f625365637572697479506172616d733a20736563757269747960648201527f20617373756d7074696f6e7320617265206e6f74206d65740000000000000000608482015260a401610ca4565b610e5c610db96040860186613537565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505088359150610e169050610dfe8880613683565b610e089080613699565b610e11906137b6565b6124d2565b604051602001610e2891815260200190565b60405160208183030381529060405280519060200120876020016020810190610e51919061388a565b63ffffffff16612502565b610e785760405162461bcd60e51b8152600401610ca4906138a7565b6000806001600160a01b038816636efb4636610ea1610e9c368b90038b018b61390b565b61251a565b610eab8980613683565b610eb59080613699565b610ec3906020810190613537565b610ed360408d0160208e0161388a565b8a6040518663ffffffff1660e01b8152600401610ef4959493929190613adc565b600060405180830381865afa158015610f11573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610f399190810190613c96565b9092509050610fdb6001600160a01b038a16632ecfe72b610f5a8980613683565b610f649080613699565b610f72906020810190612ba6565b6040516001600160e01b031960e084901b16815261ffff9091166004820152602401606060405180830381865afa158015610fb1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fd59190613408565b85610bf3565b6000805b610fe98880613683565b610ff39080613699565b611001906020810190613537565b90508110156110f357856000015160ff168460200151828151811061102857611028613d32565b602002602001015161103a9190613d48565b6001600160601b031660648560000151838151811061105b5761105b613d32565b60200260200101516001600160601b03166110769190613638565b10156110945760405162461bcd60e51b8152600401610ca490613d6e565b6110df826110a28a80613683565b6110ac9080613699565b6110ba906020810190613537565b848181106110ca576110ca613d32565b600192013560f81c9190911b91909117919050565b9150806110eb81613dfe565b915050610fdf565b5061110761110085612545565b8281161490565b6111235760405162461bcd60e51b8152600401610ca490613e19565b50505050505050505050565b6111398480613683565b6111439080613699565b611151906020810190613537565b90508251146111f15760405162461bcd60e51b815260206004820152606c602482015260008051602061483a83398151915260448201527f72696679426c6f625632466f7251756f72756d733a207365637572697479546860648201527f726573686f6c6473206c656e67746820646f6573206e6f74206d61746368207160848201526b756f72756d4e756d6265727360a01b60a482015260c401610ca4565b611201610db96040860186613537565b61121d5760405162461bcd60e51b8152600401610ca4906138a7565b6000806001600160a01b038816636efb4636611241610e9c368b90038b018b61390b565b61124b8980613683565b6112559080613699565b611263906020810190613537565b61127360408d0160208e0161388a565b8a6040518663ffffffff1660e01b8152600401611294959493929190613adc565b600060405180830381865afa1580156112b1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526112d99190810190613c96565b90925090506000806001600160a01b038b16632ecfe72b6112fa8a80613683565b6113049080613699565b611312906020810190612ba6565b6040516001600160e01b031960e084901b16815261ffff9091166004820152602401606060405180830381865afa158015611351573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113759190613408565b905060005b6113848980613683565b61138e9080613699565b61139c906020810190613537565b905081101561148d576113c8828883815181106113bb576113bb613d32565b6020026020010151610bf3565b8681815181106113da576113da613d32565b60200260200101516000015160ff16856020015182815181106113ff576113ff613d32565b60200260200101516114119190613d48565b6001600160601b031660648660000151838151811061143257611432613d32565b60200260200101516001600160601b031661144d9190613638565b101561146b5760405162461bcd60e51b8152600401610ca490613d6e565b611479836110a28b80613683565b92508061148581613dfe565b91505061137a565b506114a161149a86612545565b8381161490565b6114bd5760405162461bcd60e51b8152600401610ca490613e19565b5050505050505050505050565b6001600160a01b03841663eccbbfc96114e6602085018561388a565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611528573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061154c9190613ea3565b61156a61155c6040850185613683565b61156590613ebc565b6126d2565b146115875760405162461bcd60e51b8152600401610ca490613f93565b6116306115976060840184613537565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506115d9925050506040850185613683565b6115e39080614005565b356115f56115f08761401b565b612743565b60405160200161160791815260200190565b60405160208183030381529060405280519060200120856020016020810190610e51919061388a565b61164c5760405162461bcd60e51b8152600401610ca49061413a565b6000805b61165d606086018661419c565b90508110156119c757611673606086018661419c565b8281811061168357611683613d32565b6116999260206080909202019081019150612914565b60ff166116a96040860186613683565b6116b39080614005565b6116c1906020810190613537565b6116ce6080880188613537565b858181106116de576116de613d32565b919091013560f81c90508181106116f7576116f7613d32565b9050013560f81c60f81b60f81c60ff16146117245760405162461bcd60e51b8152600401610ca4906141e5565b611731606086018661419c565b8281811061174157611741613d32565b90506080020160200160208101906117599190612914565b60ff16611769606087018761419c565b8381811061177957611779613d32565b90506080020160400160208101906117919190612914565b60ff16116117b15760405162461bcd60e51b8152600401610ca490614248565b6001600160a01b038716631429c7c26117cd606088018861419c565b848181106117dd576117dd613d32565b6117f39260206080909202019081019150612914565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa158015611831573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061185591906133c3565b60ff16611865606087018761419c565b8381811061187557611875613d32565b905060800201604001602081019061188d9190612914565b60ff1610156118ae5760405162461bcd60e51b8152600401610ca4906142b9565b6118bb606086018661419c565b828181106118cb576118cb613d32565b90506080020160400160208101906118e39190612914565b60ff166118f36040860186613683565b6118fd9080614005565b61190b906040810190613537565b6119186080880188613537565b8581811061192857611928613d32565b919091013560f81c905081811061194157611941613d32565b9050013560f81c60f81b60f81c60ff16101561196f5760405162461bcd60e51b8152600401610ca4906142b9565b6119b382611980606088018861419c565b8481811061199057611990613d32565b6119a69260206080909202019081019150612914565b600160ff919091161b1790565b9150806119bf81613dfe565b915050611650565b506119d461110083612545565b6105bb5760405162461bcd60e51b8152600401610ca49061432a565b838214611a8d5760405162461bcd60e51b815260206004820152606b602482015260008051602061483a83398151915260448201527f72696679426c6f6273466f7251756f72756d733a20626c6f624865616465727360648201527f20616e6420626c6f62566572696669636174696f6e50726f6f6673206c656e6760848201526a0e8d040dad2e6dac2e8c6d60ab1b60a482015260c401610ca4565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015611acd573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611af591908101906134a0565b905060005b8581101561223357876001600160a01b031663eccbbfc9868684818110611b2357611b23613d32565b9050602002810190611b3591906143b2565b611b4390602081019061388a565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611b85573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ba99190613ea3565b611bde868684818110611bbe57611bbe613d32565b9050602002810190611bd091906143b2565b61155c906040810190613683565b14611bfb5760405162461bcd60e51b8152600401610ca490613f93565b611d31858583818110611c1057611c10613d32565b9050602002810190611c2291906143b2565b611c30906060810190613537565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250899250889150859050818110611c7957611c79613d32565b9050602002810190611c8b91906143b2565b611c99906040810190613683565b611ca39080614005565b35611cd48a8a86818110611cb957611cb9613d32565b9050602002810190611ccb9190614005565b6115f09061401b565b604051602001611ce691815260200190565b60405160208183030381529060405280519060200120888886818110611d0e57611d0e613d32565b9050602002810190611d2091906143b2565b610e5190604081019060200161388a565b611d4d5760405162461bcd60e51b8152600401610ca49061413a565b6000805b888884818110611d6357611d63613d32565b9050602002810190611d759190614005565b611d8390606081019061419c565b90508110156121f957888884818110611d9e57611d9e613d32565b9050602002810190611db09190614005565b611dbe90606081019061419c565b82818110611dce57611dce613d32565b611de49260206080909202019081019150612914565b60ff16878785818110611df957611df9613d32565b9050602002810190611e0b91906143b2565b611e19906040810190613683565b611e239080614005565b611e31906020810190613537565b898987818110611e4357611e43613d32565b9050602002810190611e5591906143b2565b611e63906080810190613537565b85818110611e7357611e73613d32565b919091013560f81c9050818110611e8c57611e8c613d32565b9050013560f81c60f81b60f81c60ff1614611eb95760405162461bcd60e51b8152600401610ca4906141e5565b888884818110611ecb57611ecb613d32565b9050602002810190611edd9190614005565b611eeb90606081019061419c565b82818110611efb57611efb613d32565b9050608002016020016020810190611f139190612914565b60ff16898985818110611f2857611f28613d32565b9050602002810190611f3a9190614005565b611f4890606081019061419c565b83818110611f5857611f58613d32565b9050608002016040016020810190611f709190612914565b60ff1611611f905760405162461bcd60e51b8152600401610ca490614248565b83898985818110611fa357611fa3613d32565b9050602002810190611fb59190614005565b611fc390606081019061419c565b83818110611fd357611fd3613d32565b611fe99260206080909202019081019150612914565b60ff1681518110611ffc57611ffc613d32565b016020015160f81c89898581811061201657612016613d32565b90506020028101906120289190614005565b61203690606081019061419c565b8381811061204657612046613d32565b905060800201604001602081019061205e9190612914565b60ff16101561207f5760405162461bcd60e51b8152600401610ca4906142b9565b88888481811061209157612091613d32565b90506020028101906120a39190614005565b6120b190606081019061419c565b828181106120c1576120c1613d32565b90506080020160400160208101906120d99190612914565b60ff168787858181106120ee576120ee613d32565b905060200281019061210091906143b2565b61210e906040810190613683565b6121189080614005565b612126906040810190613537565b89898781811061213857612138613d32565b905060200281019061214a91906143b2565b612158906080810190613537565b8581811061216857612168613d32565b919091013560f81c905081811061218157612181613d32565b9050013560f81c60f81b60f81c60ff1610156121af5760405162461bcd60e51b8152600401610ca4906142b9565b6121e5828a8a868181106121c5576121c5613d32565b90506020028101906121d79190614005565b61198090606081019061419c565b9150806121f181613dfe565b915050611d51565b5061220661110085612545565b6122225760405162461bcd60e51b8152600401610ca49061432a565b5061222c81613dfe565b9050611afa565b505050505050505050565b612246612859565b600061225285806143c8565b90506001600160401b0381111561226b5761226b612948565b604051908082528060200260200182016040528015612294578160200160208202803683370190505b50905060005b6122a486806143c8565b9050811015612324576122f76122ba87806143c8565b838181106122ca576122ca613d32565b9050604002018036038101906122e09190614411565b805160009081526020918201519091526040902090565b82828151811061230957612309613d32565b602090810291909101015261231d81613dfe565b905061229a565b5060006001600160a01b038816634f739f74886123496101408a016101208b0161388a565b8888876040518663ffffffff1660e01b815260040161236c95949392919061442d565b600060405180830381865afa158015612389573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526123b19190810190614582565b8051845290506123c186806143c8565b808060200260200160405190810160405280939291908181526020016000905b8282101561240d576123fe60408302860136819003810190614411565b815260200190600101906123e1565b5050505060208086019290925250612427908701876143c8565b808060200260200160405190810160405280939291908181526020016000905b828210156124735761246460408302860136819003810190614411565b81526020019060010190612447565b5050505050604084015261248f3687900387016080880161465a565b60608401526124a636879003870160408801614411565b6080840152602081015160a0840152604081015160c08401526060015160e08301525095945050505050565b6000816040516020016124e59190614676565b604051602081830303815290604052805190602001209050919050565b600083612510868585612756565b1495945050505050565b6000816040516020016124e591908151815260209182015163ffffffff169181019190915260400190565b6000610100825111156125ce5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610ca4565b81516125dc57506000919050565b600080836000815181106125f2576125f2613d32565b0160200151600160f89190911c81901b92505b84518110156126c95784818151811061262057612620613d32565b0160200151600160f89190911c1b91508282116126b55760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610ca4565b918117916126c281613dfe565b9050612605565b50909392505050565b600061049782600001516040516020016126ec9190614708565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b6000816040516020016124e59190614768565b600060208451612766919061480d565b156127ed5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610ca4565b8260205b855181116128505761280460028561480d565b6128255781600052808601516020526040600020915060028404935061283e565b8086015160005281602052604060002091506002840493505b612849602082614821565b90506127f1565b50949350505050565b6040518061010001604052806060815260200160608152602001606081526020016128826128bf565b81526020016128a4604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b60405180604001604052806128d26128e4565b81526020016128df6128e4565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff8116811461291157600080fd5b50565b60006020828403121561292657600080fd5b8135610a2a81612902565b803561ffff8116811461294357600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561298057612980612948565b60405290565b604051606081016001600160401b038111828210171561298057612980612948565b604051608081016001600160401b038111828210171561298057612980612948565b604051601f8201601f191681016001600160401b03811182821017156129f2576129f2612948565b604052919050565b600060408284031215612a0c57600080fd5b612a1461295e565b90508135612a2181612902565b81526020820135612a3181612902565b602082015292915050565b60008060608385031215612a4f57600080fd5b612a5883612931565b9150612a6784602085016129fa565b90509250929050565b600060408284031215612a8257600080fd5b50919050565b600060608284031215612a8257600080fd5b60006101808284031215612a8257600080fd5b60008083601f840112612abf57600080fd5b5081356001600160401b03811115612ad657600080fd5b602083019150836020828501011115612aee57600080fd5b9250929050565b60008060008060008060e08789031215612b0e57600080fd5b612b188888612a70565b955060408701356001600160401b0380821115612b3457600080fd5b612b408a838b01612a88565b96506060890135915080821115612b5657600080fd5b612b628a838b01612a9a565b9550612b718a60808b016129fa565b945060c0890135915080821115612b8757600080fd5b50612b9489828a01612aad565b979a9699509497509295939492505050565b600060208284031215612bb857600080fd5b610a2a82612931565b600080600080600060a08688031215612bd957600080fd5b612be38787612a70565b945060408601356001600160401b0380821115612bff57600080fd5b612c0b89838a01612a88565b95506060880135915080821115612c2157600080fd5b612c2d89838a01612a9a565b94506080880135915080821115612c4357600080fd5b50612c5088828901612aad565b969995985093965092949392505050565b60006001600160401b03821115612c7a57612c7a612948565b5060051b60200190565b60008060008060008060c08789031215612c9d57600080fd5b612ca78888612a70565b95506040808801356001600160401b0380821115612cc457600080fd5b612cd08b838c01612a88565b975060608a0135915080821115612ce657600080fd5b612cf28b838c01612a9a565b965060808a0135915080821115612d0857600080fd5b818a0191508a601f830112612d1c57600080fd5b81356020612d31612d2c83612c61565b6129ca565b82815260069290921b8401810191818101908e841115612d5057600080fd5b948201945b83861015612d7657612d678f876129fa565b82529486019490820190612d55565b985050505060a08a0135925080831115612d8f57600080fd5b5050612b9489828a01612aad565b60005b83811015612db8578181015183820152602001612da0565b838111156108985750506000910152565b60008151808452612de1816020860160208601612d9d565b601f01601f19169290920160200192915050565b602081526000610a2a6020830184612dc9565b600060808284031215612a8257600080fd5b600060a08284031215612a8257600080fd5b60008060408385031215612e3f57600080fd5b82356001600160401b0380821115612e5657600080fd5b612e6286838701612e08565b93506020850135915080821115612e7857600080fd5b50612e8585828601612e1a565b9150509250929050565b60008060008060608587031215612ea557600080fd5b84356001600160401b0380821115612ebc57600080fd5b612ec888838901612e08565b95506020870135915080821115612ede57600080fd5b612eea88838901612e1a565b94506040870135915080821115612f0057600080fd5b50612f0d87828801612aad565b95989497509550505050565b63ffffffff8116811461291157600080fd5b803561294381612f19565b60008082840360a0811215612f4a57600080fd5b6060811215612f5857600080fd5b50612f61612986565b8335612f6c81612f19565b81526020840135612f7c81612f19565b60208201526040840135612f8f81612902565b60408201529150612a6784606085016129fa565b60008083601f840112612fb557600080fd5b5081356001600160401b03811115612fcc57600080fd5b6020830191508360208260051b8501011115612aee57600080fd5b60008060008060408587031215612ffd57600080fd5b84356001600160401b038082111561301457600080fd5b61302088838901612fa3565b9096509450602087013591508082111561303957600080fd5b50612f0d87828801612fa3565b6000806040838503121561305957600080fd5b82356001600160401b038082111561307057600080fd5b61307c86838701612a88565b9350602085013591508082111561309257600080fd5b5083016101c081860312156130a657600080fd5b809150509250929050565b600081518084526020808501945080840160005b838110156130e757815163ffffffff16875295820195908201906001016130c5565b509495945050505050565b600081518084526020808501945080840160005b838110156130e75761312387835180518252602090810151910152565b6040969096019590820190600101613106565b8060005b600281101561089857815184526020938401939091019060010161313a565b613164828251613136565b6020810151610bee6040840182613136565b600081518084526020808501808196508360051b8101915082860160005b858110156131be5782840389526131ac8483516130b1565b98850198935090840190600101613194565b5091979650505050505050565b60208152600082516101808060208501526131ea6101a08501836130b1565b91506020850151601f198086850301604087015261320884836130f2565b9350604087015191508086850301606087015261322584836130f2565b93506060870151915061323b6080870183613159565b608087015180516101008801526020015161012087015260a08701518685038201610140880152915061326e84836130b1565b935060c08701519150808685030161016087015261328c84836130b1565b935060e08701519150808685030183870152506132a98382613176565b9695505050505050565b600080600080600080606087890312156132cc57600080fd5b86356001600160401b03808211156132e357600080fd5b6132ef8a838b01612fa3565b9098509650602089013591508082111561330857600080fd5b6133148a838b01612fa3565b90965094506040890135915080821115612b8757600080fd5b60008060006080848603121561334257600080fd5b61334c8585612a70565b925060408401356001600160401b038082111561336857600080fd5b61337487838801612a88565b9350606086013591508082111561338a57600080fd5b5061339786828701612a9a565b9150509250925092565b6000602082840312156133b357600080fd5b81518015158114610a2a57600080fd5b6000602082840312156133d557600080fd5b8151610a2a81612902565b600084516133f2818460208901612d9d565b8201838582376000930192835250909392505050565b60006060828403121561341a57600080fd5b604051606081018181106001600160401b038211171561343c5761343c612948565b604052825161344a81612f19565b8152602083015161345a81612f19565b6020820152604083015161346d81612902565b60408201529392505050565b60006001600160401b0382111561349257613492612948565b50601f01601f191660200190565b6000602082840312156134b257600080fd5b81516001600160401b038111156134c857600080fd5b8201601f810184136134d957600080fd5b80516134e7612d2c82613479565b8181528560208385010111156134fc57600080fd5b61350d826020830160208601612d9d565b95945050505050565b6000823561013e1983360301811261352d57600080fd5b9190910192915050565b6000808335601e1984360301811261354e57600080fd5b8301803591506001600160401b0382111561356857600080fd5b602001915036819003821315612aee57600080fd5b60006040828403121561358f57600080fd5b61359761295e565b82516135a281612902565b815260208301516135b281612902565b60208201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff8416808210156135ee576135ee6135be565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b60008261361c5761361c6135f7565b500490565b600082821015613633576136336135be565b500390565b6000816000190483118215151615613652576136526135be565b500290565b600063ffffffff8083168185168183048111821515161561367a5761367a6135be565b02949350505050565b60008235605e1983360301811261352d57600080fd5b600082356101be1983360301811261352d57600080fd5b600082601f8301126136c157600080fd5b81356136cf612d2c82613479565b8181528460208386010111156136e457600080fd5b816020850160208301376000918101602001919091529392505050565b60006040828403121561371357600080fd5b61371b61295e565b9050813581526020820135602082015292915050565b600082601f83011261374257600080fd5b61374a61295e565b80604084018581111561375c57600080fd5b845b8181101561377657803584526020938401930161375e565b509095945050505050565b60006080828403121561379357600080fd5b61379b61295e565b90506137a78383613731565b8152612a318360408401613731565b60008136036101c08112156137ca57600080fd5b6137d26129a8565b6137db84612931565b815260208401356001600160401b038111156137f657600080fd5b613802368287016136b0565b602083015250610160603f198301121561381b57600080fd5b6138236129a8565b91506138323660408601613701565b82526138413660808601613781565b6020830152613854366101008601613781565b604083015261018084013561386881612f19565b806060840152508160408201526101a084013560608201528092505050919050565b60006020828403121561389c57600080fd5b8135610a2a81612f19565b602080825260509082015260008051602061483a83398151915260408201527f72696679426c6f625632466f7251756f72756d733a20696e636c7573696f6e2060608201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608082015260a00190565b60006040828403121561391d57600080fd5b61392561295e565b8235815260208301356135b281612f19565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6000808335601e1984360301811261397757600080fd5b83016020810192503590506001600160401b0381111561399657600080fd5b8060051b3603831315612aee57600080fd5b8183526000602080850194508260005b858110156130e75781356139cb81612f19565b63ffffffff16875295820195908201906001016139b8565b6000808335601e198436030181126139fa57600080fd5b83016020810192503590506001600160401b03811115613a1957600080fd5b8060061b3603831315612aee57600080fd5b81835260208301925060008160005b84811015613a6257813586526020808301359087015260409586019590910190600101613a3a565b5093949350505050565b604081833760408201600081526040808301823750600060808301525050565b81835260006020808501808196508560051b810191508460005b878110156131be578284038952613abd8288613960565b613ac88682846139a8565b9a87019a9550505090840190600101613aa6565b858152608060208201526000613af6608083018688613937565b63ffffffff851660408401528281036060840152610180613b178586613960565b828452613b2783850182846139a8565b92505050613b3860208601866139e3565b8383036020850152613b4b838284613a2b565b92505050613b5c60408601866139e3565b8383036040850152613b6f838284613a2b565b92505050613b836060830160608701613a6c565b613b9d60e0830160e0870180358252602090810135910152565b610120613bac81870187613960565b84840383860152613bbe8482846139a8565b9350505050610140613bd281870187613960565b84840383860152613be48482846139a8565b9350505050610160613bf881870187613960565b84840383860152613c0a848284613a8c565b9c9b505050505050505050505050565b600082601f830112613c2b57600080fd5b81516020613c3b612d2c83612c61565b82815260059290921b84018101918181019086841115613c5a57600080fd5b8286015b84811015613c8b5780516001600160601b0381168114613c7e5760008081fd5b8352918301918301613c5e565b509695505050505050565b60008060408385031215613ca957600080fd5b82516001600160401b0380821115613cc057600080fd5b9084019060408287031215613cd457600080fd5b613cdc61295e565b825182811115613ceb57600080fd5b613cf788828601613c1a565b825250602083015182811115613d0c57600080fd5b613d1888828601613c1a565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b038083168185168183048111821515161561367a5761367a6135be565b602080825260769082015260008051602061483a83398151915260408201527f72696679426c6f625632466f7251756f72756d733a207369676e61746f72696560608201527f7320646f206e6f74206f776e206174206c65617374207468726573686f6c642060808201527570657263656e74616765206f6620612071756f72756d60501b60a082015260c00190565b6000600019821415613e1257613e126135be565b5060010190565b602080825260709082015260008051602061483a83398151915260408201527f72696679426c6f625632466f7251756f72756d733a207265717569726564207160608201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360808201526f6f6e6669726d65642071756f72756d7360801b60a082015260c00190565b600060208284031215613eb557600080fd5b5051919050565b600060608236031215613ece57600080fd5b613ed6612986565b82356001600160401b0380821115613eed57600080fd5b818501915060808236031215613f0257600080fd5b613f0a6129a8565b82358152602083013582811115613f2057600080fd5b613f2c368286016136b0565b602083015250604083013582811115613f4457600080fd5b613f50368286016136b0565b60408301525060608301359250613f6683612f19565b82606082015280845250505060208301356020820152613f8860408401612f2b565b604082015292915050565b6020808252606090820181905260008051602061483a83398151915260408301527f72696679426c6f62466f7251756f72756d733a2062617463684d657461646174908201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608082015260a00190565b60008235607e1983360301811261352d57600080fd5b6000608080833603121561402e57600080fd5b614036612986565b6140403685613701565b815260408085013561405181612f19565b6020818185015260609150818701356001600160401b0381111561407457600080fd5b870136601f82011261408557600080fd5b8035614093612d2c82612c61565b81815260079190911b820183019083810190368311156140b257600080fd5b928401925b82841015614126578884360312156140cf5760008081fd5b6140d76129a8565b84356140e281612902565b8152848601356140f181612902565b818701528488013561410281612902565b818901528487013561411381612f19565b81880152825292880192908401906140b7565b958701959095525093979650505050505050565b6020808252604e9082015260008051602061483a83398151915260408201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260608201526d1bdbd9881a5cc81a5b9d985b1a5960921b608082015260a00190565b6000808335601e198436030181126141b357600080fd5b8301803591506001600160401b038211156141cd57600080fd5b6020019150600781901b3603821315612aee57600080fd5b6020808252604f9082015260008051602061483a83398151915260408201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260608201526e040c8decae640dcdee840dac2e8c6d608b1b608082015260a00190565b602080825260579082015260008051602061483a83398151915260408201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560608201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608082015260a00190565b6020808252605e9082015260008051602061483a83398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b6020808252606e9082015260008051602061483a83398151915260408201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60608201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60808201526d6669726d65642071756f72756d7360901b60a082015260c00190565b60008235609e1983360301811261352d57600080fd5b6000808335601e198436030181126143df57600080fd5b8301803591506001600160401b038211156143f957600080fd5b6020019150600681901b3603821315612aee57600080fd5b60006040828403121561442357600080fd5b610a2a8383613701565b60018060a01b03861681526000602063ffffffff8716818401526080604084015261445c608084018688613937565b838103606085015284518082528286019183019060005b8181101561448f57835183529284019291840191600101614473565b50909a9950505050505050505050565b600082601f8301126144b057600080fd5b815160206144c0612d2c83612c61565b82815260059290921b840181019181810190868411156144df57600080fd5b8286015b84811015613c8b5780516144f681612f19565b83529183019183016144e3565b600082601f83011261451457600080fd5b81516020614524612d2c83612c61565b82815260059290921b8401810191818101908684111561454357600080fd5b8286015b84811015613c8b5780516001600160401b038111156145665760008081fd5b6145748986838b010161449f565b845250918301918301614547565b60006020828403121561459457600080fd5b81516001600160401b03808211156145ab57600080fd5b90830190608082860312156145bf57600080fd5b6145c76129a8565b8251828111156145d657600080fd5b6145e28782860161449f565b8252506020830151828111156145f757600080fd5b6146038782860161449f565b60208301525060408301518281111561461b57600080fd5b6146278782860161449f565b60408301525060608301518281111561463f57600080fd5b61464b87828601614503565b60608301525095945050505050565b60006080828403121561466c57600080fd5b610a2a8383613781565b6020815261ffff8251166020820152600060208301516101c08060408501526146a36101e0850183612dc9565b604086015180518051606088015260200151608087015290925060208101516146cf60a0870182613159565b5060408101516146e3610120870182613159565b5060609081015163ffffffff166101a086015294909401519390920192909252919050565b6020815281516020820152600060208301516080604084015261472e60a0840182612dc9565b90506040840151601f1984830301606085015261474b8282612dc9565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b808410156147fe578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906147bb565b509a9950505050505050505050565b60008261481c5761481c6135f7565b500690565b60008219821115614834576148346135be565b50019056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a2646970667358221220ef933533d7339fe045defeff0cb8d1b5bdbeb5a2a78af4b2c897d04532eef9fb64736f6c634300080c0033",
}

// ContractEigenDABlobVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.ABI instead.
var ContractEigenDABlobVerifierABI = ContractEigenDABlobVerifierMetaData.ABI

// ContractEigenDABlobVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDABlobVerifierMetaData.Bin instead.
var ContractEigenDABlobVerifierBin = ContractEigenDABlobVerifierMetaData.Bin

// DeployContractEigenDABlobVerifier deploys a new Ethereum contract, binding an instance of ContractEigenDABlobVerifier to it.
func DeployContractEigenDABlobVerifier(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDABatchMetadataStorage common.Address, _eigenDASignatureVerifier common.Address, _operatorStateRetriever common.Address, _registryCoordinator common.Address) (common.Address, *types.Transaction, *ContractEigenDABlobVerifier, error) {
	parsed, err := ContractEigenDABlobVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDABlobVerifierBin), backend, _eigenDAThresholdRegistry, _eigenDABatchMetadataStorage, _eigenDASignatureVerifier, _operatorStateRetriever, _registryCoordinator)
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

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xed4d5a47.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32) blobHeader) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch, blobHeader BlobHeaderV2) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch, blobHeader)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xed4d5a47.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32) blobHeader) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch, blobHeader BlobHeaderV2) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDABlobVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobHeader)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xed4d5a47.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32) blobHeader) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch, blobHeader BlobHeaderV2) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDABlobVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobHeader)
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

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV10(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV10", blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV10(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlobV10 is a free data retrieval call binding the contract method 0x8f3a8f32.
//
// Solidity: function verifyBlobV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV10(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x2229cfdb.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x2229cfdb.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, additionalQuorumNumbersRequired)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x2229cfdb.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, additionalQuorumNumbersRequired)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV20(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV20", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV20(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV20(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)
}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV21(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV21", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV21(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV21(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, additionalQuorumNumbersRequired)
}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV21(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV21(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, additionalQuorumNumbersRequired)
}

// VerifyBlobV22 is a free data retrieval call binding the contract method 0xfe727205.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV22(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV22", batchHeader, blobVerificationProof, nonSignerStakesAndSignature)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV22 is a free data retrieval call binding the contract method 0xfe727205.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV22(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV22(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
}

// VerifyBlobV22 is a free data retrieval call binding the contract method 0xfe727205.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV22(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV22(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature)
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

// VerifyBlobsV10 is a free data retrieval call binding the contract method 0xf12afea6.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobsV10(opts *bind.CallOpts, blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobsV10", blobHeaders, blobVerificationProofs, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobsV10 is a free data retrieval call binding the contract method 0xf12afea6.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobsV10(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobsV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeaders, blobVerificationProofs, additionalQuorumNumbersRequired)
}

// VerifyBlobsV10 is a free data retrieval call binding the contract method 0xf12afea6.
//
// Solidity: function verifyBlobsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobsV10(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobsV10(&_ContractEigenDABlobVerifier.CallOpts, blobHeaders, blobVerificationProofs, additionalQuorumNumbersRequired)
}
