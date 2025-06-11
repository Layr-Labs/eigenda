// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractPaymentVault

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

// IPaymentVaultReservation is an auto generated low-level Go binding around an user-defined struct.
type IPaymentVaultReservation struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}

// ContractPaymentVaultMetaData contains all meta data concerning the ContractPaymentVault contract.
var ContractPaymentVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositOnDemand\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposit\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOnDemandTotalDeposits\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_payments\",\"type\":\"uint80[]\",\"internalType\":\"uint80[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReservations\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"_reservations\",\"type\":\"tuple[]\",\"internalType\":\"structIPaymentVault.Reservation[]\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalRatePeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalSymbolsPerPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastPriceUpdateTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minNumSymbols\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onDemandPayments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricePerSymbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceUpdateCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reservationPeriodInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reservations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGlobalRatePeriodInterval\",\"inputs\":[{\"name\":\"_globalRatePeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalSymbolsPerPeriod\",\"inputs\":[{\"name\":\"_globalSymbolsPerPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPriceParams\",\"inputs\":[{\"name\":\"_minNumSymbols\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_pricePerSymbol\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_priceUpdateCooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservation\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reservation\",\"type\":\"tuple\",\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReservationPeriodInterval\",\"inputs\":[{\"name\":\"_reservationPeriodInterval\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawERC20\",\"inputs\":[{\"name\":\"_token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalRatePeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalSymbolsPerPeriodUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnDemandPaymentUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onDemandPayment\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"},{\"name\":\"totalDeposit\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceParamsUpdated\",\"inputs\":[{\"name\":\"previousMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newMinNumSymbols\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPricePerSymbol\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"previousPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newPriceUpdateCooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationPeriodIntervalUpdated\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"newValue\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReservationUpdated\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reservation\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIPaymentVault.Reservation\",\"components\":[{\"name\":\"symbolsPerSecond\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumSplits\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false}]",
}

// ContractPaymentVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractPaymentVaultMetaData.ABI instead.
var ContractPaymentVaultABI = ContractPaymentVaultMetaData.ABI

// ContractPaymentVault is an auto generated Go binding around an Ethereum contract.
type ContractPaymentVault struct {
	ContractPaymentVaultCaller     // Read-only binding to the contract
	ContractPaymentVaultTransactor // Write-only binding to the contract
	ContractPaymentVaultFilterer   // Log filterer for contract events
}

// ContractPaymentVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractPaymentVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractPaymentVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractPaymentVaultSession struct {
	Contract     *ContractPaymentVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContractPaymentVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractPaymentVaultCallerSession struct {
	Contract *ContractPaymentVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ContractPaymentVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractPaymentVaultTransactorSession struct {
	Contract     *ContractPaymentVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractPaymentVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractPaymentVaultRaw struct {
	Contract *ContractPaymentVault // Generic contract binding to access the raw methods on
}

// ContractPaymentVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractPaymentVaultCallerRaw struct {
	Contract *ContractPaymentVaultCaller // Generic read-only contract binding to access the raw methods on
}

// ContractPaymentVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractPaymentVaultTransactorRaw struct {
	Contract *ContractPaymentVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractPaymentVault creates a new instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVault(address common.Address, backend bind.ContractBackend) (*ContractPaymentVault, error) {
	contract, err := bindContractPaymentVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVault{ContractPaymentVaultCaller: ContractPaymentVaultCaller{contract: contract}, ContractPaymentVaultTransactor: ContractPaymentVaultTransactor{contract: contract}, ContractPaymentVaultFilterer: ContractPaymentVaultFilterer{contract: contract}}, nil
}

// NewContractPaymentVaultCaller creates a new read-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultCaller(address common.Address, caller bind.ContractCaller) (*ContractPaymentVaultCaller, error) {
	contract, err := bindContractPaymentVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultCaller{contract: contract}, nil
}

// NewContractPaymentVaultTransactor creates a new write-only instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractPaymentVaultTransactor, error) {
	contract, err := bindContractPaymentVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultTransactor{contract: contract}, nil
}

// NewContractPaymentVaultFilterer creates a new log filterer instance of ContractPaymentVault, bound to a specific deployed contract.
func NewContractPaymentVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractPaymentVaultFilterer, error) {
	contract, err := bindContractPaymentVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultFilterer{contract: contract}, nil
}

// bindContractPaymentVault binds a generic wrapper to an already deployed contract.
func bindContractPaymentVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractPaymentVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.ContractPaymentVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.ContractPaymentVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractPaymentVault *ContractPaymentVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractPaymentVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractPaymentVault *ContractPaymentVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.contract.Transact(opts, method, params...)
}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandTotalDeposit(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandTotalDeposit", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandTotalDeposit(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposit(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandTotalDeposit is a free data retrieval call binding the contract method 0xd1c1fdcd.
//
// Solidity: function getOnDemandTotalDeposit(address _account) view returns(uint80)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandTotalDeposit(_account common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposit(&_ContractPaymentVault.CallOpts, _account)
}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetOnDemandTotalDeposits(opts *bind.CallOpts, _accounts []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getOnDemandTotalDeposits", _accounts)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetOnDemandTotalDeposits(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposits(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetOnDemandTotalDeposits is a free data retrieval call binding the contract method 0x4184a674.
//
// Solidity: function getOnDemandTotalDeposits(address[] _accounts) view returns(uint80[] _payments)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetOnDemandTotalDeposits(_accounts []common.Address) ([]*big.Int, error) {
	return _ContractPaymentVault.Contract.GetOnDemandTotalDeposits(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservation(opts *bind.CallOpts, _account common.Address) (IPaymentVaultReservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservation", _account)

	if err != nil {
		return *new(IPaymentVaultReservation), err
	}

	out0 := *abi.ConvertType(out[0], new(IPaymentVaultReservation)).(*IPaymentVaultReservation)

	return out0, err

}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservation(_account common.Address) (IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, _account)
}

// GetReservation is a free data retrieval call binding the contract method 0xb2066f80.
//
// Solidity: function getReservation(address _account) view returns((uint64,uint64,uint64,bytes,bytes))
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservation(_account common.Address) (IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservation(&_ContractPaymentVault.CallOpts, _account)
}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GetReservations(opts *bind.CallOpts, _accounts []common.Address) ([]IPaymentVaultReservation, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "getReservations", _accounts)

	if err != nil {
		return *new([]IPaymentVaultReservation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IPaymentVaultReservation)).(*[]IPaymentVaultReservation)

	return out0, err

}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultSession) GetReservations(_accounts []common.Address) ([]IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservations(&_ContractPaymentVault.CallOpts, _accounts)
}

// GetReservations is a free data retrieval call binding the contract method 0x109f8fe5.
//
// Solidity: function getReservations(address[] _accounts) view returns((uint64,uint64,uint64,bytes,bytes)[] _reservations)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GetReservations(_accounts []common.Address) ([]IPaymentVaultReservation, error) {
	return _ContractPaymentVault.Contract.GetReservations(&_ContractPaymentVault.CallOpts, _accounts)
}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalRatePeriodInterval(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalRatePeriodInterval")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalRatePeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalRatePeriodInterval(&_ContractPaymentVault.CallOpts)
}

// GlobalRatePeriodInterval is a free data retrieval call binding the contract method 0xbff8a3d4.
//
// Solidity: function globalRatePeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalRatePeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalRatePeriodInterval(&_ContractPaymentVault.CallOpts)
}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) GlobalSymbolsPerPeriod(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "globalSymbolsPerPeriod")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) GlobalSymbolsPerPeriod() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerPeriod(&_ContractPaymentVault.CallOpts)
}

// GlobalSymbolsPerPeriod is a free data retrieval call binding the contract method 0xc98d97dd.
//
// Solidity: function globalSymbolsPerPeriod() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) GlobalSymbolsPerPeriod() (uint64, error) {
	return _ContractPaymentVault.Contract.GlobalSymbolsPerPeriod(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) LastPriceUpdateTime(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "lastPriceUpdateTime")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) LastPriceUpdateTime() (uint64, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// LastPriceUpdateTime is a free data retrieval call binding the contract method 0x49b9a7af.
//
// Solidity: function lastPriceUpdateTime() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) LastPriceUpdateTime() (uint64, error) {
	return _ContractPaymentVault.Contract.LastPriceUpdateTime(&_ContractPaymentVault.CallOpts)
}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) MinNumSymbols(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "minNumSymbols")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) MinNumSymbols() (uint64, error) {
	return _ContractPaymentVault.Contract.MinNumSymbols(&_ContractPaymentVault.CallOpts)
}

// MinNumSymbols is a free data retrieval call binding the contract method 0x761dab89.
//
// Solidity: function minNumSymbols() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) MinNumSymbols() (uint64, error) {
	return _ContractPaymentVault.Contract.MinNumSymbols(&_ContractPaymentVault.CallOpts)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultCaller) OnDemandPayments(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "onDemandPayments", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultSession) OnDemandPayments(arg0 common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.OnDemandPayments(&_ContractPaymentVault.CallOpts, arg0)
}

// OnDemandPayments is a free data retrieval call binding the contract method 0xd996dc99.
//
// Solidity: function onDemandPayments(address ) view returns(uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) OnDemandPayments(arg0 common.Address) (*big.Int, error) {
	return _ContractPaymentVault.Contract.OnDemandPayments(&_ContractPaymentVault.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultSession) Owner() (common.Address, error) {
	return _ContractPaymentVault.Contract.Owner(&_ContractPaymentVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) Owner() (common.Address, error) {
	return _ContractPaymentVault.Contract.Owner(&_ContractPaymentVault.CallOpts)
}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PricePerSymbol(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "pricePerSymbol")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) PricePerSymbol() (uint64, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PricePerSymbol is a free data retrieval call binding the contract method 0xf323726a.
//
// Solidity: function pricePerSymbol() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PricePerSymbol() (uint64, error) {
	return _ContractPaymentVault.Contract.PricePerSymbol(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) PriceUpdateCooldown(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "priceUpdateCooldown")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) PriceUpdateCooldown() (uint64, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// PriceUpdateCooldown is a free data retrieval call binding the contract method 0x039f091c.
//
// Solidity: function priceUpdateCooldown() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) PriceUpdateCooldown() (uint64, error) {
	return _ContractPaymentVault.Contract.PriceUpdateCooldown(&_ContractPaymentVault.CallOpts)
}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCaller) ReservationPeriodInterval(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservationPeriodInterval")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultSession) ReservationPeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.ReservationPeriodInterval(&_ContractPaymentVault.CallOpts)
}

// ReservationPeriodInterval is a free data retrieval call binding the contract method 0x72228ab2.
//
// Solidity: function reservationPeriodInterval() view returns(uint64)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) ReservationPeriodInterval() (uint64, error) {
	return _ContractPaymentVault.Contract.ReservationPeriodInterval(&_ContractPaymentVault.CallOpts)
}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCaller) Reservations(opts *bind.CallOpts, arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	var out []interface{}
	err := _ContractPaymentVault.contract.Call(opts, &out, "reservations", arg0)

	outstruct := new(struct {
		SymbolsPerSecond uint64
		StartTimestamp   uint64
		EndTimestamp     uint64
		QuorumNumbers    []byte
		QuorumSplits     []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SymbolsPerSecond = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.StartTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.EndTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.QuorumNumbers = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.QuorumSplits = *abi.ConvertType(out[4], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultSession) Reservations(arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	return _ContractPaymentVault.Contract.Reservations(&_ContractPaymentVault.CallOpts, arg0)
}

// Reservations is a free data retrieval call binding the contract method 0xfd3dc53a.
//
// Solidity: function reservations(address ) view returns(uint64 symbolsPerSecond, uint64 startTimestamp, uint64 endTimestamp, bytes quorumNumbers, bytes quorumSplits)
func (_ContractPaymentVault *ContractPaymentVaultCallerSession) Reservations(arg0 common.Address) (struct {
	SymbolsPerSecond uint64
	StartTimestamp   uint64
	EndTimestamp     uint64
	QuorumNumbers    []byte
	QuorumSplits     []byte
}, error) {
	return _ContractPaymentVault.Contract.Reservations(&_ContractPaymentVault.CallOpts, arg0)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) DepositOnDemand(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "depositOnDemand", _account)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) DepositOnDemand(_account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, _account)
}

// DepositOnDemand is a paid mutator transaction binding the contract method 0x8bec7d02.
//
// Solidity: function depositOnDemand(address _account) payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) DepositOnDemand(_account common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.DepositOnDemand(&_ContractPaymentVault.TransactOpts, _account)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "initialize", _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Initialize(_initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a1bbf37.
//
// Solidity: function initialize(address _initialOwner, uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown, uint64 _globalSymbolsPerPeriod, uint64 _reservationPeriodInterval, uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Initialize(_initialOwner common.Address, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64, _globalSymbolsPerPeriod uint64, _reservationPeriodInterval uint64, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Initialize(&_ContractPaymentVault.TransactOpts, _initialOwner, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown, _globalSymbolsPerPeriod, _reservationPeriodInterval, _globalRatePeriodInterval)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.RenounceOwnership(&_ContractPaymentVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.RenounceOwnership(&_ContractPaymentVault.TransactOpts)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalRatePeriodInterval(opts *bind.TransactOpts, _globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalRatePeriodInterval", _globalRatePeriodInterval)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalRatePeriodInterval(_globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRatePeriodInterval(&_ContractPaymentVault.TransactOpts, _globalRatePeriodInterval)
}

// SetGlobalRatePeriodInterval is a paid mutator transaction binding the contract method 0xaa788bd7.
//
// Solidity: function setGlobalRatePeriodInterval(uint64 _globalRatePeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalRatePeriodInterval(_globalRatePeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalRatePeriodInterval(&_ContractPaymentVault.TransactOpts, _globalRatePeriodInterval)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetGlobalSymbolsPerPeriod(opts *bind.TransactOpts, _globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setGlobalSymbolsPerPeriod", _globalSymbolsPerPeriod)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetGlobalSymbolsPerPeriod(_globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerPeriod(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerPeriod)
}

// SetGlobalSymbolsPerPeriod is a paid mutator transaction binding the contract method 0xa16cf884.
//
// Solidity: function setGlobalSymbolsPerPeriod(uint64 _globalSymbolsPerPeriod) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetGlobalSymbolsPerPeriod(_globalSymbolsPerPeriod uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetGlobalSymbolsPerPeriod(&_ContractPaymentVault.TransactOpts, _globalSymbolsPerPeriod)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetPriceParams(opts *bind.TransactOpts, _minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setPriceParams", _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetPriceParams(_minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetPriceParams is a paid mutator transaction binding the contract method 0xfba2b1d1.
//
// Solidity: function setPriceParams(uint64 _minNumSymbols, uint64 _pricePerSymbol, uint64 _priceUpdateCooldown) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetPriceParams(_minNumSymbols uint64, _pricePerSymbol uint64, _priceUpdateCooldown uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetPriceParams(&_ContractPaymentVault.TransactOpts, _minNumSymbols, _pricePerSymbol, _priceUpdateCooldown)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservation(opts *bind.TransactOpts, _account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservation", _account, _reservation)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservation(_account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservation(&_ContractPaymentVault.TransactOpts, _account, _reservation)
}

// SetReservation is a paid mutator transaction binding the contract method 0x9aec8640.
//
// Solidity: function setReservation(address _account, (uint64,uint64,uint64,bytes,bytes) _reservation) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservation(_account common.Address, _reservation IPaymentVaultReservation) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservation(&_ContractPaymentVault.TransactOpts, _account, _reservation)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) SetReservationPeriodInterval(opts *bind.TransactOpts, _reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "setReservationPeriodInterval", _reservationPeriodInterval)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) SetReservationPeriodInterval(_reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationPeriodInterval(&_ContractPaymentVault.TransactOpts, _reservationPeriodInterval)
}

// SetReservationPeriodInterval is a paid mutator transaction binding the contract method 0x897218fc.
//
// Solidity: function setReservationPeriodInterval(uint64 _reservationPeriodInterval) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) SetReservationPeriodInterval(_reservationPeriodInterval uint64) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.SetReservationPeriodInterval(&_ContractPaymentVault.TransactOpts, _reservationPeriodInterval)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.TransferOwnership(&_ContractPaymentVault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Withdraw(&_ContractPaymentVault.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Withdraw(&_ContractPaymentVault.TransactOpts, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) WithdrawERC20(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.Transact(opts, "withdrawERC20", _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.WithdrawERC20(&_ContractPaymentVault.TransactOpts, _token, _amount)
}

// WithdrawERC20 is a paid mutator transaction binding the contract method 0xa1db9782.
//
// Solidity: function withdrawERC20(address _token, uint256 _amount) returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) WithdrawERC20(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.WithdrawERC20(&_ContractPaymentVault.TransactOpts, _token, _amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Fallback(&_ContractPaymentVault.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Fallback(&_ContractPaymentVault.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractPaymentVault.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultSession) Receive() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Receive(&_ContractPaymentVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ContractPaymentVault *ContractPaymentVaultTransactorSession) Receive() (*types.Transaction, error) {
	return _ContractPaymentVault.Contract.Receive(&_ContractPaymentVault.TransactOpts)
}

// ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator is returned from FilterGlobalRatePeriodIntervalUpdated and is used to iterate over the raw logs and unpacked data for GlobalRatePeriodIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalRatePeriodIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
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
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalRatePeriodIntervalUpdated represents a GlobalRatePeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalRatePeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalRatePeriodIntervalUpdated is a free log retrieval operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalRatePeriodIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalRatePeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalRatePeriodIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalRatePeriodIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalRatePeriodIntervalUpdated is a free log subscription operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalRatePeriodIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalRatePeriodIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalRatePeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRatePeriodIntervalUpdated", log); err != nil {
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

// ParseGlobalRatePeriodIntervalUpdated is a log parse operation binding the contract event 0x833819c38214ef9f462f88b5c27a21bf201f394572a14da3e63c77ee15f0e93a.
//
// Solidity: event GlobalRatePeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalRatePeriodIntervalUpdated(log types.Log) (*ContractPaymentVaultGlobalRatePeriodIntervalUpdated, error) {
	event := new(ContractPaymentVaultGlobalRatePeriodIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalRatePeriodIntervalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator is returned from FilterGlobalSymbolsPerPeriodUpdated and is used to iterate over the raw logs and unpacked data for GlobalSymbolsPerPeriodUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator struct {
	Event *ContractPaymentVaultGlobalSymbolsPerPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
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
		it.Event = new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
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
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultGlobalSymbolsPerPeriodUpdated represents a GlobalSymbolsPerPeriodUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultGlobalSymbolsPerPeriodUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGlobalSymbolsPerPeriodUpdated is a free log retrieval operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterGlobalSymbolsPerPeriodUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "GlobalSymbolsPerPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultGlobalSymbolsPerPeriodUpdatedIterator{contract: _ContractPaymentVault.contract, event: "GlobalSymbolsPerPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalSymbolsPerPeriodUpdated is a free log subscription operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchGlobalSymbolsPerPeriodUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultGlobalSymbolsPerPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "GlobalSymbolsPerPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerPeriodUpdated", log); err != nil {
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

// ParseGlobalSymbolsPerPeriodUpdated is a log parse operation binding the contract event 0x3edf3b79e74d9e583ff51df95fbabefe15f504d33475b2cc77cffba292268aae.
//
// Solidity: event GlobalSymbolsPerPeriodUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseGlobalSymbolsPerPeriodUpdated(log types.Log) (*ContractPaymentVaultGlobalSymbolsPerPeriodUpdated, error) {
	event := new(ContractPaymentVaultGlobalSymbolsPerPeriodUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "GlobalSymbolsPerPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractPaymentVault contract.
type ContractPaymentVaultInitializedIterator struct {
	Event *ContractPaymentVaultInitialized // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultInitialized)
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
		it.Event = new(ContractPaymentVaultInitialized)
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
func (it *ContractPaymentVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultInitialized represents a Initialized event raised by the ContractPaymentVault contract.
type ContractPaymentVaultInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractPaymentVaultInitializedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultInitializedIterator{contract: _ContractPaymentVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultInitialized)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseInitialized(log types.Log) (*ContractPaymentVaultInitialized, error) {
	event := new(ContractPaymentVaultInitialized)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultOnDemandPaymentUpdatedIterator is returned from FilterOnDemandPaymentUpdated and is used to iterate over the raw logs and unpacked data for OnDemandPaymentUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultOnDemandPaymentUpdatedIterator struct {
	Event *ContractPaymentVaultOnDemandPaymentUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultOnDemandPaymentUpdated)
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
		it.Event = new(ContractPaymentVaultOnDemandPaymentUpdated)
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
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultOnDemandPaymentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultOnDemandPaymentUpdated represents a OnDemandPaymentUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOnDemandPaymentUpdated struct {
	Account         common.Address
	OnDemandPayment *big.Int
	TotalDeposit    *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOnDemandPaymentUpdated is a free log retrieval operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterOnDemandPaymentUpdated(opts *bind.FilterOpts, account []common.Address) (*ContractPaymentVaultOnDemandPaymentUpdatedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "OnDemandPaymentUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultOnDemandPaymentUpdatedIterator{contract: _ContractPaymentVault.contract, event: "OnDemandPaymentUpdated", logs: logs, sub: sub}, nil
}

// WatchOnDemandPaymentUpdated is a free log subscription operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchOnDemandPaymentUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultOnDemandPaymentUpdated, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "OnDemandPaymentUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultOnDemandPaymentUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "OnDemandPaymentUpdated", log); err != nil {
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

// ParseOnDemandPaymentUpdated is a log parse operation binding the contract event 0x6fbb447a2c09b8901d70b0d5b9fbce159ee8fda4460e5af2570cab3fe0adf268.
//
// Solidity: event OnDemandPaymentUpdated(address indexed account, uint80 onDemandPayment, uint80 totalDeposit)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseOnDemandPaymentUpdated(log types.Log) (*ContractPaymentVaultOnDemandPaymentUpdated, error) {
	event := new(ContractPaymentVaultOnDemandPaymentUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "OnDemandPaymentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractPaymentVault contract.
type ContractPaymentVaultOwnershipTransferredIterator struct {
	Event *ContractPaymentVaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultOwnershipTransferred)
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
		it.Event = new(ContractPaymentVaultOwnershipTransferred)
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
func (it *ContractPaymentVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultOwnershipTransferred represents a OwnershipTransferred event raised by the ContractPaymentVault contract.
type ContractPaymentVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractPaymentVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultOwnershipTransferredIterator{contract: _ContractPaymentVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultOwnershipTransferred)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseOwnershipTransferred(log types.Log) (*ContractPaymentVaultOwnershipTransferred, error) {
	event := new(ContractPaymentVaultOwnershipTransferred)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultPriceParamsUpdatedIterator is returned from FilterPriceParamsUpdated and is used to iterate over the raw logs and unpacked data for PriceParamsUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultPriceParamsUpdatedIterator struct {
	Event *ContractPaymentVaultPriceParamsUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultPriceParamsUpdated)
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
		it.Event = new(ContractPaymentVaultPriceParamsUpdated)
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
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultPriceParamsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultPriceParamsUpdated represents a PriceParamsUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultPriceParamsUpdated struct {
	PreviousMinNumSymbols       uint64
	NewMinNumSymbols            uint64
	PreviousPricePerSymbol      uint64
	NewPricePerSymbol           uint64
	PreviousPriceUpdateCooldown uint64
	NewPriceUpdateCooldown      uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterPriceParamsUpdated is a free log retrieval operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterPriceParamsUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultPriceParamsUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "PriceParamsUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultPriceParamsUpdatedIterator{contract: _ContractPaymentVault.contract, event: "PriceParamsUpdated", logs: logs, sub: sub}, nil
}

// WatchPriceParamsUpdated is a free log subscription operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchPriceParamsUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultPriceParamsUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "PriceParamsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultPriceParamsUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "PriceParamsUpdated", log); err != nil {
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

// ParsePriceParamsUpdated is a log parse operation binding the contract event 0x9b97ed982ea5820e21bfc9578505e78068a5333487583460ad56ff72defef77a.
//
// Solidity: event PriceParamsUpdated(uint64 previousMinNumSymbols, uint64 newMinNumSymbols, uint64 previousPricePerSymbol, uint64 newPricePerSymbol, uint64 previousPriceUpdateCooldown, uint64 newPriceUpdateCooldown)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParsePriceParamsUpdated(log types.Log) (*ContractPaymentVaultPriceParamsUpdated, error) {
	event := new(ContractPaymentVaultPriceParamsUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "PriceParamsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultReservationPeriodIntervalUpdatedIterator is returned from FilterReservationPeriodIntervalUpdated and is used to iterate over the raw logs and unpacked data for ReservationPeriodIntervalUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationPeriodIntervalUpdatedIterator struct {
	Event *ContractPaymentVaultReservationPeriodIntervalUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultReservationPeriodIntervalUpdated)
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
		it.Event = new(ContractPaymentVaultReservationPeriodIntervalUpdated)
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
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultReservationPeriodIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultReservationPeriodIntervalUpdated represents a ReservationPeriodIntervalUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationPeriodIntervalUpdated struct {
	PreviousValue uint64
	NewValue      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterReservationPeriodIntervalUpdated is a free log retrieval operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterReservationPeriodIntervalUpdated(opts *bind.FilterOpts) (*ContractPaymentVaultReservationPeriodIntervalUpdatedIterator, error) {

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "ReservationPeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultReservationPeriodIntervalUpdatedIterator{contract: _ContractPaymentVault.contract, event: "ReservationPeriodIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchReservationPeriodIntervalUpdated is a free log subscription operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchReservationPeriodIntervalUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultReservationPeriodIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "ReservationPeriodIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultReservationPeriodIntervalUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationPeriodIntervalUpdated", log); err != nil {
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

// ParseReservationPeriodIntervalUpdated is a log parse operation binding the contract event 0x1ef4a1ce7d8e50959d15578b346bb20a5b049e5ee1978014a4ba66476265c957.
//
// Solidity: event ReservationPeriodIntervalUpdated(uint64 previousValue, uint64 newValue)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseReservationPeriodIntervalUpdated(log types.Log) (*ContractPaymentVaultReservationPeriodIntervalUpdated, error) {
	event := new(ContractPaymentVaultReservationPeriodIntervalUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationPeriodIntervalUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractPaymentVaultReservationUpdatedIterator is returned from FilterReservationUpdated and is used to iterate over the raw logs and unpacked data for ReservationUpdated events raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationUpdatedIterator struct {
	Event *ContractPaymentVaultReservationUpdated // Event containing the contract specifics and raw log

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
func (it *ContractPaymentVaultReservationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPaymentVaultReservationUpdated)
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
		it.Event = new(ContractPaymentVaultReservationUpdated)
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
func (it *ContractPaymentVaultReservationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPaymentVaultReservationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPaymentVaultReservationUpdated represents a ReservationUpdated event raised by the ContractPaymentVault contract.
type ContractPaymentVaultReservationUpdated struct {
	Account     common.Address
	Reservation IPaymentVaultReservation
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReservationUpdated is a free log retrieval operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) FilterReservationUpdated(opts *bind.FilterOpts, account []common.Address) (*ContractPaymentVaultReservationUpdatedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.FilterLogs(opts, "ReservationUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractPaymentVaultReservationUpdatedIterator{contract: _ContractPaymentVault.contract, event: "ReservationUpdated", logs: logs, sub: sub}, nil
}

// WatchReservationUpdated is a free log subscription operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) WatchReservationUpdated(opts *bind.WatchOpts, sink chan<- *ContractPaymentVaultReservationUpdated, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractPaymentVault.contract.WatchLogs(opts, "ReservationUpdated", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPaymentVaultReservationUpdated)
				if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationUpdated", log); err != nil {
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

// ParseReservationUpdated is a log parse operation binding the contract event 0xff3054d138559c39b4c0826c43e94b2b2c6bc9a33ea1d0b74f16c916c7b73ec1.
//
// Solidity: event ReservationUpdated(address indexed account, (uint64,uint64,uint64,bytes,bytes) reservation)
func (_ContractPaymentVault *ContractPaymentVaultFilterer) ParseReservationUpdated(log types.Log) (*ContractPaymentVaultReservationUpdated, error) {
	event := new(ContractPaymentVaultReservationUpdated)
	if err := _ContractPaymentVault.contract.UnpackLog(event, "ReservationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
