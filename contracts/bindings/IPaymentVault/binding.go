// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIPaymentVault

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

// PaymentVaultTypesQuorumConfig is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesQuorumConfig struct {
	Token                       common.Address
	Recipient                   common.Address
	ReservationSymbolsPerSecond uint64
	OnDemandSymbolsPerSecond    uint64
	OnDemandPricePerSymbol      uint64
}

// PaymentVaultTypesQuorumProtocolConfig is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesQuorumProtocolConfig struct {
	MinNumSymbols              uint64
	ReservationAdvanceWindow   uint64
	ReservationRateLimitWindow uint64
	OnDemandRateLimitWindow    uint64
	OnDemandEnabled            bool
}

// PaymentVaultTypesReservation is an auto generated low-level Go binding around an user-defined struct.
type PaymentVaultTypesReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
}

// ContractIPaymentVaultMetaData contains all meta data concerning the ContractIPaymentVault contract.
var ContractIPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getOnDemandDeposit\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumReservedSymbols\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"period\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"}]",
}

// ContractIPaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIPaymentVaultMetaData.ABI instead.
var ContractIPaymentVaultABI = ContractIPaymentVaultMetaData.ABI

// ContractIPaymentVault is an auto generated Go binding around an Ethereum contract.
type ContractIPaymentVault struct {
	ContractIPaymentVaultCaller     // Read-only binding to the contract
	ContractIPaymentVaultTransactor // Write-only binding to the contract
	ContractIPaymentVaultFilterer   // Log filterer for contract events
}

// ContractIPaymentVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIPaymentVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIPaymentVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIPaymentVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIPaymentVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIPaymentVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIPaymentVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIPaymentVaultSession struct {
	Contract     *ContractIPaymentVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ContractIPaymentVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIPaymentVaultCallerSession struct {
	Contract *ContractIPaymentVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ContractIPaymentVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIPaymentVaultTransactorSession struct {
	Contract     *ContractIPaymentVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractIPaymentVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIPaymentVaultRaw struct {
	Contract *ContractIPaymentVault // Generic contract binding to access the raw methods on
}

// ContractIPaymentVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIPaymentVaultCallerRaw struct {
	Contract *ContractIPaymentVaultCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIPaymentVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIPaymentVaultTransactorRaw struct {
	Contract *ContractIPaymentVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIPaymentVault creates a new instance of ContractIPaymentVault, bound to a specific deployed contract.
func NewContractIPaymentVault(address common.Address, backend bind.ContractBackend) (*ContractIPaymentVault, error) {
	contract, err := bindContractIPaymentVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIPaymentVault{ContractIPaymentVaultCaller: ContractIPaymentVaultCaller{contract: contract}, ContractIPaymentVaultTransactor: ContractIPaymentVaultTransactor{contract: contract}, ContractIPaymentVaultFilterer: ContractIPaymentVaultFilterer{contract: contract}}, nil
}

// NewContractIPaymentVaultCaller creates a new read-only instance of ContractIPaymentVault, bound to a specific deployed contract.
func NewContractIPaymentVaultCaller(address common.Address, caller bind.ContractCaller) (*ContractIPaymentVaultCaller, error) {
	contract, err := bindContractIPaymentVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIPaymentVaultCaller{contract: contract}, nil
}

// NewContractIPaymentVaultTransactor creates a new write-only instance of ContractIPaymentVault, bound to a specific deployed contract.
func NewContractIPaymentVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIPaymentVaultTransactor, error) {
	contract, err := bindContractIPaymentVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIPaymentVaultTransactor{contract: contract}, nil
}

// NewContractIPaymentVaultFilterer creates a new log filterer instance of ContractIPaymentVault, bound to a specific deployed contract.
func NewContractIPaymentVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIPaymentVaultFilterer, error) {
	contract, err := bindContractIPaymentVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIPaymentVaultFilterer{contract: contract}, nil
}

// bindContractIPaymentVault binds a generic wrapper to an already deployed contract.
func bindContractIPaymentVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIPaymentVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIPaymentVault *ContractIPaymentVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIPaymentVault.Contract.ContractIPaymentVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIPaymentVault *ContractIPaymentVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIPaymentVault.Contract.ContractIPaymentVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIPaymentVault *ContractIPaymentVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIPaymentVault.Contract.ContractIPaymentVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIPaymentVault *ContractIPaymentVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIPaymentVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIPaymentVault *ContractIPaymentVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIPaymentVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIPaymentVault *ContractIPaymentVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIPaymentVault.Contract.contract.Transact(opts, method, params...)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIPaymentVault *ContractIPaymentVaultCaller) GetOnDemandDeposit(opts *bind.CallOpts, quorumId uint64, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractIPaymentVault.contract.Call(opts, &out, "getOnDemandDeposit", quorumId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIPaymentVault *ContractIPaymentVaultSession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractIPaymentVault.Contract.GetOnDemandDeposit(&_ContractIPaymentVault.CallOpts, quorumId, account)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIPaymentVault *ContractIPaymentVaultCallerSession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractIPaymentVault.Contract.GetOnDemandDeposit(&_ContractIPaymentVault.CallOpts, quorumId, account)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultCaller) GetQuorumPaymentConfig(opts *bind.CallOpts, quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	var out []interface{}
	err := _ContractIPaymentVault.contract.Call(opts, &out, "getQuorumPaymentConfig", quorumId)

	if err != nil {
		return *new(PaymentVaultTypesQuorumConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesQuorumConfig)).(*PaymentVaultTypesQuorumConfig)

	return out0, err

}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultSession) GetQuorumPaymentConfig(quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	return _ContractIPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractIPaymentVault.CallOpts, quorumId)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultCallerSession) GetQuorumPaymentConfig(quorumId uint64) (PaymentVaultTypesQuorumConfig, error) {
	return _ContractIPaymentVault.Contract.GetQuorumPaymentConfig(&_ContractIPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIPaymentVault *ContractIPaymentVaultCaller) GetQuorumProtocolConfig(opts *bind.CallOpts, quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
	var out []interface{}
	err := _ContractIPaymentVault.contract.Call(opts, &out, "getQuorumProtocolConfig", quorumId)

	if err != nil {
		return *new(PaymentVaultTypesQuorumProtocolConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesQuorumProtocolConfig)).(*PaymentVaultTypesQuorumProtocolConfig)

	return out0, err

}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIPaymentVault *ContractIPaymentVaultSession) GetQuorumProtocolConfig(quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
	return _ContractIPaymentVault.Contract.GetQuorumProtocolConfig(&_ContractIPaymentVault.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIPaymentVault *ContractIPaymentVaultCallerSession) GetQuorumProtocolConfig(quorumId uint64) (PaymentVaultTypesQuorumProtocolConfig, error) {
	return _ContractIPaymentVault.Contract.GetQuorumProtocolConfig(&_ContractIPaymentVault.CallOpts, quorumId)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIPaymentVault *ContractIPaymentVaultCaller) GetQuorumReservedSymbols(opts *bind.CallOpts, quorumId uint64, period uint64) (uint64, error) {
	var out []interface{}
	err := _ContractIPaymentVault.contract.Call(opts, &out, "getQuorumReservedSymbols", quorumId, period)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIPaymentVault *ContractIPaymentVaultSession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractIPaymentVault.Contract.GetQuorumReservedSymbols(&_ContractIPaymentVault.CallOpts, quorumId, period)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIPaymentVault *ContractIPaymentVaultCallerSession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractIPaymentVault.Contract.GetQuorumReservedSymbols(&_ContractIPaymentVault.CallOpts, quorumId, period)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultCaller) GetReservation(opts *bind.CallOpts, quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	var out []interface{}
	err := _ContractIPaymentVault.contract.Call(opts, &out, "getReservation", quorumId, account)

	if err != nil {
		return *new(PaymentVaultTypesReservation), err
	}

	out0 := *abi.ConvertType(out[0], new(PaymentVaultTypesReservation)).(*PaymentVaultTypesReservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultSession) GetReservation(quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	return _ContractIPaymentVault.Contract.GetReservation(&_ContractIPaymentVault.CallOpts, quorumId, account)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIPaymentVault *ContractIPaymentVaultCallerSession) GetReservation(quorumId uint64, account common.Address) (PaymentVaultTypesReservation, error) {
	return _ContractIPaymentVault.Contract.GetReservation(&_ContractIPaymentVault.CallOpts, quorumId, account)
}
