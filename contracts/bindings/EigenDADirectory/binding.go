// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDADirectory

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

// ContractEigenDADirectoryMetaData contains all meta data concerning the ContractEigenDADirectory contract.
var ContractEigenDADirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"replaceAddress\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddressAdded\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressRemoved\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AddressReplaced\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"key\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyExists\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"AddressDoesNotExist\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NewValueIsOldValue\",\"inputs\":[{\"name\":\"value\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// ContractEigenDADirectoryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDADirectoryMetaData.ABI instead.
var ContractEigenDADirectoryABI = ContractEigenDADirectoryMetaData.ABI

// ContractEigenDADirectory is an auto generated Go binding around an Ethereum contract.
type ContractEigenDADirectory struct {
	ContractEigenDADirectoryCaller     // Read-only binding to the contract
	ContractEigenDADirectoryTransactor // Write-only binding to the contract
	ContractEigenDADirectoryFilterer   // Log filterer for contract events
}

// ContractEigenDADirectoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDADirectoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADirectoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDADirectoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADirectoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDADirectoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDADirectorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDADirectorySession struct {
	Contract     *ContractEigenDADirectory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ContractEigenDADirectoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDADirectoryCallerSession struct {
	Contract *ContractEigenDADirectoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// ContractEigenDADirectoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDADirectoryTransactorSession struct {
	Contract     *ContractEigenDADirectoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ContractEigenDADirectoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDADirectoryRaw struct {
	Contract *ContractEigenDADirectory // Generic contract binding to access the raw methods on
}

// ContractEigenDADirectoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDADirectoryCallerRaw struct {
	Contract *ContractEigenDADirectoryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDADirectoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDADirectoryTransactorRaw struct {
	Contract *ContractEigenDADirectoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDADirectory creates a new instance of ContractEigenDADirectory, bound to a specific deployed contract.
func NewContractEigenDADirectory(address common.Address, backend bind.ContractBackend) (*ContractEigenDADirectory, error) {
	contract, err := bindContractEigenDADirectory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectory{ContractEigenDADirectoryCaller: ContractEigenDADirectoryCaller{contract: contract}, ContractEigenDADirectoryTransactor: ContractEigenDADirectoryTransactor{contract: contract}, ContractEigenDADirectoryFilterer: ContractEigenDADirectoryFilterer{contract: contract}}, nil
}

// NewContractEigenDADirectoryCaller creates a new read-only instance of ContractEigenDADirectory, bound to a specific deployed contract.
func NewContractEigenDADirectoryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDADirectoryCaller, error) {
	contract, err := bindContractEigenDADirectory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryCaller{contract: contract}, nil
}

// NewContractEigenDADirectoryTransactor creates a new write-only instance of ContractEigenDADirectory, bound to a specific deployed contract.
func NewContractEigenDADirectoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDADirectoryTransactor, error) {
	contract, err := bindContractEigenDADirectory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryTransactor{contract: contract}, nil
}

// NewContractEigenDADirectoryFilterer creates a new log filterer instance of ContractEigenDADirectory, bound to a specific deployed contract.
func NewContractEigenDADirectoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDADirectoryFilterer, error) {
	contract, err := bindContractEigenDADirectory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryFilterer{contract: contract}, nil
}

// bindContractEigenDADirectory binds a generic wrapper to an already deployed contract.
func bindContractEigenDADirectory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDADirectoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADirectory *ContractEigenDADirectoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADirectory.Contract.ContractEigenDADirectoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADirectory *ContractEigenDADirectoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.ContractEigenDADirectoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADirectory *ContractEigenDADirectoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.ContractEigenDADirectoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDADirectory *ContractEigenDADirectoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDADirectory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.contract.Transact(opts, method, params...)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCaller) GetAddress(opts *bind.CallOpts, key [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADirectory.contract.Call(opts, &out, "getAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) GetAddress(key [32]byte) (common.Address, error) {
	return _ContractEigenDADirectory.Contract.GetAddress(&_ContractEigenDADirectory.CallOpts, key)
}

// GetAddress is a free data retrieval call binding the contract method 0x21f8a721.
//
// Solidity: function getAddress(bytes32 key) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCallerSession) GetAddress(key [32]byte) (common.Address, error) {
	return _ContractEigenDADirectory.Contract.GetAddress(&_ContractEigenDADirectory.CallOpts, key)
}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCaller) GetAddress0(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADirectory.contract.Call(opts, &out, "getAddress0", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) GetAddress0(name string) (common.Address, error) {
	return _ContractEigenDADirectory.Contract.GetAddress0(&_ContractEigenDADirectory.CallOpts, name)
}

// GetAddress0 is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string name) view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCallerSession) GetAddress0(name string) (common.Address, error) {
	return _ContractEigenDADirectory.Contract.GetAddress0(&_ContractEigenDADirectory.CallOpts, name)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDADirectory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) Owner() (common.Address, error) {
	return _ContractEigenDADirectory.Contract.Owner(&_ContractEigenDADirectory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDADirectory *ContractEigenDADirectoryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDADirectory.Contract.Owner(&_ContractEigenDADirectory.CallOpts)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) AddAddress(opts *bind.TransactOpts, name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "addAddress", name, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) AddAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.AddAddress(&_ContractEigenDADirectory.TransactOpts, name, value)
}

// AddAddress is a paid mutator transaction binding the contract method 0xceb35b0f.
//
// Solidity: function addAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) AddAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.AddAddress(&_ContractEigenDADirectory.TransactOpts, name, value)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "initialize", _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.Initialize(&_ContractEigenDADirectory.TransactOpts, _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.Initialize(&_ContractEigenDADirectory.TransactOpts, _initialOwner)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) RemoveAddress(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "removeAddress", name)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) RemoveAddress(name string) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.RemoveAddress(&_ContractEigenDADirectory.TransactOpts, name)
}

// RemoveAddress is a paid mutator transaction binding the contract method 0xf94d1312.
//
// Solidity: function removeAddress(string name) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) RemoveAddress(name string) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.RemoveAddress(&_ContractEigenDADirectory.TransactOpts, name)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.RenounceOwnership(&_ContractEigenDADirectory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.RenounceOwnership(&_ContractEigenDADirectory.TransactOpts)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) ReplaceAddress(opts *bind.TransactOpts, name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "replaceAddress", name, value)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) ReplaceAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.ReplaceAddress(&_ContractEigenDADirectory.TransactOpts, name, value)
}

// ReplaceAddress is a paid mutator transaction binding the contract method 0x1d7762e7.
//
// Solidity: function replaceAddress(string name, address value) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) ReplaceAddress(name string, value common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.ReplaceAddress(&_ContractEigenDADirectory.TransactOpts, name, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.TransferOwnership(&_ContractEigenDADirectory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDADirectory *ContractEigenDADirectoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDADirectory.Contract.TransferOwnership(&_ContractEigenDADirectory.TransactOpts, newOwner)
}

// ContractEigenDADirectoryAddressAddedIterator is returned from FilterAddressAdded and is used to iterate over the raw logs and unpacked data for AddressAdded events raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressAddedIterator struct {
	Event *ContractEigenDADirectoryAddressAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADirectoryAddressAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADirectoryAddressAdded)
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
		it.Event = new(ContractEigenDADirectoryAddressAdded)
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
func (it *ContractEigenDADirectoryAddressAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADirectoryAddressAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADirectoryAddressAdded represents a AddressAdded event raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressAdded struct {
	Name  string
	Key   [32]byte
	Value common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddressAdded is a free log retrieval operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) FilterAddressAdded(opts *bind.FilterOpts, key [][32]byte, value []common.Address) (*ContractEigenDADirectoryAddressAddedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.FilterLogs(opts, "AddressAdded", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryAddressAddedIterator{contract: _ContractEigenDADirectory.contract, event: "AddressAdded", logs: logs, sub: sub}, nil
}

// WatchAddressAdded is a free log subscription operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) WatchAddressAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDADirectoryAddressAdded, key [][32]byte, value []common.Address) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.WatchLogs(opts, "AddressAdded", keyRule, valueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADirectoryAddressAdded)
				if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressAdded", log); err != nil {
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

// ParseAddressAdded is a log parse operation binding the contract event 0x6db5569d223c840fb38a83e4a556cb60a251b9680de393e47777870cdbac26e6.
//
// Solidity: event AddressAdded(string name, bytes32 indexed key, address indexed value)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) ParseAddressAdded(log types.Log) (*ContractEigenDADirectoryAddressAdded, error) {
	event := new(ContractEigenDADirectoryAddressAdded)
	if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADirectoryAddressRemovedIterator is returned from FilterAddressRemoved and is used to iterate over the raw logs and unpacked data for AddressRemoved events raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressRemovedIterator struct {
	Event *ContractEigenDADirectoryAddressRemoved // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADirectoryAddressRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADirectoryAddressRemoved)
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
		it.Event = new(ContractEigenDADirectoryAddressRemoved)
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
func (it *ContractEigenDADirectoryAddressRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADirectoryAddressRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADirectoryAddressRemoved represents a AddressRemoved event raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressRemoved struct {
	Name string
	Key  [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddressRemoved is a free log retrieval operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) FilterAddressRemoved(opts *bind.FilterOpts, key [][32]byte) (*ContractEigenDADirectoryAddressRemovedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.FilterLogs(opts, "AddressRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryAddressRemovedIterator{contract: _ContractEigenDADirectory.contract, event: "AddressRemoved", logs: logs, sub: sub}, nil
}

// WatchAddressRemoved is a free log subscription operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) WatchAddressRemoved(opts *bind.WatchOpts, sink chan<- *ContractEigenDADirectoryAddressRemoved, key [][32]byte) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.WatchLogs(opts, "AddressRemoved", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADirectoryAddressRemoved)
				if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressRemoved", log); err != nil {
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

// ParseAddressRemoved is a log parse operation binding the contract event 0xabb104e9a16f893503445ca24334a10468322f797b67092c3f53021fc4ee5022.
//
// Solidity: event AddressRemoved(string name, bytes32 indexed key)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) ParseAddressRemoved(log types.Log) (*ContractEigenDADirectoryAddressRemoved, error) {
	event := new(ContractEigenDADirectoryAddressRemoved)
	if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADirectoryAddressReplacedIterator is returned from FilterAddressReplaced and is used to iterate over the raw logs and unpacked data for AddressReplaced events raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressReplacedIterator struct {
	Event *ContractEigenDADirectoryAddressReplaced // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADirectoryAddressReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADirectoryAddressReplaced)
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
		it.Event = new(ContractEigenDADirectoryAddressReplaced)
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
func (it *ContractEigenDADirectoryAddressReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADirectoryAddressReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADirectoryAddressReplaced represents a AddressReplaced event raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryAddressReplaced struct {
	Name     string
	Key      [32]byte
	OldValue common.Address
	NewValue common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddressReplaced is a free log retrieval operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) FilterAddressReplaced(opts *bind.FilterOpts, key [][32]byte, oldValue []common.Address, newValue []common.Address) (*ContractEigenDADirectoryAddressReplacedIterator, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var oldValueRule []interface{}
	for _, oldValueItem := range oldValue {
		oldValueRule = append(oldValueRule, oldValueItem)
	}
	var newValueRule []interface{}
	for _, newValueItem := range newValue {
		newValueRule = append(newValueRule, newValueItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.FilterLogs(opts, "AddressReplaced", keyRule, oldValueRule, newValueRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryAddressReplacedIterator{contract: _ContractEigenDADirectory.contract, event: "AddressReplaced", logs: logs, sub: sub}, nil
}

// WatchAddressReplaced is a free log subscription operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) WatchAddressReplaced(opts *bind.WatchOpts, sink chan<- *ContractEigenDADirectoryAddressReplaced, key [][32]byte, oldValue []common.Address, newValue []common.Address) (event.Subscription, error) {

	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}
	var oldValueRule []interface{}
	for _, oldValueItem := range oldValue {
		oldValueRule = append(oldValueRule, oldValueItem)
	}
	var newValueRule []interface{}
	for _, newValueItem := range newValue {
		newValueRule = append(newValueRule, newValueItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.WatchLogs(opts, "AddressReplaced", keyRule, oldValueRule, newValueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADirectoryAddressReplaced)
				if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressReplaced", log); err != nil {
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

// ParseAddressReplaced is a log parse operation binding the contract event 0x236883d8e01cc81c0167947f15527771a12a5a51c0670674b60e2b9794a3647f.
//
// Solidity: event AddressReplaced(string name, bytes32 indexed key, address indexed oldValue, address indexed newValue)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) ParseAddressReplaced(log types.Log) (*ContractEigenDADirectoryAddressReplaced, error) {
	event := new(ContractEigenDADirectoryAddressReplaced)
	if err := _ContractEigenDADirectory.contract.UnpackLog(event, "AddressReplaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADirectoryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryInitializedIterator struct {
	Event *ContractEigenDADirectoryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADirectoryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADirectoryInitialized)
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
		it.Event = new(ContractEigenDADirectoryInitialized)
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
func (it *ContractEigenDADirectoryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADirectoryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADirectoryInitialized represents a Initialized event raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDADirectoryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDADirectory.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryInitializedIterator{contract: _ContractEigenDADirectory.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDADirectoryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDADirectory.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADirectoryInitialized)
				if err := _ContractEigenDADirectory.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) ParseInitialized(log types.Log) (*ContractEigenDADirectoryInitialized, error) {
	event := new(ContractEigenDADirectoryInitialized)
	if err := _ContractEigenDADirectory.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDADirectoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryOwnershipTransferredIterator struct {
	Event *ContractEigenDADirectoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDADirectoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDADirectoryOwnershipTransferred)
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
		it.Event = new(ContractEigenDADirectoryOwnershipTransferred)
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
func (it *ContractEigenDADirectoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDADirectoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDADirectoryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDADirectory contract.
type ContractEigenDADirectoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDADirectoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDADirectoryOwnershipTransferredIterator{contract: _ContractEigenDADirectory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDADirectoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDADirectory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDADirectoryOwnershipTransferred)
				if err := _ContractEigenDADirectory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDADirectory *ContractEigenDADirectoryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDADirectoryOwnershipTransferred, error) {
	event := new(ContractEigenDADirectoryOwnershipTransferred)
	if err := _ContractEigenDADirectory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
