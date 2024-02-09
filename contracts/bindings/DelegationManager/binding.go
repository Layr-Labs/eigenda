// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractDelegationManager

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

// IDelegationManagerOperatorDetails is an auto generated low-level Go binding around an user-defined struct.
type IDelegationManagerOperatorDetails struct {
	EarningsReceiver         common.Address
	DelegationApprover       common.Address
	StakerOptOutWindowBlocks uint32
}

// IDelegationManagerQueuedWithdrawalParams is an auto generated low-level Go binding around an user-defined struct.
type IDelegationManagerQueuedWithdrawalParams struct {
	Strategies []common.Address
	Shares     []*big.Int
	Withdrawer common.Address
}

// IDelegationManagerWithdrawal is an auto generated low-level Go binding around an user-defined struct.
type IDelegationManagerWithdrawal struct {
	Staker      common.Address
	DelegatedTo common.Address
	Withdrawer  common.Address
	Nonce       *big.Int
	StartBlock  uint32
	Strategies  []common.Address
	Shares      []*big.Int
}

// ISignatureUtilsSignatureWithExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithExpiry struct {
	Signature []byte
	Expiry    *big.Int
}

// IStrategyManagerDeprecatedStructQueuedWithdrawal is an auto generated low-level Go binding around an user-defined struct.
type IStrategyManagerDeprecatedStructQueuedWithdrawal struct {
	Strategies           []common.Address
	Shares               []*big.Int
	Staker               common.Address
	WithdrawerAndNonce   IStrategyManagerDeprecatedStructWithdrawerAndNonce
	WithdrawalStartBlock uint32
	DelegatedAddress     common.Address
}

// IStrategyManagerDeprecatedStructWithdrawerAndNonce is an auto generated low-level Go binding around an user-defined struct.
type IStrategyManagerDeprecatedStructWithdrawerAndNonce struct {
	Withdrawer common.Address
	Nonce      *big.Int
}

// ContractDelegationManagerMetaData contains all meta data concerning the ContractDelegationManager contract.
var ContractDelegationManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_strategyManager\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"_slasher\",\"type\":\"address\",\"internalType\":\"contractISlasher\"},{\"name\":\"_eigenPodManager\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DELEGATION_APPROVAL_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_STAKER_OPT_OUT_WINDOW_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_WITHDRAWAL_DELAY_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKER_DELEGATION_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beaconChainETHStrategy\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateCurrentStakerDelegationDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateDelegationApprovalDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateStakerDelegationDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_stakerNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateWithdrawalRoot\",\"inputs\":[{\"name\":\"withdrawal\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"completeQueuedWithdrawal\",\"inputs\":[{\"name\":\"withdrawal\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"contractIERC20[]\"},{\"name\":\"middlewareTimesIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiveAsTokens\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeQueuedWithdrawals\",\"inputs\":[{\"name\":\"withdrawals\",\"type\":\"tuple[]\",\"internalType\":\"structIDelegationManager.Withdrawal[]\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"tokens\",\"type\":\"address[][]\",\"internalType\":\"contractIERC20[][]\"},{\"name\":\"middlewareTimesIndexes\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"receiveAsTokens\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cumulativeWithdrawalsQueued\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseDelegatedShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateTo\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approverSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateToBySignature\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegatedTo\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegationApprover\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegationApproverSaltIsSpent\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"domainSeparator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"earningsReceiver\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenPodManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDelegatableShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorShares\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWithdrawalDelay\",\"inputs\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseDelegatedShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minWithdrawalDelayBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"_withdrawalDelayBlocks\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isDelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateQueuedWithdrawals\",\"inputs\":[{\"name\":\"withdrawalsToMigrate\",\"type\":\"tuple[]\",\"internalType\":\"structIStrategyManager.DeprecatedStruct_QueuedWithdrawal[]\",\"components\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawerAndNonce\",\"type\":\"tuple\",\"internalType\":\"structIStrategyManager.DeprecatedStruct_WithdrawerAndNonce\",\"components\":[{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]},{\"name\":\"withdrawalStartBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"delegatedAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minWithdrawalDelayBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"modifyOperatorDetails\",\"inputs\":[{\"name\":\"newOperatorDetails\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorDetails\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorShares\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingWithdrawals\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queueWithdrawals\",\"inputs\":[{\"name\":\"queuedWithdrawalParams\",\"type\":\"tuple[]\",\"internalType\":\"structIDelegationManager.QueuedWithdrawalParams[]\",\"components\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerAsOperator\",\"inputs\":[{\"name\":\"registeringOperatorDetails\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStrategyWithdrawalDelayBlocks\",\"inputs\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"withdrawalDelayBlocks\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slasher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISlasher\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakerNonce\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakerOptOutWindowBlocks\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyWithdrawalDelayBlocks\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"undelegate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"withdrawalRoots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperatorMetadataURI\",\"inputs\":[{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinWithdrawalDelayBlocksSet\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDetailsModified\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOperatorDetails\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorMetadataURIUpdated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorDetails\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSharesDecreased\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSharesIncreased\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerDelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerForceUndelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerUndelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyWithdrawalDelayBlocksSet\",\"inputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalCompleted\",\"inputs\":[{\"name\":\"withdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalMigrated\",\"inputs\":[{\"name\":\"oldWithdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"newWithdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalQueued\",\"inputs\":[{\"name\":\"withdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"withdrawal\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"anonymous\":false}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200616138038062006161833981016040819052620000359162000140565b6001600160a01b0380841660805280821660c052821660a0526200005862000065565b50504660e0525062000194565b600054610100900460ff1615620000d25760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116101562000125576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b03811681146200013d57600080fd5b50565b6000806000606084860312156200015657600080fd5b8351620001638162000127565b6020850151909350620001768162000127565b6040850151909250620001898162000127565b809150509250925092565b60805160a05160c05160e051615f3162000230600039600061274d0152600081816105cc015281816110930152818161140f01528181611de101528181612aa6015281816141190152614605015260006107a901526000818161051401528181611061015281816113dd0152818161170501528181611e7501528181612b7301528181612cf60152818161423f01526146ab0152615f316000f3fe608060405234801561001057600080fd5b506004361061035d5760003560e01c80635f966f14116101d3578063b7f06ebe11610104578063cf80873e116100a2578063f16172b01161007c578063f16172b01461094f578063f2fde38b14610962578063f698da2514610975578063fabc1cbc1461097d57600080fd5b8063cf80873e14610908578063da8be86414610929578063eea9064b1461093c57600080fd5b8063c488375a116100de578063c488375a14610825578063c5e480db14610845578063c94b5111146108eb578063ca661c04146108fe57600080fd5b8063b7f06ebe146107cb578063bb45fef2146107ee578063c448feb81461081c57600080fd5b8063886f1195116101715780639104c3191161014b5780639104c3191461075657806399be81c814610771578063a178848414610784578063b1344271146107a457600080fd5b8063886f1195146107125780638da5cb5b14610725578063900413471461073657600080fd5b80636d70f7ae116101ad5780636d70f7ae146106b9578063715018a6146106cc578063778e55f3146106d45780637f548071146106ff57600080fd5b80635f966f141461065157806360d7faed1461067d57806365da12641461069057600080fd5b806328a573ae116102ad5780634665bcda1161024b578063597b36da11610225578063597b36da146106005780635ac86ab7146106135780635c975abb146106365780635cfe8d2c1461063e57600080fd5b80634665bcda146105c75780634fc40b61146105ee578063595c6a67146105f857600080fd5b806339b70e381161028757806339b70e381461050f5780633cdeb5e01461054e5780633e28391d1461057d57806343377382146105a057600080fd5b806328a573ae146104c957806329c77d4f146104dc57806333404396146104fc57600080fd5b8063132d49671161031a57806316928365116102f457806316928365146104435780631bbce0911461047c57806320606b701461048f57806322bf40e4146104b657600080fd5b8063132d49671461040a578063136439dd1461041d5780631522bf021461043057600080fd5b80630449ca391461036257806304a4f979146103885780630b9f487a146103af5780630dd8dd02146103c25780630f589e59146103e257806310d67a2f146103f7575b600080fd5b610375610370366004614abb565b610990565b6040519081526020015b60405180910390f35b6103757f3b89fca151cbe5122d58acee86cf184413d751d585779bdd19d3bbfa3a306dce81565b6103756103bd366004614b21565b610a15565b6103d56103d0366004614abb565b610ad7565b60405161037f9190614b7c565b6103f56103f0366004614c19565b610e53565b005b6103f5610405366004614c6c565b610fa3565b6103f5610418366004614c90565b611056565b6103f561042b366004614cd1565b61110d565b6103f561043e366004614cea565b61124c565b610375610451366004614c6c565b6001600160a01b0316600090815260996020526040902060010154600160a01b900463ffffffff1690565b61037561048a366004614c90565b611260565b6103757f8cad95687ba82c2ce50e74f7b754645e5117c3a5bec8151c0726d5857980a86681565b6103f56104c4366004614d55565b61128e565b6103f56104d7366004614c90565b6113d2565b6103756104ea366004614c6c565b609b6020526000908152604090205481565b6103f561050a366004614dfc565b611482565b6105367f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161037f565b61053661055c366004614c6c565b6001600160a01b039081166000908152609960205260409020600101541690565b61059061058b366004614c6c565b6115bf565b604051901515815260200161037f565b6103757f39111bc4a4d688e1f685123d7497d4615370152a8ee4a0593e647bd06ad8bb0b81565b6105367f000000000000000000000000000000000000000000000000000000000000000081565b6103756213c68081565b6103f56115df565b61037561060e36600461511b565b6116a6565b610590610621366004615157565b606654600160ff9092169190911b9081161490565b606654610375565b6103f561064c3660046151cd565b6116d6565b61053661065f366004614c6c565b6001600160a01b039081166000908152609960205260409020541690565b6103f561068b36600461532d565b611981565b61053661069e366004614c6c565b609a602052600090815260409020546001600160a01b031681565b6105906106c7366004614c6c565b611a1c565b6103f5611a3c565b6103756106e23660046153bc565b609860209081526000928352604080842090915290825290205481565b6103f561070d36600461549d565b611a50565b606554610536906001600160a01b031681565b6033546001600160a01b0316610536565b61074961074436600461552d565b611b55565b60405161037f91906155b7565b61053673beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac081565b6103f561077f3660046155ca565b611c2f565b610375610792366004614c6c565b609f6020526000908152604090205481565b6105367f000000000000000000000000000000000000000000000000000000000000000081565b6105906107d9366004614cd1565b609e6020526000908152604090205460ff1681565b6105906107fc3660046155ff565b609c60209081526000928352604080842090915290825290205460ff1681565b610375609d5481565b610375610833366004614c6c565b60a16020526000908152604090205481565b6108b5610853366004614c6c565b6040805160608082018352600080835260208084018290529284018190526001600160a01b03948516815260998352839020835191820184528054851682526001015493841691810191909152600160a01b90920463ffffffff169082015290565b6040805182516001600160a01b039081168252602080850151909116908201529181015163ffffffff169082015260600161037f565b6103756108f936600461562b565b611d01565b61037562034bc081565b61091b610916366004614c6c565b611dba565b60405161037f9291906156ac565b6103d5610937366004614c6c565b612172565b6103f561094a3660046156d1565b612636565b6103f561095d366004615729565b612642565b6103f5610970366004614c6c565b6126d3565b610375612749565b6103f561098b366004614cd1565b612787565b609d54600090815b83811015610a0d57600060a160008787858181106109b8576109b8615745565b90506020020160208101906109cd9190614c6c565b6001600160a01b03166001600160a01b03168152602001908152602001600020549050828111156109fc578092505b50610a0681615771565b9050610998565b509392505050565b604080517f3b89fca151cbe5122d58acee86cf184413d751d585779bdd19d3bbfa3a306dce6020808301919091526001600160a01b038681168385015288811660608401528716608083015260a0820185905260c08083018590528351808403909101815260e0909201909252805191012060009081610a93612749565b60405161190160f01b602082015260228101919091526042810183905260620160408051808303601f19018152919052805160209091012098975050505050505050565b60665460609060019060029081161415610b0c5760405162461bcd60e51b8152600401610b039061578c565b60405180910390fd5b6000836001600160401b03811115610b2657610b26614e9e565b604051908082528060200260200182016040528015610b4f578160200160208202803683370190505b50336000908152609a60205260408120549192506001600160a01b03909116905b85811015610e4857868682818110610b8a57610b8a615745565b9050602002810190610b9c91906157c3565b610baa9060208101906157e3565b9050878783818110610bbe57610bbe615745565b9050602002810190610bd091906157c3565b610bda90806157e3565b905014610c4f5760405162461bcd60e51b815260206004820152603860248201527f44656c65676174696f6e4d616e616765722e717565756557697468647261776160448201527f6c3a20696e707574206c656e677468206d69736d6174636800000000000000006064820152608401610b03565b6000878783818110610c6357610c63615745565b9050602002810190610c7591906157c3565b610c86906060810190604001614c6c565b6001600160a01b03161415610d145760405162461bcd60e51b815260206004820152604860248201527f44656c65676174696f6e4d616e616765722e717565756557697468647261776160448201527f6c3a206d7573742070726f766964652076616c6964207769746864726177616c606482015267206164647265737360c01b608482015260a401610b03565b610e193383898985818110610d2b57610d2b615745565b9050602002810190610d3d91906157c3565b610d4e906060810190604001614c6c565b8a8a86818110610d6057610d60615745565b9050602002810190610d7291906157c3565b610d7c90806157e3565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152508e92508d9150889050818110610dc257610dc2615745565b9050602002810190610dd491906157c3565b610de29060208101906157e3565b808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506128e392505050565b838281518110610e2b57610e2b615745565b602090810291909101015280610e4081615771565b915050610b70565b509095945050505050565b336000908152609960205260409020546001600160a01b031615610eed5760405162461bcd60e51b815260206004820152604560248201527f44656c65676174696f6e4d616e616765722e726567697374657241734f70657260448201527f61746f723a206f70657261746f722068617320616c72656164792072656769736064820152641d195c995960da1b608482015260a401610b03565b610ef73384612ea3565b604080518082019091526060815260006020820152610f19338083600061313f565b336001600160a01b03167f8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e285604051610f52919061582c565b60405180910390a2336001600160a01b03167f02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b67080908484604051610f9592919061587e565b60405180910390a250505050565b606560009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610ff6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061101a91906158ad565b6001600160a01b0316336001600160a01b03161461104a5760405162461bcd60e51b8152600401610b03906158ca565b611053816134e4565b50565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806110b55750336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016145b6110d15760405162461bcd60e51b8152600401610b0390615914565b6110da836115bf565b15611108576001600160a01b038084166000908152609a602052604090205416611106818585856135db565b505b505050565b60655460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015611155573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111799190615971565b6111955760405162461bcd60e51b8152600401610b039061598e565b6066548181161461120e5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e70617573653a20696e76616c696420617474656d70742060448201527f746f20756e70617573652066756e6374696f6e616c69747900000000000000006064820152608401610b03565b606681905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d906020015b60405180910390a250565b611254613656565b611106848484846136b0565b6001600160a01b0383166000908152609b602052604081205461128585828686611d01565b95945050505050565b600054610100900460ff16158080156112ae5750600054600160ff909116105b806112c85750303b1580156112c8575060005460ff166001145b61132b5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610b03565b6000805460ff19166001179055801561134e576000805461ff0019166101001790555b61135888886138d6565b6113606139bc565b60975561136c89613a53565b61137586613aa5565b611381858585856136b0565b80156113c7576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050505050565b336001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806114315750336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016145b61144d5760405162461bcd60e51b8152600401610b0390615914565b611456836115bf565b15611108576001600160a01b038084166000908152609a60205260409020541661110681858585613bab565b606654600290600490811614156114ab5760405162461bcd60e51b8152600401610b039061578c565b600260c95414156114fe5760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610b03565b600260c95560005b888110156115ae5761159e8a8a8381811061152357611523615745565b905060200281019061153591906159d6565b89898481811061154757611547615745565b905060200281019061155991906157e3565b89898681811061156b5761156b615745565b9050602002013588888781811061158457611584615745565b905060200201602081019061159991906159ec565b613c26565b6115a781615771565b9050611506565b5050600160c9555050505050505050565b6001600160a01b039081166000908152609a602052604090205416151590565b60655460405163237dfb4760e11b81523360048201526001600160a01b03909116906346fbf68e90602401602060405180830381865afa158015611627573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061164b9190615971565b6116675760405162461bcd60e51b8152600401610b039061598e565b600019606681905560405190815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a2565b6000816040516020016116b99190615a7d565b604051602081830303815290604052805190602001209050919050565b60005b815181101561197d5760008282815181106116f6576116f6615745565b602002602001015190506000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663cd293f6f846040518263ffffffff1660e01b815260040161174f9190615a90565b60408051808303816000875af115801561176d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117919190615b3c565b91509150811561196f576040808401516001600160a01b0381166000908152609f602052918220805491928291906117c883615771565b919050555060006040518060e00160405280846001600160a01b031681526020018760a001516001600160a01b031681526020018760600151600001516001600160a01b03168152602001838152602001876080015163ffffffff16815260200187600001518152602001876020015181525090506000611848826116a6565b6000818152609e602052604090205490915060ff16156118de5760405162461bcd60e51b815260206004820152604560248201527f44656c65676174696f6e4d616e616765722e6d6967726174655175657565645760448201527f69746864726177616c733a207769746864726177616c20616c72656164792065606482015264786973747360d81b608482015260a401610b03565b6000818152609e602052604090819020805460ff19166001179055517f9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f9906119299083908590615b6a565b60405180910390a160408051868152602081018390527fdc00758b65eef71dc3780c04ebe36cab6bdb266c3a698187e29e0f0dca012630910160405180910390a1505050505b8360010193505050506116d9565b5050565b606654600290600490811614156119aa5760405162461bcd60e51b8152600401610b039061578c565b600260c95414156119fd5760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610b03565b600260c955611a0f8686868686613c26565b5050600160c95550505050565b6001600160a01b0390811660009081526099602052604090205416151590565b611a44613656565b611a4e6000613a53565b565b4283602001511015611ad45760405162461bcd60e51b815260206004820152604160248201527f44656c65676174696f6e4d616e616765722e64656c6567617465546f4279536960448201527f676e61747572653a207374616b6572207369676e6174757265206578706972656064820152601960fa1b608482015260a401610b03565b6000609b6000876001600160a01b03166001600160a01b031681526020019081526020016000205490506000611b108783888860200151611d01565b6001600160a01b0388166000908152609b602052604090206001840190558551909150611b409088908390614410565b611b4c8787868661313f565b50505050505050565b6060600082516001600160401b03811115611b7257611b72614e9e565b604051908082528060200260200182016040528015611b9b578160200160208202803683370190505b50905060005b8351811015610a0d576001600160a01b03851660009081526098602052604081208551909190869084908110611bd957611bd9615745565b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002054828281518110611c1457611c14615745565b6020908102919091010152611c2881615771565b9050611ba1565b611c3833611a1c565b611cba5760405162461bcd60e51b815260206004820152604760248201527f44656c65676174696f6e4d616e616765722e7570646174654f70657261746f7260448201527f4d657461646174615552493a2063616c6c6572206d75737420626520616e206f6064820152663832b930ba37b960c91b608482015260a401610b03565b336001600160a01b03167f02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b67080908383604051611cf592919061587e565b60405180910390a25050565b604080517f39111bc4a4d688e1f685123d7497d4615370152a8ee4a0593e647bd06ad8bb0b6020808301919091526001600160a01b0387811683850152851660608301526080820186905260a08083018590528351808403909101815260c0909201909252805191012060009081611d77612749565b60405161190160f01b602082015260228101919091526042810183905260620160408051808303601f190181529190528051602090910120979650505050505050565b6040516360f4062b60e01b81526001600160a01b03828116600483015260609182916000917f0000000000000000000000000000000000000000000000000000000000000000909116906360f4062b90602401602060405180830381865afa158015611e2a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e4e9190615b83565b6040516394f649dd60e01b81526001600160a01b03868116600483015291925060009182917f0000000000000000000000000000000000000000000000000000000000000000909116906394f649dd90602401600060405180830381865afa158015611ebe573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611ee69190810190615bf7565b9150915060008313611efd57909590945092505050565b606080835160001415611fb7576040805160018082528183019092529060208083019080368337505060408051600180825281830190925292945090506020808301908036833701905050905073beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac082600081518110611f7257611f72615745565b60200260200101906001600160a01b031690816001600160a01b0316815250508481600081518110611fa657611fa6615745565b602002602001018181525050612165565b8351611fc4906001615cb1565b6001600160401b03811115611fdb57611fdb614e9e565b604051908082528060200260200182016040528015612004578160200160208202803683370190505b50915081516001600160401b0381111561202057612020614e9e565b604051908082528060200260200182016040528015612049578160200160208202803683370190505b50905060005b84518110156120e35784818151811061206a5761206a615745565b602002602001015183828151811061208457612084615745565b60200260200101906001600160a01b031690816001600160a01b0316815250508381815181106120b6576120b6615745565b60200260200101518282815181106120d0576120d0615745565b602090810291909101015260010161204f565b5073beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac082600184516121089190615cc9565b8151811061211857612118615745565b60200260200101906001600160a01b031690816001600160a01b0316815250508481600184516121489190615cc9565b8151811061215857612158615745565b6020026020010181815250505b9097909650945050505050565b6066546060906001906002908116141561219e5760405162461bcd60e51b8152600401610b039061578c565b6121a7836115bf565b6122275760405162461bcd60e51b8152602060048201526044602482018190527f44656c65676174696f6e4d616e616765722e756e64656c65676174653a207374908201527f616b6572206d7573742062652064656c65676174656420746f20756e64656c656064820152636761746560e01b608482015260a401610b03565b61223083611a1c565b156122a35760405162461bcd60e51b815260206004820152603d60248201527f44656c65676174696f6e4d616e616765722e756e64656c65676174653a206f7060448201527f657261746f72732063616e6e6f7420626520756e64656c6567617465640000006064820152608401610b03565b6001600160a01b03831661231f5760405162461bcd60e51b815260206004820152603c60248201527f44656c65676174696f6e4d616e616765722e756e64656c65676174653a20636160448201527f6e6e6f7420756e64656c6567617465207a65726f2061646472657373000000006064820152608401610b03565b6001600160a01b038084166000818152609a6020526040902054909116903314806123525750336001600160a01b038216145b8061237957506001600160a01b038181166000908152609960205260409020600101541633145b6123eb5760405162461bcd60e51b815260206004820152603d60248201527f44656c65676174696f6e4d616e616765722e756e64656c65676174653a20636160448201527f6c6c65722063616e6e6f7420756e64656c6567617465207374616b65720000006064820152608401610b03565b6000806123f786611dba565b9092509050336001600160a01b0387161461244d57826001600160a01b0316866001600160a01b03167ff0eddf07e6ea14f388b47e1e94a0f464ecbd9eed4171130e0fc0e99fb4030a8a60405160405180910390a35b826001600160a01b0316866001600160a01b03167ffee30966a256b71e14bc0ebfc94315e28ef4a97a7131a9e2b7a310a73af4467660405160405180910390a36001600160a01b0386166000908152609a6020526040902080546001600160a01b031916905581516124cf57604080516000815260208101909152945061262d565b81516001600160401b038111156124e8576124e8614e9e565b604051908082528060200260200182016040528015612511578160200160208202803683370190505b50945060005b825181101561262b5760408051600180825281830190925260009160208083019080368337505060408051600180825281830190925292935060009291506020808301908036833701905050905084838151811061257757612577615745565b60200260200101518260008151811061259257612592615745565b60200260200101906001600160a01b031690816001600160a01b0316815250508383815181106125c4576125c4615745565b6020026020010151816000815181106125df576125df615745565b6020026020010181815250506125f889878b85856128e3565b88848151811061260a5761260a615745565b6020026020010181815250505050808061262390615771565b915050612517565b505b50505050919050565b6111083384848461313f565b61264b33611a1c565b6126c95760405162461bcd60e51b815260206004820152604360248201527f44656c65676174696f6e4d616e616765722e6d6f646966794f70657261746f7260448201527f44657461696c733a2063616c6c6572206d75737420626520616e206f706572616064820152623a37b960e91b608482015260a401610b03565b6110533382612ea3565b6126db613656565b6001600160a01b0381166127405760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610b03565b61105381613a53565b60007f000000000000000000000000000000000000000000000000000000000000000046141561277a575060975490565b6127826139bc565b905090565b606560009054906101000a90046001600160a01b03166001600160a01b031663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156127da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127fe91906158ad565b6001600160a01b0316336001600160a01b03161461282e5760405162461bcd60e51b8152600401610b03906158ca565b6066541981196066541916146128ac5760405162461bcd60e51b815260206004820152603860248201527f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060448201527f7420746f2070617573652066756e6374696f6e616c69747900000000000000006064820152608401610b03565b606681905560405181815233907f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c90602001611241565b60006001600160a01b03861661297a5760405162461bcd60e51b815260206004820152605060248201527f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160448201527f6e6451756575655769746864726177616c3a207374616b65722063616e6e6f7460648201526f206265207a65726f206164647265737360801b608482015260a401610b03565b8251612a045760405162461bcd60e51b815260206004820152604d60248201527f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160448201527f6e6451756575655769746864726177616c3a207374726174656769657320636160648201526c6e6e6f7420626520656d70747960981b608482015260a401610b03565b60005b8351811015612db1576001600160a01b03861615612a5d57612a5d8688868481518110612a3657612a36615745565b6020026020010151868581518110612a5057612a50615745565b60200260200101516135db565b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac06001600160a01b0316848281518110612a8d57612a8d615745565b60200260200101516001600160a01b03161415612b56577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663beffbb8988858481518110612ae657612ae6615745565b60200260200101516040518363ffffffff1660e01b8152600401612b1f9291906001600160a01b03929092168252602082015260400190565b600060405180830381600087803b158015612b3957600080fd5b505af1158015612b4d573d6000803e3d6000fd5b50505050612da9565b846001600160a01b0316876001600160a01b03161480612c2857507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316639b4da03d858381518110612bb257612bb2615745565b60200260200101516040518263ffffffff1660e01b8152600401612be591906001600160a01b0391909116815260200190565b602060405180830381865afa158015612c02573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c269190615971565b155b612cf45760405162461bcd60e51b8152602060048201526084602482018190527f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160448301527f6e6451756575655769746864726177616c3a2077697468647261776572206d7560648301527f73742062652073616d652061646472657373206173207374616b657220696620908201527f746869726450617274795472616e7366657273466f7262696464656e2061726560a482015263081cd95d60e21b60c482015260e401610b03565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316638c80d4e588868481518110612d3657612d36615745565b6020026020010151868581518110612d5057612d50615745565b60200260200101516040518463ffffffff1660e01b8152600401612d7693929190615ce0565b600060405180830381600087803b158015612d9057600080fd5b505af1158015612da4573d6000803e3d6000fd5b505050505b600101612a07565b506001600160a01b0386166000908152609f60205260408120805491829190612dd983615771565b919050555060006040518060e00160405280896001600160a01b03168152602001886001600160a01b03168152602001876001600160a01b031681526020018381526020014363ffffffff1681526020018681526020018581525090506000612e41826116a6565b6000818152609e602052604090819020805460ff19166001179055519091507f9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f990612e8f9083908590615b6a565b60405180910390a198975050505050505050565b6000612eb26020830183614c6c565b6001600160a01b03161415612f4c5760405162461bcd60e51b815260206004820152605460248201527f44656c65676174696f6e4d616e616765722e5f7365744f70657261746f72446560448201527f7461696c733a2063616e6e6f742073657420606561726e696e677352656365696064820152737665726020746f207a65726f206164647265737360601b608482015260a401610b03565b6213c680612f606060830160408401615d04565b63ffffffff1611156130155760405162461bcd60e51b815260206004820152606c60248201527f44656c65676174696f6e4d616e616765722e5f7365744f70657261746f72446560448201527f7461696c733a207374616b65724f70744f757457696e646f77426c6f636b732060648201527f63616e6e6f74206265203e204d41585f5354414b45525f4f50545f4f55545f5760848201526b494e444f575f424c4f434b5360a01b60a482015260c401610b03565b6001600160a01b0382166000908152609960205260409081902060010154600160a01b900463ffffffff16906130519060608401908401615d04565b63ffffffff1610156130e75760405162461bcd60e51b815260206004820152605360248201527f44656c65676174696f6e4d616e616765722e5f7365744f70657261746f72446560448201527f7461696c733a207374616b65724f70744f757457696e646f77426c6f636b732060648201527218d85b9b9bdd08189948191958dc99585cd959606a1b608482015260a401610b03565b6001600160a01b0382166000908152609960205260409020819061310b8282615d41565b505060405133907ffebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac90611cf590849061582c565b606654600090600190811614156131685760405162461bcd60e51b8152600401610b039061578c565b613171856115bf565b156131ee5760405162461bcd60e51b815260206004820152604160248201527f44656c65676174696f6e4d616e616765722e5f64656c65676174653a2073746160448201527f6b657220697320616c7265616479206163746976656c792064656c65676174656064820152601960fa1b608482015260a401610b03565b6131f784611a1c565b6132775760405162461bcd60e51b815260206004820152604560248201527f44656c65676174696f6e4d616e616765722e5f64656c65676174653a206f706560448201527f7261746f72206973206e6f74207265676973746572656420696e20456967656e6064820152642630bcb2b960d91b608482015260a401610b03565b6001600160a01b038085166000908152609960205260409020600101541680158015906132ad5750336001600160a01b03821614155b80156132c25750336001600160a01b03861614155b1561342f5742846020015110156133415760405162461bcd60e51b815260206004820152603760248201527f44656c65676174696f6e4d616e616765722e5f64656c65676174653a2061707060448201527f726f766572207369676e617475726520657870697265640000000000000000006064820152608401610b03565b6001600160a01b0381166000908152609c6020908152604080832086845290915290205460ff16156133db5760405162461bcd60e51b815260206004820152603760248201527f44656c65676174696f6e4d616e616765722e5f64656c65676174653a2061707060448201527f726f76657253616c7420616c7265616479207370656e740000000000000000006064820152608401610b03565b6001600160a01b0381166000908152609c6020908152604080832086845282528220805460ff1916600117905585015161341c908890889085908890610a15565b905061342d82828760000151614410565b505b6001600160a01b038681166000818152609a602052604080822080546001600160a01b031916948a169485179055517fc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d87433049190a360008061348e88611dba565b9150915060005b82518110156113c7576134dc888a8584815181106134b5576134b5615745565b60200260200101518585815181106134cf576134cf615745565b6020026020010151613bab565b600101613495565b6001600160a01b0381166135725760405162461bcd60e51b815260206004820152604960248201527f5061757361626c652e5f73657450617573657252656769737472793a206e657760448201527f50617573657252656769737472792063616e6e6f7420626520746865207a65726064820152686f206164647265737360b81b608482015260a401610b03565b606554604080516001600160a01b03928316815291831660208301527f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6910160405180910390a1606580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b03808516600090815260986020908152604080832093861683529290529081208054839290613612908490615cc9565b92505081905550836001600160a01b03167f6909600037b75d7b4733aedd815442b5ec018a827751c832aaff64eba5d6d2dd848484604051610f9593929190615ce0565b6033546001600160a01b03163314611a4e5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610b03565b8281146137385760405162461bcd60e51b815260206004820152604a60248201527f44656c65676174696f6e4d616e616765722e5f7365745374726174656779576960448201527f746864726177616c44656c6179426c6f636b733a20696e707574206c656e67746064820152690d040dad2e6dac2e8c6d60b31b608482015260a401610b03565b8260005b818110156138ce57600086868381811061375857613758615745565b905060200201602081019061376d9190614c6c565b6001600160a01b038116600090815260a1602052604081205491925086868581811061379b5761379b615745565b90506020020135905062034bc081111561385f5760405162461bcd60e51b815260206004820152607360248201527f44656c65676174696f6e4d616e616765722e5f7365745374726174656779576960448201527f746864726177616c44656c6179426c6f636b733a205f7769746864726177616c60648201527f44656c6179426c6f636b732063616e6e6f74206265203e204d41585f5749544860848201527244524157414c5f44454c41595f424c4f434b5360681b60a482015260c401610b03565b6001600160a01b038316600081815260a160209081526040918290208490558151928352820184905281018290527f0e7efa738e8b0ce6376a0c1af471655540d2e9a81647d7b09ed823018426576d9060600160405180910390a1505050806138c790615771565b905061373c565b505050505050565b6065546001600160a01b03161580156138f757506001600160a01b03821615155b6139795760405162461bcd60e51b815260206004820152604760248201527f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960448201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c6064820152666564206f6e636560c81b608482015260a401610b03565b606681905560405181815233907fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d9060200160405180910390a261197d826134e4565b604080518082018252600a81526922b4b3b2b72630bcb2b960b11b60209182015281517f8cad95687ba82c2ce50e74f7b754645e5117c3a5bec8151c0726d5857980a866818301527f71b625cfad44bac63b13dba07f2e1d6084ee04b6f8752101ece6126d584ee6ea81840152466060820152306080808301919091528351808303909101815260a0909101909252815191012090565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b62034bc0811115613b6a5760405162461bcd60e51b815260206004820152607860248201527f44656c65676174696f6e4d616e616765722e5f696e697469616c697a654d696e60448201527f5769746864726177616c44656c6179426c6f636b733a205f6d696e576974686460648201527f726177616c44656c6179426c6f636b732063616e6e6f74206265203e204d415860848201527f5f5749544844524157414c5f44454c41595f424c4f434b53000000000000000060a482015260c401610b03565b609d5460408051918252602082018390527fafa003cd76f87ff9d62b35beea889920f33c0c42b8d45b74954d61d50f4b6b69910160405180910390a1609d55565b6001600160a01b03808516600090815260986020908152604080832093861683529290529081208054839290613be2908490615cb1565b92505081905550836001600160a01b03167f1ec042c965e2edd7107b51188ee0f383e22e76179041ab3a9d18ff151405166c848484604051610f9593929190615ce0565b6000613c3461060e87615da4565b6000818152609e602052604090205490915060ff16613cb55760405162461bcd60e51b81526020600482015260436024820152600080516020615edc83398151915260448201527f645769746864726177616c3a20616374696f6e206973206e6f7420696e20717560648201526265756560e81b608482015260a401610b03565b609d544390613cca60a0890160808a01615d04565b63ffffffff16613cda9190615cb1565b1115613d625760405162461bcd60e51b815260206004820152605f6024820152600080516020615edc83398151915260448201527f645769746864726177616c3a206d696e5769746864726177616c44656c61794260648201527f6c6f636b7320706572696f6420686173206e6f74207965742070617373656400608482015260a401610b03565b613d726060870160408801614c6c565b6001600160a01b0316336001600160a01b031614613dff5760405162461bcd60e51b81526020600482015260506024820152600080516020615edc83398151915260448201527f645769746864726177616c3a206f6e6c7920776974686472617765722063616e60648201526f1031b7b6b83632ba329030b1ba34b7b760811b608482015260a401610b03565b8115613e8157613e1260a08701876157e3565b85149050613e815760405162461bcd60e51b81526020600482015260426024820152600080516020615edc83398151915260448201527f645769746864726177616c3a20696e707574206c656e677468206d69736d61746064820152610c6d60f31b608482015260a401610b03565b6000818152609e60205260409020805460ff191690558115613fe65760005b613ead60a08801886157e3565b9050811015613fe0574360a16000613ec860a08b018b6157e3565b85818110613ed857613ed8615745565b9050602002016020810190613eed9190614c6c565b6001600160a01b03168152602081019190915260400160002054613f1760a08a0160808b01615d04565b63ffffffff16613f279190615cb1565b1115613f455760405162461bcd60e51b8152600401610b0390615db6565b613fd8613f556020890189614c6c565b33613f6360a08b018b6157e3565b85818110613f7357613f73615745565b9050602002016020810190613f889190614c6c565b613f9560c08c018c6157e3565b86818110613fa557613fa5615745565b905060200201358a8a87818110613fbe57613fbe615745565b9050602002016020810190613fd39190614c6c565b6145ca565b600101613ea0565b506143d5565b336000908152609a60205260408120546001600160a01b0316905b61400e60a08901896157e3565b90508110156143d2574360a1600061402960a08c018c6157e3565b8581811061403957614039615745565b905060200201602081019061404e9190614c6c565b6001600160a01b0316815260208101919091526040016000205461407860a08b0160808c01615d04565b63ffffffff166140889190615cb1565b11156140a65760405162461bcd60e51b8152600401610b0390615db6565b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac06140c860a08a018a6157e3565b838181106140d8576140d8615745565b90506020020160208101906140ed9190614c6c565b6001600160a01b0316141561423d57600061410b60208a018a614c6c565b905060006001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016630e81073c8361414c60c08e018e6157e3565b8781811061415c5761415c615745565b6040516001600160e01b031960e087901b1681526001600160a01b03909416600485015260200291909101356024830152506044016020604051808303816000875af11580156141b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141d49190615b83565b6001600160a01b038084166000908152609a602052604090205491925016801561423557614235818461420a60a08f018f6157e3565b8881811061421a5761421a615745565b905060200201602081019061422f9190614c6c565b85613bab565b5050506143ca565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663c4623ea13389898581811061427f5761427f615745565b90506020020160208101906142949190614c6c565b6142a160a08d018d6157e3565b868181106142b1576142b1615745565b90506020020160208101906142c69190614c6c565b6142d360c08e018e6157e3565b878181106142e3576142e3615745565b60405160e088901b6001600160e01b03191681526001600160a01b03968716600482015294861660248601529290941660448401526020909102013560648201526084019050600060405180830381600087803b15801561434357600080fd5b505af1158015614357573d6000803e3d6000fd5b505050506001600160a01b038216156143ca576143ca823361437c60a08c018c6157e3565b8581811061438c5761438c615745565b90506020020160208101906143a19190614c6c565b6143ae60c08d018d6157e3565b868181106143be576143be615745565b90506020020135613bab565b600101614001565b50505b6040518181527fc97098c2f658800b4df29001527f7324bcdffcf6e8751a699ab920a1eced5b1d9060200160405180910390a1505050505050565b6001600160a01b0383163b1561452a57604051630b135d3f60e11b808252906001600160a01b03851690631626ba7e906144509086908690600401615e3e565b602060405180830381865afa15801561446d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906144919190615e9b565b6001600160e01b031916146111085760405162461bcd60e51b815260206004820152605360248201527f454950313237315369676e61747572655574696c732e636865636b5369676e6160448201527f747572655f454950313237313a2045524331323731207369676e6174757265206064820152721d995c9a599a58d85d1a5bdb8819985a5b1959606a1b608482015260a401610b03565b826001600160a01b031661453e838361470a565b6001600160a01b0316146111085760405162461bcd60e51b815260206004820152604760248201527f454950313237315369676e61747572655574696c732e636865636b5369676e6160448201527f747572655f454950313237313a207369676e6174757265206e6f742066726f6d6064820152661039b4b3b732b960c91b608482015260a401610b03565b6001600160a01b03831673beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac014156146755760405162387b1360e81b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063387b13009061463e90889088908790600401615ce0565b600060405180830381600087803b15801561465857600080fd5b505af115801561466c573d6000803e3d6000fd5b50505050614703565b60405163c608c7f360e01b81526001600160a01b03858116600483015284811660248301526044820184905282811660648301527f0000000000000000000000000000000000000000000000000000000000000000169063c608c7f390608401600060405180830381600087803b1580156146ef57600080fd5b505af11580156113c7573d6000803e3d6000fd5b5050505050565b60008060006147198585614726565b91509150610a0d81614796565b60008082516041141561475d5760208301516040840151606085015160001a61475187828585614951565b9450945050505061478f565b825160401415614787576020830151604084015161477c868383614a3e565b93509350505061478f565b506000905060025b9250929050565b60008160048111156147aa576147aa615ec5565b14156147b35750565b60018160048111156147c7576147c7615ec5565b14156148155760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e617475726500000000000000006044820152606401610b03565b600281600481111561482957614829615ec5565b14156148775760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606401610b03565b600381600481111561488b5761488b615ec5565b14156148e45760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b6064820152608401610b03565b60048160048111156148f8576148f8615ec5565b14156110535760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b6064820152608401610b03565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08311156149885750600090506003614a35565b8460ff16601b141580156149a057508460ff16601c14155b156149b15750600090506004614a35565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015614a05573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116614a2e57600060019250925050614a35565b9150600090505b94509492505050565b6000806001600160ff1b03831681614a5b60ff86901c601b615cb1565b9050614a6987828885614951565b935093505050935093915050565b60008083601f840112614a8957600080fd5b5081356001600160401b03811115614aa057600080fd5b6020830191508360208260051b850101111561478f57600080fd5b60008060208385031215614ace57600080fd5b82356001600160401b03811115614ae457600080fd5b614af085828601614a77565b90969095509350505050565b6001600160a01b038116811461105357600080fd5b8035614b1c81614afc565b919050565b600080600080600060a08688031215614b3957600080fd5b8535614b4481614afc565b94506020860135614b5481614afc565b93506040860135614b6481614afc565b94979396509394606081013594506080013592915050565b6020808252825182820181905260009190848201906040850190845b81811015614bb457835183529284019291840191600101614b98565b50909695505050505050565b600060608284031215614bd257600080fd5b50919050565b60008083601f840112614bea57600080fd5b5081356001600160401b03811115614c0157600080fd5b60208301915083602082850101111561478f57600080fd5b600080600060808486031215614c2e57600080fd5b614c388585614bc0565b925060608401356001600160401b03811115614c5357600080fd5b614c5f86828701614bd8565b9497909650939450505050565b600060208284031215614c7e57600080fd5b8135614c8981614afc565b9392505050565b600080600060608486031215614ca557600080fd5b8335614cb081614afc565b92506020840135614cc081614afc565b929592945050506040919091013590565b600060208284031215614ce357600080fd5b5035919050565b60008060008060408587031215614d0057600080fd5b84356001600160401b0380821115614d1757600080fd5b614d2388838901614a77565b90965094506020870135915080821115614d3c57600080fd5b50614d4987828801614a77565b95989497509550505050565b60008060008060008060008060c0898b031215614d7157600080fd5b8835614d7c81614afc565b97506020890135614d8c81614afc565b9650604089013595506060890135945060808901356001600160401b0380821115614db657600080fd5b614dc28c838d01614a77565b909650945060a08b0135915080821115614ddb57600080fd5b50614de88b828c01614a77565b999c989b5096995094979396929594505050565b6000806000806000806000806080898b031215614e1857600080fd5b88356001600160401b0380821115614e2f57600080fd5b614e3b8c838d01614a77565b909a50985060208b0135915080821115614e5457600080fd5b614e608c838d01614a77565b909850965060408b0135915080821115614e7957600080fd5b614e858c838d01614a77565b909650945060608b0135915080821115614ddb57600080fd5b634e487b7160e01b600052604160045260246000fd5b60405160e081016001600160401b0381118282101715614ed657614ed6614e9e565b60405290565b604080519081016001600160401b0381118282101715614ed657614ed6614e9e565b60405160c081016001600160401b0381118282101715614ed657614ed6614e9e565b604051601f8201601f191681016001600160401b0381118282101715614f4857614f48614e9e565b604052919050565b63ffffffff8116811461105357600080fd5b8035614b1c81614f50565b60006001600160401b03821115614f8657614f86614e9e565b5060051b60200190565b600082601f830112614fa157600080fd5b81356020614fb6614fb183614f6d565b614f20565b82815260059290921b84018101918181019086841115614fd557600080fd5b8286015b84811015614ff9578035614fec81614afc565b8352918301918301614fd9565b509695505050505050565b600082601f83011261501557600080fd5b81356020615025614fb183614f6d565b82815260059290921b8401810191818101908684111561504457600080fd5b8286015b84811015614ff95780358352918301918301615048565b600060e0828403121561507157600080fd5b615079614eb4565b905061508482614b11565b815261509260208301614b11565b60208201526150a360408301614b11565b6040820152606082013560608201526150be60808301614f62565b608082015260a08201356001600160401b03808211156150dd57600080fd5b6150e985838601614f90565b60a084015260c084013591508082111561510257600080fd5b5061510f84828501615004565b60c08301525092915050565b60006020828403121561512d57600080fd5b81356001600160401b0381111561514357600080fd5b61514f8482850161505f565b949350505050565b60006020828403121561516957600080fd5b813560ff81168114614c8957600080fd5b60006040828403121561518c57600080fd5b615194614edc565b905081356151a181614afc565b815260208201356bffffffffffffffffffffffff811681146151c257600080fd5b602082015292915050565b600060208083850312156151e057600080fd5b82356001600160401b03808211156151f757600080fd5b818501915085601f83011261520b57600080fd5b8135615219614fb182614f6d565b81815260059190911b8301840190848101908883111561523857600080fd5b8585015b83811015615312578035858111156152545760008081fd5b860160e0818c03601f190181131561526c5760008081fd5b615274614efe565b89830135888111156152865760008081fd5b6152948e8c83870101614f90565b825250604080840135898111156152ab5760008081fd5b6152b98f8d83880101615004565b8c8401525060606152cb818601614b11565b82840152608091506152df8f83870161517a565b908301526152ef60c08501614f62565b908201526152fe838301614b11565b60a08201528552505091860191860161523c565b5098975050505050505050565b801515811461105357600080fd5b60008060008060006080868803121561534557600080fd5b85356001600160401b038082111561535c57600080fd5b9087019060e0828a03121561537057600080fd5b9095506020870135908082111561538657600080fd5b5061539388828901614a77565b9095509350506040860135915060608601356153ae8161531f565b809150509295509295909350565b600080604083850312156153cf57600080fd5b82356153da81614afc565b915060208301356153ea81614afc565b809150509250929050565b60006040828403121561540757600080fd5b61540f614edc565b905081356001600160401b038082111561542857600080fd5b818401915084601f83011261543c57600080fd5b813560208282111561545057615450614e9e565b615462601f8301601f19168201614f20565b9250818352868183860101111561547857600080fd5b8181850182850137600081838501015282855280860135818601525050505092915050565b600080600080600060a086880312156154b557600080fd5b85356154c081614afc565b945060208601356154d081614afc565b935060408601356001600160401b03808211156154ec57600080fd5b6154f889838a016153f5565b9450606088013591508082111561550e57600080fd5b5061551b888289016153f5565b95989497509295608001359392505050565b6000806040838503121561554057600080fd5b823561554b81614afc565b915060208301356001600160401b0381111561556657600080fd5b61557285828601614f90565b9150509250929050565b600081518084526020808501945080840160005b838110156155ac57815187529582019590820190600101615590565b509495945050505050565b602081526000614c89602083018461557c565b600080602083850312156155dd57600080fd5b82356001600160401b038111156155f357600080fd5b614af085828601614bd8565b6000806040838503121561561257600080fd5b823561561d81614afc565b946020939093013593505050565b6000806000806080858703121561564157600080fd5b843561564c81614afc565b935060208501359250604085013561566381614afc565b9396929550929360600135925050565b600081518084526020808501945080840160005b838110156155ac5781516001600160a01b031687529582019590820190600101615687565b6040815260006156bf6040830185615673565b8281036020840152611285818561557c565b6000806000606084860312156156e657600080fd5b83356156f181614afc565b925060208401356001600160401b0381111561570c57600080fd5b615718868287016153f5565b925050604084013590509250925092565b60006060828403121561573b57600080fd5b614c898383614bc0565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006000198214156157855761578561575b565b5060010190565b60208082526019908201527f5061757361626c653a20696e6465782069732070617573656400000000000000604082015260600190565b60008235605e198336030181126157d957600080fd5b9190910192915050565b6000808335601e198436030181126157fa57600080fd5b8301803591506001600160401b0382111561581457600080fd5b6020019150600581901b360382131561478f57600080fd5b60608101823561583b81614afc565b6001600160a01b03908116835260208401359061585782614afc565b166020830152604083013561586b81614f50565b63ffffffff811660408401525092915050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b6000602082840312156158bf57600080fd5b8151614c8981614afc565b6020808252602a908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526939903ab73830bab9b2b960b11b606082015260800190565b60208082526037908201527f44656c65676174696f6e4d616e616765723a206f6e6c7953747261746567794d60408201527f616e616765724f72456967656e506f644d616e61676572000000000000000000606082015260800190565b60006020828403121561598357600080fd5b8151614c898161531f565b60208082526028908201527f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160408201526739903830bab9b2b960c11b606082015260800190565b6000823560de198336030181126157d957600080fd5b6000602082840312156159fe57600080fd5b8135614c898161531f565b600060018060a01b03808351168452806020840151166020850152806040840151166040850152506060820151606084015263ffffffff608083015116608084015260a082015160e060a0850152615a6460e0850182615673565b905060c083015184820360c0860152611285828261557c565b602081526000614c896020830184615a09565b602081526000825160e06020840152615aad610100840182615673565b90506020840151601f19848303016040850152615aca828261557c565b915050604084015160018060a01b03808216606086015260608601519150808251166080860152506bffffffffffffffffffffffff60208201511660a0850152506080840151615b2260c085018263ffffffff169052565b5060a08401516001600160a01b03811660e0850152610a0d565b60008060408385031215615b4f57600080fd5b8251615b5a8161531f565b6020939093015192949293505050565b82815260406020820152600061514f6040830184615a09565b600060208284031215615b9557600080fd5b5051919050565b600082601f830112615bad57600080fd5b81516020615bbd614fb183614f6d565b82815260059290921b84018101918181019086841115615bdc57600080fd5b8286015b84811015614ff95780518352918301918301615be0565b60008060408385031215615c0a57600080fd5b82516001600160401b0380821115615c2157600080fd5b818501915085601f830112615c3557600080fd5b81516020615c45614fb183614f6d565b82815260059290921b84018101918181019089841115615c6457600080fd5b948201945b83861015615c8b578551615c7c81614afc565b82529482019490820190615c69565b91880151919650909350505080821115615ca457600080fd5b5061557285828601615b9c565b60008219821115615cc457615cc461575b565b500190565b600082821015615cdb57615cdb61575b565b500390565b6001600160a01b039384168152919092166020820152604081019190915260600190565b600060208284031215615d1657600080fd5b8135614c8981614f50565b80546001600160a01b0319166001600160a01b0392909216919091179055565b8135615d4c81614afc565b615d568183615d21565b50600181016020830135615d6981614afc565b615d738183615d21565b506040830135615d8281614f50565b815463ffffffff60a01b191660a09190911b63ffffffff60a01b161790555050565b6000615db0368361505f565b92915050565b6020808252606e90820152600080516020615edc83398151915260408201527f645769746864726177616c3a207769746864726177616c44656c6179426c6f6360608201527f6b7320706572696f6420686173206e6f74207965742070617373656420666f7260808201526d207468697320737472617465677960901b60a082015260c00190565b82815260006020604081840152835180604085015260005b81811015615e7257858101830151858201606001528201615e56565b81811115615e84576000606083870101525b50601f01601f191692909201606001949350505050565b600060208284031215615ead57600080fd5b81516001600160e01b031981168114614c8957600080fd5b634e487b7160e01b600052602160045260246000fdfe44656c65676174696f6e4d616e616765722e5f636f6d706c6574655175657565a264697066735822122037eb7270f460c1e6fe189c4a6cd1560f73ba29417e87d480df46dfe110c1b5fc64736f6c634300080c0033",
}

// ContractDelegationManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractDelegationManagerMetaData.ABI instead.
var ContractDelegationManagerABI = ContractDelegationManagerMetaData.ABI

// ContractDelegationManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractDelegationManagerMetaData.Bin instead.
var ContractDelegationManagerBin = ContractDelegationManagerMetaData.Bin

// DeployContractDelegationManager deploys a new Ethereum contract, binding an instance of ContractDelegationManager to it.
func DeployContractDelegationManager(auth *bind.TransactOpts, backend bind.ContractBackend, _strategyManager common.Address, _slasher common.Address, _eigenPodManager common.Address) (common.Address, *types.Transaction, *ContractDelegationManager, error) {
	parsed, err := ContractDelegationManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractDelegationManagerBin), backend, _strategyManager, _slasher, _eigenPodManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractDelegationManager{ContractDelegationManagerCaller: ContractDelegationManagerCaller{contract: contract}, ContractDelegationManagerTransactor: ContractDelegationManagerTransactor{contract: contract}, ContractDelegationManagerFilterer: ContractDelegationManagerFilterer{contract: contract}}, nil
}

// ContractDelegationManager is an auto generated Go binding around an Ethereum contract.
type ContractDelegationManager struct {
	ContractDelegationManagerCaller     // Read-only binding to the contract
	ContractDelegationManagerTransactor // Write-only binding to the contract
	ContractDelegationManagerFilterer   // Log filterer for contract events
}

// ContractDelegationManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractDelegationManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDelegationManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractDelegationManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDelegationManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractDelegationManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractDelegationManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractDelegationManagerSession struct {
	Contract     *ContractDelegationManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ContractDelegationManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractDelegationManagerCallerSession struct {
	Contract *ContractDelegationManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ContractDelegationManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractDelegationManagerTransactorSession struct {
	Contract     *ContractDelegationManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ContractDelegationManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractDelegationManagerRaw struct {
	Contract *ContractDelegationManager // Generic contract binding to access the raw methods on
}

// ContractDelegationManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractDelegationManagerCallerRaw struct {
	Contract *ContractDelegationManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ContractDelegationManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractDelegationManagerTransactorRaw struct {
	Contract *ContractDelegationManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractDelegationManager creates a new instance of ContractDelegationManager, bound to a specific deployed contract.
func NewContractDelegationManager(address common.Address, backend bind.ContractBackend) (*ContractDelegationManager, error) {
	contract, err := bindContractDelegationManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManager{ContractDelegationManagerCaller: ContractDelegationManagerCaller{contract: contract}, ContractDelegationManagerTransactor: ContractDelegationManagerTransactor{contract: contract}, ContractDelegationManagerFilterer: ContractDelegationManagerFilterer{contract: contract}}, nil
}

// NewContractDelegationManagerCaller creates a new read-only instance of ContractDelegationManager, bound to a specific deployed contract.
func NewContractDelegationManagerCaller(address common.Address, caller bind.ContractCaller) (*ContractDelegationManagerCaller, error) {
	contract, err := bindContractDelegationManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerCaller{contract: contract}, nil
}

// NewContractDelegationManagerTransactor creates a new write-only instance of ContractDelegationManager, bound to a specific deployed contract.
func NewContractDelegationManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractDelegationManagerTransactor, error) {
	contract, err := bindContractDelegationManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerTransactor{contract: contract}, nil
}

// NewContractDelegationManagerFilterer creates a new log filterer instance of ContractDelegationManager, bound to a specific deployed contract.
func NewContractDelegationManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractDelegationManagerFilterer, error) {
	contract, err := bindContractDelegationManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerFilterer{contract: contract}, nil
}

// bindContractDelegationManager binds a generic wrapper to an already deployed contract.
func bindContractDelegationManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractDelegationManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractDelegationManager *ContractDelegationManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractDelegationManager.Contract.ContractDelegationManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractDelegationManager *ContractDelegationManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.ContractDelegationManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractDelegationManager *ContractDelegationManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.ContractDelegationManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractDelegationManager *ContractDelegationManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractDelegationManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractDelegationManager *ContractDelegationManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractDelegationManager *ContractDelegationManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.contract.Transact(opts, method, params...)
}

// DELEGATIONAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0x04a4f979.
//
// Solidity: function DELEGATION_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DELEGATIONAPPROVALTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "DELEGATION_APPROVAL_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DELEGATIONAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0x04a4f979.
//
// Solidity: function DELEGATION_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) DELEGATIONAPPROVALTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DELEGATIONAPPROVALTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// DELEGATIONAPPROVALTYPEHASH is a free data retrieval call binding the contract method 0x04a4f979.
//
// Solidity: function DELEGATION_APPROVAL_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DELEGATIONAPPROVALTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DELEGATIONAPPROVALTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// DOMAINTYPEHASH is a free data retrieval call binding the contract method 0x20606b70.
//
// Solidity: function DOMAIN_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DOMAINTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "DOMAIN_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINTYPEHASH is a free data retrieval call binding the contract method 0x20606b70.
//
// Solidity: function DOMAIN_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) DOMAINTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DOMAINTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// DOMAINTYPEHASH is a free data retrieval call binding the contract method 0x20606b70.
//
// Solidity: function DOMAIN_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DOMAINTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DOMAINTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// MAXSTAKEROPTOUTWINDOWBLOCKS is a free data retrieval call binding the contract method 0x4fc40b61.
//
// Solidity: function MAX_STAKER_OPT_OUT_WINDOW_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) MAXSTAKEROPTOUTWINDOWBLOCKS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "MAX_STAKER_OPT_OUT_WINDOW_BLOCKS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSTAKEROPTOUTWINDOWBLOCKS is a free data retrieval call binding the contract method 0x4fc40b61.
//
// Solidity: function MAX_STAKER_OPT_OUT_WINDOW_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) MAXSTAKEROPTOUTWINDOWBLOCKS() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MAXSTAKEROPTOUTWINDOWBLOCKS(&_ContractDelegationManager.CallOpts)
}

// MAXSTAKEROPTOUTWINDOWBLOCKS is a free data retrieval call binding the contract method 0x4fc40b61.
//
// Solidity: function MAX_STAKER_OPT_OUT_WINDOW_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) MAXSTAKEROPTOUTWINDOWBLOCKS() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MAXSTAKEROPTOUTWINDOWBLOCKS(&_ContractDelegationManager.CallOpts)
}

// MAXWITHDRAWALDELAYBLOCKS is a free data retrieval call binding the contract method 0xca661c04.
//
// Solidity: function MAX_WITHDRAWAL_DELAY_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) MAXWITHDRAWALDELAYBLOCKS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "MAX_WITHDRAWAL_DELAY_BLOCKS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXWITHDRAWALDELAYBLOCKS is a free data retrieval call binding the contract method 0xca661c04.
//
// Solidity: function MAX_WITHDRAWAL_DELAY_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) MAXWITHDRAWALDELAYBLOCKS() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MAXWITHDRAWALDELAYBLOCKS(&_ContractDelegationManager.CallOpts)
}

// MAXWITHDRAWALDELAYBLOCKS is a free data retrieval call binding the contract method 0xca661c04.
//
// Solidity: function MAX_WITHDRAWAL_DELAY_BLOCKS() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) MAXWITHDRAWALDELAYBLOCKS() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MAXWITHDRAWALDELAYBLOCKS(&_ContractDelegationManager.CallOpts)
}

// STAKERDELEGATIONTYPEHASH is a free data retrieval call binding the contract method 0x43377382.
//
// Solidity: function STAKER_DELEGATION_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) STAKERDELEGATIONTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "STAKER_DELEGATION_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// STAKERDELEGATIONTYPEHASH is a free data retrieval call binding the contract method 0x43377382.
//
// Solidity: function STAKER_DELEGATION_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) STAKERDELEGATIONTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.STAKERDELEGATIONTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// STAKERDELEGATIONTYPEHASH is a free data retrieval call binding the contract method 0x43377382.
//
// Solidity: function STAKER_DELEGATION_TYPEHASH() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) STAKERDELEGATIONTYPEHASH() ([32]byte, error) {
	return _ContractDelegationManager.Contract.STAKERDELEGATIONTYPEHASH(&_ContractDelegationManager.CallOpts)
}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) BeaconChainETHStrategy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "beaconChainETHStrategy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) BeaconChainETHStrategy() (common.Address, error) {
	return _ContractDelegationManager.Contract.BeaconChainETHStrategy(&_ContractDelegationManager.CallOpts)
}

// BeaconChainETHStrategy is a free data retrieval call binding the contract method 0x9104c319.
//
// Solidity: function beaconChainETHStrategy() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) BeaconChainETHStrategy() (common.Address, error) {
	return _ContractDelegationManager.Contract.BeaconChainETHStrategy(&_ContractDelegationManager.CallOpts)
}

// CalculateCurrentStakerDelegationDigestHash is a free data retrieval call binding the contract method 0x1bbce091.
//
// Solidity: function calculateCurrentStakerDelegationDigestHash(address staker, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) CalculateCurrentStakerDelegationDigestHash(opts *bind.CallOpts, staker common.Address, operator common.Address, expiry *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "calculateCurrentStakerDelegationDigestHash", staker, operator, expiry)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateCurrentStakerDelegationDigestHash is a free data retrieval call binding the contract method 0x1bbce091.
//
// Solidity: function calculateCurrentStakerDelegationDigestHash(address staker, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) CalculateCurrentStakerDelegationDigestHash(staker common.Address, operator common.Address, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateCurrentStakerDelegationDigestHash(&_ContractDelegationManager.CallOpts, staker, operator, expiry)
}

// CalculateCurrentStakerDelegationDigestHash is a free data retrieval call binding the contract method 0x1bbce091.
//
// Solidity: function calculateCurrentStakerDelegationDigestHash(address staker, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) CalculateCurrentStakerDelegationDigestHash(staker common.Address, operator common.Address, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateCurrentStakerDelegationDigestHash(&_ContractDelegationManager.CallOpts, staker, operator, expiry)
}

// CalculateDelegationApprovalDigestHash is a free data retrieval call binding the contract method 0x0b9f487a.
//
// Solidity: function calculateDelegationApprovalDigestHash(address staker, address operator, address _delegationApprover, bytes32 approverSalt, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) CalculateDelegationApprovalDigestHash(opts *bind.CallOpts, staker common.Address, operator common.Address, _delegationApprover common.Address, approverSalt [32]byte, expiry *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "calculateDelegationApprovalDigestHash", staker, operator, _delegationApprover, approverSalt, expiry)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateDelegationApprovalDigestHash is a free data retrieval call binding the contract method 0x0b9f487a.
//
// Solidity: function calculateDelegationApprovalDigestHash(address staker, address operator, address _delegationApprover, bytes32 approverSalt, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) CalculateDelegationApprovalDigestHash(staker common.Address, operator common.Address, _delegationApprover common.Address, approverSalt [32]byte, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateDelegationApprovalDigestHash(&_ContractDelegationManager.CallOpts, staker, operator, _delegationApprover, approverSalt, expiry)
}

// CalculateDelegationApprovalDigestHash is a free data retrieval call binding the contract method 0x0b9f487a.
//
// Solidity: function calculateDelegationApprovalDigestHash(address staker, address operator, address _delegationApprover, bytes32 approverSalt, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) CalculateDelegationApprovalDigestHash(staker common.Address, operator common.Address, _delegationApprover common.Address, approverSalt [32]byte, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateDelegationApprovalDigestHash(&_ContractDelegationManager.CallOpts, staker, operator, _delegationApprover, approverSalt, expiry)
}

// CalculateStakerDelegationDigestHash is a free data retrieval call binding the contract method 0xc94b5111.
//
// Solidity: function calculateStakerDelegationDigestHash(address staker, uint256 _stakerNonce, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) CalculateStakerDelegationDigestHash(opts *bind.CallOpts, staker common.Address, _stakerNonce *big.Int, operator common.Address, expiry *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "calculateStakerDelegationDigestHash", staker, _stakerNonce, operator, expiry)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateStakerDelegationDigestHash is a free data retrieval call binding the contract method 0xc94b5111.
//
// Solidity: function calculateStakerDelegationDigestHash(address staker, uint256 _stakerNonce, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) CalculateStakerDelegationDigestHash(staker common.Address, _stakerNonce *big.Int, operator common.Address, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateStakerDelegationDigestHash(&_ContractDelegationManager.CallOpts, staker, _stakerNonce, operator, expiry)
}

// CalculateStakerDelegationDigestHash is a free data retrieval call binding the contract method 0xc94b5111.
//
// Solidity: function calculateStakerDelegationDigestHash(address staker, uint256 _stakerNonce, address operator, uint256 expiry) view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) CalculateStakerDelegationDigestHash(staker common.Address, _stakerNonce *big.Int, operator common.Address, expiry *big.Int) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateStakerDelegationDigestHash(&_ContractDelegationManager.CallOpts, staker, _stakerNonce, operator, expiry)
}

// CalculateWithdrawalRoot is a free data retrieval call binding the contract method 0x597b36da.
//
// Solidity: function calculateWithdrawalRoot((address,address,address,uint256,uint32,address[],uint256[]) withdrawal) pure returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) CalculateWithdrawalRoot(opts *bind.CallOpts, withdrawal IDelegationManagerWithdrawal) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "calculateWithdrawalRoot", withdrawal)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateWithdrawalRoot is a free data retrieval call binding the contract method 0x597b36da.
//
// Solidity: function calculateWithdrawalRoot((address,address,address,uint256,uint32,address[],uint256[]) withdrawal) pure returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) CalculateWithdrawalRoot(withdrawal IDelegationManagerWithdrawal) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateWithdrawalRoot(&_ContractDelegationManager.CallOpts, withdrawal)
}

// CalculateWithdrawalRoot is a free data retrieval call binding the contract method 0x597b36da.
//
// Solidity: function calculateWithdrawalRoot((address,address,address,uint256,uint32,address[],uint256[]) withdrawal) pure returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) CalculateWithdrawalRoot(withdrawal IDelegationManagerWithdrawal) ([32]byte, error) {
	return _ContractDelegationManager.Contract.CalculateWithdrawalRoot(&_ContractDelegationManager.CallOpts, withdrawal)
}

// CumulativeWithdrawalsQueued is a free data retrieval call binding the contract method 0xa1788484.
//
// Solidity: function cumulativeWithdrawalsQueued(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) CumulativeWithdrawalsQueued(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "cumulativeWithdrawalsQueued", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CumulativeWithdrawalsQueued is a free data retrieval call binding the contract method 0xa1788484.
//
// Solidity: function cumulativeWithdrawalsQueued(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) CumulativeWithdrawalsQueued(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.CumulativeWithdrawalsQueued(&_ContractDelegationManager.CallOpts, arg0)
}

// CumulativeWithdrawalsQueued is a free data retrieval call binding the contract method 0xa1788484.
//
// Solidity: function cumulativeWithdrawalsQueued(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) CumulativeWithdrawalsQueued(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.CumulativeWithdrawalsQueued(&_ContractDelegationManager.CallOpts, arg0)
}

// DelegatedTo is a free data retrieval call binding the contract method 0x65da1264.
//
// Solidity: function delegatedTo(address ) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DelegatedTo(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "delegatedTo", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelegatedTo is a free data retrieval call binding the contract method 0x65da1264.
//
// Solidity: function delegatedTo(address ) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) DelegatedTo(arg0 common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.DelegatedTo(&_ContractDelegationManager.CallOpts, arg0)
}

// DelegatedTo is a free data retrieval call binding the contract method 0x65da1264.
//
// Solidity: function delegatedTo(address ) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DelegatedTo(arg0 common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.DelegatedTo(&_ContractDelegationManager.CallOpts, arg0)
}

// DelegationApprover is a free data retrieval call binding the contract method 0x3cdeb5e0.
//
// Solidity: function delegationApprover(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DelegationApprover(opts *bind.CallOpts, operator common.Address) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "delegationApprover", operator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelegationApprover is a free data retrieval call binding the contract method 0x3cdeb5e0.
//
// Solidity: function delegationApprover(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) DelegationApprover(operator common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.DelegationApprover(&_ContractDelegationManager.CallOpts, operator)
}

// DelegationApprover is a free data retrieval call binding the contract method 0x3cdeb5e0.
//
// Solidity: function delegationApprover(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DelegationApprover(operator common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.DelegationApprover(&_ContractDelegationManager.CallOpts, operator)
}

// DelegationApproverSaltIsSpent is a free data retrieval call binding the contract method 0xbb45fef2.
//
// Solidity: function delegationApproverSaltIsSpent(address , bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DelegationApproverSaltIsSpent(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "delegationApproverSaltIsSpent", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DelegationApproverSaltIsSpent is a free data retrieval call binding the contract method 0xbb45fef2.
//
// Solidity: function delegationApproverSaltIsSpent(address , bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerSession) DelegationApproverSaltIsSpent(arg0 common.Address, arg1 [32]byte) (bool, error) {
	return _ContractDelegationManager.Contract.DelegationApproverSaltIsSpent(&_ContractDelegationManager.CallOpts, arg0, arg1)
}

// DelegationApproverSaltIsSpent is a free data retrieval call binding the contract method 0xbb45fef2.
//
// Solidity: function delegationApproverSaltIsSpent(address , bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DelegationApproverSaltIsSpent(arg0 common.Address, arg1 [32]byte) (bool, error) {
	return _ContractDelegationManager.Contract.DelegationApproverSaltIsSpent(&_ContractDelegationManager.CallOpts, arg0, arg1)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCaller) DomainSeparator(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "domainSeparator")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerSession) DomainSeparator() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DomainSeparator(&_ContractDelegationManager.CallOpts)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) DomainSeparator() ([32]byte, error) {
	return _ContractDelegationManager.Contract.DomainSeparator(&_ContractDelegationManager.CallOpts)
}

// EarningsReceiver is a free data retrieval call binding the contract method 0x5f966f14.
//
// Solidity: function earningsReceiver(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) EarningsReceiver(opts *bind.CallOpts, operator common.Address) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "earningsReceiver", operator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EarningsReceiver is a free data retrieval call binding the contract method 0x5f966f14.
//
// Solidity: function earningsReceiver(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) EarningsReceiver(operator common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.EarningsReceiver(&_ContractDelegationManager.CallOpts, operator)
}

// EarningsReceiver is a free data retrieval call binding the contract method 0x5f966f14.
//
// Solidity: function earningsReceiver(address operator) view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) EarningsReceiver(operator common.Address) (common.Address, error) {
	return _ContractDelegationManager.Contract.EarningsReceiver(&_ContractDelegationManager.CallOpts, operator)
}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) EigenPodManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "eigenPodManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) EigenPodManager() (common.Address, error) {
	return _ContractDelegationManager.Contract.EigenPodManager(&_ContractDelegationManager.CallOpts)
}

// EigenPodManager is a free data retrieval call binding the contract method 0x4665bcda.
//
// Solidity: function eigenPodManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) EigenPodManager() (common.Address, error) {
	return _ContractDelegationManager.Contract.EigenPodManager(&_ContractDelegationManager.CallOpts)
}

// GetDelegatableShares is a free data retrieval call binding the contract method 0xcf80873e.
//
// Solidity: function getDelegatableShares(address staker) view returns(address[], uint256[])
func (_ContractDelegationManager *ContractDelegationManagerCaller) GetDelegatableShares(opts *bind.CallOpts, staker common.Address) ([]common.Address, []*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "getDelegatableShares", staker)

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

// GetDelegatableShares is a free data retrieval call binding the contract method 0xcf80873e.
//
// Solidity: function getDelegatableShares(address staker) view returns(address[], uint256[])
func (_ContractDelegationManager *ContractDelegationManagerSession) GetDelegatableShares(staker common.Address) ([]common.Address, []*big.Int, error) {
	return _ContractDelegationManager.Contract.GetDelegatableShares(&_ContractDelegationManager.CallOpts, staker)
}

// GetDelegatableShares is a free data retrieval call binding the contract method 0xcf80873e.
//
// Solidity: function getDelegatableShares(address staker) view returns(address[], uint256[])
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) GetDelegatableShares(staker common.Address) ([]common.Address, []*big.Int, error) {
	return _ContractDelegationManager.Contract.GetDelegatableShares(&_ContractDelegationManager.CallOpts, staker)
}

// GetOperatorShares is a free data retrieval call binding the contract method 0x90041347.
//
// Solidity: function getOperatorShares(address operator, address[] strategies) view returns(uint256[])
func (_ContractDelegationManager *ContractDelegationManagerCaller) GetOperatorShares(opts *bind.CallOpts, operator common.Address, strategies []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "getOperatorShares", operator, strategies)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetOperatorShares is a free data retrieval call binding the contract method 0x90041347.
//
// Solidity: function getOperatorShares(address operator, address[] strategies) view returns(uint256[])
func (_ContractDelegationManager *ContractDelegationManagerSession) GetOperatorShares(operator common.Address, strategies []common.Address) ([]*big.Int, error) {
	return _ContractDelegationManager.Contract.GetOperatorShares(&_ContractDelegationManager.CallOpts, operator, strategies)
}

// GetOperatorShares is a free data retrieval call binding the contract method 0x90041347.
//
// Solidity: function getOperatorShares(address operator, address[] strategies) view returns(uint256[])
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) GetOperatorShares(operator common.Address, strategies []common.Address) ([]*big.Int, error) {
	return _ContractDelegationManager.Contract.GetOperatorShares(&_ContractDelegationManager.CallOpts, operator, strategies)
}

// GetWithdrawalDelay is a free data retrieval call binding the contract method 0x0449ca39.
//
// Solidity: function getWithdrawalDelay(address[] strategies) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) GetWithdrawalDelay(opts *bind.CallOpts, strategies []common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "getWithdrawalDelay", strategies)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWithdrawalDelay is a free data retrieval call binding the contract method 0x0449ca39.
//
// Solidity: function getWithdrawalDelay(address[] strategies) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) GetWithdrawalDelay(strategies []common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.GetWithdrawalDelay(&_ContractDelegationManager.CallOpts, strategies)
}

// GetWithdrawalDelay is a free data retrieval call binding the contract method 0x0449ca39.
//
// Solidity: function getWithdrawalDelay(address[] strategies) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) GetWithdrawalDelay(strategies []common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.GetWithdrawalDelay(&_ContractDelegationManager.CallOpts, strategies)
}

// IsDelegated is a free data retrieval call binding the contract method 0x3e28391d.
//
// Solidity: function isDelegated(address staker) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCaller) IsDelegated(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "isDelegated", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDelegated is a free data retrieval call binding the contract method 0x3e28391d.
//
// Solidity: function isDelegated(address staker) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerSession) IsDelegated(staker common.Address) (bool, error) {
	return _ContractDelegationManager.Contract.IsDelegated(&_ContractDelegationManager.CallOpts, staker)
}

// IsDelegated is a free data retrieval call binding the contract method 0x3e28391d.
//
// Solidity: function isDelegated(address staker) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) IsDelegated(staker common.Address) (bool, error) {
	return _ContractDelegationManager.Contract.IsDelegated(&_ContractDelegationManager.CallOpts, staker)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address operator) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCaller) IsOperator(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "isOperator", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address operator) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerSession) IsOperator(operator common.Address) (bool, error) {
	return _ContractDelegationManager.Contract.IsOperator(&_ContractDelegationManager.CallOpts, operator)
}

// IsOperator is a free data retrieval call binding the contract method 0x6d70f7ae.
//
// Solidity: function isOperator(address operator) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) IsOperator(operator common.Address) (bool, error) {
	return _ContractDelegationManager.Contract.IsOperator(&_ContractDelegationManager.CallOpts, operator)
}

// MinWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc448feb8.
//
// Solidity: function minWithdrawalDelayBlocks() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) MinWithdrawalDelayBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "minWithdrawalDelayBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc448feb8.
//
// Solidity: function minWithdrawalDelayBlocks() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) MinWithdrawalDelayBlocks() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MinWithdrawalDelayBlocks(&_ContractDelegationManager.CallOpts)
}

// MinWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc448feb8.
//
// Solidity: function minWithdrawalDelayBlocks() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) MinWithdrawalDelayBlocks() (*big.Int, error) {
	return _ContractDelegationManager.Contract.MinWithdrawalDelayBlocks(&_ContractDelegationManager.CallOpts)
}

// OperatorDetails is a free data retrieval call binding the contract method 0xc5e480db.
//
// Solidity: function operatorDetails(address operator) view returns((address,address,uint32))
func (_ContractDelegationManager *ContractDelegationManagerCaller) OperatorDetails(opts *bind.CallOpts, operator common.Address) (IDelegationManagerOperatorDetails, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "operatorDetails", operator)

	if err != nil {
		return *new(IDelegationManagerOperatorDetails), err
	}

	out0 := *abi.ConvertType(out[0], new(IDelegationManagerOperatorDetails)).(*IDelegationManagerOperatorDetails)

	return out0, err

}

// OperatorDetails is a free data retrieval call binding the contract method 0xc5e480db.
//
// Solidity: function operatorDetails(address operator) view returns((address,address,uint32))
func (_ContractDelegationManager *ContractDelegationManagerSession) OperatorDetails(operator common.Address) (IDelegationManagerOperatorDetails, error) {
	return _ContractDelegationManager.Contract.OperatorDetails(&_ContractDelegationManager.CallOpts, operator)
}

// OperatorDetails is a free data retrieval call binding the contract method 0xc5e480db.
//
// Solidity: function operatorDetails(address operator) view returns((address,address,uint32))
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) OperatorDetails(operator common.Address) (IDelegationManagerOperatorDetails, error) {
	return _ContractDelegationManager.Contract.OperatorDetails(&_ContractDelegationManager.CallOpts, operator)
}

// OperatorShares is a free data retrieval call binding the contract method 0x778e55f3.
//
// Solidity: function operatorShares(address , address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) OperatorShares(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "operatorShares", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperatorShares is a free data retrieval call binding the contract method 0x778e55f3.
//
// Solidity: function operatorShares(address , address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) OperatorShares(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.OperatorShares(&_ContractDelegationManager.CallOpts, arg0, arg1)
}

// OperatorShares is a free data retrieval call binding the contract method 0x778e55f3.
//
// Solidity: function operatorShares(address , address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) OperatorShares(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.OperatorShares(&_ContractDelegationManager.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) Owner() (common.Address, error) {
	return _ContractDelegationManager.Contract.Owner(&_ContractDelegationManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) Owner() (common.Address, error) {
	return _ContractDelegationManager.Contract.Owner(&_ContractDelegationManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCaller) Paused(opts *bind.CallOpts, index uint8) (bool, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "paused", index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerSession) Paused(index uint8) (bool, error) {
	return _ContractDelegationManager.Contract.Paused(&_ContractDelegationManager.CallOpts, index)
}

// Paused is a free data retrieval call binding the contract method 0x5ac86ab7.
//
// Solidity: function paused(uint8 index) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) Paused(index uint8) (bool, error) {
	return _ContractDelegationManager.Contract.Paused(&_ContractDelegationManager.CallOpts, index)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) Paused0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "paused0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) Paused0() (*big.Int, error) {
	return _ContractDelegationManager.Contract.Paused0(&_ContractDelegationManager.CallOpts)
}

// Paused0 is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) Paused0() (*big.Int, error) {
	return _ContractDelegationManager.Contract.Paused0(&_ContractDelegationManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) PauserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "pauserRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) PauserRegistry() (common.Address, error) {
	return _ContractDelegationManager.Contract.PauserRegistry(&_ContractDelegationManager.CallOpts)
}

// PauserRegistry is a free data retrieval call binding the contract method 0x886f1195.
//
// Solidity: function pauserRegistry() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) PauserRegistry() (common.Address, error) {
	return _ContractDelegationManager.Contract.PauserRegistry(&_ContractDelegationManager.CallOpts)
}

// PendingWithdrawals is a free data retrieval call binding the contract method 0xb7f06ebe.
//
// Solidity: function pendingWithdrawals(bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCaller) PendingWithdrawals(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "pendingWithdrawals", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PendingWithdrawals is a free data retrieval call binding the contract method 0xb7f06ebe.
//
// Solidity: function pendingWithdrawals(bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerSession) PendingWithdrawals(arg0 [32]byte) (bool, error) {
	return _ContractDelegationManager.Contract.PendingWithdrawals(&_ContractDelegationManager.CallOpts, arg0)
}

// PendingWithdrawals is a free data retrieval call binding the contract method 0xb7f06ebe.
//
// Solidity: function pendingWithdrawals(bytes32 ) view returns(bool)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) PendingWithdrawals(arg0 [32]byte) (bool, error) {
	return _ContractDelegationManager.Contract.PendingWithdrawals(&_ContractDelegationManager.CallOpts, arg0)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "slasher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) Slasher() (common.Address, error) {
	return _ContractDelegationManager.Contract.Slasher(&_ContractDelegationManager.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) Slasher() (common.Address, error) {
	return _ContractDelegationManager.Contract.Slasher(&_ContractDelegationManager.CallOpts)
}

// StakerNonce is a free data retrieval call binding the contract method 0x29c77d4f.
//
// Solidity: function stakerNonce(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) StakerNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "stakerNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakerNonce is a free data retrieval call binding the contract method 0x29c77d4f.
//
// Solidity: function stakerNonce(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) StakerNonce(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StakerNonce(&_ContractDelegationManager.CallOpts, arg0)
}

// StakerNonce is a free data retrieval call binding the contract method 0x29c77d4f.
//
// Solidity: function stakerNonce(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) StakerNonce(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StakerNonce(&_ContractDelegationManager.CallOpts, arg0)
}

// StakerOptOutWindowBlocks is a free data retrieval call binding the contract method 0x16928365.
//
// Solidity: function stakerOptOutWindowBlocks(address operator) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) StakerOptOutWindowBlocks(opts *bind.CallOpts, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "stakerOptOutWindowBlocks", operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakerOptOutWindowBlocks is a free data retrieval call binding the contract method 0x16928365.
//
// Solidity: function stakerOptOutWindowBlocks(address operator) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) StakerOptOutWindowBlocks(operator common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StakerOptOutWindowBlocks(&_ContractDelegationManager.CallOpts, operator)
}

// StakerOptOutWindowBlocks is a free data retrieval call binding the contract method 0x16928365.
//
// Solidity: function stakerOptOutWindowBlocks(address operator) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) StakerOptOutWindowBlocks(operator common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StakerOptOutWindowBlocks(&_ContractDelegationManager.CallOpts, operator)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCaller) StrategyManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "strategyManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerSession) StrategyManager() (common.Address, error) {
	return _ContractDelegationManager.Contract.StrategyManager(&_ContractDelegationManager.CallOpts)
}

// StrategyManager is a free data retrieval call binding the contract method 0x39b70e38.
//
// Solidity: function strategyManager() view returns(address)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) StrategyManager() (common.Address, error) {
	return _ContractDelegationManager.Contract.StrategyManager(&_ContractDelegationManager.CallOpts)
}

// StrategyWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc488375a.
//
// Solidity: function strategyWithdrawalDelayBlocks(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCaller) StrategyWithdrawalDelayBlocks(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractDelegationManager.contract.Call(opts, &out, "strategyWithdrawalDelayBlocks", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StrategyWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc488375a.
//
// Solidity: function strategyWithdrawalDelayBlocks(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerSession) StrategyWithdrawalDelayBlocks(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StrategyWithdrawalDelayBlocks(&_ContractDelegationManager.CallOpts, arg0)
}

// StrategyWithdrawalDelayBlocks is a free data retrieval call binding the contract method 0xc488375a.
//
// Solidity: function strategyWithdrawalDelayBlocks(address ) view returns(uint256)
func (_ContractDelegationManager *ContractDelegationManagerCallerSession) StrategyWithdrawalDelayBlocks(arg0 common.Address) (*big.Int, error) {
	return _ContractDelegationManager.Contract.StrategyWithdrawalDelayBlocks(&_ContractDelegationManager.CallOpts, arg0)
}

// CompleteQueuedWithdrawal is a paid mutator transaction binding the contract method 0x60d7faed.
//
// Solidity: function completeQueuedWithdrawal((address,address,address,uint256,uint32,address[],uint256[]) withdrawal, address[] tokens, uint256 middlewareTimesIndex, bool receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) CompleteQueuedWithdrawal(opts *bind.TransactOpts, withdrawal IDelegationManagerWithdrawal, tokens []common.Address, middlewareTimesIndex *big.Int, receiveAsTokens bool) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "completeQueuedWithdrawal", withdrawal, tokens, middlewareTimesIndex, receiveAsTokens)
}

// CompleteQueuedWithdrawal is a paid mutator transaction binding the contract method 0x60d7faed.
//
// Solidity: function completeQueuedWithdrawal((address,address,address,uint256,uint32,address[],uint256[]) withdrawal, address[] tokens, uint256 middlewareTimesIndex, bool receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) CompleteQueuedWithdrawal(withdrawal IDelegationManagerWithdrawal, tokens []common.Address, middlewareTimesIndex *big.Int, receiveAsTokens bool) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.CompleteQueuedWithdrawal(&_ContractDelegationManager.TransactOpts, withdrawal, tokens, middlewareTimesIndex, receiveAsTokens)
}

// CompleteQueuedWithdrawal is a paid mutator transaction binding the contract method 0x60d7faed.
//
// Solidity: function completeQueuedWithdrawal((address,address,address,uint256,uint32,address[],uint256[]) withdrawal, address[] tokens, uint256 middlewareTimesIndex, bool receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) CompleteQueuedWithdrawal(withdrawal IDelegationManagerWithdrawal, tokens []common.Address, middlewareTimesIndex *big.Int, receiveAsTokens bool) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.CompleteQueuedWithdrawal(&_ContractDelegationManager.TransactOpts, withdrawal, tokens, middlewareTimesIndex, receiveAsTokens)
}

// CompleteQueuedWithdrawals is a paid mutator transaction binding the contract method 0x33404396.
//
// Solidity: function completeQueuedWithdrawals((address,address,address,uint256,uint32,address[],uint256[])[] withdrawals, address[][] tokens, uint256[] middlewareTimesIndexes, bool[] receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) CompleteQueuedWithdrawals(opts *bind.TransactOpts, withdrawals []IDelegationManagerWithdrawal, tokens [][]common.Address, middlewareTimesIndexes []*big.Int, receiveAsTokens []bool) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "completeQueuedWithdrawals", withdrawals, tokens, middlewareTimesIndexes, receiveAsTokens)
}

// CompleteQueuedWithdrawals is a paid mutator transaction binding the contract method 0x33404396.
//
// Solidity: function completeQueuedWithdrawals((address,address,address,uint256,uint32,address[],uint256[])[] withdrawals, address[][] tokens, uint256[] middlewareTimesIndexes, bool[] receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) CompleteQueuedWithdrawals(withdrawals []IDelegationManagerWithdrawal, tokens [][]common.Address, middlewareTimesIndexes []*big.Int, receiveAsTokens []bool) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.CompleteQueuedWithdrawals(&_ContractDelegationManager.TransactOpts, withdrawals, tokens, middlewareTimesIndexes, receiveAsTokens)
}

// CompleteQueuedWithdrawals is a paid mutator transaction binding the contract method 0x33404396.
//
// Solidity: function completeQueuedWithdrawals((address,address,address,uint256,uint32,address[],uint256[])[] withdrawals, address[][] tokens, uint256[] middlewareTimesIndexes, bool[] receiveAsTokens) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) CompleteQueuedWithdrawals(withdrawals []IDelegationManagerWithdrawal, tokens [][]common.Address, middlewareTimesIndexes []*big.Int, receiveAsTokens []bool) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.CompleteQueuedWithdrawals(&_ContractDelegationManager.TransactOpts, withdrawals, tokens, middlewareTimesIndexes, receiveAsTokens)
}

// DecreaseDelegatedShares is a paid mutator transaction binding the contract method 0x132d4967.
//
// Solidity: function decreaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) DecreaseDelegatedShares(opts *bind.TransactOpts, staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "decreaseDelegatedShares", staker, strategy, shares)
}

// DecreaseDelegatedShares is a paid mutator transaction binding the contract method 0x132d4967.
//
// Solidity: function decreaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) DecreaseDelegatedShares(staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DecreaseDelegatedShares(&_ContractDelegationManager.TransactOpts, staker, strategy, shares)
}

// DecreaseDelegatedShares is a paid mutator transaction binding the contract method 0x132d4967.
//
// Solidity: function decreaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) DecreaseDelegatedShares(staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DecreaseDelegatedShares(&_ContractDelegationManager.TransactOpts, staker, strategy, shares)
}

// DelegateTo is a paid mutator transaction binding the contract method 0xeea9064b.
//
// Solidity: function delegateTo(address operator, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) DelegateTo(opts *bind.TransactOpts, operator common.Address, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "delegateTo", operator, approverSignatureAndExpiry, approverSalt)
}

// DelegateTo is a paid mutator transaction binding the contract method 0xeea9064b.
//
// Solidity: function delegateTo(address operator, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) DelegateTo(operator common.Address, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DelegateTo(&_ContractDelegationManager.TransactOpts, operator, approverSignatureAndExpiry, approverSalt)
}

// DelegateTo is a paid mutator transaction binding the contract method 0xeea9064b.
//
// Solidity: function delegateTo(address operator, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) DelegateTo(operator common.Address, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DelegateTo(&_ContractDelegationManager.TransactOpts, operator, approverSignatureAndExpiry, approverSalt)
}

// DelegateToBySignature is a paid mutator transaction binding the contract method 0x7f548071.
//
// Solidity: function delegateToBySignature(address staker, address operator, (bytes,uint256) stakerSignatureAndExpiry, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) DelegateToBySignature(opts *bind.TransactOpts, staker common.Address, operator common.Address, stakerSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "delegateToBySignature", staker, operator, stakerSignatureAndExpiry, approverSignatureAndExpiry, approverSalt)
}

// DelegateToBySignature is a paid mutator transaction binding the contract method 0x7f548071.
//
// Solidity: function delegateToBySignature(address staker, address operator, (bytes,uint256) stakerSignatureAndExpiry, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) DelegateToBySignature(staker common.Address, operator common.Address, stakerSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DelegateToBySignature(&_ContractDelegationManager.TransactOpts, staker, operator, stakerSignatureAndExpiry, approverSignatureAndExpiry, approverSalt)
}

// DelegateToBySignature is a paid mutator transaction binding the contract method 0x7f548071.
//
// Solidity: function delegateToBySignature(address staker, address operator, (bytes,uint256) stakerSignatureAndExpiry, (bytes,uint256) approverSignatureAndExpiry, bytes32 approverSalt) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) DelegateToBySignature(staker common.Address, operator common.Address, stakerSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSignatureAndExpiry ISignatureUtilsSignatureWithExpiry, approverSalt [32]byte) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.DelegateToBySignature(&_ContractDelegationManager.TransactOpts, staker, operator, stakerSignatureAndExpiry, approverSignatureAndExpiry, approverSalt)
}

// IncreaseDelegatedShares is a paid mutator transaction binding the contract method 0x28a573ae.
//
// Solidity: function increaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) IncreaseDelegatedShares(opts *bind.TransactOpts, staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "increaseDelegatedShares", staker, strategy, shares)
}

// IncreaseDelegatedShares is a paid mutator transaction binding the contract method 0x28a573ae.
//
// Solidity: function increaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) IncreaseDelegatedShares(staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.IncreaseDelegatedShares(&_ContractDelegationManager.TransactOpts, staker, strategy, shares)
}

// IncreaseDelegatedShares is a paid mutator transaction binding the contract method 0x28a573ae.
//
// Solidity: function increaseDelegatedShares(address staker, address strategy, uint256 shares) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) IncreaseDelegatedShares(staker common.Address, strategy common.Address, shares *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.IncreaseDelegatedShares(&_ContractDelegationManager.TransactOpts, staker, strategy, shares)
}

// Initialize is a paid mutator transaction binding the contract method 0x22bf40e4.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 initialPausedStatus, uint256 _minWithdrawalDelayBlocks, address[] _strategies, uint256[] _withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _pauserRegistry common.Address, initialPausedStatus *big.Int, _minWithdrawalDelayBlocks *big.Int, _strategies []common.Address, _withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "initialize", initialOwner, _pauserRegistry, initialPausedStatus, _minWithdrawalDelayBlocks, _strategies, _withdrawalDelayBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0x22bf40e4.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 initialPausedStatus, uint256 _minWithdrawalDelayBlocks, address[] _strategies, uint256[] _withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) Initialize(initialOwner common.Address, _pauserRegistry common.Address, initialPausedStatus *big.Int, _minWithdrawalDelayBlocks *big.Int, _strategies []common.Address, _withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Initialize(&_ContractDelegationManager.TransactOpts, initialOwner, _pauserRegistry, initialPausedStatus, _minWithdrawalDelayBlocks, _strategies, _withdrawalDelayBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0x22bf40e4.
//
// Solidity: function initialize(address initialOwner, address _pauserRegistry, uint256 initialPausedStatus, uint256 _minWithdrawalDelayBlocks, address[] _strategies, uint256[] _withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) Initialize(initialOwner common.Address, _pauserRegistry common.Address, initialPausedStatus *big.Int, _minWithdrawalDelayBlocks *big.Int, _strategies []common.Address, _withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Initialize(&_ContractDelegationManager.TransactOpts, initialOwner, _pauserRegistry, initialPausedStatus, _minWithdrawalDelayBlocks, _strategies, _withdrawalDelayBlocks)
}

// MigrateQueuedWithdrawals is a paid mutator transaction binding the contract method 0x5cfe8d2c.
//
// Solidity: function migrateQueuedWithdrawals((address[],uint256[],address,(address,uint96),uint32,address)[] withdrawalsToMigrate) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) MigrateQueuedWithdrawals(opts *bind.TransactOpts, withdrawalsToMigrate []IStrategyManagerDeprecatedStructQueuedWithdrawal) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "migrateQueuedWithdrawals", withdrawalsToMigrate)
}

// MigrateQueuedWithdrawals is a paid mutator transaction binding the contract method 0x5cfe8d2c.
//
// Solidity: function migrateQueuedWithdrawals((address[],uint256[],address,(address,uint96),uint32,address)[] withdrawalsToMigrate) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) MigrateQueuedWithdrawals(withdrawalsToMigrate []IStrategyManagerDeprecatedStructQueuedWithdrawal) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.MigrateQueuedWithdrawals(&_ContractDelegationManager.TransactOpts, withdrawalsToMigrate)
}

// MigrateQueuedWithdrawals is a paid mutator transaction binding the contract method 0x5cfe8d2c.
//
// Solidity: function migrateQueuedWithdrawals((address[],uint256[],address,(address,uint96),uint32,address)[] withdrawalsToMigrate) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) MigrateQueuedWithdrawals(withdrawalsToMigrate []IStrategyManagerDeprecatedStructQueuedWithdrawal) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.MigrateQueuedWithdrawals(&_ContractDelegationManager.TransactOpts, withdrawalsToMigrate)
}

// ModifyOperatorDetails is a paid mutator transaction binding the contract method 0xf16172b0.
//
// Solidity: function modifyOperatorDetails((address,address,uint32) newOperatorDetails) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) ModifyOperatorDetails(opts *bind.TransactOpts, newOperatorDetails IDelegationManagerOperatorDetails) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "modifyOperatorDetails", newOperatorDetails)
}

// ModifyOperatorDetails is a paid mutator transaction binding the contract method 0xf16172b0.
//
// Solidity: function modifyOperatorDetails((address,address,uint32) newOperatorDetails) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) ModifyOperatorDetails(newOperatorDetails IDelegationManagerOperatorDetails) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.ModifyOperatorDetails(&_ContractDelegationManager.TransactOpts, newOperatorDetails)
}

// ModifyOperatorDetails is a paid mutator transaction binding the contract method 0xf16172b0.
//
// Solidity: function modifyOperatorDetails((address,address,uint32) newOperatorDetails) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) ModifyOperatorDetails(newOperatorDetails IDelegationManagerOperatorDetails) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.ModifyOperatorDetails(&_ContractDelegationManager.TransactOpts, newOperatorDetails)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) Pause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "pause", newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Pause(&_ContractDelegationManager.TransactOpts, newPausedStatus)
}

// Pause is a paid mutator transaction binding the contract method 0x136439dd.
//
// Solidity: function pause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) Pause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Pause(&_ContractDelegationManager.TransactOpts, newPausedStatus)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) PauseAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "pauseAll")
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) PauseAll() (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.PauseAll(&_ContractDelegationManager.TransactOpts)
}

// PauseAll is a paid mutator transaction binding the contract method 0x595c6a67.
//
// Solidity: function pauseAll() returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) PauseAll() (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.PauseAll(&_ContractDelegationManager.TransactOpts)
}

// QueueWithdrawals is a paid mutator transaction binding the contract method 0x0dd8dd02.
//
// Solidity: function queueWithdrawals((address[],uint256[],address)[] queuedWithdrawalParams) returns(bytes32[])
func (_ContractDelegationManager *ContractDelegationManagerTransactor) QueueWithdrawals(opts *bind.TransactOpts, queuedWithdrawalParams []IDelegationManagerQueuedWithdrawalParams) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "queueWithdrawals", queuedWithdrawalParams)
}

// QueueWithdrawals is a paid mutator transaction binding the contract method 0x0dd8dd02.
//
// Solidity: function queueWithdrawals((address[],uint256[],address)[] queuedWithdrawalParams) returns(bytes32[])
func (_ContractDelegationManager *ContractDelegationManagerSession) QueueWithdrawals(queuedWithdrawalParams []IDelegationManagerQueuedWithdrawalParams) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.QueueWithdrawals(&_ContractDelegationManager.TransactOpts, queuedWithdrawalParams)
}

// QueueWithdrawals is a paid mutator transaction binding the contract method 0x0dd8dd02.
//
// Solidity: function queueWithdrawals((address[],uint256[],address)[] queuedWithdrawalParams) returns(bytes32[])
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) QueueWithdrawals(queuedWithdrawalParams []IDelegationManagerQueuedWithdrawalParams) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.QueueWithdrawals(&_ContractDelegationManager.TransactOpts, queuedWithdrawalParams)
}

// RegisterAsOperator is a paid mutator transaction binding the contract method 0x0f589e59.
//
// Solidity: function registerAsOperator((address,address,uint32) registeringOperatorDetails, string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) RegisterAsOperator(opts *bind.TransactOpts, registeringOperatorDetails IDelegationManagerOperatorDetails, metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "registerAsOperator", registeringOperatorDetails, metadataURI)
}

// RegisterAsOperator is a paid mutator transaction binding the contract method 0x0f589e59.
//
// Solidity: function registerAsOperator((address,address,uint32) registeringOperatorDetails, string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) RegisterAsOperator(registeringOperatorDetails IDelegationManagerOperatorDetails, metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.RegisterAsOperator(&_ContractDelegationManager.TransactOpts, registeringOperatorDetails, metadataURI)
}

// RegisterAsOperator is a paid mutator transaction binding the contract method 0x0f589e59.
//
// Solidity: function registerAsOperator((address,address,uint32) registeringOperatorDetails, string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) RegisterAsOperator(registeringOperatorDetails IDelegationManagerOperatorDetails, metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.RegisterAsOperator(&_ContractDelegationManager.TransactOpts, registeringOperatorDetails, metadataURI)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.RenounceOwnership(&_ContractDelegationManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.RenounceOwnership(&_ContractDelegationManager.TransactOpts)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) SetPauserRegistry(opts *bind.TransactOpts, newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "setPauserRegistry", newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetPauserRegistry(&_ContractDelegationManager.TransactOpts, newPauserRegistry)
}

// SetPauserRegistry is a paid mutator transaction binding the contract method 0x10d67a2f.
//
// Solidity: function setPauserRegistry(address newPauserRegistry) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) SetPauserRegistry(newPauserRegistry common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetPauserRegistry(&_ContractDelegationManager.TransactOpts, newPauserRegistry)
}

// SetStrategyWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x1522bf02.
//
// Solidity: function setStrategyWithdrawalDelayBlocks(address[] strategies, uint256[] withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) SetStrategyWithdrawalDelayBlocks(opts *bind.TransactOpts, strategies []common.Address, withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "setStrategyWithdrawalDelayBlocks", strategies, withdrawalDelayBlocks)
}

// SetStrategyWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x1522bf02.
//
// Solidity: function setStrategyWithdrawalDelayBlocks(address[] strategies, uint256[] withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) SetStrategyWithdrawalDelayBlocks(strategies []common.Address, withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetStrategyWithdrawalDelayBlocks(&_ContractDelegationManager.TransactOpts, strategies, withdrawalDelayBlocks)
}

// SetStrategyWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x1522bf02.
//
// Solidity: function setStrategyWithdrawalDelayBlocks(address[] strategies, uint256[] withdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) SetStrategyWithdrawalDelayBlocks(strategies []common.Address, withdrawalDelayBlocks []*big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetStrategyWithdrawalDelayBlocks(&_ContractDelegationManager.TransactOpts, strategies, withdrawalDelayBlocks)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.TransferOwnership(&_ContractDelegationManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.TransferOwnership(&_ContractDelegationManager.TransactOpts, newOwner)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address staker) returns(bytes32[] withdrawalRoots)
func (_ContractDelegationManager *ContractDelegationManagerTransactor) Undelegate(opts *bind.TransactOpts, staker common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "undelegate", staker)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address staker) returns(bytes32[] withdrawalRoots)
func (_ContractDelegationManager *ContractDelegationManagerSession) Undelegate(staker common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Undelegate(&_ContractDelegationManager.TransactOpts, staker)
}

// Undelegate is a paid mutator transaction binding the contract method 0xda8be864.
//
// Solidity: function undelegate(address staker) returns(bytes32[] withdrawalRoots)
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) Undelegate(staker common.Address) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Undelegate(&_ContractDelegationManager.TransactOpts, staker)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) Unpause(opts *bind.TransactOpts, newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "unpause", newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Unpause(&_ContractDelegationManager.TransactOpts, newPausedStatus)
}

// Unpause is a paid mutator transaction binding the contract method 0xfabc1cbc.
//
// Solidity: function unpause(uint256 newPausedStatus) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) Unpause(newPausedStatus *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.Unpause(&_ContractDelegationManager.TransactOpts, newPausedStatus)
}

// UpdateOperatorMetadataURI is a paid mutator transaction binding the contract method 0x99be81c8.
//
// Solidity: function updateOperatorMetadataURI(string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) UpdateOperatorMetadataURI(opts *bind.TransactOpts, metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "updateOperatorMetadataURI", metadataURI)
}

// UpdateOperatorMetadataURI is a paid mutator transaction binding the contract method 0x99be81c8.
//
// Solidity: function updateOperatorMetadataURI(string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) UpdateOperatorMetadataURI(metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.UpdateOperatorMetadataURI(&_ContractDelegationManager.TransactOpts, metadataURI)
}

// UpdateOperatorMetadataURI is a paid mutator transaction binding the contract method 0x99be81c8.
//
// Solidity: function updateOperatorMetadataURI(string metadataURI) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) UpdateOperatorMetadataURI(metadataURI string) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.UpdateOperatorMetadataURI(&_ContractDelegationManager.TransactOpts, metadataURI)
}

// ContractDelegationManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContractDelegationManager contract.
type ContractDelegationManagerInitializedIterator struct {
	Event *ContractDelegationManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerInitialized)
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
		it.Event = new(ContractDelegationManagerInitialized)
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
func (it *ContractDelegationManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerInitialized represents a Initialized event raised by the ContractDelegationManager contract.
type ContractDelegationManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractDelegationManagerInitializedIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerInitializedIterator{contract: _ContractDelegationManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerInitialized)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseInitialized(log types.Log) (*ContractDelegationManagerInitialized, error) {
	event := new(ContractDelegationManagerInitialized)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator is returned from FilterMinWithdrawalDelayBlocksSet and is used to iterate over the raw logs and unpacked data for MinWithdrawalDelayBlocksSet events raised by the ContractDelegationManager contract.
type ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator struct {
	Event *ContractDelegationManagerMinWithdrawalDelayBlocksSet // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerMinWithdrawalDelayBlocksSet)
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
		it.Event = new(ContractDelegationManagerMinWithdrawalDelayBlocksSet)
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
func (it *ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerMinWithdrawalDelayBlocksSet represents a MinWithdrawalDelayBlocksSet event raised by the ContractDelegationManager contract.
type ContractDelegationManagerMinWithdrawalDelayBlocksSet struct {
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMinWithdrawalDelayBlocksSet is a free log retrieval operation binding the contract event 0xafa003cd76f87ff9d62b35beea889920f33c0c42b8d45b74954d61d50f4b6b69.
//
// Solidity: event MinWithdrawalDelayBlocksSet(uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterMinWithdrawalDelayBlocksSet(opts *bind.FilterOpts) (*ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "MinWithdrawalDelayBlocksSet")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerMinWithdrawalDelayBlocksSetIterator{contract: _ContractDelegationManager.contract, event: "MinWithdrawalDelayBlocksSet", logs: logs, sub: sub}, nil
}

// WatchMinWithdrawalDelayBlocksSet is a free log subscription operation binding the contract event 0xafa003cd76f87ff9d62b35beea889920f33c0c42b8d45b74954d61d50f4b6b69.
//
// Solidity: event MinWithdrawalDelayBlocksSet(uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchMinWithdrawalDelayBlocksSet(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerMinWithdrawalDelayBlocksSet) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "MinWithdrawalDelayBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerMinWithdrawalDelayBlocksSet)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "MinWithdrawalDelayBlocksSet", log); err != nil {
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

// ParseMinWithdrawalDelayBlocksSet is a log parse operation binding the contract event 0xafa003cd76f87ff9d62b35beea889920f33c0c42b8d45b74954d61d50f4b6b69.
//
// Solidity: event MinWithdrawalDelayBlocksSet(uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseMinWithdrawalDelayBlocksSet(log types.Log) (*ContractDelegationManagerMinWithdrawalDelayBlocksSet, error) {
	event := new(ContractDelegationManagerMinWithdrawalDelayBlocksSet)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "MinWithdrawalDelayBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOperatorDetailsModifiedIterator is returned from FilterOperatorDetailsModified and is used to iterate over the raw logs and unpacked data for OperatorDetailsModified events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorDetailsModifiedIterator struct {
	Event *ContractDelegationManagerOperatorDetailsModified // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOperatorDetailsModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOperatorDetailsModified)
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
		it.Event = new(ContractDelegationManagerOperatorDetailsModified)
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
func (it *ContractDelegationManagerOperatorDetailsModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOperatorDetailsModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOperatorDetailsModified represents a OperatorDetailsModified event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorDetailsModified struct {
	Operator           common.Address
	NewOperatorDetails IDelegationManagerOperatorDetails
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOperatorDetailsModified is a free log retrieval operation binding the contract event 0xfebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac.
//
// Solidity: event OperatorDetailsModified(address indexed operator, (address,address,uint32) newOperatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOperatorDetailsModified(opts *bind.FilterOpts, operator []common.Address) (*ContractDelegationManagerOperatorDetailsModifiedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OperatorDetailsModified", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOperatorDetailsModifiedIterator{contract: _ContractDelegationManager.contract, event: "OperatorDetailsModified", logs: logs, sub: sub}, nil
}

// WatchOperatorDetailsModified is a free log subscription operation binding the contract event 0xfebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac.
//
// Solidity: event OperatorDetailsModified(address indexed operator, (address,address,uint32) newOperatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOperatorDetailsModified(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOperatorDetailsModified, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OperatorDetailsModified", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOperatorDetailsModified)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorDetailsModified", log); err != nil {
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

// ParseOperatorDetailsModified is a log parse operation binding the contract event 0xfebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac.
//
// Solidity: event OperatorDetailsModified(address indexed operator, (address,address,uint32) newOperatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOperatorDetailsModified(log types.Log) (*ContractDelegationManagerOperatorDetailsModified, error) {
	event := new(ContractDelegationManagerOperatorDetailsModified)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorDetailsModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOperatorMetadataURIUpdatedIterator is returned from FilterOperatorMetadataURIUpdated and is used to iterate over the raw logs and unpacked data for OperatorMetadataURIUpdated events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorMetadataURIUpdatedIterator struct {
	Event *ContractDelegationManagerOperatorMetadataURIUpdated // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOperatorMetadataURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOperatorMetadataURIUpdated)
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
		it.Event = new(ContractDelegationManagerOperatorMetadataURIUpdated)
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
func (it *ContractDelegationManagerOperatorMetadataURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOperatorMetadataURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOperatorMetadataURIUpdated represents a OperatorMetadataURIUpdated event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorMetadataURIUpdated struct {
	Operator    common.Address
	MetadataURI string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorMetadataURIUpdated is a free log retrieval operation binding the contract event 0x02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b6708090.
//
// Solidity: event OperatorMetadataURIUpdated(address indexed operator, string metadataURI)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOperatorMetadataURIUpdated(opts *bind.FilterOpts, operator []common.Address) (*ContractDelegationManagerOperatorMetadataURIUpdatedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OperatorMetadataURIUpdated", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOperatorMetadataURIUpdatedIterator{contract: _ContractDelegationManager.contract, event: "OperatorMetadataURIUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorMetadataURIUpdated is a free log subscription operation binding the contract event 0x02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b6708090.
//
// Solidity: event OperatorMetadataURIUpdated(address indexed operator, string metadataURI)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOperatorMetadataURIUpdated(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOperatorMetadataURIUpdated, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OperatorMetadataURIUpdated", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOperatorMetadataURIUpdated)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorMetadataURIUpdated", log); err != nil {
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

// ParseOperatorMetadataURIUpdated is a log parse operation binding the contract event 0x02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b6708090.
//
// Solidity: event OperatorMetadataURIUpdated(address indexed operator, string metadataURI)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOperatorMetadataURIUpdated(log types.Log) (*ContractDelegationManagerOperatorMetadataURIUpdated, error) {
	event := new(ContractDelegationManagerOperatorMetadataURIUpdated)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorMetadataURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorRegisteredIterator struct {
	Event *ContractDelegationManagerOperatorRegistered // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOperatorRegistered)
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
		it.Event = new(ContractDelegationManagerOperatorRegistered)
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
func (it *ContractDelegationManagerOperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOperatorRegistered represents a OperatorRegistered event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorRegistered struct {
	Operator        common.Address
	OperatorDetails IDelegationManagerOperatorDetails
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0x8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e2.
//
// Solidity: event OperatorRegistered(address indexed operator, (address,address,uint32) operatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address) (*ContractDelegationManagerOperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOperatorRegisteredIterator{contract: _ContractDelegationManager.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0x8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e2.
//
// Solidity: event OperatorRegistered(address indexed operator, (address,address,uint32) operatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOperatorRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOperatorRegistered)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
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

// ParseOperatorRegistered is a log parse operation binding the contract event 0x8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e2.
//
// Solidity: event OperatorRegistered(address indexed operator, (address,address,uint32) operatorDetails)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOperatorRegistered(log types.Log) (*ContractDelegationManagerOperatorRegistered, error) {
	event := new(ContractDelegationManagerOperatorRegistered)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOperatorSharesDecreasedIterator is returned from FilterOperatorSharesDecreased and is used to iterate over the raw logs and unpacked data for OperatorSharesDecreased events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorSharesDecreasedIterator struct {
	Event *ContractDelegationManagerOperatorSharesDecreased // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOperatorSharesDecreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOperatorSharesDecreased)
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
		it.Event = new(ContractDelegationManagerOperatorSharesDecreased)
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
func (it *ContractDelegationManagerOperatorSharesDecreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOperatorSharesDecreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOperatorSharesDecreased represents a OperatorSharesDecreased event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorSharesDecreased struct {
	Operator common.Address
	Staker   common.Address
	Strategy common.Address
	Shares   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorSharesDecreased is a free log retrieval operation binding the contract event 0x6909600037b75d7b4733aedd815442b5ec018a827751c832aaff64eba5d6d2dd.
//
// Solidity: event OperatorSharesDecreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOperatorSharesDecreased(opts *bind.FilterOpts, operator []common.Address) (*ContractDelegationManagerOperatorSharesDecreasedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OperatorSharesDecreased", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOperatorSharesDecreasedIterator{contract: _ContractDelegationManager.contract, event: "OperatorSharesDecreased", logs: logs, sub: sub}, nil
}

// WatchOperatorSharesDecreased is a free log subscription operation binding the contract event 0x6909600037b75d7b4733aedd815442b5ec018a827751c832aaff64eba5d6d2dd.
//
// Solidity: event OperatorSharesDecreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOperatorSharesDecreased(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOperatorSharesDecreased, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OperatorSharesDecreased", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOperatorSharesDecreased)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorSharesDecreased", log); err != nil {
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

// ParseOperatorSharesDecreased is a log parse operation binding the contract event 0x6909600037b75d7b4733aedd815442b5ec018a827751c832aaff64eba5d6d2dd.
//
// Solidity: event OperatorSharesDecreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOperatorSharesDecreased(log types.Log) (*ContractDelegationManagerOperatorSharesDecreased, error) {
	event := new(ContractDelegationManagerOperatorSharesDecreased)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorSharesDecreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOperatorSharesIncreasedIterator is returned from FilterOperatorSharesIncreased and is used to iterate over the raw logs and unpacked data for OperatorSharesIncreased events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorSharesIncreasedIterator struct {
	Event *ContractDelegationManagerOperatorSharesIncreased // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOperatorSharesIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOperatorSharesIncreased)
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
		it.Event = new(ContractDelegationManagerOperatorSharesIncreased)
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
func (it *ContractDelegationManagerOperatorSharesIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOperatorSharesIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOperatorSharesIncreased represents a OperatorSharesIncreased event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOperatorSharesIncreased struct {
	Operator common.Address
	Staker   common.Address
	Strategy common.Address
	Shares   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorSharesIncreased is a free log retrieval operation binding the contract event 0x1ec042c965e2edd7107b51188ee0f383e22e76179041ab3a9d18ff151405166c.
//
// Solidity: event OperatorSharesIncreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOperatorSharesIncreased(opts *bind.FilterOpts, operator []common.Address) (*ContractDelegationManagerOperatorSharesIncreasedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OperatorSharesIncreased", operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOperatorSharesIncreasedIterator{contract: _ContractDelegationManager.contract, event: "OperatorSharesIncreased", logs: logs, sub: sub}, nil
}

// WatchOperatorSharesIncreased is a free log subscription operation binding the contract event 0x1ec042c965e2edd7107b51188ee0f383e22e76179041ab3a9d18ff151405166c.
//
// Solidity: event OperatorSharesIncreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOperatorSharesIncreased(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOperatorSharesIncreased, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OperatorSharesIncreased", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOperatorSharesIncreased)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorSharesIncreased", log); err != nil {
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

// ParseOperatorSharesIncreased is a log parse operation binding the contract event 0x1ec042c965e2edd7107b51188ee0f383e22e76179041ab3a9d18ff151405166c.
//
// Solidity: event OperatorSharesIncreased(address indexed operator, address staker, address strategy, uint256 shares)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOperatorSharesIncreased(log types.Log) (*ContractDelegationManagerOperatorSharesIncreased, error) {
	event := new(ContractDelegationManagerOperatorSharesIncreased)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OperatorSharesIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ContractDelegationManager contract.
type ContractDelegationManagerOwnershipTransferredIterator struct {
	Event *ContractDelegationManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerOwnershipTransferred)
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
		it.Event = new(ContractDelegationManagerOwnershipTransferred)
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
func (it *ContractDelegationManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ContractDelegationManager contract.
type ContractDelegationManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractDelegationManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerOwnershipTransferredIterator{contract: _ContractDelegationManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerOwnershipTransferred)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ContractDelegationManagerOwnershipTransferred, error) {
	event := new(ContractDelegationManagerOwnershipTransferred)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ContractDelegationManager contract.
type ContractDelegationManagerPausedIterator struct {
	Event *ContractDelegationManagerPaused // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerPaused)
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
		it.Event = new(ContractDelegationManagerPaused)
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
func (it *ContractDelegationManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerPaused represents a Paused event raised by the ContractDelegationManager contract.
type ContractDelegationManagerPaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*ContractDelegationManagerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerPausedIterator{contract: _ContractDelegationManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d.
//
// Solidity: event Paused(address indexed account, uint256 newPausedStatus)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerPaused)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParsePaused(log types.Log) (*ContractDelegationManagerPaused, error) {
	event := new(ContractDelegationManagerPaused)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerPauserRegistrySetIterator is returned from FilterPauserRegistrySet and is used to iterate over the raw logs and unpacked data for PauserRegistrySet events raised by the ContractDelegationManager contract.
type ContractDelegationManagerPauserRegistrySetIterator struct {
	Event *ContractDelegationManagerPauserRegistrySet // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerPauserRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerPauserRegistrySet)
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
		it.Event = new(ContractDelegationManagerPauserRegistrySet)
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
func (it *ContractDelegationManagerPauserRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerPauserRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerPauserRegistrySet represents a PauserRegistrySet event raised by the ContractDelegationManager contract.
type ContractDelegationManagerPauserRegistrySet struct {
	PauserRegistry    common.Address
	NewPauserRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPauserRegistrySet is a free log retrieval operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterPauserRegistrySet(opts *bind.FilterOpts) (*ContractDelegationManagerPauserRegistrySetIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerPauserRegistrySetIterator{contract: _ContractDelegationManager.contract, event: "PauserRegistrySet", logs: logs, sub: sub}, nil
}

// WatchPauserRegistrySet is a free log subscription operation binding the contract event 0x6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6.
//
// Solidity: event PauserRegistrySet(address pauserRegistry, address newPauserRegistry)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchPauserRegistrySet(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerPauserRegistrySet) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "PauserRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerPauserRegistrySet)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
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
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParsePauserRegistrySet(log types.Log) (*ContractDelegationManagerPauserRegistrySet, error) {
	event := new(ContractDelegationManagerPauserRegistrySet)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "PauserRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerStakerDelegatedIterator is returned from FilterStakerDelegated and is used to iterate over the raw logs and unpacked data for StakerDelegated events raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerDelegatedIterator struct {
	Event *ContractDelegationManagerStakerDelegated // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerStakerDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerStakerDelegated)
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
		it.Event = new(ContractDelegationManagerStakerDelegated)
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
func (it *ContractDelegationManagerStakerDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerStakerDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerStakerDelegated represents a StakerDelegated event raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerDelegated struct {
	Staker   common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakerDelegated is a free log retrieval operation binding the contract event 0xc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d8743304.
//
// Solidity: event StakerDelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterStakerDelegated(opts *bind.FilterOpts, staker []common.Address, operator []common.Address) (*ContractDelegationManagerStakerDelegatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "StakerDelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerStakerDelegatedIterator{contract: _ContractDelegationManager.contract, event: "StakerDelegated", logs: logs, sub: sub}, nil
}

// WatchStakerDelegated is a free log subscription operation binding the contract event 0xc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d8743304.
//
// Solidity: event StakerDelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchStakerDelegated(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerStakerDelegated, staker []common.Address, operator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "StakerDelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerStakerDelegated)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerDelegated", log); err != nil {
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

// ParseStakerDelegated is a log parse operation binding the contract event 0xc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d8743304.
//
// Solidity: event StakerDelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseStakerDelegated(log types.Log) (*ContractDelegationManagerStakerDelegated, error) {
	event := new(ContractDelegationManagerStakerDelegated)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerDelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerStakerForceUndelegatedIterator is returned from FilterStakerForceUndelegated and is used to iterate over the raw logs and unpacked data for StakerForceUndelegated events raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerForceUndelegatedIterator struct {
	Event *ContractDelegationManagerStakerForceUndelegated // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerStakerForceUndelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerStakerForceUndelegated)
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
		it.Event = new(ContractDelegationManagerStakerForceUndelegated)
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
func (it *ContractDelegationManagerStakerForceUndelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerStakerForceUndelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerStakerForceUndelegated represents a StakerForceUndelegated event raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerForceUndelegated struct {
	Staker   common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakerForceUndelegated is a free log retrieval operation binding the contract event 0xf0eddf07e6ea14f388b47e1e94a0f464ecbd9eed4171130e0fc0e99fb4030a8a.
//
// Solidity: event StakerForceUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterStakerForceUndelegated(opts *bind.FilterOpts, staker []common.Address, operator []common.Address) (*ContractDelegationManagerStakerForceUndelegatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "StakerForceUndelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerStakerForceUndelegatedIterator{contract: _ContractDelegationManager.contract, event: "StakerForceUndelegated", logs: logs, sub: sub}, nil
}

// WatchStakerForceUndelegated is a free log subscription operation binding the contract event 0xf0eddf07e6ea14f388b47e1e94a0f464ecbd9eed4171130e0fc0e99fb4030a8a.
//
// Solidity: event StakerForceUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchStakerForceUndelegated(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerStakerForceUndelegated, staker []common.Address, operator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "StakerForceUndelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerStakerForceUndelegated)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerForceUndelegated", log); err != nil {
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

// ParseStakerForceUndelegated is a log parse operation binding the contract event 0xf0eddf07e6ea14f388b47e1e94a0f464ecbd9eed4171130e0fc0e99fb4030a8a.
//
// Solidity: event StakerForceUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseStakerForceUndelegated(log types.Log) (*ContractDelegationManagerStakerForceUndelegated, error) {
	event := new(ContractDelegationManagerStakerForceUndelegated)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerForceUndelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerStakerUndelegatedIterator is returned from FilterStakerUndelegated and is used to iterate over the raw logs and unpacked data for StakerUndelegated events raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerUndelegatedIterator struct {
	Event *ContractDelegationManagerStakerUndelegated // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerStakerUndelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerStakerUndelegated)
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
		it.Event = new(ContractDelegationManagerStakerUndelegated)
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
func (it *ContractDelegationManagerStakerUndelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerStakerUndelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerStakerUndelegated represents a StakerUndelegated event raised by the ContractDelegationManager contract.
type ContractDelegationManagerStakerUndelegated struct {
	Staker   common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakerUndelegated is a free log retrieval operation binding the contract event 0xfee30966a256b71e14bc0ebfc94315e28ef4a97a7131a9e2b7a310a73af44676.
//
// Solidity: event StakerUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterStakerUndelegated(opts *bind.FilterOpts, staker []common.Address, operator []common.Address) (*ContractDelegationManagerStakerUndelegatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "StakerUndelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerStakerUndelegatedIterator{contract: _ContractDelegationManager.contract, event: "StakerUndelegated", logs: logs, sub: sub}, nil
}

// WatchStakerUndelegated is a free log subscription operation binding the contract event 0xfee30966a256b71e14bc0ebfc94315e28ef4a97a7131a9e2b7a310a73af44676.
//
// Solidity: event StakerUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchStakerUndelegated(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerStakerUndelegated, staker []common.Address, operator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "StakerUndelegated", stakerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerStakerUndelegated)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerUndelegated", log); err != nil {
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

// ParseStakerUndelegated is a log parse operation binding the contract event 0xfee30966a256b71e14bc0ebfc94315e28ef4a97a7131a9e2b7a310a73af44676.
//
// Solidity: event StakerUndelegated(address indexed staker, address indexed operator)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseStakerUndelegated(log types.Log) (*ContractDelegationManagerStakerUndelegated, error) {
	event := new(ContractDelegationManagerStakerUndelegated)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "StakerUndelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator is returned from FilterStrategyWithdrawalDelayBlocksSet and is used to iterate over the raw logs and unpacked data for StrategyWithdrawalDelayBlocksSet events raised by the ContractDelegationManager contract.
type ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator struct {
	Event *ContractDelegationManagerStrategyWithdrawalDelayBlocksSet // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerStrategyWithdrawalDelayBlocksSet)
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
		it.Event = new(ContractDelegationManagerStrategyWithdrawalDelayBlocksSet)
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
func (it *ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerStrategyWithdrawalDelayBlocksSet represents a StrategyWithdrawalDelayBlocksSet event raised by the ContractDelegationManager contract.
type ContractDelegationManagerStrategyWithdrawalDelayBlocksSet struct {
	Strategy      common.Address
	PreviousValue *big.Int
	NewValue      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStrategyWithdrawalDelayBlocksSet is a free log retrieval operation binding the contract event 0x0e7efa738e8b0ce6376a0c1af471655540d2e9a81647d7b09ed823018426576d.
//
// Solidity: event StrategyWithdrawalDelayBlocksSet(address strategy, uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterStrategyWithdrawalDelayBlocksSet(opts *bind.FilterOpts) (*ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "StrategyWithdrawalDelayBlocksSet")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerStrategyWithdrawalDelayBlocksSetIterator{contract: _ContractDelegationManager.contract, event: "StrategyWithdrawalDelayBlocksSet", logs: logs, sub: sub}, nil
}

// WatchStrategyWithdrawalDelayBlocksSet is a free log subscription operation binding the contract event 0x0e7efa738e8b0ce6376a0c1af471655540d2e9a81647d7b09ed823018426576d.
//
// Solidity: event StrategyWithdrawalDelayBlocksSet(address strategy, uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchStrategyWithdrawalDelayBlocksSet(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerStrategyWithdrawalDelayBlocksSet) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "StrategyWithdrawalDelayBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerStrategyWithdrawalDelayBlocksSet)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "StrategyWithdrawalDelayBlocksSet", log); err != nil {
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

// ParseStrategyWithdrawalDelayBlocksSet is a log parse operation binding the contract event 0x0e7efa738e8b0ce6376a0c1af471655540d2e9a81647d7b09ed823018426576d.
//
// Solidity: event StrategyWithdrawalDelayBlocksSet(address strategy, uint256 previousValue, uint256 newValue)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseStrategyWithdrawalDelayBlocksSet(log types.Log) (*ContractDelegationManagerStrategyWithdrawalDelayBlocksSet, error) {
	event := new(ContractDelegationManagerStrategyWithdrawalDelayBlocksSet)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "StrategyWithdrawalDelayBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ContractDelegationManager contract.
type ContractDelegationManagerUnpausedIterator struct {
	Event *ContractDelegationManagerUnpaused // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerUnpaused)
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
		it.Event = new(ContractDelegationManagerUnpaused)
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
func (it *ContractDelegationManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerUnpaused represents a Unpaused event raised by the ContractDelegationManager contract.
type ContractDelegationManagerUnpaused struct {
	Account         common.Address
	NewPausedStatus *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*ContractDelegationManagerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerUnpausedIterator{contract: _ContractDelegationManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c.
//
// Solidity: event Unpaused(address indexed account, uint256 newPausedStatus)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerUnpaused)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseUnpaused(log types.Log) (*ContractDelegationManagerUnpaused, error) {
	event := new(ContractDelegationManagerUnpaused)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerWithdrawalCompletedIterator is returned from FilterWithdrawalCompleted and is used to iterate over the raw logs and unpacked data for WithdrawalCompleted events raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalCompletedIterator struct {
	Event *ContractDelegationManagerWithdrawalCompleted // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerWithdrawalCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerWithdrawalCompleted)
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
		it.Event = new(ContractDelegationManagerWithdrawalCompleted)
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
func (it *ContractDelegationManagerWithdrawalCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerWithdrawalCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerWithdrawalCompleted represents a WithdrawalCompleted event raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalCompleted struct {
	WithdrawalRoot [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalCompleted is a free log retrieval operation binding the contract event 0xc97098c2f658800b4df29001527f7324bcdffcf6e8751a699ab920a1eced5b1d.
//
// Solidity: event WithdrawalCompleted(bytes32 withdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterWithdrawalCompleted(opts *bind.FilterOpts) (*ContractDelegationManagerWithdrawalCompletedIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "WithdrawalCompleted")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerWithdrawalCompletedIterator{contract: _ContractDelegationManager.contract, event: "WithdrawalCompleted", logs: logs, sub: sub}, nil
}

// WatchWithdrawalCompleted is a free log subscription operation binding the contract event 0xc97098c2f658800b4df29001527f7324bcdffcf6e8751a699ab920a1eced5b1d.
//
// Solidity: event WithdrawalCompleted(bytes32 withdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchWithdrawalCompleted(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerWithdrawalCompleted) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "WithdrawalCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerWithdrawalCompleted)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalCompleted", log); err != nil {
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

// ParseWithdrawalCompleted is a log parse operation binding the contract event 0xc97098c2f658800b4df29001527f7324bcdffcf6e8751a699ab920a1eced5b1d.
//
// Solidity: event WithdrawalCompleted(bytes32 withdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseWithdrawalCompleted(log types.Log) (*ContractDelegationManagerWithdrawalCompleted, error) {
	event := new(ContractDelegationManagerWithdrawalCompleted)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerWithdrawalMigratedIterator is returned from FilterWithdrawalMigrated and is used to iterate over the raw logs and unpacked data for WithdrawalMigrated events raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalMigratedIterator struct {
	Event *ContractDelegationManagerWithdrawalMigrated // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerWithdrawalMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerWithdrawalMigrated)
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
		it.Event = new(ContractDelegationManagerWithdrawalMigrated)
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
func (it *ContractDelegationManagerWithdrawalMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerWithdrawalMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerWithdrawalMigrated represents a WithdrawalMigrated event raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalMigrated struct {
	OldWithdrawalRoot [32]byte
	NewWithdrawalRoot [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalMigrated is a free log retrieval operation binding the contract event 0xdc00758b65eef71dc3780c04ebe36cab6bdb266c3a698187e29e0f0dca012630.
//
// Solidity: event WithdrawalMigrated(bytes32 oldWithdrawalRoot, bytes32 newWithdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterWithdrawalMigrated(opts *bind.FilterOpts) (*ContractDelegationManagerWithdrawalMigratedIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "WithdrawalMigrated")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerWithdrawalMigratedIterator{contract: _ContractDelegationManager.contract, event: "WithdrawalMigrated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalMigrated is a free log subscription operation binding the contract event 0xdc00758b65eef71dc3780c04ebe36cab6bdb266c3a698187e29e0f0dca012630.
//
// Solidity: event WithdrawalMigrated(bytes32 oldWithdrawalRoot, bytes32 newWithdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchWithdrawalMigrated(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerWithdrawalMigrated) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "WithdrawalMigrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerWithdrawalMigrated)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalMigrated", log); err != nil {
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

// ParseWithdrawalMigrated is a log parse operation binding the contract event 0xdc00758b65eef71dc3780c04ebe36cab6bdb266c3a698187e29e0f0dca012630.
//
// Solidity: event WithdrawalMigrated(bytes32 oldWithdrawalRoot, bytes32 newWithdrawalRoot)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseWithdrawalMigrated(log types.Log) (*ContractDelegationManagerWithdrawalMigrated, error) {
	event := new(ContractDelegationManagerWithdrawalMigrated)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractDelegationManagerWithdrawalQueuedIterator is returned from FilterWithdrawalQueued and is used to iterate over the raw logs and unpacked data for WithdrawalQueued events raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalQueuedIterator struct {
	Event *ContractDelegationManagerWithdrawalQueued // Event containing the contract specifics and raw log

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
func (it *ContractDelegationManagerWithdrawalQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDelegationManagerWithdrawalQueued)
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
		it.Event = new(ContractDelegationManagerWithdrawalQueued)
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
func (it *ContractDelegationManagerWithdrawalQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDelegationManagerWithdrawalQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDelegationManagerWithdrawalQueued represents a WithdrawalQueued event raised by the ContractDelegationManager contract.
type ContractDelegationManagerWithdrawalQueued struct {
	WithdrawalRoot [32]byte
	Withdrawal     IDelegationManagerWithdrawal
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalQueued is a free log retrieval operation binding the contract event 0x9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f9.
//
// Solidity: event WithdrawalQueued(bytes32 withdrawalRoot, (address,address,address,uint256,uint32,address[],uint256[]) withdrawal)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) FilterWithdrawalQueued(opts *bind.FilterOpts) (*ContractDelegationManagerWithdrawalQueuedIterator, error) {

	logs, sub, err := _ContractDelegationManager.contract.FilterLogs(opts, "WithdrawalQueued")
	if err != nil {
		return nil, err
	}
	return &ContractDelegationManagerWithdrawalQueuedIterator{contract: _ContractDelegationManager.contract, event: "WithdrawalQueued", logs: logs, sub: sub}, nil
}

// WatchWithdrawalQueued is a free log subscription operation binding the contract event 0x9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f9.
//
// Solidity: event WithdrawalQueued(bytes32 withdrawalRoot, (address,address,address,uint256,uint32,address[],uint256[]) withdrawal)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) WatchWithdrawalQueued(opts *bind.WatchOpts, sink chan<- *ContractDelegationManagerWithdrawalQueued) (event.Subscription, error) {

	logs, sub, err := _ContractDelegationManager.contract.WatchLogs(opts, "WithdrawalQueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDelegationManagerWithdrawalQueued)
				if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalQueued", log); err != nil {
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

// ParseWithdrawalQueued is a log parse operation binding the contract event 0x9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f9.
//
// Solidity: event WithdrawalQueued(bytes32 withdrawalRoot, (address,address,address,uint256,uint32,address[],uint256[]) withdrawal)
func (_ContractDelegationManager *ContractDelegationManagerFilterer) ParseWithdrawalQueued(log types.Log) (*ContractDelegationManagerWithdrawalQueued, error) {
	event := new(ContractDelegationManagerWithdrawalQueued)
	if err := _ContractDelegationManager.contract.UnpackLog(event, "WithdrawalQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
