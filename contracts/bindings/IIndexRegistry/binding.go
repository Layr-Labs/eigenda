// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIIndexRegistry

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

// IIndexRegistryOperatorUpdate is an auto generated low-level Go binding around an user-defined struct.
type IIndexRegistryOperatorUpdate struct {
	FromBlockNumber uint32
	OperatorId      [32]byte
}

// IIndexRegistryQuorumUpdate is an auto generated low-level Go binding around an user-defined struct.
type IIndexRegistryQuorumUpdate struct {
	FromBlockNumber uint32
	NumOperators    uint32
}

// ContractIIndexRegistryMetaData contains all meta data concerning the ContractIIndexRegistry contract.
var ContractIIndexRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getLatestOperatorUpdate\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operatorIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.OperatorUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestQuorumUpdate\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.QuorumUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorListAtBlockNumber\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operatorIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"arrayIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.OperatorUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quorumIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIIndexRegistry.QuorumUpdate\",\"components\":[{\"name\":\"fromBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalOperatorsForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"QuorumIndexUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"newOperatorIndex\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false}]",
}

// ContractIIndexRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIIndexRegistryMetaData.ABI instead.
var ContractIIndexRegistryABI = ContractIIndexRegistryMetaData.ABI

// ContractIIndexRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIIndexRegistry struct {
	ContractIIndexRegistryCaller     // Read-only binding to the contract
	ContractIIndexRegistryTransactor // Write-only binding to the contract
	ContractIIndexRegistryFilterer   // Log filterer for contract events
}

// ContractIIndexRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIIndexRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIIndexRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIIndexRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIIndexRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIIndexRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIIndexRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIIndexRegistrySession struct {
	Contract     *ContractIIndexRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ContractIIndexRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIIndexRegistryCallerSession struct {
	Contract *ContractIIndexRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ContractIIndexRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIIndexRegistryTransactorSession struct {
	Contract     *ContractIIndexRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractIIndexRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIIndexRegistryRaw struct {
	Contract *ContractIIndexRegistry // Generic contract binding to access the raw methods on
}

// ContractIIndexRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIIndexRegistryCallerRaw struct {
	Contract *ContractIIndexRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIIndexRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIIndexRegistryTransactorRaw struct {
	Contract *ContractIIndexRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIIndexRegistry creates a new instance of ContractIIndexRegistry, bound to a specific deployed contract.
func NewContractIIndexRegistry(address common.Address, backend bind.ContractBackend) (*ContractIIndexRegistry, error) {
	contract, err := bindContractIIndexRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistry{ContractIIndexRegistryCaller: ContractIIndexRegistryCaller{contract: contract}, ContractIIndexRegistryTransactor: ContractIIndexRegistryTransactor{contract: contract}, ContractIIndexRegistryFilterer: ContractIIndexRegistryFilterer{contract: contract}}, nil
}

// NewContractIIndexRegistryCaller creates a new read-only instance of ContractIIndexRegistry, bound to a specific deployed contract.
func NewContractIIndexRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIIndexRegistryCaller, error) {
	contract, err := bindContractIIndexRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistryCaller{contract: contract}, nil
}

// NewContractIIndexRegistryTransactor creates a new write-only instance of ContractIIndexRegistry, bound to a specific deployed contract.
func NewContractIIndexRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIIndexRegistryTransactor, error) {
	contract, err := bindContractIIndexRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistryTransactor{contract: contract}, nil
}

// NewContractIIndexRegistryFilterer creates a new log filterer instance of ContractIIndexRegistry, bound to a specific deployed contract.
func NewContractIIndexRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIIndexRegistryFilterer, error) {
	contract, err := bindContractIIndexRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistryFilterer{contract: contract}, nil
}

// bindContractIIndexRegistry binds a generic wrapper to an already deployed contract.
func bindContractIIndexRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIIndexRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIIndexRegistry *ContractIIndexRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIIndexRegistry.Contract.ContractIIndexRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIIndexRegistry *ContractIIndexRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.ContractIIndexRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIIndexRegistry *ContractIIndexRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.ContractIIndexRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIIndexRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetLatestOperatorUpdate is a free data retrieval call binding the contract method 0x12d1d74d.
//
// Solidity: function getLatestOperatorUpdate(uint8 quorumNumber, uint32 operatorIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetLatestOperatorUpdate(opts *bind.CallOpts, quorumNumber uint8, operatorIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getLatestOperatorUpdate", quorumNumber, operatorIndex)

	if err != nil {
		return *new(IIndexRegistryOperatorUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryOperatorUpdate)).(*IIndexRegistryOperatorUpdate)

	return out0, err

}

// GetLatestOperatorUpdate is a free data retrieval call binding the contract method 0x12d1d74d.
//
// Solidity: function getLatestOperatorUpdate(uint8 quorumNumber, uint32 operatorIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetLatestOperatorUpdate(quorumNumber uint8, operatorIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetLatestOperatorUpdate(&_ContractIIndexRegistry.CallOpts, quorumNumber, operatorIndex)
}

// GetLatestOperatorUpdate is a free data retrieval call binding the contract method 0x12d1d74d.
//
// Solidity: function getLatestOperatorUpdate(uint8 quorumNumber, uint32 operatorIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetLatestOperatorUpdate(quorumNumber uint8, operatorIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetLatestOperatorUpdate(&_ContractIIndexRegistry.CallOpts, quorumNumber, operatorIndex)
}

// GetLatestQuorumUpdate is a free data retrieval call binding the contract method 0x8121906f.
//
// Solidity: function getLatestQuorumUpdate(uint8 quorumNumber) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetLatestQuorumUpdate(opts *bind.CallOpts, quorumNumber uint8) (IIndexRegistryQuorumUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getLatestQuorumUpdate", quorumNumber)

	if err != nil {
		return *new(IIndexRegistryQuorumUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryQuorumUpdate)).(*IIndexRegistryQuorumUpdate)

	return out0, err

}

// GetLatestQuorumUpdate is a free data retrieval call binding the contract method 0x8121906f.
//
// Solidity: function getLatestQuorumUpdate(uint8 quorumNumber) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetLatestQuorumUpdate(quorumNumber uint8) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetLatestQuorumUpdate(&_ContractIIndexRegistry.CallOpts, quorumNumber)
}

// GetLatestQuorumUpdate is a free data retrieval call binding the contract method 0x8121906f.
//
// Solidity: function getLatestQuorumUpdate(uint8 quorumNumber) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetLatestQuorumUpdate(quorumNumber uint8) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetLatestQuorumUpdate(&_ContractIIndexRegistry.CallOpts, quorumNumber)
}

// GetOperatorListAtBlockNumber is a free data retrieval call binding the contract method 0x89026245.
//
// Solidity: function getOperatorListAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetOperatorListAtBlockNumber(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getOperatorListAtBlockNumber", quorumNumber, blockNumber)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetOperatorListAtBlockNumber is a free data retrieval call binding the contract method 0x89026245.
//
// Solidity: function getOperatorListAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetOperatorListAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorListAtBlockNumber(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetOperatorListAtBlockNumber is a free data retrieval call binding the contract method 0x89026245.
//
// Solidity: function getOperatorListAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetOperatorListAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorListAtBlockNumber(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetOperatorUpdateAtIndex is a free data retrieval call binding the contract method 0x2ed583e5.
//
// Solidity: function getOperatorUpdateAtIndex(uint8 quorumNumber, uint32 operatorIndex, uint32 arrayIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetOperatorUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, operatorIndex uint32, arrayIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getOperatorUpdateAtIndex", quorumNumber, operatorIndex, arrayIndex)

	if err != nil {
		return *new(IIndexRegistryOperatorUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryOperatorUpdate)).(*IIndexRegistryOperatorUpdate)

	return out0, err

}

// GetOperatorUpdateAtIndex is a free data retrieval call binding the contract method 0x2ed583e5.
//
// Solidity: function getOperatorUpdateAtIndex(uint8 quorumNumber, uint32 operatorIndex, uint32 arrayIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetOperatorUpdateAtIndex(quorumNumber uint8, operatorIndex uint32, arrayIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorUpdateAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, operatorIndex, arrayIndex)
}

// GetOperatorUpdateAtIndex is a free data retrieval call binding the contract method 0x2ed583e5.
//
// Solidity: function getOperatorUpdateAtIndex(uint8 quorumNumber, uint32 operatorIndex, uint32 arrayIndex) view returns((uint32,bytes32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetOperatorUpdateAtIndex(quorumNumber uint8, operatorIndex uint32, arrayIndex uint32) (IIndexRegistryOperatorUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorUpdateAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, operatorIndex, arrayIndex)
}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 quorumIndex) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetQuorumUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, quorumIndex uint32) (IIndexRegistryQuorumUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getQuorumUpdateAtIndex", quorumNumber, quorumIndex)

	if err != nil {
		return *new(IIndexRegistryQuorumUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryQuorumUpdate)).(*IIndexRegistryQuorumUpdate)

	return out0, err

}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 quorumIndex) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetQuorumUpdateAtIndex(quorumNumber uint8, quorumIndex uint32) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetQuorumUpdateAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, quorumIndex)
}

// GetQuorumUpdateAtIndex is a free data retrieval call binding the contract method 0xa48bb0ac.
//
// Solidity: function getQuorumUpdateAtIndex(uint8 quorumNumber, uint32 quorumIndex) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetQuorumUpdateAtIndex(quorumNumber uint8, quorumIndex uint32) (IIndexRegistryQuorumUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetQuorumUpdateAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, quorumIndex)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractIIndexRegistry.Contract.RegistryCoordinator(&_ContractIIndexRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractIIndexRegistry.Contract.RegistryCoordinator(&_ContractIIndexRegistry.CallOpts)
}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) TotalOperatorsForQuorum(opts *bind.CallOpts, quorumNumber uint8) (uint32, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "totalOperatorsForQuorum", quorumNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) TotalOperatorsForQuorum(quorumNumber uint8) (uint32, error) {
	return _ContractIIndexRegistry.Contract.TotalOperatorsForQuorum(&_ContractIIndexRegistry.CallOpts, quorumNumber)
}

// TotalOperatorsForQuorum is a free data retrieval call binding the contract method 0xf3410922.
//
// Solidity: function totalOperatorsForQuorum(uint8 quorumNumber) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) TotalOperatorsForQuorum(quorumNumber uint8) (uint32, error) {
	return _ContractIIndexRegistry.Contract.TotalOperatorsForQuorum(&_ContractIIndexRegistry.CallOpts, quorumNumber)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.contract.Transact(opts, "deregisterOperator", operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.DeregisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorSession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.DeregisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumNumber uint8) (*types.Transaction, error) {
	return _ContractIIndexRegistry.contract.Transact(opts, "initializeQuorum", quorumNumber)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) InitializeQuorum(quorumNumber uint8) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.InitializeQuorum(&_ContractIIndexRegistry.TransactOpts, quorumNumber)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x26d941f2.
//
// Solidity: function initializeQuorum(uint8 quorumNumber) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorSession) InitializeQuorum(quorumNumber uint8) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.InitializeQuorum(&_ContractIIndexRegistry.TransactOpts, quorumNumber)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.contract.Transact(opts, "registerOperator", operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) RegisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.RegisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x00bff04d.
//
// Solidity: function registerOperator(bytes32 operatorId, bytes quorumNumbers) returns(uint32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorSession) RegisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.RegisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers)
}

// ContractIIndexRegistryQuorumIndexUpdateIterator is returned from FilterQuorumIndexUpdate and is used to iterate over the raw logs and unpacked data for QuorumIndexUpdate events raised by the ContractIIndexRegistry contract.
type ContractIIndexRegistryQuorumIndexUpdateIterator struct {
	Event *ContractIIndexRegistryQuorumIndexUpdate // Event containing the contract specifics and raw log

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
func (it *ContractIIndexRegistryQuorumIndexUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIIndexRegistryQuorumIndexUpdate)
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
		it.Event = new(ContractIIndexRegistryQuorumIndexUpdate)
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
func (it *ContractIIndexRegistryQuorumIndexUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIIndexRegistryQuorumIndexUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIIndexRegistryQuorumIndexUpdate represents a QuorumIndexUpdate event raised by the ContractIIndexRegistry contract.
type ContractIIndexRegistryQuorumIndexUpdate struct {
	OperatorId       [32]byte
	QuorumNumber     uint8
	NewOperatorIndex uint32
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterQuorumIndexUpdate is a free log retrieval operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newOperatorIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) FilterQuorumIndexUpdate(opts *bind.FilterOpts, operatorId [][32]byte) (*ContractIIndexRegistryQuorumIndexUpdateIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIIndexRegistry.contract.FilterLogs(opts, "QuorumIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistryQuorumIndexUpdateIterator{contract: _ContractIIndexRegistry.contract, event: "QuorumIndexUpdate", logs: logs, sub: sub}, nil
}

// WatchQuorumIndexUpdate is a free log subscription operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newOperatorIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) WatchQuorumIndexUpdate(opts *bind.WatchOpts, sink chan<- *ContractIIndexRegistryQuorumIndexUpdate, operatorId [][32]byte) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIIndexRegistry.contract.WatchLogs(opts, "QuorumIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIIndexRegistryQuorumIndexUpdate)
				if err := _ContractIIndexRegistry.contract.UnpackLog(event, "QuorumIndexUpdate", log); err != nil {
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

// ParseQuorumIndexUpdate is a log parse operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newOperatorIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) ParseQuorumIndexUpdate(log types.Log) (*ContractIIndexRegistryQuorumIndexUpdate, error) {
	event := new(ContractIIndexRegistryQuorumIndexUpdate)
	if err := _ContractIIndexRegistry.contract.UnpackLog(event, "QuorumIndexUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
