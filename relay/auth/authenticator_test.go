package auth

import (
	"context"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestMockSigning is a meta-test to verify that
// the test framework's BLS keys are functioning correctly.
func TestMockSigning(t *testing.T) {
	tu.InitializeRandom()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)

	operators, err := ics.GetIndexedOperators(context.Background(), 0)
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

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)

	timeout := 10 * time.Second

	authenticator := NewRequestAuthenticator(ics, timeout)

	request := randomGetChunksRequest()
	request.RequesterId = operatorID[:]
	hash := HashGetChunksRequest(request)
	signature := ics.KeyPairs[operatorID].SignMessage([32]byte(hash))
	request.RequesterSignature = signature.G1Point.Serialize()

	now := time.Now()

	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)
	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		request,
		now)
	require.NoError(t, err)

	// Making additional requests before timeout elapses should not trigger authentication for the address "foobar".
	// To probe at this, intentionally make a request that would be considered invalid if it were authenticated.
	invalidRequest := randomGetChunksRequest()
	invalidRequest.RequesterId = operatorID[:]
	invalidRequest.RequesterSignature = signature.G1Point.Serialize() // the previous signature is invalid here

	start := now
	for now.Before(start.Add(timeout)) {
		err = authenticator.AuthenticateGetChunksRequest(
			"foobar",
			invalidRequest,
			now)
		require.NoError(t, err)

		err = authenticator.AuthenticateGetChunksRequest(
			"baz",
			invalidRequest,
			now)
		require.Error(t, err)

		now = now.Add(time.Second)
	}

	// After the timeout elapses, new requests should trigger authentication.
	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		invalidRequest,
		now)
	require.Error(t, err)
}

func TestNonExistingClient(t *testing.T) {
	tu.InitializeRandom()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)

	timeout := 10 * time.Second

	authenticator := NewRequestAuthenticator(ics, timeout)

	invalidOperatorID := tu.RandomBytes(32)

	request := randomGetChunksRequest()
	request.RequesterId = invalidOperatorID

	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)
	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		request,
		time.Now())
	require.Error(t, err)
}

func TestBadSignature(t *testing.T) {
	tu.InitializeRandom()

	operatorID := mock.MakeOperatorId(0)
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		core.QuorumID(0): {
			operatorID: 1,
		},
	}
	ics, err := mock.NewChainDataMock(stakes)
	require.NoError(t, err)

	timeout := 10 * time.Second

	authenticator := NewRequestAuthenticator(ics, timeout)

	request := randomGetChunksRequest()
	request.RequesterId = operatorID[:]
	hash := HashGetChunksRequest(request)
	signature := ics.KeyPairs[operatorID].SignMessage([32]byte(hash))
	request.RequesterSignature = signature.G1Point.Serialize()

	now := time.Now()

	ics.Mock.On("GetCurrentBlockNumber").Return(uint(0), nil)
	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		request,
		now)
	require.NoError(t, err)

	// move time forward to wipe out previous authentication
	now = now.Add(timeout)

	// Change a byte in the signature to make it invalid
	request.RequesterSignature[0] = request.RequesterSignature[0] ^ 1

	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		request,
		now)
	require.Error(t, err)

	// Sign different data with the same key.
	signature = ics.KeyPairs[operatorID].SignMessage([32]byte(tu.RandomBytes(32)))
	request.RequesterSignature = signature.G1Point.Serialize()
	err = authenticator.AuthenticateGetChunksRequest(
		"foobar",
		request,
		now)
	require.Error(t, err)
}
