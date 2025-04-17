// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV1_V2

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

// ContractEigenDACertVerifierV1V2MetaData contains all meta data concerning the ContractEigenDACertVerifierV1V2 contract.
var ContractEigenDACertVerifierV1V2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetrieverV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsAdversary\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsConfirmation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1ForZkProof\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structBlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structBlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BatchMetadataMismatch\",\"inputs\":[{\"name\":\"actualHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expectedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfirmationThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuorumNumberMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"errParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StakeThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}]",
	Bin: "0x6101a06040523480156200001257600080fd5b5060405162003ff638038062003ff68339810160408190526200003591620001d9565b6001600160a01b03808916608081905281891660a05260c05280871660e05280861661010052808516610120528316610140526020820151825189918891889188918891889160ff918216911611620000a1576040516308a6997560e01b815260040160405180910390fd5b805160ff90811661016052602090910151166101805250620002fe9b505050505050505050505050565b6001600160a01b0381168114620000e157600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620001255762000125620000e4565b604052919050565b805160ff811681146200013f57600080fd5b919050565b600082601f8301126200015657600080fd5b81516001600160401b03811115620001725762000172620000e4565b602062000188601f8301601f19168201620000fa565b82815285828487010111156200019d57600080fd5b60005b83811015620001bd578581018301518282018401528201620001a0565b83811115620001cf5760008385840101525b5095945050505050565b600080600080600080600080888a03610120811215620001f857600080fd5b89516200020581620000cb565b60208b01519099506200021881620000cb565b60408b01519098506200022b81620000cb565b60608b01519097506200023e81620000cb565b60808b01519096506200025181620000cb565b60a08b01519095506200026481620000cb565b9350604060bf19820112156200027957600080fd5b50604080519081016001600160401b038082118383101715620002a057620002a0620000e4565b81604052620002b260c08d016200012d565b8352620002c260e08d016200012d565b60208401526101008c015192945080831115620002de57600080fd5b5050620002ee8b828c0162000144565b9150509295985092959890939650565b60805160a05160c05160e0516101005161012051610140516101605161018051613c10620003e6600039600081816102f00152610b630152600081816103710152610b3e015260008181610242015281816105a6015261073a01526000818161021b01528181610585015261071901526000818161028f0152818161048f01526105640152600081816101400152818161046e01526105430152600081816101840152818161044c01528181610522015281816106280152818161066a01528181610b8e0152610de10152600081816102b60152610e41015260006101f40152613c106000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c80637d644cad116100ad578063dea610a911610071578063dea610a9146102eb578063efd57acb14610324578063f25de3f814610339578063f88adbba14610359578063fd1744841461036c57600080fd5b80637d644cad14610264578063813c2eb01461027757806382c216e71461028a578063a9c823e1146102b1578063ccb7cd0d146102d857600080fd5b8063415ef614116100f4578063415ef614146101b9578063421c0222146101dc5780634cff90c4146101ef5780635df1f618146102165780635fafa4821461023d57600080fd5b8063143eb4d914610126578063154b9e861461013b57806317f3578e1461017f57806331a3479a146101a6575b600080fd5b6101396101343660046123ac565b610393565b005b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101396101b4366004612466565b6103b4565b6101cc6101c7366004612566565b610444565b6040519015158152602001610176565b6101396101ea36600461260f565b61051a565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b610139610272366004612672565b6105ec565b610139610285366004612566565b610620565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101627f000000000000000000000000000000000000000000000000000000000000000081565b6101396102e63660046126f6565b610662565b6103127f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff9091168152602001610176565b61032c6106fd565b6040516101769190612779565b61034c610347366004612793565b61070c565b60405161017691906129b6565b6101cc610367366004612672565b61076e565b6103127f000000000000000000000000000000000000000000000000000000000000000081565b6000806103a084846107af565b915091506103ae828261088f565b50505050565b8281146103d7576040516001621398b960e31b0319815260040160405180910390fd5b60005b8381101561043d5761042d8585838181106103f7576103f76129c9565b905060200281019061040991906129df565b84848481811061041b5761041b6129c9565b905060200281019061027291906129ff565b61043681612a2b565b90506103da565b5050505050565b6000806104e57f00000000000000000000000000000000000000000000000000000000000000005b7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006104bd368b90038b018b612a74565b6104c68a612bd7565b6104cf8a612e65565b6104d7610b1e565b6104df610b8a565b8b610c12565b50905060008160068111156104fc576104fc612f8d565b141561050c576001915050610512565b60009150505b949350505050565b6000806103a07f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006105ce8a612fa3565b6105d78a612bd7565b6105df610b1e565b6105e7610b8a565b610db2565b6000806106126105fa610ddd565b61060385610e3d565b868661060d610b8a565b610ee1565b915091506103ae8282610f8c565b60008061064c7f000000000000000000000000000000000000000000000000000000000000000061046c565b9150915061065a828261088f565b505050505050565b6000806103a07f0000000000000000000000000000000000000000000000000000000000000000604051632ecfe72b60e01b815261ffff871660048201526001600160a01b039190911690632ecfe72b90602401606060405180830381865afa1580156106d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106f7919061309a565b846107af565b6060610707610b8a565b905090565b610714612196565b6107677f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061076285612fa3565b61125c565b5092915050565b60008061077c6105fa610ddd565b509050600081600a81111561079357610793612f8d565b14156107a35760019150506107a9565b60009150505b92915050565b600060606000836020015184600001516107c9919061310b565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406107f29190613144565b6107fc9190613144565b61080890612710613158565b610812919061316f565b86519091506000906108269061271061318e565b63ffffffff1690508082106108535760006040518060200160405280600081525094509450505050610888565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050505b9250929050565b60008260068111156108a3576108a3612f8d565b14156108ad575050565b60018260068111156108c1576108c1612f8d565b1415610917576000806000838060200190518101906108e091906131ba565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b600282600681111561092b5761092b612f8d565b1415610973576000808280602001905181019061094891906131e8565b604051631b00235d60e01b815260ff808416600483015282166024820152919350915060440161090e565b600382600681111561098757610987612f8d565b14156109db576000806000838060200190518101906109a691906131ba565b6040516001626dc9ad60e11b03198152600481018490526024810183905260448101829052929550909350915060640161090e565b60048260068111156109ef576109ef612f8d565b1415610a345760008082806020019051810190610a0c9190613217565b604051634a47030360e11b81526004810183905260248101829052919350915060440161090e565b6005826006811115610a4857610a48612f8d565b1415610a8d5760008082806020019051810190610a659190613217565b60405163114b085b60e21b81526004810183905260248101829052919350915060440161090e565b6006826006811115610aa157610aa1612f8d565b1415610ae157600081806020019051810190610abd919061323b565b6040516309efaa0b60e41b815263ffffffff8216600482015290915060240161090e565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b604482015260640161090e565b6040805180820182526000808252602091820152815180830190925260ff7f0000000000000000000000000000000000000000000000000000000000000000811683527f0000000000000000000000000000000000000000000000000000000000000000169082015290565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610bea573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107079190810190613258565b60006060610c208888611471565b90925090506000826006811115610c3957610c39612f8d565b14610c4357610da4565b610c5589886000015160400151611548565b90925090506000826006811115610c6e57610c6e612f8d565b14610c7857610da4565b86515151604051632ecfe72b60e01b815261ffff9091166004820152610cf3906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa158015610cc9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ced919061309a565b866107af565b90925090506000826006811115610d0c57610d0c612f8d565b14610d1657610da4565b6000610d328b610d258b611680565b868c602001518b8b6116c8565b919450925090506000836006811115610d4d57610d4d612f8d565b14610d585750610da4565b87515160200151600090610d6c9083611829565b919550935090506000846006811115610d8757610d87612f8d565b14610d93575050610da4565b610d9d868261188f565b9350935050505b995099975050505050505050565b60006060600080610dc48a8a8a61125c565b91509150610d9d8d8d8d8b600001518b878c8c89610c12565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610bea573d6000803e3d6000fd5b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663eccbbfc9610e7b60208501856132c5565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015610ebd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107a991906132e2565b60006060610eef86856118df565b9092509050600082600a811115610f0857610f08612f8d565b14610f1257610f82565b610f1c8585611941565b9092509050600082600a811115610f3557610f35612f8d565b14610f3f57610f82565b6000610f4c888787611a69565b91945092509050600083600a811115610f6757610f67612f8d565b14610f725750610f82565b610f7c8482611b0f565b92509250505b9550959350505050565b600082600a811115610fa057610fa0612f8d565b1415610faa575050565b600182600a811115610fbe57610fbe612f8d565b14156110035760008082806020019051810190610fdb9190613217565b6040516315d615c960e21b81526004810183905260248101829052919350915060440161090e565b600282600a81111561101757611017612f8d565b1415611036576000806000838060200190518101906108e091906131ba565b600382600a81111561104a5761104a612f8d565b1415611092576000808280602001905181019061106791906131e8565b6040516314fa310760e31b815260ff808416600483015282166024820152919350915060440161090e565b600482600a8111156110a6576110a6612f8d565b14156110c3576000808280602001905181019061094891906131e8565b600582600a8111156110d7576110d7612f8d565b141561112c576000806000838060200190518101906110f691906132fb565b604051638aa11c4360e01b815260ff8085166004830152808416602483015282166044820152929550909350915060640161090e565b600682600a81111561114057611140612f8d565b14156111955760008060008380602001905181019061115f91906132fb565b60405163a4ad875560e01b815260ff8085166004830152808416602483015282166044820152929550909350915060640161090e565b600782600a8111156111a9576111a9612f8d565b14156111c65760008082806020019051810190610a659190613217565b600a82600a8111156111da576111da612f8d565b14156111f657600081806020019051810190610abd919061323b565b600882600a81111561120a5761120a612f8d565b141561122b5780604051638c59c92f60e01b815260040161090e9190612779565b600982600a81111561123f5761123f612f8d565b1415610ae15760008082806020019051810190610a0c9190613217565b611264612196565b60606000836020015160000151516001600160401b038111156112895761128961223f565b6040519080825280602002602001820160405280156112b2578160200160208202803683370190505b50905060005b6020850151515181101561132f5761130285602001516000015182815181106112e3576112e36129c9565b6020026020010151805160009081526020918201519091526040902090565b828281518110611314576113146129c9565b602090810291909101015261132881612a2b565b90506112b8565b5060005b8460200151608001515181101561139a5782856020015160800151828151811061135f5761135f6129c9565b6020026020010151604051602001611378929190613348565b60405160208183030381529060405292508061139390612a2b565b9050611333565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f74916113d5918a91908890889060040161337a565b600060405180830381865afa1580156113f2573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261141a91908101906134cd565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b6000606060006114848460000151611b5f565b905060008160405160200161149b91815260200190565b60405160208183030381529060405280519060200120905060008660000151905060006114d8876040015183858a6020015163ffffffff16611b87565b905080156114ff576000604051806020016040528060008152509550955050505050610888565b6020808801516040805163ffffffff90921692820192909252908101849052606081018390526001906080015b6040516020818303038152906040529550955050505050610888565b6000606060005b83518110156116655760006001600160a01b0316856001600160a01b031663b5a872da868481518110611584576115846129c9565b60200260200101516040518263ffffffff1660e01b81526004016115b4919063ffffffff91909116815260200190565b602060405180830381865afa1580156115d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115f591906135a5565b6001600160a01b03161415611655576006848281518110611618576116186129c9565b602002602001015160405160200161163c919063ffffffff91909116815260200190565b6040516020818303038152906040529250925050610888565b61165e81612a2b565b905061154f565b50506040805160208101909152600080825291509250929050565b6000816040516020016116ab91908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b815260040161170194939291906135ce565b600060405180830381865afa15801561171e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526117469190810190613681565b5090506000915060005b885181101561180757856000015160ff1682602001518281518110611777576117776129c9565b6020026020010151611789919061371d565b6001600160601b03166064836000015183815181106117aa576117aa6129c9565b60200260200101516001600160601b03166117c5919061316f565b106117f5576117f2838a83815181106117e0576117e06129c9565b0160200151600160f89190911c1b1790565b92505b806117ff81612a2b565b915050611750565b5050604080516020810190915260008082529350915096509650969350505050565b60006060600061183885611b9f565b905083811681141561185d576040805160208101909152600080825293509150611888565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b60006060600061189e85611b9f565b90508381168114156118c3575050604080516020810190915260008082529150610888565b604080516020810183905290810185905260059060600161163c565b60006060816119026118f46040860186613743565b6118fd90613759565b611d2c565b905084811415611925575050604080516020810190915260008082529150610888565b604080516020810183905290810186905260019060600161163c565b600060608161195761195286613830565b611d9d565b905060008160405160200161196e91815260200190565b60405160208183030381529060405280519060200120905060008580604001906119989190613743565b6119a290806129df565b3590506000611a0b6119b7606089018961394f565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250869250879150611a00905060408c0160208d016132c5565b63ffffffff16611b87565b90508015611a32576000604051806020016040528060008152509550955050505050610888565b6002611a446040890160208a016132c5565b6040805163ffffffff909216602083015281018590526060810184905260800161152c565b600060608180611a7b86840187613995565b9050905060005b81811015611aef576000806000611a9c8b8b8b8788611db0565b91945092509050600083600a811115611ab757611ab7612f8d565b14611ace5750909550935060009250611b06915050565b600160ff82161b861795505050508080611ae790612a2b565b915050611a82565b505060408051602081019091526000808252935091505b93509350939050565b600060606000611b1e85611b9f565b9050838116811415611b43575050604080516020810190915260008082529150610888565b604080516020810183905290810185905260079060600161163c565b6000611b6e8260000151612041565b60208084015160408086015190516116ab9493016139de565b600083611b95868585612093565b1495945050505050565b600061010082511115611c285760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a40161090e565b8151611c3657506000919050565b60008083600081518110611c4c57611c4c6129c9565b0160200151600160f89190911c81901b92505b8451811015611d2357848181518110611c7a57611c7a6129c9565b0160200151600160f89190911c1b9150828211611d0f5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a40161090e565b91811791611d1c81612a2b565b9050611c5f565b50909392505050565b60006107a98260000151604051602001611d469190613a13565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b6000816040516020016116ab9190613a73565b600060608136611dc288840189613995565b87818110611dd257611dd26129c9565b608002919091019150611dea90506020820182613b18565b91506000611dfb608089018961394f565b87818110611e0b57611e0b6129c9565b919091013560f81c915060009050611e2660408a018a613743565b611e3090806129df565b611e3e90602081019061394f565b8360ff16818110611e5157611e516129c9565b919091013560f81c91505060ff84168114611e9e576040805160ff958616602082015291909416818501528351808203850181526060909101909352506003935090915060009050612036565b6000611eb06060850160408601613b18565b90506000611ec46040860160208701613b18565b90508060ff168260ff1611611f0d576040805160ff9384166020820152919092168183015281518082038301815260609091019091526004965094506000935061203692505050565b60008d8760ff1681518110611f2457611f246129c9565b016020015160f81c905060ff8316811115611f81576040805160ff808a1660208301528084169282019290925290841660608201526005906080016040516020818303038152906040526000985098509850505050505050612036565b6000611f9060408e018e613743565b611f9a90806129df565b611fa890604081019061394f565b8760ff16818110611fbb57611fbb6129c9565b919091013560f81c91505060ff841681101561201a576040805160ff808b166020830152808716928201929092529082166060820152600690608001604051602081830303815290604052600099509950995050505050505050612036565b5050604080516020810190915260008082529850965050505050505b955095509592505050565b600081600001518260200151836040015160405160200161206493929190613b35565b60408051601f1981840301815282825280516020918201206060808701519285019190915291830152016116ab565b6000602084516120a39190613bae565b1561212a5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a40161090e565b8260205b8551811161218d57612141600285613bae565b6121625781600052808601516020526040600020915060028404935061217b565b8086015160005281602052604060002091506002840493505b612186602082613bc2565b905061212e565b50949350505050565b6040518061010001604052806060815260200160608152602001606081526020016121bf6121fc565b81526020016121e1604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b604051806040016040528061220f612221565b815260200161221c612221565b905290565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156122775761227761223f565b60405290565b604051606081016001600160401b03811182821017156122775761227761223f565b604051608081016001600160401b03811182821017156122775761227761223f565b60405161010081016001600160401b03811182821017156122775761227761223f565b60405160a081016001600160401b03811182821017156122775761227761223f565b604051601f8201601f191681016001600160401b038111828210171561232e5761232e61223f565b604052919050565b63ffffffff8116811461234857600080fd5b50565b803561235681612336565b919050565b60ff8116811461234857600080fd5b60006040828403121561237c57600080fd5b612384612255565b905081356123918161235b565b815260208201356123a18161235b565b602082015292915050565b60008082840360a08112156123c057600080fd5b60608112156123ce57600080fd5b506123d761227d565b83356123e281612336565b815260208401356123f281612336565b602082015260408401356124058161235b565b60408201529150612419846060850161236a565b90509250929050565b60008083601f84011261243457600080fd5b5081356001600160401b0381111561244b57600080fd5b6020830191508360208260051b850101111561088857600080fd5b6000806000806040858703121561247c57600080fd5b84356001600160401b038082111561249357600080fd5b61249f88838901612422565b909650945060208701359150808211156124b857600080fd5b506124c587828801612422565b95989497509550505050565b6000606082840312156124e357600080fd5b50919050565b60006001600160401b038211156125025761250261223f565b50601f01601f191660200190565b600082601f83011261252157600080fd5b813561253461252f826124e9565b612306565b81815284602083860101111561254957600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a081121561257d57600080fd5b604081121561258b57600080fd5b5084935060408501356001600160401b03808211156125a957600080fd5b6125b5888389016124d1565b945060608701359150808211156125cb57600080fd5b9086019061018082890312156125e057600080fd5b909250608086013590808211156125f657600080fd5b5061260387828801612510565b91505092959194509250565b6000806040838503121561262257600080fd5b82356001600160401b038082111561263957600080fd5b612645868387016124d1565b9350602085013591508082111561265b57600080fd5b50612668858286016124d1565b9150509250929050565b6000806040838503121561268557600080fd5b82356001600160401b038082111561269c57600080fd5b90840190608082870312156126b057600080fd5b909250602084013590808211156126c657600080fd5b50830160a081860312156126d957600080fd5b809150509250929050565b803561ffff8116811461235657600080fd5b6000806060838503121561270957600080fd5b612712836126e4565b9150612419846020850161236a565b60005b8381101561273c578181015183820152602001612724565b838111156103ae5750506000910152565b60008151808452612765816020860160208601612721565b601f01601f19169290920160200192915050565b60208152600061278c602083018461274d565b9392505050565b6000602082840312156127a557600080fd5b81356001600160401b038111156127bb57600080fd5b610512848285016124d1565b600081518084526020808501945080840160005b838110156127fd57815163ffffffff16875295820195908201906001016127db565b509495945050505050565b600081518084526020808501945080840160005b838110156127fd5761283987835180518252602090810151910152565b604096909601959082019060010161281c565b8060005b60028110156103ae578151845260209384019390910190600101612850565b61287a82825161284c565b602081015161288c604084018261284c565b505050565b600082825180855260208086019550808260051b84010181860160005b848110156128dc57601f198684030189526128ca8383516127c7565b988401989250908301906001016128ae565b5090979650505050505050565b600061018082518185526128ff828601826127c7565b915050602083015184820360208601526129198282612808565b915050604083015184820360408601526129338282612808565b9150506060830151612948606086018261286f565b506080830151805160e08601526020015161010085015260a083015184820361012086015261297782826127c7565b91505060c083015184820361014086015261299282826127c7565b91505060e08301518482036101608601526129ad8282612891565b95945050505050565b60208152600061278c60208301846128e9565b634e487b7160e01b600052603260045260246000fd5b60008235607e198336030181126129f557600080fd5b9190910192915050565b60008235609e198336030181126129f557600080fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415612a3f57612a3f612a15565b5060010190565b600060408284031215612a5857600080fd5b612a60612255565b90508135815260208201356123a181612336565b600060408284031215612a8657600080fd5b61278c8383612a46565b600060408284031215612aa257600080fd5b612aaa612255565b9050813581526020820135602082015292915050565b600082601f830112612ad157600080fd5b612ad9612255565b806040840185811115612aeb57600080fd5b845b81811015612b05578035845260209384019301612aed565b509095945050505050565b600060808284031215612b2257600080fd5b612b2a612255565b9050612b368383612ac0565b81526123a18360408401612ac0565b60006001600160401b03821115612b5e57612b5e61223f565b5060051b60200190565b600082601f830112612b7957600080fd5b81356020612b8961252f83612b45565b82815260059290921b84018101918181019086841115612ba857600080fd5b8286015b84811015612bcc578035612bbf81612336565b8352918301918301612bac565b509695505050505050565b600060608236031215612be957600080fd5b612bf161227d565b82356001600160401b0380821115612c0857600080fd5b818501915060608236031215612c1d57600080fd5b612c2561227d565b823582811115612c3457600080fd5b8301368190036101c0811215612c4957600080fd5b612c5161229f565b612c5a836126e4565b815260208084013586811115612c6f57600080fd5b612c7b36828701612510565b83830152506040610160603f1985011215612c9557600080fd5b612c9d61229f565b9350612cab36828701612a90565b8452612cba3660808701612b10565b82850152612ccc366101008701612b10565b81850152610180850135612cdf81612336565b8060608601525083818401526101a0850135606084015282865281880135945086851115612d0c57600080fd5b612d1836868a01612510565b8287015280880135945086851115612d2f57600080fd5b612d3b36868a01612b68565b81870152858952612d4d828c0161234b565b828a0152808b0135975086881115612d6457600080fd5b612d7036898d01612510565b90890152509598975050505050505050565b600082601f830112612d9357600080fd5b81356020612da361252f83612b45565b82815260069290921b84018101918181019086841115612dc257600080fd5b8286015b84811015612bcc57612dd88882612a90565b835291830191604001612dc6565b600082601f830112612df757600080fd5b81356020612e0761252f83612b45565b82815260059290921b84018101918181019086841115612e2657600080fd5b8286015b84811015612bcc5780356001600160401b03811115612e495760008081fd5b612e578986838b0101612b68565b845250918301918301612e2a565b60006101808236031215612e7857600080fd5b612e806122c1565b82356001600160401b0380821115612e9757600080fd5b612ea336838701612b68565b83526020850135915080821115612eb957600080fd5b612ec536838701612d82565b60208401526040850135915080821115612ede57600080fd5b612eea36838701612d82565b6040840152612efc3660608701612b10565b6060840152612f0e3660e08701612a90565b6080840152610120850135915080821115612f2857600080fd5b612f3436838701612b68565b60a0840152610140850135915080821115612f4e57600080fd5b612f5a36838701612b68565b60c0840152610160850135915080821115612f7457600080fd5b50612f8136828601612de6565b60e08301525092915050565b634e487b7160e01b600052602160045260246000fd5b600060608236031215612fb557600080fd5b612fbd612255565b612fc73684612a46565b815260408301356001600160401b0380821115612fe357600080fd5b81850191506101208236031215612ff957600080fd5b6130016122e4565b82358281111561301057600080fd5b61301c36828601612d82565b82525060208301358281111561303157600080fd5b61303d36828601612d82565b6020830152506130503660408501612a90565b60408201526130623660808501612b10565b60608201526101008301358281111561307a57600080fd5b61308636828601612b68565b608083015250602084015250909392505050565b6000606082840312156130ac57600080fd5b604051606081018181106001600160401b03821117156130ce576130ce61223f565b60405282516130dc81612336565b815260208301516130ec81612336565b602082015260408301516130ff8161235b565b60408201529392505050565b600060ff821660ff84168082101561312557613125612a15565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826131535761315361312e565b500490565b60008282101561316a5761316a612a15565b500390565b600081600019048311821515161561318957613189612a15565b500290565b600063ffffffff808316818516818304811182151516156131b1576131b1612a15565b02949350505050565b6000806000606084860312156131cf57600080fd5b8351925060208401519150604084015190509250925092565b600080604083850312156131fb57600080fd5b82516132068161235b565b60208401519092506126d98161235b565b6000806040838503121561322a57600080fd5b505080516020909101519092909150565b60006020828403121561324d57600080fd5b815161278c81612336565b60006020828403121561326a57600080fd5b81516001600160401b0381111561328057600080fd5b8201601f8101841361329157600080fd5b805161329f61252f826124e9565b8181528560208385010111156132b457600080fd5b6129ad826020830160208601612721565b6000602082840312156132d757600080fd5b813561278c81612336565b6000602082840312156132f457600080fd5b5051919050565b60008060006060848603121561331057600080fd5b835161331b8161235b565b602085015190935061332c8161235b565b604085015190925061333d8161235b565b809150509250925092565b6000835161335a818460208801612721565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526133a8608084018661274d565b838103606085015284518082528286019183019060005b818110156133db578351835292840192918401916001016133bf565b50909998505050505050505050565b600082601f8301126133fb57600080fd5b8151602061340b61252f83612b45565b82815260059290921b8401810191818101908684111561342a57600080fd5b8286015b84811015612bcc57805161344181612336565b835291830191830161342e565b600082601f83011261345f57600080fd5b8151602061346f61252f83612b45565b82815260059290921b8401810191818101908684111561348e57600080fd5b8286015b84811015612bcc5780516001600160401b038111156134b15760008081fd5b6134bf8986838b01016133ea565b845250918301918301613492565b6000602082840312156134df57600080fd5b81516001600160401b03808211156134f657600080fd5b908301906080828603121561350a57600080fd5b61351261229f565b82518281111561352157600080fd5b61352d878286016133ea565b82525060208301518281111561354257600080fd5b61354e878286016133ea565b60208301525060408301518281111561356657600080fd5b613572878286016133ea565b60408301525060608301518281111561358a57600080fd5b6135968782860161344e565b60608301525095945050505050565b6000602082840312156135b757600080fd5b81516001600160a01b038116811461278c57600080fd5b8481526080602082015260006135e7608083018661274d565b63ffffffff85166040840152828103606084015261360581856128e9565b979650505050505050565b600082601f83011261362157600080fd5b8151602061363161252f83612b45565b82815260059290921b8401810191818101908684111561365057600080fd5b8286015b84811015612bcc5780516001600160601b03811681146136745760008081fd5b8352918301918301613654565b6000806040838503121561369457600080fd5b82516001600160401b03808211156136ab57600080fd5b90840190604082870312156136bf57600080fd5b6136c7612255565b8251828111156136d657600080fd5b6136e288828601613610565b8252506020830151828111156136f757600080fd5b61370388828601613610565b602083015250809450505050602083015190509250929050565b60006001600160601b03808316818516818304811182151516156131b1576131b1612a15565b60008235605e198336030181126129f557600080fd5b60006060823603121561376b57600080fd5b61377361227d565b82356001600160401b038082111561378a57600080fd5b81850191506080823603121561379f57600080fd5b6137a761229f565b823581526020830135828111156137bd57600080fd5b6137c936828601612510565b6020830152506040830135828111156137e157600080fd5b6137ed36828601612510565b6040830152506060830135925061380383612336565b826060820152808452505050602083013560208201526138256040840161234b565b604082015292915050565b6000608080833603121561384357600080fd5b61384b61227d565b6138553685612a90565b815260408085013561386681612336565b6020818185015260609150818701356001600160401b0381111561388957600080fd5b870136601f82011261389a57600080fd5b80356138a861252f82612b45565b81815260079190911b820183019083810190368311156138c757600080fd5b928401925b8284101561393b578884360312156138e45760008081fd5b6138ec61229f565b84356138f78161235b565b8152848601356139068161235b565b81870152848801356139178161235b565b818901528487013561392881612336565b81880152825292880192908401906138cc565b958701959095525093979650505050505050565b6000808335601e1984360301811261396657600080fd5b8301803591506001600160401b0382111561398057600080fd5b60200191503681900382131561088857600080fd5b6000808335601e198436030181126139ac57600080fd5b8301803591506001600160401b038211156139c657600080fd5b6020019150600781901b360382131561088857600080fd5b8381526060602082015260006139f7606083018561274d565b8281036040840152613a0981856127c7565b9695505050505050565b60208152815160208201526000602083015160806040840152613a3960a084018261274d565b90506040840151601f19848303016060850152613a56828261274d565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b80841015613b09578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190613ac6565b509a9950505050505050505050565b600060208284031215613b2a57600080fd5b813561278c8161235b565b60006101a061ffff86168352806020840152613b538184018661274d565b8451805160408601526020015160608501529150613b6e9050565b6020830151613b80608084018261286f565b506040830151613b9461010084018261286f565b5063ffffffff606084015116610180830152949350505050565b600082613bbd57613bbd61312e565b500690565b60008219821115613bd557613bd5612a15565b50019056fea2646970667358221220883d380ab0798977e8af10abeb302b0bf378cce0a67e61d0abbb9d0325ceacf664736f6c634300080c0033",
}

// ContractEigenDACertVerifierV1V2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV1V2MetaData.ABI instead.
var ContractEigenDACertVerifierV1V2ABI = ContractEigenDACertVerifierV1V2MetaData.ABI

// ContractEigenDACertVerifierV1V2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV1V2MetaData.Bin instead.
var ContractEigenDACertVerifierV1V2Bin = ContractEigenDACertVerifierV1V2MetaData.Bin

// DeployContractEigenDACertVerifierV1V2 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV1V2 to it.
func DeployContractEigenDACertVerifierV1V2(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistry common.Address, _eigenDABatchMetadataStorage common.Address, _eigenDASignatureVerifier common.Address, _eigenDARelayRegistry common.Address, _operatorStateRetriever common.Address, _registryCoordinator common.Address, _securityThresholdsV2 SecurityThresholds, _quorumNumbersRequiredV2 []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV1V2, error) {
	parsed, err := ContractEigenDACertVerifierV1V2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV1V2Bin), backend, _eigenDAThresholdRegistry, _eigenDABatchMetadataStorage, _eigenDASignatureVerifier, _eigenDARelayRegistry, _operatorStateRetriever, _registryCoordinator, _securityThresholdsV2, _quorumNumbersRequiredV2)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifierV1V2{ContractEigenDACertVerifierV1V2Caller: ContractEigenDACertVerifierV1V2Caller{contract: contract}, ContractEigenDACertVerifierV1V2Transactor: ContractEigenDACertVerifierV1V2Transactor{contract: contract}, ContractEigenDACertVerifierV1V2Filterer: ContractEigenDACertVerifierV1V2Filterer{contract: contract}}, nil
}

// ContractEigenDACertVerifierV1V2 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2 struct {
	ContractEigenDACertVerifierV1V2Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV1V2Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV1V2Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV1V2Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1V2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1V2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV1V2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1V2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV1V2Session struct {
	Contract     *ContractEigenDACertVerifierV1V2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1V2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV1V2CallerSession struct {
	Contract *ContractEigenDACertVerifierV1V2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// ContractEigenDACertVerifierV1V2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV1V2TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV1V2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1V2Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2Raw struct {
	Contract *ContractEigenDACertVerifierV1V2 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1V2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2CallerRaw struct {
	Contract *ContractEigenDACertVerifierV1V2Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1V2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1V2TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV1V2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV1V2 creates a new instance of ContractEigenDACertVerifierV1V2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1V2(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV1V2, error) {
	contract, err := bindContractEigenDACertVerifierV1V2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1V2{ContractEigenDACertVerifierV1V2Caller: ContractEigenDACertVerifierV1V2Caller{contract: contract}, ContractEigenDACertVerifierV1V2Transactor: ContractEigenDACertVerifierV1V2Transactor{contract: contract}, ContractEigenDACertVerifierV1V2Filterer: ContractEigenDACertVerifierV1V2Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV1V2Caller creates a new read-only instance of ContractEigenDACertVerifierV1V2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1V2Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV1V2Caller, error) {
	contract, err := bindContractEigenDACertVerifierV1V2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1V2Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1V2Transactor creates a new write-only instance of ContractEigenDACertVerifierV1V2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1V2Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV1V2Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV1V2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1V2Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1V2Filterer creates a new log filterer instance of ContractEigenDACertVerifierV1V2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1V2Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV1V2Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV1V2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1V2Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV1V2 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV1V2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV1V2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1V2.Contract.ContractEigenDACertVerifierV1V2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.ContractEigenDACertVerifierV1V2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.ContractEigenDACertVerifierV1V2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1V2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) EigenDABatchMetadataStorageV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "eigenDABatchMetadataStorageV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) EigenDARelayRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "eigenDARelayRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) EigenDASignatureVerifierV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "eigenDASignatureVerifierV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) EigenDAThresholdRegistryV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "eigenDAThresholdRegistryV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) EigenDAThresholdRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "eigenDAThresholdRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifierV1V2.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifierV1V2.CallOpts, signedBatch)
}

// GetQuorumNumbersRequired is a free data retrieval call binding the contract method 0xefd57acb.
//
// Solidity: function getQuorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) GetQuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "getQuorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetQuorumNumbersRequired is a free data retrieval call binding the contract method 0xefd57acb.
//
// Solidity: function getQuorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) GetQuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.GetQuorumNumbersRequired(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// GetQuorumNumbersRequired is a free data retrieval call binding the contract method 0xefd57acb.
//
// Solidity: function getQuorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) GetQuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.GetQuorumNumbersRequired(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) OperatorStateRetrieverV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "operatorStateRetrieverV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) RegistryCoordinatorV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "registryCoordinatorV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) SecurityThresholdsAdversary(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "securityThresholdsAdversary")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) SecurityThresholdsAdversary() (uint8, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.SecurityThresholdsAdversary(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) SecurityThresholdsAdversary() (uint8, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.SecurityThresholdsAdversary(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) SecurityThresholdsConfirmation(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "securityThresholdsConfirmation")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) SecurityThresholdsConfirmation() (uint8, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.SecurityThresholdsConfirmation(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) SecurityThresholdsConfirmation() (uint8, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.SecurityThresholdsConfirmation(&_ContractEigenDACertVerifierV1V2.CallOpts)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertSecurityParams(opts *bind.CallOpts, blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertSecurityParams", blobParams, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifierV1V2.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifierV1V2.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertSecurityParams0(opts *bind.CallOpts, version uint16, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertSecurityParams0", version, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifierV1V2.CallOpts, version, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifierV1V2.CallOpts, version, securityThresholds)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertV1ForZkProof(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertV1ForZkProof", blobHeader, blobVerificationProof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertV1ForZkProof(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV1ForZkProof(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1ForZkProof is a free data retrieval call binding the contract method 0xf88adbba.
//
// Solidity: function verifyDACertV1ForZkProof(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertV1ForZkProof(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV1ForZkProof(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertV2ForZKProof(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertV2ForZKProof", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertV2ForZKProof(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertV2ForZKProof(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertV2FromSignedBatch", signedBatch, blobInclusionInfo)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV1V2.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV1V2.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "verifyDACertsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) VerifyDACertsV1(blobHeaders []BlobHeader, blobVerificationProofs []BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1V2.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeaders, blobVerificationProofs)
}
