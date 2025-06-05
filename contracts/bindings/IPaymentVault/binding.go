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
	ABI: "[{\"type\":\"function\",\"name\":\"getOnDemandDeposit\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumReservedSymbols\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"period\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ReservationCreated\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structPaymentVaultTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AmountTooLarge\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidReservationPeriod\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidStartTimestamp\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuorumOwnerAlreadySet\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationMustDecrease\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationMustIncrease\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationStillActive\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationTooLong\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxLength\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SchedulePeriodCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StartTimestampMustMatch\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TimestampSchedulePeriodMismatch\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"schedulePeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]",
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

// ContractIPaymentVaultReservationCreatedIterator is returned from FilterReservationCreated and is used to iterate over the raw logs and unpacked data for ReservationCreated events raised by the ContractIPaymentVault contract.
type ContractIPaymentVaultReservationCreatedIterator struct {
	Event *ContractIPaymentVaultReservationCreated // Event containing the contract specifics and raw log

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
func (it *ContractIPaymentVaultReservationCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIPaymentVaultReservationCreated)
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
		it.Event = new(ContractIPaymentVaultReservationCreated)
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
func (it *ContractIPaymentVaultReservationCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIPaymentVaultReservationCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIPaymentVaultReservationCreated represents a ReservationCreated event raised by the ContractIPaymentVault contract.
type ContractIPaymentVaultReservationCreated struct {
	QuorumId    uint64
	Account     common.Address
	Reservation PaymentVaultTypesReservation
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReservationCreated is a free log retrieval operation binding the contract event 0x187bf155c286cdeb9b324f712e7d65b28d230a9f5a32a526e9056bd7671cbc5e.
//
// Solidity: event ReservationCreated(uint64 indexed quorumId, address indexed account, (uint64,uint64,uint64) reservation)
func (_ContractIPaymentVault *ContractIPaymentVaultFilterer) FilterReservationCreated(opts *bind.FilterOpts, quorumId []uint64, account []common.Address) (*ContractIPaymentVaultReservationCreatedIterator, error) {

	var quorumIdRule []interface{}
	for _, quorumIdItem := range quorumId {
		quorumIdRule = append(quorumIdRule, quorumIdItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractIPaymentVault.contract.FilterLogs(opts, "ReservationCreated", quorumIdRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractIPaymentVaultReservationCreatedIterator{contract: _ContractIPaymentVault.contract, event: "ReservationCreated", logs: logs, sub: sub}, nil
}

// WatchReservationCreated is a free log subscription operation binding the contract event 0x187bf155c286cdeb9b324f712e7d65b28d230a9f5a32a526e9056bd7671cbc5e.
//
// Solidity: event ReservationCreated(uint64 indexed quorumId, address indexed account, (uint64,uint64,uint64) reservation)
func (_ContractIPaymentVault *ContractIPaymentVaultFilterer) WatchReservationCreated(opts *bind.WatchOpts, sink chan<- *ContractIPaymentVaultReservationCreated, quorumId []uint64, account []common.Address) (event.Subscription, error) {

	var quorumIdRule []interface{}
	for _, quorumIdItem := range quorumId {
		quorumIdRule = append(quorumIdRule, quorumIdItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractIPaymentVault.contract.WatchLogs(opts, "ReservationCreated", quorumIdRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIPaymentVaultReservationCreated)
				if err := _ContractIPaymentVault.contract.UnpackLog(event, "ReservationCreated", log); err != nil {
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

// ParseReservationCreated is a log parse operation binding the contract event 0x187bf155c286cdeb9b324f712e7d65b28d230a9f5a32a526e9056bd7671cbc5e.
//
// Solidity: event ReservationCreated(uint64 indexed quorumId, address indexed account, (uint64,uint64,uint64) reservation)
func (_ContractIPaymentVault *ContractIPaymentVaultFilterer) ParseReservationCreated(log types.Log) (*ContractIPaymentVaultReservationCreated, error) {
	event := new(ContractIPaymentVaultReservationCreated)
	if err := _ContractIPaymentVault.contract.UnpackLog(event, "ReservationCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
