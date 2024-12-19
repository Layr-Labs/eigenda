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

// BatchHeader is an auto generated low-level Go binding around an user-defined struct.
type BatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
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

// IRewardsCoordinatorOperatorDirectedRewardsSubmission is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorOperatorDirectedRewardsSubmission struct {
	StrategiesAndMultipliers []IRewardsCoordinatorStrategyAndMultiplier
	Token                    common.Address
	OperatorRewards          []IRewardsCoordinatorOperatorReward
	StartTimestamp           uint32
	Duration                 uint32
	Description              string
}

// IRewardsCoordinatorOperatorReward is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorOperatorReward struct {
	Operator common.Address
	Amount   *big.Int
}

// IRewardsCoordinatorRewardsSubmission is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorRewardsSubmission struct {
	StrategiesAndMultipliers []IRewardsCoordinatorStrategyAndMultiplier
	Token                    common.Address
	Amount                   *big.Int
	StartTimestamp           uint32
	Duration                 uint32
}

// IRewardsCoordinatorStrategyAndMultiplier is an auto generated low-level Go binding around an user-defined struct.
type IRewardsCoordinatorStrategyAndMultiplier struct {
	Strategy   common.Address
	Multiplier *big.Int
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// SecurityThresholds is an auto generated low-level Go binding around an user-defined struct.
type SecurityThresholds struct {
	ConfirmationThreshold uint8
	AdversaryThreshold    uint8
}

// VersionedBlobParams is an auto generated low-level Go binding around an user-defined struct.
type VersionedBlobParams struct {
	MaxNumOperators uint32
	NumChunks       uint32
	CodingRate      uint8
}

// ContractIEigenDAServiceManagerMetaData contains all meta data concerning the ContractIEigenDAServiceManager contract.
var ContractIEigenDAServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAVSRewardsSubmission\",\"inputs\":[{\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createOperatorDirectedAVSRewardsSubmission\",\"inputs\":[{\"name\":\"operatorDirectedRewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.OperatorDirectedRewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"operatorRewards\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.OperatorReward[]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestServeUntilBlock\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setClaimerFor\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchConfirmerStatusChanged\",\"inputs\":[{\"name\":\"batchConfirmer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsInitiatorUpdated\",\"inputs\":[{\"name\":\"prevRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
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

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(VersionedBlobParams)).(*VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractIEigenDAServiceManager.Contract.GetBlobParams(&_ContractIEigenDAServiceManager.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractIEigenDAServiceManager.Contract.GetBlobParams(&_ContractIEigenDAServiceManager.CallOpts, version)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetDefaultSecurityThresholdsV2(opts *bind.CallOpts) (SecurityThresholds, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getDefaultSecurityThresholdsV2")

	if err != nil {
		return *new(SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(SecurityThresholds)).(*SecurityThresholds)

	return out0, err

}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractIEigenDAServiceManager.Contract.GetDefaultSecurityThresholdsV2(&_ContractIEigenDAServiceManager.CallOpts)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractIEigenDAServiceManager.Contract.GetDefaultSecurityThresholdsV2(&_ContractIEigenDAServiceManager.CallOpts)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractIEigenDAServiceManager.Contract.GetIsQuorumRequired(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractIEigenDAServiceManager.Contract.GetIsQuorumRequired(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
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

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDAServiceManager.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractIEigenDAServiceManager.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractIEigenDAServiceManager.CallOpts, quorumNumber)
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

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) LatestServeUntilBlock(opts *bind.CallOpts, referenceBlockNumber uint32) (uint32, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "latestServeUntilBlock", referenceBlockNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractIEigenDAServiceManager.CallOpts, referenceBlockNumber)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _ContractIEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractIEigenDAServiceManager.CallOpts, referenceBlockNumber)
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

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractIEigenDAServiceManager.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumNumbersRequired(&_ContractIEigenDAServiceManager.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractIEigenDAServiceManager.Contract.QuorumNumbersRequired(&_ContractIEigenDAServiceManager.CallOpts)
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
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) ConfirmBatch(opts *bind.TransactOpts, batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "confirmBatch", batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) ConfirmBatch(batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ConfirmBatch(&_ContractIEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) ConfirmBatch(batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.ConfirmBatch(&_ContractIEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) CreateAVSRewardsSubmission(opts *bind.TransactOpts, rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "createAVSRewardsSubmission", rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractIEigenDAServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractIEigenDAServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) CreateOperatorDirectedAVSRewardsSubmission(opts *bind.TransactOpts, operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "createOperatorDirectedAVSRewardsSubmission", operatorDirectedRewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) CreateOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.CreateOperatorDirectedAVSRewardsSubmission(&_ContractIEigenDAServiceManager.TransactOpts, operatorDirectedRewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) CreateOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.CreateOperatorDirectedAVSRewardsSubmission(&_ContractIEigenDAServiceManager.TransactOpts, operatorDirectedRewardsSubmissions)
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

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) SetClaimerFor(opts *bind.TransactOpts, claimer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "setClaimerFor", claimer)
}

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) SetClaimerFor(claimer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetClaimerFor(&_ContractIEigenDAServiceManager.TransactOpts, claimer)
}

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) SetClaimerFor(claimer common.Address) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.SetClaimerFor(&_ContractIEigenDAServiceManager.TransactOpts, claimer)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactor) UpdateAVSMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.contract.Transact(opts, "updateAVSMetadataURI", _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.UpdateAVSMetadataURI(&_ContractIEigenDAServiceManager.TransactOpts, _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerTransactorSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractIEigenDAServiceManager.Contract.UpdateAVSMetadataURI(&_ContractIEigenDAServiceManager.TransactOpts, _metadataURI)
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

// ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator is returned from FilterBatchConfirmerStatusChanged and is used to iterate over the raw logs and unpacked data for BatchConfirmerStatusChanged events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator struct {
	Event *ContractIEigenDAServiceManagerBatchConfirmerStatusChanged // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerBatchConfirmerStatusChanged)
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
		it.Event = new(ContractIEigenDAServiceManagerBatchConfirmerStatusChanged)
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
func (it *ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerBatchConfirmerStatusChanged represents a BatchConfirmerStatusChanged event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerBatchConfirmerStatusChanged struct {
	BatchConfirmer common.Address
	Status         bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmerStatusChanged is a free log retrieval operation binding the contract event 0x5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d414.
//
// Solidity: event BatchConfirmerStatusChanged(address batchConfirmer, bool status)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterBatchConfirmerStatusChanged(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "BatchConfirmerStatusChanged")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerBatchConfirmerStatusChangedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "BatchConfirmerStatusChanged", logs: logs, sub: sub}, nil
}

// WatchBatchConfirmerStatusChanged is a free log subscription operation binding the contract event 0x5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d414.
//
// Solidity: event BatchConfirmerStatusChanged(address batchConfirmer, bool status)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchBatchConfirmerStatusChanged(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerBatchConfirmerStatusChanged) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "BatchConfirmerStatusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerBatchConfirmerStatusChanged)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerStatusChanged", log); err != nil {
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

// ParseBatchConfirmerStatusChanged is a log parse operation binding the contract event 0x5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d414.
//
// Solidity: event BatchConfirmerStatusChanged(address batchConfirmer, bool status)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseBatchConfirmerStatusChanged(log types.Log) (*ContractIEigenDAServiceManagerBatchConfirmerStatusChanged, error) {
	event := new(ContractIEigenDAServiceManagerBatchConfirmerStatusChanged)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
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
func (it *ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 SecurityThresholds
	NewDefaultSecurityThresholdsV2      SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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

// ParseDefaultSecurityThresholdsV2Updated is a log parse operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractIEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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

// ParseQuorumAdversaryThresholdPercentagesUpdated is a log parse operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractIEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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

// ParseQuorumConfirmationThresholdPercentagesUpdated is a log parse operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractIEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated)
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
func (it *ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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

// ParseQuorumNumbersRequiredUpdated is a log parse operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated, error) {
	event := new(ContractIEigenDAServiceManagerQuorumNumbersRequiredUpdated)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator is returned from FilterRewardsInitiatorUpdated and is used to iterate over the raw logs and unpacked data for RewardsInitiatorUpdated events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator struct {
	Event *ContractIEigenDAServiceManagerRewardsInitiatorUpdated // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerRewardsInitiatorUpdated)
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
		it.Event = new(ContractIEigenDAServiceManagerRewardsInitiatorUpdated)
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
func (it *ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerRewardsInitiatorUpdated represents a RewardsInitiatorUpdated event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerRewardsInitiatorUpdated struct {
	PrevRewardsInitiator common.Address
	NewRewardsInitiator  common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRewardsInitiatorUpdated is a free log retrieval operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterRewardsInitiatorUpdated(opts *bind.FilterOpts) (*ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerRewardsInitiatorUpdatedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "RewardsInitiatorUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardsInitiatorUpdated is a free log subscription operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchRewardsInitiatorUpdated(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerRewardsInitiatorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerRewardsInitiatorUpdated)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
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

// ParseRewardsInitiatorUpdated is a log parse operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseRewardsInitiatorUpdated(log types.Log) (*ContractIEigenDAServiceManagerRewardsInitiatorUpdated, error) {
	event := new(ContractIEigenDAServiceManagerRewardsInitiatorUpdated)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator struct {
	Event *ContractIEigenDAServiceManagerVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIEigenDAServiceManagerVersionedBlobParamsAdded)
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
		it.Event = new(ContractIEigenDAServiceManagerVersionedBlobParamsAdded)
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
func (it *ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIEigenDAServiceManagerVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractIEigenDAServiceManager contract.
type ContractIEigenDAServiceManagerVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractIEigenDAServiceManager.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractIEigenDAServiceManagerVersionedBlobParamsAddedIterator{contract: _ContractIEigenDAServiceManager.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractIEigenDAServiceManagerVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractIEigenDAServiceManager.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIEigenDAServiceManagerVersionedBlobParamsAdded)
				if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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

// ParseVersionedBlobParamsAdded is a log parse operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractIEigenDAServiceManager *ContractIEigenDAServiceManagerFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractIEigenDAServiceManagerVersionedBlobParamsAdded, error) {
	event := new(ContractIEigenDAServiceManagerVersionedBlobParamsAdded)
	if err := _ContractIEigenDAServiceManager.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
