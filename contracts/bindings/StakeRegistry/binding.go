// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractStakeRegistry

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

// IStakeRegistryStakeUpdate is an auto generated low-level Go binding around an user-defined struct.
type IStakeRegistryStakeUpdate struct {
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
	Stake                 *big.Int
}

// IStakeRegistryStrategyParams is an auto generated low-level Go binding around an user-defined struct.
type IStakeRegistryStrategyParams struct {
	Strategy   common.Address
	Multiplier *big.Int
}

// ContractStakeRegistryMetaData contains all meta data concerning the ContractStakeRegistry contract.
var ContractStakeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_delegationManager\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_WEIGHING_FUNCTION_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WEIGHTING_DIVISOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addStrategies\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"_strategyParams\",\"type\":\"tuple[]\",\"internalType\":\"structIStakeRegistry.StrategyParams[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCurrentStake\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentTotalStake\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestStakeUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIStakeRegistry.StakeUpdate\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeAtBlockNumber\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeAtBlockNumberAndIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeHistory\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIStakeRegistry.StakeUpdate[]\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeHistoryLength\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIStakeRegistry.StakeUpdate\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakeUpdateIndexAtBlockNumber\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalStakeAtBlockNumberFromIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalStakeHistoryLength\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalStakeIndicesAtBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalStakeUpdateAtIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIStakeRegistry.StakeUpdate\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"stake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minimumStake\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"_strategyParams\",\"type\":\"tuple[]\",\"internalType\":\"structIStakeRegistry.StrategyParams[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minimumStakeForQuorum\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"modifyStrategyParams\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"strategyIndices\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"newMultipliers\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeStrategies\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"indicesToRemove\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinimumStakeForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"minimumStake\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"strategiesPerQuorum\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyParams\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyParamsByIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIStakeRegistry.StrategyParams\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyParamsLength\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateOperatorStake\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"weightOfOperatorForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"MinimumStakeForQuorumUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"minimumStake\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorStakeUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"stake\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumCreated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyAddedToQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyMultiplierUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyRemovedFromQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"}],\"anonymous\":false}]",
}

// ContractStakeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractStakeRegistryMetaData.ABI instead.
var ContractStakeRegistryABI = ContractStakeRegistryMetaData.ABI

// ContractStakeRegistry is an auto generated Go binding around an Ethereum contract.
type ContractStakeRegistry struct {
	ContractStakeRegistryCaller     // Read-only binding to the contract
	ContractStakeRegistryTransactor // Write-only binding to the contract
	ContractStakeRegistryFilterer   // Log filterer for contract events
}

// ContractStakeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractStakeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractStakeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractStakeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractStakeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractStakeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractStakeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractStakeRegistrySession struct {
	Contract     *ContractStakeRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ContractStakeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractStakeRegistryCallerSession struct {
	Contract *ContractStakeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ContractStakeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractStakeRegistryTransactorSession struct {
	Contract     *ContractStakeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ContractStakeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractStakeRegistryRaw struct {
	Contract *ContractStakeRegistry // Generic contract binding to access the raw methods on
}

// ContractStakeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractStakeRegistryCallerRaw struct {
	Contract *ContractStakeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractStakeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractStakeRegistryTransactorRaw struct {
	Contract *ContractStakeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractStakeRegistry creates a new instance of ContractStakeRegistry, bound to a specific deployed contract.
func NewContractStakeRegistry(address common.Address, backend bind.ContractBackend) (*ContractStakeRegistry, error) {
	contract, err := bindContractStakeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistry{ContractStakeRegistryCaller: ContractStakeRegistryCaller{contract: contract}, ContractStakeRegistryTransactor: ContractStakeRegistryTransactor{contract: contract}, ContractStakeRegistryFilterer: ContractStakeRegistryFilterer{contract: contract}}, nil
}

// NewContractStakeRegistryCaller creates a new read-only instance of ContractStakeRegistry, bound to a specific deployed contract.
func NewContractStakeRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractStakeRegistryCaller, error) {
	contract, err := bindContractStakeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryCaller{contract: contract}, nil
}

// NewContractStakeRegistryTransactor creates a new write-only instance of ContractStakeRegistry, bound to a specific deployed contract.
func NewContractStakeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractStakeRegistryTransactor, error) {
	contract, err := bindContractStakeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryTransactor{contract: contract}, nil
}

// NewContractStakeRegistryFilterer creates a new log filterer instance of ContractStakeRegistry, bound to a specific deployed contract.
func NewContractStakeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractStakeRegistryFilterer, error) {
	contract, err := bindContractStakeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryFilterer{contract: contract}, nil
}

// bindContractStakeRegistry binds a generic wrapper to an already deployed contract.
func bindContractStakeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractStakeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractStakeRegistry *ContractStakeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractStakeRegistry.Contract.ContractStakeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractStakeRegistry *ContractStakeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.ContractStakeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractStakeRegistry *ContractStakeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.ContractStakeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractStakeRegistry *ContractStakeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractStakeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractStakeRegistry *ContractStakeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractStakeRegistry *ContractStakeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.contract.Transact(opts, method, params...)
}

// MAXWEIGHINGFUNCTIONLENGTH is a free data retrieval call binding the contract method 0x7c172347.
//
// Solidity: function MAX_WEIGHING_FUNCTION_LENGTH() view returns(uint8)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) MAXWEIGHINGFUNCTIONLENGTH(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "MAX_WEIGHING_FUNCTION_LENGTH")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MAXWEIGHINGFUNCTIONLENGTH is a free data retrieval call binding the contract method 0x7c172347.
//
// Solidity: function MAX_WEIGHING_FUNCTION_LENGTH() view returns(uint8)
func (_ContractStakeRegistry *ContractStakeRegistrySession) MAXWEIGHINGFUNCTIONLENGTH() (uint8, error) {
	return _ContractStakeRegistry.Contract.MAXWEIGHINGFUNCTIONLENGTH(&_ContractStakeRegistry.CallOpts)
}

// MAXWEIGHINGFUNCTIONLENGTH is a free data retrieval call binding the contract method 0x7c172347.
//
// Solidity: function MAX_WEIGHING_FUNCTION_LENGTH() view returns(uint8)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) MAXWEIGHINGFUNCTIONLENGTH() (uint8, error) {
	return _ContractStakeRegistry.Contract.MAXWEIGHINGFUNCTIONLENGTH(&_ContractStakeRegistry.CallOpts)
}

// WEIGHTINGDIVISOR is a free data retrieval call binding the contract method 0x5e5a6775.
//
// Solidity: function WEIGHTING_DIVISOR() view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) WEIGHTINGDIVISOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "WEIGHTING_DIVISOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WEIGHTINGDIVISOR is a free data retrieval call binding the contract method 0x5e5a6775.
//
// Solidity: function WEIGHTING_DIVISOR() view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistrySession) WEIGHTINGDIVISOR() (*big.Int, error) {
	return _ContractStakeRegistry.Contract.WEIGHTINGDIVISOR(&_ContractStakeRegistry.CallOpts)
}

// WEIGHTINGDIVISOR is a free data retrieval call binding the contract method 0x5e5a6775.
//
// Solidity: function WEIGHTING_DIVISOR() view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) WEIGHTINGDIVISOR() (*big.Int, error) {
	return _ContractStakeRegistry.Contract.WEIGHTINGDIVISOR(&_ContractStakeRegistry.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) Delegation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "delegation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistrySession) Delegation() (common.Address, error) {
	return _ContractStakeRegistry.Contract.Delegation(&_ContractStakeRegistry.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) Delegation() (common.Address, error) {
	return _ContractStakeRegistry.Contract.Delegation(&_ContractStakeRegistry.CallOpts)
}

// GetCurrentStake is a free data retrieval call binding the contract method 0x5401ed27.
//
// Solidity: function getCurrentStake(bytes32 operatorId, uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetCurrentStake(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getCurrentStake", operatorId, quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentStake is a free data retrieval call binding the contract method 0x5401ed27.
//
// Solidity: function getCurrentStake(bytes32 operatorId, uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetCurrentStake(operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetCurrentStake(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetCurrentStake is a free data retrieval call binding the contract method 0x5401ed27.
//
// Solidity: function getCurrentStake(bytes32 operatorId, uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetCurrentStake(operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetCurrentStake(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetCurrentTotalStake is a free data retrieval call binding the contract method 0xd5eccc05.
//
// Solidity: function getCurrentTotalStake(uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetCurrentTotalStake(opts *bind.CallOpts, quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getCurrentTotalStake", quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentTotalStake is a free data retrieval call binding the contract method 0xd5eccc05.
//
// Solidity: function getCurrentTotalStake(uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetCurrentTotalStake(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetCurrentTotalStake(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// GetCurrentTotalStake is a free data retrieval call binding the contract method 0xd5eccc05.
//
// Solidity: function getCurrentTotalStake(uint8 quorumNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetCurrentTotalStake(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetCurrentTotalStake(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// GetLatestStakeUpdate is a free data retrieval call binding the contract method 0xf851e198.
//
// Solidity: function getLatestStakeUpdate(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetLatestStakeUpdate(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8) (IStakeRegistryStakeUpdate, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getLatestStakeUpdate", operatorId, quorumNumber)

	if err != nil {
		return *new(IStakeRegistryStakeUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakeRegistryStakeUpdate)).(*IStakeRegistryStakeUpdate)

	return out0, err

}

// GetLatestStakeUpdate is a free data retrieval call binding the contract method 0xf851e198.
//
// Solidity: function getLatestStakeUpdate(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetLatestStakeUpdate(operatorId [32]byte, quorumNumber uint8) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetLatestStakeUpdate(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetLatestStakeUpdate is a free data retrieval call binding the contract method 0xf851e198.
//
// Solidity: function getLatestStakeUpdate(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetLatestStakeUpdate(operatorId [32]byte, quorumNumber uint8) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetLatestStakeUpdate(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetStakeAtBlockNumber is a free data retrieval call binding the contract method 0xfa28c627.
//
// Solidity: function getStakeAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeAtBlockNumber(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeAtBlockNumber", operatorId, quorumNumber, blockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeAtBlockNumber is a free data retrieval call binding the contract method 0xfa28c627.
//
// Solidity: function getStakeAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeAtBlockNumber(operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeAtBlockNumber(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber, blockNumber)
}

// GetStakeAtBlockNumber is a free data retrieval call binding the contract method 0xfa28c627.
//
// Solidity: function getStakeAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeAtBlockNumber(operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeAtBlockNumber(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber, blockNumber)
}

// GetStakeAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0xf2be94ae.
//
// Solidity: function getStakeAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, bytes32 operatorId, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeAtBlockNumberAndIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, operatorId [32]byte, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeAtBlockNumberAndIndex", quorumNumber, blockNumber, operatorId, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0xf2be94ae.
//
// Solidity: function getStakeAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, bytes32 operatorId, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeAtBlockNumberAndIndex(quorumNumber uint8, blockNumber uint32, operatorId [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeAtBlockNumberAndIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, blockNumber, operatorId, index)
}

// GetStakeAtBlockNumberAndIndex is a free data retrieval call binding the contract method 0xf2be94ae.
//
// Solidity: function getStakeAtBlockNumberAndIndex(uint8 quorumNumber, uint32 blockNumber, bytes32 operatorId, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeAtBlockNumberAndIndex(quorumNumber uint8, blockNumber uint32, operatorId [32]byte, index *big.Int) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeAtBlockNumberAndIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, blockNumber, operatorId, index)
}

// GetStakeHistory is a free data retrieval call binding the contract method 0x2cd95940.
//
// Solidity: function getStakeHistory(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96)[])
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeHistory(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8) ([]IStakeRegistryStakeUpdate, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeHistory", operatorId, quorumNumber)

	if err != nil {
		return *new([]IStakeRegistryStakeUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new([]IStakeRegistryStakeUpdate)).(*[]IStakeRegistryStakeUpdate)

	return out0, err

}

// GetStakeHistory is a free data retrieval call binding the contract method 0x2cd95940.
//
// Solidity: function getStakeHistory(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96)[])
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeHistory(operatorId [32]byte, quorumNumber uint8) ([]IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetStakeHistory(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetStakeHistory is a free data retrieval call binding the contract method 0x2cd95940.
//
// Solidity: function getStakeHistory(bytes32 operatorId, uint8 quorumNumber) view returns((uint32,uint32,uint96)[])
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeHistory(operatorId [32]byte, quorumNumber uint8) ([]IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetStakeHistory(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetStakeHistoryLength is a free data retrieval call binding the contract method 0x4bd26e09.
//
// Solidity: function getStakeHistoryLength(bytes32 operatorId, uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeHistoryLength(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeHistoryLength", operatorId, quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakeHistoryLength is a free data retrieval call binding the contract method 0x4bd26e09.
//
// Solidity: function getStakeHistoryLength(bytes32 operatorId, uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeHistoryLength(operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeHistoryLength(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetStakeHistoryLength is a free data retrieval call binding the contract method 0x4bd26e09.
//
// Solidity: function getStakeHistoryLength(bytes32 operatorId, uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeHistoryLength(operatorId [32]byte, quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetStakeHistoryLength(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber)
}

// GetStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xac6bfb03.
//
// Solidity: function getStakeUpdateAtIndex(uint8 quorumNumber, bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, operatorId [32]byte, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeUpdateAtIndex", quorumNumber, operatorId, index)

	if err != nil {
		return *new(IStakeRegistryStakeUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakeRegistryStakeUpdate)).(*IStakeRegistryStakeUpdate)

	return out0, err

}

// GetStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xac6bfb03.
//
// Solidity: function getStakeUpdateAtIndex(uint8 quorumNumber, bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeUpdateAtIndex(quorumNumber uint8, operatorId [32]byte, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetStakeUpdateAtIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, operatorId, index)
}

// GetStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xac6bfb03.
//
// Solidity: function getStakeUpdateAtIndex(uint8 quorumNumber, bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeUpdateAtIndex(quorumNumber uint8, operatorId [32]byte, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetStakeUpdateAtIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, operatorId, index)
}

// GetStakeUpdateIndexAtBlockNumber is a free data retrieval call binding the contract method 0xdd9846b9.
//
// Solidity: function getStakeUpdateIndexAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint32)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetStakeUpdateIndexAtBlockNumber(opts *bind.CallOpts, operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (uint32, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getStakeUpdateIndexAtBlockNumber", operatorId, quorumNumber, blockNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetStakeUpdateIndexAtBlockNumber is a free data retrieval call binding the contract method 0xdd9846b9.
//
// Solidity: function getStakeUpdateIndexAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint32)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetStakeUpdateIndexAtBlockNumber(operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (uint32, error) {
	return _ContractStakeRegistry.Contract.GetStakeUpdateIndexAtBlockNumber(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber, blockNumber)
}

// GetStakeUpdateIndexAtBlockNumber is a free data retrieval call binding the contract method 0xdd9846b9.
//
// Solidity: function getStakeUpdateIndexAtBlockNumber(bytes32 operatorId, uint8 quorumNumber, uint32 blockNumber) view returns(uint32)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetStakeUpdateIndexAtBlockNumber(operatorId [32]byte, quorumNumber uint8, blockNumber uint32) (uint32, error) {
	return _ContractStakeRegistry.Contract.GetStakeUpdateIndexAtBlockNumber(&_ContractStakeRegistry.CallOpts, operatorId, quorumNumber, blockNumber)
}

// GetTotalStakeAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc8294c56.
//
// Solidity: function getTotalStakeAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetTotalStakeAtBlockNumberFromIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getTotalStakeAtBlockNumberFromIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalStakeAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc8294c56.
//
// Solidity: function getTotalStakeAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetTotalStakeAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeAtBlockNumberFromIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetTotalStakeAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc8294c56.
//
// Solidity: function getTotalStakeAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetTotalStakeAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeAtBlockNumberFromIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetTotalStakeHistoryLength is a free data retrieval call binding the contract method 0x0491b41c.
//
// Solidity: function getTotalStakeHistoryLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetTotalStakeHistoryLength(opts *bind.CallOpts, quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getTotalStakeHistoryLength", quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalStakeHistoryLength is a free data retrieval call binding the contract method 0x0491b41c.
//
// Solidity: function getTotalStakeHistoryLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetTotalStakeHistoryLength(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeHistoryLength(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// GetTotalStakeHistoryLength is a free data retrieval call binding the contract method 0x0491b41c.
//
// Solidity: function getTotalStakeHistoryLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetTotalStakeHistoryLength(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeHistoryLength(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// GetTotalStakeIndicesAtBlockNumber is a free data retrieval call binding the contract method 0x81c07502.
//
// Solidity: function getTotalStakeIndicesAtBlockNumber(uint32 blockNumber, bytes quorumNumbers) view returns(uint32[])
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetTotalStakeIndicesAtBlockNumber(opts *bind.CallOpts, blockNumber uint32, quorumNumbers []byte) ([]uint32, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getTotalStakeIndicesAtBlockNumber", blockNumber, quorumNumbers)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// GetTotalStakeIndicesAtBlockNumber is a free data retrieval call binding the contract method 0x81c07502.
//
// Solidity: function getTotalStakeIndicesAtBlockNumber(uint32 blockNumber, bytes quorumNumbers) view returns(uint32[])
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetTotalStakeIndicesAtBlockNumber(blockNumber uint32, quorumNumbers []byte) ([]uint32, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeIndicesAtBlockNumber(&_ContractStakeRegistry.CallOpts, blockNumber, quorumNumbers)
}

// GetTotalStakeIndicesAtBlockNumber is a free data retrieval call binding the contract method 0x81c07502.
//
// Solidity: function getTotalStakeIndicesAtBlockNumber(uint32 blockNumber, bytes quorumNumbers) view returns(uint32[])
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetTotalStakeIndicesAtBlockNumber(blockNumber uint32, quorumNumbers []byte) ([]uint32, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeIndicesAtBlockNumber(&_ContractStakeRegistry.CallOpts, blockNumber, quorumNumbers)
}

// GetTotalStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xb6904b78.
//
// Solidity: function getTotalStakeUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCaller) GetTotalStakeUpdateAtIndex(opts *bind.CallOpts, quorumNumber uint8, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "getTotalStakeUpdateAtIndex", quorumNumber, index)

	if err != nil {
		return *new(IStakeRegistryStakeUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakeRegistryStakeUpdate)).(*IStakeRegistryStakeUpdate)

	return out0, err

}

// GetTotalStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xb6904b78.
//
// Solidity: function getTotalStakeUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistrySession) GetTotalStakeUpdateAtIndex(quorumNumber uint8, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeUpdateAtIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, index)
}

// GetTotalStakeUpdateAtIndex is a free data retrieval call binding the contract method 0xb6904b78.
//
// Solidity: function getTotalStakeUpdateAtIndex(uint8 quorumNumber, uint256 index) view returns((uint32,uint32,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) GetTotalStakeUpdateAtIndex(quorumNumber uint8, index *big.Int) (IStakeRegistryStakeUpdate, error) {
	return _ContractStakeRegistry.Contract.GetTotalStakeUpdateAtIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, index)
}

// MinimumStakeForQuorum is a free data retrieval call binding the contract method 0xc46778a5.
//
// Solidity: function minimumStakeForQuorum(uint8 ) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) MinimumStakeForQuorum(opts *bind.CallOpts, arg0 uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "minimumStakeForQuorum", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumStakeForQuorum is a free data retrieval call binding the contract method 0xc46778a5.
//
// Solidity: function minimumStakeForQuorum(uint8 ) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) MinimumStakeForQuorum(arg0 uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.MinimumStakeForQuorum(&_ContractStakeRegistry.CallOpts, arg0)
}

// MinimumStakeForQuorum is a free data retrieval call binding the contract method 0xc46778a5.
//
// Solidity: function minimumStakeForQuorum(uint8 ) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) MinimumStakeForQuorum(arg0 uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.MinimumStakeForQuorum(&_ContractStakeRegistry.CallOpts, arg0)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractStakeRegistry.Contract.RegistryCoordinator(&_ContractStakeRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractStakeRegistry.Contract.RegistryCoordinator(&_ContractStakeRegistry.CallOpts)
}

// StrategiesPerQuorum is a free data retrieval call binding the contract method 0x9f3ccf65.
//
// Solidity: function strategiesPerQuorum(uint8 , uint256 ) view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) StrategiesPerQuorum(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "strategiesPerQuorum", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StrategiesPerQuorum is a free data retrieval call binding the contract method 0x9f3ccf65.
//
// Solidity: function strategiesPerQuorum(uint8 , uint256 ) view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistrySession) StrategiesPerQuorum(arg0 uint8, arg1 *big.Int) (common.Address, error) {
	return _ContractStakeRegistry.Contract.StrategiesPerQuorum(&_ContractStakeRegistry.CallOpts, arg0, arg1)
}

// StrategiesPerQuorum is a free data retrieval call binding the contract method 0x9f3ccf65.
//
// Solidity: function strategiesPerQuorum(uint8 , uint256 ) view returns(address)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) StrategiesPerQuorum(arg0 uint8, arg1 *big.Int) (common.Address, error) {
	return _ContractStakeRegistry.Contract.StrategiesPerQuorum(&_ContractStakeRegistry.CallOpts, arg0, arg1)
}

// StrategyParams is a free data retrieval call binding the contract method 0x08732461.
//
// Solidity: function strategyParams(uint8 , uint256 ) view returns(address strategy, uint96 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) StrategyParams(opts *bind.CallOpts, arg0 uint8, arg1 *big.Int) (struct {
	Strategy   common.Address
	Multiplier *big.Int
}, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "strategyParams", arg0, arg1)

	outstruct := new(struct {
		Strategy   common.Address
		Multiplier *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Strategy = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Multiplier = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StrategyParams is a free data retrieval call binding the contract method 0x08732461.
//
// Solidity: function strategyParams(uint8 , uint256 ) view returns(address strategy, uint96 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistrySession) StrategyParams(arg0 uint8, arg1 *big.Int) (struct {
	Strategy   common.Address
	Multiplier *big.Int
}, error) {
	return _ContractStakeRegistry.Contract.StrategyParams(&_ContractStakeRegistry.CallOpts, arg0, arg1)
}

// StrategyParams is a free data retrieval call binding the contract method 0x08732461.
//
// Solidity: function strategyParams(uint8 , uint256 ) view returns(address strategy, uint96 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) StrategyParams(arg0 uint8, arg1 *big.Int) (struct {
	Strategy   common.Address
	Multiplier *big.Int
}, error) {
	return _ContractStakeRegistry.Contract.StrategyParams(&_ContractStakeRegistry.CallOpts, arg0, arg1)
}

// StrategyParamsByIndex is a free data retrieval call binding the contract method 0xadc804da.
//
// Solidity: function strategyParamsByIndex(uint8 quorumNumber, uint256 index) view returns((address,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCaller) StrategyParamsByIndex(opts *bind.CallOpts, quorumNumber uint8, index *big.Int) (IStakeRegistryStrategyParams, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "strategyParamsByIndex", quorumNumber, index)

	if err != nil {
		return *new(IStakeRegistryStrategyParams), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakeRegistryStrategyParams)).(*IStakeRegistryStrategyParams)

	return out0, err

}

// StrategyParamsByIndex is a free data retrieval call binding the contract method 0xadc804da.
//
// Solidity: function strategyParamsByIndex(uint8 quorumNumber, uint256 index) view returns((address,uint96))
func (_ContractStakeRegistry *ContractStakeRegistrySession) StrategyParamsByIndex(quorumNumber uint8, index *big.Int) (IStakeRegistryStrategyParams, error) {
	return _ContractStakeRegistry.Contract.StrategyParamsByIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, index)
}

// StrategyParamsByIndex is a free data retrieval call binding the contract method 0xadc804da.
//
// Solidity: function strategyParamsByIndex(uint8 quorumNumber, uint256 index) view returns((address,uint96))
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) StrategyParamsByIndex(quorumNumber uint8, index *big.Int) (IStakeRegistryStrategyParams, error) {
	return _ContractStakeRegistry.Contract.StrategyParamsByIndex(&_ContractStakeRegistry.CallOpts, quorumNumber, index)
}

// StrategyParamsLength is a free data retrieval call binding the contract method 0x3ca5a5f5.
//
// Solidity: function strategyParamsLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) StrategyParamsLength(opts *bind.CallOpts, quorumNumber uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "strategyParamsLength", quorumNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StrategyParamsLength is a free data retrieval call binding the contract method 0x3ca5a5f5.
//
// Solidity: function strategyParamsLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistrySession) StrategyParamsLength(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.StrategyParamsLength(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// StrategyParamsLength is a free data retrieval call binding the contract method 0x3ca5a5f5.
//
// Solidity: function strategyParamsLength(uint8 quorumNumber) view returns(uint256)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) StrategyParamsLength(quorumNumber uint8) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.StrategyParamsLength(&_ContractStakeRegistry.CallOpts, quorumNumber)
}

// WeightOfOperatorForQuorum is a free data retrieval call binding the contract method 0x1f9b74e0.
//
// Solidity: function weightOfOperatorForQuorum(uint8 quorumNumber, address operator) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCaller) WeightOfOperatorForQuorum(opts *bind.CallOpts, quorumNumber uint8, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractStakeRegistry.contract.Call(opts, &out, "weightOfOperatorForQuorum", quorumNumber, operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WeightOfOperatorForQuorum is a free data retrieval call binding the contract method 0x1f9b74e0.
//
// Solidity: function weightOfOperatorForQuorum(uint8 quorumNumber, address operator) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistrySession) WeightOfOperatorForQuorum(quorumNumber uint8, operator common.Address) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.WeightOfOperatorForQuorum(&_ContractStakeRegistry.CallOpts, quorumNumber, operator)
}

// WeightOfOperatorForQuorum is a free data retrieval call binding the contract method 0x1f9b74e0.
//
// Solidity: function weightOfOperatorForQuorum(uint8 quorumNumber, address operator) view returns(uint96)
func (_ContractStakeRegistry *ContractStakeRegistryCallerSession) WeightOfOperatorForQuorum(quorumNumber uint8, operator common.Address) (*big.Int, error) {
	return _ContractStakeRegistry.Contract.WeightOfOperatorForQuorum(&_ContractStakeRegistry.CallOpts, quorumNumber, operator)
}

// AddStrategies is a paid mutator transaction binding the contract method 0xc601527d.
//
// Solidity: function addStrategies(uint8 quorumNumber, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) AddStrategies(opts *bind.TransactOpts, quorumNumber uint8, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "addStrategies", quorumNumber, _strategyParams)
}

// AddStrategies is a paid mutator transaction binding the contract method 0xc601527d.
//
// Solidity: function addStrategies(uint8 quorumNumber, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) AddStrategies(quorumNumber uint8, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.AddStrategies(&_ContractStakeRegistry.TransactOpts, quorumNumber, _strategyParams)
}

// AddStrategies is a paid mutator transaction binding the contract method 0xc601527d.
//
// Solidity: function addStrategies(uint8 quorumNumber, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) AddStrategies(quorumNumber uint8, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.AddStrategies(&_ContractStakeRegistry.TransactOpts, quorumNumber, _strategyParams)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "deregisterOperator", operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.DeregisterOperator(&_ContractStakeRegistry.TransactOpts, operatorId, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xbd29b8cd.
//
// Solidity: function deregisterOperator(bytes32 operatorId, bytes quorumNumbers) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) DeregisterOperator(operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.DeregisterOperator(&_ContractStakeRegistry.TransactOpts, operatorId, quorumNumbers)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xff694a77.
//
// Solidity: function initializeQuorum(uint8 quorumNumber, uint96 minimumStake, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) InitializeQuorum(opts *bind.TransactOpts, quorumNumber uint8, minimumStake *big.Int, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "initializeQuorum", quorumNumber, minimumStake, _strategyParams)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xff694a77.
//
// Solidity: function initializeQuorum(uint8 quorumNumber, uint96 minimumStake, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) InitializeQuorum(quorumNumber uint8, minimumStake *big.Int, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.InitializeQuorum(&_ContractStakeRegistry.TransactOpts, quorumNumber, minimumStake, _strategyParams)
}

// InitializeQuorum is a paid mutator transaction binding the contract method 0xff694a77.
//
// Solidity: function initializeQuorum(uint8 quorumNumber, uint96 minimumStake, (address,uint96)[] _strategyParams) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) InitializeQuorum(quorumNumber uint8, minimumStake *big.Int, _strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.InitializeQuorum(&_ContractStakeRegistry.TransactOpts, quorumNumber, minimumStake, _strategyParams)
}

// ModifyStrategyParams is a paid mutator transaction binding the contract method 0x20b66298.
//
// Solidity: function modifyStrategyParams(uint8 quorumNumber, uint256[] strategyIndices, uint96[] newMultipliers) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) ModifyStrategyParams(opts *bind.TransactOpts, quorumNumber uint8, strategyIndices []*big.Int, newMultipliers []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "modifyStrategyParams", quorumNumber, strategyIndices, newMultipliers)
}

// ModifyStrategyParams is a paid mutator transaction binding the contract method 0x20b66298.
//
// Solidity: function modifyStrategyParams(uint8 quorumNumber, uint256[] strategyIndices, uint96[] newMultipliers) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) ModifyStrategyParams(quorumNumber uint8, strategyIndices []*big.Int, newMultipliers []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.ModifyStrategyParams(&_ContractStakeRegistry.TransactOpts, quorumNumber, strategyIndices, newMultipliers)
}

// ModifyStrategyParams is a paid mutator transaction binding the contract method 0x20b66298.
//
// Solidity: function modifyStrategyParams(uint8 quorumNumber, uint256[] strategyIndices, uint96[] newMultipliers) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) ModifyStrategyParams(quorumNumber uint8, strategyIndices []*big.Int, newMultipliers []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.ModifyStrategyParams(&_ContractStakeRegistry.TransactOpts, quorumNumber, strategyIndices, newMultipliers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x25504777.
//
// Solidity: function registerOperator(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint96[], uint96[])
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "registerOperator", operator, operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x25504777.
//
// Solidity: function registerOperator(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint96[], uint96[])
func (_ContractStakeRegistry *ContractStakeRegistrySession) RegisterOperator(operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.RegisterOperator(&_ContractStakeRegistry.TransactOpts, operator, operatorId, quorumNumbers)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x25504777.
//
// Solidity: function registerOperator(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint96[], uint96[])
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) RegisterOperator(operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.RegisterOperator(&_ContractStakeRegistry.TransactOpts, operator, operatorId, quorumNumbers)
}

// RemoveStrategies is a paid mutator transaction binding the contract method 0x5f1f2d77.
//
// Solidity: function removeStrategies(uint8 quorumNumber, uint256[] indicesToRemove) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) RemoveStrategies(opts *bind.TransactOpts, quorumNumber uint8, indicesToRemove []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "removeStrategies", quorumNumber, indicesToRemove)
}

// RemoveStrategies is a paid mutator transaction binding the contract method 0x5f1f2d77.
//
// Solidity: function removeStrategies(uint8 quorumNumber, uint256[] indicesToRemove) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) RemoveStrategies(quorumNumber uint8, indicesToRemove []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.RemoveStrategies(&_ContractStakeRegistry.TransactOpts, quorumNumber, indicesToRemove)
}

// RemoveStrategies is a paid mutator transaction binding the contract method 0x5f1f2d77.
//
// Solidity: function removeStrategies(uint8 quorumNumber, uint256[] indicesToRemove) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) RemoveStrategies(quorumNumber uint8, indicesToRemove []*big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.RemoveStrategies(&_ContractStakeRegistry.TransactOpts, quorumNumber, indicesToRemove)
}

// SetMinimumStakeForQuorum is a paid mutator transaction binding the contract method 0xbc9a40c3.
//
// Solidity: function setMinimumStakeForQuorum(uint8 quorumNumber, uint96 minimumStake) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) SetMinimumStakeForQuorum(opts *bind.TransactOpts, quorumNumber uint8, minimumStake *big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "setMinimumStakeForQuorum", quorumNumber, minimumStake)
}

// SetMinimumStakeForQuorum is a paid mutator transaction binding the contract method 0xbc9a40c3.
//
// Solidity: function setMinimumStakeForQuorum(uint8 quorumNumber, uint96 minimumStake) returns()
func (_ContractStakeRegistry *ContractStakeRegistrySession) SetMinimumStakeForQuorum(quorumNumber uint8, minimumStake *big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.SetMinimumStakeForQuorum(&_ContractStakeRegistry.TransactOpts, quorumNumber, minimumStake)
}

// SetMinimumStakeForQuorum is a paid mutator transaction binding the contract method 0xbc9a40c3.
//
// Solidity: function setMinimumStakeForQuorum(uint8 quorumNumber, uint96 minimumStake) returns()
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) SetMinimumStakeForQuorum(quorumNumber uint8, minimumStake *big.Int) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.SetMinimumStakeForQuorum(&_ContractStakeRegistry.TransactOpts, quorumNumber, minimumStake)
}

// UpdateOperatorStake is a paid mutator transaction binding the contract method 0x66acfefe.
//
// Solidity: function updateOperatorStake(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint192)
func (_ContractStakeRegistry *ContractStakeRegistryTransactor) UpdateOperatorStake(opts *bind.TransactOpts, operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.contract.Transact(opts, "updateOperatorStake", operator, operatorId, quorumNumbers)
}

// UpdateOperatorStake is a paid mutator transaction binding the contract method 0x66acfefe.
//
// Solidity: function updateOperatorStake(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint192)
func (_ContractStakeRegistry *ContractStakeRegistrySession) UpdateOperatorStake(operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.UpdateOperatorStake(&_ContractStakeRegistry.TransactOpts, operator, operatorId, quorumNumbers)
}

// UpdateOperatorStake is a paid mutator transaction binding the contract method 0x66acfefe.
//
// Solidity: function updateOperatorStake(address operator, bytes32 operatorId, bytes quorumNumbers) returns(uint192)
func (_ContractStakeRegistry *ContractStakeRegistryTransactorSession) UpdateOperatorStake(operator common.Address, operatorId [32]byte, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractStakeRegistry.Contract.UpdateOperatorStake(&_ContractStakeRegistry.TransactOpts, operator, operatorId, quorumNumbers)
}

// ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator is returned from FilterMinimumStakeForQuorumUpdated and is used to iterate over the raw logs and unpacked data for MinimumStakeForQuorumUpdated events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator struct {
	Event *ContractStakeRegistryMinimumStakeForQuorumUpdated // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryMinimumStakeForQuorumUpdated)
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
		it.Event = new(ContractStakeRegistryMinimumStakeForQuorumUpdated)
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
func (it *ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryMinimumStakeForQuorumUpdated represents a MinimumStakeForQuorumUpdated event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryMinimumStakeForQuorumUpdated struct {
	QuorumNumber uint8
	MinimumStake *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinimumStakeForQuorumUpdated is a free log retrieval operation binding the contract event 0x26eecff2b70b0a71104ff4d940ba7162d23a95c248771fc487a7be17a596b3cf.
//
// Solidity: event MinimumStakeForQuorumUpdated(uint8 indexed quorumNumber, uint96 minimumStake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterMinimumStakeForQuorumUpdated(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "MinimumStakeForQuorumUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryMinimumStakeForQuorumUpdatedIterator{contract: _ContractStakeRegistry.contract, event: "MinimumStakeForQuorumUpdated", logs: logs, sub: sub}, nil
}

// WatchMinimumStakeForQuorumUpdated is a free log subscription operation binding the contract event 0x26eecff2b70b0a71104ff4d940ba7162d23a95c248771fc487a7be17a596b3cf.
//
// Solidity: event MinimumStakeForQuorumUpdated(uint8 indexed quorumNumber, uint96 minimumStake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchMinimumStakeForQuorumUpdated(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryMinimumStakeForQuorumUpdated, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "MinimumStakeForQuorumUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryMinimumStakeForQuorumUpdated)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "MinimumStakeForQuorumUpdated", log); err != nil {
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

// ParseMinimumStakeForQuorumUpdated is a log parse operation binding the contract event 0x26eecff2b70b0a71104ff4d940ba7162d23a95c248771fc487a7be17a596b3cf.
//
// Solidity: event MinimumStakeForQuorumUpdated(uint8 indexed quorumNumber, uint96 minimumStake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseMinimumStakeForQuorumUpdated(log types.Log) (*ContractStakeRegistryMinimumStakeForQuorumUpdated, error) {
	event := new(ContractStakeRegistryMinimumStakeForQuorumUpdated)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "MinimumStakeForQuorumUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakeRegistryOperatorStakeUpdateIterator is returned from FilterOperatorStakeUpdate and is used to iterate over the raw logs and unpacked data for OperatorStakeUpdate events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryOperatorStakeUpdateIterator struct {
	Event *ContractStakeRegistryOperatorStakeUpdate // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryOperatorStakeUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryOperatorStakeUpdate)
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
		it.Event = new(ContractStakeRegistryOperatorStakeUpdate)
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
func (it *ContractStakeRegistryOperatorStakeUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryOperatorStakeUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryOperatorStakeUpdate represents a OperatorStakeUpdate event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryOperatorStakeUpdate struct {
	OperatorId   [32]byte
	QuorumNumber uint8
	Stake        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOperatorStakeUpdate is a free log retrieval operation binding the contract event 0x2f527d527e95d8fe40aec55377743bb779087da3f6d0d08f12e36444da62327d.
//
// Solidity: event OperatorStakeUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint96 stake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterOperatorStakeUpdate(opts *bind.FilterOpts, operatorId [][32]byte) (*ContractStakeRegistryOperatorStakeUpdateIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "OperatorStakeUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryOperatorStakeUpdateIterator{contract: _ContractStakeRegistry.contract, event: "OperatorStakeUpdate", logs: logs, sub: sub}, nil
}

// WatchOperatorStakeUpdate is a free log subscription operation binding the contract event 0x2f527d527e95d8fe40aec55377743bb779087da3f6d0d08f12e36444da62327d.
//
// Solidity: event OperatorStakeUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint96 stake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchOperatorStakeUpdate(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryOperatorStakeUpdate, operatorId [][32]byte) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "OperatorStakeUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryOperatorStakeUpdate)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "OperatorStakeUpdate", log); err != nil {
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

// ParseOperatorStakeUpdate is a log parse operation binding the contract event 0x2f527d527e95d8fe40aec55377743bb779087da3f6d0d08f12e36444da62327d.
//
// Solidity: event OperatorStakeUpdate(bytes32 indexed operatorId, uint8 quorumNumber, uint96 stake)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseOperatorStakeUpdate(log types.Log) (*ContractStakeRegistryOperatorStakeUpdate, error) {
	event := new(ContractStakeRegistryOperatorStakeUpdate)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "OperatorStakeUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakeRegistryQuorumCreatedIterator is returned from FilterQuorumCreated and is used to iterate over the raw logs and unpacked data for QuorumCreated events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryQuorumCreatedIterator struct {
	Event *ContractStakeRegistryQuorumCreated // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryQuorumCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryQuorumCreated)
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
		it.Event = new(ContractStakeRegistryQuorumCreated)
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
func (it *ContractStakeRegistryQuorumCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryQuorumCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryQuorumCreated represents a QuorumCreated event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryQuorumCreated struct {
	QuorumNumber uint8
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterQuorumCreated is a free log retrieval operation binding the contract event 0x831a9c86c45bb303caf3f064be2bc2b9fd4ecf19e47c4ac02a61e75dabfe55b4.
//
// Solidity: event QuorumCreated(uint8 indexed quorumNumber)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterQuorumCreated(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractStakeRegistryQuorumCreatedIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "QuorumCreated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryQuorumCreatedIterator{contract: _ContractStakeRegistry.contract, event: "QuorumCreated", logs: logs, sub: sub}, nil
}

// WatchQuorumCreated is a free log subscription operation binding the contract event 0x831a9c86c45bb303caf3f064be2bc2b9fd4ecf19e47c4ac02a61e75dabfe55b4.
//
// Solidity: event QuorumCreated(uint8 indexed quorumNumber)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchQuorumCreated(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryQuorumCreated, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "QuorumCreated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryQuorumCreated)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "QuorumCreated", log); err != nil {
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

// ParseQuorumCreated is a log parse operation binding the contract event 0x831a9c86c45bb303caf3f064be2bc2b9fd4ecf19e47c4ac02a61e75dabfe55b4.
//
// Solidity: event QuorumCreated(uint8 indexed quorumNumber)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseQuorumCreated(log types.Log) (*ContractStakeRegistryQuorumCreated, error) {
	event := new(ContractStakeRegistryQuorumCreated)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "QuorumCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakeRegistryStrategyAddedToQuorumIterator is returned from FilterStrategyAddedToQuorum and is used to iterate over the raw logs and unpacked data for StrategyAddedToQuorum events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyAddedToQuorumIterator struct {
	Event *ContractStakeRegistryStrategyAddedToQuorum // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryStrategyAddedToQuorumIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryStrategyAddedToQuorum)
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
		it.Event = new(ContractStakeRegistryStrategyAddedToQuorum)
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
func (it *ContractStakeRegistryStrategyAddedToQuorumIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryStrategyAddedToQuorumIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryStrategyAddedToQuorum represents a StrategyAddedToQuorum event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyAddedToQuorum struct {
	QuorumNumber uint8
	Strategy     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStrategyAddedToQuorum is a free log retrieval operation binding the contract event 0x10565e56cacbf32eca267945f054fec02e59750032d113d3302182ad967f5404.
//
// Solidity: event StrategyAddedToQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterStrategyAddedToQuorum(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractStakeRegistryStrategyAddedToQuorumIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "StrategyAddedToQuorum", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryStrategyAddedToQuorumIterator{contract: _ContractStakeRegistry.contract, event: "StrategyAddedToQuorum", logs: logs, sub: sub}, nil
}

// WatchStrategyAddedToQuorum is a free log subscription operation binding the contract event 0x10565e56cacbf32eca267945f054fec02e59750032d113d3302182ad967f5404.
//
// Solidity: event StrategyAddedToQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchStrategyAddedToQuorum(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryStrategyAddedToQuorum, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "StrategyAddedToQuorum", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryStrategyAddedToQuorum)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyAddedToQuorum", log); err != nil {
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

// ParseStrategyAddedToQuorum is a log parse operation binding the contract event 0x10565e56cacbf32eca267945f054fec02e59750032d113d3302182ad967f5404.
//
// Solidity: event StrategyAddedToQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseStrategyAddedToQuorum(log types.Log) (*ContractStakeRegistryStrategyAddedToQuorum, error) {
	event := new(ContractStakeRegistryStrategyAddedToQuorum)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyAddedToQuorum", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakeRegistryStrategyMultiplierUpdatedIterator is returned from FilterStrategyMultiplierUpdated and is used to iterate over the raw logs and unpacked data for StrategyMultiplierUpdated events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyMultiplierUpdatedIterator struct {
	Event *ContractStakeRegistryStrategyMultiplierUpdated // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryStrategyMultiplierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryStrategyMultiplierUpdated)
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
		it.Event = new(ContractStakeRegistryStrategyMultiplierUpdated)
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
func (it *ContractStakeRegistryStrategyMultiplierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryStrategyMultiplierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryStrategyMultiplierUpdated represents a StrategyMultiplierUpdated event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyMultiplierUpdated struct {
	QuorumNumber uint8
	Strategy     common.Address
	Multiplier   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStrategyMultiplierUpdated is a free log retrieval operation binding the contract event 0x11a5641322da1dff56a4b66eaac31ffa465295ece907cd163437793b4d009a75.
//
// Solidity: event StrategyMultiplierUpdated(uint8 indexed quorumNumber, address strategy, uint256 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterStrategyMultiplierUpdated(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractStakeRegistryStrategyMultiplierUpdatedIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "StrategyMultiplierUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryStrategyMultiplierUpdatedIterator{contract: _ContractStakeRegistry.contract, event: "StrategyMultiplierUpdated", logs: logs, sub: sub}, nil
}

// WatchStrategyMultiplierUpdated is a free log subscription operation binding the contract event 0x11a5641322da1dff56a4b66eaac31ffa465295ece907cd163437793b4d009a75.
//
// Solidity: event StrategyMultiplierUpdated(uint8 indexed quorumNumber, address strategy, uint256 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchStrategyMultiplierUpdated(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryStrategyMultiplierUpdated, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "StrategyMultiplierUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryStrategyMultiplierUpdated)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyMultiplierUpdated", log); err != nil {
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

// ParseStrategyMultiplierUpdated is a log parse operation binding the contract event 0x11a5641322da1dff56a4b66eaac31ffa465295ece907cd163437793b4d009a75.
//
// Solidity: event StrategyMultiplierUpdated(uint8 indexed quorumNumber, address strategy, uint256 multiplier)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseStrategyMultiplierUpdated(log types.Log) (*ContractStakeRegistryStrategyMultiplierUpdated, error) {
	event := new(ContractStakeRegistryStrategyMultiplierUpdated)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyMultiplierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakeRegistryStrategyRemovedFromQuorumIterator is returned from FilterStrategyRemovedFromQuorum and is used to iterate over the raw logs and unpacked data for StrategyRemovedFromQuorum events raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyRemovedFromQuorumIterator struct {
	Event *ContractStakeRegistryStrategyRemovedFromQuorum // Event containing the contract specifics and raw log

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
func (it *ContractStakeRegistryStrategyRemovedFromQuorumIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakeRegistryStrategyRemovedFromQuorum)
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
		it.Event = new(ContractStakeRegistryStrategyRemovedFromQuorum)
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
func (it *ContractStakeRegistryStrategyRemovedFromQuorumIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakeRegistryStrategyRemovedFromQuorumIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakeRegistryStrategyRemovedFromQuorum represents a StrategyRemovedFromQuorum event raised by the ContractStakeRegistry contract.
type ContractStakeRegistryStrategyRemovedFromQuorum struct {
	QuorumNumber uint8
	Strategy     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStrategyRemovedFromQuorum is a free log retrieval operation binding the contract event 0x31fa2e2cd280c9375e13ffcf3d81e2378100186e4058f8d3ddb690b82dcd31f7.
//
// Solidity: event StrategyRemovedFromQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) FilterStrategyRemovedFromQuorum(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractStakeRegistryStrategyRemovedFromQuorumIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.FilterLogs(opts, "StrategyRemovedFromQuorum", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractStakeRegistryStrategyRemovedFromQuorumIterator{contract: _ContractStakeRegistry.contract, event: "StrategyRemovedFromQuorum", logs: logs, sub: sub}, nil
}

// WatchStrategyRemovedFromQuorum is a free log subscription operation binding the contract event 0x31fa2e2cd280c9375e13ffcf3d81e2378100186e4058f8d3ddb690b82dcd31f7.
//
// Solidity: event StrategyRemovedFromQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) WatchStrategyRemovedFromQuorum(opts *bind.WatchOpts, sink chan<- *ContractStakeRegistryStrategyRemovedFromQuorum, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractStakeRegistry.contract.WatchLogs(opts, "StrategyRemovedFromQuorum", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakeRegistryStrategyRemovedFromQuorum)
				if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyRemovedFromQuorum", log); err != nil {
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

// ParseStrategyRemovedFromQuorum is a log parse operation binding the contract event 0x31fa2e2cd280c9375e13ffcf3d81e2378100186e4058f8d3ddb690b82dcd31f7.
//
// Solidity: event StrategyRemovedFromQuorum(uint8 indexed quorumNumber, address strategy)
func (_ContractStakeRegistry *ContractStakeRegistryFilterer) ParseStrategyRemovedFromQuorum(log types.Log) (*ContractStakeRegistryStrategyRemovedFromQuorum, error) {
	event := new(ContractStakeRegistryStrategyRemovedFromQuorum)
	if err := _ContractStakeRegistry.contract.UnpackLog(event, "StrategyRemovedFromQuorum", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
