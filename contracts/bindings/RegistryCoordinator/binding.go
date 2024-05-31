// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractRegistryCoordinator

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

// IBLSApkRegistryPubkeyRegistrationParams is an auto generated low-level Go binding around an user-defined struct.
type IBLSApkRegistryPubkeyRegistrationParams struct {
	PubkeyRegistrationSignature BN254G1Point
	PubkeyG1                    BN254G1Point
	PubkeyG2                    BN254G2Point
}

// IRegistryCoordinatorOperatorInfo is an auto generated low-level Go binding around an user-defined struct.
type IRegistryCoordinatorOperatorInfo struct {
	OperatorId [32]byte
	Status     uint8
}

// IRegistryCoordinatorOperatorKickParam is an auto generated low-level Go binding around an user-defined struct.
type IRegistryCoordinatorOperatorKickParam struct {
	QuorumNumber uint8
	Operator     common.Address
}

// IRegistryCoordinatorOperatorSetParam is an auto generated low-level Go binding around an user-defined struct.
type IRegistryCoordinatorOperatorSetParam struct {
	MaxOperatorCount        uint32
	KickBIPsOfOperatorStake uint16
	KickBIPsOfTotalStake    uint16
}

// IRegistryCoordinatorQuorumBitmapUpdate is an auto generated low-level Go binding around an user-defined struct.
type IRegistryCoordinatorQuorumBitmapUpdate struct {
	UpdateBlockNumber     uint32
	NextUpdateBlockNumber uint32
	QuorumBitmap          *big.Int
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// IStakeRegistryStrategyParams is an auto generated low-level Go binding around an user-defined struct.
type IStakeRegistryStrategyParams struct {
	Strategy   common.Address
	Multiplier *big.Int
}

// ContractRegistryCoordinatorMetaData contains all meta data concerning the ContractRegistryCoordinator contract.
var ContractRegistryCoordinatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_serviceManager\",\"type\":\"address\",\"internalType\":\"contractIServiceManager\"},{\"name\":\"_stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"},{\"name\":\"_blsApkRegistry\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"},{\"name\":\"_indexRegistry\",\"type\":\"address\",\"internalType\":\"contractIIndexRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"OPERATOR_CHURN_APPROVAL_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PUBKEY_REGISTRATION_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateOperatorChurnApprovalDigestHash\",\"inputs\":[{\"name\":\"registeringOperator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registeringOperatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"operatorKickParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorKickParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"churnApprover\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createQuorum\",\"inputs\":[{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"name\":\"minimumStake\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"strategyParams\",\"type\":\"tuple[]\",\"internalType\":\"structIStakeRegistry.StrategyParams[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectionCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejector\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentQuorumBitmap\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorInfo\",\"components\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRegistryCoordinator.OperatorStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorFromId\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorId\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorSetParams\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorStatus\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIRegistryCoordinator.OperatorStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapAtBlockNumberByIndex\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapHistoryLength\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapIndicesAtBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapUpdateByIndex\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.QuorumBitmapUpdate\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBitmap\",\"type\":\"uint192\",\"internalType\":\"uint192\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"indexRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIndexRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_churnApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_ejector\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_operatorSetParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam[]\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"name\":\"_minimumStakes\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"_strategyParams\",\"type\":\"tuple[][]\",\"internalType\":\"structIStakeRegistry.StrategyParams[][]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isChurnApproverSaltUsed\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastEjectionTimestamp\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"numRegistries\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyRegistrationMessageHash\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumUpdateBlockNumber\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.PubkeyRegistrationParams\",\"components\":[{\"name\":\"pubkeyRegistrationSignature\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}]},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperatorWithChurn\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.PubkeyRegistrationParams\",\"components\":[{\"name\":\"pubkeyRegistrationSignature\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}]},{\"name\":\"operatorKickParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorKickParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"churnApproverSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registries\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"serviceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setChurnApprover\",\"inputs\":[{\"name\":\"_churnApprover\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEjectionCooldown\",\"inputs\":[{\"name\":\"_ejectionCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEjector\",\"inputs\":[{\"name\":\"_ejector\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorSetParams\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperatorsForQuorum\",\"inputs\":[{\"name\":\"operatorsPerQuorum\",\"type\":\"address[][]\",\"internalType\":\"address[][]\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSocket\",\"inputs\":[{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ChurnApproverUpdated\",\"inputs\":[{\"name\":\"prevChurnApprover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newChurnApprover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EjectorUpdated\",\"inputs\":[{\"name\":\"prevEjector\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newEjector\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSetParamsUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSocketUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"socket\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumBlockNumberUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"blocknumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6101c06040523480156200001257600080fd5b506040516200638938038062006389833981016040819052620000359162000254565b604080518082018252601681527f4156535265676973747279436f6f7264696e61746f720000000000000000000060208083019182528351808501909452600684526576302e302e3160d01b908401528151902060e08190527f6bda7e3f385e48841048390444cced5cc795af87758af67622e5f4f0882c4a996101008190524660a05287938793879387939192917f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f620001358184846040805160208101859052908101839052606081018290524660808201523060a082015260009060c0016040516020818303038152906040528051906020012090509392505050565b6080523060c05261012052505050506001600160a01b039384166101405291831661018052821661016052166101a0526200016f62000179565b50505050620002bc565b600054610100900460ff1615620001e65760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000239576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200025157600080fd5b50565b600080600080608085870312156200026b57600080fd5b845162000278816200023b565b60208601519094506200028b816200023b565b60408601519093506200029e816200023b565b6060860151909250620002b1816200023b565b939692955090935050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a051615fc5620003c4600039600081816106ab0152818161119d0152818161217301528181612ea50152818161375c0152613d340152600081816105f0015281816120fe015281816125a601528181612e25015281816136b3015281816139090152613cb30152600081816105b601528181610f380152818161213c01528181612da701528181612f8d01528181613003015281816136330152613db00152600081816104fa01528181612cfd015261357b01526000613fb70152600061400601526000613fe101526000613f3a01526000613f6401526000613f8e0152615fc56000f3fe608060405234801561001057600080fd5b50600436106102d55760003560e01c80635df45946116101825780639feab859116100e9578063d72d8dd6116100a2578063e65797ad1161007c578063e65797ad14610798578063f2fde38b1461083b578063fabc1cbc1461084e578063fd39105a1461086157600080fd5b8063d72d8dd61461076a578063d75b4c8814610772578063dd8283f31461078557600080fd5b80639feab859146106cd578063a50857bf146106f4578063a96f783e14610707578063c391425e14610710578063ca0de88214610730578063ca4f2d971461075757600080fd5b8063871ef0491161013b578063871ef04914610640578063886f1195146106535780638da5cb5b1461066c5780639aa1653d146106745780639b5d177b146106935780639e9923c2146106a657600080fd5b80635df45946146105b15780636347c900146105d857806368304835146105eb5780636e3b17db14610612578063715018a61461062557806384ca52131461062d57600080fd5b8063249a0c42116102415780633c2a7f4c116101fa578063595c6a67116101d4578063595c6a671461056f5780635ac86ab7146105775780635b0b829f146105965780635c975abb146105a957600080fd5b80633c2a7f4c1461051c5780635140a5481461053c5780635865c60c1461054f57600080fd5b8063249a0c421461048957806328f61b31146104a9578063296bb064146104bc57806329d1e0c3146104cf5780632cdd1e86146104e25780633998fdd3146104f557600080fd5b806310d67a2f1161029357806310d67a2f1461039e578063125e0584146103b157806313542a4e146103d1578063136439dd146103fa5780631478851f1461040d5780631eb812da1461044057600080fd5b8062cf2ab5146102da57806303fd3492146102ef57806304ec635114610322578063054310e61461034d5780630cf4b767146103785780630d3f21341461038b575b600080fd5b6102ed6102e8366004614ab7565b61089d565b005b61030f6102fd366004614af8565b60009081526098602052604090205490565b6040519081526020015b60405180910390f35b610335610330366004614b23565b6109b3565b6040516001600160c01b039091168152602001610319565b609d54610360906001600160a01b031681565b6040516001600160a01b039091168152602001610319565b6102ed610386366004614c42565b610ba9565b6102ed610399366004614af8565b610c91565b6102ed6103ac366004614cb7565b610c9e565b61030f6103bf366004614cb7565b609f6020526000908152604090205481565b61030f6103df366004614cb7565b6001600160a01b031660009081526099602052604090205490565b6102ed610408366004614af8565b610d51565b61043061041b366004614af8565b609a6020526000908152604090205460ff1681565b6040519015158152602001610319565b61045361044e366004614cd4565b610e8e565b60408051825163ffffffff908116825260208085015190911690820152918101516001600160c01b031690820152606001610319565b61030f610497366004614d07565b609b6020526000908152604090205481565b609e54610360906001600160a01b031681565b6103606104ca366004614af8565b610f1f565b6102ed6104dd366004614cb7565b610fab565b6102ed6104f0366004614cb7565b610fbc565b6103607f000000000000000000000000000000000000000000000000000000000000000081565b61052f61052a366004614cb7565b610fcd565b6040516103199190614d22565b6102ed61054a366004614d7a565b61104c565b61056261055d366004614cb7565b61155d565b6040516103199190614e1d565b6102ed6115d1565b610430610585366004614d07565b6001805460ff9092161b9081161490565b6102ed6105a4366004614ea2565b61169d565b60015461030f565b6103607f000000000000000000000000000000000000000000000000000000000000000081565b6103606105e6366004614af8565b611734565b6103607f000000000000000000000000000000000000000000000000000000000000000081565b6102ed610620366004614ed6565b61175e565b6102ed6118ec565b61030f61063b366004614f8d565b611900565b61033561064e366004614af8565b61194a565b600054610360906201000090046001600160a01b031681565b610360611955565b6096546106819060ff1681565b60405160ff9091168152602001610319565b6102ed6106a1366004615126565b61196e565b6103607f000000000000000000000000000000000000000000000000000000000000000081565b61030f7f2bd82124057f0913bc3b772ce7b83e8057c1ad1f3510fc83778be20f10ec5de681565b6102ed61070236600461521f565b611ca6565b61030f60a05481565b61072361071e3660046152c7565b611e2a565b604051610319919061536c565b61030f7f4d404e3276e7ac2163d8ee476afa6a41d1f68fb71f2d8b6546b24e55ce01b72a81565b6102ed6107653660046153b6565b611ee3565b609c5461030f565b6102ed61078036600461549c565b611f4a565b6102ed61079336600461564f565b611f5d565b6108076107a6366004614d07565b60408051606080820183526000808352602080840182905292840181905260ff9490941684526097825292829020825193840183525463ffffffff8116845261ffff600160201b8204811692850192909252600160301b9004169082015290565b60408051825163ffffffff16815260208084015161ffff908116918301919091529282015190921690820152606001610319565b6102ed610849366004614cb7565b612261565b6102ed61085c366004614af8565b6122d7565b61089061086f366004614cb7565b6001600160a01b031660009081526099602052604090206001015460ff1690565b6040516103199190615723565b600154600290600490811614156108cf5760405162461bcd60e51b81526004016108c690615731565b60405180910390fd5b60005b828110156109ad5760008484838181106108ee576108ee615768565b90506020020160208101906109039190614cb7565b6001600160a01b03811660009081526099602090815260408083208151808301909252805482526001810154949550929390929183019060ff16600281111561094e5761094e614de5565b600281111561095f5761095f614de5565b9052508051909150600061097282612433565b90506000610988826001600160c01b031661249c565b9050610995858583612568565b505050505080806109a590615794565b9150506108d2565b50505050565b60008381526098602052604081208054829190849081106109d6576109d6615768565b600091825260209182902060408051606081018252929091015463ffffffff808216808552600160201b8304821695850195909552600160401b9091046001600160c01b03169183019190915290925085161015610ad05760405162461bcd60e51b815260206004820152606560248201527f5265676973747279436f6f7264696e61746f722e67657451756f72756d42697460448201527f6d61704174426c6f636b4e756d6265724279496e6465783a2071756f72756d4260648201527f69746d61705570646174652069732066726f6d20616674657220626c6f636b4e6084820152643ab6b132b960d91b60a482015260c4016108c6565b602081015163ffffffff161580610af65750806020015163ffffffff168463ffffffff16105b610b9d5760405162461bcd60e51b815260206004820152606660248201527f5265676973747279436f6f7264696e61746f722e67657451756f72756d42697460448201527f6d61704174426c6f636b4e756d6265724279496e6465783a2071756f72756d4260648201527f69746d61705570646174652069732066726f6d206265666f726520626c6f636b608482015265273ab6b132b960d11b60a482015260c4016108c6565b60400151949350505050565b60013360009081526099602052604090206001015460ff166002811115610bd257610bd2614de5565b14610c455760405162461bcd60e51b815260206004820152603c60248201527f5265676973747279436f6f7264696e61746f722e757064617465536f636b657460448201527f3a206f70657261746f72206973206e6f7420726567697374657265640000000060648201526084016108c6565b33600090815260996020526040908190205490517fec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa90610c869084906157fc565b60405180910390a250565b610c99612655565b60a055565b600060029054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610cf1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d15919061580f565b6001600160a01b0316336001600160a01b031614610d455760405162461bcd60e51b81526004016108c69061582c565b610d4e816126b4565b50565b60005460405163237dfb4760e11b8152336004820152620100009091046001600160a01b0316906346fbf68e90602401602060405180830381865afa158015610d9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dc29190615876565b610dde5760405162461bcd60e51b81526004016108c690615898565b60015481811614610e575760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c697479000000000000000060648201526084016108c6565b600181905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d90602001610c86565b60408051606081018252600080825260208201819052918101919091526000838152609860205260409020805483908110610ecb57610ecb615768565b600091825260209182902060408051606081018252919092015463ffffffff8082168352600160201b820416938201939093526001600160c01b03600160401b909304929092169082015290505b92915050565b6040516308f6629d60e31b8152600481018290526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316906347b314e890602401602060405180830381865afa158015610f87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f19919061580f565b610fb3612655565b610d4e816127b9565b610fc4612655565b610d4e81612822565b6040805180820190915260008082526020820152610f196110477f2bd82124057f0913bc3b772ce7b83e8057c1ad1f3510fc83778be20f10ec5de68460405160200161102c9291909182526001600160a01b0316602082015260400190565b6040516020818303038152906040528051906020012061288b565b6128d9565b600154600290600490811614156110755760405162461bcd60e51b81526004016108c690615731565b60006110bd84848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060965460ff1691506129699050565b905084831461112e5760405162461bcd60e51b81526020600482015260436024820152600080516020615f3083398151915260448201527f6f7273466f7251756f72756d3a20696e707574206c656e677468206d69736d616064820152620e8c6d60eb1b608482015260a4016108c6565b60005b8381101561155457600085858381811061114d5761114d615768565b919091013560f81c9150369050600089898581811061116e5761116e615768565b905060200281019061118091906158e0565b6040516379a0849160e11b815260ff8616600482015291935091507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063f341092290602401602060405180830381865afa1580156111ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112109190615929565b63ffffffff1681146112ac5760405162461bcd60e51b81526020600482015260656024820152600080516020615f3083398151915260448201527f6f7273466f7251756f72756d3a206e756d626572206f6620757064617465642060648201527f6f70657261746f727320646f6573206e6f74206d617463682071756f72756d206084820152641d1bdd185b60da1b60a482015260c4016108c6565b6000805b828110156114f35760008484838181106112cc576112cc615768565b90506020020160208101906112e19190614cb7565b6001600160a01b03811660009081526099602090815260408083208151808301909252805482526001810154949550929390929183019060ff16600281111561132c5761132c614de5565b600281111561133d5761133d614de5565b9052508051909150600061135082612433565b905060016001600160c01b03821660ff8b161c8116146113d45760405162461bcd60e51b815260206004820152604460248201819052600080516020615f30833981519152908201527f6f7273466f7251756f72756d3a206f70657261746f72206e6f7420696e2071756064820152636f72756d60e01b608482015260a4016108c6565b856001600160a01b0316846001600160a01b03161161147f5760405162461bcd60e51b81526020600482015260676024820152600080516020615f3083398151915260448201527f6f7273466f7251756f72756d3a206f70657261746f7273206172726179206d7560648201527f737420626520736f7274656420696e20617363656e64696e6720616464726573608482015266399037b93232b960c91b60a482015260c4016108c6565b506114dd83838f8f8d908e60016114969190615946565b926114a39392919061595e565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061256892505050565b509092506114ec905081615794565b90506112b0565b5060ff84166000818152609b6020908152604091829020439081905591519182527f46077d55330763f16269fd75e5761663f4192d2791747c0189b16ad31db07db4910160405180910390a2505050508061154d90615794565b9050611131565b50505050505050565b60408051808201909152600080825260208201526001600160a01b0382166000908152609960209081526040918290208251808401909352805483526001810154909183019060ff1660028111156115b7576115b7614de5565b60028111156115c8576115c8614de5565b90525092915050565b60005460405163237dfb4760e11b8152336004820152620100009091046001600160a01b0316906346fbf68e90602401602060405180830381865afa15801561161e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116429190615876565b61165e5760405162461bcd60e51b81526004016108c690615898565b600019600181905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b6116a5612655565b609654829060ff908116908216106117255760405162461bcd60e51b815260206004820152603760248201527f5265676973747279436f6f7264696e61746f722e71756f72756d45786973747360448201527f3a2071756f72756d20646f6573206e6f7420657869737400000000000000000060648201526084016108c6565b61172f83836129fa565b505050565b609c818154811061174457600080fd5b6000918252602090912001546001600160a01b0316905081565b609e546001600160a01b031633146117de5760405162461bcd60e51b815260206004820152603a60248201527f5265676973747279436f6f7264696e61746f722e6f6e6c79456a6563746f723a60448201527f2063616c6c6572206973206e6f742074686520656a6563746f7200000000000060648201526084016108c6565b6001600160a01b0383166000908152609f602090815260408083204290556099825280832080548251601f870185900485028101850190935285835290939092909161184b9187908790819084018382808284376000920191909152505060965460ff1691506129699050565b9050600061185883612433565b905060018085015460ff16600281111561187457611874614de5565b14801561188957506001600160c01b03821615155b80156118a757506118a76001600160c01b0383811690831681161490565b15611554576115548787878080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612aa792505050565b6118f4612655565b6118fe6000612f19565b565b60006119407f4d404e3276e7ac2163d8ee476afa6a41d1f68fb71f2d8b6546b24e55ce01b72a878787878760405160200161102c96959493929190615988565b9695505050505050565b6000610f1982612433565b60006119696064546001600160a01b031690565b905090565b6001805460009190811614156119965760405162461bcd60e51b81526004016108c690615731565b838914611a195760405162461bcd60e51b8152602060048201526044602482018190527f5265676973747279436f6f7264696e61746f722e72656769737465724f706572908201527f61746f7257697468436875726e3a20696e707574206c656e677468206d69736d6064820152630c2e8c6d60e31b608482015260a4016108c6565b6000611a253388612f6b565b9050611a8533828888808060200260200160405190810160405280939291908181526020016000905b82821015611a7a57611a6b60408302860136819003810190615a0d565b81526020019060010190611a4e565b50505050508761309c565b6000611acc33838e8e8e8e8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c9250613229915050565b905060005b8b811015611c97576000609760008f8f85818110611af157611af1615768565b919091013560f81c82525060208082019290925260409081016000208151606081018352905463ffffffff811680835261ffff600160201b8304811695840195909552600160301b90910490931691810191909152845180519193509084908110611b5e57611b5e615768565b602002602001015163ffffffff161115611c8457611bff8e8e84818110611b8757611b87615768565b9050013560f81c60f81b60f81c84604001518481518110611baa57611baa615768565b60200260200101513386602001518681518110611bc957611bc9615768565b60200260200101518d8d88818110611be357611be3615768565b905060400201803603810190611bf99190615a0d565b866137ea565b611c84898984818110611c1457611c14615768565b9050604002016020016020810190611c2c9190614cb7565b8f8f8590866001611c3d9190615946565b92611c4a9392919061595e565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612aa792505050565b5080611c8f81615794565b915050611ad1565b50505050505050505050505050565b600180546000919081161415611cce5760405162461bcd60e51b81526004016108c690615731565b6000611cda3385612f6b565b90506000611d2333838b8b8b8b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c9250613229915050565b51905060005b88811015611e1e5760008a8a83818110611d4557611d45615768565b919091013560f81c600081815260976020526040902054855191935063ffffffff169150849084908110611d7b57611d7b615768565b602002602001015163ffffffff161115611e0b5760405162461bcd60e51b8152602060048201526044602482018190527f5265676973747279436f6f7264696e61746f722e72656769737465724f706572908201527f61746f723a206f70657261746f7220636f756e742065786365656473206d6178606482015263696d756d60e01b608482015260a4016108c6565b5080611e1681615794565b915050611d29565b50505050505050505050565b6060600082516001600160401b03811115611e4757611e47614b5b565b604051908082528060200260200182016040528015611e70578160200160208202803683370190505b50905060005b8351811015611edb57611ea285858381518110611e9557611e95615768565b6020026020010151613abf565b828281518110611eb457611eb4615768565b63ffffffff9092166020928302919091019091015280611ed381615794565b915050611e76565b509392505050565b6001805460029081161415611f0a5760405162461bcd60e51b81526004016108c690615731565b61172f3384848080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250612aa792505050565b611f52612655565b61172f838383613bfb565b600054610100900460ff1615808015611f7d5750600054600160ff909116105b80611f975750303b158015611f97575060005460ff166001145b611ffa5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016108c6565b6000805460ff19166001179055801561201d576000805461ff0019166101001790555b8251845114801561202f575081518351145b6120995760405162461bcd60e51b815260206004820152603560248201527f5265676973747279436f6f7264696e61746f722e696e697469616c697a653a206044820152740d2dce0eae840d8cadccee8d040dad2e6dac2e8c6d605b1b60648201526084016108c6565b6120a289612f19565b6120ac8686613e12565b6120b5886127b9565b6120be87612822565b609c80546001818101835560008381527faf85b9071dfafeac1409d3f1d19bafc9bc7c37974cde8df0ee6168f0086e539c92830180546001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081166001600160a01b03199283161790925585548085018755850180547f0000000000000000000000000000000000000000000000000000000000000000841690831617905585549384019095559190920180547f000000000000000000000000000000000000000000000000000000000000000090921691909316179091555b845181101561220f576121fd8582815181106121bc576121bc615768565b60200260200101518583815181106121d6576121d6615768565b60200260200101518584815181106121f0576121f0615768565b6020026020010151613bfb565b8061220781615794565b91505061219e565b508015612256576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050505050565b612269612655565b6001600160a01b0381166122ce5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016108c6565b610d4e81612f19565b600060029054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561232a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061234e919061580f565b6001600160a01b0316336001600160a01b03161461237e5760405162461bcd60e51b81526004016108c69061582c565b6001541981196001541916146123fc5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c697479000000000000000060648201526084016108c6565b600181905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c90602001610c86565b600081815260986020526040812054806124505750600092915050565b6000838152609860205260409020612469600183615a29565b8154811061247957612479615768565b600091825260209091200154600160401b90046001600160c01b03169392505050565b60606000806124aa84613f02565b61ffff166001600160401b038111156124c5576124c5614b5b565b6040519080825280601f01601f1916602001820160405280156124ef576020820181803683370190505b5090506000805b825182108015612507575061010081105b1561255e576001811b93508584161561254e578060f81b83838151811061253057612530615768565b60200101906001600160f81b031916908160001a9053508160010191505b61255781615794565b90506124f6565b5090949350505050565b60018260200151600281111561258057612580614de5565b1461258a57505050565b81516040516333567f7f60e11b81526000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906366acfefe906125df90889086908890600401615a40565b6020604051808303816000875af11580156125fe573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126229190615a70565b90506001600160c01b0381161561264e5761264e85612649836001600160c01b031661249c565b612aa7565b5050505050565b3361265e611955565b6001600160a01b0316146118fe5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016108c6565b6001600160a01b0381166127425760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a4016108c6565b600054604080516001600160a01b03620100009093048316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a1600080546001600160a01b03909216620100000262010000600160b01b0319909216919091179055565b609d54604080516001600160a01b03928316815291831660208301527f315457d8a8fe60f04af17c16e2f5a5e1db612b31648e58030360759ef8f3528c910160405180910390a1609d80546001600160a01b0319166001600160a01b0392909216919091179055565b609e54604080516001600160a01b03928316815291831660208301527f8f30ab09f43a6c157d7fce7e0a13c003042c1c95e8a72e7a146a21c0caa24dc9910160405180910390a1609e80546001600160a01b0319166001600160a01b0392909216919091179055565b6000610f19612898613f2d565b8360405161190160f01b6020820152602281018390526042810182905260009060620160405160208183030381529060405280519060200120905092915050565b604080518082019091526000808252602082015260008080612909600080516020615f7083398151915286615aaf565b90505b61291581614054565b9093509150600080516020615f7083398151915282830983141561294f576040805180820190915290815260208101919091529392505050565b600080516020615f7083398151915260018208905061290c565b600080612975846140d6565b9050808360ff166001901b116129f35760405162461bcd60e51b815260206004820152603f60248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206269746d61702065786365656473206d61782076616c75650060648201526084016108c6565b9392505050565b60ff8216600081815260976020908152604091829020845181548684018051888701805163ffffffff90951665ffffffffffff199094168417600160201b61ffff938416021767ffff0000000000001916600160301b95831695909502949094179094558551918252518316938101939093525116918101919091527f3ee6fe8d54610244c3e9d3c066ae4aee997884aa28f10616ae821925401318ac9060600160405180910390a25050565b6001600160a01b0382166000908152609960205260409020805460018083015460ff166002811115612adb57612adb614de5565b14612b5a5760405162461bcd60e51b815260206004820152604360248201527f5265676973747279436f6f7264696e61746f722e5f646572656769737465724f60448201527f70657261746f723a206f70657261746f72206973206e6f7420726567697374656064820152621c995960ea1b608482015260a4016108c6565b609654600090612b6e90859060ff16612969565b90506000612b7b83612433565b90506001600160c01b038216612bf95760405162461bcd60e51b815260206004820152603b60248201527f5265676973747279436f6f7264696e61746f722e5f646572656769737465724f60448201527f70657261746f723a206269746d61702063616e6e6f742062652030000000000060648201526084016108c6565b612c106001600160c01b0383811690831681161490565b612ca85760405162461bcd60e51b815260206004820152605960248201527f5265676973747279436f6f7264696e61746f722e5f646572656769737465724f60448201527f70657261746f723a206f70657261746f72206973206e6f74207265676973746560648201527f72656420666f72207370656369666965642071756f72756d7300000000000000608482015260a4016108c6565b6001600160c01b0382811619821616612cc18482614263565b6001600160c01b038116612d905760018501805460ff191660021790556040516351b27a6d60e11b81526001600160a01b0388811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063a364f4da90602401600060405180830381600087803b158015612d4157600080fd5b505af1158015612d55573d6000803e3d6000fd5b50506040518692506001600160a01b038a1691507f396fdcb180cb0fea26928113fb0fd1c3549863f9cd563e6a184f1d578116c8e490600090a35b60405163f4e24fe560e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063f4e24fe590612dde908a908a90600401615ac3565b600060405180830381600087803b158015612df857600080fd5b505af1158015612e0c573d6000803e3d6000fd5b505060405163bd29b8cd60e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016925063bd29b8cd9150612e5e9087908a90600401615ae7565b600060405180830381600087803b158015612e7857600080fd5b505af1158015612e8c573d6000803e3d6000fd5b505060405163bd29b8cd60e01b81526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016925063bd29b8cd9150612ede9087908a90600401615ae7565b600060405180830381600087803b158015612ef857600080fd5b505af1158015612f0c573d6000803e3d6000fd5b5050505050505050505050565b606480546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6040516309aa152760e11b81526001600160a01b0383811660048301526000917f0000000000000000000000000000000000000000000000000000000000000000909116906313542a4e90602401602060405180830381865afa158015612fd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612ffa9190615b00565b905080610f19577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663bf79ce58848461303b87610fcd565b6040518463ffffffff1660e01b815260040161305993929190615b19565b6020604051808303816000875af1158015613078573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129f39190615b00565b6020808201516000908152609a909152604090205460ff16156131425760405162461bcd60e51b815260206004820152605260248201527f5265676973747279436f6f7264696e61746f722e5f766572696679436875726e60448201527f417070726f7665725369676e61747572653a20636875726e417070726f766572606482015271081cd85b1d08185b1c9958591e481d5cd95960721b608482015260a4016108c6565b42816040015110156131d75760405162461bcd60e51b815260206004820152605260248201527f5265676973747279436f6f7264696e61746f722e5f766572696679436875726e60448201527f417070726f7665725369676e61747572653a20636875726e417070726f766572606482015271081cda59db985d1d5c9948195e1c1a5c995960721b608482015260a4016108c6565b602080820180516000908152609a909252604091829020805460ff19166001179055609d549051918301516109ad926001600160a01b03909216916132229188918891889190611900565b8351614423565b61324d60405180606001604052806060815260200160608152602001606081525090565b600061329586868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060965460ff1691506129699050565b905060006132a288612433565b90506001600160c01b0382166133205760405162461bcd60e51b815260206004820152603960248201527f5265676973747279436f6f7264696e61746f722e5f72656769737465724f706560448201527f7261746f723a206269746d61702063616e6e6f7420626520300000000000000060648201526084016108c6565b8082166001600160c01b0316156133d65760405162461bcd60e51b815260206004820152606860248201527f5265676973747279436f6f7264696e61746f722e5f72656769737465724f706560448201527f7261746f723a206f70657261746f7220616c726561647920726567697374657260648201527f656420666f7220736f6d652071756f72756d73206265696e672072656769737460848201526732b932b2103337b960c11b60a482015260c4016108c6565b60a0546001600160a01b038a166000908152609f60205260409020546001600160c01b038381169085161791429161340e9190615946565b1061348f5760405162461bcd60e51b815260206004820152604560248201527f5265676973747279436f6f7264696e61746f722e5f72656769737465724f706560448201527f7261746f723a206f70657261746f722063616e6e6f74207265726567697374656064820152641c881e595d60da1b608482015260a4016108c6565b6134998982614263565b887fec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa876040516134c991906157fc565b60405180910390a260016001600160a01b038b1660009081526099602052604090206001015460ff16600281111561350357613503614de5565b1461361c576040805180820182528a8152600160208083018281526001600160a01b038f166000908152609990925293902082518155925183820180549394939192909160ff19169083600281111561355e5761355e614de5565b021790555050604051639926ee7d60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169150639926ee7d906135b3908d908990600401615b98565b600060405180830381600087803b1580156135cd57600080fd5b505af11580156135e1573d6000803e3d6000fd5b50506040518b92506001600160a01b038d1691507fe8e68cef1c3a761ed7be7e8463a375f27f7bc335e51824223cacce636ec5c3fe90600090a35b604051631fd93ca960e11b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690633fb279529061366c908d908c908c90600401615c0c565b600060405180830381600087803b15801561368657600080fd5b505af115801561369a573d6000803e3d6000fd5b5050604051632550477760e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169250632550477791506136f0908d908d908d908d90600401615c31565b6000604051808303816000875af115801561370f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526137379190810190615cbd565b60408087019190915260208601919091525162bff04d60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169062bff04d90613794908c908c908c90600401615d20565b6000604051808303816000875af11580156137b3573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526137db9190810190615d3a565b84525050509695505050505050565b6020808301516001600160a01b03808216600081815260999094526040909320549192908716141561386a5760405162461bcd60e51b81526020600482015260356024820152600080516020615f5083398151915260448201527439371d1031b0b73737ba1031b43ab9371039b2b63360591b60648201526084016108c6565b8760ff16846000015160ff16146138e75760405162461bcd60e51b81526020600482015260476024820152600080516020615f5083398151915260448201527f726e3a2071756f72756d4e756d626572206e6f74207468652073616d65206173606482015266081cda59db995960ca1b608482015260a4016108c6565b604051635401ed2760e01b81526004810182905260ff891660248201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690635401ed2790604401602060405180830381865afa158015613958573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061397c9190615dd3565b905061398881856145dd565b6001600160601b0316866001600160601b031611613a1b5760405162461bcd60e51b81526020600482015260566024820152600080516020615f5083398151915260448201527f726e3a20696e636f6d696e67206f70657261746f722068617320696e7375666660648201527534b1b4b2b73a1039ba30b5b2903337b91031b43ab93760511b608482015260a4016108c6565b613a258885614601565b6001600160601b0316816001600160601b0316106122565760405162461bcd60e51b815260206004820152605c6024820152600080516020615f5083398151915260448201527f726e3a2063616e6e6f74206b69636b206f70657261746f722077697468206d6f60648201527f7265207468616e206b69636b424950734f66546f74616c5374616b6500000000608482015260a4016108c6565b600081815260986020526040812054815b81811015613b51576001613ae48284615a29565b613aee9190615a29565b92508463ffffffff16609860008681526020019081526020016000208463ffffffff1681548110613b2157613b21615768565b60009182526020909120015463ffffffff1611613b3f575050610f19565b80613b4981615794565b915050613ad0565b5060405162461bcd60e51b815260206004820152606c60248201527f5265676973747279436f6f7264696e61746f722e67657451756f72756d42697460448201527f6d6170496e6465784174426c6f636b4e756d6265723a206e6f206269746d617060648201527f2075706461746520666f756e6420666f72206f70657261746f7249642061742060848201526b313637b1b590373ab6b132b960a11b60a482015260c4016108c6565b60965460ff1660c08110613c6f5760405162461bcd60e51b815260206004820152603560248201527f5265676973747279436f6f7264696e61746f722e63726561746551756f72756d6044820152740e881b585e081c5d5bdc9d5b5cc81c995858da1959605a1b60648201526084016108c6565b613c7a816001615df0565b6096805460ff191660ff9290921691909117905580613c9981866129fa565b60405160016296b58960e01b031981526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ff694a7790613cec90849088908890600401615e15565b600060405180830381600087803b158015613d0657600080fd5b505af1158015613d1a573d6000803e3d6000fd5b505060405163136ca0f960e11b815260ff841660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031692506326d941f29150602401600060405180830381600087803b158015613d8257600080fd5b505af1158015613d96573d6000803e3d6000fd5b505060405163136ca0f960e11b815260ff841660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031692506326d941f29150602401600060405180830381600087803b158015613dfe57600080fd5b505af1158015612256573d6000803e3d6000fd5b6000546201000090046001600160a01b0316158015613e3957506001600160a01b03821615155b613ebb5760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a4016108c6565b600181905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2613efe826126b4565b5050565b6000805b8215610f1957613f17600184615a29565b9092169180613f2581615e8e565b915050613f06565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148015613f8657507f000000000000000000000000000000000000000000000000000000000000000046145b15613fb057507f000000000000000000000000000000000000000000000000000000000000000090565b50604080517f00000000000000000000000000000000000000000000000000000000000000006020808301919091527f0000000000000000000000000000000000000000000000000000000000000000828401527f000000000000000000000000000000000000000000000000000000000000000060608301524660808301523060a0808401919091528351808403909101815260c0909201909252805191012090565b60008080600080516020615f708339815191526003600080516020615f7083398151915286600080516020615f708339815191528889090908905060006140ca827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52600080516020615f7083398151915261461b565b91959194509092505050565b60006101008251111561415f5760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016108c6565b815161416d57506000919050565b6000808360008151811061418357614183615768565b0160200151600160f89190911c81901b92505b845181101561425a578481815181106141b1576141b1615768565b0160200151600160f89190911c1b91508282116142465760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016108c6565b9181179161425381615794565b9050614196565b50909392505050565b60008281526098602052604090205480614308576000838152609860209081526040808320815160608101835263ffffffff43811682528185018681526001600160c01b03808a16958401958652845460018101865594885295909620915191909201805495519351909416600160401b026001600160401b03938316600160201b0267ffffffffffffffff1990961691909216179390931716919091179055505050565b6000838152609860205260408120614321600184615a29565b8154811061433157614331615768565b600091825260209091200180549091504363ffffffff908116911614156143755780546001600160401b0316600160401b6001600160c01b038516021781556109ad565b805463ffffffff438116600160201b81810267ffffffff0000000019909416939093178455600087815260986020908152604080832081516060810183529485528483018481526001600160c01b03808c1693870193845282546001810184559286529390942094519401805493519151909216600160401b026001600160401b0391861690960267ffffffffffffffff199093169390941692909217179190911691909117905550505050565b6001600160a01b0383163b1561453d57604051630b135d3f60e11b808252906001600160a01b03851690631626ba7e906144639086908690600401615ae7565b602060405180830381865afa158015614480573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144a49190615eb0565b6001600160e01b0319161461172f5760405162461bcd60e51b815260206004820152605360248201527f454950313237315369676e61747572655574696c732e636865636b5369676e6160448201527f747572655f454950313237313a2045524331323731207369676e6174757265206064820152721d995c9a599a58d85d1a5bdb8819985a5b1959606a1b608482015260a4016108c6565b826001600160a01b031661455183836146ca565b6001600160a01b03161461172f5760405162461bcd60e51b815260206004820152604760248201527f454950313237315369676e61747572655574696c732e636865636b5369676e6160448201527f747572655f454950313237313a207369676e6174757265206e6f742066726f6d6064820152661039b4b3b732b960c91b608482015260a4016108c6565b6020810151600090612710906145f79061ffff1685615eda565b6129f39190615f09565b6040810151600090612710906145f79061ffff1685615eda565b600080614626614a37565b61462e614a55565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa925082801561466f57614671565bfe5b50826146bf5760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c75726500000000000060448201526064016108c6565b505195945050505050565b60008060006146d985856146e6565b91509150611edb81614756565b60008082516041141561471d5760208301516040840151606085015160001a61471187828585614911565b9450945050505061474f565b825160401415614747576020830151604084015161473c8683836149fe565b93509350505061474f565b506000905060025b9250929050565b600081600481111561476a5761476a614de5565b14156147735750565b600181600481111561478757614787614de5565b14156147d55760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016108c6565b60028160048111156147e9576147e9614de5565b14156148375760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016108c6565b600381600481111561484b5761484b614de5565b14156148a45760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016108c6565b60048160048111156148b8576148b8614de5565b1415610d4e5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016108c6565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111561494857506000905060036149f5565b8460ff16601b1415801561496057508460ff16601c14155b1561497157506000905060046149f5565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa1580156149c5573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166149ee576000600192509250506149f5565b9150600090505b94509492505050565b6000806001600160ff1b03831681614a1b60ff86901c601b615946565b9050614a2987828885614911565b935093505050935093915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60008083601f840112614a8557600080fd5b5081356001600160401b03811115614a9c57600080fd5b6020830191508360208260051b850101111561474f57600080fd5b60008060208385031215614aca57600080fd5b82356001600160401b03811115614ae057600080fd5b614aec85828601614a73565b90969095509350505050565b600060208284031215614b0a57600080fd5b5035919050565b63ffffffff81168114610d4e57600080fd5b600080600060608486031215614b3857600080fd5b833592506020840135614b4a81614b11565b929592945050506040919091013590565b634e487b7160e01b600052604160045260246000fd5b604051606081016001600160401b0381118282101715614b9357614b93614b5b565b60405290565b604080519081016001600160401b0381118282101715614b9357614b93614b5b565b604051601f8201601f191681016001600160401b0381118282101715614be357614be3614b5b565b604052919050565b60006001600160401b03831115614c0457614c04614b5b565b614c17601f8401601f1916602001614bbb565b9050828152838383011115614c2b57600080fd5b828260208301376000602084830101529392505050565b600060208284031215614c5457600080fd5b81356001600160401b03811115614c6a57600080fd5b8201601f81018413614c7b57600080fd5b614c8a84823560208401614beb565b949350505050565b6001600160a01b0381168114610d4e57600080fd5b8035614cb281614c92565b919050565b600060208284031215614cc957600080fd5b81356129f381614c92565b60008060408385031215614ce757600080fd5b50508035926020909101359150565b803560ff81168114614cb257600080fd5b600060208284031215614d1957600080fd5b6129f382614cf6565b815181526020808301519082015260408101610f19565b60008083601f840112614d4b57600080fd5b5081356001600160401b03811115614d6257600080fd5b60208301915083602082850101111561474f57600080fd5b60008060008060408587031215614d9057600080fd5b84356001600160401b0380821115614da757600080fd5b614db388838901614a73565b90965094506020870135915080821115614dcc57600080fd5b50614dd987828801614d39565b95989497509550505050565b634e487b7160e01b600052602160045260246000fd5b60038110614e1957634e487b7160e01b600052602160045260246000fd5b9052565b815181526020808301516040830191614e3890840182614dfb565b5092915050565b803561ffff81168114614cb257600080fd5b600060608284031215614e6357600080fd5b614e6b614b71565b90508135614e7881614b11565b8152614e8660208301614e3f565b6020820152614e9760408301614e3f565b604082015292915050565b60008060808385031215614eb557600080fd5b614ebe83614cf6565b9150614ecd8460208501614e51565b90509250929050565b600080600060408486031215614eeb57600080fd5b8335614ef681614c92565b925060208401356001600160401b03811115614f1157600080fd5b614f1d86828701614d39565b9497909650939450505050565b60006001600160401b03821115614f4357614f43614b5b565b5060051b60200190565b600060408284031215614f5f57600080fd5b614f67614b99565b9050614f7282614cf6565b81526020820135614f8281614c92565b602082015292915050565b600080600080600060a08688031215614fa557600080fd5b8535614fb081614c92565b945060208681013594506040808801356001600160401b03811115614fd457600080fd5b8801601f81018a13614fe557600080fd5b8035614ff8614ff382614f2a565b614bbb565b81815260069190911b8201840190848101908c83111561501757600080fd5b928501925b8284101561503d5761502e8d85614f4d565b8252928401929085019061501c565b999c989b5098996060810135995060800135979650505050505050565b6000610100828403121561506d57600080fd5b50919050565b60008083601f84011261508557600080fd5b5081356001600160401b0381111561509c57600080fd5b6020830191508360208260061b850101111561474f57600080fd5b6000606082840312156150c957600080fd5b6150d1614b71565b905081356001600160401b038111156150e957600080fd5b8201601f810184136150fa57600080fd5b61510984823560208401614beb565b825250602082013560208201526040820135604082015292915050565b60008060008060008060008060006101a08a8c03121561514557600080fd5b89356001600160401b038082111561515c57600080fd5b6151688d838e01614d39565b909b50995060208c013591508082111561518157600080fd5b61518d8d838e01614d39565b90995097508791506151a28d60408e0161505a565b96506101408c01359150808211156151b957600080fd5b6151c58d838e01615073565b90965094506101608c01359150808211156151df57600080fd5b6151eb8d838e016150b7565b93506101808c013591508082111561520257600080fd5b5061520f8c828d016150b7565b9150509295985092959850929598565b600080600080600080610160878903121561523957600080fd5b86356001600160401b038082111561525057600080fd5b61525c8a838b01614d39565b9098509650602089013591508082111561527557600080fd5b6152818a838b01614d39565b90965094508491506152968a60408b0161505a565b93506101408901359150808211156152ad57600080fd5b506152ba89828a016150b7565b9150509295509295509295565b600080604083850312156152da57600080fd5b82356152e581614b11565b91506020838101356001600160401b0381111561530157600080fd5b8401601f8101861361531257600080fd5b8035615320614ff382614f2a565b81815260059190911b8201830190838101908883111561533f57600080fd5b928401925b8284101561535d57833582529284019290840190615344565b80955050505050509250929050565b6020808252825182820181905260009190848201906040850190845b818110156153aa57835163ffffffff1683529284019291840191600101615388565b50909695505050505050565b600080602083850312156153c957600080fd5b82356001600160401b038111156153df57600080fd5b614aec85828601614d39565b6001600160601b0381168114610d4e57600080fd5b600082601f83011261541157600080fd5b81356020615421614ff383614f2a565b82815260069290921b8401810191818101908684111561544057600080fd5b8286015b84811015615491576040818903121561545d5760008081fd5b615465614b99565b813561547081614c92565b81528185013561547f816153eb565b81860152835291830191604001615444565b509695505050505050565b600080600060a084860312156154b157600080fd5b6154bb8585614e51565b925060608401356154cb816153eb565b915060808401356001600160401b038111156154e657600080fd5b6154f286828701615400565b9150509250925092565b600082601f83011261550d57600080fd5b8135602061551d614ff383614f2a565b8281526060928302850182019282820191908785111561553c57600080fd5b8387015b8581101561555f576155528982614e51565b8452928401928101615540565b5090979650505050505050565b600082601f83011261557d57600080fd5b8135602061558d614ff383614f2a565b82815260059290921b840181019181810190868411156155ac57600080fd5b8286015b848110156154915780356155c3816153eb565b83529183019183016155b0565b600082601f8301126155e157600080fd5b813560206155f1614ff383614f2a565b82815260059290921b8401810191818101908684111561561057600080fd5b8286015b848110156154915780356001600160401b038111156156335760008081fd5b6156418986838b0101615400565b845250918301918301615614565b600080600080600080600080610100898b03121561566c57600080fd5b61567589614ca7565b975061568360208a01614ca7565b965061569160408a01614ca7565b955061569f60608a01614ca7565b94506080890135935060a08901356001600160401b03808211156156c257600080fd5b6156ce8c838d016154fc565b945060c08b01359150808211156156e457600080fd5b6156f08c838d0161556c565b935060e08b013591508082111561570657600080fd5b506157138b828c016155d0565b9150509295985092959890939650565b60208101610f198284614dfb565b60208082526019908201527f5061757361626c653a20696e6465782069732070617573656400000000000000604082015260600190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006000198214156157a8576157a861577e565b5060010190565b6000815180845260005b818110156157d5576020818501810151868301820152016157b9565b818111156157e7576000602083870101525b50601f01601f19169290920160200192915050565b6020815260006129f360208301846157af565b60006020828403121561582157600080fd5b81516129f381614c92565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b60006020828403121561588857600080fd5b815180151581146129f357600080fd5b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b6000808335601e198436030181126158f757600080fd5b8301803591506001600160401b0382111561591157600080fd5b6020019150600581901b360382131561474f57600080fd5b60006020828403121561593b57600080fd5b81516129f381614b11565b600082198211156159595761595961577e565b500190565b6000808585111561596e57600080fd5b8386111561597b57600080fd5b5050820193919092039150565b600060c08201888352602060018060a01b03808a16828601526040898187015260c0606087015283895180865260e088019150848b01955060005b818110156159ed578651805160ff16845286015185168684015295850195918301916001016159c3565b505060808701989098525050505060a09091019190915250949350505050565b600060408284031215615a1f57600080fd5b6129f38383614f4d565b600082821015615a3b57615a3b61577e565b500390565b60018060a01b0384168152826020820152606060408201526000615a6760608301846157af565b95945050505050565b600060208284031215615a8257600080fd5b81516001600160c01b03811681146129f357600080fd5b634e487b7160e01b600052601260045260246000fd5b600082615abe57615abe615a99565b500690565b6001600160a01b0383168152604060208201819052600090614c8a908301846157af565b828152604060208201526000614c8a60408301846157af565b600060208284031215615b1257600080fd5b5051919050565b6001600160a01b03841681526101608101615b41602083018580358252602090810135910152565b615b5b606083016040860180358252602090810135910152565b60406080850160a084013760e0820160008152604060c0860182375060006101208301908152835190526020909201516101409091015292915050565b60018060a01b0383168152604060208201526000825160606040840152615bc260a08401826157af565b90506020840151606084015260408401516080840152809150509392505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0384168152604060208201819052600090615a679083018486615be3565b60018060a01b0385168152836020820152606060408201526000611940606083018486615be3565b600082601f830112615c6a57600080fd5b81516020615c7a614ff383614f2a565b82815260059290921b84018101918181019086841115615c9957600080fd5b8286015b84811015615491578051615cb0816153eb565b8352918301918301615c9d565b60008060408385031215615cd057600080fd5b82516001600160401b0380821115615ce757600080fd5b615cf386838701615c59565b93506020850151915080821115615d0957600080fd5b50615d1685828601615c59565b9150509250929050565b838152604060208201526000615a67604083018486615be3565b60006020808385031215615d4d57600080fd5b82516001600160401b03811115615d6357600080fd5b8301601f81018513615d7457600080fd5b8051615d82614ff382614f2a565b81815260059190911b82018301908381019087831115615da157600080fd5b928401925b82841015615dc8578351615db981614b11565b82529284019290840190615da6565b979650505050505050565b600060208284031215615de557600080fd5b81516129f3816153eb565b600060ff821660ff84168060ff03821115615e0d57615e0d61577e565b019392505050565b60006060820160ff8616835260206001600160601b03808716828601526040606081870152838751808652608088019150848901955060005b81811015615e7e57865180516001600160a01b031684528601518516868401529585019591830191600101615e4e565b50909a9950505050505050505050565b600061ffff80831681811415615ea657615ea661577e565b6001019392505050565b600060208284031215615ec257600080fd5b81516001600160e01b0319811681146129f357600080fd5b60006001600160601b0380831681851681830481118215151615615f0057615f0061577e565b02949350505050565b60006001600160601b0380841680615f2357615f23615a99565b9216919091049291505056fe5265676973747279436f6f7264696e61746f722e7570646174654f70657261745265676973747279436f6f7264696e61746f722e5f76616c696461746543687530644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47a26469706673582212200dbf2862cacc1a9e8f8da4619494f53fcd304f82a034cf9e948fa181a2d6e8a064736f6c634300080c0033",
}

// ContractRegistryCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractRegistryCoordinatorMetaData.ABI instead.
var ContractRegistryCoordinatorABI = ContractRegistryCoordinatorMetaData.ABI

// ContractRegistryCoordinatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractRegistryCoordinatorMetaData.Bin instead.
var ContractRegistryCoordinatorBin = ContractRegistryCoordinatorMetaData.Bin

// DeployContractRegistryCoordinator deploys a new Ethereum contract, binding an instance of ContractRegistryCoordinator to it.
func DeployContractRegistryCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, _serviceManager common.Address, _stakeRegistry common.Address, _blsApkRegistry common.Address, _indexRegistry common.Address) (common.Address, *types.Transaction, *ContractRegistryCoordinator, error) {
	parsed, err := ContractRegistryCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractRegistryCoordinatorBin), backend, _serviceManager, _stakeRegistry, _blsApkRegistry, _indexRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractRegistryCoordinator{ContractRegistryCoordinatorCaller: ContractRegistryCoordinatorCaller{contract: contract}, ContractRegistryCoordinatorTransactor: ContractRegistryCoordinatorTransactor{contract: contract}, ContractRegistryCoordinatorFilterer: ContractRegistryCoordinatorFilterer{contract: contract}}, nil
}

// ContractRegistryCoordinator is an auto generated Go binding around an Ethereum contract.
type ContractRegistryCoordinator struct {
	ContractRegistryCoordinatorCaller     // Read-only binding to the contract
	ContractRegistryCoordinatorTransactor // Write-only binding to the contract
	ContractRegistryCoordinatorFilterer   // Log filterer for contract events
}

// ContractRegistryCoordinatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractRegistryCoordinatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractRegistryCoordinatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractRegistryCoordinatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractRegistryCoordinatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractRegistryCoordinatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractRegistryCoordinatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractRegistryCoordinatorSession struct {
	Contract     *ContractRegistryCoordinator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ContractRegistryCoordinatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractRegistryCoordinatorCallerSession struct {
	Contract *ContractRegistryCoordinatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// ContractRegistryCoordinatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractRegistryCoordinatorTransactorSession struct {
	Contract     *ContractRegistryCoordinatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// ContractRegistryCoordinatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRegistryCoordinatorRaw struct {
	Contract *ContractRegistryCoordinator // Generic contract binding to access the raw methods on
}

// ContractRegistryCoordinatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractRegistryCoordinatorCallerRaw struct {
	Contract *ContractRegistryCoordinatorCaller // Generic read-only contract binding to access the raw methods on
}

// ContractRegistryCoordinatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractRegistryCoordinatorTransactorRaw struct {
	Contract *ContractRegistryCoordinatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractRegistryCoordinator creates a new instance of ContractRegistryCoordinator, bound to a specific deployed contract.
func NewContractRegistryCoordinator(address common.Address, backend bind.ContractBackend) (*ContractRegistryCoordinator, error) {
	contract, err := bindContractRegistryCoordinator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinator{ContractRegistryCoordinatorCaller: ContractRegistryCoordinatorCaller{contract: contract}, ContractRegistryCoordinatorTransactor: ContractRegistryCoordinatorTransactor{contract: contract}, ContractRegistryCoordinatorFilterer: ContractRegistryCoordinatorFilterer{contract: contract}}, nil
}

// NewContractRegistryCoordinatorCaller creates a new read-only instance of ContractRegistryCoordinator, bound to a specific deployed contract.
func NewContractRegistryCoordinatorCaller(address common.Address, caller bind.ContractCaller) (*ContractRegistryCoordinatorCaller, error) {
	contract, err := bindContractRegistryCoordinator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorCaller{contract: contract}, nil
}

// NewContractRegistryCoordinatorTransactor creates a new write-only instance of ContractRegistryCoordinator, bound to a specific deployed contract.
func NewContractRegistryCoordinatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractRegistryCoordinatorTransactor, error) {
	contract, err := bindContractRegistryCoordinator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorTransactor{contract: contract}, nil
}

// NewContractRegistryCoordinatorFilterer creates a new log filterer instance of ContractRegistryCoordinator, bound to a specific deployed contract.
func NewContractRegistryCoordinatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractRegistryCoordinatorFilterer, error) {
	contract, err := bindContractRegistryCoordinator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorFilterer{contract: contract}, nil
}

// bindContractRegistryCoordinator binds a generic wrapper to an already deployed contract.
func bindContractRegistryCoordinator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractRegistryCoordinatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractRegistryCoordinator.Contract.ContractRegistryCoordinatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.ContractRegistryCoordinatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.ContractRegistryCoordinatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractRegistryCoordinator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.contract.Transact(opts, method, params...)
}

// OPERATORCHURNAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0xca0de882.
//
// Solidity: function OPERATOR_CHURN_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) OPERATORCHURNAPPROVALTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "OPERATOR_CHURN_APPROVAL_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OPERATORCHURNAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0xca0de882.
//
// Solidity: function OPERATOR_CHURN_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) OPERATORCHURNAPPROVALTYPEHASH() ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.OPERATORCHURNAPPROVALTYPEHASH(&_ContractRegistryCoordinator.CallOpts)
}

// OPERATORCHURNAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0xca0de882.
//
// Solidity: function OPERATOR_CHURN_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) OPERATORCHURNAPPROVALTYPEHASH() ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.OPERATORCHURNAPPROVALTYPEHASH(&_ContractRegistryCoordinator.CallOpts)
}

// PUBKEYREGISTRATIONTYPEHASH is a free data retrieval call binding the contract method 0x9feab859.
//
// Solidity: function PUBKEY_REGISTRATION_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) PUBKEYREGISTRATIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "PUBKEY_REGISTRATION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PUBKEYREGISTRATIONTYPEHASH is a free data retrieval call binding the contract method 0x9feab859.
//
// Solidity: function PUBKEY_REGISTRATION_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) PUBKEYREGISTRATIONTYPEHASH() ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.PUBKEYREGISTRATIONTYPEHASH(&_ContractRegistryCoordinator.CallOpts)
}

// PUBKEYREGISTRATIONTYPEHASH is a free data retrieval call binding the contract method 0x9feab859.
//
// Solidity: function PUBKEY_REGISTRATION_TYPEHASH() view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) PUBKEYREGISTRATIONTYPEHASH() ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.PUBKEYREGISTRATIONTYPEHASH(&_ContractRegistryCoordinator.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) BlsApkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "blsApkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) BlsApkRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.BlsApkRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.BlsApkRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// CalculateOperatorChurnApprovalDigestHash is a free data retrieval call binding the contract method 0x84ca5213.
//
// Solidity: function calculateOperatorChurnApprovalDigestHash(address registeringOperator, bytes32 registeringOperatorId, (uint8,address)[] operatorKickParams, bytes32 salt, uint256 expiry) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) CalculateOperatorChurnApprovalDigestHash(opts *bind.CallOpts, registeringOperator common.Address, registeringOperatorId [32]byte, operatorKickParams []IRegistryCoordinatorOperatorKickParam, salt [32]byte, expiry *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "calculateOperatorChurnApprovalDigestHash", registeringOperator, registeringOperatorId, operatorKickParams, salt, expiry)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateOperatorChurnApprovalDigestHash is a free data retrieval call binding the contract method 0x84ca5213.
//
// Solidity: function calculateOperatorChurnApprovalDigestHash(address registeringOperator, bytes32 registeringOperatorId, (uint8,address)[] operatorKickParams, bytes32 salt, uint256 expiry) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) CalculateOperatorChurnApprovalDigestHash(registeringOperator common.Address, registeringOperatorId [32]byte, operatorKickParams []IRegistryCoordinatorOperatorKickParam, salt [32]byte, expiry *big.Int) ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.CalculateOperatorChurnApprovalDigestHash(&_ContractRegistryCoordinator.CallOpts, registeringOperator, registeringOperatorId, operatorKickParams, salt, expiry)
}

// CalculateOperatorChurnApprovalDigestHash is a free data retrieval call binding the contract method 0x84ca5213.
//
// Solidity: function calculateOperatorChurnApprovalDigestHash(address registeringOperator, bytes32 registeringOperatorId, (uint8,address)[] operatorKickParams, bytes32 salt, uint256 expiry) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) CalculateOperatorChurnApprovalDigestHash(registeringOperator common.Address, registeringOperatorId [32]byte, operatorKickParams []IRegistryCoordinatorOperatorKickParam, salt [32]byte, expiry *big.Int) ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.CalculateOperatorChurnApprovalDigestHash(&_ContractRegistryCoordinator.CallOpts, registeringOperator, registeringOperatorId, operatorKickParams, salt, expiry)
}

// ChurnApprover is a free data retrieval call binding the contract method 0x054310e6.
//
// Solidity: function churnApprover() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) ChurnApprover(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "churnApprover")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ChurnApprover is a free data retrieval call binding the contract method 0x054310e6.
//
// Solidity: function churnApprover() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) ChurnApprover() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.ChurnApprover(&_ContractRegistryCoordinator.CallOpts)
}

// ChurnApprover is a free data retrieval call binding the contract method 0x054310e6.
//
// Solidity: function churnApprover() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) ChurnApprover() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.ChurnApprover(&_ContractRegistryCoordinator.CallOpts)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) EjectionCooldown(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "ejectionCooldown")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) EjectionCooldown() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.EjectionCooldown(&_ContractRegistryCoordinator.CallOpts)
}

// EjectionCooldown is a free data retrieval call binding the contract method 0xa96f783e.
//
// Solidity: function ejectionCooldown() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) EjectionCooldown() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.EjectionCooldown(&_ContractRegistryCoordinator.CallOpts)
}

// Ejector is a free data retrieval call binding the contract method 0x28f61b31.
//
// Solidity: function ejector() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) Ejector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "ejector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Ejector is a free data retrieval call binding the contract method 0x28f61b31.
//
// Solidity: function ejector() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Ejector() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Ejector(&_ContractRegistryCoordinator.CallOpts)
}

// Ejector is a free data retrieval call binding the contract method 0x28f61b31.
//
// Solidity: function ejector() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) Ejector() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Ejector(&_ContractRegistryCoordinator.CallOpts)
}

// GetCurrentQuorumBitmap is a free data retrieval call binding the contract method 0x871ef049.
//
// Solidity: function getCurrentQuorumBitmap(bytes32 operatorId) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetCurrentQuorumBitmap(opts *bind.CallOpts, operatorId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getCurrentQuorumBitmap", operatorId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentQuorumBitmap is a free data retrieval call binding the contract method 0x871ef049.
//
// Solidity: function getCurrentQuorumBitmap(bytes32 operatorId) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetCurrentQuorumBitmap(operatorId [32]byte) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetCurrentQuorumBitmap(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetCurrentQuorumBitmap is a free data retrieval call binding the contract method 0x871ef049.
//
// Solidity: function getCurrentQuorumBitmap(bytes32 operatorId) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetCurrentQuorumBitmap(operatorId [32]byte) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetCurrentQuorumBitmap(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetOperator is a free data retrieval call binding the contract method 0x5865c60c.
//
// Solidity: function getOperator(address operator) view returns((bytes32,uint8))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetOperator(opts *bind.CallOpts, operator common.Address) (IRegistryCoordinatorOperatorInfo, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getOperator", operator)

	if err != nil {
		return *new(IRegistryCoordinatorOperatorInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IRegistryCoordinatorOperatorInfo)).(*IRegistryCoordinatorOperatorInfo)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0x5865c60c.
//
// Solidity: function getOperator(address operator) view returns((bytes32,uint8))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetOperator(operator common.Address) (IRegistryCoordinatorOperatorInfo, error) {
	return _ContractRegistryCoordinator.Contract.GetOperator(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetOperator is a free data retrieval call binding the contract method 0x5865c60c.
//
// Solidity: function getOperator(address operator) view returns((bytes32,uint8))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetOperator(operator common.Address) (IRegistryCoordinatorOperatorInfo, error) {
	return _ContractRegistryCoordinator.Contract.GetOperator(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetOperatorFromId is a free data retrieval call binding the contract method 0x296bb064.
//
// Solidity: function getOperatorFromId(bytes32 operatorId) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetOperatorFromId(opts *bind.CallOpts, operatorId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getOperatorFromId", operatorId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperatorFromId is a free data retrieval call binding the contract method 0x296bb064.
//
// Solidity: function getOperatorFromId(bytes32 operatorId) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetOperatorFromId(operatorId [32]byte) (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorFromId(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetOperatorFromId is a free data retrieval call binding the contract method 0x296bb064.
//
// Solidity: function getOperatorFromId(bytes32 operatorId) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetOperatorFromId(operatorId [32]byte) (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorFromId(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetOperatorId(opts *bind.CallOpts, operator common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getOperatorId", operator)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetOperatorId(operator common.Address) ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorId(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetOperatorId is a free data retrieval call binding the contract method 0x13542a4e.
//
// Solidity: function getOperatorId(address operator) view returns(bytes32)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetOperatorId(operator common.Address) ([32]byte, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorId(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetOperatorSetParams is a free data retrieval call binding the contract method 0xe65797ad.
//
// Solidity: function getOperatorSetParams(uint8 quorumNumber) view returns((uint32,uint16,uint16))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetOperatorSetParams(opts *bind.CallOpts, quorumNumber uint8) (IRegistryCoordinatorOperatorSetParam, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getOperatorSetParams", quorumNumber)

	if err != nil {
		return *new(IRegistryCoordinatorOperatorSetParam), err
	}

	out0 := *abi.ConvertType(out[0], new(IRegistryCoordinatorOperatorSetParam)).(*IRegistryCoordinatorOperatorSetParam)

	return out0, err

}

// GetOperatorSetParams is a free data retrieval call binding the contract method 0xe65797ad.
//
// Solidity: function getOperatorSetParams(uint8 quorumNumber) view returns((uint32,uint16,uint16))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetOperatorSetParams(quorumNumber uint8) (IRegistryCoordinatorOperatorSetParam, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorSetParams(&_ContractRegistryCoordinator.CallOpts, quorumNumber)
}

// GetOperatorSetParams is a free data retrieval call binding the contract method 0xe65797ad.
//
// Solidity: function getOperatorSetParams(uint8 quorumNumber) view returns((uint32,uint16,uint16))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetOperatorSetParams(quorumNumber uint8) (IRegistryCoordinatorOperatorSetParam, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorSetParams(&_ContractRegistryCoordinator.CallOpts, quorumNumber)
}

// GetOperatorStatus is a free data retrieval call binding the contract method 0xfd39105a.
//
// Solidity: function getOperatorStatus(address operator) view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetOperatorStatus(opts *bind.CallOpts, operator common.Address) (uint8, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getOperatorStatus", operator)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetOperatorStatus is a free data retrieval call binding the contract method 0xfd39105a.
//
// Solidity: function getOperatorStatus(address operator) view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetOperatorStatus(operator common.Address) (uint8, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorStatus(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetOperatorStatus is a free data retrieval call binding the contract method 0xfd39105a.
//
// Solidity: function getOperatorStatus(address operator) view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetOperatorStatus(operator common.Address) (uint8, error) {
	return _ContractRegistryCoordinator.Contract.GetOperatorStatus(&_ContractRegistryCoordinator.CallOpts, operator)
}

// GetQuorumBitmapAtBlockNumberByIndex is a free data retrieval call binding the contract method 0x04ec6351.
//
// Solidity: function getQuorumBitmapAtBlockNumberByIndex(bytes32 operatorId, uint32 blockNumber, uint256 index) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetQuorumBitmapAtBlockNumberByIndex(opts *bind.CallOpts, operatorId [32]byte, blockNumber uint32, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getQuorumBitmapAtBlockNumberByIndex", operatorId, blockNumber, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetQuorumBitmapAtBlockNumberByIndex is a free data retrieval call binding the contract method 0x04ec6351.
//
// Solidity: function getQuorumBitmapAtBlockNumberByIndex(bytes32 operatorId, uint32 blockNumber, uint256 index) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetQuorumBitmapAtBlockNumberByIndex(operatorId [32]byte, blockNumber uint32, index *big.Int) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapAtBlockNumberByIndex(&_ContractRegistryCoordinator.CallOpts, operatorId, blockNumber, index)
}

// GetQuorumBitmapAtBlockNumberByIndex is a free data retrieval call binding the contract method 0x04ec6351.
//
// Solidity: function getQuorumBitmapAtBlockNumberByIndex(bytes32 operatorId, uint32 blockNumber, uint256 index) view returns(uint192)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetQuorumBitmapAtBlockNumberByIndex(operatorId [32]byte, blockNumber uint32, index *big.Int) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapAtBlockNumberByIndex(&_ContractRegistryCoordinator.CallOpts, operatorId, blockNumber, index)
}

// GetQuorumBitmapHistoryLength is a free data retrieval call binding the contract method 0x03fd3492.
//
// Solidity: function getQuorumBitmapHistoryLength(bytes32 operatorId) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetQuorumBitmapHistoryLength(opts *bind.CallOpts, operatorId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getQuorumBitmapHistoryLength", operatorId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetQuorumBitmapHistoryLength is a free data retrieval call binding the contract method 0x03fd3492.
//
// Solidity: function getQuorumBitmapHistoryLength(bytes32 operatorId) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetQuorumBitmapHistoryLength(operatorId [32]byte) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapHistoryLength(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetQuorumBitmapHistoryLength is a free data retrieval call binding the contract method 0x03fd3492.
//
// Solidity: function getQuorumBitmapHistoryLength(bytes32 operatorId) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetQuorumBitmapHistoryLength(operatorId [32]byte) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapHistoryLength(&_ContractRegistryCoordinator.CallOpts, operatorId)
}

// GetQuorumBitmapIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xc391425e.
//
// Solidity: function getQuorumBitmapIndicesAtBlockNumber(uint32 blockNumber, bytes32[] operatorIds) view returns(uint32[])
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetQuorumBitmapIndicesAtBlockNumber(opts *bind.CallOpts, blockNumber uint32, operatorIds [][32]byte) ([]uint32, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getQuorumBitmapIndicesAtBlockNumber", blockNumber, operatorIds)

	if err != nil {
		return *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint32)).(*[]uint32)

	return out0, err

}

// GetQuorumBitmapIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xc391425e.
//
// Solidity: function getQuorumBitmapIndicesAtBlockNumber(uint32 blockNumber, bytes32[] operatorIds) view returns(uint32[])
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetQuorumBitmapIndicesAtBlockNumber(blockNumber uint32, operatorIds [][32]byte) ([]uint32, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapIndicesAtBlockNumber(&_ContractRegistryCoordinator.CallOpts, blockNumber, operatorIds)
}

// GetQuorumBitmapIndicesAtBlockNumber is a free data retrieval call binding the contract method 0xc391425e.
//
// Solidity: function getQuorumBitmapIndicesAtBlockNumber(uint32 blockNumber, bytes32[] operatorIds) view returns(uint32[])
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetQuorumBitmapIndicesAtBlockNumber(blockNumber uint32, operatorIds [][32]byte) ([]uint32, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapIndicesAtBlockNumber(&_ContractRegistryCoordinator.CallOpts, blockNumber, operatorIds)
}

// GetQuorumBitmapUpdateByIndex is a free data retrieval call binding the contract method 0x1eb812da.
//
// Solidity: function getQuorumBitmapUpdateByIndex(bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint192))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) GetQuorumBitmapUpdateByIndex(opts *bind.CallOpts, operatorId [32]byte, index *big.Int) (IRegistryCoordinatorQuorumBitmapUpdate, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "getQuorumBitmapUpdateByIndex", operatorId, index)

	if err != nil {
		return *new(IRegistryCoordinatorQuorumBitmapUpdate), err
	}

	out0 := *abi.ConvertType(out[0], new(IRegistryCoordinatorQuorumBitmapUpdate)).(*IRegistryCoordinatorQuorumBitmapUpdate)

	return out0, err

}

// GetQuorumBitmapUpdateByIndex is a free data retrieval call binding the contract method 0x1eb812da.
//
// Solidity: function getQuorumBitmapUpdateByIndex(bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint192))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) GetQuorumBitmapUpdateByIndex(operatorId [32]byte, index *big.Int) (IRegistryCoordinatorQuorumBitmapUpdate, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapUpdateByIndex(&_ContractRegistryCoordinator.CallOpts, operatorId, index)
}

// GetQuorumBitmapUpdateByIndex is a free data retrieval call binding the contract method 0x1eb812da.
//
// Solidity: function getQuorumBitmapUpdateByIndex(bytes32 operatorId, uint256 index) view returns((uint32,uint32,uint192))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) GetQuorumBitmapUpdateByIndex(operatorId [32]byte, index *big.Int) (IRegistryCoordinatorQuorumBitmapUpdate, error) {
	return _ContractRegistryCoordinator.Contract.GetQuorumBitmapUpdateByIndex(&_ContractRegistryCoordinator.CallOpts, operatorId, index)
}

// IndexRegistry is a free data retrieval call binding the contract method 0x9e9923c2.
//
// Solidity: function indexRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) IndexRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "indexRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IndexRegistry is a free data retrieval call binding the contract method 0x9e9923c2.
//
// Solidity: function indexRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) IndexRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.IndexRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// IndexRegistry is a free data retrieval call binding the contract method 0x9e9923c2.
//
// Solidity: function indexRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) IndexRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.IndexRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// IsChurnApproverSaltUsed is a free data retrieval call binding the contract method 0x1478851f.
//
// Solidity: function isChurnApproverSaltUsed(bytes32 ) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) IsChurnApproverSaltUsed(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "isChurnApproverSaltUsed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsChurnApproverSaltUsed is a free data retrieval call binding the contract method 0x1478851f.
//
// Solidity: function isChurnApproverSaltUsed(bytes32 ) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) IsChurnApproverSaltUsed(arg0 [32]byte) (bool, error) {
	return _ContractRegistryCoordinator.Contract.IsChurnApproverSaltUsed(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// IsChurnApproverSaltUsed is a free data retrieval call binding the contract method 0x1478851f.
//
// Solidity: function isChurnApproverSaltUsed(bytes32 ) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) IsChurnApproverSaltUsed(arg0 [32]byte) (bool, error) {
	return _ContractRegistryCoordinator.Contract.IsChurnApproverSaltUsed(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// LastEjectionTimestamp is a free data retrieval call binding the contract method 0x125e0584.
//
// Solidity: function lastEjectionTimestamp(address ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) LastEjectionTimestamp(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "lastEjectionTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastEjectionTimestamp is a free data retrieval call binding the contract method 0x125e0584.
//
// Solidity: function lastEjectionTimestamp(address ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) LastEjectionTimestamp(arg0 common.Address) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.LastEjectionTimestamp(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// LastEjectionTimestamp is a free data retrieval call binding the contract method 0x125e0584.
//
// Solidity: function lastEjectionTimestamp(address ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) LastEjectionTimestamp(arg0 common.Address) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.LastEjectionTimestamp(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// NumRegistries is a free data retrieval call binding the contract method 0xd72d8dd6.
//
// Solidity: function numRegistries() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) NumRegistries(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "numRegistries")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumRegistries is a free data retrieval call binding the contract method 0xd72d8dd6.
//
// Solidity: function numRegistries() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) NumRegistries() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.NumRegistries(&_ContractRegistryCoordinator.CallOpts)
}

// NumRegistries is a free data retrieval call binding the contract method 0xd72d8dd6.
//
// Solidity: function numRegistries() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) NumRegistries() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.NumRegistries(&_ContractRegistryCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Owner() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Owner(&_ContractRegistryCoordinator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) Owner() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Owner(&_ContractRegistryCoordinator.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Paused(index uint8) (bool, error) {
	return _ContractRegistryCoordinator.Contract.Paused(&_ContractRegistryCoordinator.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) Paused(index uint8) (bool, error) {
	return _ContractRegistryCoordinator.Contract.Paused(&_ContractRegistryCoordinator.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Paused0() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.Paused0(&_ContractRegistryCoordinator.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) Paused0() (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.Paused0(&_ContractRegistryCoordinator.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) PauserRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.PauserRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) PauserRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.PauserRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// PubkeyRegistrationMessageHash is a free data retrieval call binding the contract method 0x3c2a7f4c.
//
// Solidity: function pubkeyRegistrationMessageHash(address operator) view returns((uint256,uint256))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) PubkeyRegistrationMessageHash(opts *bind.CallOpts, operator common.Address) (BN254G1Point, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "pubkeyRegistrationMessageHash", operator)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// PubkeyRegistrationMessageHash is a free data retrieval call binding the contract method 0x3c2a7f4c.
//
// Solidity: function pubkeyRegistrationMessageHash(address operator) view returns((uint256,uint256))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) PubkeyRegistrationMessageHash(operator common.Address) (BN254G1Point, error) {
	return _ContractRegistryCoordinator.Contract.PubkeyRegistrationMessageHash(&_ContractRegistryCoordinator.CallOpts, operator)
}

// PubkeyRegistrationMessageHash is a free data retrieval call binding the contract method 0x3c2a7f4c.
//
// Solidity: function pubkeyRegistrationMessageHash(address operator) view returns((uint256,uint256))
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) PubkeyRegistrationMessageHash(operator common.Address) (BN254G1Point, error) {
	return _ContractRegistryCoordinator.Contract.PubkeyRegistrationMessageHash(&_ContractRegistryCoordinator.CallOpts, operator)
}

// QuorumCount is a free data retrieval call binding the contract method 0x9aa1653d.
//
// Solidity: function quorumCount() view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) QuorumCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "quorumCount")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// QuorumCount is a free data retrieval call binding the contract method 0x9aa1653d.
//
// Solidity: function quorumCount() view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) QuorumCount() (uint8, error) {
	return _ContractRegistryCoordinator.Contract.QuorumCount(&_ContractRegistryCoordinator.CallOpts)
}

// QuorumCount is a free data retrieval call binding the contract method 0x9aa1653d.
//
// Solidity: function quorumCount() view returns(uint8)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) QuorumCount() (uint8, error) {
	return _ContractRegistryCoordinator.Contract.QuorumCount(&_ContractRegistryCoordinator.CallOpts)
}

// QuorumUpdateBlockNumber is a free data retrieval call binding the contract method 0x249a0c42.
//
// Solidity: function quorumUpdateBlockNumber(uint8 ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) QuorumUpdateBlockNumber(opts *bind.CallOpts, arg0 uint8) (*big.Int, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "quorumUpdateBlockNumber", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumUpdateBlockNumber is a free data retrieval call binding the contract method 0x249a0c42.
//
// Solidity: function quorumUpdateBlockNumber(uint8 ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) QuorumUpdateBlockNumber(arg0 uint8) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.QuorumUpdateBlockNumber(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// QuorumUpdateBlockNumber is a free data retrieval call binding the contract method 0x249a0c42.
//
// Solidity: function quorumUpdateBlockNumber(uint8 ) view returns(uint256)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) QuorumUpdateBlockNumber(arg0 uint8) (*big.Int, error) {
	return _ContractRegistryCoordinator.Contract.QuorumUpdateBlockNumber(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) Registries(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "registries", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Registries(arg0 *big.Int) (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Registries(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) Registries(arg0 *big.Int) (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.Registries(&_ContractRegistryCoordinator.CallOpts, arg0)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) ServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "serviceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) ServiceManager() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.ServiceManager(&_ContractRegistryCoordinator.CallOpts)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) ServiceManager() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.ServiceManager(&_ContractRegistryCoordinator.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) StakeRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.StakeRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.StakeRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// CreateQuorum is a paid mutator transaction binding the contract method 0xd75b4c88.
//
// Solidity: function createQuorum((uint32,uint16,uint16) operatorSetParams, uint96 minimumStake, (address,uint96)[] strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) CreateQuorum(opts *bind.TransactOpts, operatorSetParams IRegistryCoordinatorOperatorSetParam, minimumStake *big.Int, strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "createQuorum", operatorSetParams, minimumStake, strategyParams)
}

// CreateQuorum is a paid mutator transaction binding the contract method 0xd75b4c88.
//
// Solidity: function createQuorum((uint32,uint16,uint16) operatorSetParams, uint96 minimumStake, (address,uint96)[] strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) CreateQuorum(operatorSetParams IRegistryCoordinatorOperatorSetParam, minimumStake *big.Int, strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.CreateQuorum(&_ContractRegistryCoordinator.TransactOpts, operatorSetParams, minimumStake, strategyParams)
}

// CreateQuorum is a paid mutator transaction binding the contract method 0xd75b4c88.
//
// Solidity: function createQuorum((uint32,uint16,uint16) operatorSetParams, uint96 minimumStake, (address,uint96)[] strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) CreateQuorum(operatorSetParams IRegistryCoordinatorOperatorSetParam, minimumStake *big.Int, strategyParams []IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.CreateQuorum(&_ContractRegistryCoordinator.TransactOpts, operatorSetParams, minimumStake, strategyParams)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xca4f2d97.
//
// Solidity: function deregisterOperator(bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) DeregisterOperator(opts *bind.TransactOpts, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "deregisterOperator", quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xca4f2d97.
//
// Solidity: function deregisterOperator(bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) DeregisterOperator(quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.DeregisterOperator(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xca4f2d97.
//
// Solidity: function deregisterOperator(bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) DeregisterOperator(quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.DeregisterOperator(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers)
}

// EjectOperator is a paid mutator transaction binding the contract method 0x6e3b17db.
//
// Solidity: function ejectOperator(address operator, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) EjectOperator(opts *bind.TransactOpts, operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "ejectOperator", operator, quorumNumbers)
}

// EjectOperator is a paid mutator transaction binding the contract method 0x6e3b17db.
//
// Solidity: function ejectOperator(address operator, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) EjectOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.EjectOperator(&_ContractRegistryCoordinator.TransactOpts, operator, quorumNumbers)
}

// EjectOperator is a paid mutator transaction binding the contract method 0x6e3b17db.
//
// Solidity: function ejectOperator(address operator, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) EjectOperator(operator common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.EjectOperator(&_ContractRegistryCoordinator.TransactOpts, operator, quorumNumbers)
}

// Initialize is a paid mutator transaction binding the contract method 0xdd8283f3.
//
// Solidity: function initialize(address _initialOwner, address _churnApprover, address _ejector, address _pauserRegistry, uint256 _initialPausedStatus, (uint32,uint16,uint16)[] _operatorSetParams, uint96[] _minimumStakes, (address,uint96)[][] _strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) Initialize(opts *bind.TransactOpts, _initialOwner common.Address, _churnApprover common.Address, _ejector common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int, _operatorSetParams []IRegistryCoordinatorOperatorSetParam, _minimumStakes []*big.Int, _strategyParams [][]IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "initialize", _initialOwner, _churnApprover, _ejector, _pauserRegistry, _initialPausedStatus, _operatorSetParams, _minimumStakes, _strategyParams)
}

// Initialize is a paid mutator transaction binding the contract method 0xdd8283f3.
//
// Solidity: function initialize(address _initialOwner, address _churnApprover, address _ejector, address _pauserRegistry, uint256 _initialPausedStatus, (uint32,uint16,uint16)[] _operatorSetParams, uint96[] _minimumStakes, (address,uint96)[][] _strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Initialize(_initialOwner common.Address, _churnApprover common.Address, _ejector common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int, _operatorSetParams []IRegistryCoordinatorOperatorSetParam, _minimumStakes []*big.Int, _strategyParams [][]IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Initialize(&_ContractRegistryCoordinator.TransactOpts, _initialOwner, _churnApprover, _ejector, _pauserRegistry, _initialPausedStatus, _operatorSetParams, _minimumStakes, _strategyParams)
}

// Initialize is a paid mutator transaction binding the contract method 0xdd8283f3.
//
// Solidity: function initialize(address _initialOwner, address _churnApprover, address _ejector, address _pauserRegistry, uint256 _initialPausedStatus, (uint32,uint16,uint16)[] _operatorSetParams, uint96[] _minimumStakes, (address,uint96)[][] _strategyParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) Initialize(_initialOwner common.Address, _churnApprover common.Address, _ejector common.Address, _pauserRegistry common.Address, _initialPausedStatus *big.Int, _operatorSetParams []IRegistryCoordinatorOperatorSetParam, _minimumStakes []*big.Int, _strategyParams [][]IStakeRegistryStrategyParams) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Initialize(&_ContractRegistryCoordinator.TransactOpts, _initialOwner, _churnApprover, _ejector, _pauserRegistry, _initialPausedStatus, _operatorSetParams, _minimumStakes, _strategyParams)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Pause(&_ContractRegistryCoordinator.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Pause(&_ContractRegistryCoordinator.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) PauseAll() (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.PauseAll(&_ContractRegistryCoordinator.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) PauseAll() (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.PauseAll(&_ContractRegistryCoordinator.TransactOpts)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0xa50857bf.
//
// Solidity: function registerOperator(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) RegisterOperator(opts *bind.TransactOpts, quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "registerOperator", quorumNumbers, socket, params, operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0xa50857bf.
//
// Solidity: function registerOperator(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) RegisterOperator(quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RegisterOperator(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers, socket, params, operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0xa50857bf.
//
// Solidity: function registerOperator(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) RegisterOperator(quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RegisterOperator(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers, socket, params, operatorSignature)
}

// RegisterOperatorWithChurn is a paid mutator transaction binding the contract method 0x9b5d177b.
//
// Solidity: function registerOperatorWithChurn(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint8,address)[] operatorKickParams, (bytes,bytes32,uint256) churnApproverSignature, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) RegisterOperatorWithChurn(opts *bind.TransactOpts, quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorKickParams []IRegistryCoordinatorOperatorKickParam, churnApproverSignature ISignatureUtilsSignatureWithSaltAndExpiry, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "registerOperatorWithChurn", quorumNumbers, socket, params, operatorKickParams, churnApproverSignature, operatorSignature)
}

// RegisterOperatorWithChurn is a paid mutator transaction binding the contract method 0x9b5d177b.
//
// Solidity: function registerOperatorWithChurn(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint8,address)[] operatorKickParams, (bytes,bytes32,uint256) churnApproverSignature, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) RegisterOperatorWithChurn(quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorKickParams []IRegistryCoordinatorOperatorKickParam, churnApproverSignature ISignatureUtilsSignatureWithSaltAndExpiry, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RegisterOperatorWithChurn(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers, socket, params, operatorKickParams, churnApproverSignature, operatorSignature)
}

// RegisterOperatorWithChurn is a paid mutator transaction binding the contract method 0x9b5d177b.
//
// Solidity: function registerOperatorWithChurn(bytes quorumNumbers, string socket, ((uint256,uint256),(uint256,uint256),(uint256[2],uint256[2])) params, (uint8,address)[] operatorKickParams, (bytes,bytes32,uint256) churnApproverSignature, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) RegisterOperatorWithChurn(quorumNumbers []byte, socket string, params IBLSApkRegistryPubkeyRegistrationParams, operatorKickParams []IRegistryCoordinatorOperatorKickParam, churnApproverSignature ISignatureUtilsSignatureWithSaltAndExpiry, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RegisterOperatorWithChurn(&_ContractRegistryCoordinator.TransactOpts, quorumNumbers, socket, params, operatorKickParams, churnApproverSignature, operatorSignature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RenounceOwnership(&_ContractRegistryCoordinator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.RenounceOwnership(&_ContractRegistryCoordinator.TransactOpts)
}

// SetChurnApprover is a paid mutator transaction binding the contract method 0x29d1e0c3.
//
// Solidity: function setChurnApprover(address _churnApprover) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) SetChurnApprover(opts *bind.TransactOpts, _churnApprover common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "setChurnApprover", _churnApprover)
}

// SetChurnApprover is a paid mutator transaction binding the contract method 0x29d1e0c3.
//
// Solidity: function setChurnApprover(address _churnApprover) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SetChurnApprover(_churnApprover common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetChurnApprover(&_ContractRegistryCoordinator.TransactOpts, _churnApprover)
}

// SetChurnApprover is a paid mutator transaction binding the contract method 0x29d1e0c3.
//
// Solidity: function setChurnApprover(address _churnApprover) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) SetChurnApprover(_churnApprover common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetChurnApprover(&_ContractRegistryCoordinator.TransactOpts, _churnApprover)
}

// SetEjectionCooldown is a paid mutator transaction binding the contract method 0x0d3f2134.
//
// Solidity: function setEjectionCooldown(uint256 _ejectionCooldown) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) SetEjectionCooldown(opts *bind.TransactOpts, _ejectionCooldown *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "setEjectionCooldown", _ejectionCooldown)
}

// SetEjectionCooldown is a paid mutator transaction binding the contract method 0x0d3f2134.
//
// Solidity: function setEjectionCooldown(uint256 _ejectionCooldown) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SetEjectionCooldown(_ejectionCooldown *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetEjectionCooldown(&_ContractRegistryCoordinator.TransactOpts, _ejectionCooldown)
}

// SetEjectionCooldown is a paid mutator transaction binding the contract method 0x0d3f2134.
//
// Solidity: function setEjectionCooldown(uint256 _ejectionCooldown) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) SetEjectionCooldown(_ejectionCooldown *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetEjectionCooldown(&_ContractRegistryCoordinator.TransactOpts, _ejectionCooldown)
}

// SetEjector is a paid mutator transaction binding the contract method 0x2cdd1e86.
//
// Solidity: function setEjector(address _ejector) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) SetEjector(opts *bind.TransactOpts, _ejector common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "setEjector", _ejector)
}

// SetEjector is a paid mutator transaction binding the contract method 0x2cdd1e86.
//
// Solidity: function setEjector(address _ejector) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SetEjector(_ejector common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetEjector(&_ContractRegistryCoordinator.TransactOpts, _ejector)
}

// SetEjector is a paid mutator transaction binding the contract method 0x2cdd1e86.
//
// Solidity: function setEjector(address _ejector) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) SetEjector(_ejector common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetEjector(&_ContractRegistryCoordinator.TransactOpts, _ejector)
}

// SetOperatorSetParams is a paid mutator transaction binding the contract method 0x5b0b829f.
//
// Solidity: function setOperatorSetParams(uint8 quorumNumber, (uint32,uint16,uint16) operatorSetParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) SetOperatorSetParams(opts *bind.TransactOpts, quorumNumber uint8, operatorSetParams IRegistryCoordinatorOperatorSetParam) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "setOperatorSetParams", quorumNumber, operatorSetParams)
}

// SetOperatorSetParams is a paid mutator transaction binding the contract method 0x5b0b829f.
//
// Solidity: function setOperatorSetParams(uint8 quorumNumber, (uint32,uint16,uint16) operatorSetParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SetOperatorSetParams(quorumNumber uint8, operatorSetParams IRegistryCoordinatorOperatorSetParam) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetOperatorSetParams(&_ContractRegistryCoordinator.TransactOpts, quorumNumber, operatorSetParams)
}

// SetOperatorSetParams is a paid mutator transaction binding the contract method 0x5b0b829f.
//
// Solidity: function setOperatorSetParams(uint8 quorumNumber, (uint32,uint16,uint16) operatorSetParams) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) SetOperatorSetParams(quorumNumber uint8, operatorSetParams IRegistryCoordinatorOperatorSetParam) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetOperatorSetParams(&_ContractRegistryCoordinator.TransactOpts, quorumNumber, operatorSetParams)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetPauserRegistry(&_ContractRegistryCoordinator.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.SetPauserRegistry(&_ContractRegistryCoordinator.TransactOpts, newPauserRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.TransferOwnership(&_ContractRegistryCoordinator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.TransferOwnership(&_ContractRegistryCoordinator.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Unpause(&_ContractRegistryCoordinator.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.Unpause(&_ContractRegistryCoordinator.TransactOpts, newPausedStatus)
}

// UpdateOperators is a paid mutator transaction binding the contract method 0x00cf2ab5.
//
// Solidity: function updateOperators(address[] operators) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) UpdateOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "updateOperators", operators)
}

// UpdateOperators is a paid mutator transaction binding the contract method 0x00cf2ab5.
//
// Solidity: function updateOperators(address[] operators) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) UpdateOperators(operators []common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateOperators(&_ContractRegistryCoordinator.TransactOpts, operators)
}

// UpdateOperators is a paid mutator transaction binding the contract method 0x00cf2ab5.
//
// Solidity: function updateOperators(address[] operators) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) UpdateOperators(operators []common.Address) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateOperators(&_ContractRegistryCoordinator.TransactOpts, operators)
}

// UpdateOperatorsForQuorum is a paid mutator transaction binding the contract method 0x5140a548.
//
// Solidity: function updateOperatorsForQuorum(address[][] operatorsPerQuorum, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) UpdateOperatorsForQuorum(opts *bind.TransactOpts, operatorsPerQuorum [][]common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "updateOperatorsForQuorum", operatorsPerQuorum, quorumNumbers)
}

// UpdateOperatorsForQuorum is a paid mutator transaction binding the contract method 0x5140a548.
//
// Solidity: function updateOperatorsForQuorum(address[][] operatorsPerQuorum, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) UpdateOperatorsForQuorum(operatorsPerQuorum [][]common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateOperatorsForQuorum(&_ContractRegistryCoordinator.TransactOpts, operatorsPerQuorum, quorumNumbers)
}

// UpdateOperatorsForQuorum is a paid mutator transaction binding the contract method 0x5140a548.
//
// Solidity: function updateOperatorsForQuorum(address[][] operatorsPerQuorum, bytes quorumNumbers) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) UpdateOperatorsForQuorum(operatorsPerQuorum [][]common.Address, quorumNumbers []byte) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateOperatorsForQuorum(&_ContractRegistryCoordinator.TransactOpts, operatorsPerQuorum, quorumNumbers)
}

// UpdateSocket is a paid mutator transaction binding the contract method 0x0cf4b767.
//
// Solidity: function updateSocket(string socket) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactor) UpdateSocket(opts *bind.TransactOpts, socket string) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.contract.Transact(opts, "updateSocket", socket)
}

// UpdateSocket is a paid mutator transaction binding the contract method 0x0cf4b767.
//
// Solidity: function updateSocket(string socket) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) UpdateSocket(socket string) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateSocket(&_ContractRegistryCoordinator.TransactOpts, socket)
}

// UpdateSocket is a paid mutator transaction binding the contract method 0x0cf4b767.
//
// Solidity: function updateSocket(string socket) returns()
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorTransactorSession) UpdateSocket(socket string) (*types.Transaction, error) {
	return _ContractRegistryCoordinator.Contract.UpdateSocket(&_ContractRegistryCoordinator.TransactOpts, socket)
}

// ContractRegistryCoordinatorChurnApproverUpdatedIterator is returned from FilterChurnApproverUpdated and is used to iterate over the raw logs and unpacked data for ChurnApproverUpdated events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorChurnApproverUpdatedIterator struct {
	Event *ContractRegistryCoordinatorChurnApproverUpdated // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorChurnApproverUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorChurnApproverUpdated)
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
		it.Event = new(ContractRegistryCoordinatorChurnApproverUpdated)
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
func (it *ContractRegistryCoordinatorChurnApproverUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorChurnApproverUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorChurnApproverUpdated represents a ChurnApproverUpdated event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorChurnApproverUpdated struct {
	PrevChurnApprover common.Address
	NewChurnApprover  common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterChurnApproverUpdated is a free log retrieval operation binding the contract event 0x315457d8a8fe60f04af17c16e2f5a5e1db612b31648e58030360759ef8f3528c.
//
// Solidity: event ChurnApproverUpdated(address prevChurnApprover, address newChurnApprover)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterChurnApproverUpdated(opts *bind.FilterOpts) (*ContractRegistryCoordinatorChurnApproverUpdatedIterator, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "ChurnApproverUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorChurnApproverUpdatedIterator{contract: _ContractRegistryCoordinator.contract, event: "ChurnApproverUpdated", logs: logs, sub: sub}, nil
}

// WatchChurnApproverUpdated is a free log subscription operation binding the contract event 0x315457d8a8fe60f04af17c16e2f5a5e1db612b31648e58030360759ef8f3528c.
//
// Solidity: event ChurnApproverUpdated(address prevChurnApprover, address newChurnApprover)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchChurnApproverUpdated(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorChurnApproverUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "ChurnApproverUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorChurnApproverUpdated)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "ChurnApproverUpdated", log); err != nil {
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

// ParseChurnApproverUpdated is a log parse operation binding the contract event 0x315457d8a8fe60f04af17c16e2f5a5e1db612b31648e58030360759ef8f3528c.
//
// Solidity: event ChurnApproverUpdated(address prevChurnApprover, address newChurnApprover)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseChurnApproverUpdated(log types.Log) (*ContractRegistryCoordinatorChurnApproverUpdated, error) {
	event := new(ContractRegistryCoordinatorChurnApproverUpdated)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "ChurnApproverUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorEjectorUpdatedIterator is returned from FilterEjectorUpdated and is used to iterate over the raw logs and unpacked data for EjectorUpdated events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorEjectorUpdatedIterator struct {
	Event *ContractRegistryCoordinatorEjectorUpdated // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorEjectorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorEjectorUpdated)
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
		it.Event = new(ContractRegistryCoordinatorEjectorUpdated)
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
func (it *ContractRegistryCoordinatorEjectorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorEjectorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorEjectorUpdated represents a EjectorUpdated event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorEjectorUpdated struct {
	PrevEjector common.Address
	NewEjector  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEjectorUpdated is a free log retrieval operation binding the contract event 0x8f30ab09f43a6c157d7fce7e0a13c003042c1c95e8a72e7a146a21c0caa24dc9.
//
// Solidity: event EjectorUpdated(address prevEjector, address newEjector)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterEjectorUpdated(opts *bind.FilterOpts) (*ContractRegistryCoordinatorEjectorUpdatedIterator, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorEjectorUpdatedIterator{contract: _ContractRegistryCoordinator.contract, event: "EjectorUpdated", logs: logs, sub: sub}, nil
}

// WatchEjectorUpdated is a free log subscription operation binding the contract event 0x8f30ab09f43a6c157d7fce7e0a13c003042c1c95e8a72e7a146a21c0caa24dc9.
//
// Solidity: event EjectorUpdated(address prevEjector, address newEjector)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchEjectorUpdated(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorEjectorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "EjectorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorEjectorUpdated)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
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

// ParseEjectorUpdated is a log parse operation binding the contract event 0x8f30ab09f43a6c157d7fce7e0a13c003042c1c95e8a72e7a146a21c0caa24dc9.
//
// Solidity: event EjectorUpdated(address prevEjector, address newEjector)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseEjectorUpdated(log types.Log) (*ContractRegistryCoordinatorEjectorUpdated, error) {
	event := new(ContractRegistryCoordinatorEjectorUpdated)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "EjectorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorInitializedIterator struct {
	Event *ContractRegistryCoordinatorInitialized // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorInitialized)
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
		it.Event = new(ContractRegistryCoordinatorInitialized)
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
func (it *ContractRegistryCoordinatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorInitialized represents a Initialized event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractRegistryCoordinatorInitializedIterator, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorInitializedIterator{contract: _ContractRegistryCoordinator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorInitialized)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseInitialized(log types.Log) (*ContractRegistryCoordinatorInitialized, error) {
	event := new(ContractRegistryCoordinatorInitialized)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorOperatorDeregisteredIterator is returned from FilterOperatorDeregistered and is used to iterate over the raw logs and unpacked data for OperatorDeregistered events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorDeregisteredIterator struct {
	Event *ContractRegistryCoordinatorOperatorDeregistered // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorOperatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorOperatorDeregistered)
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
		it.Event = new(ContractRegistryCoordinatorOperatorDeregistered)
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
func (it *ContractRegistryCoordinatorOperatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorOperatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorOperatorDeregistered represents a OperatorDeregistered event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorDeregistered struct {
	Operator   common.Address
	OperatorId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistered is a free log retrieval operation binding the contract event 0x396fdcb180cb0fea26928113fb0fd1c3549863f9cd563e6a184f1d578116c8e4.
//
// Solidity: event OperatorDeregistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterOperatorDeregistered(opts *bind.FilterOpts, operator []common.Address, operatorId [][32]byte) (*ContractRegistryCoordinatorOperatorDeregisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "OperatorDeregistered", operatorRule, operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorOperatorDeregisteredIterator{contract: _ContractRegistryCoordinator.contract, event: "OperatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistered is a free log subscription operation binding the contract event 0x396fdcb180cb0fea26928113fb0fd1c3549863f9cd563e6a184f1d578116c8e4.
//
// Solidity: event OperatorDeregistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchOperatorDeregistered(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorOperatorDeregistered, operator []common.Address, operatorId [][32]byte) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "OperatorDeregistered", operatorRule, operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorOperatorDeregistered)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
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

// ParseOperatorDeregistered is a log parse operation binding the contract event 0x396fdcb180cb0fea26928113fb0fd1c3549863f9cd563e6a184f1d578116c8e4.
//
// Solidity: event OperatorDeregistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseOperatorDeregistered(log types.Log) (*ContractRegistryCoordinatorOperatorDeregistered, error) {
	event := new(ContractRegistryCoordinatorOperatorDeregistered)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorOperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorRegisteredIterator struct {
	Event *ContractRegistryCoordinatorOperatorRegistered // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorOperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorOperatorRegistered)
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
		it.Event = new(ContractRegistryCoordinatorOperatorRegistered)
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
func (it *ContractRegistryCoordinatorOperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorOperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorOperatorRegistered represents a OperatorRegistered event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorRegistered struct {
	Operator   common.Address
	OperatorId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0xe8e68cef1c3a761ed7be7e8463a375f27f7bc335e51824223cacce636ec5c3fe.
//
// Solidity: event OperatorRegistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address, operatorId [][32]byte) (*ContractRegistryCoordinatorOperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "OperatorRegistered", operatorRule, operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorOperatorRegisteredIterator{contract: _ContractRegistryCoordinator.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0xe8e68cef1c3a761ed7be7e8463a375f27f7bc335e51824223cacce636ec5c3fe.
//
// Solidity: event OperatorRegistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorOperatorRegistered, operator []common.Address, operatorId [][32]byte) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "OperatorRegistered", operatorRule, operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorOperatorRegistered)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
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

// ParseOperatorRegistered is a log parse operation binding the contract event 0xe8e68cef1c3a761ed7be7e8463a375f27f7bc335e51824223cacce636ec5c3fe.
//
// Solidity: event OperatorRegistered(address indexed operator, bytes32 indexed operatorId)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseOperatorRegistered(log types.Log) (*ContractRegistryCoordinatorOperatorRegistered, error) {
	event := new(ContractRegistryCoordinatorOperatorRegistered)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator is returned from FilterOperatorSetParamsUpdated and is used to iterate over the raw logs and unpacked data for OperatorSetParamsUpdated events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator struct {
	Event *ContractRegistryCoordinatorOperatorSetParamsUpdated // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorOperatorSetParamsUpdated)
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
		it.Event = new(ContractRegistryCoordinatorOperatorSetParamsUpdated)
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
func (it *ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorOperatorSetParamsUpdated represents a OperatorSetParamsUpdated event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorSetParamsUpdated struct {
	QuorumNumber      uint8
	OperatorSetParams IRegistryCoordinatorOperatorSetParam
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterOperatorSetParamsUpdated is a free log retrieval operation binding the contract event 0x3ee6fe8d54610244c3e9d3c066ae4aee997884aa28f10616ae821925401318ac.
//
// Solidity: event OperatorSetParamsUpdated(uint8 indexed quorumNumber, (uint32,uint16,uint16) operatorSetParams)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterOperatorSetParamsUpdated(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "OperatorSetParamsUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorOperatorSetParamsUpdatedIterator{contract: _ContractRegistryCoordinator.contract, event: "OperatorSetParamsUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorSetParamsUpdated is a free log subscription operation binding the contract event 0x3ee6fe8d54610244c3e9d3c066ae4aee997884aa28f10616ae821925401318ac.
//
// Solidity: event OperatorSetParamsUpdated(uint8 indexed quorumNumber, (uint32,uint16,uint16) operatorSetParams)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchOperatorSetParamsUpdated(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorOperatorSetParamsUpdated, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "OperatorSetParamsUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorOperatorSetParamsUpdated)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorSetParamsUpdated", log); err != nil {
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

// ParseOperatorSetParamsUpdated is a log parse operation binding the contract event 0x3ee6fe8d54610244c3e9d3c066ae4aee997884aa28f10616ae821925401318ac.
//
// Solidity: event OperatorSetParamsUpdated(uint8 indexed quorumNumber, (uint32,uint16,uint16) operatorSetParams)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseOperatorSetParamsUpdated(log types.Log) (*ContractRegistryCoordinatorOperatorSetParamsUpdated, error) {
	event := new(ContractRegistryCoordinatorOperatorSetParamsUpdated)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorSetParamsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorOperatorSocketUpdateIterator is returned from FilterOperatorSocketUpdate and is used to iterate over the raw logs and unpacked data for OperatorSocketUpdate events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorSocketUpdateIterator struct {
	Event *ContractRegistryCoordinatorOperatorSocketUpdate // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorOperatorSocketUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorOperatorSocketUpdate)
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
		it.Event = new(ContractRegistryCoordinatorOperatorSocketUpdate)
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
func (it *ContractRegistryCoordinatorOperatorSocketUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorOperatorSocketUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorOperatorSocketUpdate represents a OperatorSocketUpdate event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOperatorSocketUpdate struct {
	OperatorId [32]byte
	Socket     string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorSocketUpdate is a free log retrieval operation binding the contract event 0xec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa.
//
// Solidity: event OperatorSocketUpdate(bytes32 indexed operatorId, string socket)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterOperatorSocketUpdate(opts *bind.FilterOpts, operatorId [][32]byte) (*ContractRegistryCoordinatorOperatorSocketUpdateIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "OperatorSocketUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorOperatorSocketUpdateIterator{contract: _ContractRegistryCoordinator.contract, event: "OperatorSocketUpdate", logs: logs, sub: sub}, nil
}

// WatchOperatorSocketUpdate is a free log subscription operation binding the contract event 0xec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa.
//
// Solidity: event OperatorSocketUpdate(bytes32 indexed operatorId, string socket)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchOperatorSocketUpdate(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorOperatorSocketUpdate, operatorId [][32]byte) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "OperatorSocketUpdate", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorOperatorSocketUpdate)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorSocketUpdate", log); err != nil {
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

// ParseOperatorSocketUpdate is a log parse operation binding the contract event 0xec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa.
//
// Solidity: event OperatorSocketUpdate(bytes32 indexed operatorId, string socket)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseOperatorSocketUpdate(log types.Log) (*ContractRegistryCoordinatorOperatorSocketUpdate, error) {
	event := new(ContractRegistryCoordinatorOperatorSocketUpdate)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OperatorSocketUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOwnershipTransferredIterator struct {
	Event *ContractRegistryCoordinatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorOwnershipTransferred)
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
		it.Event = new(ContractRegistryCoordinatorOwnershipTransferred)
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
func (it *ContractRegistryCoordinatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorOwnershipTransferred represents a OwnershipTransferred event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractRegistryCoordinatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorOwnershipTransferredIterator{contract: _ContractRegistryCoordinator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorOwnershipTransferred)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseOwnershipTransferred(log types.Log) (*ContractRegistryCoordinatorOwnershipTransferred, error) {
	event := new(ContractRegistryCoordinatorOwnershipTransferred)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorPausedIterator struct {
	Event *ContractRegistryCoordinatorPaused // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorPaused)
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
		it.Event = new(ContractRegistryCoordinatorPaused)
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
func (it *ContractRegistryCoordinatorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorPaused represents a Paused event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*ContractRegistryCoordinatorPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorPausedIterator{contract: _ContractRegistryCoordinator.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorPaused)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParsePaused(log types.Log) (*ContractRegistryCoordinatorPaused, error) {
	event := new(ContractRegistryCoordinatorPaused)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorPauserRegistrySetIterator struct {
	Event *ContractRegistryCoordinatorPauserRegistrySet // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorPauserRegistrySet)
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
		it.Event = new(ContractRegistryCoordinatorPauserRegistrySet)
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
func (it *ContractRegistryCoordinatorPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorPauserRegistrySet represents a PauserRegistrySet event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*ContractRegistryCoordinatorPauserRegistrySetIterator, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorPauserRegistrySetIterator{contract: _ContractRegistryCoordinator.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorPauserRegistrySet)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
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

// ParsePauserRegistrySet is a log parse operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParsePauserRegistrySet(log types.Log) (*ContractRegistryCoordinatorPauserRegistrySet, error) {
	event := new(ContractRegistryCoordinatorPauserRegistrySet)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator is returned from FilterQuorumBlockNumberUpdated and is used to iterate over the raw logs and unpacked data for QuorumBlockNumberUpdated events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator struct {
	Event *ContractRegistryCoordinatorQuorumBlockNumberUpdated // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorQuorumBlockNumberUpdated)
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
		it.Event = new(ContractRegistryCoordinatorQuorumBlockNumberUpdated)
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
func (it *ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorQuorumBlockNumberUpdated represents a QuorumBlockNumberUpdated event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorQuorumBlockNumberUpdated struct {
	QuorumNumber uint8
	Blocknumber  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterQuorumBlockNumberUpdated is a free log retrieval operation binding the contract event 0x46077d55330763f16269fd75e5761663f4192d2791747c0189b16ad31db07db4.
//
// Solidity: event QuorumBlockNumberUpdated(uint8 indexed quorumNumber, uint256 blocknumber)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterQuorumBlockNumberUpdated(opts *bind.FilterOpts, quorumNumber []uint8) (*ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "QuorumBlockNumberUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorQuorumBlockNumberUpdatedIterator{contract: _ContractRegistryCoordinator.contract, event: "QuorumBlockNumberUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumBlockNumberUpdated is a free log subscription operation binding the contract event 0x46077d55330763f16269fd75e5761663f4192d2791747c0189b16ad31db07db4.
//
// Solidity: event QuorumBlockNumberUpdated(uint8 indexed quorumNumber, uint256 blocknumber)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchQuorumBlockNumberUpdated(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorQuorumBlockNumberUpdated, quorumNumber []uint8) (event.Subscription, error) {

	var quorumNumberRule []interface{}
	for _, quorumNumberItem := range quorumNumber {
		quorumNumberRule = append(quorumNumberRule, quorumNumberItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "QuorumBlockNumberUpdated", quorumNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorQuorumBlockNumberUpdated)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "QuorumBlockNumberUpdated", log); err != nil {
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

// ParseQuorumBlockNumberUpdated is a log parse operation binding the contract event 0x46077d55330763f16269fd75e5761663f4192d2791747c0189b16ad31db07db4.
//
// Solidity: event QuorumBlockNumberUpdated(uint8 indexed quorumNumber, uint256 blocknumber)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseQuorumBlockNumberUpdated(log types.Log) (*ContractRegistryCoordinatorQuorumBlockNumberUpdated, error) {
	event := new(ContractRegistryCoordinatorQuorumBlockNumberUpdated)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "QuorumBlockNumberUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRegistryCoordinatorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorUnpausedIterator struct {
	Event *ContractRegistryCoordinatorUnpaused // Event containing the contract specifics and raw log

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
func (it *ContractRegistryCoordinatorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRegistryCoordinatorUnpaused)
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
		it.Event = new(ContractRegistryCoordinatorUnpaused)
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
func (it *ContractRegistryCoordinatorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRegistryCoordinatorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRegistryCoordinatorUnpaused represents a Unpaused event raised by the ContractRegistryCoordinator contract.
type ContractRegistryCoordinatorUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*ContractRegistryCoordinatorUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractRegistryCoordinatorUnpausedIterator{contract: _ContractRegistryCoordinator.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractRegistryCoordinatorUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractRegistryCoordinator.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRegistryCoordinatorUnpaused)
				if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorFilterer) ParseUnpaused(log types.Log) (*ContractRegistryCoordinatorUnpaused, error) {
	event := new(ContractRegistryCoordinatorUnpaused)
	if err := _ContractRegistryCoordinator.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
