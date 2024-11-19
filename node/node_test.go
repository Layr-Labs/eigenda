package node_test

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	privateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	opID       = [32]byte{0}
)

type components struct {
	node        *node.Node
	tx          *coremock.MockWriter
	relayClient *clientsmock.MockRelayClient
}

func newComponents(t *testing.T) *components {
	dbPath := t.TempDir()
	keyPair, err := core.GenRandomBlsKeys()
	if err != nil {
		panic("failed to create a BLS Key")
	}
	config := &node.Config{
		Timeout:                   10 * time.Second,
		ExpirationPollIntervalSec: 1,
		QuorumIDList:              []core.QuorumID{0},
		DbPath:                    dbPath,
		ID:                        opID,
		NumBatchValidators:        runtime.GOMAXPROCS(0),
		EnableNodeApi:             false,
		EnableMetrics:             false,
		RegisterNodeAtStart:       false,
	}
	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic("failed to create a logger")
	}

	err = os.MkdirAll(config.DbPath, os.ModePerm)
	if err != nil {
		panic("failed to create a directory for db")
	}
	tx := &coremock.MockWriter{}

	mockVal := coremock.NewMockShardValidator()
	mockVal.On("ValidateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	chainState, _ := coremock.MakeChainDataMock(map[uint8]int{
		0: 4,
		1: 4,
		2: 4,
	})

	store, err := node.NewLevelDBStore(dbPath, logger, nil, 1e9, 1e9)
	if err != nil {
		panic("failed to create a new levelDB store")
	}
	defer os.Remove(dbPath)
	relayClient := clientsmock.NewRelayClient()
	return &components{
		node: &node.Node{
			Config:      config,
			Logger:      logger,
			KeyPair:     keyPair,
			Metrics:     nil,
			Store:       store,
			ChainState:  chainState,
			Validator:   mockVal,
			Transactor:  tx,
			RelayClient: relayClient,
		},
		tx:          tx,
		relayClient: relayClient,
	}
}

func TestNodeStartNoAddress(t *testing.T) {
	c := newComponents(t)
	c.node.Config.RegisterNodeAtStart = false

	err := c.node.Start(context.Background())
	assert.NoError(t, err)
}

func TestNodeStartOperatorIDMatch(t *testing.T) {
	c := newComponents(t)
	c.node.Config.RegisterNodeAtStart = true
	c.node.Config.EthClientConfig = geth.EthClientConfig{
		RPCURLs:          []string{"http://localhost:8545"},
		PrivateKeyString: privateKey,
		NumConfirmations: 1,
	}
	c.tx.On("GetRegisteredQuorumIdsForOperator", mock.Anything).Return([]core.QuorumID{}, nil)
	c.tx.On("GetOperatorSetParams", mock.Anything, mock.Anything).Return(&core.OperatorSetParam{
		MaxOperatorCount:         uint32(4),
		ChurnBIPsOfOperatorStake: uint16(1000),
		ChurnBIPsOfTotalStake:    uint16(10),
	}, nil)
	c.tx.On("GetNumberOfRegisteredOperatorForQuorum", mock.Anything, mock.Anything).Return(uint32(0), nil)
	c.tx.On("RegisterOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	c.tx.On("OperatorAddressToID", mock.Anything).Return(core.OperatorID(opID), nil)

	err := c.node.Start(context.Background())
	assert.NoError(t, err)
}

func TestNodeStartOperatorIDDoesNotMatch(t *testing.T) {
	c := newComponents(t)
	c.node.Config.RegisterNodeAtStart = true
	c.node.Config.EthClientConfig = geth.EthClientConfig{
		RPCURLs:          []string{"http://localhost:8545"},
		PrivateKeyString: privateKey,
		NumConfirmations: 1,
	}
	c.tx.On("GetRegisteredQuorumIdsForOperator", mock.Anything).Return([]core.QuorumID{}, nil)
	c.tx.On("GetOperatorSetParams", mock.Anything, mock.Anything).Return(&core.OperatorSetParam{
		MaxOperatorCount:         uint32(4),
		ChurnBIPsOfOperatorStake: uint16(1000),
		ChurnBIPsOfTotalStake:    uint16(10),
	}, nil)
	c.tx.On("GetNumberOfRegisteredOperatorForQuorum", mock.Anything, mock.Anything).Return(uint32(0), nil)
	c.tx.On("RegisterOperator", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	c.tx.On("OperatorAddressToID", mock.Anything).Return(core.OperatorID{1}, nil)

	err := c.node.Start(context.Background())
	assert.ErrorContains(t, err, "operator ID mismatch")
}

func TestGetReachabilityURL(t *testing.T) {
	url, err := node.GetReachabilityURL("https://dataapi.eigenda.xyz/", "123123123")
	assert.NoError(t, err)
	assert.Equal(t, "https://dataapi.eigenda.xyz/api/v1/operators-info/port-check?operator_id=123123123", url)
	url, err = node.GetReachabilityURL("https://dataapi.eigenda.xyz", "123123123")
	assert.NoError(t, err)
	assert.Equal(t, "https://dataapi.eigenda.xyz/api/v1/operators-info/port-check?operator_id=123123123", url)
}
