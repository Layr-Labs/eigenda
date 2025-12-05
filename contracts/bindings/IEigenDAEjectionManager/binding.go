// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDAEjectionManager

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

// ContractIEigenDAEjectionManagerMetaData contains all meta data concerning the ContractIEigenDAEjectionManager contract.
var ContractIEigenDAEjectionManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addEjectorBalance\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjection\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjectionByEjector\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelEjectionWithSig\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeEjection\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectionCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejectionTime\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDepositToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEjector\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEjectorBalance\",\"inputs\":[{\"name\":\"ejector\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastEjectionInitiated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setCooldown\",\"inputs\":[{\"name\":\"cooldown\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelay\",\"inputs\":[{\"name\":\"delay\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startEjection\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawEjectorBalance\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// ContractIEigenDAEjectionManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDAEjectionManagerMetaData.ABI instead.
var ContractIEigenDAEjectionManagerABI = ContractIEigenDAEjectionManagerMetaData.ABI

// ContractIEigenDAEjectionManager is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManager struct {
	ContractIEigenDAEjectionManagerCaller     // Read-only binding to the contract
	ContractIEigenDAEjectionManagerTransactor // Write-only binding to the contract
	ContractIEigenDAEjectionManagerFilterer   // Log filterer for contract events
}

// ContractIEigenDAEjectionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAEjectionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAEjectionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDAEjectionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAEjectionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDAEjectionManagerSession struct {
	Contract     *ContractIEigenDAEjectionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractIEigenDAEjectionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDAEjectionManagerCallerSession struct {
	Contract *ContractIEigenDAEjectionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// ContractIEigenDAEjectionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDAEjectionManagerTransactorSession struct {
	Contract     *ContractIEigenDAEjectionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// ContractIEigenDAEjectionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManagerRaw struct {
	Contract *ContractIEigenDAEjectionManager // Generic contract binding to access the raw methods on
}

// ContractIEigenDAEjectionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManagerCallerRaw struct {
	Contract *ContractIEigenDAEjectionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDAEjectionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDAEjectionManagerTransactorRaw struct {
	Contract *ContractIEigenDAEjectionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDAEjectionManager creates a new instance of ContractIEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractIEigenDAEjectionManager(address common.Address, backend bind.ContractBackend) (*ContractIEigenDAEjectionManager, error) {
	contract, err := bindContractIEigenDAEjectionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAEjectionManager{ContractIEigenDAEjectionManagerCaller: ContractIEigenDAEjectionManagerCaller{contract: contract}, ContractIEigenDAEjectionManagerTransactor: ContractIEigenDAEjectionManagerTransactor{contract: contract}, ContractIEigenDAEjectionManagerFilterer: ContractIEigenDAEjectionManagerFilterer{contract: contract}}, nil
}

// NewContractIEigenDAEjectionManagerCaller creates a new read-only instance of ContractIEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractIEigenDAEjectionManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDAEjectionManagerCaller, error) {
	contract, err := bindContractIEigenDAEjectionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAEjectionManagerCaller{contract: contract}, nil
}

// NewContractIEigenDAEjectionManagerTransactor creates a new write-only instance of ContractIEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractIEigenDAEjectionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDAEjectionManagerTransactor, error) {
	contract, err := bindContractIEigenDAEjectionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAEjectionManagerTransactor{contract: contract}, nil
}

// NewContractIEigenDAEjectionManagerFilterer creates a new log filterer instance of ContractIEigenDAEjectionManager, bound to a specific deployed contract.
func NewContractIEigenDAEjectionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDAEjectionManagerFilterer, error) {
	contract, err := bindContractIEigenDAEjectionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAEjectionManagerFilterer{contract: contract}, nil
}

// bindContractIEigenDAEjectionManager binds a generic wrapper to an already deployed contract.
func bindContractIEigenDAEjectionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDAEjectionManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDAEjectionManager.Contract.ContractIEigenDAEjectionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.ContractIEigenDAEjectionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.ContractIEigenDAEjectionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDAEjectionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.contract.Transact(opts, method, params...)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) EjectionCooldown(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "ejectionCooldown")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) EjectionCooldown() (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionCooldown(&_ContractIEigenDAEjectionManager.CallOpts)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) EjectionCooldown() (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionCooldown(&_ContractIEigenDAEjectionManager.CallOpts)
}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) EjectionDelay(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "ejectionDelay")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) EjectionDelay() (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionDelay(&_ContractIEigenDAEjectionManager.CallOpts)
}

// EjectionDelay is a free data retrieval call binding the contract method 0x4f8c9a28.
//
// Solidity: function ejectionDelay() view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) EjectionDelay() (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionDelay(&_ContractIEigenDAEjectionManager.CallOpts)
}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) EjectionQuorums(opts *bind.CallOpts, operator common.Address) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "ejectionQuorums", operator)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) EjectionQuorums(operator common.Address) ([]byte, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionQuorums(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// EjectionQuorums is a free data retrieval call binding the contract method 0xe4049007.
//
// Solidity: function ejectionQuorums(address operator) view returns(bytes)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) EjectionQuorums(operator common.Address) ([]byte, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionQuorums(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) EjectionTime(opts *bind.CallOpts, operator common.Address) (uint64, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "ejectionTime", operator)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) EjectionTime(operator common.Address) (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionTime(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// EjectionTime is a free data retrieval call binding the contract method 0x156570ff.
//
// Solidity: function ejectionTime(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) EjectionTime(operator common.Address) (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.EjectionTime(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) GetDepositToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "getDepositToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) GetDepositToken() (common.Address, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetDepositToken(&_ContractIEigenDAEjectionManager.CallOpts)
}

// GetDepositToken is a free data retrieval call binding the contract method 0xfb1b5c7b.
//
// Solidity: function getDepositToken() view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) GetDepositToken() (common.Address, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetDepositToken(&_ContractIEigenDAEjectionManager.CallOpts)
}

// GetEjector is a free data retrieval call binding the contract method 0xc412ef3b.
//
// Solidity: function getEjector(address operator) view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) GetEjector(opts *bind.CallOpts, operator common.Address) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "getEjector", operator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEjector is a free data retrieval call binding the contract method 0xc412ef3b.
//
// Solidity: function getEjector(address operator) view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) GetEjector(operator common.Address) (common.Address, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetEjector(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// GetEjector is a free data retrieval call binding the contract method 0xc412ef3b.
//
// Solidity: function getEjector(address operator) view returns(address)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) GetEjector(operator common.Address) (common.Address, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetEjector(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// GetEjectorBalance is a free data retrieval call binding the contract method 0x7c292e47.
//
// Solidity: function getEjectorBalance(address ejector) view returns(uint256)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) GetEjectorBalance(opts *bind.CallOpts, ejector common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "getEjectorBalance", ejector)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEjectorBalance is a free data retrieval call binding the contract method 0x7c292e47.
//
// Solidity: function getEjectorBalance(address ejector) view returns(uint256)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) GetEjectorBalance(ejector common.Address) (*big.Int, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetEjectorBalance(&_ContractIEigenDAEjectionManager.CallOpts, ejector)
}

// GetEjectorBalance is a free data retrieval call binding the contract method 0x7c292e47.
//
// Solidity: function getEjectorBalance(address ejector) view returns(uint256)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) GetEjectorBalance(ejector common.Address) (*big.Int, error) {
	return _ContractIEigenDAEjectionManager.Contract.GetEjectorBalance(&_ContractIEigenDAEjectionManager.CallOpts, ejector)
}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCaller) LastEjectionInitiated(opts *bind.CallOpts, operator common.Address) (uint64, error) {
	var out []interface{}
	err := _ContractIEigenDAEjectionManager.contract.Call(opts, &out, "lastEjectionInitiated", operator)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) LastEjectionInitiated(operator common.Address) (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.LastEjectionInitiated(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// LastEjectionInitiated is a free data retrieval call binding the contract method 0xe6f51414.
//
// Solidity: function lastEjectionInitiated(address operator) view returns(uint64)
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerCallerSession) LastEjectionInitiated(operator common.Address) (uint64, error) {
	return _ContractIEigenDAEjectionManager.Contract.LastEjectionInitiated(&_ContractIEigenDAEjectionManager.CallOpts, operator)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) AddEjectorBalance(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "addEjectorBalance", amount)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) AddEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.AddEjectorBalance(&_ContractIEigenDAEjectionManager.TransactOpts, amount)
}

// AddEjectorBalance is a paid mutator transaction binding the contract method 0x3b115362.
//
// Solidity: function addEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) AddEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.AddEjectorBalance(&_ContractIEigenDAEjectionManager.TransactOpts, amount)
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) CancelEjection(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "cancelEjection")
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) CancelEjection() (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjection(&_ContractIEigenDAEjectionManager.TransactOpts)
}

// CancelEjection is a paid mutator transaction binding the contract method 0x39ff1868.
//
// Solidity: function cancelEjection() returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) CancelEjection() (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjection(&_ContractIEigenDAEjectionManager.TransactOpts)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) CancelEjectionByEjector(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "cancelEjectionByEjector", operator)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) CancelEjectionByEjector(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjectionByEjector(&_ContractIEigenDAEjectionManager.TransactOpts, operator)
}

// CancelEjectionByEjector is a paid mutator transaction binding the contract method 0xb0f0ba46.
//
// Solidity: function cancelEjectionByEjector(address operator) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) CancelEjectionByEjector(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjectionByEjector(&_ContractIEigenDAEjectionManager.TransactOpts, operator)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) CancelEjectionWithSig(opts *bind.TransactOpts, operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "cancelEjectionWithSig", operator, apkG2, sigma, recipient)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) CancelEjectionWithSig(operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjectionWithSig(&_ContractIEigenDAEjectionManager.TransactOpts, operator, apkG2, sigma, recipient)
}

// CancelEjectionWithSig is a paid mutator transaction binding the contract method 0x222abf86.
//
// Solidity: function cancelEjectionWithSig(address operator, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma, address recipient) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) CancelEjectionWithSig(operator common.Address, apkG2 BN254G2Point, sigma BN254G1Point, recipient common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CancelEjectionWithSig(&_ContractIEigenDAEjectionManager.TransactOpts, operator, apkG2, sigma, recipient)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) CompleteEjection(opts *bind.TransactOpts, operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "completeEjection", operator, quorums)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) CompleteEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CompleteEjection(&_ContractIEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// CompleteEjection is a paid mutator transaction binding the contract method 0x2d716fbc.
//
// Solidity: function completeEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) CompleteEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.CompleteEjection(&_ContractIEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) SetCooldown(opts *bind.TransactOpts, cooldown uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "setCooldown", cooldown)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) SetCooldown(cooldown uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.SetCooldown(&_ContractIEigenDAEjectionManager.TransactOpts, cooldown)
}

// SetCooldown is a paid mutator transaction binding the contract method 0x4b11982e.
//
// Solidity: function setCooldown(uint64 cooldown) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) SetCooldown(cooldown uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.SetCooldown(&_ContractIEigenDAEjectionManager.TransactOpts, cooldown)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) SetDelay(opts *bind.TransactOpts, delay uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "setDelay", delay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) SetDelay(delay uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.SetDelay(&_ContractIEigenDAEjectionManager.TransactOpts, delay)
}

// SetDelay is a paid mutator transaction binding the contract method 0xc1073302.
//
// Solidity: function setDelay(uint64 delay) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) SetDelay(delay uint64) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.SetDelay(&_ContractIEigenDAEjectionManager.TransactOpts, delay)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) StartEjection(opts *bind.TransactOpts, operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "startEjection", operator, quorums)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) StartEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.StartEjection(&_ContractIEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// StartEjection is a paid mutator transaction binding the contract method 0xb756c6fb.
//
// Solidity: function startEjection(address operator, bytes quorums) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) StartEjection(operator common.Address, quorums []byte) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.StartEjection(&_ContractIEigenDAEjectionManager.TransactOpts, operator, quorums)
}

// WithdrawEjectorBalance is a paid mutator transaction binding the contract method 0xd0ac9cdf.
//
// Solidity: function withdrawEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactor) WithdrawEjectorBalance(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.contract.Transact(opts, "withdrawEjectorBalance", amount)
}

// WithdrawEjectorBalance is a paid mutator transaction binding the contract method 0xd0ac9cdf.
//
// Solidity: function withdrawEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerSession) WithdrawEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.WithdrawEjectorBalance(&_ContractIEigenDAEjectionManager.TransactOpts, amount)
}

// WithdrawEjectorBalance is a paid mutator transaction binding the contract method 0xd0ac9cdf.
//
// Solidity: function withdrawEjectorBalance(uint256 amount) returns()
func (_ContractIEigenDAEjectionManager *ContractIEigenDAEjectionManagerTransactorSession) WithdrawEjectorBalance(amount *big.Int) (*types.Transaction, error) {
	return _ContractIEigenDAEjectionManager.Contract.WithdrawEjectorBalance(&_ContractIEigenDAEjectionManager.TransactOpts, amount)
}
