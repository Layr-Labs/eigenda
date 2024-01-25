// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractmockrollup

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

// EigenDABlobUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDABlobUtilsBlobVerificationProof struct {
	BatchId                uint32
	BlobIndex              uint8
	BatchMetadata          IEigenDAServiceManagerBatchMetadata
	InclusionProof         []byte
	QuorumThresholdIndexes []byte
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot            [32]byte
	QuorumNumbers              []byte
	QuorumThresholdPercentages []byte
	ReferenceBlockNumber       uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
	Fee                     *big.Int
	ConfirmationBlockNumber uint32
}

// IEigenDAServiceManagerBlobHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []IEigenDAServiceManagerQuorumBlobParam
}

// IEigenDAServiceManagerQuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerQuorumBlobParam struct {
	QuorumNumber                 uint8
	AdversaryThresholdPercentage uint8
	QuorumThresholdPercentage    uint8
	QuantizationParameter        uint8
}

// ContractmockrollupMetaData contains all meta data concerning the Contractmockrollup contract.
var ContractmockrollupMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIEigenDAServiceManager\",\"name\":\"_eigenDAServiceManager\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"_tau\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_illegalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakeRequired\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"blacklist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"point\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"proof\",\"type\":\"tuple\"}],\"name\":\"challengeCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"dataLength\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"polynomialCommitment\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eigenDAServiceManager\",\"outputs\":[{\"internalType\":\"contractIEigenDAServiceManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"illegalValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"dataLength\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"quorumNumber\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"quorumThresholdPercentage\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"quantizationParameter\",\"type\":\"uint8\"}],\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\"}],\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"name\":\"blobHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"batchId\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"blobIndex\",\"type\":\"uint8\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"quorumNumbers\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"quorumThresholdPercentages\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"referenceBlockNumber\",\"type\":\"uint32\"}],\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"name\":\"batchHeader\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint96\",\"name\":\"fee\",\"type\":\"uint96\"},{\"internalType\":\"uint32\",\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\"}],\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"name\":\"batchMetadata\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"inclusionProof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"quorumThresholdIndexes\",\"type\":\"bytes\"}],\"internalType\":\"structEigenDABlobUtils.BlobVerificationProof\",\"name\":\"blobVerificationProof\",\"type\":\"tuple\"}],\"name\":\"postCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeRequired\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tau\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ContractmockrollupABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractmockrollupMetaData.ABI instead.
var ContractmockrollupABI = ContractmockrollupMetaData.ABI

// Contractmockrollup is an auto generated Go binding around an Ethereum contract.
type Contractmockrollup struct {
	ContractmockrollupCaller     // Read-only binding to the contract
	ContractmockrollupTransactor // Write-only binding to the contract
	ContractmockrollupFilterer   // Log filterer for contract events
}

// ContractmockrollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractmockrollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractmockrollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractmockrollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractmockrollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractmockrollupSession struct {
	Contract     *Contractmockrollup // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractmockrollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractmockrollupCallerSession struct {
	Contract *ContractmockrollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContractmockrollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractmockrollupTransactorSession struct {
	Contract     *ContractmockrollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContractmockrollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractmockrollupRaw struct {
	Contract *Contractmockrollup // Generic contract binding to access the raw methods on
}

// ContractmockrollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractmockrollupCallerRaw struct {
	Contract *ContractmockrollupCaller // Generic read-only contract binding to access the raw methods on
}

// ContractmockrollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractmockrollupTransactorRaw struct {
	Contract *ContractmockrollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractmockrollup creates a new instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollup(address common.Address, backend bind.ContractBackend) (*Contractmockrollup, error) {
	contract, err := bindContractmockrollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contractmockrollup{ContractmockrollupCaller: ContractmockrollupCaller{contract: contract}, ContractmockrollupTransactor: ContractmockrollupTransactor{contract: contract}, ContractmockrollupFilterer: ContractmockrollupFilterer{contract: contract}}, nil
}

// NewContractmockrollupCaller creates a new read-only instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupCaller(address common.Address, caller bind.ContractCaller) (*ContractmockrollupCaller, error) {
	contract, err := bindContractmockrollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupCaller{contract: contract}, nil
}

// NewContractmockrollupTransactor creates a new write-only instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractmockrollupTransactor, error) {
	contract, err := bindContractmockrollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupTransactor{contract: contract}, nil
}

// NewContractmockrollupFilterer creates a new log filterer instance of Contractmockrollup, bound to a specific deployed contract.
func NewContractmockrollupFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractmockrollupFilterer, error) {
	contract, err := bindContractmockrollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractmockrollupFilterer{contract: contract}, nil
}

// bindContractmockrollup binds a generic wrapper to an already deployed contract.
func bindContractmockrollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractmockrollupMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contractmockrollup *ContractmockrollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contractmockrollup.Contract.ContractmockrollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contractmockrollup *ContractmockrollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ContractmockrollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contractmockrollup *ContractmockrollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ContractmockrollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contractmockrollup *ContractmockrollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contractmockrollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contractmockrollup *ContractmockrollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contractmockrollup *ContractmockrollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.contract.Transact(opts, method, params...)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCaller) Blacklist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "blacklist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupSession) Blacklist(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Blacklist(&_Contractmockrollup.CallOpts, arg0)
}

// Blacklist is a free data retrieval call binding the contract method 0xf9f92be4.
//
// Solidity: function blacklist(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCallerSession) Blacklist(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Blacklist(&_Contractmockrollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_Contractmockrollup *ContractmockrollupCaller) Commitments(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Validator            common.Address
		DataLength           uint32
		PolynomialCommitment BN254G1Point
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataLength = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.PolynomialCommitment = *abi.ConvertType(out[2], new(BN254G1Point)).(*BN254G1Point)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_Contractmockrollup *ContractmockrollupSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _Contractmockrollup.Contract.Commitments(&_Contractmockrollup.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x49ce8997.
//
// Solidity: function commitments(uint256 ) view returns(address validator, uint32 dataLength, (uint256,uint256) polynomialCommitment)
func (_Contractmockrollup *ContractmockrollupCallerSession) Commitments(arg0 *big.Int) (struct {
	Validator            common.Address
	DataLength           uint32
	PolynomialCommitment BN254G1Point
}, error) {
	return _Contractmockrollup.Contract.Commitments(&_Contractmockrollup.CallOpts, arg0)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupCaller) EigenDAServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "eigenDAServiceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupSession) EigenDAServiceManager() (common.Address, error) {
	return _Contractmockrollup.Contract.EigenDAServiceManager(&_Contractmockrollup.CallOpts)
}

// EigenDAServiceManager is a free data retrieval call binding the contract method 0xfc30cad0.
//
// Solidity: function eigenDAServiceManager() view returns(address)
func (_Contractmockrollup *ContractmockrollupCallerSession) EigenDAServiceManager() (common.Address, error) {
	return _Contractmockrollup.Contract.EigenDAServiceManager(&_Contractmockrollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCaller) IllegalValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "illegalValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupSession) IllegalValue() (*big.Int, error) {
	return _Contractmockrollup.Contract.IllegalValue(&_Contractmockrollup.CallOpts)
}

// IllegalValue is a free data retrieval call binding the contract method 0x4440bc5c.
//
// Solidity: function illegalValue() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCallerSession) IllegalValue() (*big.Int, error) {
	return _Contractmockrollup.Contract.IllegalValue(&_Contractmockrollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCaller) StakeRequired(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "stakeRequired")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupSession) StakeRequired() (*big.Int, error) {
	return _Contractmockrollup.Contract.StakeRequired(&_Contractmockrollup.CallOpts)
}

// StakeRequired is a free data retrieval call binding the contract method 0x05010105.
//
// Solidity: function stakeRequired() view returns(uint256)
func (_Contractmockrollup *ContractmockrollupCallerSession) StakeRequired() (*big.Int, error) {
	return _Contractmockrollup.Contract.StakeRequired(&_Contractmockrollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_Contractmockrollup *ContractmockrollupCaller) Tau(opts *bind.CallOpts) (struct {
	X *big.Int
	Y *big.Int
}, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "tau")

	outstruct := new(struct {
		X *big.Int
		Y *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.X = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Y = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_Contractmockrollup *ContractmockrollupSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _Contractmockrollup.Contract.Tau(&_Contractmockrollup.CallOpts)
}

// Tau is a free data retrieval call binding the contract method 0xcfc4af55.
//
// Solidity: function tau() view returns(uint256 X, uint256 Y)
func (_Contractmockrollup *ContractmockrollupCallerSession) Tau() (struct {
	X *big.Int
	Y *big.Int
}, error) {
	return _Contractmockrollup.Contract.Tau(&_Contractmockrollup.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Contractmockrollup.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupSession) Validators(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Validators(&_Contractmockrollup.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_Contractmockrollup *ContractmockrollupCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _Contractmockrollup.Contract.Validators(&_Contractmockrollup.CallOpts, arg0)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupTransactor) ChallengeCommitment(opts *bind.TransactOpts, timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.contract.Transact(opts, "challengeCommitment", timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ChallengeCommitment(&_Contractmockrollup.TransactOpts, timestamp, point, proof)
}

// ChallengeCommitment is a paid mutator transaction binding the contract method 0x6281e63b.
//
// Solidity: function challengeCommitment(uint256 timestamp, uint256 point, (uint256[2],uint256[2]) proof) returns()
func (_Contractmockrollup *ContractmockrollupTransactorSession) ChallengeCommitment(timestamp *big.Int, point *big.Int, proof BN254G2Point) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.ChallengeCommitment(&_Contractmockrollup.TransactOpts, timestamp, point, proof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x4114c4ca.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_Contractmockrollup *ContractmockrollupTransactor) PostCommitment(opts *bind.TransactOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _Contractmockrollup.contract.Transact(opts, "postCommitment", blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x4114c4ca.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_Contractmockrollup *ContractmockrollupSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.PostCommitment(&_Contractmockrollup.TransactOpts, blobHeader, blobVerificationProof)
}

// PostCommitment is a paid mutator transaction binding the contract method 0x4114c4ca.
//
// Solidity: function postCommitment(((uint256,uint256),uint32,(uint8,uint8,uint8,uint8)[]) blobHeader, (uint32,uint8,((bytes32,bytes,bytes,uint32),bytes32,uint96,uint32),bytes,bytes) blobVerificationProof) returns()
func (_Contractmockrollup *ContractmockrollupTransactorSession) PostCommitment(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof EigenDABlobUtilsBlobVerificationProof) (*types.Transaction, error) {
	return _Contractmockrollup.Contract.PostCommitment(&_Contractmockrollup.TransactOpts, blobHeader, blobVerificationProof)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupTransactor) RegisterValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contractmockrollup.contract.Transact(opts, "registerValidator")
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupSession) RegisterValidator() (*types.Transaction, error) {
	return _Contractmockrollup.Contract.RegisterValidator(&_Contractmockrollup.TransactOpts)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbcc6587f.
//
// Solidity: function registerValidator() payable returns()
func (_Contractmockrollup *ContractmockrollupTransactorSession) RegisterValidator() (*types.Transaction, error) {
	return _Contractmockrollup.Contract.RegisterValidator(&_Contractmockrollup.TransactOpts)
}
