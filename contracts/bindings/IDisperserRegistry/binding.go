// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIDisperserRegistry

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

// DisperserRegistryTypesLockedDisperserDeposit is an auto generated low-level Go binding around an user-defined struct.
type DisperserRegistryTypesLockedDisperserDeposit struct {
	Deposit    *big.Int
	Refund     *big.Int
	Token      common.Address
	LockPeriod uint64
}

// ContractIDisperserRegistryMetaData contains all meta data concerning the ContractIDisperserRegistry contract.
var ContractIDisperserRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterDisperser\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDepositParams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserAddress\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserDepositParams\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserDepositUnlockTime\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserOwner\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserURL\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerDisperser\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferDisperserOwnership\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDisperserInfo\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"disperser\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawDisperserDeposit\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"DepositMustBeAtLeastRefund\",\"inputs\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DisperserNotDeregistered\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"DisperserNotRegistered\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidDisperserAddress\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidNewOwner\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidTokenAddress\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotDisperserOwner\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RefundLocked\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unlockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroRefund\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
}

// ContractIDisperserRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIDisperserRegistryMetaData.ABI instead.
var ContractIDisperserRegistryABI = ContractIDisperserRegistryMetaData.ABI

// ContractIDisperserRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIDisperserRegistry struct {
	ContractIDisperserRegistryCaller     // Read-only binding to the contract
	ContractIDisperserRegistryTransactor // Write-only binding to the contract
	ContractIDisperserRegistryFilterer   // Log filterer for contract events
}

// ContractIDisperserRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIDisperserRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIDisperserRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIDisperserRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIDisperserRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIDisperserRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIDisperserRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIDisperserRegistrySession struct {
	Contract     *ContractIDisperserRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ContractIDisperserRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIDisperserRegistryCallerSession struct {
	Contract *ContractIDisperserRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ContractIDisperserRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIDisperserRegistryTransactorSession struct {
	Contract     *ContractIDisperserRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ContractIDisperserRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIDisperserRegistryRaw struct {
	Contract *ContractIDisperserRegistry // Generic contract binding to access the raw methods on
}

// ContractIDisperserRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIDisperserRegistryCallerRaw struct {
	Contract *ContractIDisperserRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIDisperserRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIDisperserRegistryTransactorRaw struct {
	Contract *ContractIDisperserRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIDisperserRegistry creates a new instance of ContractIDisperserRegistry, bound to a specific deployed contract.
func NewContractIDisperserRegistry(address common.Address, backend bind.ContractBackend) (*ContractIDisperserRegistry, error) {
	contract, err := bindContractIDisperserRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIDisperserRegistry{ContractIDisperserRegistryCaller: ContractIDisperserRegistryCaller{contract: contract}, ContractIDisperserRegistryTransactor: ContractIDisperserRegistryTransactor{contract: contract}, ContractIDisperserRegistryFilterer: ContractIDisperserRegistryFilterer{contract: contract}}, nil
}

// NewContractIDisperserRegistryCaller creates a new read-only instance of ContractIDisperserRegistry, bound to a specific deployed contract.
func NewContractIDisperserRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIDisperserRegistryCaller, error) {
	contract, err := bindContractIDisperserRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIDisperserRegistryCaller{contract: contract}, nil
}

// NewContractIDisperserRegistryTransactor creates a new write-only instance of ContractIDisperserRegistry, bound to a specific deployed contract.
func NewContractIDisperserRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIDisperserRegistryTransactor, error) {
	contract, err := bindContractIDisperserRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIDisperserRegistryTransactor{contract: contract}, nil
}

// NewContractIDisperserRegistryFilterer creates a new log filterer instance of ContractIDisperserRegistry, bound to a specific deployed contract.
func NewContractIDisperserRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIDisperserRegistryFilterer, error) {
	contract, err := bindContractIDisperserRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIDisperserRegistryFilterer{contract: contract}, nil
}

// bindContractIDisperserRegistry binds a generic wrapper to an already deployed contract.
func bindContractIDisperserRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIDisperserRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIDisperserRegistry.Contract.ContractIDisperserRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.ContractIDisperserRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.ContractIDisperserRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIDisperserRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDepositParams(opts *bind.CallOpts) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDepositParams")

	if err != nil {
		return *new(DisperserRegistryTypesLockedDisperserDeposit), err
	}

	out0 := *abi.ConvertType(out[0], new(DisperserRegistryTypesLockedDisperserDeposit)).(*DisperserRegistryTypesLockedDisperserDeposit)

	return out0, err

}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDepositParams() (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractIDisperserRegistry.Contract.GetDepositParams(&_ContractIDisperserRegistry.CallOpts)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDepositParams() (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractIDisperserRegistry.Contract.GetDepositParams(&_ContractIDisperserRegistry.CallOpts)
}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDisperserAddress(opts *bind.CallOpts, disperserKey uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDisperserAddress", disperserKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDisperserAddress(disperserKey uint32) (common.Address, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserAddress(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDisperserAddress(disperserKey uint32) (common.Address, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserAddress(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDisperserDepositParams(opts *bind.CallOpts, disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDisperserDepositParams", disperserKey)

	if err != nil {
		return *new(DisperserRegistryTypesLockedDisperserDeposit), err
	}

	out0 := *abi.ConvertType(out[0], new(DisperserRegistryTypesLockedDisperserDeposit)).(*DisperserRegistryTypesLockedDisperserDeposit)

	return out0, err

}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDisperserDepositParams(disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserDepositParams(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDisperserDepositParams(disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserDepositParams(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDisperserDepositUnlockTime(opts *bind.CallOpts, disperserKey uint32) (uint64, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDisperserDepositUnlockTime", disperserKey)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDisperserDepositUnlockTime(disperserKey uint32) (uint64, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserDepositUnlockTime(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDisperserDepositUnlockTime(disperserKey uint32) (uint64, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserDepositUnlockTime(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDisperserOwner(opts *bind.CallOpts, disperserKey uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDisperserOwner", disperserKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDisperserOwner(disperserKey uint32) (common.Address, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserOwner(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDisperserOwner(disperserKey uint32) (common.Address, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserOwner(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCaller) GetDisperserURL(opts *bind.CallOpts, disperserKey uint32) (string, error) {
	var out []interface{}
	err := _ContractIDisperserRegistry.contract.Call(opts, &out, "getDisperserURL", disperserKey)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) GetDisperserURL(disperserKey uint32) (string, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserURL(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryCallerSession) GetDisperserURL(disperserKey uint32) (string, error) {
	return _ContractIDisperserRegistry.Contract.GetDisperserURL(&_ContractIDisperserRegistry.CallOpts, disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactor) DeregisterDisperser(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.contract.Transact(opts, "deregisterDisperser", disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.DeregisterDisperser(&_ContractIDisperserRegistry.TransactOpts, disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorSession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.DeregisterDisperser(&_ContractIDisperserRegistry.TransactOpts, disperserKey)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactor) RegisterDisperser(opts *bind.TransactOpts, disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.contract.Transact(opts, "registerDisperser", disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.RegisterDisperser(&_ContractIDisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorSession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.RegisterDisperser(&_ContractIDisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address newOwner) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactor) TransferDisperserOwnership(opts *bind.TransactOpts, disperserKey uint32, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.contract.Transact(opts, "transferDisperserOwnership", disperserKey, newOwner)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address newOwner) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) TransferDisperserOwnership(disperserKey uint32, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.TransferDisperserOwnership(&_ContractIDisperserRegistry.TransactOpts, disperserKey, newOwner)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address newOwner) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorSession) TransferDisperserOwnership(disperserKey uint32, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.TransferDisperserOwnership(&_ContractIDisperserRegistry.TransactOpts, disperserKey, newOwner)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactor) UpdateDisperserInfo(opts *bind.TransactOpts, disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.contract.Transact(opts, "updateDisperserInfo", disperserKey, disperser, disperserURL)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) UpdateDisperserInfo(disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.UpdateDisperserInfo(&_ContractIDisperserRegistry.TransactOpts, disperserKey, disperser, disperserURL)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorSession) UpdateDisperserInfo(disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.UpdateDisperserInfo(&_ContractIDisperserRegistry.TransactOpts, disperserKey, disperser, disperserURL)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactor) WithdrawDisperserDeposit(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.contract.Transact(opts, "withdrawDisperserDeposit", disperserKey)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistrySession) WithdrawDisperserDeposit(disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.WithdrawDisperserDeposit(&_ContractIDisperserRegistry.TransactOpts, disperserKey)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractIDisperserRegistry *ContractIDisperserRegistryTransactorSession) WithdrawDisperserDeposit(disperserKey uint32) (*types.Transaction, error) {
	return _ContractIDisperserRegistry.Contract.WithdrawDisperserDeposit(&_ContractIDisperserRegistry.TransactOpts, disperserKey)
}
