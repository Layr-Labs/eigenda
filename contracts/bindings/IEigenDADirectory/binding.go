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

// ConfigRegistryTypesBlockNumberCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ConfigRegistryTypesBlockNumberCheckpoint struct {
	ActivationBlock *big.Int
	Value           []byte
}

// ConfigRegistryTypesTimeStampCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ConfigRegistryTypesTimeStampCheckpoint struct {
	ActivationTime *big.Int
	Value          []byte
}

// ContractIEigenDADirectoryMetaData contains all meta data concerning the ContractIEigenDADirectory contract.
var ContractIEigenDADirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigBlockNumber\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"abn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addConfigTimeStamp\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"activationTS\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActivationBlockNumber\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActivationTimeStamp\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllConfigNamesBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllConfigNamesTimeStamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNames\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCheckpointBlockNumber\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structConfigRegistryTypes.BlockNumberCheckpoint\",\"components\":[{\"name\":\"activationBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCheckpointTimeStamp\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structConfigRegistryTypes.TimeStampCheckpoint\",\"components\":[{\"name\":\"activationTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigBlockNumber\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigNameBlockNumber\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigNameTimeStamp\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigTimeStamp\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getName\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumCheckpointsBlockNumber\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumCheckpointsTimeStamp\",\"inputs\":[{\"name\":\"nameDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddressAdded\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressReplaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"AddressDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NewValueIsOldValue\",\"inputs\":[{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
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

// GetActivationBlockNumber is a free data retrieval call binding the contract method 0xa78735a2.
//
// Solidity: function getActivationBlockNumber(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetActivationBlockNumber(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getActivationBlockNumber", nameDigest, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActivationBlockNumber is a free data retrieval call binding the contract method 0xa78735a2.
//
// Solidity: function getActivationBlockNumber(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetActivationBlockNumber(nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetActivationBlockNumber is a free data retrieval call binding the contract method 0xa78735a2.
//
// Solidity: function getActivationBlockNumber(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetActivationBlockNumber(nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetActivationTimeStamp is a free data retrieval call binding the contract method 0x16e34391.
//
// Solidity: function getActivationTimeStamp(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetActivationTimeStamp(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getActivationTimeStamp", nameDigest, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActivationTimeStamp is a free data retrieval call binding the contract method 0x16e34391.
//
// Solidity: function getActivationTimeStamp(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetActivationTimeStamp(nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetActivationTimeStamp is a free data retrieval call binding the contract method 0x16e34391.
//
// Solidity: function getActivationTimeStamp(bytes32 nameDigest, uint256 index) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetActivationTimeStamp(nameDigest [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetActivationTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
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

// GetAllConfigNamesBlockNumber is a free data retrieval call binding the contract method 0xda1a8a0a.
//
// Solidity: function getAllConfigNamesBlockNumber() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetAllConfigNamesBlockNumber(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getAllConfigNamesBlockNumber")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAllConfigNamesBlockNumber is a free data retrieval call binding the contract method 0xda1a8a0a.
//
// Solidity: function getAllConfigNamesBlockNumber() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetAllConfigNamesBlockNumber() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllConfigNamesBlockNumber(&_ContractIEigenDADirectory.CallOpts)
}

// GetAllConfigNamesBlockNumber is a free data retrieval call binding the contract method 0xda1a8a0a.
//
// Solidity: function getAllConfigNamesBlockNumber() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetAllConfigNamesBlockNumber() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllConfigNamesBlockNumber(&_ContractIEigenDADirectory.CallOpts)
}

// GetAllConfigNamesTimeStamp is a free data retrieval call binding the contract method 0x4420027c.
//
// Solidity: function getAllConfigNamesTimeStamp() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetAllConfigNamesTimeStamp(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getAllConfigNamesTimeStamp")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAllConfigNamesTimeStamp is a free data retrieval call binding the contract method 0x4420027c.
//
// Solidity: function getAllConfigNamesTimeStamp() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetAllConfigNamesTimeStamp() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllConfigNamesTimeStamp(&_ContractIEigenDADirectory.CallOpts)
}

// GetAllConfigNamesTimeStamp is a free data retrieval call binding the contract method 0x4420027c.
//
// Solidity: function getAllConfigNamesTimeStamp() view returns(string[])
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetAllConfigNamesTimeStamp() ([]string, error) {
	return _ContractIEigenDADirectory.Contract.GetAllConfigNamesTimeStamp(&_ContractIEigenDADirectory.CallOpts)
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

// GetCheckpointBlockNumber is a free data retrieval call binding the contract method 0x723e08c8.
//
// Solidity: function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetCheckpointBlockNumber(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesBlockNumberCheckpoint, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getCheckpointBlockNumber", nameDigest, index)

	if err != nil {
		return *new(ConfigRegistryTypesBlockNumberCheckpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(ConfigRegistryTypesBlockNumberCheckpoint)).(*ConfigRegistryTypesBlockNumberCheckpoint)

	return out0, err

}

// GetCheckpointBlockNumber is a free data retrieval call binding the contract method 0x723e08c8.
//
// Solidity: function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetCheckpointBlockNumber(nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesBlockNumberCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetCheckpointBlockNumber is a free data retrieval call binding the contract method 0x723e08c8.
//
// Solidity: function getCheckpointBlockNumber(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetCheckpointBlockNumber(nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesBlockNumberCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetCheckpointTimeStamp is a free data retrieval call binding the contract method 0xc4fd4234.
//
// Solidity: function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetCheckpointTimeStamp(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesTimeStampCheckpoint, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getCheckpointTimeStamp", nameDigest, index)

	if err != nil {
		return *new(ConfigRegistryTypesTimeStampCheckpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(ConfigRegistryTypesTimeStampCheckpoint)).(*ConfigRegistryTypesTimeStampCheckpoint)

	return out0, err

}

// GetCheckpointTimeStamp is a free data retrieval call binding the contract method 0xc4fd4234.
//
// Solidity: function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetCheckpointTimeStamp(nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesTimeStampCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetCheckpointTimeStamp is a free data retrieval call binding the contract method 0xc4fd4234.
//
// Solidity: function getCheckpointTimeStamp(bytes32 nameDigest, uint256 index) view returns((uint256,bytes))
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetCheckpointTimeStamp(nameDigest [32]byte, index *big.Int) (ConfigRegistryTypesTimeStampCheckpoint, error) {
	return _ContractIEigenDADirectory.Contract.GetCheckpointTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetConfigBlockNumber is a free data retrieval call binding the contract method 0xf4a56be3.
//
// Solidity: function getConfigBlockNumber(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigBlockNumber(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigBlockNumber", nameDigest, index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigBlockNumber is a free data retrieval call binding the contract method 0xf4a56be3.
//
// Solidity: function getConfigBlockNumber(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigBlockNumber(nameDigest [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetConfigBlockNumber is a free data retrieval call binding the contract method 0xf4a56be3.
//
// Solidity: function getConfigBlockNumber(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigBlockNumber(nameDigest [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetConfigNameBlockNumber is a free data retrieval call binding the contract method 0xb0465b5f.
//
// Solidity: function getConfigNameBlockNumber(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigNameBlockNumber(opts *bind.CallOpts, nameDigest [32]byte) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigNameBlockNumber", nameDigest)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigNameBlockNumber is a free data retrieval call binding the contract method 0xb0465b5f.
//
// Solidity: function getConfigNameBlockNumber(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigNameBlockNumber(nameDigest [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigNameBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetConfigNameBlockNumber is a free data retrieval call binding the contract method 0xb0465b5f.
//
// Solidity: function getConfigNameBlockNumber(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigNameBlockNumber(nameDigest [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigNameBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetConfigNameTimeStamp is a free data retrieval call binding the contract method 0xe2c53d48.
//
// Solidity: function getConfigNameTimeStamp(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigNameTimeStamp(opts *bind.CallOpts, nameDigest [32]byte) (string, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigNameTimeStamp", nameDigest)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetConfigNameTimeStamp is a free data retrieval call binding the contract method 0xe2c53d48.
//
// Solidity: function getConfigNameTimeStamp(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigNameTimeStamp(nameDigest [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigNameTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetConfigNameTimeStamp is a free data retrieval call binding the contract method 0xe2c53d48.
//
// Solidity: function getConfigNameTimeStamp(bytes32 nameDigest) view returns(string)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigNameTimeStamp(nameDigest [32]byte) (string, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigNameTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetConfigTimeStamp is a free data retrieval call binding the contract method 0xd8e62afb.
//
// Solidity: function getConfigTimeStamp(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetConfigTimeStamp(opts *bind.CallOpts, nameDigest [32]byte, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getConfigTimeStamp", nameDigest, index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetConfigTimeStamp is a free data retrieval call binding the contract method 0xd8e62afb.
//
// Solidity: function getConfigTimeStamp(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetConfigTimeStamp(nameDigest [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
}

// GetConfigTimeStamp is a free data retrieval call binding the contract method 0xd8e62afb.
//
// Solidity: function getConfigTimeStamp(bytes32 nameDigest, uint256 index) view returns(bytes)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetConfigTimeStamp(nameDigest [32]byte, index *big.Int) ([]byte, error) {
	return _ContractIEigenDADirectory.Contract.GetConfigTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest, index)
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

// GetNumCheckpointsBlockNumber is a free data retrieval call binding the contract method 0xac1cc0c0.
//
// Solidity: function getNumCheckpointsBlockNumber(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumCheckpointsBlockNumber(opts *bind.CallOpts, nameDigest [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumCheckpointsBlockNumber", nameDigest)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumCheckpointsBlockNumber is a free data retrieval call binding the contract method 0xac1cc0c0.
//
// Solidity: function getNumCheckpointsBlockNumber(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumCheckpointsBlockNumber(nameDigest [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetNumCheckpointsBlockNumber is a free data retrieval call binding the contract method 0xac1cc0c0.
//
// Solidity: function getNumCheckpointsBlockNumber(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumCheckpointsBlockNumber(nameDigest [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsBlockNumber(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetNumCheckpointsTimeStamp is a free data retrieval call binding the contract method 0x69393318.
//
// Solidity: function getNumCheckpointsTimeStamp(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCaller) GetNumCheckpointsTimeStamp(opts *bind.CallOpts, nameDigest [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDADirectory.contract.Call(opts, &out, "getNumCheckpointsTimeStamp", nameDigest)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumCheckpointsTimeStamp is a free data retrieval call binding the contract method 0x69393318.
//
// Solidity: function getNumCheckpointsTimeStamp(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) GetNumCheckpointsTimeStamp(nameDigest [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest)
}

// GetNumCheckpointsTimeStamp is a free data retrieval call binding the contract method 0x69393318.
//
// Solidity: function getNumCheckpointsTimeStamp(bytes32 nameDigest) view returns(uint256)
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryCallerSession) GetNumCheckpointsTimeStamp(nameDigest [32]byte) (*big.Int, error) {
	return _ContractIEigenDADirectory.Contract.GetNumCheckpointsTimeStamp(&_ContractIEigenDADirectory.CallOpts, nameDigest)
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

// AddConfigBlockNumber is a paid mutator transaction binding the contract method 0x3a45bc4f.
//
// Solidity: function addConfigBlockNumber(string name, uint256 abn, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigBlockNumber(opts *bind.TransactOpts, name string, abn *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigBlockNumber", name, abn, value)
}

// AddConfigBlockNumber is a paid mutator transaction binding the contract method 0x3a45bc4f.
//
// Solidity: function addConfigBlockNumber(string name, uint256 abn, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigBlockNumber(name string, abn *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBlockNumber(&_ContractIEigenDADirectory.TransactOpts, name, abn, value)
}

// AddConfigBlockNumber is a paid mutator transaction binding the contract method 0x3a45bc4f.
//
// Solidity: function addConfigBlockNumber(string name, uint256 abn, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigBlockNumber(name string, abn *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigBlockNumber(&_ContractIEigenDADirectory.TransactOpts, name, abn, value)
}

// AddConfigTimeStamp is a paid mutator transaction binding the contract method 0xa2e91eb9.
//
// Solidity: function addConfigTimeStamp(string name, uint256 activationTS, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactor) AddConfigTimeStamp(opts *bind.TransactOpts, name string, activationTS *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.contract.Transact(opts, "addConfigTimeStamp", name, activationTS, value)
}

// AddConfigTimeStamp is a paid mutator transaction binding the contract method 0xa2e91eb9.
//
// Solidity: function addConfigTimeStamp(string name, uint256 activationTS, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectorySession) AddConfigTimeStamp(name string, activationTS *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigTimeStamp(&_ContractIEigenDADirectory.TransactOpts, name, activationTS, value)
}

// AddConfigTimeStamp is a paid mutator transaction binding the contract method 0xa2e91eb9.
//
// Solidity: function addConfigTimeStamp(string name, uint256 activationTS, bytes value) returns()
func (_ContractIEigenDADirectory *ContractIEigenDADirectoryTransactorSession) AddConfigTimeStamp(name string, activationTS *big.Int, value []byte) (*types.Transaction, error) {
	return _ContractIEigenDADirectory.Contract.AddConfigTimeStamp(&_ContractIEigenDADirectory.TransactOpts, name, activationTS, value)
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
