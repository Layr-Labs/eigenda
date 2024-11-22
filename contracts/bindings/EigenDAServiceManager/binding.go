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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"__avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"__rewardsCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRewardsCoordinator\"},{\"name\":\"__registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"__stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"},{\"name\":\"__eigenDAThresholdRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"},{\"name\":\"__eigenDARelayRegistry\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STORE_DURATION_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"THRESHOLD_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkSignatures\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"components\":[{\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structBatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAVSRewardsSubmission\",\"inputs\":[{\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eigenDARelayRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDARelayRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenDAThresholdRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenDAThresholdRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlobParams\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVersionedBlobParams\",\"components\":[{\"name\":\"maxNumOperators\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"numChunks\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"codingRate\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDefaultSecurityThresholdsV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSecurityThresholds\",\"components\":[{\"name\":\"confirmationThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThreshold\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchConfirmers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_rewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isBatchConfirmer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestServeUntilBlock\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardsInitiator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBatchConfirmer\",\"inputs\":[{\"name\":\"_batchConfirmer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsInitiator\",\"inputs\":[{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStaleStakesForbidden\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staleStakesForbidden\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trySignatureAndApkVerification\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"apk\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"pairingSuccessful\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"siganatureIsValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchConfirmerStatusChanged\",\"inputs\":[{\"name\":\"batchConfirmer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsInitiatorUpdated\",\"inputs\":[{\"name\":\"prevRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleStakesForbiddenUpdate\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6101c06040523480156200001257600080fd5b5060405162005dc438038062005dc48339810160408190526200003591620002f5565b6001600160a01b0380831660805280821660a05280871660c05280861660e052808516610100528316610120528386868286620000716200021a565b50505050806001600160a01b0316610140816001600160a01b031681525050806001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000cf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000f5919062000389565b6001600160a01b0316610160816001600160a01b031681525050806001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200014e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000174919062000389565b6001600160a01b0316610180816001600160a01b031681525050610160516001600160a01b031663df5cf7236040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001d0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001f6919062000389565b6001600160a01b03166101a052506200020e6200021a565b505050505050620003b0565b603254610100900460ff1615620002875760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff9081161015620002da576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b0381168114620002f257600080fd5b50565b60008060008060008060c087890312156200030f57600080fd5b86516200031c81620002dc565b60208801519096506200032f81620002dc565b60408801519095506200034281620002dc565b60608801519094506200035581620002dc565b60808801519093506200036881620002dc565b60a08801519092506200037b81620002dc565b809150509295509295509295565b6000602082840312156200039c57600080fd5b8151620003a981620002dc565b9392505050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a0516158cd620004f7600039600081816106520152611aa90152600081816104340152611c8b01526000818161048601528181611e6101526120230152600081816104d3015281816111c5015281816117740152818161190c0152611b46015260008181610edd01528181611038015281816110cf01528181612c9401528181612e170152612eb6015260008181610d0801528181610d9701528181610e17015281816129f801528181612abc01528181612bd20152612d720152600081816133cc0152818161348801526135740152600081816104aa01528181612a4c01528181612b180152612b970152600061052301526000818161073b015281816107b101528181610a4001528181610c7801528181612fcc015261301901526158cd6000f3fe608060405234801561001057600080fd5b50600436106102a05760003560e01c80637794965a11610167578063e481af9d116100ce578063f122098311610087578063f122098314610710578063f2fde38b14610723578063f8c6681414610736578063fabc1cbc1461075d578063fc299dee14610770578063fce36c7d1461078357600080fd5b8063e481af9d14610691578063eaefd27d14610699578063eccbbfc9146106ac578063ee6c3bcf146106cc578063ef024458146106df578063ef635529146106e757600080fd5b8063a5b7890a11610120578063a5b7890a146105eb578063a98fb3551461060e578063b98d090814610621578063bafa91071461062e578063df5cf7231461064d578063e15234ff1461067457600080fd5b80637794965a146105665780638687feae14610579578063886f1195146105a15780638da5cb5b146105b45780639926ee7d146105c5578063a364f4da146105d857600080fd5b80635c975abb1161020b5780636d14a987116101c45780636d14a987146104ce5780636efb4636146104f5578063715018a614610516578063722764431461051e57806372d18e8d14610545578063775bbcb51461055357600080fd5b80635c975abb1461041d5780635df459461461042f5780635e0334761461046e5780635e8b3f2d1461047857806368304835146104815780636b3aa72e146104a857600080fd5b806333cfb7b71161025d57806333cfb7b7146103875780633bc28c8c146103a7578063416c7e5e146103ba5780634972134a146103cd578063595c6a67146103f25780635ac86ab7146103fa57600080fd5b8063048886d2146102a557806310d67a2f146102cd578063136439dd146102e25780631429c7c2146102f5578063171f1d5b1461031a5780632ecfe72b14610344575b600080fd5b6102b86102b33660046145ab565b610796565b60405190151581526020015b60405180910390f35b6102e06102db3660046145dd565b61082a565b005b6102e06102f03660046145fa565b6108e6565b6103086103033660046145ab565b610a25565b60405160ff90911681526020016102c4565b61032d610328366004614764565b610ab4565b6040805192151583529015156020830152016102c4565b6103576103523660046147b5565b610c3e565b60408051825163ffffffff9081168252602080850151909116908201529181015160ff16908201526060016102c4565b61039a6103953660046145dd565b610ce3565b6040516102c491906147e4565b6102e06103b53660046145dd565b6111b2565b6102e06103c836600461483f565b6111c3565b6000546103dd9063ffffffff1681565b60405163ffffffff90911681526020016102c4565b6102e06112fa565b6102b86104083660046145ab565b60fc54600160ff9092169190911b9081161490565b60fc545b6040519081526020016102c4565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016102c4565b6103dd620189c081565b6103dd61012c81565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b7f0000000000000000000000000000000000000000000000000000000000000000610456565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b610508610503366004614b1d565b6113c1565b6040516102c4929190614c10565b6102e06122d8565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b60005463ffffffff166103dd565b6102e0610561366004614c59565b6122ec565b6102e0610574366004614d34565b612455565b60408051808201909152600381526221212160e81b60208201525b6040516102c49190614dec565b60fb54610456906001600160a01b031681565b6065546001600160a01b0316610456565b6102e06105d3366004614e76565b6129ed565b6102e06105e63660046145dd565b612ab1565b6102b86105f93660046145dd565b60026020526000908152604090205460ff1681565b6102e061061c366004614f21565b612b78565b60c9546102b89060ff1681565b60408051808201909152600381526237373760e81b6020820152610594565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b6040805180820190915260028152600160f01b6020820152610594565b61039a612bcc565b6103dd6106a7366004614f71565b612f95565b6104216106ba366004614f71565b60016020526000908152604090205481565b6103086106da3660046145ab565b612fb1565b610421606481565b6106ef613003565b60408051825160ff90811682526020938401511692810192909252016102c4565b6102e061071e3660046145dd565b61309d565b6102e06107313660046145dd565b6130ae565b6104567f000000000000000000000000000000000000000000000000000000000000000081565b6102e061076b3660046145fa565b613124565b609754610456906001600160a01b031681565b6102e0610791366004614f8e565b613280565b604051630244436960e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063048886d290602401602060405180830381865afa158015610800573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108249190615002565b92915050565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561087d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a1919061501f565b6001600160a01b0316336001600160a01b0316146108da5760405162461bcd60e51b81526004016108d19061503c565b60405180910390fd5b6108e3816135ab565b50565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa15801561092e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109529190615002565b61096e5760405162461bcd60e51b81526004016108d190615086565b60fc54818116146109e75760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c697479000000000000000060648201526084016108d1565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d906020015b60405180910390a250565b604051630a14e3e160e11b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690631429c7c2906024015b602060405180830381865afa158015610a90573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061082491906150ce565b60008060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000187876000015188602001518860000151600060028110610afc57610afc6150eb565b60200201518951600160200201518a60200151600060028110610b2157610b216150eb565b60200201518b60200151600160028110610b3d57610b3d6150eb565b602090810291909101518c518d830151604051610b9a9a99989796959401988952602089019790975260408801959095526060870193909352608086019190915260a085015260c084015260e08301526101008201526101200190565b6040516020818303038152906040528051906020012060001c610bbd9190615101565b9050610c30610bd6610bcf88846136a2565b8690613739565b610bde6137cd565b610c26610c1785610c11604080518082018252600080825260209182015281518083019092526001825260029082015290565b906136a2565b610c208c61388d565b90613739565b886201d4c061391d565b909890975095505050505050565b60408051606081018252600080825260208201819052818301529051632ecfe72b60e01b815261ffff831660048201526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632ecfe72b90602401606060405180830381865afa158015610cbf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108249190615123565b6040516309aa152760e11b81526001600160a01b0382811660048301526060916000917f000000000000000000000000000000000000000000000000000000000000000016906313542a4e90602401602060405180830381865afa158015610d4f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d739190615194565b60405163871ef04960e01b8152600481018290529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063871ef04990602401602060405180830381865afa158015610dde573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e0291906151ad565b90506001600160c01b0381161580610e9c57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e73573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e9791906150ce565b60ff16155b15610eb857505060408051600081526020810190915292915050565b6000610ecc826001600160c01b0316613b41565b90506000805b8251811015610fa2577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316633ca5a5f5848381518110610f1c57610f1c6150eb565b01602001516040516001600160e01b031960e084901b16815260f89190911c6004820152602401602060405180830381865afa158015610f60573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f849190615194565b610f8e90836151ec565b915080610f9a81615204565b915050610ed2565b506000816001600160401b03811115610fbd57610fbd614613565b604051908082528060200260200182016040528015610fe6578160200160208202803683370190505b5090506000805b84518110156111a557600085828151811061100a5761100a6150eb565b0160200151604051633ca5a5f560e01b815260f89190911c6004820181905291506000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690633ca5a5f590602401602060405180830381865afa15801561107f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110a39190615194565b905060005b8181101561118f576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa15801561111d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111419190615234565b60000151868681518110611157576111576150eb565b6001600160a01b03909216602092830291909101909101528461117981615204565b955050808061118790615204565b9150506110a8565b505050808061119d90615204565b915050610fed565b5090979650505050505050565b6111ba613c03565b6108e381613c5d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611221573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611245919061501f565b6001600160a01b0316336001600160a01b0316146112f15760405162461bcd60e51b815260206004820152605c60248201527f424c535369676e6174757265436865636b65722e6f6e6c79436f6f7264696e6160448201527f746f724f776e65723a2063616c6c6572206973206e6f7420746865206f776e6560648201527f72206f6620746865207265676973747279436f6f7264696e61746f7200000000608482015260a4016108d1565b6108e381613cc6565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015611342573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113669190615002565b6113825760405162461bcd60e51b81526004016108d190615086565b60001960fc81905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b60408051808201909152606080825260208201526000846114385760405162461bcd60e51b8152602060048201526037602482015260008051602061587883398151915260448201527f7265733a20656d7074792071756f72756d20696e70757400000000000000000060648201526084016108d1565b60408301515185148015611450575060a08301515185145b8015611460575060c08301515185145b8015611470575060e08301515185145b6114da5760405162461bcd60e51b8152602060048201526041602482015260008051602061587883398151915260448201527f7265733a20696e7075742071756f72756d206c656e677468206d69736d6174636064820152600d60fb1b608482015260a4016108d1565b825151602084015151146115525760405162461bcd60e51b815260206004820152604460248201819052600080516020615878833981519152908201527f7265733a20696e707574206e6f6e7369676e6572206c656e677468206d69736d6064820152630c2e8c6d60e31b608482015260a4016108d1565b4363ffffffff168463ffffffff16106115c15760405162461bcd60e51b815260206004820152603c602482015260008051602061587883398151915260448201527f7265733a20696e76616c6964207265666572656e636520626c6f636b0000000060648201526084016108d1565b6040805180820182526000808252602080830191909152825180840190935260608084529083015290866001600160401b0381111561160257611602614613565b60405190808252806020026020018201604052801561162b578160200160208202803683370190505b506020820152866001600160401b0381111561164957611649614613565b604051908082528060200260200182016040528015611672578160200160208202803683370190505b50815260408051808201909152606080825260208201528560200151516001600160401b038111156116a6576116a6614613565b6040519080825280602002602001820160405280156116cf578160200160208202803683370190505b5081526020860151516001600160401b038111156116ef576116ef614613565b604051908082528060200260200182016040528015611718578160200160208202803683370190505b50816020018190525060006117ea8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051639aa1653d60e01b815290516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169350639aa1653d925060048083019260209291908290030181865afa1580156117c1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117e591906150ce565b613d0e565b905060005b876020015151811015611a855761183488602001518281518110611815576118156150eb565b6020026020010151805160009081526020918201519091526040902090565b8360200151828151811061184a5761184a6150eb565b6020908102919091010152801561190a57602083015161186b600183615275565b8151811061187b5761187b6150eb565b602002602001015160001c8360200151828151811061189c5761189c6150eb565b602002602001015160001c1161190a576040805162461bcd60e51b815260206004820152602481019190915260008051602061587883398151915260448201527f7265733a206e6f6e5369676e65725075626b657973206e6f7420736f7274656460648201526084016108d1565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166304ec63518460200151838151811061194f5761194f6150eb565b60200260200101518b8b60000151858151811061196e5761196e6150eb565b60200260200101516040518463ffffffff1660e01b81526004016119ab9392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa1580156119c8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119ec91906151ad565b6001600160c01b031683600001518281518110611a0b57611a0b6150eb565b602002602001018181525050611a71610bcf611a458486600001518581518110611a3757611a376150eb565b602002602001015116613d9f565b8a602001518481518110611a5b57611a5b6150eb565b6020026020010151613dca90919063ffffffff16565b945080611a7d81615204565b9150506117ef565b5050611a9083613eae565b60c95490935060ff16600081611aa7576000611b29565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c448feb86040518163ffffffff1660e01b8152600401602060405180830381865afa158015611b05573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b299190615194565b905060005b8a8110156121a7578215611c89578963ffffffff16827f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663249a0c428f8f86818110611b8557611b856150eb565b60405160e085901b6001600160e01b031916815292013560f81c600483015250602401602060405180830381865afa158015611bc5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611be99190615194565b611bf391906151ec565b11611c895760405162461bcd60e51b8152602060048201526066602482015260008051602061587883398151915260448201527f7265733a205374616b6552656769737472792075706461746573206d7573742060648201527f62652077697468696e207769746864726177616c44656c6179426c6f636b732060848201526577696e646f7760d01b60a482015260c4016108d1565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166368bccaac8d8d84818110611cca57611cca6150eb565b9050013560f81c60f81b60f81c8c8c60a001518581518110611cee57611cee6150eb565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015611d4a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d6e919061528c565b6001600160401b031916611d918a604001518381518110611815576118156150eb565b67ffffffffffffffff191614611e2d5760405162461bcd60e51b8152602060048201526061602482015260008051602061587883398151915260448201527f7265733a2071756f72756d41706b206861736820696e2073746f72616765206460648201527f6f6573206e6f74206d617463682070726f76696465642071756f72756d2061706084820152606b60f81b60a482015260c4016108d1565b611e5d89604001518281518110611e4657611e466150eb565b60200260200101518761373990919063ffffffff16565b95507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c8294c568d8d84818110611ea057611ea06150eb565b9050013560f81c60f81b60f81c8c8c60c001518581518110611ec457611ec46150eb565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015611f20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f4491906152b7565b85602001518281518110611f5a57611f5a6150eb565b6001600160601b03909216602092830291909101820152850151805182908110611f8657611f866150eb565b602002602001015185600001518281518110611fa457611fa46150eb565b60200260200101906001600160601b031690816001600160601b0316815250506000805b8a60200151518110156121925761201c86600001518281518110611fee57611fee6150eb565b60200260200101518f8f86818110612008576120086150eb565b600192013560f81c9290921c811614919050565b15612180577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663f2be94ae8f8f86818110612062576120626150eb565b9050013560f81c60f81b60f81c8e89602001518581518110612086576120866150eb565b60200260200101518f60e0015188815181106120a4576120a46150eb565b602002602001015187815181106120bd576120bd6150eb565b60209081029190910101516040516001600160e01b031960e087901b16815260ff909416600485015263ffffffff92831660248501526044840191909152166064820152608401602060405180830381865afa158015612121573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061214591906152b7565b8751805185908110612159576121596150eb565b6020026020010181815161216d91906152d4565b6001600160601b03169052506001909101905b8061218a81615204565b915050611fc8565b5050808061219f90615204565b915050611b2e565b5050506000806121c18c868a606001518b60800151610ab4565b91509150816122325760405162461bcd60e51b8152602060048201526043602482015260008051602061587883398151915260448201527f7265733a2070616972696e6720707265636f6d70696c652063616c6c206661696064820152621b195960ea1b608482015260a4016108d1565b806122935760405162461bcd60e51b8152602060048201526039602482015260008051602061587883398151915260448201527f7265733a207369676e617475726520697320696e76616c69640000000000000060648201526084016108d1565b505060008782602001516040516020016122ae9291906152fc565b60408051808303601f190181529190528051602090910120929b929a509198505050505050505050565b6122e0613c03565b6122ea6000613f49565b565b603254610100900460ff161580801561230c5750603254600160ff909116105b806123265750303b158015612326575060325460ff166001145b6123895760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016108d1565b6032805460ff1916600117905580156123ac576032805461ff0019166101001790555b6123b68686613f9b565b6123bf84613f49565b6123c882613c5d565b60005b8351811015612406576123f68482815181106123e9576123e96150eb565b6020026020010151614085565b6123ff81615204565b90506123cb565b50801561244d576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050565b60fc54600090600190811614156124ae5760405162461bcd60e51b815260206004820152601960248201527f5061757361626c653a20696e646578206973207061757365640000000000000060448201526064016108d1565b3360009081526002602052604090205460ff166125225760405162461bcd60e51b815260206004820152602c60248201527f6f6e6c794261746368436f6e6669726d65723a206e6f742066726f6d2062617460448201526b31b41031b7b73334b936b2b960a11b60648201526084016108d1565b32331461259f5760405162461bcd60e51b8152602060048201526051602482015260008051602061585883398151915260448201527f63683a2068656164657220616e64206e6f6e7369676e65722064617461206d75606482015270737420626520696e2063616c6c6461746160781b608482015260a4016108d1565b436125b06080850160608601614f71565b63ffffffff161061262f5760405162461bcd60e51b815260206004820152604f602482015260008051602061585883398151915260448201527f63683a20737065636966696564207265666572656e6365426c6f636b4e756d6260648201526e657220697320696e2066757475726560881b608482015260a4016108d1565b63ffffffff431661012c6126496080860160608701614f71565b6126539190615344565b63ffffffff1610156126d95760405162461bcd60e51b8152602060048201526055602482015260008051602061585883398151915260448201527f63683a20737065636966696564207265666572656e6365426c6f636b4e756d62606482015274195c881a5cc81d1bdbc819985c881a5b881c185cdd605a1b608482015260a4016108d1565b6126e6604084018461536c565b90506126f5602085018561536c565b90501461278d5760405162461bcd60e51b8152602060048201526066602482015260008051602061585883398151915260448201527f63683a2071756f72756d4e756d6265727320616e64207369676e65645374616b60648201527f65466f7251756f72756d73206d757374206265206f66207468652073616d65206084820152650d8cadccee8d60d31b60a482015260c4016108d1565b60006127a061279b856153b9565b6140e8565b90506000806127cc836127b6602089018961536c565b6127c660808b0160608c01614f71565b896113c1565b9150915060005b6127e0604088018861536c565b9050811015612922576127f6604088018861536c565b82818110612806576128066150eb565b9050013560f81c60f81b60f81c60ff168360200151828151811061282c5761282c6150eb565b602002602001015161283e919061545b565b6001600160601b031660648460000151838151811061285f5761285f6150eb565b60200260200101516001600160601b031661287a919061548a565b10156129105760405162461bcd60e51b81526020600482015260646024820181905260008051602061585883398151915260448301527f63683a207369676e61746f7269657320646f206e6f74206f776e206174206c65908201527f617374207468726573686f6c642070657263656e74616765206f6620612071756084820152636f72756d60e01b60a482015260c4016108d1565b8061291a81615204565b9150506127d3565b506000805463ffffffff169061293788614163565b6040805160208082018490528183018790524360e01b6001600160e01b0319166060830152825160448184030181526064830180855281519183019190912063ffffffff881660008181526001909452928590205552905191925086917fc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a9181900360840190a26129c9826001615344565b6000805463ffffffff191663ffffffff929092169190911790555050505050505050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614612a355760405162461bcd60e51b81526004016108d1906154a9565b604051639926ee7d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690639926ee7d90612a839085908590600401615521565b600060405180830381600087803b158015612a9d57600080fd5b505af115801561244d573d6000803e3d6000fd5b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614612af95760405162461bcd60e51b81526004016108d1906154a9565b6040516351b27a6d60e11b81526001600160a01b0382811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063a364f4da906024015b600060405180830381600087803b158015612b5d57600080fd5b505af1158015612b71573d6000803e3d6000fd5b5050505050565b612b80613c03565b60405163a98fb35560e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a98fb35590612b43908490600401614dec565b606060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c2e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c5291906150ce565b60ff16905080612c7057505060408051600081526020810190915290565b6000805b82811015612d2557604051633ca5a5f560e01b815260ff821660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015612ce3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d079190615194565b612d1190836151ec565b915080612d1d81615204565b915050612c74565b506000816001600160401b03811115612d4057612d40614613565b604051908082528060200260200182016040528015612d69578160200160208202803683370190505b5090506000805b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612dce573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612df291906150ce565b60ff16811015612f8b57604051633ca5a5f560e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015612e66573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e8a9190615194565b905060005b81811015612f76576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa158015612f04573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f289190615234565b60000151858581518110612f3e57612f3e6150eb565b6001600160a01b039092166020928302919091019091015283612f6081615204565b9450508080612f6e90615204565b915050612e8f565b50508080612f8390615204565b915050612d70565b5090949350505050565b600061012c612fa7620189c084615344565b6108249190615344565b60405163ee6c3bcf60e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ee6c3bcf90602401610a73565b60408051808201909152600080825260208201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ef6355296040518163ffffffff1660e01b81526004016040805180830381865afa158015613074573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613098919061556c565b905090565b6130a5613c03565b6108e381614085565b6130b6613c03565b6001600160a01b03811661311b5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016108d1565b6108e381613f49565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015613177573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061319b919061501f565b6001600160a01b0316336001600160a01b0316146131cb5760405162461bcd60e51b81526004016108d19061503c565b60fc5419811960fc541916146132495760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c697479000000000000000060648201526084016108d1565b60fc81905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c90602001610a1a565b613288614176565b60005b8181101561355c578282828181106132a5576132a56150eb565b90506020028101906132b791906155a1565b6132c89060408101906020016145dd565b6001600160a01b03166323b872dd33308686868181106132ea576132ea6150eb565b90506020028101906132fc91906155a1565b604080516001600160e01b031960e087901b1681526001600160a01b039485166004820152939092166024840152013560448201526064016020604051808303816000875af1158015613353573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133779190615002565b50600083838381811061338c5761338c6150eb565b905060200281019061339e91906155a1565b6133af9060408101906020016145dd565b604051636eb1769f60e11b81523060048201526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081166024830152919091169063dd62ed3e90604401602060405180830381865afa15801561341d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906134419190615194565b9050838383818110613455576134556150eb565b905060200281019061346791906155a1565b6134789060408101906020016145dd565b6001600160a01b031663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000838787878181106134ba576134ba6150eb565b90506020028101906134cc91906155a1565b604001356134da91906151ec565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015613525573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135499190615002565b50508061355590615204565b905061328b565b5060405163fce36c7d60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fce36c7d90612a83908590859060040161561c565b6001600160a01b0381166136395760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a4016108d1565b60fb54604080516001600160a01b03928316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a160fb80546001600160a01b0319166001600160a01b0392909216919091179055565b60408051808201909152600080825260208201526136be6144c2565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080156136f1576136f3565bfe5b50806137315760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b60448201526064016108d1565b505092915050565b60408051808201909152600080825260208201526137556144e0565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080156136f15750806137315760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b60448201526064016108d1565b6137d56144fe565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b6040805180820190915260008082526020820152600080806138bd60008051602061583883398151915286615101565b90505b6138c98161420b565b9093509150600080516020615838833981519152828309831415613903576040805180820190915290815260208101919091529392505050565b6000805160206158388339815191526001820890506138c0565b60408051808201825286815260208082018690528251808401909352868352820184905260009182919061394f614523565b60005b6002811015613b1457600061396882600661548a565b905084826002811061397c5761397c6150eb565b6020020151518361398e8360006151ec565b600c811061399e5761399e6150eb565b60200201528482600281106139b5576139b56150eb565b602002015160200151838260016139cc91906151ec565b600c81106139dc576139dc6150eb565b60200201528382600281106139f3576139f36150eb565b6020020151515183613a068360026151ec565b600c8110613a1657613a166150eb565b6020020152838260028110613a2d57613a2d6150eb565b6020020151516001602002015183613a468360036151ec565b600c8110613a5657613a566150eb565b6020020152838260028110613a6d57613a6d6150eb565b602002015160200151600060028110613a8857613a886150eb565b602002015183613a998360046151ec565b600c8110613aa957613aa96150eb565b6020020152838260028110613ac057613ac06150eb565b602002015160200151600160028110613adb57613adb6150eb565b602002015183613aec8360056151ec565b600c8110613afc57613afc6150eb565b60200201525080613b0c81615204565b915050613952565b50613b1d614542565b60006020826101808560088cfa9151919c9115159b50909950505050505050505050565b6060600080613b4f84613d9f565b61ffff166001600160401b03811115613b6a57613b6a614613565b6040519080825280601f01601f191660200182016040528015613b94576020820181803683370190505b5090506000805b825182108015613bac575061010081105b15612f8b576001811b935085841615613bf3578060f81b838381518110613bd557613bd56150eb565b60200101906001600160f81b031916908160001a9053508160010191505b613bfc81615204565b9050613b9b565b6065546001600160a01b031633146122ea5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016108d1565b609754604080516001600160a01b03928316815291831660208301527fe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3910160405180910390a1609780546001600160a01b0319166001600160a01b0392909216919091179055565b60c9805460ff19168215159081179091556040519081527f40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc906020015b60405180910390a150565b600080613d1a8461428d565b9050808360ff166001901b11613d985760405162461bcd60e51b815260206004820152603f60248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206269746d61702065786365656473206d61782076616c75650060648201526084016108d1565b9392505050565b6000805b821561082457613db4600184615275565b9092169180613dc281615729565b915050613da3565b60408051808201909152600080825260208201526102008261ffff1610613e265760405162461bcd60e51b815260206004820152601060248201526f7363616c61722d746f6f2d6c6172676560801b60448201526064016108d1565b8161ffff1660011415613e3a575081610824565b6040805180820190915260008082526020820181905284906001905b8161ffff168661ffff1610613ea357600161ffff871660ff83161c81161415613e8657613e838484613739565b93505b613e908384613739565b92506201fffe600192831b169101613e56565b509195945050505050565b60408051808201909152600080825260208201528151158015613ed357506020820151155b15613ef1575050604080518082019091526000808252602082015290565b6040518060400160405280836000015181526020016000805160206158388339815191528460200151613f249190615101565b613f3c90600080516020615838833981519152615275565b905292915050565b919050565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60fb546001600160a01b0316158015613fbc57506001600160a01b03821615155b61403e5760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a4016108d1565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2614081826135ab565b5050565b6001600160a01b038116600081815260026020908152604091829020805460ff8082161560ff1990921682179092558351948552161515908301527f5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d4149101613d03565b600061412582604080518082019091526000808252602082015250604080518082019091528151815260609091015163ffffffff16602082015290565b6040805182516020808301919091529092015163ffffffff16908201526060015b604051602081830303815290604052805190602001209050919050565b60008160405160200161414691906157b9565b6097546001600160a01b031633146122ea5760405162461bcd60e51b815260206004820152604c60248201527f536572766963654d616e61676572426173652e6f6e6c7952657761726473496e60448201527f69746961746f723a2063616c6c6572206973206e6f742074686520726577617260648201526b32399034b734ba34b0ba37b960a11b608482015260a4016108d1565b60008080600080516020615838833981519152600360008051602061583883398151915286600080516020615838833981519152888909090890506000614281827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f5260008051602061583883398151915261441a565b91959194509092505050565b6000610100825111156143165760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a4016108d1565b815161432457506000919050565b6000808360008151811061433a5761433a6150eb565b0160200151600160f89190911c81901b92505b845181101561441157848181518110614368576143686150eb565b0160200151600160f89190911c1b91508282116143fd5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a4016108d1565b9181179161440a81615204565b905061434d565b50909392505050565b600080614425614542565b61442d614560565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa92508280156136f15750826144b75760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c75726500000000000060448201526064016108d1565b505195945050505050565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b604051806040016040528061451161457e565b815260200161451e61457e565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60ff811681146108e357600080fd5b6000602082840312156145bd57600080fd5b8135613d988161459c565b6001600160a01b03811681146108e357600080fd5b6000602082840312156145ef57600080fd5b8135613d98816145c8565b60006020828403121561460c57600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561464b5761464b614613565b60405290565b60405161010081016001600160401b038111828210171561464b5761464b614613565b604051601f8201601f191681016001600160401b038111828210171561469c5761469c614613565b604052919050565b6000604082840312156146b657600080fd5b6146be614629565b9050813581526020820135602082015292915050565b600082601f8301126146e557600080fd5b6146ed614629565b8060408401858111156146ff57600080fd5b845b81811015614719578035845260209384019301614701565b509095945050505050565b60006080828403121561473657600080fd5b61473e614629565b905061474a83836146d4565b815261475983604084016146d4565b602082015292915050565b600080600080610120858703121561477b57600080fd5b8435935061478c86602087016146a4565b925061479b8660608701614724565b91506147aa8660e087016146a4565b905092959194509250565b6000602082840312156147c757600080fd5b813561ffff81168114613d9857600080fd5b8035613f44816145c8565b6020808252825182820181905260009190848201906040850190845b818110156148255783516001600160a01b031683529284019291840191600101614800565b50909695505050505050565b80151581146108e357600080fd5b60006020828403121561485157600080fd5b8135613d9881614831565b63ffffffff811681146108e357600080fd5b8035613f448161485c565b60006001600160401b0382111561489257614892614613565b5060051b60200190565b600082601f8301126148ad57600080fd5b813560206148c26148bd83614879565b614674565b82815260059290921b840181019181810190868411156148e157600080fd5b8286015b848110156149055780356148f88161485c565b83529183019183016148e5565b509695505050505050565b600082601f83011261492157600080fd5b813560206149316148bd83614879565b82815260069290921b8401810191818101908684111561495057600080fd5b8286015b848110156149055761496688826146a4565b835291830191604001614954565b600082601f83011261498557600080fd5b813560206149956148bd83614879565b82815260059290921b840181019181810190868411156149b457600080fd5b8286015b848110156149055780356001600160401b038111156149d75760008081fd5b6149e58986838b010161489c565b8452509183019183016149b8565b60006101808284031215614a0657600080fd5b614a0e614651565b905081356001600160401b0380821115614a2757600080fd5b614a338583860161489c565b83526020840135915080821115614a4957600080fd5b614a5585838601614910565b60208401526040840135915080821115614a6e57600080fd5b614a7a85838601614910565b6040840152614a8c8560608601614724565b6060840152614a9e8560e086016146a4565b6080840152610120840135915080821115614ab857600080fd5b614ac48583860161489c565b60a0840152610140840135915080821115614ade57600080fd5b614aea8583860161489c565b60c0840152610160840135915080821115614b0457600080fd5b50614b1184828501614974565b60e08301525092915050565b600080600080600060808688031215614b3557600080fd5b8535945060208601356001600160401b0380821115614b5357600080fd5b818801915088601f830112614b6757600080fd5b813581811115614b7657600080fd5b896020828501011115614b8857600080fd5b6020830196509450614b9c6040890161486e565b93506060880135915080821115614bb257600080fd5b50614bbf888289016149f3565b9150509295509295909350565b600081518084526020808501945080840160005b83811015614c055781516001600160601b031687529582019590820190600101614be0565b509495945050505050565b6040815260008351604080840152614c2b6080840182614bcc565b90506020850151603f19848303016060850152614c488282614bcc565b925050508260208301529392505050565b600080600080600060a08688031215614c7157600080fd5b8535614c7c816145c8565b945060208681013594506040870135614c94816145c8565b935060608701356001600160401b03811115614caf57600080fd5b8701601f81018913614cc057600080fd5b8035614cce6148bd82614879565b81815260059190911b8201830190838101908b831115614ced57600080fd5b928401925b82841015614d14578335614d05816145c8565b82529284019290840190614cf2565b8096505050505050614d28608087016147d9565b90509295509295909350565b60008060408385031215614d4757600080fd5b82356001600160401b0380821115614d5e57600080fd5b9084019060808287031215614d7257600080fd5b90925060208401359080821115614d8857600080fd5b50614d95858286016149f3565b9150509250929050565b6000815180845260005b81811015614dc557602081850181015186830182015201614da9565b81811115614dd7576000602083870101525b50601f01601f19169290920160200192915050565b602081526000613d986020830184614d9f565b60006001600160401b03831115614e1857614e18614613565b614e2b601f8401601f1916602001614674565b9050828152838383011115614e3f57600080fd5b828260208301376000602084830101529392505050565b600082601f830112614e6757600080fd5b613d9883833560208501614dff565b60008060408385031215614e8957600080fd5b8235614e94816145c8565b915060208301356001600160401b0380821115614eb057600080fd5b9084019060608287031215614ec457600080fd5b604051606081018181108382111715614edf57614edf614613565b604052823582811115614ef157600080fd5b614efd88828601614e56565b82525060208301356020820152604083013560408201528093505050509250929050565b600060208284031215614f3357600080fd5b81356001600160401b03811115614f4957600080fd5b8201601f81018413614f5a57600080fd5b614f6984823560208401614dff565b949350505050565b600060208284031215614f8357600080fd5b8135613d988161485c565b60008060208385031215614fa157600080fd5b82356001600160401b0380821115614fb857600080fd5b818501915085601f830112614fcc57600080fd5b813581811115614fdb57600080fd5b8660208260051b8501011115614ff057600080fd5b60209290920196919550909350505050565b60006020828403121561501457600080fd5b8151613d9881614831565b60006020828403121561503157600080fd5b8151613d98816145c8565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b6000602082840312156150e057600080fd5b8151613d988161459c565b634e487b7160e01b600052603260045260246000fd5b60008261511e57634e487b7160e01b600052601260045260246000fd5b500690565b60006060828403121561513557600080fd5b604051606081018181106001600160401b038211171561515757615157614613565b60405282516151658161485c565b815260208301516151758161485c565b602082015260408301516151888161459c565b60408201529392505050565b6000602082840312156151a657600080fd5b5051919050565b6000602082840312156151bf57600080fd5b81516001600160c01b0381168114613d9857600080fd5b634e487b7160e01b600052601160045260246000fd5b600082198211156151ff576151ff6151d6565b500190565b6000600019821415615218576152186151d6565b5060010190565b6001600160601b03811681146108e357600080fd5b60006040828403121561524657600080fd5b61524e614629565b8251615259816145c8565b815260208301516152698161521f565b60208201529392505050565b600082821015615287576152876151d6565b500390565b60006020828403121561529e57600080fd5b815167ffffffffffffffff1981168114613d9857600080fd5b6000602082840312156152c957600080fd5b8151613d988161521f565b60006001600160601b03838116908316818110156152f4576152f46151d6565b039392505050565b63ffffffff60e01b8360e01b1681526000600482018351602080860160005b838110156153375781518552938201939082019060010161531b565b5092979650505050505050565b600063ffffffff808316818516808303821115615363576153636151d6565b01949350505050565b6000808335601e1984360301811261538357600080fd5b8301803591506001600160401b0382111561539d57600080fd5b6020019150368190038213156153b257600080fd5b9250929050565b6000608082360312156153cb57600080fd5b604051608081016001600160401b0382821081831117156153ee576153ee614613565b8160405284358352602085013591508082111561540a57600080fd5b61541636838701614e56565b6020840152604085013591508082111561542f57600080fd5b5061543c36828601614e56565b60408301525060608301356154508161485c565b606082015292915050565b60006001600160601b0380831681851681830481118215151615615481576154816151d6565b02949350505050565b60008160001904831182151516156154a4576154a46151d6565b500290565b60208082526052908201527f536572766963654d616e61676572426173652e6f6e6c7952656769737472794360408201527f6f6f7264696e61746f723a2063616c6c6572206973206e6f742074686520726560608201527133b4b9ba393c9031b7b7b93234b730ba37b960711b608082015260a00190565b60018060a01b038316815260406020820152600082516060604084015261554b60a0840182614d9f565b90506020840151606084015260408401516080840152809150509392505050565b60006040828403121561557e57600080fd5b615586614629565b82516155918161459c565b815260208301516152698161459c565b60008235609e198336030181126155b757600080fd5b9190910192915050565b8183526000602080850194508260005b85811015614c055781356155e4816145c8565b6001600160a01b03168752818301356155fc8161521f565b6001600160601b03168784015260409687019691909101906001016155d1565b60208082528181018390526000906040808401600586901b8501820187855b8881101561571b57878303603f190184528135368b9003609e1901811261566157600080fd5b8a0160a0813536839003601e1901811261567a57600080fd5b820180356001600160401b0381111561569257600080fd5b8060061b36038413156156a457600080fd5b8287526156b6838801828c85016155c1565b925050506156c58883016147d9565b6001600160a01b031688860152818701358786015260606156e781840161486e565b63ffffffff169086015260806156fe83820161486e565b63ffffffff1695019490945250928501929085019060010161563b565b509098975050505050505050565b600061ffff80831681811415615741576157416151d6565b6001019392505050565b6000808335601e1984360301811261576257600080fd5b83016020810192503590506001600160401b0381111561578157600080fd5b8036038313156153b257600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b602081528135602082015260006157d3602084018461574b565b608060408501526157e860a085018284615790565b9150506157f8604085018561574b565b848303601f1901606086015261580f838284615790565b9250505060608401356158218161485c565b63ffffffff16608093909301929092525091905056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47456967656e4441536572766963654d616e616765722e636f6e6669726d426174424c535369676e6174757265436865636b65722e636865636b5369676e617475a264697066735822122022b793b54d30654647995b796aab4ad6fb84e318c90d8daa14f38eb714197db664736f6c634300080c0033",
}

// ContractEigenDAServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.ABI instead.
var ContractEigenDAServiceManagerABI = ContractEigenDAServiceManagerMetaData.ABI

// ContractEigenDAServiceManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.Bin instead.
var ContractEigenDAServiceManagerBin = ContractEigenDAServiceManagerMetaData.Bin

// DeployContractEigenDAServiceManager deploys a new Ethereum contract, binding an instance of ContractEigenDAServiceManager to it.
func DeployContractEigenDAServiceManager(auth *bind.TransactOpts, backend bind.ContractBackend, __avsDirectory common.Address, __rewardsCoordinator common.Address, __registryCoordinator common.Address, __stakeRegistry common.Address, __eigenDAThresholdRegistry common.Address, __eigenDARelayRegistry common.Address) (common.Address, *types.Transaction, *ContractEigenDAServiceManager, error) {
	parsed, err := ContractEigenDAServiceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDAServiceManagerBin), backend, __avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry, __eigenDAThresholdRegistry, __eigenDARelayRegistry)
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

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) GetDefaultSecurityThresholdsV2(opts *bind.CallOpts) (SecurityThresholds, error) {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "getDefaultSecurityThresholdsV2")

	if err != nil {
		return *new(SecurityThresholds), err
	}

	out0 := *abi.ConvertType(out[0], new(SecurityThresholds)).(*SecurityThresholds)

	return out0, err

}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDAServiceManager.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDAServiceManager.CallOpts)
}

// GetDefaultSecurityThresholdsV2 is a free data retrieval call binding the contract method 0xef635529.
//
// Solidity: function getDefaultSecurityThresholdsV2() view returns((uint8,uint8))
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetDefaultSecurityThresholdsV2() (SecurityThresholds, error) {
	return _ContractEigenDAServiceManager.Contract.GetDefaultSecurityThresholdsV2(&_ContractEigenDAServiceManager.CallOpts)
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
