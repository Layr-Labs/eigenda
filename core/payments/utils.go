package payments

// Computes the number of symbols to bill for a blob dispersal.
//
// If the actual symbol count is less than the minimum billable threshold, returns the minimum. Otherwise, returns the
// input symbol count.
//
// minNumSymbols is a parameter defined in the PaymentVault contract
func CalculateBillableSymbols(symbolCount uint32, minNumSymbols uint32) uint32 {
	if symbolCount < minNumSymbols {
		return minNumSymbols
	}
	return symbolCount
}
