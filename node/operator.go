package node

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/crypto"
)

type Operator struct {
	Address    string
	Socket     string
	Timeout    time.Duration
	PrivKey    *ecdsa.PrivateKey
	KeyPair    *core.KeyPair
	OperatorId core.OperatorID
	QuorumIDs  []core.QuorumID
}

// RegisterOperator operator registers the operator with the given public key for the given quorum IDs.
func RegisterOperator(ctx context.Context, operator *Operator, transactor core.Transactor, churnerClient ChurnerClient, logger logging.Logger) error {
	quorumsToRegister, err := operator.getQuorumIdsToRegister(ctx, transactor)
	if err != nil {
		return fmt.Errorf("failed to get quorum ids to register: %w", err)
	}
	if len(quorumsToRegister) == 0 {
		return nil
	}

	logger.Info("Quorums to register for", "quorums", quorumsToRegister)

	// register for quorums
	shouldCallChurner := false
	// check if one of the quorums to register for is full
	for _, quorumID := range quorumsToRegister {
		operatorSetParams, err := transactor.GetOperatorSetParams(ctx, quorumID)
		if err != nil {
			return err
		}

		numberOfRegisteredOperators, err := transactor.GetNumberOfRegisteredOperatorForQuorum(ctx, quorumID)
		if err != nil {
			return err
		}

		// if the quorum is full, we need to call the churner
		if operatorSetParams.MaxOperatorCount == numberOfRegisteredOperators {
			shouldCallChurner = true
			break
		}
	}

	logger.Info("Should call churner", "shouldCallChurner", shouldCallChurner)

	// Generate salt and expiry

	privateKeyBytes := []byte(operator.KeyPair.PrivKey.String())
	salt := [32]byte{}
	copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), quorumsToRegister, privateKeyBytes))

	// Get the current block number
	expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())

	// if we should call the churner, call it
	if shouldCallChurner {
		churnReply, err := churnerClient.Churn(ctx, operator.Address, operator.KeyPair, quorumsToRegister)
		if err != nil {
			return fmt.Errorf("failed to request churn approval: %w", err)
		}

		return transactor.RegisterOperatorWithChurn(ctx, operator.KeyPair, operator.Socket, quorumsToRegister, operator.PrivKey, salt, expiry, churnReply)
	} else {
		// other wise just register normally
		return transactor.RegisterOperator(ctx, operator.KeyPair, operator.Socket, quorumsToRegister, operator.PrivKey, salt, expiry)
	}
}

// DeregisterOperator deregisters the operator with the given public key from the specified quorums that it is registered with at the supplied block number.
// If the operator isn't registered with any of the specified quorums, this function will return error, and no quorum will be deregistered.
func DeregisterOperator(ctx context.Context, operator *Operator, KeyPair *core.KeyPair, transactor core.Transactor) error {
	blockNumber, err := transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	return transactor.DeregisterOperator(ctx, KeyPair.GetPubKeyG1(), blockNumber, operator.QuorumIDs)
}

// UpdateOperatorSocket updates the socket for the given operator
func UpdateOperatorSocket(ctx context.Context, transactor core.Transactor, socket string) error {
	return transactor.UpdateOperatorSocket(ctx, socket)
}

// getQuorumIdsToRegister returns the quorum ids that the operator is not registered in.
func (c *Operator) getQuorumIdsToRegister(ctx context.Context, transactor core.Transactor) ([]core.QuorumID, error) {
	if len(c.QuorumIDs) == 0 {
		return nil, fmt.Errorf("an operator should be in at least one quorum to be useful")
	}

	registeredQuorumIds, err := transactor.GetRegisteredQuorumIdsForOperator(ctx, c.OperatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered quorum ids for an operator: %w", err)
	}

	quorumIdsToRegister := make([]core.QuorumID, 0, len(c.QuorumIDs))
	for _, quorumID := range c.QuorumIDs {
		if !slices.Contains(registeredQuorumIds, quorumID) {
			quorumIdsToRegister = append(quorumIdsToRegister, quorumID)
		} else {
			return nil, fmt.Errorf("the operator already registered for quorum %d", quorumID)
		}
	}

	return quorumIdsToRegister, nil
}
