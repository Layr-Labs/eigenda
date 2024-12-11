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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetriever\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProofV2\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlobsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6101406040523480156200001257600080fd5b50604051620043853803806200438583398101604081905262000035916200007e565b6001600160a01b0395861660805293851660a05291841660c052831660e052821661010052166101205262000112565b6001600160a01b03811681146200007b57600080fd5b50565b60008060008060008060c087890312156200009857600080fd5b8651620000a58162000065565b6020880151909650620000b88162000065565b6040880151909550620000cb8162000065565b6060880151909450620000de8162000065565b6080880151909350620000f18162000065565b60a0880151909250620001048162000065565b809150509295509295509295565b60805160a05160c05160e0516101005161012051614189620001fc6000396000818161027f015281816105cb01526108da015260008181610219015281816105aa01526108b90152600081816102a60152818161049e015261089801526000818161037a0152818161047d01526108770152600081816102580152818161075d01526107bb0152600081816103a1015281816103de0152818161045c015281816105290152818161063f015281816106ae0152818161073c0152818161079a015281816107f501528181610856015281816109240152818161099b01526109e801526141896000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c80638687feae116100b8578063cc392c251161007c578063cc392c251461031e578063e15234ff14610331578063ee6c3bcf14610339578063ef6355291461034c578063efd4532b14610375578063f8c668141461039c57600080fd5b80638687feae146102c85780638d67b909146102dd57806392ce4ab2146102f0578063b60e966214610303578063bafa91071461031657600080fd5b80632ecfe72b116100ff5780632ecfe72b146101d15780634ca22c3f14610214578063640f65d9146102535780636d14a9871461027a57806372276443146102a157600080fd5b8063048886d21461013c5780630a4715f414610164578063127af44d146101795780631429c7c21461018c5780632e29ee19146101b1575b600080fd5b61014f61014a36600461259a565b6103c3565b60405190151581526020015b60405180910390f35b6101776101723660046125d6565b610457565b005b6101776101873660046127a9565b6104f8565b61019f61019a36600461259a565b61050e565b60405160ff909116815260200161015b565b6101c46101bf3660046127dd565b61059d565b60405161015b9190612a00565b6101e46101df366004612a13565b610605565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff169082015260600161015b565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161015b565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b6102d06106aa565b60405161015b9190612a86565b6101776102eb366004612a99565b610737565b6101776102fe366004612b28565b61078b565b610177610311366004612be0565b610795565b6102d06107f1565b61017761032c366004612c4b565b610851565b6102d0610920565b61019f61034736600461259a565b610980565b6103546109d2565b60408051825160ff908116825260209384015116928101929092520161015b565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b61023b7f000000000000000000000000000000000000000000000000000000000000000081565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa15801561042d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190612cae565b92915050565b6104f37f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006104cc36889003880188612cfe565b6104d587612ede565b6104de87613147565b6104e66109d2565b6104ee610920565b610a67565b505050565b61050a61050483610605565b82610e84565b5050565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa158015610579573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610451919061326f565b6105a56124df565b6104517f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006105f7604086018661328c565b61060090613378565b611035565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa158015610686573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190613384565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561070a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261073291908101906133f5565b905090565b61050a7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000008484610786610920565b611230565b61050a8282610e84565b6107eb7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000868686866107e6610920565b611769565b50505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa15801561070a573d6000803e3d6000fd5b61050a7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061090288613462565b61090b88612ede565b6109136109d2565b61091b610920565b611fb7565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561070a573d6000803e3d6000fd5b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf9060240161055c565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa158015610a43573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073291906134b9565b60408401518551855151610ab8929190610a8090611fea565b604051602001610a9291815260200190565b60405160208183030381529060405280519060200120876020015163ffffffff16612055565b610b365760405162461bcd60e51b8152602060048201526050602482015260008051602061413483398151915260448201527f72696679426c6f625632466f7251756f72756d733a20696e636c7573696f6e2060648201526f1c1c9bdbd9881a5cc81a5b9d985b1a5960821b608482015260a4015b60405180910390fd5b600080886001600160a01b0316636efb4636610b518961206d565b885151602090810151908b01516040516001600160e01b031960e086901b168152610b83939291908b906004016134fa565b600060405180830381865afa158015610ba0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610bc891908101906135ad565b91509150610bde88876000015160200151612098565b85515151604051632ecfe72b60e01b815261ffff9091166004820152610c59906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa158015610c2f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c539190613384565b85610e84565b6000805b8751516020015151811015610dc457856000015160ff1684602001518281518110610c8a57610c8a613649565b6020026020010151610c9c9190613675565b6001600160601b0316606485600001518381518110610cbd57610cbd613649565b60200260200101516001600160601b0316610cd891906136a4565b1015610d7f5760405162461bcd60e51b8152602060048201526076602482015260008051602061413483398151915260448201527f72696679426c6f625632466f7251756f72756d733a207369676e61746f72696560648201527f7320646f206e6f74206f776e206174206c65617374207468726573686f6c642060848201527570657263656e74616765206f6620612071756f72756d60501b60a482015260c401610b2d565b875151602001518051610db091849184908110610d9e57610d9e613649565b0160200151600160f89190911c1b1790565b915080610dbc816136c3565b915050610c5d565b50610dd8610dd1856121cb565b8281161490565b610e775760405162461bcd60e51b8152602060048201526070602482015260008051602061413483398151915260448201527f72696679426c6f625632466f7251756f72756d733a207265717569726564207160648201527f756f72756d7320617265206e6f74206120737562736574206f6620746865206360848201526f6f6e6669726d65642071756f72756d7360801b60a482015260c401610b2d565b5050505050505050505050565b806020015160ff16816000015160ff1611610f395760405162461bcd60e51b8152602060048201526075602482015260008051602061413483398151915260448201527f72696679426c6f625365637572697479506172616d733a20636f6e6669726d6160648201527f74696f6e5468726573686f6c64206d7573742062652067726561746572207468608482015274185b8818591d995c9cd85c9e551a1c995cda1bdb19605a1b60a482015260c401610b2d565b60208101518151600091610f4c916136de565b60ff1690506000836020015163ffffffff16846040015160ff1683620f4240610f759190613717565b610f7f9190613717565b610f8b9061271061372b565b610f9591906136a4565b8451909150610fa690612710613742565b63ffffffff168110156107eb5760405162461bcd60e51b8152602060048201526058602482015260008051602061413483398151915260448201527f72696679426c6f625365637572697479506172616d733a20736563757269747960648201527f20617373756d7074696f6e7320617265206e6f74206d65740000000000000000608482015260a401610b2d565b61103d6124df565b8151516000906001600160401b0381111561105a5761105a612670565b604051908082528060200260200182016040528015611083578160200160208202803683370190505b50905060005b8351518110156110f8576110cb846000015182815181106110ac576110ac613649565b6020026020010151805160009081526020918201519091526040902090565b8282815181106110dd576110dd613649565b60209081029190910101526110f1816136c3565b9050611089565b50606060005b84608001515181101561115d57818560800151828151811061112257611122613649565b602002602001015160405160200161113b929190613765565b604051602081830303815290604052915080611156906136c3565b90506110fe565b5060a08401516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491611197918a919087908990600401613797565b600060405180830381865afa1580156111b4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526111dc91908101906138ea565b80518552855160208087019190915280870151604080880191909152606080890151818901529781015160808801529082015160a087015281015160c0860152949094015160e08401525090949350505050565b6001600160a01b03841663eccbbfc961124c60208501856139c2565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa15801561128e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b291906139df565b6112d06112c260408501856139f8565b6112cb90613a0e565b612358565b146112ed5760405162461bcd60e51b8152600401610b2d90613ae5565b6113a16112fd6060840184613b57565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061133f9250505060408501856139f8565b6113499080613b9d565b3561135b61135687613bb3565b6123c9565b60405160200161136d91815260200190565b6040516020818303038152906040528051906020012085602001602081019061139691906139c2565b63ffffffff16612055565b6113bd5760405162461bcd60e51b8152600401610b2d90613cd2565b6000805b6113ce6060860186613d34565b9050811015611738576113e46060860186613d34565b828181106113f4576113f4613649565b61140a926020608090920201908101915061259a565b60ff1661141a60408601866139f8565b6114249080613b9d565b611432906020810190613b57565b61143f6080880188613b57565b8581811061144f5761144f613649565b919091013560f81c905081811061146857611468613649565b9050013560f81c60f81b60f81c60ff16146114955760405162461bcd60e51b8152600401610b2d90613d7d565b6114a26060860186613d34565b828181106114b2576114b2613649565b90506080020160200160208101906114ca919061259a565b60ff166114da6060870187613d34565b838181106114ea576114ea613649565b9050608002016040016020810190611502919061259a565b60ff16116115225760405162461bcd60e51b8152600401610b2d90613de0565b6001600160a01b038716631429c7c261153e6060880188613d34565b8481811061154e5761154e613649565b611564926020608090920201908101915061259a565b6040516001600160e01b031960e084901b16815260ff9091166004820152602401602060405180830381865afa1580156115a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115c6919061326f565b60ff166115d66060870187613d34565b838181106115e6576115e6613649565b90506080020160400160208101906115fe919061259a565b60ff16101561161f5760405162461bcd60e51b8152600401610b2d90613e51565b61162c6060860186613d34565b8281811061163c5761163c613649565b9050608002016040016020810190611654919061259a565b60ff1661166460408601866139f8565b61166e9080613b9d565b61167c906040810190613b57565b6116896080880188613b57565b8581811061169957611699613649565b919091013560f81c90508181106116b2576116b2613649565b9050013560f81c60f81b60f81c60ff1610156116e05760405162461bcd60e51b8152600401610b2d90613e51565b611724826116f16060880188613d34565b8481811061170157611701613649565b611717926020608090920201908101915061259a565b600160ff919091161b1790565b915080611730816136c3565b9150506113c1565b50611745610dd1836121cb565b6117615760405162461bcd60e51b8152600401610b2d90613ec2565b505050505050565b8382146118065760405162461bcd60e51b815260206004820152606b602482015260008051602061413483398151915260448201527f72696679426c6f6273466f7251756f72756d733a20626c6f624865616465727360648201527f20616e6420626c6f62566572696669636174696f6e50726f6f6673206c656e6760848201526a0e8d040dad2e6dac2e8c6d60ab1b60a482015260c401610b2d565b6000876001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015611846573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261186e91908101906133f5565b905060005b85811015611fac57876001600160a01b031663eccbbfc986868481811061189c5761189c613649565b90506020028101906118ae9190613f4a565b6118bc9060208101906139c2565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156118fe573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061192291906139df565b61195786868481811061193757611937613649565b90506020028101906119499190613f4a565b6112c29060408101906139f8565b146119745760405162461bcd60e51b8152600401610b2d90613ae5565b611aaa85858381811061198957611989613649565b905060200281019061199b9190613f4a565b6119a9906060810190613b57565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508992508891508590508181106119f2576119f2613649565b9050602002810190611a049190613f4a565b611a129060408101906139f8565b611a1c9080613b9d565b35611a4d8a8a86818110611a3257611a32613649565b9050602002810190611a449190613b9d565b61135690613bb3565b604051602001611a5f91815260200190565b60405160208183030381529060405280519060200120888886818110611a8757611a87613649565b9050602002810190611a999190613f4a565b6113969060408101906020016139c2565b611ac65760405162461bcd60e51b8152600401610b2d90613cd2565b6000805b888884818110611adc57611adc613649565b9050602002810190611aee9190613b9d565b611afc906060810190613d34565b9050811015611f7257888884818110611b1757611b17613649565b9050602002810190611b299190613b9d565b611b37906060810190613d34565b82818110611b4757611b47613649565b611b5d926020608090920201908101915061259a565b60ff16878785818110611b7257611b72613649565b9050602002810190611b849190613f4a565b611b929060408101906139f8565b611b9c9080613b9d565b611baa906020810190613b57565b898987818110611bbc57611bbc613649565b9050602002810190611bce9190613f4a565b611bdc906080810190613b57565b85818110611bec57611bec613649565b919091013560f81c9050818110611c0557611c05613649565b9050013560f81c60f81b60f81c60ff1614611c325760405162461bcd60e51b8152600401610b2d90613d7d565b888884818110611c4457611c44613649565b9050602002810190611c569190613b9d565b611c64906060810190613d34565b82818110611c7457611c74613649565b9050608002016020016020810190611c8c919061259a565b60ff16898985818110611ca157611ca1613649565b9050602002810190611cb39190613b9d565b611cc1906060810190613d34565b83818110611cd157611cd1613649565b9050608002016040016020810190611ce9919061259a565b60ff1611611d095760405162461bcd60e51b8152600401610b2d90613de0565b83898985818110611d1c57611d1c613649565b9050602002810190611d2e9190613b9d565b611d3c906060810190613d34565b83818110611d4c57611d4c613649565b611d62926020608090920201908101915061259a565b60ff1681518110611d7557611d75613649565b016020015160f81c898985818110611d8f57611d8f613649565b9050602002810190611da19190613b9d565b611daf906060810190613d34565b83818110611dbf57611dbf613649565b9050608002016040016020810190611dd7919061259a565b60ff161015611df85760405162461bcd60e51b8152600401610b2d90613e51565b888884818110611e0a57611e0a613649565b9050602002810190611e1c9190613b9d565b611e2a906060810190613d34565b82818110611e3a57611e3a613649565b9050608002016040016020810190611e52919061259a565b60ff16878785818110611e6757611e67613649565b9050602002810190611e799190613f4a565b611e879060408101906139f8565b611e919080613b9d565b611e9f906040810190613b57565b898987818110611eb157611eb1613649565b9050602002810190611ec39190613f4a565b611ed1906080810190613b57565b85818110611ee157611ee1613649565b919091013560f81c9050818110611efa57611efa613649565b9050013560f81c60f81b60f81c60ff161015611f285760405162461bcd60e51b8152600401610b2d90613e51565b611f5e828a8a86818110611f3e57611f3e613649565b9050602002810190611f509190613b9d565b6116f1906060810190613d34565b915080611f6a816136c3565b915050611aca565b50611f7f610dd1856121cb565b611f9b5760405162461bcd60e51b8152600401610b2d90613ec2565b50611fa5816136c3565b9050611873565b505050505050505050565b6000611fc887878760200151611035565b9050611fde8a8a8a886000015188868989610a67565b50505050505050505050565b600081600001518260200151836040015160405160200161200d93929190613f60565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152015b604051602081830303815290604052805190602001209050919050565b6000836120638685856123dc565b1495945050505050565b60008160405160200161203891908151815260209182015163ffffffff169181019190915260400190565b60005b81518110156104f35760006001600160a01b0316836001600160a01b031663b5a872da8484815181106120d0576120d0613649565b60200260200101516040518263ffffffff1660e01b8152600401612100919063ffffffff91909116815260200190565b602060405180830381865afa15801561211d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121419190613fd9565b6001600160a01b031614156121bb5760405162461bcd60e51b8152602060048201526046602482015260008051602061413483398151915260448201527f7269667952656c61794b6579735365743a2072656c6179206b6579206973206e6064820152651bdd081cd95d60d21b608482015260a401610b2d565b6121c4816136c3565b905061209b565b6000610100825111156122545760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610b2d565b815161226257506000919050565b6000808360008151811061227857612278613649565b0160200151600160f89190911c81901b92505b845181101561234f578481815181106122a6576122a6613649565b0160200151600160f89190911c1b915082821161233b5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610b2d565b91811791612348816136c3565b905061228b565b50909392505050565b600061045182600001516040516020016123729190614002565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b6000816040516020016120389190614062565b6000602084516123ec9190614107565b156124735760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610b2d565b8260205b855181116124d65761248a600285614107565b6124ab578160005280860151602052604060002091506002840493506124c4565b8086015160005281602052604060002091506002840493505b6124cf60208261411b565b9050612477565b50949350505050565b604051806101000160405280606081526020016060815260200160608152602001612508612545565b815260200161252a604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b604051806040016040528061255861256a565b815260200161256561256a565b905290565b60405180604001604052806002906020820280368337509192915050565b60ff8116811461259757600080fd5b50565b6000602082840312156125ac57600080fd5b81356125b781612588565b9392505050565b6000606082840312156125d057600080fd5b50919050565b600080600083850360808112156125ec57600080fd5b60408112156125fa57600080fd5b5083925060408401356001600160401b038082111561261857600080fd5b612624878388016125be565b9350606086013591508082111561263a57600080fd5b508401610180818703121561264e57600080fd5b809150509250925092565b803561ffff8116811461266b57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156126a8576126a8612670565b60405290565b604051606081016001600160401b03811182821017156126a8576126a8612670565b604051608081016001600160401b03811182821017156126a8576126a8612670565b60405161010081016001600160401b03811182821017156126a8576126a8612670565b60405160c081016001600160401b03811182821017156126a8576126a8612670565b604051601f8201601f191681016001600160401b038111828210171561275f5761275f612670565b604052919050565b60006040828403121561277957600080fd5b612781612686565b9050813561278e81612588565b8152602082013561279e81612588565b602082015292915050565b600080606083850312156127bc57600080fd5b6127c583612659565b91506127d48460208501612767565b90509250929050565b6000602082840312156127ef57600080fd5b81356001600160401b0381111561280557600080fd5b612811848285016125be565b949350505050565b600081518084526020808501945080840160005b8381101561284f57815163ffffffff168752958201959082019060010161282d565b509495945050505050565b600081518084526020808501945080840160005b8381101561284f5761288b87835180518252602090810151910152565b604096909601959082019060010161286e565b8060005b60028110156107eb5781518452602093840193909101906001016128a2565b6128cc82825161289e565b60208101516104f3604084018261289e565b600081518084526020808501808196508360051b8101915082860160005b85811015612926578284038952612914848351612819565b988501989350908401906001016128fc565b5091979650505050505050565b6000610180825181855261294982860182612819565b91505060208301518482036020860152612963828261285a565b9150506040830151848203604086015261297d828261285a565b915050606083015161299260608601826128c1565b506080830151805160e08601526020015161010085015260a08301518482036101208601526129c18282612819565b91505060c08301518482036101408601526129dc8282612819565b91505060e08301518482036101608601526129f782826128de565b95945050505050565b6020815260006125b76020830184612933565b600060208284031215612a2557600080fd5b6125b782612659565b60005b83811015612a49578181015183820152602001612a31565b838111156107eb5750506000910152565b60008151808452612a72816020860160208601612a2e565b601f01601f19169290920160200192915050565b6020815260006125b76020830184612a5a565b60008060408385031215612aac57600080fd5b82356001600160401b0380821115612ac357600080fd5b9084019060808287031215612ad757600080fd5b90925060208401359080821115612aed57600080fd5b50830160a08186031215612b0057600080fd5b809150509250929050565b63ffffffff8116811461259757600080fd5b803561266b81612b0b565b60008082840360a0811215612b3c57600080fd5b6060811215612b4a57600080fd5b50612b536126ae565b8335612b5e81612b0b565b81526020840135612b6e81612b0b565b60208201526040840135612b8181612588565b604082015291506127d48460608501612767565b60008083601f840112612ba757600080fd5b5081356001600160401b03811115612bbe57600080fd5b6020830191508360208260051b8501011115612bd957600080fd5b9250929050565b60008060008060408587031215612bf657600080fd5b84356001600160401b0380821115612c0d57600080fd5b612c1988838901612b95565b90965094506020870135915080821115612c3257600080fd5b50612c3f87828801612b95565b95989497509550505050565b60008060408385031215612c5e57600080fd5b82356001600160401b0380821115612c7557600080fd5b612c81868387016125be565b93506020850135915080821115612c9757600080fd5b50612ca4858286016125be565b9150509250929050565b600060208284031215612cc057600080fd5b815180151581146125b757600080fd5b600060408284031215612ce257600080fd5b612cea612686565b905081358152602082013561279e81612b0b565b600060408284031215612d1057600080fd5b6125b78383612cd0565b60006001600160401b03821115612d3357612d33612670565b50601f01601f191660200190565b600082601f830112612d5257600080fd5b8135612d65612d6082612d1a565b612737565b818152846020838601011115612d7a57600080fd5b816020850160208301376000918101602001919091529392505050565b600060408284031215612da957600080fd5b612db1612686565b9050813581526020820135602082015292915050565b600082601f830112612dd857600080fd5b612de0612686565b806040840185811115612df257600080fd5b845b81811015612e0c578035845260209384019301612df4565b509095945050505050565b600060808284031215612e2957600080fd5b612e31612686565b9050612e3d8383612dc7565b815261279e8360408401612dc7565b60006001600160401b03821115612e6557612e65612670565b5060051b60200190565b600082601f830112612e8057600080fd5b81356020612e90612d6083612e4c565b82815260059290921b84018101918181019086841115612eaf57600080fd5b8286015b84811015612ed3578035612ec681612b0b565b8352918301918301612eb3565b509695505050505050565b600060608236031215612ef057600080fd5b612ef86126ae565b82356001600160401b0380821115612f0f57600080fd5b81850191506040808336031215612f2557600080fd5b612f2d612686565b833583811115612f3c57600080fd5b8401368190036101c0811215612f5157600080fd5b612f596126d0565b612f6283612659565b815260208084013587811115612f7757600080fd5b612f8336828701612d41565b8383015250610160603f1984011215612f9b57600080fd5b612fa36126d0565b9250612fb136878601612d97565b8352612fc03660808601612e17565b81840152612fd2366101008601612e17565b86840152610180840135612fe581612b0b565b8060608501525082868301526101a084013560608301528185528088013593508684111561301257600080fd5b61301e36858a01612e6f565b81860152848952613030818c01612b1d565b90890152505050508581013592508183111561304b57600080fd5b61305736848801612d41565b9084015250909392505050565b600082601f83011261307557600080fd5b81356020613085612d6083612e4c565b82815260069290921b840181019181810190868411156130a457600080fd5b8286015b84811015612ed3576130ba8882612d97565b8352918301916040016130a8565b600082601f8301126130d957600080fd5b813560206130e9612d6083612e4c565b82815260059290921b8401810191818101908684111561310857600080fd5b8286015b84811015612ed35780356001600160401b0381111561312b5760008081fd5b6131398986838b0101612e6f565b84525091830191830161310c565b6000610180823603121561315a57600080fd5b6131626126f2565b82356001600160401b038082111561317957600080fd5b61318536838701612e6f565b8352602085013591508082111561319b57600080fd5b6131a736838701613064565b602084015260408501359150808211156131c057600080fd5b6131cc36838701613064565b60408401526131de3660608701612e17565b60608401526131f03660e08701612d97565b608084015261012085013591508082111561320a57600080fd5b61321636838701612e6f565b60a084015261014085013591508082111561323057600080fd5b61323c36838701612e6f565b60c084015261016085013591508082111561325657600080fd5b50613263368286016130c8565b60e08301525092915050565b60006020828403121561328157600080fd5b81516125b781612588565b6000823561013e198336030181126132a357600080fd5b9190910192915050565b600061014082840312156132c057600080fd5b6132c8612715565b905081356001600160401b03808211156132e157600080fd5b6132ed85838601613064565b8352602084013591508082111561330357600080fd5b61330f85838601613064565b60208401526133218560408601612d97565b60408401526133338560808601612e17565b606084015261010084013591508082111561334d57600080fd5b5061335a84828501612e6f565b60808301525061336d6101208301612b1d565b60a082015292915050565b600061045136836132ad565b60006060828403121561339657600080fd5b604051606081018181106001600160401b03821117156133b8576133b8612670565b60405282516133c681612b0b565b815260208301516133d681612b0b565b602082015260408301516133e981612588565b60408201529392505050565b60006020828403121561340757600080fd5b81516001600160401b0381111561341d57600080fd5b8201601f8101841361342e57600080fd5b805161343c612d6082612d1a565b81815285602083850101111561345157600080fd5b6129f7826020830160208601612a2e565b60006060823603121561347457600080fd5b61347c612686565b6134863684612cd0565b815260408301356001600160401b038111156134a157600080fd5b6134ad368286016132ad565b60208301525092915050565b6000604082840312156134cb57600080fd5b6134d3612686565b82516134de81612588565b815260208301516134ee81612588565b60208201529392505050565b8481526080602082015260006135136080830186612a5a565b63ffffffff8516604084015282810360608401526135318185612933565b979650505050505050565b600082601f83011261354d57600080fd5b8151602061355d612d6083612e4c565b82815260059290921b8401810191818101908684111561357c57600080fd5b8286015b84811015612ed35780516001600160601b03811681146135a05760008081fd5b8352918301918301613580565b600080604083850312156135c057600080fd5b82516001600160401b03808211156135d757600080fd5b90840190604082870312156135eb57600080fd5b6135f3612686565b82518281111561360257600080fd5b61360e8882860161353c565b82525060208301518281111561362357600080fd5b61362f8882860161353c565b602083015250809450505050602083015190509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006001600160601b038083168185168183048111821515161561369b5761369b61365f565b02949350505050565b60008160001904831182151516156136be576136be61365f565b500290565b60006000198214156136d7576136d761365f565b5060010190565b600060ff821660ff8416808210156136f8576136f861365f565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b60008261372657613726613701565b500490565b60008282101561373d5761373d61365f565b500390565b600063ffffffff8083168185168183048111821515161561369b5761369b61365f565b60008351613777818460208801612a2e565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526137c56080840186612a5a565b838103606085015284518082528286019183019060005b818110156137f8578351835292840192918401916001016137dc565b50909998505050505050505050565b600082601f83011261381857600080fd5b81516020613828612d6083612e4c565b82815260059290921b8401810191818101908684111561384757600080fd5b8286015b84811015612ed357805161385e81612b0b565b835291830191830161384b565b600082601f83011261387c57600080fd5b8151602061388c612d6083612e4c565b82815260059290921b840181019181810190868411156138ab57600080fd5b8286015b84811015612ed35780516001600160401b038111156138ce5760008081fd5b6138dc8986838b0101613807565b8452509183019183016138af565b6000602082840312156138fc57600080fd5b81516001600160401b038082111561391357600080fd5b908301906080828603121561392757600080fd5b61392f6126d0565b82518281111561393e57600080fd5b61394a87828601613807565b82525060208301518281111561395f57600080fd5b61396b87828601613807565b60208301525060408301518281111561398357600080fd5b61398f87828601613807565b6040830152506060830151828111156139a757600080fd5b6139b38782860161386b565b60608301525095945050505050565b6000602082840312156139d457600080fd5b81356125b781612b0b565b6000602082840312156139f157600080fd5b5051919050565b60008235605e198336030181126132a357600080fd5b600060608236031215613a2057600080fd5b613a286126ae565b82356001600160401b0380821115613a3f57600080fd5b818501915060808236031215613a5457600080fd5b613a5c6126d0565b82358152602083013582811115613a7257600080fd5b613a7e36828601612d41565b602083015250604083013582811115613a9657600080fd5b613aa236828601612d41565b60408301525060608301359250613ab883612b0b565b82606082015280845250505060208301356020820152613ada60408401612b1d565b604082015292915050565b6020808252606090820181905260008051602061413483398151915260408301527f72696679426c6f62466f7251756f72756d733a2062617463684d657461646174908201527f6120646f6573206e6f74206d617463682073746f726564206d65746164617461608082015260a00190565b6000808335601e19843603018112613b6e57600080fd5b8301803591506001600160401b03821115613b8857600080fd5b602001915036819003821315612bd957600080fd5b60008235607e198336030181126132a357600080fd5b60006080808336031215613bc657600080fd5b613bce6126ae565b613bd83685612d97565b8152604080850135613be981612b0b565b6020818185015260609150818701356001600160401b03811115613c0c57600080fd5b870136601f820112613c1d57600080fd5b8035613c2b612d6082612e4c565b81815260079190911b82018301908381019036831115613c4a57600080fd5b928401925b82841015613cbe57888436031215613c675760008081fd5b613c6f6126d0565b8435613c7a81612588565b815284860135613c8981612588565b8187015284880135613c9a81612588565b8189015284870135613cab81612b0b565b8188015282529288019290840190613c4f565b958701959095525093979650505050505050565b6020808252604e9082015260008051602061413483398151915260408201527f72696679426c6f62466f7251756f72756d733a20696e636c7573696f6e20707260608201526d1bdbd9881a5cc81a5b9d985b1a5960921b608082015260a00190565b6000808335601e19843603018112613d4b57600080fd5b8301803591506001600160401b03821115613d6557600080fd5b6020019150600781901b3603821315612bd957600080fd5b6020808252604f9082015260008051602061413483398151915260408201527f72696679426c6f62466f7251756f72756d733a2071756f72756d4e756d62657260608201526e040c8decae640dcdee840dac2e8c6d608b1b608082015260a00190565b602080825260579082015260008051602061413483398151915260408201527f72696679426c6f62466f7251756f72756d733a207468726573686f6c6420706560608201527f7263656e746167657320617265206e6f742076616c6964000000000000000000608082015260a00190565b6020808252605e9082015260008051602061413483398151915260408201527f72696679426c6f62466f7251756f72756d733a20636f6e6669726d6174696f6e60608201527f5468726573686f6c6450657263656e74616765206973206e6f74206d65740000608082015260a00190565b6020808252606e9082015260008051602061413483398151915260408201527f72696679426c6f62466f7251756f72756d733a2072657175697265642071756f60608201527f72756d7320617265206e6f74206120737562736574206f662074686520636f6e60808201526d6669726d65642071756f72756d7360901b60a082015260c00190565b60008235609e198336030181126132a357600080fd5b60006101a061ffff86168352806020840152613f7e81840186612a5a565b8451805160408601526020015160608501529150613f999050565b6020830151613fab60808401826128c1565b506040830151613fbf6101008401826128c1565b5063ffffffff606084015116610180830152949350505050565b600060208284031215613feb57600080fd5b81516001600160a01b03811681146125b757600080fd5b6020815281516020820152600060208301516080604084015261402860a0840182612a5a565b90506040840151601f198483030160608501526140458282612a5a565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b808410156140f8578751805160ff90811684528a82015181168b8501528882015116888401528601518516868301529688019660019390930192908201906140b5565b509a9950505050505050505050565b60008261411657614116613701565b500690565b6000821982111561412e5761412e61365f565b50019056fe456967656e4441426c6f62566572696669636174696f6e5574696c732e5f7665a2646970667358221220d7ca5d1839a0cbbd77a4e2f2f1f3ddc3b0172168ff2052c6cd2b0b89ba3ce09564736f6c634300080c0033",
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

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xcc392c25.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierCaller) VerifyBlobV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	var out []interface{}
	err := _ContractEigenDABlobVerifier.contract.Call(opts, &out, "verifyBlobV2FromSignedBatch", signedBatch, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xcc392c25.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
func (_ContractEigenDABlobVerifier *ContractEigenDABlobVerifierSession) VerifyBlobV2FromSignedBatch(signedBatch SignedBatch, blobVerificationProof BlobVerificationProofV2) error {
	return _ContractEigenDABlobVerifier.Contract.VerifyBlobV2FromSignedBatch(&_ContractEigenDABlobVerifier.CallOpts, signedBatch, blobVerificationProof)
}

// VerifyBlobV2FromSignedBatch is a free data retrieval call binding the contract method 0xcc392c25.
//
// Solidity: function verifyBlobV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[],uint32)) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),uint32[]),uint32,bytes) blobVerificationProof) view returns()
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
