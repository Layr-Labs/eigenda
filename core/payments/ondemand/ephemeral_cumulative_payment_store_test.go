package ondemand

import (
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

		setValue := big.NewInt(500)
		err := store.SetCumulativePayment(setValue)
		require.NoError(t, err)
		value, err := store.GetCumulativePayment()
		require.NoError(t, err)
		assert.Equal(t, setValue, value)

		setValue = big.NewInt(2000)
		err = store.SetCumulativePayment(setValue)
		require.NoError(t, err)
		value, err = store.GetCumulativePayment()
		require.NoError(t, err)
		assert.Equal(t, setValue, value)
	})

	t.Run("make sure values are copies", func(t *testing.T) {
		store := NewEphemeralCumulativePaymentStore()

		original := big.NewInt(300)
		err := store.SetCumulativePayment(original)
		require.NoError(t, err)

		original.Add(original, big.NewInt(100))

		value, err := store.GetCumulativePayment()
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(300), value)

		value.Add(value, big.NewInt(50))
		value2, err := store.GetCumulativePayment()
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(300), value2)
	})
}
