package clients

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
)

var minNumBins uint32 = 3
var requiredQuorums = []uint8{0, 1}

type Accountant interface {
	AccountBlob(ctx context.Context, numSymbols uint64, quorums []uint8) (*commonpb.PaymentHeader, []byte, error)
}

var _ Accountant = &accountant{}

type accountant struct {
	// on-chain states
	reservation       *core.ActiveReservation
	onDemand          *core.OnDemandPayment
	reservationWindow uint32
	pricePerSymbol    uint32
	minNumSymbols     uint32

	// local accounting
	// contains 3 bins; circular wrapping of indices
	binRecords        []BinRecord
	usageLock         sync.Mutex
	cumulativePayment *big.Int

	paymentSigner core.PaymentSigner
	numBins       uint32
}

type BinRecord struct {
	Index uint32
	Usage uint64
}

func NewAccountant(reservation *core.ActiveReservation, onDemand *core.OnDemandPayment, reservationWindow uint32, pricePerSymbol uint32, minNumSymbols uint32, paymentSigner core.PaymentSigner, numBins uint32) *accountant {
	//TODO: client storage; currently every instance starts fresh but on-chain or a small store makes more sense
	// Also client is currently responsible for supplying network params, we need to add RPC in order to be automatic
	// There's a subsequent PR that handles populating the accountant with on-chain state from the disperser
	binRecords := make([]BinRecord, numBins)
	for i := range binRecords {
		binRecords[i] = BinRecord{Index: uint32(i), Usage: 0}
	}
	a := accountant{
		reservation:       reservation,
		onDemand:          onDemand,
		reservationWindow: reservationWindow,
		pricePerSymbol:    pricePerSymbol,
		minNumSymbols:     minNumSymbols,
		binRecords:        binRecords,
		cumulativePayment: big.NewInt(0),
		paymentSigner:     paymentSigner,
		numBins:           max(numBins, minNumBins),
	}
	// TODO: add a routine to refresh the on-chain state occasionally?
	return &a
}

// BlobPaymentInfo calculates and records payment information. The accountant
// will attempt to use the active reservation first and check for quorum settings,
// then on-demand if the reservation is not available. The returned values are
// bin index for reservation payments and cumulative payment for on-demand payments,
// and both fields are used to create the payment header and signature
func (a *accountant) BlobPaymentInfo(ctx context.Context, numSymbols uint64, quorumNumbers []uint8) (uint32, *big.Int, error) {
	now := time.Now().Unix()
	currentBinIndex := meterer.GetBinIndex(uint64(now), a.reservationWindow)

	a.usageLock.Lock()
	defer a.usageLock.Unlock()
	relativeBinRecord := a.GetRelativeBinRecord(currentBinIndex)
	relativeBinRecord.Usage += numSymbols

	// first attempt to use the active reservation
	binLimit := a.reservation.SymbolsPerSec * uint64(a.reservationWindow)
	if relativeBinRecord.Usage <= binLimit {
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentBinIndex, big.NewInt(0), nil
	}

	overflowBinRecord := a.GetRelativeBinRecord(currentBinIndex + 2)
	// Allow one overflow when the overflow bin is empty, the current usage and new length are both less than the limit
	if overflowBinRecord.Usage == 0 && relativeBinRecord.Usage-numSymbols < binLimit && numSymbols <= binLimit {
		overflowBinRecord.Usage += relativeBinRecord.Usage - binLimit
		if err := QuorumCheck(quorumNumbers, a.reservation.QuorumNumbers); err != nil {
			return 0, big.NewInt(0), err
		}
		return currentBinIndex, big.NewInt(0), nil
	}

	// reservation not available, attempt on-demand
	//todo: rollback later if disperser respond with some type of rejection?
	relativeBinRecord.Usage -= numSymbols
	incrementRequired := big.NewInt(int64(a.PaymentCharged(uint(numSymbols))))
	a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
	if a.cumulativePayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
			return 0, big.NewInt(0), err
		}
		return 0, a.cumulativePayment, nil
	}
	return 0, big.NewInt(0), fmt.Errorf("neither reservation nor on-demand payment is available")
}

// AccountBlob accountant provides and records payment information
func (a *accountant) AccountBlob(ctx context.Context, numSymbols uint64, quorums []uint8) (*commonpb.PaymentHeader, []byte, error) {
	binIndex, cumulativePayment, err := a.BlobPaymentInfo(ctx, numSymbols, quorums)
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

	signature, err := a.paymentSigner.SignBlobPayment(pm)
	if err != nil {
		return nil, nil, err
	}

	return protoPaymentHeader, signature, nil
}

// TODO: PaymentCharged and SymbolsCharged copied from meterer, should be refactored
// PaymentCharged returns the chargeable price for a given data length
func (a *accountant) PaymentCharged(numSymbols uint) uint64 {
	return uint64(a.SymbolsCharged(numSymbols)) * uint64(a.pricePerSymbol)
}

// SymbolsCharged returns the number of symbols charged for a given data length
// being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
func (a *accountant) SymbolsCharged(numSymbols uint) uint32 {
	if numSymbols <= uint(a.minNumSymbols) {
		return a.minNumSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return uint32(core.RoundUpDivide(uint(numSymbols), uint(a.minNumSymbols))) * a.minNumSymbols
}

func (a *accountant) GetRelativeBinRecord(index uint32) *BinRecord {
	relativeIndex := index % a.numBins
	if a.binRecords[relativeIndex].Index != uint32(index) {
		a.binRecords[relativeIndex] = BinRecord{
			Index: uint32(index),
			Usage: 0,
		}
	}

	return &a.binRecords[relativeIndex]
}

// QuorumCheck eagerly returns error if the check finds a quorum number not an element of the allowed quorum numbers
func QuorumCheck(quorumNumbers []uint8, allowedNumbers []uint8) error {
	if len(quorumNumbers) == 0 {
		return fmt.Errorf("no quorum numbers provided")
	}
	for _, quorum := range quorumNumbers {
		if !slices.Contains(allowedNumbers, quorum) {
			return fmt.Errorf("provided quorum number %v not allowed", quorum)
		}
	}
	return nil
}
