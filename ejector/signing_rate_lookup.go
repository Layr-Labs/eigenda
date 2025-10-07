package ejector

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
)

// Signals whether we are using protocol version v1 or v2.
type ProtocolVersion int

const (
	ProtocolVersionV1 ProtocolVersion = 1
	ProtocolVersionV2 ProtocolVersion = 2
)

// A tool for looking up signing rates for validators.
type SigningRateLookup interface {
	// GetSigningRates returns signing rate information for all validators over the given time span. This method
	// is not required to return data in any particular order.
	GetSigningRates(
		// The time span in the past over which to calculate signing rates.
		timeSpan time.Duration,
		// A list of quorums to include. If empty, all quorums are included. If more than one quorum is given,
		// the results for each quorum are "summed" together. That is to say, each validator will only be returned in
		// a single result, and its signing rate will be equal to the sum of its signing rates across the all 
		// given quorums.
		quorums []core.QuorumID,
		// Whether to collect signing rates for protocol version v1 or v2. Not all implementations may support both.
		version ProtocolVersion,
		// If true, omit validators with perfect signing rates (i.e. 100% signed). Some implementations
		// may ignore this flag (i.e. data API lookup).
		omitPerfectSigners bool,
	) ([]*validator.ValidatorSigningRate, error)
}
