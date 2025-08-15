// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAEjectionLib

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

// ContractEigenDAEjectionLibMetaData contains all meta data concerning the ContractEigenDAEjectionLib contract.
var ContractEigenDAEjectionLibMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"event\",\"name\":\"CooldownSet\",\"inputs\":[{\"name\":\"cooldown\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelaySet\",\"inputs\":[{\"name\":\"delay\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EjectionCancelled\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EjectionCompleted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EjectionStarted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"timestampStarted\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"ejectionTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false}]",
}

// ContractEigenDAEjectionLibABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAEjectionLibMetaData.ABI instead.
var ContractEigenDAEjectionLibABI = ContractEigenDAEjectionLibMetaData.ABI

// ContractEigenDAEjectionLib is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAEjectionLib struct {
	ContractEigenDAEjectionLibCaller     // Read-only binding to the contract
	ContractEigenDAEjectionLibTransactor // Write-only binding to the contract
	ContractEigenDAEjectionLibFilterer   // Log filterer for contract events
}

// ContractEigenDAEjectionLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAEjectionLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAEjectionLibSession struct {
	Contract     *ContractEigenDAEjectionLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ContractEigenDAEjectionLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAEjectionLibCallerSession struct {
	Contract *ContractEigenDAEjectionLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ContractEigenDAEjectionLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAEjectionLibTransactorSession struct {
	Contract     *ContractEigenDAEjectionLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ContractEigenDAEjectionLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAEjectionLibRaw struct {
	Contract *ContractEigenDAEjectionLib // Generic contract binding to access the raw methods on
}

// ContractEigenDAEjectionLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionLibCallerRaw struct {
	Contract *ContractEigenDAEjectionLibCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAEjectionLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionLibTransactorRaw struct {
	Contract *ContractEigenDAEjectionLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAEjectionLib creates a new instance of ContractEigenDAEjectionLib, bound to a specific deployed contract.
func NewContractEigenDAEjectionLib(address common.Address, backend bind.ContractBackend) (*ContractEigenDAEjectionLib, error) {
	contract, err := bindContractEigenDAEjectionLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLib{ContractEigenDAEjectionLibCaller: ContractEigenDAEjectionLibCaller{contract: contract}, ContractEigenDAEjectionLibTransactor: ContractEigenDAEjectionLibTransactor{contract: contract}, ContractEigenDAEjectionLibFilterer: ContractEigenDAEjectionLibFilterer{contract: contract}}, nil
}

// NewContractEigenDAEjectionLibCaller creates a new read-only instance of ContractEigenDAEjectionLib, bound to a specific deployed contract.
func NewContractEigenDAEjectionLibCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAEjectionLibCaller, error) {
	contract, err := bindContractEigenDAEjectionLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibCaller{contract: contract}, nil
}

// NewContractEigenDAEjectionLibTransactor creates a new write-only instance of ContractEigenDAEjectionLib, bound to a specific deployed contract.
func NewContractEigenDAEjectionLibTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAEjectionLibTransactor, error) {
	contract, err := bindContractEigenDAEjectionLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibTransactor{contract: contract}, nil
}

// NewContractEigenDAEjectionLibFilterer creates a new log filterer instance of ContractEigenDAEjectionLib, bound to a specific deployed contract.
func NewContractEigenDAEjectionLibFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAEjectionLibFilterer, error) {
	contract, err := bindContractEigenDAEjectionLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibFilterer{contract: contract}, nil
}

// bindContractEigenDAEjectionLib binds a generic wrapper to an already deployed contract.
func bindContractEigenDAEjectionLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAEjectionLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAEjectionLib.Contract.ContractEigenDAEjectionLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAEjectionLib.Contract.ContractEigenDAEjectionLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAEjectionLib.Contract.ContractEigenDAEjectionLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAEjectionLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAEjectionLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAEjectionLib.Contract.contract.Transact(opts, method, params...)
}

// ContractEigenDAEjectionLibCooldownSetIterator is returned from FilterCooldownSet and is used to iterate over the raw logs and unpacked data for CooldownSet events raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibCooldownSetIterator struct {
	Event *ContractEigenDAEjectionLibCooldownSet // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAEjectionLibCooldownSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAEjectionLibCooldownSet)
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
		it.Event = new(ContractEigenDAEjectionLibCooldownSet)
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
func (it *ContractEigenDAEjectionLibCooldownSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAEjectionLibCooldownSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAEjectionLibCooldownSet represents a CooldownSet event raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibCooldownSet struct {
	Cooldown uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCooldownSet is a free log retrieval operation binding the contract event 0xd3c86bd8150fcc5f376375097a2a9b313bf89addc501ee368c9eebe117328497.
//
// Solidity: event CooldownSet(uint64 cooldown)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) FilterCooldownSet(opts *bind.FilterOpts) (*ContractEigenDAEjectionLibCooldownSetIterator, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.FilterLogs(opts, "CooldownSet")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibCooldownSetIterator{contract: _ContractEigenDAEjectionLib.contract, event: "CooldownSet", logs: logs, sub: sub}, nil
}

// WatchCooldownSet is a free log subscription operation binding the contract event 0xd3c86bd8150fcc5f376375097a2a9b313bf89addc501ee368c9eebe117328497.
//
// Solidity: event CooldownSet(uint64 cooldown)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) WatchCooldownSet(opts *bind.WatchOpts, sink chan<- *ContractEigenDAEjectionLibCooldownSet) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.WatchLogs(opts, "CooldownSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAEjectionLibCooldownSet)
				if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "CooldownSet", log); err != nil {
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

// ParseCooldownSet is a log parse operation binding the contract event 0xd3c86bd8150fcc5f376375097a2a9b313bf89addc501ee368c9eebe117328497.
//
// Solidity: event CooldownSet(uint64 cooldown)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) ParseCooldownSet(log types.Log) (*ContractEigenDAEjectionLibCooldownSet, error) {
	event := new(ContractEigenDAEjectionLibCooldownSet)
	if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "CooldownSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAEjectionLibDelaySetIterator is returned from FilterDelaySet and is used to iterate over the raw logs and unpacked data for DelaySet events raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibDelaySetIterator struct {
	Event *ContractEigenDAEjectionLibDelaySet // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAEjectionLibDelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAEjectionLibDelaySet)
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
		it.Event = new(ContractEigenDAEjectionLibDelaySet)
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
func (it *ContractEigenDAEjectionLibDelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAEjectionLibDelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAEjectionLibDelaySet represents a DelaySet event raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibDelaySet struct {
	Delay uint64
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDelaySet is a free log retrieval operation binding the contract event 0xc03c652cce63e6a4eb932e2608f14d717836f14b42a735017b4fbbdc5b5da27e.
//
// Solidity: event DelaySet(uint64 delay)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) FilterDelaySet(opts *bind.FilterOpts) (*ContractEigenDAEjectionLibDelaySetIterator, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.FilterLogs(opts, "DelaySet")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibDelaySetIterator{contract: _ContractEigenDAEjectionLib.contract, event: "DelaySet", logs: logs, sub: sub}, nil
}

// WatchDelaySet is a free log subscription operation binding the contract event 0xc03c652cce63e6a4eb932e2608f14d717836f14b42a735017b4fbbdc5b5da27e.
//
// Solidity: event DelaySet(uint64 delay)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) WatchDelaySet(opts *bind.WatchOpts, sink chan<- *ContractEigenDAEjectionLibDelaySet) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.WatchLogs(opts, "DelaySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAEjectionLibDelaySet)
				if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "DelaySet", log); err != nil {
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

// ParseDelaySet is a log parse operation binding the contract event 0xc03c652cce63e6a4eb932e2608f14d717836f14b42a735017b4fbbdc5b5da27e.
//
// Solidity: event DelaySet(uint64 delay)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) ParseDelaySet(log types.Log) (*ContractEigenDAEjectionLibDelaySet, error) {
	event := new(ContractEigenDAEjectionLibDelaySet)
	if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "DelaySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAEjectionLibEjectionCancelledIterator is returned from FilterEjectionCancelled and is used to iterate over the raw logs and unpacked data for EjectionCancelled events raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionCancelledIterator struct {
	Event *ContractEigenDAEjectionLibEjectionCancelled // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAEjectionLibEjectionCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAEjectionLibEjectionCancelled)
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
		it.Event = new(ContractEigenDAEjectionLibEjectionCancelled)
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
func (it *ContractEigenDAEjectionLibEjectionCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAEjectionLibEjectionCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAEjectionLibEjectionCancelled represents a EjectionCancelled event raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionCancelled struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEjectionCancelled is a free log retrieval operation binding the contract event 0x3123488031678febf5655aec782d1252ab6d922921498edd621d5ea339e292e1.
//
// Solidity: event EjectionCancelled(address operator)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) FilterEjectionCancelled(opts *bind.FilterOpts) (*ContractEigenDAEjectionLibEjectionCancelledIterator, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.FilterLogs(opts, "EjectionCancelled")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibEjectionCancelledIterator{contract: _ContractEigenDAEjectionLib.contract, event: "EjectionCancelled", logs: logs, sub: sub}, nil
}

// WatchEjectionCancelled is a free log subscription operation binding the contract event 0x3123488031678febf5655aec782d1252ab6d922921498edd621d5ea339e292e1.
//
// Solidity: event EjectionCancelled(address operator)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) WatchEjectionCancelled(opts *bind.WatchOpts, sink chan<- *ContractEigenDAEjectionLibEjectionCancelled) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.WatchLogs(opts, "EjectionCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAEjectionLibEjectionCancelled)
				if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionCancelled", log); err != nil {
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

// ParseEjectionCancelled is a log parse operation binding the contract event 0x3123488031678febf5655aec782d1252ab6d922921498edd621d5ea339e292e1.
//
// Solidity: event EjectionCancelled(address operator)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) ParseEjectionCancelled(log types.Log) (*ContractEigenDAEjectionLibEjectionCancelled, error) {
	event := new(ContractEigenDAEjectionLibEjectionCancelled)
	if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAEjectionLibEjectionCompletedIterator is returned from FilterEjectionCompleted and is used to iterate over the raw logs and unpacked data for EjectionCompleted events raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionCompletedIterator struct {
	Event *ContractEigenDAEjectionLibEjectionCompleted // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAEjectionLibEjectionCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAEjectionLibEjectionCompleted)
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
		it.Event = new(ContractEigenDAEjectionLibEjectionCompleted)
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
func (it *ContractEigenDAEjectionLibEjectionCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAEjectionLibEjectionCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAEjectionLibEjectionCompleted represents a EjectionCompleted event raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionCompleted struct {
	Operator common.Address
	Quorums  []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEjectionCompleted is a free log retrieval operation binding the contract event 0x172c16b12b5cfa6a611b21c65ea2c32bcac50585932f62c4544db6fca0a9b9bd.
//
// Solidity: event EjectionCompleted(address operator, bytes quorums)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) FilterEjectionCompleted(opts *bind.FilterOpts) (*ContractEigenDAEjectionLibEjectionCompletedIterator, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.FilterLogs(opts, "EjectionCompleted")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibEjectionCompletedIterator{contract: _ContractEigenDAEjectionLib.contract, event: "EjectionCompleted", logs: logs, sub: sub}, nil
}

// WatchEjectionCompleted is a free log subscription operation binding the contract event 0x172c16b12b5cfa6a611b21c65ea2c32bcac50585932f62c4544db6fca0a9b9bd.
//
// Solidity: event EjectionCompleted(address operator, bytes quorums)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) WatchEjectionCompleted(opts *bind.WatchOpts, sink chan<- *ContractEigenDAEjectionLibEjectionCompleted) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.WatchLogs(opts, "EjectionCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAEjectionLibEjectionCompleted)
				if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionCompleted", log); err != nil {
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

// ParseEjectionCompleted is a log parse operation binding the contract event 0x172c16b12b5cfa6a611b21c65ea2c32bcac50585932f62c4544db6fca0a9b9bd.
//
// Solidity: event EjectionCompleted(address operator, bytes quorums)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) ParseEjectionCompleted(log types.Log) (*ContractEigenDAEjectionLibEjectionCompleted, error) {
	event := new(ContractEigenDAEjectionLibEjectionCompleted)
	if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAEjectionLibEjectionStartedIterator is returned from FilterEjectionStarted and is used to iterate over the raw logs and unpacked data for EjectionStarted events raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionStartedIterator struct {
	Event *ContractEigenDAEjectionLibEjectionStarted // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAEjectionLibEjectionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAEjectionLibEjectionStarted)
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
		it.Event = new(ContractEigenDAEjectionLibEjectionStarted)
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
func (it *ContractEigenDAEjectionLibEjectionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAEjectionLibEjectionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAEjectionLibEjectionStarted represents a EjectionStarted event raised by the ContractEigenDAEjectionLib contract.
type ContractEigenDAEjectionLibEjectionStarted struct {
	Operator         common.Address
	Quorums          []byte
	TimestampStarted uint64
	EjectionTime     uint64
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterEjectionStarted is a free log retrieval operation binding the contract event 0xca0b1f475104769e928fe6cb95b114ff6aded77892f03b9c9201f8e18c6510de.
//
// Solidity: event EjectionStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) FilterEjectionStarted(opts *bind.FilterOpts) (*ContractEigenDAEjectionLibEjectionStartedIterator, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.FilterLogs(opts, "EjectionStarted")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionLibEjectionStartedIterator{contract: _ContractEigenDAEjectionLib.contract, event: "EjectionStarted", logs: logs, sub: sub}, nil
}

// WatchEjectionStarted is a free log subscription operation binding the contract event 0xca0b1f475104769e928fe6cb95b114ff6aded77892f03b9c9201f8e18c6510de.
//
// Solidity: event EjectionStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) WatchEjectionStarted(opts *bind.WatchOpts, sink chan<- *ContractEigenDAEjectionLibEjectionStarted) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAEjectionLib.contract.WatchLogs(opts, "EjectionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAEjectionLibEjectionStarted)
				if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionStarted", log); err != nil {
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

// ParseEjectionStarted is a log parse operation binding the contract event 0xca0b1f475104769e928fe6cb95b114ff6aded77892f03b9c9201f8e18c6510de.
//
// Solidity: event EjectionStarted(address operator, bytes quorums, uint64 timestampStarted, uint64 ejectionTime)
func (_ContractEigenDAEjectionLib *ContractEigenDAEjectionLibFilterer) ParseEjectionStarted(log types.Log) (*ContractEigenDAEjectionLibEjectionStarted, error) {
	event := new(ContractEigenDAEjectionLibEjectionStarted)
	if err := _ContractEigenDAEjectionLib.contract.UnpackLog(event, "EjectionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
