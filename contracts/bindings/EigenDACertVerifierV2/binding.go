// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV2

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

// EigenDATypesV1SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
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

// ContractEigenDACertVerifierV2MetaData contains all meta data concerning the ContractEigenDACertVerifierV2 contract.
var ContractEigenDACertVerifierV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDASignatureVerifierV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_operatorStateRetrieverV2\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinatorV2\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"_quorumNumbersRequiredV2\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.SignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.Attestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetrieverV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequiredV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2ForZKProof\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.SignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.Attestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
}

// ContractEigenDACertVerifierV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.ABI instead.
var ContractEigenDACertVerifierV2ABI = ContractEigenDACertVerifierV2MetaData.ABI

// ContractEigenDACertVerifierV2 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2 struct {
	ContractEigenDACertVerifierV2Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV2Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV2Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV2Session struct {
	Contract     *ContractEigenDACertVerifierV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV2CallerSession struct {
	Contract *ContractEigenDACertVerifierV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDACertVerifierV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV2TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2Raw struct {
	Contract *ContractEigenDACertVerifierV2 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2CallerRaw struct {
	Contract *ContractEigenDACertVerifierV2Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV2TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV2 creates a new instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV2, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2{ContractEigenDACertVerifierV2Caller: ContractEigenDACertVerifierV2Caller{contract: contract}, ContractEigenDACertVerifierV2Transactor: ContractEigenDACertVerifierV2Transactor{contract: contract}, ContractEigenDACertVerifierV2Filterer: ContractEigenDACertVerifierV2Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV2Caller creates a new read-only instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV2Caller, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV2Transactor creates a new write-only instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV2Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV2Filterer creates a new log filterer instance of ContractEigenDACertVerifierV2, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV2Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV2Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV2Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV2 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.ContractEigenDACertVerifierV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV2.Contract.contract.Transact(opts, method, params...)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) EigenDASignatureVerifierV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "eigenDASignatureVerifierV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDASignatureVerifierV2 is a free data retrieval call binding the contract method 0x154b9e86.
//
// Solidity: function eigenDASignatureVerifierV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) EigenDASignatureVerifierV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDASignatureVerifierV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) EigenDAThresholdRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "eigenDAThresholdRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDAThresholdRegistryV2 is a free data retrieval call binding the contract method 0x17f3578e.
//
// Solidity: function eigenDAThresholdRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) EigenDAThresholdRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDAThresholdRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(EigenDATypesV1NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1NonSignerStakesAndSignature)).(*EigenDATypesV1NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) GetNonSignerStakesAndSignature(signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifierV2.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]))
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) GetNonSignerStakesAndSignature(signedBatch EigenDATypesV2SignedBatch) (EigenDATypesV1NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifierV2.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) OperatorStateRetrieverV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "operatorStateRetrieverV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// OperatorStateRetrieverV2 is a free data retrieval call binding the contract method 0x5df1f618.
//
// Solidity: function operatorStateRetrieverV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) OperatorStateRetrieverV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.OperatorStateRetrieverV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) QuorumNumbersRequiredV2(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "quorumNumbersRequiredV2")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifierV2.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// QuorumNumbersRequiredV2 is a free data retrieval call binding the contract method 0xb74d7871.
//
// Solidity: function quorumNumbersRequiredV2() view returns(bytes)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) QuorumNumbersRequiredV2() ([]byte, error) {
	return _ContractEigenDACertVerifierV2.Contract.QuorumNumbersRequiredV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) RegistryCoordinatorV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "registryCoordinatorV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// RegistryCoordinatorV2 is a free data retrieval call binding the contract method 0x5fafa482.
//
// Solidity: function registryCoordinatorV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) RegistryCoordinatorV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.RegistryCoordinatorV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) SecurityThresholdsV2(opts *bind.CallOpts) (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "securityThresholdsV2")

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
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsV2 is a free data retrieval call binding the contract method 0xed0450ae.
//
// Solidity: function securityThresholdsV2() view returns(uint8 confirmationThreshold, uint8 adversaryThreshold)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) SecurityThresholdsV2() (struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2(opts *bind.CallOpts, batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2ForZKProof(opts *bind.CallOpts, batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertV2ForZKProof", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2ForZKProof(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2ForZKProof is a free data retrieval call binding the contract method 0x415ef614.
//
// Solidity: function verifyDACertV2ForZKProof((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2ForZKProof(batchHeader EigenDATypesV2BatchHeaderV2, blobInclusionInfo EigenDATypesV2BlobInclusionInfo, nonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2ForZKProof(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2FromSignedBatch(opts *bind.CallOpts, signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertV2FromSignedBatch", signedBatch, blobInclusionInfo)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2FromSignedBatch(signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2FromSignedBatch(signedBatch EigenDATypesV2SignedBatch, blobInclusionInfo EigenDATypesV2BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch, blobInclusionInfo)
}
