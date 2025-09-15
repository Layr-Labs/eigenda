package signingrate

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Sort buckets by start time. Modifies the input slice.
func sortValidatorSigningRateBuckets(buckets []*validator.SigningRateBucket) {
	sort.Slice(buckets, func(i int, j int) bool {
		return buckets[i].GetStartTimestamp() < buckets[j].GetStartTimestamp()
	})
}

// Sort validator signing rates by validator ID. Modifies the input slice.
func sortValidatorSigningRates(rates []*validator.ValidatorSigningRate) {
	sort.Slice(rates, func(i int, j int) bool {
		return bytes.Compare(rates[i].GetId(), rates[j].GetId()) < 0
	})
}

// Sort quorum signing rates by quorum ID. Modifies the input slice.
func sortQuorumSigningRates(quorums []*validator.QuorumSigningRate) {
	sort.Slice(quorums, func(i int, j int) bool {
		return quorums[i].GetQuorumId() < quorums[j].GetQuorumId()
	})
}

// Performs a deep copy of a ValidatorSigningRate.
func cloneValidatorSigningRate(info *validator.ValidatorSigningRate) *validator.ValidatorSigningRate {
	return &validator.ValidatorSigningRate{
		Id:              info.GetId(),
		SignedBatches:   info.GetSignedBatches(),
		SignedBytes:     info.GetSignedBytes(),
		UnsignedBatches: info.GetUnsignedBatches(),
		UnsignedBytes:   info.GetUnsignedBytes(),
		SigningLatency:  info.GetSigningLatency(),
	}
}

// Given a timestamp, finds the start timestamp of the bucket that contains that timestamp (inclusive).
// The "primary key" of a bucket is the start timestamp, so this function effectively maps an arbitrary timestamp
// to the key of the bucket that contains data for this timestamp.
//
// Bucket timestamps are aligned with to clean multiples of the bucket span. If the bucket span is 10 minutes, then
// the first bucket will start at the epoch, the second bucket will start exactly 10 minutes after the epoch, and so on.
//
// Bucket timestamps are always reported at second granularity (i.e. no fractional seconds).
//
// Although humans consuming information in a UI may find it aesthetically pleasing to have buckets aligned to
// wall-clock boundaries (e.g. a 1-hour bucket starting at the top of the hour in UTC), such alignment is
// disproportionately complex for the benefit it provides.
func bucketStartTimestamp(bucketSpan time.Duration, targetTime time.Time) (time.Time, error) {
	spanSeconds := uint64(bucketSpan.Seconds())
	if spanSeconds == 0 {
		return time.Time{}, fmt.Errorf("bucket span must be at least one second, got %s", bucketSpan)
	}

	targetSeconds := uint64(targetTime.Unix())

	startTimestampSeconds := (targetSeconds / spanSeconds) * spanSeconds
	return time.Unix(int64(startTimestampSeconds), 0), nil
}

// Given a timestamp, finds the end timestamp of the bucket that contains that timestamp (exclusive).
func bucketEndTimestamp(bucketSpan time.Duration, targetTime time.Time) (time.Time, error) {
	startTimestamp, err := bucketStartTimestamp(bucketSpan, targetTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("bucket start timestamp: %w", err)
	}
	return time.Unix(startTimestamp.Unix()+int64(bucketSpan.Seconds()), 0), nil
}
