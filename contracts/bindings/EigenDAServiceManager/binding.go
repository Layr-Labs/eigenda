// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractEigenDAServiceManager

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

// IBLSSignatureCheckerQuorumStakeTotals is an auto generated low-level Go binding around an user-defined struct.
type IBLSSignatureCheckerQuorumStakeTotals struct {
	SignedStakeForQuorum []*big.Int
	TotalStakeForQuorum  []*big.Int
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

// ContractEigenDAServiceManagerMetaData contains all meta data concerning the ContractEigenDAServiceManager contract.
var ContractEigenDAServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"__avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"__rewardsCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRewardsCoordinator\"},{\"name\":\"__registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"__stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"},{\"name\":\"__eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"__eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"},{\"name\":\"__paymentVault\",\"type\":\"address\",\"internalType\":\"contractIPaymentVault\"},{\"name\":\"__eigenDADisperserRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDADisperserRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STORE_DURATION_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"THRESHOLD_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkSignatures\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"components\":[{\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAVSRewardsSubmission\",\"inputs\":[{\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createOperatorDirectedAVSRewardsSubmission\",\"inputs\":[{\"name\":\"operatorDirectedRewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.OperatorDirectedRewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"operatorRewards\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.OperatorReward[]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDADisperserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDADisperserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchConfirmers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_rewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isBatchConfirmer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestServeUntilBlock\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paymentVault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPaymentVault\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardsInitiator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBatchConfirmer\",\"inputs\":[{\"name\":\"_batchConfirmer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setClaimerFor\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsInitiator\",\"inputs\":[{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStaleStakesForbidden\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staleStakesForbidden\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trySignatureAndApkVerification\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"apk\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"pairingSuccessful\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"siganatureIsValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchConfirmerStatusChanged\",\"inputs\":[{\"name\":\"batchConfirmer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultSecurityThresholdsV2Updated\",\"inputs\":[{\"name\":\"previousDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]},{\"name\":\"newDefaultSecurityThresholdsV2\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumAdversaryThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumAdversaryThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumConfirmationThresholdPercentagesUpdated\",\"inputs\":[{\"name\":\"previousQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumConfirmationThresholdPercentages\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuorumNumbersRequiredUpdated\",\"inputs\":[{\"name\":\"previousQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"newQuorumNumbersRequired\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsInitiatorUpdated\",\"inputs\":[{\"name\":\"prevRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleStakesForbiddenUpdate\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VersionedBlobParamsAdded\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"},{\"name\":\"versionedBlobParams\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"anonymous\":false}]",
	Bin: "0x6102006040523480156200001257600080fd5b506040516200abf03803806200abf0833981810160405281019062000038919062000779565b8588888888838383838c8c8c8c8373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff1660c08173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660e08173ffffffffffffffffffffffffffffffffffffffff1681525050505050508373ffffffffffffffffffffffffffffffffffffffff166101008173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff166101208173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff166101408173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff166101608173ffffffffffffffffffffffffffffffffffffffff168152505050505050620002016200044e60201b60201c565b505050508073ffffffffffffffffffffffffffffffffffffffff166101808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000286573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002ac919062000842565b73ffffffffffffffffffffffffffffffffffffffff166101a08173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff16635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200032c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003529190620008b9565b73ffffffffffffffffffffffffffffffffffffffff166101c08173ffffffffffffffffffffffffffffffffffffffff16815250506101a05173ffffffffffffffffffffffffffffffffffffffff1663df5cf7236040518163ffffffff1660e01b8152600401602060405180830381865afa158015620003d5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620003fb919062000930565b73ffffffffffffffffffffffffffffffffffffffff166101e08173ffffffffffffffffffffffffffffffffffffffff168152505050620004406200044e60201b60201c565b505050505050505062000a46565b603260019054906101000a900460ff1615620004a1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200049890620009e9565b60405180910390fd5b60ff8016603260009054906101000a900460ff1660ff161015620005165760ff603260006101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff6040516200050d919062000a29565b60405180910390a15b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200054a826200051d565b9050919050565b60006200055e826200053d565b9050919050565b620005708162000551565b81146200057c57600080fd5b50565b600081519050620005908162000565565b92915050565b6000620005a3826200053d565b9050919050565b620005b58162000596565b8114620005c157600080fd5b50565b600081519050620005d581620005aa565b92915050565b6000620005e8826200053d565b9050919050565b620005fa81620005db565b81146200060657600080fd5b50565b6000815190506200061a81620005ef565b92915050565b60006200062d826200053d565b9050919050565b6200063f8162000620565b81146200064b57600080fd5b50565b6000815190506200065f8162000634565b92915050565b600062000672826200053d565b9050919050565b620006848162000665565b81146200069057600080fd5b50565b600081519050620006a48162000679565b92915050565b6000620006b7826200053d565b9050919050565b620006c981620006aa565b8114620006d557600080fd5b50565b600081519050620006e981620006be565b92915050565b6000620006fc826200053d565b9050919050565b6200070e81620006ef565b81146200071a57600080fd5b50565b6000815190506200072e8162000703565b92915050565b600062000741826200053d565b9050919050565b620007538162000734565b81146200075f57600080fd5b50565b600081519050620007738162000748565b92915050565b600080600080600080600080610100898b0312156200079d576200079c62000518565b5b6000620007ad8b828c016200057f565b9850506020620007c08b828c01620005c4565b9750506040620007d38b828c0162000609565b9650506060620007e68b828c016200064e565b9550506080620007f98b828c0162000693565b94505060a06200080c8b828c01620006d8565b93505060c06200081f8b828c016200071d565b92505060e0620008328b828c0162000762565b9150509295985092959890939650565b6000602082840312156200085b576200085a62000518565b5b60006200086b848285016200064e565b91505092915050565b600062000881826200053d565b9050919050565b620008938162000874565b81146200089f57600080fd5b50565b600081519050620008b38162000888565b92915050565b600060208284031215620008d257620008d162000518565b5b6000620008e284828501620008a2565b91505092915050565b6000620008f8826200053d565b9050919050565b6200090a81620008eb565b81146200091657600080fd5b50565b6000815190506200092a81620008ff565b92915050565b60006020828403121562000949576200094862000518565b5b6000620009598482850162000919565b91505092915050565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b6000620009d160278362000962565b9150620009de8262000973565b604082019050919050565b6000602082019050818103600083015262000a0481620009c2565b9050919050565b600060ff82169050919050565b62000a238162000a0b565b82525050565b600060208201905062000a40600083018462000a18565b92915050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516101c0516101e05161a03162000bbf60003960008181611f0001526134440152600081816117da01526120b201526000818161180b0152818161223d01526124160152600081816115260152818161185701528181611bb101528181611d3a0152611fae0152600081816111dc01528181611338015281816113df01528181613606015281816137ad0152613854015260008181610f5801528181610ff6015281816110b801528181612dc0015281816131c4015281816135070152613712015260008181612ee7015281816130a00152818161313101528181613dd60152613e8e01526000818161183101528181612e4e01528181613252015261330801526000613a86015260006139bf015260006126a301526000818161097601528181610c9e01528181610eb301528181612cd7015281816133ab0152818161346a015281816139e50152613b47015261a0316000f3fe608060405234801561001057600080fd5b50600436106102f15760003560e01c80638687feae1161019d578063e481af9d116100e9578063ef024458116100a2578063f8c668141161007c578063f8c66814146108fe578063fabc1cbc1461091c578063fc299dee14610938578063fce36c7d14610956576102f1565b8063ef024458146108a8578063f1220983146108c6578063f2fde38b146108e2576102f1565b8063e481af9d146107be578063eaefd27d146107dc578063eccbbfc91461080c578063ed3916f71461083c578063ee6c3bcf1461085a578063eeae17f61461088a576102f1565b8063a364f4da11610156578063b98d090811610130578063b98d090814610746578063bafa910714610764578063df5cf72314610782578063e15234ff146107a0576102f1565b8063a364f4da146106de578063a5b7890a146106fa578063a98fb3551461072a576102f1565b80638687feae14610630578063886f11951461064e5780638da5cb5b1461066c5780639926ee7d1461068a578063a0169ddd146106a6578063a20b99bf146106c2576102f1565b80635c975abb1161025c5780636d14a9871161021557806372276443116101ef57806372276443146105bc57806372d18e8d146105da578063775bbcb5146105f85780637794965a14610614576102f1565b80636d14a987146105635780636efb463614610581578063715018a6146105b2576102f1565b80635c975abb146104af5780635df45946146104cd5780635e033476146104eb5780635e8b3f2d1461050957806368304835146105275780636b3aa72e14610545576102f1565b806333cfb7b7116102ae57806333cfb7b7146103ef5780633bc28c8c1461041f578063416c7e5e1461043b5780634972134a14610457578063595c6a67146104755780635ac86ab71461047f576102f1565b8063048886d2146102f657806310d67a2f14610326578063136439dd146103425780631429c7c21461035e578063171f1d5b1461038e5780632ecfe72b146103bf575b600080fd5b610310600480360381019061030b91906159b1565b610972565b60405161031d91906159f9565b60405180910390f35b610340600480360381019061033b9190615a84565b610a15565b005b61035c60048036038101906103579190615ae7565b610b1f565b005b610378600480360381019061037391906159b1565b610c9a565b6040516103859190615b23565b60405180910390f35b6103a860048036038101906103a39190615d65565b610d3d565b6040516103b6929190615dcd565b60405180910390f35b6103d960048036038101906103d49190615e30565b610ea9565b6040516103e69190615ecd565b60405180910390f35b61040960048036038101906104049190615f14565b610f52565b6040516104169190615fff565b60405180910390f35b61043960048036038101906104349190615f14565b611510565b005b6104556004803603810190610450919061604d565b611524565b005b61045f61162c565b60405161046c9190616089565b60405180910390f35b61047d611640565b005b610499600480360381019061049491906159b1565b6117b2565b6040516104a691906159f9565b60405180910390f35b6104b76117ce565b6040516104c491906160b3565b60405180910390f35b6104d56117d8565b6040516104e2919061612d565b60405180910390f35b6104f36117fc565b6040516105009190616089565b60405180910390f35b610511611803565b60405161051e9190616089565b60405180910390f35b61052f611809565b60405161053c9190616169565b60405180910390f35b61054d61182d565b60405161055a9190616193565b60405180910390f35b61056b611855565b60405161057891906161cf565b60405180910390f35b61059b6004803603810190610596919061664d565b611879565b6040516105a992919061681a565b60405180910390f35b6105ba61268d565b005b6105c46126a1565b6040516105d1919061686b565b60405180910390f35b6105e26126c5565b6040516105ef9190616089565b60405180910390f35b610612600480360381019061060d9190616949565b6126de565b005b61062e60048036038101906106299190616a04565b61287d565b005b610638612cd3565b6040516106459190616b04565b60405180910390f35b610656612d6e565b6040516106639190616b47565b60405180910390f35b610674612d94565b6040516106819190616193565b60405180910390f35b6106a4600480360381019061069f9190616c97565b612dbe565b005b6106c060048036038101906106bb9190615f14565b612edd565b005b6106dc60048036038101906106d79190616d49565b612f73565b005b6106f860048036038101906106f39190615f14565b6131c2565b005b610714600480360381019061070f9190615f14565b6132de565b60405161072191906159f9565b60405180910390f35b610744600480360381019061073f9190616e37565b6132fe565b005b61074e613394565b60405161075b91906159f9565b60405180910390f35b61076c6133a7565b6040516107799190616b04565b60405180910390f35b61078a613442565b6040516107979190616ea1565b60405180910390f35b6107a8613466565b6040516107b59190616b04565b60405180910390f35b6107c6613501565b6040516107d39190615fff565b60405180910390f35b6107f660048036038101906107f19190616ebc565b613980565b6040516108039190616089565b60405180910390f35b61082660048036038101906108219190616ebc565b6139a5565b6040516108339190616ee9565b60405180910390f35b6108446139bd565b6040516108519190616f25565b60405180910390f35b610874600480360381019061086f91906159b1565b6139e1565b6040516108819190615b23565b60405180910390f35b610892613a84565b60405161089f9190616f61565b60405180910390f35b6108b0613aa8565b6040516108bd91906160b3565b60405180910390f35b6108e060048036038101906108db9190615f14565b613aad565b005b6108fc60048036038101906108f79190615f14565b613ac1565b005b610906613b45565b6040516109139190616f9d565b60405180910390f35b61093660048036038101906109319190615ae7565b613b69565b005b610940613d0a565b60405161094d9190616193565b60405180910390f35b610970600480360381019061096b919061700e565b613d30565b005b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663048886d2836040518263ffffffff1660e01b81526004016109cd9190615b23565b602060405180830381865afa1580156109ea573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a0e9190617070565b9050919050565b60fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a82573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aa691906170b2565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610b13576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b0a90617162565b60405180910390fd5b610b1c81613f1d565b50565b60fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b8152600401610b7a9190616193565b602060405180830381865afa158015610b97573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bbb9190617070565b610bfa576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bf1906171f4565b60405180910390fd5b60fc548160fc541614610c42576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c3990617286565b60405180910390fd5b8060fc819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d82604051610c8f91906160b3565b60405180910390a250565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631429c7c2836040518263ffffffff1660e01b8152600401610cf59190615b23565b602060405180830381865afa158015610d12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d3691906172bb565b9050919050565b60008060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000187876000015188602001518860000151600060028110610d8657610d856172e8565b5b60200201518960000151600160028110610da357610da26172e8565b5b60200201518a60200151600060028110610dc057610dbf6172e8565b5b60200201518b60200151600160028110610ddd57610ddc6172e8565b5b60200201518b600001518c60200151604051602001610e0499989796959493929190617359565b6040516020818303038152906040528051906020012060001c610e27919061742b565b9050610e97610e51610e42838961402c90919063ffffffff16565b8661410990919063ffffffff16565b610e5961420c565b610e8d610e7685610e686142d6565b61402c90919063ffffffff16565b610e7f8c6142fa565b61410990919063ffffffff16565b886201d4c0614405565b80935081945050505094509492505050565b610eb16157d0565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632ecfe72b836040518263ffffffff1660e01b8152600401610f0a919061746b565b606060405180830381865afa158015610f27573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f4b91906174ff565b9050919050565b606060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166313542a4e846040518263ffffffff1660e01b8152600401610faf9190616193565b602060405180830381865afa158015610fcc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ff09190617541565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663871ef049836040518263ffffffff1660e01b815260040161104d9190616ee9565b602060405180830381865afa15801561106a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108e91906175be565b905060008177ffffffffffffffffffffffffffffffffffffffffffffffff16148061114a575060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611121573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061114591906172bb565b60ff16145b156111a257600067ffffffffffffffff81111561116a57611169615b8a565b5b6040519080825280602002602001820160405280156111985781602001602082028036833780820191505090505b509250505061150b565b60006111c78277ffffffffffffffffffffffffffffffffffffffffffffffff166146b6565b9050600080600090505b82518110156112b4577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ca5a5f5848381518110611229576112286172e8565b5b602001015160f81c60f81b60f81c6040518263ffffffff1660e01b81526004016112539190615b23565b602060405180830381865afa158015611270573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112949190617600565b8261129f919061765c565b915080806112ac906176b2565b9150506111d1565b5060008167ffffffffffffffff8111156112d1576112d0615b8a565b5b6040519080825280602002602001820160405280156112ff5781602001602082028036833780820191505090505b5090506000805b8451811015611500576000858281518110611324576113236172e8565b5b602001015160f81c60f81b60f81c905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ca5a5f5836040518263ffffffff1660e01b815260040161138f9190615b23565b602060405180830381865afa1580156113ac573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113d09190617600565b905060005b818110156114ea577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663adc804da84836040518363ffffffff1660e01b81526004016114389291906176fb565b6040805180830381865afa158015611454573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061147891906177de565b6000015186868151811061148f5761148e6172e8565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505084806114d4906176b2565b95505080806114e2906176b2565b9150506113d5565b50505080806114f8906176b2565b915050611306565b508196505050505050505b919050565b6115186147ab565b61152181614829565b50565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561158f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115b391906170b2565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611620576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611617906178a3565b60405180910390fd5b611629816148c8565b50565b60008054906101000a900463ffffffff1681565b60fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b815260040161169b9190616193565b602060405180830381865afa1580156116b8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116dc9190617070565b61171b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611712906171f4565b60405180910390fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60fc819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040516117a891906160b3565b60405180910390a2565b6000808260ff166001901b9050808160fc541614915050919050565b600060fc54905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b620189c081565b61012c81565b7f000000000000000000000000000000000000000000000000000000000000000081565b60007f0000000000000000000000000000000000000000000000000000000000000000905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b611881615800565b6000808686905014156118c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118c090617935565b60405180910390fd5b826040015151868690501480156118e757508260a001515186869050145b80156118fa57508260c001515186869050145b801561190d57508260e001515186869050145b61194c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611943906179ed565b60405180910390fd5b82600001515183602001515114611998576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161198f90617aa5565b60405180910390fd5b4363ffffffff168463ffffffff16106119e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119dd90617b37565b60405180910390fd5b600060405180604001604052806000815260200160008152509050611a09615800565b8787905067ffffffffffffffff811115611a2657611a25615b8a565b5b604051908082528060200260200182016040528015611a545781602001602082028036833780820191505090505b5081602001819052508787905067ffffffffffffffff811115611a7a57611a79615b8a565b5b604051908082528060200260200182016040528015611aa85781602001602082028036833780820191505090505b508160000181905250611ab961581a565b85602001515167ffffffffffffffff811115611ad857611ad7615b8a565b5b604051908082528060200260200182016040528015611b065781602001602082028036833780820191505090505b50816000018190525085602001515167ffffffffffffffff811115611b2e57611b2d615b8a565b5b604051908082528060200260200182016040528015611b5c5781602001602082028036833780820191505090505b5081602001819052506000611c438a8a8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611c1a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c3e91906172bb565b61491c565b905060005b876020015151811015611ecf57611c7c88602001518281518110611c6f57611c6e6172e8565b5b602002602001015161497d565b83602001518281518110611c9357611c926172e8565b5b60200260200101818152505060008114611d38578260200151600182611cb99190617b57565b81518110611cca57611cc96172e8565b5b602002602001015160001c83602001518281518110611cec57611ceb6172e8565b5b602002602001015160001c11611d37576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611d2e90617bfd565b60405180910390fd5b5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166304ec635184602001518381518110611d8b57611d8a6172e8565b5b60200260200101518b8b600001518581518110611dab57611daa6172e8565b5b60200260200101516040518463ffffffff1660e01b8152600401611dd193929190617c4e565b602060405180830381865afa158015611dee573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e1291906175be565b77ffffffffffffffffffffffffffffffffffffffffffffffff1683600001518281518110611e4357611e426172e8565b5b602002602001018181525050611eba611eab611e7e8486600001518581518110611e7057611e6f6172e8565b5b602002602001015116614998565b8a602001518481518110611e9557611e946172e8565b5b60200260200101516149d690919063ffffffff16565b8661410990919063ffffffff16565b94508080611ec7906176b2565b915050611c48565b5050611eda83614ac9565b9250600060c960009054906101000a900460ff169050600081611efe576000611f8e565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c448feb86040518163ffffffff1660e01b8152600401602060405180830381865afa158015611f69573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f8d9190617600565b5b905060005b8b8b90508110156125a85782156120b0578963ffffffff16827f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663249a0c428f8f86818110611ffb57611ffa6172e8565b5b9050013560f81c60f81b60f81c6040518263ffffffff1660e01b81526004016120249190615b23565b602060405180830381865afa158015612041573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120659190617600565b61206f919061765c565b116120af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016120a690617d43565b60405180910390fd5b5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166368bccaac8d8d848181106120ff576120fe6172e8565b5b9050013560f81c60f81b60f81c8c8c60a001518581518110612124576121236172e8565b5b60200260200101516040518463ffffffff1660e01b815260040161214a93929190617d63565b602060405180830381865afa158015612167573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061218b9190617df2565b67ffffffffffffffff19166121bd8a6040015183815181106121b0576121af6172e8565b5b602002602001015161497d565b67ffffffffffffffff191614612208576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016121ff90617edd565b60405180910390fd5b61223989604001518281518110612222576122216172e8565b5b60200260200101518761410990919063ffffffff16565b95507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c8294c568d8d8481811061228a576122896172e8565b5b9050013560f81c60f81b60f81c8c8c60c0015185815181106122af576122ae6172e8565b5b60200260200101516040518463ffffffff1660e01b81526004016122d593929190617d63565b602060405180830381865afa1580156122f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123169190617efd565b8560200151828151811061232d5761232c6172e8565b5b60200260200101906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff16815250508460200151818151811061236e5761236d6172e8565b5b60200260200101518560000151828151811061238d5761238c6172e8565b5b60200260200101906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff16815250506000805b8a60200151518110156125935761240f866000015182815181106123e2576123e16172e8565b5b60200260200101518f8f868181106123fd576123fc6172e8565b5b9050013560f81c60f81b60f81c614b87565b15612580577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f2be94ae8f8f86818110612463576124626172e8565b5b9050013560f81c60f81b60f81c8e89602001518581518110612488576124876172e8565b5b60200260200101518f60e0015188815181106124a7576124a66172e8565b5b602002602001015187815181106124c1576124c06172e8565b5b60200260200101516040518563ffffffff1660e01b81526004016124e89493929190617f2a565b602060405180830381865afa158015612505573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125299190617efd565b876000015184815181106125405761253f6172e8565b5b602002602001018181516125549190617f6f565b9150906bffffffffffffffffffffffff1690816bffffffffffffffffffffffff16815250508160010191505b808061258b906176b2565b9150506123bb565b505080806125a0906176b2565b915050611f93565b5050506000806125c28c868a606001518b60800151610d3d565b9150915081612606576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016125fd9061803b565b60405180910390fd5b80612646576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161263d906180cd565b60405180910390fd5b505060008782602001516040516020016126619291906181db565b604051602081830303815290604052805190602001209050828195509550505050509550959350505050565b6126956147ab565b61269f6000614b9e565b565b7f000000000000000000000000000000000000000000000000000000000000000081565b60008060009054906101000a900463ffffffff16905090565b6000603260019054906101000a900460ff1615905080801561271257506001603260009054906101000a900460ff1660ff16105b80612741575061272130614c64565b15801561274057506001603260009054906101000a900460ff1660ff16145b5b612780576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161277790618275565b60405180910390fd5b6001603260006101000a81548160ff021916908360ff16021790555080156127be576001603260016101000a81548160ff0219169083151502179055505b6127c88686614c87565b6127d184614b9e565b6127da82614829565b60005b835181101561281a576128098482815181106127fc576127fb6172e8565b5b6020026020010151614db3565b80612813906176b2565b90506127dd565b508015612875576000603260016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498600160405161286c91906182d0565b60405180910390a15b505050505050565b6000612888816117b2565b156128c8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016128bf90618337565b60405180910390fd5b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661291e57600080fd5b3373ffffffffffffffffffffffffffffffffffffffff163273ffffffffffffffffffffffffffffffffffffffff161461298c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612983906183c9565b60405180910390fd5b438360600160208101906129a09190616ebc565b63ffffffff16106129e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016129dd9061845b565b60405180910390fd5b4363ffffffff1661012c846060016020810190612a039190616ebc565b612a0d919061847b565b63ffffffff161015612a54576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612a4b90618527565b60405180910390fd5b828060400190612a649190618556565b9050838060200190612a769190618556565b905014612ab8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612aaf9061862b565b60405180910390fd5b6000612acc84612ac7906186fb565b614edf565b9050600080612afe83878060200190612ae59190618556565b896060016020810190612af89190616ebc565b89611879565b9150915060005b868060400190612b159190618556565b9050811015612c1357868060400190612b2e9190618556565b82818110612b3f57612b3e6172e8565b5b9050013560f81c60f81b60f81c60ff1683602001518281518110612b6657612b656172e8565b5b6020026020010151612b78919061870e565b6bffffffffffffffffffffffff16606484600001518381518110612b9f57612b9e6172e8565b5b60200260200101516bffffffffffffffffffffffff16612bbf9190618754565b1015612c00576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612bf790618820565b60405180910390fd5b8080612c0b906176b2565b915050612b05565b5060008060009054906101000a900463ffffffff1690506000612c3588614f17565b9050612c42818443614f47565b600160008463ffffffff1663ffffffff16815260200190815260200160002081905550847fc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a83604051612c959190616089565b60405180910390a2600182612caa919061847b565b6000806101000a81548163ffffffff021916908363ffffffff1602179055505050505050505050565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638687feae6040518163ffffffff1660e01b8152600401600060405180830381865afa158015612d40573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190612d6991906188b0565b905090565b60fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612e4c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612e4390618991565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639926ee7d83836040518363ffffffff1660e01b8152600401612ea7929190618a69565b600060405180830381600087803b158015612ec157600080fd5b505af1158015612ed5573d6000803e3d6000fd5b505050505050565b612ee56147ab565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a0169ddd826040518263ffffffff1660e01b8152600401612f3e9190616193565b600060405180830381600087803b158015612f5857600080fd5b505af1158015612f6c573d6000803e3d6000fd5b5050505050565b612f7b614f7d565b60005b8282905081101561312e576000805b848484818110612fa057612f9f6172e8565b5b9050602002810190612fb29190618a99565b8060400190612fc19190618ac1565b905081101561303757848484818110612fdd57612fdc6172e8565b5b9050602002810190612fef9190618a99565b8060400190612ffe9190618ac1565b8281811061300f5761300e6172e8565b5b9050604002016020013582613024919061765c565b915080613030906176b2565b9050612f8d565b5061309b333083878787818110613051576130506172e8565b5b90506020028101906130639190618a99565b60200160208101906130759190618b62565b73ffffffffffffffffffffffffffffffffffffffff1661500f909392919063ffffffff16565b61311c7f0000000000000000000000000000000000000000000000000000000000000000828686868181106130d3576130d26172e8565b5b90506020028101906130e59190618a99565b60200160208101906130f79190618b62565b73ffffffffffffffffffffffffffffffffffffffff166150989092919063ffffffff16565b5080613127906176b2565b9050612f7e565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639cb9a5fa3084846040518463ffffffff1660e01b815260040161318c9392919061916e565b600060405180830381600087803b1580156131a657600080fd5b505af11580156131ba573d6000803e3d6000fd5b505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613250576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161324790618991565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a364f4da826040518263ffffffff1660e01b81526004016132a99190616193565b600060405180830381600087803b1580156132c357600080fd5b505af11580156132d7573d6000803e3d6000fd5b5050505050565b60026020528060005260406000206000915054906101000a900460ff1681565b6133066147ab565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a98fb355826040518263ffffffff1660e01b815260040161335f91906191e4565b600060405180830381600087803b15801561337957600080fd5b505af115801561338d573d6000803e3d6000fd5b5050505050565b60c960009054906101000a900460ff1681565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663bafa91076040518163ffffffff1660e01b8152600401600060405180830381865afa158015613414573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061343d91906188b0565b905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b60607f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e15234ff6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156134d3573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906134fc91906188b0565b905090565b606060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613570573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061359491906172bb565b60ff16905060008114156135f457600067ffffffffffffffff8111156135bd576135bc615b8a565b5b6040519080825280602002602001820160405280156135eb5781602001602082028036833780820191505090505b5091505061397d565b600080600090505b828110156136be577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ca5a5f5826040518263ffffffff1660e01b815260040161365d9190615b23565b602060405180830381865afa15801561367a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061369e9190617600565b826136a9919061765c565b915080806136b6906176b2565b9150506135fc565b5060008167ffffffffffffffff8111156136db576136da615b8a565b5b6040519080825280602002602001820160405280156137095781602001602082028036833780820191505090505b5090506000805b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561377b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061379f91906172bb565b60ff168110156139745760007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ca5a5f5836040518263ffffffff1660e01b81526004016138049190615b23565b602060405180830381865afa158015613821573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138459190617600565b905060005b8181101561395f577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663adc804da84836040518363ffffffff1660e01b81526004016138ad9291906176fb565b6040805180830381865afa1580156138c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906138ed91906177de565b60000151858581518110613904576139036172e8565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508380613949906176b2565b9450508080613957906176b2565b91505061384a565b5050808061396c906176b2565b915050613710565b50819450505050505b90565b600061012c620189c083613994919061847b565b61399e919061847b565b9050919050565b60016020528060005260406000206000915090505481565b7f000000000000000000000000000000000000000000000000000000000000000081565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663ee6c3bcf836040518263ffffffff1660e01b8152600401613a3c9190615b23565b602060405180830381865afa158015613a59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613a7d91906172bb565b9050919050565b7f000000000000000000000000000000000000000000000000000000000000000081565b606481565b613ab56147ab565b613abe81614db3565b50565b613ac96147ab565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415613b39576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613b3090619278565b60405180910390fd5b613b4281614b9e565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b60fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613bd6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613bfa91906170b2565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614613c67576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613c5e90617162565b60405180910390fd5b60fc5419811960fc54191614613cb2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613ca99061930a565b60405180910390fd5b8060fc819055503373ffffffffffffffffffffffffffffffffffffffff167f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c82604051613cff91906160b3565b60405180910390a250565b609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b613d38614f7d565b60005b82829050811015613e8b57613dd13330858585818110613d5e57613d5d6172e8565b5b9050602002810190613d70919061932a565b60400135868686818110613d8757613d866172e8565b5b9050602002810190613d99919061932a565b6020016020810190613dab9190618b62565b73ffffffffffffffffffffffffffffffffffffffff1661500f909392919063ffffffff16565b613e7a7f0000000000000000000000000000000000000000000000000000000000000000848484818110613e0857613e076172e8565b5b9050602002810190613e1a919061932a565b60400135858585818110613e3157613e306172e8565b5b9050602002810190613e43919061932a565b6020016020810190613e559190618b62565b73ffffffffffffffffffffffffffffffffffffffff166150989092919063ffffffff16565b80613e84906176b2565b9050613d3b565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fce36c7d83836040518363ffffffff1660e01b8152600401613ee79291906194ca565b600060405180830381600087803b158015613f0157600080fd5b505af1158015613f15573d6000803e3d6000fd5b505050505050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415613f8d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613f8490619586565b60405180910390fd5b7f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb660fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682604051613fe09291906195a6565b60405180910390a18060fb60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b614034615834565b61403c61584e565b836000015181600060038110614055576140546172e8565b5b602002018181525050836020015181600160038110614077576140766172e8565b5b6020020181815250508281600260038110614095576140946172e8565b5b602002018181525050600060408360608460076107d05a03fa905080600081146140be576140c0565bfe5b5080614101576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016140f89061961b565b60405180910390fd5b505092915050565b614111615834565b614119615870565b836000015181600060048110614132576141316172e8565b5b602002018181525050836020015181600160048110614154576141536172e8565b5b602002018181525050826000015181600260048110614176576141756172e8565b5b602002018181525050826020015181600360048110614198576141976172e8565b5b602002018181525050600060408360808460066107d05a03fa905080600081146141c1576141c3565bfe5b5080614204576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016141fb90619687565b60405180910390fd5b505092915050565b614214615892565b604051806040016040528060405180604001604052807f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed815250815260200160405180604001604052807f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d815250815250905090565b6142de615834565b6040518060400160405280600181526020016002815250905090565b614302615834565b60008060007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478560001c614336919061742b565b90505b6001156143e357614349816151aa565b80935081945050507f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47806143805761437f6173fc565b5b8283098314156143a9576040518060400160405280828152602001838152509350505050614400565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47806143d8576143d76173fc565b5b600182089050614339565b604051806040016040528060008152602001600081525093505050505b919050565b60008060006040518060400160405280898152602001878152509050600060405180604001604052808981526020018781525090506144426158b8565b60005b600281101561466b57600060068261445d9190618754565b9050848260028110614472576144716172e8565b5b60200201516000015183600083614489919061765c565b600c811061449a576144996172e8565b5b6020020181815250508482600281106144b6576144b56172e8565b5b602002015160200151836001836144cd919061765c565b600c81106144de576144dd6172e8565b5b6020020181815250508382600281106144fa576144f96172e8565b5b602002015160000151600060028110614516576145156172e8565b5b602002015183600283614529919061765c565b600c811061453a576145396172e8565b5b602002018181525050838260028110614556576145556172e8565b5b602002015160000151600160028110614572576145716172e8565b5b602002015183600383614585919061765c565b600c8110614596576145956172e8565b5b6020020181815250508382600281106145b2576145b16172e8565b5b6020020151602001516000600281106145ce576145cd6172e8565b5b6020020151836004836145e1919061765c565b600c81106145f2576145f16172e8565b5b60200201818152505083826002811061460e5761460d6172e8565b5b60200201516020015160016002811061462a576146296172e8565b5b60200201518360058361463d919061765c565b600c811061464e5761464d6172e8565b5b602002018181525050508080614663906176b2565b915050614445565b506146746158db565b60006020826020600c028560088cfa90508060008360006001811061469c5761469b6172e8565b5b602002015114159650965050505050509550959350505050565b60606000806146c484614998565b61ffff1667ffffffffffffffff8111156146e1576146e0615b8a565b5b6040519080825280601f01601f1916602001820160405280156147135781602001600182028036833780820191505090505b5090506000805b82518210801561472b575061010081105b1561479f57806001901b935060008487161461478e578060f81b838381518110614758576147576172e8565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508160010191505b80614798906176b2565b905061471a565b50819350505050919050565b6147b36152a2565b73ffffffffffffffffffffffffffffffffffffffff166147d1612d94565b73ffffffffffffffffffffffffffffffffffffffff1614614827576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161481e906196f3565b60405180910390fd5b565b7fe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168260405161487c929190619713565b60405180910390a180609760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b8060c960006101000a81548160ff0219169083151502179055507f40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc8160405161491191906159f9565b60405180910390a150565b600080614928846152aa565b9050808360ff166001901b11614973576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161496a906197ae565b60405180910390fd5b8091505092915050565b60008151600052816020015160205260406000209050919050565b600080600090505b60008311156149cd576001836149b69190617b57565b8316925080806149c5906197ce565b9150506149a0565b80915050919050565b6149de615834565b6102008261ffff1610614a26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614a1d90619845565b60405180910390fd5b60018261ffff161415614a3b57829050614ac3565b600060405180604001604052806000815260200160008152509050600084905060006001905060005b8161ffff168661ffff1610614abb576001808260ff168861ffff16901c1661ffff161415614a9957614a968484614109565b93505b614aa38384614109565b925060018261ffff16901b9150806001019050614a64565b839450505050505b92915050565b614ad1615834565b60008260000151148015614ae9575060008260200151145b15614b0c5760405180604001604052806000815260200160008152509050614b82565b6040518060400160405280836000015181526020017f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478460200151614b51919061742b565b7f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47614b7c9190617b57565b81525090505b919050565b600060018260ff1684901c16600114905092915050565b6000606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081606560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600073ffffffffffffffffffffffffffffffffffffffff1660fb60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16148015614d125750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b614d51576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614d48906198fd565b60405180910390fd5b8060fc819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d82604051614d9e91906160b3565b60405180910390a2614daf82613f1d565b5050565b600260008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507f5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d41481600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16604051614ed492919061991d565b60405180910390a150565b6000614eea826153d1565b604051602001614efa9190619975565b604051602081830303815290604052805190602001209050919050565b600081604051602001614f2a9190619ac5565b604051602081830303815290604052805190602001209050919050565b6000838383604051602001614f5e93929190619ae7565b6040516020818303038152906040528051906020012090509392505050565b609760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461500d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161500490619bbc565b60405180910390fd5b565b615092846323b872dd60e01b85858560405160240161503093929190619bdc565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050615403565b50505050565b6000818473ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e30866040518363ffffffff1660e01b81526004016150d6929190619713565b602060405180830381865afa1580156150f3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906151179190617600565b615121919061765c565b90506151a48463095ea7b360e01b8584604051602401615142929190619c13565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050615403565b50505050565b60008060007f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47806151de576151dd6173fc565b5b60037f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478061520f5761520e6173fc565b5b867f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd478061523f5761523e6173fc565b5b888909090890506000615293827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f527f30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd476154ca565b90508181935093505050915091565b600033905090565b6000610100825111156152f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016152e990619cd4565b60405180910390fd5b60008251141561530557600090506153cc565b6000808360008151811061531c5761531b6172e8565b5b602001015160f81c60f81b60f81c60ff166001901b91506000600190505b84518110156153c557848181518110615356576153556172e8565b5b602001015160f81c60f81b60f81c60ff166001901b91508282116153af576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016153a690619d8c565b60405180910390fd5b8183179250806153be906176b2565b905061533a565b5081925050505b919050565b6153d96158fd565b604051806040016040528083600001518152602001836060015163ffffffff168152509050919050565b6000615465826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661561a9092919063ffffffff16565b90506000815111156154c557808060200190518101906154859190617070565b6154c4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016154bb90619e1e565b60405180910390fd5b5b505050565b6000806154d56158db565b6154dd615920565b6020816000600681106154f3576154f26172e8565b5b602002018181525050602081600160068110615512576155116172e8565b5b602002018181525050602081600260068110615531576155306172e8565b5b602002018181525050868160036006811061554f5761554e6172e8565b5b602002018181525050858160046006811061556d5761556c6172e8565b5b602002018181525050848160056006811061558b5761558a6172e8565b5b60200201818152505060208260c08360056107d05a03fa925082600081146155b2576155b4565bfe5b50826155f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016155ec90619e8a565b60405180910390fd5b81600060018110615609576156086172e8565b5b602002015193505050509392505050565b60606156298484600085615632565b90509392505050565b606082471015615677576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161566e90619f1c565b60405180910390fd5b61568085615746565b6156bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016156b690619f88565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516156e89190619fe4565b60006040518083038185875af1925050503d8060008114615725576040519150601f19603f3d011682016040523d82523d6000602084013e61572a565b606091505b509150915061573a828286615769565b92505050949350505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b60608315615779578290506157c9565b60008351111561578c5782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016157c091906191e4565b60405180910390fd5b9392505050565b6040518060600160405280600063ffffffff168152602001600063ffffffff168152602001600060ff1681525090565b604051806040016040528060608152602001606081525090565b604051806040016040528060608152602001606081525090565b604051806040016040528060008152602001600081525090565b6040518060600160405280600390602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b60405180604001604052806158a5615942565b81526020016158b2615942565b81525090565b604051806101800160405280600c90602082028036833780820191505090505090565b6040518060200160405280600190602082028036833780820191505090505090565b604051806040016040528060008019168152602001600063ffffffff1681525090565b6040518060c00160405280600690602082028036833780820191505090505090565b6040518060400160405280600290602082028036833780820191505090505090565b6000604051905090565b600080fd5b600080fd5b600060ff82169050919050565b61598e81615978565b811461599957600080fd5b50565b6000813590506159ab81615985565b92915050565b6000602082840312156159c7576159c661596e565b5b60006159d58482850161599c565b91505092915050565b60008115159050919050565b6159f3816159de565b82525050565b6000602082019050615a0e60008301846159ea565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000615a3f82615a14565b9050919050565b6000615a5182615a34565b9050919050565b615a6181615a46565b8114615a6c57600080fd5b50565b600081359050615a7e81615a58565b92915050565b600060208284031215615a9a57615a9961596e565b5b6000615aa884828501615a6f565b91505092915050565b6000819050919050565b615ac481615ab1565b8114615acf57600080fd5b50565b600081359050615ae181615abb565b92915050565b600060208284031215615afd57615afc61596e565b5b6000615b0b84828501615ad2565b91505092915050565b615b1d81615978565b82525050565b6000602082019050615b386000830184615b14565b92915050565b6000819050919050565b615b5181615b3e565b8114615b5c57600080fd5b50565b600081359050615b6e81615b48565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b615bc282615b79565b810181811067ffffffffffffffff82111715615be157615be0615b8a565b5b80604052505050565b6000615bf4615964565b9050615c008282615bb9565b919050565b600080fd5b600060408284031215615c2057615c1f615b74565b5b615c2a6040615bea565b90506000615c3a84828501615ad2565b6000830152506020615c4e84828501615ad2565b60208301525092915050565b600080fd5b600067ffffffffffffffff821115615c7a57615c79615b8a565b5b602082029050919050565b600080fd5b6000615c9d615c9884615c5f565b615bea565b90508060208402830185811115615cb757615cb6615c85565b5b835b81811015615ce05780615ccc8882615ad2565b845260208401935050602081019050615cb9565b5050509392505050565b600082601f830112615cff57615cfe615c5a565b5b6002615d0c848285615c8a565b91505092915050565b600060808284031215615d2b57615d2a615b74565b5b615d356040615bea565b90506000615d4584828501615cea565b6000830152506040615d5984828501615cea565b60208301525092915050565b6000806000806101208587031215615d8057615d7f61596e565b5b6000615d8e87828801615b5f565b9450506020615d9f87828801615c0a565b9350506060615db087828801615d15565b92505060e0615dc187828801615c0a565b91505092959194509250565b6000604082019050615de260008301856159ea565b615def60208301846159ea565b9392505050565b600061ffff82169050919050565b615e0d81615df6565b8114615e1857600080fd5b50565b600081359050615e2a81615e04565b92915050565b600060208284031215615e4657615e4561596e565b5b6000615e5484828501615e1b565b91505092915050565b600063ffffffff82169050919050565b615e7681615e5d565b82525050565b615e8581615978565b82525050565b606082016000820151615ea16000850182615e6d565b506020820151615eb46020850182615e6d565b506040820151615ec76040850182615e7c565b50505050565b6000606082019050615ee26000830184615e8b565b92915050565b615ef181615a34565b8114615efc57600080fd5b50565b600081359050615f0e81615ee8565b92915050565b600060208284031215615f2a57615f2961596e565b5b6000615f3884828501615eff565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b615f7681615a34565b82525050565b6000615f888383615f6d565b60208301905092915050565b6000602082019050919050565b6000615fac82615f41565b615fb68185615f4c565b9350615fc183615f5d565b8060005b83811015615ff2578151615fd98882615f7c565b9750615fe483615f94565b925050600181019050615fc5565b5085935050505092915050565b600060208201905081810360008301526160198184615fa1565b905092915050565b61602a816159de565b811461603557600080fd5b50565b60008135905061604781616021565b92915050565b6000602082840312156160635761606261596e565b5b600061607184828501616038565b91505092915050565b61608381615e5d565b82525050565b600060208201905061609e600083018461607a565b92915050565b6160ad81615ab1565b82525050565b60006020820190506160c860008301846160a4565b92915050565b6000819050919050565b60006160f36160ee6160e984615a14565b6160ce565b615a14565b9050919050565b6000616105826160d8565b9050919050565b6000616117826160fa565b9050919050565b6161278161610c565b82525050565b6000602082019050616142600083018461611e565b92915050565b6000616153826160fa565b9050919050565b61616381616148565b82525050565b600060208201905061617e600083018461615a565b92915050565b61618d81615a34565b82525050565b60006020820190506161a86000830184616184565b92915050565b60006161b9826160fa565b9050919050565b6161c9816161ae565b82525050565b60006020820190506161e460008301846161c0565b92915050565b600080fd5b60008083601f84011261620557616204615c5a565b5b8235905067ffffffffffffffff811115616222576162216161ea565b5b60208301915083600182028301111561623e5761623d615c85565b5b9250929050565b61624e81615e5d565b811461625957600080fd5b50565b60008135905061626b81616245565b92915050565b600067ffffffffffffffff82111561628c5761628b615b8a565b5b602082029050602081019050919050565b60006162b06162ab84616271565b615bea565b905080838252602082019050602084028301858111156162d3576162d2615c85565b5b835b818110156162fc57806162e8888261625c565b8452602084019350506020810190506162d5565b5050509392505050565b600082601f83011261631b5761631a615c5a565b5b813561632b84826020860161629d565b91505092915050565b600067ffffffffffffffff82111561634f5761634e615b8a565b5b602082029050602081019050919050565b600061637361636e84616334565b615bea565b9050808382526020820190506040840283018581111561639657616395615c85565b5b835b818110156163bf57806163ab8882615c0a565b845260208401935050604081019050616398565b5050509392505050565b600082601f8301126163de576163dd615c5a565b5b81356163ee848260208601616360565b91505092915050565b600067ffffffffffffffff82111561641257616411615b8a565b5b602082029050602081019050919050565b6000616436616431846163f7565b615bea565b9050808382526020820190506020840283018581111561645957616458615c85565b5b835b818110156164a057803567ffffffffffffffff81111561647e5761647d615c5a565b5b80860161648b8982616306565b8552602085019450505060208101905061645b565b5050509392505050565b600082601f8301126164bf576164be615c5a565b5b81356164cf848260208601616423565b91505092915050565b600061018082840312156164ef576164ee615b74565b5b6164fa610100615bea565b9050600082013567ffffffffffffffff81111561651a57616519615c05565b5b61652684828501616306565b600083015250602082013567ffffffffffffffff81111561654a57616549615c05565b5b616556848285016163c9565b602083015250604082013567ffffffffffffffff81111561657a57616579615c05565b5b616586848285016163c9565b604083015250606061659a84828501615d15565b60608301525060e06165ae84828501615c0a565b60808301525061012082013567ffffffffffffffff8111156165d3576165d2615c05565b5b6165df84828501616306565b60a08301525061014082013567ffffffffffffffff81111561660457616603615c05565b5b61661084828501616306565b60c08301525061016082013567ffffffffffffffff81111561663557616634615c05565b5b616641848285016164aa565b60e08301525092915050565b6000806000806000608086880312156166695761666861596e565b5b600061667788828901615b5f565b955050602086013567ffffffffffffffff81111561669857616697615973565b5b6166a4888289016161ef565b945094505060406166b78882890161625c565b925050606086013567ffffffffffffffff8111156166d8576166d7615973565b5b6166e4888289016164d8565b9150509295509295909350565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b60006bffffffffffffffffffffffff82169050919050565b61673e8161671d565b82525050565b60006167508383616735565b60208301905092915050565b6000602082019050919050565b6000616774826166f1565b61677e81856166fc565b93506167898361670d565b8060005b838110156167ba5781516167a18882616744565b97506167ac8361675c565b92505060018101905061678d565b5085935050505092915050565b600060408301600083015184820360008601526167e48282616769565b915050602083015184820360208601526167fe8282616769565b9150508091505092915050565b61681481615b3e565b82525050565b6000604082019050818103600083015261683481856167c7565b9050616843602083018461680b565b9392505050565b6000616855826160fa565b9050919050565b6168658161684a565b82525050565b6000602082019050616880600083018461685c565b92915050565b600067ffffffffffffffff8211156168a1576168a0615b8a565b5b602082029050602081019050919050565b60006168c56168c084616886565b615bea565b905080838252602082019050602084028301858111156168e8576168e7615c85565b5b835b8181101561691157806168fd8882615eff565b8452602084019350506020810190506168ea565b5050509392505050565b600082601f8301126169305761692f615c5a565b5b81356169408482602086016168b2565b91505092915050565b600080600080600060a086880312156169655761696461596e565b5b600061697388828901615a6f565b955050602061698488828901615ad2565b945050604061699588828901615eff565b935050606086013567ffffffffffffffff8111156169b6576169b5615973565b5b6169c28882890161691b565b92505060806169d388828901615eff565b9150509295509295909350565b600080fd5b6000608082840312156169fb576169fa6169e0565b5b81905092915050565b60008060408385031215616a1b57616a1a61596e565b5b600083013567ffffffffffffffff811115616a3957616a38615973565b5b616a45858286016169e5565b925050602083013567ffffffffffffffff811115616a6657616a65615973565b5b616a72858286016164d8565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b83811015616ab6578082015181840152602081019050616a9b565b83811115616ac5576000848401525b50505050565b6000616ad682616a7c565b616ae08185616a87565b9350616af0818560208601616a98565b616af981615b79565b840191505092915050565b60006020820190508181036000830152616b1e8184616acb565b905092915050565b6000616b31826160fa565b9050919050565b616b4181616b26565b82525050565b6000602082019050616b5c6000830184616b38565b92915050565b600080fd5b600067ffffffffffffffff821115616b8257616b81615b8a565b5b616b8b82615b79565b9050602081019050919050565b82818337600083830152505050565b6000616bba616bb584616b67565b615bea565b905082815260208101848484011115616bd657616bd5616b62565b5b616be1848285616b98565b509392505050565b600082601f830112616bfe57616bfd615c5a565b5b8135616c0e848260208601616ba7565b91505092915050565b600060608284031215616c2d57616c2c615b74565b5b616c376060615bea565b9050600082013567ffffffffffffffff811115616c5757616c56615c05565b5b616c6384828501616be9565b6000830152506020616c7784828501615b5f565b6020830152506040616c8b84828501615ad2565b60408301525092915050565b60008060408385031215616cae57616cad61596e565b5b6000616cbc85828601615eff565b925050602083013567ffffffffffffffff811115616cdd57616cdc615973565b5b616ce985828601616c17565b9150509250929050565b60008083601f840112616d0957616d08615c5a565b5b8235905067ffffffffffffffff811115616d2657616d256161ea565b5b602083019150836020820283011115616d4257616d41615c85565b5b9250929050565b60008060208385031215616d6057616d5f61596e565b5b600083013567ffffffffffffffff811115616d7e57616d7d615973565b5b616d8a85828601616cf3565b92509250509250929050565b600067ffffffffffffffff821115616db157616db0615b8a565b5b616dba82615b79565b9050602081019050919050565b6000616dda616dd584616d96565b615bea565b905082815260208101848484011115616df657616df5616b62565b5b616e01848285616b98565b509392505050565b600082601f830112616e1e57616e1d615c5a565b5b8135616e2e848260208601616dc7565b91505092915050565b600060208284031215616e4d57616e4c61596e565b5b600082013567ffffffffffffffff811115616e6b57616e6a615973565b5b616e7784828501616e09565b91505092915050565b6000616e8b826160fa565b9050919050565b616e9b81616e80565b82525050565b6000602082019050616eb66000830184616e92565b92915050565b600060208284031215616ed257616ed161596e565b5b6000616ee08482850161625c565b91505092915050565b6000602082019050616efe600083018461680b565b92915050565b6000616f0f826160fa565b9050919050565b616f1f81616f04565b82525050565b6000602082019050616f3a6000830184616f16565b92915050565b6000616f4b826160fa565b9050919050565b616f5b81616f40565b82525050565b6000602082019050616f766000830184616f52565b92915050565b6000616f87826160fa565b9050919050565b616f9781616f7c565b82525050565b6000602082019050616fb26000830184616f8e565b92915050565b60008083601f840112616fce57616fcd615c5a565b5b8235905067ffffffffffffffff811115616feb57616fea6161ea565b5b60208301915083602082028301111561700757617006615c85565b5b9250929050565b600080602083850312156170255761702461596e565b5b600083013567ffffffffffffffff81111561704357617042615973565b5b61704f85828601616fb8565b92509250509250929050565b60008151905061706a81616021565b92915050565b6000602082840312156170865761708561596e565b5b60006170948482850161705b565b91505092915050565b6000815190506170ac81615ee8565b92915050565b6000602082840312156170c8576170c761596e565b5b60006170d68482850161709d565b91505092915050565b600082825260208201905092915050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320756e70617573657200000000000000000000000000000000000000000000602082015250565b600061714c602a836170df565b9150617157826170f0565b604082019050919050565b6000602082019050818103600083015261717b8161713f565b9050919050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320706175736572000000000000000000000000000000000000000000000000602082015250565b60006171de6028836170df565b91506171e982617182565b604082019050919050565b6000602082019050818103600083015261720d816171d1565b9050919050565b7f5061757361626c652e70617573653a20696e76616c696420617474656d70742060008201527f746f20756e70617573652066756e6374696f6e616c6974790000000000000000602082015250565b60006172706038836170df565b915061727b82617214565b604082019050919050565b6000602082019050818103600083015261729f81617263565b9050919050565b6000815190506172b581615985565b92915050565b6000602082840312156172d1576172d061596e565b5b60006172df848285016172a6565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b61733261732d82615b3e565b617317565b82525050565b6000819050919050565b61735361734e82615ab1565b617338565b82525050565b6000617365828c617321565b602082019150617375828b617342565b602082019150617385828a617342565b6020820191506173958289617342565b6020820191506173a58288617342565b6020820191506173b58287617342565b6020820191506173c58286617342565b6020820191506173d58285617342565b6020820191506173e58284617342565b6020820191508190509a9950505050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061743682615ab1565b915061744183615ab1565b925082617451576174506173fc565b5b828206905092915050565b61746581615df6565b82525050565b6000602082019050617480600083018461745c565b92915050565b60008151905061749581616245565b92915050565b6000606082840312156174b1576174b0615b74565b5b6174bb6060615bea565b905060006174cb84828501617486565b60008301525060206174df84828501617486565b60208301525060406174f3848285016172a6565b60408301525092915050565b6000606082840312156175155761751461596e565b5b60006175238482850161749b565b91505092915050565b60008151905061753b81615b48565b92915050565b6000602082840312156175575761755661596e565b5b60006175658482850161752c565b91505092915050565b600077ffffffffffffffffffffffffffffffffffffffffffffffff82169050919050565b61759b8161756e565b81146175a657600080fd5b50565b6000815190506175b881617592565b92915050565b6000602082840312156175d4576175d361596e565b5b60006175e2848285016175a9565b91505092915050565b6000815190506175fa81615abb565b92915050565b6000602082840312156176165761761561596e565b5b6000617624848285016175eb565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061766782615ab1565b915061767283615ab1565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156176a7576176a661762d565b5b828201905092915050565b60006176bd82615ab1565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156176f0576176ef61762d565b5b600182019050919050565b60006040820190506177106000830185615b14565b61771d60208301846160a4565b9392505050565b600061772f82615a34565b9050919050565b61773f81617724565b811461774a57600080fd5b50565b60008151905061775c81617736565b92915050565b61776b8161671d565b811461777657600080fd5b50565b60008151905061778881617762565b92915050565b6000604082840312156177a4576177a3615b74565b5b6177ae6040615bea565b905060006177be8482850161774d565b60008301525060206177d284828501617779565b60208301525092915050565b6000604082840312156177f4576177f361596e565b5b60006178028482850161778e565b91505092915050565b7f424c535369676e6174757265436865636b65722e6f6e6c79436f6f7264696e6160008201527f746f724f776e65723a2063616c6c6572206973206e6f7420746865206f776e6560208201527f72206f6620746865207265676973747279436f6f7264696e61746f7200000000604082015250565b600061788d605c836170df565b91506178988261780b565b606082019050919050565b600060208201905081810360008301526178bc81617880565b9050919050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a20656d7074792071756f72756d20696e707574000000000000000000602082015250565b600061791f6037836170df565b915061792a826178c3565b604082019050919050565b6000602082019050818103600083015261794e81617912565b9050919050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a20696e7075742071756f72756d206c656e677468206d69736d61746360208201527f6800000000000000000000000000000000000000000000000000000000000000604082015250565b60006179d76041836170df565b91506179e282617955565b606082019050919050565b60006020820190508181036000830152617a06816179ca565b9050919050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a20696e707574206e6f6e7369676e6572206c656e677468206d69736d60208201527f6174636800000000000000000000000000000000000000000000000000000000604082015250565b6000617a8f6044836170df565b9150617a9a82617a0d565b606082019050919050565b60006020820190508181036000830152617abe81617a82565b9050919050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a20696e76616c6964207265666572656e636520626c6f636b00000000602082015250565b6000617b21603c836170df565b9150617b2c82617ac5565b604082019050919050565b60006020820190508181036000830152617b5081617b14565b9050919050565b6000617b6282615ab1565b9150617b6d83615ab1565b925082821015617b8057617b7f61762d565b5b828203905092915050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a206e6f6e5369676e65725075626b657973206e6f7420736f72746564602082015250565b6000617be76040836170df565b9150617bf282617b8b565b604082019050919050565b60006020820190508181036000830152617c1681617bda565b9050919050565b6000617c38617c33617c2e84615e5d565b6160ce565b615ab1565b9050919050565b617c4881617c1d565b82525050565b6000606082019050617c63600083018661680b565b617c70602083018561607a565b617c7d6040830184617c3f565b949350505050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a205374616b6552656769737472792075706461746573206d7573742060208201527f62652077697468696e207769746864726177616c44656c6179426c6f636b732060408201527f77696e646f770000000000000000000000000000000000000000000000000000606082015250565b6000617d2d6066836170df565b9150617d3882617c85565b608082019050919050565b60006020820190508181036000830152617d5c81617d20565b9050919050565b6000606082019050617d786000830186615b14565b617d85602083018561607a565b617d926040830184617c3f565b949350505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000082169050919050565b617dcf81617d9a565b8114617dda57600080fd5b50565b600081519050617dec81617dc6565b92915050565b600060208284031215617e0857617e0761596e565b5b6000617e1684828501617ddd565b91505092915050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a2071756f72756d41706b206861736820696e2073746f72616765206460208201527f6f6573206e6f74206d617463682070726f76696465642071756f72756d20617060408201527f6b00000000000000000000000000000000000000000000000000000000000000606082015250565b6000617ec76061836170df565b9150617ed282617e1f565b608082019050919050565b60006020820190508181036000830152617ef681617eba565b9050919050565b600060208284031215617f1357617f1261596e565b5b6000617f2184828501617779565b91505092915050565b6000608082019050617f3f6000830187615b14565b617f4c602083018661607a565b617f59604083018561680b565b617f666060830184617c3f565b95945050505050565b6000617f7a8261671d565b9150617f858361671d565b925082821015617f9857617f9761762d565b5b828203905092915050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a2070616972696e6720707265636f6d70696c652063616c6c2066616960208201527f6c65640000000000000000000000000000000000000000000000000000000000604082015250565b60006180256043836170df565b915061803082617fa3565b606082019050919050565b6000602082019050818103600083015261805481618018565b9050919050565b7f424c535369676e6174757265436865636b65722e636865636b5369676e61747560008201527f7265733a207369676e617475726520697320696e76616c696400000000000000602082015250565b60006180b76039836170df565b91506180c28261805b565b604082019050919050565b600060208201905081810360008301526180e6816180aa565b9050919050565b60008160e01b9050919050565b6000618105826180ed565b9050919050565b61811d61811882615e5d565b6180fa565b82525050565b600081519050919050565b600081905092915050565b6000819050602082019050919050565b61815281615b3e565b82525050565b60006181648383618149565b60208301905092915050565b6000602082019050919050565b600061818882618123565b618192818561812e565b935061819d83618139565b8060005b838110156181ce5781516181b58882618158565b97506181c083618170565b9250506001810190506181a1565b5085935050505092915050565b60006181e7828561810c565b6004820191506181f7828461817d565b91508190509392505050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b600061825f602e836170df565b915061826a82618203565b604082019050919050565b6000602082019050818103600083015261828e81618252565b9050919050565b6000819050919050565b60006182ba6182b56182b084618295565b6160ce565b615978565b9050919050565b6182ca8161829f565b82525050565b60006020820190506182e560008301846182c1565b92915050565b7f5061757361626c653a20696e6465782069732070617573656400000000000000600082015250565b60006183216019836170df565b915061832c826182eb565b602082019050919050565b6000602082019050818103600083015261835081618314565b9050919050565b7f68656164657220616e64206e6f6e7369676e65722064617461206d757374206260008201527f6520696e2063616c6c6461746100000000000000000000000000000000000000602082015250565b60006183b3602d836170df565b91506183be82618357565b604082019050919050565b600060208201905081810360008301526183e2816183a6565b9050919050565b7f737065636966696564207265666572656e6365426c6f636b4e756d626572206960008201527f7320696e20667574757265000000000000000000000000000000000000000000602082015250565b6000618445602b836170df565b9150618450826183e9565b604082019050919050565b6000602082019050818103600083015261847481618438565b9050919050565b600061848682615e5d565b915061849183615e5d565b92508263ffffffff038211156184aa576184a961762d565b5b828201905092915050565b7f737065636966696564207265666572656e6365426c6f636b4e756d626572206960008201527f7320746f6f2066617220696e2070617374000000000000000000000000000000602082015250565b60006185116031836170df565b915061851c826184b5565b604082019050919050565b6000602082019050818103600083015261854081618504565b9050919050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261857357618572618547565b5b80840192508235915067ffffffffffffffff8211156185955761859461854c565b5b6020830192506001820236038313156185b1576185b0618551565b5b509250929050565b7f71756f72756d4e756d6265727320616e64207369676e65645374616b65466f7260008201527f51756f72756d73206d7573742062652073616d65206c656e6774680000000000602082015250565b6000618615603b836170df565b9150618620826185b9565b604082019050919050565b6000602082019050818103600083015261864481618608565b9050919050565b60006080828403121561866157618660615b74565b5b61866b6080615bea565b9050600061867b84828501615b5f565b600083015250602082013567ffffffffffffffff81111561869f5761869e615c05565b5b6186ab84828501616be9565b602083015250604082013567ffffffffffffffff8111156186cf576186ce615c05565b5b6186db84828501616be9565b60408301525060606186ef8482850161625c565b60608301525092915050565b6000618707368361864b565b9050919050565b60006187198261671d565b91506187248361671d565b9250816bffffffffffffffffffffffff04831182151516156187495761874861762d565b5b828202905092915050565b600061875f82615ab1565b915061876a83615ab1565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156187a3576187a261762d565b5b828202905092915050565b7f7369676e61746f7269657320646f206e6f74206f776e207468726573686f6c6460008201527f2070657263656e74616765206f6620612071756f72756d000000000000000000602082015250565b600061880a6037836170df565b9150618815826187ae565b604082019050919050565b60006020820190508181036000830152618839816187fd565b9050919050565b600061885361884e84616b67565b615bea565b90508281526020810184848401111561886f5761886e616b62565b5b61887a848285616a98565b509392505050565b600082601f83011261889757618896615c5a565b5b81516188a7848260208601618840565b91505092915050565b6000602082840312156188c6576188c561596e565b5b600082015167ffffffffffffffff8111156188e4576188e3615973565b5b6188f084828501618882565b91505092915050565b7f536572766963654d616e61676572426173652e6f6e6c7952656769737472794360008201527f6f6f7264696e61746f723a2063616c6c6572206973206e6f742074686520726560208201527f67697374727920636f6f7264696e61746f720000000000000000000000000000604082015250565b600061897b6052836170df565b9150618986826188f9565b606082019050919050565b600060208201905081810360008301526189aa8161896e565b9050919050565b600082825260208201905092915050565b60006189cd82616a7c565b6189d781856189b1565b93506189e7818560208601616a98565b6189f081615b79565b840191505092915050565b618a0481615b3e565b82525050565b618a1381615ab1565b82525050565b60006060830160008301518482036000860152618a3682826189c2565b9150506020830151618a4b60208601826189fb565b506040830151618a5e6040860182618a0a565b508091505092915050565b6000604082019050618a7e6000830185616184565b8181036020830152618a908184618a19565b90509392505050565b60008235600160c003833603038112618ab557618ab4618547565b5b80830191505092915050565b60008083356001602003843603038112618ade57618add618547565b5b80840192508235915067ffffffffffffffff821115618b0057618aff61854c565b5b602083019250604082023603831315618b1c57618b1b618551565b5b509250929050565b6000618b2f82615a34565b9050919050565b618b3f81618b24565b8114618b4a57600080fd5b50565b600081359050618b5c81618b36565b92915050565b600060208284031215618b7857618b7761596e565b5b6000618b8684828501618b4d565b91505092915050565b600082825260208201905092915050565b6000819050919050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112618bd657618bd5618bb4565b5b83810192508235915060208301925067ffffffffffffffff821115618bfe57618bfd618baa565b5b604082023603841315618c1457618c13618baf565b5b509250929050565b600082825260208201905092915050565b6000819050919050565b600081359050618c4681617736565b92915050565b6000618c5b6020840184618c37565b905092915050565b6000618c6e826160fa565b9050919050565b618c7e81618c63565b82525050565b600081359050618c9381617762565b92915050565b6000618ca86020840184618c84565b905092915050565b60408201618cc16000830183618c4c565b618cce6000850182618c75565b50618cdc6020830183618c99565b618ce96020850182616735565b50505050565b6000618cfb8383618cb0565b60408301905092915050565b600082905092915050565b6000604082019050919050565b6000618d2b8385618c1c565b9350618d3682618c2d565b8060005b85811015618d6f57618d4c8284618d07565b618d568882618cef565b9750618d6183618d12565b925050600181019050618d3a565b5085925050509392505050565b6000618d8b6020840184618b4d565b905092915050565b6000618d9e826160fa565b9050919050565b618dae81618d93565b82525050565b60008083356001602003843603038112618dd157618dd0618bb4565b5b83810192508235915060208301925067ffffffffffffffff821115618df957618df8618baa565b5b604082023603841315618e0f57618e0e618baf565b5b509250929050565b600082825260208201905092915050565b6000819050919050565b6000618e416020840184615eff565b905092915050565b6000618e586020840184615ad2565b905092915050565b60408201618e716000830183618e32565b618e7e6000850182615f6d565b50618e8c6020830183618e49565b618e996020850182618a0a565b50505050565b6000618eab8383618e60565b60408301905092915050565b600082905092915050565b6000604082019050919050565b6000618edb8385618e17565b9350618ee682618e28565b8060005b85811015618f1f57618efc8284618eb7565b618f068882618e9f565b9750618f1183618ec2565b925050600181019050618eea565b5085925050509392505050565b6000618f3b602084018461625c565b905092915050565b60008083356001602003843603038112618f6057618f5f618bb4565b5b83810192508235915060208301925067ffffffffffffffff821115618f8857618f87618baa565b5b600182023603841315618f9e57618f9d618baf565b5b509250929050565b600082825260208201905092915050565b6000618fc38385618fa6565b9350618fd0838584616b98565b618fd983615b79565b840190509392505050565b600060c08301618ff76000840184618bb9565b858303600087015261900a838284618d1f565b9250505061901b6020840184618d7c565b6190286020860182618da5565b506190366040840184618db4565b8583036040870152619049838284618ecf565b9250505061905a6060840184618f2c565b6190676060860182615e6d565b506190756080840184618f2c565b6190826080860182615e6d565b5061909060a0840184618f43565b85830360a08701526190a3838284618fb7565b925050508091505092915050565b60006190bd8383618fe4565b905092915050565b60008235600160c0038336030381126190e1576190e0618bb4565b5b82810191505092915050565b6000602082019050919050565b60006191068385618b8f565b93508360208402850161911884618ba0565b8060005b8781101561915c57848403895261913382846190c5565b61913d85826190b1565b9450619148836190ed565b925060208a0199505060018101905061911c565b50829750879450505050509392505050565b60006040820190506191836000830186616184565b81810360208301526191968184866190fa565b9050949350505050565b600081519050919050565b60006191b6826191a0565b6191c081856170df565b93506191d0818560208601616a98565b6191d981615b79565b840191505092915050565b600060208201905081810360008301526191fe81846191ab565b905092915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006192626026836170df565b915061926d82619206565b604082019050919050565b6000602082019050818103600083015261929181619255565b9050919050565b7f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060008201527f7420746f2070617573652066756e6374696f6e616c6974790000000000000000602082015250565b60006192f46038836170df565b91506192ff82619298565b604082019050919050565b60006020820190508181036000830152619323816192e7565b9050919050565b60008235600160a00383360303811261934657619345618547565b5b80830191505092915050565b600082825260208201905092915050565b6000819050919050565b600060a083016193806000840184618bb9565b8583036000870152619393838284618d1f565b925050506193a46020840184618d7c565b6193b16020860182618da5565b506193bf6040840184618e49565b6193cc6040860182618a0a565b506193da6060840184618f2c565b6193e76060860182615e6d565b506193f56080840184618f2c565b6194026080860182615e6d565b508091505092915050565b6000619419838361936d565b905092915050565b60008235600160a00383360303811261943d5761943c618bb4565b5b82810191505092915050565b6000602082019050919050565b60006194628385619352565b93508360208402850161947484619363565b8060005b878110156194b857848403895261948f8284619421565b619499858261940d565b94506194a483619449565b925060208a01995050600181019050619478565b50829750879450505050509392505050565b600060208201905081810360008301526194e5818486619456565b90509392505050565b7f5061757361626c652e5f73657450617573657252656769737472793a206e657760008201527f50617573657252656769737472792063616e6e6f7420626520746865207a657260208201527f6f20616464726573730000000000000000000000000000000000000000000000604082015250565b60006195706049836170df565b915061957b826194ee565b606082019050919050565b6000602082019050818103600083015261959f81619563565b9050919050565b60006040820190506195bb6000830185616b38565b6195c86020830184616b38565b9392505050565b7f65632d6d756c2d6661696c656400000000000000000000000000000000000000600082015250565b6000619605600d836170df565b9150619610826195cf565b602082019050919050565b60006020820190508181036000830152619634816195f8565b9050919050565b7f65632d6164642d6661696c656400000000000000000000000000000000000000600082015250565b6000619671600d836170df565b915061967c8261963b565b602082019050919050565b600060208201905081810360008301526196a081619664565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006196dd6020836170df565b91506196e8826196a7565b602082019050919050565b6000602082019050818103600083015261970c816196d0565b9050919050565b60006040820190506197286000830185616184565b6197356020830184616184565b9392505050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206269746d61702065786365656473206d61782076616c756500602082015250565b6000619798603f836170df565b91506197a38261973c565b604082019050919050565b600060208201905081810360008301526197c78161978b565b9050919050565b60006197d982615df6565b915061ffff8214156197ee576197ed61762d565b5b600182019050919050565b7f7363616c61722d746f6f2d6c6172676500000000000000000000000000000000600082015250565b600061982f6010836170df565b915061983a826197f9565b602082019050919050565b6000602082019050818103600083015261985e81619822565b9050919050565b7f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960008201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c60208201527f6564206f6e636500000000000000000000000000000000000000000000000000604082015250565b60006198e76047836170df565b91506198f282619865565b606082019050919050565b60006020820190508181036000830152619916816198da565b9050919050565b60006040820190506199326000830185616184565b61993f60208301846159ea565b9392505050565b60408201600082015161995c60008501826189fb565b50602082015161996f6020850182615e6d565b50505050565b600060408201905061998a6000830184619946565b92915050565b600061999f6020840184615b5f565b905092915050565b600080833560016020038436030381126199c4576199c3618bb4565b5b83810192508235915060208301925067ffffffffffffffff8211156199ec576199eb618baa565b5b600182023603841315619a0257619a01618baf565b5b509250929050565b6000619a1683856189b1565b9350619a23838584616b98565b619a2c83615b79565b840190509392505050565b600060808301619a4a6000840184619990565b619a5760008601826189fb565b50619a6560208401846199a7565b8583036020870152619a78838284619a0a565b92505050619a8960408401846199a7565b8583036040870152619a9c838284619a0a565b92505050619aad6060840184618f2c565b619aba6060860182615e6d565b508091505092915050565b60006020820190508181036000830152619adf8184619a37565b905092915050565b6000619af38286617321565b602082019150619b038285617321565b602082019150619b13828461810c565b600482019150819050949350505050565b7f536572766963654d616e61676572426173652e6f6e6c7952657761726473496e60008201527f69746961746f723a2063616c6c6572206973206e6f742074686520726577617260208201527f647320696e69746961746f720000000000000000000000000000000000000000604082015250565b6000619ba6604c836170df565b9150619bb182619b24565b606082019050919050565b60006020820190508181036000830152619bd581619b99565b9050919050565b6000606082019050619bf16000830186616184565b619bfe6020830185616184565b619c0b60408301846160a4565b949350505050565b6000604082019050619c286000830185616184565b619c3560208301846160a4565b9392505050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f7264657265644279746573417272617920697320746f6f2060208201527f6c6f6e6700000000000000000000000000000000000000000000000000000000604082015250565b6000619cbe6044836170df565b9150619cc982619c3c565b606082019050919050565b60006020820190508181036000830152619ced81619cb1565b9050919050565b7f4269746d61705574696c732e6f72646572656442797465734172726179546f4260008201527f69746d61703a206f72646572656442797465734172726179206973206e6f742060208201527f6f72646572656400000000000000000000000000000000000000000000000000604082015250565b6000619d766047836170df565b9150619d8182619cf4565b606082019050919050565b60006020820190508181036000830152619da581619d69565b9050919050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000619e08602a836170df565b9150619e1382619dac565b604082019050919050565b60006020820190508181036000830152619e3781619dfb565b9050919050565b7f424e3235342e6578704d6f643a2063616c6c206661696c757265000000000000600082015250565b6000619e74601a836170df565b9150619e7f82619e3e565b602082019050919050565b60006020820190508181036000830152619ea381619e67565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b6000619f066026836170df565b9150619f1182619eaa565b604082019050919050565b60006020820190508181036000830152619f3581619ef9565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b6000619f72601d836170df565b9150619f7d82619f3c565b602082019050919050565b60006020820190508181036000830152619fa181619f65565b9050919050565b600081905092915050565b6000619fbe82616a7c565b619fc88185619fa8565b9350619fd8818560208601616a98565b80840191505092915050565b6000619ff08284619fb3565b91508190509291505056fea26469706673582212203e5a3821771bff0ed98cbf1f7b5134bc979fb59c69a7d37bdd326a44a50d2aba64736f6c634300080c0033",
}

// ContractEigenDAServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.ABI instead.
var ContractEigenDAServiceManagerABI = ContractEigenDAServiceManagerMetaData.ABI

// ContractEigenDAServiceManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.Bin instead.
var ContractEigenDAServiceManagerBin = ContractEigenDAServiceManagerMetaData.Bin

// DeployContractEigenDAServiceManager deploys a new Ethereum contract, binding an instance of ContractEigenDAServiceManager to it.
func DeployContractEigenDAServiceManager(auth *bind.TransactOpts, backend bind.ContractBackend, __avsDirectory common.Address, __rewardsCoordinator common.Address, __registryCoordinator common.Address, __stakeRegistry common.Address, __eigenDAThresholdRegistry common.Address, __eigenDARelayRegistry common.Address, __paymentVault common.Address, __eigenDADisperserRegistry common.Address) (common.Address, *types.Transaction, *ContractEigenDAServiceManager, error) {
	parsed, err := ContractEigenDAServiceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDAServiceManagerBin), backend, __avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry, __eigenDAThresholdRegistry, __eigenDARelayRegistry, __paymentVault, __eigenDADisperserRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractEigenDAServiceManager{ContractEigenDAServiceManagerCaller: ContractEigenDAServiceManagerCaller{contract: contract}, ContractEigenDAServiceManagerTransactor: ContractEigenDAServiceManagerTransactor{contract: contract}, ContractEigenDAServiceManagerFilterer: ContractEigenDAServiceManagerFilterer{contract: contract}}, nil
}

// ContractEigenDAServiceManager is an auto generated Go binding around an Ethereum contract.
type ContractEigenDAServiceManager struct {
	ContractEigenDAServiceManagerCaller     // Read-only binding to the contract
	ContractEigenDAServiceManagerTransactor // Write-only binding to the contract
	ContractEigenDAServiceManagerFilterer   // Log filterer for contract events
}

// ContractEigenDAServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractEigenDAServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractEigenDAServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractEigenDAServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractEigenDAServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractEigenDAServiceManagerSession struct {
	Contract     *ContractEigenDAServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ContractEigenDAServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractEigenDAServiceManagerCallerSession struct {
	Contract *ContractEigenDAServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ContractEigenDAServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractEigenDAServiceManagerTransactorSession struct {
	Contract     *ContractEigenDAServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ContractEigenDAServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractEigenDAServiceManagerRaw struct {
	Contract *ContractEigenDAServiceManager // Generic contract binding to access the raw methods on
}

// ContractEigenDAServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractEigenDAServiceManagerCallerRaw struct {
	Contract *ContractEigenDAServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractEigenDAServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractEigenDAServiceManagerTransactorRaw struct {
	Contract *ContractEigenDAServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractEigenDAServiceManager creates a new instance of ContractEigenDAServiceManager, bound to a specific deployed contract.
func NewContractEigenDAServiceManager(address common.Address, backend bind.ContractBackend) (*ContractEigenDAServiceManager, error) {
	contract, err := bindContractEigenDAServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManager{ContractEigenDAServiceManagerCaller: ContractEigenDAServiceManagerCaller{contract: contract}, ContractEigenDAServiceManagerTransactor: ContractEigenDAServiceManagerTransactor{contract: contract}, ContractEigenDAServiceManagerFilterer: ContractEigenDAServiceManagerFilterer{contract: contract}}, nil
}

// NewContractEigenDAServiceManagerCaller creates a new read-only instance of ContractEigenDAServiceManager, bound to a specific deployed contract.
func NewContractEigenDAServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractEigenDAServiceManagerCaller, error) {
	contract, err := bindContractEigenDAServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerCaller{contract: contract}, nil
}

// NewContractEigenDAServiceManagerTransactor creates a new write-only instance of ContractEigenDAServiceManager, bound to a specific deployed contract.
func NewContractEigenDAServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractEigenDAServiceManagerTransactor, error) {
	contract, err := bindContractEigenDAServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerTransactor{contract: contract}, nil
}

// NewContractEigenDAServiceManagerFilterer creates a new log filterer instance of ContractEigenDAServiceManager, bound to a specific deployed contract.
func NewContractEigenDAServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractEigenDAServiceManagerFilterer, error) {
	contract, err := bindContractEigenDAServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerFilterer{contract: contract}, nil
}

// bindContractEigenDAServiceManager binds a generic wrapper to an already deployed contract.
func bindContractEigenDAServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractEigenDAServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAServiceManager.Contract.ContractEigenDAServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.ContractEigenDAServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.ContractEigenDAServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractEigenDAServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.contract.Transact(opts, method, params...)
}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) BLOCKSTALEMEASURE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "BLOCK_STALE_MEASURE")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) BLOCKSTALEMEASURE() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.BLOCKSTALEMEASURE(&_ContractEigenDAServiceManager.CallOpts)
}

// BLOCKSTALEMEASURE is a free data retrieval call binding the contract method 0x5e8b3f2d.
//
// Solidity: function BLOCK_STALE_MEASURE() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) BLOCKSTALEMEASURE() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.BLOCKSTALEMEASURE(&_ContractEigenDAServiceManager.CallOpts)
}

// STOREDURATIONBLOCKS is a free data retrieval call binding the contract method 0x5e033476.
//
// Solidity: function STORE_DURATION_BLOCKS() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) STOREDURATIONBLOCKS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "STORE_DURATION_BLOCKS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// STOREDURATIONBLOCKS is a free data retrieval call binding the contract method 0x5e033476.
//
// Solidity: function STORE_DURATION_BLOCKS() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) STOREDURATIONBLOCKS() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.STOREDURATIONBLOCKS(&_ContractEigenDAServiceManager.CallOpts)
}

// STOREDURATIONBLOCKS is a free data retrieval call binding the contract method 0x5e033476.
//
// Solidity: function STORE_DURATION_BLOCKS() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) STOREDURATIONBLOCKS() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.STOREDURATIONBLOCKS(&_ContractEigenDAServiceManager.CallOpts)
}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) THRESHOLDDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "THRESHOLD_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) THRESHOLDDENOMINATOR() (*big.Int, error) {
	return _ContractEigenDAServiceManager.Contract.THRESHOLDDENOMINATOR(&_ContractEigenDAServiceManager.CallOpts)
}

// THRESHOLDDENOMINATOR is a free data retrieval call binding the contract method 0xef024458.
//
// Solidity: function THRESHOLD_DENOMINATOR() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) THRESHOLDDENOMINATOR() (*big.Int, error) {
	return _ContractEigenDAServiceManager.Contract.THRESHOLDDENOMINATOR(&_ContractEigenDAServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.AvsDirectory(&_ContractEigenDAServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.AvsDirectory(&_ContractEigenDAServiceManager.CallOpts)
}

// BatchId is a free data retrieval call binding the contract method 0x4972134a.
//
// Solidity: function batchId() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) BatchId(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "batchId")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// BatchId is a free data retrieval call binding the contract method 0x4972134a.
//
// Solidity: function batchId() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) BatchId() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.BatchId(&_ContractEigenDAServiceManager.CallOpts)
}

// BatchId is a free data retrieval call binding the contract method 0x4972134a.
//
// Solidity: function batchId() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) BatchId() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.BatchId(&_ContractEigenDAServiceManager.CallOpts)
}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 ) view returns(bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) BatchIdToBatchMetadataHash(opts *bind.CallOpts, arg0 uint32) ([32]byte, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "batchIdToBatchMetadataHash", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 ) view returns(bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) BatchIdToBatchMetadataHash(arg0 uint32) ([32]byte, error) {
	return _ContractEigenDAServiceManager.Contract.BatchIdToBatchMetadataHash(&_ContractEigenDAServiceManager.CallOpts, arg0)
}

// BatchIdToBatchMetadataHash is a free data retrieval call binding the contract method 0xeccbbfc9.
//
// Solidity: function batchIdToBatchMetadataHash(uint32 ) view returns(bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) BatchIdToBatchMetadataHash(arg0 uint32) ([32]byte, error) {
	return _ContractEigenDAServiceManager.Contract.BatchIdToBatchMetadataHash(&_ContractEigenDAServiceManager.CallOpts, arg0)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) BlsApkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "blsApkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.BlsApkRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// BlsApkRegistry is a free data retrieval call binding the contract method 0x5df45946.
//
// Solidity: function blsApkRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) BlsApkRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.BlsApkRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) CheckSignatures(opts *bind.CallOpts, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "checkSignatures", msgHash, quorumNumbers, referenceBlockNumber, params)

	if err != nil {
		return *new(IBLSSignatureCheckerQuorumStakeTotals), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(IBLSSignatureCheckerQuorumStakeTotals)).(*IBLSSignatureCheckerQuorumStakeTotals)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return out0, out1, err

}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractEigenDAServiceManager.Contract.CheckSignatures(&_ContractEigenDAServiceManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// CheckSignatures is a free data retrieval call binding the contract method 0x6efb4636.
//
// Solidity: function checkSignatures(bytes32 msgHash, bytes quorumNumbers, uint32 referenceBlockNumber, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) params) view returns((uint96[],uint96[]), bytes32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) CheckSignatures(msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, params IBLSSignatureCheckerNonSignerStakesAndSignature) (IBLSSignatureCheckerQuorumStakeTotals, [32]byte, error) {
	return _ContractEigenDAServiceManager.Contract.CheckSignatures(&_ContractEigenDAServiceManager.CallOpts, msgHash, quorumNumbers, referenceBlockNumber, params)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) Delegation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "delegation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Delegation() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.Delegation(&_ContractEigenDAServiceManager.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) Delegation() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.Delegation(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDADisperserRegistry is a free data retrieval call binding the contract method 0xeeae17f6.
//
// Solidity: function eigenDADisperserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) EigenDADisperserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "eigenDADisperserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDADisperserRegistry is a free data retrieval call binding the contract method 0xeeae17f6.
//
// Solidity: function eigenDADisperserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) EigenDADisperserRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDADisperserRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDADisperserRegistry is a free data retrieval call binding the contract method 0xeeae17f6.
//
// Solidity: function eigenDADisperserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) EigenDADisperserRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDADisperserRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) EigenDARelayRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "eigenDARelayRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDARelayRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDARelayRegistry is a free data retrieval call binding the contract method 0x72276443.
//
// Solidity: function eigenDARelayRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) EigenDARelayRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDARelayRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) EigenDAThresholdRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "eigenDAThresholdRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDAThresholdRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// EigenDAThresholdRegistry is a free data retrieval call binding the contract method 0xf8c66814.
//
// Solidity: function eigenDAThresholdRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) EigenDAThresholdRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.EigenDAThresholdRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetBlobParams(opts *bind.CallOpts, version uint16) (VersionedBlobParams, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getBlobParams", version)

	if err != nil {
		return *new(VersionedBlobParams), err
	}

	out0 := *abi.ConvertType(out[0], new(VersionedBlobParams)).(*VersionedBlobParams)

	return out0, err

}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDAServiceManager.Contract.GetBlobParams(&_ContractEigenDAServiceManager.CallOpts, version)
}

// GetBlobParams is a free data retrieval call binding the contract method 0x2ecfe72b.
//
// Solidity: function getBlobParams(uint16 version) view returns((uint32,uint32,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetBlobParams(version uint16) (VersionedBlobParams, error) {
	return _ContractEigenDAServiceManager.Contract.GetBlobParams(&_ContractEigenDAServiceManager.CallOpts, version)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetIsQuorumRequired(opts *bind.CallOpts, quorumNumber uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getIsQuorumRequired", quorumNumber)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.GetIsQuorumRequired(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetIsQuorumRequired is a free data retrieval call binding the contract method 0x048886d2.
//
// Solidity: function getIsQuorumRequired(uint8 quorumNumber) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetIsQuorumRequired(quorumNumber uint8) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.GetIsQuorumRequired(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractEigenDAServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.GetOperatorRestakedStrategies(&_ContractEigenDAServiceManager.CallOpts, operator)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetQuorumAdversaryThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getQuorumAdversaryThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetQuorumConfirmationThresholdPercentage(opts *bind.CallOpts, quorumNumber uint8) (uint8, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getQuorumConfirmationThresholdPercentage", quorumNumber)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.GetRestakeableStrategies(&_ContractEigenDAServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.GetRestakeableStrategies(&_ContractEigenDAServiceManager.CallOpts)
}

// IsBatchConfirmer is a free data retrieval call binding the contract method 0xa5b7890a.
//
// Solidity: function isBatchConfirmer(address ) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) IsBatchConfirmer(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "isBatchConfirmer", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBatchConfirmer is a free data retrieval call binding the contract method 0xa5b7890a.
//
// Solidity: function isBatchConfirmer(address ) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) IsBatchConfirmer(arg0 common.Address) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.IsBatchConfirmer(&_ContractEigenDAServiceManager.CallOpts, arg0)
}

// IsBatchConfirmer is a free data retrieval call binding the contract method 0xa5b7890a.
//
// Solidity: function isBatchConfirmer(address ) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) IsBatchConfirmer(arg0 common.Address) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.IsBatchConfirmer(&_ContractEigenDAServiceManager.CallOpts, arg0)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) LatestServeUntilBlock(opts *bind.CallOpts, referenceBlockNumber uint32) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "latestServeUntilBlock", referenceBlockNumber)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractEigenDAServiceManager.CallOpts, referenceBlockNumber)
}

// LatestServeUntilBlock is a free data retrieval call binding the contract method 0xeaefd27d.
//
// Solidity: function latestServeUntilBlock(uint32 referenceBlockNumber) view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) LatestServeUntilBlock(referenceBlockNumber uint32) (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.LatestServeUntilBlock(&_ContractEigenDAServiceManager.CallOpts, referenceBlockNumber)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Owner() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.Owner(&_ContractEigenDAServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) Owner() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.Owner(&_ContractEigenDAServiceManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Paused(index uint8) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.Paused(&_ContractEigenDAServiceManager.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) Paused(index uint8) (bool, error) {
	return _ContractEigenDAServiceManager.Contract.Paused(&_ContractEigenDAServiceManager.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Paused0() (*big.Int, error) {
	return _ContractEigenDAServiceManager.Contract.Paused0(&_ContractEigenDAServiceManager.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) Paused0() (*big.Int, error) {
	return _ContractEigenDAServiceManager.Contract.Paused0(&_ContractEigenDAServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) PauserRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.PauserRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) PauserRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.PauserRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// PaymentVault is a free data retrieval call binding the contract method 0xed3916f7.
//
// Solidity: function paymentVault() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) PaymentVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "paymentVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentVault is a free data retrieval call binding the contract method 0xed3916f7.
//
// Solidity: function paymentVault() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) PaymentVault() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.PaymentVault(&_ContractEigenDAServiceManager.CallOpts)
}

// PaymentVault is a free data retrieval call binding the contract method 0xed3916f7.
//
// Solidity: function paymentVault() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) PaymentVault() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.PaymentVault(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) QuorumAdversaryThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "quorumAdversaryThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumAdversaryThresholdPercentages is a free data retrieval call binding the contract method 0x8687feae.
//
// Solidity: function quorumAdversaryThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) QuorumAdversaryThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumAdversaryThresholdPercentages(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) QuorumConfirmationThresholdPercentages(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "quorumConfirmationThresholdPercentages")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumConfirmationThresholdPercentages is a free data retrieval call binding the contract method 0xbafa9107.
//
// Solidity: function quorumConfirmationThresholdPercentages() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) QuorumConfirmationThresholdPercentages() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumConfirmationThresholdPercentages(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) QuorumNumbersRequired(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "quorumNumbersRequired")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumNumbersRequired(&_ContractEigenDAServiceManager.CallOpts)
}

// QuorumNumbersRequired is a free data retrieval call binding the contract method 0xe15234ff.
//
// Solidity: function quorumNumbersRequired() view returns(bytes)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) QuorumNumbersRequired() ([]byte, error) {
	return _ContractEigenDAServiceManager.Contract.QuorumNumbersRequired(&_ContractEigenDAServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) RegistryCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "registryCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.RegistryCoordinator(&_ContractEigenDAServiceManager.CallOpts)
}

// RegistryCoordinator is a free data retrieval call binding the contract method 0x6d14a987.
//
// Solidity: function registryCoordinator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) RegistryCoordinator() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.RegistryCoordinator(&_ContractEigenDAServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) RewardsInitiator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "rewardsInitiator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) RewardsInitiator() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.RewardsInitiator(&_ContractEigenDAServiceManager.CallOpts)
}

// RewardsInitiator is a free data retrieval call binding the contract method 0xfc299dee.
//
// Solidity: function rewardsInitiator() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) RewardsInitiator() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.RewardsInitiator(&_ContractEigenDAServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) StakeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "stakeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) StakeRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.StakeRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// StakeRegistry is a free data retrieval call binding the contract method 0x68304835.
//
// Solidity: function stakeRegistry() view returns(address)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) StakeRegistry() (common.Address, error) {
	return _ContractEigenDAServiceManager.Contract.StakeRegistry(&_ContractEigenDAServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) StaleStakesForbidden(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "staleStakesForbidden")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) StaleStakesForbidden() (bool, error) {
	return _ContractEigenDAServiceManager.Contract.StaleStakesForbidden(&_ContractEigenDAServiceManager.CallOpts)
}

// StaleStakesForbidden is a free data retrieval call binding the contract method 0xb98d0908.
//
// Solidity: function staleStakesForbidden() view returns(bool)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) StaleStakesForbidden() (bool, error) {
	return _ContractEigenDAServiceManager.Contract.StaleStakesForbidden(&_ContractEigenDAServiceManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) TaskNumber(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "taskNumber")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) TaskNumber() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.TaskNumber(&_ContractEigenDAServiceManager.CallOpts)
}

// TaskNumber is a free data retrieval call binding the contract method 0x72d18e8d.
//
// Solidity: function taskNumber() view returns(uint32)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) TaskNumber() (uint32, error) {
	return _ContractEigenDAServiceManager.Contract.TaskNumber(&_ContractEigenDAServiceManager.CallOpts)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) TrySignatureAndApkVerification(opts *bind.CallOpts, msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "trySignatureAndApkVerification", msgHash, apk, apkG2, sigma)

	outstruct := new(struct {
		PairingSuccessful bool
		SiganatureIsValid bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PairingSuccessful = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.SiganatureIsValid = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractEigenDAServiceManager.Contract.TrySignatureAndApkVerification(&_ContractEigenDAServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// TrySignatureAndApkVerification is a free data retrieval call binding the contract method 0x171f1d5b.
//
// Solidity: function trySignatureAndApkVerification(bytes32 msgHash, (uint256,uint256) apk, (uint256[2],uint256[2]) apkG2, (uint256,uint256) sigma) view returns(bool pairingSuccessful, bool siganatureIsValid)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) TrySignatureAndApkVerification(msgHash [32]byte, apk BN254G1Point, apkG2 BN254G2Point, sigma BN254G1Point) (struct {
	PairingSuccessful bool
	SiganatureIsValid bool
}, error) {
	return _ContractEigenDAServiceManager.Contract.TrySignatureAndApkVerification(&_ContractEigenDAServiceManager.CallOpts, msgHash, apk, apkG2, sigma)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) ConfirmBatch(opts *bind.TransactOpts, batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "confirmBatch", batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) ConfirmBatch(batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.ConfirmBatch(&_ContractEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) ConfirmBatch(batchHeader BatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.ConfirmBatch(&_ContractEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) CreateAVSRewardsSubmission(opts *bind.TransactOpts, rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "createAVSRewardsSubmission", rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractEigenDAServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xfce36c7d.
//
// Solidity: function createAVSRewardsSubmission(((address,uint96)[],address,uint256,uint32,uint32)[] rewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) CreateAVSRewardsSubmission(rewardsSubmissions []IRewardsCoordinatorRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.CreateAVSRewardsSubmission(&_ContractEigenDAServiceManager.TransactOpts, rewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) CreateOperatorDirectedAVSRewardsSubmission(opts *bind.TransactOpts, operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "createOperatorDirectedAVSRewardsSubmission", operatorDirectedRewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) CreateOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.CreateOperatorDirectedAVSRewardsSubmission(&_ContractEigenDAServiceManager.TransactOpts, operatorDirectedRewardsSubmissions)
}

// CreateOperatorDirectedAVSRewardsSubmission is a paid mutator transaction binding the contract method 0xa20b99bf.
//
// Solidity: function createOperatorDirectedAVSRewardsSubmission(((address,uint96)[],address,(address,uint256)[],uint32,uint32,string)[] operatorDirectedRewardsSubmissions) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) CreateOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions []IRewardsCoordinatorOperatorDirectedRewardsSubmission) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.CreateOperatorDirectedAVSRewardsSubmission(&_ContractEigenDAServiceManager.TransactOpts, operatorDirectedRewardsSubmissions)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractEigenDAServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.DeregisterOperatorFromAVS(&_ContractEigenDAServiceManager.TransactOpts, operator)
}

// Initialize is a paid mutator transaction binding the contract method 0x775bbcb5.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address[] _batchConfirmers, address _rewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) Initialize(opts *bind.TransactOpts, _pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _batchConfirmers []common.Address, _rewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "initialize", _pauserRegistry, _initialPausedStatus, _initialOwner, _batchConfirmers, _rewardsInitiator)
}

// Initialize is a paid mutator transaction binding the contract method 0x775bbcb5.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address[] _batchConfirmers, address _rewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Initialize(_pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _batchConfirmers []common.Address, _rewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Initialize(&_ContractEigenDAServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus, _initialOwner, _batchConfirmers, _rewardsInitiator)
}

// Initialize is a paid mutator transaction binding the contract method 0x775bbcb5.
//
// Solidity: function initialize(address _pauserRegistry, uint256 _initialPausedStatus, address _initialOwner, address[] _batchConfirmers, address _rewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) Initialize(_pauserRegistry common.Address, _initialPausedStatus *big.Int, _initialOwner common.Address, _batchConfirmers []common.Address, _rewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Initialize(&_ContractEigenDAServiceManager.TransactOpts, _pauserRegistry, _initialPausedStatus, _initialOwner, _batchConfirmers, _rewardsInitiator)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Pause(&_ContractEigenDAServiceManager.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Pause(&_ContractEigenDAServiceManager.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) PauseAll() (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.PauseAll(&_ContractEigenDAServiceManager.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) PauseAll() (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.PauseAll(&_ContractEigenDAServiceManager.TransactOpts)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.RegisterOperatorToAVS(&_ContractEigenDAServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.RegisterOperatorToAVS(&_ContractEigenDAServiceManager.TransactOpts, operator, operatorSignature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.RenounceOwnership(&_ContractEigenDAServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.RenounceOwnership(&_ContractEigenDAServiceManager.TransactOpts)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) SetBatchConfirmer(opts *bind.TransactOpts, _batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "setBatchConfirmer", _batchConfirmer)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) SetBatchConfirmer(_batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetBatchConfirmer(&_ContractEigenDAServiceManager.TransactOpts, _batchConfirmer)
}

// SetBatchConfirmer is a paid mutator transaction binding the contract method 0xf1220983.
//
// Solidity: function setBatchConfirmer(address _batchConfirmer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) SetBatchConfirmer(_batchConfirmer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetBatchConfirmer(&_ContractEigenDAServiceManager.TransactOpts, _batchConfirmer)
}

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) SetClaimerFor(opts *bind.TransactOpts, claimer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "setClaimerFor", claimer)
}

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) SetClaimerFor(claimer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetClaimerFor(&_ContractEigenDAServiceManager.TransactOpts, claimer)
}

// SetClaimerFor is a paid mutator transaction binding the contract method 0xa0169ddd.
//
// Solidity: function setClaimerFor(address claimer) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) SetClaimerFor(claimer common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetClaimerFor(&_ContractEigenDAServiceManager.TransactOpts, claimer)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetPauserRegistry(&_ContractEigenDAServiceManager.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetPauserRegistry(&_ContractEigenDAServiceManager.TransactOpts, newPauserRegistry)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) SetRewardsInitiator(opts *bind.TransactOpts, newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "setRewardsInitiator", newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetRewardsInitiator(&_ContractEigenDAServiceManager.TransactOpts, newRewardsInitiator)
}

// SetRewardsInitiator is a paid mutator transaction binding the contract method 0x3bc28c8c.
//
// Solidity: function setRewardsInitiator(address newRewardsInitiator) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) SetRewardsInitiator(newRewardsInitiator common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetRewardsInitiator(&_ContractEigenDAServiceManager.TransactOpts, newRewardsInitiator)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) SetStaleStakesForbidden(opts *bind.TransactOpts, value bool) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "setStaleStakesForbidden", value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetStaleStakesForbidden(&_ContractEigenDAServiceManager.TransactOpts, value)
}

// SetStaleStakesForbidden is a paid mutator transaction binding the contract method 0x416c7e5e.
//
// Solidity: function setStaleStakesForbidden(bool value) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) SetStaleStakesForbidden(value bool) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.SetStaleStakesForbidden(&_ContractEigenDAServiceManager.TransactOpts, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.TransferOwnership(&_ContractEigenDAServiceManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.TransferOwnership(&_ContractEigenDAServiceManager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Unpause(&_ContractEigenDAServiceManager.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.Unpause(&_ContractEigenDAServiceManager.TransactOpts, newPausedStatus)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) UpdateAVSMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "updateAVSMetadataURI", _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.UpdateAVSMetadataURI(&_ContractEigenDAServiceManager.TransactOpts, _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.UpdateAVSMetadataURI(&_ContractEigenDAServiceManager.TransactOpts, _metadataURI)
}

// ContractEigenDAServiceManagerBatchConfirmedIterator is returned from FilterBatchConfirmed and is used to iterate over the raw logs and unpacked data for BatchConfirmed events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerBatchConfirmedIterator struct {
	Event *ContractEigenDAServiceManagerBatchConfirmed // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerBatchConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerBatchConfirmed)
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
		it.Event = new(ContractEigenDAServiceManagerBatchConfirmed)
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
func (it *ContractEigenDAServiceManagerBatchConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerBatchConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerBatchConfirmed represents a BatchConfirmed event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerBatchConfirmed struct {
	BatchHeaderHash [32]byte
	BatchId         uint32
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmed is a free log retrieval operation binding the contract event 0xc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterBatchConfirmed(opts *bind.FilterOpts, batchHeaderHash [][32]byte) (*ContractEigenDAServiceManagerBatchConfirmedIterator, error) {

	var batchHeaderHashRule []interface{}
	for _, batchHeaderHashItem := range batchHeaderHash {
		batchHeaderHashRule = append(batchHeaderHashRule, batchHeaderHashItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "BatchConfirmed", batchHeaderHashRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerBatchConfirmedIterator{contract: _ContractEigenDAServiceManager.contract, event: "BatchConfirmed", logs: logs, sub: sub}, nil
}

// WatchBatchConfirmed is a free log subscription operation binding the contract event 0xc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a.
//
// Solidity: event BatchConfirmed(bytes32 indexed batchHeaderHash, uint32 batchId)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchBatchConfirmed(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerBatchConfirmed, batchHeaderHash [][32]byte) (event.Subscription, error) {

	var batchHeaderHashRule []interface{}
	for _, batchHeaderHashItem := range batchHeaderHash {
		batchHeaderHashRule = append(batchHeaderHashRule, batchHeaderHashItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "BatchConfirmed", batchHeaderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerBatchConfirmed)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmed", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseBatchConfirmed(log types.Log) (*ContractEigenDAServiceManagerBatchConfirmed, error) {
	event := new(ContractEigenDAServiceManagerBatchConfirmed)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator is returned from FilterBatchConfirmerStatusChanged and is used to iterate over the raw logs and unpacked data for BatchConfirmerStatusChanged events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator struct {
	Event *ContractEigenDAServiceManagerBatchConfirmerStatusChanged // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerBatchConfirmerStatusChanged)
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
		it.Event = new(ContractEigenDAServiceManagerBatchConfirmerStatusChanged)
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
func (it *ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerBatchConfirmerStatusChanged represents a BatchConfirmerStatusChanged event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerBatchConfirmerStatusChanged struct {
	BatchConfirmer common.Address
	Status         bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBatchConfirmerStatusChanged is a free log retrieval operation binding the contract event 0x5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d414.
//
// Solidity: event BatchConfirmerStatusChanged(address batchConfirmer, bool status)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterBatchConfirmerStatusChanged(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "BatchConfirmerStatusChanged")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerBatchConfirmerStatusChangedIterator{contract: _ContractEigenDAServiceManager.contract, event: "BatchConfirmerStatusChanged", logs: logs, sub: sub}, nil
}

// WatchBatchConfirmerStatusChanged is a free log subscription operation binding the contract event 0x5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d414.
//
// Solidity: event BatchConfirmerStatusChanged(address batchConfirmer, bool status)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchBatchConfirmerStatusChanged(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerBatchConfirmerStatusChanged) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "BatchConfirmerStatusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerBatchConfirmerStatusChanged)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerStatusChanged", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseBatchConfirmerStatusChanged(log types.Log) (*ContractEigenDAServiceManagerBatchConfirmerStatusChanged, error) {
	event := new(ContractEigenDAServiceManagerBatchConfirmerStatusChanged)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "BatchConfirmerStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator is returned from FilterDefaultSecurityThresholdsV2Updated and is used to iterate over the raw logs and unpacked data for DefaultSecurityThresholdsV2Updated events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator struct {
	Event *ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
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
		it.Event = new(ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
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
func (it *ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated represents a DefaultSecurityThresholdsV2Updated event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated struct {
	PreviousDefaultSecurityThresholdsV2 SecurityThresholds
	NewDefaultSecurityThresholdsV2      SecurityThresholds
	Raw                                 types.Log // Blockchain specific contextual infos
}

// FilterDefaultSecurityThresholdsV2Updated is a free log retrieval operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterDefaultSecurityThresholdsV2Updated(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerDefaultSecurityThresholdsV2UpdatedIterator{contract: _ContractEigenDAServiceManager.contract, event: "DefaultSecurityThresholdsV2Updated", logs: logs, sub: sub}, nil
}

// WatchDefaultSecurityThresholdsV2Updated is a free log subscription operation binding the contract event 0xfe03afd62c76a6aed7376ae995cc55d073ba9d83d83ac8efc5446f8da4d50997.
//
// Solidity: event DefaultSecurityThresholdsV2Updated((uint8,uint8) previousDefaultSecurityThresholdsV2, (uint8,uint8) newDefaultSecurityThresholdsV2)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchDefaultSecurityThresholdsV2Updated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "DefaultSecurityThresholdsV2Updated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseDefaultSecurityThresholdsV2Updated(log types.Log) (*ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated, error) {
	event := new(ContractEigenDAServiceManagerDefaultSecurityThresholdsV2Updated)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "DefaultSecurityThresholdsV2Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerInitializedIterator struct {
	Event *ContractEigenDAServiceManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerInitialized)
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
		it.Event = new(ContractEigenDAServiceManagerInitialized)
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
func (it *ContractEigenDAServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerInitialized represents a Initialized event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerInitializedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerInitializedIterator{contract: _ContractEigenDAServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerInitialized)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseInitialized(log types.Log) (*ContractEigenDAServiceManagerInitialized, error) {
	event := new(ContractEigenDAServiceManagerInitialized)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerOwnershipTransferredIterator struct {
	Event *ContractEigenDAServiceManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerOwnershipTransferred)
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
		it.Event = new(ContractEigenDAServiceManagerOwnershipTransferred)
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
func (it *ContractEigenDAServiceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractEigenDAServiceManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerOwnershipTransferredIterator{contract: _ContractEigenDAServiceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerOwnershipTransferred)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractEigenDAServiceManagerOwnershipTransferred, error) {
	event := new(ContractEigenDAServiceManagerOwnershipTransferred)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerPausedIterator struct {
	Event *ContractEigenDAServiceManagerPaused // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerPaused)
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
		it.Event = new(ContractEigenDAServiceManagerPaused)
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
func (it *ContractEigenDAServiceManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerPaused represents a Paused event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*ContractEigenDAServiceManagerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerPausedIterator{contract: _ContractEigenDAServiceManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerPaused)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParsePaused(log types.Log) (*ContractEigenDAServiceManagerPaused, error) {
	event := new(ContractEigenDAServiceManagerPaused)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerPauserRegistrySetIterator struct {
	Event *ContractEigenDAServiceManagerPauserRegistrySet // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerPauserRegistrySet)
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
		it.Event = new(ContractEigenDAServiceManagerPauserRegistrySet)
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
func (it *ContractEigenDAServiceManagerPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerPauserRegistrySet represents a PauserRegistrySet event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerPauserRegistrySetIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerPauserRegistrySetIterator{contract: _ContractEigenDAServiceManager.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerPauserRegistrySet)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParsePauserRegistrySet(log types.Log) (*ContractEigenDAServiceManagerPauserRegistrySet, error) {
	event := new(ContractEigenDAServiceManagerPauserRegistrySet)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator is returned from FilterQuorumAdversaryThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumAdversaryThresholdPercentagesUpdated events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
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
func (it *ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated represents a QuorumAdversaryThresholdPercentagesUpdated event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated struct {
	PreviousQuorumAdversaryThresholdPercentages []byte
	NewQuorumAdversaryThresholdPercentages      []byte
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterQuorumAdversaryThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterQuorumAdversaryThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAServiceManager.contract, event: "QuorumAdversaryThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumAdversaryThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0xf73542111561dc551cbbe9111c4dd3a040d53d7bc0339a53290f4d7f9a95c3cc.
//
// Solidity: event QuorumAdversaryThresholdPercentagesUpdated(bytes previousQuorumAdversaryThresholdPercentages, bytes newQuorumAdversaryThresholdPercentages)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchQuorumAdversaryThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "QuorumAdversaryThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseQuorumAdversaryThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAServiceManagerQuorumAdversaryThresholdPercentagesUpdated)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumAdversaryThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator is returned from FilterQuorumConfirmationThresholdPercentagesUpdated and is used to iterate over the raw logs and unpacked data for QuorumConfirmationThresholdPercentagesUpdated events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator struct {
	Event *ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
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
		it.Event = new(ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
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
func (it *ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated represents a QuorumConfirmationThresholdPercentagesUpdated event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated struct {
	PreviousQuorumConfirmationThresholdPercentages []byte
	NewQuorumConfirmationThresholdPercentages      []byte
	Raw                                            types.Log // Blockchain specific contextual infos
}

// FilterQuorumConfirmationThresholdPercentagesUpdated is a free log retrieval operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterQuorumConfirmationThresholdPercentagesUpdated(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdatedIterator{contract: _ContractEigenDAServiceManager.contract, event: "QuorumConfirmationThresholdPercentagesUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumConfirmationThresholdPercentagesUpdated is a free log subscription operation binding the contract event 0x9f1ea99a8363f2964c53c763811648354a8437441b30b39465f9d26118d6a5a0.
//
// Solidity: event QuorumConfirmationThresholdPercentagesUpdated(bytes previousQuorumConfirmationThresholdPercentages, bytes newQuorumConfirmationThresholdPercentages)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchQuorumConfirmationThresholdPercentagesUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "QuorumConfirmationThresholdPercentagesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseQuorumConfirmationThresholdPercentagesUpdated(log types.Log) (*ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated, error) {
	event := new(ContractEigenDAServiceManagerQuorumConfirmationThresholdPercentagesUpdated)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumConfirmationThresholdPercentagesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator is returned from FilterQuorumNumbersRequiredUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumbersRequiredUpdated events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator struct {
	Event *ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated)
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
		it.Event = new(ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated)
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
func (it *ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated represents a QuorumNumbersRequiredUpdated event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated struct {
	PreviousQuorumNumbersRequired []byte
	NewQuorumNumbersRequired      []byte
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumbersRequiredUpdated is a free log retrieval operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterQuorumNumbersRequiredUpdated(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerQuorumNumbersRequiredUpdatedIterator{contract: _ContractEigenDAServiceManager.contract, event: "QuorumNumbersRequiredUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumbersRequiredUpdated is a free log subscription operation binding the contract event 0x60c0ba1da794fcbbf549d370512442cb8f3f3f774cb557205cc88c6f842cb36a.
//
// Solidity: event QuorumNumbersRequiredUpdated(bytes previousQuorumNumbersRequired, bytes newQuorumNumbersRequired)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchQuorumNumbersRequiredUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "QuorumNumbersRequiredUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseQuorumNumbersRequiredUpdated(log types.Log) (*ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated, error) {
	event := new(ContractEigenDAServiceManagerQuorumNumbersRequiredUpdated)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "QuorumNumbersRequiredUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator is returned from FilterRewardsInitiatorUpdated and is used to iterate over the raw logs and unpacked data for RewardsInitiatorUpdated events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator struct {
	Event *ContractEigenDAServiceManagerRewardsInitiatorUpdated // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerRewardsInitiatorUpdated)
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
		it.Event = new(ContractEigenDAServiceManagerRewardsInitiatorUpdated)
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
func (it *ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerRewardsInitiatorUpdated represents a RewardsInitiatorUpdated event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerRewardsInitiatorUpdated struct {
	PrevRewardsInitiator common.Address
	NewRewardsInitiator  common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRewardsInitiatorUpdated is a free log retrieval operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterRewardsInitiatorUpdated(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerRewardsInitiatorUpdatedIterator{contract: _ContractEigenDAServiceManager.contract, event: "RewardsInitiatorUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardsInitiatorUpdated is a free log subscription operation binding the contract event 0xe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3.
//
// Solidity: event RewardsInitiatorUpdated(address prevRewardsInitiator, address newRewardsInitiator)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchRewardsInitiatorUpdated(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerRewardsInitiatorUpdated) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "RewardsInitiatorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerRewardsInitiatorUpdated)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseRewardsInitiatorUpdated(log types.Log) (*ContractEigenDAServiceManagerRewardsInitiatorUpdated, error) {
	event := new(ContractEigenDAServiceManagerRewardsInitiatorUpdated)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "RewardsInitiatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator is returned from FilterStaleStakesForbiddenUpdate and is used to iterate over the raw logs and unpacked data for StaleStakesForbiddenUpdate events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator struct {
	Event *ContractEigenDAServiceManagerStaleStakesForbiddenUpdate // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerStaleStakesForbiddenUpdate)
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
		it.Event = new(ContractEigenDAServiceManagerStaleStakesForbiddenUpdate)
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
func (it *ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerStaleStakesForbiddenUpdate represents a StaleStakesForbiddenUpdate event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerStaleStakesForbiddenUpdate struct {
	Value bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterStaleStakesForbiddenUpdate is a free log retrieval operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterStaleStakesForbiddenUpdate(opts *bind.FilterOpts) (*ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerStaleStakesForbiddenUpdateIterator{contract: _ContractEigenDAServiceManager.contract, event: "StaleStakesForbiddenUpdate", logs: logs, sub: sub}, nil
}

// WatchStaleStakesForbiddenUpdate is a free log subscription operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchStaleStakesForbiddenUpdate(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerStaleStakesForbiddenUpdate) (event.Subscription, error) {

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "StaleStakesForbiddenUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerStaleStakesForbiddenUpdate)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
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

// ParseStaleStakesForbiddenUpdate is a log parse operation binding the contract event 0x40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc.
//
// Solidity: event StaleStakesForbiddenUpdate(bool value)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseStaleStakesForbiddenUpdate(log types.Log) (*ContractEigenDAServiceManagerStaleStakesForbiddenUpdate, error) {
	event := new(ContractEigenDAServiceManagerStaleStakesForbiddenUpdate)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "StaleStakesForbiddenUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerUnpausedIterator struct {
	Event *ContractEigenDAServiceManagerUnpaused // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerUnpaused)
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
		it.Event = new(ContractEigenDAServiceManagerUnpaused)
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
func (it *ContractEigenDAServiceManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerUnpaused represents a Unpaused event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*ContractEigenDAServiceManagerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerUnpausedIterator{contract: _ContractEigenDAServiceManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerUnpaused)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseUnpaused(log types.Log) (*ContractEigenDAServiceManagerUnpaused, error) {
	event := new(ContractEigenDAServiceManagerUnpaused)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator is returned from FilterVersionedBlobParamsAdded and is used to iterate over the raw logs and unpacked data for VersionedBlobParamsAdded events raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator struct {
	Event *ContractEigenDAServiceManagerVersionedBlobParamsAdded // Event containing the contract specifics and raw log

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
func (it *ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractEigenDAServiceManagerVersionedBlobParamsAdded)
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
		it.Event = new(ContractEigenDAServiceManagerVersionedBlobParamsAdded)
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
func (it *ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractEigenDAServiceManagerVersionedBlobParamsAdded represents a VersionedBlobParamsAdded event raised by the ContractEigenDAServiceManager contract.
type ContractEigenDAServiceManagerVersionedBlobParamsAdded struct {
	Version             uint16
	VersionedBlobParams VersionedBlobParams
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterVersionedBlobParamsAdded is a free log retrieval operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) FilterVersionedBlobParamsAdded(opts *bind.FilterOpts, version []uint16) (*ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.FilterLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return &ContractEigenDAServiceManagerVersionedBlobParamsAddedIterator{contract: _ContractEigenDAServiceManager.contract, event: "VersionedBlobParamsAdded", logs: logs, sub: sub}, nil
}

// WatchVersionedBlobParamsAdded is a free log subscription operation binding the contract event 0xdbee9d337a6e5fde30966e157673aaeeb6a0134afaf774a4b6979b7c79d07da4.
//
// Solidity: event VersionedBlobParamsAdded(uint16 indexed version, (uint32,uint32,uint8) versionedBlobParams)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) WatchVersionedBlobParamsAdded(opts *bind.WatchOpts, sink chan<- *ContractEigenDAServiceManagerVersionedBlobParamsAdded, version []uint16) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _ContractEigenDAServiceManager.contract.WatchLogs(opts, "VersionedBlobParamsAdded", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractEigenDAServiceManagerVersionedBlobParamsAdded)
				if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
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
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerFilterer) ParseVersionedBlobParamsAdded(log types.Log) (*ContractEigenDAServiceManagerVersionedBlobParamsAdded, error) {
	event := new(ContractEigenDAServiceManagerVersionedBlobParamsAdded)
	if err := _ContractEigenDAServiceManager.contract.UnpackLog(event, "VersionedBlobParamsAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
