package reservation

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var testTime = time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

func TestNewReservationPaymentValidatorInvalidParams(t *testing.T) {
	ctx := context.Background()
	testVault := vault.NewTestPaymentVault()

	validator, err := NewReservationPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		nil, // nil time source
		OverfillOncePermitted,
		10*time.Second,
		time.Second,
	)
	require.Error(t, err)
	require.Nil(t, validator)
}

func TestDebitMultipleAccounts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	testVault := vault.NewTestPaymentVault()
	testVault.SetGlobalSymbolsPerSecond(1000)
	testVault.SetMinNumSymbols(1)

	testVault.SetReservation(accountA, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 100,
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	testVault.SetReservation(accountB, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 200,
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	mockTimeSource := func() time.Time { return testTime }

	paymentValidator, err := NewReservationPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		mockTimeSource,
		OverfillOncePermitted,
		10*time.Second,
		time.Second,
	)
	require.NoError(t, err)
	require.NotNil(t, paymentValidator)

	err = paymentValidator.Debit(ctx, accountA, uint32(50), []uint8{}, testTime)
	require.NoError(t, err, "first debit from account A should succeed")

	err = paymentValidator.Debit(ctx, accountB, uint32(75), []uint8{}, testTime)
	require.NoError(t, err, "first debit from account B should succeed")

	// should reuse cached ledger
	err = paymentValidator.Debit(ctx, accountA, uint32(25), []uint8{}, testTime)
	require.NoError(t, err, "second debit from account A should succeed")
}

func TestDebitInsufficientCapacity(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetGlobalSymbolsPerSecond(1000)
	testVault.SetMinNumSymbols(1)

	// Set up a small reservation
	testVault.SetReservation(accountID, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 10, // Very low rate
		StartTimestamp:   uint64(testTime.Unix()),
		EndTimestamp:     uint64(testTime.Add(24 * time.Hour).Unix()),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	mockTimeSource := func() time.Time { return testTime }

	paymentValidator, err := NewReservationPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		mockTimeSource,
		OverfillOncePermitted,
		1*time.Second,
		time.Second,
	)
	require.NoError(t, err)

	// First debit exceeding capacity should succeed with OverfillOncePermitted
	err = paymentValidator.Debit(ctx, accountID, uint32(20), []uint8{}, testTime)
	require.NoError(t, err, "first debit should succeed with OverfillOncePermitted even when exceeding capacity")

	// Second debit should fail since bucket is overfilled
	err = paymentValidator.Debit(ctx, accountID, uint32(1), []uint8{}, testTime)
	require.Error(t, err, "second debit should fail when bucket is overfilled")
}

func TestDebitNoReservation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetGlobalSymbolsPerSecond(1000)
	testVault.SetMinNumSymbols(1)

	mockTimeSource := func() time.Time { return testTime }

	paymentValidator, err := NewReservationPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		mockTimeSource,
		OverfillOncePermitted,
		10*time.Second,
		time.Second,
	)
	require.NoError(t, err)

	err = paymentValidator.Debit(ctx, accountID, uint32(10), []uint8{}, testTime)
	require.Error(t, err, "debit should fail when no reservation exists")
}
