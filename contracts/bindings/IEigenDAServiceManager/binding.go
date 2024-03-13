// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractIEigenDAServiceManager

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

// IBLSSignatureCheckerNonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerNonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// ContractIEigenDAServiceManagerMetaData contains all meta data concerning the ContractIEigenDAServiceManager contract.
var ContractIEigenDAServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestServeUntilBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBatchConfirmer\",\"inputs\":[{\"name\":\"_batchConfirmer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchConfirmerChanged\",\"inputs\":[{\"name\":\"previousAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// ContractIEigenDAServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractIEigenDAServiceManagerMetaData.ABI instead.
var ContractIEigenDAServiceManagerABI = ContractIEigenDAServiceManagerMetaData.ABI

// ContractIEigenDAServiceManager is an auto generated Go binding around an Ethereum contract.
type ContractIEigenDAServiceManager struct {
	ContractIEigenDAServiceManagerCaller     // Read-only binding to the contract
	ContractIEigenDAServiceManagerTransactor // Write-only binding to the contract
	ContractIEigenDAServiceManagerFilterer   // Log filterer for contract events
}

// ContractIEigenDAServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractIEigenDAServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractIEigenDAServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractIEigenDAServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractIEigenDAServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractIEigenDAServiceManagerSession struct {
	Contract     *ContractIEigenDAServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractIEigenDAServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractIEigenDAServiceManagerCallerSession struct {
	Contract *ContractIEigenDAServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractIEigenDAServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractIEigenDAServiceManagerTransactorSession struct {
	Contract     *ContractIEigenDAServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractIEigenDAServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractIEigenDAServiceManagerRaw struct {
	Contract *ContractIEigenDAServiceManager // Generic contract binding to access the raw methods on
}

// ContractIEigenDAServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractIEigenDAServiceManagerCallerRaw struct {
	Contract *ContractIEigenDAServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractIEigenDAServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractIEigenDAServiceManagerTransactorRaw struct {
	Contract *ContractIEigenDAServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractIEigenDAServiceManager creates a new instance of ContractIEigenDAServiceManager, bound to a specific deployed contract.
func NewContractIEigenDAServiceManager(address common.Address, backend bind.ContractBackend) (*ContractIEigenDAServiceManager, error) {
	contract, err := bindContractIEigenDAServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManager{ContractIEigenDAServiceManagerCaller: ContractIEigenDAServiceManagerCaller{contract: contract}, ContractIEigenDAServiceManagerTransactor: ContractIEigenDAServiceManagerTransactor{contract: contract}, ContractIEigenDAServiceManagerFilterer: ContractIEigenDAServiceManagerFilterer{contract: contract}}, nil
}

// NewContractIEigenDAServiceManagerCaller creates a new read-only instance of ContractIEigenDAServiceManager, bound to a specific deployed contract.
func NewContractIEigenDAServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractIEigenDAServiceManagerCaller, error) {
	contract, err := bindContractIEigenDAServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerCaller{contract: contract}, nil
}

// NewContractIEigenDAServiceManagerTransactor creates a new write-only instance of ContractIEigenDAServiceManager, bound to a specific deployed contract.
func NewContractIEigenDAServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractIEigenDAServiceManagerTransactor, error) {
	contract, err := bindContractIEigenDAServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerTransactor{contract: contract}, nil
}

// NewContractIEigenDAServiceManagerFilterer creates a new log filterer instance of ContractIEigenDAServiceManager, bound to a specific deployed contract.
func NewContractIEigenDAServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractIEigenDAServiceManagerFilterer, error) {
	contract, err := bindContractIEigenDAServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerFilterer{contract: contract}, nil
}

// bindContractIEigenDAServiceManager binds a generic wrapper to an already deployed contract.
func bindContractIEigenDAServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractIEigenDAServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDAServiceManager.Contract.ContractIEigenDAServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ContractIEigenDAServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ContractIEigenDAServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractIEigenDAServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.contract.Transact(opts, method, params...)
}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) BLOCKSTALEMEASURE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "BLOCK_STALE_MEASURE")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) BLOCKSTALEMEASURE() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.BLOCKSTALEMEASURE(&_ContractIEigenDAServiceManager.CallOpts)
}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) BLOCKSTALEMEASURE() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.BLOCKSTALEMEASURE(&_ContractIEigenDAServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.AvsDirectory(&_ContractIEigenDAServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.AvsDirectory(&_ContractIEigenDAServiceManager.CallOpts)
}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 batchId) view returns(bytes32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) BatchIdToBatchMetadataHash(opts *bind.CallOpts, batchId uint32) ([32]byte, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "batchIdToBatchMetadataHash", batchId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 batchId) view returns(bytes32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) BatchIdToBatchMetadataHash(batchId uint32) ([32]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.BatchIdToBatchMetadataHash(&_ContractIEigenDAServiceManager.CallOpts, batchId)
}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 batchId) view returns(bytes32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) BatchIdToBatchMetadataHash(batchId uint32) ([32]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.BatchIdToBatchMetadataHash(&_ContractIEigenDAServiceManager.CallOpts, batchId)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractIEigenDAServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractIEigenDAServiceManager.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.GetRestakeableStrategies(&_ContractIEigenDAServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.GetRestakeableStrategies(&_ContractIEigenDAServiceManager.CallOpts)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0x758f8dba.
//
// Solidity: function latestServeUntilBlock() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) LatestServeUntilBlock(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "latestServeUntilBlock")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0x758f8dba.
//
// Solidity: function latestServeUntilBlock() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) LatestServeUntilBlock() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractIEigenDAServiceManager.CallOpts)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0x758f8dba.
//
// Solidity: function latestServeUntilBlock() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) LatestServeUntilBlock() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractIEigenDAServiceManager.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumAdversaryThresholdPercentages(&_ContractIEigenDAServiceManager.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumAdversaryThresholdPercentages(&_ContractIEigenDAServiceManager.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumConfirmationThresholdPercentages(&_ContractIEigenDAServiceManager.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumConfirmationThresholdPercentages(&_ContractIEigenDAServiceManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) TaskNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "taskNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) TaskNumber() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.TaskNumber(&_ContractIEigenDAServiceManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) TaskNumber() (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.TaskNumber(&_ContractIEigenDAServiceManager.CallOpts)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) ConfirmBatch(opts *bind.TransactOpts, batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "confirmBatch", batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) ConfirmBatch(batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ConfirmBatch(&_ContractIEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) ConfirmBatch(batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ConfirmBatch(&_ContractIEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractIEigenDAServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractIEigenDAServiceManager.TransactOpts, operator)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.RegisterOperatorToAVS(&_ContractIEigenDAServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.RegisterOperatorToAVS(&_ContractIEigenDAServiceManager.TransactOpts, operator, operatorSignature)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) SetBatchConfirmer(opts *bind.TransactOpts, _batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "setBatchConfirmer", _batchConfirmer)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) SetBatchConfirmer(_batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetBatchConfirmer(&_ContractIEigenDAServiceManager.TransactOpts, _batchConfirmer)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) SetBatchConfirmer(_batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetBatchConfirmer(&_ContractIEigenDAServiceManager.TransactOpts, _batchConfirmer)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) SetMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "setMetadataURI", _metadataURI)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) SetMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetMetadataURI(&_ContractIEigenDAServiceManager.TransactOpts, _metadataURI)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x750521f5.
//
// Solidity: function setMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) SetMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetMetadataURI(&_ContractIEigenDAServiceManager.TransactOpts, _metadataURI)
}

// ContractIEigenDAServiceManagerBatchConfirmedIterator is returned from FilterBatchConfirmed and is used to iterate over the raw logs and unpacked data for BatchConfirmed events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmedIterator struct {
	Event *ContractIEigenDAServiceManagerBatchConfirmed // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerBatchConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerBatchConfirmed)
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
		it.Event = new(ContractIEigenDAServiceManagerBatchConfirmed)
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
func (it *ContractIEigenDAServiceManagerBatchConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerBatchConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerBatchConfirmed represents a BatchConfirmed event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmed struct {
	BatchHeaderHash [32]byte
	BatchId         uint32
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmed is a free log retrieval operation binding the contract event 0xc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterBatchConfirmed(opts *bind.FilterOpts, batchHeaderHash [][32]byte) (*ContractIEigenDAServiceManagerBatchConfirmedIterator, error) {

	var batchHeaderHashRule []interface{}
	for _, batchHeaderHashItem := range batchHeaderHash {
		batchHeaderHashRule = append(batchHeaderHashRule, batchHeaderHashItem)
	}

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "BatchConfirmed", batchHeaderHashRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerBatchConfirmedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "BatchConfirmed", logs: logs, sub: sub}, nil
}

// WatchBatchConfirmed is a free log subscription operation binding the contract event 0xc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchBatchConfirmed(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerBatchConfirmed, batchHeaderHash [][32]byte) (event.Subscription, error) {

	var batchHeaderHashRule []interface{}
	for _, batchHeaderHashItem := range batchHeaderHash {
		batchHeaderHashRule = append(batchHeaderHashRule, batchHeaderHashItem)
	}

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "BatchConfirmed", batchHeaderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerBatchConfirmed)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmed", log); err != nil {
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

// ParseBatchConfirmed is a log parse operation binding the contract event 0xc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseBatchConfirmed(log types.Log) (*ContractIEigenDAServiceManagerBatchConfirmed, error) {
	event := new(ContractIEigenDAServiceManagerBatchConfirmed)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerBatchConfirmerChangedIterator is returned from FilterBatchConfirmerChanged and is used to iterate over the raw logs and unpacked data for BatchConfirmerChanged events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmerChangedIterator struct {
	Event *ContractIEigenDAServiceManagerBatchConfirmerChanged // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerBatchConfirmerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerBatchConfirmerChanged)
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
		it.Event = new(ContractIEigenDAServiceManagerBatchConfirmerChanged)
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
func (it *ContractIEigenDAServiceManagerBatchConfirmerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerBatchConfirmerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerBatchConfirmerChanged represents a BatchConfirmerChanged event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmerChanged struct {
	PreviousAddress common.Address
	NewAddress      common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmerChanged is a free log retrieval operation binding the contract event 0xf024af0387e1367ceb1c6a3b6a00db4e9917e56bfb22a289808100f8e2b2c085.
//
// Solidity: event BatchConfirmerChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterBatchConfirmerChanged(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerBatchConfirmerChangedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "BatchConfirmerChanged")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerBatchConfirmerChangedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "BatchConfirmerChanged", logs: logs, sub: sub}, nil
}

// WatchBatchConfirmerChanged is a free log subscription operation binding the contract event 0xf024af0387e1367ceb1c6a3b6a00db4e9917e56bfb22a289808100f8e2b2c085.
//
// Solidity: event BatchConfirmerChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchBatchConfirmerChanged(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerBatchConfirmerChanged) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "BatchConfirmerChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerBatchConfirmerChanged)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerChanged", log); err != nil {
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

// ParseBatchConfirmerChanged is a log parse operation binding the contract event 0xf024af0387e1367ceb1c6a3b6a00db4e9917e56bfb22a289808100f8e2b2c085.
//
// Solidity: event BatchConfirmerChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseBatchConfirmerChanged(log types.Log) (*ContractIEigenDAServiceManagerBatchConfirmerChanged, error) {
	event := new(ContractIEigenDAServiceManagerBatchConfirmerChanged)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
