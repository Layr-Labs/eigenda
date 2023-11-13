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
)

// IIndexRegistryOperatorIndexUpdate is an auto generated low-level Go binding around an user-defined struct.
type IIndexRegistryOperatorIndexUpdate struct {
	FromBlockNumber uint32
	Index           uint32
}

// ContractIIndexRegistryMetaData contains all meta data concerning the ContractIIndexRegistry contract.
var ContractIIndexRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newIndex\",\"type\":\"uint32\"}],\"name\":\"GlobalIndexUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newIndex\",\"type\":\"uint32\"}],\"name\":\"QuorumIndexUpdate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"operatorIdsToSwap\",\"type\":\"bytes32[]\"}],\"name\":\"deregisterOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"getOperatorIndexForQuorumAtBlockNumberByIndex\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"getOperatorIndexUpdateOfOperatorIdForQuorumAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fromBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"internalType\":\"structIIndexRegistry.OperatorIndexUpdate\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"}],\"name\":\"getOperatorListForQuorumAtBlockNumber\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"getTotalOperatorsForQuorumAtBlockNumberByIndex\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"getTotalOperatorsUpdateForQuorumAtIndex\",\"outputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fromBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"internalType\":\"structIIndexRegistry.OperatorIndexUpdate\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"}],\"name\":\"registerOperator\",\"outputs\":[{\"internalType\":\"uint32[]\",\"name\":\"\",\"type\":\"uint32[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registryCoordinator\",\"outputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"}],\"name\":\"totalOperatorsForQuorum\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
	parsed, err := abi.JSON(strings.NewReader(ContractIIndexRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

// GetOperatorIndexForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xb9717755.
//
// Solidity: function getOperatorIndexForQuorumAtBlockNumberByIndex(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetOperatorIndexForQuorumAtBlockNumberByIndex(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getOperatorIndexForQuorumAtBlockNumberByIndex", operatorId, quorumNumber, blockNumber, index)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetOperatorIndexForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xb9717755.
//
// Solidity: function getOperatorIndexForQuorumAtBlockNumberByIndex(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetOperatorIndexForQuorumAtBlockNumberByIndex(operatorId [32]byte, quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorIndexForQuorumAtBlockNumberByIndex(&_ContractIIndexRegistry.CallOpts, operatorId, quorumNumber, blockNumber, index)
}

// GetOperatorIndexForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xb9717755.
//
// Solidity: function getOperatorIndexForQuorumAtBlockNumberByIndex(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetOperatorIndexForQuorumAtBlockNumberByIndex(operatorId [32]byte, quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorIndexForQuorumAtBlockNumberByIndex(&_ContractIIndexRegistry.CallOpts, operatorId, quorumNumber, blockNumber, index)
}

// GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex is a free data retrieval call binding the contract method 0x97faf021.
//
// Solidity: function getOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(bytes32 operatorId, uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getOperatorIndexUpdateOfOperatorIdForQuorumAtIndex", operatorId, quorumNumber, index)

	if err != nil {
		return *new(IIndexRegistryOperatorIndexUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryOperatorIndexUpdate)).(*IIndexRegistryOperatorIndexUpdate)

	return out0, err

}

// GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex is a free data retrieval call binding the contract method 0x97faf021.
//
// Solidity: function getOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(bytes32 operatorId, uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(operatorId [32]byte, quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(&_ContractIIndexRegistry.CallOpts, operatorId, quorumNumber, index)
}

// GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex is a free data retrieval call binding the contract method 0x97faf021.
//
// Solidity: function getOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(bytes32 operatorId, uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(operatorId [32]byte, quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorIndexUpdateOfOperatorIdForQuorumAtIndex(&_ContractIIndexRegistry.CallOpts, operatorId, quorumNumber, index)
}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetOperatorListForQuorumAtBlockNumber(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getOperatorListForQuorumAtBlockNumber", quorumNumber, blockNumber)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetOperatorListForQuorumAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorListForQuorumAtBlockNumber(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetOperatorListForQuorumAtBlockNumber is a free data retrieval call binding the contract method 0x889ae3e5.
//
// Solidity: function getOperatorListForQuorumAtBlockNumber(uint8 quorumNumber, uint32 blockNumber) view returns(bytes32[])
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetOperatorListForQuorumAtBlockNumber(quorumNumber uint8, blockNumber uint32) ([][32]byte, error) {
	return _ContractIIndexRegistry.Contract.GetOperatorListForQuorumAtBlockNumber(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber)
}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetTotalOperatorsForQuorumAtBlockNumberByIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getTotalOperatorsForQuorumAtBlockNumberByIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetTotalOperatorsForQuorumAtBlockNumberByIndex(quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIIndexRegistry.Contract.GetTotalOperatorsForQuorumAtBlockNumberByIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetTotalOperatorsForQuorumAtBlockNumberByIndex is a free data retrieval call binding the contract method 0xa454b3be.
//
// Solidity: function getTotalOperatorsForQuorumAtBlockNumberByIndex(uint8 quorumNumber, uint32 blockNumber, uint32 index) view returns(uint32)
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetTotalOperatorsForQuorumAtBlockNumberByIndex(quorumNumber uint8, blockNumber uint32, index uint32) (uint32, error) {
	return _ContractIIndexRegistry.Contract.GetTotalOperatorsForQuorumAtBlockNumberByIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetTotalOperatorsUpdateForQuorumAtIndex is a free data retrieval call binding the contract method 0x2d1a8b78.
//
// Solidity: function getTotalOperatorsUpdateForQuorumAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCaller) GetTotalOperatorsUpdateForQuorumAtIndex(opts *bind.CallOpts, quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	var out []interface{}
	err := _ContractIIndexRegistry.contract.Call(opts, &out, "getTotalOperatorsUpdateForQuorumAtIndex", quorumNumber, index)

	if err != nil {
		return *new(IIndexRegistryOperatorIndexUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IIndexRegistryOperatorIndexUpdate)).(*IIndexRegistryOperatorIndexUpdate)

	return out0, err

}

// GetTotalOperatorsUpdateForQuorumAtIndex is a free data retrieval call binding the contract method 0x2d1a8b78.
//
// Solidity: function getTotalOperatorsUpdateForQuorumAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) GetTotalOperatorsUpdateForQuorumAtIndex(quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetTotalOperatorsUpdateForQuorumAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, index)
}

// GetTotalOperatorsUpdateForQuorumAtIndex is a free data retrieval call binding the contract method 0x2d1a8b78.
//
// Solidity: function getTotalOperatorsUpdateForQuorumAtIndex(uint8 quorumNumber, uint32 index) view returns((uint32,uint32))
func (_ContractIIndexRegistry *ContractIIndexRegistryCallerSession) GetTotalOperatorsUpdateForQuorumAtIndex(quorumNumber uint8, index uint32) (IIndexRegistryOperatorIndexUpdate, error) {
	return _ContractIIndexRegistry.Contract.GetTotalOperatorsUpdateForQuorumAtIndex(&_ContractIIndexRegistry.CallOpts, quorumNumber, index)
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

// DeregisterOperator is a paid mutator transaction binding the contract method 0x0854259f.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers, bytes32[] operatorIdsToSwap) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte, operatorIdsToSwap [][32]byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.contract.Transact(opts, "deregisterOperator", operatorId, quorumNumbers, operatorIdsToSwap)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x0854259f.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers, bytes32[] operatorIdsToSwap) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistrySession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte, operatorIdsToSwap [][32]byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.DeregisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers, operatorIdsToSwap)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x0854259f.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers, bytes32[] operatorIdsToSwap) returns()
func (_ContractIIndexRegistry *ContractIIndexRegistryTransactorSession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte, operatorIdsToSwap [][32]byte) (*types.Transaction, error) {
	return _ContractIIndexRegistry.Contract.DeregisterOperator(&_ContractIIndexRegistry.TransactOpts, operatorId, quorumNumbers, operatorIdsToSwap)
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

// ContractIIndexRegistryGlobalIndexUpdateIterator is returned from FilterGlobalIndexUpdate and is used to iterate over the raw logs and unpacked data for GlobalIndexUpdate events raised by the ContractIIndexRegistry contract.
type ContractIIndexRegistryGlobalIndexUpdateIterator struct {
	Event *ContractIIndexRegistryGlobalIndexUpdate // Event containing the contract specifics and raw log

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
func (it *ContractIIndexRegistryGlobalIndexUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIIndexRegistryGlobalIndexUpdate)
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
		it.Event = new(ContractIIndexRegistryGlobalIndexUpdate)
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
func (it *ContractIIndexRegistryGlobalIndexUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIIndexRegistryGlobalIndexUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIIndexRegistryGlobalIndexUpdate represents a GlobalIndexUpdate event raised by the ContractIIndexRegistry contract.
type ContractIIndexRegistryGlobalIndexUpdate struct {
	OperatorId [32]byte
	NewIndex   uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterGlobalIndexUpdate is a free log retrieval operation binding the contract event 0x5bbd5cd0459a4845f4d6da0ce59566bd41aa74c4e71e2a71b9edd1ec25433c87.
//
// Solidity: event GlobalIndexUpdate(bytes32 indexed operatorId, uint32 newIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) FilterGlobalIndexUpdate(opts *bind.FilterOpts, operatorId [][32]byte) (*ContractIIndexRegistryGlobalIndexUpdateIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIIndexRegistry.contract.FilterLogs(opts, "GlobalIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractIIndexRegistryGlobalIndexUpdateIterator{contract: _ContractIIndexRegistry.contract, event: "GlobalIndexUpdate", logs: logs, sub: sub}, nil
}

// WatchGlobalIndexUpdate is a free log subscription operation binding the contract event 0x5bbd5cd0459a4845f4d6da0ce59566bd41aa74c4e71e2a71b9edd1ec25433c87.
//
// Solidity: event GlobalIndexUpdate(bytes32 indexed operatorId, uint32 newIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) WatchGlobalIndexUpdate(opts *bind.WatchOpts, sink chan<- *ContractIIndexRegistryGlobalIndexUpdate, operatorId [][32]byte) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractIIndexRegistry.contract.WatchLogs(opts, "GlobalIndexUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIIndexRegistryGlobalIndexUpdate)
				if err := _ContractIIndexRegistry.contract.UnpackLog(event, "GlobalIndexUpdate", log); err != nil {
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

// ParseGlobalIndexUpdate is a log parse operation binding the contract event 0x5bbd5cd0459a4845f4d6da0ce59566bd41aa74c4e71e2a71b9edd1ec25433c87.
//
// Solidity: event GlobalIndexUpdate(bytes32 indexed operatorId, uint32 newIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) ParseGlobalIndexUpdate(log types.Log) (*ContractIIndexRegistryGlobalIndexUpdate, error) {
	event := new(ContractIIndexRegistryGlobalIndexUpdate)
	if err := _ContractIIndexRegistry.contract.UnpackLog(event, "GlobalIndexUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	OperatorId   [32]byte
	QuorumNumber uint8
	NewIndex     uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterQuorumIndexUpdate is a free log retrieval operation binding the contract event 0x6ee1e4f4075f3d067176140d34e87874244dd273294c05b2218133e49a2ba6f6.
//
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
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
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
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
// Solidity: event QuorumIndexUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint32 newIndex)
func (_ContractIIndexRegistry *ContractIIndexRegistryFilterer) ParseQuorumIndexUpdate(log types.Log) (*ContractIIndexRegistryQuorumIndexUpdate, error) {
	event := new(ContractIIndexRegistryQuorumIndexUpdate)
	if err := _ContractIIndexRegistry.contract.UnpackLog(event, "QuorumIndexUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
