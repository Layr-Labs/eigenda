package payments

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
)

// TODO: we need to keep track of how many in flight dispersals there are, and not let that number exceed a certain
// value. The account ledger will need to check whether the on demand ledger is available before trying to debit,
// and do a wait if it isn't. We also need to consider how to "time out" an old request that was made to the disperser
// which was never responded to. We can't wait forever, eventually we need to declare a dispersal "failed", and move on

// EXTRACTED FROM accountant.go - ON-DEMAND PAYMENT RELATED CODE

// Variables related to on-demand payments
// OnDemandQuorums contains the required quorum numbers for on-demand payments as a set
var OnDemandQuorums = map[core.QuorumID]bool{
	0: true,
	1: true,
}

type OnDemandLedger struct {
}

func NewOnDemandLedger() (*OnDemandLedger, error) {
	return &OnDemandLedger{}, nil
}

// TODO: reconsider int64
func (odl *OnDemandLedger) Debit(symbolCount int64, quorums []core.QuorumID) error {
	if symbolCount <= 0 {
		return fmt.Errorf("symbolCount must be > 0, got %d", symbolCount)
	}

	for _, quorum := range quorums {
		if !OnDemandQuorums[quorum] {
			return fmt.Errorf("quorum %d cannot be dispersed to with on-demand payments", quorum)
		}
	}

	// TODO continue work here
	return nil
}

// Accountant struct fields related to on-demand payments:
// - onDemand          *core.OnDemandPayment
// - pricePerSymbol    uint64
// - minNumSymbols     uint64
// - cumulativePayment *big.Int

// From NewAccountant constructor - on-demand related initialization:
// onDemand:          onDemand,
// pricePerSymbol:    pricePerSymbol,
// minNumSymbols:     minNumSymbols,
// cumulativePayment: big.NewInt(0),

// From blobPaymentInfo - on-demand payment logic (lines 106-124):
// This section handles when reservation is not available and falls back to on-demand payment
/*
	// reservation not available, rollback reservation records, attempt on-demand
	//todo: rollback on-demand if disperser respond with some type of rejection?
	relativePeriodRecord.Usage -= symbolUsage
	incrementRequired := big.NewInt(int64(a.paymentCharged(numSymbols)))

	resultingPayment := big.NewInt(0)
	resultingPayment.Add(a.cumulativePayment, incrementRequired)
	if resultingPayment.Cmp(a.onDemand.CumulativePayment) <= 0 {
		if err := QuorumCheck(quorumNumbers, requiredQuorums); err != nil {
			return big.NewInt(0), err
		}
		a.cumulativePayment.Add(a.cumulativePayment, incrementRequired)
		return a.cumulativePayment, nil
	}
	return big.NewInt(0), fmt.Errorf(
		"invalid payments: no available bandwidth reservation found for account %s, and current cumulativePayment balance insufficient "+
			"to make an on-demand dispersal. Consider increasing reservation or cumulative payment on-chain. "+
			"For more details, see https://docs.eigenda.xyz/core-concepts/payments#disperser-client-requirements", a.accountID.Hex())
*/

// // paymentCharged returns the chargeable price for a given data length (lines 154-156)
// func (a *Accountant) paymentCharged(numSymbols uint64) uint64 {
// 	return a.symbolsCharged(numSymbols) * a.pricePerSymbol
// }

// // symbolsCharged returns the number of symbols charged for a given data length (lines 159-166)
// // being at least MinNumSymbols or the nearest rounded-up multiple of MinNumSymbols.
// func (a *Accountant) symbolsCharged(numSymbols uint64) uint64 {
// 	if numSymbols <= a.minNumSymbols {
// 		return a.minNumSymbols
// 	}
// 	// Round up to the nearest multiple of MinNumSymbols
// 	return core.RoundUpDivide(numSymbols, a.minNumSymbols) * a.minNumSymbols
// }

// From SetPaymentState - on-demand payment state setting (lines 198-217):
/*
	a.minNumSymbols = paymentState.GetPaymentGlobalParams().GetMinNumSymbols()
	a.pricePerSymbol = paymentState.GetPaymentGlobalParams().GetPricePerSymbol()

	if paymentState.GetOnchainCumulativePayment() == nil {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: big.NewInt(0),
		}
	} else {
		a.onDemand = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).SetBytes(paymentState.GetOnchainCumulativePayment()),
		}
	}

	if paymentState.GetCumulativePayment() == nil {
		a.cumulativePayment = big.NewInt(0)
	} else {
		a.cumulativePayment = new(big.Int).SetBytes(paymentState.GetCumulativePayment())
	}
*/

// QuorumCheck function used for on-demand payments (lines 260-270)
// This is used to check quorums for on-demand payments against requiredQuorums
// func QuorumCheck(quorumNumbers []uint8, allowedNumbers []uint8) error {
// 	if len(quorumNumbers) == 0 {
// 		return fmt.Errorf("no quorum numbers provided")
// 	}
// 	for _, quorum := range quorumNumbers {
// 		if !slices.Contains(allowedNumbers, quorum) {
// 			return fmt.Errorf("provided quorum number %v not allowed", quorum)
// 		}
// 	}
// 	return nil
// }

// NOTES ON ON-DEMAND PAYMENT FLOW:
// 1. When blobPaymentInfo is called, it first tries to use reservation
// 2. If reservation is not available (usage exceeds limit), it falls back to on-demand
// 3. For on-demand:
//    - It calculates incrementRequired based on paymentCharged(numSymbols)
//    - Checks if cumulativePayment + incrementRequired <= onDemand.CumulativePayment
//    - If yes, updates cumulativePayment and returns it (non-zero indicates on-demand)
//    - If no, returns error about insufficient balance
// 4. On-demand payments must use requiredQuorums (0, 1) instead of reservation quorums
// 5. Payment calculation uses pricePerSymbol * symbolsCharged
// 6. symbolsCharged rounds up to nearest minNumSymbols multiple
