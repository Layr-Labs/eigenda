package corev2

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/core/types"
)

type Writer interface {
	Reader

	// RegisterOperator registers a new operator with the given public key and socket with the provided quorum ids.
	// If the operator is already registered with a given quorum id, the transaction will fail (noop) and an error
	// will be returned.
	RegisterOperator(
		ctx context.Context,
		keypair *bn254.KeyPair,
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
		keypair *bn254.KeyPair,
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
	DeregisterOperator(ctx context.Context, pubkeyG1 *bn254.G1Point, blockNumber uint32, quorumIds []QuorumID) error

	// UpdateOperatorSocket updates the socket of the operator in all the quorums that it is registered with.
	UpdateOperatorSocket(ctx context.Context, socket string) error

	// BuildEjectOperatorsTxn returns a transaction that ejects operators from AVS registryCoordinator.
	// The operatorsByQuorum provides a list of operators for each quorum. Within a quorum,
	// the operators are ordered; in case of rate limiting, the first operators will be ejected.
	BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][]OperatorID) (*types.Transaction, error)
}
