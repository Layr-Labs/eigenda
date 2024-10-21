package mock

import (
	"context"
	"math/big"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

var _ corev2.Reader = (*MockReader)(nil)

func (t *MockReader) GetBlockStaleMeasure(ctx context.Context) (uint32, error) {
	args := t.Called()
	return *new(uint32), args.Error(0)
}

func (t *MockReader) GetStoreDurationBlocks(ctx context.Context) (uint32, error) {
	args := t.Called()
	return *new(uint32), args.Error(0)
}

func (t *MockReader) GetRegisteredQuorumIdsForOperator(ctx context.Context, operator corev2.OperatorID) ([]corev2.QuorumID, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]corev2.QuorumID), args.Error(1)
}

func (t *MockReader) GetOperatorStakes(ctx context.Context, operatorId corev2.OperatorID, blockNumber uint32) (corev2.OperatorStakes, []corev2.QuorumID, error) {
	args := t.Called()
	result0 := args.Get(0)
	result1 := args.Get(1)
	return result0.(corev2.OperatorStakes), result1.([]corev2.QuorumID), args.Error(1)
}

func (t *MockReader) GetOperatorStakesForQuorums(ctx context.Context, quorums []corev2.QuorumID, blockNumber uint32) (corev2.OperatorStakes, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(corev2.OperatorStakes), args.Error(1)
}

func (t *MockReader) StakeRegistry(ctx context.Context) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockReader) OperatorIDToAddress(ctx context.Context, operatorId corev2.OperatorID) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockReader) OperatorAddressToID(ctx context.Context, address gethcommon.Address) (corev2.OperatorID, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(corev2.OperatorID), args.Error(1)
}

func (t *MockReader) BatchOperatorIDToAddress(ctx context.Context, operatorIds []corev2.OperatorID) ([]gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]gethcommon.Address), args.Error(1)
}

func (t *MockReader) GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds []corev2.OperatorID, blockNumber uint32) ([]*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]*big.Int), args.Error(1)
}

func (t *MockReader) GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId corev2.OperatorID) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockReader) GetOperatorSetParams(ctx context.Context, quorumID corev2.QuorumID) (*corev2.OperatorSetParam, error) {
	args := t.Called(ctx, quorumID)
	result := args.Get(0)
	return result.(*corev2.OperatorSetParam), args.Error(1)
}

func (t *MockReader) GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID corev2.QuorumID) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockReader) WeightOfOperatorForQuorum(ctx context.Context, quorumID corev2.QuorumID, operator gethcommon.Address) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockReader) CalculateOperatorChurnApprovalDigestHash(
	ctx context.Context,
	operatorAddress gethcommon.Address,
	operatorId corev2.OperatorID,
	operatorsToChurn []corev2.OperatorToChurn,
	salt [32]byte,
	expiry *big.Int,
) ([32]byte, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([32]byte), args.Error(1)
}

func (t *MockReader) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockReader) GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint8), args.Error(1)
}

func (t *MockReader) GetRequiredQuorumNumbers(ctx context.Context, blockNumber uint32) ([]uint8, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]uint8), args.Error(1)
}

func (t *MockReader) PubkeyHashToOperator(ctx context.Context, operatorId corev2.OperatorID) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}
