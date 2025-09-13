package ondemand_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDebit(t *testing.T) {
	t.Run("successful debit", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestDebitSuccessful")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		cumulativePayment, err := ledger.Debit(t.Context(), 50, []core.QuorumID{0})
		require.NoError(t, err)
		require.NotNil(t, cumulativePayment)
		require.Equal(t, big.NewInt(50), cumulativePayment)

		// verify the store was updated
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(50), storedPayment)
	})

	t.Run("invalid quorum", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestDebitInvalidQuorum")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// quorum 5 not supported
		cumulativePayment, err := ledger.Debit(t.Context(), 50, []core.QuorumID{0, 1, 5})

		require.Error(t, err)
		require.Nil(t, cumulativePayment)

		var quorumNotSupportedError *ondemand.QuorumNotSupportedError
		require.True(t, errors.As(err, &quorumNotSupportedError))

		// verify the store was not updated
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), storedPayment)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestDebitInsufficientFunds")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(100), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// attempt to debit more than total deposits
		cumulativePayment, err := ledger.Debit(t.Context(), 2000, []core.QuorumID{0})
		require.Error(t, err)
		require.Nil(t, cumulativePayment)
		var insufficientFundsError *ondemand.InsufficientFundsError
		require.True(t, errors.As(err, &insufficientFundsError))

		// verify the store was not updated
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), storedPayment)
	})

	t.Run("minimum symbols applied", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestDebitMinimumSymbols")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// debit 5 symbols, but minNumSymbols is 10
		cumulativePayment, err := ledger.Debit(t.Context(), 5, []core.QuorumID{0})
		require.NoError(t, err)
		require.NotNil(t, cumulativePayment)
		require.Equal(t, big.NewInt(10), cumulativePayment)

		// verify the store was updated with minimum charge
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(10), storedPayment)
	})
}

func TestRevertDebit(t *testing.T) {
	t.Run("successful revert", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestRevertDebitSuccessful")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// debit first
		cumulativePayment, err := ledger.Debit(t.Context(), 100, []core.QuorumID{0})
		require.NoError(t, err)
		require.Equal(t, big.NewInt(100), cumulativePayment)

		// verify the store has the debit
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(100), storedPayment)

		// revert the debit
		cumulativePayment, err = ledger.RevertDebit(t.Context(), 50)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(50), cumulativePayment)

		// verify the store was updated after revert
		storedPayment, err = store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(50), storedPayment)
	})

	t.Run("minimum symbols applied", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestRevertDebitMinimum")
		defer cleanup()

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// debit 5 (charged 10 due to minimum)
		cumulativePayment, err := ledger.Debit(t.Context(), 5, []core.QuorumID{0})
		require.NoError(t, err)
		require.Equal(t, big.NewInt(10), cumulativePayment)

		// verify the store has the minimum charge
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(10), storedPayment)

		// revert 5 (should revert 10 due to minimum)
		cumulativePayment, err = ledger.RevertDebit(t.Context(), 5)
		require.NoError(t, err)
		require.Equal(t, 0, cumulativePayment.Cmp(big.NewInt(0)))

		// verify the store was updated to 0
		storedPayment, err = store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), storedPayment)
	})
}

func TestUpdateTotalDeposits(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		ledger, err := ondemand.OnDemandLedgerFromValue(big.NewInt(1000), big.NewInt(1), 10, big.NewInt(0))
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// update to a new value
		err = ledger.UpdateTotalDeposits(big.NewInt(2000))
		require.NoError(t, err)

		// verify the update
		totalDeposits := ledger.GetTotalDeposits()
		require.Equal(t, big.NewInt(2000), totalDeposits)
	})

	t.Run("nil total deposits", func(t *testing.T) {
		ledger, err := ondemand.OnDemandLedgerFromValue(big.NewInt(1000), big.NewInt(1), 10, big.NewInt(0))
		require.NoError(t, err)
		require.NotNil(t, ledger)

		err = ledger.UpdateTotalDeposits(nil)
		require.Error(t, err)
	})

	t.Run("negative total deposits", func(t *testing.T) {
		ledger, err := ondemand.OnDemandLedgerFromValue(big.NewInt(1000), big.NewInt(1), 10, big.NewInt(0))
		require.NoError(t, err)
		require.NotNil(t, ledger)

		err = ledger.UpdateTotalDeposits(big.NewInt(-100))
		require.Error(t, err)
	})
}

func TestOnDemandLedgerFromStore(t *testing.T) {
	t.Run("preexisting store value", func(t *testing.T) {
		store, cleanup := createTestStore(t, "TestFromPreexistingStore")
		defer cleanup()

		// set initial cumulative payment in store
		err := store.StoreCumulativePayment(t.Context(), big.NewInt(500))
		require.NoError(t, err)

		ledger, err := ondemand.OnDemandLedgerFromStore(
			t.Context(), big.NewInt(1000), big.NewInt(1), 10, store)
		require.NoError(t, err)
		require.NotNil(t, ledger)

		// verify ledger works with the initial cumulative payment
		cumulativePayment, err := ledger.Debit(t.Context(), 100, []core.QuorumID{0})
		require.NoError(t, err)
		require.Equal(t, big.NewInt(600), cumulativePayment)

		// verify the store was updated
		storedPayment, err := store.GetCumulativePayment(t.Context())
		require.NoError(t, err)
		require.Equal(t, big.NewInt(600), storedPayment)
	})

	t.Run("nil store", func(t *testing.T) {
		ledger, err := ondemand.OnDemandLedgerFromStore(t.Context(), big.NewInt(1000), big.NewInt(1), 10, nil)
		require.Error(t, err)
		require.Nil(t, ledger)
	})
}

// Creates a payment table and store for testing, returning the store and a cleanup function
func createTestStore(t *testing.T, tableNameSuffix string) (*ondemand.CumulativePaymentStore, func()) {
	tableName := createPaymentTable(t, tableNameSuffix)
	testAccountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	store, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, testAccountID)
	require.NoError(t, err)
	require.NotNil(t, store)

	cleanup := func() {
		deleteTable(t, tableName)
	}

	return store, cleanup
}
