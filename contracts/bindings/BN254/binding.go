// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractBN254

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

// ContractBN254MetaData contains all meta data concerning the ContractBN254 contract.
var ContractBN254MetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220d8e32ecdd91faecffe0e35def77376afe88a1f8652d0b73d03366bb065b22a3264736f6c634300080c0033",
}

// ContractBN254ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractBN254MetaData.ABI instead.
var ContractBN254ABI = ContractBN254MetaData.ABI

// ContractBN254Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractBN254MetaData.Bin instead.
var ContractBN254Bin = ContractBN254MetaData.Bin

// DeployContractBN254 deploys a new Ethereum contract, binding an instance of ContractBN254 to it.
func DeployContractBN254(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractBN254, error) {
	parsed, err := ContractBN254MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBN254Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractBN254{ContractBN254Caller: ContractBN254Caller{contract: contract}, ContractBN254Transactor: ContractBN254Transactor{contract: contract}, ContractBN254Filterer: ContractBN254Filterer{contract: contract}}, nil
}

// ContractBN254 is an auto generated Go binding around an Ethereum contract.
type ContractBN254 struct {
	ContractBN254Caller     // Read-only binding to the contract
	ContractBN254Transactor // Write-only binding to the contract
	ContractBN254Filterer   // Log filterer for contract events
}

// ContractBN254Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractBN254Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBN254Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractBN254Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBN254Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractBN254Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractBN254Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractBN254Session struct {
	Contract     *ContractBN254    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractBN254CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractBN254CallerSession struct {
	Contract *ContractBN254Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// ContractBN254TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractBN254TransactorSession struct {
	Contract     *ContractBN254Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ContractBN254Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractBN254Raw struct {
	Contract *ContractBN254 // Generic contract binding to access the raw methods on
}

// ContractBN254CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractBN254CallerRaw struct {
	Contract *ContractBN254Caller // Generic read-only contract binding to access the raw methods on
}

// ContractBN254TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractBN254TransactorRaw struct {
	Contract *ContractBN254Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractBN254 creates a new instance of ContractBN254, bound to a specific deployed contract.
func NewContractBN254(address common.Address, backend bind.ContractBackend) (*ContractBN254, error) {
	contract, err := bindContractBN254(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractBN254{ContractBN254Caller: ContractBN254Caller{contract: contract}, ContractBN254Transactor: ContractBN254Transactor{contract: contract}, ContractBN254Filterer: ContractBN254Filterer{contract: contract}}, nil
}

// NewContractBN254Caller creates a new read-only instance of ContractBN254, bound to a specific deployed contract.
func NewContractBN254Caller(address common.Address, caller bind.ContractCaller) (*ContractBN254Caller, error) {
	contract, err := bindContractBN254(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBN254Caller{contract: contract}, nil
}

// NewContractBN254Transactor creates a new write-only instance of ContractBN254, bound to a specific deployed contract.
func NewContractBN254Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractBN254Transactor, error) {
	contract, err := bindContractBN254(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractBN254Transactor{contract: contract}, nil
}

// NewContractBN254Filterer creates a new log filterer instance of ContractBN254, bound to a specific deployed contract.
func NewContractBN254Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractBN254Filterer, error) {
	contract, err := bindContractBN254(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractBN254Filterer{contract: contract}, nil
}

// bindContractBN254 binds a generic wrapper to an already deployed contract.
func bindContractBN254(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractBN254MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBN254 *ContractBN254Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBN254.Contract.ContractBN254Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBN254 *ContractBN254Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBN254.Contract.ContractBN254Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBN254 *ContractBN254Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBN254.Contract.ContractBN254Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractBN254 *ContractBN254CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractBN254.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractBN254 *ContractBN254TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractBN254.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractBN254 *ContractBN254TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractBN254.Contract.contract.Transact(opts, method, params...)
}
