// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEjectionManager

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

// IEjectionManagerQuorumEjectionParams is an auto generated low-level Go binding around an user-defined struct.
type IEjectionManagerQuorumEjectionParams struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}

// ContractEjectionManagerMetaData contains all meta data concerning the ContractEjectionManager contract.
var ContractEjectionManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"_registryCoordinator\",\"type\":\"address\"},{\"internalType\":\"contractIStakeRegistry\",\"name\":\"_stakeRegistry\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"ejector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"EjectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"FailedOperatorEjection\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"operatorId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"}],\"name\":\"OperatorEjected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"rateLimitWindow\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"ejectableStakePercent\",\"type\":\"uint16\"}],\"name\":\"QuorumEjectionParamsSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_quorumNumber\",\"type\":\"uint8\"}],\"name\":\"amountEjectableForQuorum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[][]\",\"name\":\"_operatorIds\",\"type\":\"bytes32[][]\"}],\"name\":\"ejectOperators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_ejectors\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"rateLimitWindow\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ejectableStakePercent\",\"type\":\"uint16\"}],\"internalType\":\"structIEjectionManager.QuorumEjectionParams[]\",\"name\":\"_quorumEjectionParams\",\"type\":\"tuple[]\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isEjector\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"quorumEjectionParams\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rateLimitWindow\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ejectableStakePercent\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registryCoordinator\",\"outputs\":[{\"internalType\":\"contractIRegistryCoordinator\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_ejector\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_status\",\"type\":\"bool\"}],\"name\":\"setEjector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_quorumNumber\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"rateLimitWindow\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"ejectableStakePercent\",\"type\":\"uint16\"}],\"internalType\":\"structIEjectionManager.QuorumEjectionParams\",\"name\":\"_quorumEjectionParams\",\"type\":\"tuple\"}],\"name\":\"setQuorumEjectionParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stakeEjectedForQuorum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakeEjected\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeRegistry\",\"outputs\":[{\"internalType\":\"contractIStakeRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ContractEjectionManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEjectionManagerMetaData.ABI instead.
var ContractEjectionManagerABI = ContractEjectionManagerMetaData.ABI

// ContractEjectionManager is an auto generated Go binding around an Ethereum contract.
type ContractEjectionManager struct {
	ContractEjectionManagerCaller     // Read-only binding to the contract
	ContractEjectionManagerTransactor // Write-only binding to the contract
	ContractEjectionManagerFilterer   // Log filterer for contract events
}

// ContractEjectionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEjectionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEjectionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEjectionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEjectionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEjectionManagerSession struct {
	Contract     *ContractEjectionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ContractEjectionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEjectionManagerCallerSession struct {
	Contract *ContractEjectionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// ContractEjectionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEjectionManagerTransactorSession struct {
	Contract     *ContractEjectionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ContractEjectionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEjectionManagerRaw struct {
	Contract *ContractEjectionManager // Generic contract binding to access the raw methods on
}

// ContractEjectionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEjectionManagerCallerRaw struct {
	Contract *ContractEjectionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEjectionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEjectionManagerTransactorRaw struct {
	Contract *ContractEjectionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEjectionManager creates a new instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManager(address common.Address, backend bind.ContractBackend) (*ContractEjectionManager, error) {
	contract, err := bindContractEjectionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManager{ContractEjectionManagerCaller: ContractEjectionManagerCaller{contract: contract}, ContractEjectionManagerTransactor: ContractEjectionManagerTransactor{contract: contract}, ContractEjectionManagerFilterer: ContractEjectionManagerFilterer{contract: contract}}, nil
}

// NewContractEjectionManagerCaller creates a new read-only instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractEjectionManagerCaller, error) {
	contract, err := bindContractEjectionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerCaller{contract: contract}, nil
}

// NewContractEjectionManagerTransactor creates a new write-only instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEjectionManagerTransactor, error) {
	contract, err := bindContractEjectionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerTransactor{contract: contract}, nil
}

// NewContractEjectionManagerFilterer creates a new log filterer instance of ContractEjectionManager, bound to a specific deployed contract.
func NewContractEjectionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEjectionManagerFilterer, error) {
	contract, err := bindContractEjectionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerFilterer{contract: contract}, nil
}

// bindContractEjectionManager binds a generic wrapper to an already deployed contract.
func bindContractEjectionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEjectionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEjectionManager.Contract.ContractEjectionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.ContractEjectionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEjectionManager *ContractEjectionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.ContractEjectionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEjectionManager *ContractEjectionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEjectionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEjectionManager *ContractEjectionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEjectionManager *ContractEjectionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.contract.Transact(opts, method, params...)
}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerCaller) AmountEjectableForQuorum(opts *bind.CallOpts, _quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "amountEjectableForQuorum", _quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerSession) AmountEjectableForQuorum(_quorumNumber uint8) (*big.Int, error) {
	return _ContractEjectionManager.Contract.AmountEjectableForQuorum(&_ContractEjectionManager.CallOpts, _quorumNumber)
}

// AmountEjectableForQuorum is a free data retrieval call binding the contract method 0xb13f4504.
//
// Solidity: function amountEjectableForQuorum(uint8 _quorumNumber) view returns(uint256)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) AmountEjectableForQuorum(_quorumNumber uint8) (*big.Int, error) {
	return _ContractEjectionManager.Contract.AmountEjectableForQuorum(&_ContractEjectionManager.CallOpts, _quorumNumber)
}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerCaller) IsEjector(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "isEjector", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerSession) IsEjector(arg0 common.Address) (bool, error) {
	return _ContractEjectionManager.Contract.IsEjector(&_ContractEjectionManager.CallOpts, arg0)
}

// IsEjector is a free data retrieval call binding the contract method 0x6c08a879.
//
// Solidity: function isEjector(address ) view returns(bool)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) IsEjector(arg0 common.Address) (bool, error) {
	return _ContractEjectionManager.Contract.IsEjector(&_ContractEjectionManager.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) Owner() (common.Address, error) {
	return _ContractEjectionManager.Contract.Owner(&_ContractEjectionManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) Owner() (common.Address, error) {
	return _ContractEjectionManager.Contract.Owner(&_ContractEjectionManager.CallOpts)
}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerCaller) QuorumEjectionParams(opts *bind.CallOpts, arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "quorumEjectionParams", arg0)

	outstruct := new(struct {
		RateLimitWindow       uint32
		EjectableStakePercent uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RateLimitWindow = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.EjectableStakePercent = *abi.ConvertType(out[1], new(uint16)).(*uint16)

	return *outstruct, err

}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerSession) QuorumEjectionParams(arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	return _ContractEjectionManager.Contract.QuorumEjectionParams(&_ContractEjectionManager.CallOpts, arg0)
}

// QuorumEjectionParams is a free data retrieval call binding the contract method 0x00482569.
//
// Solidity: function quorumEjectionParams(uint8 ) view returns(uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) QuorumEjectionParams(arg0 uint8) (struct {
	RateLimitWindow       uint32
	EjectableStakePercent uint16
}, error) {
	return _ContractEjectionManager.Contract.QuorumEjectionParams(&_ContractEjectionManager.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEjectionManager.Contract.RegistryCoordinator(&_ContractEjectionManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEjectionManager.Contract.RegistryCoordinator(&_ContractEjectionManager.CallOpts)
}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerCaller) StakeEjectedForQuorum(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "stakeEjectedForQuorum", arg0, arg1)

	outstruct := new(struct {
		Timestamp    *big.Int
		StakeEjected *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.StakeEjected = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerSession) StakeEjectedForQuorum(arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	return _ContractEjectionManager.Contract.StakeEjectedForQuorum(&_ContractEjectionManager.CallOpts, arg0, arg1)
}

// StakeEjectedForQuorum is a free data retrieval call binding the contract method 0x3a0b0ddd.
//
// Solidity: function stakeEjectedForQuorum(uint8 , uint256 ) view returns(uint256 timestamp, uint256 stakeEjected)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) StakeEjectedForQuorum(arg0 uint8, arg1 *big.Int) (struct {
	Timestamp    *big.Int
	StakeEjected *big.Int
}, error) {
	return _ContractEjectionManager.Contract.StakeEjectedForQuorum(&_ContractEjectionManager.CallOpts, arg0, arg1)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEjectionManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerSession) StakeRegistry() (common.Address, error) {
	return _ContractEjectionManager.Contract.StakeRegistry(&_ContractEjectionManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEjectionManager *ContractEjectionManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractEjectionManager.Contract.StakeRegistry(&_ContractEjectionManager.CallOpts)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) EjectOperators(opts *bind.TransactOpts, _operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "ejectOperators", _operatorIds)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) EjectOperators(_operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.EjectOperators(&_ContractEjectionManager.TransactOpts, _operatorIds)
}

// EjectOperators is a paid mutator transaction binding the contract method 0x0a0593d1.
//
// Solidity: function ejectOperators(bytes32[][] _operatorIds) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) EjectOperators(_operatorIds [][][32]byte) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.EjectOperators(&_ContractEjectionManager.TransactOpts, _operatorIds)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "initialize", _owner, _ejectors, _quorumEjectionParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) Initialize(_owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.Initialize(&_ContractEjectionManager.TransactOpts, _owner, _ejectors, _quorumEjectionParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x8b88a024.
//
// Solidity: function initialize(address _owner, address[] _ejectors, (uint32,uint16)[] _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) Initialize(_owner common.Address, _ejectors []common.Address, _quorumEjectionParams []IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.Initialize(&_ContractEjectionManager.TransactOpts, _owner, _ejectors, _quorumEjectionParams)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.RenounceOwnership(&_ContractEjectionManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.RenounceOwnership(&_ContractEjectionManager.TransactOpts)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) SetEjector(opts *bind.TransactOpts, _ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "setEjector", _ejector, _status)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) SetEjector(_ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetEjector(&_ContractEjectionManager.TransactOpts, _ejector, _status)
}

// SetEjector is a paid mutator transaction binding the contract method 0x10ea4f8a.
//
// Solidity: function setEjector(address _ejector, bool _status) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) SetEjector(_ejector common.Address, _status bool) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetEjector(&_ContractEjectionManager.TransactOpts, _ejector, _status)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) SetQuorumEjectionParams(opts *bind.TransactOpts, _quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "setQuorumEjectionParams", _quorumNumber, _quorumEjectionParams)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) SetQuorumEjectionParams(_quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetQuorumEjectionParams(&_ContractEjectionManager.TransactOpts, _quorumNumber, _quorumEjectionParams)
}

// SetQuorumEjectionParams is a paid mutator transaction binding the contract method 0x77d17586.
//
// Solidity: function setQuorumEjectionParams(uint8 _quorumNumber, (uint32,uint16) _quorumEjectionParams) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) SetQuorumEjectionParams(_quorumNumber uint8, _quorumEjectionParams IEjectionManagerQuorumEjectionParams) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.SetQuorumEjectionParams(&_ContractEjectionManager.TransactOpts, _quorumNumber, _quorumEjectionParams)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.TransferOwnership(&_ContractEjectionManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEjectionManager *ContractEjectionManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEjectionManager.Contract.TransferOwnership(&_ContractEjectionManager.TransactOpts, newOwner)
}

// ContractEjectionManagerEjectorUpdatedIterator is returned from FilterEjectorUpdated and is used to iterate over the raw logs and unpacked data for EjectorUpdated events raised by the ContractEjectionManager contract.
type ContractEjectionManagerEjectorUpdatedIterator struct {
	Event *ContractEjectionManagerEjectorUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerEjectorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerEjectorUpdated)
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
		it.Event = new(ContractEjectionManagerEjectorUpdated)
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
func (it *ContractEjectionManagerEjectorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerEjectorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerEjectorUpdated represents a EjectorUpdated event raised by the ContractEjectionManager contract.
type ContractEjectionManagerEjectorUpdated struct {
	Ejector common.Address
	Status  bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEjectorUpdated is a free log retrieval operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterEjectorUpdated(opts *bind.FilterOpts) (*ContractEjectionManagerEjectorUpdatedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerEjectorUpdatedIterator{contract: _ContractEjectionManager.contract, event: "EjectorUpdated", logs: logs, sub: sub}, nil
}

// WatchEjectorUpdated is a free log subscription operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchEjectorUpdated(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerEjectorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerEjectorUpdated)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
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

// ParseEjectorUpdated is a log parse operation binding the contract event 0x7676686b6d22e112412bd874d70177e011ab06602c26063f19f0386c9a3cee42.
//
// Solidity: event EjectorUpdated(address ejector, bool status)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseEjectorUpdated(log types.Log) (*ContractEjectionManagerEjectorUpdated, error) {
	event := new(ContractEjectionManagerEjectorUpdated)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerFailedOperatorEjectionIterator is returned from FilterFailedOperatorEjection and is used to iterate over the raw logs and unpacked data for FailedOperatorEjection events raised by the ContractEjectionManager contract.
type ContractEjectionManagerFailedOperatorEjectionIterator struct {
	Event *ContractEjectionManagerFailedOperatorEjection // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerFailedOperatorEjectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerFailedOperatorEjection)
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
		it.Event = new(ContractEjectionManagerFailedOperatorEjection)
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
func (it *ContractEjectionManagerFailedOperatorEjectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerFailedOperatorEjectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerFailedOperatorEjection represents a FailedOperatorEjection event raised by the ContractEjectionManager contract.
type ContractEjectionManagerFailedOperatorEjection struct {
	OperatorId   [32]byte
	QuorumNumber uint8
	Err          []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFailedOperatorEjection is a free log retrieval operation binding the contract event 0xae1dcabe5fd19643522b5e06189c4f844e5ba5d3bf0c17e47a22c68dd585b6ef.
//
// Solidity: event FailedOperatorEjection(bytes32 operatorId, uint8 quorumNumber, bytes err)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterFailedOperatorEjection(opts *bind.FilterOpts) (*ContractEjectionManagerFailedOperatorEjectionIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "FailedOperatorEjection")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerFailedOperatorEjectionIterator{contract: _ContractEjectionManager.contract, event: "FailedOperatorEjection", logs: logs, sub: sub}, nil
}

// WatchFailedOperatorEjection is a free log subscription operation binding the contract event 0xae1dcabe5fd19643522b5e06189c4f844e5ba5d3bf0c17e47a22c68dd585b6ef.
//
// Solidity: event FailedOperatorEjection(bytes32 operatorId, uint8 quorumNumber, bytes err)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchFailedOperatorEjection(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerFailedOperatorEjection) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "FailedOperatorEjection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerFailedOperatorEjection)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "FailedOperatorEjection", log); err != nil {
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

// ParseFailedOperatorEjection is a log parse operation binding the contract event 0xae1dcabe5fd19643522b5e06189c4f844e5ba5d3bf0c17e47a22c68dd585b6ef.
//
// Solidity: event FailedOperatorEjection(bytes32 operatorId, uint8 quorumNumber, bytes err)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseFailedOperatorEjection(log types.Log) (*ContractEjectionManagerFailedOperatorEjection, error) {
	event := new(ContractEjectionManagerFailedOperatorEjection)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "FailedOperatorEjection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEjectionManager contract.
type ContractEjectionManagerInitializedIterator struct {
	Event *ContractEjectionManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerInitialized)
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
		it.Event = new(ContractEjectionManagerInitialized)
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
func (it *ContractEjectionManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerInitialized represents a Initialized event raised by the ContractEjectionManager contract.
type ContractEjectionManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEjectionManagerInitializedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerInitializedIterator{contract: _ContractEjectionManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerInitialized)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseInitialized(log types.Log) (*ContractEjectionManagerInitialized, error) {
	event := new(ContractEjectionManagerInitialized)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerOperatorEjectedIterator is returned from FilterOperatorEjected and is used to iterate over the raw logs and unpacked data for OperatorEjected events raised by the ContractEjectionManager contract.
type ContractEjectionManagerOperatorEjectedIterator struct {
	Event *ContractEjectionManagerOperatorEjected // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerOperatorEjectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerOperatorEjected)
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
		it.Event = new(ContractEjectionManagerOperatorEjected)
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
func (it *ContractEjectionManagerOperatorEjectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerOperatorEjectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerOperatorEjected represents a OperatorEjected event raised by the ContractEjectionManager contract.
type ContractEjectionManagerOperatorEjected struct {
	OperatorId   [32]byte
	QuorumNumber uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOperatorEjected is a free log retrieval operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterOperatorEjected(opts *bind.FilterOpts) (*ContractEjectionManagerOperatorEjectedIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "OperatorEjected")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerOperatorEjectedIterator{contract: _ContractEjectionManager.contract, event: "OperatorEjected", logs: logs, sub: sub}, nil
}

// WatchOperatorEjected is a free log subscription operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchOperatorEjected(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerOperatorEjected) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "OperatorEjected")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerOperatorEjected)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "OperatorEjected", log); err != nil {
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

// ParseOperatorEjected is a log parse operation binding the contract event 0x97ddb711c61a9d2d7effcba3e042a33862297f898d555655cca39ec4451f53b4.
//
// Solidity: event OperatorEjected(bytes32 operatorId, uint8 quorumNumber)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseOperatorEjected(log types.Log) (*ContractEjectionManagerOperatorEjected, error) {
	event := new(ContractEjectionManagerOperatorEjected)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "OperatorEjected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEjectionManager contract.
type ContractEjectionManagerOwnershipTransferredIterator struct {
	Event *ContractEjectionManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerOwnershipTransferred)
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
		it.Event = new(ContractEjectionManagerOwnershipTransferred)
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
func (it *ContractEjectionManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEjectionManager contract.
type ContractEjectionManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEjectionManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerOwnershipTransferredIterator{contract: _ContractEjectionManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerOwnershipTransferred)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEjectionManagerOwnershipTransferred, error) {
	event := new(ContractEjectionManagerOwnershipTransferred)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEjectionManagerQuorumEjectionParamsSetIterator is returned from FilterQuorumEjectionParamsSet and is used to iterate over the raw logs and unpacked data for QuorumEjectionParamsSet events raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjectionParamsSetIterator struct {
	Event *ContractEjectionManagerQuorumEjectionParamsSet // Event containing the contract specifics and raw log

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
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEjectionManagerQuorumEjectionParamsSet)
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
		it.Event = new(ContractEjectionManagerQuorumEjectionParamsSet)
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
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEjectionManagerQuorumEjectionParamsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEjectionManagerQuorumEjectionParamsSet represents a QuorumEjectionParamsSet event raised by the ContractEjectionManager contract.
type ContractEjectionManagerQuorumEjectionParamsSet struct {
	QuorumNumber          uint8
	RateLimitWindow       uint32
	EjectableStakePercent uint16
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterQuorumEjectionParamsSet is a free log retrieval operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) FilterQuorumEjectionParamsSet(opts *bind.FilterOpts) (*ContractEjectionManagerQuorumEjectionParamsSetIterator, error) {

	logs, sub, err := _ContractEjectionManager.contract.FilterLogs(opts, "QuorumEjectionParamsSet")
	if err != nil {
		return nil, err
	}
	return &ContractEjectionManagerQuorumEjectionParamsSetIterator{contract: _ContractEjectionManager.contract, event: "QuorumEjectionParamsSet", logs: logs, sub: sub}, nil
}

// WatchQuorumEjectionParamsSet is a free log subscription operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) WatchQuorumEjectionParamsSet(opts *bind.WatchOpts, sink chan<- *ContractEjectionManagerQuorumEjectionParamsSet) (event.Subscription, error) {

	logs, sub, err := _ContractEjectionManager.contract.WatchLogs(opts, "QuorumEjectionParamsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEjectionManagerQuorumEjectionParamsSet)
				if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjectionParamsSet", log); err != nil {
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

// ParseQuorumEjectionParamsSet is a log parse operation binding the contract event 0xe69c2827a1e2fdd32265ebb4eeea5ee564f0551cf5dfed4150f8e116a67209eb.
//
// Solidity: event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent)
func (_ContractEjectionManager *ContractEjectionManagerFilterer) ParseQuorumEjectionParamsSet(log types.Log) (*ContractEjectionManagerQuorumEjectionParamsSet, error) {
	event := new(ContractEjectionManagerQuorumEjectionParamsSet)
	if err := _ContractEjectionManager.contract.UnpackLog(event, "QuorumEjectionParamsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
