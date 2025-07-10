package clients

import (
	"encoding/hex"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create PaymentVaultParams for testing
func createTestPaymentVaultParams(reservationWindow, pricePerSymbol, minNumSymbols uint64) *meterer.PaymentVaultParams {
	quorums := []core.QuorumID{0, 1}
	quorumPaymentConfigs := make(map[core.QuorumID]*core.PaymentQuorumConfig)
	quorumProtocolConfigs := make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig)

	for _, quorumID := range quorums {
		quorumPaymentConfigs[quorumID] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: 2000,
			OnDemandSymbolsPerSecond:    1000,
			OnDemandPricePerSymbol:      pricePerSymbol,
		}

		quorumProtocolConfigs[quorumID] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              minNumSymbols,
			ReservationAdvanceWindow:   10,
			ReservationRateLimitWindow: reservationWindow,
			OnDemandRateLimitWindow:    30,
			OnDemandEnabled:            true,
		}
	}

	return &meterer.PaymentVaultParams{
		QuorumPaymentConfigs:  quorumPaymentConfigs,
		QuorumProtocolConfigs: quorumProtocolConfigs,
		OnDemandQuorumNumbers: quorums,
	}
}

// Helper to create accountant with reservations
func createAccountantWithReservations(symbolsPerSecond uint64) *Accountant {
	privateKey, _ := crypto.GenerateKey()
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey.D.Bytes()))
	accountant := NewAccountant(accountId)

	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: symbolsPerSecond,
			StartTimestamp:   uint64(now.Add(-time.Minute).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: symbolsPerSecond,
			StartTimestamp:   uint64(now.Add(-time.Minute).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}

	err := accountant.SetPaymentState(
		createTestPaymentVaultParams(2, 1, 1),
		reservations,
		big.NewInt(400),
		big.NewInt(500),
		make(meterer.QuorumPeriodRecords),
	)
	if err != nil {
		panic(err) // Test helper, panic on setup failure
	}
	return accountant
}

// Helper to create accountant for on-demand only
func createAccountantOnDemandOnly(balance int64) *Accountant {
	privateKey, _ := crypto.GenerateKey()
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey.D.Bytes()))
	accountant := NewAccountant(accountId)

	err := accountant.SetPaymentState(
		createTestPaymentVaultParams(5, 1, 100),
		map[core.QuorumID]*core.ReservedPayment{},
		big.NewInt(0),
		big.NewInt(balance),
		make(meterer.QuorumPeriodRecords),
	)
	if err != nil {
		panic(err) // Test helper, panic on setup failure
	}
	return accountant
}

func TestAccountBlob_ErrorCases(t *testing.T) {
	acc := createAccountantWithReservations(100)
	now := time.Now().UnixNano()

	// Zero symbols
	header, err := acc.AccountBlob(now, 0, []uint8{0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zero symbols requested")
	assert.Nil(t, header)

	// Empty quorums
	header, err = acc.AccountBlob(now, 50, []uint8{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no quorums provided")
	assert.Nil(t, header)

	// Max symbols (insufficient balance)
	header, err = acc.AccountBlob(now, math.MaxUint64, []uint8{0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient ondemand payment")
	assert.Nil(t, header)
}

func TestAccountBlob_ReservationOverflow(t *testing.T) {
	acc := createAccountantWithReservations(50) // 50 symbols/sec, window=2, so limit=100
	now := time.Now().UnixNano()
	quorums := []uint8{0, 1}

	// First call: 80 symbols (within limit)
	header, err := acc.AccountBlob(now, 80, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment) // reservation

	// Second call: 30 symbols (80+30=110, should overflow)
	header, err = acc.AccountBlob(now, 30, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), header.CumulativePayment) // still reservation

	// Third call: should use on-demand
	header, err = acc.AccountBlob(now, 20, quorums)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(420), header.CumulativePayment) // 400 + 20
}

func TestAccountBlob_OnDemandOnly(t *testing.T) {
	acc := createAccountantOnDemandOnly(1500)
	now := time.Now().UnixNano()

	header, err := acc.AccountBlob(now, 1500, []uint8{0, 1})
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1500), header.CumulativePayment)
}

func TestAccountBlob_InsufficientBalance(t *testing.T) {
	acc := createAccountantOnDemandOnly(500)
	now := time.Now().UnixNano()

	header, err := acc.AccountBlob(now, 2000, []uint8{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient ondemand payment")
	assert.Nil(t, header)
}

func TestAccountBlob_BinRotation(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(1)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	now := time.Now()
	// Start the reservation 1 hour in the past to allow for past period access
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Add(-time.Hour).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Add(-time.Hour).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}
	err = accountant.SetPaymentState(
		createTestPaymentVaultParams(reservationWindow, 1, 100),
		reservations,
		big.NewInt(0),    // Start with no cumulative payment used
		big.NewInt(2000), // On-chain deposit sufficient for potential on-demand
		make(meterer.QuorumPeriodRecords),
	)
	require.NoError(t, err)

	quorums := []uint8{0, 1}

	// Use the same base time to ensure consistency with reservations
	baseTime := now.UnixNano()

	// First call - use current period
	currentPeriod := payment_logic.GetReservationPeriodByNanosecond(baseTime, reservationWindow)
	_, err = accountant.AccountBlob(baseTime, 800, quorums)
	assert.NoError(t, err)

	// Check bin 0 has usage 800 for charged quorums
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
	}

	// Second call - use previous period (which should be allowed by validation)
	prevTime := baseTime - int64(reservationWindow)*time.Second.Nanoseconds()
	prevPeriod := payment_logic.GetReservationPeriodByNanosecond(prevTime, reservationWindow)
	_, err = accountant.AccountBlob(prevTime, 300, quorums)
	assert.NoError(t, err)

	// Check previous period has usage 300
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(prevPeriod, quorumNumber)
		assert.Equal(t, uint64(300), record.Usage)
	}

	// Third call - same period as second call
	_, err = accountant.AccountBlob(prevTime, 500, quorums)
	assert.NoError(t, err)

	// Check previous period now has usage 800; overflow period is 0
	overflowPeriod := payment_logic.GetOverflowPeriod(currentPeriod, reservationWindow)
	for _, quorumNumber := range quorums {
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(prevPeriod, quorumNumber)
		assert.Equal(t, uint64(800), record.Usage)
		record = accountant.periodRecords.GetRelativePeriodRecord(overflowPeriod, quorumNumber)
		assert.Equal(t, uint64(0), record.Usage)
	}
}

func TestAccountant_Concurrent(t *testing.T) {
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	reservationWindow := uint64(1)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}
	err = accountant.SetPaymentState(
		createTestPaymentVaultParams(reservationWindow, 1, 100),
		reservations,
		big.NewInt(1000),
		big.NewInt(1000),
		make(meterer.QuorumPeriodRecords),
	)
	require.NoError(t, err)

	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	nowNano := time.Now().UnixNano()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := accountant.AccountBlob(nowNano, 100, quorums)
			assert.NoError(t, err)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state
	for _, quorumNumber := range quorums {
		currentPeriod := payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow)
		record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, quorumNumber)
		assert.Equal(t, uint64(1000), record.Usage)
	}
}

// This test is now covered by TestAccountBlob_ReservationOverflow

// This test is now covered by TestAccountBlob_ReservationOverflow

func TestAccountant_SetPaymentState(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	acc := NewAccountant(accountID)

	tests := []struct {
		name    string
		state   *v2.GetPaymentStateForAllQuorumsReply
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil payment state",
			state:   nil,
			wantErr: true,
			errMsg:  "payment state cannot be nil",
		},
		{
			name:    "nil payment vault params",
			state:   &v2.GetPaymentStateForAllQuorumsReply{},
			wantErr: true,
			errMsg:  "payment vault params cannot be nil",
		},
		{
			name: "successful state update",
			state: &v2.GetPaymentStateForAllQuorumsReply{
				PaymentVaultParams: &v2.PaymentVaultParams{
					QuorumPaymentConfigs: map[uint32]*v2.PaymentQuorumConfig{
						0: {
							ReservationSymbolsPerSecond: 100,
							OnDemandSymbolsPerSecond:    200,
							OnDemandPricePerSymbol:      10,
						},
					},
					QuorumProtocolConfigs: map[uint32]*v2.PaymentQuorumProtocolConfig{
						0: {
							MinNumSymbols:              1,
							ReservationAdvanceWindow:   2,
							ReservationRateLimitWindow: 3,
							OnDemandRateLimitWindow:    4,
							OnDemandEnabled:            true,
						},
					},
				},
				Reservations: map[uint32]*v2.QuorumReservation{
					0: {
						SymbolsPerSecond: 100,
						StartTimestamp:   1000,
						EndTimestamp:     2000,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.state == nil {
				// Test nil case directly since the conversion function expects non-nil
				_, _, _, _, _, err := meterer.ConvertPaymentStateFromProtobuf(tt.state)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			// Convert protobuf to native types
			paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords, err := meterer.ConvertPaymentStateFromProtobuf(tt.state)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)

			// Now test SetPaymentState with converted parameters
			err = acc.SetPaymentState(paymentVaultParams, reservations, cumulativePayment, onchainCumulativePayment, periodRecords)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, acc.paymentVaultParams.QuorumPaymentConfigs[0])
				assert.NotNil(t, acc.paymentVaultParams.QuorumProtocolConfigs[0])
				assert.NotNil(t, acc.reservations[0])
			}
		})
	}
}

func TestAccountant_ReservationUsage(t *testing.T) {
	acc := createAccountantWithReservations(200)
	now := time.Now().UnixNano()

	// Success case
	err := acc.reservationUsage(50, []core.QuorumID{0, 1}, now)
	assert.NoError(t, err)

	// Invalid quorum
	err = acc.reservationUsage(50, []core.QuorumID{0, 2}, now)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum number mismatch")

	// Limit exceeded - 100 symbols/sec * 2 sec window = 200 limit, so 201 should fail
	acc2 := createAccountantWithReservations(100)
	err = acc2.reservationUsage(201, []core.QuorumID{0}, now)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds bin limit")
}

func TestAccountant_OnDemandUsage(t *testing.T) {
	acc := createAccountantOnDemandOnly(1000)

	// Success case
	payment, err := acc.onDemandUsage(50, []uint8{0, 1})
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.Cmp(big.NewInt(0)) > 0)

	// Insufficient balance
	payment, err = acc.onDemandUsage(2000, []uint8{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient ondemand payment")
	assert.Nil(t, payment)

	// Invalid quorum
	payment, err = acc.onDemandUsage(50, []uint8{2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum number mismatch: 2")
	assert.Nil(t, payment)
}

func TestAccountant_MixedReservationStates(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x123")
	now := time.Now()
	acc := NewAccountant(accountID)

	// Set up payment vault params
	vaultParams := &meterer.PaymentVaultParams{
		QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}
	for i := core.QuorumID(0); i < 3; i++ {
		vaultParams.QuorumProtocolConfigs[i] = &core.PaymentQuorumProtocolConfig{
			ReservationRateLimitWindow: 1,
			OnDemandRateLimitWindow:    1,
			MinNumSymbols:              1,
		}
		vaultParams.QuorumPaymentConfigs[i] = &core.PaymentQuorumConfig{
			ReservationSymbolsPerSecond: 0,
			OnDemandSymbolsPerSecond:    100,
			OnDemandPricePerSymbol:      1,
		}
	}

	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint64(now.Add(time.Hour * 1).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour * 2).Unix()), // Future
		},
		1: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint64(now.Add(time.Hour * -2).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour * -1).Unix()), // Expired
		},
		2: {
			SymbolsPerSecond: 100,
			StartTimestamp:   uint64(now.Add(time.Hour * -1).Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour * 1).Unix()), // Active
		},
	}

	err := acc.SetPaymentState(
		vaultParams,
		reservations,
		big.NewInt(0),    // Start with no cumulative payment used
		big.NewInt(1000), // On-chain deposit of 1000
		make(meterer.QuorumPeriodRecords),
	)
	require.NoError(t, err)

	// Reservations and OnDemand are not sufficient for all three quorums
	payment, err := acc.AccountBlob(now.UnixNano(), 50, []uint8{0, 1, 2})
	assert.Nil(t, payment)
	assert.Contains(t, err.Error(), "cannot create payment information")

	// Separate reservation dispersal is sufficient for quorum 2
	payment, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{2})
	// 1749697512 1749701112 1749693912.770014000
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.CumulativePayment.Cmp(big.NewInt(0)) == 0)

	// Alternatively use ondemand for quorum 0 or/and 1
	payment, err = acc.AccountBlob(now.UnixNano(), 50, []uint8{0, 1})
	require.NoError(t, err)
	assert.NotNil(t, payment)
	assert.True(t, payment.CumulativePayment.Cmp(big.NewInt(0)) > 0)
}

func TestAccountant_ReservationRollback(t *testing.T) {
	now := time.Now()
	reservation := &core.ReservedPayment{
		SymbolsPerSecond: 50,
		StartTimestamp:   uint64(now.Unix()),
		EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
	}
	reservationWindow := uint64(2)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountId := gethcommon.HexToAddress(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(accountId)

	// Create payment state with test configurations
	vaultParams := createTestPaymentVaultParams(reservationWindow, 1, 1)

	reservations := map[core.QuorumID]*core.ReservedPayment{0: reservation, 1: reservation}
	err = accountant.SetPaymentState(
		vaultParams,
		reservations,
		big.NewInt(0),
		big.NewInt(0),
		make(meterer.QuorumPeriodRecords),
	)
	require.NoError(t, err)

	// Test rollback when a later quorum fails
	nowNano := time.Now().UnixNano()
	currentPeriod := payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow)

	// First update should succeed
	moreUsedQuorum := uint8(1)
	lessUsedQuorum := uint8(0)
	_, err = accountant.AccountBlob(nowNano, 50, []uint8{moreUsedQuorum})
	assert.NoError(t, err)

	// Verify first quorum was updated
	record := accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(50), record.Usage)

	// Use both quorums, more used quorum overflows
	_, err = accountant.AccountBlob(nowNano, 60, []uint8{moreUsedQuorum, lessUsedQuorum})
	assert.NoError(t, err)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(100), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(payment_logic.GetOverflowPeriod(currentPeriod, reservationWindow), moreUsedQuorum)
	assert.Equal(t, uint64(10), record.Usage)

	// Use both quorums, more used quorum cannot overflow again
	_, err = accountant.AccountBlob(nowNano, 60, []uint8{moreUsedQuorum, lessUsedQuorum})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reservation limit exceeded")

	// No reservation updates were made
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, moreUsedQuorum)
	assert.Equal(t, uint64(100), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
	record = accountant.periodRecords.GetRelativePeriodRecord(payment_logic.GetOverflowPeriod(currentPeriod, reservationWindow), moreUsedQuorum)
	assert.Equal(t, uint64(10), record.Usage)

	// Test rollback when a quorum doesn't exist
	_, err = accountant.AccountBlob(nowNano, 50, []uint8{lessUsedQuorum, 2})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum number mismatch")
	// quorum usage rolled back
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)

	// Test rollback when config is missing
	// Remove config for quorum 1
	delete(accountant.paymentVaultParams.QuorumProtocolConfigs, 1)
	_, err = accountant.AccountBlob(nowNano, 50, []uint8{0, 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum config not found")
	// quorum usage stays the same
	record = accountant.periodRecords.GetRelativePeriodRecord(currentPeriod, lessUsedQuorum)
	assert.Equal(t, uint64(60), record.Usage)
}
