package mock

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type MockWriter struct {
	MockReader
	mock.Mock
}

var _ chainio.Writer = (*MockWriter)(nil)

func (t *MockWriter) RegisterOperator(
	ctx context.Context,
	keypair *bn254.KeyPair,
	socket string,
	quorumIds []chainio.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) error {
	args := t.Called(ctx, keypair, socket, quorumIds, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	return args.Error(0)
}

func (t *MockWriter) RegisterOperatorWithChurn(
	ctx context.Context,
	keypair *bn254.KeyPair,
	socket string,
	quorumIds []chainio.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
	churnReply *churner.ChurnReply) error {
	args := t.Called(ctx, keypair, socket, quorumIds, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry, churnReply)
	return args.Error(0)
}

func (t *MockWriter) DeregisterOperator(ctx context.Context, pubkeyG1 *bn254.G1Point, blockNumber uint32, quorumIds []chainio.QuorumID) error {
	args := t.Called()
	return args.Error(0)
}

func (t *MockWriter) UpdateOperatorSocket(ctx context.Context, socket string) error {
	args := t.Called()
	return args.Error(0)
}

func (t *MockWriter) BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][]chainio.OperatorID) (*types.Transaction, error) {
	args := t.Called(ctx, operatorsByQuorum)
	result := args.Get(0)
	return result.(*types.Transaction), args.Error(1)
}
