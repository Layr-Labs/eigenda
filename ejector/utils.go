package ejector

import (
	"sort"
	"strings"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Combines two ValidatorSigningRate reports. Signed/unsigned batches and bytes are summed. Latency is taken
// as a weighed average (by batch count). If one of the rates is nil, the other is returned directly.
func combineSigningRates(
	rateA *validator.ValidatorSigningRate,
	rateB *validator.ValidatorSigningRate,
) *validator.ValidatorSigningRate {

	if rateA == nil {
		return rateB
	}
	if rateB == nil {
		return rateA
	}

	totalSignedBatches := rateA.GetSignedBatches() + rateB.GetSignedBatches()
	var latency uint64
	if totalSignedBatches > 0 {
		latency = (rateA.GetSigningLatency()*rateA.GetSignedBatches() +
			rateB.GetSigningLatency()*rateB.GetSignedBatches()) / totalSignedBatches
	}

	return &validator.ValidatorSigningRate{
		ValidatorId:     rateA.GetValidatorId(),
		SignedBatches:   rateA.GetSignedBatches() + rateB.GetSignedBatches(),
		UnsignedBatches: rateA.GetUnsignedBatches() + rateB.GetUnsignedBatches(),
		SignedBytes:     rateA.GetSignedBytes() + rateB.GetSignedBytes(),
		UnsignedBytes:   rateA.GetUnsignedBytes() + rateB.GetUnsignedBytes(),
		SigningLatency:  latency,
	}
}

// Sorts the given signing rates in place by unsigned bytes in descending order. The first entry will
// have the highest number of unsigned bytes, the last entry the lowest. Breaks ties by ordering by
// number of unsigned batches, also in descending order. Breaks further ties by ordering by validator ID
// in lexicographical order.
func sortByUnsignedBytesDescending(rates []*validator.ValidatorSigningRate) {
	sort.Slice(rates, func(i, j int) bool {
		// Primary sort: unsigned bytes (descending)
		if rates[i].GetUnsignedBytes() != rates[j].GetUnsignedBytes() {
			return rates[i].GetUnsignedBytes() > rates[j].GetUnsignedBytes()
		}

		// Tie breaker 1: unsigned batches (descending)
		if rates[i].GetUnsignedBatches() != rates[j].GetUnsignedBatches() {
			return rates[i].GetUnsignedBatches() > rates[j].GetUnsignedBatches()
		}

		// Tie breaker 2: validator ID (lexicographical ascending)
		return strings.Compare(string(rates[i].GetValidatorId()), string(rates[j].GetValidatorId())) < 0
	})
}
