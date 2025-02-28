// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDARelayRegistry

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

// RelayInfo is an auto generated low-level Go binding around an user-defined struct.
type RelayInfo struct {
	RelayAddress common.Address
	RelayURL     string
}

// ContractIEigenDARelayRegistryMetaData contains all meta data concerning the ContractIEigenDARelayRegistry contract.
var ContractIEigenDARelayRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addRelayInfo\",\"inputs\":[{\"name\":\"relayInfo\",\"type\":\"tuple\",\"internalType\":\"structRelayInfo\",\"components\":[{\"name\":\"relayAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayURL\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"relayKeyToAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToUrl\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"RelayAdded\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"key\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"relayURL\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false}]",
}

// ContractIEigenDARelayRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDARelayRegistryMetaData.ABI instead.
var ContractIEigenDARelayRegistryABI = ContractIEigenDARelayRegistryMetaData.ABI

// ContractIEigenDARelayRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistry struct {
	ContractIEigenDARelayRegistryCaller     // Read-only binding to the contract
	ContractIEigenDARelayRegistryTransactor // Write-only binding to the contract
	ContractIEigenDARelayRegistryFilterer   // Log filterer for contract events
}

// ContractIEigenDARelayRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDARelayRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDARelayRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDARelayRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDARelayRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDARelayRegistrySession struct {
	Contract     *ContractIEigenDARelayRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractIEigenDARelayRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDARelayRegistryCallerSession struct {
	Contract *ContractIEigenDARelayRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractIEigenDARelayRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDARelayRegistryTransactorSession struct {
	Contract     *ContractIEigenDARelayRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractIEigenDARelayRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistryRaw struct {
	Contract *ContractIEigenDARelayRegistry // Generic contract binding to access the raw methods on
}

// ContractIEigenDARelayRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistryCallerRaw struct {
	Contract *ContractIEigenDARelayRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDARelayRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDARelayRegistryTransactorRaw struct {
	Contract *ContractIEigenDARelayRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDARelayRegistry creates a new instance of ContractIEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractIEigenDARelayRegistry(address common.Address, backend bind.ContractBackend) (*ContractIEigenDARelayRegistry, error) {
	contract, err := bindContractIEigenDARelayRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDARelayRegistry{ContractIEigenDARelayRegistryCaller: ContractIEigenDARelayRegistryCaller{contract: contract}, ContractIEigenDARelayRegistryTransactor: ContractIEigenDARelayRegistryTransactor{contract: contract}, ContractIEigenDARelayRegistryFilterer: ContractIEigenDARelayRegistryFilterer{contract: contract}}, nil
}

// NewContractIEigenDARelayRegistryCaller creates a new read-only instance of ContractIEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractIEigenDARelayRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDARelayRegistryCaller, error) {
	contract, err := bindContractIEigenDARelayRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDARelayRegistryCaller{contract: contract}, nil
}

// NewContractIEigenDARelayRegistryTransactor creates a new write-only instance of ContractIEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractIEigenDARelayRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDARelayRegistryTransactor, error) {
	contract, err := bindContractIEigenDARelayRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDARelayRegistryTransactor{contract: contract}, nil
}

// NewContractIEigenDARelayRegistryFilterer creates a new log filterer instance of ContractIEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractIEigenDARelayRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDARelayRegistryFilterer, error) {
	contract, err := bindContractIEigenDARelayRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDARelayRegistryFilterer{contract: contract}, nil
}

// bindContractIEigenDARelayRegistry binds a generic wrapper to an already deployed contract.
func bindContractIEigenDARelayRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDARelayRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDARelayRegistry.Contract.ContractIEigenDARelayRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.ContractIEigenDARelayRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.ContractIEigenDARelayRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDARelayRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.contract.Transact(opts, method, params...)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryCaller) RelayKeyToAddress(opts *bind.CallOpts, key uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistrySession) RelayKeyToAddress(key uint32) (common.Address, error) {
	return _ContractIEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractIEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 key) view returns(address)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryCallerSession) RelayKeyToAddress(key uint32) (common.Address, error) {
	return _ContractIEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractIEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryCaller) RelayKeyToUrl(opts *bind.CallOpts, key uint32) (string, error) {
	var out []interface{}
	err := _ContractIEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToUrl", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistrySession) RelayKeyToUrl(key uint32) (string, error) {
	return _ContractIEigenDARelayRegistry.Contract.RelayKeyToUrl(&_ContractIEigenDARelayRegistry.CallOpts, key)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 key) view returns(string)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryCallerSession) RelayKeyToUrl(key uint32) (string, error) {
	return _ContractIEigenDARelayRegistry.Contract.RelayKeyToUrl(&_ContractIEigenDARelayRegistry.CallOpts, key)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryTransactor) AddRelayInfo(opts *bind.TransactOpts, relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.contract.Transact(opts, "addRelayInfo", relayInfo)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistrySession) AddRelayInfo(relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.AddRelayInfo(&_ContractIEigenDARelayRegistry.TransactOpts, relayInfo)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x2fc35013.
//
// Solidity: function addRelayInfo((address,string) relayInfo) returns(uint32)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryTransactorSession) AddRelayInfo(relayInfo RelayInfo) (*types.Transaction, error) {
	return _ContractIEigenDARelayRegistry.Contract.AddRelayInfo(&_ContractIEigenDARelayRegistry.TransactOpts, relayInfo)
}

// ContractIEigenDARelayRegistryRelayAddedIterator is returned from FilterRelayAdded and is used to iterate over the raw logs and unpacked data for RelayAdded events raised by the ContractIEigenDARelayRegistry contract.
type ContractIEigenDARelayRegistryRelayAddedIterator struct {
	Event *ContractIEigenDARelayRegistryRelayAdded // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDARelayRegistryRelayAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDARelayRegistryRelayAdded)
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
		it.Event = new(ContractIEigenDARelayRegistryRelayAdded)
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
func (it *ContractIEigenDARelayRegistryRelayAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDARelayRegistryRelayAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDARelayRegistryRelayAdded represents a RelayAdded event raised by the ContractIEigenDARelayRegistry contract.
type ContractIEigenDARelayRegistryRelayAdded struct {
	Relay    common.Address
	Key      uint32
	RelayURL string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRelayAdded is a free log retrieval operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryFilterer) FilterRelayAdded(opts *bind.FilterOpts, relay []common.Address, key []uint32) (*ContractIEigenDARelayRegistryRelayAddedIterator, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDARelayRegistry.contract.FilterLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDARelayRegistryRelayAddedIterator{contract: _ContractIEigenDARelayRegistry.contract, event: "RelayAdded", logs: logs, sub: sub}, nil
}

// WatchRelayAdded is a free log subscription operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryFilterer) WatchRelayAdded(opts *bind.WatchOpts, sink chan<- *ContractIEigenDARelayRegistryRelayAdded, relay []common.Address, key []uint32) (event.Subscription, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractIEigenDARelayRegistry.contract.WatchLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDARelayRegistryRelayAdded)
				if err := _ContractIEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
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

// ParseRelayAdded is a log parse operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractIEigenDARelayRegistry *ContractIEigenDARelayRegistryFilterer) ParseRelayAdded(log types.Log) (*ContractIEigenDARelayRegistryRelayAdded, error) {
	event := new(ContractIEigenDARelayRegistryRelayAdded)
	if err := _ContractIEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
