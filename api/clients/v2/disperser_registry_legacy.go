package clients

import (
	"context"
	"fmt"
)

var _ DisperserRegistry = (*LegacyDisperserRegistry)(nil)

// LegacyDisperserRegistry implements [DisperserRegistry] without interacting with the on-chain registry.
// This is a temporary implementation until the new DisperserRegistry contract is ready.
//
// TODO(litt3): We are currently working on a new DisperserRegistry contract which will support multiplexed dispersal,
// but it's not ready yet. For now, we have a legacy implementation that uses hardcoded values that match the current
// state of the network, before having deployed any additional dispersers.
type LegacyDisperserRegistry struct {
	connectionInfo *DisperserConnectionInfo
}

// NewLegacyDisperserRegistry creates a new legacy disperser registry.
// The connectionInfo parameter specifies how to connect to disperser ID 0.
func NewLegacyDisperserRegistry(connectionInfo *DisperserConnectionInfo) *LegacyDisperserRegistry {
	return &LegacyDisperserRegistry{
		connectionInfo: connectionInfo,
	}
}

// GetDefaultDispersers implements [DisperserRegistry].
//
// Return a single default disperser with ID 0, which is the only disperser currently deployed on the network.
func (r *LegacyDisperserRegistry) GetDefaultDispersers(ctx context.Context) ([]uint32, error) {
	return []uint32{0}, nil
}

// GetOnDemandDispersers implements [DisperserRegistry].
//
// Return a single on-demand disperser with ID 0, which is the only disperser currently deployed on the network.
func (r *LegacyDisperserRegistry) GetOnDemandDispersers(ctx context.Context) ([]uint32, error) {
	return []uint32{0}, nil
}

// GetDisperserConnectionInfo implements [DisperserRegistry].
//
// Returns the connection info for disperser ID 0. All other IDs return an error.
func (r *LegacyDisperserRegistry) GetDisperserConnectionInfo(
	ctx context.Context,
	disperserID uint32,
) (*DisperserConnectionInfo, error) {
	if disperserID != 0 {
		return nil, fmt.Errorf("legacy registry only supports disperser ID 0, got %d", disperserID)
	}
	return r.connectionInfo, nil
}
