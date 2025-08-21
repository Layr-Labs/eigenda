package ephemeral

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	store := NewEphemeralCumulativePaymentStore()
	require.NotNil(t, store)
	require.NotNil(t, store.cumulativePayment)
	require.Equal(t, big.NewInt(0), store.cumulativePayment)
}

func TestAddCumulativePayment(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()
		maxCumulativePayment := big.NewInt(1000)

		newValue, err := store.AddCumulativePayment(ctx, big.NewInt(500), maxCumulativePayment)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), newValue)
		require.Equal(t, big.NewInt(500), store.cumulativePayment)

		newValue, err = store.AddCumulativePayment(ctx, big.NewInt(300), maxCumulativePayment)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(800), newValue)
		require.Equal(t, big.NewInt(800), store.cumulativePayment)
	})

	t.Run("subtraction", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()
		maxPayment := big.NewInt(1000)

		newValue, err := store.AddCumulativePayment(ctx, big.NewInt(500), maxPayment)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), newValue)

		newValue, err = store.AddCumulativePayment(ctx, big.NewInt(-200), maxPayment)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(300), newValue)
		require.Equal(t, big.NewInt(300), store.cumulativePayment)
	})

	t.Run("values are copies", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()
		maxCumulativePayment := big.NewInt(1000)

		amount := big.NewInt(300)
		newValue, err := store.AddCumulativePayment(ctx, amount, maxCumulativePayment)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(300), newValue)

		amount.Add(amount, big.NewInt(100))
		require.Equal(t, big.NewInt(300), store.cumulativePayment)

		newValue.Add(newValue, big.NewInt(50))
		require.Equal(t, big.NewInt(300), store.cumulativePayment)
	})
}

func TestAddCumulativePaymentErrorCases(t *testing.T) {
	t.Run("input validation errors", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()

		newValue, err := store.AddCumulativePayment(ctx, nil, big.NewInt(1000))
		require.Error(t, err)
		require.Nil(t, newValue)

		newValue, err = store.AddCumulativePayment(ctx, big.NewInt(100), nil)
		require.Error(t, err)
		require.Nil(t, newValue)

		newValue, err = store.AddCumulativePayment(ctx, big.NewInt(100), big.NewInt(-1000))
		require.Error(t, err)
		require.Nil(t, newValue)
	})

	t.Run("decrement below zero returns error", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()
		maxPayment := big.NewInt(1000)

		_, err := store.AddCumulativePayment(ctx, big.NewInt(500), maxPayment)
		require.NoError(t, err)

		newValue, err := store.AddCumulativePayment(ctx, big.NewInt(-600), maxPayment)
		require.Error(t, err)
		require.Nil(t, newValue)
		require.Equal(t, big.NewInt(500), store.cumulativePayment, "cumulative payment should not change on error")
	})

	t.Run("exceeds max returns InsufficientFundsError", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()
		maxPayment := big.NewInt(400)

		newValue, err := store.AddCumulativePayment(ctx, big.NewInt(500), maxPayment)
		require.Error(t, err)
		require.Nil(t, newValue)

		var insufficientFundsErr *ondemand.InsufficientFundsError
		require.ErrorAs(t, err, &insufficientFundsErr)
		require.Equal(t, big.NewInt(0), insufficientFundsErr.CurrentCumulativePayment)
		require.Equal(t, big.NewInt(400), insufficientFundsErr.MaxCumulativePayment)
		require.Equal(t, big.NewInt(500), insufficientFundsErr.BlobCost)

		require.Equal(t, big.NewInt(0), store.cumulativePayment, "cumulative payment should not change on error")
	})
}
