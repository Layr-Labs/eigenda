package churner_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/stretchr/testify/assert"

	dacore "github.com/Layr-Labs/eigenda/core"
	indexermock "github.com/Layr-Labs/eigenda/core/mock"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestProcessChurnRequest(t *testing.T) {
	setupMockWriter()
	mockIndexer := &indexermock.MockIndexedChainState{}
	config := &churner.Config{
		LoggerConfig: common.DefaultLoggerConfig(),
		EthClientConfig: geth.EthClientConfig{
			PrivateKeyString: churnerPrivateKeyHex,
			NumConfirmations: 0,
			NumRetries:       numRetries,
		},
	}
	metrics := churner.NewMetrics("9001", logger)
	cn, err := churner.NewChurner(config, mockIndexer, transactorMock, logger, metrics)
	assert.NoError(t, err)
	assert.NotNil(t, cn)

	ctx := context.Background()

	keyPair, err := dacore.GenRandomBlsKeys()
	assert.NoError(t, err)

	salt := [32]byte{1, 2, 3}
	request := &churner.ChurnRequest{
		OperatorToRegisterPubkeyG1: keyPair.PubKey,
		OperatorToRegisterPubkeyG2: keyPair.GetPubKeyG2(),
		Salt:                       salt,
		QuorumIDs:                  []dacore.QuorumID{0, 1},
	}

	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		request.OperatorToRegisterPubkeyG1.Serialize(),
		request.OperatorToRegisterPubkeyG2.Serialize(),
		request.Salt[:],
	)
	copy(requestHash[:], requestHashBytes)

	request.OperatorRequestSignature = keyPair.SignMessage(requestHash)

	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: keyPair.PubKey,
	}, nil)

	response, err := cn.ProcessChurnRequest(ctx, gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"), request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.SignatureWithSaltAndExpiry.Salt)
	assert.NotNil(t, response.SignatureWithSaltAndExpiry.Expiry)
	assert.Equal(t, expectedReplySignature, response.SignatureWithSaltAndExpiry.Signature)
	assert.Equal(t, 2, len(response.OperatorsToChurn))
	actualQuorums := make([]dacore.QuorumID, 0)
	for _, o := range response.OperatorsToChurn {
		actualQuorums = append(actualQuorums, o.QuorumId)
		if o.QuorumId == 0 {
			// no churning for quorum 0
			assert.Equal(t, gethcommon.HexToAddress("0x"), o.Operator)
			assert.Nil(t, o.Pubkey)
		}
		if o.QuorumId == 1 {
			// churn the operator with quorum 1
			assert.Equal(t, gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"), o.Operator)
			assert.Equal(t, keyPair.PubKey, o.Pubkey)
		}
	}
	assert.ElementsMatch(t, []dacore.QuorumID{0, 1}, actualQuorums)
}
