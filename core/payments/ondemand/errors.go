package ondemand

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
)

// Indicates that a requested quorum is not supported for on-demand payments.
type QuorumNotSupportedError struct {
	RequestedQuorum  core.QuorumID
	SupportedQuorums []core.QuorumID
}

func (e *QuorumNotSupportedError) Error() string {
	return fmt.Sprintf("quorum %v not supported for on-demand payments, supported quorums: %v",
		e.RequestedQuorum, e.SupportedQuorums)
}

// InsufficientFundsError indicates that the debit would exceed the total deposits available in the on-demand account.
type InsufficientFundsError struct {
	CurrentCumulativePayment *big.Int
	TotalDeposits            *big.Int
	BlobCost                 *big.Int
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf(
		"insufficient on-demand funds: current cumulative payment: %s wei, total deposits: %s wei, blob cost: %s wei",
		e.CurrentCumulativePayment.String(),
		e.TotalDeposits.String(),
		e.BlobCost.String())
}
