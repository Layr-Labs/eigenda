// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDADisperserRegistry

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

// EigenDATypesV2DisperserInfo is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV2DisperserInfo struct {
	DisperserAddress common.Address
}

// ContractEigenDADisperserRegistryMetaData contains all meta data concerning the ContractEigenDADisperserRegistry contract.
var ContractEigenDADisperserRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"disperserKeyToAddress\",\"inputs\":[{\"name\":\"_key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disperserKeyToInfo\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDisperserInfo\",\"inputs\":[{\"name\":\"_disperserKey\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"_disperserInfo\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV2.DisperserInfo\",\"components\":[{\"name\":\"disperserAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DisperserAdded\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"disperser\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// ContractEigenDADisperserRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDADisperserRegistryMetaData.ABI instead.
var ContractEigenDADisperserRegistryABI = ContractEigenDADisperserRegistryMetaData.ABI

// ContractEigenDADisperserRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistry struct {
	ContractEigenDADisperserRegistryCaller     // Read-only binding to the contract
	ContractEigenDADisperserRegistryTransactor // Write-only binding to the contract
	ContractEigenDADisperserRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDADisperserRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDADisperserRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADisperserRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDADisperserRegistrySession struct {
	Contract     *ContractEigenDADisperserRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractEigenDADisperserRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDADisperserRegistryCallerSession struct {
	Contract *ContractEigenDADisperserRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// ContractEigenDADisperserRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDADisperserRegistryTransactorSession struct {
	Contract     *ContractEigenDADisperserRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// ContractEigenDADisperserRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryRaw struct {
	Contract *ContractEigenDADisperserRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDADisperserRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryCallerRaw struct {
	Contract *ContractEigenDADisperserRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDADisperserRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDADisperserRegistryTransactorRaw struct {
	Contract *ContractEigenDADisperserRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDADisperserRegistry creates a new instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDADisperserRegistry, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistry{ContractEigenDADisperserRegistryCaller: ContractEigenDADisperserRegistryCaller{contract: contract}, ContractEigenDADisperserRegistryTransactor: ContractEigenDADisperserRegistryTransactor{contract: contract}, ContractEigenDADisperserRegistryFilterer: ContractEigenDADisperserRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDADisperserRegistryCaller creates a new read-only instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDADisperserRegistryCaller, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryCaller{contract: contract}, nil
}

// NewContractEigenDADisperserRegistryTransactor creates a new write-only instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDADisperserRegistryTransactor, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDADisperserRegistryFilterer creates a new log filterer instance of ContractEigenDADisperserRegistry, bound to a specific deployed contract.
func NewContractEigenDADisperserRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDADisperserRegistryFilterer, error) {
	contract, err := bindContractEigenDADisperserRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDADisperserRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDADisperserRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDADisperserRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.ContractEigenDADisperserRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADisperserRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.contract.Transact(opts, method, params...)
}

// DisperserKeyToAddress is a free data retrieval call binding the contract method 0x07d69fad.
//
// Solidity: function disperserKeyToAddress(uint32 _key) view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) DisperserKeyToAddress(opts *bind.CallOpts, _key uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "disperserKeyToAddress", _key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DisperserKeyToAddress is a free data retrieval call binding the contract method 0x07d69fad.
//
// Solidity: function disperserKeyToAddress(uint32 _key) view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) DisperserKeyToAddress(_key uint32) (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.DisperserKeyToAddress(&_ContractEigenDADisperserRegistry.CallOpts, _key)
}

// DisperserKeyToAddress is a free data retrieval call binding the contract method 0x07d69fad.
//
// Solidity: function disperserKeyToAddress(uint32 _key) view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) DisperserKeyToAddress(_key uint32) (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.DisperserKeyToAddress(&_ContractEigenDADisperserRegistry.CallOpts, _key)
}

// DisperserKeyToInfo is a free data retrieval call binding the contract method 0x1e0bf73c.
//
// Solidity: function disperserKeyToInfo(uint32 ) view returns(address disperserAddress)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) DisperserKeyToInfo(opts *bind.CallOpts, arg0 uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "disperserKeyToInfo", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DisperserKeyToInfo is a free data retrieval call binding the contract method 0x1e0bf73c.
//
// Solidity: function disperserKeyToInfo(uint32 ) view returns(address disperserAddress)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) DisperserKeyToInfo(arg0 uint32) (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.DisperserKeyToInfo(&_ContractEigenDADisperserRegistry.CallOpts, arg0)
}

// DisperserKeyToInfo is a free data retrieval call binding the contract method 0x1e0bf73c.
//
// Solidity: function disperserKeyToInfo(uint32 ) view returns(address disperserAddress)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) DisperserKeyToInfo(arg0 uint32) (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.DisperserKeyToInfo(&_ContractEigenDADisperserRegistry.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADisperserRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) Owner() (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.Owner(&_ContractEigenDADisperserRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDADisperserRegistry.Contract.Owner(&_ContractEigenDADisperserRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "initialize", _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Initialize(&_ContractEigenDADisperserRegistry.TransactOpts, _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.Initialize(&_ContractEigenDADisperserRegistry.TransactOpts, _initialOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.RenounceOwnership(&_ContractEigenDADisperserRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.RenounceOwnership(&_ContractEigenDADisperserRegistry.TransactOpts)
}

// SetDisperserInfo is a paid mutator transaction binding the contract method 0x9a0f62a0.
//
// Solidity: function setDisperserInfo(uint32 _disperserKey, (address) _disperserInfo) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) SetDisperserInfo(opts *bind.TransactOpts, _disperserKey uint32, _disperserInfo EigenDATypesV2DisperserInfo) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "setDisperserInfo", _disperserKey, _disperserInfo)
}

// SetDisperserInfo is a paid mutator transaction binding the contract method 0x9a0f62a0.
//
// Solidity: function setDisperserInfo(uint32 _disperserKey, (address) _disperserInfo) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) SetDisperserInfo(_disperserKey uint32, _disperserInfo EigenDATypesV2DisperserInfo) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.SetDisperserInfo(&_ContractEigenDADisperserRegistry.TransactOpts, _disperserKey, _disperserInfo)
}

// SetDisperserInfo is a paid mutator transaction binding the contract method 0x9a0f62a0.
//
// Solidity: function setDisperserInfo(uint32 _disperserKey, (address) _disperserInfo) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) SetDisperserInfo(_disperserKey uint32, _disperserInfo EigenDATypesV2DisperserInfo) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.SetDisperserInfo(&_ContractEigenDADisperserRegistry.TransactOpts, _disperserKey, _disperserInfo)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADisperserRegistry.Contract.TransferOwnership(&_ContractEigenDADisperserRegistry.TransactOpts, newOwner)
}

// ContractEigenDADisperserRegistryDisperserAddedIterator is returned from FilterDisperserAdded and is used to iterate over the raw logs and unpacked data for DisperserAdded events raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryDisperserAddedIterator struct {
	Event *ContractEigenDADisperserRegistryDisperserAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADisperserRegistryDisperserAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADisperserRegistryDisperserAdded)
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
		it.Event = new(ContractEigenDADisperserRegistryDisperserAdded)
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
func (it *ContractEigenDADisperserRegistryDisperserAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADisperserRegistryDisperserAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADisperserRegistryDisperserAdded represents a DisperserAdded event raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryDisperserAdded struct {
	Key       uint32
	Disperser common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDisperserAdded is a free log retrieval operation binding the contract event 0x97fb4432fef273711f9ccc876095cf8e22b00f159658bbd807a8ea80a4c3c859.
//
// Solidity: event DisperserAdded(uint32 indexed key, address indexed disperser)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) FilterDisperserAdded(opts *bind.FilterOpts, key []uint32, disperser []common.Address) (*ContractEigenDADisperserRegistryDisperserAddedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var disperserRule []interface{}
	for _, disperserItem := range disperser {
		disperserRule = append(disperserRule, disperserItem)
	}

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.FilterLogs(opts, "DisperserAdded", keyRule, disperserRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryDisperserAddedIterator{contract: _ContractEigenDADisperserRegistry.contract, event: "DisperserAdded", logs: logs, sub: sub}, nil
}

// WatchDisperserAdded is a free log subscription operation binding the contract event 0x97fb4432fef273711f9ccc876095cf8e22b00f159658bbd807a8ea80a4c3c859.
//
// Solidity: event DisperserAdded(uint32 indexed key, address indexed disperser)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) WatchDisperserAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDADisperserRegistryDisperserAdded, key []uint32, disperser []common.Address) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var disperserRule []interface{}
	for _, disperserItem := range disperser {
		disperserRule = append(disperserRule, disperserItem)
	}

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.WatchLogs(opts, "DisperserAdded", keyRule, disperserRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADisperserRegistryDisperserAdded)
				if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "DisperserAdded", log); err != nil {
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

// ParseDisperserAdded is a log parse operation binding the contract event 0x97fb4432fef273711f9ccc876095cf8e22b00f159658bbd807a8ea80a4c3c859.
//
// Solidity: event DisperserAdded(uint32 indexed key, address indexed disperser)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) ParseDisperserAdded(log types.Log) (*ContractEigenDADisperserRegistryDisperserAdded, error) {
	event := new(ContractEigenDADisperserRegistryDisperserAdded)
	if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "DisperserAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADisperserRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryInitializedIterator struct {
	Event *ContractEigenDADisperserRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADisperserRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADisperserRegistryInitialized)
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
		it.Event = new(ContractEigenDADisperserRegistryInitialized)
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
func (it *ContractEigenDADisperserRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADisperserRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADisperserRegistryInitialized represents a Initialized event raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDADisperserRegistryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryInitializedIterator{contract: _ContractEigenDADisperserRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDADisperserRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADisperserRegistryInitialized)
				if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) ParseInitialized(log types.Log) (*ContractEigenDADisperserRegistryInitialized, error) {
	event := new(ContractEigenDADisperserRegistryInitialized)
	if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADisperserRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryOwnershipTransferredIterator struct {
	Event *ContractEigenDADisperserRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADisperserRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADisperserRegistryOwnershipTransferred)
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
		it.Event = new(ContractEigenDADisperserRegistryOwnershipTransferred)
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
func (it *ContractEigenDADisperserRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADisperserRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADisperserRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDADisperserRegistry contract.
type ContractEigenDADisperserRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDADisperserRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADisperserRegistryOwnershipTransferredIterator{contract: _ContractEigenDADisperserRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDADisperserRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDADisperserRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADisperserRegistryOwnershipTransferred)
				if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDADisperserRegistry *ContractEigenDADisperserRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDADisperserRegistryOwnershipTransferred, error) {
	event := new(ContractEigenDADisperserRegistryOwnershipTransferred)
	if err := _ContractEigenDADisperserRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
