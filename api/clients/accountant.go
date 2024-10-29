package clients

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
)

type IAccountant interface {
	AccountBlob(ctx context.Context, data []byte, quorums []uint8) (uint32, uint64, error)
}

type Accountant struct {
	// on-chain states
	reservation       core.ActiveReservation
	onDemand          core.OnDemandPayment
	reservationWindow uint32
	pricePerSymbol    uint32
	minNumSymbols     uint32

	// local accounting
	// contains 3 bins; index 0 for current bin, 1 for next bin, 2 for overflowed bin
	binRecords        []BinRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int
	stopRotation      chan struct{}

	paymentSigner core.PaymentSigner
}

type BinRecord struct {
	Index uint32
	Usage uint64
}

func NewAccountant(reservation core.ActiveReservation, onDemand core.OnDemandPayment, reservationWindow uint32, pricePerSymbol uint32, minNumSymbols uint32, paymentSigner core.PaymentSigner) *Accountant {
	//TODO: client storage; currently every instance starts fresh but on-chain or a small store makes more sense
	// Also client is currently responsible for supplying network params, we need to add RPC in order to be automatic
	// There's a subsequent PR that handles populating the accountant with on-chain state from the disperser
	a := Accountant{
		reservation:       reservation,
		onDemand:          onDemand,
		reservationWindow: reservationWindow,
		pricePerSymbol:    pricePerSymbol,
		minNumSymbols:     minNumSymbols,
		binRecords:        []BinRecord{{Index: 0, Usage: 0}, {Index: 1, Usage: 0}, {Index: 2, Usage: 0}},
		cumulativePayment: big.NewInt(0),
		stopRotation:      make(chan struct{}),
		paymentSigner:     paymentSigner,
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// accountant calculates and records payment information
func (a *Accountant) BlobPaymentInfo(ctx context.Context, dataLength uint64) (uint32, *big.Int, error) {
	now := time.Now().Unix()
	currentBinIndex := meterer.GetBinIndex(uint64(now), a.reservationWindow)
	// index := time.Now().Unix() / int64(a.reservationWindow)

	a.usageLock.Lock()
	defer a.usageLock.Unlock()
	relativeBinRecord := a.GetRelativeBinRecord(currentBinIndex)
	relativeBinRecord.Usage += dataLength

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSec * uint64(a.reservationWindow)
	if relativeBinRecord.Usage <= binLimit {
		return currentBinIndex, big.NewInt(0), nil
	}

	overflowBinRecord := a.GetOverflowBinRecord(currentBinIndex)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowBinRecord.Usage == 0 && relativeBinRecord.Usage-dataLength < binLimit && dataLength <= binLimit {
		overflowBinRecord.Usage += relativeBinRecord.Usage - binLimit
		return currentBinIndex, big.NewInt(0), nil
	}

	// reservation not available, attempt on-demand
	//todo: rollback if disperser respond with some type of rejection?
	relativeBinRecord.Usage -= dataLength
	incrementRequired := big.NewInt(int64(a.PaymentCharged(uint(dataLength))))
	a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
	if a.cumulativePayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		return 0, a.cumulativePayment, nil
	}
	return 0, big.NewInt(0), fmt.Errorf("Accountant cannot approve payment for this blob")
}

// accountant provides and records payment information
func (a *Accountant) AccountBlob(ctx context.Context, dataLength uint64, quorums []uint8) (*commonpb.PaymentHeader, []byte, error) {
	binIndex, cumulativePayment, err := a.BlobPaymentInfo(ctx, dataLength)
	if err != nil {
		return nil, nil, err
	}

	accountID := a.paymentSigner.GetAccountID()
	pm := &core.PaymentMetadata{
		AccountID:         accountID,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
	}
	protoPaymentHeader := pm.ConvertToProtoPaymentHeader()

	signature, err := a.paymentSigner.SignBlobPayment(protoPaymentHeader)
	if err != nil {
		return nil, nil, err
	}

	return protoPaymentHeader, signature, nil
}

// TODO: PaymentCharged and SymbolsCharged copied from meterer, should be refactored
// PaymentCharged returns the chargeable price for a given data length
func (a *Accountant) PaymentCharged(dataLength uint) uint64 {
	return uint64(a.SymbolsCharged(dataLength)) * uint64(a.pricePerSymbol)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *Accountant) SymbolsCharged(dataLength uint) uint32 {
	if dataLength <= uint(a.minNumSymbols) {
		return a.minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return uint32(core.RoundUpDivide(uint(dataLength), uint(a.minNumSymbols))) * a.minNumSymbols
}

func (a *Accountant) GetRelativeBinRecord(index uint32) BinRecord {
	relativeIndex := index % 3

	if a.binRecords[relativeIndex].Index != uint32(index) {
		a.binRecords[relativeIndex] = BinRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}

	return a.binRecords[relativeIndex]
}

func (a *Accountant) GetOverflowBinRecord(index uint32) BinRecord {
	relativeIndex := (index + 2) % 3

	if a.binRecords[relativeIndex].Index != uint32(index+2) {
		a.binRecords[relativeIndex] = BinRecord{
			Index: uint32(index + 2),
			Usage: 0,
		}
	}

	return a.binRecords[relativeIndex]
}
