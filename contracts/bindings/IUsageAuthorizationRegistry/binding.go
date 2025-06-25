// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIUsageAuthorizationRegistry

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

// UsageAuthorizationTypesQuorumConfig is an auto generated low-level Go binding around an user-defined struct.
type UsageAuthorizationTypesQuorumConfig struct {
	Token                       common.Address
	Recipient                   common.Address
	ReservationSymbolsPerSecond uint64
	OnDemandSymbolsPerSecond    uint64
	OnDemandPricePerSymbol      uint64
}

// UsageAuthorizationTypesQuorumProtocolConfig is an auto generated low-level Go binding around an user-defined struct.
type UsageAuthorizationTypesQuorumProtocolConfig struct {
	MinNumSymbols              uint64
	ReservationAdvanceWindow   uint64
	ReservationRateLimitWindow uint64
	OnDemandRateLimitWindow    uint64
	OnDemandEnabled            bool
}

// UsageAuthorizationTypesReservation is an auto generated low-level Go binding around an user-defined struct.
type UsageAuthorizationTypesReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
}

// ContractIUsageAuthorizationRegistryMetaData contains all meta data concerning the ContractIUsageAuthorizationRegistry contract.
var ContractIUsageAuthorizationRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decreaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOnDemandDeposit\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumPaymentConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumReservedSymbols\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"period\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseReservation\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"cfg\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.QuorumConfig\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"reservationSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandSymbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandPricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumProtocolConfig\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"protocolCfg\",\"type\":\"tuple\",\"internalType\":\"structUsageAuthorizationTypes.QuorumProtocolConfig\",\"components\":[{\"name\":\"minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationAdvanceWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reservationRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandRateLimitWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"onDemandEnabled\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferQuorumOwnership\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AmountTooLarge\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidReservationPeriod\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidStartTimestamp\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotEnoughSymbolsAvailable\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"requiredSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"availableSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OnDemandDisabled\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"OwnerIsZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuorumOwnerAlreadySet\",\"inputs\":[{\"name\":\"quorumId\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationMustDecrease\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationMustIncrease\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationStillActive\",\"inputs\":[{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ReservationTooLong\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"maxLength\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SchedulePeriodCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StartTimestampMustMatch\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"TimestampSchedulePeriodMismatch\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"schedulePeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]",
}

// ContractIUsageAuthorizationRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIUsageAuthorizationRegistryMetaData.ABI instead.
var ContractIUsageAuthorizationRegistryABI = ContractIUsageAuthorizationRegistryMetaData.ABI

// ContractIUsageAuthorizationRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistry struct {
	ContractIUsageAuthorizationRegistryCaller     // Read-only binding to the contract
	ContractIUsageAuthorizationRegistryTransactor // Write-only binding to the contract
	ContractIUsageAuthorizationRegistryFilterer   // Log filterer for contract events
}

// ContractIUsageAuthorizationRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIUsageAuthorizationRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIUsageAuthorizationRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIUsageAuthorizationRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIUsageAuthorizationRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIUsageAuthorizationRegistrySession struct {
	Contract     *ContractIUsageAuthorizationRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                        // Call options to use throughout this session
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ContractIUsageAuthorizationRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIUsageAuthorizationRegistryCallerSession struct {
	Contract *ContractIUsageAuthorizationRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                              // Call options to use throughout this session
}

// ContractIUsageAuthorizationRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIUsageAuthorizationRegistryTransactorSession struct {
	Contract     *ContractIUsageAuthorizationRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                              // Transaction auth options to use throughout this session
}

// ContractIUsageAuthorizationRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistryRaw struct {
	Contract *ContractIUsageAuthorizationRegistry // Generic contract binding to access the raw methods on
}

// ContractIUsageAuthorizationRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistryCallerRaw struct {
	Contract *ContractIUsageAuthorizationRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIUsageAuthorizationRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIUsageAuthorizationRegistryTransactorRaw struct {
	Contract *ContractIUsageAuthorizationRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIUsageAuthorizationRegistry creates a new instance of ContractIUsageAuthorizationRegistry, bound to a specific deployed contract.
func NewContractIUsageAuthorizationRegistry(address common.Address, backend bind.ContractBackend) (*ContractIUsageAuthorizationRegistry, error) {
	contract, err := bindContractIUsageAuthorizationRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIUsageAuthorizationRegistry{ContractIUsageAuthorizationRegistryCaller: ContractIUsageAuthorizationRegistryCaller{contract: contract}, ContractIUsageAuthorizationRegistryTransactor: ContractIUsageAuthorizationRegistryTransactor{contract: contract}, ContractIUsageAuthorizationRegistryFilterer: ContractIUsageAuthorizationRegistryFilterer{contract: contract}}, nil
}

// NewContractIUsageAuthorizationRegistryCaller creates a new read-only instance of ContractIUsageAuthorizationRegistry, bound to a specific deployed contract.
func NewContractIUsageAuthorizationRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIUsageAuthorizationRegistryCaller, error) {
	contract, err := bindContractIUsageAuthorizationRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIUsageAuthorizationRegistryCaller{contract: contract}, nil
}

// NewContractIUsageAuthorizationRegistryTransactor creates a new write-only instance of ContractIUsageAuthorizationRegistry, bound to a specific deployed contract.
func NewContractIUsageAuthorizationRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIUsageAuthorizationRegistryTransactor, error) {
	contract, err := bindContractIUsageAuthorizationRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIUsageAuthorizationRegistryTransactor{contract: contract}, nil
}

// NewContractIUsageAuthorizationRegistryFilterer creates a new log filterer instance of ContractIUsageAuthorizationRegistry, bound to a specific deployed contract.
func NewContractIUsageAuthorizationRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIUsageAuthorizationRegistryFilterer, error) {
	contract, err := bindContractIUsageAuthorizationRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIUsageAuthorizationRegistryFilterer{contract: contract}, nil
}

// bindContractIUsageAuthorizationRegistry binds a generic wrapper to an already deployed contract.
func bindContractIUsageAuthorizationRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIUsageAuthorizationRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIUsageAuthorizationRegistry.Contract.ContractIUsageAuthorizationRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.ContractIUsageAuthorizationRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.ContractIUsageAuthorizationRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIUsageAuthorizationRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCaller) GetOnDemandDeposit(opts *bind.CallOpts, quorumId uint64, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractIUsageAuthorizationRegistry.contract.Call(opts, &out, "getOnDemandDeposit", quorumId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetOnDemandDeposit(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, account)
}

// GetOnDemandDeposit is a free data retrieval call binding the contract method 0x400c2322.
//
// Solidity: function getOnDemandDeposit(uint64 quorumId, address account) view returns(uint256)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerSession) GetOnDemandDeposit(quorumId uint64, account common.Address) (*big.Int, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetOnDemandDeposit(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, account)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCaller) GetQuorumPaymentConfig(opts *bind.CallOpts, quorumId uint64) (UsageAuthorizationTypesQuorumConfig, error) {
	var out []interface{}
	err := _ContractIUsageAuthorizationRegistry.contract.Call(opts, &out, "getQuorumPaymentConfig", quorumId)

	if err != nil {
		return *new(UsageAuthorizationTypesQuorumConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(UsageAuthorizationTypesQuorumConfig)).(*UsageAuthorizationTypesQuorumConfig)

	return out0, err

}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) GetQuorumPaymentConfig(quorumId uint64) (UsageAuthorizationTypesQuorumConfig, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumPaymentConfig(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId)
}

// GetQuorumPaymentConfig is a free data retrieval call binding the contract method 0x7a9426ca.
//
// Solidity: function getQuorumPaymentConfig(uint64 quorumId) view returns((address,address,uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerSession) GetQuorumPaymentConfig(quorumId uint64) (UsageAuthorizationTypesQuorumConfig, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumPaymentConfig(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCaller) GetQuorumProtocolConfig(opts *bind.CallOpts, quorumId uint64) (UsageAuthorizationTypesQuorumProtocolConfig, error) {
	var out []interface{}
	err := _ContractIUsageAuthorizationRegistry.contract.Call(opts, &out, "getQuorumProtocolConfig", quorumId)

	if err != nil {
		return *new(UsageAuthorizationTypesQuorumProtocolConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(UsageAuthorizationTypesQuorumProtocolConfig)).(*UsageAuthorizationTypesQuorumProtocolConfig)

	return out0, err

}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) GetQuorumProtocolConfig(quorumId uint64) (UsageAuthorizationTypesQuorumProtocolConfig, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumProtocolConfig(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId)
}

// GetQuorumProtocolConfig is a free data retrieval call binding the contract method 0x89a06b35.
//
// Solidity: function getQuorumProtocolConfig(uint64 quorumId) view returns((uint64,uint64,uint64,uint64,bool))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerSession) GetQuorumProtocolConfig(quorumId uint64) (UsageAuthorizationTypesQuorumProtocolConfig, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumProtocolConfig(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCaller) GetQuorumReservedSymbols(opts *bind.CallOpts, quorumId uint64, period uint64) (uint64, error) {
	var out []interface{}
	err := _ContractIUsageAuthorizationRegistry.contract.Call(opts, &out, "getQuorumReservedSymbols", quorumId, period)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumReservedSymbols(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, period)
}

// GetQuorumReservedSymbols is a free data retrieval call binding the contract method 0x4023a200.
//
// Solidity: function getQuorumReservedSymbols(uint64 quorumId, uint64 period) view returns(uint64)
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerSession) GetQuorumReservedSymbols(quorumId uint64, period uint64) (uint64, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetQuorumReservedSymbols(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, period)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCaller) GetReservation(opts *bind.CallOpts, quorumId uint64, account common.Address) (UsageAuthorizationTypesReservation, error) {
	var out []interface{}
	err := _ContractIUsageAuthorizationRegistry.contract.Call(opts, &out, "getReservation", quorumId, account)

	if err != nil {
		return *new(UsageAuthorizationTypesReservation), err
	}

	out0 := *abi.ConvertType(out[0], new(UsageAuthorizationTypesReservation)).(*UsageAuthorizationTypesReservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) GetReservation(quorumId uint64, account common.Address) (UsageAuthorizationTypesReservation, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetReservation(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, account)
}

// GetReservation is a free data retrieval call binding the contract method 0x00e691aa.
//
// Solidity: function getReservation(uint64 quorumId, address account) view returns((uint64,uint64,uint64))
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryCallerSession) GetReservation(quorumId uint64, account common.Address) (UsageAuthorizationTypesReservation, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.GetReservation(&_ContractIUsageAuthorizationRegistry.CallOpts, quorumId, account)
}

// AddReservation is a paid mutator transaction binding the contract method 0x11d86be0.
//
// Solidity: function addReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) AddReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "addReservation", quorumId, account, reservation)
}

// AddReservation is a paid mutator transaction binding the contract method 0x11d86be0.
//
// Solidity: function addReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) AddReservation(quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.AddReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, reservation)
}

// AddReservation is a paid mutator transaction binding the contract method 0x11d86be0.
//
// Solidity: function addReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) AddReservation(quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.AddReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) DecreaseReservation(opts *bind.TransactOpts, quorumId uint64, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "decreaseReservation", quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) DecreaseReservation(quorumId uint64, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.DecreaseReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, reservation)
}

// DecreaseReservation is a paid mutator transaction binding the contract method 0x3b50e8b8.
//
// Solidity: function decreaseReservation(uint64 quorumId, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) DecreaseReservation(quorumId uint64, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.DecreaseReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, reservation)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x5c0ffe4c.
//
// Solidity: function depositOnDemand(uint64 quorumId, address account, uint256 amount) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) DepositOnDemand(opts *bind.TransactOpts, quorumId uint64, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "depositOnDemand", quorumId, account, amount)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x5c0ffe4c.
//
// Solidity: function depositOnDemand(uint64 quorumId, address account, uint256 amount) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) DepositOnDemand(quorumId uint64, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.DepositOnDemand(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, amount)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x5c0ffe4c.
//
// Solidity: function depositOnDemand(uint64 quorumId, address account, uint256 amount) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) DepositOnDemand(quorumId uint64, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.DepositOnDemand(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, amount)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) IncreaseReservation(opts *bind.TransactOpts, quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "increaseReservation", quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) IncreaseReservation(quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.IncreaseReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, reservation)
}

// IncreaseReservation is a paid mutator transaction binding the contract method 0xe453a058.
//
// Solidity: function increaseReservation(uint64 quorumId, address account, (uint64,uint64,uint64) reservation) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) IncreaseReservation(quorumId uint64, account common.Address, reservation UsageAuthorizationTypesReservation) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.IncreaseReservation(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, account, reservation)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumId uint64, newOwner common.Address, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "initializeQuorum", quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.InitializeQuorum(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, newOwner, protocolCfg)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0x063a5d56.
//
// Solidity: function initializeQuorum(uint64 quorumId, address newOwner, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) InitializeQuorum(quorumId uint64, newOwner common.Address, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.InitializeQuorum(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, newOwner, protocolCfg)
}

// SetQuorumConfig is a paid mutator transaction binding the contract method 0xfcb7ec34.
//
// Solidity: function setQuorumConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) cfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) SetQuorumConfig(opts *bind.TransactOpts, quorumId uint64, cfg UsageAuthorizationTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "setQuorumConfig", quorumId, cfg)
}

// SetQuorumConfig is a paid mutator transaction binding the contract method 0xfcb7ec34.
//
// Solidity: function setQuorumConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) cfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) SetQuorumConfig(quorumId uint64, cfg UsageAuthorizationTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.SetQuorumConfig(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, cfg)
}

// SetQuorumConfig is a paid mutator transaction binding the contract method 0xfcb7ec34.
//
// Solidity: function setQuorumConfig(uint64 quorumId, (address,address,uint64,uint64,uint64) cfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) SetQuorumConfig(quorumId uint64, cfg UsageAuthorizationTypesQuorumConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.SetQuorumConfig(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, cfg)
}

// SetQuorumProtocolConfig is a paid mutator transaction binding the contract method 0x04cac6e5.
//
// Solidity: function setQuorumProtocolConfig(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) SetQuorumProtocolConfig(opts *bind.TransactOpts, quorumId uint64, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "setQuorumProtocolConfig", quorumId, protocolCfg)
}

// SetQuorumProtocolConfig is a paid mutator transaction binding the contract method 0x04cac6e5.
//
// Solidity: function setQuorumProtocolConfig(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) SetQuorumProtocolConfig(quorumId uint64, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.SetQuorumProtocolConfig(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, protocolCfg)
}

// SetQuorumProtocolConfig is a paid mutator transaction binding the contract method 0x04cac6e5.
//
// Solidity: function setQuorumProtocolConfig(uint64 quorumId, (uint64,uint64,uint64,uint64,bool) protocolCfg) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) SetQuorumProtocolConfig(quorumId uint64, protocolCfg UsageAuthorizationTypesQuorumProtocolConfig) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.SetQuorumProtocolConfig(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, protocolCfg)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.TransferOwnership(&_ContractIUsageAuthorizationRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.TransferOwnership(&_ContractIUsageAuthorizationRegistry.TransactOpts, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactor) TransferQuorumOwnership(opts *bind.TransactOpts, quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.contract.Transact(opts, "transferQuorumOwnership", quorumId, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistrySession) TransferQuorumOwnership(quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.TransferQuorumOwnership(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, newOwner)
}

// TransferQuorumOwnership is a paid mutator transaction binding the contract method 0xede57a2d.
//
// Solidity: function transferQuorumOwnership(uint64 quorumId, address newOwner) returns()
func (_ContractIUsageAuthorizationRegistry *ContractIUsageAuthorizationRegistryTransactorSession) TransferQuorumOwnership(quorumId uint64, newOwner common.Address) (*types.Transaction, error) {
	return _ContractIUsageAuthorizationRegistry.Contract.TransferQuorumOwnership(&_ContractIUsageAuthorizationRegistry.TransactOpts, quorumId, newOwner)
}
