package core

import (
	"context"
	"crypto/ecdsa"
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

type Reader interface {

	// GetRegisteredQuorumIdsForOperator returns the quorum ids that the operator is registered in with the given public key.
	GetRegisteredQuorumIdsForOperator(ctx context.Context, operatorID OperatorID) ([]QuorumID, error)

	// GetOperatorStakes returns the stakes of all operators within the quorums that the operator represented by operatorId
	//  is registered with. The returned stakes are for the block number supplied. The indices of the operators within each quorum
	// are also returned.
	GetOperatorStakes(ctx context.Context, operatorID OperatorID, blockNumber uint32) (OperatorStakes, []QuorumID, error)

	// GetOperatorStakesForQuorums returns the stakes of all operators within the supplied quorums. The returned stakes are for the block number supplied.
	// The indices of the operators within each quorum are also returned.
	GetOperatorStakesForQuorums(ctx context.Context, quorums []QuorumID, blockNumber uint32) (OperatorStakes, error)

	// GetBlockStaleMeasure returns the BLOCK_STALE_MEASURE defined onchain.
	GetBlockStaleMeasure(ctx context.Context) (uint32, error)

	// GetStoreDurationBlocks returns the STORE_DURATION_BLOCKS defined onchain.
	GetStoreDurationBlocks(ctx context.Context) (uint32, error)

	// StakeRegistry returns the address of the stake registry contract.
	StakeRegistry(ctx context.Context) (gethcommon.Address, error)

	// OperatorIDToAddress returns the address of the operator from the operator id.
	OperatorIDToAddress(ctx context.Context, operatorId OperatorID) (gethcommon.Address, error)

	// OperatorAddressToID returns the operator id from the operator address.
	OperatorAddressToID(ctx context.Context, operatorAddress gethcommon.Address) (OperatorID, error)

	// BatchOperatorIDToAddress returns the addresses of the operators from the operator id.
	BatchOperatorIDToAddress(ctx context.Context, operatorIds []OperatorID) ([]gethcommon.Address, error)

	// GetCurrentQuorumBitmapByOperatorId returns the current quorum bitmap for the operator.
	GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId OperatorID) (*big.Int, error)

	// GetQuorumBitmapForOperatorsAtBlockNumber returns the quorum bitmaps for the operators at the given block number.
	// The result slice will be of same length as "operatorIds", with the i-th entry be the result for the operatorIds[i].
	// If an operator failed to find bitmap, the corresponding result entry will be an empty bitmap.
	GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds []OperatorID, blockNumber uint32) ([]*big.Int, error)

	// GetOperatorSetParams returns operator set params for the quorum.
	GetOperatorSetParams(ctx context.Context, quorumID QuorumID) (*OperatorSetParam, error)

	// GetNumberOfRegisteredOperatorForQuorum returns the number of registered operators for the quorum.
	GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID QuorumID) (uint32, error)

	// WeightOfOperatorForQuorum returns the weight of the operator for the quorum view.
	WeightOfOperatorForQuorum(ctx context.Context, quorumID QuorumID, operator gethcommon.Address) (*big.Int, error)

	// CalculateOperatorChurnApprovalDigestHash returns calculated operator churn approval digest hash.
	CalculateOperatorChurnApprovalDigestHash(
		ctx context.Context,
		operatorAddress gethcommon.Address,
		operatorId OperatorID,
		operatorsToChurn []OperatorToChurn,
		salt [32]byte,
		expiry *big.Int,
	) ([32]byte, error)

	// GetCurrentBlockNumber returns the current block number.
	GetCurrentBlockNumber(ctx context.Context) (uint32, error)

	// GetQuorumCount returns the number of quorums registered at given block number.
	GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error)

	// GetQuorumSecurityParams returns the security params for the registered quorums.
	GetQuorumSecurityParams(ctx context.Context, blockNumber uint32) ([]SecurityParam, error)

	// GetRequiredQuorumNumbers returns set of required quorum numbers
	GetRequiredQuorumNumbers(ctx context.Context, blockNumber uint32) ([]QuorumID, error)

	// GetActiveReservations returns active reservations (end timestamp > current timestamp)
	GetActiveReservations(ctx context.Context, blockNumber uint32, accountIDs []string) (map[string]ActiveReservation, error)

	// GetActiveReservationByAccount returns active reservation by account ID
	GetActiveReservationByAccount(ctx context.Context, blockNumber uint32, accountID string) (ActiveReservation, error)

	// GetOnDemandPayments returns all on-demand payments
	GetOnDemandPayments(ctx context.Context, blockNumber uint32, accountIDs []string) (map[string]OnDemandPayment, error)

	// GetOnDemandPaymentByAccount returns on-demand payment of an account
	GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint32, accountID string) (OnDemandPayment, error)
}

type Writer interface {
	Reader

	// RegisterOperator registers a new operator with the given public key and socket with the provided quorum ids.
	// If the operator is already registered with a given quorum id, the transaction will fail (noop) and an error
	// will be returned.
	RegisterOperator(
		ctx context.Context,
		keypair *KeyPair,
		socket string,
		quorumIds []QuorumID,
		operatorEcdsaPrivateKey *ecdsa.PrivateKey,
		operatorToAvsRegistrationSigSalt [32]byte,
		operatorToAvsRegistrationSigExpiry *big.Int,
	) error

	// RegisterOperatorWithChurn registers a new operator with the given public key and socket with the provided quorum ids
	// with the provided signature from the churner
	RegisterOperatorWithChurn(
		ctx context.Context,
		keypair *KeyPair,
		socket string,
		quorumIds []QuorumID,
		operatorEcdsaPrivateKey *ecdsa.PrivateKey,
		operatorToAvsRegistrationSigSalt [32]byte,
		operatorToAvsRegistrationSigExpiry *big.Int,
		churnReply *churner.ChurnReply,
	) error

	// DeregisterOperator deregisters an operator with the given public key from the all the quorums that it is
	// registered with at the supplied block number. To fully deregister an operator, this function should be called
	// with the current block number.
	// If the operator isn't registered with any of the specified quorums, this function will return error, and
	// no quorum will be deregistered.
	DeregisterOperator(ctx context.Context, pubkeyG1 *G1Point, blockNumber uint32, quorumIds []QuorumID) error

	// UpdateOperatorSocket updates the socket of the operator in all the quorums that it is registered with.
	UpdateOperatorSocket(ctx context.Context, socket string) error

	// BuildEjectOperatorsTxn returns a transaction that ejects operators from AVS registryCoordinator.
	// The operatorsByQuorum provides a list of operators for each quorum. Within a quorum,
	// the operators are ordered; in case of rate limiting, the first operators will be ejected.
	BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][]OperatorID) (*types.Transaction, error)

	// BuildConfirmBatchTxn builds a transaction to confirm a batch header and signature aggregation.
	BuildConfirmBatchTxn(ctx context.Context, batchHeader *BatchHeader, quorums map[QuorumID]*QuorumResult, signatureAggregation *SignatureAggregation) (*types.Transaction, error)

	// ConfirmBatch confirms a batch header and signature aggregation. The signature aggregation must satisfy the quorum thresholds
	// specified in the batch header. If the signature aggregation does not satisfy the quorum thresholds, the transaction will fail.
	ConfirmBatch(ctx context.Context, batchHeader *BatchHeader, quorums map[QuorumID]*QuorumResult, signatureAggregation *SignatureAggregation) (*types.Receipt, error)
}
