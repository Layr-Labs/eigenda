package signingrate

import (
	"bytes"
	"sort"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Sort buckets by start time. Modifies the input slice.
func sortValidatorSigningRateBuckets(buckets []*validator.SigningRateBucket) {
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].GetStartTimestamp() < buckets[j].GetStartTimestamp()
	})
}

// Sort validator signing rates by ID. Modifies the input slice.
func sortValidatorSigningRate(rates []*validator.ValidatorSigningRate) {
	sort.Slice(rates, func(i, j int) bool {
		return bytes.Compare(rates[i].GetId(), rates[j].GetId()) < 0
	})
}
