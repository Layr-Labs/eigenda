// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDACertVerifierV1

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

// EigenDATypesV1BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// EigenDATypesV1BatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BatchMetadata struct {
	BatchHeader             EigenDATypesV1BatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// EigenDATypesV1BlobHeader is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []EigenDATypesV1QuorumBlobParam
}

// EigenDATypesV1BlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1BlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  EigenDATypesV1BatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// EigenDATypesV1QuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1QuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// EigenDATypesV1VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type EigenDATypesV1VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// ContractEigenDACertVerifierV1MetaData contains all meta data concerning the ContractEigenDACertVerifierV1 contract.
var ContractEigenDACertVerifierV1MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_eigenDAThresholdRegistryV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"_eigenDABatchMetadataStorageV1\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDABatchMetadataStorageV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDABatchMetadataStorage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistryV1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.VersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertV1\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyDACertsV1\",\"inputs\":[{\"name\":\"blobHeaders\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobHeader[]\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProofs\",\"type\":\"tuple[]\",\"internalType\":\"structEigenDATypesV1.BlobVerificationProof[]\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structEigenDATypesV1.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"}]",
}

// ContractEigenDACertVerifierV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDACertVerifierV1MetaData.ABI instead.
var ContractEigenDACertVerifierV1ABI = ContractEigenDACertVerifierV1MetaData.ABI

// ContractEigenDACertVerifierV1 is an auto generated Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1 struct {
	ContractEigenDACertVerifierV1Caller     // Read-only binding to the contract
	ContractEigenDACertVerifierV1Transactor // Write-only binding to the contract
	ContractEigenDACertVerifierV1Filterer   // Log filterer for contract events
}

// ContractEigenDACertVerifierV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDACertVerifierV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDACertVerifierV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDACertVerifierV1Session struct {
	Contract     *ContractEigenDACertVerifierV1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDACertVerifierV1CallerSession struct {
	Contract *ContractEigenDACertVerifierV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDACertVerifierV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDACertVerifierV1TransactorSession struct {
	Contract     *ContractEigenDACertVerifierV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDACertVerifierV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1Raw struct {
	Contract *ContractEigenDACertVerifierV1 // Generic contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1CallerRaw struct {
	Contract *ContractEigenDACertVerifierV1Caller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDACertVerifierV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDACertVerifierV1TransactorRaw struct {
	Contract *ContractEigenDACertVerifierV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDACertVerifierV1 creates a new instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1(address common.Address, backend bind.ContractBackend) (*ContractEigenDACertVerifierV1, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1{ContractEigenDACertVerifierV1Caller: ContractEigenDACertVerifierV1Caller{contract: contract}, ContractEigenDACertVerifierV1Transactor: ContractEigenDACertVerifierV1Transactor{contract: contract}, ContractEigenDACertVerifierV1Filterer: ContractEigenDACertVerifierV1Filterer{contract: contract}}, nil
}

// NewContractEigenDACertVerifierV1Caller creates a new read-only instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Caller(address common.Address, caller bind.ContractCaller) (*ContractEigenDACertVerifierV1Caller, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Caller{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1Transactor creates a new write-only instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Transactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDACertVerifierV1Transactor, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Transactor{contract: contract}, nil
}

// NewContractEigenDACertVerifierV1Filterer creates a new log filterer instance of ContractEigenDACertVerifierV1, bound to a specific deployed contract.
func NewContractEigenDACertVerifierV1Filterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDACertVerifierV1Filterer, error) {
	contract, err := bindContractEigenDACertVerifierV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDACertVerifierV1Filterer{contract: contract}, nil
}

// bindContractEigenDACertVerifierV1 binds a generic wrapper to an already deployed contract.
func bindContractEigenDACertVerifierV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDACertVerifierV1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.ContractEigenDACertVerifierV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDACertVerifierV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDACertVerifierV1.Contract.contract.Transact(opts, method, params...)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) EigenDABatchMetadataStorageV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "eigenDABatchMetadataStorageV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDABatchMetadataStorageV1 is a free data retrieval call binding the contract method 0xa9c823e1.
//
// Solidity: function eigenDABatchMetadataStorageV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) EigenDABatchMetadataStorageV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDABatchMetadataStorageV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) EigenDAThresholdRegistryV1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "eigenDAThresholdRegistryV1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// EigenDAThresholdRegistryV1 is a free data retrieval call binding the contract method 0x4cff90c4.
//
// Solidity: function eigenDAThresholdRegistryV1() view returns(address)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) EigenDAThresholdRegistryV1() (common.Address, error) {
	return _ContractEigenDACertVerifierV1.Contract.EigenDAThresholdRegistryV1(&_ContractEigenDACertVerifierV1.CallOpts)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetBlobParams(opts *bind.CallOpts, version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(EigenDATypesV1VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDATypesV1VersionedBlobParams)).(*EigenDATypesV1VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetBlobParams(&_ContractEigenDACertVerifierV1.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetBlobParams(version uint16) (EigenDATypesV1VersionedBlobParams, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetBlobParams(&_ContractEigenDACertVerifierV1.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetIsQuorumRequired(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDACertVerifierV1.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDACertVerifierV1.CallOpts, quorumNumber)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV1.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDACertVerifierV1.Contract.QuorumNumbersRequired(&_ContractEigenDACertVerifierV1.CallOpts)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertV1(opts *bind.CallOpts, blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "verifyDACertV1", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertV1 is a free data retrieval call binding the contract method 0x7d644cad.
//
// Solidity: function verifyDACertV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertV1(blobHeader EigenDATypesV1BlobHeader, blobVerificationProof EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Caller) VerifyDACertsV1(opts *bind.CallOpts, blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDACertVerifierV1.contract.Call(opts, &out, "verifyDACertsV1", blobHeaders, blobVerificationProofs)

	if err != nil {
		return err
	}

	return err

}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1Session) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}

// VerifyDACertsV1 is a free data retrieval call binding the contract method 0x31a3479a.
//
// Solidity: function verifyDACertsV1(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[])[] blobHeaders, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)[] blobVerificationProofs) view returns()
func (_ContractEigenDACertVerifierV1 *ContractEigenDACertVerifierV1CallerSession) VerifyDACertsV1(blobHeaders []EigenDATypesV1BlobHeader, blobVerificationProofs []EigenDATypesV1BlobVerificationProof) error {
	return _ContractEigenDACertVerifierV1.Contract.VerifyDACertsV1(&_ContractEigenDACertVerifierV1.CallOpts, blobHeaders, blobVerificationProofs)
}
