package churner_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	dacore "github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	indexermock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
)

var (
	keyPair                        *dacore.KeyPair
	quorumIds                      = []uint32{0, 1}
	logger                         = logging.NewNoopLogger()
	transactorMock                 = &coremock.MockWriter{}
	mockIndexer                    = &indexermock.MockIndexedChainState{}
	operatorAddr                   = gethcommon.HexToAddress("0x0000000000000000000000000000000000000001")
	operatorToChurnInPrivateKeyHex = "0000000000000000000000000000000000000000000000000000000000000020"
	churnerPrivateKeyHex           = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expectedReplySignature         = []byte{0x4, 0xc, 0x2b, 0xd1, 0xce, 0xde, 0xb8, 0xbf, 0xb6, 0xba, 0x99, 0x3, 0x96, 0x57, 0x86, 0xcc, 0x4c, 0xf4, 0xed, 0xcf, 0x2f, 0xdb, 0x64, 0xf1, 0xca, 0x6, 0x80, 0x37, 0xd6, 0x6a, 0xf5, 0x92, 0x64, 0x49, 0x1c, 0xcb, 0x7d, 0xa5, 0x11, 0x9a, 0xb2, 0xab, 0x3, 0x11, 0x87, 0x31, 0x84, 0xd8, 0xff, 0xd, 0xd5, 0xd, 0x75, 0x93, 0xbd, 0x7, 0xf4, 0x2b, 0x2, 0x32, 0xa6, 0xf2, 0xb, 0xf1, 0x1c}
	numRetries                     = 0
)

func TestChurn(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()

	salt := crypto.Keccak256([]byte(operatorToChurnInPrivateKeyHex), []byte("ChurnRequest"))
	request := &pb.ChurnRequest{
		OperatorAddress:            operatorAddr.Hex(),
		OperatorToRegisterPubkeyG1: keyPair.PubKey.Serialize(),
		OperatorToRegisterPubkeyG2: keyPair.GetPubKeyG2().Serialize(),
		Salt:                       salt,
		QuorumIds:                  quorumIds,
	}

	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		[]byte(request.OperatorAddress),
		request.OperatorToRegisterPubkeyG1,
		request.OperatorToRegisterPubkeyG2,
		request.Salt,
	)
	copy(requestHash[:], requestHashBytes)

	signature := keyPair.SignMessage(requestHash)
	request.OperatorRequestSignature = signature.Serialize()

	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: keyPair.PubKey,
	}, nil)

	reply, err := s.Churn(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotNil(t, reply.SignatureWithSaltAndExpiry.GetSalt())
	assert.NotNil(t, reply.SignatureWithSaltAndExpiry.GetExpiry())
	assert.Equal(t, expectedReplySignature, reply.SignatureWithSaltAndExpiry.GetSignature())
	assert.Equal(t, 2, len(reply.OperatorsToChurn))
	actualQuorums := make([]uint32, 0)
	for _, param := range reply.OperatorsToChurn {
		actualQuorums = append(actualQuorums, param.GetQuorumId())
		if param.GetQuorumId() == 0 {
			// no churning for quorum 0
			assert.Equal(t, gethcommon.HexToAddress("0x").Bytes(), param.GetOperator())
			assert.Nil(t, param.GetPubkey())
		}
		if param.GetQuorumId() == 1 {
			// churn the operator with quorum 1
			assert.Equal(t, operatorAddr.Bytes(), param.GetOperator())
			assert.Equal(t, keyPair.PubKey.Serialize(), param.GetPubkey())
		}
	}
	assert.ElementsMatch(t, quorumIds, actualQuorums)

	// retry prior to expiry should fail
	_, err = s.Churn(ctx, request)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = ResourceExhausted desc = previous approval not expired, retry in 900 seconds")
}

func TestChurnWithInvalidQuorum(t *testing.T) {
	s := newTestServer(t)
	ctx := context.Background()

	salt := crypto.Keccak256([]byte(operatorToChurnInPrivateKeyHex), []byte("ChurnRequest"))
	request := &pb.ChurnRequest{
		OperatorToRegisterPubkeyG1: keyPair.PubKey.Serialize(),
		OperatorToRegisterPubkeyG2: keyPair.GetPubKeyG2().Serialize(),
		Salt:                       salt,
		QuorumIds:                  []uint32{0, 1, 2},
	}

	var requestHash [32]byte
	requestHashBytes := crypto.Keccak256(
		[]byte("ChurnRequest"),
		request.OperatorToRegisterPubkeyG1,
		request.OperatorToRegisterPubkeyG2,
		request.Salt,
	)
	copy(requestHash[:], requestHashBytes)

	signature := keyPair.SignMessage(requestHash)
	request.OperatorRequestSignature = signature.Serialize()

	mockIndexer.On("GetIndexedOperatorInfoByOperatorId").Return(&core.IndexedOperatorInfo{
		PubkeyG1: keyPair.PubKey,
	}, nil)

	_, err := s.Churn(ctx, request)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = InvalidArgument desc = invalid request: invalid request: the quorum_id must be in range [0, 1], but found 2")
}

func setupMockWriter() {
	transactorMock.On("StakeRegistry").Return(gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"), nil).Once()
	transactorMock.On("OperatorIDToAddress").Return(operatorAddr, nil)
	transactorMock.On("GetCurrentQuorumBitmapByOperatorId").Return(big.NewInt(0), nil)
	transactorMock.On("GetCurrentBlockNumber").Return(uint32(2), nil)
	transactorMock.On("GetQuorumCount").Return(uint8(2), nil)
	transactorMock.On("GetOperatorStakesForQuorums").Return(dacore.OperatorStakes{
		0: {
			0: {
				OperatorID: makeOperatorId(1),
				Stake:      big.NewInt(2),
			},
		},
		1: {
			0: {
				OperatorID: makeOperatorId(1),
				Stake:      big.NewInt(2),
			},
		},
	}, nil)
	transactorMock.On("GetOperatorSetParams", mock.Anything, uint8(0)).Return(&dacore.OperatorSetParam{
		MaxOperatorCount:         2,
		ChurnBIPsOfOperatorStake: 20,
		ChurnBIPsOfTotalStake:    20000,
	}, nil)
	transactorMock.On("GetOperatorSetParams", mock.Anything, uint8(1)).Return(&dacore.OperatorSetParam{
		MaxOperatorCount:         1,
		ChurnBIPsOfOperatorStake: 20,
		ChurnBIPsOfTotalStake:    20000,
	}, nil)
	transactorMock.On("WeightOfOperatorForQuorum").Return(big.NewInt(1), nil)
	transactorMock.On("CalculateOperatorChurnApprovalDigestHash").Return([32]byte{1, 2, 3}, nil)
}

func newTestServer(t *testing.T) *churner.Server {
	config := &churner.Config{
		LoggerConfig: common.DefaultLoggerConfig(),
		EthClientConfig: geth.EthClientConfig{
			PrivateKeyString: churnerPrivateKeyHex,
			NumRetries:       numRetries,
		},
		ChurnApprovalInterval: 15 * time.Minute,
	}

	var err error
	keyPair, err = dacore.GenRandomBlsKeys()
	if err != nil {
		t.Fatalf("Generating random BLS keys Error: %s", err.Error())
	}

	setupMockWriter()

	metrics := churner.NewMetrics("9001", logger)
	cn, err := churner.NewChurner(config, mockIndexer, transactorMock, logger, metrics)
	if err != nil {
		log.Fatalln("cannot create churner", err)
	}

	return churner.NewServer(config, cn, logger, metrics)
}

func makeOperatorId(id int) dacore.OperatorID {
	data := [32]byte{}
	copy(data[:], []byte(fmt.Sprintf("%d", id)))
	return data
}
