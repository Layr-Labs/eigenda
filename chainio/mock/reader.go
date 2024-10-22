package mock

import (
	"context"
	"math/big"

	"github.com/Layr-Labs/eigenda/chainio"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

var _ chainio.Reader = (*MockReader)(nil)

func (t *MockReader) GetBlockStaleMeasure(ctx context.Context) (uint32, error) {
	args := t.Called()
	return *new(uint32), args.Error(0)
}

func (t *MockReader) GetStoreDurationBlocks(ctx context.Context) (uint32, error) {
	args := t.Called()
	return *new(uint32), args.Error(0)
}

func (t *MockReader) GetRegisteredQuorumIdsForOperator(ctx context.Context, operator [32]byte) ([]chainio.QuorumID, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]chainio.QuorumID), args.Error(1)
}

func (t *MockReader) GetOperatorStakes(ctx context.Context, operatorId [32]byte, blockNumber uint32) (chainio.OperatorStakes, []chainio.QuorumID, error) {
	args := t.Called()
	result0 := args.Get(0)
	result1 := args.Get(1)
	return result0.(chainio.OperatorStakes), result1.([]chainio.QuorumID), args.Error(1)
}

func (t *MockReader) GetOperatorStakesForQuorums(ctx context.Context, quorums []chainio.QuorumID, blockNumber uint32) (chainio.OperatorStakes, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(chainio.OperatorStakes), args.Error(1)
}

func (t *MockReader) StakeRegistry(ctx context.Context) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockReader) OperatorIDToAddress(ctx context.Context, operatorId [32]byte) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockReader) OperatorAddressToID(ctx context.Context, address gethcommon.Address) ([32]byte, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([32]byte), args.Error(1)
}

func (t *MockReader) BatchOperatorIDToAddress(ctx context.Context, operatorIds [][32]byte) ([]gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]gethcommon.Address), args.Error(1)
}

func (t *MockReader) GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds [][32]byte, blockNumber uint32) ([]*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]*big.Int), args.Error(1)
}

func (t *MockReader) GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId [32]byte) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockReader) GetOperatorSetParams(ctx context.Context, quorumID chainio.QuorumID) (*chainio.OperatorSetParam, error) {
	args := t.Called(ctx, quorumID)
	result := args.Get(0)
	return result.(*chainio.OperatorSetParam), args.Error(1)
}

func (t *MockReader) GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID chainio.QuorumID) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockReader) WeightOfOperatorForQuorum(ctx context.Context, quorumID chainio.QuorumID, operator gethcommon.Address) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockReader) CalculateOperatorChurnApprovalDigestHash(
	ctx context.Context,
	operatorAddress gethcommon.Address,
	operatorId [32]byte,
	operatorsToChurn []chainio.OperatorToChurn,
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

func (t *MockReader) PubkeyHashToOperator(ctx context.Context, operatorId [32]byte) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}
