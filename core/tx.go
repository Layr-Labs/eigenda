package core

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type OperatorStake struct {
	OperatorID OperatorID
	Stake      *big.Int
}

type OperatorToChurn struct {
	QuorumId QuorumID
	Operator gethcommon.Address
	Pubkey   *G1Point
}

type OperatorSetParam struct {
	MaxOperatorCount         uint32
	ChurnBIPsOfOperatorStake uint16
	ChurnBIPsOfTotalStake    uint16
}

type OperatorStakes map[QuorumID]map[OperatorIndex]OperatorStake

type Transactor interface {

	// GetRegisteredQuorumIdsForOperator returns the quorum ids that the operator is registered in with the given public key.
	GetRegisteredQuorumIdsForOperator(ctx context.Context, operatorID OperatorID) ([]QuorumID, error)

	// RegisterOperator registers a new operator with the given public key and socket with the provided quorum ids.
	// If the operator is already registered with a given quorum id, the transaction will fail (noop) and an error
	// will be returned.
	RegisterOperator(ctx context.Context, keypair *KeyPair, socket string, quorumIds []QuorumID) error

	// RegisterOperatorWithChurn registers a new operator with the given public key and socket with the provided quorum ids
	// with the provided signature from the churner
	RegisterOperatorWithChurn(ctx context.Context, keypair *KeyPair, socket string, quorumIds []QuorumID, churnReply *churner.ChurnReply) error

	// DeregisterOperator deregisters an operator with the given public key from the all the quorums that it is
	// registered with at the supplied block number. To fully deregister an operator, this function should be called
	// with the current block number.
	DeregisterOperator(ctx context.Context, pubkeyG1 *G1Point, blockNumber uint32) error

	// UpdateOperatorSocket updates the socket of the operator in all the quorums that it is registered with.
	UpdateOperatorSocket(ctx context.Context, socket string) error

	// GetOperatorStakes returns the stakes of all operators within the quorums that the operator represented by operatorId
	//  is registered with. The returned stakes are for the block number supplied. The indices of the operators within each quorum
	// are also returned.
	GetOperatorStakes(ctx context.Context, operatorID OperatorID, blockNumber uint32) (OperatorStakes, []QuorumID, error)

	// GetOperatorStakes returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
	// The indices of the operators within each quorum are also returned.
	GetOperatorStakesForQuorums(ctx context.Context, quorums []QuorumID, blockNumber uint32) (OperatorStakes, error)

	// BuildConfirmBatchTxn builds a transaction to confirm a batch header and signature aggregation.
	BuildConfirmBatchTxn(ctx context.Context, batchHeader *BatchHeader, quorums map[QuorumID]*QuorumResult, signatureAggregation *SignatureAggregation) (*types.Transaction, error)

	// ConfirmBatch confirms a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
	// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
	ConfirmBatch(ctx context.Context, batchHeader *BatchHeader, quorums map[QuorumID]*QuorumResult, signatureAggregation *SignatureAggregation) (*types.Receipt, error)

	// GetBlockStaleMeasure returns the BLOCK_STALE_MEASURE defined onchain.
	GetBlockStaleMeasure(ctx context.Context) (uint32, error)
	// GetStoreDurationBlocks returns the STORE_DURATION_BLOCKS defined onchain.
	GetStoreDurationBlocks(ctx context.Context) (uint32, error)

	// StakeRegistry returns the address of the stake registry contract.
	StakeRegistry(ctx context.Context) (gethcommon.Address, error)

	// OperatorIDToAddress returns the address of the operator from the operator id.
	OperatorIDToAddress(ctx context.Context, operatorId OperatorID) (gethcommon.Address, error)

	// GetCurrentQuorumBitmapByOperatorId returns the current quorum bitmap for the operator.
	GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId OperatorID) (*big.Int, error)

	// GetOperatorSetParams returns operator set params for the quorum.
	GetOperatorSetParams(ctx context.Context, quorumID QuorumID) (*OperatorSetParam, error)

	// GetNumberOfRegisteredOperatorForQuorum returns the number of registered operators for the quorum.
	GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID QuorumID) (uint32, error)

	// WeightOfOperatorForQuorum returns the weight of the operator for the quorum view.
	WeightOfOperatorForQuorum(ctx context.Context, quorumID QuorumID, operator gethcommon.Address) (*big.Int, error)

	// CalculateOperatorChurnApprovalDigestHash returns calculated operator churn approval digest hash.
	CalculateOperatorChurnApprovalDigestHash(
		ctx context.Context,
		operatorId OperatorID,
		operatorsToChurn []OperatorToChurn,
		salt [32]byte,
		expiry *big.Int,
	) ([32]byte, error)

	// GetCurrentBlockNumber returns the current block number.
	GetCurrentBlockNumber(ctx context.Context) (uint32, error)

	// GetQuorumCount returns the number of quorums registered at given block number.
	GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error)
}
