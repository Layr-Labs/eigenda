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
}

// ContractRegistryCoordinatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractRegistryCoordinatorMetaData.ABI instead.
var ContractRegistryCoordinatorABI = ContractRegistryCoordinatorMetaData.ABI

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
