// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractDisperserRegistry

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

// ContractDisperserRegistryMetaData contains all meta data concerning the ContractDisperserRegistry contract.
var ContractDisperserRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterDisperser\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDepositParams\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserAddress\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserDepositParams\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserDepositUnlockTime\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserOwner\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDisperserURL\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExcessBalance\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextDisperserKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUpdateFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"name\":\"updateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerDisperser\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositParams\",\"inputs\":[{\"name\":\"depositParams\",\"type\":\"tuple\",\"internalType\":\"structDisperserRegistryTypes.LockedDisperserDeposit\",\"components\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lockPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUpdateFee\",\"inputs\":[{\"name\":\"updateFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferDisperserOwnership\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateDisperserInfo\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"disperser\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"disperserURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawDisperserDeposit\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositMustBeAtLeastRefund\",\"inputs\":[{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DisperserNotDeregistered\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"DisperserNotRegistered\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidDisperserAddress\",\"inputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidNewOwner\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidTokenAddress\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MissingRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotDisperserOwner\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RefundLocked\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"unlockTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroRefund\",\"inputs\":[{\"name\":\"disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
}

// ContractDisperserRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractDisperserRegistryMetaData.ABI instead.
var ContractDisperserRegistryABI = ContractDisperserRegistryMetaData.ABI

// ContractDisperserRegistry is an auto generated Go binding around an Ethereum contract.
type ContractDisperserRegistry struct {
	ContractDisperserRegistryCaller     // Read-only binding to the contract
	ContractDisperserRegistryTransactor // Write-only binding to the contract
	ContractDisperserRegistryFilterer   // Log filterer for contract events
}

// ContractDisperserRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractDisperserRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDisperserRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractDisperserRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDisperserRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractDisperserRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDisperserRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractDisperserRegistrySession struct {
	Contract     *ContractDisperserRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ContractDisperserRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractDisperserRegistryCallerSession struct {
	Contract *ContractDisperserRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ContractDisperserRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractDisperserRegistryTransactorSession struct {
	Contract     *ContractDisperserRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ContractDisperserRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractDisperserRegistryRaw struct {
	Contract *ContractDisperserRegistry // Generic contract binding to access the raw methods on
}

// ContractDisperserRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractDisperserRegistryCallerRaw struct {
	Contract *ContractDisperserRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractDisperserRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractDisperserRegistryTransactorRaw struct {
	Contract *ContractDisperserRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractDisperserRegistry creates a new instance of ContractDisperserRegistry, bound to a specific deployed contract.
func NewContractDisperserRegistry(address common.Address, backend bind.ContractBackend) (*ContractDisperserRegistry, error) {
	contract, err := bindContractDisperserRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractDisperserRegistry{ContractDisperserRegistryCaller: ContractDisperserRegistryCaller{contract: contract}, ContractDisperserRegistryTransactor: ContractDisperserRegistryTransactor{contract: contract}, ContractDisperserRegistryFilterer: ContractDisperserRegistryFilterer{contract: contract}}, nil
}

// NewContractDisperserRegistryCaller creates a new read-only instance of ContractDisperserRegistry, bound to a specific deployed contract.
func NewContractDisperserRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractDisperserRegistryCaller, error) {
	contract, err := bindContractDisperserRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractDisperserRegistryCaller{contract: contract}, nil
}

// NewContractDisperserRegistryTransactor creates a new write-only instance of ContractDisperserRegistry, bound to a specific deployed contract.
func NewContractDisperserRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractDisperserRegistryTransactor, error) {
	contract, err := bindContractDisperserRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractDisperserRegistryTransactor{contract: contract}, nil
}

// NewContractDisperserRegistryFilterer creates a new log filterer instance of ContractDisperserRegistry, bound to a specific deployed contract.
func NewContractDisperserRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractDisperserRegistryFilterer, error) {
	contract, err := bindContractDisperserRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractDisperserRegistryFilterer{contract: contract}, nil
}

// bindContractDisperserRegistry binds a generic wrapper to an already deployed contract.
func bindContractDisperserRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractDisperserRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractDisperserRegistry *ContractDisperserRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractDisperserRegistry.Contract.ContractDisperserRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractDisperserRegistry *ContractDisperserRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.ContractDisperserRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractDisperserRegistry *ContractDisperserRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.ContractDisperserRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractDisperserRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDepositParams(opts *bind.CallOpts) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDepositParams")

	if err != nil {
		return *new(DisperserRegistryTypesLockedDisperserDeposit), err
	}

	out0 := *abi.ConvertType(out[0], new(DisperserRegistryTypesLockedDisperserDeposit)).(*DisperserRegistryTypesLockedDisperserDeposit)

	return out0, err

}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDepositParams() (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractDisperserRegistry.Contract.GetDepositParams(&_ContractDisperserRegistry.CallOpts)
}

// GetDepositParams is a free data retrieval call binding the contract method 0x8e1e4829.
//
// Solidity: function getDepositParams() view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDepositParams() (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractDisperserRegistry.Contract.GetDepositParams(&_ContractDisperserRegistry.CallOpts)
}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDisperserAddress(opts *bind.CallOpts, disperserKey uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDisperserAddress", disperserKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDisperserAddress(disperserKey uint32) (common.Address, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserAddress(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserAddress is a free data retrieval call binding the contract method 0xbf9142b3.
//
// Solidity: function getDisperserAddress(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDisperserAddress(disperserKey uint32) (common.Address, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserAddress(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDisperserDepositParams(opts *bind.CallOpts, disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDisperserDepositParams", disperserKey)

	if err != nil {
		return *new(DisperserRegistryTypesLockedDisperserDeposit), err
	}

	out0 := *abi.ConvertType(out[0], new(DisperserRegistryTypesLockedDisperserDeposit)).(*DisperserRegistryTypesLockedDisperserDeposit)

	return out0, err

}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDisperserDepositParams(disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserDepositParams(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositParams is a free data retrieval call binding the contract method 0x3613d42e.
//
// Solidity: function getDisperserDepositParams(uint32 disperserKey) view returns((uint256,uint256,address,uint64))
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDisperserDepositParams(disperserKey uint32) (DisperserRegistryTypesLockedDisperserDeposit, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserDepositParams(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDisperserDepositUnlockTime(opts *bind.CallOpts, disperserKey uint32) (uint64, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDisperserDepositUnlockTime", disperserKey)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDisperserDepositUnlockTime(disperserKey uint32) (uint64, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserDepositUnlockTime(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserDepositUnlockTime is a free data retrieval call binding the contract method 0x00926944.
//
// Solidity: function getDisperserDepositUnlockTime(uint32 disperserKey) view returns(uint64)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDisperserDepositUnlockTime(disperserKey uint32) (uint64, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserDepositUnlockTime(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDisperserOwner(opts *bind.CallOpts, disperserKey uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDisperserOwner", disperserKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDisperserOwner(disperserKey uint32) (common.Address, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserOwner(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserOwner is a free data retrieval call binding the contract method 0x518e40f7.
//
// Solidity: function getDisperserOwner(uint32 disperserKey) view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDisperserOwner(disperserKey uint32) (common.Address, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserOwner(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetDisperserURL(opts *bind.CallOpts, disperserKey uint32) (string, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getDisperserURL", disperserKey)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetDisperserURL(disperserKey uint32) (string, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserURL(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetDisperserURL is a free data retrieval call binding the contract method 0x85c11f1b.
//
// Solidity: function getDisperserURL(uint32 disperserKey) view returns(string)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetDisperserURL(disperserKey uint32) (string, error) {
	return _ContractDisperserRegistry.Contract.GetDisperserURL(&_ContractDisperserRegistry.CallOpts, disperserKey)
}

// GetExcessBalance is a free data retrieval call binding the contract method 0xace5592c.
//
// Solidity: function getExcessBalance(address token) view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetExcessBalance(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getExcessBalance", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExcessBalance is a free data retrieval call binding the contract method 0xace5592c.
//
// Solidity: function getExcessBalance(address token) view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetExcessBalance(token common.Address) (*big.Int, error) {
	return _ContractDisperserRegistry.Contract.GetExcessBalance(&_ContractDisperserRegistry.CallOpts, token)
}

// GetExcessBalance is a free data retrieval call binding the contract method 0xace5592c.
//
// Solidity: function getExcessBalance(address token) view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetExcessBalance(token common.Address) (*big.Int, error) {
	return _ContractDisperserRegistry.Contract.GetExcessBalance(&_ContractDisperserRegistry.CallOpts, token)
}

// GetNextDisperserKey is a free data retrieval call binding the contract method 0x1af77b47.
//
// Solidity: function getNextDisperserKey() view returns(uint32)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetNextDisperserKey(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getNextDisperserKey")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetNextDisperserKey is a free data retrieval call binding the contract method 0x1af77b47.
//
// Solidity: function getNextDisperserKey() view returns(uint32)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetNextDisperserKey() (uint32, error) {
	return _ContractDisperserRegistry.Contract.GetNextDisperserKey(&_ContractDisperserRegistry.CallOpts)
}

// GetNextDisperserKey is a free data retrieval call binding the contract method 0x1af77b47.
//
// Solidity: function getNextDisperserKey() view returns(uint32)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetNextDisperserKey() (uint32, error) {
	return _ContractDisperserRegistry.Contract.GetNextDisperserKey(&_ContractDisperserRegistry.CallOpts)
}

// GetUpdateFee is a free data retrieval call binding the contract method 0xb7b0ccde.
//
// Solidity: function getUpdateFee() view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) GetUpdateFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "getUpdateFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUpdateFee is a free data retrieval call binding the contract method 0xb7b0ccde.
//
// Solidity: function getUpdateFee() view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) GetUpdateFee() (*big.Int, error) {
	return _ContractDisperserRegistry.Contract.GetUpdateFee(&_ContractDisperserRegistry.CallOpts)
}

// GetUpdateFee is a free data retrieval call binding the contract method 0xb7b0ccde.
//
// Solidity: function getUpdateFee() view returns(uint256)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) GetUpdateFee() (*big.Int, error) {
	return _ContractDisperserRegistry.Contract.GetUpdateFee(&_ContractDisperserRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDisperserRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) Owner() (common.Address, error) {
	return _ContractDisperserRegistry.Contract.Owner(&_ContractDisperserRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDisperserRegistry *ContractDisperserRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractDisperserRegistry.Contract.Owner(&_ContractDisperserRegistry.CallOpts)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) DeregisterDisperser(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "deregisterDisperser", disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.DeregisterDisperser(&_ContractDisperserRegistry.TransactOpts, disperserKey)
}

// DeregisterDisperser is a paid mutator transaction binding the contract method 0x4fe5de32.
//
// Solidity: function deregisterDisperser(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) DeregisterDisperser(disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.DeregisterDisperser(&_ContractDisperserRegistry.TransactOpts, disperserKey)
}

// Initialize is a paid mutator transaction binding the contract method 0xa0794864.
//
// Solidity: function initialize(address initialOwner, (uint256,uint256,address,uint64) depositParams, uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, depositParams DisperserRegistryTypesLockedDisperserDeposit, updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "initialize", initialOwner, depositParams, updateFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xa0794864.
//
// Solidity: function initialize(address initialOwner, (uint256,uint256,address,uint64) depositParams, uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) Initialize(initialOwner common.Address, depositParams DisperserRegistryTypesLockedDisperserDeposit, updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.Initialize(&_ContractDisperserRegistry.TransactOpts, initialOwner, depositParams, updateFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xa0794864.
//
// Solidity: function initialize(address initialOwner, (uint256,uint256,address,uint64) depositParams, uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) Initialize(initialOwner common.Address, depositParams DisperserRegistryTypesLockedDisperserDeposit, updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.Initialize(&_ContractDisperserRegistry.TransactOpts, initialOwner, depositParams, updateFee)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) RegisterDisperser(opts *bind.TransactOpts, disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "registerDisperser", disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.RegisterDisperser(&_ContractDisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// RegisterDisperser is a paid mutator transaction binding the contract method 0x1b4b940a.
//
// Solidity: function registerDisperser(address disperserAddress, string disperserURL) returns(uint32 disperserKey)
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) RegisterDisperser(disperserAddress common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.RegisterDisperser(&_ContractDisperserRegistry.TransactOpts, disperserAddress, disperserURL)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) SetDepositParams(opts *bind.TransactOpts, depositParams DisperserRegistryTypesLockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "setDepositParams", depositParams)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) SetDepositParams(depositParams DisperserRegistryTypesLockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.SetDepositParams(&_ContractDisperserRegistry.TransactOpts, depositParams)
}

// SetDepositParams is a paid mutator transaction binding the contract method 0x6923ef3a.
//
// Solidity: function setDepositParams((uint256,uint256,address,uint64) depositParams) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) SetDepositParams(depositParams DisperserRegistryTypesLockedDisperserDeposit) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.SetDepositParams(&_ContractDisperserRegistry.TransactOpts, depositParams)
}

// SetUpdateFee is a paid mutator transaction binding the contract method 0x6a2e770f.
//
// Solidity: function setUpdateFee(uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) SetUpdateFee(opts *bind.TransactOpts, updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "setUpdateFee", updateFee)
}

// SetUpdateFee is a paid mutator transaction binding the contract method 0x6a2e770f.
//
// Solidity: function setUpdateFee(uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) SetUpdateFee(updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.SetUpdateFee(&_ContractDisperserRegistry.TransactOpts, updateFee)
}

// SetUpdateFee is a paid mutator transaction binding the contract method 0x6a2e770f.
//
// Solidity: function setUpdateFee(uint256 updateFee) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) SetUpdateFee(updateFee *big.Int) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.SetUpdateFee(&_ContractDisperserRegistry.TransactOpts, updateFee)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) TransferDisperserOwnership(opts *bind.TransactOpts, disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "transferDisperserOwnership", disperserKey, disperserAddress)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) TransferDisperserOwnership(disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.TransferDisperserOwnership(&_ContractDisperserRegistry.TransactOpts, disperserKey, disperserAddress)
}

// TransferDisperserOwnership is a paid mutator transaction binding the contract method 0xbb4c7dda.
//
// Solidity: function transferDisperserOwnership(uint32 disperserKey, address disperserAddress) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) TransferDisperserOwnership(disperserKey uint32, disperserAddress common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.TransferDisperserOwnership(&_ContractDisperserRegistry.TransactOpts, disperserKey, disperserAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.TransferOwnership(&_ContractDisperserRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.TransferOwnership(&_ContractDisperserRegistry.TransactOpts, newOwner)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) UpdateDisperserInfo(opts *bind.TransactOpts, disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "updateDisperserInfo", disperserKey, disperser, disperserURL)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) UpdateDisperserInfo(disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.UpdateDisperserInfo(&_ContractDisperserRegistry.TransactOpts, disperserKey, disperser, disperserURL)
}

// UpdateDisperserInfo is a paid mutator transaction binding the contract method 0xc4dc08bc.
//
// Solidity: function updateDisperserInfo(uint32 disperserKey, address disperser, string disperserURL) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) UpdateDisperserInfo(disperserKey uint32, disperser common.Address, disperserURL string) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.UpdateDisperserInfo(&_ContractDisperserRegistry.TransactOpts, disperserKey, disperser, disperserURL)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactor) WithdrawDisperserDeposit(opts *bind.TransactOpts, disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.contract.Transact(opts, "withdrawDisperserDeposit", disperserKey)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistrySession) WithdrawDisperserDeposit(disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.WithdrawDisperserDeposit(&_ContractDisperserRegistry.TransactOpts, disperserKey)
}

// WithdrawDisperserDeposit is a paid mutator transaction binding the contract method 0xad839358.
//
// Solidity: function withdrawDisperserDeposit(uint32 disperserKey) returns()
func (_ContractDisperserRegistry *ContractDisperserRegistryTransactorSession) WithdrawDisperserDeposit(disperserKey uint32) (*types.Transaction, error) {
	return _ContractDisperserRegistry.Contract.WithdrawDisperserDeposit(&_ContractDisperserRegistry.TransactOpts, disperserKey)
}
