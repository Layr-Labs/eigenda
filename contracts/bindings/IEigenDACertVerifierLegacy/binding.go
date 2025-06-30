// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDACertVerifierLegacy

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

// EigenDATypesV1BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// EigenDATypesV1BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchMetadata struct {
	BatchHeader             EigenDATypesV1BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// EigenDATypesV1BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []EigenDATypesV1QuorumBlobParam
}

// EigenDATypesV1BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  EigenDATypesV1BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// EigenDATypesV1NonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1NonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// EigenDATypesV1QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// EigenDATypesV1SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// EigenDATypesV1VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// EigenDATypesV2Attestation is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2Attestation struct {
	NonSignerPubkeys []BN254G1Point
	QuorumApks       []BN254G1Point
	Sigma            BN254G1Point
	ApkG2            BN254G2Point
	QuorumNumbers    []uint32
}

// EigenDATypesV2BatchHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2BatchHeaderV2 struct {
	BatchRoot            [32]byte
	ReferenceBlockNumber uint32
}

// EigenDATypesV2BlobCertificate is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2BlobCertificate struct {
	BlobHeader EigenDATypesV2BlobHeaderV2
	Signature  []byte
	RelayKeys  []uint32
}

// EigenDATypesV2BlobCommitment is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2BlobCommitment struct {
	Commitment       BN254G1Point
	LengthCommitment BN254G2Point
	LengthProof      BN254G2Point
	Length           uint32
}

// EigenDATypesV2BlobHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2BlobHeaderV2 struct {
	Version           uint16
	QuorumNumbers     []byte
	Commitment        EigenDATypesV2BlobCommitment
	PaymentHeaderHash [32]byte
}

// EigenDATypesV2BlobInclusionInfo is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2BlobInclusionInfo struct {
	BlobCertificate EigenDATypesV2BlobCertificate
	BlobIndex       uint32
	InclusionProof  []byte
}

// EigenDATypesV2SignedBatch is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2SignedBatch struct {
	BatchHeader EigenDATypesV2BatchHeaderV2
	Attestation EigenDATypesV2Attestation
}

// ContractIEigenDACertVerifierLegacyMetaData contains all meta data concerning the ContractIEigenDACertVerifierLegacy contract.
var ContractIEigenDACertVerifierLegacyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.SignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.Attestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.SignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.Attestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
}

// ContractIEigenDACertVerifierLegacyABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDACertVerifierLegacyMetaData.ABI instead.
var ContractIEigenDACertVerifierLegacyABI = ContractIEigenDACertVerifierLegacyMetaData.ABI

// ContractIEigenDACertVerifierLegacy is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacy struct {
	ContractIEigenDACertVerifierLegacyCaller     // Read-only binding to the contract
	ContractIEigenDACertVerifierLegacyTransactor // Write-only binding to the contract
	ContractIEigenDACertVerifierLegacyFilterer   // Log filterer for contract events
}

// ContractIEigenDACertVerifierLegacyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertVerifierLegacyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertVerifierLegacyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDACertVerifierLegacyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertVerifierLegacySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDACertVerifierLegacySession struct {
	Contract     *ContractIEigenDACertVerifierLegacy // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ContractIEigenDACertVerifierLegacyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDACertVerifierLegacyCallerSession struct {
	Contract *ContractIEigenDACertVerifierLegacyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// ContractIEigenDACertVerifierLegacyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDACertVerifierLegacyTransactorSession struct {
	Contract     *ContractIEigenDACertVerifierLegacyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// ContractIEigenDACertVerifierLegacyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacyRaw struct {
	Contract *ContractIEigenDACertVerifierLegacy // Generic contract binding to access the raw methods on
}

// ContractIEigenDACertVerifierLegacyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacyCallerRaw struct {
	Contract *ContractIEigenDACertVerifierLegacyCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDACertVerifierLegacyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDACertVerifierLegacyTransactorRaw struct {
	Contract *ContractIEigenDACertVerifierLegacyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDACertVerifierLegacy creates a new instance of ContractIEigenDACertVerifierLegacy, bound to a specific deployed contract.
func NewContractIEigenDACertVerifierLegacy(address common.Address, backend bind.ContractBackend) (*ContractIEigenDACertVerifierLegacy, error) {
	contract, err := bindContractIEigenDACertVerifierLegacy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacy{ContractIEigenDACertVerifierLegacyCaller: ContractIEigenDACertVerifierLegacyCaller{contract: contract}, ContractIEigenDACertVerifierLegacyTransactor: ContractIEigenDACertVerifierLegacyTransactor{contract: contract}, ContractIEigenDACertVerifierLegacyFilterer: ContractIEigenDACertVerifierLegacyFilterer{contract: contract}}, nil
}

// NewContractIEigenDACertVerifierLegacyCaller creates a new read-only instance of ContractIEigenDACertVerifierLegacy, bound to a specific deployed contract.
func NewContractIEigenDACertVerifierLegacyCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDACertVerifierLegacyCaller, error) {
	contract, err := bindContractIEigenDACertVerifierLegacy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyCaller{contract: contract}, nil
}

// NewContractIEigenDACertVerifierLegacyTransactor creates a new write-only instance of ContractIEigenDACertVerifierLegacy, bound to a specific deployed contract.
func NewContractIEigenDACertVerifierLegacyTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDACertVerifierLegacyTransactor, error) {
	contract, err := bindContractIEigenDACertVerifierLegacy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyTransactor{contract: contract}, nil
}

// NewContractIEigenDACertVerifierLegacyFilterer creates a new log filterer instance of ContractIEigenDACertVerifierLegacy, bound to a specific deployed contract.
func NewContractIEigenDACertVerifierLegacyFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDACertVerifierLegacyFilterer, error) {
	contract, err := bindContractIEigenDACertVerifierLegacy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyFilterer{contract: contract}, nil
}

// bindContractIEigenDACertVerifierLegacy binds a generic wrapper to an already deployed contract.
func bindContractIEigenDACertVerifierLegacy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDACertVerifierLegacyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.ContractIEigenDACertVerifierLegacyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.ContractIEigenDACertVerifierLegacyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.ContractIEigenDACertVerifierLegacyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.contract.Transact(opts, method, params...)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(EigenDATypesV1VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1VersionedBlobParams)).(*EigenDATypesV1VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetBlobParams(&_ContractIEigenDACertVerifierLegacy.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetBlobParams(&_ContractIEigenDACertVerifierLegacy.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetIsQuorumRequired(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetIsQuorumRequired(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(EigenDATypesV1NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1NonSignerStakesAndSignature)).(*EigenDATypesV1NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) GetNonSignerStakesAndSignature(signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetNonSignerStakesAndSignature(&_ContractIEigenDACertVerifierLegacy.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) GetNonSignerStakesAndSignature(signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetNonSignerStakesAndSignature(&_ContractIEigenDACertVerifierLegacy.CallOpts, signedBatch)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractIEigenDACertVerifierLegacy.CallOpts, quorumNumber)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumAdversaryThresholdPercentages(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumAdversaryThresholdPercentages(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumConfirmationThresholdPercentages(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumConfirmationThresholdPercentages(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumNumbersRequired(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.QuorumNumbersRequired(&_ContractIEigenDACertVerifierLegacy.CallOpts)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertSecurityParams(opts *bind.CallOpts, blobParams EigenDATypesV1VersionedBlobParams, securityThresholds EigenDATypesV1SecurityThresholds) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertSecurityParams", blobParams, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertSecurityParams(blobParams EigenDATypesV1VersionedBlobParams, securityThresholds EigenDATypesV1SecurityThresholds) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertSecurityParams(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertSecurityParams(blobParams EigenDATypesV1VersionedBlobParams, securityThresholds EigenDATypesV1SecurityThresholds) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertSecurityParams(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertSecurityParams0(opts *bind.CallOpts, version uint16, securityThresholds EigenDATypesV1SecurityThresholds) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertSecurityParams0", version, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertSecurityParams0(version uint16, securityThresholds EigenDATypesV1SecurityThresholds) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertSecurityParams0(&_ContractIEigenDACertVerifierLegacy.CallOpts, version, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertSecurityParams0(version uint16, securityThresholds EigenDATypesV1SecurityThresholds) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertSecurityParams0(&_ContractIEigenDACertVerifierLegacy.CallOpts, version, securityThresholds)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertV1(opts *bind.CallOpts, blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV1(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV1(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertV2(opts *bind.CallOpts, batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertV2(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2(&_ContractIEigenDACertVerifierLegacy.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertV2(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2(&_ContractIEigenDACertVerifierLegacy.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertV2ForZKProof(opts *bind.CallOpts, batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertV2ForZKProof", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertV2ForZKProof(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2ForZKProof(&_ContractIEigenDACertVerifierLegacy.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertV2ForZKProof(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2ForZKProof(&_ContractIEigenDACertVerifierLegacy.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertV2FromSignedBatch(opts *bind.CallOpts, signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertV2FromSignedBatch", signedBatch, blobInclusionInfo)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertV2FromSignedBatch(signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2FromSignedBatch(&_ContractIEigenDACertVerifierLegacy.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertV2FromSignedBatch(signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertV2FromSignedBatch(&_ContractIEigenDACertVerifierLegacy.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCaller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractIEigenDACertVerifierLegacy.contract.Call(opts, &out, "verifyDACertsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacySession) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertsV1(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyCallerSession) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractIEigenDACertVerifierLegacy.Contract.VerifyDACertsV1(&_ContractIEigenDACertVerifierLegacy.CallOpts, blobHeaders, blobVerificationProofs)
}

// ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated)
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
func (it *ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 EigenDATypesV1SecurityThresholds
	NewDefaultSecurityThresholdsV2      EigenDATypesV1SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractIEigenDACertVerifierLegacy.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated)
				if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractIEigenDACertVerifierLegacyDefaultSecurityThresholdsV2Updated)
	if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractIEigenDACertVerifierLegacy.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractIEigenDACertVerifierLegacyQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractIEigenDACertVerifierLegacy.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractIEigenDACertVerifierLegacyQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated)
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
func (it *ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdatedIterator{contract: _ContractIEigenDACertVerifierLegacy.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated)
				if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated, error) {
	event := new(ContractIEigenDACertVerifierLegacyQuorumNumbersRequiredUpdated)
	if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator struct {
	Event *ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded)
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
		it.Event = new(ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded)
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
func (it *ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractIEigenDACertVerifierLegacy contract.
type ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams EigenDATypesV1VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertVerifierLegacyVersionedBlobParamsAddedIterator{contract: _ContractIEigenDACertVerifierLegacy.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractIEigenDACertVerifierLegacy.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded)
				if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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
func (_ContractIEigenDACertVerifierLegacy *ContractIEigenDACertVerifierLegacyFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded, error) {
	event := new(ContractIEigenDACertVerifierLegacyVersionedBlobParamsAdded)
	if err := _ContractIEigenDACertVerifierLegacy.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
