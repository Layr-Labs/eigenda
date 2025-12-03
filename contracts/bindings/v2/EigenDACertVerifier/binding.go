// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifier

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
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

// EigenDACertTypesEigenDACertV4 is an auto generated low-level Go binding around an user-defined struct.
type EigenDACertTypesEigenDACertV4 struct {
	BatchHeader                 EigenDATypesV2BatchHeaderV2
	BlobInclusionInfo           EigenDATypesV2BlobInclusionInfo
	NonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
	OffchainDerivationVersion   uint16
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

// ContractEigenDACertVerifierMetaData contains all meta data concerning the ContractEigenDACertVerifier contract.
var ContractEigenDACertVerifierMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initEigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"initEigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"initSecurityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"initQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initOffchainDerivationVersion\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_decodeCert\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"cert\",\"type\":\"tuple\",\"internalType\":\"structEigenDACertTypes.EigenDACertV4\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainDerivationVersion\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkDACertReverts\",\"inputs\":[{\"name\":\"daCert\",\"type\":\"tuple\",\"internalType\":\"structEigenDACertTypes.EigenDACertV4\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"offchainDerivationVersion\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offchainDerivationVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"semver\",\"inputs\":[],\"outputs\":[{\"name\":\"major\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minor\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"patch\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"error\",\"name\":\"BlobQuorumsNotSubset\",\"inputs\":[{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmedQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBlobVersion\",\"inputs\":[{\"name\":\"blobVersion\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"nextBlobVersion\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidInclusionProof\",\"inputs\":[{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"rootHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidQuorumNumbersRequired\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequiredQuorumsNotSubset\",\"inputs\":[{\"name\":\"requiredQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blobQuorumsBitmap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SecurityAssumptionsNotMet\",\"inputs\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
	ID:  "ContractEigenDACertVerifier",
}

// ContractEigenDACertVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifier struct {
	abi abi.ABI
}

// NewContractEigenDACertVerifier creates a new instance of ContractEigenDACertVerifier.
func NewContractEigenDACertVerifier() *ContractEigenDACertVerifier {
	parsed, err := ContractEigenDACertVerifierMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ContractEigenDACertVerifier{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ContractEigenDACertVerifier) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address initEigenDAThresholdRegistry, address initEigenDASignatureVerifier, (uint8,uint8) initSecurityThresholds, bytes initQuorumNumbersRequired, uint16 initOffchainDerivationVersion) returns()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackConstructor(initEigenDAThresholdRegistry common.Address, initEigenDASignatureVerifier common.Address, initSecurityThresholds EigenDATypesV1SecurityThresholds, initQuorumNumbersRequired []byte, initOffchainDerivationVersion uint16) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("", initEigenDAThresholdRegistry, initEigenDASignatureVerifier, initSecurityThresholds, initQuorumNumbersRequired, initOffchainDerivationVersion)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDecodeCert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x693194fa.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function _decodeCert(bytes data) pure returns(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes,uint16) cert)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackDecodeCert(data []byte) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("_decodeCert", data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecodeCert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x693194fa.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function _decodeCert(bytes data) pure returns(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes,uint16) cert)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackDecodeCert(data []byte) ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("_decodeCert", data)
}

// UnpackDecodeCert is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x693194fa.
//
// Solidity: function _decodeCert(bytes data) pure returns(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes,uint16) cert)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackDecodeCert(data []byte) (EigenDACertTypesEigenDACertV4, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("_decodeCert", data)
	if err != nil {
		return *new(EigenDACertTypesEigenDACertV4), err
	}
	out0 := *abi.ConvertType(out[0], new(EigenDACertTypesEigenDACertV4)).(*EigenDACertTypesEigenDACertV4)
	return out0, nil
}

// PackCertVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ead0b96.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function certVersion() pure returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackCertVersion() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("certVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCertVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2ead0b96.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function certVersion() pure returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackCertVersion() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("certVersion")
}

// UnpackCertVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackCertVersion(data []byte) (uint8, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("certVersion", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackCheckDACert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9077193b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackCheckDACert(abiEncodedCert []byte) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("checkDACert", abiEncodedCert)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckDACert is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9077193b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackCheckDACert(abiEncodedCert []byte) ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("checkDACert", abiEncodedCert)
}

// UnpackCheckDACert is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackCheckDACert(data []byte) (uint8, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("checkDACert", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackCheckDACertReverts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb31cd5e6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function checkDACertReverts(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes,uint16) daCert) view returns()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackCheckDACertReverts(daCert EigenDACertTypesEigenDACertV4) []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("checkDACertReverts", daCert)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCheckDACertReverts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb31cd5e6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function checkDACertReverts(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes,uint16) daCert) view returns()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackCheckDACertReverts(daCert EigenDACertTypesEigenDACertV4) ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("checkDACertReverts", daCert)
}

// PackEigenDASignatureVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefd4532b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackEigenDASignatureVerifier() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("eigenDASignatureVerifier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEigenDASignatureVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefd4532b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackEigenDASignatureVerifier() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("eigenDASignatureVerifier")
}

// UnpackEigenDASignatureVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackEigenDASignatureVerifier(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("eigenDASignatureVerifier", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEigenDAThresholdRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c66814.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackEigenDAThresholdRegistry() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("eigenDAThresholdRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEigenDAThresholdRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c66814.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackEigenDAThresholdRegistry() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("eigenDAThresholdRegistry")
}

// UnpackEigenDAThresholdRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackEigenDAThresholdRegistry(data []byte) (common.Address, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("eigenDAThresholdRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOffchainDerivationVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb326e37f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function offchainDerivationVersion() view returns(uint16)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackOffchainDerivationVersion() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("offchainDerivationVersion")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOffchainDerivationVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb326e37f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function offchainDerivationVersion() view returns(uint16)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackOffchainDerivationVersion() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("offchainDerivationVersion")
}

// UnpackOffchainDerivationVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb326e37f.
//
// Solidity: function offchainDerivationVersion() view returns(uint16)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackOffchainDerivationVersion(data []byte) (uint16, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("offchainDerivationVersion", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackQuorumNumbersRequired is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe15234ff.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackQuorumNumbersRequired() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("quorumNumbersRequired")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQuorumNumbersRequired is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe15234ff.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackQuorumNumbersRequired() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("quorumNumbersRequired")
}

// UnpackQuorumNumbersRequired is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackQuorumNumbersRequired(data []byte) ([]byte, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("quorumNumbersRequired", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackSecurityThresholds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21b9b2fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackSecurityThresholds() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("securityThresholds")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSecurityThresholds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21b9b2fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackSecurityThresholds() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("securityThresholds")
}

// UnpackSecurityThresholds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackSecurityThresholds(data []byte) (EigenDATypesV1SecurityThresholds, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("securityThresholds", data)
	if err != nil {
		return *new(EigenDATypesV1SecurityThresholds), err
	}
	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1SecurityThresholds)).(*EigenDATypesV1SecurityThresholds)
	return out0, nil
}

// PackSemver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcda493c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function semver() pure returns(uint8 major, uint8 minor, uint8 patch)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) PackSemver() []byte {
	enc, err := contractEigenDACertVerifier.abi.Pack("semver")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSemver is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcda493c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function semver() pure returns(uint8 major, uint8 minor, uint8 patch)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) TryPackSemver() ([]byte, error) {
	return contractEigenDACertVerifier.abi.Pack("semver")
}

// SemverOutput serves as a container for the return parameters of contract
// method Semver.
type SemverOutput struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// UnpackSemver is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcda493c8.
//
// Solidity: function semver() pure returns(uint8 major, uint8 minor, uint8 patch)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackSemver(data []byte) (SemverOutput, error) {
	out, err := contractEigenDACertVerifier.abi.Unpack("semver", data)
	outstruct := new(SemverOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Major = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Minor = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Patch = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	return *outstruct, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["BlobQuorumsNotSubset"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackBlobQuorumsNotSubsetError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["InvalidBlobVersion"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackInvalidBlobVersionError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["InvalidInclusionProof"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackInvalidInclusionProofError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["InvalidQuorumNumbersRequired"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackInvalidQuorumNumbersRequiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["InvalidSecurityThresholds"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackInvalidSecurityThresholdsError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["RequiredQuorumsNotSubset"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackRequiredQuorumsNotSubsetError(raw[4:])
	}
	if bytes.Equal(raw[:4], contractEigenDACertVerifier.abi.Errors["SecurityAssumptionsNotMet"].ID.Bytes()[:4]) {
		return contractEigenDACertVerifier.UnpackSecurityAssumptionsNotMetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ContractEigenDACertVerifierBlobQuorumsNotSubset represents a BlobQuorumsNotSubset error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierBlobQuorumsNotSubset struct {
	BlobQuorumsBitmap      *big.Int
	ConfirmedQuorumsBitmap *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BlobQuorumsNotSubset(uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap)
func ContractEigenDACertVerifierBlobQuorumsNotSubsetErrorID() common.Hash {
	return common.HexToHash("0x948e0606890e7792a2da364dbeff7a3f50d7c3f2cf3f5e874bfb0d7276e9b328")
}

// UnpackBlobQuorumsNotSubsetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BlobQuorumsNotSubset(uint256 blobQuorumsBitmap, uint256 confirmedQuorumsBitmap)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackBlobQuorumsNotSubsetError(raw []byte) (*ContractEigenDACertVerifierBlobQuorumsNotSubset, error) {
	out := new(ContractEigenDACertVerifierBlobQuorumsNotSubset)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "BlobQuorumsNotSubset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierInvalidBlobVersion represents a InvalidBlobVersion error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierInvalidBlobVersion struct {
	BlobVersion     uint16
	NextBlobVersion uint16
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidBlobVersion(uint16 blobVersion, uint16 nextBlobVersion)
func ContractEigenDACertVerifierInvalidBlobVersionErrorID() common.Hash {
	return common.HexToHash("0xd6531e7f8a6d92d8e0a5809fddb3accf2cd3b01e5aa4b96867e98835d2185ce2")
}

// UnpackInvalidBlobVersionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidBlobVersion(uint16 blobVersion, uint16 nextBlobVersion)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackInvalidBlobVersionError(raw []byte) (*ContractEigenDACertVerifierInvalidBlobVersion, error) {
	out := new(ContractEigenDACertVerifierInvalidBlobVersion)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "InvalidBlobVersion", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierInvalidInclusionProof represents a InvalidInclusionProof error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierInvalidInclusionProof struct {
	BlobIndex uint32
	BlobHash  [32]byte
	RootHash  [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInclusionProof(uint32 blobIndex, bytes32 blobHash, bytes32 rootHash)
func ContractEigenDACertVerifierInvalidInclusionProofErrorID() common.Hash {
	return common.HexToHash("0x2e547424af90adc34cfc67b4edba519a979d7fc073924797703294a133b1ce11")
}

// UnpackInvalidInclusionProofError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInclusionProof(uint32 blobIndex, bytes32 blobHash, bytes32 rootHash)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackInvalidInclusionProofError(raw []byte) (*ContractEigenDACertVerifierInvalidInclusionProof, error) {
	out := new(ContractEigenDACertVerifierInvalidInclusionProof)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "InvalidInclusionProof", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierInvalidQuorumNumbersRequired represents a InvalidQuorumNumbersRequired error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierInvalidQuorumNumbersRequired struct {
	Length *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidQuorumNumbersRequired(uint256 length)
func ContractEigenDACertVerifierInvalidQuorumNumbersRequiredErrorID() common.Hash {
	return common.HexToHash("0x0008b88edf63cb97efb816fa31f6075f3b46147cf438761a53a85665ce52113a")
}

// UnpackInvalidQuorumNumbersRequiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidQuorumNumbersRequired(uint256 length)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackInvalidQuorumNumbersRequiredError(raw []byte) (*ContractEigenDACertVerifierInvalidQuorumNumbersRequired, error) {
	out := new(ContractEigenDACertVerifierInvalidQuorumNumbersRequired)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "InvalidQuorumNumbersRequired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierInvalidSecurityThresholds represents a InvalidSecurityThresholds error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierInvalidSecurityThresholds struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSecurityThresholds()
func ContractEigenDACertVerifierInvalidSecurityThresholdsErrorID() common.Hash {
	return common.HexToHash("0x08a69975c4c065dd20db258fd793a9eb4231cd659928ecfc755e5cc8047fe11b")
}

// UnpackInvalidSecurityThresholdsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSecurityThresholds()
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackInvalidSecurityThresholdsError(raw []byte) (*ContractEigenDACertVerifierInvalidSecurityThresholds, error) {
	out := new(ContractEigenDACertVerifierInvalidSecurityThresholds)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "InvalidSecurityThresholds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierRequiredQuorumsNotSubset represents a RequiredQuorumsNotSubset error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierRequiredQuorumsNotSubset struct {
	RequiredQuorumsBitmap *big.Int
	BlobQuorumsBitmap     *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RequiredQuorumsNotSubset(uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap)
func ContractEigenDACertVerifierRequiredQuorumsNotSubsetErrorID() common.Hash {
	return common.HexToHash("0x452c216cac89a98c729d0974371a87b40868dd87073b3418ab1bf6e938db3f16")
}

// UnpackRequiredQuorumsNotSubsetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RequiredQuorumsNotSubset(uint256 requiredQuorumsBitmap, uint256 blobQuorumsBitmap)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackRequiredQuorumsNotSubsetError(raw []byte) (*ContractEigenDACertVerifierRequiredQuorumsNotSubset, error) {
	out := new(ContractEigenDACertVerifierRequiredQuorumsNotSubset)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "RequiredQuorumsNotSubset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractEigenDACertVerifierSecurityAssumptionsNotMet represents a SecurityAssumptionsNotMet error raised by the ContractEigenDACertVerifier contract.
type ContractEigenDACertVerifierSecurityAssumptionsNotMet struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
	CodingRate            uint8
	NumChunks             uint32
	MaxNumOperators       uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SecurityAssumptionsNotMet(uint8 confirmationThreshold, uint8 adversaryThreshold, uint8 codingRate, uint32 numChunks, uint32 maxNumOperators)
func ContractEigenDACertVerifierSecurityAssumptionsNotMetErrorID() common.Hash {
	return common.HexToHash("0xf6a44993484a4a6b12403f546a5fe315b0c0c33758393492fac6fbb2a437bd9a")
}

// UnpackSecurityAssumptionsNotMetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SecurityAssumptionsNotMet(uint8 confirmationThreshold, uint8 adversaryThreshold, uint8 codingRate, uint32 numChunks, uint32 maxNumOperators)
func (contractEigenDACertVerifier *ContractEigenDACertVerifier) UnpackSecurityAssumptionsNotMetError(raw []byte) (*ContractEigenDACertVerifierSecurityAssumptionsNotMet, error) {
	out := new(ContractEigenDACertVerifierSecurityAssumptionsNotMet)
	if err := contractEigenDACertVerifier.abi.UnpackIntoInterface(out, "SecurityAssumptionsNotMet", raw); err != nil {
		return nil, err
	}
	return out, nil
}
