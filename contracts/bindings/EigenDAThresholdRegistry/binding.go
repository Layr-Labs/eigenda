// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAThresholdRegistry

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

// ContractEigenDAThresholdRegistryMetaData contains all meta data concerning the ContractEigenDAThresholdRegistry contract.
var ContractEigenDAThresholdRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVersionedBlobParams\",\"inputs\":[{\"name\":\"_versionedBlobParams\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_quorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_quorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_versionedBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams[]\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextBlobVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"versionedBlobParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
}

// ContractEigenDAThresholdRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAThresholdRegistryMetaData.ABI instead.
var ContractEigenDAThresholdRegistryABI = ContractEigenDAThresholdRegistryMetaData.ABI

// ContractEigenDAThresholdRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistry struct {
	ContractEigenDAThresholdRegistryCaller     // Read-only binding to the contract
	ContractEigenDAThresholdRegistryTransactor // Write-only binding to the contract
	ContractEigenDAThresholdRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDAThresholdRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAThresholdRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAThresholdRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAThresholdRegistrySession struct {
	Contract     *ContractEigenDAThresholdRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractEigenDAThresholdRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAThresholdRegistryCallerSession struct {
	Contract *ContractEigenDAThresholdRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractEigenDAThresholdRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAThresholdRegistryTransactorSession struct {
	Contract     *ContractEigenDAThresholdRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractEigenDAThresholdRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryRaw struct {
	Contract *ContractEigenDAThresholdRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDAThresholdRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryCallerRaw struct {
	Contract *ContractEigenDAThresholdRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAThresholdRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAThresholdRegistryTransactorRaw struct {
	Contract *ContractEigenDAThresholdRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAThresholdRegistry creates a new instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDAThresholdRegistry, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistry{ContractEigenDAThresholdRegistryCaller: ContractEigenDAThresholdRegistryCaller{contract: contract}, ContractEigenDAThresholdRegistryTransactor: ContractEigenDAThresholdRegistryTransactor{contract: contract}, ContractEigenDAThresholdRegistryFilterer: ContractEigenDAThresholdRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDAThresholdRegistryCaller creates a new read-only instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAThresholdRegistryCaller, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryCaller{contract: contract}, nil
}

// NewContractEigenDAThresholdRegistryTransactor creates a new write-only instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAThresholdRegistryTransactor, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDAThresholdRegistryFilterer creates a new log filterer instance of ContractEigenDAThresholdRegistry, bound to a specific deployed contract.
func NewContractEigenDAThresholdRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAThresholdRegistryFilterer, error) {
	contract, err := bindContractEigenDAThresholdRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDAThresholdRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDAThresholdRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAThresholdRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.ContractEigenDAThresholdRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(EigenDATypesV1VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1VersionedBlobParams)).(*EigenDATypesV1VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetIsQuorumRequired(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetIsQuorumRequired(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAThresholdRegistry.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAThresholdRegistry.CallOpts, quorumNumber)
}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) NextBlobVersion(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "nextBlobVersion")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) NextBlobVersion() (uint16, error) {
	return _ContractEigenDAThresholdRegistry.Contract.NextBlobVersion(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// NextBlobVersion is a free data retrieval call binding the contract method 0x32430f14.
//
// Solidity: function nextBlobVersion() view returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) NextBlobVersion() (uint16, error) {
	return _ContractEigenDAThresholdRegistry.Contract.NextBlobVersion(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) Owner() (common.Address, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Owner(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Owner(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAThresholdRegistry.Contract.QuorumNumbersRequired(&_ContractEigenDAThresholdRegistry.CallOpts)
}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCaller) VersionedBlobParams(opts *bind.CallOpts, arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	var out []interface{}
	err := _ContractEigenDAThresholdRegistry.contract.Call(opts, &out, "versionedBlobParams", arg0)

	outstruct := new(struct {
		MaxNumOperators uint32
		NumChunks       uint32
		CodingRate      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxNumOperators = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.NumChunks = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.CodingRate = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) VersionedBlobParams(arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.VersionedBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, arg0)
}

// VersionedBlobParams is a free data retrieval call binding the contract method 0xf74e363c.
//
// Solidity: function versionedBlobParams(uint16 ) view returns(uint32 maxNumOperators, uint32 numChunks, uint8 codingRate)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryCallerSession) VersionedBlobParams(arg0 uint16) (struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}, error) {
	return _ContractEigenDAThresholdRegistry.Contract.VersionedBlobParams(&_ContractEigenDAThresholdRegistry.CallOpts, arg0)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) AddVersionedBlobParams(opts *bind.TransactOpts, _versionedBlobParams EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "addVersionedBlobParams", _versionedBlobParams)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) AddVersionedBlobParams(_versionedBlobParams EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.AddVersionedBlobParams(&_ContractEigenDAThresholdRegistry.TransactOpts, _versionedBlobParams)
}

// AddVersionedBlobParams is a paid mutator transaction binding the contract method 0x8a476982.
//
// Solidity: function addVersionedBlobParams((uint32,uint32,uint8) _versionedBlobParams) returns(uint16)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) AddVersionedBlobParams(_versionedBlobParams EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.AddVersionedBlobParams(&_ContractEigenDAThresholdRegistry.TransactOpts, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "initialize", _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8491bad6.
//
// Solidity: function initialize(address _initialOwner, bytes _quorumAdversaryThresholdPercentages, bytes _quorumConfirmationThresholdPercentages, bytes _quorumNumbersRequired, (uint32,uint32,uint8)[] _versionedBlobParams) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) Initialize(_initialOwner common.Address, _quorumAdversaryThresholdPercentages []byte, _quorumConfirmationThresholdPercentages []byte, _quorumNumbersRequired []byte, _versionedBlobParams []EigenDATypesV1VersionedBlobParams) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.Initialize(&_ContractEigenDAThresholdRegistry.TransactOpts, _initialOwner, _quorumAdversaryThresholdPercentages, _quorumConfirmationThresholdPercentages, _quorumNumbersRequired, _versionedBlobParams)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.RenounceOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.RenounceOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.TransferOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAThresholdRegistry.Contract.TransferOwnership(&_ContractEigenDAThresholdRegistry.TransactOpts, newOwner)
}

// ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
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
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 EigenDATypesV1SecurityThresholds
	NewDefaultSecurityThresholdsV2      EigenDATypesV1SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractEigenDAThresholdRegistryDefaultSecurityThresholdsV2Updated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryInitializedIterator struct {
	Event *ContractEigenDAThresholdRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryInitialized)
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
		it.Event = new(ContractEigenDAThresholdRegistryInitialized)
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
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryInitialized represents a Initialized event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryInitializedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryInitialized)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseInitialized(log types.Log) (*ContractEigenDAThresholdRegistryInitialized, error) {
	event := new(ContractEigenDAThresholdRegistryInitialized)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryOwnershipTransferredIterator struct {
	Event *ContractEigenDAThresholdRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryOwnershipTransferred)
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
		it.Event = new(ContractEigenDAThresholdRegistryOwnershipTransferred)
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
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDAThresholdRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryOwnershipTransferredIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryOwnershipTransferred)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDAThresholdRegistryOwnershipTransferred, error) {
	event := new(ContractEigenDAThresholdRegistryOwnershipTransferred)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
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
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdatedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated, error) {
	event := new(ContractEigenDAThresholdRegistryQuorumNumbersRequiredUpdated)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator struct {
	Event *ContractEigenDAThresholdRegistryVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
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
		it.Event = new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
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
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAThresholdRegistryVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractEigenDAThresholdRegistry contract.
type ContractEigenDAThresholdRegistryVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams EigenDATypesV1VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAThresholdRegistryVersionedBlobParamsAddedIterator{contract: _ContractEigenDAThresholdRegistry.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDAThresholdRegistryVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAThresholdRegistry.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
				if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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
func (_ContractEigenDAThresholdRegistry *ContractEigenDAThresholdRegistryFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractEigenDAThresholdRegistryVersionedBlobParamsAdded, error) {
	event := new(ContractEigenDAThresholdRegistryVersionedBlobParamsAdded)
	if err := _ContractEigenDAThresholdRegistry.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
