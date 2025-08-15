package ondemand

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	store := NewEphemeralCumulativePaymentStore()
	require.NotNil(t, store)
	require.NotNil(t, store.cumulativePayment)
	assert.Equal(t, big.NewInt(0), store.cumulativePayment)
}

func TestSetAndGet(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()

		ctx := context.Background()
		setValue := big.NewInt(500)
		err := store.SetCumulativePayment(ctx, setValue)
		require.NoError(t, err)
		value, err := store.GetCumulativePayment(ctx)
		require.NoError(t, err)
		assert.Equal(t, setValue, value)

		setValue = big.NewInt(2000)
		err = store.SetCumulativePayment(ctx, setValue)
		require.NoError(t, err)
		value, err = store.GetCumulativePayment(ctx)
		require.NoError(t, err)
		assert.Equal(t, setValue, value)
	})

	t.Run("make sure values are copies", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()
		ctx := context.Background()

		original := big.NewInt(300)
		err := store.SetCumulativePayment(ctx, original)
		require.NoError(t, err)

		original.Add(original, big.NewInt(100))

		value, err := store.GetCumulativePayment(ctx)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(300), value)

		value.Add(value, big.NewInt(50))
		value2, err := store.GetCumulativePayment(ctx)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(300), value2)
	})
}

func TestSetCumulativePaymentErrorCases(t *testing.T) {
	store := NewEphemeralCumulativePaymentStore()
	ctx := context.Background()

	err := store.SetCumulativePayment(ctx, nil)
	require.Error(t, err)

	err = store.SetCumulativePayment(ctx, big.NewInt(-100))
	require.Error(t, err)

	err = store.SetCumulativePayment(ctx, big.NewInt(0))
	require.NoError(t, err, "setting 0 should work")

	err = store.SetCumulativePayment(ctx, big.NewInt(1000))
	require.NoError(t, err, "invalid sets shouldn't break anything")
}
