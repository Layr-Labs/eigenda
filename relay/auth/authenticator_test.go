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

	timeout := 10 * time.Second

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024, timeout)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = operatorID[:]
	signature := SignGetChunksRequest(ics.KeyPairs[operatorID], request)
	request.OperatorSignature = signature

	now := time.Now()

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		now)
	require.NoError(t, err)

	// Making additional requests before timeout elapses should not trigger authentication for the address "foobar".
	// To probe at this, intentionally make a request that would be considered invalid if it were authenticated.
	invalidRequest := randomGetChunksRequest()
	invalidRequest.OperatorId = operatorID[:]
	invalidRequest.OperatorSignature = signature // the previous signature is invalid here

	start := now
	for now.Before(start.Add(timeout)) {
		err = authenticator.AuthenticateGetChunksRequest(
			ctx,
			"foobar",
			invalidRequest,
			now)
		require.NoError(t, err)

		err = authenticator.AuthenticateGetChunksRequest(
			ctx,
			"baz",
			invalidRequest,
			now)
		require.Error(t, err)

		now = now.Add(time.Second)
	}

	// After the timeout elapses, new requests should trigger authentication.
	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		invalidRequest,
		now)
	require.Error(t, err)
}

func TestAuthenticationSavingDisabled(t *testing.T) {
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

	// This disables saving of authentication results.
	timeout := time.Duration(0)

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024, timeout)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = operatorID[:]
	signature := SignGetChunksRequest(ics.KeyPairs[operatorID], request)
	request.OperatorSignature = signature

	now := time.Now()

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		now)
	require.NoError(t, err)

	// There is no authentication timeout, so a new request should trigger authentication.
	// To probe at this, intentionally make a request that would be considered invalid if it were authenticated.
	invalidRequest := randomGetChunksRequest()
	invalidRequest.OperatorId = operatorID[:]
	invalidRequest.OperatorSignature = signature // the previous signature is invalid here

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		invalidRequest,
		now)
	require.Error(t, err)
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

	timeout := 10 * time.Second

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024, timeout)
	require.NoError(t, err)

	invalidOperatorID := tu.RandomBytes(32)

	request := randomGetChunksRequest()
	request.OperatorId = invalidOperatorID

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		time.Now())
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

	timeout := 10 * time.Second

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024, timeout)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = operatorID[:]
	request.OperatorSignature = SignGetChunksRequest(ics.KeyPairs[operatorID], request)

	now := time.Now()

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		now)
	require.NoError(t, err)

	// move time forward to wipe out previous authentication
	now = now.Add(timeout)

	// Change a byte in the signature to make it invalid
	request.OperatorSignature[0] = request.OperatorSignature[0] ^ 1

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		now)
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

	timeout := 10 * time.Second

	authenticator, err := NewRequestAuthenticator(ctx, ics, 1024, timeout)
	require.NoError(t, err)

	request := randomGetChunksRequest()
	request.OperatorId = nil

	err = authenticator.AuthenticateGetChunksRequest(
		ctx,
		"foobar",
		request,
		time.Now())
	require.Error(t, err)
}
