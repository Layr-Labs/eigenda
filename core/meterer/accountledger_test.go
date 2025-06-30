package meterer

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccountLedgerInterface verifies that LocalAccountLedger implements AccountLedger
func TestAccountLedgerInterface(t *testing.T) {
	var _ AccountLedger = &LocalAccountLedger{}
}

// TestNewLocalAccountLedger tests the creation of a new LocalAccountLedger
func TestNewLocalAccountLedger(t *testing.T) {
	ledger := NewLocalAccountLedger()
	assert.NotNil(t, ledger)

	state := ledger.GetAccountState()
	assert.NotNil(t, state.Reservations)
	assert.Empty(t, state.Reservations)
	assert.NotNil(t, state.OnDemand)
	assert.Equal(t, big.NewInt(0), state.OnDemand.CumulativePayment)
	assert.NotNil(t, state.PeriodRecords)
	assert.Empty(t, state.PeriodRecords)
	assert.NotNil(t, state.CumulativePayment)
	assert.Equal(t, big.NewInt(0), state.CumulativePayment)
}

// TestSetAccountState tests setting complete account state
func TestSetAccountState(t *testing.T) {
	ledger := NewLocalAccountLedger()

	// Test setting nil state
	ledger.SetAccountState(AccountState{})
	state := ledger.GetAccountState()
	assert.NotNil(t, state.Reservations)
	assert.Empty(t, state.Reservations)
	assert.NotNil(t, state.OnDemand)
	assert.Equal(t, big.NewInt(0), state.OnDemand.CumulativePayment)

	// Test setting complete state
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(time.Now().Unix()),
			EndTimestamp:     uint64(time.Now().Add(time.Hour).Unix()),
		},
	}
	onDemand := &core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	periodRecords := make(QuorumPeriodRecords)
	periodRecords[0] = []*PeriodRecord{{Index: 1, Usage: 100}}
	cumulativePayment := big.NewInt(250)

	inputState := AccountState{
		Reservations:      reservations,
		OnDemand:          onDemand,
		PeriodRecords:     periodRecords,
		CumulativePayment: cumulativePayment,
	}

	ledger.SetAccountState(inputState)
	resultState := ledger.GetAccountState()

	// Verify reservations
	assert.Len(t, resultState.Reservations, 1)
	assert.Equal(t, uint64(1000), resultState.Reservations[0].SymbolsPerSecond)

	// Verify on-demand
	assert.Equal(t, big.NewInt(500), resultState.OnDemand.CumulativePayment)

	// Verify period records
	assert.Len(t, resultState.PeriodRecords, 1)
	assert.Equal(t, uint64(100), resultState.PeriodRecords[0][0].Usage)

	// Verify cumulative payment
	assert.Equal(t, big.NewInt(250), resultState.CumulativePayment)

	// Verify deep copying - modifications to input should not affect ledger
	reservations[0].SymbolsPerSecond = 2000
	onDemand.CumulativePayment.SetInt64(1000)
	cumulativePayment.SetInt64(500)

	finalState := ledger.GetAccountState()
	assert.Equal(t, uint64(1000), finalState.Reservations[0].SymbolsPerSecond)
	assert.Equal(t, big.NewInt(500), finalState.OnDemand.CumulativePayment)
	assert.Equal(t, big.NewInt(250), finalState.CumulativePayment)
}

// TestGetAccountState tests retrieving account state with deep copying
func TestGetAccountState(t *testing.T) {
	ledger := NewLocalAccountLedger()

	// Set up initial state
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {SymbolsPerSecond: 1000, StartTimestamp: 100, EndTimestamp: 200},
	}
	onDemand := &core.OnDemandPayment{CumulativePayment: big.NewInt(500)}

	ledger.SetAccountState(AccountState{
		Reservations:      reservations,
		OnDemand:          onDemand,
		CumulativePayment: big.NewInt(250),
	})

	// Get state and verify deep copying
	state1 := ledger.GetAccountState()
	state2 := ledger.GetAccountState()

	// Modify state1 - should not affect state2 or ledger
	state1.Reservations[0].SymbolsPerSecond = 2000
	state1.OnDemand.CumulativePayment.SetInt64(1000)
	state1.CumulativePayment.SetInt64(500)

	// Verify state2 is unchanged
	assert.Equal(t, uint64(1000), state2.Reservations[0].SymbolsPerSecond)
	assert.Equal(t, big.NewInt(500), state2.OnDemand.CumulativePayment)
	assert.Equal(t, big.NewInt(250), state2.CumulativePayment)

	// Verify ledger is unchanged
	currentState := ledger.GetAccountState()
	assert.Equal(t, uint64(1000), currentState.Reservations[0].SymbolsPerSecond)
	assert.Equal(t, big.NewInt(500), currentState.OnDemand.CumulativePayment)
	assert.Equal(t, big.NewInt(250), currentState.CumulativePayment)
}

// createTestPaymentVaultParams creates test payment vault parameters
func createTestPaymentVaultParams() *PaymentVaultParams {
	return &PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {
				ReservationSymbolsPerSecond: 1000,
				OnDemandSymbolsPerSecond:    500,
				OnDemandPricePerSymbol:      1,
			},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {
				MinNumSymbols:              1,
				ReservationAdvanceWindow:   10,
				ReservationRateLimitWindow: 5,
				OnDemandRateLimitWindow:    5,
				OnDemandEnabled:            true,
			},
		},
		OnDemandQuorumNumbers: []uint8{0},
	}
}

// TestRecordReservationUsage tests recording reservation usage
func TestRecordReservationUsage(t *testing.T) {
	ledger := NewLocalAccountLedger()
	params := createTestPaymentVaultParams()
	accountID := gethcommon.HexToAddress("0x123")

	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}

	ledger.SetAccountState(AccountState{
		Reservations: reservations,
		OnDemand:     &core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
	})

	// Test successful usage recording
	err := ledger.RecordReservationUsage(
		context.Background(),
		accountID,
		now.UnixNano(),
		100,
		[]core.QuorumID{0},
		params,
	)
	assert.NoError(t, err)

	// Verify usage was recorded
	state := ledger.GetAccountState()
	assert.NotEmpty(t, state.PeriodRecords[0])

	// Test usage with non-existent reservation
	err = ledger.RecordReservationUsage(
		context.Background(),
		accountID,
		now.UnixNano(),
		100,
		[]core.QuorumID{1}, // No reservation for quorum 1
		params,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quorum number mismatch")
}

// TestRecordReservationUsage_Rollback tests that multi-quorum operations rollback on failure
func TestRecordReservationUsage_Rollback(t *testing.T) {
	ledger := NewLocalAccountLedger()
	params := createTestPaymentVaultParams()
	accountID := gethcommon.HexToAddress("0x123")

	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}

	ledger.SetAccountState(AccountState{
		Reservations: reservations,
		OnDemand:     &core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
	})

	// Record some initial usage
	err := ledger.RecordReservationUsage(
		context.Background(),
		accountID,
		now.UnixNano(),
		100,
		[]core.QuorumID{0},
		params,
	)
	require.NoError(t, err)

	initialState := ledger.GetAccountState()
	// Calculate the correct period index and relative index for the timestamp
	currentPeriod := GetReservationPeriodByNanosecond(now.UnixNano(), params.QuorumProtocolConfigs[0].ReservationRateLimitWindow)
	relativeIndex := currentPeriod % uint64(MinNumBins)
	initialUsage := initialState.PeriodRecords[0][relativeIndex].Usage

	// Attempt to record usage for multiple quorums where one fails
	err = ledger.RecordReservationUsage(
		context.Background(),
		accountID,
		now.UnixNano(),
		100,
		[]core.QuorumID{0, 1}, // Quorum 1 has no reservation
		params,
	)
	assert.Error(t, err)

	// Verify that no changes were made (rollback occurred)
	finalState := ledger.GetAccountState()
	assert.Equal(t, initialUsage, finalState.PeriodRecords[0][relativeIndex].Usage)
}

// TestRecordOnDemandUsage tests recording on-demand usage
func TestRecordOnDemandUsage(t *testing.T) {
	ledger := NewLocalAccountLedger()
	params := createTestPaymentVaultParams()
	accountID := gethcommon.HexToAddress("0x123")

	// Set up on-demand payment
	ledger.SetAccountState(AccountState{
		Reservations:      make(map[core.QuorumID]*core.ReservedPayment),
		OnDemand:          &core.OnDemandPayment{CumulativePayment: big.NewInt(1000)},
		CumulativePayment: big.NewInt(0),
	})

	// Test successful on-demand usage
	newPayment, err := ledger.RecordOnDemandUsage(
		context.Background(),
		accountID,
		100,
		[]core.QuorumID{0},
		params,
	)
	assert.NoError(t, err)
	assert.NotNil(t, newPayment)
	assert.Equal(t, big.NewInt(100), newPayment) // 100 symbols * 1 price = 100

	// Verify cumulative payment was updated
	state := ledger.GetAccountState()
	assert.Equal(t, big.NewInt(100), state.CumulativePayment)

	// Test insufficient balance
	_, err = ledger.RecordOnDemandUsage(
		context.Background(),
		accountID,
		1000, // This would require 1000 payment, but only 900 remaining
		[]core.QuorumID{0},
		params,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient ondemand payment")

	// Test with invalid quorum
	_, err = ledger.RecordOnDemandUsage(
		context.Background(),
		accountID,
		100,
		[]core.QuorumID{1}, // Not in OnDemandQuorumNumbers
		params,
	)
	assert.Error(t, err)
}

// TestConcurrentAccess tests thread safety of LocalAccountLedger
func TestConcurrentAccess(t *testing.T) {
	ledger := NewLocalAccountLedger()
	params := createTestPaymentVaultParams()
	accountID := gethcommon.HexToAddress("0x123")

	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 10000,
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}

	ledger.SetAccountState(AccountState{
		Reservations:      reservations,
		OnDemand:          &core.OnDemandPayment{CumulativePayment: big.NewInt(10000)},
		CumulativePayment: big.NewInt(0),
	})

	const numGoroutines = 10
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	// Launch concurrent reservation usage operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				err := ledger.RecordReservationUsage(
					context.Background(),
					accountID,
					now.UnixNano(),
					1,
					[]core.QuorumID{0},
					params,
				)
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	// Launch concurrent on-demand usage operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				_, err := ledger.RecordOnDemandUsage(
					context.Background(),
					accountID,
					1,
					[]core.QuorumID{0},
					params,
				)
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	// Launch concurrent read operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				state := ledger.GetAccountState()
				assert.NotNil(t, state.Reservations)
				assert.NotNil(t, state.OnDemand)
				assert.NotNil(t, state.CumulativePayment)
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent operation failed: %v", err)
	}

	// Verify final state consistency
	finalState := ledger.GetAccountState()
	assert.NotNil(t, finalState.CumulativePayment)
	assert.True(t, finalState.CumulativePayment.Cmp(big.NewInt(0)) >= 0)
	assert.True(t, finalState.CumulativePayment.Cmp(big.NewInt(10000)) <= 0)
}

// TestAccountState_DeepCopy tests that AccountState fields are properly deep copied
func TestAccountState_DeepCopy(t *testing.T) {
	ledger := NewLocalAccountLedger()

	// Create complex state with multiple elements
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {SymbolsPerSecond: 1000, StartTimestamp: 100, EndTimestamp: 200},
		1: {SymbolsPerSecond: 2000, StartTimestamp: 150, EndTimestamp: 250},
	}
	onDemand := &core.OnDemandPayment{CumulativePayment: big.NewInt(500)}
	periodRecords := make(QuorumPeriodRecords)
	periodRecords[0] = []*PeriodRecord{{Index: 1, Usage: 100}, {Index: 2, Usage: 200}}
	periodRecords[1] = []*PeriodRecord{{Index: 1, Usage: 150}}

	ledger.SetAccountState(AccountState{
		Reservations:      reservations,
		OnDemand:          onDemand,
		PeriodRecords:     periodRecords,
		CumulativePayment: big.NewInt(250),
	})

	// Get state and modify it
	state := ledger.GetAccountState()

	// Modify reservations
	state.Reservations[0].SymbolsPerSecond = 9999
	delete(state.Reservations, 1)

	// Modify on-demand
	state.OnDemand.CumulativePayment.SetInt64(9999)

	// Modify period records
	state.PeriodRecords[0][0].Usage = 9999

	// Modify cumulative payment
	state.CumulativePayment.SetInt64(9999)

	// Verify original ledger state is unchanged
	originalState := ledger.GetAccountState()
	assert.Equal(t, uint64(1000), originalState.Reservations[0].SymbolsPerSecond)
	assert.Contains(t, originalState.Reservations, core.QuorumID(1))
	assert.Equal(t, big.NewInt(500), originalState.OnDemand.CumulativePayment)
	assert.Equal(t, uint64(100), originalState.PeriodRecords[0][0].Usage)
	assert.Equal(t, big.NewInt(250), originalState.CumulativePayment)
}

// TestAccountState_NilHandling tests handling of nil fields in AccountState
func TestAccountState_NilHandling(t *testing.T) {
	ledger := NewLocalAccountLedger()

	// Test with all nil fields
	ledger.SetAccountState(AccountState{
		Reservations:      nil,
		OnDemand:          nil,
		PeriodRecords:     nil,
		CumulativePayment: nil,
	})

	state := ledger.GetAccountState()
	assert.NotNil(t, state.Reservations)
	assert.Empty(t, state.Reservations)
	assert.NotNil(t, state.OnDemand)
	assert.Equal(t, big.NewInt(0), state.OnDemand.CumulativePayment)
	assert.NotNil(t, state.PeriodRecords)
	assert.Empty(t, state.PeriodRecords)
	assert.NotNil(t, state.CumulativePayment)
	assert.Equal(t, big.NewInt(0), state.CumulativePayment)
}

// BenchmarkLocalAccountLedger_GetAccountState benchmarks the GetAccountState operation
func BenchmarkLocalAccountLedger_GetAccountState(b *testing.B) {
	ledger := NewLocalAccountLedger()

	// Set up complex state
	reservations := make(map[core.QuorumID]*core.ReservedPayment)
	periodRecords := make(QuorumPeriodRecords)

	for i := 0; i < 10; i++ {
		reservations[core.QuorumID(i)] = &core.ReservedPayment{
			SymbolsPerSecond: 1000,
			StartTimestamp:   100,
			EndTimestamp:     200,
		}

		periodRecords[core.QuorumID(i)] = make([]*PeriodRecord, MinNumBins)
		for j := 0; j < int(MinNumBins); j++ {
			periodRecords[core.QuorumID(i)][j] = &PeriodRecord{
				Index: uint32(j),
				Usage: uint64(j * 100),
			}
		}
	}

	ledger.SetAccountState(AccountState{
		Reservations:      reservations,
		OnDemand:          &core.OnDemandPayment{CumulativePayment: big.NewInt(1000)},
		PeriodRecords:     periodRecords,
		CumulativePayment: big.NewInt(500),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state := ledger.GetAccountState()
		_ = state
	}
}

// BenchmarkLocalAccountLedger_RecordReservationUsage benchmarks reservation usage recording
func BenchmarkLocalAccountLedger_RecordReservationUsage(b *testing.B) {
	ledger := NewLocalAccountLedger()
	params := createTestPaymentVaultParams()
	accountID := gethcommon.HexToAddress("0x123")

	now := time.Now()
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000000, // Large limit to avoid hitting limits
			StartTimestamp:   uint64(now.Unix()),
			EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
		},
	}

	ledger.SetAccountState(AccountState{
		Reservations: reservations,
		OnDemand:     &core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ledger.RecordReservationUsage(
			context.Background(),
			accountID,
			now.UnixNano(),
			1,
			[]core.QuorumID{0},
			params,
		)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
