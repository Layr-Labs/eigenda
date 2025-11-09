// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAProofValidator

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

// ContractEigenDAProofValidatorMetaData contains all meta data concerning the ContractEigenDAProofValidator contract.
var ContractEigenDAProofValidatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDACertVerifierRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validateCertificate\",\"inputs\":[{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateReadPreimage\",\"inputs\":[{\"name\":\"certHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"preimageChunk\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"}]",
}

// ContractEigenDAProofValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAProofValidatorMetaData.ABI instead.
var ContractEigenDAProofValidatorABI = ContractEigenDAProofValidatorMetaData.ABI

// ContractEigenDAProofValidator is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAProofValidator struct {
	ContractEigenDAProofValidatorCaller     // Read-only binding to the contract
	ContractEigenDAProofValidatorTransactor // Write-only binding to the contract
	ContractEigenDAProofValidatorFilterer   // Log filterer for contract events
}

// ContractEigenDAProofValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAProofValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAProofValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAProofValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAProofValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAProofValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAProofValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAProofValidatorSession struct {
	Contract     *ContractEigenDAProofValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDAProofValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAProofValidatorCallerSession struct {
	Contract *ContractEigenDAProofValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDAProofValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAProofValidatorTransactorSession struct {
	Contract     *ContractEigenDAProofValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDAProofValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAProofValidatorRaw struct {
	Contract *ContractEigenDAProofValidator // Generic contract binding to access the raw methods on
}

// ContractEigenDAProofValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAProofValidatorCallerRaw struct {
	Contract *ContractEigenDAProofValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAProofValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAProofValidatorTransactorRaw struct {
	Contract *ContractEigenDAProofValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAProofValidator creates a new instance of ContractEigenDAProofValidator, bound to a specific deployed contract.
func NewContractEigenDAProofValidator(address common.Address, backend bind.ContractBackend) (*ContractEigenDAProofValidator, error) {
	contract, err := bindContractEigenDAProofValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAProofValidator{ContractEigenDAProofValidatorCaller: ContractEigenDAProofValidatorCaller{contract: contract}, ContractEigenDAProofValidatorTransactor: ContractEigenDAProofValidatorTransactor{contract: contract}, ContractEigenDAProofValidatorFilterer: ContractEigenDAProofValidatorFilterer{contract: contract}}, nil
}

// NewContractEigenDAProofValidatorCaller creates a new read-only instance of ContractEigenDAProofValidator, bound to a specific deployed contract.
func NewContractEigenDAProofValidatorCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAProofValidatorCaller, error) {
	contract, err := bindContractEigenDAProofValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAProofValidatorCaller{contract: contract}, nil
}

// NewContractEigenDAProofValidatorTransactor creates a new write-only instance of ContractEigenDAProofValidator, bound to a specific deployed contract.
func NewContractEigenDAProofValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAProofValidatorTransactor, error) {
	contract, err := bindContractEigenDAProofValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAProofValidatorTransactor{contract: contract}, nil
}

// NewContractEigenDAProofValidatorFilterer creates a new log filterer instance of ContractEigenDAProofValidator, bound to a specific deployed contract.
func NewContractEigenDAProofValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAProofValidatorFilterer, error) {
	contract, err := bindContractEigenDAProofValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAProofValidatorFilterer{contract: contract}, nil
}

// bindContractEigenDAProofValidator binds a generic wrapper to an already deployed contract.
func bindContractEigenDAProofValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAProofValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAProofValidator.Contract.ContractEigenDAProofValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAProofValidator.Contract.ContractEigenDAProofValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAProofValidator.Contract.ContractEigenDAProofValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAProofValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAProofValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAProofValidator.Contract.contract.Transact(opts, method, params...)
}

// ValidateCertificate is a free data retrieval call binding the contract method 0xe667d8aa.
//
// Solidity: function validateCertificate(bytes proof) view returns(bool isValid)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorCaller) ValidateCertificate(opts *bind.CallOpts, proof []byte) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAProofValidator.contract.Call(opts, &out, "validateCertificate", proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidateCertificate is a free data retrieval call binding the contract method 0xe667d8aa.
//
// Solidity: function validateCertificate(bytes proof) view returns(bool isValid)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorSession) ValidateCertificate(proof []byte) (bool, error) {
	return _ContractEigenDAProofValidator.Contract.ValidateCertificate(&_ContractEigenDAProofValidator.CallOpts, proof)
}

// ValidateCertificate is a free data retrieval call binding the contract method 0xe667d8aa.
//
// Solidity: function validateCertificate(bytes proof) view returns(bool isValid)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorCallerSession) ValidateCertificate(proof []byte) (bool, error) {
	return _ContractEigenDAProofValidator.Contract.ValidateCertificate(&_ContractEigenDAProofValidator.CallOpts, proof)
}

// ValidateReadPreimage is a free data retrieval call binding the contract method 0x08273c74.
//
// Solidity: function validateReadPreimage(bytes32 certHash, uint256 offset, bytes proof) pure returns(bytes preimageChunk)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorCaller) ValidateReadPreimage(opts *bind.CallOpts, certHash [32]byte, offset *big.Int, proof []byte) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAProofValidator.contract.Call(opts, &out, "validateReadPreimage", certHash, offset, proof)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ValidateReadPreimage is a free data retrieval call binding the contract method 0x08273c74.
//
// Solidity: function validateReadPreimage(bytes32 certHash, uint256 offset, bytes proof) pure returns(bytes preimageChunk)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorSession) ValidateReadPreimage(certHash [32]byte, offset *big.Int, proof []byte) ([]byte, error) {
	return _ContractEigenDAProofValidator.Contract.ValidateReadPreimage(&_ContractEigenDAProofValidator.CallOpts, certHash, offset, proof)
}

// ValidateReadPreimage is a free data retrieval call binding the contract method 0x08273c74.
//
// Solidity: function validateReadPreimage(bytes32 certHash, uint256 offset, bytes proof) pure returns(bytes preimageChunk)
func (_ContractEigenDAProofValidator *ContractEigenDAProofValidatorCallerSession) ValidateReadPreimage(certHash [32]byte, offset *big.Int, proof []byte) ([]byte, error) {
	return _ContractEigenDAProofValidator.Contract.ValidateReadPreimage(&_ContractEigenDAProofValidator.CallOpts, certHash, offset, proof)
}
