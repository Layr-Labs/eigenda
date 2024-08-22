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

// IEigenDABlobVerifierBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type IEigenDABlobVerifierBlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  IEigenDAServiceManagerBatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
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
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// IEigenDASignatureVerifierNonSignerStakesAndSignature is an auto generated low-level Go binding around an user-defined struct.
type IEigenDASignatureVerifierNonSignerStakesAndSignature struct {
	NonSignerQuorumBitmapIndices []uint32
	NonSignerPubkeys             []BN254G1Point
	QuorumApks                   []BN254G1Point
	ApkG2                        BN254G2Point
	Sigma                        BN254G1Point
	QuorumApkIndices             []uint32
	TotalStakeIndices            []uint32
	NonSignerStakeIndices        [][]uint32
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

// ContractEigenDAServiceManagerMetaData contains all meta data concerning the ContractEigenDAServiceManager contract.
var ContractEigenDAServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"__avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"__rewardsCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRewardsCoordinator\"},{\"name\":\"__registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"__stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BLOCK_STALE_MEASURE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STORE_DURATION_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"THRESHOLD_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchIdToBatchMetadataHash\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blsApkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBLSApkRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkSignatures\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.QuorumStakeTotals\",\"components\":[{\"name\":\"signedStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"totalStakeForQuorum\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmBatch\",\"inputs\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIBLSSignatureChecker.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAVSRewardsSubmission\",\"inputs\":[{\"name\":\"rewardsSubmissions\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.RewardsSubmission[]\",\"components\":[{\"name\":\"strategiesAndMultipliers\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardsCoordinator.StrategyAndMultiplier[]\",\"components\":[{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"multiplier\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getIsQuorumRequired\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumAdversaryThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuorumConfirmationThresholdPercentage\",\"inputs\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"_initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_batchConfirmers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_rewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isBatchConfirmer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestServeUntilBlock\",\"inputs\":[{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumAdversaryThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumConfirmationThresholdPercentages\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumNumbersRequired\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryCoordinator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardsInitiator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBatchConfirmer\",\"inputs\":[{\"name\":\"_batchConfirmer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardsInitiator\",\"inputs\":[{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStaleStakesForbidden\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staleStakesForbidden\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"taskNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trySignatureAndApkVerification\",\"inputs\":[{\"name\":\"msgHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"apk\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"pairingSuccessful\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"siganatureIsValid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyBlob\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structIEigenDABlobVerifier.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyBlob\",\"inputs\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structIEigenDABlobVerifier.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyPreconfirmation\",\"inputs\":[{\"name\":\"miniBatchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIEigenDASignatureVerifier.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyPreconfirmation\",\"inputs\":[{\"name\":\"miniBatchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"nonSignerStakesAndSignature\",\"type\":\"tuple\",\"internalType\":\"structIEigenDASignatureVerifier.NonSignerStakesAndSignature\",\"components\":[{\"name\":\"nonSignerQuorumBitmapIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerPubkeys\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApks\",\"type\":\"tuple[]\",\"internalType\":\"structBN254.G1Point[]\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"apkG2\",\"type\":\"tuple\",\"internalType\":\"structBN254.G2Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"Y\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"}]},{\"name\":\"sigma\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"quorumApkIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"totalStakeIndices\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"nonSignerStakeIndices\",\"type\":\"uint32[][]\",\"internalType\":\"uint32[][]\"}]},{\"name\":\"additionalQuorumNumbersRequired\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BatchConfirmed\",\"inputs\":[{\"name\":\"batchHeaderHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"batchId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BatchConfirmerStatusChanged\",\"inputs\":[{\"name\":\"batchConfirmer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsInitiatorUpdated\",\"inputs\":[{\"name\":\"prevRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newRewardsInitiator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StaleStakesForbiddenUpdate\",\"inputs\":[{\"name\":\"value\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6101806040523480156200001257600080fd5b5060405162007408380380620074088339810160408190526200003591620002e5565b6001600160a01b0380851660805280841660a05280831660c052811660e0528184848284620000636200020a565b50505050806001600160a01b0316610100816001600160a01b031681525050806001600160a01b031663683048356040518163ffffffff1660e01b8152600401602060405180830381865afa158015620000c1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000e791906200034d565b6001600160a01b0316610120816001600160a01b031681525050806001600160a01b0316635df459466040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000140573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200016691906200034d565b6001600160a01b0316610140816001600160a01b031681525050610120516001600160a01b031663df5cf7236040518163ffffffff1660e01b8152600401602060405180830381865afa158015620001c2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001e891906200034d565b6001600160a01b03166101605250620002006200020a565b5050505062000374565b603254610100900460ff1615620002775760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60325460ff9081161015620002ca576032805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b0381168114620002e257600080fd5b50565b60008060008060808587031215620002fc57600080fd5b84516200030981620002cc565b60208601519094506200031c81620002cc565b60408601519093506200032f81620002cc565b60608601519092506200034281620002cc565b939692955090935050565b6000602082840312156200036057600080fd5b81516200036d81620002cc565b9392505050565b60805160a05160c05160e05161010051610120516101405161016051616f86620004826000396000818161063c01526119220152600081816103f10152611b0401526000818161044301528181611cda0152611e9c0152600081816104900152818161103e015281816115ed0152818161178501526119bf015260008181610d5601528181610eb101528181610f4801528181612bc801528181612d4b0152612dea015260008181610b7d01528181610c0c01528181610c8c015281816128b60152818161297a01528181612b060152612ca60152600081816132fe015281816133ba01526134a60152600081816104670152818161290a015281816129d60152612a550152616f866000f3fe608060405234801561001057600080fd5b50600436106102a05760003560e01c80637dc21e2611610167578063df5cf723116100ce578063ef02445811610087578063ef024458146106cc578063f1220983146106d4578063f2fde38b146106e7578063fabc1cbc146106fa578063fc299dee1461070d578063fce36c7d1461072057600080fd5b8063df5cf72314610637578063e15234ff1461065e578063e481af9d1461067e578063eaefd27d14610686578063eccbbfc914610699578063ee6c3bcf146106b957600080fd5b8063a5b7890a11610120578063a5b7890a146105ad578063a98fb355146105d0578063b98d0908146105e3578063bafa9107146105f0578063c18e773814610611578063da388de41461062457600080fd5b80637dc21e26146105225780638687feae14610535578063886f1195146105635780638da5cb5b146105765780639926ee7d14610587578063a364f4da1461059a57600080fd5b80635df459461161020b5780636efb4636116101c45780636efb4636146104b2578063715018a6146104d357806372d18e8d146104db578063775bbcb5146104e95780637794965a146104fc5780637b204a141461050f57600080fd5b80635df45946146103ec5780635e0334761461042b5780635e8b3f2d14610435578063683048351461043e5780636b3aa72e146104655780636d14a9871461048b57600080fd5b80633bc28c8c1161025d5780633bc28c8c14610364578063416c7e5e146103775780634972134a1461038a578063595c6a67146103af5780635ac86ab7146103b75780635c975abb146103da57600080fd5b8063048886d2146102a557806310d67a2f146102cd578063136439dd146102e25780631429c7c2146102f5578063171f1d5b1461031a57806333cfb7b714610344575b600080fd5b6102b86102b3366004615400565b610733565b60405190151581526020015b60405180910390f35b6102e06102db366004615432565b61076c565b005b6102e06102f036600461544f565b610828565b610308610303366004615400565b610967565b60405160ff90911681526020016102c4565b61032d6103283660046155fd565b6109ce565b6040805192151583529015156020830152016102c4565b610357610352366004615432565b610b58565b6040516102c49190615659565b6102e0610372366004615432565b61102b565b6102e06103853660046156b4565b61103c565b60005461039a9063ffffffff1681565b60405163ffffffff90911681526020016102c4565b6102e0611173565b6102b86103c5366004615400565b60fc54600160ff9092169190911b9081161490565b60fc545b6040519081526020016102c4565b6104137f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016102c4565b61039a620189c081565b61039a61012c81565b6104137f000000000000000000000000000000000000000000000000000000000000000081565b7f0000000000000000000000000000000000000000000000000000000000000000610413565b6104137f000000000000000000000000000000000000000000000000000000000000000081565b6104c56104c0366004615987565b61123a565b6040516102c4929190615a7a565b6102e0612151565b60005463ffffffff1661039a565b6102e06104f7366004615ac3565b612165565b6102e061050a366004615bb6565b6122ce565b6102e061051d366004615ca2565b61283a565b6102e0610530366004615d29565b612885565b61055660405180604001604052806002815260200161212160f01b81525081565b6040516102c49190615dfe565b60fb54610413906001600160a01b031681565b6065546001600160a01b0316610413565b6102e0610595366004615e11565b6128ab565b6102e06105a8366004615432565b61296f565b6102b86105bb366004615432565b60026020526000908152604090205460ff1681565b6102e06105de366004615ea6565b612a36565b60c9546102b89060ff1681565b61055660405180604001604052806002815260200161373760f01b81525081565b6102e061061f366004615ef6565b612a8a565b6102e0610632366004615f4f565b612ab3565b6104137f000000000000000000000000000000000000000000000000000000000000000081565b610556604051806040016040528060028152602001600160f01b81525081565b610357612b00565b61039a610694366004615ffb565b612ec9565b6103de6106a7366004615ffb565b60016020526000908152604090205481565b6103086106c7366004615400565b612eeb565b6103de606481565b6102e06106e2366004615432565b612f42565b6102e06106f5366004615432565b612f53565b6102e061070836600461544f565b612fc9565b609754610413906001600160a01b031681565b6102e061072e366004616016565b613125565b600080600160ff84161b905080610762604051806040016040528060028152602001600160f01b8152506134dd565b9091161492915050565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107e3919061608a565b6001600160a01b0316336001600160a01b03161461081c5760405162461bcd60e51b8152600401610813906160a7565b60405180910390fd5b6108258161366a565b50565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015610870573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061089491906160f1565b6108b05760405162461bcd60e51b81526004016108139061610e565b60fc54818116146109295760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c69747900000000000000006064820152608401610813565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d906020015b60405180910390a250565b60008160ff1660405180604001604052806002815260200161373760f01b8152505111156109c95760405180604001604052806002815260200161373760f01b8152508260ff16815181106109be576109be616156565b016020015160f81c90505b919050565b60008060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000187876000015188602001518860000151600060028110610a1657610a16616156565b60200201518951600160200201518a60200151600060028110610a3b57610a3b616156565b60200201518b60200151600160028110610a5757610a57616156565b602090810291909101518c518d830151604051610ab49a99989796959401988952602089019790975260408801959095526060870193909352608086019190915260a085015260c084015260e08301526101008201526101200190565b6040516020818303038152906040528051906020012060001c610ad7919061616c565b9050610b4a610af0610ae98884613761565b86906137f8565b610af861388c565b610b40610b3185610b2b604080518082018252600080825260209182015281518083019092526001825260029082015290565b90613761565b610b3a8c61394c565b906137f8565b886201d4c06139dc565b909890975095505050505050565b6040516309aa152760e11b81526001600160a01b0382811660048301526060916000917f000000000000000000000000000000000000000000000000000000000000000016906313542a4e90602401602060405180830381865afa158015610bc4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be8919061618e565b60405163871ef04960e01b8152600481018290529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063871ef04990602401602060405180830381865afa158015610c53573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c7791906161a7565b90506001600160c01b0381161580610d1157507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ce8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d0c91906161d0565b60ff16155b15610d315760408051600080825260208201909252905b50949350505050565b6000610d45826001600160c01b0316613c00565b90506000805b8251811015610e1b577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316633ca5a5f5848381518110610d9557610d95616156565b01602001516040516001600160e01b031960e084901b16815260f89190911c6004820152602401602060405180830381865afa158015610dd9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dfd919061618e565b610e079083616203565b915080610e138161621b565b915050610d4b565b506000816001600160401b03811115610e3657610e36615468565b604051908082528060200260200182016040528015610e5f578160200160208202803683370190505b5090506000805b845181101561101e576000858281518110610e8357610e83616156565b0160200151604051633ca5a5f560e01b815260f89190911c6004820181905291506000906001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690633ca5a5f590602401602060405180830381865afa158015610ef8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f1c919061618e565b905060005b81811015611008576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa158015610f96573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fba919061624b565b60000151868681518110610fd057610fd0616156565b6001600160a01b039092166020928302919091019091015284610ff28161621b565b95505080806110009061621b565b915050610f21565b50505080806110169061621b565b915050610e66565b5090979650505050505050565b611033613cc2565b61082581613d1c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638da5cb5b6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561109a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110be919061608a565b6001600160a01b0316336001600160a01b03161461116a5760405162461bcd60e51b815260206004820152605c60248201527f424c535369676e6174757265436865636b65722e6f6e6c79436f6f7264696e6160448201527f746f724f776e65723a2063616c6c6572206973206e6f7420746865206f776e6560648201527f72206f6620746865207265676973747279436f6f7264696e61746f7200000000608482015260a401610813565b61082581613d85565b60fb5460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa1580156111bb573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111df91906160f1565b6111fb5760405162461bcd60e51b81526004016108139061610e565b60001960fc81905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b60408051808201909152606080825260208201526000846112b15760405162461bcd60e51b81526020600482015260376024820152600080516020616f3183398151915260448201527f7265733a20656d7074792071756f72756d20696e7075740000000000000000006064820152608401610813565b604083015151851480156112c9575060a08301515185145b80156112d9575060c08301515185145b80156112e9575060e08301515185145b6113535760405162461bcd60e51b81526020600482015260416024820152600080516020616f3183398151915260448201527f7265733a20696e7075742071756f72756d206c656e677468206d69736d6174636064820152600d60fb1b608482015260a401610813565b825151602084015151146113cb5760405162461bcd60e51b815260206004820152604460248201819052600080516020616f31833981519152908201527f7265733a20696e707574206e6f6e7369676e6572206c656e677468206d69736d6064820152630c2e8c6d60e31b608482015260a401610813565b4363ffffffff168463ffffffff161061143a5760405162461bcd60e51b815260206004820152603c6024820152600080516020616f3183398151915260448201527f7265733a20696e76616c6964207265666572656e636520626c6f636b000000006064820152608401610813565b6040805180820182526000808252602080830191909152825180840190935260608084529083015290866001600160401b0381111561147b5761147b615468565b6040519080825280602002602001820160405280156114a4578160200160208202803683370190505b506020820152866001600160401b038111156114c2576114c2615468565b6040519080825280602002602001820160405280156114eb578160200160208202803683370190505b50815260408051808201909152606080825260208201528560200151516001600160401b0381111561151f5761151f615468565b604051908082528060200260200182016040528015611548578160200160208202803683370190505b5081526020860151516001600160401b0381111561156857611568615468565b604051908082528060200260200182016040528015611591578160200160208202803683370190505b50816020018190525060006116638a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505060408051639aa1653d60e01b815290516001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169350639aa1653d925060048083019260209291908290030181865afa15801561163a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061165e91906161d0565b613dcd565b905060005b8760200151518110156118fe576116ad8860200151828151811061168e5761168e616156565b6020026020010151805160009081526020918201519091526040902090565b836020015182815181106116c3576116c3616156565b602090810291909101015280156117835760208301516116e460018361628c565b815181106116f4576116f4616156565b602002602001015160001c8360200151828151811061171557611715616156565b602002602001015160001c11611783576040805162461bcd60e51b8152602060048201526024810191909152600080516020616f3183398151915260448201527f7265733a206e6f6e5369676e65725075626b657973206e6f7420736f727465646064820152608401610813565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166304ec6351846020015183815181106117c8576117c8616156565b60200260200101518b8b6000015185815181106117e7576117e7616156565b60200260200101516040518463ffffffff1660e01b81526004016118249392919092835263ffffffff918216602084015216604082015260600190565b602060405180830381865afa158015611841573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061186591906161a7565b6001600160c01b03168360000151828151811061188457611884616156565b6020026020010181815250506118ea610ae96118be84866000015185815181106118b0576118b0616156565b602002602001015116613e5e565b8a6020015184815181106118d4576118d4616156565b6020026020010151613e8990919063ffffffff16565b9450806118f68161621b565b915050611668565b505061190983613f6d565b60c95490935060ff166000816119205760006119a2565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c448feb86040518163ffffffff1660e01b8152600401602060405180830381865afa15801561197e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119a2919061618e565b905060005b8a811015612020578215611b02578963ffffffff16827f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663249a0c428f8f868181106119fe576119fe616156565b60405160e085901b6001600160e01b031916815292013560f81c600483015250602401602060405180830381865afa158015611a3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a62919061618e565b611a6c9190616203565b11611b025760405162461bcd60e51b81526020600482015260666024820152600080516020616f3183398151915260448201527f7265733a205374616b6552656769737472792075706461746573206d7573742060648201527f62652077697468696e207769746864726177616c44656c6179426c6f636b732060848201526577696e646f7760d01b60a482015260c401610813565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166368bccaac8d8d84818110611b4357611b43616156565b9050013560f81c60f81b60f81c8c8c60a001518581518110611b6757611b67616156565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015611bc3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611be791906162a3565b6001600160401b031916611c0a8a60400151838151811061168e5761168e616156565b67ffffffffffffffff191614611ca65760405162461bcd60e51b81526020600482015260616024820152600080516020616f3183398151915260448201527f7265733a2071756f72756d41706b206861736820696e2073746f72616765206460648201527f6f6573206e6f74206d617463682070726f76696465642071756f72756d2061706084820152606b60f81b60a482015260c401610813565b611cd689604001518281518110611cbf57611cbf616156565b6020026020010151876137f890919063ffffffff16565b95507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c8294c568d8d84818110611d1957611d19616156565b9050013560f81c60f81b60f81c8c8c60c001518581518110611d3d57611d3d616156565b60209081029190910101516040516001600160e01b031960e086901b16815260ff909316600484015263ffffffff9182166024840152166044820152606401602060405180830381865afa158015611d99573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611dbd91906162ce565b85602001518281518110611dd357611dd3616156565b6001600160601b03909216602092830291909101820152850151805182908110611dff57611dff616156565b602002602001015185600001518281518110611e1d57611e1d616156565b60200260200101906001600160601b031690816001600160601b0316815250506000805b8a602001515181101561200b57611e9586600001518281518110611e6757611e67616156565b60200260200101518f8f86818110611e8157611e81616156565b600192013560f81c9290921c811614919050565b15611ff9577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663f2be94ae8f8f86818110611edb57611edb616156565b9050013560f81c60f81b60f81c8e89602001518581518110611eff57611eff616156565b60200260200101518f60e001518881518110611f1d57611f1d616156565b60200260200101518781518110611f3657611f36616156565b60209081029190910101516040516001600160e01b031960e087901b16815260ff909416600485015263ffffffff92831660248501526044840191909152166064820152608401602060405180830381865afa158015611f9a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611fbe91906162ce565b8751805185908110611fd257611fd2616156565b60200260200101818151611fe691906162eb565b6001600160601b03169052506001909101905b806120038161621b565b915050611e41565b505080806120189061621b565b9150506119a7565b50505060008061203a8c868a606001518b608001516109ce565b91509150816120ab5760405162461bcd60e51b81526020600482015260436024820152600080516020616f3183398151915260448201527f7265733a2070616972696e6720707265636f6d70696c652063616c6c206661696064820152621b195960ea1b608482015260a401610813565b8061210c5760405162461bcd60e51b81526020600482015260396024820152600080516020616f3183398151915260448201527f7265733a207369676e617475726520697320696e76616c6964000000000000006064820152608401610813565b50506000878260200151604051602001612127929190616313565b60408051808303601f190181529190528051602090910120929b929a509198505050505050505050565b612159613cc2565b6121636000614003565b565b603254610100900460ff16158080156121855750603254600160ff909116105b8061219f5750303b15801561219f575060325460ff166001145b6122025760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610813565b6032805460ff191660011790558015612225576032805461ff0019166101001790555b61222f8686614055565b61223884614003565b61224182613d1c565b60005b835181101561227f5761226f84828151811061226257612262616156565b602002602001015161413b565b6122788161621b565b9050612244565b5080156122c6576032805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050565b60fc54600090600190811614156123275760405162461bcd60e51b815260206004820152601960248201527f5061757361626c653a20696e64657820697320706175736564000000000000006044820152606401610813565b3360009081526002602052604090205460ff1661239b5760405162461bcd60e51b815260206004820152602c60248201527f6f6e6c794261746368436f6e6669726d65723a206e6f742066726f6d2062617460448201526b31b41031b7b73334b936b2b960a11b6064820152608401610813565b3233146124185760405162461bcd60e51b81526020600482015260516024820152600080516020616ef183398151915260448201527f63683a2068656164657220616e64206e6f6e7369676e65722064617461206d75606482015270737420626520696e2063616c6c6461746160781b608482015260a401610813565b436124296080850160608601615ffb565b63ffffffff16106124a85760405162461bcd60e51b815260206004820152604f6024820152600080516020616ef183398151915260448201527f63683a20737065636966696564207265666572656e6365426c6f636b4e756d6260648201526e657220697320696e2066757475726560881b608482015260a401610813565b63ffffffff431661012c6124c26080860160608701615ffb565b6124cc919061635b565b63ffffffff1610156125525760405162461bcd60e51b81526020600482015260556024820152600080516020616ef183398151915260448201527f63683a20737065636966696564207265666572656e6365426c6f636b4e756d62606482015274195c881a5cc81d1bdbc819985c881a5b881c185cdd605a1b608482015260a401610813565b61255f6040840184616383565b905061256e6020850185616383565b9050146126065760405162461bcd60e51b81526020600482015260666024820152600080516020616ef183398151915260448201527f63683a2071756f72756d4e756d6265727320616e64207369676e65645374616b60648201527f65466f7251756f72756d73206d757374206265206f66207468652073616d65206084820152650d8cadccee8d60d31b60a482015260c401610813565b600061261961261485616459565b61419e565b90506000806126458361262f6020890189616383565b61263f60808b0160608c01615ffb565b8961123a565b9150915060005b6126596040880188616383565b905081101561279b5761266f6040880188616383565b8281811061267f5761267f616156565b9050013560f81c60f81b60f81c60ff16836020015182815181106126a5576126a5616156565b60200260200101516126b79190616465565b6001600160601b03166064846000015183815181106126d8576126d8616156565b60200260200101516001600160601b03166126f39190616494565b10156127895760405162461bcd60e51b815260206004820152606460248201819052600080516020616ef183398151915260448301527f63683a207369676e61746f7269657320646f206e6f74206f776e206174206c65908201527f617374207468726573686f6c642070657263656e74616765206f6620612071756084820152636f72756d60e01b60a482015260c401610813565b806127938161621b565b91505061264c565b506000805463ffffffff16906127b088614219565b90506127bd81844361422c565b63ffffffff8316600081815260016020908152604091829020939093555190815286917fc75557c4ad49697e231449688be13ef11cb6be8ed0d18819d8dde074a5a16f8a910160405180910390a261281682600161635b565b6000805463ffffffff191663ffffffff929092169190911790555050505050505050565b612880308484604051806040016040528060028152602001600160f01b8152508560405160200161286c9291906164b3565b60405160208183030381529060405261426e565b505050565b61288030848484604051806040016040528060028152602001600160f01b815250614973565b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146128f35760405162461bcd60e51b8152600401610813906164d9565b604051639926ee7d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690639926ee7d906129419085908590600401616551565b600060405180830381600087803b15801561295b57600080fd5b505af11580156122c6573d6000803e3d6000fd5b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146129b75760405162461bcd60e51b8152600401610813906164d9565b6040516351b27a6d60e11b81526001600160a01b0382811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063a364f4da906024015b600060405180830381600087803b158015612a1b57600080fd5b505af1158015612a2f573d6000803e3d6000fd5b5050505050565b612a3e613cc2565b60405163a98fb35560e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063a98fb35590612a01908490600401615dfe565b612aaf308383604051806040016040528060028152602001600160f01b81525061426e565b5050565b612afa30858585604051806040016040528060028152602001600160f01b81525086604051602001612ae69291906164b3565b604051602081830303815290604052614973565b50505050565b606060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612b62573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b8691906161d0565b60ff16905080612ba457505060408051600081526020810190915290565b6000805b82811015612c5957604051633ca5a5f560e01b815260ff821660048201527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015612c17573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c3b919061618e565b612c459083616203565b915080612c518161621b565b915050612ba8565b506000816001600160401b03811115612c7457612c74615468565b604051908082528060200260200182016040528015612c9d578160200160208202803683370190505b5090506000805b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639aa1653d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612d02573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d2691906161d0565b60ff16811015612ebf57604051633ca5a5f560e01b815260ff821660048201526000907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690633ca5a5f590602401602060405180830381865afa158015612d9a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612dbe919061618e565b905060005b81811015612eaa576040516356e4026d60e11b815260ff84166004820152602481018290527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063adc804da906044016040805180830381865afa158015612e38573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612e5c919061624b565b60000151858581518110612e7257612e72616156565b6001600160a01b039092166020928302919091019091015283612e948161621b565b9450508080612ea29061621b565b915050612dc3565b50508080612eb79061621b565b915050612ca4565b5090949350505050565b600061012c612edb620189c08461635b565b612ee5919061635b565b92915050565b60008160ff1660405180604001604052806002815260200161212160f01b8152505111156109c95760405180604001604052806002815260200161212160f01b8152508260ff16815181106109be576109be616156565b612f4a613cc2565b6108258161413b565b612f5b613cc2565b6001600160a01b038116612fc05760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610813565b61082581614003565b60fb60009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561301c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613040919061608a565b6001600160a01b0316336001600160a01b0316146130705760405162461bcd60e51b8152600401610813906160a7565b60fc5419811960fc541916146130ee5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c69747900000000000000006064820152608401610813565b60fc81905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c9060200161095c565b6097546001600160a01b031633146131ba5760405162461bcd60e51b815260206004820152604c60248201527f536572766963654d616e61676572426173652e6f6e6c7952657761726473496e60448201527f69746961746f723a2063616c6c6572206973206e6f742074686520726577617260648201526b32399034b734ba34b0ba37b960a11b608482015260a401610813565b60005b8181101561348e578282828181106131d7576131d7616156565b90506020028101906131e9919061659c565b6131fa906040810190602001615432565b6001600160a01b03166323b872dd333086868681811061321c5761321c616156565b905060200281019061322e919061659c565b604080516001600160e01b031960e087901b1681526001600160a01b039485166004820152939092166024840152013560448201526064016020604051808303816000875af1158015613285573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906132a991906160f1565b5060008383838181106132be576132be616156565b90506020028101906132d0919061659c565b6132e1906040810190602001615432565b604051636eb1769f60e11b81523060048201526001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081166024830152919091169063dd62ed3e90604401602060405180830381865afa15801561334f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613373919061618e565b905083838381811061338757613387616156565b9050602002810190613399919061659c565b6133aa906040810190602001615432565b6001600160a01b031663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000838787878181106133ec576133ec616156565b90506020028101906133fe919061659c565b6040013561340c9190616203565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015613457573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061347b91906160f1565b5050806134879061621b565b90506131bd565b5060405163fce36c7d60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063fce36c7d906129419085908590600401616617565b6000610100825111156135665760405162461bcd60e51b8152602060048201526044602482018190527f4269746d61705574696c732e6f72646572656442797465734172726179546f42908201527f69746d61703a206f7264657265644279746573417272617920697320746f6f206064820152636c6f6e6760e01b608482015260a401610813565b815161357457506000919050565b6000808360008151811061358a5761358a616156565b0160200151600160f89190911c81901b92505b8451811015613661578481815181106135b8576135b8616156565b0160200151600160f89190911c1b915082821161364d5760405162461bcd60e51b815260206004820152604760248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206f72646572656442797465734172726179206973206e6f74206064820152661bdc99195c995960ca1b608482015260a401610813565b9181179161365a8161621b565b905061359d565b50909392505050565b6001600160a01b0381166136f85760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a401610813565b60fb54604080516001600160a01b03928316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a160fb80546001600160a01b0319166001600160a01b0392909216919091179055565b604080518082019091526000808252602082015261377d615317565b835181526020808501519082015260408082018490526000908360608460076107d05a03fa90508080156137b0576137b2565bfe5b50806137f05760405162461bcd60e51b815260206004820152600d60248201526c1958cb5b5d5b0b59985a5b1959609a1b6044820152606401610813565b505092915050565b6040805180820190915260008082526020820152613814615335565b835181526020808501518183015283516040808401919091529084015160608301526000908360808460066107d05a03fa90508080156137b05750806137f05760405162461bcd60e51b815260206004820152600d60248201526c1958cb5859190b59985a5b1959609a1b6044820152606401610813565b613894615353565b50604080516080810182527f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28183019081527f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed6060830152815281518083019092527f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec82527f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d60208381019190915281019190915290565b60408051808201909152600080825260208201526000808061397c600080516020616eb18339815191528661616c565b90505b61398881615089565b9093509150600080516020616eb18339815191528283098314156139c2576040805180820190915290815260208101919091529392505050565b600080516020616eb183398151915260018208905061397f565b604080518082018252868152602080820186905282518084019093528683528201849052600091829190613a0e615378565b60005b6002811015613bd3576000613a27826006616494565b9050848260028110613a3b57613a3b616156565b60200201515183613a4d836000616203565b600c8110613a5d57613a5d616156565b6020020152848260028110613a7457613a74616156565b60200201516020015183826001613a8b9190616203565b600c8110613a9b57613a9b616156565b6020020152838260028110613ab257613ab2616156565b6020020151515183613ac5836002616203565b600c8110613ad557613ad5616156565b6020020152838260028110613aec57613aec616156565b6020020151516001602002015183613b05836003616203565b600c8110613b1557613b15616156565b6020020152838260028110613b2c57613b2c616156565b602002015160200151600060028110613b4757613b47616156565b602002015183613b58836004616203565b600c8110613b6857613b68616156565b6020020152838260028110613b7f57613b7f616156565b602002015160200151600160028110613b9a57613b9a616156565b602002015183613bab836005616203565b600c8110613bbb57613bbb616156565b60200201525080613bcb8161621b565b915050613a11565b50613bdc615397565b60006020826101808560088cfa9151919c9115159b50909950505050505050505050565b6060600080613c0e84613e5e565b61ffff166001600160401b03811115613c2957613c29615468565b6040519080825280601f01601f191660200182016040528015613c53576020820181803683370190505b5090506000805b825182108015613c6b575061010081105b15612ebf576001811b935085841615613cb2578060f81b838381518110613c9457613c94616156565b60200101906001600160f81b031916908160001a9053508160010191505b613cbb8161621b565b9050613c5a565b6065546001600160a01b031633146121635760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610813565b609754604080516001600160a01b03928316815291831660208301527fe11cddf1816a43318ca175bbc52cd0185436e9cbead7c83acc54a73e461717e3910160405180910390a1609780546001600160a01b0319166001600160a01b0392909216919091179055565b60c9805460ff19168215159081179091556040519081527f40e4ed880a29e0f6ddce307457fb75cddf4feef7d3ecb0301bfdf4976a0e2dfc906020015b60405180910390a150565b600080613dd9846134dd565b9050808360ff166001901b11613e575760405162461bcd60e51b815260206004820152603f60248201527f4269746d61705574696c732e6f72646572656442797465734172726179546f4260448201527f69746d61703a206269746d61702065786365656473206d61782076616c7565006064820152608401610813565b9392505050565b6000805b8215612ee557613e7360018461628c565b9092169180613e8181616724565b915050613e62565b60408051808201909152600080825260208201526102008261ffff1610613ee55760405162461bcd60e51b815260206004820152601060248201526f7363616c61722d746f6f2d6c6172676560801b6044820152606401610813565b8161ffff1660011415613ef9575081612ee5565b6040805180820190915260008082526020820181905284906001905b8161ffff168661ffff1610613f6257600161ffff871660ff83161c81161415613f4557613f4284846137f8565b93505b613f4f83846137f8565b92506201fffe600192831b169101613f15565b509195945050505050565b60408051808201909152600080825260208201528151158015613f9257506020820151155b15613fb0575050604080518082019091526000808252602082015290565b604051806040016040528083600001518152602001600080516020616eb18339815191528460200151613fe3919061616c565b613ffb90600080516020616eb183398151915261628c565b905292915050565b606580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b60fb546001600160a01b031615801561407657506001600160a01b03821615155b6140f85760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a401610813565b60fc81905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2612aaf8261366a565b6001600160a01b038116600081815260026020908152604091829020805460ff8082161560ff1990921682179092558351948552161515908301527f5c3265f5fb462ef4930fe47beaa183647c97f19ba545b761f41bc8cd4621d4149101613dc2565b60006141db82604080518082019091526000808252602082015250604080518082019091528151815260609091015163ffffffff16602082015290565b6040805182516020808301919091529092015163ffffffff16908201526060015b604051602081830303815290604052805190602001209050919050565b6000816040516020016141fc91906167b4565b604080516020808201959095528082019390935260e09190911b6001600160e01b03191660608301528051604481840301815260649092019052805191012090565b6001600160a01b03841663eccbbfc961428a6020850185615ffb565b6040516001600160e01b031960e084901b16815263ffffffff919091166004820152602401602060405180830381865afa1580156142cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906142f0919061618e565b61430e614300604085018561682f565b61430990616845565b61510b565b146143955760405162461bcd60e51b81526020600482015260576024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a2062617463684d6574616461746120646f6573206e6f60648201527f74206d617463682073746f726564206d65746164617461000000000000000000608482015260a401610813565b6144496143a56060840184616383565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506143e792505050604085018561682f565b6143f190806168a5565b356144036143fe876168bb565b61514a565b60405160200161441591815260200190565b6040516020818303038152906040528051906020012085602001602081019061443e9190615ffb565b63ffffffff1661515d565b6144b75760405162461bcd60e51b81526020600482015260456024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a20696e636c7573696f6e2070726f6f6620697320696e6064820152641d985b1a5960da1b608482015260a401610813565b6000805b6144c860608601866169d6565b9050811015614943576144de60608601866169d6565b828181106144ee576144ee616156565b6145049260206080909202019081019150615400565b60ff16614514604086018661682f565b61451e90806168a5565b61452c906020810190616383565b6145396080880188616383565b8581811061454957614549616156565b919091013560f81c905081811061456257614562616156565b9050013560f81c60f81b60f81c60ff16146145e25760405162461bcd60e51b81526020600482015260466024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a2071756f72756d4e756d62657220646f6573206e6f74606482015265040dac2e8c6d60d31b608482015260a401610813565b6145ef60608601866169d6565b828181106145ff576145ff616156565b90506080020160400160208101906146179190615400565b60ff1661462760608701876169d6565b8381811061463757614637616156565b905060800201602001602081019061464f9190615400565b60ff16106146ca5760405162461bcd60e51b815260206004820152604e6024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a207468726573686f6c642070657263656e746167657360648201526d08185c99481b9bdd081d985b1a5960921b608482015260a401610813565b60006147026146dc60608801886169d6565b848181106146ec576146ec616156565b6106c79260206080909202019081019150615400565b905060ff8116156147c65760ff811661471e60608801886169d6565b8481811061472e5761472e616156565b90506080020160200160208101906147469190615400565b60ff1610156147c65760405162461bcd60e51b81526020600482015260526024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a206164766572736172795468726573686f6c6450657260648201527118d95b9d1859d9481a5cc81b9bdd081b595d60721b608482015260a401610813565b6147d360608701876169d6565b838181106147e3576147e3616156565b90506080020160400160208101906147fb9190615400565b60ff1661480b604087018761682f565b61481590806168a5565b614823906040810190616383565b6148306080890189616383565b8681811061484057614840616156565b919091013560f81c905081811061485957614859616156565b9050013560f81c60f81b60f81c60ff1610156148e95760405162461bcd60e51b81526020600482015260556024820152600080516020616ed183398151915260448201527f6f7251756f72756d733a20636f6e6669726d6174696f6e5468726573686f6c6460648201527414195c98d95b9d1859d9481a5cc81b9bdd081b595d605a1b608482015260a401610813565b61492d836148fa60608901896169d6565b8581811061490a5761490a616156565b6149209260206080909202019081019150615400565b600160ff919091161b1790565b925050808061493b9061621b565b9150506144bb565b50614957614950836134dd565b8281161490565b612a2f5760405162461bcd60e51b815260040161081390616a1f565b6149806040850185616383565b905061498f6020860186616383565b905014614a3e5760405162461bcd60e51b81526020600482015260786024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a2071756f72756d4e756d6260648201527f65727320616e64207369676e65645374616b65466f7251756f72756d73206d7560848201527f7374206265206f66207468652073616d65206c656e677468000000000000000060a482015260c401610813565b614a4a6143fe846168bb565b604051602001614a5c91815260200190565b60405160208183030381529060405280519060200120846000013514614afe5760405162461bcd60e51b815260206004820152605f6024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a20626c6f6248656164657260648201527f73526f6f7420646f6573206e6f74206d6174636820626c6f6248656164657200608482015260a401610813565b6000614b0c61261486616459565b90506000806001600160a01b038816636efb463684614b2e60208b018b616383565b614b3e60808d0160608e01615ffb565b8a6040518663ffffffff1660e01b8152600401614b5f959493929190616bad565b600060405180830381865afa158015614b7c573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052614ba49190810190616d0f565b9150915060005b614bb86040890189616383565b9050811015614d2357614bf0614bd160208a018a616383565b83818110614be157614be1616156565b919091013560f81c9050610967565b60ff1683602001518281518110614c0957614c09616156565b6020026020010151614c1b9190616465565b6001600160601b0316606484600001518381518110614c3c57614c3c616156565b60200260200101516001600160601b0316614c579190616494565b1015614d115760405162461bcd60e51b81526020600482015260836024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a207369676e61746f72696560648201527f7320646f206e6f74206f776e206174206c6561737420636f6e6669726d61746960848201527f6f6e207468726573686f6c642070657263656e74616765206f6620612071756f60a48201526272756d60e81b60c482015260e401610813565b80614d1b8161621b565b915050614bab565b506000805b614d3560608901896169d6565b905081101561505557614d4b60608901896169d6565b82818110614d5b57614d5b616156565b9050608002016040016020810190614d739190615400565b60ff16614d8360608a018a6169d6565b83818110614d9357614d93616156565b9050608002016020016020810190614dab9190615400565b60ff1610614e355760405162461bcd60e51b81526020600482015260596024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a207468726573686f6c642060648201527f70657263656e746167657320617265206e6f742076616c696400000000000000608482015260a401610813565b6000614e476146dc60608b018b6169d6565b905060ff811615614f165760ff8116614e6360608b018b6169d6565b84818110614e7357614e73616156565b9050608002016020016020810190614e8b9190615400565b60ff161015614f165760405162461bcd60e51b815260206004820152605d6024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a206164766572736172795460648201527f68726573686f6c6450657263656e74616765206973206e6f74206d6574000000608482015260a401610813565b6000614f4e614f2860608c018c6169d6565b85818110614f3857614f38616156565b6103039260206080909202019081019150615400565b905060ff81161561501d5760ff8116614f6a60608c018c6169d6565b85818110614f7a57614f7a616156565b9050608002016040016020810190614f929190615400565b60ff16101561501d5760405162461bcd60e51b81526020600482015260606024820152600080516020616f1183398151915260448201527f6e6669726d6174696f6e466f7251756f72756d733a20636f6e6669726d61746960648201527f6f6e5468726573686f6c6450657263656e74616765206973206e6f74206d6574608482015260a401610813565b61503e8461502e60608d018d6169d6565b8681811061490a5761490a616156565b93505050808061504d9061621b565b915050614d28565b50615062614950866134dd565b61507e5760405162461bcd60e51b815260040161081390616a1f565b505050505050505050565b60008080600080516020616eb18339815191526003600080516020616eb183398151915286600080516020616eb18339815191528889090908905060006150ff827f0c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52600080516020616eb1833981519152615175565b91959194509092505050565b6000612ee582600001516040516020016151259190616dab565b604051602081830303815290604052805190602001208360200151846040015161422c565b6000816040516020016141fc9190616e0b565b60008361516b86858561521d565b1495945050505050565b600080615180615397565b6151886153b5565b602080825281810181905260408201819052606082018890526080820187905260a082018690528260c08360056107d05a03fa92508280156137b05750826152125760405162461bcd60e51b815260206004820152601a60248201527f424e3235342e6578704d6f643a2063616c6c206661696c7572650000000000006044820152606401610813565b505195945050505050565b60006020845161522d919061616c565b156152b45760405162461bcd60e51b815260206004820152604b60248201527f4d65726b6c652e70726f63657373496e636c7573696f6e50726f6f664b65636360448201527f616b3a2070726f6f66206c656e6774682073686f756c642062652061206d756c60648201526a3a34b836329037b310199960a91b608482015260a401610813565b8260205b85518111610d28576152cb60028561616c565b6152ec57816000528086015160205260406000209150600284049350615305565b8086015160005281602052604060002091506002840493505b615310602082616203565b90506152b8565b60405180606001604052806003906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b60405180604001604052806153666153d3565b81526020016153736153d3565b905290565b604051806101800160405280600c906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b60ff8116811461082557600080fd5b60006020828403121561541257600080fd5b8135613e57816153f1565b6001600160a01b038116811461082557600080fd5b60006020828403121561544457600080fd5b8135613e578161541d565b60006020828403121561546157600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b03811182821017156154a0576154a0615468565b60405290565b60405161010081016001600160401b03811182821017156154a0576154a0615468565b604051606081016001600160401b03811182821017156154a0576154a0615468565b604051608081016001600160401b03811182821017156154a0576154a0615468565b604051601f8201601f191681016001600160401b038111828210171561553557615535615468565b604052919050565b60006040828403121561554f57600080fd5b61555761547e565b9050813581526020820135602082015292915050565b600082601f83011261557e57600080fd5b61558661547e565b80604084018581111561559857600080fd5b845b818110156155b257803584526020938401930161559a565b509095945050505050565b6000608082840312156155cf57600080fd5b6155d761547e565b90506155e3838361556d565b81526155f2836040840161556d565b602082015292915050565b600080600080610120858703121561561457600080fd5b84359350615625866020870161553d565b925061563486606087016155bd565b91506156438660e0870161553d565b905092959194509250565b80356109c98161541d565b6020808252825182820181905260009190848201906040850190845b8181101561569a5783516001600160a01b031683529284019291840191600101615675565b50909695505050505050565b801515811461082557600080fd5b6000602082840312156156c657600080fd5b8135613e57816156a6565b803563ffffffff811681146109c957600080fd5b60006001600160401b038211156156fe576156fe615468565b5060051b60200190565b600082601f83011261571957600080fd5b8135602061572e615729836156e5565b61550d565b82815260059290921b8401810191818101908684111561574d57600080fd5b8286015b8481101561576f57615762816156d1565b8352918301918301615751565b509695505050505050565b600082601f83011261578b57600080fd5b8135602061579b615729836156e5565b82815260069290921b840181019181810190868411156157ba57600080fd5b8286015b8481101561576f576157d0888261553d565b8352918301916040016157be565b600082601f8301126157ef57600080fd5b813560206157ff615729836156e5565b82815260059290921b8401810191818101908684111561581e57600080fd5b8286015b8481101561576f5780356001600160401b038111156158415760008081fd5b61584f8986838b0101615708565b845250918301918301615822565b6000610180828403121561587057600080fd5b6158786154a6565b905081356001600160401b038082111561589157600080fd5b61589d85838601615708565b835260208401359150808211156158b357600080fd5b6158bf8583860161577a565b602084015260408401359150808211156158d857600080fd5b6158e48583860161577a565b60408401526158f685606086016155bd565b60608401526159088560e0860161553d565b608084015261012084013591508082111561592257600080fd5b61592e85838601615708565b60a084015261014084013591508082111561594857600080fd5b61595485838601615708565b60c084015261016084013591508082111561596e57600080fd5b5061597b848285016157de565b60e08301525092915050565b60008060008060006080868803121561599f57600080fd5b8535945060208601356001600160401b03808211156159bd57600080fd5b818801915088601f8301126159d157600080fd5b8135818111156159e057600080fd5b8960208285010111156159f257600080fd5b6020830196509450615a06604089016156d1565b93506060880135915080821115615a1c57600080fd5b50615a298882890161585d565b9150509295509295909350565b600081518084526020808501945080840160005b83811015615a6f5781516001600160601b031687529582019590820190600101615a4a565b509495945050505050565b6040815260008351604080840152615a956080840182615a36565b90506020850151603f19848303016060850152615ab28282615a36565b925050508260208301529392505050565b600080600080600060a08688031215615adb57600080fd5b8535615ae68161541d565b945060208681013594506040870135615afe8161541d565b935060608701356001600160401b03811115615b1957600080fd5b8701601f81018913615b2a57600080fd5b8035615b38615729826156e5565b81815260059190911b8201830190838101908b831115615b5757600080fd5b928401925b82841015615b7e578335615b6f8161541d565b82529284019290840190615b5c565b8096505050505050615b926080870161564e565b90509295509295909350565b600060808284031215615bb057600080fd5b50919050565b60008060408385031215615bc957600080fd5b82356001600160401b0380821115615be057600080fd5b615bec86838701615b9e565b93506020850135915080821115615c0257600080fd5b50615c0f8582860161585d565b9150509250929050565b600060a08284031215615bb057600080fd5b60006001600160401b03831115615c4457615c44615468565b615c57601f8401601f191660200161550d565b9050828152838383011115615c6b57600080fd5b828260208301376000602084830101529392505050565b600082601f830112615c9357600080fd5b613e5783833560208501615c2b565b600080600060608486031215615cb757600080fd5b83356001600160401b0380821115615cce57600080fd5b615cda87838801615b9e565b94506020860135915080821115615cf057600080fd5b615cfc87838801615c19565b93506040860135915080821115615d1257600080fd5b50615d1f86828701615c82565b9150509250925092565b600080600060608486031215615d3e57600080fd5b83356001600160401b0380821115615d5557600080fd5b615d6187838801615b9e565b94506020860135915080821115615d7757600080fd5b615d8387838801615b9e565b93506040860135915080821115615d9957600080fd5b50615d1f8682870161585d565b60005b83811015615dc1578181015183820152602001615da9565b83811115612afa5750506000910152565b60008151808452615dea816020860160208601615da6565b601f01601f19169290920160200192915050565b602081526000613e576020830184615dd2565b60008060408385031215615e2457600080fd5b8235615e2f8161541d565b915060208301356001600160401b0380821115615e4b57600080fd5b9084019060608287031215615e5f57600080fd5b615e676154c9565b823582811115615e7657600080fd5b615e8288828601615c82565b82525060208301356020820152604083013560408201528093505050509250929050565b600060208284031215615eb857600080fd5b81356001600160401b03811115615ece57600080fd5b8201601f81018413615edf57600080fd5b615eee84823560208401615c2b565b949350505050565b60008060408385031215615f0957600080fd5b82356001600160401b0380821115615f2057600080fd5b615f2c86838701615b9e565b93506020850135915080821115615f4257600080fd5b50615c0f85828601615c19565b60008060008060808587031215615f6557600080fd5b84356001600160401b0380821115615f7c57600080fd5b615f8888838901615b9e565b95506020870135915080821115615f9e57600080fd5b615faa88838901615b9e565b94506040870135915080821115615fc057600080fd5b615fcc8883890161585d565b93506060870135915080821115615fe257600080fd5b50615fef87828801615c82565b91505092959194509250565b60006020828403121561600d57600080fd5b613e57826156d1565b6000806020838503121561602957600080fd5b82356001600160401b038082111561604057600080fd5b818501915085601f83011261605457600080fd5b81358181111561606357600080fd5b8660208260051b850101111561607857600080fd5b60209290920196919550909350505050565b60006020828403121561609c57600080fd5b8151613e578161541d565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b60006020828403121561610357600080fd5b8151613e57816156a6565b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b60008261618957634e487b7160e01b600052601260045260246000fd5b500690565b6000602082840312156161a057600080fd5b5051919050565b6000602082840312156161b957600080fd5b81516001600160c01b0381168114613e5757600080fd5b6000602082840312156161e257600080fd5b8151613e57816153f1565b634e487b7160e01b600052601160045260246000fd5b60008219821115616216576162166161ed565b500190565b600060001982141561622f5761622f6161ed565b5060010190565b6001600160601b038116811461082557600080fd5b60006040828403121561625d57600080fd5b61626561547e565b82516162708161541d565b8152602083015161628081616236565b60208201529392505050565b60008282101561629e5761629e6161ed565b500390565b6000602082840312156162b557600080fd5b815167ffffffffffffffff1981168114613e5757600080fd5b6000602082840312156162e057600080fd5b8151613e5781616236565b60006001600160601b038381169083168181101561630b5761630b6161ed565b039392505050565b63ffffffff60e01b8360e01b1681526000600482018351602080860160005b8381101561634e57815185529382019390820190600101616332565b5092979650505050505050565b600063ffffffff80831681851680830382111561637a5761637a6161ed565b01949350505050565b6000808335601e1984360301811261639a57600080fd5b8301803591506001600160401b038211156163b457600080fd5b6020019150368190038213156163c957600080fd5b9250929050565b6000608082840312156163e257600080fd5b6163ea6154eb565b90508135815260208201356001600160401b038082111561640a57600080fd5b61641685838601615c82565b6020840152604084013591508082111561642f57600080fd5b5061643c84828501615c82565b60408301525061644e606083016156d1565b606082015292915050565b6000612ee536836163d0565b60006001600160601b038083168185168183048111821515161561648b5761648b6161ed565b02949350505050565b60008160001904831182151516156164ae576164ae6161ed565b500290565b600083516164c5818460208801615da6565b83519083019061637a818360208801615da6565b60208082526052908201527f536572766963654d616e61676572426173652e6f6e6c7952656769737472794360408201527f6f6f7264696e61746f723a2063616c6c6572206973206e6f742074686520726560608201527133b4b9ba393c9031b7b7b93234b730ba37b960711b608082015260a00190565b60018060a01b038316815260406020820152600082516060604084015261657b60a0840182615dd2565b90506020840151606084015260408401516080840152809150509392505050565b60008235609e198336030181126165b257600080fd5b9190910192915050565b8183526000602080850194508260005b85811015615a6f5781356165df8161541d565b6001600160a01b03168752818301356165f781616236565b6001600160601b03168784015260409687019691909101906001016165cc565b60208082528181018390526000906040808401600586901b8501820187855b8881101561671657878303603f190184528135368b9003609e1901811261665c57600080fd5b8a0160a0813536839003601e1901811261667557600080fd5b820180356001600160401b0381111561668d57600080fd5b8060061b360384131561669f57600080fd5b8287526166b1838801828c85016165bc565b925050506166c088830161564e565b6001600160a01b031688860152818701358786015260606166e28184016156d1565b63ffffffff169086015260806166f98382016156d1565b63ffffffff16950194909452509285019290850190600101616636565b509098975050505050505050565b600061ffff8083168181141561673c5761673c6161ed565b6001019392505050565b6000808335601e1984360301811261675d57600080fd5b83016020810192503590506001600160401b0381111561677c57600080fd5b8036038313156163c957600080fd5b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b602081528135602082015260006167ce6020840184616746565b608060408501526167e360a08501828461678b565b9150506167f36040850185616746565b848303601f1901606086015261680a83828461678b565b9250505063ffffffff61681f606086016156d1565b1660808401528091505092915050565b60008235605e198336030181126165b257600080fd5b60006060823603121561685757600080fd5b61685f6154c9565b82356001600160401b0381111561687557600080fd5b616881368286016163d0565b8252506020830135602082015261689a604084016156d1565b604082015292915050565b60008235607e198336030181126165b257600080fd5b600060808083360312156168ce57600080fd5b6168d66154c9565b6168e0368561553d565b815260406168ef8186016156d1565b6020818185015260609150818701356001600160401b0381111561691257600080fd5b870136601f82011261692357600080fd5b8035616931615729826156e5565b81815260079190911b8201830190838101903683111561695057600080fd5b928401925b828410156169c25788843603121561696d5760008081fd5b6169756154eb565b8435616980816153f1565b81528486013561698f816153f1565b81870152848801356169a0816153f1565b818901526169af8588016156d1565b8188015282529288019290840190616955565b958701959095525093979650505050505050565b6000808335601e198436030181126169ed57600080fd5b8301803591506001600160401b03821115616a0757600080fd5b6020019150600781901b36038213156163c957600080fd5b6020808252606590820152600080516020616ed183398151915260408201527f6f7251756f72756d733a2072657175697265642071756f72756d73206172652060608201527f6e6f74206120737562736574206f662074686520636f6e6669726d65642071756080820152646f72756d7360d81b60a082015260c00190565b600081518084526020808501945080840160005b83811015615a6f57815163ffffffff1687529582019590820190600101616ab2565b600081518084526020808501945080840160005b83811015615a6f57616b0587835180518252602090810151910152565b6040969096019590820190600101616ae8565b8060005b6002811015612afa578151845260209384019390910190600101616b1c565b616b46828251616b18565b60208101516128806040840182616b18565b600081518084526020808501808196508360051b8101915082860160005b85811015616ba0578284038952616b8e848351616a9e565b98850198935090840190600101616b76565b5091979650505050505050565b858152608060208201526000616bc760808301868861678b565b63ffffffff8516604084015282810360608401526101808451818352616bef82840182616a9e565b91505060208501518282036020840152616c098282616ad4565b91505060408501518282036040840152616c238282616ad4565b9150506060850151616c386060840182616b3b565b506080850151805160e08401526020015161010083015260a0850151828203610120840152616c678282616a9e565b91505060c0850151828203610140840152616c828282616a9e565b91505060e0850151828203610160840152616c9d8282616b58565b9a9950505050505050505050565b600082601f830112616cbc57600080fd5b81516020616ccc615729836156e5565b82815260059290921b84018101918181019086841115616ceb57600080fd5b8286015b8481101561576f578051616d0281616236565b8352918301918301616cef565b60008060408385031215616d2257600080fd5b82516001600160401b0380821115616d3957600080fd5b9084019060408287031215616d4d57600080fd5b616d5561547e565b825182811115616d6457600080fd5b616d7088828601616cab565b825250602083015182811115616d8557600080fd5b616d9188828601616cab565b602083015250809450505050602083015190509250929050565b60208152815160208201526000602083015160806040840152616dd160a0840182615dd2565b90506040840151601f19848303016060850152616dee8282615dd2565b91505063ffffffff60608501511660808401528091505092915050565b60208082528251805183830152810151604083015260009060a0830181850151606063ffffffff808316828801526040925082880151608080818a015285825180885260c08b0191508884019750600093505b80841015616ea1578751805160ff90811684528a82015181168b850152888201511688840152860151851686830152968801966001939093019290820190616e5e565b509a995050505050505050505056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47456967656e4441426c6f6256657269666965722e5f766572696679426c6f6246456967656e4441536572766963654d616e616765722e636f6e6669726d426174456967656e4441426c6f6256657269666965722e5f766572696679507265636f424c535369676e6174757265436865636b65722e636865636b5369676e617475a26469706673582212206dfa5cecb52c22fc660911bf130e0863fc55547c7774b917e2a2c6e2e4b5557364736f6c634300080c0033",
}

// ContractEigenDAServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.ABI instead.
var ContractEigenDAServiceManagerABI = ContractEigenDAServiceManagerMetaData.ABI

// ContractEigenDAServiceManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractEigenDAServiceManagerMetaData.Bin instead.
var ContractEigenDAServiceManagerBin = ContractEigenDAServiceManagerMetaData.Bin

// DeployContractEigenDAServiceManager deploys a new Ethereum contract, binding an instance of ContractEigenDAServiceManager to it.
func DeployContractEigenDAServiceManager(auth *bind.TransactOpts, backend bind.ContractBackend, __avsDirectory common.Address, __rewardsCoordinator common.Address, __registryCoordinator common.Address, __stakeRegistry common.Address) (common.Address, *types.Transaction, *ContractEigenDAServiceManager, error) {
	parsed, err := ContractEigenDAServiceManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractEigenDAServiceManagerBin), backend, __avsDirectory, __rewardsCoordinator, __registryCoordinator, __stakeRegistry)
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
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
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
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumAdversaryThresholdPercentage is a free data retrieval call binding the contract method 0xee6c3bcf.
//
// Solidity: function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) view returns(uint8 adversaryThresholdPercentage)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) GetQuorumAdversaryThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumAdversaryThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
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
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) GetQuorumConfirmationThresholdPercentage(quorumNumber uint8) (uint8, error) {
	return _ContractEigenDAServiceManager.Contract.GetQuorumConfirmationThresholdPercentage(&_ContractEigenDAServiceManager.CallOpts, quorumNumber)
}

// GetQuorumConfirmationThresholdPercentage is a free data retrieval call binding the contract method 0x1429c7c2.
//
// Solidity: function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) view returns(uint8 confirmationThresholdPercentage)
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

// VerifyBlob is a free data retrieval call binding the contract method 0x7b204a14.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) VerifyBlob(opts *bind.CallOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "verifyBlob", blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlob is a free data retrieval call binding the contract method 0x7b204a14.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) VerifyBlob(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDAServiceManager.Contract.VerifyBlob(&_ContractEigenDAServiceManager.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlob is a free data retrieval call binding the contract method 0x7b204a14.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) VerifyBlob(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDAServiceManager.Contract.VerifyBlob(&_ContractEigenDAServiceManager.CallOpts, blobHeader, blobVerificationProof, additionalQuorumNumbersRequired)
}

// VerifyBlob0 is a free data retrieval call binding the contract method 0xc18e7738.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) VerifyBlob0(opts *bind.CallOpts, blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof) error {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "verifyBlob0", blobHeader, blobVerificationProof)

	if err != nil {
		return err
	}

	return err

}

// VerifyBlob0 is a free data retrieval call binding the contract method 0xc18e7738.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) VerifyBlob0(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof) error {
	return _ContractEigenDAServiceManager.Contract.VerifyBlob0(&_ContractEigenDAServiceManager.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyBlob0 is a free data retrieval call binding the contract method 0xc18e7738.
//
// Solidity: function verifyBlob(((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes) blobVerificationProof) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) VerifyBlob0(blobHeader IEigenDAServiceManagerBlobHeader, blobVerificationProof IEigenDABlobVerifierBlobVerificationProof) error {
	return _ContractEigenDAServiceManager.Contract.VerifyBlob0(&_ContractEigenDAServiceManager.CallOpts, blobHeader, blobVerificationProof)
}

// VerifyPreconfirmation is a free data retrieval call binding the contract method 0x7dc21e26.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) VerifyPreconfirmation(opts *bind.CallOpts, miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature) error {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "verifyPreconfirmation", miniBatchHeader, blobHeader, nonSignerStakesAndSignature)

	if err != nil {
		return err
	}

	return err

}

// VerifyPreconfirmation is a free data retrieval call binding the contract method 0x7dc21e26.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) VerifyPreconfirmation(miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature) error {
	return _ContractEigenDAServiceManager.Contract.VerifyPreconfirmation(&_ContractEigenDAServiceManager.CallOpts, miniBatchHeader, blobHeader, nonSignerStakesAndSignature)
}

// VerifyPreconfirmation is a free data retrieval call binding the contract method 0x7dc21e26.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) VerifyPreconfirmation(miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature) error {
	return _ContractEigenDAServiceManager.Contract.VerifyPreconfirmation(&_ContractEigenDAServiceManager.CallOpts, miniBatchHeader, blobHeader, nonSignerStakesAndSignature)
}

// VerifyPreconfirmation0 is a free data retrieval call binding the contract method 0xda388de4.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCaller) VerifyPreconfirmation0(opts *bind.CallOpts, miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	var out []interface{}
	err := _ContractEigenDAServiceManager.contract.Call(opts, &out, "verifyPreconfirmation0", miniBatchHeader, blobHeader, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)

	if err != nil {
		return err
	}

	return err

}

// VerifyPreconfirmation0 is a free data retrieval call binding the contract method 0xda388de4.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) VerifyPreconfirmation0(miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDAServiceManager.Contract.VerifyPreconfirmation0(&_ContractEigenDAServiceManager.CallOpts, miniBatchHeader, blobHeader, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)
}

// VerifyPreconfirmation0 is a free data retrieval call binding the contract method 0xda388de4.
//
// Solidity: function verifyPreconfirmation((bytes32,bytes,bytes,uint32) miniBatchHeader, ((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]) blobHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature, bytes additionalQuorumNumbersRequired) view returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerCallerSession) VerifyPreconfirmation0(miniBatchHeader IEigenDAServiceManagerBatchHeader, blobHeader IEigenDAServiceManagerBlobHeader, nonSignerStakesAndSignature IEigenDASignatureVerifierNonSignerStakesAndSignature, additionalQuorumNumbersRequired []byte) error {
	return _ContractEigenDAServiceManager.Contract.VerifyPreconfirmation0(&_ContractEigenDAServiceManager.CallOpts, miniBatchHeader, blobHeader, nonSignerStakesAndSignature, additionalQuorumNumbersRequired)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactor) ConfirmBatch(opts *bind.TransactOpts, batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.contract.Transact(opts, "confirmBatch", batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerSession) ConfirmBatch(batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
	return _ContractEigenDAServiceManager.Contract.ConfirmBatch(&_ContractEigenDAServiceManager.TransactOpts, batchHeader, nonSignerStakesAndSignature)
}

// ConfirmBatch is a paid mutator transaction binding the contract method 0x7794965a.
//
// Solidity: function confirmBatch((bytes32,bytes,bytes,uint32) batchHeader, (uint32[],(uint256,uint256)[],(uint256,uint256)[],(uint256[2],uint256[2]),(uint256,uint256),uint32[],uint32[],uint32[][]) nonSignerStakesAndSignature) returns()
func (_ContractEigenDAServiceManager *ContractEigenDAServiceManagerTransactorSession) ConfirmBatch(batchHeader IEigenDAServiceManagerBatchHeader, nonSignerStakesAndSignature IBLSSignatureCheckerNonSignerStakesAndSignature) (*types.Transaction, error) {
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
