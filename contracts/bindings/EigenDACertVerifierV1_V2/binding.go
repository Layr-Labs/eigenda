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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorage\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"},{\"name\":\"_eigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetriever\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetrieverV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsAdversary\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsConfirmation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structQuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structBlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structBatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BatchMetadataMismatch\",\"inputs\":[{\"name\":\"actualHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expectedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfirmationThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"QuorumNumberMismatch\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"errParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StakeThresholdNotMet\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"requiredThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"actualThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}]",
	Bin: "0x6101a06040523480156200001257600080fd5b5060405162003e7c38038062003e7c8339810160408190526200003591620001d9565b6001600160a01b03808916608081905281891660a05260c05280871660e05280861661010052808516610120528316610140526020820151825189918891889188918891889160ff918216911611620000a1576040516308a6997560e01b815260040160405180910390fd5b805160ff90811661016052602090910151166101805250620002fe9b505050505050505050505050565b6001600160a01b0381168114620000e157600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620001255762000125620000e4565b604052919050565b805160ff811681146200013f57600080fd5b919050565b600082601f8301126200015657600080fd5b81516001600160401b03811115620001725762000172620000e4565b602062000188601f8301601f19168201620000fa565b82815285828487010111156200019d57600080fd5b60005b83811015620001bd578581018301518282018401528201620001a0565b83811115620001cf5760008385840101525b5095945050505050565b600080600080600080600080888a03610120811215620001f857600080fd5b89516200020581620000cb565b60208b01519099506200021881620000cb565b60408b01519098506200022b81620000cb565b60608b01519097506200023e81620000cb565b60808b01519096506200025181620000cb565b60a08b01519095506200026481620000cb565b9350604060bf19820112156200027957600080fd5b50604080519081016001600160401b038082118383101715620002a057620002a0620000e4565b81604052620002b260c08d016200012d565b8352620002c260e08d016200012d565b60208401526101008c015192945080831115620002de57600080fd5b5050620002ee8b828c0162000144565b9150509295985092959890939650565b60805160a05160c05160e0516101005161012051610140516101605161018051613a96620003e6600039600081816102e50152610ab50152600081816103530152610a900152600081816102010152818161042201526106cd0152600081816101da0152818161040101526106ac01526000818161024e015281816103e001526104e7015260008181610135015281816103bf01526104c60152600081816101790152818161039e015281816104a40152818161055b015281816105fd01528181610ae00152610ba80152600081816102980152610c08015260006101b30152613a966000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c806382c216e7116100a2578063ccb7cd0d11610071578063ccb7cd0d146102cd578063dea610a9146102e0578063efd57acb14610319578063f25de3f81461032e578063fd1744841461034e57600080fd5b806382c216e7146102495780638ec28be614610270578063a9c823e114610293578063c084bcbf146102ba57600080fd5b80634cff90c4116100e95780634cff90c4146101ae5780635df1f618146101d55780635fafa482146101fc5780637d644cad14610223578063813c2eb01461023657600080fd5b8063143eb4d91461011b578063154b9e861461013057806317f3578e14610174578063421c02221461019b575b600080fd5b61012e6101293660046122fa565b610375565b005b6101577f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6101577f000000000000000000000000000000000000000000000000000000000000000081565b61012e6101a9366004612388565b610396565b6101577f000000000000000000000000000000000000000000000000000000000000000081565b6101577f000000000000000000000000000000000000000000000000000000000000000081565b6101577f000000000000000000000000000000000000000000000000000000000000000081565b61012e6102313660046123eb565b610468565b61012e6102443660046124da565b61049c565b6101577f000000000000000000000000000000000000000000000000000000000000000081565b61028361027e3660046124da565b610553565b604051901515815260200161016b565b6101577f000000000000000000000000000000000000000000000000000000000000000081565b6102836102c83660046123eb565b6105b4565b61012e6102db366004612595565b6105f5565b6103077f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff909116815260200161016b565b610321610690565b60405161016b9190612618565b61034161033c366004612632565b61069f565b60405161016b9190612852565b6103077f000000000000000000000000000000000000000000000000000000000000000081565b6000806103828484610701565b9150915061039082826107e1565b50505050565b6000806103827f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061044a8a612a3e565b6104538a612b35565b61045b610a70565b610463610adc565b610b64565b60008061048e610476610ba4565b61047f85610c04565b8686610489610adc565b610ca8565b915091506103908282610d53565b60008061053d7f00000000000000000000000000000000000000000000000000000000000000005b7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610515368b90038b018b612ce0565b61051e8a612b35565b6105278a612d7b565b61052f610a70565b610537610adc565b8b611023565b9150915061054b82826107e1565b505050505050565b60008061057f7f00000000000000000000000000000000000000000000000000000000000000006104c4565b509050600081600681111561059657610596612ea3565b14156105a65760019150506105ac565b60009150505b949350505050565b6000806105c2610476610ba4565b509050600081600a8111156105d9576105d9612ea3565b14156105e95760019150506105ef565b60009150505b92915050565b6000806103827f0000000000000000000000000000000000000000000000000000000000000000604051632ecfe72b60e01b815261ffff871660048201526001600160a01b039190911690632ecfe72b90602401606060405180830381865afa158015610666573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061068a9190612eb9565b84610701565b606061069a610adc565b905090565b6106a76120e4565b6106fa7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006106f585612a3e565b6111ae565b5092915050565b6000606060008360200151846000015161071b9190612f40565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406107449190612f79565b61074e9190612f79565b61075a90612710612f8d565b6107649190612fa4565b865190915060009061077890612710612fc3565b63ffffffff1690508082106107a557600060405180602001604052806000815250945094505050506107da565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050505b9250929050565b60008260068111156107f5576107f5612ea3565b14156107ff575050565b600182600681111561081357610813612ea3565b1415610869576000806000838060200190518101906108329190612fef565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b600282600681111561087d5761087d612ea3565b14156108c5576000808280602001905181019061089a919061301d565b604051631b00235d60e01b815260ff8084166004830152821660248201529193509150604401610860565b60038260068111156108d9576108d9612ea3565b141561092d576000806000838060200190518101906108f89190612fef565b6040516001626dc9ad60e11b031981526004810184905260248101839052604481018290529295509093509150606401610860565b600482600681111561094157610941612ea3565b1415610986576000808280602001905181019061095e919061304c565b604051634a47030360e11b815260048101839052602481018290529193509150604401610860565b600582600681111561099a5761099a612ea3565b14156109df57600080828060200190518101906109b7919061304c565b60405163114b085b60e21b815260048101839052602481018290529193509150604401610860565b60068260068111156109f3576109f3612ea3565b1415610a3357600081806020019051810190610a0f9190613070565b6040516309efaa0b60e41b815263ffffffff82166004820152909150602401610860565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b6044820152606401610860565b6040805180820182526000808252602091820152815180830190925260ff7f0000000000000000000000000000000000000000000000000000000000000000811683527f0000000000000000000000000000000000000000000000000000000000000000169082015290565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b3c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261069a919081019061308d565b60006060600080610b768a8a8a6111ae565b91509150610b8f8d8d8d8b600001518b878c8c89611023565b9350935050505b995099975050505050505050565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b3c573d6000803e3d6000fd5b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663eccbbfc9610c4260208501856130fa565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa158015610c84573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ef9190613117565b60006060610cb686856113c3565b9092509050600082600a811115610ccf57610ccf612ea3565b14610cd957610d49565b610ce3858561143a565b9092509050600082600a811115610cfc57610cfc612ea3565b14610d0657610d49565b6000610d1388878761157a565b91945092509050600083600a811115610d2e57610d2e612ea3565b14610d395750610d49565b610d438482611620565b92509250505b9550959350505050565b600082600a811115610d6757610d67612ea3565b1415610d71575050565b600182600a811115610d8557610d85612ea3565b1415610dca5760008082806020019051810190610da2919061304c565b6040516315d615c960e21b815260048101839052602481018290529193509150604401610860565b600282600a811115610dde57610dde612ea3565b1415610dfd576000806000838060200190518101906108329190612fef565b600382600a811115610e1157610e11612ea3565b1415610e595760008082806020019051810190610e2e919061301d565b6040516314fa310760e31b815260ff8084166004830152821660248201529193509150604401610860565b600482600a811115610e6d57610e6d612ea3565b1415610e8a576000808280602001905181019061089a919061301d565b600582600a811115610e9e57610e9e612ea3565b1415610ef357600080600083806020019051810190610ebd9190613130565b604051638aa11c4360e01b815260ff80851660048301528084166024830152821660448201529295509093509150606401610860565b600682600a811115610f0757610f07612ea3565b1415610f5c57600080600083806020019051810190610f269190613130565b60405163a4ad875560e01b815260ff80851660048301528084166024830152821660448201529295509093509150606401610860565b600782600a811115610f7057610f70612ea3565b1415610f8d57600080828060200190518101906109b7919061304c565b600a82600a811115610fa157610fa1612ea3565b1415610fbd57600081806020019051810190610a0f9190613070565b600882600a811115610fd157610fd1612ea3565b1415610ff25780604051638c59c92f60e01b81526004016108609190612618565b600982600a81111561100657611006612ea3565b1415610a33576000808280602001905181019061095e919061304c565b600060606110318888611670565b9092509050600082600681111561104a5761104a612ea3565b1461105457610b96565b6110668988600001516040015161172f565b9092509050600082600681111561107f5761107f612ea3565b1461108957610b96565b86515151604051632ecfe72b60e01b815261ffff9091166004820152611104906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa1580156110da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110fe9190612eb9565b86610701565b9092509050600082600681111561111d5761111d612ea3565b1461112757610b96565b60006111438b6111368b61184e565b868c602001518b8b611896565b91945092509050600083600681111561115e5761115e612ea3565b146111695750610b96565b8751516020015160009061117d90836119f7565b91955093509050600084600681111561119857611198612ea3565b146111a4575050610b96565b610b8f8682611a5d565b6111b66120e4565b60606000836020015160000151516001600160401b038111156111db576111db61218d565b604051908082528060200260200182016040528015611204578160200160208202803683370190505b50905060005b602085015151518110156112815761125485602001516000015182815181106112355761123561317d565b6020026020010151805160009081526020918201519091526040902090565b8282815181106112665761126661317d565b602090810291909101015261127a81613193565b905061120a565b5060005b846020015160800151518110156112ec578285602001516080015182815181106112b1576112b161317d565b60200260200101516040516020016112ca9291906131ae565b6040516020818303038152906040529250806112e590613193565b9050611285565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491611327918a9190889088906004016131e0565b600060405180830381865afa158015611344573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261136c9190810190613333565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b60006060816113e66113d8604086018661340b565b6113e19061342b565b611aad565b9050848114156114095750506040805160208101909152600080825291506107da565b60408051602081018390529081018690526001906060015b60405160208183030381529060405292509250506107da565b600060608161145061144b86613502565b611b1e565b905060008160405160200161146791815260200190565b6040516020818303038152906040528051906020012090506000858060400190611491919061340b565b61149b9080613621565b35905060006115046114b06060890189613637565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508692508791506114f9905060408c0160208d016130fa565b63ffffffff16611b31565b9050801561152b5760006040518060200160405280600081525095509550505050506107da565b600261153d6040890160208a016130fa565b6040805163ffffffff90921660208301528101859052606081018490526080015b60405160208183030381529060405295509550505050506107da565b60006060818061158c8684018761367d565b9050905060005b818110156116005760008060006115ad8b8b8b8788611b49565b91945092509050600083600a8111156115c8576115c8612ea3565b146115df5750909550935060009250611617915050565b600160ff82161b8617955050505080806115f890613193565b915050611593565b505060408051602081019091526000808252935091505b93509350939050565b60006060600061162f85611dda565b90508381168114156116545750506040805160208101909152600080825291506107da565b6040805160208101839052908101859052600790606001611421565b6000606060006116838460000151611f67565b905060008160405160200161169a91815260200190565b60405160208183030381529060405280519060200120905060008660000151905060006116d7876040015183858a6020015163ffffffff16611b31565b905080156116fe5760006040518060200160405280600081525095509550505050506107da565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260019060800161155e565b6000606060005b83518110156118335760006001600160a01b0316856001600160a01b031663b5a872da86848151811061176b5761176b61317d565b60200260200101516040518263ffffffff1660e01b815260040161179b919063ffffffff91909116815260200190565b602060405180830381865afa1580156117b8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117dc91906136c6565b6001600160a01b031614156118235760068482815181106117ff576117ff61317d565b6020026020010151604051602001611421919063ffffffff91909116815260200190565b61182c81613193565b9050611736565b50506040805160208101909152600080825291509250929050565b60008160405160200161187991908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b81526004016118cf94939291906136ef565b600060405180830381865afa1580156118ec573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261191491908101906137a2565b5090506000915060005b88518110156119d557856000015160ff16826020015182815181106119455761194561317d565b6020026020010151611957919061383e565b6001600160601b03166064836000015183815181106119785761197861317d565b60200260200101516001600160601b03166119939190612fa4565b106119c3576119c0838a83815181106119ae576119ae61317d565b0160200151600160f89190911c1b1790565b92505b806119cd81613193565b91505061191e565b5050604080516020810190915260008082529350915096509650969350505050565b600060606000611a0685611dda565b9050838116811415611a2b576040805160208101909152600080825293509150611a56565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b600060606000611a6c85611dda565b9050838116811415611a915750506040805160208101909152600080825291506107da565b6040805160208101839052908101859052600590606001611421565b60006105ef8260000151604051602001611ac79190613864565b60408051808303601f1901815282825280516020918201208682015187840151838601929092528484015260e01b6001600160e01b0319166060840152815160448185030181526064909301909152815191012090565b60008160405160200161187991906138c4565b600083611b3f868585611f8f565b1495945050505050565b600060608136611b5b8884018961367d565b87818110611b6b57611b6b61317d565b608002919091019150611b8390506020820182613969565b91506000611b946080890189613637565b87818110611ba457611ba461317d565b919091013560f81c915060009050611bbf60408a018a61340b565b611bc99080613621565b611bd7906020810190613637565b8360ff16818110611bea57611bea61317d565b919091013560f81c91505060ff84168114611c37576040805160ff958616602082015291909416818501528351808203850181526060909101909352506003935090915060009050611dcf565b6000611c496060850160408601613969565b90506000611c5d6040860160208701613969565b90508060ff168260ff1611611ca6576040805160ff93841660208201529190921681830152815180820383018152606090910190915260049650945060009350611dcf92505050565b60008d8760ff1681518110611cbd57611cbd61317d565b016020015160f81c905060ff8316811115611d1a576040805160ff808a1660208301528084169282019290925290841660608201526005906080016040516020818303038152906040526000985098509850505050505050611dcf565b6000611d2960408e018e61340b565b611d339080613621565b611d41906040810190613637565b8760ff16818110611d5457611d5461317d565b919091013560f81c91505060ff8416811015611db3576040805160ff808b166020830152808716928201929092529082166060820152600690608001604051602081830303815290604052600099509950995050505050505050611dcf565b5050604080516020810190915260008082529850965050505050505b955095509592505050565b600061010082511115611e635760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610860565b8151611e7157506000919050565b60008083600081518110611e8757611e8761317d565b0160200151600160f89190911c81901b92505b8451811015611f5e57848181518110611eb557611eb561317d565b0160200151600160f89190911c1b9150828211611f4a5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610860565b91811791611f5781613193565b9050611e9a565b50909392505050565b6000611f768260000151612092565b6020808401516040808601519051611879949301613986565b600060208451611f9f91906139bb565b156120265760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610860565b8260205b855181116120895761203d6002856139bb565b61205e57816000528086015160205260406000209150600284049350612077565b8086015160005281602052604060002091506002840493505b6120826020826139cf565b905061202a565b50949350505050565b60008160000151826020015183604001516040516020016120b5939291906139e7565b60408051601f198184030181528282528051602091820120606080870151928501919091529183015201611879565b60405180610100016040528060608152602001606081526020016060815260200161210d61214a565b815260200161212f604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b604051806040016040528061215d61216f565b815260200161216a61216f565b905290565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156121c5576121c561218d565b60405290565b604051606081016001600160401b03811182821017156121c5576121c561218d565b60405160a081016001600160401b03811182821017156121c5576121c561218d565b604051608081016001600160401b03811182821017156121c5576121c561218d565b60405161010081016001600160401b03811182821017156121c5576121c561218d565b604051601f8201601f191681016001600160401b038111828210171561227c5761227c61218d565b604052919050565b63ffffffff8116811461229657600080fd5b50565b80356122a481612284565b919050565b60ff8116811461229657600080fd5b6000604082840312156122ca57600080fd5b6122d26121a3565b905081356122df816122a9565b815260208201356122ef816122a9565b602082015292915050565b60008082840360a081121561230e57600080fd5b606081121561231c57600080fd5b506123256121cb565b833561233081612284565b8152602084013561234081612284565b60208201526040840135612353816122a9565b6040820152915061236784606085016122b8565b90509250929050565b60006060828403121561238257600080fd5b50919050565b6000806040838503121561239b57600080fd5b82356001600160401b03808211156123b257600080fd5b6123be86838701612370565b935060208501359150808211156123d457600080fd5b506123e185828601612370565b9150509250929050565b600080604083850312156123fe57600080fd5b82356001600160401b038082111561241557600080fd5b908401906080828703121561242957600080fd5b9092506020840135908082111561243f57600080fd5b50830160a0818603121561245257600080fd5b809150509250929050565b60006001600160401b038211156124765761247661218d565b50601f01601f191660200190565b600082601f83011261249557600080fd5b81356124a86124a38261245d565b612254565b8181528460208386010111156124bd57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a08112156124f157600080fd5b60408112156124ff57600080fd5b5084935060408501356001600160401b038082111561251d57600080fd5b61252988838901612370565b9450606087013591508082111561253f57600080fd5b90860190610180828903121561255457600080fd5b9092506080860135908082111561256a57600080fd5b5061257787828801612484565b91505092959194509250565b803561ffff811681146122a457600080fd5b600080606083850312156125a857600080fd5b6125b183612583565b915061236784602085016122b8565b60005b838110156125db5781810151838201526020016125c3565b838111156103905750506000910152565b600081518084526126048160208601602086016125c0565b601f01601f19169290920160200192915050565b60208152600061262b60208301846125ec565b9392505050565b60006020828403121561264457600080fd5b81356001600160401b0381111561265a57600080fd5b6105ac84828501612370565b600081518084526020808501945080840160005b8381101561269c57815163ffffffff168752958201959082019060010161267a565b509495945050505050565b600081518084526020808501945080840160005b8381101561269c576126d887835180518252602090810151910152565b60409690960195908201906001016126bb565b8060005b60028110156103905781518452602093840193909101906001016126ef565b6127198282516126eb565b602081015161272b60408401826126eb565b505050565b600081518084526020808501808196508360051b8101915082860160005b85811015612778578284038952612766848351612666565b9885019893509084019060010161274e565b5091979650505050505050565b6000610180825181855261279b82860182612666565b915050602083015184820360208601526127b582826126a7565b915050604083015184820360408601526127cf82826126a7565b91505060608301516127e4606086018261270e565b506080830151805160e08601526020015161010085015260a08301518482036101208601526128138282612666565b91505060c083015184820361014086015261282e8282612666565b91505060e08301518482036101608601526128498282612730565b95945050505050565b60208152600061262b6020830184612785565b60006040828403121561287757600080fd5b61287f6121a3565b90508135815260208201356122ef81612284565b60006001600160401b038211156128ac576128ac61218d565b5060051b60200190565b6000604082840312156128c857600080fd5b6128d06121a3565b9050813581526020820135602082015292915050565b600082601f8301126128f757600080fd5b813560206129076124a383612893565b82815260069290921b8401810191818101908684111561292657600080fd5b8286015b8481101561294a5761293c88826128b6565b83529183019160400161292a565b509695505050505050565b600082601f83011261296657600080fd5b61296e6121a3565b80604084018581111561298057600080fd5b845b8181101561299a578035845260209384019301612982565b509095945050505050565b6000608082840312156129b757600080fd5b6129bf6121a3565b90506129cb8383612955565b81526122ef8360408401612955565b600082601f8301126129eb57600080fd5b813560206129fb6124a383612893565b82815260059290921b84018101918181019086841115612a1a57600080fd5b8286015b8481101561294a578035612a3181612284565b8352918301918301612a1e565b600060608236031215612a5057600080fd5b612a586121a3565b612a623684612865565b815260408301356001600160401b0380821115612a7e57600080fd5b81850191506101208236031215612a9457600080fd5b612a9c6121ed565b823582811115612aab57600080fd5b612ab7368286016128e6565b825250602083013582811115612acc57600080fd5b612ad8368286016128e6565b602083015250612aeb36604085016128b6565b6040820152612afd36608085016129a5565b606082015261010083013582811115612b1557600080fd5b612b21368286016129da565b608083015250602084015250909392505050565b600060608236031215612b4757600080fd5b612b4f6121cb565b82356001600160401b0380821115612b6657600080fd5b818501915060608236031215612b7b57600080fd5b612b836121cb565b823582811115612b9257600080fd5b8301368190036101c0811215612ba757600080fd5b612baf61220f565b612bb883612583565b815260208084013586811115612bcd57600080fd5b612bd936828701612484565b83830152506040610160603f1985011215612bf357600080fd5b612bfb61220f565b9350612c09368287016128b6565b8452612c1836608087016129a5565b82850152612c2a3661010087016129a5565b81850152610180850135612c3d81612284565b8060608601525083818401526101a0850135606084015282865281880135945086851115612c6a57600080fd5b612c7636868a01612484565b8287015280880135945086851115612c8d57600080fd5b612c9936868a016129da565b81870152858952612cab828c01612299565b828a0152808b0135975086881115612cc257600080fd5b612cce36898d01612484565b90890152509598975050505050505050565b600060408284031215612cf257600080fd5b61262b8383612865565b600082601f830112612d0d57600080fd5b81356020612d1d6124a383612893565b82815260059290921b84018101918181019086841115612d3c57600080fd5b8286015b8481101561294a5780356001600160401b03811115612d5f5760008081fd5b612d6d8986838b01016129da565b845250918301918301612d40565b60006101808236031215612d8e57600080fd5b612d96612231565b82356001600160401b0380821115612dad57600080fd5b612db9368387016129da565b83526020850135915080821115612dcf57600080fd5b612ddb368387016128e6565b60208401526040850135915080821115612df457600080fd5b612e00368387016128e6565b6040840152612e1236606087016129a5565b6060840152612e243660e087016128b6565b6080840152610120850135915080821115612e3e57600080fd5b612e4a368387016129da565b60a0840152610140850135915080821115612e6457600080fd5b612e70368387016129da565b60c0840152610160850135915080821115612e8a57600080fd5b50612e9736828601612cfc565b60e08301525092915050565b634e487b7160e01b600052602160045260246000fd5b600060608284031215612ecb57600080fd5b604051606081018181106001600160401b0382111715612eed57612eed61218d565b6040528251612efb81612284565b81526020830151612f0b81612284565b60208201526040830151612f1e816122a9565b60408201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff841680821015612f5a57612f5a612f2a565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b600082612f8857612f88612f63565b500490565b600082821015612f9f57612f9f612f2a565b500390565b6000816000190483118215151615612fbe57612fbe612f2a565b500290565b600063ffffffff80831681851681830481118215151615612fe657612fe6612f2a565b02949350505050565b60008060006060848603121561300457600080fd5b8351925060208401519150604084015190509250925092565b6000806040838503121561303057600080fd5b825161303b816122a9565b6020840151909250612452816122a9565b6000806040838503121561305f57600080fd5b505080516020909101519092909150565b60006020828403121561308257600080fd5b815161262b81612284565b60006020828403121561309f57600080fd5b81516001600160401b038111156130b557600080fd5b8201601f810184136130c657600080fd5b80516130d46124a38261245d565b8181528560208385010111156130e957600080fd5b6128498260208301602086016125c0565b60006020828403121561310c57600080fd5b813561262b81612284565b60006020828403121561312957600080fd5b5051919050565b60008060006060848603121561314557600080fd5b8351613150816122a9565b6020850151909350613161816122a9565b6040850151909250613172816122a9565b809150509250925092565b634e487b7160e01b600052603260045260246000fd5b60006000198214156131a7576131a7612f2a565b5060010190565b600083516131c08184602088016125c0565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff8616818401526080604084015261320e60808401866125ec565b838103606085015284518082528286019183019060005b8181101561324157835183529284019291840191600101613225565b50909998505050505050505050565b600082601f83011261326157600080fd5b815160206132716124a383612893565b82815260059290921b8401810191818101908684111561329057600080fd5b8286015b8481101561294a5780516132a781612284565b8352918301918301613294565b600082601f8301126132c557600080fd5b815160206132d56124a383612893565b82815260059290921b840181019181810190868411156132f457600080fd5b8286015b8481101561294a5780516001600160401b038111156133175760008081fd5b6133258986838b0101613250565b8452509183019183016132f8565b60006020828403121561334557600080fd5b81516001600160401b038082111561335c57600080fd5b908301906080828603121561337057600080fd5b61337861220f565b82518281111561338757600080fd5b61339387828601613250565b8252506020830151828111156133a857600080fd5b6133b487828601613250565b6020830152506040830151828111156133cc57600080fd5b6133d887828601613250565b6040830152506060830151828111156133f057600080fd5b6133fc878286016132b4565b60608301525095945050505050565b60008235605e1983360301811261342157600080fd5b9190910192915050565b60006060823603121561343d57600080fd5b6134456121cb565b82356001600160401b038082111561345c57600080fd5b81850191506080823603121561347157600080fd5b61347961220f565b8235815260208301358281111561348f57600080fd5b61349b36828601612484565b6020830152506040830135828111156134b357600080fd5b6134bf36828601612484565b604083015250606083013592506134d583612284565b826060820152808452505050602083013560208201526134f760408401612299565b604082015292915050565b6000608080833603121561351557600080fd5b61351d6121cb565b61352736856128b6565b815260408085013561353881612284565b6020818185015260609150818701356001600160401b0381111561355b57600080fd5b870136601f82011261356c57600080fd5b803561357a6124a382612893565b81815260079190911b8201830190838101903683111561359957600080fd5b928401925b8284101561360d578884360312156135b65760008081fd5b6135be61220f565b84356135c9816122a9565b8152848601356135d8816122a9565b81870152848801356135e9816122a9565b81890152848701356135fa81612284565b818801528252928801929084019061359e565b958701959095525093979650505050505050565b60008235607e1983360301811261342157600080fd5b6000808335601e1984360301811261364e57600080fd5b8301803591506001600160401b0382111561366857600080fd5b6020019150368190038213156107da57600080fd5b6000808335601e1984360301811261369457600080fd5b8301803591506001600160401b038211156136ae57600080fd5b6020019150600781901b36038213156107da57600080fd5b6000602082840312156136d857600080fd5b81516001600160a01b038116811461262b57600080fd5b84815260806020820152600061370860808301866125ec565b63ffffffff8516604084015282810360608401526137268185612785565b979650505050505050565b600082601f83011261374257600080fd5b815160206137526124a383612893565b82815260059290921b8401810191818101908684111561377157600080fd5b8286015b8481101561294a5780516001600160601b03811681146137955760008081fd5b8352918301918301613775565b600080604083850312156137b557600080fd5b82516001600160401b03808211156137cc57600080fd5b90840190604082870312156137e057600080fd5b6137e86121a3565b8251828111156137f757600080fd5b61380388828601613731565b82525060208301518281111561381857600080fd5b61382488828601613731565b602083015250809450505050602083015190509250929050565b60006001600160601b0380831681851681830481118215151615612fe657612fe6612f2a565b6020815281516020820152600060208301516080604084015261388a60a08401826125ec565b90506040840151601f198483030160608501526138a782826125ec565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b8084101561395a578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190613917565b509a9950505050505050505050565b60006020828403121561397b57600080fd5b813561262b816122a9565b83815260606020820152600061399f60608301856125ec565b82810360408401526139b18185612666565b9695505050505050565b6000826139ca576139ca612f63565b500690565b600082198211156139e2576139e2612f2a565b500190565b60006101a061ffff86168352806020840152613a05818401866125ec565b8451805160408601526020015160608501529150613a209050565b6020830151613a32608084018261270e565b506040830151613a4661010084018261270e565b5063ffffffff60608401511661018083015294935050505056fea264697066735822122029557e8fd819dbe40e6fc12819935aa5c98b1992eed866d96cf31fec32842bac64736f6c634300080c0033",
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

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) CheckDACertV1(opts *bind.CallOpts, blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "checkDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) CheckDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.CheckDACertV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// CheckDACertV1 is a free data retrieval call binding the contract method 0xc084bcbf.
//
// Solidity: function checkDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) CheckDACertV1(blobHeader BlobHeader, blobVerificationProof BlobVerificationProof) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.CheckDACertV1(&_ContractEigenDACertVerifierV1V2.CallOpts, blobHeader, blobVerificationProof)
}

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Caller) CheckDACertV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1V2.contract.Call(opts, &out, "checkDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2Session) CheckDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.CheckDACertV2(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV1V2 *ContractEigenDACertVerifierV1V2CallerSession) CheckDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV1V2.Contract.CheckDACertV2(&_ContractEigenDACertVerifierV1V2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
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
