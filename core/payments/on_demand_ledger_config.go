package payments

import (
	"math/big"
)

type OnDemandLedgerConfig struct {
	totalDeposits  *big.Int
	pricePerSymbol uint64
	minNumSymbols  uint64
}
