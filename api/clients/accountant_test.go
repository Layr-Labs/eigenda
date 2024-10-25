package clients

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountant(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  100,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint32(6)
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)

	assert.NotNil(t, accountant)
	assert.Equal(t, reservation, accountant.reservation)
	assert.Equal(t, onDemand, accountant.onDemand)
	assert.Equal(t, reservationWindow, accountant.reservationWindow)
	assert.Equal(t, pricePerSymbol, accountant.pricePerSymbol)
	assert.Equal(t, minNumSymbols, accountant.minNumSymbols)
	assert.Equal(t, []uint64{0, 0, 0}, accountant.binUsages)
	assert.Equal(t, big.NewInt(0), accountant.cumulativePayment)
}

func TestAccountBlob_Reservation(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  200,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint32(5)
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	dataLength := uint64(500)
	quorums := []uint8{0, 1}

	header, _, err := accountant.AccountBlob(ctx, dataLength, quorums)
	metadata := core.ConvertPaymentHeader(header)

	assert.NoError(t, err)
	assert.Equal(t, meterer.GetBinIndex(uint64(time.Now().Unix()), reservationWindow), header.BinIndex)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)
	assert.Equal(t, []uint64{500, 0, 0}, accountant.binUsages)

	dataLength = uint64(700)

	header, _, err = accountant.AccountBlob(ctx, dataLength, quorums)
	metadata = core.ConvertPaymentHeader(header)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, header.BinIndex)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)
	assert.Equal(t, []uint64{1200, 0, 200}, accountant.binUsages)

	// Second call should use on-demand payment
	header, _, err = accountant.AccountBlob(ctx, 300, quorums)
	metadata = core.ConvertPaymentHeader(header)

	assert.NoError(t, err)
	assert.Equal(t, uint32(0), header.BinIndex)
	assert.Equal(t, big.NewInt(3), metadata.CumulativePayment)
}

func TestAccountBlob_OnDemand(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  200,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint32(5)
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	dataLength := uint64(1500)
	quorums := []uint8{0, 1}

	header, _, err := accountant.AccountBlob(ctx, dataLength, quorums)
	metadata := core.ConvertPaymentHeader(header)
	expectedPayment := big.NewInt(int64(dataLength * uint64(pricePerSymbol) / uint64(minNumSymbols)))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), header.BinIndex)
	assert.Equal(t, expectedPayment, metadata.CumulativePayment)
	assert.Equal(t, []uint64{0, 0, 0}, accountant.binUsages)
	assert.Equal(t, expectedPayment, accountant.cumulativePayment)
}

func TestAccountBlob_InsufficientOnDemand(t *testing.T) {
	reservation := core.ActiveReservation{}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(500),
	}
	reservationWindow := uint32(60)
	pricePerSymbol := uint32(100)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	dataLength := uint64(2000)
	quorums := []uint8{0, 1}

	_, _, err = accountant.AccountBlob(ctx, dataLength, quorums)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Accountant cannot approve payment for this blob")
}

func TestAccountBlobCallSeries(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  200,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(5)
	pricePerSymbol := uint32(100)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}
	now := time.Now().Unix()

	// First call: Use reservation
	header, _, err := accountant.AccountBlob(ctx, 800, quorums)
	metadata := core.ConvertPaymentHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetBinIndex(uint64(now), reservationWindow), header.BinIndex)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)

	// Second call: Use remaining reservation + overflow
	header, _, err = accountant.AccountBlob(ctx, 300, quorums)
	metadata = core.ConvertPaymentHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetBinIndex(uint64(now), reservationWindow), header.BinIndex)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)

	// Third call: Use on-demand
	header, _, err = accountant.AccountBlob(ctx, 500, quorums)
	metadata = core.ConvertPaymentHeader(header)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), header.BinIndex)
	assert.Equal(t, big.NewInt(500), metadata.CumulativePayment)

	// Fourth call: Insufficient on-demand
	_, _, err = accountant.AccountBlob(ctx, 600, quorums)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Accountant cannot approve payment for this blob")
}

func TestAccountBlob_BinRotation(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  1000,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	_, _, err = accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Wait for bin rotation
	time.Sleep(2 * time.Second)

	// Second call after bin rotation
	_, _, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{300, 0, 0}, accountant.binUsages)

	// Third call
	_, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)
}

func TestBinRotation(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  1000,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// First call
	_, _, err = accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Second call for overflow
	_, _, err = accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1600, 0, 600}, accountant.binUsages)

	// Wait for bin rotation
	time.Sleep(1200 * time.Millisecond)

	_, _, err = accountant.AccountBlob(ctx, 300, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{300, 600, 0}, accountant.binUsages)

	// another bin rotation
	time.Sleep(1200 * time.Millisecond)

	_, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1100, 0, 100}, accountant.binUsages)
}

func TestConcurrentBinRotationAndAccountBlob(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  1000,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// Start concurrent AccountBlob calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, _, err := accountant.AccountBlob(ctx, 100, quorums)
				assert.NoError(t, err)
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check final state
	assert.Equal(t, uint64(1000), accountant.binUsages[0]+accountant.binUsages[1]+accountant.binUsages[2])
}

func TestAccountBlob_ReservationWithOneOverflow(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  200,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(5)
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()
	ctx := context.Background()
	quorums := []uint8{0, 1}
	now := time.Now().Unix()

	// Okay reservation
	header, _, err := accountant.AccountBlob(ctx, 800, quorums)
	assert.NoError(t, err)
	assert.Equal(t, meterer.GetBinIndex(uint64(now), reservationWindow), header.BinIndex)
	metadata := core.ConvertPaymentHeader(header)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)
	assert.Equal(t, []uint64{800, 0, 0}, accountant.binUsages)

	// Second call: Allow one overflow
	header, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	metadata = core.ConvertPaymentHeader(header)
	assert.Equal(t, big.NewInt(0), metadata.CumulativePayment)
	assert.Equal(t, []uint64{1300, 0, 300}, accountant.binUsages)

	// Third call: Should use on-demand payment
	header, _, err = accountant.AccountBlob(ctx, 200, quorums)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), header.BinIndex)
	metadata = core.ConvertPaymentHeader(header)
	assert.Equal(t, big.NewInt(2), metadata.CumulativePayment)
	assert.Equal(t, []uint64{1300, 0, 300}, accountant.binUsages)
}

func TestAccountBlob_ReservationOverflowReset(t *testing.T) {
	reservation := core.ActiveReservation{
		SymbolsPerSec:  1000,
		StartTimestamp: 100,
		EndTimestamp:   200,
		QuorumSplit:    []byte{50, 50},
		QuorumNumbers:  []uint8{0, 1},
	}
	onDemand := core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	reservationWindow := uint32(1) // Set to 1 second for testing
	pricePerSymbol := uint32(1)
	minNumSymbols := uint32(100)

	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	paymentSigner := auth.NewPaymentSigner(hex.EncodeToString(privateKey1.D.Bytes()))
	accountant := NewAccountant(reservation, onDemand, reservationWindow, pricePerSymbol, minNumSymbols, paymentSigner)
	defer accountant.Stop()

	ctx := context.Background()
	quorums := []uint8{0, 1}

	// full reservation
	_, _, err = accountant.AccountBlob(ctx, 1000, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1000, 0, 0}, accountant.binUsages)

	// no overflow
	header, _, err := accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1000, 0, 0}, accountant.binUsages)
	metadata := core.ConvertPaymentHeader(header)
	assert.Equal(t, big.NewInt(5), metadata.CumulativePayment)

	// Wait for bin rotation
	time.Sleep(1500 * time.Millisecond)

	// Third call: Should use new bin and allow overflow again
	header, _, err = accountant.AccountBlob(ctx, 500, quorums)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{500, 0, 0}, accountant.binUsages)
}
