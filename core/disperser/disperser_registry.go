package disperser

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
)

// DisperserRegistry provides access to disperser information from the DisperserRegistry contract.
type DisperserRegistry interface {
	// Returns the set of dispersers that network participants should interact with by default.
	GetDefaultDispersers(ctx context.Context) (map[uint32]struct{}, error)
	// Returns the set of dispersers that support on-demand payments
	GetOnDemandDispersers(ctx context.Context) (map[uint32]struct{}, error)
	// Returns the network address for a specific disperser
	GetDisperserNetworkAddress(ctx context.Context, disperserID uint32) (*common.NetworkAddress, error)
}
