// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDARelayRegistry

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

// ContractEigenDARelayRegistryMetaData contains all meta data concerning the ContractEigenDARelayRegistry contract.
var ContractEigenDARelayRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRelayURL\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRelayAddress\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRelayKey\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRelayURL\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextRelayKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayAddressToKey\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayKeyToURL\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RelayAdded\",\"inputs\":[{\"name\":\"relay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"key\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"uint32\"},{\"name\":\"relayURL\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5061001961001e565b6100de565b600054610100900460ff161561008a5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff90811610156100dc576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b61098f806100ed6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80639561f9ea116100715780639561f9ea14610194578063b5a872da146101ba578063c4d66de8146101e3578063f2fde38b146101f6578063f6f9c09714610209578063fbe839a31461021c57600080fd5b806309ce3ce2146100b957806315ddaa5d14610102578063715018a6146101125780638050a8991461011c5780638898cf74146101635780638da5cb5b14610183575b600080fd5b6100e86100c7366004610777565b6001600160a01b031660009081526066602052604090205463ffffffff1690565b60405163ffffffff90911681526020015b60405180910390f35b6068546100e89063ffffffff1681565b61011a61022f565b005b61014b61012a366004610799565b63ffffffff166000908152606760205260409020546001600160a01b031690565b6040516001600160a01b0390911681526020016100f9565b610176610171366004610799565b610243565b6040516100f991906107bf565b6033546001600160a01b031661014b565b6100e86101a2366004610777565b60666020526000908152604090205463ffffffff1681565b61014b6101c8366004610799565b6067602052600090815260409020546001600160a01b031681565b61011a6101f1366004610777565b6102ec565b61011a610204366004610777565b610404565b610176610217366004610799565b61047d565b6100e861022a36600461082a565b610517565b610237610616565b6102416000610670565b565b63ffffffff81166000908152606560205260409020805460609190610267906108ec565b80601f0160208091040260200160405190810160405280929190818152602001828054610293906108ec565b80156102e05780601f106102b5576101008083540402835291602001916102e0565b820191906000526020600020905b8154815290600101906020018083116102c357829003601f168201915b50505050509050919050565b600054610100900460ff161580801561030c5750600054600160ff909116105b806103265750303b158015610326575060005460ff166001145b61038e5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff1916600117905580156103b1576000805461ff0019166101001790555b6103ba82610670565b8015610400576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b61040c610616565b6001600160a01b0381166104715760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610385565b61047a81610670565b50565b60656020526000908152604090208054610496906108ec565b80601f01602080910402602001604051908101604052809291908181526020018280546104c2906108ec565b801561050f5780601f106104e45761010080835404028352916020019161050f565b820191906000526020600020905b8154815290600101906020018083116104f257829003601f168201915b505050505081565b6000610521610616565b60685463ffffffff1660009081526065602090815260409091208351610549928501906106c2565b50606880546001600160a01b0385166000818152606660209081526040808320805463ffffffff191663ffffffff96871617905585548516835260679091529081902080546001600160a01b031916831790559254925192909116917f01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99906105d29086906107bf565b60405180910390a36068805463ffffffff169060006105f083610927565b91906101000a81548163ffffffff021916908363ffffffff160217905550905092915050565b6033546001600160a01b031633146102415760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610385565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b8280546106ce906108ec565b90600052602060002090601f0160209004810192826106f05760008555610736565b82601f1061070957805160ff1916838001178555610736565b82800160010185558215610736579182015b8281111561073657825182559160200191906001019061071b565b50610742929150610746565b5090565b5b808211156107425760008155600101610747565b80356001600160a01b038116811461077257600080fd5b919050565b60006020828403121561078957600080fd5b6107928261075b565b9392505050565b6000602082840312156107ab57600080fd5b813563ffffffff8116811461079257600080fd5b600060208083528351808285015260005b818110156107ec578581018301518582016040015282016107d0565b818111156107fe576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561083d57600080fd5b6108468361075b565b9150602083013567ffffffffffffffff8082111561086357600080fd5b818501915085601f83011261087757600080fd5b81358181111561088957610889610814565b604051601f8201601f19908116603f011681019083821181831017156108b1576108b1610814565b816040528281528860208487010111156108ca57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b600181811c9082168061090057607f821691505b6020821081141561092157634e487b7160e01b600052602260045260246000fd5b50919050565b600063ffffffff8083168181141561094f57634e487b7160e01b600052601160045260246000fd5b600101939250505056fea26469706673582212202b6b39b905ee133307fe7588cf884ba0c039027692990b4de147bf5fd2b5a1a664736f6c634300080c0033",
}

// ContractEigenDARelayRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDARelayRegistryMetaData.ABI instead.
var ContractEigenDARelayRegistryABI = ContractEigenDARelayRegistryMetaData.ABI

// ContractEigenDARelayRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDARelayRegistryMetaData.Bin instead.
var ContractEigenDARelayRegistryBin = ContractEigenDARelayRegistryMetaData.Bin

// DeployContractEigenDARelayRegistry deploys a new Ethereum contract, binding an instance of ContractEigenDARelayRegistry to it.
func DeployContractEigenDARelayRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContractEigenDARelayRegistry, error) {
	parsed, err := ContractEigenDARelayRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDARelayRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDARelayRegistry{ContractEigenDARelayRegistryCaller: ContractEigenDARelayRegistryCaller{contract: contract}, ContractEigenDARelayRegistryTransactor: ContractEigenDARelayRegistryTransactor{contract: contract}, ContractEigenDARelayRegistryFilterer: ContractEigenDARelayRegistryFilterer{contract: contract}}, nil
}

// ContractEigenDARelayRegistry is an auto generated Go binding around an Ethereum contract.
type ContractEigenDARelayRegistry struct {
	ContractEigenDARelayRegistryCaller     // Read-only binding to the contract
	ContractEigenDARelayRegistryTransactor // Write-only binding to the contract
	ContractEigenDARelayRegistryFilterer   // Log filterer for contract events
}

// ContractEigenDARelayRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDARelayRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDARelayRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDARelayRegistrySession struct {
	Contract     *ContractEigenDARelayRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractEigenDARelayRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDARelayRegistryCallerSession struct {
	Contract *ContractEigenDARelayRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// ContractEigenDARelayRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDARelayRegistryTransactorSession struct {
	Contract     *ContractEigenDARelayRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// ContractEigenDARelayRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryRaw struct {
	Contract *ContractEigenDARelayRegistry // Generic contract binding to access the raw methods on
}

// ContractEigenDARelayRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryCallerRaw struct {
	Contract *ContractEigenDARelayRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDARelayRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDARelayRegistryTransactorRaw struct {
	Contract *ContractEigenDARelayRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDARelayRegistry creates a new instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistry(address common.Address, backend bind.ContractBackend) (*ContractEigenDARelayRegistry, error) {
	contract, err := bindContractEigenDARelayRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistry{ContractEigenDARelayRegistryCaller: ContractEigenDARelayRegistryCaller{contract: contract}, ContractEigenDARelayRegistryTransactor: ContractEigenDARelayRegistryTransactor{contract: contract}, ContractEigenDARelayRegistryFilterer: ContractEigenDARelayRegistryFilterer{contract: contract}}, nil
}

// NewContractEigenDARelayRegistryCaller creates a new read-only instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDARelayRegistryCaller, error) {
	contract, err := bindContractEigenDARelayRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryCaller{contract: contract}, nil
}

// NewContractEigenDARelayRegistryTransactor creates a new write-only instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDARelayRegistryTransactor, error) {
	contract, err := bindContractEigenDARelayRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryTransactor{contract: contract}, nil
}

// NewContractEigenDARelayRegistryFilterer creates a new log filterer instance of ContractEigenDARelayRegistry, bound to a specific deployed contract.
func NewContractEigenDARelayRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDARelayRegistryFilterer, error) {
	contract, err := bindContractEigenDARelayRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryFilterer{contract: contract}, nil
}

// bindContractEigenDARelayRegistry binds a generic wrapper to an already deployed contract.
func bindContractEigenDARelayRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDARelayRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.ContractEigenDARelayRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDARelayRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetRelayAddress is a free data retrieval call binding the contract method 0x8050a899.
//
// Solidity: function getRelayAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) GetRelayAddress(opts *bind.CallOpts, key uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "getRelayAddress", key)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRelayAddress is a free data retrieval call binding the contract method 0x8050a899.
//
// Solidity: function getRelayAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) GetRelayAddress(key uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayAddress(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// GetRelayAddress is a free data retrieval call binding the contract method 0x8050a899.
//
// Solidity: function getRelayAddress(uint32 key) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) GetRelayAddress(key uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayAddress(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// GetRelayKey is a free data retrieval call binding the contract method 0x09ce3ce2.
//
// Solidity: function getRelayKey(address relay) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) GetRelayKey(opts *bind.CallOpts, relay common.Address) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "getRelayKey", relay)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetRelayKey is a free data retrieval call binding the contract method 0x09ce3ce2.
//
// Solidity: function getRelayKey(address relay) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) GetRelayKey(relay common.Address) (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayKey(&_ContractEigenDARelayRegistry.CallOpts, relay)
}

// GetRelayKey is a free data retrieval call binding the contract method 0x09ce3ce2.
//
// Solidity: function getRelayKey(address relay) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) GetRelayKey(relay common.Address) (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayKey(&_ContractEigenDARelayRegistry.CallOpts, relay)
}

// GetRelayURL is a free data retrieval call binding the contract method 0x8898cf74.
//
// Solidity: function getRelayURL(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) GetRelayURL(opts *bind.CallOpts, key uint32) (string, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "getRelayURL", key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRelayURL is a free data retrieval call binding the contract method 0x8898cf74.
//
// Solidity: function getRelayURL(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) GetRelayURL(key uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayURL(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// GetRelayURL is a free data retrieval call binding the contract method 0x8898cf74.
//
// Solidity: function getRelayURL(uint32 key) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) GetRelayURL(key uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.GetRelayURL(&_ContractEigenDARelayRegistry.CallOpts, key)
}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) NextRelayKey(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "nextRelayKey")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) NextRelayKey() (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.NextRelayKey(&_ContractEigenDARelayRegistry.CallOpts)
}

// NextRelayKey is a free data retrieval call binding the contract method 0x15ddaa5d.
//
// Solidity: function nextRelayKey() view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) NextRelayKey() (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.NextRelayKey(&_ContractEigenDARelayRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) Owner() (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.Owner(&_ContractEigenDARelayRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.Owner(&_ContractEigenDARelayRegistry.CallOpts)
}

// RelayAddressToKey is a free data retrieval call binding the contract method 0x9561f9ea.
//
// Solidity: function relayAddressToKey(address ) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayAddressToKey(opts *bind.CallOpts, arg0 common.Address) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayAddressToKey", arg0)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RelayAddressToKey is a free data retrieval call binding the contract method 0x9561f9ea.
//
// Solidity: function relayAddressToKey(address ) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayAddressToKey(arg0 common.Address) (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayAddressToKey(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayAddressToKey is a free data retrieval call binding the contract method 0x9561f9ea.
//
// Solidity: function relayAddressToKey(address ) view returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayAddressToKey(arg0 common.Address) (uint32, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayAddressToKey(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 ) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayKeyToAddress(opts *bind.CallOpts, arg0 uint32) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 ) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayKeyToAddress(arg0 uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToAddress is a free data retrieval call binding the contract method 0xb5a872da.
//
// Solidity: function relayKeyToAddress(uint32 ) view returns(address)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayKeyToAddress(arg0 uint32) (common.Address, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToAddress(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToURL is a free data retrieval call binding the contract method 0xf6f9c097.
//
// Solidity: function relayKeyToURL(uint32 ) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCaller) RelayKeyToURL(opts *bind.CallOpts, arg0 uint32) (string, error) {
	var out []interface{}
	err := _ContractEigenDARelayRegistry.contract.Call(opts, &out, "relayKeyToURL", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RelayKeyToURL is a free data retrieval call binding the contract method 0xf6f9c097.
//
// Solidity: function relayKeyToURL(uint32 ) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RelayKeyToURL(arg0 uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToURL(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// RelayKeyToURL is a free data retrieval call binding the contract method 0xf6f9c097.
//
// Solidity: function relayKeyToURL(uint32 ) view returns(string)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryCallerSession) RelayKeyToURL(arg0 uint32) (string, error) {
	return _ContractEigenDARelayRegistry.Contract.RelayKeyToURL(&_ContractEigenDARelayRegistry.CallOpts, arg0)
}

// AddRelayURL is a paid mutator transaction binding the contract method 0xfbe839a3.
//
// Solidity: function addRelayURL(address relay, string relayURL) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) AddRelayURL(opts *bind.TransactOpts, relay common.Address, relayURL string) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "addRelayURL", relay, relayURL)
}

// AddRelayURL is a paid mutator transaction binding the contract method 0xfbe839a3.
//
// Solidity: function addRelayURL(address relay, string relayURL) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) AddRelayURL(relay common.Address, relayURL string) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.AddRelayURL(&_ContractEigenDARelayRegistry.TransactOpts, relay, relayURL)
}

// AddRelayURL is a paid mutator transaction binding the contract method 0xfbe839a3.
//
// Solidity: function addRelayURL(address relay, string relayURL) returns(uint32)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) AddRelayURL(relay common.Address, relayURL string) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.AddRelayURL(&_ContractEigenDARelayRegistry.TransactOpts, relay, relayURL)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "initialize", _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.Initialize(&_ContractEigenDARelayRegistry.TransactOpts, _initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _initialOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) Initialize(_initialOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.Initialize(&_ContractEigenDARelayRegistry.TransactOpts, _initialOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.RenounceOwnership(&_ContractEigenDARelayRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.RenounceOwnership(&_ContractEigenDARelayRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.TransferOwnership(&_ContractEigenDARelayRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDARelayRegistry.Contract.TransferOwnership(&_ContractEigenDARelayRegistry.TransactOpts, newOwner)
}

// ContractEigenDARelayRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryInitializedIterator struct {
	Event *ContractEigenDARelayRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryInitialized)
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
		it.Event = new(ContractEigenDARelayRegistryInitialized)
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
func (it *ContractEigenDARelayRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryInitialized represents a Initialized event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDARelayRegistryInitializedIterator, error) {

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryInitializedIterator{contract: _ContractEigenDARelayRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryInitialized)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseInitialized(log types.Log) (*ContractEigenDARelayRegistryInitialized, error) {
	event := new(ContractEigenDARelayRegistryInitialized)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDARelayRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryOwnershipTransferredIterator struct {
	Event *ContractEigenDARelayRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryOwnershipTransferred)
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
		it.Event = new(ContractEigenDARelayRegistryOwnershipTransferred)
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
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDARelayRegistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryOwnershipTransferredIterator{contract: _ContractEigenDARelayRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryOwnershipTransferred)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDARelayRegistryOwnershipTransferred, error) {
	event := new(ContractEigenDARelayRegistryOwnershipTransferred)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDARelayRegistryRelayAddedIterator is returned from FilterRelayAdded and is used to iterate over the raw logs and unpacked data for RelayAdded events raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryRelayAddedIterator struct {
	Event *ContractEigenDARelayRegistryRelayAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDARelayRegistryRelayAdded)
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
		it.Event = new(ContractEigenDARelayRegistryRelayAdded)
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
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDARelayRegistryRelayAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDARelayRegistryRelayAdded represents a RelayAdded event raised by the ContractEigenDARelayRegistry contract.
type ContractEigenDARelayRegistryRelayAdded struct {
	Relay    common.Address
	Key      uint32
	RelayURL string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRelayAdded is a free log retrieval operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) FilterRelayAdded(opts *bind.FilterOpts, relay []common.Address, key []uint32) (*ContractEigenDARelayRegistryRelayAddedIterator, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.FilterLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDARelayRegistryRelayAddedIterator{contract: _ContractEigenDARelayRegistry.contract, event: "RelayAdded", logs: logs, sub: sub}, nil
}

// WatchRelayAdded is a free log subscription operation binding the contract event 0x01c289e409d41a712a615bf286126433da55c193bbe64fc8e77af5f1ff13db99.
//
// Solidity: event RelayAdded(address indexed relay, uint32 indexed key, string relayURL)
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) WatchRelayAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDARelayRegistryRelayAdded, relay []common.Address, key []uint32) (event.Subscription, error) {

	var relayRule []interface{}
	for _, relayItem := range relay {
		relayRule = append(relayRule, relayItem)
	}
	var keyRule []interface{}
	for _, keyItem := range key {
		keyRule = append(keyRule, keyItem)
	}

	logs, sub, err := _ContractEigenDARelayRegistry.contract.WatchLogs(opts, "RelayAdded", relayRule, keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDARelayRegistryRelayAdded)
				if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
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
func (_ContractEigenDARelayRegistry *ContractEigenDARelayRegistryFilterer) ParseRelayAdded(log types.Log) (*ContractEigenDARelayRegistryRelayAdded, error) {
	event := new(ContractEigenDARelayRegistryRelayAdded)
	if err := _ContractEigenDARelayRegistry.contract.UnpackLog(event, "RelayAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
