// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifier

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

// EigenDATypesV1SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// ContractEigenDACertVerifierMetaData contains all meta data concerning the ContractEigenDACertVerifier contract.
var ContractEigenDACertVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initEigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"initEigenDASignatureVerifier\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"},{\"name\":\"initSecurityThresholds\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"initQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"certVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"checkDACert\",\"inputs\":[{\"name\":\"abiEncodedCert\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDASignatureVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDASignatureVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"securityThresholds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.SecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"InvalidSecurityThresholds\",\"inputs\":[]}]",
}

// ContractEigenDACertVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierMetaData.ABI instead.
var ContractEigenDACertVerifierABI = ContractEigenDACertVerifierMetaData.ABI

// ContractEigenDACertVerifier is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifier struct {
	ContractEigenDACertVerifierCaller     // Read-only binding to the contract
	ContractEigenDACertVerifierTransactor // Write-only binding to the contract
	ContractEigenDACertVerifierFilterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierSession struct {
	Contract     *ContractEigenDACertVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierCallerSession struct {
	Contract *ContractEigenDACertVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractEigenDACertVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierTransactorSession struct {
	Contract     *ContractEigenDACertVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierRaw struct {
	Contract *ContractEigenDACertVerifier // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierCallerRaw struct {
	Contract *ContractEigenDACertVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierTransactorRaw struct {
	Contract *ContractEigenDACertVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifier creates a new instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifier(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifier, error) {
	contract, err := bindContractEigenDACertVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifier{ContractEigenDACertVerifierCaller: ContractEigenDACertVerifierCaller{contract: contract}, ContractEigenDACertVerifierTransactor: ContractEigenDACertVerifierTransactor{contract: contract}, ContractEigenDACertVerifierFilterer: ContractEigenDACertVerifierFilterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierCaller creates a new read-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierCaller, error) {
	contract, err := bindContractEigenDACertVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierCaller{contract: contract}, nil
}

// NewContractEigenDACertVerifierTransactor creates a new write-only instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierTransactor, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierTransactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierFilterer creates a new log filterer instance of ContractEigenDACertVerifier, bound to a specific deployed contract.
func NewContractEigenDACertVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierFilterer, error) {
	contract, err := bindContractEigenDACertVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierFilterer{contract: contract}, nil
}

// bindContractEigenDACertVerifier binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.ContractEigenDACertVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifier.Contract.contract.Transact(opts, method, params...)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) CertVersion(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "certVersion")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) CertVersion() (uint64, error) {
	return _ContractEigenDACertVerifier.Contract.CertVersion(&_ContractEigenDACertVerifier.CallOpts)
}

// CertVersion is a free data retrieval call binding the contract method 0x2ead0b96.
//
// Solidity: function certVersion() pure returns(uint64)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) CertVersion() (uint64, error) {
	return _ContractEigenDACertVerifier.Contract.CertVersion(&_ContractEigenDACertVerifier.CallOpts)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) CheckDACert(opts *bind.CallOpts, abiEncodedCert []byte) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "checkDACert", abiEncodedCert)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.CheckDACert(&_ContractEigenDACertVerifier.CallOpts, abiEncodedCert)
}

// CheckDACert is a free data retrieval call binding the contract method 0x9077193b.
//
// Solidity: function checkDACert(bytes abiEncodedCert) view returns(uint8)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) CheckDACert(abiEncodedCert []byte) (uint8, error) {
	return _ContractEigenDACertVerifier.Contract.CheckDACert(&_ContractEigenDACertVerifier.CallOpts, abiEncodedCert)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDASignatureVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDASignatureVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDASignatureVerifier is a free data retrieval call binding the contract method 0xefd4532b.
//
// Solidity: function eigenDASignatureVerifier() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDASignatureVerifier() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDASignatureVerifier(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDACertVerifier.Contract.EigenDAThresholdRegistry(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifier.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCaller) SecurityThresholds(opts *bind.CallOpts) (EigenDATypesV1SecurityThresholds, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifier.contract.Call(opts, &out, "securityThresholds")

	if err != nil {
		return *new(EigenDATypesV1SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1SecurityThresholds)).(*EigenDATypesV1SecurityThresholds)

	return out0, err

}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierSession) SecurityThresholds() (EigenDATypesV1SecurityThresholds, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholds(&_ContractEigenDACertVerifier.CallOpts)
}

// SecurityThresholds is a free data retrieval call binding the contract method 0x21b9b2fb.
//
// Solidity: function securityThresholds() view returns((uint8,uint8))
func (_ContractEigenDACertVerifier *ContractEigenDACertVerifierCallerSession) SecurityThresholds() (EigenDATypesV1SecurityThresholds, error) {
	return _ContractEigenDACertVerifier.Contract.SecurityThresholds(&_ContractEigenDACertVerifier.CallOpts)
}
