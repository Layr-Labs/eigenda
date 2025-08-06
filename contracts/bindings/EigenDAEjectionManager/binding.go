// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAEjectionManager

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

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// ContractEigenDAEjectionManagerMetaData contains all meta data concerning the ContractEigenDAEjectionManager contract.
var ContractEigenDAEjectionManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"depositToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositAmount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"addressDirectory_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"estimatedGasUsed_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addEjectorBalance\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjection\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjectionByEjector\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjectionWithSig\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeEjection\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectionCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionInitiated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionTime\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDepositAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDepositToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastEjectionInitiated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setCooldown\",\"inputs\":[{\"name\":\"cooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelay\",\"inputs\":[{\"name\":\"delay\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startEjection\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// ContractEigenDAEjectionManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAEjectionManagerMetaData.ABI instead.
var ContractEigenDAEjectionManagerABI = ContractEigenDAEjectionManagerMetaData.ABI

// ContractEigenDAEjectionManager is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAEjectionManager struct {
	ContractEigenDAEjectionManagerCaller     // Read-only binding to the contract
	ContractEigenDAEjectionManagerTransactor // Write-only binding to the contract
	ContractEigenDAEjectionManagerFilterer   // Log filterer for contract events
}

// ContractEigenDAEjectionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAEjectionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAEjectionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAEjectionManagerSession struct {
	Contract     *ContractEigenDAEjectionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractEigenDAEjectionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAEjectionManagerCallerSession struct {
	Contract *ContractEigenDAEjectionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractEigenDAEjectionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAEjectionManagerTransactorSession struct {
	Contract     *ContractEigenDAEjectionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractEigenDAEjectionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAEjectionManagerRaw struct {
	Contract *ContractEigenDAEjectionManager // Generic contract binding to access the raw methods on
}

// ContractEigenDAEjectionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionManagerCallerRaw struct {
	Contract *ContractEigenDAEjectionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAEjectionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAEjectionManagerTransactorRaw struct {
	Contract *ContractEigenDAEjectionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAEjectionManager creates a new instance of ContractEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractEigenDAEjectionManager(address common.Address, backend bind.ContractBackend) (*ContractEigenDAEjectionManager, error) {
	contract, err := bindContractEigenDAEjectionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionManager{ContractEigenDAEjectionManagerCaller: ContractEigenDAEjectionManagerCaller{contract: contract}, ContractEigenDAEjectionManagerTransactor: ContractEigenDAEjectionManagerTransactor{contract: contract}, ContractEigenDAEjectionManagerFilterer: ContractEigenDAEjectionManagerFilterer{contract: contract}}, nil
}

// NewContractEigenDAEjectionManagerCaller creates a new read-only instance of ContractEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractEigenDAEjectionManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAEjectionManagerCaller, error) {
	contract, err := bindContractEigenDAEjectionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionManagerCaller{contract: contract}, nil
}

// NewContractEigenDAEjectionManagerTransactor creates a new write-only instance of ContractEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractEigenDAEjectionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAEjectionManagerTransactor, error) {
	contract, err := bindContractEigenDAEjectionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionManagerTransactor{contract: contract}, nil
}

// NewContractEigenDAEjectionManagerFilterer creates a new log filterer instance of ContractEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractEigenDAEjectionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAEjectionManagerFilterer, error) {
	contract, err := bindContractEigenDAEjectionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAEjectionManagerFilterer{contract: contract}, nil
}

// bindContractEigenDAEjectionManager binds a generic wrapper to an already deployed contract.
func bindContractEigenDAEjectionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAEjectionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAEjectionManager.Contract.ContractEigenDAEjectionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.ContractEigenDAEjectionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.ContractEigenDAEjectionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAEjectionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.contract.Transact(opts, method, params...)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) EjectionCooldown(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "ejectionCooldown")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) EjectionCooldown() (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionCooldown(&_ContractEigenDAEjectionManager.CallOpts)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) EjectionCooldown() (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionCooldown(&_ContractEigenDAEjectionManager.CallOpts)
}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) EjectionDelay(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "ejectionDelay")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) EjectionDelay() (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionDelay(&_ContractEigenDAEjectionManager.CallOpts)
}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) EjectionDelay() (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionDelay(&_ContractEigenDAEjectionManager.CallOpts)
}

// EjectionInitiated is a free data retrieval call binding the contract method 0x85b62b1a.
//
// Solidity: function ejectionInitiated(address operator) view returns(bool)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) EjectionInitiated(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "ejectionInitiated", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// EjectionInitiated is a free data retrieval call binding the contract method 0x85b62b1a.
//
// Solidity: function ejectionInitiated(address operator) view returns(bool)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) EjectionInitiated(operator common.Address) (bool, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionInitiated(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// EjectionInitiated is a free data retrieval call binding the contract method 0x85b62b1a.
//
// Solidity: function ejectionInitiated(address operator) view returns(bool)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) EjectionInitiated(operator common.Address) (bool, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionInitiated(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) EjectionQuorums(opts *bind.CallOpts, operator common.Address) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "ejectionQuorums", operator)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) EjectionQuorums(operator common.Address) ([]byte, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionQuorums(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) EjectionQuorums(operator common.Address) ([]byte, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionQuorums(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) EjectionTime(opts *bind.CallOpts, operator common.Address) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "ejectionTime", operator)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) EjectionTime(operator common.Address) (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionTime(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) EjectionTime(operator common.Address) (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.EjectionTime(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// GetDepositAmount is a free data retrieval call binding the contract method 0x7d96f693.
//
// Solidity: function getDepositAmount() view returns(uint256)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) GetDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "getDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositAmount is a free data retrieval call binding the contract method 0x7d96f693.
//
// Solidity: function getDepositAmount() view returns(uint256)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) GetDepositAmount() (*big.Int, error) {
	return _ContractEigenDAEjectionManager.Contract.GetDepositAmount(&_ContractEigenDAEjectionManager.CallOpts)
}

// GetDepositAmount is a free data retrieval call binding the contract method 0x7d96f693.
//
// Solidity: function getDepositAmount() view returns(uint256)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) GetDepositAmount() (*big.Int, error) {
	return _ContractEigenDAEjectionManager.Contract.GetDepositAmount(&_ContractEigenDAEjectionManager.CallOpts)
}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) GetDepositToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "getDepositToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) GetDepositToken() (common.Address, error) {
	return _ContractEigenDAEjectionManager.Contract.GetDepositToken(&_ContractEigenDAEjectionManager.CallOpts)
}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) GetDepositToken() (common.Address, error) {
	return _ContractEigenDAEjectionManager.Contract.GetDepositToken(&_ContractEigenDAEjectionManager.CallOpts)
}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCaller) LastEjectionInitiated(opts *bind.CallOpts, operator common.Address) (uint64, error) {
	var out []interface{}
	err := _ContractEigenDAEjectionManager.contract.Call(opts, &out, "lastEjectionInitiated", operator)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) LastEjectionInitiated(operator common.Address) (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.LastEjectionInitiated(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerCallerSession) LastEjectionInitiated(operator common.Address) (uint64, error) {
	return _ContractEigenDAEjectionManager.Contract.LastEjectionInitiated(&_ContractEigenDAEjectionManager.CallOpts, operator)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) AddEjectorBalance(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "addEjectorBalance", amount)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) AddEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.AddEjectorBalance(&_ContractEigenDAEjectionManager.TransactOpts, amount)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) AddEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.AddEjectorBalance(&_ContractEigenDAEjectionManager.TransactOpts, amount)
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) CancelEjection(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "cancelEjection")
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) CancelEjection() (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjection(&_ContractEigenDAEjectionManager.TransactOpts)
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) CancelEjection() (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjection(&_ContractEigenDAEjectionManager.TransactOpts)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) CancelEjectionByEjector(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "cancelEjectionByEjector", operator)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) CancelEjectionByEjector(operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjectionByEjector(&_ContractEigenDAEjectionManager.TransactOpts, operator)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) CancelEjectionByEjector(operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjectionByEjector(&_ContractEigenDAEjectionManager.TransactOpts, operator)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) CancelEjectionWithSig(opts *bind.TransactOpts, operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "cancelEjectionWithSig", operator, apkG2, sigma, recipient)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) CancelEjectionWithSig(operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjectionWithSig(&_ContractEigenDAEjectionManager.TransactOpts, operator, apkG2, sigma, recipient)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) CancelEjectionWithSig(operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CancelEjectionWithSig(&_ContractEigenDAEjectionManager.TransactOpts, operator, apkG2, sigma, recipient)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) CompleteEjection(opts *bind.TransactOpts, operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "completeEjection", operator, quorums)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) CompleteEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CompleteEjection(&_ContractEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) CompleteEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.CompleteEjection(&_ContractEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) SetCooldown(opts *bind.TransactOpts, cooldown uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "setCooldown", cooldown)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) SetCooldown(cooldown uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.SetCooldown(&_ContractEigenDAEjectionManager.TransactOpts, cooldown)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) SetCooldown(cooldown uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.SetCooldown(&_ContractEigenDAEjectionManager.TransactOpts, cooldown)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) SetDelay(opts *bind.TransactOpts, delay uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "setDelay", delay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) SetDelay(delay uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.SetDelay(&_ContractEigenDAEjectionManager.TransactOpts, delay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) SetDelay(delay uint64) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.SetDelay(&_ContractEigenDAEjectionManager.TransactOpts, delay)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactor) StartEjection(opts *bind.TransactOpts, operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.contract.Transact(opts, "startEjection", operator, quorums)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerSession) StartEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.StartEjection(&_ContractEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractEigenDAEjectionManager *ContractEigenDAEjectionManagerTransactorSession) StartEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractEigenDAEjectionManager.Contract.StartEjection(&_ContractEigenDAEjectionManager.TransactOpts, operator, quorums)
}
