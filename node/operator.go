package node

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"
	"github.com/ethereum/go-ethereum/crypto"
)

type Operator struct {
	Address             string
	Socket              string
	Timeout             time.Duration
	PrivKey             *ecdsa.PrivateKey
	Signer              blssigner.Signer
	OperatorId          core.OperatorID
	QuorumIDs           []core.QuorumID
	RegisterNodeAtStart bool
}

// RegisterOperator operator registers the operator with the given public key for the given quorum IDs.
func RegisterOperator(ctx context.Context, operator *Operator, transactor core.Writer, logger logging.Logger) error {
	if len(operator.QuorumIDs) > 1+core.MaxQuorumID {
		return fmt.Errorf("cannot provide more than %d quorums", 1+core.MaxQuorumID)
	}
	quorumsToRegister, err := operator.getQuorumIdsToRegister(ctx, transactor)
	if err != nil {
		return fmt.Errorf("failed to get quorum ids to register: %w", err)
	}
	if !operator.RegisterNodeAtStart {
		// For operator-initiated registration, the supplied quorums must be not registered yet.
		if len(quorumsToRegister) != len(operator.QuorumIDs) {
			return errors.New("quorums to register must be not registered yet")
		}
	}
	if len(quorumsToRegister) == 0 {
		return nil
	}

	logger.Info("Quorums to register for", "quorums", fmt.Sprint(quorumsToRegister)) //nolint:staticcheck // printing byte slices is fine here

	// Generate salt and expiry
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return err
	}
	salt := [32]byte{}
	copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), quorumsToRegister, bytes))

	// Get the current block number
	expiry := big.NewInt((time.Now().Add(10 * time.Minute)).Unix())

	err = transactor.RegisterOperator(
		ctx,
		operator.Signer,
		operator.Socket,
		quorumsToRegister,
		operator.PrivKey,
		salt,
		expiry)
	if err != nil {
		return fmt.Errorf("failed to register operator: %w", err)
	}
	return nil
}

// DeregisterOperator deregisters the operator with the given public key from the specified quorums that it is registered with at the supplied block number.
// If the operator isn't registered with any of the specified quorums, this function will return error, and no quorum will be deregistered.
func DeregisterOperator(ctx context.Context, operator *Operator, pubKeyG1 *core.G1Point, transactor core.Writer) error {
	if len(operator.QuorumIDs) > 1+core.MaxQuorumID {
		return fmt.Errorf("cannot provide more than %d quorums", 1+core.MaxQuorumID)
	}
	blockNumber, err := transactor.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	return transactor.DeregisterOperator(ctx, pubKeyG1, blockNumber, operator.QuorumIDs)
}

// UpdateOperatorSocket updates the socket for the given operator
func UpdateOperatorSocket(ctx context.Context, transactor core.Writer, socket string) error {
	return transactor.UpdateOperatorSocket(ctx, socket)
}

// getQuorumIdsToRegister returns the quorum ids that the operator is not registered in.
func (c *Operator) getQuorumIdsToRegister(ctx context.Context, transactor core.Writer) ([]core.QuorumID, error) {
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
		}
	}

	return quorumIdsToRegister, nil
}
