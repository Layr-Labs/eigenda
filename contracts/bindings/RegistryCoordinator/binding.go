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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_serviceManager\",\"type\":\"address\",\"internalType\":\"contractIServiceManager\"},{\"name\":\"_stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"},{\"name\":\"_blsApkRegistry\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"},{\"name\":\"_indexRegistry\",\"type\":\"address\",\"internalType\":\"contractIIndexRegistry\"},{\"name\":\"_socketRegistry\",\"type\":\"address\",\"internalType\":\"contractISocketRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"OPERATOR_CHURN_APPROVAL_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PUBKEY_REGISTRATION_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateOperatorChurnApprovalDigestHash\",\"inputs\":[{\"name\":\"registeringOperator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registeringOperatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"operatorKickParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorKickParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"churnApprover\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createQuorum\",\"inputs\":[{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"name\":\"minimumStake\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"strategyParams\",\"type\":\"tuple[]\",\"internalType\":\"structIStakeRegistry.StrategyParams[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ejectionCooldown\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ejector\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentQuorumBitmap\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorInfo\",\"components\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRegistryCoordinator.OperatorStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorFromId\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorId\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorSetParams\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorStatus\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIRegistryCoordinator.OperatorStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapAtBlockNumberByIndex\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint192\",\"internalType\":\"uint192\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapHistoryLength\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapIndicesAtBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"operatorIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumBitmapUpdateByIndex\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.QuorumBitmapUpdate\",\"components\":[{\"name\":\"updateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"nextUpdateBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBitmap\",\"type\":\"uint192\",\"internalType\":\"uint192\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"indexRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIIndexRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_churnApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_ejector\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_operatorSetParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam[]\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"name\":\"_minimumStakes\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"_strategyParams\",\"type\":\"tuple[][]\",\"internalType\":\"structIStakeRegistry.StrategyParams[][]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isChurnApproverSaltUsed\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastEjectionTimestamp\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"numRegistries\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyRegistrationMessageHash\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumUpdateBlockNumber\",\"inputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.PubkeyRegistrationParams\",\"components\":[{\"name\":\"pubkeyRegistrationSignature\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}]},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperatorWithChurn\",\"inputs\":[{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSApkRegistry.PubkeyRegistrationParams\",\"components\":[{\"name\":\"pubkeyRegistrationSignature\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG1\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"pubkeyG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]}]},{\"name\":\"operatorKickParams\",\"type\":\"tuple[]\",\"internalType\":\"structIRegistryCoordinator.OperatorKickParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"name\":\"churnApproverSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registries\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"serviceManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIServiceManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setChurnApprover\",\"inputs\":[{\"name\":\"_churnApprover\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEjectionCooldown\",\"inputs\":[{\"name\":\"_ejectionCooldown\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEjector\",\"inputs\":[{\"name\":\"_ejector\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorSetParams\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"socketRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISocketRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperatorsForQuorum\",\"inputs\":[{\"name\":\"operatorsPerQuorum\",\"type\":\"address[][]\",\"internalType\":\"address[][]\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSocket\",\"inputs\":[{\"name\":\"socket\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ChurnApproverUpdated\",\"inputs\":[{\"name\":\"prevChurnApprover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newChurnApprover\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EjectorUpdated\",\"inputs\":[{\"name\":\"prevEjector\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newEjector\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSetParamsUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"operatorSetParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIRegistryCoordinator.OperatorSetParam\",\"components\":[{\"name\":\"maxOperatorCount\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"kickBIPsOfOperatorStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"kickBIPsOfTotalStake\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSocketUpdate\",\"inputs\":[{\"name\":\"operatorId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"socket\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumBlockNumberUpdated\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"blocknumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6101e06040523480156200001257600080fd5b506040516200abfe3803806200abfe833981810160405281019062000038919062000519565b84848484846040518060400160405280601681526020017f4156535265676973747279436f6f7264696e61746f72000000000000000000008152506040518060400160405280600681526020017f76302e302e31000000000000000000000000000000000000000000000000000081525060008280519060200120905060008280519060200120905060007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f90508260e081815250508161010081815250504660a08181525050620001128184846200028460201b60201c565b608081815250503073ffffffffffffffffffffffffffffffffffffffff1660c08173ffffffffffffffffffffffffffffffffffffffff168152505080610120818152505050505050508473ffffffffffffffffffffffffffffffffffffffff166101408173ffffffffffffffffffffffffffffffffffffffff16815250508373ffffffffffffffffffffffffffffffffffffffff166101808173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff166101608173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff166101a08173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff166101c08173ffffffffffffffffffffffffffffffffffffffff1681525050505050505062000279620002c060201b60201c565b505050505062000729565b60008383834630604051602001620002a1959493929190620005e8565b6040516020818303038152906040528051906020012090509392505050565b600060019054906101000a900460ff161562000313576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200030a90620006cc565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff161015620003855760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff6040516200037c91906200070c565b60405180910390a15b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620003b9826200038c565b9050919050565b6000620003cd82620003ac565b9050919050565b620003df81620003c0565b8114620003eb57600080fd5b50565b600081519050620003ff81620003d4565b92915050565b60006200041282620003ac565b9050919050565b620004248162000405565b81146200043057600080fd5b50565b600081519050620004448162000419565b92915050565b60006200045782620003ac565b9050919050565b62000469816200044a565b81146200047557600080fd5b50565b60008151905062000489816200045e565b92915050565b60006200049c82620003ac565b9050919050565b620004ae816200048f565b8114620004ba57600080fd5b50565b600081519050620004ce81620004a3565b92915050565b6000620004e182620003ac565b9050919050565b620004f381620004d4565b8114620004ff57600080fd5b50565b6000815190506200051381620004e8565b92915050565b600080600080600060a0868803121562000538576200053762000387565b5b60006200054888828901620003ee565b95505060206200055b8882890162000433565b94505060406200056e8882890162000478565b93505060606200058188828901620004bd565b9250506080620005948882890162000502565b9150509295509295909350565b6000819050919050565b620005b681620005a1565b82525050565b6000819050919050565b620005d181620005bc565b82525050565b620005e281620003ac565b82525050565b600060a082019050620005ff6000830188620005ab565b6200060e6020830187620005ab565b6200061d6040830186620005ab565b6200062c6060830185620005c6565b6200063b6080830184620005d7565b9695505050505050565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b6000620006b460278362000645565b9150620006c18262000656565b604082019050919050565b60006020820190508181036000830152620006e781620006a5565b9050919050565b600060ff82169050919050565b6200070681620006ee565b82525050565b6000602082019050620007236000830184620006fb565b92915050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516101c05161a3bb6200084360003960008181612a2b0152612f9b015260008181611552015281816122420152818161283401528181613b070152818161449101526149e8015260008181611ba00152818161272e01528181612e9d01528181613a7a015281816143da0152818161465701526149590152600081816112ab01528181611b3d015281816127b1015281816139ed01528181613c8f01528181613d360152818161434b0152614a730152600081816113740152818161391d015261426f01526000614d4601526000614d8801526000614d6701526000614c9c01526000614cf201526000614d1b015261a3bb6000f3fe608060405234801561001057600080fd5b50600436106102f05760003560e01c80635df459461161019d578063a50857bf116100e9578063d75b4c88116100a2578063ea32afae1161007c578063ea32afae14610951578063f2fde38b1461096f578063fabc1cbc1461098b578063fd39105a146109a7576102f0565b8063d75b4c88146108e9578063dd8283f314610905578063e65797ad14610921576102f0565b8063a50857bf14610827578063a96f783e14610843578063c391425e14610861578063ca0de88214610891578063ca4f2d97146108af578063d72d8dd6146108cb576102f0565b8063871ef049116101565780639aa1653d116101305780639aa1653d146107b15780639b5d177b146107cf5780639e9923c2146107eb5780639feab85914610809576102f0565b8063871ef04914610745578063886f1195146107755780638da5cb5b14610793576102f0565b80635df45946146106835780636347c900146106a157806368304835146106d15780636e3b17db146106ef578063715018a61461070b57806384ca521314610715576102f0565b8063249a0c421161025c5780633c2a7f4c11610215578063595c6a67116101ef578063595c6a671461060f5780635ac86ab7146106195780635b0b829f146106495780635c975abb14610665576102f0565b80633c2a7f4c146105935780635140a548146105c35780635865c60c146105df576102f0565b8063249a0c42146104bf57806328f61b31146104ef578063296bb0641461050d57806329d1e0c31461053d5780632cdd1e86146105595780633998fdd314610575576102f0565b806310d67a2f116102ae57806310d67a2f146103c7578063125e0584146103e357806313542a4e14610413578063136439dd146104435780631478851f1461045f5780631eb812da1461048f576102f0565b8062cf2ab5146102f557806303fd34921461031157806304ec635114610341578063054310e6146103715780630cf4b7671461038f5780630d3f2134146103ab575b600080fd5b61030f600480360381019061030a9190615c90565b6109d7565b005b61032b60048036038101906103269190615d13565b610b55565b6040516103389190615d59565b60405180910390f35b61035b60048036038101906103569190615ddc565b610b75565b6040516103689190615e62565b60405180910390f35b610379610d2b565b6040516103869190615ebe565b60405180910390f35b6103a960048036038101906103a4919061601a565b610d51565b005b6103c560048036038101906103c09190616063565b610e56565b005b6103e160048036038101906103dc91906160ce565b610e68565b005b6103fd60048036038101906103f89190616127565b610f72565b60405161040a9190615d59565b60405180910390f35b61042d60048036038101906104289190616127565b610f8a565b60405161043a9190616163565b60405180910390f35b61045d60048036038101906104589190616063565b610fd6565b005b61047960048036038101906104749190615d13565b611151565b6040516104869190616199565b60405180910390f35b6104a960048036038101906104a491906161b4565b611171565b6040516104b69190616254565b60405180910390f35b6104d960048036038101906104d491906162a8565b611269565b6040516104e69190615d59565b60405180910390f35b6104f7611281565b6040516105049190615ebe565b60405180910390f35b61052760048036038101906105229190615d13565b6112a7565b6040516105349190615ebe565b60405180910390f35b61055760048036038101906105529190616127565b61134a565b005b610573600480360381019061056e9190616127565b61135e565b005b61057d611372565b60405161058a9190616334565b60405180910390f35b6105ad60048036038101906105a89190616127565b611396565b6040516105ba919061638d565b60405180910390f35b6105dd60048036038101906105d89190616454565b6113fe565b005b6105f960048036038101906105f49190616127565b6118e4565b604051610606919061658a565b60405180910390f35b610617611982565b005b610633600480360381019061062e91906162a8565b611af4565b6040516106409190616199565b60405180910390f35b610663600480360381019061065e919061664d565b611b10565b005b61066d611b31565b60405161067a9190615d59565b60405180910390f35b61068b611b3b565b60405161069891906166ae565b60405180910390f35b6106bb60048036038101906106b69190616063565b611b5f565b6040516106c89190615ebe565b60405180910390f35b6106d9611b9e565b6040516106e691906166ea565b60405180910390f35b61070960048036038101906107049190616705565b611bc2565b005b610713611dd9565b005b61072f600480360381019061072a9190616878565b611ded565b60405161073c9190616163565b60405180910390f35b61075f600480360381019061075a9190615d13565b611e53565b60405161076c9190615e62565b60405180910390f35b61077d611e65565b60405161078a9190616930565b60405180910390f35b61079b611e8b565b6040516107a89190615ebe565b60405180910390f35b6107b9611e9a565b6040516107c6919061695a565b60405180910390f35b6107e960048036038101906107e49190616b67565b611ead565b005b6107f3612240565b6040516108009190616cb3565b60405180910390f35b610811612264565b60405161081e9190616163565b60405180910390f35b610841600480360381019061083c9190616cce565b612288565b005b61084b612422565b6040516108589190615d59565b60405180910390f35b61087b60048036038101906108769190616e56565b612428565b6040516108889190616f61565b60405180910390f35b6108996124f3565b6040516108a69190616163565b60405180910390f35b6108c960048036038101906108c49190616f83565b612517565b005b6108d36125b5565b6040516108e09190615d59565b60405180910390f35b61090360048036038101906108fe9190617165565b6125c2565b005b61091f600480360381019061091a919061743b565b6125da565b005b61093b600480360381019061093691906162a8565b612990565b6040516109489190617596565b60405180910390f35b610959612a29565b60405161096691906175d2565b60405180910390f35b61098960048036038101906109849190616127565b612a4d565b005b6109a560048036038101906109a09190616063565b612ad1565b005b6109c160048036038101906109bc9190616127565b612c72565b6040516109ce91906175fc565b60405180910390f35b60026109e281611af4565b15610a22576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a1990617674565b60405180910390fd5b60005b83839050811015610b4f576000848483818110610a4557610a44617694565b5b9050602002016020810190610a5a9190616127565b90506000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020604051806040016040529081600082015481526020016001820160009054906101000a900460ff166002811115610ad757610ad66164e4565b5b6002811115610ae957610ae86164e4565b5b8152505090506000816000015190506000610b0382612ccb565b90506000610b2a8277ffffffffffffffffffffffffffffffffffffffffffffffff16612d66565b9050610b37858583612e5b565b50505050508080610b47906176f2565b915050610a25565b50505050565b600060986000838152602001908152602001600020805490509050919050565b600080609860008681526020019081526020016000208381548110610b9d57610b9c617694565b5b906000526020600020016040518060600160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900477ffffffffffffffffffffffffffffffffffffffffffffffff1677ffffffffffffffffffffffffffffffffffffffffffffffff1677ffffffffffffffffffffffffffffffffffffffffffffffff16815250509050806000015163ffffffff168463ffffffff161015610cb5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cac906177d3565b60405180910390fd5b6000816020015163ffffffff161480610cdd5750806020015163ffffffff168463ffffffff16105b610d1c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d139061788b565b60405180910390fd5b80604001519150509392505050565b609d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60016002811115610d6557610d646164e4565b5b609960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff166002811115610dc757610dc66164e4565b5b14610e07576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dfe9061791d565b60405180910390fd5b610e53609960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000015482612f99565b50565b610e5e613062565b8060a08190555050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ed5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ef99190617952565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f66576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f5d906179f1565b60405180910390fd5b610f6f816130e0565b50565b609f6020528060005260406000206000915090505481565b6000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001549050919050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b81526004016110319190615ebe565b602060405180830381865afa15801561104e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110729190617a3d565b6110b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110a890617adc565b60405180910390fd5b6001548160015416146110f9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110f090617b6e565b60405180910390fd5b806001819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d826040516111469190615d59565b60405180910390a250565b609a6020528060005260406000206000915054906101000a900460ff1681565b611179615af3565b60986000848152602001908152602001600020828154811061119e5761119d617694565b5b906000526020600020016040518060600160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160089054906101000a900477ffffffffffffffffffffffffffffffffffffffffffffffff1677ffffffffffffffffffffffffffffffffffffffffffffffff1677ffffffffffffffffffffffffffffffffffffffffffffffff1681525050905092915050565b609b6020528060005260406000206000915090505481565b609e60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166347b314e8836040518263ffffffff1660e01b81526004016113029190616163565b602060405180830381865afa15801561131f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113439190617952565b9050919050565b611352613062565b61135b816131ef565b50565b611366613062565b61136f8161328e565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b61139e615b3a565b6113f76113f27f2bd82124057f0913bc3b772ce7b83e8057c1ad1f3510fc83778be20f10ec5de6846040516020016113d7929190617b8e565b6040516020818303038152906040528051906020012061332d565b613347565b9050919050565b600261140981611af4565b15611449576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161144090617674565b60405180910390fd5b60006114a884848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050609660009054906101000a900460ff16613452565b90508383905086869050146114f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114e990617c29565b60405180910390fd5b60005b848490508110156118db57600085858381811061151557611514617694565b5b9050013560f81c60f81b60f81c905036600089898581811061153a57611539617694565b5b905060200281019061154c9190617c58565b915091507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f3410922846040518263ffffffff1660e01b81526004016115a9919061695a565b602060405180830381865afa1580156115c6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115ea9190617cd0565b63ffffffff168282905014611634576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161162b90617d95565b60405180910390fd5b6000805b8383905081101561186c57600084848381811061165857611657617694565b5b905060200201602081019061166d9190616127565b90506000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020604051806040016040529081600082015481526020016001820160009054906101000a900460ff1660028111156116ea576116e96164e4565b5b60028111156116fc576116fb6164e4565b5b815250509050600081600001519050600061171682612ccb565b905061173c8177ffffffffffffffffffffffffffffffffffffffffffffffff168a6134b3565b61177b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161177290617e27565b60405180910390fd5b8573ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16116117e9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016117e090617edf565b60405180910390fd5b5061185583838f8f8d9060018f6118009190617eff565b9261180d93929190617f5f565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050612e5b565b82945050505080611865906176f2565b9050611638565b5043609b60008660ff1660ff168152602001908152602001600020819055508360ff167f46077d55330763f16269fd75e5761663f4192d2791747c0189b16ad31db07db4436040516118be9190615d59565b60405180910390a250505050806118d4906176f2565b90506114f5565b50505050505050565b6118ec615b54565b609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020604051806040016040529081600082015481526020016001820160009054906101000a900460ff166002811115611965576119646164e4565b5b6002811115611977576119766164e4565b5b815250509050919050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b81526004016119dd9190615ebe565b602060405180830381865afa1580156119fa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a1e9190617a3d565b611a5d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611a5490617adc565b60405180910390fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6001819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff604051611aea9190615d59565b60405180910390a2565b6000808260ff166001901b905080816001541614915050919050565b611b18613062565b81611b22816134ca565b611b2c8383613524565b505050565b6000600154905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b609c8181548110611b6f57600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7f000000000000000000000000000000000000000000000000000000000000000081565b611bca6135ed565b42609f60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000609960008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002090506000816000015490506000611cb985858080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050609660009054906101000a900460ff16613452565b90506000611cc683612ccb565b905060016002811115611cdc57611cdb6164e4565b5b8460010160009054906101000a900460ff166002811115611d0057611cff6164e4565b5b148015611d2d5750611d2b8277ffffffffffffffffffffffffffffffffffffffffffffffff1661367f565b155b8015611d7c5750611d7b8177ffffffffffffffffffffffffffffffffffffffffffffffff168377ffffffffffffffffffffffffffffffffffffffffffffffff1661368b90919063ffffffff16565b5b15611dd057611dcf8787878080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061369a565b5b50505050505050565b611de1613062565b611deb6000613b9b565b565b6000611e487f4d404e3276e7ac2163d8ee476afa6a41d1f68fb71f2d8b6546b24e55ce01b72a8787878787604051602001611e2d96959493929190618096565b6040516020818303038152906040528051906020012061332d565b905095945050505050565b6000611e5e82612ccb565b9050919050565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000611e95613c61565b905090565b609660009054906101000a900460ff1681565b6000611eb881611af4565b15611ef8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611eef90617674565b60405180910390fd5b898990508585905014611f40576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611f3790618170565b60405180910390fd5b6000611f4c3388613c8b565b9050611fb033828888808060200260200160405190810160405280939291908181526020016000905b82821015611fa557848483905060400201803603810190611f969190618190565b81526020019060010190611f75565b505050505087613de5565b600061200433838e8e8e8e8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505089613f0c565b905060005b8c8c9050811015612231576000609760008f8f8581811061202d5761202c617694565b5b9050013560f81c60f81b60f81c60ff1660ff1681526020019081526020016000206040518060600160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900461ffff1661ffff1661ffff1681526020016000820160069054906101000a900461ffff1661ffff1661ffff16815250509050806000015163ffffffff16836000015183815181106120e3576120e2617694565b5b602002602001015163ffffffff16111561221d576121888e8e8481811061210d5761210c617694565b5b9050013560f81c60f81b60f81c8460400151848151811061213157612130617694565b5b6020026020010151338660200151868151811061215157612150617694565b5b60200260200101518d8d8881811061216c5761216b617694565b5b9050604002018036038101906121829190618190565b86614548565b61221c89898481811061219e5761219d617694565b5b90506040020160200160208101906121b69190616127565b8f8f85906001876121c79190617eff565b926121d493929190617f5f565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061369a565b5b508080612229906176f2565b915050612009565b50505050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f2bd82124057f0913bc3b772ce7b83e8057c1ad1f3510fc83778be20f10ec5de681565b600061229381611af4565b156122d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122ca90617674565b60405180910390fd5b60006122df3385613c8b565b9050600061233533838b8b8b8b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505089613f0c565b60000151905060005b898990508110156124165760008a8a8381811061235e5761235d617694565b5b9050013560f81c60f81b60f81c9050609760008260ff1660ff16815260200190815260200160002060000160009054906101000a900463ffffffff1663ffffffff168383815181106123b3576123b2617694565b5b602002602001015163ffffffff161115612402576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016123f99061822f565b60405180910390fd5b50808061240e906176f2565b91505061233e565b50505050505050505050565b60a05481565b60606000825167ffffffffffffffff81111561244757612446615eef565b5b6040519080825280602002602001820160405280156124755781602001602082028036833780820191505090505b50905060005b83518110156124e8576124a88585838151811061249b5761249a617694565b5b60200260200101516147cc565b8282815181106124bb576124ba617694565b5b602002602001019063ffffffff16908163ffffffff168152505080806124e0906176f2565b91505061247b565b508091505092915050565b7f4d404e3276e7ac2163d8ee476afa6a41d1f68fb71f2d8b6546b24e55ce01b72a81565b600161252281611af4565b15612562576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161255990617674565b60405180910390fd5b6125b03384848080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061369a565b505050565b6000609c80549050905090565b6125ca613062565b6125d58383836148c4565b505050565b60008060019054906101000a900460ff1615905080801561260b5750600160008054906101000a900460ff1660ff16105b80612638575061261a30614b03565b1580156126375750600160008054906101000a900460ff1660ff16145b5b612677576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161266e906182c1565b60405180910390fd5b60016000806101000a81548160ff021916908360ff16021790555080156126b4576001600060016101000a81548160ff0219169083151502179055505b825184511480156126c6575081518351145b612705576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016126fc90618353565b60405180910390fd5b61270e89613b9b565b6127188686614b26565b612721886131ef565b61272a8761328e565b609c7f00000000000000000000000000000000000000000000000000000000000000009080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550609c7f00000000000000000000000000000000000000000000000000000000000000009080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550609c7f00000000000000000000000000000000000000000000000000000000000000009080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060005b845181101561292b576129188582815181106128d5576128d4617694565b5b60200260200101518583815181106128f0576128ef617694565b5b602002602001015185848151811061290b5761290a617694565b5b60200260200101516148c4565b8080612923906176f2565b9150506128b6565b5080156129855760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498600160405161297c91906183ae565b60405180910390a15b505050505050505050565b612998615b83565b609760008360ff1660ff1681526020019081526020016000206040518060600160405290816000820160009054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016000820160049054906101000a900461ffff1661ffff1661ffff1681526020016000820160069054906101000a900461ffff1661ffff1661ffff16815250509050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b612a55613062565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415612ac5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612abc9061843b565b60405180910390fd5b612ace81613b9b565b50565b600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612b3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b629190617952565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612bcf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612bc6906179f1565b60405180910390fd5b600154198119600154191614612c1a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612c11906184cd565b60405180910390fd5b806001819055503373ffffffffffffffffffffffffffffffffffffffff167f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c82604051612c679190615d59565b60405180910390a250565b6000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff169050919050565b600080609860008481526020019081526020016000208054905090506000811415612cfa576000915050612d61565b60986000848152602001908152602001600020600182612d1a91906184ed565b81548110612d2b57612d2a617694565b5b9060005260206000200160000160089054906101000a900477ffffffffffffffffffffffffffffffffffffffffffffffff169150505b919050565b6060600080612d7484614c52565b61ffff1667ffffffffffffffff811115612d9157612d90615eef565b5b6040519080825280601f01601f191660200182016040528015612dc35781602001600182028036833780820191505090505b5090506000805b825182108015612ddb575061010081105b15612e4f57806001901b9350600084871614612e3e578060f81b838381518110612e0857612e07617694565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508160010191505b80612e48906176f2565b9050612dca565b50819350505050919050565b60016002811115612e6f57612e6e6164e4565b5b82602001516002811115612e8657612e856164e4565b5b14612e9057612f94565b60008260000151905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166366acfefe8684866040518463ffffffff1660e01b8152600401612ef8939291906185a9565b6020604051808303816000875af1158015612f17573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f3b9190618613565b9050612f608177ffffffffffffffffffffffffffffffffffffffffffffffff1661367f565b612f9157612f9085612f8b8377ffffffffffffffffffffffffffffffffffffffffffffffff16612d66565b61369a565b5b50505b505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f043367e83836040518363ffffffff1660e01b8152600401612ff4929190618684565b600060405180830381600087803b15801561300e57600080fd5b505af1158015613022573d6000803e3d6000fd5b50505050817fec2963ab21c1e50e1e582aa542af2e4bf7bf38e6e1403c27b42e1c5d6e621eaa8260405161305691906186b4565b60405180910390a25050565b61306a614c90565b73ffffffffffffffffffffffffffffffffffffffff16613088611e8b565b73ffffffffffffffffffffffffffffffffffffffff16146130de576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016130d590618722565b60405180910390fd5b565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415613150576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613147906187da565b60405180910390fd5b7f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff16826040516131a39291906187fa565b60405180910390a180600060026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b7f315457d8a8fe60f04af17c16e2f5a5e1db612b31648e58030360759ef8f3528c609d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682604051613242929190618823565b60405180910390a180609d60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b7f8f30ab09f43a6c157d7fce7e0a13c003042c1c95e8a72e7a146a21c0caa24dc9609e60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16826040516132e1929190618823565b60405180910390a180609e60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600061334061333a614c98565b83614db2565b9050919050565b61334f615b3a565b60008060007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478560001c613383919061887b565b90505b6001156134305761339681614de5565b80935081945050507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47806133cd576133cc61884c565b5b8283098314156133f657604051806040016040528082815260200183815250935050505061344d565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47806134255761342461884c565b5b600182089050613386565b604051806040016040528060008152602001600081525093505050505b919050565b60008061345e84614edd565b9050808360ff166001901b116134a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016134a09061891e565b60405180910390fd5b8091505092915050565b600060018260ff1684901c16600114905092915050565b609660009054906101000a900460ff1660ff168160ff1610613521576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613518906189b0565b60405180910390fd5b50565b80609760008460ff1660ff16815260200190815260200160002060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548161ffff021916908361ffff16021790555060408201518160000160066101000a81548161ffff021916908361ffff1602179055509050508160ff167f3ee6fe8d54610244c3e9d3c066ae4aee997884aa28f10616ae821925401318ac826040516135e19190617596565b60405180910390a25050565b609e60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461367d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161367490618a42565b60405180910390fd5b565b60008082149050919050565b60008282841614905092915050565b6000609960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050600081600001549050600160028111156136fa576136f96164e4565b5b8260010160009054906101000a900460ff16600281111561371e5761371d6164e4565b5b1461375e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161375590618ad4565b60405180910390fd5b600061377984609660009054906101000a900460ff16613452565b9050600061378683612ccb565b90506137ab8277ffffffffffffffffffffffffffffffffffffffffffffffff1661367f565b156137eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016137e290618b66565b60405180910390fd5b6138328177ffffffffffffffffffffffffffffffffffffffffffffffff168377ffffffffffffffffffffffffffffffffffffffffffffffff1661368b90919063ffffffff16565b613871576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161386890618c1e565b60405180910390fd5b60006138ba8377ffffffffffffffffffffffffffffffffffffffffffffffff168377ffffffffffffffffffffffffffffffffffffffffffffffff1661500490919063ffffffff16565b90506138c68482615012565b6138e98177ffffffffffffffffffffffffffffffffffffffffffffffff1661367f565b156139eb5760028560010160006101000a81548160ff02191690836002811115613916576139156164e4565b5b02179055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a364f4da886040518263ffffffff1660e01b81526004016139749190615ebe565b600060405180830381600087803b15801561398e57600080fd5b505af11580156139a2573d6000803e3d6000fd5b50505050838773ffffffffffffffffffffffffffffffffffffffff167f396fdcb180cb0fea26928113fb0fd1c3549863f9cd563e6a184f1d578116c8e460405160405180910390a35b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f4e24fe588886040518363ffffffff1660e01b8152600401613a46929190618c3e565b600060405180830381600087803b158015613a6057600080fd5b505af1158015613a74573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663bd29b8cd85886040518363ffffffff1660e01b8152600401613ad3929190618c6e565b600060405180830381600087803b158015613aed57600080fd5b505af1158015613b01573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663bd29b8cd85886040518363ffffffff1660e01b8152600401613b60929190618c6e565b600060405180830381600087803b158015613b7a57600080fd5b505af1158015613b8e573d6000803e3d6000fd5b5050505050505050505050565b6000606460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081606460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000606460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166313542a4e846040518263ffffffff1660e01b8152600401613ce69190615ebe565b602060405180830381865afa158015613d03573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613d279190618cb3565b90506000801b811415613ddf577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663bf79ce588484613d7b87611396565b6040518463ffffffff1660e01b8152600401613d9993929190618e01565b6020604051808303816000875af1158015613db8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ddc9190618cb3565b90505b92915050565b609a60008260200151815260200190815260200160002060009054906101000a900460ff1615613e4a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613e4190618ed2565b60405180910390fd5b4281604001511015613e91576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613e8890618f8a565b60405180910390fd5b6001609a60008360200151815260200190815260200160002060006101000a81548160ff021916908315150217905550613f06609d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16613efc86868686602001518760400151611ded565b836000015161534f565b50505050565b613f14615bb2565b6000613f7386868080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050609660009054906101000a900460ff16613452565b90506000613f8088612ccb565b9050613fa58277ffffffffffffffffffffffffffffffffffffffffffffffff1661367f565b15613fe5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613fdc9061901c565b60405180910390fd5b61402c8177ffffffffffffffffffffffffffffffffffffffffffffffff168377ffffffffffffffffffffffffffffffffffffffffffffffff166154e190919063ffffffff16565b61406b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614062906190d4565b60405180910390fd5b60006140b48377ffffffffffffffffffffffffffffffffffffffffffffffff168377ffffffffffffffffffffffffffffffffffffffffffffffff166154f090919063ffffffff16565b90504260a054609f60008d73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546141049190617eff565b10614144576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161413b90619166565b60405180910390fd5b61414e8982615012565b60016002811115614162576141616164e4565b5b609960008c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff1660028111156141c4576141c36164e4565b5b146143495760405180604001604052808a8152602001600160028111156141ee576141ed6164e4565b5b815250609960008c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000820151816000015560208201518160010160006101000a81548160ff02191690836002811115614265576142646164e4565b5b02179055509050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639926ee7d8b876040518363ffffffff1660e01b81526004016142c8929190619220565b600060405180830381600087803b1580156142e257600080fd5b505af11580156142f6573d6000803e3d6000fd5b505050506143048987612f99565b888a73ffffffffffffffffffffffffffffffffffffffff167fe8e68cef1c3a761ed7be7e8463a375f27f7bc335e51824223cacce636ec5c3fe60405160405180910390a35b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633fb279528b8a8a6040518463ffffffff1660e01b81526004016143a69392919061927d565b600060405180830381600087803b1580156143c057600080fd5b505af11580156143d4573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663255047778b8b8b8b6040518563ffffffff1660e01b815260040161443794939291906192af565b6000604051808303816000875af1158015614456573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061447f919061939b565b856020018660400182905282905250507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1662bff04d8a8a8a6040518463ffffffff1660e01b81526004016144eb93929190619413565b6000604051808303816000875af115801561450a573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906145339190619508565b84600001819052505050509695505050505050565b6000826020015190506000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000015490508173ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161415614607576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016145fe906195c3565b60405180910390fd5b8760ff16846000015160ff1614614653576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161464a90619655565b60405180910390fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635401ed27838b6040518363ffffffff1660e01b81526004016146b0929190619675565b602060405180830381865afa1580156146cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906146f1919061969e565b90506146fd81856154fd565b6bffffffffffffffffffffffff16866bffffffffffffffffffffffff161161475a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161475190619763565b60405180910390fd5b614764888561552c565b6bffffffffffffffffffffffff16816bffffffffffffffffffffffff16106147c1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016147b89061981b565b60405180910390fd5b505050505050505050565b6000806098600084815260200190815260200160002080549050905060005b81811015614882576001818361480191906184ed565b61480b91906184ed565b92508463ffffffff16609860008681526020019081526020016000208463ffffffff168154811061483f5761483e617694565b5b9060005260206000200160000160009054906101000a900463ffffffff1663ffffffff161161486f5750506148be565b808061487a906176f2565b9150506147eb565b506040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016148b5906198d3565b60405180910390fd5b92915050565b6000609660009054906101000a900460ff16905060c060ff168160ff1610614921576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161491890619965565b60405180910390fd5b60018161492e9190619985565b609660006101000a81548160ff021916908360ff16021790555060008190506149578186613524565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663ff694a778286866040518463ffffffff1660e01b81526004016149b493929190619ad9565b600060405180830381600087803b1580156149ce57600080fd5b505af11580156149e2573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166326d941f2826040518263ffffffff1660e01b8152600401614a3f919061695a565b600060405180830381600087803b158015614a5957600080fd5b505af1158015614a6d573d6000803e3d6000fd5b505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166326d941f2826040518263ffffffff1660e01b8152600401614aca919061695a565b600060405180830381600087803b158015614ae457600080fd5b505af1158015614af8573d6000803e3d6000fd5b505050505050505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600073ffffffffffffffffffffffffffffffffffffffff16600060029054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16148015614bb15750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b614bf0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614be790619baf565b60405180910390fd5b806001819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d82604051614c3d9190615d59565b60405180910390a2614c4e826130e0565b5050565b600080600090505b6000831115614c8757600183614c7091906184ed565b831692508080614c7f90619bcf565b915050614c5a565b80915050919050565b600033905090565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163073ffffffffffffffffffffffffffffffffffffffff16148015614d1457507f000000000000000000000000000000000000000000000000000000000000000046145b15614d41577f00000000000000000000000000000000000000000000000000000000000000009050614daf565b614dac7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061555b565b90505b90565b60008282604051602001614dc7929190619c72565b60405160208183030381529060405280519060200120905092915050565b60008060007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4780614e1957614e1861884c565b5b60037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4780614e4a57614e4961884c565b5b867f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd4780614e7a57614e7961884c565b5b888909090890506000614ece827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f527f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47615595565b90508181935093505050915091565b600061010082511115614f25576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614f1c90619d41565b60405180910390fd5b600082511415614f385760009050614fff565b60008083600081518110614f4f57614f4e617694565b5b602001015160f81c60f81b60f81c60ff166001901b91506000600190505b8451811015614ff857848181518110614f8957614f88617694565b5b602001015160f81c60f81b60f81c60ff166001901b9150828211614fe2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614fd990619df9565b60405180910390fd5b818317925080614ff1906176f2565b9050614f6d565b5081925050505b919050565b600081198316905092915050565b6000609860008481526020019081526020016000208054905090506000811415615154576098600084815260200190815260200160002060405180606001604052804363ffffffff168152602001600063ffffffff1681526020018477ffffffffffffffffffffffffffffffffffffffffffffffff168152509080600181540180825580915050600190039060005260206000200160009091909190915060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908377ffffffffffffffffffffffffffffffffffffffffffffffff160217905550505061534a565b60006098600085815260200190815260200160002060018361517691906184ed565b8154811061518757615186617694565b5b9060005260206000200190504363ffffffff168160000160009054906101000a900463ffffffff1663ffffffff16141561520b57828160000160086101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908377ffffffffffffffffffffffffffffffffffffffffffffffff160217905550615348565b438160000160046101000a81548163ffffffff021916908363ffffffff1602179055506098600085815260200190815260200160002060405180606001604052804363ffffffff168152602001600063ffffffff1681526020018577ffffffffffffffffffffffffffffffffffffffffffffffff168152509080600181540180825580915050600190039060005260206000200160009091909190915060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548163ffffffff021916908363ffffffff16021790555060408201518160000160086101000a81548177ffffffffffffffffffffffffffffffffffffffffffffffff021916908377ffffffffffffffffffffffffffffffffffffffffffffffff16021790555050505b505b505050565b615358836156e5565b1561546457631626ba7e60e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168373ffffffffffffffffffffffffffffffffffffffff16631626ba7e84846040518363ffffffff1660e01b81526004016153bf929190618c6e565b602060405180830381865afa1580156153dc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906154009190619e71565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461545f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161545690619f36565b60405180910390fd5b6154dc565b8273ffffffffffffffffffffffffffffffffffffffff166154858383615708565b73ffffffffffffffffffffffffffffffffffffffff16146154db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016154d290619fee565b60405180910390fd5b5b505050565b60008082841614905092915050565b6000818317905092915050565b600061271061ffff16826020015161ffff168461551a919061a00e565b615524919061a054565b905092915050565b600061271061ffff16826040015161ffff1684615549919061a00e565b615553919061a054565b905092915050565b6000838383463060405160200161557695949392919061a085565b6040516020818303038152906040528051906020012090509392505050565b6000806155a0615bd3565b6155a8615bf5565b6020816000600681106155be576155bd617694565b5b6020020181815250506020816001600681106155dd576155dc617694565b5b6020020181815250506020816002600681106155fc576155fb617694565b5b602002018181525050868160036006811061561a57615619617694565b5b602002018181525050858160046006811061563857615637617694565b5b602002018181525050848160056006811061565657615655617694565b5b60200201818152505060208260c08360056107d05a03fa9250826000811461567d5761567f565bfe5b50826156c0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016156b79061a124565b60405180910390fd5b816000600181106156d4576156d3617694565b5b602002015193505050509392505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6000806000615717858561572f565b91509150615724816157b2565b819250505092915050565b6000806041835114156157715760008060006020860151925060408601519150606086015160001a905061576587828585615987565b945094505050506157ab565b6040835114156157a2576000806020850151915060408501519050615797868383615a94565b9350935050506157ab565b60006002915091505b9250929050565b600060048111156157c6576157c56164e4565b5b8160048111156157d9576157d86164e4565b5b14156157e457615984565b600160048111156157f8576157f76164e4565b5b81600481111561580b5761580a6164e4565b5b141561584c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016158439061a190565b60405180910390fd5b600260048111156158605761585f6164e4565b5b816004811115615873576158726164e4565b5b14156158b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016158ab9061a1fc565b60405180910390fd5b600360048111156158c8576158c76164e4565b5b8160048111156158db576158da6164e4565b5b141561591c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016159139061a28e565b60405180910390fd5b60048081111561592f5761592e6164e4565b5b816004811115615942576159416164e4565b5b1415615983576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161597a9061a320565b60405180910390fd5b5b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c11156159c2576000600391509150615a8b565b601b8560ff16141580156159da5750601c8560ff1614155b156159ec576000600491509150615a8b565b600060018787878760405160008152602001604052604051615a11949392919061a340565b6020604051602081039080840390855afa158015615a33573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415615a8257600060019250925050615a8b565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c615ad79190617eff565b9050615ae587828885615987565b935093505050935093915050565b6040518060600160405280600063ffffffff168152602001600063ffffffff168152602001600077ffffffffffffffffffffffffffffffffffffffffffffffff1681525090565b604051806040016040528060008152602001600081525090565b60405180604001604052806000801916815260200160006002811115615b7d57615b7c6164e4565b5b81525090565b6040518060600160405280600063ffffffff168152602001600061ffff168152602001600061ffff1681525090565b60405180606001604052806060815260200160608152602001606081525090565b6040518060200160405280600190602082028036833780820191505090505090565b6040518060c00160405280600690602082028036833780820191505090505090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f840112615c5057615c4f615c2b565b5b8235905067ffffffffffffffff811115615c6d57615c6c615c30565b5b602083019150836020820283011115615c8957615c88615c35565b5b9250929050565b60008060208385031215615ca757615ca6615c21565b5b600083013567ffffffffffffffff811115615cc557615cc4615c26565b5b615cd185828601615c3a565b92509250509250929050565b6000819050919050565b615cf081615cdd565b8114615cfb57600080fd5b50565b600081359050615d0d81615ce7565b92915050565b600060208284031215615d2957615d28615c21565b5b6000615d3784828501615cfe565b91505092915050565b6000819050919050565b615d5381615d40565b82525050565b6000602082019050615d6e6000830184615d4a565b92915050565b600063ffffffff82169050919050565b615d8d81615d74565b8114615d9857600080fd5b50565b600081359050615daa81615d84565b92915050565b615db981615d40565b8114615dc457600080fd5b50565b600081359050615dd681615db0565b92915050565b600080600060608486031215615df557615df4615c21565b5b6000615e0386828701615cfe565b9350506020615e1486828701615d9b565b9250506040615e2586828701615dc7565b9150509250925092565b600077ffffffffffffffffffffffffffffffffffffffffffffffff82169050919050565b615e5c81615e2f565b82525050565b6000602082019050615e776000830184615e53565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000615ea882615e7d565b9050919050565b615eb881615e9d565b82525050565b6000602082019050615ed36000830184615eaf565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b615f2782615ede565b810181811067ffffffffffffffff82111715615f4657615f45615eef565b5b80604052505050565b6000615f59615c17565b9050615f658282615f1e565b919050565b600067ffffffffffffffff821115615f8557615f84615eef565b5b615f8e82615ede565b9050602081019050919050565b82818337600083830152505050565b6000615fbd615fb884615f6a565b615f4f565b905082815260208101848484011115615fd957615fd8615ed9565b5b615fe4848285615f9b565b509392505050565b600082601f83011261600157616000615c2b565b5b8135616011848260208601615faa565b91505092915050565b6000602082840312156160305761602f615c21565b5b600082013567ffffffffffffffff81111561604e5761604d615c26565b5b61605a84828501615fec565b91505092915050565b60006020828403121561607957616078615c21565b5b600061608784828501615dc7565b91505092915050565b600061609b82615e9d565b9050919050565b6160ab81616090565b81146160b657600080fd5b50565b6000813590506160c8816160a2565b92915050565b6000602082840312156160e4576160e3615c21565b5b60006160f2848285016160b9565b91505092915050565b61610481615e9d565b811461610f57600080fd5b50565b600081359050616121816160fb565b92915050565b60006020828403121561613d5761613c615c21565b5b600061614b84828501616112565b91505092915050565b61615d81615cdd565b82525050565b60006020820190506161786000830184616154565b92915050565b60008115159050919050565b6161938161617e565b82525050565b60006020820190506161ae600083018461618a565b92915050565b600080604083850312156161cb576161ca615c21565b5b60006161d985828601615cfe565b92505060206161ea85828601615dc7565b9150509250929050565b6161fd81615d74565b82525050565b61620c81615e2f565b82525050565b60608201600082015161622860008501826161f4565b50602082015161623b60208501826161f4565b50604082015161624e6040850182616203565b50505050565b60006060820190506162696000830184616212565b92915050565b600060ff82169050919050565b6162858161626f565b811461629057600080fd5b50565b6000813590506162a28161627c565b92915050565b6000602082840312156162be576162bd615c21565b5b60006162cc84828501616293565b91505092915050565b6000819050919050565b60006162fa6162f56162f084615e7d565b6162d5565b615e7d565b9050919050565b600061630c826162df565b9050919050565b600061631e82616301565b9050919050565b61632e81616313565b82525050565b60006020820190506163496000830184616325565b92915050565b61635881615d40565b82525050565b604082016000820151616374600085018261634f565b506020820151616387602085018261634f565b50505050565b60006040820190506163a2600083018461635e565b92915050565b60008083601f8401126163be576163bd615c2b565b5b8235905067ffffffffffffffff8111156163db576163da615c30565b5b6020830191508360208202830111156163f7576163f6615c35565b5b9250929050565b60008083601f84011261641457616413615c2b565b5b8235905067ffffffffffffffff81111561643157616430615c30565b5b60208301915083600182028301111561644d5761644c615c35565b5b9250929050565b6000806000806040858703121561646e5761646d615c21565b5b600085013567ffffffffffffffff81111561648c5761648b615c26565b5b616498878288016163a8565b9450945050602085013567ffffffffffffffff8111156164bb576164ba615c26565b5b6164c7878288016163fe565b925092505092959194509250565b6164de81615cdd565b82525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60038110616524576165236164e4565b5b50565b600081905061653582616513565b919050565b600061654582616527565b9050919050565b6165558161653a565b82525050565b60408201600082015161657160008501826164d5565b506020820151616584602085018261654c565b50505050565b600060408201905061659f600083018461655b565b92915050565b600080fd5b600080fd5b600061ffff82169050919050565b6165c6816165af565b81146165d157600080fd5b50565b6000813590506165e3816165bd565b92915050565b6000606082840312156165ff576165fe6165a5565b5b6166096060615f4f565b9050600061661984828501615d9b565b600083015250602061662d848285016165d4565b6020830152506040616641848285016165d4565b60408301525092915050565b6000806080838503121561666457616663615c21565b5b600061667285828601616293565b9250506020616683858286016165e9565b9150509250929050565b600061669882616301565b9050919050565b6166a88161668d565b82525050565b60006020820190506166c3600083018461669f565b92915050565b60006166d482616301565b9050919050565b6166e4816166c9565b82525050565b60006020820190506166ff60008301846166db565b92915050565b60008060006040848603121561671e5761671d615c21565b5b600061672c86828701616112565b935050602084013567ffffffffffffffff81111561674d5761674c615c26565b5b616759868287016163fe565b92509250509250925092565b600067ffffffffffffffff8211156167805761677f615eef565b5b602082029050602081019050919050565b6000604082840312156167a7576167a66165a5565b5b6167b16040615f4f565b905060006167c184828501616293565b60008301525060206167d584828501616112565b60208301525092915050565b60006167f46167ef84616765565b615f4f565b9050808382526020820190506040840283018581111561681757616816615c35565b5b835b81811015616840578061682c8882616791565b845260208401935050604081019050616819565b5050509392505050565b600082601f83011261685f5761685e615c2b565b5b813561686f8482602086016167e1565b91505092915050565b600080600080600060a0868803121561689457616893615c21565b5b60006168a288828901616112565b95505060206168b388828901615cfe565b945050604086013567ffffffffffffffff8111156168d4576168d3615c26565b5b6168e08882890161684a565b93505060606168f188828901615cfe565b925050608061690288828901615dc7565b9150509295509295909350565b600061691a82616301565b9050919050565b61692a8161690f565b82525050565b60006020820190506169456000830184616921565b92915050565b6169548161626f565b82525050565b600060208201905061696f600083018461694b565b92915050565b60008083601f84011261698b5761698a615c2b565b5b8235905067ffffffffffffffff8111156169a8576169a7615c30565b5b6020830191508360018202830111156169c4576169c3615c35565b5b9250929050565b600080fd5b600061010082840312156169e7576169e66169cb565b5b81905092915050565b60008083601f840112616a0657616a05615c2b565b5b8235905067ffffffffffffffff811115616a2357616a22615c30565b5b602083019150836040820283011115616a3f57616a3e615c35565b5b9250929050565b600067ffffffffffffffff821115616a6157616a60615eef565b5b616a6a82615ede565b9050602081019050919050565b6000616a8a616a8584616a46565b615f4f565b905082815260208101848484011115616aa657616aa5615ed9565b5b616ab1848285615f9b565b509392505050565b600082601f830112616ace57616acd615c2b565b5b8135616ade848260208601616a77565b91505092915050565b600060608284031215616afd57616afc6165a5565b5b616b076060615f4f565b9050600082013567ffffffffffffffff811115616b2757616b266165aa565b5b616b3384828501616ab9565b6000830152506020616b4784828501615cfe565b6020830152506040616b5b84828501615dc7565b60408301525092915050565b60008060008060008060008060006101a08a8c031215616b8a57616b89615c21565b5b60008a013567ffffffffffffffff811115616ba857616ba7615c26565b5b616bb48c828d016163fe565b995099505060208a013567ffffffffffffffff811115616bd757616bd6615c26565b5b616be38c828d01616975565b97509750506040616bf68c828d016169d0565b9550506101408a013567ffffffffffffffff811115616c1857616c17615c26565b5b616c248c828d016169f0565b94509450506101608a013567ffffffffffffffff811115616c4857616c47615c26565b5b616c548c828d01616ae7565b9250506101808a013567ffffffffffffffff811115616c7657616c75615c26565b5b616c828c828d01616ae7565b9150509295985092959850929598565b6000616c9d82616301565b9050919050565b616cad81616c92565b82525050565b6000602082019050616cc86000830184616ca4565b92915050565b6000806000806000806101608789031215616cec57616ceb615c21565b5b600087013567ffffffffffffffff811115616d0a57616d09615c26565b5b616d1689828a016163fe565b9650965050602087013567ffffffffffffffff811115616d3957616d38615c26565b5b616d4589828a01616975565b94509450506040616d5889828a016169d0565b92505061014087013567ffffffffffffffff811115616d7a57616d79615c26565b5b616d8689828a01616ae7565b9150509295509295509295565b600067ffffffffffffffff821115616dae57616dad615eef565b5b602082029050602081019050919050565b6000616dd2616dcd84616d93565b615f4f565b90508083825260208201905060208402830185811115616df557616df4615c35565b5b835b81811015616e1e5780616e0a8882615cfe565b845260208401935050602081019050616df7565b5050509392505050565b600082601f830112616e3d57616e3c615c2b565b5b8135616e4d848260208601616dbf565b91505092915050565b60008060408385031215616e6d57616e6c615c21565b5b6000616e7b85828601615d9b565b925050602083013567ffffffffffffffff811115616e9c57616e9b615c26565b5b616ea885828601616e28565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6000616eea83836161f4565b60208301905092915050565b6000602082019050919050565b6000616f0e82616eb2565b616f188185616ebd565b9350616f2383616ece565b8060005b83811015616f54578151616f3b8882616ede565b9750616f4683616ef6565b925050600181019050616f27565b5085935050505092915050565b60006020820190508181036000830152616f7b8184616f03565b905092915050565b60008060208385031215616f9a57616f99615c21565b5b600083013567ffffffffffffffff811115616fb857616fb7615c26565b5b616fc4858286016163fe565b92509250509250929050565b60006bffffffffffffffffffffffff82169050919050565b616ff181616fd0565b8114616ffc57600080fd5b50565b60008135905061700e81616fe8565b92915050565b600067ffffffffffffffff82111561702f5761702e615eef565b5b602082029050602081019050919050565b600061704b82615e9d565b9050919050565b61705b81617040565b811461706657600080fd5b50565b60008135905061707881617052565b92915050565b600060408284031215617094576170936165a5565b5b61709e6040615f4f565b905060006170ae84828501617069565b60008301525060206170c284828501616fff565b60208301525092915050565b60006170e16170dc84617014565b615f4f565b9050808382526020820190506040840283018581111561710457617103615c35565b5b835b8181101561712d5780617119888261707e565b845260208401935050604081019050617106565b5050509392505050565b600082601f83011261714c5761714b615c2b565b5b813561715c8482602086016170ce565b91505092915050565b600080600060a0848603121561717e5761717d615c21565b5b600061718c868287016165e9565b935050606061719d86828701616fff565b925050608084013567ffffffffffffffff8111156171be576171bd615c26565b5b6171ca86828701617137565b9150509250925092565b600067ffffffffffffffff8211156171ef576171ee615eef565b5b602082029050602081019050919050565b600061721361720e846171d4565b615f4f565b9050808382526020820190506060840283018581111561723657617235615c35565b5b835b8181101561725f578061724b88826165e9565b845260208401935050606081019050617238565b5050509392505050565b600082601f83011261727e5761727d615c2b565b5b813561728e848260208601617200565b91505092915050565b600067ffffffffffffffff8211156172b2576172b1615eef565b5b602082029050602081019050919050565b60006172d66172d184617297565b615f4f565b905080838252602082019050602084028301858111156172f9576172f8615c35565b5b835b81811015617322578061730e8882616fff565b8452602084019350506020810190506172fb565b5050509392505050565b600082601f83011261734157617340615c2b565b5b81356173518482602086016172c3565b91505092915050565b600067ffffffffffffffff82111561737557617374615eef565b5b602082029050602081019050919050565b60006173996173948461735a565b615f4f565b905080838252602082019050602084028301858111156173bc576173bb615c35565b5b835b8181101561740357803567ffffffffffffffff8111156173e1576173e0615c2b565b5b8086016173ee8982617137565b855260208501945050506020810190506173be565b5050509392505050565b600082601f83011261742257617421615c2b565b5b8135617432848260208601617386565b91505092915050565b600080600080600080600080610100898b03121561745c5761745b615c21565b5b600061746a8b828c01616112565b985050602061747b8b828c01616112565b975050604061748c8b828c01616112565b965050606061749d8b828c016160b9565b95505060806174ae8b828c01615dc7565b94505060a089013567ffffffffffffffff8111156174cf576174ce615c26565b5b6174db8b828c01617269565b93505060c089013567ffffffffffffffff8111156174fc576174fb615c26565b5b6175088b828c0161732c565b92505060e089013567ffffffffffffffff81111561752957617528615c26565b5b6175358b828c0161740d565b9150509295985092959890939650565b61754e816165af565b82525050565b60608201600082015161756a60008501826161f4565b50602082015161757d6020850182617545565b5060408201516175906040850182617545565b50505050565b60006060820190506175ab6000830184617554565b92915050565b60006175bc82616301565b9050919050565b6175cc816175b1565b82525050565b60006020820190506175e760008301846175c3565b92915050565b6175f68161653a565b82525050565b600060208201905061761160008301846175ed565b92915050565b600082825260208201905092915050565b7f5061757361626c653a20696e6465782069732070617573656400000000000000600082015250565b600061765e601983617617565b915061766982617628565b602082019050919050565b6000602082019050818103600083015261768d81617651565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006176fd82615d40565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156177305761772f6176c3565b5b600182019050919050565b7f526567436f6f72642e67657451756f72756d4269746d61704174426c6f636b4e60008201527f756d6265724279496e6465783a2071756f72756d4269746d617055706461746560208201527f2069732066726f6d20616674657220626c6f636b4e756d626572000000000000604082015250565b60006177bd605a83617617565b91506177c88261773b565b606082019050919050565b600060208201905081810360008301526177ec816177b0565b9050919050565b7f526567436f6f72642e67657451756f72756d4269746d61704174426c6f636b4e60008201527f756d6265724279496e6465783a2071756f72756d4269746d617055706461746560208201527f2069732066726f6d206265666f726520626c6f636b4e756d6265720000000000604082015250565b6000617875605b83617617565b9150617880826177f3565b606082019050919050565b600060208201905081810360008301526178a481617868565b9050919050565b7f526567436f6f72642e757064617465536f636b65743a206f70657261746f722060008201527f6e6f742072656769737465726564000000000000000000000000000000000000602082015250565b6000617907602e83617617565b9150617912826178ab565b604082019050919050565b60006020820190508181036000830152617936816178fa565b9050919050565b60008151905061794c816160fb565b92915050565b60006020828403121561796857617967615c21565b5b60006179768482850161793d565b91505092915050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320756e70617573657200000000000000000000000000000000000000000000602082015250565b60006179db602a83617617565b91506179e68261797f565b604082019050919050565b60006020820190508181036000830152617a0a816179ce565b9050919050565b617a1a8161617e565b8114617a2557600080fd5b50565b600081519050617a3781617a11565b92915050565b600060208284031215617a5357617a52615c21565b5b6000617a6184828501617a28565b91505092915050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320706175736572000000000000000000000000000000000000000000000000602082015250565b6000617ac6602883617617565b9150617ad182617a6a565b604082019050919050565b60006020820190508181036000830152617af581617ab9565b9050919050565b7f5061757361626c652e70617573653a20696e76616c696420617474656d70742060008201527f746f20756e70617573652066756e6374696f6e616c6974790000000000000000602082015250565b6000617b58603883617617565b9150617b6382617afc565b604082019050919050565b60006020820190508181036000830152617b8781617b4b565b9050919050565b6000604082019050617ba36000830185616154565b617bb06020830184615eaf565b9392505050565b7f526567436f6f72642e7570646174654f70657261746f7273466f7251756f727560008201527f6d3a20696e707574206c656e677468206d69736d617463680000000000000000602082015250565b6000617c13603883617617565b9150617c1e82617bb7565b604082019050919050565b60006020820190508181036000830152617c4281617c06565b9050919050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112617c7557617c74617c49565b5b80840192508235915067ffffffffffffffff821115617c9757617c96617c4e565b5b602083019250602082023603831315617cb357617cb2617c53565b5b509250929050565b600081519050617cca81615d84565b92915050565b600060208284031215617ce657617ce5615c21565b5b6000617cf484828501617cbb565b91505092915050565b7f526567436f6f72642e7570646174654f70657261746f7273466f7251756f727560008201527f6d3a206e756d626572206f662075706461746564206f70657261746f7273206460208201527f6f6573206e6f74206d617463682071756f72756d20746f74616c000000000000604082015250565b6000617d7f605a83617617565b9150617d8a82617cfd565b606082019050919050565b60006020820190508181036000830152617dae81617d72565b9050919050565b7f526567436f6f72642e7570646174654f70657261746f7273466f7251756f727560008201527f6d3a206f70657261746f72206e6f7420696e2071756f72756d00000000000000602082015250565b6000617e11603983617617565b9150617e1c82617db5565b604082019050919050565b60006020820190508181036000830152617e4081617e04565b9050919050565b7f526567436f6f72642e7570646174654f70657261746f7273466f7251756f727560008201527f6d3a206f70657261746f7273206172726179206d75737420626520736f72746560208201527f6420696e20617363656e64696e672061646472657373206f7264657200000000604082015250565b6000617ec9605c83617617565b9150617ed482617e47565b606082019050919050565b60006020820190508181036000830152617ef881617ebc565b9050919050565b6000617f0a82615d40565b9150617f1583615d40565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115617f4a57617f496176c3565b5b828201905092915050565b600080fd5b600080fd5b60008085851115617f7357617f72617f55565b5b83861115617f8457617f83617f5a565b5b6001850283019150848603905094509492505050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b617fcf8161626f565b82525050565b617fde81615e9d565b82525050565b604082016000820151617ffa6000850182617fc6565b50602082015161800d6020850182617fd5565b50505050565b600061801f8383617fe4565b60408301905092915050565b6000602082019050919050565b600061804382617f9a565b61804d8185617fa5565b935061805883617fb6565b8060005b838110156180895781516180708882618013565b975061807b8361802b565b92505060018101905061805c565b5085935050505092915050565b600060c0820190506180ab6000830189616154565b6180b86020830188615eaf565b6180c56040830187616154565b81810360608301526180d78186618038565b90506180e66080830185616154565b6180f360a0830184615d4a565b979650505050505050565b7f526567436f6f72642e72656769737465724f70657261746f725769746843687560008201527f726e3a20696e707574206c656e677468206d69736d6174636800000000000000602082015250565b600061815a603983617617565b9150618165826180fe565b604082019050919050565b600060208201905081810360008301526181898161814d565b9050919050565b6000604082840312156181a6576181a5615c21565b5b60006181b484828501616791565b91505092915050565b7f526567436f6f72642e72656769737465724f70657261746f723a206f7065726160008201527f746f7220636f756e742065786365656473206d6178696d756d00000000000000602082015250565b6000618219603983617617565b9150618224826181bd565b604082019050919050565b600060208201905081810360008301526182488161820c565b9050919050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b60006182ab602e83617617565b91506182b68261824f565b604082019050919050565b600060208201905081810360008301526182da8161829e565b9050919050565b7f526567436f6f72642e696e697469616c697a653a20696e707574206c656e677460008201527f68206d69736d6174636800000000000000000000000000000000000000000000602082015250565b600061833d602a83617617565b9150618348826182e1565b604082019050919050565b6000602082019050818103600083015261836c81618330565b9050919050565b6000819050919050565b600061839861839361838e84618373565b6162d5565b61626f565b9050919050565b6183a88161837d565b82525050565b60006020820190506183c3600083018461839f565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000618425602683617617565b9150618430826183c9565b604082019050919050565b6000602082019050818103600083015261845481618418565b9050919050565b7f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060008201527f7420746f2070617573652066756e6374696f6e616c6974790000000000000000602082015250565b60006184b7603883617617565b91506184c28261845b565b604082019050919050565b600060208201905081810360008301526184e6816184aa565b9050919050565b60006184f882615d40565b915061850383615d40565b925082821015618516576185156176c3565b5b828203905092915050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561855b578082015181840152602081019050618540565b8381111561856a576000848401525b50505050565b600061857b82618521565b618585818561852c565b935061859581856020860161853d565b61859e81615ede565b840191505092915050565b60006060820190506185be6000830186615eaf565b6185cb6020830185616154565b81810360408301526185dd8184618570565b9050949350505050565b6185f081615e2f565b81146185fb57600080fd5b50565b60008151905061860d816185e7565b92915050565b60006020828403121561862957618628615c21565b5b6000618637848285016185fe565b91505092915050565b600081519050919050565b600061865682618640565b6186608185617617565b935061867081856020860161853d565b61867981615ede565b840191505092915050565b60006040820190506186996000830185616154565b81810360208301526186ab818461864b565b90509392505050565b600060208201905081810360008301526186ce818461864b565b905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061870c602083617617565b9150618717826186d6565b602082019050919050565b6000602082019050818103600083015261873b816186ff565b9050919050565b7f5061757361626c652e5f73657450617573657252656769737472793a206e657760008201527f50617573657252656769737472792063616e6e6f7420626520746865207a657260208201527f6f20616464726573730000000000000000000000000000000000000000000000604082015250565b60006187c4604983617617565b91506187cf82618742565b606082019050919050565b600060208201905081810360008301526187f3816187b7565b9050919050565b600060408201905061880f6000830185616921565b61881c6020830184616921565b9392505050565b60006040820190506188386000830185615eaf565b6188456020830184615eaf565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061888682615d40565b915061889183615d40565b9250826188a1576188a061884c565b5b828206905092915050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206269746d61702065786365656473206d61782076616c756500602082015250565b6000618908603f83617617565b9150618913826188ac565b604082019050919050565b60006020820190508181036000830152618937816188fb565b9050919050565b7f526567436f6f72642e71756f72756d4578697374733a2071756f72756d20646f60008201527f6573206e6f742065786973740000000000000000000000000000000000000000602082015250565b600061899a602c83617617565b91506189a58261893e565b604082019050919050565b600060208201905081810360008301526189c98161898d565b9050919050565b7f526567436f6f72642e6f6e6c79456a6563746f723a2063616c6c65722069732060008201527f6e6f742074686520656a6563746f720000000000000000000000000000000000602082015250565b6000618a2c602f83617617565b9150618a37826189d0565b604082019050919050565b60006020820190508181036000830152618a5b81618a1f565b9050919050565b7f526567436f6f72642e5f646572656769737465724f70657261746f723a206f7060008201527f657261746f72206973206e6f7420726567697374657265640000000000000000602082015250565b6000618abe603883617617565b9150618ac982618a62565b604082019050919050565b60006020820190508181036000830152618aed81618ab1565b9050919050565b7f526567436f6f72642e5f646572656769737465724f70657261746f723a20626960008201527f746d61702063616e6e6f74206265203000000000000000000000000000000000602082015250565b6000618b50603083617617565b9150618b5b82618af4565b604082019050919050565b60006020820190508181036000830152618b7f81618b43565b9050919050565b7f526567436f6f72642e5f646572656769737465724f70657261746f723a206f7060008201527f657261746f72206973206e6f74207265676973746572656420666f722071756f60208201527f72756d7300000000000000000000000000000000000000000000000000000000604082015250565b6000618c08604483617617565b9150618c1382618b86565b606082019050919050565b60006020820190508181036000830152618c3781618bfb565b9050919050565b6000604082019050618c536000830185615eaf565b8181036020830152618c658184618570565b90509392505050565b6000604082019050618c836000830185616154565b8181036020830152618c958184618570565b90509392505050565b600081519050618cad81615ce7565b92915050565b600060208284031215618cc957618cc8615c21565b5b6000618cd784828501618c9e565b91505092915050565b600082905092915050565b6000618cfa6020840184615dc7565b905092915050565b60408201618d136000830183618ceb565b618d20600085018261634f565b50618d2e6020830183618ceb565b618d3b602085018261634f565b50505050565b600082905092915050565b600082905092915050565b618d6360408383615f9b565b5050565b60808201618d786000830183618d4c565b618d856000850182618d57565b50618d936040830183618d4c565b618da06040850182618d57565b50505050565b6101008201618db86000830183618ce0565b618dc56000850182618d02565b50618dd36040830183618ce0565b618de06040850182618d02565b50618dee6080830183618d41565b618dfb6080850182618d67565b50505050565b600061016082019050618e176000830186615eaf565b618e246020830185618da6565b618e3261012083018461635e565b949350505050565b7f526567436f6f72642e5f766572696679436875726e417070726f76657253696760008201527f6e61747572653a20636875726e417070726f7665722073616c7420616c72656160208201527f6479207573656400000000000000000000000000000000000000000000000000604082015250565b6000618ebc604783617617565b9150618ec782618e3a565b606082019050919050565b60006020820190508181036000830152618eeb81618eaf565b9050919050565b7f526567436f6f72642e5f766572696679436875726e417070726f76657253696760008201527f6e61747572653a20636875726e417070726f766572207369676e61747572652060208201527f6578706972656400000000000000000000000000000000000000000000000000604082015250565b6000618f74604783617617565b9150618f7f82618ef2565b606082019050919050565b60006020820190508181036000830152618fa381618f67565b9050919050565b7f526567436f6f72642e5f72656769737465724f70657261746f723a206269746d60008201527f61702063616e6e6f742062652030000000000000000000000000000000000000602082015250565b6000619006602e83617617565b915061901182618faa565b604082019050919050565b6000602082019050818103600083015261903581618ff9565b9050919050565b7f526567436f6f72642e5f72656769737465724f70657261746f723a206f70657260008201527f61746f7220616c7265616479207265676973746572656420666f7220736f6d6560208201527f2071756f72756d73000000000000000000000000000000000000000000000000604082015250565b60006190be604883617617565b91506190c98261903c565b606082019050919050565b600060208201905081810360008301526190ed816190b1565b9050919050565b7f526567436f6f72642e5f72656769737465724f70657261746f723a206f70657260008201527f61746f722063616e6e6f74207265726567697374657220796574000000000000602082015250565b6000619150603a83617617565b915061915b826190f4565b604082019050919050565b6000602082019050818103600083015261917f81619143565b9050919050565b600082825260208201905092915050565b60006191a282618521565b6191ac8185619186565b93506191bc81856020860161853d565b6191c581615ede565b840191505092915050565b600060608301600083015184820360008601526191ed8282619197565b915050602083015161920260208601826164d5565b506040830151619215604086018261634f565b508091505092915050565b60006040820190506192356000830185615eaf565b818103602083015261924781846191d0565b90509392505050565b600061925c838561852c565b9350619269838584615f9b565b61927283615ede565b840190509392505050565b60006040820190506192926000830186615eaf565b81810360208301526192a5818486619250565b9050949350505050565b60006060820190506192c46000830187615eaf565b6192d16020830186616154565b81810360408301526192e4818486619250565b905095945050505050565b6000815190506192fe81616fe8565b92915050565b600061931761931284617297565b615f4f565b9050808382526020820190506020840283018581111561933a57619339615c35565b5b835b81811015619363578061934f88826192ef565b84526020840193505060208101905061933c565b5050509392505050565b600082601f83011261938257619381615c2b565b5b8151619392848260208601619304565b91505092915050565b600080604083850312156193b2576193b1615c21565b5b600083015167ffffffffffffffff8111156193d0576193cf615c26565b5b6193dc8582860161936d565b925050602083015167ffffffffffffffff8111156193fd576193fc615c26565b5b6194098582860161936d565b9150509250929050565b60006040820190506194286000830186616154565b818103602083015261943b818486619250565b9050949350505050565b600067ffffffffffffffff8211156194605761945f615eef565b5b602082029050602081019050919050565b600061948461947f84619445565b615f4f565b905080838252602082019050602084028301858111156194a7576194a6615c35565b5b835b818110156194d057806194bc8882617cbb565b8452602084019350506020810190506194a9565b5050509392505050565b600082601f8301126194ef576194ee615c2b565b5b81516194ff848260208601619471565b91505092915050565b60006020828403121561951e5761951d615c21565b5b600082015167ffffffffffffffff81111561953c5761953b615c26565b5b619548848285016194da565b91505092915050565b7f526567436f6f72642e5f76616c6964617465436875726e3a2063616e6e6f742060008201527f636875726e2073656c6600000000000000000000000000000000000000000000602082015250565b60006195ad602a83617617565b91506195b882619551565b604082019050919050565b600060208201905081810360008301526195dc816195a0565b9050919050565b7f526567436f6f72642e5f76616c6964617465436875726e3a2071756f72756d4e60008201527f756d626572206e6f74207468652073616d65206173207369676e656400000000602082015250565b600061963f603c83617617565b915061964a826195e3565b604082019050919050565b6000602082019050818103600083015261966e81619632565b9050919050565b600060408201905061968a6000830185616154565b619697602083018461694b565b9392505050565b6000602082840312156196b4576196b3615c21565b5b60006196c2848285016192ef565b91505092915050565b7f526567436f6f72642e5f76616c6964617465436875726e3a20696e636f6d696e60008201527f67206f70657261746f722068617320696e73756666696369656e74207374616b60208201527f6520666f7220636875726e000000000000000000000000000000000000000000604082015250565b600061974d604b83617617565b9150619758826196cb565b606082019050919050565b6000602082019050818103600083015261977c81619740565b9050919050565b7f526567436f6f72642e5f76616c6964617465436875726e3a2063616e6e6f742060008201527f6b69636b206f70657261746f722077697468206d6f7265207468616e206b696360208201527f6b424950734f66546f74616c5374616b65000000000000000000000000000000604082015250565b6000619805605183617617565b915061981082619783565b606082019050919050565b60006020820190508181036000830152619834816197f8565b9050919050565b7f526567436f6f72642e67657451756f72756d4269746d6170496e64657841744260008201527f6c6f636b4e756d6265723a206e6f206269746d61702075706461746520666f7560208201527f6e6420666f72206f70657261746f7220617420626c6f636b4e756d6265720000604082015250565b60006198bd605e83617617565b91506198c88261983b565b606082019050919050565b600060208201905081810360008301526198ec816198b0565b9050919050565b7f526567436f6f72642e63726561746551756f72756d3a206d61782071756f727560008201527f6d73207265616368656400000000000000000000000000000000000000000000602082015250565b600061994f602a83617617565b915061995a826198f3565b604082019050919050565b6000602082019050818103600083015261997e81619942565b9050919050565b60006199908261626f565b915061999b8361626f565b92508260ff038211156199b1576199b06176c3565b5b828201905092915050565b6199c581616fd0565b82525050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6000619a0282616301565b9050919050565b619a12816199f7565b82525050565b619a2181616fd0565b82525050565b604082016000820151619a3d6000850182619a09565b506020820151619a506020850182619a18565b50505050565b6000619a628383619a27565b60408301905092915050565b6000602082019050919050565b6000619a86826199cb565b619a9081856199d6565b9350619a9b836199e7565b8060005b83811015619acc578151619ab38882619a56565b9750619abe83619a6e565b925050600181019050619a9f565b5085935050505092915050565b6000606082019050619aee600083018661694b565b619afb60208301856199bc565b8181036040830152619b0d8184619a7b565b9050949350505050565b7f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960008201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c60208201527f6564206f6e636500000000000000000000000000000000000000000000000000604082015250565b6000619b99604783617617565b9150619ba482619b17565b606082019050919050565b60006020820190508181036000830152619bc881619b8c565b9050919050565b6000619bda826165af565b915061ffff821415619bef57619bee6176c3565b5b600182019050919050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b6000619c3b600283619bfa565b9150619c4682619c05565b600282019050919050565b6000819050919050565b619c6c619c6782615cdd565b619c51565b82525050565b6000619c7d82619c2e565b9150619c898285619c5b565b602082019150619c998284619c5b565b6020820191508190509392505050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f7264657265644279746573417272617920697320746f6f2060208201527f6c6f6e6700000000000000000000000000000000000000000000000000000000604082015250565b6000619d2b604483617617565b9150619d3682619ca9565b606082019050919050565b60006020820190508181036000830152619d5a81619d1e565b9050919050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f72646572656442797465734172726179206973206e6f742060208201527f6f72646572656400000000000000000000000000000000000000000000000000604082015250565b6000619de3604783617617565b9150619dee82619d61565b606082019050919050565b60006020820190508181036000830152619e1281619dd6565b9050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b619e4e81619e19565b8114619e5957600080fd5b50565b600081519050619e6b81619e45565b92915050565b600060208284031215619e8757619e86615c21565b5b6000619e9584828501619e5c565b91505092915050565b7f454950313237315369676e61747572655574696c732e636865636b5369676e6160008201527f747572655f454950313237313a2045524331323731207369676e61747572652060208201527f766572696669636174696f6e206661696c656400000000000000000000000000604082015250565b6000619f20605383617617565b9150619f2b82619e9e565b606082019050919050565b60006020820190508181036000830152619f4f81619f13565b9050919050565b7f454950313237315369676e61747572655574696c732e636865636b5369676e6160008201527f747572655f454950313237313a207369676e6174757265206e6f742066726f6d60208201527f207369676e657200000000000000000000000000000000000000000000000000604082015250565b6000619fd8604783617617565b9150619fe382619f56565b606082019050919050565b6000602082019050818103600083015261a00781619fcb565b9050919050565b600061a01982616fd0565b915061a02483616fd0565b9250816bffffffffffffffffffffffff048311821515161561a0495761a0486176c3565b5b828202905092915050565b600061a05f82616fd0565b915061a06a83616fd0565b92508261a07a5761a07961884c565b5b828204905092915050565b600060a08201905061a09a6000830188616154565b61a0a76020830187616154565b61a0b46040830186616154565b61a0c16060830185615d4a565b61a0ce6080830184615eaf565b9695505050505050565b7f424e3235342e6578704d6f643a2063616c6c206661696c757265000000000000600082015250565b600061a10e601a83617617565b915061a1198261a0d8565b602082019050919050565b6000602082019050818103600083015261a13d8161a101565b9050919050565b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b600061a17a601883617617565b915061a1858261a144565b602082019050919050565b6000602082019050818103600083015261a1a98161a16d565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b600061a1e6601f83617617565b915061a1f18261a1b0565b602082019050919050565b6000602082019050818103600083015261a2158161a1d9565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061a278602283617617565b915061a2838261a21c565b604082019050919050565b6000602082019050818103600083015261a2a78161a26b565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b600061a30a602283617617565b915061a3158261a2ae565b604082019050919050565b6000602082019050818103600083015261a3398161a2fd565b9050919050565b600060808201905061a3556000830187616154565b61a362602083018661694b565b61a36f6040830185616154565b61a37c6060830184616154565b9594505050505056fea2646970667358221220bcd60a898a30942eeceaa0b690e952b2725c942a1471150379831535e9cf0e8b64736f6c634300080c0033",
}

// ContractRegistryCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractRegistryCoordinatorMetaData.ABI instead.
var ContractRegistryCoordinatorABI = ContractRegistryCoordinatorMetaData.ABI

// ContractRegistryCoordinatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractRegistryCoordinatorMetaData.Bin instead.
var ContractRegistryCoordinatorBin = ContractRegistryCoordinatorMetaData.Bin

// DeployContractRegistryCoordinator deploys a new Ethereum contract, binding an instance of ContractRegistryCoordinator to it.
func DeployContractRegistryCoordinator(auth *bind.TransactOpts, backend bind.ContractBackend, _serviceManager common.Address, _stakeRegistry common.Address, _blsApkRegistry common.Address, _indexRegistry common.Address, _socketRegistry common.Address) (common.Address, *types.Transaction, *ContractRegistryCoordinator, error) {
	parsed, err := ContractRegistryCoordinatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractRegistryCoordinatorBin), backend, _serviceManager, _stakeRegistry, _blsApkRegistry, _indexRegistry, _socketRegistry)
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

// SocketRegistry is a free data retrieval call binding the contract method 0xea32afae.
//
// Solidity: function socketRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCaller) SocketRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractRegistryCoordinator.contract.Call(opts, &out, "socketRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SocketRegistry is a free data retrieval call binding the contract method 0xea32afae.
//
// Solidity: function socketRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorSession) SocketRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.SocketRegistry(&_ContractRegistryCoordinator.CallOpts)
}

// SocketRegistry is a free data retrieval call binding the contract method 0xea32afae.
//
// Solidity: function socketRegistry() view returns(address)
func (_ContractRegistryCoordinator *ContractRegistryCoordinatorCallerSession) SocketRegistry() (common.Address, error) {
	return _ContractRegistryCoordinator.Contract.SocketRegistry(&_ContractRegistryCoordinator.CallOpts)
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
