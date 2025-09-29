package ejector

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Signals whether we are using protocol version v1 or v2.
type ProtocolVersion int

const (
	ProtocolVersionV1 ProtocolVersion = 1
	ProtocolVersionV2 ProtocolVersion = 2
)

// A tool for looking up signing rates for validators.
type SigningRateLookup interface {
	// GetSigningRates returns signing rate information for all validators over the given time span.
	GetSigningRates(
		// The time span in the past over which to calculate signing rates.
		timeSpan time.Duration, 
		// Whether to collect signing rates for protocol version v1 or v2. Not all implementations may support both.
		version ProtocolVersion,
		// If true, omit validators with perfect signing rates (i.e. 100% signed). Some implementations
		// may ignore this flag (i.e. data API lookup).
		omitPerfectSigners bool,
	) ([]*validator.ValidatorSigningRate, error)
}
