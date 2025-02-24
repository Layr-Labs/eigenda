package test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestRequiredQuorumsCache verifies basic cache functionality, ensuring correct values returned and expected reuse
func TestRequiredQuorumsCache(t *testing.T) {
	testRandom := random.NewTestRandom()

	mockCertVerifier := clientsmock.MockCertVerifier{}
	store, err := clients.NewRequiredQuorumsStore(&mockCertVerifier)
	require.NoError(t, err)

	expectedRequiredQuorums1 := testRandom.Bytes(100)
	expectedRequiredQuorums2 := testRandom.Bytes(200)
	address1 := "asdf"
	address2 := "qwert"

	mockCertVerifier.On("GetQuorumNumbersRequired", mock.Anything, address1).Return(
		expectedRequiredQuorums1, nil).Once()
	mockCertVerifier.On("GetQuorumNumbersRequired", mock.Anything, address2).Return(
		expectedRequiredQuorums2, nil).Once()

	quorumNumbers, err := store.GetQuorumNumbersRequired(context.Background(), address1)
	require.NoError(t, err)
	require.Equal(t, expectedRequiredQuorums1, quorumNumbers)

	quorumNumbers, err = store.GetQuorumNumbersRequired(context.Background(), address1)
	require.NoError(t, err)
	require.Equal(t, expectedRequiredQuorums1, quorumNumbers)

	quorumNumbers, err = store.GetQuorumNumbersRequired(context.Background(), address2)
	require.NoError(t, err)
	require.Equal(t, expectedRequiredQuorums2, quorumNumbers)

	quorumNumbers, err = store.GetQuorumNumbersRequired(context.Background(), address1)
	require.NoError(t, err)
	require.Equal(t, expectedRequiredQuorums1, quorumNumbers)
}
