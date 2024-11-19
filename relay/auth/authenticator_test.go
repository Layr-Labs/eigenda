package auth

import (
	"context"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
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

	err = authenticator.AuthenticateGetChunksRequest(
		string(invalidOperatorID),
		&pb.GetChunksRequest{
			RequesterId: invalidOperatorID,
		},
		time.Now())
}
