package mock

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type MockWriter struct {
	mock.Mock
}

var _ core.Writer = (*MockWriter)(nil)

func (t *MockWriter) GetBlockStaleMeasure(ctx context.Context) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockWriter) GetStoreDurationBlocks(ctx context.Context) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockWriter) GetRegisteredQuorumIdsForOperator(ctx context.Context, operator core.OperatorID) ([]core.QuorumID, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]core.QuorumID), args.Error(1)
}

func (t *MockWriter) RegisterOperator(
	ctx context.Context,
	keypair *core.KeyPair,
	socket string,
	quorumIds []core.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
) error {
	args := t.Called(ctx, keypair, socket, quorumIds, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)
	return args.Error(0)
}

func (t *MockWriter) RegisterOperatorWithChurn(
	ctx context.Context,
	keypair *core.KeyPair,
	socket string,
	quorumIds []core.QuorumID,
	operatorEcdsaPrivateKey *ecdsa.PrivateKey,
	operatorToAvsRegistrationSigSalt [32]byte,
	operatorToAvsRegistrationSigExpiry *big.Int,
	churnReply *churner.ChurnReply) error {
	args := t.Called(ctx, keypair, socket, quorumIds, operatorEcdsaPrivateKey, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry, churnReply)
	return args.Error(0)
}

func (t *MockWriter) DeregisterOperator(ctx context.Context, pubkeyG1 *core.G1Point, blockNumber uint32, quorumIds []core.QuorumID) error {
	args := t.Called()
	return args.Error(0)
}

func (t *MockWriter) UpdateOperatorSocket(ctx context.Context, socket string) error {
	args := t.Called()
	return args.Error(0)
}

func (t *MockWriter) BuildEjectOperatorsTxn(ctx context.Context, operatorsByQuorum [][]core.OperatorID) (*types.Transaction, error) {
	args := t.Called(ctx, operatorsByQuorum)
	result := args.Get(0)
	return result.(*types.Transaction), args.Error(1)
}

func (t *MockWriter) GetOperatorStakes(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (core.OperatorStakes, []core.QuorumID, error) {
	args := t.Called()
	result0 := args.Get(0)
	result1 := args.Get(1)
	return result0.(core.OperatorStakes), result1.([]core.QuorumID), args.Error(1)
}

func (t *MockWriter) GetOperatorStakesForQuorums(ctx context.Context, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStakes, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(core.OperatorStakes), args.Error(1)
}

func (t *MockWriter) BuildConfirmBatchTxn(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Transaction, error) {
	args := t.Called(ctx, batchHeader, quorums, signatureAggregation)
	result := args.Get(0)
	return result.(*types.Transaction), args.Error(1)
}

func (t *MockWriter) ConfirmBatch(ctx context.Context, batchHeader *core.BatchHeader, quorums map[core.QuorumID]*core.QuorumResult, signatureAggregation *core.SignatureAggregation) (*types.Receipt, error) {
	args := t.Called()
	var receipt *types.Receipt
	if args.Get(0) != nil {
		receipt = args.Get(0).(*types.Receipt)
	}
	return receipt, args.Error(1)
}

func (t *MockWriter) StakeRegistry(ctx context.Context) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockWriter) OperatorIDToAddress(ctx context.Context, operatorId core.OperatorID) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockWriter) OperatorAddressToID(ctx context.Context, address gethcommon.Address) (core.OperatorID, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(core.OperatorID), args.Error(1)
}

func (t *MockWriter) BatchOperatorIDToAddress(ctx context.Context, operatorIds []core.OperatorID) ([]gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]gethcommon.Address), args.Error(1)
}

func (t *MockWriter) GetQuorumBitmapForOperatorsAtBlockNumber(ctx context.Context, operatorIds []core.OperatorID, blockNumber uint32) ([]*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]*big.Int), args.Error(1)
}

func (t *MockWriter) GetCurrentQuorumBitmapByOperatorId(ctx context.Context, operatorId core.OperatorID) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockWriter) GetOperatorSetParams(ctx context.Context, quorumID core.QuorumID) (*core.OperatorSetParam, error) {
	args := t.Called(ctx, quorumID)
	result := args.Get(0)
	return result.(*core.OperatorSetParam), args.Error(1)
}

func (t *MockWriter) GetNumberOfRegisteredOperatorForQuorum(ctx context.Context, quorumID core.QuorumID) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockWriter) WeightOfOperatorForQuorum(ctx context.Context, quorumID core.QuorumID, operator gethcommon.Address) (*big.Int, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*big.Int), args.Error(1)
}

func (t *MockWriter) CalculateOperatorChurnApprovalDigestHash(
	ctx context.Context,
	operatorAddress gethcommon.Address,
	operatorId core.OperatorID,
	operatorsToChurn []core.OperatorToChurn,
	salt [32]byte,
	expiry *big.Int,
) ([32]byte, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([32]byte), args.Error(1)
}

func (t *MockWriter) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockWriter) GetQuorumCount(ctx context.Context, blockNumber uint32) (uint8, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint8), args.Error(1)
}

func (t *MockWriter) GetQuorumSecurityParams(ctx context.Context, blockNumber uint32) ([]core.SecurityParam, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]core.SecurityParam), args.Error(1)
}

func (t *MockWriter) GetRequiredQuorumNumbers(ctx context.Context, blockNumber uint32) ([]uint8, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]uint8), args.Error(1)
}

func (t *MockWriter) GetNumBlobVersions(ctx context.Context) (uint16, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint16), args.Error(1)
}

func (t *MockWriter) GetVersionedBlobParams(ctx context.Context, blobVersion uint16) (*core.BlobVersionParameters, error) {
	args := t.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result := args.Get(0)
	return result.(*core.BlobVersionParameters), args.Error(1)
}

func (t *MockWriter) GetAllVersionedBlobParams(ctx context.Context) (map[uint16]*core.BlobVersionParameters, error) {
	args := t.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(map[uint16]*core.BlobVersionParameters), args.Error(1)
}

func (t *MockWriter) PubkeyHashToOperator(ctx context.Context, operatorId core.OperatorID) (gethcommon.Address, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(gethcommon.Address), args.Error(1)
}

func (t *MockWriter) GetReservedPayments(ctx context.Context, accountIDs []gethcommon.Address) (map[gethcommon.Address]*core.ReservedPayment, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(map[gethcommon.Address]*core.ReservedPayment), args.Error(1)
}

func (t *MockWriter) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*core.ReservedPayment), args.Error(1)
}

func (t *MockWriter) GetOnDemandPayments(ctx context.Context, accountIDs []gethcommon.Address) (map[gethcommon.Address]*core.OnDemandPayment, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(map[gethcommon.Address]*core.OnDemandPayment), args.Error(1)
}

func (t *MockWriter) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(*core.OnDemandPayment), args.Error(1)
}

func (t *MockWriter) GetOperatorSocket(ctx context.Context, operatorID core.OperatorID) (string, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(string), args.Error(1)
}

func (t *MockWriter) GetNumRelays(ctx context.Context) (uint32, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint32), args.Error(1)
}

func (t *MockWriter) GetRelayURL(ctx context.Context, key uint32) (string, error) {
	args := t.Called()
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	result := args.Get(0)
	return result.(string), args.Error(1)
}

func (t *MockWriter) GetRelayURLs(ctx context.Context) (map[uint32]string, error) {
	args := t.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}

	return result.(map[uint32]string), args.Error(1)
}

func (t *MockWriter) GetDisperserAddress(ctx context.Context, disperserID uint32) (gethcommon.Address, error) {
	args := t.Called(disperserID)
	result := args.Get(0)
	if result == nil {
		var zeroValue gethcommon.Address
		return zeroValue, args.Error(1)
	}

	return result.(gethcommon.Address), args.Error(1)
}
