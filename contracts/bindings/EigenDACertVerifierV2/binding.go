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

// BatchHeaderV2 is an auto generated low-level Go binding around an user-defined struct.
type BatchHeaderV2 struct {
	BatchRoot            [32]byte
	ReferenceBlockNumber uint32
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

// ContractEigenDACertVerifierV2MetaData contains all meta data concerning the ContractEigenDACertVerifierV2 contract.
var ContractEigenDACertVerifierV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDASignatureVerifierV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"_eigenDARelayRegistryV2\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"_operatorStateRetrieverV2\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"},{\"name\":\"_registryCoordinatorV2\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_securityThresholdsV2\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"checkDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifierV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNonSignerStakesAndSignature\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]}],\"outputs\":[{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorStateRetrieverV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractOperatorStateRetriever\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinatorV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsAdversary\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholdsConfirmation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"blobParams\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyDACertSecurityParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"securityThresholds\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structNonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV2FromSignedBatch\",\"inputs\":[{\"name\":\"signedBatch\",\"type\":\"tuple\",\"internalType\":\"structSignedBatch\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"attestation\",\"type\":\"tuple\",\"internalType\":\"structAttestation\",\"components\":[{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"quorumNumbers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structBlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structBlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structBlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidThresholdPercentages\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"type\":\"error\",\"name\":\"RelayKeyNotSet\",\"inputs\":[{\"name\":\"relayKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"gamma\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	Bin: "0x6101606040523480156200001257600080fd5b5060405162002d1438038062002d148339810160408190526200003591620000da565b6001600160a01b0380871660805285811660a05284811660c05283811660e0528216610100526020810151815160ff91821691161162000088576040516308a6997560e01b815260040160405180910390fd5b805160ff90811661012052602090910151166101405250620001c99350505050565b6001600160a01b0381168114620000c057600080fd5b50565b805160ff81168114620000d557600080fd5b919050565b60008060008060008086880360e0811215620000f557600080fd5b87516200010281620000aa565b60208901519097506200011581620000aa565b60408901519096506200012881620000aa565b60608901519095506200013b81620000aa565b60808901519094506200014e81620000aa565b92506040609f19820112156200016357600080fd5b50604080519081016001600160401b03811182821017156200019557634e487b7160e01b600052604160045260246000fd5b604052620001a660a08901620000c3565b8152620001b660c08901620000c3565b6020820152809150509295509295509295565b60805160a05160c05160e051610100516101205161014051612a886200028c6000396000818161022a0152610961015260008181610283015261093c01526000818161019301528181610352015261057901526000818161016c0152818161033101526105580152600081816101cd0152818161031001526103e301526000818160ee015281816102ef01526103c2015260008181610132015281816102ce015281816103a001528181610457015281816104b8015261098c0152612a886000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c8063813c2eb01161008c578063ccb7cd0d11610066578063ccb7cd0d14610212578063dea610a914610225578063f25de3f81461025e578063fd1744841461027e57600080fd5b8063813c2eb0146101b557806382c216e7146101c85780638ec28be6146101ef57600080fd5b8063143eb4d9146100d4578063154b9e86146100e957806317f3578e1461012d578063421c0222146101545780635df1f618146101675780635fafa4821461018e575b600080fd5b6100e76100e236600461179a565b6102a5565b005b6101107f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6101107f000000000000000000000000000000000000000000000000000000000000000081565b6100e7610162366004611828565b6102c6565b6101107f000000000000000000000000000000000000000000000000000000000000000081565b6101107f000000000000000000000000000000000000000000000000000000000000000081565b6100e76101c3366004611908565b610398565b6101107f000000000000000000000000000000000000000000000000000000000000000081565b6102026101fd366004611908565b61044f565b6040519015158152602001610124565b6100e76102203660046119c3565b6104b0565b61024c7f000000000000000000000000000000000000000000000000000000000000000081565b60405160ff9091168152602001610124565b61027161026c3660046119ee565b61054b565b6040516101249190611c0e565b61024c7f000000000000000000000000000000000000000000000000000000000000000081565b6000806102b284846105ad565b915091506102c0828261068d565b50505050565b6000806102b27f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061037a8a611e01565b6103838a611ef8565b61038b61091c565b610393610988565b610a15565b6000806104397f00000000000000000000000000000000000000000000000000000000000000005b7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610411368b90038b018b6120a3565b61041a8a611ef8565b6104238a61213e565b61042b61091c565b610433610988565b8b610a55565b91509150610447828261068d565b505050505050565b60008061047b7f00000000000000000000000000000000000000000000000000000000000000006103c0565b509050600081600681111561049257610492612266565b14156104a25760019150506104a8565b60009150505b949350505050565b6000806102b27f0000000000000000000000000000000000000000000000000000000000000000604051632ecfe72b60e01b815261ffff871660048201526001600160a01b039190911690632ecfe72b90602401606060405180830381865afa158015610521573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610545919061227c565b846105ad565b610553611584565b6105a67f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006105a185611e01565b610be0565b5092915050565b600060606000836020015184600001516105c79190612303565b60ff1690506000856020015163ffffffff16866040015160ff1683620f42406105f0919061233c565b6105fa919061233c565b61060690612710612350565b6106109190612367565b865190915060009061062490612710612386565b63ffffffff1690508082106106515760006040518060200160405280600081525094509450505050610686565b604080516020810185905290810183905260608101829052600390608001604051602081830303815290604052945094505050505b9250929050565b60008260068111156106a1576106a1612266565b14156106ab575050565b60018260068111156106bf576106bf612266565b1415610715576000806000838060200190518101906106de91906123b2565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b600282600681111561072957610729612266565b1415610771576000808280602001905181019061074691906123e0565b604051631b00235d60e01b815260ff808416600483015282166024820152919350915060440161070c565b600382600681111561078557610785612266565b14156107d9576000806000838060200190518101906107a491906123b2565b6040516001626dc9ad60e11b03198152600481018490526024810183905260448101829052929550909350915060640161070c565b60048260068111156107ed576107ed612266565b1415610832576000808280602001905181019061080a919061241a565b604051634a47030360e11b81526004810183905260248101829052919350915060440161070c565b600582600681111561084657610846612266565b141561088b5760008082806020019051810190610863919061241a565b60405163114b085b60e21b81526004810183905260248101829052919350915060440161070c565b600682600681111561089f5761089f612266565b14156108df576000818060200190518101906108bb919061243e565b6040516309efaa0b60e41b815263ffffffff8216600482015290915060240161070c565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b604482015260640161070c565b6040805180820182526000808252602091820152815180830190925260ff7f0000000000000000000000000000000000000000000000000000000000000000811683527f0000000000000000000000000000000000000000000000000000000000000000169082015290565b60607f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156109e8573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610a109190810190612487565b905090565b60006060600080610a278a8a8a610be0565b91509150610a408d8d8d8b600001518b878c8c89610a55565b9350935050505b995099975050505050505050565b60006060610a638888610df5565b90925090506000826006811115610a7c57610a7c612266565b14610a8657610a47565b610a9889886000015160400151610ecb565b90925090506000826006811115610ab157610ab1612266565b14610abb57610a47565b86515151604051632ecfe72b60e01b815261ffff9091166004820152610b36906001600160a01b038d1690632ecfe72b90602401606060405180830381865afa158015610b0c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b30919061227c565b866105ad565b90925090506000826006811115610b4f57610b4f612266565b14610b5957610a47565b6000610b758b610b688b611003565b868c602001518b8b61104b565b919450925090506000836006811115610b9057610b90612266565b14610b9b5750610a47565b87515160200151600090610baf90836111ac565b919550935090506000846006811115610bca57610bca612266565b14610bd6575050610a47565b610a408682611212565b610be8611584565b60606000836020015160000151516001600160401b03811115610c0d57610c0d61162d565b604051908082528060200260200182016040528015610c36578160200160208202803683370190505b50905060005b60208501515151811015610cb357610c868560200151600001518281518110610c6757610c676124f4565b6020026020010151805160009081526020918201519091526040902090565b828281518110610c9857610c986124f4565b6020908102919091010152610cac8161250a565b9050610c3c565b5060005b84602001516080015151811015610d1e57828560200151608001518281518110610ce357610ce36124f4565b6020026020010151604051602001610cfc929190612525565b604051602081830303815290604052925080610d179061250a565b9050610cb7565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f7491610d59918a919088908890600401612583565b600060405180830381865afa158015610d76573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610d9e91908101906126d6565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b600060606000610e088460000151611262565b9050600081604051602001610e1f91815260200190565b6040516020818303038152906040528051906020012090506000866000015190506000610e5c876040015183858a6020015163ffffffff1661128a565b90508015610e83576000604051806020016040528060008152509550955050505050610686565b6020808801516040805163ffffffff90921692820192909252908101849052606081018390526001906080016040516020818303038152906040529550955050505050610686565b6000606060005b8351811015610fe85760006001600160a01b0316856001600160a01b031663b5a872da868481518110610f0757610f076124f4565b60200260200101516040518263ffffffff1660e01b8152600401610f37919063ffffffff91909116815260200190565b602060405180830381865afa158015610f54573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f7891906127ae565b6001600160a01b03161415610fd8576006848281518110610f9b57610f9b6124f4565b6020026020010151604051602001610fbf919063ffffffff91909116815260200190565b6040516020818303038152906040529250925050610686565b610fe18161250a565b9050610ed2565b50506040805160208101909152600080825291509250929050565b60008160405160200161102e91908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b815260040161108494939291906127d7565b600060405180830381865afa1580156110a1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526110c9919081019061288a565b5090506000915060005b885181101561118a57856000015160ff16826020015182815181106110fa576110fa6124f4565b602002602001015161110c9190612926565b6001600160601b031660648360000151838151811061112d5761112d6124f4565b60200260200101516001600160601b03166111489190612367565b1061117857611175838a8381518110611163576111636124f4565b0160200151600160f89190911c1b1790565b92505b806111828161250a565b9150506110d3565b5050604080516020810190915260008082529350915096509650969350505050565b6000606060006111bb856112a2565b90508381168114156111e057604080516020810190915260008082529350915061120b565b6040805160208101929092528181018590528051808303820181526060909201905260049250905060005b9250925092565b600060606000611221856112a2565b9050838116811415611246575050604080516020810190915260008082529150610686565b6040805160208101839052908101859052600590606001610fbf565b6000611271826000015161142f565b602080840151604080860151905161102e94930161294c565b600083611298868585611481565b1495945050505050565b60006101008251111561132b5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a40161070c565b815161133957506000919050565b6000808360008151811061134f5761134f6124f4565b0160200151600160f89190911c81901b92505b84518110156114265784818151811061137d5761137d6124f4565b0160200151600160f89190911c1b91508282116114125760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a40161070c565b9181179161141f8161250a565b9050611362565b50909392505050565b6000816000015182602001518360400151604051602001611452939291906129ad565b60408051601f19818403018152828252805160209182012060608087015192850191909152918301520161102e565b6000602084516114919190612a26565b156115185760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a40161070c565b8260205b8551811161157b5761152f600285612a26565b61155057816000528086015160205260406000209150600284049350611569565b8086015160005281602052604060002091506002840493505b611574602082612a3a565b905061151c565b50949350505050565b6040518061010001604052806060815260200160608152602001606081526020016115ad6115ea565b81526020016115cf604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b60405180604001604052806115fd61160f565b815260200161160a61160f565b905290565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156116655761166561162d565b60405290565b604051606081016001600160401b03811182821017156116655761166561162d565b60405160a081016001600160401b03811182821017156116655761166561162d565b604051608081016001600160401b03811182821017156116655761166561162d565b60405161010081016001600160401b03811182821017156116655761166561162d565b604051601f8201601f191681016001600160401b038111828210171561171c5761171c61162d565b604052919050565b63ffffffff8116811461173657600080fd5b50565b803561174481611724565b919050565b60ff8116811461173657600080fd5b60006040828403121561176a57600080fd5b611772611643565b9050813561177f81611749565b8152602082013561178f81611749565b602082015292915050565b60008082840360a08112156117ae57600080fd5b60608112156117bc57600080fd5b506117c561166b565b83356117d081611724565b815260208401356117e081611724565b602082015260408401356117f381611749565b604082015291506118078460608501611758565b90509250929050565b60006060828403121561182257600080fd5b50919050565b6000806040838503121561183b57600080fd5b82356001600160401b038082111561185257600080fd5b61185e86838701611810565b9350602085013591508082111561187457600080fd5b5061188185828601611810565b9150509250929050565b60006001600160401b038211156118a4576118a461162d565b50601f01601f191660200190565b600082601f8301126118c357600080fd5b81356118d66118d18261188b565b6116f4565b8181528460208386010111156118eb57600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a081121561191f57600080fd5b604081121561192d57600080fd5b5084935060408501356001600160401b038082111561194b57600080fd5b61195788838901611810565b9450606087013591508082111561196d57600080fd5b90860190610180828903121561198257600080fd5b9092506080860135908082111561199857600080fd5b506119a5878288016118b2565b91505092959194509250565b803561ffff8116811461174457600080fd5b600080606083850312156119d657600080fd5b6119df836119b1565b91506118078460208501611758565b600060208284031215611a0057600080fd5b81356001600160401b03811115611a1657600080fd5b6104a884828501611810565b600081518084526020808501945080840160005b83811015611a5857815163ffffffff1687529582019590820190600101611a36565b509495945050505050565b600081518084526020808501945080840160005b83811015611a5857611a9487835180518252602090810151910152565b6040969096019590820190600101611a77565b8060005b60028110156102c0578151845260209384019390910190600101611aab565b611ad5828251611aa7565b6020810151611ae76040840182611aa7565b505050565b600081518084526020808501808196508360051b8101915082860160005b85811015611b34578284038952611b22848351611a22565b98850198935090840190600101611b0a565b5091979650505050505050565b60006101808251818552611b5782860182611a22565b91505060208301518482036020860152611b718282611a63565b91505060408301518482036040860152611b8b8282611a63565b9150506060830151611ba06060860182611aca565b506080830151805160e08601526020015161010085015260a0830151848203610120860152611bcf8282611a22565b91505060c0830151848203610140860152611bea8282611a22565b91505060e0830151848203610160860152611c058282611aec565b95945050505050565b602081526000611c216020830184611b41565b9392505050565b600060408284031215611c3a57600080fd5b611c42611643565b905081358152602082013561178f81611724565b60006001600160401b03821115611c6f57611c6f61162d565b5060051b60200190565b600060408284031215611c8b57600080fd5b611c93611643565b9050813581526020820135602082015292915050565b600082601f830112611cba57600080fd5b81356020611cca6118d183611c56565b82815260069290921b84018101918181019086841115611ce957600080fd5b8286015b84811015611d0d57611cff8882611c79565b835291830191604001611ced565b509695505050505050565b600082601f830112611d2957600080fd5b611d31611643565b806040840185811115611d4357600080fd5b845b81811015611d5d578035845260209384019301611d45565b509095945050505050565b600060808284031215611d7a57600080fd5b611d82611643565b9050611d8e8383611d18565b815261178f8360408401611d18565b600082601f830112611dae57600080fd5b81356020611dbe6118d183611c56565b82815260059290921b84018101918181019086841115611ddd57600080fd5b8286015b84811015611d0d578035611df481611724565b8352918301918301611de1565b600060608236031215611e1357600080fd5b611e1b611643565b611e253684611c28565b815260408301356001600160401b0380821115611e4157600080fd5b81850191506101208236031215611e5757600080fd5b611e5f61168d565b823582811115611e6e57600080fd5b611e7a36828601611ca9565b825250602083013582811115611e8f57600080fd5b611e9b36828601611ca9565b602083015250611eae3660408501611c79565b6040820152611ec03660808501611d68565b606082015261010083013582811115611ed857600080fd5b611ee436828601611d9d565b608083015250602084015250909392505050565b600060608236031215611f0a57600080fd5b611f1261166b565b82356001600160401b0380821115611f2957600080fd5b818501915060608236031215611f3e57600080fd5b611f4661166b565b823582811115611f5557600080fd5b8301368190036101c0811215611f6a57600080fd5b611f726116af565b611f7b836119b1565b815260208084013586811115611f9057600080fd5b611f9c368287016118b2565b83830152506040610160603f1985011215611fb657600080fd5b611fbe6116af565b9350611fcc36828701611c79565b8452611fdb3660808701611d68565b82850152611fed366101008701611d68565b8185015261018085013561200081611724565b8060608601525083818401526101a085013560608401528286528188013594508685111561202d57600080fd5b61203936868a016118b2565b828701528088013594508685111561205057600080fd5b61205c36868a01611d9d565b8187015285895261206e828c01611739565b828a0152808b013597508688111561208557600080fd5b61209136898d016118b2565b90890152509598975050505050505050565b6000604082840312156120b557600080fd5b611c218383611c28565b600082601f8301126120d057600080fd5b813560206120e06118d183611c56565b82815260059290921b840181019181810190868411156120ff57600080fd5b8286015b84811015611d0d5780356001600160401b038111156121225760008081fd5b6121308986838b0101611d9d565b845250918301918301612103565b6000610180823603121561215157600080fd5b6121596116d1565b82356001600160401b038082111561217057600080fd5b61217c36838701611d9d565b8352602085013591508082111561219257600080fd5b61219e36838701611ca9565b602084015260408501359150808211156121b757600080fd5b6121c336838701611ca9565b60408401526121d53660608701611d68565b60608401526121e73660e08701611c79565b608084015261012085013591508082111561220157600080fd5b61220d36838701611d9d565b60a084015261014085013591508082111561222757600080fd5b61223336838701611d9d565b60c084015261016085013591508082111561224d57600080fd5b5061225a368286016120bf565b60e08301525092915050565b634e487b7160e01b600052602160045260246000fd5b60006060828403121561228e57600080fd5b604051606081018181106001600160401b03821117156122b0576122b061162d565b60405282516122be81611724565b815260208301516122ce81611724565b602082015260408301516122e181611749565b60408201529392505050565b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff84168082101561231d5761231d6122ed565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b60008261234b5761234b612326565b500490565b600082821015612362576123626122ed565b500390565b6000816000190483118215151615612381576123816122ed565b500290565b600063ffffffff808316818516818304811182151516156123a9576123a96122ed565b02949350505050565b6000806000606084860312156123c757600080fd5b8351925060208401519150604084015190509250925092565b600080604083850312156123f357600080fd5b82516123fe81611749565b602084015190925061240f81611749565b809150509250929050565b6000806040838503121561242d57600080fd5b505080516020909101519092909150565b60006020828403121561245057600080fd5b8151611c2181611724565b60005b8381101561247657818101518382015260200161245e565b838111156102c05750506000910152565b60006020828403121561249957600080fd5b81516001600160401b038111156124af57600080fd5b8201601f810184136124c057600080fd5b80516124ce6118d18261188b565b8181528560208385010111156124e357600080fd5b611c0582602083016020860161245b565b634e487b7160e01b600052603260045260246000fd5b600060001982141561251e5761251e6122ed565b5060010190565b6000835161253781846020880161245b565b60f89390931b6001600160f81b0319169190920190815260010192915050565b6000815180845261256f81602086016020860161245b565b601f01601f19169290920160200192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526125b16080840186612557565b838103606085015284518082528286019183019060005b818110156125e4578351835292840192918401916001016125c8565b50909998505050505050505050565b600082601f83011261260457600080fd5b815160206126146118d183611c56565b82815260059290921b8401810191818101908684111561263357600080fd5b8286015b84811015611d0d57805161264a81611724565b8352918301918301612637565b600082601f83011261266857600080fd5b815160206126786118d183611c56565b82815260059290921b8401810191818101908684111561269757600080fd5b8286015b84811015611d0d5780516001600160401b038111156126ba5760008081fd5b6126c88986838b01016125f3565b84525091830191830161269b565b6000602082840312156126e857600080fd5b81516001600160401b03808211156126ff57600080fd5b908301906080828603121561271357600080fd5b61271b6116af565b82518281111561272a57600080fd5b612736878286016125f3565b82525060208301518281111561274b57600080fd5b612757878286016125f3565b60208301525060408301518281111561276f57600080fd5b61277b878286016125f3565b60408301525060608301518281111561279357600080fd5b61279f87828601612657565b60608301525095945050505050565b6000602082840312156127c057600080fd5b81516001600160a01b0381168114611c2157600080fd5b8481526080602082015260006127f06080830186612557565b63ffffffff85166040840152828103606084015261280e8185611b41565b979650505050505050565b600082601f83011261282a57600080fd5b8151602061283a6118d183611c56565b82815260059290921b8401810191818101908684111561285957600080fd5b8286015b84811015611d0d5780516001600160601b038116811461287d5760008081fd5b835291830191830161285d565b6000806040838503121561289d57600080fd5b82516001600160401b03808211156128b457600080fd5b90840190604082870312156128c857600080fd5b6128d0611643565b8251828111156128df57600080fd5b6128eb88828601612819565b82525060208301518281111561290057600080fd5b61290c88828601612819565b602083015250809450505050602083015190509250929050565b60006001600160601b03808316818516818304811182151516156123a9576123a96122ed565b838152600060206060818401526129666060840186612557565b838103604085015284518082528286019183019060005b8181101561299f57835163ffffffff168352928401929184019160010161297d565b509098975050505050505050565b60006101a061ffff861683528060208401526129cb81840186612557565b84518051604086015260200151606085015291506129e69050565b60208301516129f86080840182611aca565b506040830151612a0c610100840182611aca565b5063ffffffff606084015116610180830152949350505050565b600082612a3557612a35612326565b500690565b60008219821115612a4d57612a4d6122ed565b50019056fea2646970667358221220dc96048de5e1169a6056cbbd446f398e113e1ddeb594418ca04b29a754aa286c64736f6c634300080c0033",
}

// ContractEigenDACertVerifierV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.ABI instead.
var ContractEigenDACertVerifierV2ABI = ContractEigenDACertVerifierV2MetaData.ABI

// ContractEigenDACertVerifierV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.Bin instead.
var ContractEigenDACertVerifierV2Bin = ContractEigenDACertVerifierV2MetaData.Bin

// DeployContractEigenDACertVerifierV2 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV2 to it.
func DeployContractEigenDACertVerifierV2(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistryV2 common.Address, _eigenDASignatureVerifierV2 common.Address, _eigenDARelayRegistryV2 common.Address, _operatorStateRetrieverV2 common.Address, _registryCoordinatorV2 common.Address, _securityThresholdsV2 SecurityThresholds) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV2, error) {
	parsed, err := ContractEigenDACertVerifierV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV2Bin), backend, _eigenDAThresholdRegistryV2, _eigenDASignatureVerifierV2, _eigenDARelayRegistryV2, _operatorStateRetrieverV2, _registryCoordinatorV2, _securityThresholdsV2)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDACertVerifierV2{ContractEigenDACertVerifierV2Caller: ContractEigenDACertVerifierV2Caller{contract: contract}, ContractEigenDACertVerifierV2Transactor: ContractEigenDACertVerifierV2Transactor{contract: contract}, ContractEigenDACertVerifierV2Filterer: ContractEigenDACertVerifierV2Filterer{contract: contract}}, nil
}

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

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) CheckDACertV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "checkDACertV2", batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) CheckDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.CheckDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// CheckDACertV2 is a free data retrieval call binding the contract method 0x8ec28be6.
//
// Solidity: function checkDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns(bool success)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) CheckDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) (bool, error) {
	return _ContractEigenDACertVerifierV2.Contract.CheckDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) EigenDARelayRegistryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "eigenDARelayRegistryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
}

// EigenDARelayRegistryV2 is a free data retrieval call binding the contract method 0x82c216e7.
//
// Solidity: function eigenDARelayRegistryV2() view returns(address)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) EigenDARelayRegistryV2() (common.Address, error) {
	return _ContractEigenDACertVerifierV2.Contract.EigenDARelayRegistryV2(&_ContractEigenDACertVerifierV2.CallOpts)
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
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) GetNonSignerStakesAndSignature(opts *bind.CallOpts, signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "getNonSignerStakesAndSignature", signedBatch)

	if err != nil {
		return *new(NonSignerStakesAndSignature), err
	}

	out0 := *abi.ConvertType(out[0], new(NonSignerStakesAndSignature)).(*NonSignerStakesAndSignature)

	return out0, err

}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
	return _ContractEigenDACertVerifierV2.Contract.GetNonSignerStakesAndSignature(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch)
}

// GetNonSignerStakesAndSignature is a free data retrieval call binding the contract method 0xf25de3f8.
//
// Solidity: function getNonSignerStakesAndSignature(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch) view returns((uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) GetNonSignerStakesAndSignature(signedBatch SignedBatch) (NonSignerStakesAndSignature, error) {
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

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) SecurityThresholdsAdversary(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "securityThresholdsAdversary")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) SecurityThresholdsAdversary() (uint8, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsAdversary(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsAdversary is a free data retrieval call binding the contract method 0xdea610a9.
//
// Solidity: function securityThresholdsAdversary() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) SecurityThresholdsAdversary() (uint8, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsAdversary(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) SecurityThresholdsConfirmation(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "securityThresholdsConfirmation")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) SecurityThresholdsConfirmation() (uint8, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsConfirmation(&_ContractEigenDACertVerifierV2.CallOpts)
}

// SecurityThresholdsConfirmation is a free data retrieval call binding the contract method 0xfd174484.
//
// Solidity: function securityThresholdsConfirmation() view returns(uint8)
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) SecurityThresholdsConfirmation() (uint8, error) {
	return _ContractEigenDACertVerifierV2.Contract.SecurityThresholdsConfirmation(&_ContractEigenDACertVerifierV2.CallOpts)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertSecurityParams(opts *bind.CallOpts, blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertSecurityParams", blobParams, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifierV2.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams is a free data retrieval call binding the contract method 0x143eb4d9.
//
// Solidity: function verifyDACertSecurityParams((uint32,uint32,uint8) blobParams, (uint8,uint8) securityThresholds) pure returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertSecurityParams(blobParams VersionedBlobParams, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertSecurityParams(&_ContractEigenDACertVerifierV2.CallOpts, blobParams, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertSecurityParams0(opts *bind.CallOpts, version uint16, securityThresholds SecurityThresholds) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV2.contract.Call(opts, &out, "verifyDACertSecurityParams0", version, securityThresholds)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifierV2.CallOpts, version, securityThresholds)
}

// VerifyDACertSecurityParams0 is a free data retrieval call binding the contract method 0xccb7cd0d.
//
// Solidity: function verifyDACertSecurityParams(uint16 version, (uint8,uint8) securityThresholds) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertSecurityParams0(version uint16, securityThresholds SecurityThresholds) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertSecurityParams0(&_ContractEigenDACertVerifierV2.CallOpts, version, securityThresholds)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2(opts *bind.CallOpts, batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
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
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2 is a free data retrieval call binding the contract method 0x813c2eb0.
//
// Solidity: function verifyDACertV2((bytes32,uint32) batchHeader, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes signedQuorumNumbers) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2(batchHeader BatchHeaderV2, blobInclusionInfo BlobInclusionInfo, nonSignerStakesAndSignature NonSignerStakesAndSignature, signedQuorumNumbers []byte) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2(&_ContractEigenDACertVerifierV2.CallOpts, batchHeader, blobInclusionInfo, nonSignerStakesAndSignature, signedQuorumNumbers)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Caller) VerifyDACertV2FromSignedBatch(opts *bind.CallOpts, signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
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
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2Session) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch, blobInclusionInfo)
}

// VerifyDACertV2FromSignedBatch is a free data retrieval call binding the contract method 0x421c0222.
//
// Solidity: function verifyDACertV2FromSignedBatch(((bytes32,uint32),((uint256,uint256)[],(uint256,uint256)[],(uint256,uint256),(uint256[2],uint256[2]),uint32[])) signedBatch, (((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes) blobInclusionInfo) view returns()
func (_ContractEigenDACertVerifierV2 *ContractEigenDACertVerifierV2CallerSession) VerifyDACertV2FromSignedBatch(signedBatch SignedBatch, blobInclusionInfo BlobInclusionInfo) error {
	return _ContractEigenDACertVerifierV2.Contract.VerifyDACertV2FromSignedBatch(&_ContractEigenDACertVerifierV2.CallOpts, signedBatch, blobInclusionInfo)
}
