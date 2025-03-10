package auth

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/api/hashing"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/require"
)

// TestMockSigning is a meta-test to verify that
// the test framework's BLS keys are functioning correctly.
func TestMockSigning(t *testing.T) {
	tu.InitializeRandom()

	ctx := context.Background()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)

	operators, err := ics.GetIndexedOperators(ctx, 0)
	require.NoError(t, err)

	operator, ok := operators[operatorID]
	require.True(t, ok)

	bytesToSign := tu.RandomBytes(32)
	signature := ics.KeyPairs[operatorID].SignMessage([32]byte(bytesToSign))

	isValid := signature.Verify(operator.PubkeyG2, [32]byte(bytesToSign))
	require.True(t, isValid)

	// Changing a byte in the message should invalidate the signature
	bytesToSign[0] = bytesToSign[0] ^ 1

	isValid = signature.Verify(operator.PubkeyG2, [32]byte(bytesToSign))
	require.False(t, isValid)
}

func TestValidRequest(t *testing.T) {
	tu.InitializeRandom()

	ctx := context.Background()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)
	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = operatorID[:]
	signature, err := SignGetChunksRequest(ics.KeyPairs[operatorID], request)
	require.NoError(t, err)
	request.OperatorSignature = signature

	hash, err := authenticator.AuthenticateGetChunksRequest(ctx, request)
	require.NoError(t, err)
	expectedHash, err := hashing.HashGetChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

func TestNonExistingClient(t *testing.T) {
	tu.InitializeRandom()

	ctx := context.Background()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)
	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024)
	require.NoError(t, err)

	invalidOperatorID := tu.RandomBytes(32)

	request := randomGetChunksRequest()
	request.OperatorId = invalidOperatorID

	_, err = authenticator.AuthenticateGetChunksRequest(ctx, request)
	require.Error(t, err)
}

func TestBadSignature(t *testing.T) {
	tu.InitializeRandom()

	ctx := context.Background()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)
	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = operatorID[:]
	request.OperatorSignature, err = SignGetChunksRequest(ics.KeyPairs[operatorID], request)
	require.NoError(t, err)

	hash, err := authenticator.AuthenticateGetChunksRequest(ctx, request)
	require.NoError(t, err)
	expectedHash, err := hashing.HashGetChunksRequest(request)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)

	// Change a byte in the signature to make it invalid
	request.OperatorSignature[0] = request.OperatorSignature[0] ^ 1

	_, err = authenticator.AuthenticateGetChunksRequest(ctx, request)
	require.Error(t, err)
}

func TestMissingOperatorID(t *testing.T) {
	tu.InitializeRandom()

	ctx := context.Background()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)
	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = nil

	_, err = authenticator.AuthenticateGetChunksRequest(ctx, request)
	require.Error(t, err)
}
