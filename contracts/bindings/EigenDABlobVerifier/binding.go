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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"securityThreshold\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple[]\",\"internalType\":\"structSecurityThresholds[]\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple[]\",\"internalType\":\"structSecurityThresholds[]\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"securityThreshold\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620054523803806200545283398101604081905262000035916200007e565b6001600160a01b0395861660805293851660a05291841660c052831660e052821661010052166101205262000112565b6001600160a01b03811681146200007b57600080fd5b50565b60008060008060008060c087890312156200009857600080fd5b8651620000a58162000065565b6020880151909650620000b88162000065565b6040880151909550620000cb8162000065565b6060880151909450620000de8162000065565b6080880151909350620000f18162000065565b60a0880151909250620001048162000065565b809150509295509295509295565b60805160a05160c05160e0516101005161012051615192620002c0600039600081816103760152818161057c0152818161082a0152818161099201528181610b3301526110c60152600081816102d70152818161055b015281816108090152818161097101528181610b1201526110a501526000818161039d0152818161053a015281816107770152818161095001528181610a1f01528181610af101528181610be801528181611084015261117b01526000818161047101528181610519015281816107560152818161092f015281816109fe01528181610ad001528181610bc701528181611063015261115a01526000818161034f01528181610d1801528181610d6c01528181610ded0152610fea0152600081816104be015281816104f801528181610612015281816106bc015281816107350152818161089e0152818161090e015281816109dd01528181610aaf01528181610ba601528181610c6901528181610cf701528181610d4b01528181610dcc01528181610e2101528181610e8101528181610ef801528181610f4501528181610fc901528181611042015261113901526151926000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80637227644311610104578063e15234ff116100a2578063f12afea611610071578063f12afea614610493578063f50bd5e7146104a6578063f8c66814146104b9578063fe727205146104e057600080fd5b8063e15234ff14610428578063ee6c3bcf14610430578063ef63552914610443578063efd4532b1461046c57600080fd5b80638f3a8f32116100de5780638f3a8f32146103e757806392ce4ab2146103fa578063b60e96621461040d578063bafa91071461042057600080fd5b806372276443146103985780638687feae146103bf5780638d67b909146103d457600080fd5b8063488ce43611610171578063588ee67a1161014b578063588ee67a146103245780635f44b41814610337578063640f65d91461034a5780636d14a9871461037157600080fd5b8063488ce436146102bf5780634ca22c3f146102d2578063579e958c1461031157600080fd5b80631429c7c2116101ad5780631429c7c2146102245780632229cfdb146102495780632e29ee191461025c5780632ecfe72b1461027c57600080fd5b806301f18c77146101d4578063048886d2146101e9578063127af44d14610211575b600080fd5b6101e76101e2366004612f11565b6104f3565b005b6101fc6101f7366004612fad565b6105f7565b60405190151581526020015b60405180910390f35b6101e761021f366004613121565b61068b565b610237610232366004612fad565b6106a1565b60405160ff9091168152602001610208565b6101e761025736600461317a565b610730565b61026f61026a36600461322b565b6107fc565b604051610208919061344e565b61028f61028a366004613461565b610864565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff1690820152606001610208565b6101e76102cd36600461347c565b610909565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610208565b6101e761031f3660046134df565b6109d8565b6101e7610332366004613616565b610aaa565b6101e76103453660046136ab565b610ba1565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6103c7610c65565b60405161020891906137a8565b6101e76103e23660046137df565b610cf2565b6101e76103f5366004613838565b610d46565b6101e76104083660046138b0565b610dbd565b6101e761041b366004613961565b610dc7565b6103c7610e1d565b6103c7610e7d565b61023761043e366004612fad565b610edd565b61044b610f2f565b60408051825160ff9081168252602093840151169281019290925201610208565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6101e76104a13660046139c0565b610fc4565b6101e76104b4366004613a3a565b61103d565b6102f97f000000000000000000000000000000000000000000000000000000000000000081565b6101e76104ee366004613aa6565b611134565b6105f17f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006105a48a613d90565b6105ad8a613e5f565b6105b5610f2f565b8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506111d592505050565b50505050565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa158015610661573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106859190613ff6565b92915050565b61069d61069783610864565b82611208565b5050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa15801561070c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106859190614018565b6107f47f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006107a5368b90038b018b614035565b6107ae8a613e5f565b6107b78a6140d0565b8989898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506113be92505050565b505050505050565b610804612e08565b6106857f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061085660408601866141f8565b61085f90614219565b61166d565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa1580156108e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106859190614225565b61069d7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006109ba88613d90565b6109c388613e5f565b6109cb610f2f565b6109d3610e7d565b6111d5565b610aa37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610a4d368a90038a018a614035565b610a5689613e5f565b610a5f896140d0565b610a67610f2f565b89898080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506113be92505050565b5050505050565b610aa37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610b5b8b613d90565b610b648b613e5f565b8a8a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061186892505050565b6107f47f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610c16368b90038b018b614035565b610c1f8a613e5f565b610c288a6140d0565b8989898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061188b92505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610cc5573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610ced9190810190614296565b905090565b61069d7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610d41610e7d565b611bcf565b6105f17f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008686610d95610e7d565b8787604051602001610da993929190614303565b604051602081830303815290604052611bcf565b61069d8282611208565b6105f17f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000086868686610e18610e7d565b612100565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610cc5573d6000803e3d6000fd5b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610cc5573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf906024016106ef565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa158015610fa0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ced919061432b565b6107f47f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000088888888611015610e7d565b898960405160200161102993929190614303565b604051602081830303815290604052612100565b610aa37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006110ee8b613d90565b6110f78b613e5f565b8a8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506111d592505050565b6111d07f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006111a936889003880188614035565b6111b287613e5f565b6111bb876140d0565b6111c3610f2f565b6111cb610e7d565b6113be565b505050565b60006111e68787876020015161166d565b90506111fc8a8a8a8860000151888689896113be565b50505050505050505050565b806020015160ff16816000015160ff16116112c25760405162461bcd60e51b8152602060048201526075602482015260008051602061513d83398151915260448201527f72696679426c6f625365637572697479506172616d733a20636f6e6669726d6160648201527f74696f6e5468726573686f6c64206d7573742062652067726561746572207468608482015274185b8818591d995c9cd85c9e551a1c995cda1bdb19605a1b60a482015260c4015b60405180910390fd5b602081015181516000916112d591614382565b60ff1690506000836020015163ffffffff16846040015160ff1683620f42406112fe91906143bb565b61130891906143bb565b611314906127106143cf565b61131e91906143e6565b845190915061132f90612710614405565b63ffffffff168110156105f15760405162461bcd60e51b8152602060048201526058602482015260008051602061513d83398151915260448201527f72696679426c6f625365637572697479506172616d733a20736563757269747960648201527f20617373756d7074696f6e7320617265206e6f74206d65740000000000000000608482015260a4016112b9565b6040840151855185515161140f9291906113d79061294e565b6040516020016113e991815260200190565b60405160208183030381529060405280519060200120876020015163ffffffff1661297e565b61142b5760405162461bcd60e51b81526004016112b990614431565b600080886001600160a01b0316636efb463661144689612996565b885151602090810151908b01516040516001600160e01b031960e086901b168152611478939291908b90600401614495565b600060405180830381865afa158015611495573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526114bd9190810190614548565b915091506114d3888760000151604001516129c1565b85515151604051632ecfe72b60e01b815261ffff909116600482015261154e906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa158015611524573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115489190614225565b85611208565b6000805b875151602001515181101561163057856000015160ff168460200151828151811061157f5761157f6145e4565b602002602001015161159191906145fa565b6001600160601b03166064856000015183815181106115b2576115b26145e4565b60200260200101516001600160601b03166115cd91906143e6565b10156115eb5760405162461bcd60e51b81526004016112b990614620565b87515160200151805161161c9184918490811061160a5761160a6145e4565b0160200151600160f89190911c1b1790565b915080611628816146b0565b915050611552565b5061164461163d85612af4565b8281161490565b6116605760405162461bcd60e51b81526004016112b9906146cb565b5050505050505050505050565b611675612e08565b8151516000906001600160401b0381111561169257611692612fe8565b6040519080825280602002602001820160405280156116bb578160200160208202803683370190505b50905060005b83515181101561173057611703846000015182815181106116e4576116e46145e4565b6020026020010151805160009081526020918201519091526040902090565b828281518110611715576117156145e4565b6020908102919091010152611729816146b0565b90506116c1565b50606060005b84608001515181101561179557818560800151828151811061175a5761175a6145e4565b6020026020010151604051602001611773929190614755565b60405160208183030381529060405291508061178e906146b0565b9050611736565b5060a08401516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f74916117cf918a919087908990600401614787565b600060405180830381865afa1580156117ec573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261181491908101906148da565b80518552855160208087019190915280870151604080880191909152606080890151818901529781015160808801529082015160a087015281015160c0860152949094015160e08401525090949350505050565b60006118798787876020015161166d565b90506111fc8a8a8a8860000151888689895b83515160200151518251146119315760405162461bcd60e51b815260206004820152606c602482015260008051602061513d83398151915260448201527f72696679426c6f625632466f7251756f72756d733a207365637572697479546860648201527f726573686f6c6473206c656e67746820646f6573206e6f74206d61746368207160848201526b756f72756d4e756d6265727360a01b60a482015260c4016112b9565b6040840151855185515161194a9291906113d79061294e565b6119665760405162461bcd60e51b81526004016112b990614431565b600080886001600160a01b0316636efb463661198189612996565b885151602090810151908b01516040516001600160e01b031960e086901b1681526119b3939291908b90600401614495565b600060405180830381865afa1580156119d0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526119f89190810190614548565b91509150611a0e888760000151604001516129c1565b85515151604051632ecfe72b60e01b815261ffff909116600482015260009081906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa158015611a60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a849190614225565b905060005b8851516020015151811015611b9157611abb82888381518110611aae57611aae6145e4565b6020026020010151611208565b868181518110611acd57611acd6145e4565b60200260200101516000015160ff1685602001518281518110611af257611af26145e4565b6020026020010151611b0491906145fa565b6001600160601b0316606486600001518381518110611b2557611b256145e4565b60200260200101516001600160601b0316611b4091906143e6565b1015611b5e5760405162461bcd60e51b81526004016112b990614620565b885151602001518051611b7d9185918490811061160a5761160a6145e4565b925080611b89816146b0565b915050611a89565b50611ba5611b9e86612af4565b8381161490565b611bc15760405162461bcd60e51b81526004016112b9906146cb565b505050505050505050505050565b6001600160a01b03841663eccbbfc9611beb60208501856149b2565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015611c2d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c5191906149cf565b611c6f611c6160408501856149e8565b611c6a906149fe565b612c81565b14611c8c5760405162461bcd60e51b81526004016112b990614ad5565b611d40611c9c6060840184614b47565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611cde9250505060408501856149e8565b611ce89080614b8d565b35611cfa611cf587614ba3565b612cf2565b604051602001611d0c91815260200190565b60405160208183030381529060405280519060200120856020016020810190611d3591906149b2565b63ffffffff1661297e565b611d5c5760405162461bcd60e51b81526004016112b990614cc2565b6000805b611d6d6060860186614d24565b90508110156120d757611d836060860186614d24565b82818110611d9357611d936145e4565b611da99260206080909202019081019150612fad565b60ff16611db960408601866149e8565b611dc39080614b8d565b611dd1906020810190614b47565b611dde6080880188614b47565b85818110611dee57611dee6145e4565b919091013560f81c9050818110611e0757611e076145e4565b9050013560f81c60f81b60f81c60ff1614611e345760405162461bcd60e51b81526004016112b990614d6d565b611e416060860186614d24565b82818110611e5157611e516145e4565b9050608002016020016020810190611e699190612fad565b60ff16611e796060870187614d24565b83818110611e8957611e896145e4565b9050608002016040016020810190611ea19190612fad565b60ff1611611ec15760405162461bcd60e51b81526004016112b990614dd0565b6001600160a01b038716631429c7c2611edd6060880188614d24565b84818110611eed57611eed6145e4565b611f039260206080909202019081019150612fad565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa158015611f41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f659190614018565b60ff16611f756060870187614d24565b83818110611f8557611f856145e4565b9050608002016040016020810190611f9d9190612fad565b60ff161015611fbe5760405162461bcd60e51b81526004016112b990614e41565b611fcb6060860186614d24565b82818110611fdb57611fdb6145e4565b9050608002016040016020810190611ff39190612fad565b60ff1661200360408601866149e8565b61200d9080614b8d565b61201b906040810190614b47565b6120286080880188614b47565b85818110612038576120386145e4565b919091013560f81c9050818110612051576120516145e4565b9050013560f81c60f81b60f81c60ff16101561207f5760405162461bcd60e51b81526004016112b990614e41565b6120c3826120906060880188614d24565b848181106120a0576120a06145e4565b6120b69260206080909202019081019150612fad565b600160ff919091161b1790565b9150806120cf816146b0565b915050611d60565b506120e461163d83612af4565b6107f45760405162461bcd60e51b81526004016112b990614eb2565b83821461219d5760405162461bcd60e51b815260206004820152606b602482015260008051602061513d83398151915260448201527f72696679426c6f6273466f7251756f72756d733a20626c6f624865616465727360648201527f20616e6420626c6f62566572696669636174696f6e50726f6f6673206c656e6760848201526a0e8d040dad2e6dac2e8c6d60ab1b60a482015260c4016112b9565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156121dd573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526122059190810190614296565b905060005b8581101561294357876001600160a01b031663eccbbfc9868684818110612233576122336145e4565b90506020028101906122459190614f3a565b6122539060208101906149b2565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015612295573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122b991906149cf565b6122ee8686848181106122ce576122ce6145e4565b90506020028101906122e09190614f3a565b611c619060408101906149e8565b1461230b5760405162461bcd60e51b81526004016112b990614ad5565b612441858583818110612320576123206145e4565b90506020028101906123329190614f3a565b612340906060810190614b47565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250899250889150859050818110612389576123896145e4565b905060200281019061239b9190614f3a565b6123a99060408101906149e8565b6123b39080614b8d565b356123e48a8a868181106123c9576123c96145e4565b90506020028101906123db9190614b8d565b611cf590614ba3565b6040516020016123f691815260200190565b6040516020818303038152906040528051906020012088888681811061241e5761241e6145e4565b90506020028101906124309190614f3a565b611d359060408101906020016149b2565b61245d5760405162461bcd60e51b81526004016112b990614cc2565b6000805b888884818110612473576124736145e4565b90506020028101906124859190614b8d565b612493906060810190614d24565b9050811015612909578888848181106124ae576124ae6145e4565b90506020028101906124c09190614b8d565b6124ce906060810190614d24565b828181106124de576124de6145e4565b6124f49260206080909202019081019150612fad565b60ff16878785818110612509576125096145e4565b905060200281019061251b9190614f3a565b6125299060408101906149e8565b6125339080614b8d565b612541906020810190614b47565b898987818110612553576125536145e4565b90506020028101906125659190614f3a565b612573906080810190614b47565b85818110612583576125836145e4565b919091013560f81c905081811061259c5761259c6145e4565b9050013560f81c60f81b60f81c60ff16146125c95760405162461bcd60e51b81526004016112b990614d6d565b8888848181106125db576125db6145e4565b90506020028101906125ed9190614b8d565b6125fb906060810190614d24565b8281811061260b5761260b6145e4565b90506080020160200160208101906126239190612fad565b60ff16898985818110612638576126386145e4565b905060200281019061264a9190614b8d565b612658906060810190614d24565b83818110612668576126686145e4565b90506080020160400160208101906126809190612fad565b60ff16116126a05760405162461bcd60e51b81526004016112b990614dd0565b838989858181106126b3576126b36145e4565b90506020028101906126c59190614b8d565b6126d3906060810190614d24565b838181106126e3576126e36145e4565b6126f99260206080909202019081019150612fad565b60ff168151811061270c5761270c6145e4565b016020015160f81c898985818110612726576127266145e4565b90506020028101906127389190614b8d565b612746906060810190614d24565b83818110612756576127566145e4565b905060800201604001602081019061276e9190612fad565b60ff16101561278f5760405162461bcd60e51b81526004016112b990614e41565b8888848181106127a1576127a16145e4565b90506020028101906127b39190614b8d565b6127c1906060810190614d24565b828181106127d1576127d16145e4565b90506080020160400160208101906127e99190612fad565b60ff168787858181106127fe576127fe6145e4565b90506020028101906128109190614f3a565b61281e9060408101906149e8565b6128289080614b8d565b612836906040810190614b47565b898987818110612848576128486145e4565b905060200281019061285a9190614f3a565b612868906080810190614b47565b85818110612878576128786145e4565b919091013560f81c9050818110612891576128916145e4565b9050013560f81c60f81b60f81c60ff1610156128bf5760405162461bcd60e51b81526004016112b990614e41565b6128f5828a8a868181106128d5576128d56145e4565b90506020028101906128e79190614b8d565b612090906060810190614d24565b915080612901816146b0565b915050612461565b5061291661163d85612af4565b6129325760405162461bcd60e51b81526004016112b990614eb2565b5061293c816146b0565b905061220a565b505050505050505050565b6000816040516020016129619190614f50565b604051602081830303815290604052805190602001209050919050565b60008361298c868585612d05565b1495945050505050565b60008160405160200161296191908151815260209182015163ffffffff169181019190915260400190565b60005b81518110156111d05760006001600160a01b0316836001600160a01b0316638050a8998484815181106129f9576129f96145e4565b60200260200101516040518263ffffffff1660e01b8152600401612a29919063ffffffff91909116815260200190565b602060405180830381865afa158015612a46573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612a6a9190614fe2565b6001600160a01b03161415612ae45760405162461bcd60e51b8152602060048201526046602482015260008051602061513d83398151915260448201527f7269667952656c61794b6579735365743a2072656c6179206b6579206973206e6064820152651bdd081cd95d60d21b608482015260a4016112b9565b612aed816146b0565b90506129c4565b600061010082511115612b7d5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016112b9565b8151612b8b57506000919050565b60008083600081518110612ba157612ba16145e4565b0160200151600160f89190911c81901b92505b8451811015612c7857848181518110612bcf57612bcf6145e4565b0160200151600160f89190911c1b9150828211612c645760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016112b9565b91811791612c71816146b0565b9050612bb4565b50909392505050565b60006106858260000151604051602001612c9b919061500b565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b600081604051602001612961919061506b565b600060208451612d159190615110565b15612d9c5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a4016112b9565b8260205b85518111612dff57612db3600285615110565b612dd457816000528086015160205260406000209150600284049350612ded565b8086015160005281602052604060002091506002840493505b612df8602082615124565b9050612da0565b50949350505050565b604051806101000160405280606081526020016060815260200160608152602001612e31612e6e565b8152602001612e53604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060400160405280612e81612e93565b8152602001612e8e612e93565b905290565b60405180604001604052806002906020820280368337509192915050565b600060608284031215612ec357600080fd5b50919050565b60008083601f840112612edb57600080fd5b5081356001600160401b03811115612ef257600080fd5b602083019150836020828501011115612f0a57600080fd5b9250929050565b60008060008060608587031215612f2757600080fd5b84356001600160401b0380821115612f3e57600080fd5b612f4a88838901612eb1565b95506020870135915080821115612f6057600080fd5b612f6c88838901612eb1565b94506040870135915080821115612f8257600080fd5b50612f8f87828801612ec9565b95989497509550505050565b60ff81168114612faa57600080fd5b50565b600060208284031215612fbf57600080fd5b8135612fca81612f9b565b9392505050565b803561ffff81168114612fe357600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561302057613020612fe8565b60405290565b604051606081016001600160401b038111828210171561302057613020612fe8565b60405160c081016001600160401b038111828210171561302057613020612fe8565b604051608081016001600160401b038111828210171561302057613020612fe8565b60405161010081016001600160401b038111828210171561302057613020612fe8565b604051601f8201601f191681016001600160401b03811182821017156130d7576130d7612fe8565b604052919050565b6000604082840312156130f157600080fd5b6130f9612ffe565b9050813561310681612f9b565b8152602082013561311681612f9b565b602082015292915050565b6000806060838503121561313457600080fd5b61313d83612fd1565b915061314c84602085016130df565b90509250929050565b600060408284031215612ec357600080fd5b60006101808284031215612ec357600080fd5b60008060008060008060e0878903121561319357600080fd5b61319d8888613155565b955060408701356001600160401b03808211156131b957600080fd5b6131c58a838b01612eb1565b965060608901359150808211156131db57600080fd5b6131e78a838b01613167565b95506131f68a60808b016130df565b945060c089013591508082111561320c57600080fd5b5061321989828a01612ec9565b979a9699509497509295939492505050565b60006020828403121561323d57600080fd5b81356001600160401b0381111561325357600080fd5b61325f84828501612eb1565b949350505050565b600081518084526020808501945080840160005b8381101561329d57815163ffffffff168752958201959082019060010161327b565b509495945050505050565b600081518084526020808501945080840160005b8381101561329d576132d987835180518252602090810151910152565b60409690960195908201906001016132bc565b8060005b60028110156105f15781518452602093840193909101906001016132f0565b61331a8282516132ec565b60208101516111d060408401826132ec565b600081518084526020808501808196508360051b8101915082860160005b85811015613374578284038952613362848351613267565b9885019893509084019060010161334a565b5091979650505050505050565b6000610180825181855261339782860182613267565b915050602083015184820360208601526133b182826132a8565b915050604083015184820360408601526133cb82826132a8565b91505060608301516133e0606086018261330f565b506080830151805160e08601526020015161010085015260a083015184820361012086015261340f8282613267565b91505060c083015184820361014086015261342a8282613267565b91505060e0830151848203610160860152613445828261332c565b95945050505050565b602081526000612fca6020830184613381565b60006020828403121561347357600080fd5b612fca82612fd1565b6000806040838503121561348f57600080fd5b82356001600160401b03808211156134a657600080fd5b6134b286838701612eb1565b935060208501359150808211156134c857600080fd5b506134d585828601612eb1565b9150509250929050565b600080600080600060a086880312156134f757600080fd5b6135018787613155565b945060408601356001600160401b038082111561351d57600080fd5b61352989838a01612eb1565b9550606088013591508082111561353f57600080fd5b61354b89838a01613167565b9450608088013591508082111561356157600080fd5b5061356e88828901612ec9565b969995985093965092949392505050565b60006001600160401b0382111561359857613598612fe8565b5060051b60200190565b600082601f8301126135b357600080fd5b813560206135c86135c38361357f565b6130af565b82815260069290921b840181019181810190868411156135e757600080fd5b8286015b8481101561360b576135fd88826130df565b8352918301916040016135eb565b509695505050505050565b60008060008060006080868803121561362e57600080fd5b85356001600160401b038082111561364557600080fd5b61365189838a01612eb1565b9650602088013591508082111561366757600080fd5b61367389838a01612eb1565b9550604088013591508082111561368957600080fd5b61369589838a016135a2565b9450606088013591508082111561356157600080fd5b60008060008060008060c087890312156136c457600080fd5b6136ce8888613155565b955060408701356001600160401b03808211156136ea57600080fd5b6136f68a838b01612eb1565b9650606089013591508082111561370c57600080fd5b6137188a838b01613167565b9550608089013591508082111561372e57600080fd5b61373a8a838b016135a2565b945060a089013591508082111561320c57600080fd5b60005b8381101561376b578181015183820152602001613753565b838111156105f15750506000910152565b60008151808452613794816020860160208601613750565b601f01601f19169290920160200192915050565b602081526000612fca602083018461377c565b600060808284031215612ec357600080fd5b600060a08284031215612ec357600080fd5b600080604083850312156137f257600080fd5b82356001600160401b038082111561380957600080fd5b613815868387016137bb565b9350602085013591508082111561382b57600080fd5b506134d5858286016137cd565b6000806000806060858703121561384e57600080fd5b84356001600160401b038082111561386557600080fd5b613871888389016137bb565b9550602087013591508082111561388757600080fd5b612f6c888389016137cd565b63ffffffff81168114612faa57600080fd5b8035612fe381613893565b60008082840360a08112156138c457600080fd5b60608112156138d257600080fd5b506138db613026565b83356138e681613893565b815260208401356138f681613893565b6020820152604084013561390981612f9b565b6040820152915061314c84606085016130df565b60008083601f84011261392f57600080fd5b5081356001600160401b0381111561394657600080fd5b6020830191508360208260051b8501011115612f0a57600080fd5b6000806000806040858703121561397757600080fd5b84356001600160401b038082111561398e57600080fd5b61399a8883890161391d565b909650945060208701359150808211156139b357600080fd5b50612f8f8782880161391d565b600080600080600080606087890312156139d957600080fd5b86356001600160401b03808211156139f057600080fd5b6139fc8a838b0161391d565b90985096506020890135915080821115613a1557600080fd5b613a218a838b0161391d565b9096509450604089013591508082111561320c57600080fd5b600080600080600060a08688031215613a5257600080fd5b85356001600160401b0380821115613a6957600080fd5b613a7589838a01612eb1565b96506020880135915080821115613a8b57600080fd5b613a9789838a01612eb1565b955061354b8960408a016130df565b600080600060808486031215613abb57600080fd5b613ac58585613155565b925060408401356001600160401b0380821115613ae157600080fd5b613aed87838801612eb1565b93506060860135915080821115613b0357600080fd5b50613b1086828701613167565b9150509250925092565b600060408284031215613b2c57600080fd5b613b34612ffe565b905081358152602082013561311681613893565b600060408284031215613b5a57600080fd5b613b62612ffe565b9050813581526020820135602082015292915050565b600082601f830112613b8957600080fd5b81356020613b996135c38361357f565b82815260069290921b84018101918181019086841115613bb857600080fd5b8286015b8481101561360b57613bce8882613b48565b835291830191604001613bbc565b600082601f830112613bed57600080fd5b613bf5612ffe565b806040840185811115613c0757600080fd5b845b81811015613c21578035845260209384019301613c09565b509095945050505050565b600060808284031215613c3e57600080fd5b613c46612ffe565b9050613c528383613bdc565b81526131168360408401613bdc565b600082601f830112613c7257600080fd5b81356020613c826135c38361357f565b82815260059290921b84018101918181019086841115613ca157600080fd5b8286015b8481101561360b578035613cb881613893565b8352918301918301613ca5565b60006101408284031215613cd857600080fd5b613ce0613048565b905081356001600160401b0380821115613cf957600080fd5b613d0585838601613b78565b83526020840135915080821115613d1b57600080fd5b613d2785838601613b78565b6020840152613d398560408601613b48565b6040840152613d4b8560808601613c2c565b6060840152610100840135915080821115613d6557600080fd5b50613d7284828501613c61565b608083015250613d8561012083016138a5565b60a082015292915050565b600060608236031215613da257600080fd5b613daa612ffe565b613db43684613b1a565b815260408301356001600160401b03811115613dcf57600080fd5b613ddb36828601613cc5565b60208301525092915050565b60006001600160401b03821115613e0057613e00612fe8565b50601f01601f191660200190565b600082601f830112613e1f57600080fd5b8135613e2d6135c382613de7565b818152846020838601011115613e4257600080fd5b816020850160208301376000918101602001919091529392505050565b600060608236031215613e7157600080fd5b613e79613026565b82356001600160401b0380821115613e9057600080fd5b818501915060608236031215613ea557600080fd5b613ead613026565b823582811115613ebc57600080fd5b8301368190036101c0811215613ed157600080fd5b613ed961306a565b613ee283612fd1565b815260208084013586811115613ef757600080fd5b613f0336828701613e0e565b83830152506040610160603f1985011215613f1d57600080fd5b613f2561306a565b9350613f3336828701613b48565b8452613f423660808701613c2c565b82850152613f54366101008701613c2c565b81850152610180850135613f6781613893565b8060608601525083818401526101a08501356060840152828652613f8c8289016138a5565b8287015280880135945086851115613fa357600080fd5b613faf36868a01613c61565b81870152858952613fc1828c016138a5565b828a0152808b0135975086881115613fd857600080fd5b613fe436898d01613e0e565b90890152509598975050505050505050565b60006020828403121561400857600080fd5b81518015158114612fca57600080fd5b60006020828403121561402a57600080fd5b8151612fca81612f9b565b60006040828403121561404757600080fd5b612fca8383613b1a565b600082601f83011261406257600080fd5b813560206140726135c38361357f565b82815260059290921b8401810191818101908684111561409157600080fd5b8286015b8481101561360b5780356001600160401b038111156140b45760008081fd5b6140c28986838b0101613c61565b845250918301918301614095565b600061018082360312156140e357600080fd5b6140eb61308c565b82356001600160401b038082111561410257600080fd5b61410e36838701613c61565b8352602085013591508082111561412457600080fd5b61413036838701613b78565b6020840152604085013591508082111561414957600080fd5b61415536838701613b78565b60408401526141673660608701613c2c565b60608401526141793660e08701613b48565b608084015261012085013591508082111561419357600080fd5b61419f36838701613c61565b60a08401526101408501359150808211156141b957600080fd5b6141c536838701613c61565b60c08401526101608501359150808211156141df57600080fd5b506141ec36828601614051565b60e08301525092915050565b6000823561013e1983360301811261420f57600080fd5b9190910192915050565b60006106853683613cc5565b60006060828403121561423757600080fd5b604051606081018181106001600160401b038211171561425957614259612fe8565b604052825161426781613893565b8152602083015161427781613893565b6020820152604083015161428a81612f9b565b60408201529392505050565b6000602082840312156142a857600080fd5b81516001600160401b038111156142be57600080fd5b8201601f810184136142cf57600080fd5b80516142dd6135c382613de7565b8181528560208385010111156142f257600080fd5b613445826020830160208601613750565b60008451614315818460208901613750565b8201838582376000930192835250909392505050565b60006040828403121561433d57600080fd5b614345612ffe565b825161435081612f9b565b8152602083015161436081612f9b565b60208201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff84168082101561439c5761439c61436c565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826143ca576143ca6143a5565b500490565b6000828210156143e1576143e161436c565b500390565b60008160001904831182151516156144005761440061436c565b500290565b600063ffffffff808316818516818304811182151516156144285761442861436c565b02949350505050565b602080825260509082015260008051602061513d83398151915260408201527f72696679426c6f625632466f7251756f72756d733a20696e636c7573696f6e2060608201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608082015260a00190565b8481526080602082015260006144ae608083018661377c565b63ffffffff8516604084015282810360608401526144cc8185613381565b979650505050505050565b600082601f8301126144e857600080fd5b815160206144f86135c38361357f565b82815260059290921b8401810191818101908684111561451757600080fd5b8286015b8481101561360b5780516001600160601b038116811461453b5760008081fd5b835291830191830161451b565b6000806040838503121561455b57600080fd5b82516001600160401b038082111561457257600080fd5b908401906040828703121561458657600080fd5b61458e612ffe565b82518281111561459d57600080fd5b6145a9888286016144d7565b8252506020830151828111156145be57600080fd5b6145ca888286016144d7565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b60006001600160601b03808316818516818304811182151516156144285761442861436c565b602080825260769082015260008051602061513d83398151915260408201527f72696679426c6f625632466f7251756f72756d733a207369676e61746f72696560608201527f7320646f206e6f74206f776e206174206c65617374207468726573686f6c642060808201527570657263656e74616765206f6620612071756f72756d60501b60a082015260c00190565b60006000198214156146c4576146c461436c565b5060010190565b602080825260709082015260008051602061513d83398151915260408201527f72696679426c6f625632466f7251756f72756d733a207265717569726564207160608201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360808201526f6f6e6669726d65642071756f72756d7360801b60a082015260c00190565b60008351614767818460208801613750565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526147b5608084018661377c565b838103606085015284518082528286019183019060005b818110156147e8578351835292840192918401916001016147cc565b50909998505050505050505050565b600082601f83011261480857600080fd5b815160206148186135c38361357f565b82815260059290921b8401810191818101908684111561483757600080fd5b8286015b8481101561360b57805161484e81613893565b835291830191830161483b565b600082601f83011261486c57600080fd5b8151602061487c6135c38361357f565b82815260059290921b8401810191818101908684111561489b57600080fd5b8286015b8481101561360b5780516001600160401b038111156148be5760008081fd5b6148cc8986838b01016147f7565b84525091830191830161489f565b6000602082840312156148ec57600080fd5b81516001600160401b038082111561490357600080fd5b908301906080828603121561491757600080fd5b61491f61306a565b82518281111561492e57600080fd5b61493a878286016147f7565b82525060208301518281111561494f57600080fd5b61495b878286016147f7565b60208301525060408301518281111561497357600080fd5b61497f878286016147f7565b60408301525060608301518281111561499757600080fd5b6149a38782860161485b565b60608301525095945050505050565b6000602082840312156149c457600080fd5b8135612fca81613893565b6000602082840312156149e157600080fd5b5051919050565b60008235605e1983360301811261420f57600080fd5b600060608236031215614a1057600080fd5b614a18613026565b82356001600160401b0380821115614a2f57600080fd5b818501915060808236031215614a4457600080fd5b614a4c61306a565b82358152602083013582811115614a6257600080fd5b614a6e36828601613e0e565b602083015250604083013582811115614a8657600080fd5b614a9236828601613e0e565b60408301525060608301359250614aa883613893565b82606082015280845250505060208301356020820152614aca604084016138a5565b604082015292915050565b6020808252606090820181905260008051602061513d83398151915260408301527f72696679426c6f62466f7251756f72756d733a2062617463684d657461646174908201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608082015260a00190565b6000808335601e19843603018112614b5e57600080fd5b8301803591506001600160401b03821115614b7857600080fd5b602001915036819003821315612f0a57600080fd5b60008235607e1983360301811261420f57600080fd5b60006080808336031215614bb657600080fd5b614bbe613026565b614bc83685613b48565b8152604080850135614bd981613893565b6020818185015260609150818701356001600160401b03811115614bfc57600080fd5b870136601f820112614c0d57600080fd5b8035614c1b6135c38261357f565b81815260079190911b82018301908381019036831115614c3a57600080fd5b928401925b82841015614cae57888436031215614c575760008081fd5b614c5f61306a565b8435614c6a81612f9b565b815284860135614c7981612f9b565b8187015284880135614c8a81612f9b565b8189015284870135614c9b81613893565b8188015282529288019290840190614c3f565b958701959095525093979650505050505050565b6020808252604e9082015260008051602061513d83398151915260408201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260608201526d1bdbd9881a5cc81a5b9d985b1a5960921b608082015260a00190565b6000808335601e19843603018112614d3b57600080fd5b8301803591506001600160401b03821115614d5557600080fd5b6020019150600781901b3603821315612f0a57600080fd5b6020808252604f9082015260008051602061513d83398151915260408201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260608201526e040c8decae640dcdee840dac2e8c6d608b1b608082015260a00190565b602080825260579082015260008051602061513d83398151915260408201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560608201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608082015260a00190565b6020808252605e9082015260008051602061513d83398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b6020808252606e9082015260008051602061513d83398151915260408201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60608201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60808201526d6669726d65642071756f72756d7360901b60a082015260c00190565b60008235609e1983360301811261420f57600080fd5b6020815261ffff8251166020820152600060208301516101c0806040850152614f7d6101e085018361377c565b60408601518051805160608801526020015160808701529092506020810151614fa960a087018261330f565b506040810151614fbd61012087018261330f565b5060609081015163ffffffff166101a086015294909401519390920192909252919050565b600060208284031215614ff457600080fd5b81516001600160a01b0381168114612fca57600080fd5b6020815281516020820152600060208301516080604084015261503160a084018261377c565b90506040840151601f1984830301606085015261504e828261377c565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b80841015615101578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906150be565b509a9950505050505050505050565b60008261511f5761511f6143a5565b500690565b600082198211156151375761513761436c565b50019056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a2646970667358221220dab4fba935bc36f29d71fee6373b3bec913d05f8c097864b11e64ba3889fd8be64736f6c634300080c0033",
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

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0x2e29ee19.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0x2e29ee19.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDABlobVerifier.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDABlobVerifier.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0x2e29ee19.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
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
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x2229cfdb.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, quorumNumbersRequired)
}

// VerifyBlobV2 is a free data retrieval call binding the contract method 0x2229cfdb.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThreshold, quorumNumbersRequired)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV20(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV20", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV20(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, quorumNumbersRequired)
}

// VerifyBlobV20 is a free data retrieval call binding the contract method 0x579e958c.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV20(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV20(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, quorumNumbersRequired)
}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV21(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV21", batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV21(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV21(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, quorumNumbersRequired)
}

// VerifyBlobV21 is a free data retrieval call binding the contract method 0x5f44b418.
//
// Solidity: function verifyBlobV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV21(batchHeader BatchHeaderV2, blobVerificationProof BlobVerificationProofV2, nonSignerStakesAndSignature NonSignerStakesAndSignature, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV21(&_ContractEigenDABlobVerifier.CallOpts, batchHeader, blobVerificationProof, nonSignerStakesAndSignature, securityThresholds, quorumNumbersRequired)
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

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x01f18c77.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch", signedBatch, blobVerificationProof, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x01f18c77.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, quorumNumbersRequired)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0x01f18c77.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, quorumNumbersRequired)
}

// VerifyBlobV2FromSignedBatch0 is a free data retrieval call binding the contract method 0x488ce436.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch0(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch0", signedBatch, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch0 is a free data retrieval call binding the contract method 0x488ce436.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch0(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch0(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobV2FromSignedBatch0 is a free data retrieval call binding the contract method 0x488ce436.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2FromSignedBatch0(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch0(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobV2FromSignedBatch1 is a free data retrieval call binding the contract method 0x588ee67a.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch1(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch1", signedBatch, blobVerificationProof, securityThresholds, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch1 is a free data retrieval call binding the contract method 0x588ee67a.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch1(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch1(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, securityThresholds, quorumNumbersRequired)
}

// VerifyBlobV2FromSignedBatch1 is a free data retrieval call binding the contract method 0x588ee67a.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8)[] securityThresholds, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2FromSignedBatch1(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThresholds []SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch1(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, securityThresholds, quorumNumbersRequired)
}

// VerifyBlobV2FromSignedBatch2 is a free data retrieval call binding the contract method 0xf50bd5e7.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch2(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch2", signedBatch, blobVerificationProof, securityThreshold, quorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch2 is a free data retrieval call binding the contract method 0xf50bd5e7.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch2(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch2(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, securityThreshold, quorumNumbersRequired)
}

// VerifyBlobV2FromSignedBatch2 is a free data retrieval call binding the contract method 0xf50bd5e7.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32,uint32[]),uint32,bytes) blobVerificationProof, (uint8,uint8) securityThreshold, bytes quorumNumbersRequired) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCallerSession) VerifyBlobV2FromSignedBatch2(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2, securityThreshold SecurityThresholds, quorumNumbersRequired []byte) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch2(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof, securityThreshold, quorumNumbersRequired)
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
