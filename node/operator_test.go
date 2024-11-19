package node_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/node"
	nodemock "github.com/Layr-Labs/eigenda/node/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterOperator(t *testing.T) {
	logger := logging.NewNoopLogger()
	operatorID := [32]byte(hexutil.MustDecode("0x3fbfefcdc76462d2cdb7d0cea75f27223829481b8b4aa6881c94cb2126a316ad"))
	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	// Create a new operator
	operator := &node.Operator{
		Address:             "0xB7Ad27737D88B07De48CDc2f379917109E993Be4",
		Socket:              "localhost:50051",
		Timeout:             10 * time.Second,
		PrivKey:             nil,
		KeyPair:             keyPair,
		OperatorId:          operatorID,
		QuorumIDs:           []core.QuorumID{0, 1},
		RegisterNodeAtStart: false,
	}
	createMockTx := func(quorumIDs []uint8) *coremock.MockWriter {
		tx := &coremock.MockWriter{}
		tx.On("GetRegisteredQuorumIdsForOperator").Return(quorumIDs, nil)
		tx.On("GetOperatorSetParams", mock.Anything, mock.Anything).Return(&core.OperatorSetParam{
			MaxOperatorCount:         1,
			ChurnBIPsOfOperatorStake: 20,
			ChurnBIPsOfTotalStake:    20000,
		}, nil)
		tx.On("GetNumberOfRegisteredOperatorForQuorum").Return(uint32(0), nil)
		tx.On("RegisterOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		return tx

	}
	tx1 := createMockTx([]uint8{2})
	churnerClient := &nodemock.ChurnerClient{}
	churnerClient.On("Churn").Return(nil, nil)
	err = node.RegisterOperator(context.Background(), operator, tx1, churnerClient, logger)
	assert.NoError(t, err)
	// Try to register with a quorum that's already registered
	tx2 := createMockTx([]uint8{0})
	err = node.RegisterOperator(context.Background(), operator, tx2, churnerClient, logger)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "quorums to register must be not registered yet"))
}

func TestRegisterOperatorWithChurn(t *testing.T) {
	logger := logging.NewNoopLogger()
	operatorID := [32]byte(hexutil.MustDecode("0x3fbfefcdc76462d2cdb7d0cea75f27223829481b8b4aa6881c94cb2126a316ad"))
	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)
	// Create a new operator
	operator := &node.Operator{
		Address:    "0xB7Ad27737D88B07De48CDc2f379917109E993Be4",
		Socket:     "localhost:50051",
		Timeout:    10 * time.Second,
		PrivKey:    nil,
		KeyPair:    keyPair,
		OperatorId: operatorID,
		QuorumIDs:  []core.QuorumID{1},
	}
	tx := &coremock.MockWriter{}
	tx.On("GetRegisteredQuorumIdsForOperator").Return([]uint8{2}, nil)
	tx.On("GetOperatorSetParams", mock.Anything, mock.Anything).Return(&core.OperatorSetParam{
		MaxOperatorCount:         1,
		ChurnBIPsOfOperatorStake: 20,
		ChurnBIPsOfTotalStake:    20000,
	}, nil)
	tx.On("GetNumberOfRegisteredOperatorForQuorum").Return(uint32(1), nil)
	tx.On("RegisterOperatorWithChurn", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	churnerClient := &nodemock.ChurnerClient{}
	churnerClient.On("Churn").Return(nil, nil)
	err = node.RegisterOperator(context.Background(), operator, tx, churnerClient, logger)
	assert.NoError(t, err)
	tx.AssertCalled(t, "RegisterOperatorWithChurn", mock.Anything, mock.Anything, mock.Anything, []core.QuorumID{1}, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
