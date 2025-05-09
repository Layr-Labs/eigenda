// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDACertTypeBindings

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

// EigenDACertTypesEigenDACertV3 is an auto generated low-level Go binding around an user-defined struct.
type EigenDACertTypesEigenDACertV3 struct {
	BatchHeader                 EigenDATypesV2BatchHeaderV2
	BlobInclusionInfo           EigenDATypesV2BlobInclusionInfo
	NonSignerStakesAndSignature EigenDATypesV1NonSignerStakesAndSignature
	SignedQuorumNumbers         []byte
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

// ContractIEigenDACertTypeBindingsMetaData contains all meta data concerning the ContractIEigenDACertTypeBindings contract.
var ContractIEigenDACertTypeBindingsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"dummyFnCertV3\",\"inputs\":[{\"name\":\"cert\",\"type\":\"tuple\",\"internalType\":\"structEigenDACertTypes.EigenDACertV3\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BatchHeaderV2\",\"components\":[{\"name\":\"batchRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobInclusionInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobInclusionInfo\",\"components\":[{\"name\":\"blobCertificate\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCertificate\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobHeaderV2\",\"components\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.BlobCommitment\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"lengthCommitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"lengthProof\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"length\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"paymentHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"relayKeys\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}]},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"signedQuorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"}]",
}

// ContractIEigenDACertTypeBindingsABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDACertTypeBindingsMetaData.ABI instead.
var ContractIEigenDACertTypeBindingsABI = ContractIEigenDACertTypeBindingsMetaData.ABI

// ContractIEigenDACertTypeBindings is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindings struct {
	ContractIEigenDACertTypeBindingsCaller     // Read-only binding to the contract
	ContractIEigenDACertTypeBindingsTransactor // Write-only binding to the contract
	ContractIEigenDACertTypeBindingsFilterer   // Log filterer for contract events
}

// ContractIEigenDACertTypeBindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertTypeBindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertTypeBindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDACertTypeBindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDACertTypeBindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDACertTypeBindingsSession struct {
	Contract     *ContractIEigenDACertTypeBindings // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractIEigenDACertTypeBindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDACertTypeBindingsCallerSession struct {
	Contract *ContractIEigenDACertTypeBindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractIEigenDACertTypeBindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDACertTypeBindingsTransactorSession struct {
	Contract     *ContractIEigenDACertTypeBindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractIEigenDACertTypeBindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindingsRaw struct {
	Contract *ContractIEigenDACertTypeBindings // Generic contract binding to access the raw methods on
}

// ContractIEigenDACertTypeBindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindingsCallerRaw struct {
	Contract *ContractIEigenDACertTypeBindingsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDACertTypeBindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDACertTypeBindingsTransactorRaw struct {
	Contract *ContractIEigenDACertTypeBindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDACertTypeBindings creates a new instance of ContractIEigenDACertTypeBindings, bound to a specific deployed contract.
func NewContractIEigenDACertTypeBindings(address common.Address, backend bind.ContractBackend) (*ContractIEigenDACertTypeBindings, error) {
	contract, err := bindContractIEigenDACertTypeBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertTypeBindings{ContractIEigenDACertTypeBindingsCaller: ContractIEigenDACertTypeBindingsCaller{contract: contract}, ContractIEigenDACertTypeBindingsTransactor: ContractIEigenDACertTypeBindingsTransactor{contract: contract}, ContractIEigenDACertTypeBindingsFilterer: ContractIEigenDACertTypeBindingsFilterer{contract: contract}}, nil
}

// NewContractIEigenDACertTypeBindingsCaller creates a new read-only instance of ContractIEigenDACertTypeBindings, bound to a specific deployed contract.
func NewContractIEigenDACertTypeBindingsCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDACertTypeBindingsCaller, error) {
	contract, err := bindContractIEigenDACertTypeBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertTypeBindingsCaller{contract: contract}, nil
}

// NewContractIEigenDACertTypeBindingsTransactor creates a new write-only instance of ContractIEigenDACertTypeBindings, bound to a specific deployed contract.
func NewContractIEigenDACertTypeBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDACertTypeBindingsTransactor, error) {
	contract, err := bindContractIEigenDACertTypeBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertTypeBindingsTransactor{contract: contract}, nil
}

// NewContractIEigenDACertTypeBindingsFilterer creates a new log filterer instance of ContractIEigenDACertTypeBindings, bound to a specific deployed contract.
func NewContractIEigenDACertTypeBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDACertTypeBindingsFilterer, error) {
	contract, err := bindContractIEigenDACertTypeBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDACertTypeBindingsFilterer{contract: contract}, nil
}

// bindContractIEigenDACertTypeBindings binds a generic wrapper to an already deployed contract.
func bindContractIEigenDACertTypeBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDACertTypeBindingsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDACertTypeBindings.Contract.ContractIEigenDACertTypeBindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDACertTypeBindings.Contract.ContractIEigenDACertTypeBindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDACertTypeBindings.Contract.ContractIEigenDACertTypeBindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDACertTypeBindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDACertTypeBindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDACertTypeBindings.Contract.contract.Transact(opts, method, params...)
}

// DummyFnCertV3 is a free data retrieval call binding the contract method 0xec06d922.
//
// Solidity: function dummyFnCertV3(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsCaller) DummyFnCertV3(opts *bind.CallOpts, cert EigenDACertTypesEigenDACertV3) error {
	var out []interface{}
	err := _ContractIEigenDACertTypeBindings.contract.Call(opts, &out, "dummyFnCertV3", cert)

	if err != nil {
		return err
	}

	return err

}

// DummyFnCertV3 is a free data retrieval call binding the contract method 0xec06d922.
//
// Solidity: function dummyFnCertV3(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsSession) DummyFnCertV3(cert EigenDACertTypesEigenDACertV3) error {
	return _ContractIEigenDACertTypeBindings.Contract.DummyFnCertV3(&_ContractIEigenDACertTypeBindings.CallOpts, cert)
}

// DummyFnCertV3 is a free data retrieval call binding the contract method 0xec06d922.
//
// Solidity: function dummyFnCertV3(((bytes32,uint32),(((uint16,bytes,((uint256,uint256),(uint256[2],uint256[2]),(uint256[2],uint256[2]),uint32),bytes32),bytes,uint32[]),uint32,bytes),(uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]),bytes) cert) view returns()
func (_ContractIEigenDACertTypeBindings *ContractIEigenDACertTypeBindingsCallerSession) DummyFnCertV3(cert EigenDACertTypesEigenDACertV3) error {
	return _ContractIEigenDACertTypeBindings.Contract.DummyFnCertV3(&_ContractIEigenDACertTypeBindings.CallOpts, cert)
}
