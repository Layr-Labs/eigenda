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

// ConfigRegistryTypesBytes32Checkpoint is an auto generated low-level Go binding around an user-defined struct.
type ConfigRegistryTypesBytes32Checkpoint struct {
	ActivationKey *big.Int
	Value         [32]byte
}

// ConfigRegistryTypesBytesCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ConfigRegistryTypesBytesCheckpoint struct {
	ActivationKey *big.Int
	Value         []byte
}

// ContractIEigenDADirectoryMetaData contains all meta data concerning the ContractIEigenDADirectory contract.
var ContractIEigenDADirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigBytes\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"activationKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigBytes32\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"activationKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActivationKeyBytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActivationKeyBytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNames\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCheckpointBytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structConfigRegistryTypes.BytesCheckpoint\",\"components\":[{\"name\":\"activationKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCheckpointBytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structConfigRegistryTypes.Bytes32Checkpoint\",\"components\":[{\"name\":\"activationKey\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getName\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumCheckpointsBytes\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumCheckpointsBytes32\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddressAdded\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressReplaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"AddressDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConfigAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ConfigDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NewValueIsOldValue\",\"inputs\":[{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
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

// GetActivationKeyBytes is a free data retrieval call binding the contract method 0x19f73b7f.
//
// Solidity: function getActivationKeyBytes(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetActivationKeyBytes(opts *bind.CallOpts, key [32]byte, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getActivationKeyBytes", key, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActivationKeyBytes is a free data retrieval call binding the contract method 0x19f73b7f.
//
// Solidity: function getActivationKeyBytes(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetActivationKeyBytes(key [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationKeyBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetActivationKeyBytes is a free data retrieval call binding the contract method 0x19f73b7f.
//
// Solidity: function getActivationKeyBytes(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetActivationKeyBytes(key [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationKeyBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetActivationKeyBytes32 is a free data retrieval call binding the contract method 0x8d033b96.
//
// Solidity: function getActivationKeyBytes32(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetActivationKeyBytes32(opts *bind.CallOpts, key [32]byte, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getActivationKeyBytes32", key, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActivationKeyBytes32 is a free data retrieval call binding the contract method 0x8d033b96.
//
// Solidity: function getActivationKeyBytes32(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetActivationKeyBytes32(key [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationKeyBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetActivationKeyBytes32 is a free data retrieval call binding the contract method 0x8d033b96.
//
// Solidity: function getActivationKeyBytes32(bytes32 key, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetActivationKeyBytes32(key [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationKeyBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
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

// GetCheckpointBytes is a free data retrieval call binding the contract method 0x67134a1d.
//
// Solidity: function getCheckpointBytes(bytes32 key, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetCheckpointBytes(opts *bind.CallOpts, key [32]byte, index *big.Int) (ConfigRegistryTypesBytesCheckpoint, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getCheckpointBytes", key, index)

	if err != nil {
		return *new(ConfigRegistryTypesBytesCheckpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(ConfigRegistryTypesBytesCheckpoint)).(*ConfigRegistryTypesBytesCheckpoint)

	return out0, err

}

// GetCheckpointBytes is a free data retrieval call binding the contract method 0x67134a1d.
//
// Solidity: function getCheckpointBytes(bytes32 key, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetCheckpointBytes(key [32]byte, index *big.Int) (ConfigRegistryTypesBytesCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetCheckpointBytes is a free data retrieval call binding the contract method 0x67134a1d.
//
// Solidity: function getCheckpointBytes(bytes32 key, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetCheckpointBytes(key [32]byte, index *big.Int) (ConfigRegistryTypesBytesCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetCheckpointBytes32 is a free data retrieval call binding the contract method 0xf1e98690.
//
// Solidity: function getCheckpointBytes32(bytes32 key, uint256 index) view returns((uint256,bytes32))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetCheckpointBytes32(opts *bind.CallOpts, key [32]byte, index *big.Int) (ConfigRegistryTypesBytes32Checkpoint, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getCheckpointBytes32", key, index)

	if err != nil {
		return *new(ConfigRegistryTypesBytes32Checkpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(ConfigRegistryTypesBytes32Checkpoint)).(*ConfigRegistryTypesBytes32Checkpoint)

	return out0, err

}

// GetCheckpointBytes32 is a free data retrieval call binding the contract method 0xf1e98690.
//
// Solidity: function getCheckpointBytes32(bytes32 key, uint256 index) view returns((uint256,bytes32))
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetCheckpointBytes32(key [32]byte, index *big.Int) (ConfigRegistryTypesBytes32Checkpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetCheckpointBytes32 is a free data retrieval call binding the contract method 0xf1e98690.
//
// Solidity: function getCheckpointBytes32(bytes32 key, uint256 index) view returns((uint256,bytes32))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetCheckpointBytes32(key [32]byte, index *big.Int) (ConfigRegistryTypesBytes32Checkpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetConfigBytes is a free data retrieval call binding the contract method 0x5e28bafc.
//
// Solidity: function getConfigBytes(bytes32 key, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes(opts *bind.CallOpts, key [32]byte, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes", key, index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigBytes is a free data retrieval call binding the contract method 0x5e28bafc.
//
// Solidity: function getConfigBytes(bytes32 key, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes(key [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetConfigBytes is a free data retrieval call binding the contract method 0x5e28bafc.
//
// Solidity: function getConfigBytes(bytes32 key, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes(key [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0xb8587056.
//
// Solidity: function getConfigBytes32(bytes32 key, uint256 index) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBytes32(opts *bind.CallOpts, key [32]byte, index *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBytes32", key, index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0xb8587056.
//
// Solidity: function getConfigBytes32(bytes32 key, uint256 index) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBytes32(key [32]byte, index *big.Int) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
}

// GetConfigBytes32 is a free data retrieval call binding the contract method 0xb8587056.
//
// Solidity: function getConfigBytes32(bytes32 key, uint256 index) view returns(bytes32)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBytes32(key [32]byte, index *big.Int) ([32]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBytes32(&_ContractIEigenDADirectory.CallOpts, key, index)
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

// GetNumCheckpointsBytes is a free data retrieval call binding the contract method 0x195e8d6e.
//
// Solidity: function getNumCheckpointsBytes(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumCheckpointsBytes(opts *bind.CallOpts, key [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumCheckpointsBytes", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumCheckpointsBytes is a free data retrieval call binding the contract method 0x195e8d6e.
//
// Solidity: function getNumCheckpointsBytes(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumCheckpointsBytes(key [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBytes(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetNumCheckpointsBytes is a free data retrieval call binding the contract method 0x195e8d6e.
//
// Solidity: function getNumCheckpointsBytes(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumCheckpointsBytes(key [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBytes(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetNumCheckpointsBytes32 is a free data retrieval call binding the contract method 0x27b092a6.
//
// Solidity: function getNumCheckpointsBytes32(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumCheckpointsBytes32(opts *bind.CallOpts, key [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumCheckpointsBytes32", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumCheckpointsBytes32 is a free data retrieval call binding the contract method 0x27b092a6.
//
// Solidity: function getNumCheckpointsBytes32(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumCheckpointsBytes32(key [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBytes32(&_ContractIEigenDADirectory.CallOpts, key)
}

// GetNumCheckpointsBytes32 is a free data retrieval call binding the contract method 0x27b092a6.
//
// Solidity: function getNumCheckpointsBytes32(bytes32 key) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumCheckpointsBytes32(key [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBytes32(&_ContractIEigenDADirectory.CallOpts, key)
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

// AddConfigBytes is a paid mutator transaction binding the contract method 0x515a49b3.
//
// Solidity: function addConfigBytes(string name, uint256 activationKey, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigBytes(opts *bind.TransactOpts, name string, activationKey *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigBytes", name, activationKey, value)
}

// AddConfigBytes is a paid mutator transaction binding the contract method 0x515a49b3.
//
// Solidity: function addConfigBytes(string name, uint256 activationKey, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigBytes(name string, activationKey *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, activationKey, value)
}

// AddConfigBytes is a paid mutator transaction binding the contract method 0x515a49b3.
//
// Solidity: function addConfigBytes(string name, uint256 activationKey, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigBytes(name string, activationKey *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes(&_ContractIEigenDADirectory.TransactOpts, name, activationKey, value)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x99febdd8.
//
// Solidity: function addConfigBytes32(string name, uint256 activationKey, bytes32 value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigBytes32(opts *bind.TransactOpts, name string, activationKey *big.Int, value [32]byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigBytes32", name, activationKey, value)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x99febdd8.
//
// Solidity: function addConfigBytes32(string name, uint256 activationKey, bytes32 value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigBytes32(name string, activationKey *big.Int, value [32]byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, activationKey, value)
}

// AddConfigBytes32 is a paid mutator transaction binding the contract method 0x99febdd8.
//
// Solidity: function addConfigBytes32(string name, uint256 activationKey, bytes32 value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigBytes32(name string, activationKey *big.Int, value [32]byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBytes32(&_ContractIEigenDADirectory.TransactOpts, name, activationKey, value)
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
