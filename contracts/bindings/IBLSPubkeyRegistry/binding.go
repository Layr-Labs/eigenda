// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIBLSPubkeyRegistry

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

// IBLSPubkeyRegistryApkUpdate is an auto generated low-level Go binding around an user-defined struct.
type IBLSPubkeyRegistryApkUpdate struct {
	ApkHash               [24]byte
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
}

// ContractIBLSPubkeyRegistryMetaData contains all meta data concerning the ContractIBLSPubkeyRegistry contract.
var ContractIBLSPubkeyRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApkForQuorum\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkHashForQuorumAtBlockNumberFromIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkIndicesForQuorumsAtBlockNumber\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApkUpdateForQuorumByIndex\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSPubkeyRegistry.ApkUpdate\",\"components\":[{\"name\":\"apkHash\",\"type\":\"bytes24\",\"internalType\":\"bytes24\"},{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorFromPubkeyHash\",\"inputs\":[{\"name\":\"pubkeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pubkey\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OperatorAddedToQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRemovedFromQuorums\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
}

// ContractIBLSPubkeyRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIBLSPubkeyRegistryMetaData.ABI instead.
var ContractIBLSPubkeyRegistryABI = ContractIBLSPubkeyRegistryMetaData.ABI

// ContractIBLSPubkeyRegistry is an auto generated Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistry struct {
	ContractIBLSPubkeyRegistryCaller     // Read-only binding to the contract
	ContractIBLSPubkeyRegistryTransactor // Write-only binding to the contract
	ContractIBLSPubkeyRegistryFilterer   // Log filterer for contract events
}

// ContractIBLSPubkeyRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIBLSPubkeyRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIBLSPubkeyRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIBLSPubkeyRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIBLSPubkeyRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIBLSPubkeyRegistrySession struct {
	Contract     *ContractIBLSPubkeyRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ContractIBLSPubkeyRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIBLSPubkeyRegistryCallerSession struct {
	Contract *ContractIBLSPubkeyRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ContractIBLSPubkeyRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIBLSPubkeyRegistryTransactorSession struct {
	Contract     *ContractIBLSPubkeyRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// ContractIBLSPubkeyRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistryRaw struct {
	Contract *ContractIBLSPubkeyRegistry // Generic contract binding to access the raw methods on
}

// ContractIBLSPubkeyRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistryCallerRaw struct {
	Contract *ContractIBLSPubkeyRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIBLSPubkeyRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIBLSPubkeyRegistryTransactorRaw struct {
	Contract *ContractIBLSPubkeyRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIBLSPubkeyRegistry creates a new instance of ContractIBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractIBLSPubkeyRegistry(address common.Address, backend bind.ContractBackend) (*ContractIBLSPubkeyRegistry, error) {
	contract, err := bindContractIBLSPubkeyRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistry{ContractIBLSPubkeyRegistryCaller: ContractIBLSPubkeyRegistryCaller{contract: contract}, ContractIBLSPubkeyRegistryTransactor: ContractIBLSPubkeyRegistryTransactor{contract: contract}, ContractIBLSPubkeyRegistryFilterer: ContractIBLSPubkeyRegistryFilterer{contract: contract}}, nil
}

// NewContractIBLSPubkeyRegistryCaller creates a new read-only instance of ContractIBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractIBLSPubkeyRegistryCaller(address common.Address, caller bind.ContractCaller) (*ContractIBLSPubkeyRegistryCaller, error) {
	contract, err := bindContractIBLSPubkeyRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistryCaller{contract: contract}, nil
}

// NewContractIBLSPubkeyRegistryTransactor creates a new write-only instance of ContractIBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractIBLSPubkeyRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIBLSPubkeyRegistryTransactor, error) {
	contract, err := bindContractIBLSPubkeyRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistryTransactor{contract: contract}, nil
}

// NewContractIBLSPubkeyRegistryFilterer creates a new log filterer instance of ContractIBLSPubkeyRegistry, bound to a specific deployed contract.
func NewContractIBLSPubkeyRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIBLSPubkeyRegistryFilterer, error) {
	contract, err := bindContractIBLSPubkeyRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistryFilterer{contract: contract}, nil
}

// bindContractIBLSPubkeyRegistry binds a generic wrapper to an already deployed contract.
func bindContractIBLSPubkeyRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIBLSPubkeyRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIBLSPubkeyRegistry.Contract.ContractIBLSPubkeyRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.ContractIBLSPubkeyRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.ContractIBLSPubkeyRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIBLSPubkeyRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) GetApkForQuorum(opts *bind.CallOpts, quorumNumber uint8) (BN254G1Point, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "getApkForQuorum", quorumNumber)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) GetApkForQuorum(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkForQuorum(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// GetApkForQuorum is a free data retrieval call binding the contract method 0x63a94510.
//
// Solidity: function getApkForQuorum(uint8 quorumNumber) view returns((uint256,uint256))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) GetApkForQuorum(quorumNumber uint8) (BN254G1Point, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkForQuorum(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber)
}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) GetApkHashForQuorumAtBlockNumberFromIndex(opts *bind.CallOpts, quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "getApkHashForQuorumAtBlockNumberFromIndex", quorumNumber, blockNumber, index)

	if err != nil {
		return *new([24]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([24]byte)).(*[24]byte)

	return out0, err

}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) GetApkHashForQuorumAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkHashForQuorumAtBlockNumberFromIndex(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkHashForQuorumAtBlockNumberFromIndex is a free data retrieval call binding the contract method 0xc1af6b24.
//
// Solidity: function getApkHashForQuorumAtBlockNumberFromIndex(uint8 quorumNumber, uint32 blockNumber, uint256 index) view returns(bytes24)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) GetApkHashForQuorumAtBlockNumberFromIndex(quorumNumber uint8, blockNumber uint32, index *big.Int) ([24]byte, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkHashForQuorumAtBlockNumberFromIndex(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber, blockNumber, index)
}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) GetApkIndicesForQuorumsAtBlockNumber(opts *bind.CallOpts, quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "getApkIndicesForQuorumsAtBlockNumber", quorumNumbers, blockNumber)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) GetApkIndicesForQuorumsAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkIndicesForQuorumsAtBlockNumber(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkIndicesForQuorumsAtBlockNumber is a free data retrieval call binding the contract method 0xeda10763.
//
// Solidity: function getApkIndicesForQuorumsAtBlockNumber(bytes quorumNumbers, uint256 blockNumber) view returns(uint32[])
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) GetApkIndicesForQuorumsAtBlockNumber(quorumNumbers []byte, blockNumber *big.Int) ([]uint32, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkIndicesForQuorumsAtBlockNumber(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumbers, blockNumber)
}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) GetApkUpdateForQuorumByIndex(opts *bind.CallOpts, quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "getApkUpdateForQuorumByIndex", quorumNumber, index)

	if err != nil {
		return *new(IBLSPubkeyRegistryApkUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSPubkeyRegistryApkUpdate)).(*IBLSPubkeyRegistryApkUpdate)

	return out0, err

}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) GetApkUpdateForQuorumByIndex(quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkUpdateForQuorumByIndex(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber, index)
}

// GetApkUpdateForQuorumByIndex is a free data retrieval call binding the contract method 0x7225057e.
//
// Solidity: function getApkUpdateForQuorumByIndex(uint8 quorumNumber, uint256 index) view returns((bytes24,uint32,uint32))
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) GetApkUpdateForQuorumByIndex(quorumNumber uint8, index *big.Int) (IBLSPubkeyRegistryApkUpdate, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetApkUpdateForQuorumByIndex(&_ContractIBLSPubkeyRegistry.CallOpts, quorumNumber, index)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) GetOperatorFromPubkeyHash(opts *bind.CallOpts, pubkeyHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "getOperatorFromPubkeyHash", pubkeyHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractIBLSPubkeyRegistry.CallOpts, pubkeyHash)
}

// GetOperatorFromPubkeyHash is a free data retrieval call binding the contract method 0x47b314e8.
//
// Solidity: function getOperatorFromPubkeyHash(bytes32 pubkeyHash) view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) GetOperatorFromPubkeyHash(pubkeyHash [32]byte) (common.Address, error) {
	return _ContractIBLSPubkeyRegistry.Contract.GetOperatorFromPubkeyHash(&_ContractIBLSPubkeyRegistry.CallOpts, pubkeyHash)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIBLSPubkeyRegistry.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) RegistryCoordinator() (common.Address, error) {
	return _ContractIBLSPubkeyRegistry.Contract.RegistryCoordinator(&_ContractIBLSPubkeyRegistry.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractIBLSPubkeyRegistry.Contract.RegistryCoordinator(&_ContractIBLSPubkeyRegistry.CallOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactor) DeregisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.contract.Transact(opts, "deregisterOperator", operator, quorumNumbers, pubkey)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) DeregisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.DeregisterOperator(&_ContractIBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0x24369b2a.
//
// Solidity: function deregisterOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns()
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactorSession) DeregisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.DeregisterOperator(&_ContractIBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactor) RegisterOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.contract.Transact(opts, "registerOperator", operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistrySession) RegisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.RegisterOperator(&_ContractIBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x03ce4bad.
//
// Solidity: function registerOperator(address operator, bytes quorumNumbers, (uint256,uint256) pubkey) returns(bytes32)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryTransactorSession) RegisterOperator(operator common.Address, quorumNumbers []byte, pubkey BN254G1Point) (*types.Transaction, error) {
	return _ContractIBLSPubkeyRegistry.Contract.RegisterOperator(&_ContractIBLSPubkeyRegistry.TransactOpts, operator, quorumNumbers, pubkey)
}

// ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator is returned from FilterOperatorAddedToQuorums and is used to iterate over the raw logs and unpacked data for OperatorAddedToQuorums events raised by the ContractIBLSPubkeyRegistry contract.
type ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator struct {
	Event *ContractIBLSPubkeyRegistryOperatorAddedToQuorums // Event containing the contract specifics and raw log

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
func (it *ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIBLSPubkeyRegistryOperatorAddedToQuorums)
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
		it.Event = new(ContractIBLSPubkeyRegistryOperatorAddedToQuorums)
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
func (it *ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIBLSPubkeyRegistryOperatorAddedToQuorums represents a OperatorAddedToQuorums event raised by the ContractIBLSPubkeyRegistry contract.
type ContractIBLSPubkeyRegistryOperatorAddedToQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorAddedToQuorums is a free log retrieval operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) FilterOperatorAddedToQuorums(opts *bind.FilterOpts) (*ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator, error) {

	logs, sub, err := _ContractIBLSPubkeyRegistry.contract.FilterLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistryOperatorAddedToQuorumsIterator{contract: _ContractIBLSPubkeyRegistry.contract, event: "OperatorAddedToQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorAddedToQuorums is a free log subscription operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) WatchOperatorAddedToQuorums(opts *bind.WatchOpts, sink chan<- *ContractIBLSPubkeyRegistryOperatorAddedToQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractIBLSPubkeyRegistry.contract.WatchLogs(opts, "OperatorAddedToQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIBLSPubkeyRegistryOperatorAddedToQuorums)
				if err := _ContractIBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
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

// ParseOperatorAddedToQuorums is a log parse operation binding the contract event 0x5358c5b42179178c8fc757734ac2a3198f9071c765ee0d8389211525f5005246.
//
// Solidity: event OperatorAddedToQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) ParseOperatorAddedToQuorums(log types.Log) (*ContractIBLSPubkeyRegistryOperatorAddedToQuorums, error) {
	event := new(ContractIBLSPubkeyRegistryOperatorAddedToQuorums)
	if err := _ContractIBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorAddedToQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator is returned from FilterOperatorRemovedFromQuorums and is used to iterate over the raw logs and unpacked data for OperatorRemovedFromQuorums events raised by the ContractIBLSPubkeyRegistry contract.
type ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator struct {
	Event *ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums // Event containing the contract specifics and raw log

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
func (it *ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums)
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
		it.Event = new(ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums)
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
func (it *ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums represents a OperatorRemovedFromQuorums event raised by the ContractIBLSPubkeyRegistry contract.
type ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums struct {
	Operator      common.Address
	QuorumNumbers []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemovedFromQuorums is a free log retrieval operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) FilterOperatorRemovedFromQuorums(opts *bind.FilterOpts) (*ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator, error) {

	logs, sub, err := _ContractIBLSPubkeyRegistry.contract.FilterLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return &ContractIBLSPubkeyRegistryOperatorRemovedFromQuorumsIterator{contract: _ContractIBLSPubkeyRegistry.contract, event: "OperatorRemovedFromQuorums", logs: logs, sub: sub}, nil
}

// WatchOperatorRemovedFromQuorums is a free log subscription operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) WatchOperatorRemovedFromQuorums(opts *bind.WatchOpts, sink chan<- *ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums) (event.Subscription, error) {

	logs, sub, err := _ContractIBLSPubkeyRegistry.contract.WatchLogs(opts, "OperatorRemovedFromQuorums")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums)
				if err := _ContractIBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
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

// ParseOperatorRemovedFromQuorums is a log parse operation binding the contract event 0x14a5172b312e9d2c22b8468f9c70ec2caa9de934fe380734fbc6f3beff2b14ba.
//
// Solidity: event OperatorRemovedFromQuorums(address operator, bytes quorumNumbers)
func (_ContractIBLSPubkeyRegistry *ContractIBLSPubkeyRegistryFilterer) ParseOperatorRemovedFromQuorums(log types.Log) (*ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums, error) {
	event := new(ContractIBLSPubkeyRegistryOperatorRemovedFromQuorums)
	if err := _ContractIBLSPubkeyRegistry.contract.UnpackLog(event, "OperatorRemovedFromQuorums", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
