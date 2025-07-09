// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIRelayRegistry

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

// ContractIRelayRegistryMetaData contains all meta data concerning the ContractIRelayRegistry contract.
var ContractIRelayRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addRelayInfo\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"url\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"dispersers\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRelayInfo\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"url\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"relayKeyToAddress\",\"inputs\":[{\"name\":\"relayId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToDispersers\",\"inputs\":[{\"name\":\"relayId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToUrl\",\"inputs\":[{\"name\":\"relayId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"RelayAdded\",\"inputs\":[{\"name\":\"relayId\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"relay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"url\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"dispersers\",\"type\":\"uint32[]\",\"indexed\":false,\"internalType\":\"uint32[]\"}],\"anonymous\":false}]",
}

// ContractIRelayRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIRelayRegistryMetaData.ABI instead.
var ContractIRelayRegistryABI = ContractIRelayRegistryMetaData.ABI

// ContractIRelayRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIRelayRegistry struct {
	ContractIRelayRegistryCaller     // Read-only binding to the contract
	ContractIRelayRegistryTransactor // Write-only binding to the contract
	ContractIRelayRegistryFilterer   // Log filterer for contract events
}

// ContractIRelayRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIRelayRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIRelayRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIRelayRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIRelayRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIRelayRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIRelayRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIRelayRegistrySession struct {
	Contract     *ContractIRelayRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ContractIRelayRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIRelayRegistryCallerSession struct {
	Contract *ContractIRelayRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ContractIRelayRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIRelayRegistryTransactorSession struct {
	Contract     *ContractIRelayRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ContractIRelayRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIRelayRegistryRaw struct {
	Contract *ContractIRelayRegistry // Generic contract binding to access the raw methods on
}

// ContractIRelayRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIRelayRegistryCallerRaw struct {
	Contract *ContractIRelayRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIRelayRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIRelayRegistryTransactorRaw struct {
	Contract *ContractIRelayRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIRelayRegistry creates a new instance of ContractIRelayRegistry, bound to a specific deployed contract.
func NewContractIRelayRegistry(address common.Address, backend bind.ContractBackend) (*ContractIRelayRegistry, error) {
	contract, err := bindContractIRelayRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIRelayRegistry{ContractIRelayRegistryCaller: ContractIRelayRegistryCaller{contract: contract}, ContractIRelayRegistryTransactor: ContractIRelayRegistryTransactor{contract: contract}, ContractIRelayRegistryFilterer: ContractIRelayRegistryFilterer{contract: contract}}, nil
}

// NewContractIRelayRegistryCaller creates a new read-only instance of ContractIRelayRegistry, bound to a specific deployed contract.
func NewContractIRelayRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIRelayRegistryCaller, error) {
	contract, err := bindContractIRelayRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIRelayRegistryCaller{contract: contract}, nil
}

// NewContractIRelayRegistryTransactor creates a new write-only instance of ContractIRelayRegistry, bound to a specific deployed contract.
func NewContractIRelayRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIRelayRegistryTransactor, error) {
	contract, err := bindContractIRelayRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIRelayRegistryTransactor{contract: contract}, nil
}

// NewContractIRelayRegistryFilterer creates a new log filterer instance of ContractIRelayRegistry, bound to a specific deployed contract.
func NewContractIRelayRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIRelayRegistryFilterer, error) {
	contract, err := bindContractIRelayRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIRelayRegistryFilterer{contract: contract}, nil
}

// bindContractIRelayRegistry binds a generic wrapper to an already deployed contract.
func bindContractIRelayRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIRelayRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIRelayRegistry *ContractIRelayRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIRelayRegistry.Contract.ContractIRelayRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIRelayRegistry *ContractIRelayRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.ContractIRelayRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIRelayRegistry *ContractIRelayRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.ContractIRelayRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIRelayRegistry *ContractIRelayRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIRelayRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.contract.Transact(opts, method, params...)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 relayId) view returns(address)
func (_ContractIRelayRegistry *ContractIRelayRegistryCaller) RelayKeyToAddress(opts *bind.CallOpts, relayId uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractIRelayRegistry.contract.Call(opts, &out, "relayKeyToAddress", relayId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 relayId) view returns(address)
func (_ContractIRelayRegistry *ContractIRelayRegistrySession) RelayKeyToAddress(relayId uint32) (common.Address, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToAddress(&_ContractIRelayRegistry.CallOpts, relayId)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 relayId) view returns(address)
func (_ContractIRelayRegistry *ContractIRelayRegistryCallerSession) RelayKeyToAddress(relayId uint32) (common.Address, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToAddress(&_ContractIRelayRegistry.CallOpts, relayId)
}

// RelayKeyToDispersers is a free data retrieval call binding the contract method 0x3994eb53.
//
// Solidity: function relayKeyToDispersers(uint32 relayId) view returns(uint32[])
func (_ContractIRelayRegistry *ContractIRelayRegistryCaller) RelayKeyToDispersers(opts *bind.CallOpts, relayId uint32) ([]uint32, error) {
	var out []interface{}
	err := _ContractIRelayRegistry.contract.Call(opts, &out, "relayKeyToDispersers", relayId)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// RelayKeyToDispersers is a free data retrieval call binding the contract method 0x3994eb53.
//
// Solidity: function relayKeyToDispersers(uint32 relayId) view returns(uint32[])
func (_ContractIRelayRegistry *ContractIRelayRegistrySession) RelayKeyToDispersers(relayId uint32) ([]uint32, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToDispersers(&_ContractIRelayRegistry.CallOpts, relayId)
}

// RelayKeyToDispersers is a free data retrieval call binding the contract method 0x3994eb53.
//
// Solidity: function relayKeyToDispersers(uint32 relayId) view returns(uint32[])
func (_ContractIRelayRegistry *ContractIRelayRegistryCallerSession) RelayKeyToDispersers(relayId uint32) ([]uint32, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToDispersers(&_ContractIRelayRegistry.CallOpts, relayId)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 relayId) view returns(string)
func (_ContractIRelayRegistry *ContractIRelayRegistryCaller) RelayKeyToUrl(opts *bind.CallOpts, relayId uint32) (string, error) {
	var out []interface{}
	err := _ContractIRelayRegistry.contract.Call(opts, &out, "relayKeyToUrl", relayId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 relayId) view returns(string)
func (_ContractIRelayRegistry *ContractIRelayRegistrySession) RelayKeyToUrl(relayId uint32) (string, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToUrl(&_ContractIRelayRegistry.CallOpts, relayId)
}

// RelayKeyToUrl is a free data retrieval call binding the contract method 0x631eabb8.
//
// Solidity: function relayKeyToUrl(uint32 relayId) view returns(string)
func (_ContractIRelayRegistry *ContractIRelayRegistryCallerSession) RelayKeyToUrl(relayId uint32) (string, error) {
	return _ContractIRelayRegistry.Contract.RelayKeyToUrl(&_ContractIRelayRegistry.CallOpts, relayId)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x060176ba.
//
// Solidity: function addRelayInfo(address relay, string url, uint32[] dispersers) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactor) AddRelayInfo(opts *bind.TransactOpts, relay common.Address, url string, dispersers []uint32) (*types.Transaction, error) {
	return _ContractIRelayRegistry.contract.Transact(opts, "addRelayInfo", relay, url, dispersers)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x060176ba.
//
// Solidity: function addRelayInfo(address relay, string url, uint32[] dispersers) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistrySession) AddRelayInfo(relay common.Address, url string, dispersers []uint32) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.AddRelayInfo(&_ContractIRelayRegistry.TransactOpts, relay, url, dispersers)
}

// AddRelayInfo is a paid mutator transaction binding the contract method 0x060176ba.
//
// Solidity: function addRelayInfo(address relay, string url, uint32[] dispersers) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactorSession) AddRelayInfo(relay common.Address, url string, dispersers []uint32) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.AddRelayInfo(&_ContractIRelayRegistry.TransactOpts, relay, url, dispersers)
}

// AddRelayInfo0 is a paid mutator transaction binding the contract method 0x1c5d99ae.
//
// Solidity: function addRelayInfo(address relay, string url) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactor) AddRelayInfo0(opts *bind.TransactOpts, relay common.Address, url string) (*types.Transaction, error) {
	return _ContractIRelayRegistry.contract.Transact(opts, "addRelayInfo0", relay, url)
}

// AddRelayInfo0 is a paid mutator transaction binding the contract method 0x1c5d99ae.
//
// Solidity: function addRelayInfo(address relay, string url) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistrySession) AddRelayInfo0(relay common.Address, url string) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.AddRelayInfo0(&_ContractIRelayRegistry.TransactOpts, relay, url)
}

// AddRelayInfo0 is a paid mutator transaction binding the contract method 0x1c5d99ae.
//
// Solidity: function addRelayInfo(address relay, string url) returns(uint32)
func (_ContractIRelayRegistry *ContractIRelayRegistryTransactorSession) AddRelayInfo0(relay common.Address, url string) (*types.Transaction, error) {
	return _ContractIRelayRegistry.Contract.AddRelayInfo0(&_ContractIRelayRegistry.TransactOpts, relay, url)
}

// ContractIRelayRegistryRelayAddedIterator is returned from FilterRelayAdded and is used to iterate over the raw logs and unpacked data for RelayAdded events raised by the ContractIRelayRegistry contract.
type ContractIRelayRegistryRelayAddedIterator struct {
	Event *ContractIRelayRegistryRelayAdded // Event containing the contract specifics and raw log

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
func (it *ContractIRelayRegistryRelayAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIRelayRegistryRelayAdded)
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
		it.Event = new(ContractIRelayRegistryRelayAdded)
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
func (it *ContractIRelayRegistryRelayAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIRelayRegistryRelayAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIRelayRegistryRelayAdded represents a RelayAdded event raised by the ContractIRelayRegistry contract.
type ContractIRelayRegistryRelayAdded struct {
	RelayId    uint32
	Relay      common.Address
	Url        string
	Dispersers []uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRelayAdded is a free log retrieval operation binding the contract event 0x83a4f08e585cb8e18d60a77f861551c5791eea59d8d20a039f12b92b83794c84.
//
// Solidity: event RelayAdded(uint32 indexed relayId, address indexed relay, string url, uint32[] dispersers)
func (_ContractIRelayRegistry *ContractIRelayRegistryFilterer) FilterRelayAdded(opts *bind.FilterOpts, relayId []uint32, relay []common.Address) (*ContractIRelayRegistryRelayAddedIterator, error) {

	var relayIdRule []interface{}
	for _, relayIdItem := range relayId {
		relayIdRule = append(relayIdRule, relayIdItem)
	}
	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}

	logs, sub, err := _ContractIRelayRegistry.contract.FilterLogs(opts, "RelayAdded", relayIdRule, relayRule)
	if err != nil {
		return nil, err
	}
	return &ContractIRelayRegistryRelayAddedIterator{contract: _ContractIRelayRegistry.contract, event: "RelayAdded", logs: logs, sub: sub}, nil
}

// WatchRelayAdded is a free log subscription operation binding the contract event 0x83a4f08e585cb8e18d60a77f861551c5791eea59d8d20a039f12b92b83794c84.
//
// Solidity: event RelayAdded(uint32 indexed relayId, address indexed relay, string url, uint32[] dispersers)
func (_ContractIRelayRegistry *ContractIRelayRegistryFilterer) WatchRelayAdded(opts *bind.WatchOpts, sink chan<- *ContractIRelayRegistryRelayAdded, relayId []uint32, relay []common.Address) (event.Subscription, error) {

	var relayIdRule []interface{}
	for _, relayIdItem := range relayId {
		relayIdRule = append(relayIdRule, relayIdItem)
	}
	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}

	logs, sub, err := _ContractIRelayRegistry.contract.WatchLogs(opts, "RelayAdded", relayIdRule, relayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIRelayRegistryRelayAdded)
				if err := _ContractIRelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
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

// ParseRelayAdded is a log parse operation binding the contract event 0x83a4f08e585cb8e18d60a77f861551c5791eea59d8d20a039f12b92b83794c84.
//
// Solidity: event RelayAdded(uint32 indexed relayId, address indexed relay, string url, uint32[] dispersers)
func (_ContractIRelayRegistry *ContractIRelayRegistryFilterer) ParseRelayAdded(log types.Log) (*ContractIRelayRegistryRelayAdded, error) {
	event := new(ContractIRelayRegistryRelayAdded)
	if err := _ContractIRelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
