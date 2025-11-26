package disperser

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
)

var _ DisperserRegistry = (*LegacyDisperserRegistry)(nil)

// LegacyDisperserRegistry implements [DisperserRegistry] without actually interacting with the on-chain registry.
//
// TODO(litt3): We are currently working on a new DisperserRegistry contract which will support multiplexed dispersal,
// but it's not ready yet. For now, we have a legacy implementation that uses hardcoded values that match the current
// state of the network, before having deployed any additional dispersers.
type LegacyDisperserRegistry struct {
	networkAddress *common.NetworkAddress
}

// Creates a new legacy disperser registry.
// The networkAddress parameter specifies how to connect to disperser ID 0.
func NewLegacyDisperserRegistry(networkAddress *common.NetworkAddress) *LegacyDisperserRegistry {
	return &LegacyDisperserRegistry{
		networkAddress: networkAddress,
	}
}

// GetDefaultDispersers implements [DisperserRegistry].
//
// Return a single default disperser with ID 0, which is the only disperser currently deployed on the network.
func (r *LegacyDisperserRegistry) GetDefaultDispersers(ctx context.Context) (map[uint32]struct{}, error) {
	return map[uint32]struct{}{0: {}}, nil
}

// GetOnDemandDispersers implements [DisperserRegistry].
//
// Return a single on-demand disperser with ID 0, which is the only disperser currently deployed on the network.
func (r *LegacyDisperserRegistry) GetOnDemandDispersers(ctx context.Context) (map[uint32]struct{}, error) {
	return map[uint32]struct{}{0: {}}, nil
}

// Implements [DisperserRegistry].
//
// Returns the network address for disperser ID 0. All other IDs return an error.
func (r *LegacyDisperserRegistry) GetDisperserNetworkAddress(
	ctx context.Context,
	disperserID uint32,
) (*common.NetworkAddress, error) {
	if disperserID != 0 {
		return nil, fmt.Errorf("legacy registry only supports disperser ID 0, got %d", disperserID)
	}
	return r.networkAddress, nil
}
