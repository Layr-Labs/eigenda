package ondemand

import (
	"math/big"
)

type OnDemandLedgerConfig struct {
	totalDeposits  *big.Int
	pricePerSymbol *big.Int
	minNumSymbols  *big.Int
}
