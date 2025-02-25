package relay

import (
	"context"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// RelayUrlProvider provides relay URL strings, based on relay key
type RelayUrlProvider interface {
	// GetRelayUrl gets the URL string for a given relayKey
	GetRelayUrl(ctx context.Context, relayKey v2.RelayKey) (string, error)
	// GetRelayCount returns the number of relays in the registry
	GetRelayCount(ctx context.Context) (uint32, error)
}
