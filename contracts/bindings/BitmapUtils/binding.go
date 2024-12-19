// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBitmapUtils

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

// ContractBitmapUtilsMetaData contains all meta data concerning the ContractBitmapUtils contract.
var ContractBitmapUtilsMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220f92fdde73152a572f9b0656d02fd87c50f50841475fbb322c147e17a4e60a2b364736f6c634300080c0033",
}

// ContractBitmapUtilsABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBitmapUtilsMetaData.ABI instead.
var ContractBitmapUtilsABI = ContractBitmapUtilsMetaData.ABI

// ContractBitmapUtilsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBitmapUtilsMetaData.Bin instead.
var ContractBitmapUtilsBin = ContractBitmapUtilsMetaData.Bin

// DeployContractBitmapUtils deploys a new Ethereum contract, binding an instance of ContractBitmapUtils to it.
func DeployContractBitmapUtils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractBitmapUtils, error) {
	parsed, err := ContractBitmapUtilsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBitmapUtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBitmapUtils{ContractBitmapUtilsCaller: ContractBitmapUtilsCaller{contract: contract}, ContractBitmapUtilsTransactor: ContractBitmapUtilsTransactor{contract: contract}, ContractBitmapUtilsFilterer: ContractBitmapUtilsFilterer{contract: contract}}, nil
}

// ContractBitmapUtils is an auto generated Go binding around an Ethereum contract.
type ContractBitmapUtils struct {
	ContractBitmapUtilsCaller     // Read-only binding to the contract
	ContractBitmapUtilsTransactor // Write-only binding to the contract
	ContractBitmapUtilsFilterer   // Log filterer for contract events
}

// ContractBitmapUtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBitmapUtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBitmapUtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBitmapUtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBitmapUtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBitmapUtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBitmapUtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBitmapUtilsSession struct {
	Contract     *ContractBitmapUtils // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractBitmapUtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBitmapUtilsCallerSession struct {
	Contract *ContractBitmapUtilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ContractBitmapUtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBitmapUtilsTransactorSession struct {
	Contract     *ContractBitmapUtilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractBitmapUtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBitmapUtilsRaw struct {
	Contract *ContractBitmapUtils // Generic contract binding to access the raw methods on
}

// ContractBitmapUtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBitmapUtilsCallerRaw struct {
	Contract *ContractBitmapUtilsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractBitmapUtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBitmapUtilsTransactorRaw struct {
	Contract *ContractBitmapUtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBitmapUtils creates a new instance of ContractBitmapUtils, bound to a specific deployed contract.
func NewContractBitmapUtils(address common.Address, backend bind.ContractBackend) (*ContractBitmapUtils, error) {
	contract, err := bindContractBitmapUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBitmapUtils{ContractBitmapUtilsCaller: ContractBitmapUtilsCaller{contract: contract}, ContractBitmapUtilsTransactor: ContractBitmapUtilsTransactor{contract: contract}, ContractBitmapUtilsFilterer: ContractBitmapUtilsFilterer{contract: contract}}, nil
}

// NewContractBitmapUtilsCaller creates a new read-only instance of ContractBitmapUtils, bound to a specific deployed contract.
func NewContractBitmapUtilsCaller(address common.Address, caller bind.ContractCaller) (*ContractBitmapUtilsCaller, error) {
	contract, err := bindContractBitmapUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBitmapUtilsCaller{contract: contract}, nil
}

// NewContractBitmapUtilsTransactor creates a new write-only instance of ContractBitmapUtils, bound to a specific deployed contract.
func NewContractBitmapUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractBitmapUtilsTransactor, error) {
	contract, err := bindContractBitmapUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBitmapUtilsTransactor{contract: contract}, nil
}

// NewContractBitmapUtilsFilterer creates a new log filterer instance of ContractBitmapUtils, bound to a specific deployed contract.
func NewContractBitmapUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractBitmapUtilsFilterer, error) {
	contract, err := bindContractBitmapUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBitmapUtilsFilterer{contract: contract}, nil
}

// bindContractBitmapUtils binds a generic wrapper to an already deployed contract.
func bindContractBitmapUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBitmapUtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBitmapUtils *ContractBitmapUtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBitmapUtils.Contract.ContractBitmapUtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBitmapUtils *ContractBitmapUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBitmapUtils.Contract.ContractBitmapUtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBitmapUtils *ContractBitmapUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBitmapUtils.Contract.ContractBitmapUtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBitmapUtils *ContractBitmapUtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBitmapUtils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBitmapUtils *ContractBitmapUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBitmapUtils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBitmapUtils *ContractBitmapUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBitmapUtils.Contract.contract.Transact(opts, method, params...)
}
