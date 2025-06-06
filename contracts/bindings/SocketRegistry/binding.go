// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractSocketRegistry

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

// ContractSocketRegistryMetaData contains all meta data concerning the ContractSocketRegistry contract.
var ContractSocketRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorIdToSocket\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setOperatorSocket\",\"inputs\":[{\"name\":\"_operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_socket\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// ContractSocketRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractSocketRegistryMetaData.ABI instead.
var ContractSocketRegistryABI = ContractSocketRegistryMetaData.ABI

// ContractSocketRegistry is an auto generated Go binding around an Ethereum contract.
type ContractSocketRegistry struct {
	ContractSocketRegistryCaller     // Read-only binding to the contract
	ContractSocketRegistryTransactor // Write-only binding to the contract
	ContractSocketRegistryFilterer   // Log filterer for contract events
}

// ContractSocketRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractSocketRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractSocketRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractSocketRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSocketRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSocketRegistrySession struct {
	Contract     *ContractSocketRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ContractSocketRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractSocketRegistryCallerSession struct {
	Contract *ContractSocketRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ContractSocketRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractSocketRegistryTransactorSession struct {
	Contract     *ContractSocketRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractSocketRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractSocketRegistryRaw struct {
	Contract *ContractSocketRegistry // Generic contract binding to access the raw methods on
}

// ContractSocketRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractSocketRegistryCallerRaw struct {
	Contract *ContractSocketRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractSocketRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractSocketRegistryTransactorRaw struct {
	Contract *ContractSocketRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractSocketRegistry creates a new instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistry(address common.Address, backend bind.ContractBackend) (*ContractSocketRegistry, error) {
	contract, err := bindContractSocketRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistry{ContractSocketRegistryCaller: ContractSocketRegistryCaller{contract: contract}, ContractSocketRegistryTransactor: ContractSocketRegistryTransactor{contract: contract}, ContractSocketRegistryFilterer: ContractSocketRegistryFilterer{contract: contract}}, nil
}

// NewContractSocketRegistryCaller creates a new read-only instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractSocketRegistryCaller, error) {
	contract, err := bindContractSocketRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryCaller{contract: contract}, nil
}

// NewContractSocketRegistryTransactor creates a new write-only instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractSocketRegistryTransactor, error) {
	contract, err := bindContractSocketRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryTransactor{contract: contract}, nil
}

// NewContractSocketRegistryFilterer creates a new log filterer instance of ContractSocketRegistry, bound to a specific deployed contract.
func NewContractSocketRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractSocketRegistryFilterer, error) {
	contract, err := bindContractSocketRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractSocketRegistryFilterer{contract: contract}, nil
}

// bindContractSocketRegistry binds a generic wrapper to an already deployed contract.
func bindContractSocketRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractSocketRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSocketRegistry *ContractSocketRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.ContractSocketRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractSocketRegistry *ContractSocketRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractSocketRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractSocketRegistry *ContractSocketRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractSocketRegistry *ContractSocketRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) GetOperatorSocket(opts *bind.CallOpts, _operatorId [32]byte) (string, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "getOperatorSocket", _operatorId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistrySession) GetOperatorSocket(_operatorId [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.GetOperatorSocket(&_ContractSocketRegistry.CallOpts, _operatorId)
}

// GetOperatorSocket is a free data retrieval call binding the contract method 0x10bea0d7.
//
// Solidity: function getOperatorSocket(bytes32 _operatorId) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) GetOperatorSocket(_operatorId [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.GetOperatorSocket(&_ContractSocketRegistry.CallOpts, _operatorId)
}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) OperatorIdToSocket(opts *bind.CallOpts, arg0 [32]byte) (string, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "operatorIdToSocket", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistrySession) OperatorIdToSocket(arg0 [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.OperatorIdToSocket(&_ContractSocketRegistry.CallOpts, arg0)
}

// OperatorIdToSocket is a free data retrieval call binding the contract method 0xaf65fdfc.
//
// Solidity: function operatorIdToSocket(bytes32 ) view returns(string)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) OperatorIdToSocket(arg0 [32]byte) (string, error) {
	return _ContractSocketRegistry.Contract.OperatorIdToSocket(&_ContractSocketRegistry.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractSocketRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractSocketRegistry.Contract.RegistryCoordinator(&_ContractSocketRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractSocketRegistry *ContractSocketRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractSocketRegistry.Contract.RegistryCoordinator(&_ContractSocketRegistry.CallOpts)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactor) SetOperatorSocket(opts *bind.TransactOpts, _operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.contract.Transact(opts, "setOperatorSocket", _operatorId, _socket)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistrySession) SetOperatorSocket(_operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.SetOperatorSocket(&_ContractSocketRegistry.TransactOpts, _operatorId, _socket)
}

// SetOperatorSocket is a paid mutator transaction binding the contract method 0xf043367e.
//
// Solidity: function setOperatorSocket(bytes32 _operatorId, string _socket) returns()
func (_ContractSocketRegistry *ContractSocketRegistryTransactorSession) SetOperatorSocket(_operatorId [32]byte, _socket string) (*types.Transaction, error) {
	return _ContractSocketRegistry.Contract.SetOperatorSocket(&_ContractSocketRegistry.TransactOpts, _operatorId, _socket)
}
