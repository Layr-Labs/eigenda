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
	GetSigningRates(timeSpan time.Duration, version ProtocolVersion) ([]*validator.ValidatorSigningRate, error)
}
