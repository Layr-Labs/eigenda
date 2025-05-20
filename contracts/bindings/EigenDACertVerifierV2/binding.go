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
	Bin: "0x6101006040523480156200001257600080fd5b5060405162002a2838038062002a288339810160408190526200003591620002a8565b816020015160ff16826000015160ff161162000064576040516308a6997560e01b815260040160405180910390fd5b6001600160a01b0380871660805285811660a05284811660c052831660e05281516000805460208086015160ff9081166101000261ffff199093169416939093171790558151620000bc9160019190840190620000c9565b50505050505050620003c1565b828054620000d79062000384565b90600052602060002090601f016020900481019282620000fb576000855562000146565b82601f106200011657805160ff191683800117855562000146565b8280016001018555821562000146579182015b828111156200014657825182559160200191906001019062000129565b506200015492915062000158565b5090565b5b8082111562000154576000815560010162000159565b6001600160a01b03811681146200018557600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b0381118282101715620001c357620001c362000188565b60405290565b604051601f8201601f191681016001600160401b0381118282101715620001f457620001f462000188565b604052919050565b805160ff811681146200020e57600080fd5b919050565b600082601f8301126200022557600080fd5b81516001600160401b0381111562000241576200024162000188565b602062000257601f8301601f19168201620001c9565b82815285828487010111156200026c57600080fd5b60005b838110156200028c5785810183015182820184015282016200026f565b838111156200029e5760008385840101525b5095945050505050565b60008060008060008086880360e0811215620002c357600080fd5b8751620002d0816200016f565b6020890151909750620002e3816200016f565b6040890151909650620002f6816200016f565b606089015190955062000309816200016f565b93506040607f19820112156200031e57600080fd5b50620003296200019e565b6200033760808901620001fc565b81526200034760a08901620001fc565b602082015260c08801519092506001600160401b038111156200036957600080fd5b6200037789828a0162000213565b9150509295509295509295565b600181811c908216806200039957607f821691505b60208210811415620003bb57634e487b7160e01b600052602260045260246000fd5b50919050565b60805160a05160c05160e0516125e96200043f600039600081816101720152818161032801526104b201526000818161014b01528181610307015261049101526000818160a801528181610235015281816102e6015261039801526000818160ec01528181610214015281816102c5015261037701526125e96000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c80635fafa482116100665780635fafa4821461016d578063813c2eb014610194578063b74d7871146101a7578063ed0450ae146101bc578063f25de3f8146101ec57600080fd5b8063154b9e86146100a357806317f3578e146100e7578063415ef6141461010e578063421c0222146101315780635df1f61814610146575b600080fd5b6100ca7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020015b60405180910390f35b6100ca7f000000000000000000000000000000000000000000000000000000000000000081565b61012161011c3660046114d0565b61020c565b60405190151581526020016100de565b61014461013f366004611579565b6102c0565b005b6100ca7f000000000000000000000000000000000000000000000000000000000000000081565b6100ca7f000000000000000000000000000000000000000000000000000000000000000081565b6101446101a23660046114d0565b610372565b6101af6103f4565b6040516100de9190611634565b6000546101d29060ff8082169161010090041682565b6040805160ff9384168152929091166020830152016100de565b6101ff6101fa36600461164e565b610482565b6040516100de919061186e565b60008061028b7f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000610263368a90038a018a6118df565b61026c89611a47565b61027589611cde565b61027d6104e7565b61028561051b565b8a6105ad565b50905060008160048111156102a2576102a2611e06565b14156102b25760019150506102b8565b60009150505b949350505050565b61036e7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061035087611e1c565b61035987611a47565b6103616104e7565b61036961051b565b610717565b5050565b6103ee7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006103c6368890038801886118df565b6103cf87611a47565b6103d887611cde565b6103e06104e7565b6103e861051b565b88610749565b50505050565b6001805461040190611f13565b80601f016020809104026020016040519081016040528092919081815260200182805461042d90611f13565b801561047a5780601f1061044f5761010080835404028352916020019161047a565b820191906000526020600020905b81548152906001019060200180831161045d57829003601f168201915b505050505081565b61048a6112a9565b60006104df7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006104da86611e1c565b61076a565b509392505050565b6040805180820182526000808252602091820181905282518084019093525460ff8082168452610100909104169082015290565b60606001805461052a90611f13565b80601f016020809104026020016040519081016040528092919081815260200182805461055690611f13565b80156105a35780601f10610578576101008083540402835291602001916105a3565b820191906000526020600020905b81548152906001019060200180831161058657829003601f168201915b5050505050905090565b600060606105bb888861097f565b909250905060008260048111156105d4576105d4611e06565b146105de5761070a565b86515151604051632ecfe72b60e01b815261ffff9091166004820152610659906001600160a01b038c1690632ecfe72b90602401606060405180830381865afa15801561062f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106539190611f48565b86610a58565b9092509050600082600481111561067257610672611e06565b1461067c5761070a565b60006106988a61068b8b610b35565b868c602001518b8b610b7d565b9194509250905060008360048111156106b3576106b3611e06565b146106be575061070a565b875151602001516000906106d29083610cde565b9195509350905060008460048111156106ed576106ed611e06565b146106f957505061070a565b6107038682610d44565b9350935050505b9850989650505050505050565b60008061072588888861076a565b9150915061073d8a8a88600001518886898988610749565b50505050505050505050565b60008061075c8a8a8a8a8a8a8a8a6105ad565b9150915061073d8282610da8565b6107726112a9565b60606000836020015160000151516001600160401b038111156107975761079761136a565b6040519080825280602002602001820160405280156107c0578160200160208202803683370190505b50905060005b6020850151515181101561083d5761081085602001516000015182815181106107f1576107f1611fbf565b6020026020010151805160009081526020918201519091526040902090565b82828151811061082257610822611fbf565b602090810291909101015261083681611feb565b90506107c6565b5060005b846020015160800151518110156108a85782856020015160800151828151811061086d5761086d611fbf565b6020026020010151604051602001610886929190612006565b6040516020818303038152906040529250806108a190611feb565b9050610841565b508351602001516040516313dce7dd60e21b81526000916001600160a01b03891691634f739f74916108e3918a919088908890600401612038565b600060405180830381865afa158015610900573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610928919081019061218b565b8051855260208087018051518288015280518201516040808901919091528151606090810151818a0152915181015160808901529183015160a08801529082015160c0870152015160e08501525050935093915050565b6000606060006109928460000151610f87565b90506000816040516020016109a991815260200190565b60405160208183030381529060405280519060200120905060008660000151905060006109e6876040015183858a6020015163ffffffff16610faf565b90508015610a0d576000604051806020016040528060008152509550955050505050610a51565b6020808801516040805163ffffffff909216928201929092529081018490526060810183905260019060800160405160208183030381529060405295509550505050505b9250929050565b60006060600083602001518460000151610a729190612263565b60ff1690506000856020015163ffffffff16866040015160ff1683620f4240610a9b919061229c565b610aa5919061229c565b610ab1906127106122b0565b610abb91906122c7565b8651909150600090610acf906127106122e6565b63ffffffff169050808210610afc5760006040518060200160405280600081525094509450505050610a51565b60408051602081018590529081018390526060810182905260029060800160405160208183030381529060405294509450505050610a51565b600081604051602001610b6091908151815260209182015163ffffffff169181019190915260400190565b604051602081830303815290604052805190602001209050919050565b60006060600080896001600160a01b0316636efb46368a8a8a8a6040518563ffffffff1660e01b8152600401610bb69493929190612312565b600060405180830381865afa158015610bd3573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052610bfb91908101906123c5565b5090506000915060005b8851811015610cbc57856000015160ff1682602001518281518110610c2c57610c2c611fbf565b6020026020010151610c3e9190612461565b6001600160601b0316606483600001518381518110610c5f57610c5f611fbf565b60200260200101516001600160601b0316610c7a91906122c7565b10610caa57610ca7838a8381518110610c9557610c95611fbf565b0160200151600160f89190911c1b1790565b92505b80610cb481611feb565b915050610c05565b5050604080516020810190915260008082529350915096509650969350505050565b600060606000610ced85610fc7565b9050838116811415610d12576040805160208101909152600080825293509150610d3d565b6040805160208101929092528181018590528051808303820181526060909201905260039250905060005b9250925092565b600060606000610d5385610fc7565b9050838116811415610d78575050604080516020810190915260008082529150610a51565b60408051602081018390529081018590526004906060016040516020818303038152906040529250925050610a51565b6000826004811115610dbc57610dbc611e06565b1415610dc6575050565b6001826004811115610dda57610dda611e06565b1415610e3057600080600083806020019051810190610df99190612487565b60405163d54d727760e01b815260048101849052602481018390526044810182905292955090935091506064015b60405180910390fd5b6002826004811115610e4457610e44611e06565b1415610e9857600080600083806020019051810190610e639190612487565b6040516001626dc9ad60e11b031981526004810184905260248101839052604481018290529295509093509150606401610e27565b6003826004811115610eac57610eac611e06565b1415610ef15760008082806020019051810190610ec991906124b5565b604051634a47030360e11b815260048101839052602481018290529193509150604401610e27565b6004826004811115610f0557610f05611e06565b1415610f4a5760008082806020019051810190610f2291906124b5565b60405163114b085b60e21b815260048101839052602481018290529193509150604401610e27565b60405162461bcd60e51b8152602060048201526012602482015271556e6b6e6f776e206572726f7220636f646560701b6044820152606401610e27565b6000610f968260000151611154565b6020808401516040808601519051610b609493016124d9565b600083610fbd8685856111a6565b1495945050505050565b6000610100825111156110505760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610e27565b815161105e57506000919050565b6000808360008151811061107457611074611fbf565b0160200151600160f89190911c81901b92505b845181101561114b578481815181106110a2576110a2611fbf565b0160200151600160f89190911c1b91508282116111375760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610e27565b9181179161114481611feb565b9050611087565b50909392505050565b60008160000151826020015183604001516040516020016111779392919061250e565b60408051601f198184030181528282528051602091820120606080870151928501919091529183015201610b60565b6000602084516111b69190612587565b1561123d5760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610e27565b8260205b855181116112a057611254600285612587565b6112755781600052808601516020526040600020915060028404935061128e565b8086015160005281602052604060002091506002840493505b61129960208261259b565b9050611241565b50949350505050565b6040518061010001604052806060815260200160608152602001606081526020016112d261130f565b81526020016112f4604051806040016040528060008152602001600081525090565b81526020016060815260200160608152602001606081525090565b6040518060400160405280611322611334565b815260200161132f611334565b905290565b60405180604001604052806002906020820280368337509192915050565b60006060828403121561136457600080fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156113a2576113a261136a565b60405290565b604051606081016001600160401b03811182821017156113a2576113a261136a565b604051608081016001600160401b03811182821017156113a2576113a261136a565b60405161010081016001600160401b03811182821017156113a2576113a261136a565b60405160a081016001600160401b03811182821017156113a2576113a261136a565b604051601f8201601f191681016001600160401b03811182821017156114595761145961136a565b604052919050565b600082601f83011261147257600080fd5b81356001600160401b0381111561148b5761148b61136a565b61149e601f8201601f1916602001611431565b8181528460208386010111156114b357600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008084860360a08112156114e757600080fd5b60408112156114f557600080fd5b5084935060408501356001600160401b038082111561151357600080fd5b61151f88838901611352565b9450606087013591508082111561153557600080fd5b90860190610180828903121561154a57600080fd5b9092506080860135908082111561156057600080fd5b5061156d87828801611461565b91505092959194509250565b6000806040838503121561158c57600080fd5b82356001600160401b03808211156115a357600080fd5b6115af86838701611352565b935060208501359150808211156115c557600080fd5b506115d285828601611352565b9150509250929050565b60005b838110156115f75781810151838201526020016115df565b838111156103ee5750506000910152565b600081518084526116208160208601602086016115dc565b601f01601f19169290920160200192915050565b6020815260006116476020830184611608565b9392505050565b60006020828403121561166057600080fd5b81356001600160401b0381111561167657600080fd5b6102b884828501611352565b600081518084526020808501945080840160005b838110156116b857815163ffffffff1687529582019590820190600101611696565b509495945050505050565b600081518084526020808501945080840160005b838110156116b8576116f487835180518252602090810151910152565b60409690960195908201906001016116d7565b8060005b60028110156103ee57815184526020938401939091019060010161170b565b611735828251611707565b60208101516117476040840182611707565b505050565b600081518084526020808501808196508360051b8101915082860160005b85811015611794578284038952611782848351611682565b9885019893509084019060010161176a565b5091979650505050505050565b600061018082518185526117b782860182611682565b915050602083015184820360208601526117d182826116c3565b915050604083015184820360408601526117eb82826116c3565b9150506060830151611800606086018261172a565b506080830151805160e08601526020015161010085015260a083015184820361012086015261182f8282611682565b91505060c083015184820361014086015261184a8282611682565b91505060e0830151848203610160860152611865828261174c565b95945050505050565b60208152600061164760208301846117a1565b63ffffffff8116811461189357600080fd5b50565b80356118a181611881565b919050565b6000604082840312156118b857600080fd5b6118c0611380565b90508135815260208201356118d481611881565b602082015292915050565b6000604082840312156118f157600080fd5b61164783836118a6565b60006040828403121561190d57600080fd5b611915611380565b9050813581526020820135602082015292915050565b600082601f83011261193c57600080fd5b611944611380565b80604084018581111561195657600080fd5b845b81811015611970578035845260209384019301611958565b509095945050505050565b60006080828403121561198d57600080fd5b611995611380565b90506119a1838361192b565b81526118d4836040840161192b565b60006001600160401b038211156119c9576119c961136a565b5060051b60200190565b600082601f8301126119e457600080fd5b813560206119f96119f4836119b0565b611431565b82815260059290921b84018101918181019086841115611a1857600080fd5b8286015b84811015611a3c578035611a2f81611881565b8352918301918301611a1c565b509695505050505050565b600060608236031215611a5957600080fd5b611a616113a8565b82356001600160401b0380821115611a7857600080fd5b818501915060608236031215611a8d57600080fd5b611a956113a8565b823582811115611aa457600080fd5b8301368190036101c0811215611ab957600080fd5b611ac16113ca565b823561ffff81168114611ad357600080fd5b815260208381013586811115611ae857600080fd5b611af436828701611461565b83830152506040610160603f1985011215611b0e57600080fd5b611b166113ca565b9350611b24368287016118fb565b8452611b33366080870161197b565b82850152611b4536610100870161197b565b81850152610180850135611b5881611881565b8060608601525083818401526101a0850135606084015282865281880135945086851115611b8557600080fd5b611b9136868a01611461565b8287015280880135945086851115611ba857600080fd5b611bb436868a016119d3565b81870152858952611bc6828c01611896565b828a0152808b0135975086881115611bdd57600080fd5b611be936898d01611461565b90890152509598975050505050505050565b600082601f830112611c0c57600080fd5b81356020611c1c6119f4836119b0565b82815260069290921b84018101918181019086841115611c3b57600080fd5b8286015b84811015611a3c57611c5188826118fb565b835291830191604001611c3f565b600082601f830112611c7057600080fd5b81356020611c806119f4836119b0565b82815260059290921b84018101918181019086841115611c9f57600080fd5b8286015b84811015611a3c5780356001600160401b03811115611cc25760008081fd5b611cd08986838b01016119d3565b845250918301918301611ca3565b60006101808236031215611cf157600080fd5b611cf96113ec565b82356001600160401b0380821115611d1057600080fd5b611d1c368387016119d3565b83526020850135915080821115611d3257600080fd5b611d3e36838701611bfb565b60208401526040850135915080821115611d5757600080fd5b611d6336838701611bfb565b6040840152611d75366060870161197b565b6060840152611d873660e087016118fb565b6080840152610120850135915080821115611da157600080fd5b611dad368387016119d3565b60a0840152610140850135915080821115611dc757600080fd5b611dd3368387016119d3565b60c0840152610160850135915080821115611ded57600080fd5b50611dfa36828601611c5f565b60e08301525092915050565b634e487b7160e01b600052602160045260246000fd5b600060608236031215611e2e57600080fd5b611e36611380565b611e4036846118a6565b815260408301356001600160401b0380821115611e5c57600080fd5b81850191506101208236031215611e7257600080fd5b611e7a61140f565b823582811115611e8957600080fd5b611e9536828601611bfb565b825250602083013582811115611eaa57600080fd5b611eb636828601611bfb565b602083015250611ec936604085016118fb565b6040820152611edb366080850161197b565b606082015261010083013582811115611ef357600080fd5b611eff368286016119d3565b608083015250602084015250909392505050565b600181811c90821680611f2757607f821691505b6020821081141561136457634e487b7160e01b600052602260045260246000fd5b600060608284031215611f5a57600080fd5b604051606081018181106001600160401b0382111715611f7c57611f7c61136a565b6040528251611f8a81611881565b81526020830151611f9a81611881565b6020820152604083015160ff81168114611fb357600080fd5b60408201529392505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415611fff57611fff611fd5565b5060010190565b600083516120188184602088016115dc565b60f89390931b6001600160f81b0319169190920190815260010192915050565b60018060a01b03851681526000602063ffffffff861681840152608060408401526120666080840186611608565b838103606085015284518082528286019183019060005b818110156120995783518352928401929184019160010161207d565b50909998505050505050505050565b600082601f8301126120b957600080fd5b815160206120c96119f4836119b0565b82815260059290921b840181019181810190868411156120e857600080fd5b8286015b84811015611a3c5780516120ff81611881565b83529183019183016120ec565b600082601f83011261211d57600080fd5b8151602061212d6119f4836119b0565b82815260059290921b8401810191818101908684111561214c57600080fd5b8286015b84811015611a3c5780516001600160401b0381111561216f5760008081fd5b61217d8986838b01016120a8565b845250918301918301612150565b60006020828403121561219d57600080fd5b81516001600160401b03808211156121b457600080fd5b90830190608082860312156121c857600080fd5b6121d06113ca565b8251828111156121df57600080fd5b6121eb878286016120a8565b82525060208301518281111561220057600080fd5b61220c878286016120a8565b60208301525060408301518281111561222457600080fd5b612230878286016120a8565b60408301525060608301518281111561224857600080fd5b6122548782860161210c565b60608301525095945050505050565b600060ff821660ff84168082101561227d5761227d611fd5565b90039392505050565b634e487b7160e01b600052601260045260246000fd5b6000826122ab576122ab612286565b500490565b6000828210156122c2576122c2611fd5565b500390565b60008160001904831182151516156122e1576122e1611fd5565b500290565b600063ffffffff8083168185168183048111821515161561230957612309611fd5565b02949350505050565b84815260806020820152600061232b6080830186611608565b63ffffffff85166040840152828103606084015261234981856117a1565b979650505050505050565b600082601f83011261236557600080fd5b815160206123756119f4836119b0565b82815260059290921b8401810191818101908684111561239457600080fd5b8286015b84811015611a3c5780516001600160601b03811681146123b85760008081fd5b8352918301918301612398565b600080604083850312156123d857600080fd5b82516001600160401b03808211156123ef57600080fd5b908401906040828703121561240357600080fd5b61240b611380565b82518281111561241a57600080fd5b61242688828601612354565b82525060208301518281111561243b57600080fd5b61244788828601612354565b602083015250809450505050602083015190509250929050565b60006001600160601b038083168185168183048111821515161561230957612309611fd5565b60008060006060848603121561249c57600080fd5b8351925060208401519150604084015190509250925092565b600080604083850312156124c857600080fd5b505080516020909101519092909150565b8381526060602082015260006124f26060830185611608565b82810360408401526125048185611682565b9695505050505050565b60006101a061ffff8616835280602084015261252c81840186611608565b84518051604086015260200151606085015291506125479050565b6020830151612559608084018261172a565b50604083015161256d61010084018261172a565b5063ffffffff606084015116610180830152949350505050565b60008261259657612596612286565b500690565b600082198211156125ae576125ae611fd5565b50019056fea26469706673582212202932d34f4fa537a33d6ac86df824291714beb08b6b3d9608641c1cb5c3fa361964736f6c634300080c0033",
}

// ContractEigenDACertVerifierV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.ABI instead.
var ContractEigenDACertVerifierV2ABI = ContractEigenDACertVerifierV2MetaData.ABI

// ContractEigenDACertVerifierV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDACertVerifierV2MetaData.Bin instead.
var ContractEigenDACertVerifierV2Bin = ContractEigenDACertVerifierV2MetaData.Bin

// DeployContractEigenDACertVerifierV2 deploys a new Ethereum contract, binding an instance of ContractEigenDACertVerifierV2 to it.
func DeployContractEigenDACertVerifierV2(auth *bind.TransactOpts, backend bind.ContractBackend, _eigenDAThresholdRegistryV2 common.Address, _eigenDASignatureVerifierV2 common.Address, _operatorStateRetrieverV2 common.Address, _registryCoordinatorV2 common.Address, _securityThresholdsV2 EigenDATypesV1SecurityThresholds, _quorumNumbersRequiredV2 []byte) (common.Address, *types.Transaction, *ContractEigenDACertVerifierV2, error) {
	parsed, err := ContractEigenDACertVerifierV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDACertVerifierV2Bin), backend, _eigenDAThresholdRegistryV2, _eigenDASignatureVerifierV2, _operatorStateRetrieverV2, _registryCoordinatorV2, _securityThresholdsV2, _quorumNumbersRequiredV2)
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
