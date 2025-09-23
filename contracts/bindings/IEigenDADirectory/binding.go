// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDADirectory

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

// ContractIEigenDADirectoryMetaData contains all meta data concerning the ContractIEigenDADirectory contract.
var ContractIEigenDADirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigBytes\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigBytes32\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNames\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes32\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes32ExtraInfo\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes32ExtraInfo\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytesExtraInfo\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytesExtraInfo\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getName\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumRegisteredKeysBytes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumRegisteredKeysBytes32\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredKeyBytes\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredKeyBytes32\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeConfigBytes\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeConfigBytes32\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceConfigBytes\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceConfigBytes32\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddressAdded\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressReplaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytes32Added\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytes32Removed\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytes32Replaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newValue\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytesAdded\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytesRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConfigBytesReplaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newValue\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"extraInfo\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"AddressDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConfigAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConfigDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NewValueIsOldValue\",\"inputs\":[{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// ContractIEigenDADirectoryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDADirectoryMetaData.ABI instead.
var ContractIEigenDADirectoryABI = ContractIEigenDADirectoryMetaData.ABI

// ContractIEigenDADirectory is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDADirectory struct {
	ContractIEigenDADirectoryCaller     // Read-only binding to the contract
	ContractIEigenDADirectoryTransactor // Write-only binding to the contract
	ContractIEigenDADirectoryFilterer   // Log filterer for contract events
}

// ContractIEigenDADirectoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDADirectoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDADirectoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDADirectoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDADirectoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDADirectoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDADirectorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDADirectorySession struct {
	Contract     *ContractIEigenDADirectory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ContractIEigenDADirectoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDADirectoryCallerSession struct {
	Contract *ContractIEigenDADirectoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ContractIEigenDADirectoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDADirectoryTransactorSession struct {
	Contract     *ContractIEigenDADirectoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ContractIEigenDADirectoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDADirectoryRaw struct {
	Contract *ContractIEigenDADirectory // Generic contract binding to access the raw methods on
}

// ContractIEigenDADirectoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDADirectoryCallerRaw struct {
	Contract *ContractIEigenDADirectoryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDADirectoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDADirectoryTransactorRaw struct {
	Contract *ContractIEigenDADirectoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDADirectory creates a new instance of ContractIEigenDADirectory, bound to a specific deployed contract.
func NewContractIEigenDADirectory(address common.Address, backend bind.ContractBackend) (*ContractIEigenDADirectory, error) {
	contract, err := bindContractIEigenDADirectory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectory{ContractIEigenDADirectoryCaller: ContractIEigenDADirectoryCaller{contract: contract}, ContractIEigenDADirectoryTransactor: ContractIEigenDADirectoryTransactor{contract: contract}, ContractIEigenDADirectoryFilterer: ContractIEigenDADirectoryFilterer{contract: contract}}, nil
}

// NewContractIEigenDADirectoryCaller creates a new read-only instance of ContractIEigenDADirectory, bound to a specific deployed contract.
func NewContractIEigenDADirectoryCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDADirectoryCaller, error) {
	contract, err := bindContractIEigenDADirectory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryCaller{contract: contract}, nil
}

// NewContractIEigenDADirectoryTransactor creates a new write-only instance of ContractIEigenDADirectory, bound to a specific deployed contract.
func NewContractIEigenDADirectoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDADirectoryTransactor, error) {
	contract, err := bindContractIEigenDADirectory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryTransactor{contract: contract}, nil
}

// NewContractIEigenDADirectoryFilterer creates a new log filterer instance of ContractIEigenDADirectory, bound to a specific deployed contract.
func NewContractIEigenDADirectoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDADirectoryFilterer, error) {
	contract, err := bindContractIEigenDADirectory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryFilterer{contract: contract}, nil
}

// bindContractIEigenDADirectory binds a generic wrapper to an already deployed contract.
func bindContractIEigenDADirectory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDADirectoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDADirectory.Contract.ContractIEigenDADirectoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ContractIEigenDADirectoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ContractIEigenDADirectoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDADirectory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.contract.Transact(opts, method, params...)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetAddress(opts *bind.CallOpts, key [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetAddress(key [32]byte) (common.Address, error) {
	return _ContractIEigenDADirectory.Contract.GetAddress(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetAddress(key [32]byte) (common.Address, error) {
	return _ContractIEigenDADirectory.Contract.GetAddress(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetAddress0(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getAddress0", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetAddress0(name string) (common.Address, error) {
	return _ContractIEigenDADirectory.Contract.GetAddress0(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetAddress0(name string) (common.Address, error) {
	return _ContractIEigenDADirectory.Contract.GetAddress0(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetAllNames(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getAllNames")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetAllNames() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllNames(&_ContractIEigenDADirectory.CallOpts)
}

// GetAllNames is a free data retrieval call binding the contract method 0xfb825e5f.
//
// Solidity: function getAllNames() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetAllNames() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllNames(&_ContractIEigenDADirectory.CallOpts)
}

// GetConfigBytes is a free data retrieval call binding the contract method 0x4f93c9bf.
//
// Solidity: function getConfigBytes(string name) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes(opts *bind.CallOpts, name string) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes", name)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigBytes is a free data retrieval call binding the contract method 0x4f93c9bf.
//
// Solidity: function getConfigBytes(string name) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes(name string) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes is a free data retrieval call binding the contract method 0x4f93c9bf.
//
// Solidity: function getConfigBytes(string name) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes(name string) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes0 is a free data retrieval call binding the contract method 0x62c7855b.
//
// Solidity: function getConfigBytes(bytes32 key) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes0(opts *bind.CallOpts, key [32]byte) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes0", key)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigBytes0 is a free data retrieval call binding the contract method 0x62c7855b.
//
// Solidity: function getConfigBytes(bytes32 key) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes0(key [32]byte) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes0(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytes0 is a free data retrieval call binding the contract method 0x62c7855b.
//
// Solidity: function getConfigBytes(bytes32 key) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes0(key [32]byte) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes0(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0x6f8210b4.
//
// Solidity: function getConfigBytes32(bytes32 key) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes32(opts *bind.CallOpts, key [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes32", key)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0x6f8210b4.
//
// Solidity: function getConfigBytes32(bytes32 key) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes32(key [32]byte) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0x6f8210b4.
//
// Solidity: function getConfigBytes32(bytes32 key) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes32(key [32]byte) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytes320 is a free data retrieval call binding the contract method 0x99c2c9d3.
//
// Solidity: function getConfigBytes32(string name) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes320(opts *bind.CallOpts, name string) ([32]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes320", name)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConfigBytes320 is a free data retrieval call binding the contract method 0x99c2c9d3.
//
// Solidity: function getConfigBytes32(string name) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes320(name string) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes320(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes320 is a free data retrieval call binding the contract method 0x99c2c9d3.
//
// Solidity: function getConfigBytes32(string name) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes320(name string) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes320(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes32ExtraInfo is a free data retrieval call binding the contract method 0x2572b2ef.
//
// Solidity: function getConfigBytes32ExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes32ExtraInfo(opts *bind.CallOpts, name string) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes32ExtraInfo", name)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigBytes32ExtraInfo is a free data retrieval call binding the contract method 0x2572b2ef.
//
// Solidity: function getConfigBytes32ExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes32ExtraInfo(name string) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32ExtraInfo(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes32ExtraInfo is a free data retrieval call binding the contract method 0x2572b2ef.
//
// Solidity: function getConfigBytes32ExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes32ExtraInfo(name string) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32ExtraInfo(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytes32ExtraInfo0 is a free data retrieval call binding the contract method 0x5dc72cdb.
//
// Solidity: function getConfigBytes32ExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes32ExtraInfo0(opts *bind.CallOpts, key [32]byte) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes32ExtraInfo0", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigBytes32ExtraInfo0 is a free data retrieval call binding the contract method 0x5dc72cdb.
//
// Solidity: function getConfigBytes32ExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes32ExtraInfo0(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32ExtraInfo0(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytes32ExtraInfo0 is a free data retrieval call binding the contract method 0x5dc72cdb.
//
// Solidity: function getConfigBytes32ExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes32ExtraInfo0(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32ExtraInfo0(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytesExtraInfo is a free data retrieval call binding the contract method 0x3344b34a.
//
// Solidity: function getConfigBytesExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytesExtraInfo(opts *bind.CallOpts, key [32]byte) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytesExtraInfo", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigBytesExtraInfo is a free data retrieval call binding the contract method 0x3344b34a.
//
// Solidity: function getConfigBytesExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytesExtraInfo(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytesExtraInfo(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytesExtraInfo is a free data retrieval call binding the contract method 0x3344b34a.
//
// Solidity: function getConfigBytesExtraInfo(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytesExtraInfo(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytesExtraInfo(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetConfigBytesExtraInfo0 is a free data retrieval call binding the contract method 0x5cef4101.
//
// Solidity: function getConfigBytesExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytesExtraInfo0(opts *bind.CallOpts, name string) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytesExtraInfo0", name)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigBytesExtraInfo0 is a free data retrieval call binding the contract method 0x5cef4101.
//
// Solidity: function getConfigBytesExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytesExtraInfo0(name string) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytesExtraInfo0(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetConfigBytesExtraInfo0 is a free data retrieval call binding the contract method 0x5cef4101.
//
// Solidity: function getConfigBytesExtraInfo(string name) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytesExtraInfo0(name string) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytesExtraInfo0(&_ContractIEigenDADirectory.CallOpts, name)
}

// GetName is a free data retrieval call binding the contract method 0x54b8d5e3.
//
// Solidity: function getName(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetName(opts *bind.CallOpts, key [32]byte) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getName", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetName is a free data retrieval call binding the contract method 0x54b8d5e3.
//
// Solidity: function getName(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetName(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetName(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetName is a free data retrieval call binding the contract method 0x54b8d5e3.
//
// Solidity: function getName(bytes32 key) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetName(key [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetName(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetNumRegisteredKeysBytes is a free data retrieval call binding the contract method 0x6b8902e6.
//
// Solidity: function getNumRegisteredKeysBytes() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumRegisteredKeysBytes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumRegisteredKeysBytes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumRegisteredKeysBytes is a free data retrieval call binding the contract method 0x6b8902e6.
//
// Solidity: function getNumRegisteredKeysBytes() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumRegisteredKeysBytes() (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumRegisteredKeysBytes(&_ContractIEigenDADirectory.CallOpts)
}

// GetNumRegisteredKeysBytes is a free data retrieval call binding the contract method 0x6b8902e6.
//
// Solidity: function getNumRegisteredKeysBytes() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumRegisteredKeysBytes() (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumRegisteredKeysBytes(&_ContractIEigenDADirectory.CallOpts)
}

// GetNumRegisteredKeysBytes32 is a free data retrieval call binding the contract method 0xebcce4ab.
//
// Solidity: function getNumRegisteredKeysBytes32() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumRegisteredKeysBytes32(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumRegisteredKeysBytes32")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumRegisteredKeysBytes32 is a free data retrieval call binding the contract method 0xebcce4ab.
//
// Solidity: function getNumRegisteredKeysBytes32() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumRegisteredKeysBytes32() (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumRegisteredKeysBytes32(&_ContractIEigenDADirectory.CallOpts)
}

// GetNumRegisteredKeysBytes32 is a free data retrieval call binding the contract method 0xebcce4ab.
//
// Solidity: function getNumRegisteredKeysBytes32() view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumRegisteredKeysBytes32() (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumRegisteredKeysBytes32(&_ContractIEigenDADirectory.CallOpts)
}

// GetRegisteredKeyBytes is a free data retrieval call binding the contract method 0xb5df927b.
//
// Solidity: function getRegisteredKeyBytes(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetRegisteredKeyBytes(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getRegisteredKeyBytes", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRegisteredKeyBytes is a free data retrieval call binding the contract method 0xb5df927b.
//
// Solidity: function getRegisteredKeyBytes(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetRegisteredKeyBytes(index *big.Int) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetRegisteredKeyBytes(&_ContractIEigenDADirectory.CallOpts, index)
}

// GetRegisteredKeyBytes is a free data retrieval call binding the contract method 0xb5df927b.
//
// Solidity: function getRegisteredKeyBytes(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetRegisteredKeyBytes(index *big.Int) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetRegisteredKeyBytes(&_ContractIEigenDADirectory.CallOpts, index)
}

// GetRegisteredKeyBytes32 is a free data retrieval call binding the contract method 0x85109e2f.
//
// Solidity: function getRegisteredKeyBytes32(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetRegisteredKeyBytes32(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getRegisteredKeyBytes32", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRegisteredKeyBytes32 is a free data retrieval call binding the contract method 0x85109e2f.
//
// Solidity: function getRegisteredKeyBytes32(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetRegisteredKeyBytes32(index *big.Int) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetRegisteredKeyBytes32(&_ContractIEigenDADirectory.CallOpts, index)
}

// GetRegisteredKeyBytes32 is a free data retrieval call binding the contract method 0x85109e2f.
//
// Solidity: function getRegisteredKeyBytes32(uint256 index) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetRegisteredKeyBytes32(index *big.Int) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetRegisteredKeyBytes32(&_ContractIEigenDADirectory.CallOpts, index)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddAddress(opts *bind.TransactOpts, name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addAddress", name, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddAddress(&_ContractIEigenDADirectory.TransactOpts, name, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddAddress(&_ContractIEigenDADirectory.TransactOpts, name, value)
}

// AddConfigBytes is a paid mutator transaction binding the contract method 0xbc33a6de.
//
// Solidity: function addConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigBytes(opts *bind.TransactOpts, name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigBytes", name, value, extraInfo)
}

// AddConfigBytes is a paid mutator transaction binding the contract method 0xbc33a6de.
//
// Solidity: function addConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigBytes(name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// AddConfigBytes is a paid mutator transaction binding the contract method 0xbc33a6de.
//
// Solidity: function addConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigBytes(name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x456afa28.
//
// Solidity: function addConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigBytes32(opts *bind.TransactOpts, name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigBytes32", name, value, extraInfo)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x456afa28.
//
// Solidity: function addConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigBytes32(name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x456afa28.
//
// Solidity: function addConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigBytes32(name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) RemoveAddress(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "removeAddress", name)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) RemoveAddress(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveAddress(&_ContractIEigenDADirectory.TransactOpts, name)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) RemoveAddress(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveAddress(&_ContractIEigenDADirectory.TransactOpts, name)
}

// RemoveConfigBytes is a paid mutator transaction binding the contract method 0xc266b2b6.
//
// Solidity: function removeConfigBytes(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) RemoveConfigBytes(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "removeConfigBytes", name)
}

// RemoveConfigBytes is a paid mutator transaction binding the contract method 0xc266b2b6.
//
// Solidity: function removeConfigBytes(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) RemoveConfigBytes(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name)
}

// RemoveConfigBytes is a paid mutator transaction binding the contract method 0xc266b2b6.
//
// Solidity: function removeConfigBytes(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) RemoveConfigBytes(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name)
}

// RemoveConfigBytes32 is a paid mutator transaction binding the contract method 0x92579b17.
//
// Solidity: function removeConfigBytes32(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) RemoveConfigBytes32(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "removeConfigBytes32", name)
}

// RemoveConfigBytes32 is a paid mutator transaction binding the contract method 0x92579b17.
//
// Solidity: function removeConfigBytes32(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) RemoveConfigBytes32(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name)
}

// RemoveConfigBytes32 is a paid mutator transaction binding the contract method 0x92579b17.
//
// Solidity: function removeConfigBytes32(string name) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) RemoveConfigBytes32(name string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.RemoveConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) ReplaceAddress(opts *bind.TransactOpts, name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "replaceAddress", name, value)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) ReplaceAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceAddress(&_ContractIEigenDADirectory.TransactOpts, name, value)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) ReplaceAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceAddress(&_ContractIEigenDADirectory.TransactOpts, name, value)
}

// ReplaceConfigBytes is a paid mutator transaction binding the contract method 0xc506298c.
//
// Solidity: function replaceConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) ReplaceConfigBytes(opts *bind.TransactOpts, name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "replaceConfigBytes", name, value, extraInfo)
}

// ReplaceConfigBytes is a paid mutator transaction binding the contract method 0xc506298c.
//
// Solidity: function replaceConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) ReplaceConfigBytes(name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// ReplaceConfigBytes is a paid mutator transaction binding the contract method 0xc506298c.
//
// Solidity: function replaceConfigBytes(string name, bytes value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) ReplaceConfigBytes(name string, value []byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// ReplaceConfigBytes32 is a paid mutator transaction binding the contract method 0xdfe07f69.
//
// Solidity: function replaceConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) ReplaceConfigBytes32(opts *bind.TransactOpts, name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "replaceConfigBytes32", name, value, extraInfo)
}

// ReplaceConfigBytes32 is a paid mutator transaction binding the contract method 0xdfe07f69.
//
// Solidity: function replaceConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) ReplaceConfigBytes32(name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// ReplaceConfigBytes32 is a paid mutator transaction binding the contract method 0xdfe07f69.
//
// Solidity: function replaceConfigBytes32(string name, bytes32 value, string extraInfo) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) ReplaceConfigBytes32(name string, value [32]byte, extraInfo string) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.ReplaceConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, value, extraInfo)
}

// ContractIEigenDADirectoryAddressAddedIterator is returned from FilterAddressAdded and is used to iterate over the raw logs and unpacked data for AddressAdded events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressAddedIterator struct {
	Event *ContractIEigenDADirectoryAddressAdded // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryAddressAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryAddressAdded)
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
		it.Event = new(ContractIEigenDADirectoryAddressAdded)
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
func (it *ContractIEigenDADirectoryAddressAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryAddressAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryAddressAdded represents a AddressAdded event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressAdded struct {
	Name  string
	Key   [32]byte
	Value common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddressAdded is a free log retrieval operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterAddressAdded(opts *bind.FilterOpts, key [][32]byte, value []common.Address) (*ContractIEigenDADirectoryAddressAddedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "AddressAdded", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryAddressAddedIterator{contract: _ContractIEigenDADirectory.contract, event: "AddressAdded", logs: logs, sub: sub}, nil
}

// WatchAddressAdded is a free log subscription operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchAddressAdded(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryAddressAdded, key [][32]byte, value []common.Address) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "AddressAdded", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryAddressAdded)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressAdded", log); err != nil {
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

// ParseAddressAdded is a log parse operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseAddressAdded(log types.Log) (*ContractIEigenDADirectoryAddressAdded, error) {
	event := new(ContractIEigenDADirectoryAddressAdded)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryAddressRemovedIterator is returned from FilterAddressRemoved and is used to iterate over the raw logs and unpacked data for AddressRemoved events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressRemovedIterator struct {
	Event *ContractIEigenDADirectoryAddressRemoved // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryAddressRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryAddressRemoved)
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
		it.Event = new(ContractIEigenDADirectoryAddressRemoved)
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
func (it *ContractIEigenDADirectoryAddressRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryAddressRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryAddressRemoved represents a AddressRemoved event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressRemoved struct {
	Name string
	Key  [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddressRemoved is a free log retrieval operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterAddressRemoved(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryAddressRemovedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "AddressRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryAddressRemovedIterator{contract: _ContractIEigenDADirectory.contract, event: "AddressRemoved", logs: logs, sub: sub}, nil
}

// WatchAddressRemoved is a free log subscription operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchAddressRemoved(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryAddressRemoved, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "AddressRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryAddressRemoved)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressRemoved", log); err != nil {
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

// ParseAddressRemoved is a log parse operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseAddressRemoved(log types.Log) (*ContractIEigenDADirectoryAddressRemoved, error) {
	event := new(ContractIEigenDADirectoryAddressRemoved)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryAddressReplacedIterator is returned from FilterAddressReplaced and is used to iterate over the raw logs and unpacked data for AddressReplaced events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressReplacedIterator struct {
	Event *ContractIEigenDADirectoryAddressReplaced // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryAddressReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryAddressReplaced)
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
		it.Event = new(ContractIEigenDADirectoryAddressReplaced)
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
func (it *ContractIEigenDADirectoryAddressReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryAddressReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryAddressReplaced represents a AddressReplaced event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryAddressReplaced struct {
	Name     string
	Key      [32]byte
	OldValue common.Address
	NewValue common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddressReplaced is a free log retrieval operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterAddressReplaced(opts *bind.FilterOpts, key [][32]byte, oldValue []common.Address, newValue []common.Address) (*ContractIEigenDADirectoryAddressReplacedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var oldValueRule []interface{}
	for _, oldValueItem := range oldValue {
		oldValueRule = append(oldValueRule, oldValueItem)
	}
	var newValueRule []interface{}
	for _, newValueItem := range newValue {
		newValueRule = append(newValueRule, newValueItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "AddressReplaced", keyRule, oldValueRule, newValueRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryAddressReplacedIterator{contract: _ContractIEigenDADirectory.contract, event: "AddressReplaced", logs: logs, sub: sub}, nil
}

// WatchAddressReplaced is a free log subscription operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchAddressReplaced(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryAddressReplaced, key [][32]byte, oldValue []common.Address, newValue []common.Address) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var oldValueRule []interface{}
	for _, oldValueItem := range oldValue {
		oldValueRule = append(oldValueRule, oldValueItem)
	}
	var newValueRule []interface{}
	for _, newValueItem := range newValue {
		newValueRule = append(newValueRule, newValueItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "AddressReplaced", keyRule, oldValueRule, newValueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryAddressReplaced)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressReplaced", log); err != nil {
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

// ParseAddressReplaced is a log parse operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseAddressReplaced(log types.Log) (*ContractIEigenDADirectoryAddressReplaced, error) {
	event := new(ContractIEigenDADirectoryAddressReplaced)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "AddressReplaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytes32AddedIterator is returned from FilterConfigBytes32Added and is used to iterate over the raw logs and unpacked data for ConfigBytes32Added events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32AddedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytes32Added // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytes32AddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytes32Added)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytes32Added)
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
func (it *ContractIEigenDADirectoryConfigBytes32AddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytes32AddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytes32Added represents a ConfigBytes32Added event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32Added struct {
	Name      string
	Key       [32]byte
	Value     [32]byte
	ExtraInfo string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfigBytes32Added is a free log retrieval operation binding the contract event 0x8ca30810a87a9829b9eca24ecb258e4008a32bea8a04e2976dba7b8733fdc258.
//
// Solidity: event ConfigBytes32Added(string name, bytes32 indexed key, bytes32 value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytes32Added(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytes32AddedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytes32Added", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytes32AddedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytes32Added", logs: logs, sub: sub}, nil
}

// WatchConfigBytes32Added is a free log subscription operation binding the contract event 0x8ca30810a87a9829b9eca24ecb258e4008a32bea8a04e2976dba7b8733fdc258.
//
// Solidity: event ConfigBytes32Added(string name, bytes32 indexed key, bytes32 value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytes32Added(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytes32Added, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytes32Added", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytes32Added)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Added", log); err != nil {
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

// ParseConfigBytes32Added is a log parse operation binding the contract event 0x8ca30810a87a9829b9eca24ecb258e4008a32bea8a04e2976dba7b8733fdc258.
//
// Solidity: event ConfigBytes32Added(string name, bytes32 indexed key, bytes32 value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytes32Added(log types.Log) (*ContractIEigenDADirectoryConfigBytes32Added, error) {
	event := new(ContractIEigenDADirectoryConfigBytes32Added)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Added", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytes32RemovedIterator is returned from FilterConfigBytes32Removed and is used to iterate over the raw logs and unpacked data for ConfigBytes32Removed events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32RemovedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytes32Removed // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytes32RemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytes32Removed)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytes32Removed)
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
func (it *ContractIEigenDADirectoryConfigBytes32RemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytes32RemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytes32Removed represents a ConfigBytes32Removed event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32Removed struct {
	Name string
	Key  [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterConfigBytes32Removed is a free log retrieval operation binding the contract event 0x377bcf9a98ffa92fa8e7a9043365e65e6c8d1b0923041d33b3e58c22ab6c2d6d.
//
// Solidity: event ConfigBytes32Removed(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytes32Removed(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytes32RemovedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytes32Removed", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytes32RemovedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytes32Removed", logs: logs, sub: sub}, nil
}

// WatchConfigBytes32Removed is a free log subscription operation binding the contract event 0x377bcf9a98ffa92fa8e7a9043365e65e6c8d1b0923041d33b3e58c22ab6c2d6d.
//
// Solidity: event ConfigBytes32Removed(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytes32Removed(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytes32Removed, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytes32Removed", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytes32Removed)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Removed", log); err != nil {
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

// ParseConfigBytes32Removed is a log parse operation binding the contract event 0x377bcf9a98ffa92fa8e7a9043365e65e6c8d1b0923041d33b3e58c22ab6c2d6d.
//
// Solidity: event ConfigBytes32Removed(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytes32Removed(log types.Log) (*ContractIEigenDADirectoryConfigBytes32Removed, error) {
	event := new(ContractIEigenDADirectoryConfigBytes32Removed)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Removed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytes32ReplacedIterator is returned from FilterConfigBytes32Replaced and is used to iterate over the raw logs and unpacked data for ConfigBytes32Replaced events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32ReplacedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytes32Replaced // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytes32ReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytes32Replaced)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytes32Replaced)
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
func (it *ContractIEigenDADirectoryConfigBytes32ReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytes32ReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytes32Replaced represents a ConfigBytes32Replaced event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytes32Replaced struct {
	Name      string
	Key       [32]byte
	OldValue  [32]byte
	NewValue  [32]byte
	ExtraInfo string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfigBytes32Replaced is a free log retrieval operation binding the contract event 0x3f6edb55367576b200c679337c1833d0eea24053ed16620e22f1124885321920.
//
// Solidity: event ConfigBytes32Replaced(string name, bytes32 indexed key, bytes32 oldValue, bytes32 newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytes32Replaced(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytes32ReplacedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytes32Replaced", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytes32ReplacedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytes32Replaced", logs: logs, sub: sub}, nil
}

// WatchConfigBytes32Replaced is a free log subscription operation binding the contract event 0x3f6edb55367576b200c679337c1833d0eea24053ed16620e22f1124885321920.
//
// Solidity: event ConfigBytes32Replaced(string name, bytes32 indexed key, bytes32 oldValue, bytes32 newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytes32Replaced(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytes32Replaced, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytes32Replaced", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytes32Replaced)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Replaced", log); err != nil {
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

// ParseConfigBytes32Replaced is a log parse operation binding the contract event 0x3f6edb55367576b200c679337c1833d0eea24053ed16620e22f1124885321920.
//
// Solidity: event ConfigBytes32Replaced(string name, bytes32 indexed key, bytes32 oldValue, bytes32 newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytes32Replaced(log types.Log) (*ContractIEigenDADirectoryConfigBytes32Replaced, error) {
	event := new(ContractIEigenDADirectoryConfigBytes32Replaced)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytes32Replaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytesAddedIterator is returned from FilterConfigBytesAdded and is used to iterate over the raw logs and unpacked data for ConfigBytesAdded events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesAddedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytesAdded // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytesAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytesAdded)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytesAdded)
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
func (it *ContractIEigenDADirectoryConfigBytesAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytesAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytesAdded represents a ConfigBytesAdded event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesAdded struct {
	Name      string
	Key       [32]byte
	Value     []byte
	ExtraInfo string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfigBytesAdded is a free log retrieval operation binding the contract event 0x2a82356b15556239d600bd78e622c325f7afe7b56e52c2deedf7170c22e23f1f.
//
// Solidity: event ConfigBytesAdded(string name, bytes32 indexed key, bytes value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytesAdded(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytesAddedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytesAdded", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytesAddedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytesAdded", logs: logs, sub: sub}, nil
}

// WatchConfigBytesAdded is a free log subscription operation binding the contract event 0x2a82356b15556239d600bd78e622c325f7afe7b56e52c2deedf7170c22e23f1f.
//
// Solidity: event ConfigBytesAdded(string name, bytes32 indexed key, bytes value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytesAdded(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytesAdded, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytesAdded", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytesAdded)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesAdded", log); err != nil {
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

// ParseConfigBytesAdded is a log parse operation binding the contract event 0x2a82356b15556239d600bd78e622c325f7afe7b56e52c2deedf7170c22e23f1f.
//
// Solidity: event ConfigBytesAdded(string name, bytes32 indexed key, bytes value, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytesAdded(log types.Log) (*ContractIEigenDADirectoryConfigBytesAdded, error) {
	event := new(ContractIEigenDADirectoryConfigBytesAdded)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytesRemovedIterator is returned from FilterConfigBytesRemoved and is used to iterate over the raw logs and unpacked data for ConfigBytesRemoved events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesRemovedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytesRemoved // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytesRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytesRemoved)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytesRemoved)
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
func (it *ContractIEigenDADirectoryConfigBytesRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytesRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytesRemoved represents a ConfigBytesRemoved event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesRemoved struct {
	Name string
	Key  [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterConfigBytesRemoved is a free log retrieval operation binding the contract event 0x4eff0809880a1b2fa6fbf6a3d6800c3256abe038317258c730ef36c8e13e6081.
//
// Solidity: event ConfigBytesRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytesRemoved(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytesRemovedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytesRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytesRemovedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytesRemoved", logs: logs, sub: sub}, nil
}

// WatchConfigBytesRemoved is a free log subscription operation binding the contract event 0x4eff0809880a1b2fa6fbf6a3d6800c3256abe038317258c730ef36c8e13e6081.
//
// Solidity: event ConfigBytesRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytesRemoved(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytesRemoved, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytesRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytesRemoved)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesRemoved", log); err != nil {
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

// ParseConfigBytesRemoved is a log parse operation binding the contract event 0x4eff0809880a1b2fa6fbf6a3d6800c3256abe038317258c730ef36c8e13e6081.
//
// Solidity: event ConfigBytesRemoved(string name, bytes32 indexed key)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytesRemoved(log types.Log) (*ContractIEigenDADirectoryConfigBytesRemoved, error) {
	event := new(ContractIEigenDADirectoryConfigBytesRemoved)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDADirectoryConfigBytesReplacedIterator is returned from FilterConfigBytesReplaced and is used to iterate over the raw logs and unpacked data for ConfigBytesReplaced events raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesReplacedIterator struct {
	Event *ContractIEigenDADirectoryConfigBytesReplaced // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDADirectoryConfigBytesReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDADirectoryConfigBytesReplaced)
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
		it.Event = new(ContractIEigenDADirectoryConfigBytesReplaced)
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
func (it *ContractIEigenDADirectoryConfigBytesReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDADirectoryConfigBytesReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDADirectoryConfigBytesReplaced represents a ConfigBytesReplaced event raised by the ContractIEigenDADirectory contract.
type ContractIEigenDADirectoryConfigBytesReplaced struct {
	Name      string
	Key       [32]byte
	OldValue  []byte
	NewValue  []byte
	ExtraInfo string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfigBytesReplaced is a free log retrieval operation binding the contract event 0xfec871e3d17048fa8a49fb9d6e13a5998b704b6b39b815c3d8172066acf11d8d.
//
// Solidity: event ConfigBytesReplaced(string name, bytes32 indexed key, bytes oldValue, bytes newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) FilterConfigBytesReplaced(opts *bind.FilterOpts, key [][32]byte) (*ContractIEigenDADirectoryConfigBytesReplacedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.FilterLogs(opts, "ConfigBytesReplaced", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDADirectoryConfigBytesReplacedIterator{contract: _ContractIEigenDADirectory.contract, event: "ConfigBytesReplaced", logs: logs, sub: sub}, nil
}

// WatchConfigBytesReplaced is a free log subscription operation binding the contract event 0xfec871e3d17048fa8a49fb9d6e13a5998b704b6b39b815c3d8172066acf11d8d.
//
// Solidity: event ConfigBytesReplaced(string name, bytes32 indexed key, bytes oldValue, bytes newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) WatchConfigBytesReplaced(opts *bind.WatchOpts, sink chan<- *ContractIEigenDADirectoryConfigBytesReplaced, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDADirectory.contract.WatchLogs(opts, "ConfigBytesReplaced", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDADirectoryConfigBytesReplaced)
				if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesReplaced", log); err != nil {
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

// ParseConfigBytesReplaced is a log parse operation binding the contract event 0xfec871e3d17048fa8a49fb9d6e13a5998b704b6b39b815c3d8172066acf11d8d.
//
// Solidity: event ConfigBytesReplaced(string name, bytes32 indexed key, bytes oldValue, bytes newValue, string extraInfo)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryFilterer) ParseConfigBytesReplaced(log types.Log) (*ContractIEigenDADirectoryConfigBytesReplaced, error) {
	event := new(ContractIEigenDADirectoryConfigBytesReplaced)
	if err := _ContractIEigenDADirectory.contract.UnpackLog(event, "ConfigBytesReplaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
