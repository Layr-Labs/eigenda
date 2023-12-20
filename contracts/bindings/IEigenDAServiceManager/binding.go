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
	BlobHeadersRoot            [32]byte
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	ReferenceBlockNumber       uint32
}

// ContractIEigenDAServiceManagerMetaData contains all meta data concerning the ContractIEigenDAServiceManager contract.
var ContractIEigenDAServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slasher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISlasher\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"fee\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePerBytePerTimeSet\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSetterChanged\",\"inputs\":[{\"name\":\"previousAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentManagerSet\",\"inputs\":[{\"name\":\"previousAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false}]",
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) Owner() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.Owner(&_ContractIEigenDAServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) Owner() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.Owner(&_ContractIEigenDAServiceManager.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "slasher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) Slasher() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.Slasher(&_ContractIEigenDAServiceManager.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) Slasher() (common.Address, error) {
	return _ContractIEigenDAServiceManager.Contract.Slasher(&_ContractIEigenDAServiceManager.CallOpts)
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

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) FreezeOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "freezeOperator", operator)
}

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) FreezeOperator(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.FreezeOperator(&_ContractIEigenDAServiceManager.TransactOpts, operator)
}

// FreezeOperator is a paid mutator transaction binding the contract method 0x38c8ee64.
//
// Solidity: function freezeOperator(address operator) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) FreezeOperator(operator common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.FreezeOperator(&_ContractIEigenDAServiceManager.TransactOpts, operator)
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
	Fee             *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmed is a free log retrieval operation binding the contract event 0x2eaa707a79ac1f835863f5a6fdb5f27c0e295dc23adf970a445cd87d126c4d63.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId, uint96 fee)
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

// WatchBatchConfirmed is a free log subscription operation binding the contract event 0x2eaa707a79ac1f835863f5a6fdb5f27c0e295dc23adf970a445cd87d126c4d63.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId, uint96 fee)
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

// ParseBatchConfirmed is a log parse operation binding the contract event 0x2eaa707a79ac1f835863f5a6fdb5f27c0e295dc23adf970a445cd87d126c4d63.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId, uint96 fee)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseBatchConfirmed(log types.Log) (*ContractIEigenDAServiceManagerBatchConfirmed, error) {
	event := new(ContractIEigenDAServiceManagerBatchConfirmed)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator is returned from FilterFeePerBytePerTimeSet and is used to iterate over the raw logs and unpacked data for FeePerBytePerTimeSet events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator struct {
	Event *ContractIEigenDAServiceManagerFeePerBytePerTimeSet // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerFeePerBytePerTimeSet)
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
		it.Event = new(ContractIEigenDAServiceManagerFeePerBytePerTimeSet)
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
func (it *ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerFeePerBytePerTimeSet represents a FeePerBytePerTimeSet event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerFeePerBytePerTimeSet struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeePerBytePerTimeSet is a free log retrieval operation binding the contract event 0xcd1b2c2a220284accd1f9effd811cdecb6beaa4638618b48bbea07ce7ae16996.
//
// Solidity: event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterFeePerBytePerTimeSet(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "FeePerBytePerTimeSet")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerFeePerBytePerTimeSetIterator{contract: _ContractIEigenDAServiceManager.contract, event: "FeePerBytePerTimeSet", logs: logs, sub: sub}, nil
}

// WatchFeePerBytePerTimeSet is a free log subscription operation binding the contract event 0xcd1b2c2a220284accd1f9effd811cdecb6beaa4638618b48bbea07ce7ae16996.
//
// Solidity: event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchFeePerBytePerTimeSet(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerFeePerBytePerTimeSet) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "FeePerBytePerTimeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerFeePerBytePerTimeSet)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "FeePerBytePerTimeSet", log); err != nil {
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

// ParseFeePerBytePerTimeSet is a log parse operation binding the contract event 0xcd1b2c2a220284accd1f9effd811cdecb6beaa4638618b48bbea07ce7ae16996.
//
// Solidity: event FeePerBytePerTimeSet(uint256 previousValue, uint256 newValue)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseFeePerBytePerTimeSet(log types.Log) (*ContractIEigenDAServiceManagerFeePerBytePerTimeSet, error) {
	event := new(ContractIEigenDAServiceManagerFeePerBytePerTimeSet)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "FeePerBytePerTimeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerFeeSetterChangedIterator is returned from FilterFeeSetterChanged and is used to iterate over the raw logs and unpacked data for FeeSetterChanged events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerFeeSetterChangedIterator struct {
	Event *ContractIEigenDAServiceManagerFeeSetterChanged // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerFeeSetterChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerFeeSetterChanged)
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
		it.Event = new(ContractIEigenDAServiceManagerFeeSetterChanged)
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
func (it *ContractIEigenDAServiceManagerFeeSetterChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerFeeSetterChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerFeeSetterChanged represents a FeeSetterChanged event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerFeeSetterChanged struct {
	PreviousAddress common.Address
	NewAddress      common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterFeeSetterChanged is a free log retrieval operation binding the contract event 0x774b126b94b3cc801460a024dd575406c3ebf27affd7c36198a53ac6655f056d.
//
// Solidity: event FeeSetterChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterFeeSetterChanged(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerFeeSetterChangedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "FeeSetterChanged")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerFeeSetterChangedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "FeeSetterChanged", logs: logs, sub: sub}, nil
}

// WatchFeeSetterChanged is a free log subscription operation binding the contract event 0x774b126b94b3cc801460a024dd575406c3ebf27affd7c36198a53ac6655f056d.
//
// Solidity: event FeeSetterChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchFeeSetterChanged(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerFeeSetterChanged) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "FeeSetterChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerFeeSetterChanged)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "FeeSetterChanged", log); err != nil {
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

// ParseFeeSetterChanged is a log parse operation binding the contract event 0x774b126b94b3cc801460a024dd575406c3ebf27affd7c36198a53ac6655f056d.
//
// Solidity: event FeeSetterChanged(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseFeeSetterChanged(log types.Log) (*ContractIEigenDAServiceManagerFeeSetterChanged, error) {
	event := new(ContractIEigenDAServiceManagerFeeSetterChanged)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "FeeSetterChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerPaymentManagerSetIterator is returned from FilterPaymentManagerSet and is used to iterate over the raw logs and unpacked data for PaymentManagerSet events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerPaymentManagerSetIterator struct {
	Event *ContractIEigenDAServiceManagerPaymentManagerSet // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerPaymentManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerPaymentManagerSet)
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
		it.Event = new(ContractIEigenDAServiceManagerPaymentManagerSet)
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
func (it *ContractIEigenDAServiceManagerPaymentManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerPaymentManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerPaymentManagerSet represents a PaymentManagerSet event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerPaymentManagerSet struct {
	PreviousAddress common.Address
	NewAddress      common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaymentManagerSet is a free log retrieval operation binding the contract event 0xa3044efb81dffce20bbf49cae117f167852a973364ae504dfade51a8d022c95a.
//
// Solidity: event PaymentManagerSet(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterPaymentManagerSet(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerPaymentManagerSetIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "PaymentManagerSet")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerPaymentManagerSetIterator{contract: _ContractIEigenDAServiceManager.contract, event: "PaymentManagerSet", logs: logs, sub: sub}, nil
}

// WatchPaymentManagerSet is a free log subscription operation binding the contract event 0xa3044efb81dffce20bbf49cae117f167852a973364ae504dfade51a8d022c95a.
//
// Solidity: event PaymentManagerSet(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchPaymentManagerSet(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerPaymentManagerSet) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "PaymentManagerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerPaymentManagerSet)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "PaymentManagerSet", log); err != nil {
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

// ParsePaymentManagerSet is a log parse operation binding the contract event 0xa3044efb81dffce20bbf49cae117f167852a973364ae504dfade51a8d022c95a.
//
// Solidity: event PaymentManagerSet(address previousAddress, address newAddress)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParsePaymentManagerSet(log types.Log) (*ContractIEigenDAServiceManagerPaymentManagerSet, error) {
	event := new(ContractIEigenDAServiceManagerPaymentManagerSet)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "PaymentManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
