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
	DeprecatedEarningsReceiver common.Address
	DelegationApprover         common.Address
	StakerOptOutWindowBlocks   uint32
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

// ContractDelegationManagerMetaData contains all meta data concerning the ContractDelegationManager contract.
var ContractDelegationManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_strategyManager\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"_slasher\",\"type\":\"address\",\"internalType\":\"contractISlasher\"},{\"name\":\"_eigenPodManager\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DELEGATION_APPROVAL_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_STAKER_OPT_OUT_WINDOW_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_WITHDRAWAL_DELAY_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"STAKER_DELEGATION_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beaconChainETHStrategy\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateCurrentStakerDelegationDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateDelegationApprovalDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateStakerDelegationDigestHash\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_stakerNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateWithdrawalRoot\",\"inputs\":[{\"name\":\"withdrawal\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"completeQueuedWithdrawal\",\"inputs\":[{\"name\":\"withdrawal\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"contractIERC20[]\"},{\"name\":\"middlewareTimesIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiveAsTokens\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeQueuedWithdrawals\",\"inputs\":[{\"name\":\"withdrawals\",\"type\":\"tuple[]\",\"internalType\":\"structIDelegationManager.Withdrawal[]\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]},{\"name\":\"tokens\",\"type\":\"address[][]\",\"internalType\":\"contractIERC20[][]\"},{\"name\":\"middlewareTimesIndexes\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"receiveAsTokens\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cumulativeWithdrawalsQueued\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseDelegatedShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateTo\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approverSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateToBySignature\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSignatureAndExpiry\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"approverSalt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegatedTo\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegationApprover\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegationApproverSaltIsSpent\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"domainSeparator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eigenPodManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDelegatableShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorShares\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWithdrawalDelay\",\"inputs\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseDelegatedShares\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"initialPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minWithdrawalDelayBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"_withdrawalDelayBlocks\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isDelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minWithdrawalDelayBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"modifyOperatorDetails\",\"inputs\":[{\"name\":\"newOperatorDetails\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"__deprecated_earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorDetails\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"__deprecated_earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorShares\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseAll\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingWithdrawals\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queueWithdrawals\",\"inputs\":[{\"name\":\"queuedWithdrawalParams\",\"type\":\"tuple[]\",\"internalType\":\"structIDelegationManager.QueuedWithdrawalParams[]\",\"components\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerAsOperator\",\"inputs\":[{\"name\":\"registeringOperatorDetails\",\"type\":\"tuple\",\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"__deprecated_earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinWithdrawalDelayBlocks\",\"inputs\":[{\"name\":\"newMinWithdrawalDelayBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPauserRegistry\",\"inputs\":[{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"internalType\":\"contractIPauserRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStrategyWithdrawalDelayBlocks\",\"inputs\":[{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"withdrawalDelayBlocks\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slasher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractISlasher\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakerNonce\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakerOptOutWindowBlocks\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"strategyWithdrawalDelayBlocks\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStrategy\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"undelegate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"withdrawalRoots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperatorMetadataURI\",\"inputs\":[{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinWithdrawalDelayBlocksSet\",\"inputs\":[{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDetailsModified\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOperatorDetails\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"__deprecated_earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorMetadataURIUpdated\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operatorDetails\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.OperatorDetails\",\"components\":[{\"name\":\"__deprecated_earningsReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationApprover\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakerOptOutWindowBlocks\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSharesDecreased\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorSharesIncreased\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserRegistrySet\",\"inputs\":[{\"name\":\"pauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"},{\"name\":\"newPauserRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIPauserRegistry\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerDelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerForceUndelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerUndelegated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyWithdrawalDelayBlocksSet\",\"inputs\":[{\"name\":\"strategy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIStrategy\"},{\"name\":\"previousValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newPausedStatus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalCompleted\",\"inputs\":[{\"name\":\"withdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalQueued\",\"inputs\":[{\"name\":\"withdrawalRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"withdrawal\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDelegationManager.Withdrawal\",\"components\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatedTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"strategies\",\"type\":\"address[]\",\"internalType\":\"contractIStrategy[]\"},{\"name\":\"shares\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"anonymous\":false}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200a4f53803806200a4f58339818101604052810190620000389190620002ca565b8282828273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660c08173ffffffffffffffffffffffffffffffffffffffff16815250508173ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050505050620000ea620000fb60201b60201c565b4660e081815250505050506200040a565b600060019054906101000a900460ff16156200014e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200014590620003ad565b60405180910390fd5b60ff801660008054906101000a900460ff1660ff161015620001c05760ff6000806101000a81548160ff021916908360ff1602179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860ff604051620001b79190620003ed565b60405180910390a15b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001f482620001c7565b9050919050565b60006200020882620001e7565b9050919050565b6200021a81620001fb565b81146200022657600080fd5b50565b6000815190506200023a816200020f565b92915050565b60006200024d82620001e7565b9050919050565b6200025f8162000240565b81146200026b57600080fd5b50565b6000815190506200027f8162000254565b92915050565b60006200029282620001e7565b9050919050565b620002a48162000285565b8114620002b057600080fd5b50565b600081519050620002c48162000299565b92915050565b600080600060608486031215620002e657620002e5620001c2565b5b6000620002f68682870162000229565b935050602062000309868287016200026e565b92505060406200031c86828701620002b3565b9150509250925092565b600082825260208201905092915050565b7f496e697469616c697a61626c653a20636f6e747261637420697320696e69746960008201527f616c697a696e6700000000000000000000000000000000000000000000000000602082015250565b60006200039560278362000326565b9150620003a28262000337565b604082019050919050565b60006020820190508181036000830152620003c88162000386565b9050919050565b600060ff82169050919050565b620003e781620003cf565b82525050565b6000602082019050620004046000830184620003dc565b92915050565b60805160a05160c05160e05161a0566200049f600039600061329401526000818161127c015281816117cf01528181611ba7015281816125e10152818161360e01528181614def0152615379015260006123840152600081816112270152818161177a01528181611a5b01528181612680015281816136ef015281816137e401528181614f98015261540d015261a0566000f3fe608060405234801561001057600080fd5b50600436106103425760003560e01c8063635bbd10116101b8578063b7f06ebe11610104578063cf80873e116100a2578063f16172b01161007c578063f16172b014610a6e578063f2fde38b14610a8a578063f698da2514610aa6578063fabc1cbc14610ac457610342565b8063cf80873e146109f1578063da8be86414610a22578063eea9064b14610a5257610342565b8063c488375a116100de578063c488375a14610943578063c5e480db14610973578063c94b5111146109a3578063ca661c04146109d357610342565b8063b7f06ebe146108c5578063bb45fef2146108f5578063c448feb81461092557610342565b8063886f1195116101715780639104c3191161014b5780639104c3191461083d57806399be81c81461085b578063a178848414610877578063b1344271146108a757610342565b8063886f1195146107d15780638da5cb5b146107ef578063900413471461080d57610342565b8063635bbd10146106ff57806365da12641461071b5780636d70f7ae1461074b578063715018a61461077b578063778e55f3146107855780637f548071146107b557610342565b806328a573ae116102925780634665bcda11610230578063597b36da1161020a578063597b36da146106655780635ac86ab7146106955780635c975abb146106c557806360d7faed146106e357610342565b80634665bcda1461061f5780634fc40b611461063d578063595c6a671461065b57610342565b806339b70e381161026c57806339b70e38146105835780633cdeb5e0146105a15780633e28391d146105d1578063433773821461060157610342565b806328a573ae1461051b57806329c77d4f14610537578063334043961461056757610342565b8063132d4967116102ff57806316928365116102d957806316928365146104815780631bbce091146104b157806320606b70146104e157806322bf40e4146104ff57610342565b8063132d49671461042d578063136439dd146104495780631522bf021461046557610342565b80630449ca391461034757806304a4f979146103775780630b9f487a146103955780630dd8dd02146103c55780630f589e59146103f557806310d67a2f14610411575b600080fd5b610361600480360381019061035c9190615998565b610ae0565b60405161036e91906159fe565b60405180910390f35b61037f610b8a565b60405161038c9190615a32565b60405180910390f35b6103af60048036038101906103aa9190615b03565b610bae565b6040516103bc9190615a32565b60405180910390f35b6103df60048036038101906103da9190615bd4565b610c46565b6040516103ec9190615cdf565b60405180910390f35b61040f600480360381019061040a9190615d7b565b61100d565b005b61042b60048036038101906104269190615e19565b61111b565b005b61044760048036038101906104429190615e84565b611225565b005b610463600480360381019061045e9190615ed7565b61138e565b005b61047f600480360381019061047a9190615f5a565b611509565b005b61049b60048036038101906104969190615fdb565b611523565b6040516104a891906159fe565b60405180910390f35b6104cb60048036038101906104c69190616008565b611585565b6040516104d89190615a32565b60405180910390f35b6104e96115e0565b6040516104f69190615a32565b60405180910390f35b6105196004803603810190610514919061605b565b611604565b005b61053560048036038101906105309190615e84565b611778565b005b610551600480360381019061054c9190615fdb565b6118e1565b60405161055e91906159fe565b60405180910390f35b610581600480360381019061057c919061622c565b6118f9565b005b61058b611a59565b6040516105989190616374565b60405180910390f35b6105bb60048036038101906105b69190615fdb565b611a7d565b6040516105c8919061639e565b60405180910390f35b6105eb60048036038101906105e69190615fdb565b611ae9565b6040516105f891906163d4565b60405180910390f35b610609611b81565b6040516106169190615a32565b60405180910390f35b610627611ba5565b6040516106349190616410565b60405180910390f35b610645611bc9565b60405161065291906159fe565b60405180910390f35b610663611bd0565b005b61067f600480360381019061067a919061676f565b611d42565b60405161068c9190615a32565b60405180910390f35b6106af60048036038101906106aa91906167f1565b611d72565b6040516106bc91906163d4565b60405180910390f35b6106cd611d8e565b6040516106da91906159fe565b60405180910390f35b6106fd60048036038101906106f891906168bf565b611d98565b005b61071960048036038101906107149190615ed7565b611e4e565b005b61073560048036038101906107309190615fdb565b611e62565b604051610742919061639e565b60405180910390f35b61076560048036038101906107609190615fdb565b611e95565b60405161077291906163d4565b60405180910390f35b610783611f64565b005b61079f600480360381019061079a9190616963565b611f78565b6040516107ac91906159fe565b60405180910390f35b6107cf60048036038101906107ca9190616ac4565b611f9d565b005b6107d9612138565b6040516107e69190616b98565b60405180910390f35b6107f761215e565b604051610804919061639e565b60405180910390f35b61082760048036038101906108229190616bb3565b612188565b6040516108349190616ccd565b60405180910390f35b6108456122b6565b6040516108529190616d10565b60405180910390f35b61087560048036038101906108709190616d2b565b6122ce565b005b610891600480360381019061088c9190615fdb565b61236a565b60405161089e91906159fe565b60405180910390f35b6108af612382565b6040516108bc9190616d99565b60405180910390f35b6108df60048036038101906108da9190616db4565b6123a6565b6040516108ec91906163d4565b60405180910390f35b61090f600480360381019061090a9190616de1565b6123c6565b60405161091c91906163d4565b60405180910390f35b61092d6123f5565b60405161093a91906159fe565b60405180910390f35b61095d60048036038101906109589190616e21565b6123fb565b60405161096a91906159fe565b60405180910390f35b61098d60048036038101906109889190615fdb565b612413565b60405161099a9190616eae565b60405180910390f35b6109bd60048036038101906109b89190616ec9565b61253e565b6040516109ca9190615a32565b60405180910390f35b6109db6125d3565b6040516109e891906159fe565b60405180910390f35b610a0b6004803603810190610a069190615fdb565b6125da565b604051610a19929190616fee565b60405180910390f35b610a3c6004803603810190610a379190615fdb565b612a71565b604051610a499190615cdf565b60405180910390f35b610a6c6004803603810190610a679190617025565b613115565b005b610a886004803603810190610a839190617094565b6131b7565b005b610aa46004803603810190610a9f9190615fdb565b61320c565b005b610aae613290565b604051610abb9190615a32565b60405180910390f35b610ade6004803603810190610ad99190615ed7565b6132d2565b005b600080609d54905060005b84849050811015610b7f57600060a16000878785818110610b0f57610b0e6170c1565b5b9050602002016020810190610b249190616e21565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811115610b6d578092505b5080610b789061711f565b9050610aeb565b508091505092915050565b7f14bde674c9f64b2ad00eaaee4a8bed1fabef35c7507e3c5b9cfc9436909a2dad81565b6000807f14bde674c9f64b2ad00eaaee4a8bed1fabef35c7507e3c5b9cfc9436909a2dad8588888787604051602001610bec96959493929190617168565b6040516020818303038152906040528051906020012090506000610c0e613290565b82604051602001610c20929190617241565b604051602081830303815290604052805190602001209050809250505095945050505050565b60606001610c5381611d72565b15610c93576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c8a906172d5565b60405180910390fd5b60008484905067ffffffffffffffff811115610cb257610cb1616441565b5b604051908082528060200260200182016040528015610ce05781602001602082028036833780820191505090505b5090506000609a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060005b8686905081101561100057868682818110610d6857610d676170c1565b5b9050602002810190610d7a9190617304565b8060200190610d89919061732c565b9050878783818110610d9e57610d9d6170c1565b5b9050602002810190610db09190617304565b8060000190610dbf919061738f565b905014610e01576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610df890617464565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff16878783818110610e2b57610e2a6170c1565b5b9050602002810190610e3d9190617304565b6040016020810190610e4f9190615fdb565b73ffffffffffffffffffffffffffffffffffffffff1614610ea5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e9c906174f6565b60405180910390fd5b610fce3383898985818110610ebd57610ebc6170c1565b5b9050602002810190610ecf9190617304565b6040016020810190610ee19190615fdb565b8a8a86818110610ef457610ef36170c1565b5b9050602002810190610f069190617304565b8060000190610f15919061738f565b80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050508b8b87818110610f6857610f676170c1565b5b9050602002810190610f7a9190617304565b8060200190610f89919061732c565b80806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f82011690508083019250505050505050613473565b838281518110610fe157610fe06170c1565b5b6020026020010181815250508080610ff89061711f565b915050610d4a565b5081935050505092915050565b61101633611ae9565b15611056576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161104d906175ae565b60405180910390fd5b6110603384613a4b565b6110686158b2565b6110773333836000801b613bfd565b3373ffffffffffffffffffffffffffffffffffffffff167f8e8485583a2310d41f7c82b9427d0bd49bad74bb9cff9d3402a29d8f9b28a0e2856040516110bd9190617656565b60405180910390a23373ffffffffffffffffffffffffffffffffffffffff167f02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b6708090848460405161110d92919061769e565b60405180910390a250505050565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611188573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111ac91906176d7565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611219576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161121090617776565b60405180910390fd5b61122281614013565b50565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806112ca57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b611309576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161130090617808565b60405180910390fd5b61131283611ae9565b15611389576000609a60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905061138781858585614122565b505b505050565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b81526004016113e9919061639e565b602060405180830381865afa158015611406573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061142a919061783d565b611469576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611460906178dc565b60405180910390fd5b6066548160665416146114b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114a89061796e565b60405180910390fd5b806066819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d826040516114fe91906159fe565b60405180910390a250565b61151161420d565b61151d8484848461428b565b50505050565b6000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160149054906101000a900463ffffffff1663ffffffff169050919050565b600080609b60008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490506115d68582868661253e565b9150509392505050565b7f8cad95687ba82c2ce50e74f7b754645e5117c3a5bec8151c0726d5857980a86681565b60008060019054906101000a900460ff161590508080156116355750600160008054906101000a900460ff1660ff16105b80611662575061164430614455565b1580156116615750600160008054906101000a900460ff1660ff16145b5b6116a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161169890617a00565b60405180910390fd5b60016000806101000a81548160ff021916908360ff16021790555080156116de576001600060016101000a81548160ff0219169083151502179055505b6116e88888614478565b6116f06145a4565b6097819055506116ff89614634565b611708866146fa565b6117148585858561428b565b801561176d5760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860016040516117649190617a5b565b60405180910390a15b505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061181d57507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b61185c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161185390617808565b60405180910390fd5b61186583611ae9565b156118dc576000609a60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506118da81858585614785565b505b505050565b609b6020528060005260406000206000915090505481565b600261190481611d72565b15611944576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161193b906172d5565b60405180910390fd5b600260c954141561198a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161198190617ac2565b60405180910390fd5b600260c98190555060005b89899050811015611a4557611a348a8a838181106119b6576119b56170c1565b5b90506020028101906119c89190617ae2565b8989848181106119db576119da6170c1565b5b90506020028101906119ed9190617b0a565b898986818110611a00576119ff6170c1565b5b90506020020135888887818110611a1a57611a196170c1565b5b9050602002016020810190611a2f9190617b6d565b614870565b80611a3e9061711f565b9050611995565b50600160c981905550505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b60008073ffffffffffffffffffffffffffffffffffffffff16609a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614159050919050565b7f39111bc4a4d688e1f685123d7497d4615370152a8ee4a0593e647bd06ad8bb0b81565b7f000000000000000000000000000000000000000000000000000000000000000081565b6213c68081565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166346fbf68e336040518263ffffffff1660e01b8152600401611c2b919061639e565b602060405180830381865afa158015611c48573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c6c919061783d565b611cab576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611ca2906178dc565b60405180910390fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6066819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff604051611d3891906159fe565b60405180910390a2565b600081604051602001611d559190617d1b565b604051602081830303815290604052805190602001209050919050565b6000808260ff166001901b905080816066541614915050919050565b6000606654905090565b6002611da381611d72565b15611de3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611dda906172d5565b60405180910390fd5b600260c9541415611e29576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611e2090617ac2565b60405180910390fd5b600260c981905550611e3e8686868686614870565b600160c981905550505050505050565b611e5661420d565b611e5f816146fa565b50565b609a6020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614158015611f5d57508173ffffffffffffffffffffffffffffffffffffffff16609a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b9050919050565b611f6c61420d565b611f766000614634565b565b6098602052816000526040600020602052806000526040600020600091509150505481565b4283602001511015611fe4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611fdb90617dd5565b60405180910390fd5b611fed85611ae9565b1561202d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161202490617e8d565b60405180910390fd5b61203684611e95565b612075576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161206c90617f45565b60405180910390fd5b6000609b60008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905060006120cb878388886020015161253e565b905060018201609b60008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061212387828760000151615195565b61212f87878686613bfd565b50505050505050565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60606000825167ffffffffffffffff8111156121a7576121a6616441565b5b6040519080825280602002602001820160405280156121d55781602001602082028036833780820191505090505b50905060005b83518110156122ab57609860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000858381518110612238576122376170c1565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205482828151811061228e5761228d6170c1565b5b602002602001018181525050806122a49061711f565b90506121db565b508091505092915050565b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac081565b6122d733611e95565b612316576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161230d90617ffd565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff167f02a919ed0e2acad1dd90f17ef2fa4ae5462ee1339170034a8531cca4b6708090838360405161235e92919061769e565b60405180910390a25050565b609f6020528060005260406000206000915090505481565b7f000000000000000000000000000000000000000000000000000000000000000081565b609e6020528060005260406000206000915054906101000a900460ff1681565b609c6020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b609d5481565b60a16020528060005260406000206000915090505481565b61241b6158cc565b609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060600160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160149054906101000a900463ffffffff1663ffffffff1663ffffffff16815250509050919050565b6000807f39111bc4a4d688e1f685123d7497d4615370152a8ee4a0593e647bd06ad8bb0b8685878660405160200161257a95949392919061801d565b604051602081830303815290604052805190602001209050600061259c613290565b826040516020016125ae929190617241565b6040516020818303038152906040528051906020012090508092505050949350505050565b62034bc081565b60608060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166360f4062b856040518263ffffffff1660e01b8152600401612638919061639e565b602060405180830381865afa158015612655573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061267991906180a6565b90506000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166394f649dd876040518263ffffffff1660e01b81526004016126d7919061639e565b600060405180830381865afa1580156126f4573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061271d919061822b565b915091506000831361273757818194509450505050612a6c565b60608060008451141561286557600167ffffffffffffffff81111561275f5761275e616441565b5b60405190808252806020026020018201604052801561278d5781602001602082028036833780820191505090505b509150600167ffffffffffffffff8111156127ab576127aa616441565b5b6040519080825280602002602001820160405280156127d95781602001602082028036833780820191505090505b50905073beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac082600081518110612805576128046170c1565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508481600081518110612854576128536170c1565b5b602002602001018181525050612a60565b6001845161287391906182a3565b67ffffffffffffffff81111561288c5761288b616441565b5b6040519080825280602002602001820160405280156128ba5781602001602082028036833780820191505090505b509150815167ffffffffffffffff8111156128d8576128d7616441565b5b6040519080825280602002602001820160405280156129065781602001602082028036833780820191505090505b50905060005b84518110156129c257848181518110612928576129276170c1565b5b6020026020010151838281518110612943576129426170c1565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508381815181106129905761298f6170c1565b5b60200260200101518282815181106129ab576129aa6170c1565b5b60200260200101818152505080600101905061290c565b5073beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac082600184516129e791906182f9565b815181106129f8576129f76170c1565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050848160018451612a4291906182f9565b81518110612a5357612a526170c1565b5b6020026020010181815250505b81819650965050505050505b915091565b60606001612a7e81611d72565b15612abe576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612ab5906172d5565b60405180910390fd5b612ac783611ae9565b612b06576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612afd906183c5565b60405180910390fd5b612b0f83611e95565b15612b4f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612b4690618457565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415612bbf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612bb6906184e9565b60405180910390fd5b6000609a60008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480612c8857508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b80612d205750609960008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b612d5f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401612d569061857b565b60405180910390fd5b600080612d6b866125da565b915091508573ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614612dfd578273ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167ff0eddf07e6ea14f388b47e1e94a0f464ecbd9eed4171130e0fc0e99fb4030a8a60405160405180910390a35b8273ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167ffee30966a256b71e14bc0ebfc94315e28ef4a97a7131a9e2b7a310a73af4467660405160405180910390a36000609a60008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600082511415612f3157600067ffffffffffffffff811115612efb57612efa616441565b5b604051908082528060200260200182016040528015612f295781602001602082028036833780820191505090505b50945061310c565b815167ffffffffffffffff811115612f4c57612f4b616441565b5b604051908082528060200260200182016040528015612f7a5781602001602082028036833780820191505090505b50945060005b825181101561310a576000600167ffffffffffffffff811115612fa657612fa5616441565b5b604051908082528060200260200182016040528015612fd45781602001602082028036833780820191505090505b5090506000600167ffffffffffffffff811115612ff457612ff3616441565b5b6040519080825280602002602001820160405280156130225781602001602082028036833780820191505090505b509050848381518110613038576130376170c1565b5b602002602001015182600081518110613054576130536170c1565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508383815181106130a1576130a06170c1565b5b6020026020010151816000815181106130bd576130bc6170c1565b5b6020026020010181815250506130d689878b8585613473565b8884815181106130e9576130e86170c1565b5b602002602001018181525050505080806131029061711f565b915050612f80565b505b50505050919050565b61311e33611ae9565b1561315e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161315590618633565b60405180910390fd5b61316783611e95565b6131a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161319d906186eb565b60405180910390fd5b6131b233848484613bfd565b505050565b6131c033611e95565b6131ff576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016131f6906187a3565b60405180910390fd5b6132093382613a4b565b50565b61321461420d565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415613284576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161327b90618835565b60405180910390fd5b61328d81614634565b50565b60007f00000000000000000000000000000000000000000000000000000000000000004614156132c45760975490506132cf565b6132cc6145a4565b90505b90565b606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663eab66d7a6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561333f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061336391906176d7565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146133d0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016133c790617776565b60405180910390fd5b60665419811960665419161461341b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613412906188c7565b60405180910390fd5b806066819055503373ffffffffffffffffffffffffffffffffffffffff167f3582d1828e26bf56bd801502bc021ac0bc8afb57c826e4986b45593c8fad389c8260405161346891906159fe565b60405180910390a250565b60008073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1614156134e4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016134db9061897f565b60405180910390fd5b600083511415613529576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161352090618a37565b60405180910390fd5b60005b83518110156138b157600073ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff16146135aa576135a98688868481518110613581576135806170c1565b5b602002602001015186858151811061359c5761359b6170c1565b5b6020026020010151614122565b5b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac073ffffffffffffffffffffffffffffffffffffffff168482815181106135e8576135e76170c1565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1614156136b8577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663beffbb898885848151811061365c5761365b6170c1565b5b60200260200101516040518363ffffffff1660e01b8152600401613681929190618a57565b600060405180830381600087803b15801561369b57600080fd5b505af11580156136af573d6000803e3d6000fd5b505050506138a6565b8473ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff1614806137a357507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16639b4da03d85838151811061373c5761373b6170c1565b5b60200260200101516040518263ffffffff1660e01b81526004016137609190616d10565b602060405180830381865afa15801561377d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137a1919061783d565b155b6137e2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016137d990618b64565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16638c80d4e588868481518110613832576138316170c1565b5b602002602001015186858151811061384d5761384c6170c1565b5b60200260200101516040518463ffffffff1660e01b815260040161387393929190618b84565b600060405180830381600087803b15801561388d57600080fd5b505af11580156138a1573d6000803e3d6000fd5b505050505b80600101905061352c565b506000609f60008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050609f60008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008154809291906139469061711f565b919050555060006040518060e001604052808973ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018381526020014363ffffffff16815260200186815260200185815250905060006139d582611d42565b90506001609e600083815260200190815260200160002060006101000a81548160ff0219169083151502179055507f9009ab153e8014fbfb02f2217f5cde7aa7f9ad734ae85ca3ee3f4ca2fdd499f98183604051613a34929190618bbb565b60405180910390a180935050505095945050505050565b6213c680816040016020810190613a629190618beb565b63ffffffff161115613aa9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613aa090618cd6565b60405180910390fd5b609960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160149054906101000a900463ffffffff1663ffffffff16816040016020810190613b159190618beb565b63ffffffff161015613b5c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613b5390618d8e565b60405180910390fd5b80609960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208181613ba89190618f59565b9050503373ffffffffffffffffffffffffffffffffffffffff167ffebe5cd24b2cbc7b065b9d0fdeb904461e4afcff57dd57acda1e7832031ba7ac82604051613bf19190617656565b60405180910390a25050565b6000613c0881611d72565b15613c48576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613c3f906172d5565b60405180910390fd5b6000609960008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614158015613d1857508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b8015613d5057508473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15613ec9574284602001511015613d9c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613d9390618fd9565b60405180910390fd5b609c60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600084815260200190815260200160002060009054906101000a900460ff1615613e3a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401613e319061906b565b60405180910390fd5b6001609c60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600085815260200190815260200160002060006101000a81548160ff0219169083151502179055506000613eb6878784878960200151610bae565b9050613ec782828760000151615195565b505b84609a60008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff167fc3ee9f2e5fda98e8066a1f745b2df9285f416fe98cf2559cd21484b3d874330460405160405180910390a3600080613fad886125da565b9150915060005b825181101561400857613ffd888a858481518110613fd557613fd46170c1565b5b6020026020010151858581518110613ff057613fef6170c1565b5b6020026020010151614785565b806001019050613fb4565b505050505050505050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415614083576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161407a90619123565b60405180910390fd5b7f6e9fcd539896fca60e8b0f01dd580233e48a6b0f7df013b89ba7f565869acdb6606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16826040516140d6929190619143565b60405180910390a180606560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b80609860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546141ae91906182f9565b925050819055508373ffffffffffffffffffffffffffffffffffffffff167f6909600037b75d7b4733aedd815442b5ec018a827751c832aaff64eba5d6d2dd8484846040516141ff93929190618b84565b60405180910390a250505050565b614215615327565b73ffffffffffffffffffffffffffffffffffffffff1661423361215e565b73ffffffffffffffffffffffffffffffffffffffff1614614289576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614280906191b8565b60405180910390fd5b565b8181905084849050146142d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016142ca90619270565b60405180910390fd5b600084849050905060005b8181101561444d5760008686838181106142fb576142fa6170c1565b5b90506020020160208101906143109190616e21565b9050600060a160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600086868581811061436b5761436a6170c1565b5b90506020020135905062034bc08111156143ba576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016143b19061934e565b60405180910390fd5b8060a160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055507f0e7efa738e8b0ce6376a0c1af471655540d2e9a81647d7b09ed823018426576d8383836040516144319392919061936e565b60405180910390a1505050806144469061711f565b90506142de565b505050505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600073ffffffffffffffffffffffffffffffffffffffff16606560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161480156145035750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614155b614542576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016145399061943d565b60405180910390fd5b806066819055503373ffffffffffffffffffffffffffffffffffffffff167fab40a374bc51de372200a8bc981af8c9ecdc08dfdaef0bb6e09f88f3c616ef3d8260405161458f91906159fe565b60405180910390a26145a082614013565b5050565b60007f8cad95687ba82c2ce50e74f7b754645e5117c3a5bec8151c0726d5857980a8666040518060400160405280600a81526020017f456967656e4c6179657200000000000000000000000000000000000000000000815250805190602001204630604051602001614619949392919061945d565b60405160208183030381529060405280519060200120905090565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b62034bc0811115614740576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161473790619560565b60405180910390fd5b7fafa003cd76f87ff9d62b35beea889920f33c0c42b8d45b74954d61d50f4b6b69609d5482604051614773929190619580565b60405180910390a180609d8190555050565b80609860008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461481191906182a3565b925050819055508373ffffffffffffffffffffffffffffffffffffffff167f1ec042c965e2edd7107b51188ee0f383e22e76179041ab3a9d18ff151405166c84848460405161486293929190618b84565b60405180910390a250505050565b60006148848661487f906195a9565b611d42565b9050609e600082815260200190815260200160002060009054906101000a900460ff166148e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016148dd90619654565b60405180910390fd5b43609d548760800160208101906148fd9190618beb565b63ffffffff1661490d91906182a3565b111561494e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016149459061970c565b60405180910390fd5b8560400160208101906149619190615fdb565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146149ce576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016149c5906197c4565b60405180910390fd5b8115614a2b57858060a001906149e4919061738f565b90508585905014614a2a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614a219061987c565b60405180910390fd5b5b609e600082815260200190815260200160002060006101000a81549060ff02191690558115614bfd5760005b868060a00190614a67919061738f565b9050811015614bf7574360a16000898060a00190614a85919061738f565b85818110614a9657614a956170c1565b5b9050602002016020810190614aab9190616e21565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054886080016020810190614af99190618beb565b63ffffffff16614b0991906182a3565b1115614b4a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614b419061995a565b60405180910390fd5b614bec876000016020810190614b609190615fdb565b33898060a00190614b71919061738f565b85818110614b8257614b816170c1565b5b9050602002016020810190614b979190616e21565b8a8060c00190614ba7919061732c565b86818110614bb857614bb76170c1565b5b905060200201358a8a87818110614bd257614bd16170c1565b5b9050602002016020810190614be791906199b8565b61532f565b806001019050614a57565b50615156565b6000609a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060005b878060a00190614c74919061738f565b9050811015615153574360a160008a8060a00190614c92919061738f565b85818110614ca357614ca26170c1565b5b9050602002016020810190614cb89190616e21565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054896080016020810190614d069190618beb565b63ffffffff16614d1691906182a3565b1115614d57576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401614d4e9061995a565b60405180910390fd5b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac073ffffffffffffffffffffffffffffffffffffffff16888060a00190614d92919061738f565b83818110614da357614da26170c1565b5b9050602002016020810190614db89190616e21565b73ffffffffffffffffffffffffffffffffffffffff161415614f96576000886000016020810190614de99190615fdb565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16630e81073c838c8060c00190614e3a919061732c565b87818110614e4b57614e4a6170c1565b5b905060200201356040518363ffffffff1660e01b8152600401614e6f929190618a57565b6020604051808303816000875af1158015614e8e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614eb291906199e5565b90506000609a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614614f8e57614f8d81848d8060a00190614f61919061738f565b88818110614f7257614f716170c1565b5b9050602002016020810190614f879190616e21565b85614785565b5b505050615148565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c4623ea133898985818110614fe657614fe56170c1565b5b9050602002016020810190614ffb91906199b8565b8b8060a0019061500b919061738f565b8681811061501c5761501b6170c1565b5b90506020020160208101906150319190616e21565b8c8060c00190615041919061732c565b87818110615052576150516170c1565b5b905060200201356040518563ffffffff1660e01b81526004016150789493929190619a33565b600060405180830381600087803b15801561509257600080fd5b505af11580156150a6573d6000803e3d6000fd5b50505050600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146151475761514682338a8060a001906150f3919061738f565b85818110615104576151036170c1565b5b90506020020160208101906151199190616e21565b8b8060c00190615129919061732c565b8681811061513a576151396170c1565b5b90506020020135614785565b5b5b806001019050614c64565b50505b7fc97098c2f658800b4df29001527f7324bcdffcf6e8751a699ab920a1eced5b1d816040516151859190615a32565b60405180910390a1505050505050565b61519e836154a4565b156152aa57631626ba7e60e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168373ffffffffffffffffffffffffffffffffffffffff16631626ba7e84846040518363ffffffff1660e01b8152600401615205929190619b00565b602060405180830381865afa158015615222573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906152469190619b88565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916146152a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161529c90619c4d565b60405180910390fd5b615322565b8273ffffffffffffffffffffffffffffffffffffffff166152cb83836154c7565b73ffffffffffffffffffffffffffffffffffffffff1614615321576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161531890619d05565b60405180910390fd5b5b505050565b600033905090565b73beac0eeeeeeeeeeeeeeeeeeeeeeeeeeeeeebeac073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561540b577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663387b13008686856040518463ffffffff1660e01b81526004016153d493929190619d25565b600060405180830381600087803b1580156153ee57600080fd5b505af1158015615402573d6000803e3d6000fd5b5050505061549d565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663c608c7f3858585856040518563ffffffff1660e01b815260040161546a9493929190619d5c565b600060405180830381600087803b15801561548457600080fd5b505af1158015615498573d6000803e3d6000fd5b505050505b5050505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b60008060006154d685856154ee565b915091506154e381615571565b819250505092915050565b6000806041835114156155305760008060006020860151925060408601519150606086015160001a905061552487828585615746565b9450945050505061556a565b604083511415615561576000806020850151915060408501519050615556868383615853565b93509350505061556a565b60006002915091505b9250929050565b6000600481111561558557615584619da1565b5b81600481111561559857615597619da1565b5b14156155a357615743565b600160048111156155b7576155b6619da1565b5b8160048111156155ca576155c9619da1565b5b141561560b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161560290619e1c565b60405180910390fd5b6002600481111561561f5761561e619da1565b5b81600481111561563257615631619da1565b5b1415615673576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161566a90619e88565b60405180910390fd5b6003600481111561568757615686619da1565b5b81600481111561569a57615699619da1565b5b14156156db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016156d290619f1a565b60405180910390fd5b6004808111156156ee576156ed619da1565b5b81600481111561570157615700619da1565b5b1415615742576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161573990619fac565b60405180910390fd5b5b50565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08360001c111561578157600060039150915061584a565b601b8560ff16141580156157995750601c8560ff1614155b156157ab57600060049150915061584a565b6000600187878787604051600081526020016040526040516157d09493929190619fdb565b6020604051602081039080840390855afa1580156157f2573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156158415760006001925092505061584a565b80600092509250505b94509492505050565b60008060007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60001b841690506000601b60ff8660001c901c61589691906182a3565b90506158a487828885615746565b935093505050935093915050565b604051806040016040528060608152602001600081525090565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600063ffffffff1681525090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261595857615957615933565b5b8235905067ffffffffffffffff81111561597557615974615938565b5b6020830191508360208202830111156159915761599061593d565b5b9250929050565b600080602083850312156159af576159ae615929565b5b600083013567ffffffffffffffff8111156159cd576159cc61592e565b5b6159d985828601615942565b92509250509250929050565b6000819050919050565b6159f8816159e5565b82525050565b6000602082019050615a1360008301846159ef565b92915050565b6000819050919050565b615a2c81615a19565b82525050565b6000602082019050615a476000830184615a23565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000615a7882615a4d565b9050919050565b615a8881615a6d565b8114615a9357600080fd5b50565b600081359050615aa581615a7f565b92915050565b615ab481615a19565b8114615abf57600080fd5b50565b600081359050615ad181615aab565b92915050565b615ae0816159e5565b8114615aeb57600080fd5b50565b600081359050615afd81615ad7565b92915050565b600080600080600060a08688031215615b1f57615b1e615929565b5b6000615b2d88828901615a96565b9550506020615b3e88828901615a96565b9450506040615b4f88828901615a96565b9350506060615b6088828901615ac2565b9250506080615b7188828901615aee565b9150509295509295909350565b60008083601f840112615b9457615b93615933565b5b8235905067ffffffffffffffff811115615bb157615bb0615938565b5b602083019150836020820283011115615bcd57615bcc61593d565b5b9250929050565b60008060208385031215615beb57615bea615929565b5b600083013567ffffffffffffffff811115615c0957615c0861592e565b5b615c1585828601615b7e565b92509250509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b615c5681615a19565b82525050565b6000615c688383615c4d565b60208301905092915050565b6000602082019050919050565b6000615c8c82615c21565b615c968185615c2c565b9350615ca183615c3d565b8060005b83811015615cd2578151615cb98882615c5c565b9750615cc483615c74565b925050600181019050615ca5565b5085935050505092915050565b60006020820190508181036000830152615cf98184615c81565b905092915050565b600080fd5b600060608284031215615d1c57615d1b615d01565b5b81905092915050565b60008083601f840112615d3b57615d3a615933565b5b8235905067ffffffffffffffff811115615d5857615d57615938565b5b602083019150836001820283011115615d7457615d7361593d565b5b9250929050565b600080600060808486031215615d9457615d93615929565b5b6000615da286828701615d06565b935050606084013567ffffffffffffffff811115615dc357615dc261592e565b5b615dcf86828701615d25565b92509250509250925092565b6000615de682615a6d565b9050919050565b615df681615ddb565b8114615e0157600080fd5b50565b600081359050615e1381615ded565b92915050565b600060208284031215615e2f57615e2e615929565b5b6000615e3d84828501615e04565b91505092915050565b6000615e5182615a6d565b9050919050565b615e6181615e46565b8114615e6c57600080fd5b50565b600081359050615e7e81615e58565b92915050565b600080600060608486031215615e9d57615e9c615929565b5b6000615eab86828701615a96565b9350506020615ebc86828701615e6f565b9250506040615ecd86828701615aee565b9150509250925092565b600060208284031215615eed57615eec615929565b5b6000615efb84828501615aee565b91505092915050565b60008083601f840112615f1a57615f19615933565b5b8235905067ffffffffffffffff811115615f3757615f36615938565b5b602083019150836020820283011115615f5357615f5261593d565b5b9250929050565b60008060008060408587031215615f7457615f73615929565b5b600085013567ffffffffffffffff811115615f9257615f9161592e565b5b615f9e87828801615942565b9450945050602085013567ffffffffffffffff811115615fc157615fc061592e565b5b615fcd87828801615f04565b925092505092959194509250565b600060208284031215615ff157615ff0615929565b5b6000615fff84828501615a96565b91505092915050565b60008060006060848603121561602157616020615929565b5b600061602f86828701615a96565b935050602061604086828701615a96565b925050604061605186828701615aee565b9150509250925092565b60008060008060008060008060c0898b03121561607b5761607a615929565b5b60006160898b828c01615a96565b985050602061609a8b828c01615e04565b97505060406160ab8b828c01615aee565b96505060606160bc8b828c01615aee565b955050608089013567ffffffffffffffff8111156160dd576160dc61592e565b5b6160e98b828c01615942565b945094505060a089013567ffffffffffffffff81111561610c5761610b61592e565b5b6161188b828c01615f04565b92509250509295985092959890939650565b60008083601f8401126161405761613f615933565b5b8235905067ffffffffffffffff81111561615d5761615c615938565b5b6020830191508360208202830111156161795761617861593d565b5b9250929050565b60008083601f84011261619657616195615933565b5b8235905067ffffffffffffffff8111156161b3576161b2615938565b5b6020830191508360208202830111156161cf576161ce61593d565b5b9250929050565b60008083601f8401126161ec576161eb615933565b5b8235905067ffffffffffffffff81111561620957616208615938565b5b6020830191508360208202830111156162255761622461593d565b5b9250929050565b6000806000806000806000806080898b03121561624c5761624b615929565b5b600089013567ffffffffffffffff81111561626a5761626961592e565b5b6162768b828c0161612a565b9850985050602089013567ffffffffffffffff8111156162995761629861592e565b5b6162a58b828c01616180565b9650965050604089013567ffffffffffffffff8111156162c8576162c761592e565b5b6162d48b828c01615f04565b9450945050606089013567ffffffffffffffff8111156162f7576162f661592e565b5b6163038b828c016161d6565b92509250509295985092959890939650565b6000819050919050565b600061633a61633561633084615a4d565b616315565b615a4d565b9050919050565b600061634c8261631f565b9050919050565b600061635e82616341565b9050919050565b61636e81616353565b82525050565b60006020820190506163896000830184616365565b92915050565b61639881615a6d565b82525050565b60006020820190506163b3600083018461638f565b92915050565b60008115159050919050565b6163ce816163b9565b82525050565b60006020820190506163e960008301846163c5565b92915050565b60006163fa82616341565b9050919050565b61640a816163ef565b82525050565b60006020820190506164256000830184616401565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61647982616430565b810181811067ffffffffffffffff8211171561649857616497616441565b5b80604052505050565b60006164ab61591f565b90506164b78282616470565b919050565b600080fd5b600063ffffffff82169050919050565b6164da816164c1565b81146164e557600080fd5b50565b6000813590506164f7816164d1565b92915050565b600067ffffffffffffffff82111561651857616517616441565b5b602082029050602081019050919050565b600061653c616537846164fd565b6164a1565b9050808382526020820190506020840283018581111561655f5761655e61593d565b5b835b8181101561658857806165748882615e6f565b845260208401935050602081019050616561565b5050509392505050565b600082601f8301126165a7576165a6615933565b5b81356165b7848260208601616529565b91505092915050565b600067ffffffffffffffff8211156165db576165da616441565b5b602082029050602081019050919050565b60006165ff6165fa846165c0565b6164a1565b905080838252602082019050602084028301858111156166225761662161593d565b5b835b8181101561664b57806166378882615aee565b845260208401935050602081019050616624565b5050509392505050565b600082601f83011261666a57616669615933565b5b813561667a8482602086016165ec565b91505092915050565b600060e082840312156166995761669861642b565b5b6166a360e06164a1565b905060006166b384828501615a96565b60008301525060206166c784828501615a96565b60208301525060406166db84828501615a96565b60408301525060606166ef84828501615aee565b6060830152506080616703848285016164e8565b60808301525060a082013567ffffffffffffffff811115616727576167266164bc565b5b61673384828501616592565b60a08301525060c082013567ffffffffffffffff811115616757576167566164bc565b5b61676384828501616655565b60c08301525092915050565b60006020828403121561678557616784615929565b5b600082013567ffffffffffffffff8111156167a3576167a261592e565b5b6167af84828501616683565b91505092915050565b600060ff82169050919050565b6167ce816167b8565b81146167d957600080fd5b50565b6000813590506167eb816167c5565b92915050565b60006020828403121561680757616806615929565b5b6000616815848285016167dc565b91505092915050565b600060e0828403121561683457616833615d01565b5b81905092915050565b60008083601f84011261685357616852615933565b5b8235905067ffffffffffffffff8111156168705761686f615938565b5b60208301915083602082028301111561688c5761688b61593d565b5b9250929050565b61689c816163b9565b81146168a757600080fd5b50565b6000813590506168b981616893565b92915050565b6000806000806000608086880312156168db576168da615929565b5b600086013567ffffffffffffffff8111156168f9576168f861592e565b5b6169058882890161681e565b955050602086013567ffffffffffffffff8111156169265761692561592e565b5b6169328882890161683d565b9450945050604061694588828901615aee565b9250506060616956888289016168aa565b9150509295509295909350565b6000806040838503121561697a57616979615929565b5b600061698885828601615a96565b925050602061699985828601615e6f565b9150509250929050565b600080fd5b600067ffffffffffffffff8211156169c3576169c2616441565b5b6169cc82616430565b9050602081019050919050565b82818337600083830152505050565b60006169fb6169f6846169a8565b6164a1565b905082815260208101848484011115616a1757616a166169a3565b5b616a228482856169d9565b509392505050565b600082601f830112616a3f57616a3e615933565b5b8135616a4f8482602086016169e8565b91505092915050565b600060408284031215616a6e57616a6d61642b565b5b616a7860406164a1565b9050600082013567ffffffffffffffff811115616a9857616a976164bc565b5b616aa484828501616a2a565b6000830152506020616ab884828501615aee565b60208301525092915050565b600080600080600060a08688031215616ae057616adf615929565b5b6000616aee88828901615a96565b9550506020616aff88828901615a96565b945050604086013567ffffffffffffffff811115616b2057616b1f61592e565b5b616b2c88828901616a58565b935050606086013567ffffffffffffffff811115616b4d57616b4c61592e565b5b616b5988828901616a58565b9250506080616b6a88828901615ac2565b9150509295509295909350565b6000616b8282616341565b9050919050565b616b9281616b77565b82525050565b6000602082019050616bad6000830184616b89565b92915050565b60008060408385031215616bca57616bc9615929565b5b6000616bd885828601615a96565b925050602083013567ffffffffffffffff811115616bf957616bf861592e565b5b616c0585828601616592565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b616c44816159e5565b82525050565b6000616c568383616c3b565b60208301905092915050565b6000602082019050919050565b6000616c7a82616c0f565b616c848185616c1a565b9350616c8f83616c2b565b8060005b83811015616cc0578151616ca78882616c4a565b9750616cb283616c62565b925050600181019050616c93565b5085935050505092915050565b60006020820190508181036000830152616ce78184616c6f565b905092915050565b6000616cfa82616341565b9050919050565b616d0a81616cef565b82525050565b6000602082019050616d256000830184616d01565b92915050565b60008060208385031215616d4257616d41615929565b5b600083013567ffffffffffffffff811115616d6057616d5f61592e565b5b616d6c85828601615d25565b92509250509250929050565b6000616d8382616341565b9050919050565b616d9381616d78565b82525050565b6000602082019050616dae6000830184616d8a565b92915050565b600060208284031215616dca57616dc9615929565b5b6000616dd884828501615ac2565b91505092915050565b60008060408385031215616df857616df7615929565b5b6000616e0685828601615a96565b9250506020616e1785828601615ac2565b9150509250929050565b600060208284031215616e3757616e36615929565b5b6000616e4584828501615e6f565b91505092915050565b616e5781615a6d565b82525050565b616e66816164c1565b82525050565b606082016000820151616e826000850182616e4e565b506020820151616e956020850182616e4e565b506040820151616ea86040850182616e5d565b50505050565b6000606082019050616ec36000830184616e6c565b92915050565b60008060008060808587031215616ee357616ee2615929565b5b6000616ef187828801615a96565b9450506020616f0287828801615aee565b9350506040616f1387828801615a96565b9250506060616f2487828801615aee565b91505092959194509250565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b616f6581616cef565b82525050565b6000616f778383616f5c565b60208301905092915050565b6000602082019050919050565b6000616f9b82616f30565b616fa58185616f3b565b9350616fb083616f4c565b8060005b83811015616fe1578151616fc88882616f6b565b9750616fd383616f83565b925050600181019050616fb4565b5085935050505092915050565b600060408201905081810360008301526170088185616f90565b9050818103602083015261701c8184616c6f565b90509392505050565b60008060006060848603121561703e5761703d615929565b5b600061704c86828701615a96565b935050602084013567ffffffffffffffff81111561706d5761706c61592e565b5b61707986828701616a58565b925050604061708a86828701615ac2565b9150509250925092565b6000606082840312156170aa576170a9615929565b5b60006170b884828501615d06565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061712a826159e5565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561715d5761715c6170f0565b5b600182019050919050565b600060c08201905061717d6000830189615a23565b61718a602083018861638f565b617197604083018761638f565b6171a4606083018661638f565b6171b16080830185615a23565b6171be60a08301846159ef565b979650505050505050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b600061720a6002836171c9565b9150617215826171d4565b600282019050919050565b6000819050919050565b61723b61723682615a19565b617220565b82525050565b600061724c826171fd565b9150617258828561722a565b602082019150617268828461722a565b6020820191508190509392505050565b600082825260208201905092915050565b7f5061757361626c653a20696e6465782069732070617573656400000000000000600082015250565b60006172bf601983617278565b91506172ca82617289565b602082019050919050565b600060208201905081810360008301526172ee816172b2565b9050919050565b600080fd5b600080fd5b600080fd5b6000823560016060038336030381126173205761731f6172f5565b5b80830191505092915050565b60008083356001602003843603038112617349576173486172f5565b5b80840192508235915067ffffffffffffffff82111561736b5761736a6172fa565b5b602083019250602082023603831315617387576173866172ff565b5b509250929050565b600080833560016020038436030381126173ac576173ab6172f5565b5b80840192508235915067ffffffffffffffff8211156173ce576173cd6172fa565b5b6020830192506020820236038313156173ea576173e96172ff565b5b509250929050565b7f44656c65676174696f6e4d616e616765722e717565756557697468647261776160008201527f6c3a20696e707574206c656e677468206d69736d617463680000000000000000602082015250565b600061744e603883617278565b9150617459826173f2565b604082019050919050565b6000602082019050818103600083015261747d81617441565b9050919050565b7f44656c65676174696f6e4d616e616765722e717565756557697468647261776160008201527f6c3a2077697468647261776572206d757374206265207374616b657200000000602082015250565b60006174e0603c83617278565b91506174eb82617484565b604082019050919050565b6000602082019050818103600083015261750f816174d3565b9050919050565b7f44656c65676174696f6e4d616e616765722e726567697374657241734f70657260008201527f61746f723a2063616c6c657220697320616c7265616479206163746976656c7960208201527f2064656c65676174656400000000000000000000000000000000000000000000604082015250565b6000617598604a83617278565b91506175a382617516565b606082019050919050565b600060208201905081810360008301526175c78161758b565b9050919050565b60006175dd6020840184615a96565b905092915050565b60006175f460208401846164e8565b905092915050565b6060820161760d60008301836175ce565b61761a6000850182616e4e565b5061762860208301836175ce565b6176356020850182616e4e565b5061764360408301836175e5565b6176506040850182616e5d565b50505050565b600060608201905061766b60008301846175fc565b92915050565b600061767d8385617278565b935061768a8385846169d9565b61769383616430565b840190509392505050565b600060208201905081810360008301526176b9818486617671565b90509392505050565b6000815190506176d181615a7f565b92915050565b6000602082840312156176ed576176ec615929565b5b60006176fb848285016176c2565b91505092915050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320756e70617573657200000000000000000000000000000000000000000000602082015250565b6000617760602a83617278565b915061776b82617704565b604082019050919050565b6000602082019050818103600083015261778f81617753565b9050919050565b7f44656c65676174696f6e4d616e616765723a206f6e6c7953747261746567794d60008201527f616e616765724f72456967656e506f644d616e61676572000000000000000000602082015250565b60006177f2603783617278565b91506177fd82617796565b604082019050919050565b60006020820190508181036000830152617821816177e5565b9050919050565b60008151905061783781616893565b92915050565b60006020828403121561785357617852615929565b5b600061786184828501617828565b91505092915050565b7f6d73672e73656e646572206973206e6f74207065726d697373696f6e6564206160008201527f7320706175736572000000000000000000000000000000000000000000000000602082015250565b60006178c6602883617278565b91506178d18261786a565b604082019050919050565b600060208201905081810360008301526178f5816178b9565b9050919050565b7f5061757361626c652e70617573653a20696e76616c696420617474656d70742060008201527f746f20756e70617573652066756e6374696f6e616c6974790000000000000000602082015250565b6000617958603883617278565b9150617963826178fc565b604082019050919050565b600060208201905081810360008301526179878161794b565b9050919050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b60006179ea602e83617278565b91506179f58261798e565b604082019050919050565b60006020820190508181036000830152617a19816179dd565b9050919050565b6000819050919050565b6000617a45617a40617a3b84617a20565b616315565b6167b8565b9050919050565b617a5581617a2a565b82525050565b6000602082019050617a706000830184617a4c565b92915050565b7f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00600082015250565b6000617aac601f83617278565b9150617ab782617a76565b602082019050919050565b60006020820190508181036000830152617adb81617a9f565b9050919050565b60008235600160e003833603038112617afe57617afd6172f5565b5b80830191505092915050565b60008083356001602003843603038112617b2757617b266172f5565b5b80840192508235915067ffffffffffffffff821115617b4957617b486172fa565b5b602083019250602082023603831315617b6557617b646172ff565b5b509250929050565b600060208284031215617b8357617b82615929565b5b6000617b91848285016168aa565b91505092915050565b600082825260208201905092915050565b6000617bb682616f30565b617bc08185617b9a565b9350617bcb83616f4c565b8060005b83811015617bfc578151617be38882616f6b565b9750617bee83616f83565b925050600181019050617bcf565b5085935050505092915050565b600082825260208201905092915050565b6000617c2582616c0f565b617c2f8185617c09565b9350617c3a83616c2b565b8060005b83811015617c6b578151617c528882616c4a565b9750617c5d83616c62565b925050600181019050617c3e565b5085935050505092915050565b600060e083016000830151617c906000860182616e4e565b506020830151617ca36020860182616e4e565b506040830151617cb66040860182616e4e565b506060830151617cc96060860182616c3b565b506080830151617cdc6080860182616e5d565b5060a083015184820360a0860152617cf48282617bab565b91505060c083015184820360c0860152617d0e8282617c1a565b9150508091505092915050565b60006020820190508181036000830152617d358184617c78565b905092915050565b7f44656c65676174696f6e4d616e616765722e64656c6567617465546f4279536960008201527f676e61747572653a207374616b6572207369676e61747572652065787069726560208201527f6400000000000000000000000000000000000000000000000000000000000000604082015250565b6000617dbf604183617278565b9150617dca82617d3d565b606082019050919050565b60006020820190508181036000830152617dee81617db2565b9050919050565b7f44656c65676174696f6e4d616e616765722e64656c6567617465546f4279536960008201527f676e61747572653a207374616b657220697320616c726561647920616374697660208201527f656c792064656c65676174656400000000000000000000000000000000000000604082015250565b6000617e77604d83617278565b9150617e8282617df5565b606082019050919050565b60006020820190508181036000830152617ea681617e6a565b9050919050565b7f44656c65676174696f6e4d616e616765722e64656c6567617465546f4279536960008201527f676e61747572653a206f70657261746f72206973206e6f74207265676973746560208201527f72656420696e20456967656e4c61796572000000000000000000000000000000604082015250565b6000617f2f605183617278565b9150617f3a82617ead565b606082019050919050565b60006020820190508181036000830152617f5e81617f22565b9050919050565b7f44656c65676174696f6e4d616e616765722e7570646174654f70657261746f7260008201527f4d657461646174615552493a2063616c6c6572206d75737420626520616e206f60208201527f70657261746f7200000000000000000000000000000000000000000000000000604082015250565b6000617fe7604783617278565b9150617ff282617f65565b606082019050919050565b6000602082019050818103600083015261801681617fda565b9050919050565b600060a0820190506180326000830188615a23565b61803f602083018761638f565b61804c604083018661638f565b61805960608301856159ef565b61806660808301846159ef565b9695505050505050565b6000819050919050565b61808381618070565b811461808e57600080fd5b50565b6000815190506180a08161807a565b92915050565b6000602082840312156180bc576180bb615929565b5b60006180ca84828501618091565b91505092915050565b6000815190506180e281615e58565b92915050565b60006180fb6180f6846164fd565b6164a1565b9050808382526020820190506020840283018581111561811e5761811d61593d565b5b835b81811015618147578061813388826180d3565b845260208401935050602081019050618120565b5050509392505050565b600082601f83011261816657618165615933565b5b81516181768482602086016180e8565b91505092915050565b60008151905061818e81615ad7565b92915050565b60006181a76181a2846165c0565b6164a1565b905080838252602082019050602084028301858111156181ca576181c961593d565b5b835b818110156181f357806181df888261817f565b8452602084019350506020810190506181cc565b5050509392505050565b600082601f83011261821257618211615933565b5b8151618222848260208601618194565b91505092915050565b6000806040838503121561824257618241615929565b5b600083015167ffffffffffffffff8111156182605761825f61592e565b5b61826c85828601618151565b925050602083015167ffffffffffffffff81111561828d5761828c61592e565b5b618299858286016181fd565b9150509250929050565b60006182ae826159e5565b91506182b9836159e5565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156182ee576182ed6170f0565b5b828201905092915050565b6000618304826159e5565b915061830f836159e5565b925082821015618322576183216170f0565b5b828203905092915050565b7f44656c65676174696f6e4d616e616765722e756e64656c65676174653a20737460008201527f616b6572206d7573742062652064656c65676174656420746f20756e64656c6560208201527f6761746500000000000000000000000000000000000000000000000000000000604082015250565b60006183af604483617278565b91506183ba8261832d565b606082019050919050565b600060208201905081810360008301526183de816183a2565b9050919050565b7f44656c65676174696f6e4d616e616765722e756e64656c65676174653a206f7060008201527f657261746f72732063616e6e6f7420626520756e64656c656761746564000000602082015250565b6000618441603d83617278565b915061844c826183e5565b604082019050919050565b6000602082019050818103600083015261847081618434565b9050919050565b7f44656c65676174696f6e4d616e616765722e756e64656c65676174653a20636160008201527f6e6e6f7420756e64656c6567617465207a65726f206164647265737300000000602082015250565b60006184d3603c83617278565b91506184de82618477565b604082019050919050565b60006020820190508181036000830152618502816184c6565b9050919050565b7f44656c65676174696f6e4d616e616765722e756e64656c65676174653a20636160008201527f6c6c65722063616e6e6f7420756e64656c6567617465207374616b6572000000602082015250565b6000618565603d83617278565b915061857082618509565b604082019050919050565b6000602082019050818103600083015261859481618558565b9050919050565b7f44656c65676174696f6e4d616e616765722e64656c6567617465546f3a20737460008201527f616b657220697320616c7265616479206163746976656c792064656c6567617460208201527f6564000000000000000000000000000000000000000000000000000000000000604082015250565b600061861d604283617278565b91506186288261859b565b606082019050919050565b6000602082019050818103600083015261864c81618610565b9050919050565b7f44656c65676174696f6e4d616e616765722e64656c6567617465546f3a206f7060008201527f657261746f72206973206e6f74207265676973746572656420696e204569676560208201527f6e4c617965720000000000000000000000000000000000000000000000000000604082015250565b60006186d5604683617278565b91506186e082618653565b606082019050919050565b60006020820190508181036000830152618704816186c8565b9050919050565b7f44656c65676174696f6e4d616e616765722e6d6f646966794f70657261746f7260008201527f44657461696c733a2063616c6c6572206d75737420626520616e206f7065726160208201527f746f720000000000000000000000000000000000000000000000000000000000604082015250565b600061878d604383617278565b91506187988261870b565b606082019050919050565b600060208201905081810360008301526187bc81618780565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061881f602683617278565b915061882a826187c3565b604082019050919050565b6000602082019050818103600083015261884e81618812565b9050919050565b7f5061757361626c652e756e70617573653a20696e76616c696420617474656d7060008201527f7420746f2070617573652066756e6374696f6e616c6974790000000000000000602082015250565b60006188b1603883617278565b91506188bc82618855565b604082019050919050565b600060208201905081810360008301526188e0816188a4565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160008201527f6e6451756575655769746864726177616c3a207374616b65722063616e6e6f7460208201527f206265207a65726f206164647265737300000000000000000000000000000000604082015250565b6000618969605083617278565b9150618974826188e7565b606082019050919050565b600060208201905081810360008301526189988161895c565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160008201527f6e6451756575655769746864726177616c3a207374726174656769657320636160208201527f6e6e6f7420626520656d70747900000000000000000000000000000000000000604082015250565b6000618a21604d83617278565b9150618a2c8261899f565b606082019050919050565b60006020820190508181036000830152618a5081618a14565b9050919050565b6000604082019050618a6c600083018561638f565b618a7960208301846159ef565b9392505050565b7f44656c65676174696f6e4d616e616765722e5f72656d6f76655368617265734160008201527f6e6451756575655769746864726177616c3a2077697468647261776572206d7560208201527f73742062652073616d652061646472657373206173207374616b65722069662060408201527f746869726450617274795472616e7366657273466f7262696464656e2061726560608201527f2073657400000000000000000000000000000000000000000000000000000000608082015250565b6000618b4e608483617278565b9150618b5982618a80565b60a082019050919050565b60006020820190508181036000830152618b7d81618b41565b9050919050565b6000606082019050618b99600083018661638f565b618ba66020830185616d01565b618bb360408301846159ef565b949350505050565b6000604082019050618bd06000830185615a23565b8181036020830152618be28184617c78565b90509392505050565b600060208284031215618c0157618c00615929565b5b6000618c0f848285016164e8565b91505092915050565b7f44656c65676174696f6e4d616e616765722e5f7365744f70657261746f72446560008201527f7461696c733a207374616b65724f70744f757457696e646f77426c6f636b732060208201527f63616e6e6f74206265203e204d41585f5354414b45525f4f50545f4f55545f5760408201527f494e444f575f424c4f434b530000000000000000000000000000000000000000606082015250565b6000618cc0606c83617278565b9150618ccb82618c18565b608082019050919050565b60006020820190508181036000830152618cef81618cb3565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f7365744f70657261746f72446560008201527f7461696c733a207374616b65724f70744f757457696e646f77426c6f636b732060208201527f63616e6e6f742062652064656372656173656400000000000000000000000000604082015250565b6000618d78605383617278565b9150618d8382618cf6565b606082019050919050565b60006020820190508181036000830152618da781618d6b565b9050919050565b60008135618dbb81615a7f565b80915050919050565b60008160001b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff618df184618dc4565b9350801983169250808416831791505092915050565b6000618e1282616341565b9050919050565b6000819050919050565b618e2c82618e07565b618e3f618e3882618e19565b8354618dd1565b8255505050565b60008135618e53816164d1565b80915050919050565b60008160a01b9050919050565b600077ffffffff0000000000000000000000000000000000000000618e8d84618e5c565b9350801983169250808416831791505092915050565b6000618ebe618eb9618eb4846164c1565b616315565b6164c1565b9050919050565b6000819050919050565b618ed882618ea3565b618eeb618ee482618ec5565b8354618e69565b8255505050565b600081016000830180618f0481618dae565b9050618f108184618e23565b505050600181016020830180618f2581618dae565b9050618f318184618e23565b505050600181016040830180618f4681618e46565b9050618f528184618ecf565b5050505050565b618f638282618ef2565b5050565b7f44656c65676174696f6e4d616e616765722e5f64656c65676174653a2061707060008201527f726f766572207369676e61747572652065787069726564000000000000000000602082015250565b6000618fc3603783617278565b9150618fce82618f67565b604082019050919050565b60006020820190508181036000830152618ff281618fb6565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f64656c65676174653a2061707060008201527f726f76657253616c7420616c7265616479207370656e74000000000000000000602082015250565b6000619055603783617278565b915061906082618ff9565b604082019050919050565b6000602082019050818103600083015261908481619048565b9050919050565b7f5061757361626c652e5f73657450617573657252656769737472793a206e657760008201527f50617573657252656769737472792063616e6e6f7420626520746865207a657260208201527f6f20616464726573730000000000000000000000000000000000000000000000604082015250565b600061910d604983617278565b91506191188261908b565b606082019050919050565b6000602082019050818103600083015261913c81619100565b9050919050565b60006040820190506191586000830185616b89565b6191656020830184616b89565b9392505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006191a2602083617278565b91506191ad8261916c565b602082019050919050565b600060208201905081810360008301526191d181619195565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f7365745374726174656779576960008201527f746864726177616c44656c6179426c6f636b733a20696e707574206c656e677460208201527f68206d69736d6174636800000000000000000000000000000000000000000000604082015250565b600061925a604a83617278565b9150619265826191d8565b606082019050919050565b600060208201905081810360008301526192898161924d565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f7365745374726174656779576960008201527f746864726177616c44656c6179426c6f636b733a205f7769746864726177616c60208201527f44656c6179426c6f636b732063616e6e6f74206265203e204d41585f5749544860408201527f44524157414c5f44454c41595f424c4f434b5300000000000000000000000000606082015250565b6000619338607383617278565b915061934382619290565b608082019050919050565b600060208201905081810360008301526193678161932b565b9050919050565b60006060820190506193836000830186616d01565b61939060208301856159ef565b61939d60408301846159ef565b949350505050565b7f5061757361626c652e5f696e697469616c697a655061757365723a205f696e6960008201527f7469616c697a6550617573657228292063616e206f6e6c792062652063616c6c60208201527f6564206f6e636500000000000000000000000000000000000000000000000000604082015250565b6000619427604783617278565b9150619432826193a5565b606082019050919050565b600060208201905081810360008301526194568161941a565b9050919050565b60006080820190506194726000830187615a23565b61947f6020830186615a23565b61948c60408301856159ef565b619499606083018461638f565b95945050505050565b7f44656c65676174696f6e4d616e616765722e5f7365744d696e5769746864726160008201527f77616c44656c6179426c6f636b733a205f6d696e5769746864726177616c446560208201527f6c6179426c6f636b732063616e6e6f74206265203e204d41585f57495448445260408201527f4157414c5f44454c41595f424c4f434b53000000000000000000000000000000606082015250565b600061954a607183617278565b9150619555826194a2565b608082019050919050565b600060208201905081810360008301526195798161953d565b9050919050565b600060408201905061959560008301856159ef565b6195a260208301846159ef565b9392505050565b60006195b53683616683565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f636f6d706c657465517565756560008201527f645769746864726177616c3a20616374696f6e206973206e6f7420696e20717560208201527f6575650000000000000000000000000000000000000000000000000000000000604082015250565b600061963e604383617278565b9150619649826195bc565b606082019050919050565b6000602082019050818103600083015261966d81619631565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f636f6d706c657465517565756560008201527f645769746864726177616c3a206d696e5769746864726177616c44656c61794260208201527f6c6f636b7320706572696f6420686173206e6f74207965742070617373656400604082015250565b60006196f6605f83617278565b915061970182619674565b606082019050919050565b60006020820190508181036000830152619725816196e9565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f636f6d706c657465517565756560008201527f645769746864726177616c3a206f6e6c7920776974686472617765722063616e60208201527f20636f6d706c65746520616374696f6e00000000000000000000000000000000604082015250565b60006197ae605083617278565b91506197b98261972c565b606082019050919050565b600060208201905081810360008301526197dd816197a1565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f636f6d706c657465517565756560008201527f645769746864726177616c3a20696e707574206c656e677468206d69736d617460208201527f6368000000000000000000000000000000000000000000000000000000000000604082015250565b6000619866604283617278565b9150619871826197e4565b606082019050919050565b6000602082019050818103600083015261989581619859565b9050919050565b7f44656c65676174696f6e4d616e616765722e5f636f6d706c657465517565756560008201527f645769746864726177616c3a207769746864726177616c44656c6179426c6f6360208201527f6b7320706572696f6420686173206e6f74207965742070617373656420666f7260408201527f2074686973207374726174656779000000000000000000000000000000000000606082015250565b6000619944606e83617278565b915061994f8261989c565b608082019050919050565b6000602082019050818103600083015261997381619937565b9050919050565b600061998582615a6d565b9050919050565b6199958161997a565b81146199a057600080fd5b50565b6000813590506199b28161998c565b92915050565b6000602082840312156199ce576199cd615929565b5b60006199dc848285016199a3565b91505092915050565b6000602082840312156199fb576199fa615929565b5b6000619a098482850161817f565b91505092915050565b6000619a1d82616341565b9050919050565b619a2d81619a12565b82525050565b6000608082019050619a48600083018761638f565b619a556020830186619a24565b619a626040830185616d01565b619a6f60608301846159ef565b95945050505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015619ab2578082015181840152602081019050619a97565b83811115619ac1576000848401525b50505050565b6000619ad282619a78565b619adc8185619a83565b9350619aec818560208601619a94565b619af581616430565b840191505092915050565b6000604082019050619b156000830185615a23565b8181036020830152619b278184619ac7565b90509392505050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b619b6581619b30565b8114619b7057600080fd5b50565b600081519050619b8281619b5c565b92915050565b600060208284031215619b9e57619b9d615929565b5b6000619bac84828501619b73565b91505092915050565b7f454950313237315369676e61747572655574696c732e636865636b5369676e6160008201527f747572655f454950313237313a2045524331323731207369676e61747572652060208201527f766572696669636174696f6e206661696c656400000000000000000000000000604082015250565b6000619c37605383617278565b9150619c4282619bb5565b606082019050919050565b60006020820190508181036000830152619c6681619c2a565b9050919050565b7f454950313237315369676e61747572655574696c732e636865636b5369676e6160008201527f747572655f454950313237313a207369676e6174757265206e6f742066726f6d60208201527f207369676e657200000000000000000000000000000000000000000000000000604082015250565b6000619cef604783617278565b9150619cfa82619c6d565b606082019050919050565b60006020820190508181036000830152619d1e81619ce2565b9050919050565b6000606082019050619d3a600083018661638f565b619d47602083018561638f565b619d5460408301846159ef565b949350505050565b6000608082019050619d71600083018761638f565b619d7e6020830186616d01565b619d8b60408301856159ef565b619d986060830184619a24565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f45434453413a20696e76616c6964207369676e61747572650000000000000000600082015250565b6000619e06601883617278565b9150619e1182619dd0565b602082019050919050565b60006020820190508181036000830152619e3581619df9565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265206c656e67746800600082015250565b6000619e72601f83617278565b9150619e7d82619e3c565b602082019050919050565b60006020820190508181036000830152619ea181619e65565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202773272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b6000619f04602283617278565b9150619f0f82619ea8565b604082019050919050565b60006020820190508181036000830152619f3381619ef7565b9050919050565b7f45434453413a20696e76616c6964207369676e6174757265202776272076616c60008201527f7565000000000000000000000000000000000000000000000000000000000000602082015250565b6000619f96602283617278565b9150619fa182619f3a565b604082019050919050565b60006020820190508181036000830152619fc581619f89565b9050919050565b619fd5816167b8565b82525050565b6000608082019050619ff06000830187615a23565b619ffd6020830186619fcc565b61a00a6040830185615a23565b61a0176060830184615a23565b9594505050505056fea26469706673582212206f0e2a8551a07a0f0a56654405cd27bdb6b812846823148e5259b0710a081b6e64736f6c634300080c0033",
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

// SetMinWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x635bbd10.
//
// Solidity: function setMinWithdrawalDelayBlocks(uint256 newMinWithdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactor) SetMinWithdrawalDelayBlocks(opts *bind.TransactOpts, newMinWithdrawalDelayBlocks *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.contract.Transact(opts, "setMinWithdrawalDelayBlocks", newMinWithdrawalDelayBlocks)
}

// SetMinWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x635bbd10.
//
// Solidity: function setMinWithdrawalDelayBlocks(uint256 newMinWithdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerSession) SetMinWithdrawalDelayBlocks(newMinWithdrawalDelayBlocks *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetMinWithdrawalDelayBlocks(&_ContractDelegationManager.TransactOpts, newMinWithdrawalDelayBlocks)
}

// SetMinWithdrawalDelayBlocks is a paid mutator transaction binding the contract method 0x635bbd10.
//
// Solidity: function setMinWithdrawalDelayBlocks(uint256 newMinWithdrawalDelayBlocks) returns()
func (_ContractDelegationManager *ContractDelegationManagerTransactorSession) SetMinWithdrawalDelayBlocks(newMinWithdrawalDelayBlocks *big.Int) (*types.Transaction, error) {
	return _ContractDelegationManager.Contract.SetMinWithdrawalDelayBlocks(&_ContractDelegationManager.TransactOpts, newMinWithdrawalDelayBlocks)
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
