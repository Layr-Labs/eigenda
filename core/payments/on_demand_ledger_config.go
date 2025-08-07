package payments

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type OnDemandLedgerConfig struct {
	accountID      gethcommon.Address
	totalDeposits  *big.Int
	pricePerSymbol uint64
	minNumSymbols  uint64
}
