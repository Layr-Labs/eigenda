// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierRouter

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

// ContractEigenDACertVerifierRouterMetaData contains all meta data concerning the ContractEigenDACertVerifierRouter contract.
var ContractEigenDACertVerifierRouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addCertVerifier\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"certVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVerifierABNs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"certVerifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCertVerifierAt\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initRBNs\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"initCertVerifiers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CertVerifierAdded\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"certVerifier\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ABNNotGreaterThanLast\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"ABNNotInFuture\",\"inputs\":[{\"name\":\"activationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidCertLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RBNInFuture\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
}

// ContractEigenDACertVerifierRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierRouterMetaData.ABI instead.
var ContractEigenDACertVerifierRouterABI = ContractEigenDACertVerifierRouterMetaData.ABI

// ContractEigenDACertVerifierRouter is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouter struct {
	ContractEigenDACertVerifierRouterCaller     // Read-only binding to the contract
	ContractEigenDACertVerifierRouterTransactor // Write-only binding to the contract
	ContractEigenDACertVerifierRouterFilterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierRouterSession struct {
	Contract     *ContractEigenDACertVerifierRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                      // Call options to use throughout this session
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierRouterCallerSession struct {
	Contract *ContractEigenDACertVerifierRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                            // Call options to use throughout this session
}

// ContractEigenDACertVerifierRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierRouterTransactorSession struct {
	Contract     *ContractEigenDACertVerifierRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                            // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouterRaw struct {
	Contract *ContractEigenDACertVerifierRouter // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouterCallerRaw struct {
	Contract *ContractEigenDACertVerifierRouterCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRouterTransactorRaw struct {
	Contract *ContractEigenDACertVerifierRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierRouter creates a new instance of ContractEigenDACertVerifierRouter, bound to a specific deployed contract.
func NewContractEigenDACertVerifierRouter(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierRouter, error) {
	contract, err := bindContractEigenDACertVerifierRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouter{ContractEigenDACertVerifierRouterCaller: ContractEigenDACertVerifierRouterCaller{contract: contract}, ContractEigenDACertVerifierRouterTransactor: ContractEigenDACertVerifierRouterTransactor{contract: contract}, ContractEigenDACertVerifierRouterFilterer: ContractEigenDACertVerifierRouterFilterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierRouterCaller creates a new read-only instance of ContractEigenDACertVerifierRouter, bound to a specific deployed contract.
func NewContractEigenDACertVerifierRouterCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierRouterCaller, error) {
	contract, err := bindContractEigenDACertVerifierRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterCaller{contract: contract}, nil
}

// NewContractEigenDACertVerifierRouterTransactor creates a new write-only instance of ContractEigenDACertVerifierRouter, bound to a specific deployed contract.
func NewContractEigenDACertVerifierRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierRouterTransactor, error) {
	contract, err := bindContractEigenDACertVerifierRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterTransactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierRouterFilterer creates a new log filterer instance of ContractEigenDACertVerifierRouter, bound to a specific deployed contract.
func NewContractEigenDACertVerifierRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierRouterFilterer, error) {
	contract, err := bindContractEigenDACertVerifierRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterFilterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierRouter binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierRouter.Contract.ContractEigenDACertVerifierRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.ContractEigenDACertVerifierRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.ContractEigenDACertVerifierRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.contract.Transact(opts, method, params...)
}

// CertVerifierABNs is a free data retrieval call binding the contract method 0xf0df66df.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCaller) CertVerifierABNs(opts *bind.CallOpts, arg0 *big.Int) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierRouter.contract.Call(opts, &out, "certVerifierABNs", arg0)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// CertVerifierABNs is a free data retrieval call binding the contract method 0xf0df66df.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) CertVerifierABNs(arg0 *big.Int) (uint32, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CertVerifierABNs(&_ContractEigenDACertVerifierRouter.CallOpts, arg0)
}

// CertVerifierABNs is a free data retrieval call binding the contract method 0xf0df66df.
//
// Solidity: function certVerifierABNs(uint256 ) view returns(uint32)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerSession) CertVerifierABNs(arg0 *big.Int) (uint32, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CertVerifierABNs(&_ContractEigenDACertVerifierRouter.CallOpts, arg0)
}

// CertVerifiers is a free data retrieval call binding the contract method 0x4c046566.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCaller) CertVerifiers(opts *bind.CallOpts, arg0 uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierRouter.contract.Call(opts, &out, "certVerifiers", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CertVerifiers is a free data retrieval call binding the contract method 0x4c046566.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) CertVerifiers(arg0 uint32) (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CertVerifiers(&_ContractEigenDACertVerifierRouter.CallOpts, arg0)
}

// CertVerifiers is a free data retrieval call binding the contract method 0x4c046566.
//
// Solidity: function certVerifiers(uint32 ) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerSession) CertVerifiers(arg0 uint32) (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CertVerifiers(&_ContractEigenDACertVerifierRouter.CallOpts, arg0)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCaller) CheckDACert(opts *bind.CallOpts, abiEncodedCert []byte) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierRouter.contract.Call(opts, &out, "checkDACert", abiEncodedCert)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CheckDACert(&_ContractEigenDACertVerifierRouter.CallOpts, abiEncodedCert)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifierRouter.Contract.CheckDACert(&_ContractEigenDACertVerifierRouter.CallOpts, abiEncodedCert)
}

// GetCertVerifierAt is a free data retrieval call binding the contract method 0x4a4ae0e2.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCaller) GetCertVerifierAt(opts *bind.CallOpts, referenceBlockNumber uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierRouter.contract.Call(opts, &out, "getCertVerifierAt", referenceBlockNumber)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetCertVerifierAt is a free data retrieval call binding the contract method 0x4a4ae0e2.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) GetCertVerifierAt(referenceBlockNumber uint32) (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.GetCertVerifierAt(&_ContractEigenDACertVerifierRouter.CallOpts, referenceBlockNumber)
}

// GetCertVerifierAt is a free data retrieval call binding the contract method 0x4a4ae0e2.
//
// Solidity: function getCertVerifierAt(uint32 referenceBlockNumber) view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerSession) GetCertVerifierAt(referenceBlockNumber uint32) (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.GetCertVerifierAt(&_ContractEigenDACertVerifierRouter.CallOpts, referenceBlockNumber)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) Owner() (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.Owner(&_ContractEigenDACertVerifierRouter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDACertVerifierRouter.Contract.Owner(&_ContractEigenDACertVerifierRouter.CallOpts)
}

// AddCertVerifier is a paid mutator transaction binding the contract method 0xbfda00de.
//
// Solidity: function addCertVerifier(uint32 activationBlockNumber, address certVerifier) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactor) AddCertVerifier(opts *bind.TransactOpts, activationBlockNumber uint32, certVerifier common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.contract.Transact(opts, "addCertVerifier", activationBlockNumber, certVerifier)
}

// AddCertVerifier is a paid mutator transaction binding the contract method 0xbfda00de.
//
// Solidity: function addCertVerifier(uint32 activationBlockNumber, address certVerifier) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) AddCertVerifier(activationBlockNumber uint32, certVerifier common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.AddCertVerifier(&_ContractEigenDACertVerifierRouter.TransactOpts, activationBlockNumber, certVerifier)
}

// AddCertVerifier is a paid mutator transaction binding the contract method 0xbfda00de.
//
// Solidity: function addCertVerifier(uint32 activationBlockNumber, address certVerifier) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorSession) AddCertVerifier(activationBlockNumber uint32, certVerifier common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.AddCertVerifier(&_ContractEigenDACertVerifierRouter.TransactOpts, activationBlockNumber, certVerifier)
}

// Initialize is a paid mutator transaction binding the contract method 0x9d8ecd85.
//
// Solidity: function initialize(address _initialOwner, uint32[] initRBNs, address[] initCertVerifiers) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, initRBNs []uint32, initCertVerifiers []common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.contract.Transact(opts, "initialize", _initialOwner, initRBNs, initCertVerifiers)
}

// Initialize is a paid mutator transaction binding the contract method 0x9d8ecd85.
//
// Solidity: function initialize(address _initialOwner, uint32[] initRBNs, address[] initCertVerifiers) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) Initialize(_initialOwner common.Address, initRBNs []uint32, initCertVerifiers []common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.Initialize(&_ContractEigenDACertVerifierRouter.TransactOpts, _initialOwner, initRBNs, initCertVerifiers)
}

// Initialize is a paid mutator transaction binding the contract method 0x9d8ecd85.
//
// Solidity: function initialize(address _initialOwner, uint32[] initRBNs, address[] initCertVerifiers) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorSession) Initialize(_initialOwner common.Address, initRBNs []uint32, initCertVerifiers []common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.Initialize(&_ContractEigenDACertVerifierRouter.TransactOpts, _initialOwner, initRBNs, initCertVerifiers)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.RenounceOwnership(&_ContractEigenDACertVerifierRouter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.RenounceOwnership(&_ContractEigenDACertVerifierRouter.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.TransferOwnership(&_ContractEigenDACertVerifierRouter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierRouter.Contract.TransferOwnership(&_ContractEigenDACertVerifierRouter.TransactOpts, newOwner)
}

// ContractEigenDACertVerifierRouterCertVerifierAddedIterator is returned from FilterCertVerifierAdded and is used to iterate over the raw logs and unpacked data for CertVerifierAdded events raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterCertVerifierAddedIterator struct {
	Event *ContractEigenDACertVerifierRouterCertVerifierAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDACertVerifierRouterCertVerifierAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierRouterCertVerifierAdded)
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
		it.Event = new(ContractEigenDACertVerifierRouterCertVerifierAdded)
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
func (it *ContractEigenDACertVerifierRouterCertVerifierAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierRouterCertVerifierAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierRouterCertVerifierAdded represents a CertVerifierAdded event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterCertVerifierAdded struct {
	ActivationBlockNumber uint32
	CertVerifier          common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterCertVerifierAdded is a free log retrieval operation binding the contract event 0x3c87ded09f10478b3e4c40df4329a85dc74ce5f77d000d69a438e6af6096b0e2.
//
// Solidity: event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) FilterCertVerifierAdded(opts *bind.FilterOpts, activationBlockNumber []uint32, certVerifier []common.Address) (*ContractEigenDACertVerifierRouterCertVerifierAddedIterator, error) {

	var activationBlockNumberRule []interface{}
	for _, activationBlockNumberItem := range activationBlockNumber {
		activationBlockNumberRule = append(activationBlockNumberRule, activationBlockNumberItem)
	}
	var certVerifierRule []interface{}
	for _, certVerifierItem := range certVerifier {
		certVerifierRule = append(certVerifierRule, certVerifierItem)
	}

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.FilterLogs(opts, "CertVerifierAdded", activationBlockNumberRule, certVerifierRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterCertVerifierAddedIterator{contract: _ContractEigenDACertVerifierRouter.contract, event: "CertVerifierAdded", logs: logs, sub: sub}, nil
}

// WatchCertVerifierAdded is a free log subscription operation binding the contract event 0x3c87ded09f10478b3e4c40df4329a85dc74ce5f77d000d69a438e6af6096b0e2.
//
// Solidity: event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) WatchCertVerifierAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierRouterCertVerifierAdded, activationBlockNumber []uint32, certVerifier []common.Address) (event.Subscription, error) {

	var activationBlockNumberRule []interface{}
	for _, activationBlockNumberItem := range activationBlockNumber {
		activationBlockNumberRule = append(activationBlockNumberRule, activationBlockNumberItem)
	}
	var certVerifierRule []interface{}
	for _, certVerifierItem := range certVerifier {
		certVerifierRule = append(certVerifierRule, certVerifierItem)
	}

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.WatchLogs(opts, "CertVerifierAdded", activationBlockNumberRule, certVerifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierRouterCertVerifierAdded)
				if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "CertVerifierAdded", log); err != nil {
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

// ParseCertVerifierAdded is a log parse operation binding the contract event 0x3c87ded09f10478b3e4c40df4329a85dc74ce5f77d000d69a438e6af6096b0e2.
//
// Solidity: event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) ParseCertVerifierAdded(log types.Log) (*ContractEigenDACertVerifierRouterCertVerifierAdded, error) {
	event := new(ContractEigenDACertVerifierRouterCertVerifierAdded)
	if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "CertVerifierAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierRouterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterInitializedIterator struct {
	Event *ContractEigenDACertVerifierRouterInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDACertVerifierRouterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierRouterInitialized)
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
		it.Event = new(ContractEigenDACertVerifierRouterInitialized)
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
func (it *ContractEigenDACertVerifierRouterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierRouterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierRouterInitialized represents a Initialized event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDACertVerifierRouterInitializedIterator, error) {

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterInitializedIterator{contract: _ContractEigenDACertVerifierRouter.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierRouterInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierRouterInitialized)
				if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) ParseInitialized(log types.Log) (*ContractEigenDACertVerifierRouterInitialized, error) {
	event := new(ContractEigenDACertVerifierRouterInitialized)
	if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDACertVerifierRouterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterOwnershipTransferredIterator struct {
	Event *ContractEigenDACertVerifierRouterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDACertVerifierRouterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDACertVerifierRouterOwnershipTransferred)
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
		it.Event = new(ContractEigenDACertVerifierRouterOwnershipTransferred)
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
func (it *ContractEigenDACertVerifierRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDACertVerifierRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDACertVerifierRouterOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDACertVerifierRouter contract.
type ContractEigenDACertVerifierRouterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDACertVerifierRouterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierRouterOwnershipTransferredIterator{contract: _ContractEigenDACertVerifierRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDACertVerifierRouterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDACertVerifierRouter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDACertVerifierRouterOwnershipTransferred)
				if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDACertVerifierRouter *ContractEigenDACertVerifierRouterFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDACertVerifierRouterOwnershipTransferred, error) {
	event := new(ContractEigenDACertVerifierRouterOwnershipTransferred)
	if err := _ContractEigenDACertVerifierRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
